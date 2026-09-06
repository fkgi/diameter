package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/dictionary"
)

func init() {
	diameter.ConnectionUpNotify = func(c *diameter.Connection) {
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
	dictionary.NotifyHandlerError = func(proto, msg string) {
		log.Println("[ERROR]", "error in", proto, "with reason", msg)
	}
	diameter.TraceEvent = func(old, new, event string, err error) {
		log.Printf("[INFO] diameter state update: %s->%s by event %s: error=%v",
			old, new, event, err)
	}
	diameter.TraceMessage = func(msg diameter.Message, dct diameter.Direction, err error) {
		count(msg, dct, err)
		buf := new(strings.Builder)
		fmt.Fprintf(buf, "%s diameter message handling: error=%v", dct, err)
		fmt.Fprintln(buf)
		fmt.Fprint(buf, dictionary.TraceMessageVarbose("| ", msg))
		log.Print("[INFO] ", buf)
	}
}
