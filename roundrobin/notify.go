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

func count(msg diameter.Message, dct diameter.Direction, err error) {
	if msg.FlgR {
		if dct == diameter.Rx {
			statistics.RxReq++
			if _, ok := err.(diameter.RejectRxMessage); ok {
				statistics.TxDisc++
			}
		} else {
			statistics.TxReq++
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
				statistics.RxIvld++
			} else if code < 1000 {
				statistics.RxAnsEtc++
			} else if code < 2000 {
				statistics.RxAns1xxx++
			} else if code < 3000 {
				statistics.RxAns2xxx++
			} else if code < 4000 {
				statistics.RxAns3xxx++
			} else if code < 5000 {
				statistics.RxAns4xxx++
			} else if code < 6000 {
				statistics.RxAns5xxx++
			} else {
				statistics.RxAnsEtc++
			}
		} else {
			if code < 1000 {
				statistics.TxAnsEtc++
			} else if code < 2000 {
				statistics.TxAns1xxx++
			} else if code < 3000 {
				statistics.TxAns2xxx++
			} else if code < 4000 {
				statistics.TxAns3xxx++
			} else if code < 5000 {
				statistics.TxAns4xxx++
			} else if code < 6000 {
				statistics.TxAns5xxx++
			} else {
				statistics.TxAnsEtc++
			}
		}
	}
}
