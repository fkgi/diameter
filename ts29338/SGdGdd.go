package ts29338

import (
	"time"

	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/ts29229"
	"github.com/fkgi/diameter/ts29272"
	"github.com/fkgi/diameter/ts29336"
)

const (
	// DiameterErrorAbsentUser is Result-Code 5550
	DiameterErrorAbsentUser uint32 = 5550
	// DiameterErrorUserBusyForMtSms is Result-Code 5551
	DiameterErrorUserBusyForMtSms uint32 = 5551
	// DiameterErrorFacilityNotSupported is Result-Code 5552
	DiameterErrorFacilityNotSupported uint32 = 5552
	// DiameterErrorIlleagalUser is Result-Code 5553
	DiameterErrorIlleagalUser uint32 = 5553
	// DiameterErrorIlleagalEquipment is Result-Code 5554
	DiameterErrorIlleagalEquipment uint32 = 5554
	// DiameterErrorSmDeliveryFailure is Result-Code 5555
	DiameterErrorSmDeliveryFailure uint32 = 5555
	// DiameterErrorServiceNotSubscribed is Result-Code 5556
	DiameterErrorServiceNotSubscribed uint32 = 5556
	// DiameterErrorServiceBarred is Result-Code 5557
	DiameterErrorServiceBarred uint32 = 5557
	// DiameterErrorMwdListFull is Result-Code 5558
	DiameterErrorMwdListFull uint32 = 5558
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

const v3gpp uint32 = 10415

// SCAddress AVP contain the E164 number of the SMS-SC or MTC-IWF.
type SCAddress string

// Encode return AVP struct of this value
func (v SCAddress) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(3300), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(string(v))
	return a
}

// DecodeSCAddress get AVP value
func DecodeSCAddress(o msg.GroupedAVP) (r []SCAddress) {
	for _, a := range o {
		if a.Code == 3300 && a.VenID == v3gpp {
			s := new(string)
			a.Decode(s)
			r = append(r, SCAddress(*s))
		}
	}
	return
}

// SMRPUI AVP contain a short message transfer protocol data unit (TPDU).
// Maximum length is 200 octets.
type SMRPUI []byte

// Encode return AVP struct of this value
func (v SMRPUI) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(3301), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode([]byte(v))
	return a
}

// DecodeSMRPUI get AVP value
func DecodeSMRPUI(o msg.GroupedAVP) (r []SMRPUI) {
	for _, a := range o {
		if a.Code == 3301 && a.VenID == v3gpp {
			s := new([]byte)
			a.Decode(s)
			r = append(r, SMRPUI(*s))
		}
	}
	return
}

// TFRFlags AVP is bit mask.
// When moreMsgToSend set, the service centre has more short messages to send.
type TFRFlags struct {
	MMS bool // More message to send
}

// Encode return AVP struct of this value
func (v TFRFlags) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(3302), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	i := uint32(0)

	if v.MMS {
		i = i | 0x00000001
	}
	a.Encode(i)
	return a
}

