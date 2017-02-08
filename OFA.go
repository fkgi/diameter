package main

import (
	"log"

	"github.com/fkgi/diameter/example"
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
           [ SM-Delivery- Failure-Cause ]
           [ SM-RP-UI ]
         * [ AVP ]
         * [ Failed-AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	h := &example.Handler{}

	h.Init("responder.in", "SMSC01.test.com", "MME01.test.com")

	h.Pull(func(r *msg.Message) *msg.Message {
		m := msg.Message{}
		m.Ver = msg.DiaVer
		m.FlgR = false
		m.FlgP = true
		m.FlgE = false
		m.FlgT = false
		m.Code = uint32(8388645)
		m.AppID = 16777313

		var avps []msg.Avp
		// Session-Id
		avps = append(avps, msg.SessionID(h.SessionID))
		// Vendor-Specific-Application-Id
		avps = append(avps, msg.VendorSpecificApplicationID(
			10415, 16777313, true))
		// Result-Code AVP
		avps = append(avps, msg.ResultCode(2000))
		// Auth-Session-State
		avps = append(avps, msg.AuthSessionState(false))
		// Origin-Host
		avps = append(avps, msg.OriginHost(h.OrigHost))
		// Origin-Realm
		avps = append(avps, msg.OriginRealm(h.OrigRealm))

		// SMS data (SM-RP-UI AVP)
		avps = append(avps, ts29338.SMRPUI([]byte{
			0x21, 0x8f, 0x0b, 0x81, 0x90, 0x90, 0x99, 0x19,
			0x17, 0xf1, 0x00, 0x08, 0x06, 0x00, 0x31, 0x00,
			0x2d, 0x00, 0x31}))

		//  (SMSMI-Correlation-ID AVP)
		// avps = append(avps, msg.NewAVP_SMSMICorrelationID())
		//  (SM-Delivery-Outcome AVP)
		// avps = append(avps, msg.NewAVP_SMDeliveryOutcome())

		m.Encode(avps)
		return &m
	})
}
