package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/dictionary"
)

func init() {
	dictionary.TraceTxHttpRequest = func(p string, tx []byte, c int, rx []byte, e error) {
		traceHTTPHandling("Tx", p, tx, c, rx, e)
	}
	dictionary.TraceRxHttpRequest = func(p string, tx []byte, c int, rx []byte, e error) {
		traceHTTPHandling("Rx", p, tx, c, rx, e)
	}

	diameter.TraceEvent = func(
		c *diameter.Connection, old, event string, err error) {
		n := c.State()
		log.Printf("[INFO] diameter state update: %s->%s by event %s, peer=%s, error=%v",
			old, n, event, c.Host, err)

		if old != "open" && n == "open" {
			buf := new(strings.Builder)
			fmt.Fprint(buf, "diameter connection up")
			fmt.Fprintf(buf, "\n| local host/realm: %s/%s", diameter.Host, diameter.Realm)
			fmt.Fprintf(buf, "\n| peer  host/realm: %s/%s", c.Host, c.Realm)
			fmt.Fprint(buf, "\n| available application: ")
			for _, ap := range c.AvailableApplications() {
				fmt.Fprintf(buf, "%s(%d), ", dictionary.GetApplicationName(ap), ap)
			}
			log.Print("[INFO] ", buf)
		}
	}
	diameter.TraceTxMessage = func(
		c *diameter.Connection, msg diameter.Message, err error) {
		count(msg, false, err)
		traceMessage(c.Host.String(), "Tx", msg, err)
	}
	diameter.TraceRxMessage = func(
		c *diameter.Connection, msg diameter.Message, err error) {
		count(msg, true, err)
		traceMessage(c.Host.String(), "Rx", msg, err)
	}
}

func traceMessage(host, direction string, msg diameter.Message, err error) {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "%s diameter message handling: error=%v", direction, err)
	fmt.Fprintln(buf, "\n| Peer:", host)
	fmt.Fprint(buf, dictionary.TraceMessageVarbose("| ", msg))
	log.Print("[INFO] ", buf)
}

func traceHTTPHandling(d string, p string, tx []byte, c int, rx []byte, e error) {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "%s HTTP request handling: error=%v\n", d, e)
	fmt.Fprintln(buf, "| path:", p)
	fmt.Fprintln(buf, "| body:", string(tx))
	fmt.Fprintln(buf, "| result:", c, http.StatusText(c))
	fmt.Fprintln(buf, "| body:", string(rx))
	log.Print("[INFO] ", buf)
}