// DecodeTFRFlags get AVP value
func DecodeTFRFlags(o msg.GroupedAVP) (r []TFRFlags) {
	for _, a := range o {
		if a.Code == 3302 && a.VenID == v3gpp {
			s := new(uint32)
			a.Decode(s)
			r = append(r, TFRFlags{
				MMS: (*s)&0x00000001 == 0x00000001})
		}
	}
	return
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
	a := msg.Avp{Code: uint32(3303), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	var t []msg.Avp

	// SM-Enumerated-Delivery-Failure-Cause
	{
		a := msg.Avp{Code: uint32(3304), VenID: v3gpp,
			FlgV: true, FlgM: true, FlgP: false}
		a.Encode(v.Cause)
		t = append(t, a)
	}

	// SM-Diagnostic-Info
	if v.Diag != nil {
		a := msg.Avp{Code: uint32(3305), VenID: v3gpp,
			FlgV: true, FlgM: true, FlgP: false}
		a.Encode(v.Diag)
		t = append(t, a)
	}

	a.Encode(msg.GroupedAVP(t))
	return a
}

// DecodeSMDeliveryFailureCause get AVP value
func DecodeSMDeliveryFailureCause(o msg.GroupedAVP) (r []SMDeliveryFailureCause) {
	for _, a := range o {
		if a.Code == 3303 && a.VenID == v3gpp {
			s := SMDeliveryFailureCause{}
			o2 := new(msg.GroupedAVP)
			a.Decode(o2)
			for _, a := range *o2 {
				if a.Code == 3304 && a.VenID == v3gpp {
					a.Decode(&s.Cause)
					break
				}
			}
			for _, a := range *o2 {
				if a.Code == 3305 && a.VenID == v3gpp {
					a.Decode(&s.Diag)
					break
				}
			}
			r = append(r, s)
		}
	}
	return
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
	a := msg.Avp{Code: uint32(3306), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// DecodeSMDeliveryTimer get AVP value
func DecodeSMDeliveryTimer(o msg.GroupedAVP) (r []SMDeliveryTimer) {
	for _, a := range o {
		if a.Code == 3306 && a.VenID == v3gpp {
			s := new(uint32)
			a.Decode(s)
			r = append(r, SMDeliveryTimer(*s))
		}
	}
	return
}

// SMDeliveryStartTime AVP contain the timestamp (in UTC) at which
// the SM Delivery Supervision Timer was started.
type SMDeliveryStartTime time.Time

// Encode return AVP struct of this value
func (v SMDeliveryStartTime) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(3307), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(time.Time(v))
	return a
}

// DecodeSMDeliveryStartTime get AVP value
func DecodeSMDeliveryStartTime(o msg.GroupedAVP) (r []SMDeliveryStartTime) {
	for _, a := range o {
		if a.Code == 3307 && a.VenID == v3gpp {
			s := new(time.Time)
			a.Decode(s)
			r = append(r, SMDeliveryStartTime(*s))
		}
	}
	return
}

// OFRFlags AVP is bit mask.
// When s6as6d set, the OFR message is sent on the Gdd interface (source node is an SGSN).
// When cleared, sent on the SGd interface (source node is an MME).
type OFRFlags struct {
	S6as6d bool
}

// Encode return AVP struct of this value
func (v OFRFlags) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(3328), VenID: v3gpp,
		FlgV: true, FlgM: false, FlgP: false}
	i := uint32(0)

	if v.S6as6d {
		i = i | 0x00000001
	}

	a.Encode(i)
	return a
}

// DecodeOFRFlags get AVP value
func DecodeOFRFlags(o msg.GroupedAVP) (r []OFRFlags) {
	for _, a := range o {
		if a.Code == 3328 && a.VenID == v3gpp {
			s := new(uint32)
			a.Decode(s)
			r = append(r, OFRFlags{
				S6as6d: (*s)&0x00000001 == 0x00000001})
		}
	}
	return
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
	a := msg.Avp{Code: uint32(3324), VenID: v3gpp,
		FlgV: true, FlgM: false, FlgP: false}
	var t []msg.Avp

	// HSS-ID
	if len(v.HSSID) != 0 {
		a := msg.Avp{Code: uint32(3325), VenID: v3gpp,
			FlgV: true, FlgM: false, FlgP: false}
		a.Encode(v.HSSID)
		t = append(t, a)
	}
	// Originating-SIP-URI
	if len(v.OrigSIPURI) != 0 {
		a := msg.Avp{Code: uint32(3326), VenID: v3gpp,
			FlgV: true, FlgM: false, FlgP: false}
		a.Encode(v.OrigSIPURI)
		t = append(t, a)
	}
	// Destination-SIP-URI
	if len(v.DestSIPURI) != 0 {
		a := msg.Avp{Code: uint32(3327), VenID: v3gpp,
			FlgV: true, FlgM: false, FlgP: false}
		a.Encode(v.DestSIPURI)
		t = append(t, a)
	}

	a.Encode(msg.GroupedAVP(t))
	return a
}

// DecodeSMSMICorrelationID get AVP value
func DecodeSMSMICorrelationID(o msg.GroupedAVP) (r []SMSMICorrelationID) {
	for _, a := range o {
		if a.Code == 3324 && a.VenID == v3gpp {
			s := SMSMICorrelationID{}
			o2 := new(msg.GroupedAVP)
			a.Decode(o2)
			for _, a := range *o2 {
				if a.Code == 3325 && a.VenID == v3gpp {
					a.Decode(&s.HSSID)
					break
				}
			}
			for _, a := range *o2 {
				if a.Code == 3326 && a.VenID == v3gpp {
					a.Decode(&s.OrigSIPURI)
					break
				}
			}
			for _, a := range *o2 {
				if a.Code == 3327 && a.VenID == v3gpp {
					a.Decode(&s.DestSIPURI)
					break
				}
			}
			r = append(r, s)
		}
	}
	return
}

