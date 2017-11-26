package ts29338

import (
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
	msg.SessionID
	// DRMP
	msg.OriginHost
	msg.OriginRealm
	*msg.DestinationHost
	msg.DestinationRealm
	SCAddress
	*OFRFlags
	// SupportedFeatures
	ts29336.UserIdentifier
	SMRPUI
	// SMSMICorrelationID
	// SMDeliveryOutcome
	// Proxy-Info
	RouteRecord []msg.RouteRecord
}

// Encode return Message struct of this value
func (v *MOForwardShortMessageRequest) Encode() msg.Message {

	m := msg.Message{
		Ver:  msg.DiaVer,
		FlgR: true, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388645, AppID: 16777313}

	var avps []msg.Avp
	avps = append(avps, v.SessionID.Encode())
	avps = append(avps, msg.VendorSpecificApplicationID{
		VendorID: 10415,
		App:      msg.AuthApplicationID(16777313)}.Encode())
	avps = append(avps, msg.AuthSessionState(msg.StateNotMaintained).Encode())

	avps = append(avps, v.OriginHost.Encode())
	avps = append(avps, v.OriginRealm.Encode())

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

	for _, rr := range v.RouteRecord {
		avps = append(avps, rr.Encode())
	}

	m.Encode(avps)
	return m
}

// GetMOForwardShortMessageRequest get Message value
func GetMOForwardShortMessageRequest(o msg.Message) (
	s MOForwardShortMessageRequest, b bool) {
	if o.Code != 8388645 || o.AppID != 16777313 || !o.FlgR {
		return
	}
	avp, e := o.Decode()
	if e != nil {
		return
	}

	if s.SessionID, b = msg.GetSessionID(avp); !b {
		return
	}
	if s.OriginHost, b = msg.GetOriginHost(avp); !b {
		return
	}
	if s.OriginRealm, b = msg.GetOriginRealm(avp); !b {
		return
	}
	if tmp, ok := msg.GetDestinationHost(avp); ok {
		s.DestinationHost = &tmp
	}
	if s.DestinationRealm, b = msg.GetDestinationRealm(avp); !b {
		return
	}
	if s.SCAddress, b = GetSCAddress(avp); !b {
		return
	}
	if tmp, ok := GetOFRFlags(avp); ok {
		s.OFRFlags = &tmp
	}
	if s.UserIdentifier, b = ts29336.GetUserIdentifier(avp); !b {
		return
	}
	if s.SMRPUI, b = GetSMRPUI(avp); !b {
		return
	}
	// if s.SMSMICorrelationID, b = GetSMSMICorrelationID(avp); !b {
	// 	return
	// }
	// if s.SMDeliveryOutcome, b = GetSMDeliveryOutcome(avp); !b {
	// 	return
	// }
	s.RouteRecord = msg.GetRouteRecords(avp)

	b = true
	return
}

/*
MOForwardShortMessageAnswer is OFA message.
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
type MOForwardShortMessageAnswer struct {
	msg.SessionID
	// DRMP
	*msg.ResultCode
	*msg.ExperimentalResult
	msg.OriginHost
	msg.OriginRealm
	// SupportedFeatures
	*SMDeliveryFailureCause
	SMRPUI
	// Failed-AVP
	// Proxy-Info
	RouteRecord []msg.RouteRecord
}

// Encode return Message struct of this value
func (v *MOForwardShortMessageAnswer) Encode() msg.Message {

	m := msg.Message{
		Ver:  msg.DiaVer,
		FlgR: false, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388645, AppID: 16777313}

	var avps []msg.Avp
	avps = append(avps, v.SessionID.Encode())
	avps = append(avps, msg.VendorSpecificApplicationID{
		VendorID: 10415,
		App:      msg.AuthApplicationID(16777313)}.Encode())

	if v.ResultCode != nil {
		avps = append(avps, v.ResultCode.Encode())
	}
	if v.ExperimentalResult != nil {
		avps = append(avps, v.ExperimentalResult.Encode())
	}

	avps = append(avps, msg.AuthSessionState(msg.StateNotMaintained).Encode())
	avps = append(avps, v.OriginHost.Encode())
	avps = append(avps, v.OriginRealm.Encode())

	if v.SMDeliveryFailureCause != nil {
		avps = append(avps, v.SMDeliveryFailureCause.Encode())
	}

	if v.SMRPUI != nil {
		avps = append(avps, v.SMRPUI.Encode())
	}

	for _, rr := range v.RouteRecord {
		avps = append(avps, rr.Encode())
	}

	m.Encode(avps)
	return m
}

// GetMOForwardShortMessageAnswer get Message value
func GetMOForwardShortMessageAnswer(o msg.Message) (
	s MOForwardShortMessageAnswer, b bool) {
	if o.Code != 8388645 || o.AppID != 16777313 || o.FlgR {
		return
	}
	avp, e := o.Decode()
	if e != nil {
		return
	}

	if s.SessionID, b = msg.GetSessionID(avp); !b {
		return
	}
	if tmp, ok := msg.GetResultCode(avp); ok {
		s.ResultCode = &tmp
	}
	if tmp, ok := msg.GetExperimentalResult(avp); ok {
		s.ExperimentalResult = &tmp
	}
	if s.OriginHost, b = msg.GetOriginHost(avp); !b {
		return
	}
	if s.OriginRealm, b = msg.GetOriginRealm(avp); !b {
		return
	}

	if tmp, ok := GetSMDeliveryFailureCause(avp); ok {
		s.SMDeliveryFailureCause = &tmp
	}
	s.SMRPUI, _ = GetSMRPUI(avp)
	s.RouteRecord = msg.GetRouteRecords(avp)

	b = true
	return
}
