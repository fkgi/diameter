package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/connector"
)

func init() {
	connector.TransportInfoNotify = func(src, dst net.Addr) {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "Detected transport address")
		if src != nil {
			fmt.Fprintf(buf, "| local: %s://%s\n", src.Network(), src.String())
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
		fmt.Fprintf(buf, "| applications:    %d\n", c.AvailableApplications())
		log.Print(buf)
	}
	diameter.TraceEvent = func(old, new, event string, err error) {
		if old != new || err != nil {
			log.Println("Diameter state update:",
				old, "->", new, "by event", event, "with error", err)
		}
	}
	diameter.TraceMessage = func(msg diameter.Message, dct diameter.Direction, err error) {
		if msg.AppID != 0 {
			log.Printf("%s diameter message handling: peer=%s, error=%v\n",
				dct, msg.PeerName, err)
		}
	}

}
