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
		fmt.Fprintln(buf, "detected transport address")
		if src != nil {
			fmt.Fprintf(buf, "| local: %s://%s\n", src.Network(), src.String())
		}
		log.Print("[INFO] ", buf)
	}
	connector.TransportUpNotify = func(src, dst net.Addr) {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "transport connection up")
		fmt.Fprintf(buf, "| local: %s://%s\n", src.Network(), src.String())
		fmt.Fprintf(buf, "| peer : %s://%s\n", dst.Network(), dst.String())
		log.Print("[INFO] ", buf)
	}

	diameter.ConnectionUpNotify = func(c *diameter.Connection) {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "diameter connection up")
		fmt.Fprintln(buf, "| local host/realm:", diameter.Host, "/", diameter.Realm)
		fmt.Fprintln(buf, "| peer  host/realm:", c.Host, "/", c.Realm)
		fmt.Fprintf(buf, "| applications:    %d\n", c.AvailableApplications())
		log.Print("[INFO] ", buf)
	}
	diameter.TraceEvent = func(old, new, event string, err error) {
		if old != new || err != nil {
			log.Printf("[INFO] diameter state update: %s->%s by event %s: error=%v",
				old, new, event, err)
		}
	}
	/*
		diameter.TraceMessage = func(msg diameter.Message, dct diameter.Direction, err error) {
			if msg.AppID != 0 {
				t := "answer"
				if msg.FlgR {
					t = "request"
				}
				log.Printf("[INFO] %s diameter %s message handling: peer=%s, error=%v\n",
					dct, t, msg.PeerName, err)
			}
		}
	*/
}
