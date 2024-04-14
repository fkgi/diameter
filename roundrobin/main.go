package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/connector"
	"github.com/fkgi/diameter/dictionary"
	"github.com/fkgi/diameter/multiplexer"
)

const apipath = "/msg/v1/"

var DefaultTxHandler func(diameter.Message) diameter.Message

func main() {
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

	log.Printf("booting Round-Robin Diameter debugger <%s REV.%d>...",
		diameter.ProductName, diameter.FirmwareRev)

	log.Println("loading dictionary file", *dic)
	diameter.DefaultRxHandler = handleRx
	if data, err := os.ReadFile(*dic); err != nil {
		log.Fatalln("failed to open dictionary file:", err)
	} else if d, err := dictionary.LoadDictionary(data); err != nil {
		log.Fatalln("failed to read dictionary file:", err)
	} else {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "supported data")
		for vndname, vnd := range d {
			fmt.Fprintf(buf, "| vendor: %s(%d)", vndname, vnd.ID)
			fmt.Fprintln(buf)
			for appname, app := range vnd.Apps {
				fmt.Fprintf(buf, "| | application: %s(%d)", appname, app.ID)
				fmt.Fprintln(buf)
				fmt.Fprint(buf, "| | | command:")
				for cmdname, cmd := range app.Cmds {
					fmt.Fprintf(buf, " %s(%d)", cmdname, cmd.ID)
				}
				fmt.Fprintln(buf)
			}
			fmt.Fprint(buf, "| | AVP:")
			for name, avp := range vnd.Avps {
				fmt.Fprintf(buf, " %s(%d,%s)", name, avp.ID, avp.Type)
			}
			fmt.Fprintln(buf)
		}
		log.Print(buf)
	}

	if *hp == "" {
		log.Println("no HTTP peer defined, Rx message will reject")
	} else {
		rxPath = "http://" + *hp + apipath
		log.Println("HTTP peer:", rxPath)
	}
	if *hl == "" {
		log.Println("no HTTP local port defined, Tx message is not available")
	} else {
		log.Println("listening HTTP local port:", *hl)
		go func() {
			err := http.ListenAndServe(*hl, http.HandlerFunc(handleTx))
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
		DefaultTxHandler = connector.DefaultTxHandler
		log.Println("connecting Diameter...")
		log.Println("closed, error=", connector.DialAndServe(*dl, *dp))
	} else {
		DefaultTxHandler = multiplexer.DefaultTxHandler
		log.Println("listening Diameter...")
		log.Println("closed, error=", multiplexer.ListenAndServe(*dl))
	}
}

func formatAVPs(avps []diameter.AVP) (map[string]any, error) {
	result := make(map[string]any)
	for _, a := range avps {
		n, v, e := dictionary.DecodeAVP(a)
		if e != nil {
			return nil, e
		}
		result[n] = v
	}
	return result, nil
}

func parseAVPs(d map[string]any) ([]diameter.AVP, error) {
	avps := make([]diameter.AVP, 0, 10)
	for k, v := range d {
		a, e := dictionary.EncodeAVP(k, v)
		if e != nil {
			return nil, fmt.Errorf("%s is invalid: %v", k, e)
		}
		avps = append(avps, a)
	}
	return avps, nil
}
