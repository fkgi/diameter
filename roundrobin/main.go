package main

import (
	"flag"
	"log"
	"os"
	"syscall"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/connector"
	"github.com/fkgi/diameter/dictionary"
)

func main() {
	log.Printf("booting RR <%s REV.%d>...",
		diameter.ProductName, diameter.FirmwareRev)

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

	rxPath = "http://" + *r + "/msg/v1/"

	if data, err := os.ReadFile(*d); err != nil {
		log.Fatalln("failed to open dictionary file", err)
	} else if err = dictionary.LoadDictionary(data); err != nil {
		log.Fatalln("failed to read dictionary file", err)
	}

	diameter.DefaultRxHandler = handleRx
	listenAndServeHttp(*t)

	connector.TermSignals = []os.Signal{
		syscall.SIGINT, syscall.SIGTERM, os.Interrupt}
	if *p == "" {
		log.Println("listening...")
		log.Println("closed, error=", connector.ListenAndServe(*l, *s))
	} else {
		log.Println("connecting...")
		log.Println("closed, error=", connector.DialAndServe(*l, *p, *s))
	}
}
