package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/fkgi/diameter"
)

func init() {
	diameter.ConnectionUpNotify = func(c *diameter.Connection) {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "diameter connection up")
		fmt.Fprintln(buf, "| peer host/realm:", c.Host, "/", c.Realm)
		fmt.Fprintf(buf, "| applications:    %d\n", c.AvailableApplications())
		log.Print("[INFO] ", buf)
	}
	diameter.TraceEvent = func(old, new, event string, err error) {
		if old != new || err != nil {
			log.Printf("[INFO] diameter state update: %s->%s by event %s: error=%v",
				old, new, event, err)
		}
	}
}
