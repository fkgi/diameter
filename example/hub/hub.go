package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/fkgi/diameter/connection"
	"github.com/fkgi/diameter/msg"
)

var (
	logger  = log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds)
	conset  = make(map[msg.DestinationHost]*connection.Connection)
	gateway *string
	tw      = 20
	ew      = 3
	tp      = 30
	cp      = 3
	ts      = 100
)

func main() {
	logger.Printf("booting HUB with <%s REV.%d>...",
		connection.ProductName, connection.FirmwareRevision)

	h, e := os.Hostname()
	if e != nil {
		h = "localhost.localnetwork"
	}
	port := flag.String("p", ":3868", "port")
	fqdn := flag.String("h", h, "diameter host name")
	gateway = flag.String("g", "", "default gateway host name")
	flag.Parse()

	ln := connection.LocalNode{}
	ln.Host, e = msg.ParseDiameterIdentity(*fqdn)
	if e != nil {
		logger.Fatalln(" | invalid host name:", e)
	}
	logger.Printf(" | diameter host  =%s", ln.Host)

	ln.Realm, e = msg.ParseDiameterIdentity((*fqdn)[strings.Index(*fqdn, ".")+1:])
	if e != nil {
		logger.Fatalln(" | invalid host realm:", e)
	}
	logger.Printf(" | diameter realm =%s", ln.Realm)
	ln.InitIDs()

	ln.Properties.Tw = time.Second * time.Duration(tw)
	logger.Printf(" | watchdog timer       =%d[sec]", tw)

	ln.Properties.Ew = ew
	logger.Printf(" | watchdog retry count =%d", ew)

	ln.Properties.Tp = time.Second * time.Duration(tp)
	logger.Printf(" | msg send timer       =%d[sec]", tp)

	ln.Properties.Cp = cp
	logger.Printf(" | msg send retry count =%d", cp)

	ln.Properties.Ts = time.Millisecond * time.Duration(ts)
	logger.Printf(" | msg send transport timeout =%d[msec]", ts)

	ln.Properties.Apps = make([]connection.AuthApplication, 0)
	ln.Properties.Apps = append(ln.Properties.Apps,
		connection.AuthApplication{VendorID: 0, AppID: 0})
	ln.Properties.Apps = append(ln.Properties.Apps,
		connection.AuthApplication{VendorID: 0, AppID: 0xffffffff})

	logger.Println("listening ...")
	l, e := net.Listen("tcp", *port)
	if e != nil {
		logger.Fatalln(" | invalid address:", e)
	}
	logger.Printf(" | address =%s:%s", l.Addr().Network(), l.Addr().String())

	for {
		c, e := l.Accept()
		if e != nil {
			logger.Fatalln(e)
		}
		go handleConnection(ln.Accept(c))
	}
}

func handleConnection(c *connection.Connection) {
	logger.Println("accept new connection")
	if !c.WaitOpen() {
		return
	}
	logger.Printf("connection (host=%s, realm=%s) is open", c.PeerHost(), c.PeerRealm())
	if _, ok := conset[msg.DestinationHost(c.PeerHost())]; ok {
		c.Close(msg.Rebooting)
		return
	}
	conset[msg.DestinationHost(c.PeerHost())] = c

	for {
		m, ch, e := c.Recieve()
		if e != nil {
			logger.Println(e)
			break
		}
		if avp, e := m.Decode(); e == nil {
			t, ok := msg.GetDestinationHost(avp)
			if !ok {
				t = msg.DestinationHost(*gateway)
			}
			if dst, ok := conset[t]; ok {
				m = dst.Transmit(m)
				ch <- &m
				continue
			}
		}

		m = msg.Message{
			Ver:   msg.DiaVer,
			FlgR:  false,
			FlgP:  m.FlgP,
			FlgE:  true,
			FlgT:  false,
			HbHID: m.HbHID,
			EtEID: m.EtEID,
			Code:  m.Code,
			AppID: m.AppID}
		var avps []msg.Avp
		avps = append(avps, msg.DiameterUnableToDeliver.Encode())
		avps = append(avps, msg.OriginHost(c.LocalHost()).Encode())
		avps = append(avps, msg.OriginRealm(c.LocalRealm()).Encode())
		m.Encode(avps)

		ch <- &m
	}
	c.Close(msg.Rebooting)
}
