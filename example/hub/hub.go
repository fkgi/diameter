package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/sock"
)

var (
	logger  = log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds)
	conset  = make(map[msg.DestinationHost]*sock.Conn)
	gateway *string
)

func main() {
	logger.Printf(
		"booting HUB with <%s REV.%d>...",
		sock.ProductName, sock.FirmwareRevision)

	h, e := os.Hostname()
	if e != nil {
		h = "localhost.localnetwork"
	}
	port := flag.String("p", ":3868", "local address")
	fqdn := flag.String("h", h, "diameter host name")
	gateway = flag.String("g", "", "default gateway host name")
	flag.Parse()

	ln := sock.Local{}
	ln.Host, e = msg.ParseDiameterIdentity(*fqdn)
	if e != nil {
		logger.Fatalln(" | invalid host name:", e)
	}
	logger.Printf(" | diameter host              =%s", ln.Host)

	ln.Realm, e = msg.ParseDiameterIdentity(
		(*fqdn)[strings.Index(*fqdn, ".")+1:])
	if e != nil {
		logger.Fatalln(" | invalid host realm:", e)
	}
	logger.Printf(" | diameter realm             =%s", ln.Realm)

	ln.AuthApps = map[msg.VendorID][]msg.ApplicationID{
		0: []msg.ApplicationID{
			msg.AuthApplicationID(0xffffffff)}}

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
		con, e := ln.Accept(nil, c)
		if e != nil {
			logger.Fatalln(e)
		}
		go handleConnection(con)
	}
}

func handleConnection(c *sock.Conn) {
	logger.Println("accept new sock")
	for c.State() != "open" {
		time.Sleep(time.Millisecond * 100)
	}

	logger.Printf("sock (host=%s, realm=%s) is open",
		c.PeerHost(), c.PeerRealm())
	if _, ok := conset[msg.DestinationHost(c.PeerHost())]; ok {
		logger.Printf("error duplicate host")
		c.Close(time.Millisecond * 100)
		return
	}
	conset[msg.DestinationHost(c.PeerHost())] = c

	for c.State() != "close" {
	}
	/*
			m, ch, e := c.Recieve()
			if e != nil {
				logger.Printf("recieve message failed: %s", e)
				break
			}
			logger.Printf("recieve message from host=%s, realm=%s", c.PeerHost(), c.PeerRealm())
			if avp, e := m.Decode(); e == nil {
				t, ok := msg.GetDestinationHost(avp)
				if !ok {
					t = msg.DestinationHost(*gateway)
				}
				if dst, ok := conset[t]; ok {
					logger.Printf("message send to host=%s, realm=%s", dst.PeerHost(), dst.PeerRealm())
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

			logger.Printf("message destination not found")
			ch <- &m
		}
		c.Close(msg.Rebooting)
	*/
	delete(conset, msg.DestinationHost(c.PeerHost()))
	logger.Printf("sock (host=%s, realm=%s) is closed",
		c.PeerHost(), c.PeerRealm())
}
