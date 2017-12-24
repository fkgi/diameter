package ts29338

import (
	"time"

	dia "github.com/fkgi/diameter"
	"github.com/fkgi/sms"
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
	} else if e = a.Decode(s); e == nil {
		v, e = teldata.ParseE164(*s)
	}
	return
}

// SMRPUI AVP contain a short message transfer protocol data unit (TPDU).
func setSMRPUI(v sms.TPDU) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3301, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	a.Encode(v.Encode())
	return
}

func getSMRPUIasDeliver(a dia.RawAVP) (v sms.Deliver, e error) {
	s := new([]byte)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil {
		e = v.Decode(*s)
	}
	return
}

func getSMRPUIasDeliverReport(a dia.RawAVP) (v sms.DeliverReport, e error) {
	s := new([]byte)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil {
		e = v.Decode(*s)
	}
	return
}

// MMENumberForMTSMS AVP from ts29.272
func setMMENumberForMTSMS(v teldata.E164) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 1645, VenID: 10415, FlgV: true, FlgM: false, FlgP: false}
	a.Encode(v.Bytes())
	return
}

func getMMENumberForMTSMS(a dia.RawAVP) (v teldata.E164, e error) {
	s := new([]byte)
	if !a.FlgV || a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil {
		v, e = teldata.B2E164(*s)
	}
	return
}

// SGSNNumber AVP from ts29.272
func setSGSNNumber(v teldata.E164) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 1489, VenID: 10415, FlgV: true, FlgM: false, FlgP: false}
	a.Encode(v.String())
	return
}

func getSGSNNumber(a dia.RawAVP) (v teldata.E164, e error) {
	s := new(string)
	if !a.FlgV || a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil {
		v, e = teldata.ParseE164(*s)
	}
	return
}

// TFRFlags AVP is bit mask.
// When moreMsgToSend set, the service centre has more short messages to send.
func setTFRFlags(m bool) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3302, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	i := uint32(0)
	if m {
		i = i | 0x00000001
	}
	a.Encode(i)
	return
}

func getTFRFlags(a dia.RawAVP) (m bool, e error) {
	v := new(uint32)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(v); e == nil {
		m = (*v)&0x00000001 == 0x00000001
	}
	return
}

// DeliveryFailureCause AVP contain cause of the failure of a SM delivery with an complementary information.
type DeliveryFailureCause int

const (
	// CauseMemoryCapacityExceeded is MEMORY_CAPACITY_EXCEEDED
	CauseMemoryCapacityExceeded DeliveryFailureCause = iota
	// CauseEquipmentProtocolError is EQUIPMENT_PROTOCOL_ERROR
	CauseEquipmentProtocolError
	// CauseEquipmentNotSMEquipped is EQUIPMENT_NOT_SM-EQUIPPED
	CauseEquipmentNotSMEquipped
	// CauseUnknownServiceCenter is UNKNOWN_SERVICE_CENTRE
	CauseUnknownServiceCenter
	// CauseSCCongestion is SC-CONGESTION
	CauseSCCongestion
	// CauseInvalidSMEAddress is INVALID_SME-ADDRESS
	CauseInvalidSMEAddress
	// CauseUserNotSCUser is USER_NOT_SC-USER
	CauseUserNotSCUser
)

