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

/*
MTForwardShortMessageRequest is TFR message.
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
*/
type MTForwardShortMessageRequest struct {
	msg.SessionID
	// DRMP
	msg.OriginHost
	msg.OriginRealm
	msg.DestinationHost
	msg.DestinationRealm
	msg.UserName
	// SupportedFeatures
	// SMSMI-Correlation-ID
	SCAddress
	SMRPUI
	// MME-Number-for-MT-SMS
	*SGSNNumber
	*TFRFlags
	*SMDeliveryTimer
	*SMDeliveryStartTime
	*MaximumRetransmissionTime
	*SMSGMSCAddress
	// Proxy-Info
	RouteRecord []msg.RouteRecord
}

// Encode return Message struct of this value
func (v *MTForwardShortMessageRequest) Encode() msg.Message {

	m := msg.Message{
		Ver:  msg.DiaVer,
		FlgR: true, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388646, AppID: 16777313}

	var avps []msg.Avp
	avps = append(avps, v.SessionID.Encode())
	avps = append(avps, msg.VendorSpecificApplicationID{
		VendorID: 10415,
		App:      msg.AuthApplicationID(16777313)}.Encode())
	avps = append(avps, msg.AuthSessionState(msg.StateNotMaintained).Encode())

	avps = append(avps, v.OriginHost.Encode())
	avps = append(avps, v.OriginRealm.Encode())
	avps = append(avps, v.DestinationHost.Encode())
	avps = append(avps, v.DestinationRealm.Encode())

	avps = append(avps, v.UserName.Encode())
	avps = append(avps, v.SCAddress.Encode())
	avps = append(avps, v.SMRPUI.Encode())

	if v.SGSNNumber != nil {
		avps = append(avps, v.SGSNNumber.Encode())
	}
	if v.TFRFlags != nil {
		avps = append(avps, v.TFRFlags.Encode())
	}
	if v.SMDeliveryTimer != nil {
		avps = append(avps, v.SMDeliveryTimer.Encode())
	}
	if v.SMDeliveryStartTime != nil {
		avps = append(avps, v.SMDeliveryStartTime.Encode())
	}
	if v.MaximumRetransmissionTime != nil {
		avps = append(avps, v.MaximumRetransmissionTime.Encode())
	}
	if v.SMSGMSCAddress != nil {
		avps = append(avps, v.SMSGMSCAddress.Encode())
	}

	for _, rr := range v.RouteRecord {
		avps = append(avps, rr.Encode())
	}

	m.Encode(avps)
	return m
}

// GetMTForwardShortMessageRequestt get Message value
func GetMTForwardShortMessageRequestt(o msg.Message) (
	s MTForwardShortMessageRequest, b bool) {
	if o.Code != 8388646 || o.AppID != 16777313 || !o.FlgR {
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
	if s.DestinationHost, b = msg.GetDestinationHost(avp); !b {
		return
	}
	if s.DestinationRealm, b = msg.GetDestinationRealm(avp); !b {
		return
	}

	if s.UserName, b = msg.GetUserName(avp); !b {
		return
	}
	if s.SCAddress, b = GetSCAddress(avp); !b {
		return
	}
	if s.SMRPUI, b = GetSMRPUI(avp); !b {
		return
	}

	if tmp, ok := GetSGSNNumber(avp); ok {
		s.SGSNNumber = &tmp
	}
	if tmp, ok := GetTFRFlags(avp); ok {
		s.TFRFlags = &tmp
	}
	if tmp, ok := GetSMDeliveryTimer(avp); ok {
		s.SMDeliveryTimer = &tmp
	}
	if tmp, ok := GetSMDeliveryStartTime(avp); ok {
		s.SMDeliveryStartTime = &tmp
	}
	if tmp, ok := GetMaximumRetransmissionTime(avp); ok {
		s.MaximumRetransmissionTime = &tmp
	}
	if tmp, ok := GetSMSGMSCAddress(avp); ok {
		s.SMSGMSCAddress = &tmp
	}
	s.RouteRecord = msg.GetRouteRecords(avp)

	b = true
	return
}

/*
MTForwardShortMessageAnswer is TFA message.
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
type MTForwardShortMessageAnswer struct {
	msg.SessionID
	// DRMP
	*msg.ResultCode
	*msg.ExperimentalResult
	msg.OriginHost
	msg.OriginRealm
	// SupportedFeatures
	*AbsentUserDiagnosticSM
	*SMDeliveryFailureCause
	SMRPUI
	*RequestedRetransmissionTime
	*ts29336.UserIdentifier
	// Failed-AVP
	// Proxy-Info
	RouteRecord []msg.RouteRecord
}

// Encode return Message struct of this value
func (v *MTForwardShortMessageAnswer) Encode() msg.Message {

	m := msg.Message{
		Ver:  msg.DiaVer,
		FlgR: false, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388646, AppID: 16777313}

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

	if v.AbsentUserDiagnosticSM != nil {
		avps = append(avps, v.AbsentUserDiagnosticSM.Encode())
	}
	if v.SMDeliveryFailureCause != nil {
		avps = append(avps, v.SMDeliveryFailureCause.Encode())
	}
	if v.SMRPUI != nil {
		avps = append(avps, v.SMRPUI.Encode())
	}
	if v.RequestedRetransmissionTime != nil {
		avps = append(avps, v.RequestedRetransmissionTime.Encode())
	}
	if v.UserIdentifier != nil {
		avps = append(avps, v.UserIdentifier.Encode())
	}

	for _, rr := range v.RouteRecord {
		avps = append(avps, rr.Encode())
	}

	m.Encode(avps)
	return m
}

// GetMTForwardShortMessageAnswer get Message value
func GetMTForwardShortMessageAnswer(o msg.Message) (
	s MTForwardShortMessageAnswer, b bool) {
	if o.Code != 8388646 || o.AppID != 16777313 || o.FlgR {
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

	if tmp, ok := GetAbsentUserDiagnosticSM(avp); ok {
		s.AbsentUserDiagnosticSM = &tmp
	}
	if tmp, ok := GetSMDeliveryFailureCause(avp); ok {
		s.SMDeliveryFailureCause = &tmp
	}
	s.SMRPUI, _ = GetSMRPUI(avp)
	if tmp, ok := GetRequestedRetransmissionTime(avp); ok {
		s.RequestedRetransmissionTime = &tmp
	}
	if tmp, ok := ts29336.GetUserIdentifier(avp); ok {
		s.UserIdentifier = &tmp
	}

	s.RouteRecord = msg.GetRouteRecords(avp)

	b = true
	return
}
