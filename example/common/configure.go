package common

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fkgi/diameter/connection"
	"github.com/fkgi/diameter/msg"
)

// Log is logger
var Log *log.Logger

// GeneratePath generate file path
func GeneratePath(s string) (i, o, c string) {
	// get option flag
	isock := flag.String("i", s+".in", "input UNIX socket name")
	osock := flag.String("o", s+".out", "output UNIX socket name")
	conf := flag.String("c", s+".xml", "xml config file name")
	flag.Parse()

	// create path
	if wdir, e := os.Getwd(); e != nil {
		Log.Fatalln(e)
	} else {
		i = wdir + string(os.PathSeparator) + *isock
		o = wdir + string(os.PathSeparator) + *osock
		c = wdir + string(os.PathSeparator) + *conf
	}
	return
}

type xmlConfig struct {
	XMLName  xml.Name `xml:"config"`
	Local    xmlNode  `xml:"local"`
	Peer     xmlNode  `xml:"peer"`
	Watchdog xmlTimer `xml:"watchdog"`
	Message  xmlTimer `xml:"message"`
}

type xmlNode struct {
	NetType string `xml:"type"`
	Addr    string `xml:"addr"`
	FQDN    string `xml:"fqdn"`
}

type xmlTimer struct {
	Interval string `xml:"interval"`
	Retry    string `xml:"retry"`
}

// LoadConfig load xml config file
func LoadConfig(conf string, ln *connection.LocalNode, pn *connection.PeerNode) (
	la net.Addr, pa net.Addr) {

	Log.Println("config file is " + conf)
	var e error

	// load config file
	x := xmlConfig{}
	if data, e := ioutil.ReadFile(conf); e != nil {
		Log.Fatalln(e)
	} else if e = xml.Unmarshal([]byte(data), &x); e != nil {
		Log.Fatalln(e)
	}

	Log.Printf("local-host parameter:")
	la = loadNetAddr(x.Local.NetType, x.Local.Addr)
	if la != nil {
		Log.Printf("  address        =%s:%s", la.Network(), la)
	}
	ln.Host, ln.Realm = genHostRealm(x.Local.FQDN)
	Log.Printf("  diameter host  =%s", ln.Host)
	Log.Printf("  diameter realm =%s", ln.Realm)
	ln.InitIDs()

	if pn != nil {
		Log.Printf("peer-host parameter:")
		pa = loadNetAddr(x.Peer.NetType, x.Peer.Addr)
		Log.Printf("  address        =%s:%s", pa.Network(), pa)
		pn.Host, pn.Realm = genHostRealm(x.Peer.FQDN)
		Log.Printf("  diameter host  =%s", pn.Host)
		Log.Printf("  diameter realm =%s", pn.Realm)
	}

	Log.Printf("diameter parameter:")
	ln.Properties = connection.Properties{}

	var i int
	x.Watchdog.Interval = strings.TrimSpace(x.Watchdog.Interval)
	i, e = strconv.Atoi(x.Watchdog.Interval)
	ln.Properties.Tw = time.Second * time.Duration(i)
	if e != nil {
		Log.Fatalln("invalid DWR timer:", e)
	}
	Log.Printf("  watchdog timer       =%d[sec]", i)

	x.Watchdog.Retry = strings.TrimSpace(x.Watchdog.Retry)
	ln.Properties.Ew, e = strconv.Atoi(x.Watchdog.Retry)
	if e != nil {
		Log.Fatalln("invalid DWR retry count:", e)
	}
	Log.Printf("  watchdog retry count =%d", ln.Properties.Ew)

	x.Message.Interval = strings.TrimSpace(x.Message.Interval)
	i, e = strconv.Atoi(x.Message.Interval)
	ln.Properties.Tp = time.Second * time.Duration(i)
	if e != nil {
		Log.Fatalln("invalid Message retry timer:", e)
	}
	Log.Printf("  msg send timer       =%d[sec]", i)

	x.Message.Retry = strings.TrimSpace(x.Message.Retry)
	ln.Properties.Cp, e = strconv.Atoi(x.Message.Retry)
	if e != nil {
		Log.Fatalln("invalid Message retry count:", e)
	}
	Log.Printf("  msg send retry count =%d", ln.Properties.Cp)

	ln.Properties.Ts = time.Millisecond * time.Duration(100)
	Log.Printf("  msg send transport timeout =%d[msec]",
		ln.Properties.Ts/time.Millisecond)

	ln.Properties.Apps = make([]connection.AuthApplication, 0)
	ln.Properties.Apps = append(ln.Properties.Apps,
		connection.AuthApplication{VendorID: 0, AppID: 0})
	ln.Properties.Apps = append(ln.Properties.Apps,
		connection.AuthApplication{VendorID: 0, AppID: 0xffffffff})

	if pn != nil {
		pn.Properties = ln.Properties
	}

	return
}

func loadNetAddr(nettype, netaddr string) (addr net.Addr) {
	netaddr = strings.TrimSpace(netaddr)
	nettype = strings.TrimSpace(nettype)

	switch nettype {
	/*
		case "sctp", "sctp4", "sctp6":
			a, e := extnet.ResolveSCTPAddr(nettype, netaddr)
			if e != nil {
				Log.Fatalln("invalid sctp address", e)
			}
			addr = a
	*/
	case "tcp", "tcp4", "tcp6":
		var e error
		addr, e = net.ResolveTCPAddr(nettype, netaddr)
		if e != nil {
			addr = nil
		}
	default:
		addr = nil
	}
	return
}

func genHostRealm(fqdn string) (host, realm msg.DiameterIdentity) {
	fqdn = strings.TrimSpace(fqdn)
	var e error
	host, e = msg.ParseDiameterIdentity(fqdn)
	if e != nil {
		log.Fatalln("invalid host name:", e)
	}

	realm, e = msg.ParseDiameterIdentity(fqdn[strings.Index(fqdn, ".")+1:])
	if e != nil {
		log.Fatalln("invalid host realm:", e)
	}
	return
}
