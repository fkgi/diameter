package ts29338

import (
	"time"

	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/ts29272"
	"github.com/fkgi/diameter/ts29336"
)

const (
	// DiameterErrorAbsentUser is Result-Code 5550
	DiameterErrorAbsentUser msg.ResultCode = 5550
	// DiameterErrorUserBusyForMtSms is Result-Code 5551
	DiameterErrorUserBusyForMtSms msg.ResultCode = 5551
	// DiameterErrorFacilityNotSupported is Result-Code 5552
	DiameterErrorFacilityNotSupported msg.ResultCode = 5552
	// DiameterErrorIlleagalUser is Result-Code 5553
	DiameterErrorIlleagalUser msg.ResultCode = 5553
	// DiameterErrorIlleagalEquipment is Result-Code 5554
	DiameterErrorIlleagalEquipment msg.ResultCode = 5554
	// DiameterErrorSmDeliveryFailure is Result-Code 5555
	DiameterErrorSmDeliveryFailure msg.ResultCode = 5555
	// DiameterErrorServiceNotSubscribed is Result-Code 5556
	DiameterErrorServiceNotSubscribed msg.ResultCode = 5556
	// DiameterErrorServiceBarred is Result-Code 5557
	DiameterErrorServiceBarred msg.ResultCode = 5557
	// DiameterErrorMwdListFull is Result-Code 5558
	DiameterErrorMwdListFull msg.ResultCode = 5558
)

