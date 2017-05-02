package common

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fkgi/diameter/connection"
	"github.com/fkgi/diameter/msg"
)

// ListenUnixSock start listening unixsock
func ListenUnixSock(isock, osock string) (net.Listener, net.Listener) {
	Log.Println("open unixsock ...")
	il, e := net.Listen("unix", isock)
	if e != nil {
		Log.Fatalln(e)
	}
	Log.Printf("  input path  =%s", isock)
	ol, e := net.Listen("unix", osock)
	if e != nil {
		Log.Fatalln(e)
	}
	Log.Printf("  output path =%s", osock)
	return il, ol
}

//CloseUnixSock remove unixsock
func CloseUnixSock(isock, osock string) {
	Log.Println("closing unixsock ...")
	os.Remove(isock)
	os.Remove(osock)
}

// RunUnixsockRelay relay data between unix-socket and diameter-link
func RunUnixsockRelay(prov *connection.Connection, il, ol net.Listener) {
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)

	// recieve diameter request
	go func() {
		for {
			c, e := il.Accept()
			if e != nil {
				Log.Println(e)
				sigc <- os.Interrupt
				return
			}
			defer c.Close()

			m, ch, e := prov.Recieve()
			if e != nil {
				Log.Println(e)
				continue
			}
			if _, e = m.WriteTo(c); e != nil {
				Log.Println(e)
				continue
			}
			if _, e = m.ReadFrom(c); e != nil {
				Log.Println(e)
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
				Log.Println(e)
				sigc <- os.Interrupt
				return
			}
			defer c.Close()

			m := msg.Message{}
			if _, e = m.ReadFrom(c); e != nil {
				Log.Println(e)
				continue
			}
			if avp, e := m.Decode(); e != nil {
				Log.Println(e)
				continue
			} else {
				// add RouteRecord AVP
				var src msg.DiameterIdentity
				for _, a := range avp {
					if a.Code == 264 {
						e = a.Decode(&src)
						if e != nil {
							Log.Println(e)
						}
						break
					}
				}
				avp = append(avp, msg.RouteRecord(src).Encode())
				if e = m.Encode(avp); e != nil {
					Log.Println(e)
					continue
				}
			}
			if m, e = prov.Send(m); e != nil {
				c.Close()
				Log.Println(e)
				continue
			}
			if _, e = m.WriteTo(c); e != nil {
				c.Close()
				Log.Println(e)
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
}
