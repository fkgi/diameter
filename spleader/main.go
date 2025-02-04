package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/fkgi/diameter"
)

var upLink diameter.Identity

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "spleader.internal"
	}
	dlocal := flag.String("l", hostname,
		"Diameter local host. `[(tcp|sctp)://][realm/]hostname[:port]`")
	hlocal := flag.String("i", ":12001",
		"HTTP local interface address. `[host]:port`")
	flag.Parse()

	upLink, err = diameter.ParseIdentity(flag.Arg(0))
	if err != nil {
		log.Fatalln("invalid uplink peer hostname:", err)
	}
	log.Printf("uplink peer hostname is %s", upLink)

	log.Printf("booting spleader for Round-Robin <%s REV.%d>...",
		diameter.ProductName, diameter.FirmwareRev)

	http.HandleFunc("/diastate/v1/connection", conStateHandler)
	log.Println("listening HTTP local port:", *hlocal)
	go func() {
		err := http.ListenAndServe(*hlocal, nil)
		if err != nil {
			log.Println("failed to listen HTTP, API is not available:", err)
		}
	}()

	diameter.DefaultRxHandler = rxhandler

	log.Println("listening Diameter...")
	log.Println("closed, error=", ListenAndServe(*dlocal))
}
