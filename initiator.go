package main

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/provider"
)

var (
	isock *string
	osock *string
	prov  *provider.Provider
)

const (
	v = "01"
)

type xmlConfig struct {
	XMLName  xml.Name `xml:"config"`
	Local    xmlNode  `xml:"local"`
	Peer     xmlNode  `xml:"peer"`
	Watchdog xmlTimer `xml:"watchdog"`
	Message  xmlTimer `xml:"message"`
	//	Ttl     string   `xml:"ttl"`
}

type xmlNode struct {
	Addr string `xml:"addr"`
	FQDN string `xml:"fqdn"`
}

type xmlTimer struct {
	Interval string `xml:"interval"`
	Retry    string `xml:"retry"`
}

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	provider.Notificator = func(e error) { log.Println(e) }

	// boot log
	log.Println("Mave ver." + v + " booting ...")

	// get option flag
	isock = flag.String("i", "mave.in", "input UNIX socket name")
	osock = flag.String("o", "mave.out", "output UNIX socket name")
	conf := flag.String("c", "mave.xml", "xml config file name")
	flag.Parse()

	// create path
	if wdir, e := os.Getwd(); e != nil {
		log.Fatalln(e)
	} else {
		*isock = wdir + string(os.PathSeparator) + *isock
		*osock = wdir + string(os.PathSeparator) + *osock
		*conf = wdir + string(os.PathSeparator) + *conf
	}

	// load config file
	log.Println("config file is " + *conf)
	x := xmlConfig{}
	if data, e := ioutil.ReadFile(*conf); e != nil {
		log.Fatalln(e)
	} else if e = xml.Unmarshal([]byte(data), &x); e != nil {
		log.Fatalln(e)
	}

	x.Local.FQDN = strings.TrimSpace(x.Local.FQDN)
	lhost, e := msg.ParseDiameterIdentity(x.Local.FQDN)
	if e != nil {
		log.Fatalln("invalid local-host name config:", e)
	}
	lrealm, e := msg.ParseDiameterIdentity(
		x.Local.FQDN[strings.Index(x.Local.FQDN, ".")+1:])
	if e != nil {
		log.Fatalln("invalid local-host realm config:", e)
	}
	x.Local.Addr = strings.TrimSpace(x.Local.Addr)
	laddr, e := net.ResolveTCPAddr("tcp4", x.Local.Addr)
	if e != nil {
		log.Fatalln("invalid local-host address config:", e)
	}

	x.Peer.FQDN = strings.TrimSpace(x.Peer.FQDN)
	phost, e := msg.ParseDiameterIdentity(x.Peer.FQDN)
	if e != nil {
		log.Fatalln("invalid local-host name config:", e)
	}
	prealm, e := msg.ParseDiameterIdentity(
		x.Peer.FQDN[strings.Index(x.Peer.FQDN, ".")+1:])
	if e != nil {
		log.Fatalln("invalid local-host realm config:", e)
	}
	x.Peer.Addr = strings.TrimSpace(x.Peer.Addr)
	paddr, e := net.ResolveTCPAddr("tcp4", x.Peer.Addr)
	if e != nil {
		log.Fatalln("invalid local-host address config:", e)
	}

	x.Watchdog.Interval = strings.TrimSpace(x.Watchdog.Interval)
	wtime, e := strconv.Atoi(x.Watchdog.Interval)
	if e != nil {
		log.Fatalln("invalid DWR timer:", e)
	}
	x.Watchdog.Retry = strings.TrimSpace(x.Watchdog.Retry)
	wcount, e := strconv.Atoi(x.Watchdog.Retry)
	if e != nil {
		log.Fatalln("invalid DWR retry count:", e)
	}
	x.Message.Interval = strings.TrimSpace(x.Message.Interval)
	mtime, e := strconv.Atoi(x.Message.Interval)
	if e != nil {
		log.Fatalln("invalid Message retry timer:", e)
	}
	x.Message.Retry = strings.TrimSpace(x.Message.Retry)
	mcount, e := strconv.Atoi(x.Message.Retry)
	if e != nil {
		log.Fatalln("invalid Message retry count:", e)
	}

	log.Println("local-host parameter:")
	log.Println("  address        =", laddr)
	log.Println("  diameter host  =", lhost)
	log.Println("  diameter realm =", lrealm)
	log.Println("peer-host parameter:")
	log.Println("  address        =", paddr)
	log.Println("  diameter host  =", phost)
	log.Println("  diameter realm =", prealm)

	// open Diameter socket
	log.Println("start connecting Diameter connection")
	ln := &provider.LocalNode{}
	ln.Host = lhost
	ln.Realm = lrealm
	ln.Addr = make([]net.IP, 1)
	ln.Addr[0] = laddr.IP
	ln.InitIDs()
	pl := provider.Listen(ln)

	pn := &provider.PeerNode{}
	pn.Host = phost
	pn.Realm = prealm
	pn.Addr = make([]net.IP, 1)
	pn.Addr[0] = paddr.IP
	pn.Tw = time.Second * time.Duration(wtime)
	pn.Ew = wcount
	pn.Tp = time.Second * time.Duration(mtime)
	pn.Cp = mcount
	pn.Ts = time.Millisecond * time.Duration(100)
	pn.SupportedApps = make([][2]uint32, 0)
	pn.SupportedApps = append(pn.SupportedApps, [2]uint32{0, 0})
	pn.SupportedApps = append(pn.SupportedApps, [2]uint32{0, 0xffffffff})
	prov = pl.AddPeer(pn)

	pl.Dial(pn, laddr, paddr)
	time.Sleep(time.Second)

	// open UNIX socket
	log.Println("start listening on UNIX socket", *isock, "and", *osock)
	il, e := net.Listen("unix", *isock)
	if e != nil {
		log.Fatalln(e)
	}
	ol, e := net.Listen("unix", *osock)
	if e != nil {
		log.Fatalln(e)
	}

	// set kill-signal trap
	setKillHandler(stop)

	// read UNIX socket
	go func() {
		for {
			c, e := il.Accept()
			if e != nil {
				log.Println(e)
				stop()
			}
			defer c.Close()

			m, ch, e := prov.Recieve()
			if e != nil {
				log.Println(e)
				continue
			}
			if _, e = m.WriteTo(c); e != nil {
				log.Println(e)
				continue
			}
			if _, e = m.ReadFrom(c); e != nil {
				log.Println(e)
				continue
			}
			ch <- &m
		}
	}()

	for {
		c, e := ol.Accept()
		if e != nil {
			log.Println(e)
			stop()
		}
		defer c.Close()

		m := msg.Message{}
		if _, e = m.ReadFrom(c); e != nil {
			log.Println(e)
			continue
		}
		if avp, e := m.Decode(); e != nil {
			log.Println(e)
			continue
		} else {
			var src msg.DiameterIdentity
			for _, a := range avp {
				if a.Code == 264 {
					e = a.Decode(&src)
					if e != nil {
						log.Println(e)
					}
					break
				}
			}
			avp = append(avp, msg.RouteRecord(src))
			if e = m.Encode(avp); e != nil {
				log.Println(e)
				continue
			}
		}
		if m, e = prov.Send(m); e != nil {
			log.Println(e)
			continue
		}
		if _, e = m.WriteTo(c); e != nil {
			log.Println(e)
			continue
		}
	}
}

func stop() {
	log.Println("shutdown ...")

	prov.Close(msg.Enumerated(0))
	os.Remove(*isock)
	os.Remove(*osock)
	time.Sleep(time.Second)

	os.Exit(0)
}

func setKillHandler(f func()) {
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-sigc
		f()
	}()
}