/*
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

// SCAddress AVP contain the E164 number of the SMS-SC or MTC-IWF.
type SCAddress string

// Encode return AVP struct of this value
func (v SCAddress) Encode() msg.Avp {
	a := msg.Avp{Code: 3300, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(string(v))
	return a
}

// GetSCAddress get AVP value
func GetSCAddress(o msg.GroupedAVP) (SCAddress, bool) {
	s := new(string)
	if a, ok := o.Get(3300, 10415); ok {
		a.Decode(s)
	} else {
		return "", false
	}
	return SCAddress(*s), true
}

// SMRPUI AVP contain a short message transfer protocol data unit (TPDU).
// Maximum length is 200 octets.
type SMRPUI []byte

// Encode return AVP struct of this value
func (v SMRPUI) Encode() msg.Avp {
	a := msg.Avp{Code: 3301, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode([]byte(v))
	return a
}

// GetSMRPUI get AVP value
func GetSMRPUI(o msg.GroupedAVP) (SMRPUI, bool) {
	s := new([]byte)
	if a, ok := o.Get(3301, 10415); ok {
		a.Decode(s)
	} else {
		return nil, false
	}
	return SMRPUI(*s), true
}

// TFRFlags AVP is bit mask.
// When moreMsgToSend set, the service centre has more short messages to send.
type TFRFlags struct {
	MMS bool // More message to send
}

// Encode return AVP struct of this value
func (v TFRFlags) Encode() msg.Avp {
	a := msg.Avp{Code: 3302, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	i := uint32(0)

	if v.MMS {
		i = i | 0x00000001
	}
	a.Encode(i)
	return a
}

// GetTFRFlags get AVP value
func GetTFRFlags(o msg.GroupedAVP) (TFRFlags, bool) {
	s := new(uint32)
	if a, ok := o.Get(3302, 10415); ok {
		a.Decode(s)
	} else {
		return TFRFlags{}, false
	}
	return TFRFlags{
		MMS: (*s)&0x00000001 == 0x00000001}, true
}

// SMDeliveryFailureCause AVP contain cause of the failure of
// a SM delivery with an complementary information.
// If diag is nil, complementary information is empty.
type SMDeliveryFailureCause struct {
	Cause msg.Enumerated
	Diag  []byte
}

// Encode return AVP struct of this value
func (v SMDeliveryFailureCause) Encode() msg.Avp {
	a := msg.Avp{Code: 3303, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	var t []msg.Avp

	// SM-Enumerated-Delivery-Failure-Cause
	{
		a := msg.Avp{Code: 3304, VenID: 10415,
			FlgV: true, FlgM: true, FlgP: false}
		a.Encode(v.Cause)
		t = append(t, a)
	}

	// SM-Diagnostic-Info
	if v.Diag != nil {
		a := msg.Avp{Code: 3305, VenID: 10415,
			FlgV: true, FlgM: true, FlgP: false}
		a.Encode(v.Diag)
		t = append(t, a)
	}

	a.Encode(msg.GroupedAVP(t))
	return a
}

// GetSMDeliveryFailureCause get AVP value
func GetSMDeliveryFailureCause(o msg.GroupedAVP) (SMDeliveryFailureCause, bool) {
	s := SMDeliveryFailureCause{}
	if a, ok := o.Get(3303, 10415); ok {
		o = msg.GroupedAVP{}
		a.Decode(&o)
	} else {
		return s, false
	}
	if t, ok := o.Get(3304, 10415); ok {
		t.Decode(&s.Cause)
	}
	if t, ok := o.Get(3305, 10415); ok {
		t.Decode(&s.Diag)
	}
	return s, true
}

const (
	// MemoryCapacityExceeded is Enumerated value 0
	MemoryCapacityExceeded msg.Enumerated = 0
	// EquipmentProtocolError is Enumerated value 1
	EquipmentProtocolError msg.Enumerated = 1
	// EquipmentNotSMEquipped is Enumerated value 2
	EquipmentNotSMEquipped msg.Enumerated = 2
	// UnknownServiceCenter is Enumerated value 3
	UnknownServiceCenter msg.Enumerated = 3
	// SCCongestion is Enumerated value 4
	SCCongestion msg.Enumerated = 4
	// InvalidSMEAddress is Enumerated value 5
	InvalidSMEAddress msg.Enumerated = 5
	// UserNotSCUser is Enumerated value 6
	UserNotSCUser msg.Enumerated = 6
)

// SMDeliveryTimer AVP contain the value in seconds of the timer for SM Delivery.
type SMDeliveryTimer uint32

// Encode return AVP struct of this value
func (v SMDeliveryTimer) Encode() msg.Avp {
	a := msg.Avp{Code: 3306, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetSMDeliveryTimer get AVP value
func GetSMDeliveryTimer(o msg.GroupedAVP) (SMDeliveryTimer, bool) {
	s := new(uint32)
	if a, ok := o.Get(3306, 10415); ok {
		a.Decode(s)
	} else {
		return 0, false
	}
	return SMDeliveryTimer(*s), true
}

// SMDeliveryStartTime AVP contain the timestamp (in UTC) at which
// the SM Delivery Supervision Timer was started.
type SMDeliveryStartTime time.Time

// Encode return AVP struct of this value
func (v SMDeliveryStartTime) Encode() msg.Avp {
	a := msg.Avp{Code: 3307, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(time.Time(v))
	return a
}

// GetSMDeliveryStartTime get AVP value
func GetSMDeliveryStartTime(o msg.GroupedAVP) (SMDeliveryStartTime, bool) {
	s := new(time.Time)
	if a, ok := o.Get(3307, 10415); ok {
		a.Decode(s)
	} else {
		return SMDeliveryStartTime{}, false
	}
	return SMDeliveryStartTime(*s), true
}

// OFRFlags AVP is bit mask.
// When s6as6d set, the OFR message is sent on the Gdd interface (source node is an SGSN).
// When cleared, sent on the SGd interface (source node is an MME).
type OFRFlags struct {
	S6as6d bool
}

// Encode return AVP struct of this value
func (v OFRFlags) Encode() msg.Avp {
	a := msg.Avp{Code: 3328, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	i := uint32(0)

	if v.S6as6d {
		i = i | 0x00000001
	}

	a.Encode(i)
	return a
}

// GetOFRFlags get AVP value
func GetOFRFlags(o msg.GroupedAVP) (OFRFlags, bool) {
	s := new(uint32)
	if a, ok := o.Get(3328, 10415); ok {
		a.Decode(s)
	} else {
		return OFRFlags{}, false
	}
	return OFRFlags{
		S6as6d: (*s)&0x00000001 == 0x00000001}, true
}

// SMSMICorrelationID AVP ontain information identities used in the context
// of MSISDN-less SMS delivery in IMS
type SMSMICorrelationID struct {
	HSSID      string
	OrigSIPURI string
	DestSIPURI string
}

// Encode return AVP struct of this value
func (v SMSMICorrelationID) Encode() msg.Avp {
	a := msg.Avp{Code: 3324, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	var t []msg.Avp

	// HSS-ID
	if len(v.HSSID) != 0 {
		a := msg.Avp{Code: 3325, VenID: 10415,
			FlgV: true, FlgM: false, FlgP: false}
		a.Encode(v.HSSID)
		t = append(t, a)
	}
	// Originating-SIP-URI
	if len(v.OrigSIPURI) != 0 {
		a := msg.Avp{Code: 3326, VenID: 10415,
			FlgV: true, FlgM: false, FlgP: false}
		a.Encode(v.OrigSIPURI)
		t = append(t, a)
	}
	// Destination-SIP-URI
	if len(v.DestSIPURI) != 0 {
		a := msg.Avp{Code: 3327, VenID: 10415,
			FlgV: true, FlgM: false, FlgP: false}
		a.Encode(v.DestSIPURI)
		t = append(t, a)
	}

	a.Encode(msg.GroupedAVP(t))
	return a
}

// GetSMSMICorrelationID get AVP value
func GetSMSMICorrelationID(o msg.GroupedAVP) (SMSMICorrelationID, bool) {
	s := SMSMICorrelationID{}
	if a, ok := o.Get(3324, 10415); ok {
		o = msg.GroupedAVP{}
		a.Decode(&o)
	} else {
		return s, false
	}
	if t, ok := o.Get(3325, 10415); ok {
		t.Decode(&s.HSSID)
	}
	if t, ok := o.Get(3326, 10415); ok {
		t.Decode(&s.OrigSIPURI)
	}
	if t, ok := o.Get(3327, 10415); ok {
		t.Decode(&s.DestSIPURI)
	}
	return s, true
}

// MaximumRetransmissionTime AVP contain the maximum retransmission time (in UTC) until which
// the SMS-GMSC is capable to retransmit the MT Short Message.
type MaximumRetransmissionTime time.Time

// Encode return AVP struct of this value
func (v MaximumRetransmissionTime) Encode() msg.Avp {
	a := msg.Avp{Code: 3330, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode(time.Time(v))
	return a
}

// GetMaximumRetransmissionTime get AVP value
func GetMaximumRetransmissionTime(o msg.GroupedAVP) (MaximumRetransmissionTime, bool) {
	s := new(time.Time)
	if a, ok := o.Get(3330, 10415); ok {
		a.Decode(s)
	} else {
		return MaximumRetransmissionTime{}, false
	}
	return MaximumRetransmissionTime(*s), true
}

// RequestedRetransmissionTime AVP contain the timestamp (in UTC) at which
// the SMS-GMSC is requested to retransmit the MT Short Message.
type RequestedRetransmissionTime time.Time

// Encode return AVP struct of this value
func (v RequestedRetransmissionTime) Encode() msg.Avp {
	a := msg.Avp{Code: 3331, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode(time.Time(v))
	return a
}

// GetRequestedRetransmissionTime get AVP value
func GetRequestedRetransmissionTime(o msg.GroupedAVP) (RequestedRetransmissionTime, bool) {
	s := new(time.Time)
	if a, ok := o.Get(3331, 10415); ok {
		a.Decode(s)
	} else {
		return RequestedRetransmissionTime{}, false
	}
	return RequestedRetransmissionTime(*s), true
}

// SMSGMSCAddress AVP contain the E.164 number of the SMS-GMSC or SMS Router.
type SMSGMSCAddress string

// Encode return AVP struct of this value
func (v SMSGMSCAddress) Encode() msg.Avp {
	a := msg.Avp{Code: 3332, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode(string(v))
	return a
}

// GetSMSGMSCAddress get AVP value
func GetSMSGMSCAddress(o msg.GroupedAVP) (SMSGMSCAddress, bool) {
	s := new(string)
	if a, ok := o.Get(3332, 10415); ok {
		a.Decode(s)
	} else {
		return "", false
	}
	return SMSGMSCAddress(*s), true
}

/*
// UserName AVP from RFC6733
type UserName msg.UserName

// Encode return AVP struct of this value
func (v UserName) Encode() msg.Avp {
	return msg.UserName(v).Encode()
}

// GetUserName get AVP value
func GetUserName(o msg.GroupedAVP) (UserName, bool) {
	a, ok := msg.GetUserName(o)
	return UserName(a), ok
}

// UserIdentifier AVP from ts29.336
type UserIdentifier ts29336.UserIdentifier

// Encode return AVP struct of this value
func (v UserIdentifier) Encode() msg.Avp {
	return ts29336.UserIdentifier(v).Encode()
}

// GetUserIdentifier get AVP value
func GetUserIdentifier(o msg.GroupedAVP) (UserIdentifier, bool) {
	a, ok := ts29336.GetUserIdentifier(o)
	return UserIdentifier(a), ok
}

// MMENumberForMTSMS AVP from ts29.272
type MMENumberForMTSMS ts29272.MMENumberForMTSMS

// Encode return AVP struct of this value
func (v MMENumberForMTSMS) Encode() msg.Avp {
	return ts29272.MMENumberForMTSMS(v).Encode()
}

// GetMMENumberForMTSMS get AVP value
func GetMMENumberForMTSMS(o msg.GroupedAVP) (MMENumberForMTSMS, bool) {
	a, ok := ts29272.GetMMENumberForMTSMS(o)
	return MMENumberForMTSMS(a), ok
}
*/

