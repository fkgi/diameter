package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/connector"
	"github.com/fkgi/diameter/dictionary"
	"github.com/fkgi/diameter/multiplexer"
)

const apiPath = "/msg/v1/"

var rxPath string

func main() {
	log.Printf("booting Round-Robin Diameter debugger <%s REV.%d>...",
		diameter.ProductName, diameter.FirmwareRev)

	host, err := os.Hostname()
	if err != nil {
		host = "hub.internal"
	}
	dl := flag.String("diameter-local", host,
		"Diameter local host with format [tcp|sctp://][realm/]hostname[:port].")
	dp := flag.String("diameter-peer", "",
		"Diameter peer host to connect with format as same as -diameter-local.")
	hl := flag.String("http-local", "", "HTTP local host with format [host][:port].")
	hp := flag.String("http-peer", "", "HTTP peer host with format host[:port].")
	dic := flag.String("dictionary", "dictionary.json", "Diameter dictionary file path.")
	flag.Parse()

	log.Println("loading dictionary file", *dic)
	var dicData dictionary.Dictionary
	if data, err := os.ReadFile(*dic); err != nil {
		log.Fatalln("failed to open dictionary file:", err)
	} else if dicData, err = dictionary.LoadDictionary(data); err != nil {
		log.Fatalln("failed to read dictionary file:", err)
	}

	buf := new(strings.Builder)
	fmt.Fprintln(buf, "supported data")
	for vn, vnd := range dicData {
		fmt.Fprintf(buf, "| vendor: %s(%d)", vn, vnd.ID)
		fmt.Fprintln(buf)
		for an, app := range vnd.Apps {
			fmt.Fprintf(buf, "| | application: %s(%d)", an, app.ID)
			fmt.Fprintln(buf)
			fmt.Fprint(buf, "| | | command:")
			for cn, cmd := range app.Cmds {
				fmt.Fprintf(buf, " %s(%d)", cn, cmd.ID)
			}
			fmt.Fprintln(buf)
		}
		fmt.Fprint(buf, "| | AVP:")
		for an, avp := range vnd.Avps {
			fmt.Fprintf(buf, " %s(%d,%s)", an, avp.ID, avp.Type)
		}
		fmt.Fprintln(buf)
	}
	log.Print(buf)

	var router diameter.Router
	if *dp == "" {
		router = multiplexer.DefaultRouter
	} else {
		router = connector.DefaultRouter
	}
	for vn, vnd := range dicData {
		if vnd.ID == 0 {
			continue
		}
		for an, app := range vnd.Apps {
			for cn, cmd := range app.Cmds {
				registerHandler(apiPath+vn+"/"+an+"/"+cn,
					cmd.ID, app.ID, vnd.ID, router)
			}
		}
	}
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		httpErr("not found", "invalid URI path", http.StatusNotFound, w)
	})

	if *hp == "" {
		log.Println("no HTTP peer defined, Rx message will reject")
	} else {
		rxPath = "http://" + *hp
		_, err = url.Parse(rxPath)
		if err != nil {
			log.Println("invalid HTTP peer host, Rx message will reject")
			rxPath = ""
		} else {
			log.Println("HTTP peer:", rxPath)
		}
	}
	if *hl == "" {
		log.Println("no HTTP local port defined, Tx message is not available")
	} else {
		log.Println("listening HTTP local port:", *hl)
		go func() {
			err := http.ListenAndServe(*hl, nil)
			if err != nil {
				log.Fatalln("failed to listen HTTP:", err)
			}
		}()
	}

	connector.TermSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM, os.Interrupt}
	diameter.ConnectionUpNotify = func(c *diameter.Connection) {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "Diameter connection up")
		fmt.Fprintln(buf, "| local host/realm:", diameter.Host, "/", diameter.Realm)
		fmt.Fprintln(buf, "| peer host/realm: ", c.Host, "/", c.Realm)
		log.Print(buf)
	}
	diameter.TraceEvent = func(old, new, event string, err error) {
		log.Println("Diameter state update:",
			old, "->", new, "by event", event, "with error", err)
	}
	diameter.TraceMessage = func(msg diameter.Message, dct diameter.Direction, err error) {
		buf := new(strings.Builder)
		fmt.Fprintf(buf, "%s diameter message handling: error=%v", dct, err)
		fmt.Fprintln(buf)
		fmt.Fprint(buf, dictionary.TraceMessageVarbose("| ", msg))
		log.Print(buf)
	}

	if len(*dp) != 0 {
		log.Println("connecting Diameter...")
		log.Println("closed, error=", connector.DialAndServe(*dl, *dp))
	} else {
		log.Println("listening Diameter...")
		log.Println("closed, error=", multiplexer.ListenAndServe(*dl))
	}
}