func setSMDeliveryFailureCause(c DeliveryFailureCause, d sms.DeliverReport) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3303, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	v := make([]dia.RawAVP, 1, 2)

	// SM-Enumerated-Delivery-Failure-Cause
	v[0] = dia.RawAVP{Code: 3304, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	switch c {
	case CauseMemoryCapacityExceeded:
		v[0].Encode(dia.Enumerated(0))
	case CauseEquipmentProtocolError:
		v[0].Encode(dia.Enumerated(1))
	case CauseEquipmentNotSMEquipped:
		v[0].Encode(dia.Enumerated(2))
	case CauseUnknownServiceCenter:
		v[0].Encode(dia.Enumerated(3))
	case CauseSCCongestion:
		v[0].Encode(dia.Enumerated(4))
	case CauseInvalidSMEAddress:
		v[0].Encode(dia.Enumerated(5))
	case CauseUserNotSCUser:
		v[0].Encode(dia.Enumerated(6))
	}

	// SM-Diagnostic-Info
	if d.FCS != nil {
		t := dia.RawAVP{Code: 3305, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
		t.Encode(d.Encode())
		v = append(v, t)
	}

	a.Encode(v)
	return
}

func getSMDeliveryFailureCause(a dia.RawAVP) (c DeliveryFailureCause, d sms.DeliverReport, e error) {
	o := []dia.RawAVP{}
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else {
		e = a.Decode(&o)
	}
	for _, a := range o {
		switch a.Code {
		case 3304:
			s := new(dia.Enumerated)
			if !a.FlgV || !a.FlgM || a.FlgP {
				e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
			} else if e = a.Decode(s); e == nil {
				switch *s {
				case 0:
					c = CauseMemoryCapacityExceeded
				case 1:
					c = CauseEquipmentProtocolError
				case 2:
					c = CauseEquipmentNotSMEquipped
				case 3:
					c = CauseUnknownServiceCenter
				case 4:
					c = CauseSCCongestion
				case 5:
					c = CauseInvalidSMEAddress
				case 6:
					c = CauseUserNotSCUser
				default:
					e = dia.InvalidAVP(dia.DiameterInvalidAvpValue)
				}
			}
		case 3305:
			s := new([]byte)
			if !a.FlgV || !a.FlgM || a.FlgP {
				e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
			} else if e = a.Decode(s); e == nil {
				e = d.Decode(*s)
			}
		}
	}
	return
}

// SMDeliveryTimer AVP contain the value in seconds of the timer for SM Delivery.
func setSMDeliveryTimer(v uint32) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3306, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getSMDeliveryTimer(a dia.RawAVP) (v uint32, e error) {
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else {
		e = a.Decode(&v)
	}
	return
}

// SMDeliveryStartTime AVP contain the timestamp (in UTC) at which
// the SM Delivery Supervision Timer was started.
func setSMDeliveryStartTime(v time.Time) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3307, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	a.Encode(v)
	return a
}

func getSMDeliveryStartTime(a dia.RawAVP) (v time.Time, e error) {
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else {
		e = a.Decode(&v)
	}
	return
}

// MaximumRetransmissionTime AVP contain the maximum retransmission time (in UTC) until which
// the SMS-GMSC is capable to retransmit the MT Short Message.
func setMaximumRetransmissionTime(v time.Time) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3330, VenID: 10415, FlgV: true, FlgM: false, FlgP: false}
	a.Encode(v)
	return
}

func getMaximumRetransmissionTime(a dia.RawAVP) (v time.Time, e error) {
	if !a.FlgV || a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else {
		e = a.Decode(&v)
	}
	return
}

// SMSGMSCAddress AVP contain the E.164 number of the SMS-GMSC or SMS Router.
func setSMSGMSCAddress(v teldata.E164) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3332, VenID: 10415, FlgV: true, FlgM: false, FlgP: false}
	a.Encode(v.String())
	return
}

func getSMSGMSCAddress(a dia.RawAVP) (v teldata.E164, e error) {
	s := new(string)
	if !a.FlgV || a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil {
		v, e = teldata.ParseE164(*s)
	}
	return
}

// RequestedRetransmissionTime AVP contain the timestamp (in UTC) at which
// the SMS-GMSC is requested to retransmit the MT Short Message.
func setRequestedRetransmissionTime(v time.Time) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3331, VenID: 10415, FlgV: true, FlgM: false, FlgP: false}
	a.Encode(v)
	return
}

func getRequestedRetransmissionTime(a dia.RawAVP) (v time.Time, e error) {
	if !a.FlgV || a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else {
		e = a.Decode(&v)
	}
	return
}

/*
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
