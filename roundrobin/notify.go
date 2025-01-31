package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/connector"
	"github.com/fkgi/diameter/dictionary"
)

func init() {
	connector.TransportInfoNotify = func(src, dst net.Addr) {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "Detected transport address")
		if src != nil {
			fmt.Fprintf(buf, "| local: %s://%s\n", src.Network(), src.String())
		}
		if dst != nil {
			fmt.Fprintf(buf, "| peer : %s://%s\n", dst.Network(), dst.String())
		}
		log.Print(buf)
	}
	connector.TransportUpNotify = func(src, dst net.Addr) {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "Transport connection up")
		fmt.Fprintf(buf, "| local: %s://%s\n", src.Network(), src.String())
		fmt.Fprintf(buf, "| peer : %s://%s\n", dst.Network(), dst.String())
		log.Print(buf)
	}

	diameter.ConnectionUpNotify = func(c *diameter.Connection) {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "Diameter connection up")
		fmt.Fprintln(buf, "| local host/realm:", diameter.Host, "/", diameter.Realm)
		fmt.Fprintln(buf, "| peer  host/realm:", c.Host, "/", c.Realm)
		fmt.Fprint(buf, "| available application: ")
		for _, ap := range c.AvailableApplications() {
			for _, v := range dicData.V {
				for _, app := range v.P {
					if app.I == ap {
						fmt.Fprintf(buf, "%s(%d), ", app.N, ap)
					}
				}
			}
		}
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
	dictionary.NotifyHandlerError = func(proto, msg string) {
		log.Println("error in", proto, "with reason", msg)
	}

}
