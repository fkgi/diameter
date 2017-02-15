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
func SMRPMTI(e msg.Enumerated) msg.Avp {
	a := msg.Avp{Code: uint32(3308), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(e)
	return a
}

const (
	// SmDeliver is Enumerated value 0
	SmDeliver msg.Enumerated = 0
	// SmStatusReport is Enumerated value 1
	SmStatusReport msg.Enumerated = 1
)

// SMRPSMEA AVP contain the RP-Originating SME-address of the Short Message Entity that has originated the SM.
// It shall be formatted according to the formatting rules of the address fields described in 3GPP TS 23.040.
func SMRPSMEA(smeAddr []byte) msg.Avp {
	a := msg.Avp{Code: uint32(3309), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(smeAddr)
	return a
}

// SRRFlags AVP contain a bit mask.
// gprsIndicator shall be ture if the SMS-GMSC supports receiving of two serving nodes addresses from the HSS.
// smRpPri shall be true if the delivery of the short message shall be attempted when
// a service centre address is already contained in the Message Waiting Data file.
// singleAttempt if true indicates that only one delivery attempt shall be performed for this particular SM.
func SRRFlags(gprsIndicator, smRpPri, singleAttempt bool) msg.Avp {
	a := msg.Avp{Code: uint32(3310), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	i := uint32(0)

	if gprsIndicator {
		i = i | 0x00000001
	}
	if smRpPri {
		i = i | 0x00000002
	}
	if singleAttempt {
		i = i | 0x00000004
	}
	a.Encode(i)
	return a
}

// SMDeliveryNotIntended AVP indicate by its presence that delivery of a short message is not intended.
func SMDeliveryNotIntended(e msg.Enumerated) msg.Avp {
	a := msg.Avp{Code: uint32(3311), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(e)
	return a
}

const (
	// OnlyImsiRequested is Enumerated value 0
	OnlyImsiRequested msg.Enumerated = 0
	// OnlyMccMnCRequested is Enumerated value 1
	OnlyMccMnCRequested msg.Enumerated = 1
)

// MWDStatus AVP
func MWDStatus(scAddrNotIncluded, mnrfSet, mcefSet, mnrgSet bool) msg.Avp {
	a := msg.Avp{Code: uint32(3312), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	i := uint32(0)

	if scAddrNotIncluded {
		i = i | 0x00000001
	}
	if mnrfSet {
		i = i | 0x00000002
	}
	if mcefSet {
		i = i | 0x00000004
	}
	if mnrgSet {
		i = i | 0x00000008
	}
	a.Encode(i)
	return a
}

// MMEAbsentUserDiagnosticSM AVP
func MMEAbsentUserDiagnosticSM(i uint32) msg.Avp {
	a := msg.Avp{Code: uint32(3313), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// MSCAbsentUserDiagnosticSM AVP
func MSCAbsentUserDiagnosticSM(i uint32) msg.Avp {
	a := msg.Avp{Code: uint32(3314), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// SGSNAbsentUserDiagnosticSM AVP
func SGSNAbsentUserDiagnosticSM(i uint32) msg.Avp {
	a := msg.Avp{Code: uint32(3315), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// SMDeliveryOutcome AVP
func SMDeliveryOutcome(e msg.Enumerated, i uint32) msg.Avp {
	a := msg.Avp{Code: uint32(3316), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	// a.Encode()
	return a
}

// MMESMDeliveryOutcome AVP
func MMESMDeliveryOutcome(e msg.Enumerated, i uint32) msg.Avp {
	a := msg.Avp{Code: uint32(3317), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []msg.Avp

	t = append(t, SMDeliveryCause(e))
	t = append(t, AbsentUserDiagnosticSM(i))

	a.Encode(t)
	return a
}

// MSCSMDeliveryOutcome AVP
func MSCSMDeliveryOutcome(e msg.Enumerated, i uint32) msg.Avp {
	a := msg.Avp{Code: uint32(3318), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []msg.Avp

	t = append(t, SMDeliveryCause(e))
	t = append(t, AbsentUserDiagnosticSM(i))

	a.Encode(t)
	return a
}

// SGSNSMDeliveryOutcome AVP
func SGSNSMDeliveryOutcome(e msg.Enumerated, i uint32) msg.Avp {
	a := msg.Avp{Code: uint32(3319), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []msg.Avp

	t = append(t, SMDeliveryCause(e))
	t = append(t, AbsentUserDiagnosticSM(i))

	a.Encode(t)
	return a
}

// IPSMGWSMDeliveryOutcome AVP
func IPSMGWSMDeliveryOutcome(e msg.Enumerated, i uint32) msg.Avp {
	a := msg.Avp{Code: uint32(3320), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []msg.Avp

	t = append(t, SMDeliveryCause(e))
	t = append(t, AbsentUserDiagnosticSM(i))

	a.Encode(t)
	return a
}

// SMDeliveryCause AVP
func SMDeliveryCause(e msg.Enumerated) msg.Avp {
	a := msg.Avp{Code: uint32(3321), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(e)
	return a
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
func AbsentUserDiagnosticSM(i uint32) msg.Avp {
	a := msg.Avp{Code: uint32(3322), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// RDRFlags AVP
func RDRFlags(singleAttemptDelivery bool) msg.Avp {
	a := msg.Avp{Code: uint32(3323), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	i := uint32(0)

	if singleAttemptDelivery {
		i = i | 0x00000001
	}

	a.Encode(i)
	return a
}

// MaximumUEAvailabilityTime AVP
func MaximumUEAvailabilityTime(t time.Time) msg.Avp {
	a := msg.Avp{Code: uint32(3329), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(t)
	return a
}
