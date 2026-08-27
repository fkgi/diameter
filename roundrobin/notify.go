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
		fmt.Fprintln(buf, "diameter connection up")
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
func count(msg diameter.Message, dct diameter.Direction, err error) {
	if msg.FlgR {
		if dct == diameter.Rx {
			statistics.rxReq++
			if _, ok := err.(diameter.RejectRxMessage); ok {
				statistics.txDisc++
			}
		} else {
			statistics.txReq++
		}
	} else {
		var code uint32
		if avps, e := msg.GetAVP(); e == nil {
			for _, a := range avps {
				switch a.Code {
				case 268:
					code, _ = diameter.GetResultCode(a)
				case 297:
					code, _ = diameter.GetResultCode(a)
					code %= 10000
				}
				if code != 0 {
					break
				}
			}
		}
		if dct == diameter.Rx {
			if _, ok := err.(diameter.FailureAnswer); err != nil && !ok {
				statistics.rxIvld++
			} else if code < 1000 {
				statistics.rxAnsEtc++
			} else if code < 2000 {
				statistics.rxAns1xxx++
			} else if code < 3000 {
				statistics.rxAns2xxx++
			} else if code < 4000 {
				statistics.rxAns3xxx++
			} else if code < 5000 {
				statistics.rxAns4xxx++
			} else if code < 6000 {
				statistics.rxAns5xxx++
			} else {
				statistics.rxAnsEtc++
			}
		} else {
			if code < 1000 {
				statistics.txAnsEtc++
			} else if code < 2000 {
				statistics.txAns1xxx++
			} else if code < 3000 {
				statistics.txAns2xxx++
			} else if code < 4000 {
				statistics.txAns3xxx++
			} else if code < 5000 {
				statistics.txAns4xxx++
			} else if code < 6000 {
				statistics.txAns5xxx++
			} else {
				statistics.txAnsEtc++
			}
		}
	}
}
