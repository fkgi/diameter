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

func main() {
	common.Log = log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds)
	h := &common.Handler{}

	h.Init("initiator.out", "MME01.test.com", "SMSC01.test.com")

	h.Push(func() *msg.Message {
		msisdn, _ := ts29329.ParseMSISDN("819011112222")

		m := ts29338.MOForwardShortMessageRequest{
			SessionID:        h.SessionID,
			OriginHost:       h.OriginHost,
			OriginRealm:      h.OriginRealm,
			DestinationRealm: h.DestinationRealm,
			SCAddress:        ts29338.SCAddress("819099990000"),
			OFRFlags:         new(ts29338.OFRFlags),
			UserIdentifier: ts29336.UserIdentifier{
				UserName: msg.UserName("440019011112222"),
				MSISDN:   msisdn},
			SMRPUI: ts29338.SMRPUI([]byte{
				0x21, 0x8f, 0x0b, 0x81, 0x90, 0x90, 0x99, 0x19,
				0x17, 0xf1, 0x00, 0x08, 0x06, 0x00, 0x31, 0x00,
				0x2d, 0x00, 0x31})}
		m.OFRFlags.S6as6d = false

		//  (SMSMI-Correlation-ID AVP)
		// avps = append(avps, msg.NewAVP_SMSMICorrelationID())
		//  (SM-Delivery-Outcome AVP)
		// avps = append(avps, msg.NewAVP_SMDeliveryOutcome())

		msg := m.Encode()
		return &msg
	})
}
