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
	"time"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/connector"
	"github.com/fkgi/diameter/dictionary"
)

const apiPath = "/diamsg/v1/"

var rxPath string
var dicData dictionary.XDictionary

func main() {
	host, err := os.Hostname()
	if err != nil {
		host = "roundrobin.internal"
	}
	dlocal := flag.String("l", host,
		"Diameter local host. `[realm/]hostname[:port]`")
	hlocal := flag.String("i", ":8080",
		"HTTP local interface address. `[host]:port`")
	hpeer := flag.String("b", "localhost",
		"HTTP backend host address. `host[:port]`")
	dict := flag.String("d", "dictionary.xml",
		"Diameter dictionary file `path`.")
	cause := flag.String("c", "rebooting",
		"Disconnect cause in sending DPR. `rebooting|busy|do_not_want_to_talk_to_you`")
	server := flag.Bool("s", false, "Run as server")
	to := flag.Int("t", int(diameter.WDInterval/time.Second),
		"message timeout timer [s]")
	help := flag.Bool("h", false, "Print usage")
	flag.Parse()

	dpeer := flag.Arg(0)
	if *help || dpeer == "" {
		fmt.Printf("Usage: %s [OPTION]... DIAMETER_PEER\n", os.Args[0])
		fmt.Println("DIAMETER_PEER format is [(tcp|sctp)://][realm/]hostname[:port]")
		fmt.Println()
		flag.PrintDefaults()
		return
	}

	log.Printf("booting Round-Robin debugger for Diameter <%s REV.%d>...",
		diameter.ProductName, diameter.FirmwareRev)

	diameter.WDInterval = time.Duration(*to) * time.Second

	log.Println("loading dictionary file", *dict)
	if data, err := os.ReadFile(*dict); err != nil {
		log.Fatalln("failed to open dictionary file:", err)
	} else if dicData, err = dictionary.LoadDictionary(data); err != nil {
		log.Fatalln("failed to read dictionary file:", err)
	} else {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "supported data")
		for _, vnd := range dicData.V {
			fmt.Fprintf(buf, "| vendor: %s(%d)", vnd.N, vnd.I)
			fmt.Fprintln(buf)
			for _, app := range vnd.P {
				fmt.Fprintf(buf, "| | application: %s(%d)", app.N, app.I)
				fmt.Fprintln(buf)
				fmt.Fprint(buf, "| | | command:")
				for _, cmd := range app.C {
					fmt.Fprintf(buf, " %s(%d)", cmd.N, cmd.I)
				}
				fmt.Fprintln(buf)
			}
			fmt.Fprint(buf, "| | AVP:")
			for _, avp := range vnd.V {
				fmt.Fprintf(buf, " %s(%d,%s)", avp.N, avp.I, avp.T)
			}
			fmt.Fprintln(buf)
		}
		log.Print(buf)
	}

	http.HandleFunc("/diastate/v1/connection", conStateHandler)
	http.HandleFunc("/diastate/v1/statistics", statsHandler)

	rxPath = "http://" + *hpeer
	_, err = url.Parse(rxPath)
	if err != nil {
		log.Println("invalid HTTP backend host, Rx request will be rejected")
		rxPath = ""
	} else {
		log.Println("HTTP backend:", rxPath)
	}
	dicData.RegisterHandler(rxPath, apiPath, connector.DefaultRouter)

	log.Println("listening HTTP local port:", *hlocal)
	go func() {
		err := http.ListenAndServe(*hlocal, nil)
		if err != nil {
			log.Println("failed to listen HTTP, Tx request is not available:", err)
		}
	}()

	connector.TermSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM, os.Interrupt}
	switch *cause {
	case "rebooting":
		connector.TermCause = diameter.Rebooting
	case "busy":
		connector.TermCause = diameter.Busy
	case "do_not_want_to_talk_to_you":
		connector.TermCause = diameter.DoNotWantToTalkToYou
	default:
		connector.TermCause = diameter.Rebooting
	}

	if *server {
		log.Println("listening Diameter...")
		log.Println("closed, error=", connector.ListenAndServe(*dlocal, dpeer))
	} else {
		log.Println("connecting Diameter...")
		log.Println("closed, error=", connector.DialAndServe(*dlocal, dpeer))
	}
}
