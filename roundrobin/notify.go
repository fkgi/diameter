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
		fmt.Fprintln(buf, "detected transport address")
		if src != nil {
			fmt.Fprintf(buf, "| local: %s://%s\n", src.Network(), src.String())
		}
		if dst != nil {
			fmt.Fprintf(buf, "| peer : %s://%s\n", dst.Network(), dst.String())
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
		buf := new(strings.Builder)
		fmt.Fprintf(buf, "%s diameter message handling: error=%v", dct, err)
		fmt.Fprintln(buf)
		fmt.Fprint(buf, dictionary.TraceMessageVarbose("| ", msg))
		log.Print("[INFO] ", buf)

		if msg.FlgR {
			if dct == diameter.Rx {
				rxReq++
				if _, ok := err.(diameter.RejectRxMessage); ok {
					txDisc++
				}
			} else {
				txReq++
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
					rxIvld++
				} else if code < 1000 {
					rxAns[0]++
				} else if code < 2000 {
					rxAns[1]++
				} else if code < 3000 {
					rxAns[2]++
				} else if code < 4000 {
					rxAns[3]++
				} else if code < 5000 {
					rxAns[4]++
				} else if code < 6000 {
					rxAns[5]++
				} else {
					rxAns[0]++
				}
			} else {
				if code < 1000 {
					txAns[0]++
				} else if code < 2000 {
					txAns[1]++
				} else if code < 3000 {
					txAns[2]++
				} else if code < 4000 {
					txAns[3]++
				} else if code < 5000 {
					txAns[4]++
				} else if code < 6000 {
					txAns[5]++
				} else {
					txAns[0]++
				}
			}
		}
	}

}