// MaximumRetransmissionTime AVP contain the maximum retransmission time (in UTC) until which
// the SMS-GMSC is capable to retransmit the MT Short Message.
type MaximumRetransmissionTime time.Time

// Encode return AVP struct of this value
func (v MaximumRetransmissionTime) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(3330), VenID: v3gpp,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode(time.Time(v))
	return a
}

// DecodeMaximumRetransmissionTime get AVP value
func DecodeMaximumRetransmissionTime(o msg.GroupedAVP) (r []MaximumRetransmissionTime) {
	for _, a := range o {
		if a.Code == 3330 && a.VenID == v3gpp {
			s := new(time.Time)
			a.Decode(s)
			r = append(r, MaximumRetransmissionTime(*s))
		}
	}
	return
}

// RequestedRetransmissionTime AVP contain the timestamp (in UTC) at which
// the SMS-GMSC is requested to retransmit the MT Short Message.
type RequestedRetransmissionTime time.Time

// Encode return AVP struct of this value
func (v RequestedRetransmissionTime) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(3331), VenID: v3gpp,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode(time.Time(v))
	return a
}

// DecodeRequestedRetransmissionTime get AVP value
func DecodeRequestedRetransmissionTime(o msg.GroupedAVP) (r []RequestedRetransmissionTime) {
	for _, a := range o {
		if a.Code == 3331 && a.VenID == v3gpp {
			s := new(time.Time)
			a.Decode(s)
			r = append(r, RequestedRetransmissionTime(*s))
		}
	}
	return
}

// SMSGMSCAddress AVP contain the E.164 number of the SMS-GMSC or SMS Router.
type SMSGMSCAddress string

// Encode return AVP struct of this value
func (v SMSGMSCAddress) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(3332), VenID: v3gpp,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode(string(v))
	return a
}

// DecodeSMSGMSCAddress get AVP value
func DecodeSMSGMSCAddress(o msg.GroupedAVP) (r []SMSGMSCAddress) {
	for _, a := range o {
		if a.Code == 3332 && a.VenID == v3gpp {
			s := new(string)
			a.Decode(s)
			r = append(r, SMSGMSCAddress(*s))
		}
	}
	return
}

// UserName AVP from RFC6733
type UserName msg.UserName

// Encode return AVP struct of this value
func (v UserName) Encode() msg.Avp {
	return msg.UserName(v).Encode()
}

// DecodeUserName get AVP value
func DecodeUserName(o msg.GroupedAVP) (r []UserName) {
	return msg.DecodeUserName(o)
}

// UserIdentifier AVP from ts29.336
type UserIdentifier ts29336.UserIdentifier

// Encode return AVP struct of this value
func (v UserIdentifier) Encode() msg.Avp {
	return ts29336.UserIdentifier(v).Encode()
}

// MMENumberForMTSMS AVP from ts29.272
type MMENumberForMTSMS ts29272.MMENumberForMTSMS

// Encode return AVP struct of this value
func (v MMENumberForMTSMS) Encode() msg.Avp {
	return ts29272.MMENumberForMTSMS(v).Encode()
}

// SGSNNumber AVP from ts29.272
type SGSNNumber ts29272.SGSNNumber

// Encode return AVP struct of this value
func (v SGSNNumber) Encode() msg.Avp {
	a := ts29272.SGSNNumber(v).Encode()
	a.FlgM = false
	return a
}

// SupportedFeatures AVP from ts29.229
type SupportedFeatures ts29229.SupportedFeatures

// Encode return AVP struct of this value
func (v SupportedFeatures) Encode() msg.Avp {
	return ts29229.SupportedFeatures(v).Encode()
}

// FeatureListID AVP from ts29.229
type FeatureListID ts29229.FeatureListID

// Encode return AVP struct of this value
func (v FeatureListID) Encode() msg.Avp {
	return ts29229.FeatureListID(v).Encode()
}

// FeatureList AVP from ts29.229
type FeatureList ts29229.FeatureList

// Encode return AVP struct of this value
func (v FeatureList) Encode() msg.Avp {
	return ts29229.FeatureList(v).Encode()
}

// ExternalIdentifier AVP from ts29.336
type ExternalIdentifier ts29336.ExternalIdentifier

// Encode return AVP struct of this value
func (v ExternalIdentifier) Encode() msg.Avp {
	a := ts29336.ExternalIdentifier(v).Encode()
	a.FlgM = false
	return a
}
