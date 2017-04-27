package common

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fkgi/diameter/connection"
	"github.com/fkgi/diameter/msg"
)

// RunUnixsockRelay relay data between unix-socket and diameter-link
func RunUnixsockRelay(prov *connection.Connection, isock, osock string) {
	il, e := net.Listen("unix", isock)
	if e != nil {
		log.Fatalln(e)
	}
	ol, e := net.Listen("unix", osock)
	if e != nil {
		log.Fatalln(e)
	}
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)

	// recieve diameter request
	go func() {
		for {
			c, e := il.Accept()
			if e != nil {
				log.Println(e)
				sigc <- os.Interrupt
				return
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

	// send diameter request
	go func() {
		for {
			c, e := ol.Accept()
			if e != nil {
				log.Println(e)
				sigc <- os.Interrupt
				return
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
				// add RouteRecord AVP
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
				avp = append(avp, msg.RouteRecord(src).Encode())
				if e = m.Encode(avp); e != nil {
					log.Println(e)
					continue
				}
			}
			if m, e = prov.Send(m); e != nil {
				c.Close()
				log.Println(e)
				continue
			}
			if _, e = m.WriteTo(c); e != nil {
				c.Close()
				log.Println(e)
				continue
			}
		}
	}()

	<-sigc
	Log.Println("shutdown ...")

	prov.Close(msg.Rebooting)
	for prov.State() != "Shutdown" {
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
	os.Remove(isock)
	os.Remove(osock)
}
