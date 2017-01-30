package main

import (
	"log"

	"github.com/fkgi/diameter/example"
	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/ts29338"
)

/*
 <OFR> ::= < Diameter Header: 8388645, REQ, PXY, 16777313>
		   < Session-Id >
		   [ Vendor-Specific-Application-Id ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ Destination-Host ]
           { Destination-Realm }
           { SC-Address }
           [ OFR-Flags ]
		 * [ Supported-Features ]
           { User-Identifier }
           { SM-RP-UI }
           [ SMSMI-Correlation-ID ]
           [ SM-Delivery-Outcome ]
         * [ AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	h := &example.Handler{}

	h.Init("initiator.out", "MME01.test.com", "SMSC01.test.com")

	h.Push(func() *msg.Message {
		m := msg.Message{}
		m.Ver = msg.DiaVer
		m.FlgR = true
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
		// Auth-Session-State
		avps = append(avps, msg.AuthSessionState(false))
		// Origin-Host
		avps = append(avps, msg.OriginHost(h.OrigHost))
		// Origin-Realm
		avps = append(avps, msg.OriginRealm(h.OrigRealm))
		// Destination-Realm
		avps = append(avps, msg.DestinationRealm(h.DestRealm))

		// address of SMS-SC (SC-Address AVP)
		avps = append(avps, ts29338.ServiceCenterAddress(
			[]byte("819099990000")))

		// << functional information >>
		// capability or status of SMSC (OFR-Flags AVP)
		// avps = append(avps, msg.NewAVP_OFRFlags(true))

		// address of user (User-Identifier AVP)
		avps = append(avps, ts29338.UserIdentifier(
			"440019011112222", "819011112222", "", 0))

		// SMS data (SM-RP-UI AVP)
		avps = append(avps, ts29338.SMRPUI(
			[]byte{0x21, 0x8f, 0x0b, 0x81, 0x90, 0x90, 0x99, 0x19, 0x17, 0xf1, 0x00, 0x08, 0x06, 0x00, 0x31, 0x00, 0x2d, 0x00, 0x31}))

		//  (SMSMI-Correlation-ID AVP)
		// avps = append(avps, msg.NewAVP_SMSMICorrelationID())
		//  (SM-Delivery-Outcome AVP)
		// avps = append(avps, msg.NewAVP_SMDeliveryOutcome())

		m.Encode(avps)
		return &m
	})
}
