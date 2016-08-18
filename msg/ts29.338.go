package msg

import "time"

const (
	DIAMETER_ERROR_ABSENT_USER            uint32 = 5550
	DIAMETER_ERROR_USER_BUSY_FOR_MT_SMS   uint32 = 5551
	DIAMETER_ERROR_FACILITY_NOT_SUPPORTED uint32 = 5552
	DIAMETER_ERROR_ILLEGAL_USER           uint32 = 5553
	DIAMETER_ERROR_ILLEGAL_EQUIPMENT      uint32 = 5554
	DIAMETER_ERROR_SM_DELIVERY_FAILURE    uint32 = 5555
	DIAMETER_ERROR_SERVICE_NOT_SUBSCRIBED uint32 = 5556
	DIAMETER_ERROR_SERVICE_BARRED         uint32 = 5557
	DIAMETER_ERROR_MWD_LIST_FULL          uint32 = 5558
)

// SCAddress AVP
func SCAddress(msisdn []byte) Avp {
	a := Avp{Code: uint32(3300), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(msisdn)
	return a
}

// SMRPUI AVP
func SMRPUI(s []byte) Avp {
	a := Avp{Code: uint32(3301), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(s)
	return a
}

// TFRFlags AVP
func TFRFlags(moreMsgToSend bool) Avp {
	a := Avp{Code: uint32(3302), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	i := uint32(0)

	if moreMsgToSend {
		i = i | 0x00000001
	}
	a.Encode(i)
	return a
}

// SMDeliveryFailureCause AVP
func SMDeliveryFailureCause(cause Enumerated, diag []byte) Avp {
	a := Avp{Code: uint32(3303), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []Avp

	// SM-Enumerated-Delivery-Failure-Cause
	{
		v := Avp{Code: uint32(3304), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(cause)
		t = append(t, v)
	}

	// SM-Diagnostic-Info
	if len(diag) != 0 {
		v := Avp{Code: uint32(3305), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(diag)
		t = append(t, v)
	}

	a.Encode(t)
	return a
}

// SMDeliveryTimer AVP
func SMDeliveryTimer(i uint32) Avp {
	a := Avp{Code: uint32(3306), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// SMDeliveryStartTime AVP
func SMDeliveryStartTime(t time.Time) Avp {
	a := Avp{Code: uint32(3307), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(t)
	return a
}

// SMSMICorrelationID AVP
func SMSMICorrelationID(hssID []byte, oURI, dURI string) Avp {
	a := Avp{Code: uint32(3324), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	var t []Avp

	// HSS-ID
	if len(hssID) != 0 {
		v := Avp{Code: uint32(3325), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
		v.Encode(hssID)
		t = append(t, v)
	}
	// Originating-SIP-URI
	if len(oURI) != 0 {
		v := Avp{Code: uint32(3326), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
		v.Encode(oURI)
		t = append(t, v)
	}
	// Destination-SIP-URI
	if len(dURI) != 0 {
		v := Avp{Code: uint32(3327), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
		v.Encode(dURI)
		t = append(t, v)
	}

	a.Encode(t)
	return a
}

// OFRFlags AVP
func OFRFlags(s6as6d bool) Avp {
	a := Avp{Code: uint32(3328), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	i := uint32(0)

	if s6as6d {
		i = i | 0x00000001
	}

	a.Encode(i)
	return a
}

// MaximumRetransmissionTime AVP
func MaximumRetransmissionTime(t time.Time) Avp {
	a := Avp{Code: uint32(3330), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(t)
	return a
}

// RequestedRetransmissionTime AVP
func RequestedRetransmissionTime(t time.Time) Avp {
	a := Avp{Code: uint32(3331), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(t)
	return a
}

// SMRPMTI AVP
func SMRPMTI(e Enumerated) Avp {
	a := Avp{Code: uint32(3308), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(e)
	return a
}

const (
	SMRPMTI_SM_DELIVER       Enumerated = 0
	SMRPMTI_SM_STATUS_REPORT Enumerated = 1
)

// SMRPSMEA AVP
func SMRPSMEA(smeAddr []byte) Avp {
	a := Avp{Code: uint32(3309), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(smeAddr)
	return a
}

// SRRFlags AVP
func SRRFlags(gprsIndicator, smRpPri, singleAttemptDelivery bool) Avp {
	a := Avp{Code: uint32(3310), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	i := uint32(0)

	if gprsIndicator {
		i = i | 0x00000001
	}
	if smRpPri {
		i = i | 0x00000002
	}
	if singleAttemptDelivery {
		i = i | 0x00000004
	}
	a.Encode(i)
	return a
}

// SMDeliveryNotIntended AVP
func SMDeliveryNotIntended(e Enumerated) Avp {
	a := Avp{Code: uint32(3311), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(e)
	return a
}

const (
	SMDeliveryNotIntended_ONLY_IMSI_REQUESTED    Enumerated = 0
	SMDeliveryNotIntended_ONLY_MCC_MNC_REQUESTED Enumerated = 1
)

// MWDStatus AVP
func MWDStatus(scAddrNotIncluded, mnrfSet, mcefSet, mnrgSet bool) Avp {
	a := Avp{Code: uint32(3312), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
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
func MMEAbsentUserDiagnosticSM(i uint32) Avp {
	a := Avp{Code: uint32(3313), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// MSCAbsentUserDiagnosticSM AVP
func MSCAbsentUserDiagnosticSM(i uint32) Avp {
	a := Avp{Code: uint32(3314), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// SGSNAbsentUserDiagnosticSM AVP
func SGSNAbsentUserDiagnosticSM(i uint32) Avp {
	a := Avp{Code: uint32(3315), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// SMDeliveryOutcome AVP
func SMDeliveryOutcome(e Enumerated, i uint32) Avp {
	a := Avp{Code: uint32(3316), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode()
	return a
}

// MMESMDeliveryOutcome AVP
func MMESMDeliveryOutcome(e Enumerated, i uint32) Avp {
	a := Avp{Code: uint32(3317), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []Avp

	t = append(t, SMDeliveryCause(e))
	t = append(t, AbsentUserDiagnosticSM(i))

	a.Encode(t)
	return a
}

// MSCSMDeliveryOutcome AVP
func MSCSMDeliveryOutcome(e Enumerated, i uint32) Avp {
	a := Avp{Code: uint32(3318), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []Avp

	t = append(t, SMDeliveryCause(e))
	t = append(t, AbsentUserDiagnosticSM(i))

	a.Encode(t)
	return a
}

// SGSNSMDeliveryOutcome AVP
func SGSNSMDeliveryOutcome(e Enumerated, i uint32) Avp {
	a := Avp{Code: uint32(3319), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []Avp

	t = append(t, SMDeliveryCause(e))
	t = append(t, AbsentUserDiagnosticSM(i))

	a.Encode(t)
	return a
}

// IPSMGWSMDeliveryOutcome AVP
func IPSMGWSMDeliveryOutcome(e Enumerated, i uint32) Avp {
	a := Avp{Code: uint32(3320), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []Avp

	t = append(t, SMDeliveryCause(e))
	t = append(t, AbsentUserDiagnosticSM(i))

	a.Encode(t)
	return a
}

// SMDeliveryCause AVP
func SMDeliveryCause(e Enumerated) Avp {
	a := Avp{Code: uint32(3321), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(e)
	return a
}

const (
	SMDeliveryCause_UE_MEMORY_CAPACITY_EXCEEDED Enumerated = 0
	SMDeliveryCause_ABSENT_USER                 Enumerated = 1
	SMDeliveryCause_SUCCESSFUL_TRANSFER         Enumerated = 2
)

// AbsentUserDiagnosticSM AVP
func AbsentUserDiagnosticSM(i uint32) Avp {
	a := Avp{Code: uint32(3322), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// RDRFlags AVP
func RDRFlags(singleAttemptDelivery bool) Avp {
	a := Avp{Code: uint32(3323), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	i := uint32(0)

	if singleAttemptDelivery {
		i = i | 0x00000001
	}

	a.Encode(i)
	return a
}

// MaximumUEAvailabilityTime AVP
func MaximumUEAvailabilityTime(t time.Time) Avp {
	a := Avp{Code: uint32(3329), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(t)
	return a
}
