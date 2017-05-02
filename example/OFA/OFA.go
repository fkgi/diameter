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
		avps = append(avps, h.SessionID.Encode())
		// Vendor-Specific-Application-Id
		avps = append(avps, msg.VendorSpecificApplicationID{
			VendorID: 10415,
			App:      msg.AuthApplicationID(16777313)}.Encode())
		// Result-Code AVP
		avps = append(avps, msg.DiameterSuccess.Encode())
		// Auth-Session-State
		avps = append(avps, msg.AuthSessionState(msg.StateNotMaintained).Encode())
		// Origin-Host
		avps = append(avps, h.OriginHost.Encode())
		// Origin-Realm
		avps = append(avps, h.OriginRealm.Encode())

		// SMS data (SM-RP-UI AVP)
		avps = append(avps, ts29338.SMRPUI([]byte{
			0x21, 0x8f, 0x0b, 0x81, 0x90, 0x90, 0x99, 0x19,
			0x17, 0xf1, 0x00, 0x08, 0x06, 0x00, 0x31, 0x00,
			0x2d, 0x00, 0x31}).Encode())

		m.Encode(avps)
		return &m
	})
}
