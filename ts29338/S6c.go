package ts29338

import (
	"time"

	"github.com/fkgi/diameter/msg"
)

/*
 <SRR> ::= < Diameter Header: 8388647, REQ, PXY, 16777312 >
           < Session-Id >
		   [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ Destination-Host ]
           { Destination-Realm }
           [ MSISDN ]
           [ User-Name ]
           [ SMSMI-Correlation-ID ]
         * [ Supported-Features ]
           [ SC-Address ]
           [ SM-RP-MTI ]
           [ SM-RP-SMEA ]
           [ SRR-Flags ]
           [ SM-Delivery-Not-Intended ]
         * [ AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]

 <SRA> ::= < Diameter Header: 8388647, PXY, 16777312 >
           < Session-Id >
		   [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           [ Result-Code ]
           [ Experimental-Result ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ User-Name ]
         * [ Supported-Features ]
           [ Serving-Node ]
           [ Additional-Serving-Node ]
           [ LMSI ]
           [ User-Identifier ]
           [ MWD-Status ]
           [ MME-Absent-User-Diagnostic-SM ]
           [ MSC-Absent-User-Diagnostic-SM ]
           [ SGSN-Absent-User-Diagnostic-SM ]
         * [ AVP ]
         * [ Failed-AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/

/*
 <ALR> ::= < Diameter Header: 8388648, REQ, PXY, 16777312 >
           < Session-Id >
           [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ Destination-Host ]
           { Destination-Realm }
           { SC-Address }
           { User-Identifier }
           [ SMSMI-Correlation-ID ]
           [ Maximum-UE-Availability-Time ]
           [ SMS-GMSC-Alert-Event ]
           [ Serving-Node ]
         * [ Supported-Features ]
         * [ AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
 <ALA> ::= < Diameter Header: 8388648, PXY, 16777312 >
           < Session-Id >
           [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           [ Result-Code ]
           [ Experimental-Result ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
         * [ Supported-Features ]
         * [ AVP ]
         * [ Failed-AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/

/*
 <RDR> ::= < Diameter Header: 8388649, REQ, PXY, 16777312 >
           < Session-Id >
           [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ Destination-Host ]
           { Destination-Realm }
         * [ Supported-Features ]
           { User-Identifier }
           [ SMSMI-Correlation-ID ]
           { SC-Address }
           { SM-Delivery-Outcome }
           [ RDR-Flags ]
         * [ AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
 <RDA> ::= < Diameter Header: 8388649, PXY, 16777312 >
           < Session-Id >
           [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           [ Result-Code ]
           [ Experimental-Result ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
         * [ Supported-Features ]
           [ User-Identifier ]
         * [ AVP ]
         * [ Failed-AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/

// SMRPMTI AVP contain the RP-Message Type Indicator of the Short Message.
type SMRPMTI msg.Enumerated

// Encode return AVP struct of this value
func (v SMRPMTI) Encode() msg.Avp {
	a := msg.Avp{Code: 3308, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(msg.Enumerated(v))
	return a
}

// GetSMRPMTI get AVP value
func GetSMRPMTI(o msg.GroupedAVP) (SMRPMTI, bool) {
	s := new(msg.Enumerated)
	if a, ok := o.Get(3308, 10415); ok {
		a.Decode(s)
	} else {
		return 0, false
	}
	return SMRPMTI(*s), true
}

const (
	// SmDeliver is Enumerated value 0
	SmDeliver msg.Enumerated = 0
	// SmStatusReport is Enumerated value 1
	SmStatusReport msg.Enumerated = 1
)

// SMRPSMEA AVP contain the RP-Originating SME-address of the Short Message Entity that has originated the SM.
// It shall be formatted according to the formatting rules of the address fields described in 3GPP TS 23.040.
type SMRPSMEA []byte

// Encode return AVP struct of this value
func (v SMRPSMEA) Encode() msg.Avp {
	a := msg.Avp{Code: 3309, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode([]byte(v))
	return a
}

// GetSMRPSMEA get AVP value
func GetSMRPSMEA(o msg.GroupedAVP) (SMRPSMEA, bool) {
	s := new([]byte)
	if a, ok := o.Get(3309, 10415); ok {
		a.Decode(s)
	} else {
		return nil, false
	}
	return SMRPSMEA(*s), true
}

// SRRFlags AVP contain a bit mask.
// gprsIndicator shall be ture if the SMS-GMSC supports receiving of two serving nodes addresses from the HSS.
// smRpPri shall be true if the delivery of the short message shall be attempted when
// a service centre address is already contained in the Message Waiting Data file.
// singleAttempt if true indicates that only one delivery attempt shall be performed for this particular SM.
type SRRFlags struct {
	GprsIndicator bool
	SmRpPri       bool
	SingleAttempt bool
}

// Encode return AVP struct of this value
func (v SRRFlags) Encode() msg.Avp {
	a := msg.Avp{Code: 3310, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	i := uint32(0)

	if v.GprsIndicator {
		i = i | 0x00000001
	}
	if v.SmRpPri {
		i = i | 0x00000002
	}
	if v.SingleAttempt {
		i = i | 0x00000004
	}
	a.Encode(i)
	return a
}

// GetSRRFlags get AVP value
func GetSRRFlags(o msg.GroupedAVP) (SRRFlags, bool) {
	s := new(uint32)
	if a, ok := o.Get(3310, 10415); ok {
		a.Decode(s)
	} else {
		return SRRFlags{}, false
	}
	return SRRFlags{
		GprsIndicator: (*s)&0x00000001 == 0x00000001,
		SmRpPri:       (*s)&0x00000002 == 0x00000002,
		SingleAttempt: (*s)&0x00000004 == 0x00000004}, true
}

// SMDeliveryNotIntended AVP indicate by its presence
// that delivery of a short message is not intended.
type SMDeliveryNotIntended msg.Enumerated

// Encode return AVP struct of this value
func (v SMDeliveryNotIntended) Encode() msg.Avp {
	a := msg.Avp{Code: 3311, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(msg.Enumerated(v))
	return a
}

// GetSMDeliveryNotIntended get AVP value
func GetSMDeliveryNotIntended(o msg.GroupedAVP) (SMDeliveryNotIntended, bool) {
	s := new(msg.Enumerated)
	if a, ok := o.Get(3311, 10415); ok {
		a.Decode(s)
	} else {
		return 0, false
	}
	return SMDeliveryNotIntended(*s), true
}

const (
	// OnlyImsiRequested is Enumerated value 0
	OnlyImsiRequested msg.Enumerated = 0
	// OnlyMccMnCRequested is Enumerated value 1
	OnlyMccMnCRequested msg.Enumerated = 1
)

// MWDStatus AVP
type MWDStatus struct {
	ScAddrNotIncluded bool
	MnrfSet           bool
	McefSet           bool
	MnrgSet           bool
}

// Encode return AVP struct of this value
func (v MWDStatus) Encode() msg.Avp {
	a := msg.Avp{Code: 3312, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	i := uint32(0)

	if v.ScAddrNotIncluded {
		i = i | 0x00000001
	}
	if v.MnrfSet {
		i = i | 0x00000002
	}
	if v.McefSet {
		i = i | 0x00000004
	}
	if v.MnrgSet {
		i = i | 0x00000008
	}
	a.Encode(i)
	return a
}

// GetMWDStatus get AVP value
func GetMWDStatus(o msg.GroupedAVP) (MWDStatus, bool) {
	s := new(uint32)
	if a, ok := o.Get(3312, 10415); ok {
		a.Decode(s)
	} else {
		return MWDStatus{}, false
	}
	return MWDStatus{
		ScAddrNotIncluded: (*s)&0x00000001 == 0x00000001,
		MnrfSet:           (*s)&0x00000002 == 0x00000002,
		McefSet:           (*s)&0x00000004 == 0x00000004,
		MnrgSet:           (*s)&0x00000008 == 0x00000008}, true
}

// MMEAbsentUserDiagnosticSM AVP
type MMEAbsentUserDiagnosticSM uint32

// Encode return AVP struct of this value
func (v MMEAbsentUserDiagnosticSM) Encode() msg.Avp {
	a := msg.Avp{Code: 3313, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetMMEAbsentUserDiagnosticSM get AVP value
func GetMMEAbsentUserDiagnosticSM(o msg.GroupedAVP) (MMEAbsentUserDiagnosticSM, bool) {
	s := new(uint32)
	if a, ok := o.Get(3313, 10415); ok {
		a.Decode(s)
	} else {
		return 0, false
	}
	return MMEAbsentUserDiagnosticSM(*s), true
}

// MSCAbsentUserDiagnosticSM AVP
type MSCAbsentUserDiagnosticSM uint32

// Encode return AVP struct of this value
func (v MSCAbsentUserDiagnosticSM) Encode() msg.Avp {
	a := msg.Avp{Code: 3314, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetMSCAbsentUserDiagnosticSM get AVP value
func GetMSCAbsentUserDiagnosticSM(o msg.GroupedAVP) (MSCAbsentUserDiagnosticSM, bool) {
	s := new(uint32)
	if a, ok := o.Get(3314, 10415); ok {
		a.Decode(s)
	} else {
		return 0, false
	}
	return MSCAbsentUserDiagnosticSM(*s), true
}

// SGSNAbsentUserDiagnosticSM AVP
type SGSNAbsentUserDiagnosticSM uint32

// Encode return AVP struct of this value
func (v SGSNAbsentUserDiagnosticSM) Encode() msg.Avp {
	a := msg.Avp{Code: 3315, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetSGSNAbsentUserDiagnosticSM get AVP value
func GetSGSNAbsentUserDiagnosticSM(o msg.GroupedAVP) (SGSNAbsentUserDiagnosticSM, bool) {
	s := new(uint32)
	if a, ok := o.Get(3315, 10415); ok {
		a.Decode(s)
	} else {
		return 0, false
	}
	return SGSNAbsentUserDiagnosticSM(*s), true
}

// SMDeliveryOutcome AVP
type SMDeliveryOutcome struct {
	E msg.Enumerated
	I uint32
}

// Encode return AVP struct of this value
func (v SMDeliveryOutcome) Encode() msg.Avp {
	a := msg.Avp{Code: 3316, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	// a.Encode()
	return a
}

// MMESMDeliveryOutcome AVP
type MMESMDeliveryOutcome struct {
	SMDeliveryCause        SMDeliveryCause
	AbsentUserDiagnosticSM AbsentUserDiagnosticSM
}

// Encode return AVP struct of this value
func (v MMESMDeliveryOutcome) Encode() msg.Avp {
	a := msg.Avp{Code: 3317, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	var t []msg.Avp

	t = append(t, v.SMDeliveryCause.Encode())
	t = append(t, v.AbsentUserDiagnosticSM.Encode())

	a.Encode(msg.GroupedAVP(t))
	return a
}

// MSCSMDeliveryOutcome AVP
type MSCSMDeliveryOutcome struct {
	SMDeliveryCause        SMDeliveryCause
	AbsentUserDiagnosticSM AbsentUserDiagnosticSM
}

// Encode return AVP struct of this value
func (v MSCSMDeliveryOutcome) Encode() msg.Avp {
	a := msg.Avp{Code: 3318, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	var t []msg.Avp

	t = append(t, v.SMDeliveryCause.Encode())
	t = append(t, v.AbsentUserDiagnosticSM.Encode())

	a.Encode(msg.GroupedAVP(t))
	return a
}

// SGSNSMDeliveryOutcome AVP
type SGSNSMDeliveryOutcome struct {
	SMDeliveryCause        SMDeliveryCause
	AbsentUserDiagnosticSM AbsentUserDiagnosticSM
}

// Encode return AVP struct of this value
func (v SGSNSMDeliveryOutcome) Encode() msg.Avp {
	a := msg.Avp{Code: 3319, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	var t []msg.Avp

	t = append(t, v.SMDeliveryCause.Encode())
	t = append(t, v.AbsentUserDiagnosticSM.Encode())

	a.Encode(msg.GroupedAVP(t))
	return a
}

// IPSMGWSMDeliveryOutcome AVP
type IPSMGWSMDeliveryOutcome struct {
	SMDeliveryCause        SMDeliveryCause
	AbsentUserDiagnosticSM AbsentUserDiagnosticSM
}

// Encode return AVP struct of this value
func (v IPSMGWSMDeliveryOutcome) Encode() msg.Avp {
	a := msg.Avp{Code: 3320, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	var t []msg.Avp

	t = append(t, v.SMDeliveryCause.Encode())
	t = append(t, v.AbsentUserDiagnosticSM.Encode())

	a.Encode(msg.GroupedAVP(t))
	return a
}

// SMDeliveryCause AVP
type SMDeliveryCause msg.Enumerated

// Encode return AVP struct of this value
func (v SMDeliveryCause) Encode() msg.Avp {
	a := msg.Avp{Code: 3321, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(msg.Enumerated(v))
	return a
}

// GetSMDeliveryCause get AVP value
func GetSMDeliveryCause(o msg.GroupedAVP) (SMDeliveryCause, bool) {
	s := new(msg.Enumerated)
	if a, ok := o.Get(3321, 10415); ok {
		a.Decode(s)
	} else {
		return 0, false
	}
	return SMDeliveryCause(*s), true
}

const (
	// UeMemoryCapacityExceeded is Enumerated value 0
	UeMemoryCapacityExceeded msg.Enumerated = 0
	// AbsentUser is Enumerated value 1
	AbsentUser msg.Enumerated = 1
	// SuccessfulTransfer is Enumerated value 2
	SuccessfulTransfer msg.Enumerated = 2
)

// AbsentUserDiagnosticSM AVP
type AbsentUserDiagnosticSM uint32

// Encode return AVP struct of this value
func (v AbsentUserDiagnosticSM) Encode() msg.Avp {
	a := msg.Avp{Code: 3322, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetAbsentUserDiagnosticSM get AVP value
func GetAbsentUserDiagnosticSM(o msg.GroupedAVP) (AbsentUserDiagnosticSM, bool) {
	s := new(uint32)
	if a, ok := o.Get(3322, 10415); ok {
		a.Decode(s)
	} else {
		return 0, false
	}
	return AbsentUserDiagnosticSM(*s), true
}

// RDRFlags AVP
type RDRFlags struct {
	SingleAttemptDelivery bool
}

// Encode return AVP struct of this value
func (v RDRFlags) Encode() msg.Avp {
	a := msg.Avp{Code: 3323, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	i := uint32(0)

	if v.SingleAttemptDelivery {
		i = i | 0x00000001
	}
	a.Encode(i)
	return a
}

// GetRDRFlags get AVP value
func GetRDRFlags(o msg.GroupedAVP) (RDRFlags, bool) {
	s := new(uint32)
	if a, ok := o.Get(3323, 10415); ok {
		a.Decode(s)
	} else {
		return RDRFlags{}, false
	}
	return RDRFlags{
		SingleAttemptDelivery: (*s)&0x00000001 == 0x00000001}, true
}

// MaximumUEAvailabilityTime AVP
type MaximumUEAvailabilityTime time.Time

// Encode return AVP struct of this value
func (v MaximumUEAvailabilityTime) Encode() msg.Avp {
	a := msg.Avp{Code: 3329, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode(time.Time(v))
	return a
}

// GetMaximumUEAvailabilityTime get AVP value
func GetMaximumUEAvailabilityTime(o msg.GroupedAVP) (MaximumUEAvailabilityTime, bool) {
	s := new(time.Time)
	if a, ok := o.Get(3329, 10415); ok {
		a.Decode(s)
	} else {
		return MaximumUEAvailabilityTime{}, false
	}
	return MaximumUEAvailabilityTime(*s), true
}
