package main

import (
	"log"
	"time"

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
	if e := pl.Bind(la); e != nil {
		log.Fatalln(e)
	}
	time.Sleep(time.Second)

	// open UNIX socket
	log.Println("start listening on UNIX socket", isock, "and", osock)
	example.RunUnixsockRelay(prov, isock, osock)
	time.Sleep(time.Second)
}
