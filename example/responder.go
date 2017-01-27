package main

import (
	"log"
	"time"

	"net"

	"github.com/fkgi/diameter/example"
	"github.com/fkgi/diameter/provider"
	"github.com/fkgi/extnet"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	provider.Notificator = func(e error) { log.Println("DIAM:", e) }
	extnet.Notificator = func(e error) { log.Println("SCTP:", e) }

	// boot log
	log.Println("initiator sample booting ...")

	// get option flag
	isock, osock, conf := example.GeneratePath("responder")
	la, ln, _, pn := example.LoadConfig(conf)

	// open Diameter socket
	log.Println("start connecting Diameter connection")
	pl := provider.Listen(ln)
	prov := pl.AddPeer(pn)
	prov.Open()

	switch la.Network() {
	case "sctp", "sctp4", "sctp6":
		if a, ok := la.(*extnet.SCTPAddr); !ok {
			log.Fatalln("invalid sctp address")
		} else if lnr, e := extnet.ListenSCTP(la.Network(), a); e != nil {
			log.Fatalln(e)
		} else {
			go pl.Bind(lnr)
		}
	case "tcp", "tcp4", "tcp6":
		if a, ok := la.(*net.TCPAddr); !ok {
			log.Fatalln("invalid tcp address")
		} else if lnr, e := net.ListenTCP(la.Network(), a); e != nil {
			log.Fatalln(e)
		} else {
			go pl.Bind(lnr)
		}
	}
	time.Sleep(time.Second)

	// open UNIX socket
	log.Println("start listening on UNIX socket", isock, "and", osock)
	example.RunUnixsockRelay(prov, isock, osock)
	time.Sleep(time.Second)
}
