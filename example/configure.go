package example

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/provider"
	"github.com/fkgi/extnet"
)

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
func LoadConfig(conf string) (
	la net.Addr, ln *provider.LocalNode, pa net.Addr, pn *provider.PeerNode) {

	log.Println("config file is " + conf)
	var e error

	// load config file
	x := xmlConfig{}
	if data, e := ioutil.ReadFile(conf); e != nil {
		log.Fatalln(e)
	} else if e = xml.Unmarshal([]byte(data), &x); e != nil {
		log.Fatalln(e)
	}

	log.Printf("local-host parameter:")
	ln = &provider.LocalNode{}
	la, ln.Addr = loadNetAddr(x.Local.NetType, x.Local.Addr)
	log.Printf("  address        =%s:%s", la.Network(), la)
	ln.Host, ln.Realm = loadHostRealm(x.Local.FQDN)
	log.Printf("  diameter host  =%s", ln.Host)
	log.Printf("  diameter realm =%s", ln.Realm)

	ln.InitIDs()

	log.Printf("peer-host parameter:")
	pn = &provider.PeerNode{}
	pa, pn.Addr = loadNetAddr(x.Peer.NetType, x.Peer.Addr)
	log.Printf("  address        =%s:%s", pa.Network(), pa)
	pn.Host, pn.Realm = loadHostRealm(x.Peer.FQDN)
	log.Printf("  diameter host  =%s", pn.Host)
	log.Printf("  diameter realm =%s", pn.Realm)

	log.Printf("diameter parameter:")
	var i int
	x.Watchdog.Interval = strings.TrimSpace(x.Watchdog.Interval)
	i, e = strconv.Atoi(x.Watchdog.Interval)
	pn.Tw = time.Second * time.Duration(i)
	if e != nil {
		log.Fatalln("invalid DWR timer:", e)
	}
	log.Printf("  watchdog timer       =%d[sec]", i)

	x.Watchdog.Retry = strings.TrimSpace(x.Watchdog.Retry)
	pn.Ew, e = strconv.Atoi(x.Watchdog.Retry)
	if e != nil {
		log.Fatalln("invalid DWR retry count:", e)
	}
	log.Printf("  watchdog retry count =%d", pn.Ew)

	x.Message.Interval = strings.TrimSpace(x.Message.Interval)
	i, e = strconv.Atoi(x.Message.Interval)
	pn.Tp = time.Second * time.Duration(i)
	if e != nil {
		log.Fatalln("invalid Message retry timer:", e)
	}
	log.Printf("  msg send timer       =%d[sec]", i)

	x.Message.Retry = strings.TrimSpace(x.Message.Retry)
	pn.Cp, e = strconv.Atoi(x.Message.Retry)
	if e != nil {
		log.Fatalln("invalid Message retry count:", e)
	}
	log.Printf("  msg send retry count =%d", pn.Cp)

	pn.Ts = time.Millisecond * time.Duration(100)
	pn.SupportedApps = make([][2]uint32, 0)
	pn.SupportedApps = append(pn.SupportedApps, [2]uint32{0, 0})
	pn.SupportedApps = append(pn.SupportedApps, [2]uint32{0, 0xffffffff})

	return
}

func loadNetAddr(nettype, netaddr string) (addr net.Addr, ip []net.IP) {
	netaddr = strings.TrimSpace(netaddr)
	nettype = strings.TrimSpace(nettype)

	switch nettype {
	case "sctp", "sctp4", "sctp6":
		a, e := extnet.ResolveSCTPAddr(nettype, netaddr)
		if e != nil {
			log.Fatalln("invalid sctp address", e)
		}
		addr = a
		ip = a.IP
	case "tcp", "tcp4", "tcp6":
		a, e := net.ResolveTCPAddr(nettype, netaddr)
		if e != nil {
			log.Fatalln("invalid tcp address", e)
		}
		addr = a
		ip = []net.IP{a.IP}
	default:
		log.Fatalln("invalid network type", nettype)
	}
	return
}

func loadHostRealm(fqdn string) (host, realm msg.DiameterIdentity) {
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
