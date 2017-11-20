package ts29338

import (
	dia "github.com/fkgi/diameter"
	"github.com/fkgi/teldata"
)

func setSCAddress(v teldata.E164) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3300, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	a.Encode(v.String())
	return
}

func getSCAddress(a dia.RawAVP) (v teldata.E164, e error) {
	s := new(string)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e != nil {
		v, e = teldata.ParseE164(*s)
	}
	return
}

/*
// SMRPUI AVP contain a short message transfer protocol data unit (TPDU).
// Maximum length is 200 octets.
type SMRPUI []byte

// ToRaw return AVP struct of this value
func (v SMRPUI) ToRaw() diameter.RawAVP {
	a := diameter.RawAVP{Code: 3301, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode([]byte(v))
	return a
}

// GetSMRPUI get AVP value
func GetSMRPUI(o diameter.GroupedAVP) (SMRPUI, bool) {
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

// ToRaw return AVP struct of this value
func (v TFRFlags) ToRaw() diameter.RawAVP {
	a := diameter.RawAVP{Code: 3302, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	i := uint32(0)

	if v.MMS {
		i = i | 0x00000001
	}
	a.Encode(i)
	return a
}

// GetTFRFlags get AVP value
func GetTFRFlags(o diameter.GroupedAVP) (TFRFlags, bool) {
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
	Cause diameter.Enumerated
	Diag  []byte
}

// ToRaw return AVP struct of this value
func (v SMDeliveryFailureCause) ToRaw() diameter.RawAVP {
	a := diameter.RawAVP{Code: 3303, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	var t []diameter.RawAVP

	// SM-Enumerated-Delivery-Failure-Cause
	{
		a := diameter.RawAVP{Code: 3304, VenID: 10415,
			FlgV: true, FlgM: true, FlgP: false}
		a.Encode(v.Cause)
		t = append(t, a)
	}

	// SM-Diagnostic-Info
	if v.Diag != nil {
		a := diameter.RawAVP{Code: 3305, VenID: 10415,
			FlgV: true, FlgM: true, FlgP: false}
		a.Encode(v.Diag)
		t = append(t, a)
	}

	a.Encode(diameter.GroupedAVP(t))
	return a
}

// GetSMDeliveryFailureCause get AVP value
func GetSMDeliveryFailureCause(o diameter.GroupedAVP) (SMDeliveryFailureCause, bool) {
	s := SMDeliveryFailureCause{}
	if a, ok := o.Get(3303, 10415); ok {
		o = diameter.GroupedAVP{}
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
	MemoryCapacityExceeded diameter.Enumerated = 0
	// EquipmentProtocolError is Enumerated value 1
	EquipmentProtocolError diameter.Enumerated = 1
	// EquipmentNotSMEquipped is Enumerated value 2
	EquipmentNotSMEquipped diameter.Enumerated = 2
	// UnknownServiceCenter is Enumerated value 3
	UnknownServiceCenter diameter.Enumerated = 3
	// SCCongestion is Enumerated value 4
	SCCongestion diameter.Enumerated = 4
	// InvalidSMEAddress is Enumerated value 5
	InvalidSMEAddress diameter.Enumerated = 5
	// UserNotSCUser is Enumerated value 6
	UserNotSCUser diameter.Enumerated = 6
)

// SMDeliveryTimer AVP contain the value in seconds of the timer for SM Delivery.
type SMDeliveryTimer uint32

// ToRaw return AVP struct of this value
func (v SMDeliveryTimer) ToRaw() diameter.RawAVP {
	a := diameter.RawAVP{Code: 3306, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetSMDeliveryTimer get AVP value
func GetSMDeliveryTimer(o diameter.GroupedAVP) (SMDeliveryTimer, bool) {
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

// ToRaw return AVP struct of this value
func (v SMDeliveryStartTime) ToRaw() diameter.RawAVP {
	a := diameter.RawAVP{Code: 3307, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(time.Time(v))
	return a
}

// GetSMDeliveryStartTime get AVP value
func GetSMDeliveryStartTime(o diameter.GroupedAVP) (SMDeliveryStartTime, bool) {
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

// ToRaw return AVP struct of this value
func (v OFRFlags) ToRaw() diameter.RawAVP {
	a := diameter.RawAVP{Code: 3328, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	i := uint32(0)

	if v.S6as6d {
		i = i | 0x00000001
	}

	a.Encode(i)
	return a
}

// GetOFRFlags get AVP value
func GetOFRFlags(o diameter.GroupedAVP) (OFRFlags, bool) {
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

// ToRaw return AVP struct of this value
func (v SMSMICorrelationID) ToRaw() diameter.RawAVP {
	a := diameter.RawAVP{Code: 3324, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	var t []diameter.RawAVP

	// HSS-ID
	if len(v.HSSID) != 0 {
		a := diameter.RawAVP{Code: 3325, VenID: 10415,
			FlgV: true, FlgM: false, FlgP: false}
		a.Encode(v.HSSID)
		t = append(t, a)
	}
	// Originating-SIP-URI
	if len(v.OrigSIPURI) != 0 {
		a := diameter.RawAVP{Code: 3326, VenID: 10415,
			FlgV: true, FlgM: false, FlgP: false}
		a.Encode(v.OrigSIPURI)
		t = append(t, a)
	}
	// Destination-SIP-URI
	if len(v.DestSIPURI) != 0 {
		a := diameter.RawAVP{Code: 3327, VenID: 10415,
			FlgV: true, FlgM: false, FlgP: false}
		a.Encode(v.DestSIPURI)
		t = append(t, a)
	}

	a.Encode(diameter.GroupedAVP(t))
	return a
}

// GetSMSMICorrelationID get AVP value
func GetSMSMICorrelationID(o diameter.GroupedAVP) (SMSMICorrelationID, bool) {
	s := SMSMICorrelationID{}
	if a, ok := o.Get(3324, 10415); ok {
		o = diameter.GroupedAVP{}
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

// ToRaw return AVP struct of this value
func (v MaximumRetransmissionTime) ToRaw() diameter.RawAVP {
	a := diameter.RawAVP{Code: 3330, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode(time.Time(v))
	return a
}

// GetMaximumRetransmissionTime get AVP value
func GetMaximumRetransmissionTime(o diameter.GroupedAVP) (MaximumRetransmissionTime, bool) {
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

// ToRaw return AVP struct of this value
func (v RequestedRetransmissionTime) ToRaw() diameter.RawAVP {
	a := diameter.RawAVP{Code: 3331, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode(time.Time(v))
	return a
}

// GetRequestedRetransmissionTime get AVP value
func GetRequestedRetransmissionTime(o diameter.GroupedAVP) (RequestedRetransmissionTime, bool) {
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

// ToRaw return AVP struct of this value
func (v SMSGMSCAddress) ToRaw() diameter.RawAVP {
	a := diameter.RawAVP{Code: 3332, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode(string(v))
	return a
}

// GetSMSGMSCAddress get AVP value
func GetSMSGMSCAddress(o diameter.GroupedAVP) (SMSGMSCAddress, bool) {
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
type UserName diameter.UserName

// ToRaw return AVP struct of this value
func (v UserName) ToRaw() diameter.RawAVP {
	return diameter.UserName(v).Encode()
}

// GetUserName get AVP value
func GetUserName(o diameter.GroupedAVP) (UserName, bool) {
	a, ok := diameter.GetUserName(o)
	return UserName(a), ok
}

// UserIdentifier AVP from ts29.336
type UserIdentifier ts29336.UserIdentifier

// ToRaw return AVP struct of this value
func (v UserIdentifier) ToRaw() diameter.RawAVP {
	return ts29336.UserIdentifier(v).Encode()
}

// GetUserIdentifier get AVP value
func GetUserIdentifier(o diameter.GroupedAVP) (UserIdentifier, bool) {
	a, ok := ts29336.GetUserIdentifier(o)
	return UserIdentifier(a), ok
}

// MMENumberForMTSMS AVP from ts29.272
type MMENumberForMTSMS ts29272.MMENumberForMTSMS

// ToRaw return AVP struct of this value
func (v MMENumberForMTSMS) ToRaw() diameter.RawAVP {
	return ts29272.MMENumberForMTSMS(v).Encode()
}

// GetMMENumberForMTSMS get AVP value
func GetMMENumberForMTSMS(o diameter.GroupedAVP) (MMENumberForMTSMS, bool) {
	a, ok := ts29272.GetMMENumberForMTSMS(o)
	return MMENumberForMTSMS(a), ok
}
*/

