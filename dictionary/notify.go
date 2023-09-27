package dictionary

import (
	"fmt"
	"log"

	"github.com/fkgi/diameter"
)

func SetVarboseTrace() {
	diameter.TraceMessage = trace
}

func trace(msg diameter.Message, dct diameter.Direction, err error) {
	com, e := DecodeMessage(msg)
	if e != nil {
		com = fmt.Sprintf("unknown command(appID=%d, code=%d)", msg.AppID, msg.Code)
	}
	if msg.FlgR {
		com += "-Request"
	} else {
		com += "-Answer"
	}
	flags := ""
	if msg.FlgP {
		flags += "P"
	} else {
		flags += "-"
	}
	if msg.FlgE {
		flags += "E"
	} else {
		flags += "-"
	}
	if msg.FlgT {
		flags += "T"
	} else {
		flags += "-"
	}
	avp := ""
	avps, e := msg.GetAVP()
	if e != nil {
		avp = fmt.Sprintf(" | | HEX body=% x", msg.AVPs)
	} else {
		for _, a := range avps {
			n, v, e := DecodeAVP(a)
			if e != nil {
				avp += fmt.Sprintf(" | | unknown AVP(vendorID=%d, code=%d): % x\n",
					a.VendorID, a.Code, a.Data)
			} else {
				avp += fmt.Sprintf(" | | %s: %v\n", n, v)
			}
		}
	}

	log.Printf("%s diameter message handling: error=%v\n | %s (%s)\n | Hop-by-Hop ID=%d, End-to-End ID=%d\n%s",
		dct, err, com, flags, msg.HbHID, msg.EtEID, avp)
}
