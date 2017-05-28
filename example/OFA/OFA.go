package main

import (
	"log"
	"os"

	"github.com/fkgi/diameter/example/common"
	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/ts29338"
)

/*
 <OFA> ::= < Diameter Header: 8388645, PXY, 16777313 >
           < Session-Id >
           [ Vendor-Specific-Application-Id ]
           [ Result-Code ]
           [ Experimental-Result ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
         * [ Supported-Features ]
           [ SM-Delivery-Failure-Cause ]
           [ SM-RP-UI ]
         * [ AVP ]
         * [ Failed-AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/

func main() {
	common.Log = log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds)
	h := &common.Handler{}

	h.Init("responder.in", "SMSC01.test.com", "MME01.test.com")

	h.Pull(func(r *msg.Message) *msg.Message {
		m := ts29338.MOForwardShortMessageAnswer{
			SessionID:   h.SessionID,
			ResultCode:  new(msg.ResultCode),
			OriginHost:  h.OriginHost,
			OriginRealm: h.OriginRealm,
		}
		*m.ResultCode = msg.DiameterSuccess

		msg := m.Encode()
		return &msg
	})
}
