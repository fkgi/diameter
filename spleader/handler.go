package main

import (
	"bytes"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/connector"
	"github.com/fkgi/diameter/sctp"
)

var cons = make(map[net.Conn]*diameter.Connection)

func ListenAndServe(la string) (err error) {
	scheme, host, realm, ips, port, err := connector.ResolveIdentity(la)
	if err != nil {
		return
	}
	diameter.Host = host
	diameter.Realm = realm

	var l net.Listener
	switch scheme {
	case "sctp":
		src := &sctp.SCTPAddr{IP: ips, Port: port}
		if connector.TransportInfoNotify != nil {
			connector.TransportInfoNotify(src, nil)
		}
		l, err = sctp.ListenSCTP(src)
	default:
		src := &net.TCPAddr{IP: ips[0], Port: port}
		if connector.TransportInfoNotify != nil {
			connector.TransportInfoNotify(src, nil)
		}
		l, err = net.ListenTCP("tcp", src)
	}
	if err != nil {
		return
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-sigc
		l.Close()
	}()

	for {
		var c net.Conn
		if c, err = l.Accept(); err != nil {
			break
		}
		if connector.TransportUpNotify != nil {
			connector.TransportUpNotify(c.LocalAddr(), c.RemoteAddr())
		}

		con := diameter.Connection{}
		cons[c] = &con
		go func() {
			con.ListenAndServe(c)
			delete(cons, c)
		}()
	}
	l.Close()

	for _, con := range cons {
		con.Close(diameter.Rebooting)
	}
	for len(cons) != 0 {
		time.Sleep(time.Millisecond * 100)
	}

	return
}

func rxhandler(m diameter.Message) diameter.Message {
	var dHost diameter.Identity
	for rdr := bytes.NewReader(m.AVPs); rdr.Len() != 0; {
		a := diameter.AVP{}
		e := a.UnmarshalFrom(rdr)
		if e != nil {
			continue
		}
		if a.VendorID != 0 {
			continue
		}
		if a.Code == 293 {
			dHost, _ = diameter.GetDestinationHost(a)
		}
	}

	dcon := []*diameter.Connection{}
	if m.PeerName == upLink {
		if dHost == upLink {
			return m.GenerateAnswerBy(diameter.UnableToDeliver)
		}
		for _, con := range cons {
			if con.Host == dHost {
				dcon = append(dcon, con)
			}
		}
		if len(dcon) == 0 {
			for _, con := range cons {
				if con.Host != upLink {
					dcon = append(dcon, con)
				}
			}
		}
	} else {
		for _, con := range cons {
			if con.Host == upLink {
				dcon = append(dcon, con)
			}
		}
	}

	if len(dcon) == 0 {
		return m.GenerateAnswerBy(diameter.UnableToDeliver)
	}

	buf := bytes.NewBuffer(m.AVPs)
	diameter.SetRouteRecord(diameter.Host).MarshalTo(buf)
	m.AVPs = buf.Bytes()
	return dcon[rand.Intn(len(dcon))].DefaultTxHandler(m)
}