/*
// SGSNNumber AVP from ts29.272
type SGSNNumber ts29272.SGSNNumber

// ToRaw return AVP struct of this value
func (v SGSNNumber) ToRaw() diameter.RawAVP {
	a := ts29272.SGSNNumber(v).Encode()
	a.FlgM = false
	return a
}

// GetSGSNNumber get AVP value
func GetSGSNNumber(o diameter.GroupedAVP) (SGSNNumber, bool) {
	a, ok := ts29272.GetSGSNNumber(o)
	return SGSNNumber(a), ok
}

/*
// SupportedFeatures AVP from ts29.229
type SupportedFeatures ts29229.SupportedFeatures

// ToRaw return AVP struct of this value
func (v SupportedFeatures) ToRaw() diameter.RawAVP {
	return ts29229.SupportedFeatures(v).Encode()
}

// GetSupportedFeatures get AVP value
func GetSupportedFeatures(o diameter.GroupedAVP) (SupportedFeatures, bool) {
	a, ok := ts29229.GetSupportedFeatures(o)
	return SupportedFeatures(a), ok
}

// FeatureListID AVP from ts29.229
type FeatureListID ts29229.FeatureListID

// ToRaw return AVP struct of this value
func (v FeatureListID) ToRaw() diameter.RawAVP {
	return ts29229.FeatureListID(v).Encode()
}

// GetFeatureListID get AVP value
func GetFeatureListID(o diameter.GroupedAVP) (FeatureListID, bool) {
	a, ok := ts29229.GetFeatureListID(o)
	return FeatureListID(a), ok
}

// FeatureList AVP from ts29.229
type FeatureList ts29229.FeatureList

// ToRaw return AVP struct of this value
func (v FeatureList) ToRaw() diameter.RawAVP {
	return ts29229.FeatureList(v).Encode()
}

// GetFeatureList get AVP value
func GetFeatureList(o diameter.GroupedAVP) (FeatureList, bool) {
	a, ok := ts29229.GetFeatureList(o)
	return FeatureList(a), ok
}
*/
/*
// ExternalIdentifier AVP from ts29.336
type ExternalIdentifier ts29336.ExternalIdentifier

// ToRaw return AVP struct of this value
func (v ExternalIdentifier) ToRaw() diameter.RawAVP {
	a := ts29336.ExternalIdentifier(v).Encode()
	a.FlgM = false
	return a
}

// GetExternalIdentifier get AVP value
func GetExternalIdentifier(o diameter.GroupedAVP) (ExternalIdentifier, bool) {
	a, ok := ts29336.GetExternalIdentifier(o)
	return ExternalIdentifier(a), ok
}
*/
