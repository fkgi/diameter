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

func (v SCAddress) Avp() msg.Avp {
	a := msg.Avp{Code: uint32(3300), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(string(v))
	return a
}

// SCAddress get AVP value
func (o msg.GroupedAVP) SCAddress() (r []SCAddress) {
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

func (v SMRPUI) Avp() msg.Avp {
	a := msg.Avp{Code: uint32(3301), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode([]byte(v))
	return a
}

// SMRPUI get AVP value
func (o msg.GroupedAVP) SMRPUI() (r []SMRPUI) {
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

func (v TFRFlags) Avp() msg.Avp {
	a := msg.Avp{Code: uint32(3302), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	i := uint32(0)

	if v.MMS {
		i = i | 0x00000001
	}
	a.Encode(i)
	return a
}

// TFRFlags get AVP value
func (o msg.GroupedAVP) TFRFlags() (r []TFRFlags) {
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

// SMDeliveryFailureCause AVP contain cause of the failure of a SM delivery with an complementary information.
// If diag is nil, complementary information is empty.
type SMDeliveryFailureCause struct {
	Cause msg.Enumerated
	Diag  []byte
}

func (v SMDeliveryFailureCause) Avp() msg.Avp {
	a := msg.Avp{Code: uint32(3303), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	var t []msg.Avp

	// SM-Enumerated-Delivery-Failure-Cause
	{
		v := msg.Avp{Code: uint32(3304), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(v.Cause)
		t = append(t, v)
	}

	// SM-Diagnostic-Info
	if diag != nil {
		v := msg.Avp{Code: uint32(3305), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(v.Diag)
		t = append(t, v)
	}

	a.Encode(t)
	return a
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
func SMDeliveryTimer(i uint32) msg.Avp {
	a := msg.Avp{Code: uint32(3306), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// SMDeliveryStartTime AVP contain the timestamp (in UTC) at which
// the SM Delivery Supervision Timer was started.
func SMDeliveryStartTime(t time.Time) msg.Avp {
	a := msg.Avp{Code: uint32(3307), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(t)
	return a
}

// OFRFlags AVP is bit mask.
// When s6as6d set, the OFR message is sent on the Gdd interface (source node is an SGSN).
// When cleared, sent on the SGd interface (source node is an MME).
func OFRFlags(s6as6d bool) msg.Avp {
	a := msg.Avp{Code: uint32(3328), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	i := uint32(0)

	if s6as6d {
		i = i | 0x00000001
	}

	a.Encode(i)
	return a
}

// SMSMICorrelationID AVP ontain information identities used in the context
// of MSISDN-less SMS delivery in IMS
func SMSMICorrelationID(hssID, oURI, dURI string) msg.Avp {
	a := msg.Avp{Code: uint32(3324), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	var t []msg.Avp

	// HSS-ID
	if len(hssID) != 0 {
		v := msg.Avp{Code: uint32(3325), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
		v.Encode(hssID)
		t = append(t, v)
	}
	// Originating-SIP-URI
	if len(oURI) != 0 {
		v := msg.Avp{Code: uint32(3326), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
		v.Encode(oURI)
		t = append(t, v)
	}
	// Destination-SIP-URI
	if len(dURI) != 0 {
		v := msg.Avp{Code: uint32(3327), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
		v.Encode(dURI)
		t = append(t, v)
	}

	a.Encode(t)
	return a
}

// MaximumRetransmissionTime AVP contain the maximum retransmission time (in UTC) until which
// the SMS-GMSC is capable to retransmit the MT Short Message.
func MaximumRetransmissionTime(t time.Time) msg.Avp {
	a := msg.Avp{Code: uint32(3330), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(t)
	return a
}

// RequestedRetransmissionTime AVP contain the timestamp (in UTC) at which
// the SMS-GMSC is requested to retransmit the MT Short Message.
func RequestedRetransmissionTime(t time.Time) msg.Avp {
	a := msg.Avp{Code: uint32(3331), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(t)
	return a
}

// SMSGMSCAddress AVP contain the E.164 number of the SMS-GMSC or SMS Router.
func SMSGMSCAddress(msisdn string) msg.Avp {
	a := msg.Avp{Code: uint32(3332), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(msisdn)
	return a
}

// UserName AVP from RFC6733
func UserName(s string) msg.Avp {
	return msg.UserName(s)
}

// UserIdentifier AVP from ts29.336
func UserIdentifier(uname, msisdn, extid string, lmsi uint32) msg.Avp {
	return ts29336.UserIdentifier(uname, msisdn, extid, lmsi)
}

// MMENumberForMTSMS AVP from ts29.272
func MMENumberForMTSMS(b []byte) msg.Avp {
	return ts29272.MMENumberForMTSMS(b)
}

// SGSNNumber AVP from ts29.272
func SGSNNumber(b string) msg.Avp {
	a := ts29272.SGSNNumber(b)
	a.FlgM = false
	return a
}

// SupportedFeatures AVP from ts29.229
func SupportedFeatures(vendorID, featureID, featureList uint32) msg.Avp {
	return ts29229.SupportedFeatures(vendorID, featureID, featureList)
}

// FeatureListID AVP from ts29.229
func FeatureListID(i uint32) msg.Avp {
	return ts29229.FeatureListID(i)
}

// FeatureList AVP from ts29.229
func FeatureList(i uint32) msg.Avp {
	return ts29229.FeatureList(i)
}

// ExternalIdentifier AVP from ts29.336
func ExternalIdentifier(extid string) msg.Avp {
	a := ts29336.ExternalIdentifier(extid)
	a.FlgM = false
	return a
}
