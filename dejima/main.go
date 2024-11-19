package main

import (
	"flag"
	"fmt"
	"log"
	"net"
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
	log.Printf("booting Round-Robin debugger Connector for Diameter <%s REV.%d>...",
		diameter.ProductName, diameter.FirmwareRev)

	host, err := os.Hostname()
	if err != nil {
		host = "roundrobin.internal"
	}
	dlocal := flag.String("diameter-local", host,
		"Diameter local host uri. `[(tcp|sctp)://][realm/]hostname[:port]`")
	dpeer := flag.String("diameter-peer", "",
		"Diameter peer host uri to connect. `[(tcp|sctp)://][realm/]hostname[:port]`")
	hlocal := flag.String("http-local", "",
		"HTTP local interface. `[host][:port]`")
	hpeer := flag.String("http-peer", "",
		"HTTP backend host. `host[:port]`")
	dict := flag.String("dictionary", "dictionary.json",
		"Diameter dictionary file path.")
	cause := flag.String("cause", "",
		"Disconnect cause in sending DPR. `rebooting|busy|do_not_want_to_talk_to_you`")
	help := flag.Bool("help", false, "Print usage")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	var router diameter.Router
	if *dpeer == "" {
		router = multiplexer.DefaultRouter
	} else {
		router = connector.DefaultRouter
	}

	log.Println("loading dictionary file", *dict)
	var dicData dictionary.XDictionary
	if data, err := os.ReadFile(*dict); err != nil {
		log.Fatalln("failed to open dictionary file:", err)
	} else if dicData, err = dictionary.LoadDictionary(data); err != nil {
		log.Fatalln("failed to read dictionary file:", err)
	}

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
	dicData.RegisterHandler(rxPath, apiPath, router)

	if *hpeer == "" {
		log.Println("no HTTP peer defined, Rx request will be rejected")
	} else {
		rxPath = "http://" + *hpeer
		_, err = url.Parse(rxPath)
		if err != nil {
			log.Println("invalid HTTP peer host, Rx request will be rejected")
			rxPath = ""
		} else {
			log.Println("HTTP peer:", rxPath)
		}
	}
	if *hlocal == "" {
		log.Println("no HTTP local port defined, Tx request is not available")
	} else {
		log.Println("listening HTTP local port:", *hlocal)
		go func() {
			err := http.ListenAndServe(*hlocal, nil)
			if err != nil {
				log.Fatalln("failed to listen HTTP:", err)
			}
		}()
	}

	connector.TermSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM, os.Interrupt}
	switch *cause {
	case "rebooting":
		connector.TermCause = diameter.Rebooting
	case "busy":
		connector.TermCause = diameter.Busy
	case "do_not_want_to_talk_to_you":
		connector.TermCause = diameter.DoNotWantToTalkToYou
	default:
		connector.TermCause = connector.UndefinedCause
	}

	connector.TransportInfoNotify = func(src, dst net.Addr) {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "Transport address")
		if src != nil {
			fmt.Fprintln(buf, "| local:", src.Network(), src.String())
		}
		if dst != nil {
			fmt.Fprintln(buf, "| peer :", dst.Network(), dst.String())
		}
		log.Print(buf)
	}

	diameter.ConnectionUpNotify = func(c *diameter.Connection) {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "Diameter connection up")
		fmt.Fprintln(buf, "| local host/realm:", diameter.Host, "/", diameter.Realm)
		fmt.Fprintln(buf, "| peer  host/realm:", c.Host, "/", c.Realm)
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

	if len(*dpeer) != 0 {
		log.Println("connecting Diameter...")
		log.Println("closed, error=", connector.DialAndServe(*dlocal, *dpeer))
	} else {
		log.Println("listening Diameter...")
		log.Println("closed, error=", multiplexer.ListenAndServe(*dlocal))
	}
}
