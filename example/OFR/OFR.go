package main

import (
	"log"
	"os"

	"github.com/fkgi/diameter/example/common"
	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/ts29329"
	"github.com/fkgi/diameter/ts29336"
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
	common.Log = log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds)
	h := &common.Handler{}

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
		avps = append(avps, h.SessionID.Encode())
		// Vendor-Specific-Application-Id
		avps = append(avps, msg.VendorSpecificApplicationID{
			VendorID: 10415,
			App:      msg.AuthApplicationID(16777313)}.Encode())
		// Auth-Session-State
		avps = append(avps, msg.AuthSessionState(msg.StateNotMaintained).Encode())
		// Origin-Host
		avps = append(avps, h.OriginHost.Encode())
		// Origin-Realm
		avps = append(avps, h.OriginRealm.Encode())
		// Destination-Realm
		avps = append(avps, h.DestinationRealm.Encode())

		// address of SMS-SC (SC-Address AVP)
		avps = append(avps, ts29338.SCAddress("819099990000").Encode())

		// << functional information >>
		// capability or status of SMSC (OFR-Flags AVP)
		avps = append(avps, ts29338.OFRFlags{S6as6d: false}.Encode())

		// address of user (User-Identifier AVP)
		msisdn, _ := ts29329.ParseMSISDN("819011112222")
		avps = append(avps, ts29336.UserIdentifier{
			UserName: msg.UserName("440019011112222"),
			MSISDN:   msisdn}.Encode())

		// SMS data (SM-RP-UI AVP)
		avps = append(avps, ts29338.SMRPUI([]byte{
			0x21, 0x8f, 0x0b, 0x81, 0x90, 0x90, 0x99, 0x19,
			0x17, 0xf1, 0x00, 0x08, 0x06, 0x00, 0x31, 0x00,
			0x2d, 0x00, 0x31}).Encode())

		//  (SMSMI-Correlation-ID AVP)
		// avps = append(avps, msg.NewAVP_SMSMICorrelationID())
		//  (SM-Delivery-Outcome AVP)
		// avps = append(avps, msg.NewAVP_SMDeliveryOutcome())

		m.Encode(avps)
		return &m
	})
}
