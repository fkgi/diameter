package ts29338

import (
	"github.com/fkgi/diameter/connection"
	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/ts29336"
)

/*
MOForwardShortMessageRequest is OFR message.
 <OFR> ::= < Diameter Header: 8388645, REQ, PXY, 16777313 >
           < Session-Id >
           [ DRMP ]
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
type MOForwardShortMessageRequest struct {
	*msg.DestinationHost
	msg.DestinationRealm
	SCAddress
	*OFRFlags
	// SupportedFeatures
	ts29336.UserIdentifier
	SMRPUI
	// SMSMICorrelationID
	// SMDeliveryOutcome
}

// Encode return AVP struct of this value
func (v *MOForwardShortMessageRequest) Encode(l *connection.LocalNode) msg.Message {
	m := msg.Message{
		Ver:  msg.DiaVer,
		FlgR: true, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388645, AppID: 16777313}

	var avps []msg.Avp
	avps = append(avps, l.NextSession().Encode())
	avps = append(avps, msg.VendorSpecificApplicationID{
		VendorID: 10415,
		App:      msg.AuthApplicationID(16777313)}.Encode())
	avps = append(avps, msg.AuthSessionState(msg.StateNotMaintained).Encode())
	avps = append(avps, msg.OriginHost(l.Host).Encode())
	avps = append(avps, msg.OriginRealm(l.Realm).Encode())
	if v.DestinationHost != nil {
		avps = append(avps, v.DestinationHost.Encode())
	}
	avps = append(avps, v.DestinationRealm.Encode())

	avps = append(avps, v.SCAddress.Encode())
	if v.OFRFlags != nil {
		avps = append(avps, v.OFRFlags.Encode())
	}
	avps = append(avps, v.UserIdentifier.Encode())
	avps = append(avps, v.SMRPUI.Encode())
	// avps = append(avps, v.SMSMICorrelationID.Encode())
	// avps = append(avps, v.SMDeliveryOutcome.Encode())

	m.Encode(avps)
	return m
}

/*
 <OFA> ::= < Diameter Header: 8388645, PXY, 16777313 >
           < Session-Id >
           [ DRMP ]
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

/*
 <TFR> ::= < Diameter Header: 8388646, REQ, PXY, 16777313 >
           < Session-Id >
           [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           { Destination-Host }
           { Destination-Realm }
           { User-Name }
         * [ Supported-Features ]
           [ SMSMI-Correlation-ID ]
           { SC-Address }
           { SM-RP-UI }
           [ MME-Number-for-MT-SMS ]
           [ SGSN-Number ]
           [ TFR-Flags ]
           [ SM-Delivery-Timer ]
           [ SM-Delivery-Start-Time ]
           [ Maximum-Retransmission-Time ]
		   [ SMS-GMSC-Address ]
         * [ AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
 <TFA> ::= < Diameter Header: 8388646, PXY, 16777313 >
           < Session-Id >
           [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           [ Result-Code ]
           [ Experimental-Result ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
         * [ Supported-Features ]
           [ Absent-User-Diagnostic-SM ]
           [ SM-Delivery- Failure-Cause ]
           [ SM-RP-UI ]
           [ Requested-Retransmission-Time ]
           [ User-Identifier ]
         * [ AVP ]
         * [ Failed-AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/
