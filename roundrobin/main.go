package main

import (
	"flag"
	"log"
	"os"
	"syscall"

	"github.com/fkgi/diameter"
)

func main() {
	log.Printf("booting RR <%s REV.%d>...", diameter.ProductName, diameter.FirmwareRev)

	host, err := os.Hostname()
	if err != nil {
		host = "hub.internal"
	}
	l := flag.String("l", host+":3868", "diameter local host:port")
	p := flag.String("p", "", "diameter peer host:port for dial")
	t := flag.String("t", ":8080", "http local host:port")
	r := flag.String("r", "localhost:8081", "http peer host:port")
	s := flag.Bool("sctp", false, "flag for sctp")
	d := flag.String("d", "dictionary.json", "diameter dictionary")
	flag.Parse()

	rxPath = "http://" + *r

	if err = loadDictionary(*d); err != nil {
		log.Fatalln(err)
	}

	listenAndServeHttp(*t)

	diameter.TermSignals = []os.Signal{
		syscall.SIGINT, syscall.SIGTERM, os.Interrupt}
	diameter.Router = true
	if *p == "" {
		log.Println("listening...")
		log.Println("closed, error=", diameter.ListenAndServe(*l, *s))
	} else {
		log.Println("connecting...")
		log.Println("closed, error=", diameter.DialAndServe(*l, *p, *s))
	}
}