// SGSNNumber AVP from ts29.272
type SGSNNumber ts29272.SGSNNumber

// Encode return AVP struct of this value
func (v SGSNNumber) Encode() msg.Avp {
	a := ts29272.SGSNNumber(v).Encode()
	a.FlgM = false
	return a
}

// GetSGSNNumber get AVP value
func GetSGSNNumber(o msg.GroupedAVP) (SGSNNumber, bool) {
	a, ok := ts29272.GetSGSNNumber(o)
	return SGSNNumber(a), ok
}

/*
// SupportedFeatures AVP from ts29.229
type SupportedFeatures ts29229.SupportedFeatures

// Encode return AVP struct of this value
func (v SupportedFeatures) Encode() msg.Avp {
	return ts29229.SupportedFeatures(v).Encode()
}

// GetSupportedFeatures get AVP value
func GetSupportedFeatures(o msg.GroupedAVP) (SupportedFeatures, bool) {
	a, ok := ts29229.GetSupportedFeatures(o)
	return SupportedFeatures(a), ok
}

// FeatureListID AVP from ts29.229
type FeatureListID ts29229.FeatureListID

// Encode return AVP struct of this value
func (v FeatureListID) Encode() msg.Avp {
	return ts29229.FeatureListID(v).Encode()
}

// GetFeatureListID get AVP value
func GetFeatureListID(o msg.GroupedAVP) (FeatureListID, bool) {
	a, ok := ts29229.GetFeatureListID(o)
	return FeatureListID(a), ok
}

// FeatureList AVP from ts29.229
type FeatureList ts29229.FeatureList

// Encode return AVP struct of this value
func (v FeatureList) Encode() msg.Avp {
	return ts29229.FeatureList(v).Encode()
}

// GetFeatureList get AVP value
func GetFeatureList(o msg.GroupedAVP) (FeatureList, bool) {
	a, ok := ts29229.GetFeatureList(o)
	return FeatureList(a), ok
}
*/

// ExternalIdentifier AVP from ts29.336
type ExternalIdentifier ts29336.ExternalIdentifier

// Encode return AVP struct of this value
func (v ExternalIdentifier) Encode() msg.Avp {
	a := ts29336.ExternalIdentifier(v).Encode()
	a.FlgM = false
	return a
}

// GetExternalIdentifier get AVP value
func GetExternalIdentifier(o msg.GroupedAVP) (ExternalIdentifier, bool) {
	a, ok := ts29336.GetExternalIdentifier(o)
	return ExternalIdentifier(a), ok
}
