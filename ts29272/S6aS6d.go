package ts29272

import (
	"strconv"

	"github.com/fkgi/diameter/msg"
)

const v3gpp uint32 = 10415

const (
	// DiameterErrorUserUnknown is Result-Code 5001
	DiameterErrorUserUnknown uint32 = 5001
	// DiameterErrorUnknownEpsSubscription is Result-Code 5420
	DiameterErrorUnknownEpsSubscription uint32 = 5420
	// DiameterErrorRatNotAllowed is Result-Code 5421
	DiameterErrorRatNotAllowed uint32 = 5421
	// DiameterErrorRoamingNotAllowed is Result-Code 5004
	DiameterErrorRoamingNotAllowed uint32 = 5004
	// DiameterErrorEquipmentUnknown is Result-Code 5422
	DiameterErrorEquipmentUnknown uint32 = 5422
	// DiameterErrorUnknownServingNode is Result-Code 5423
	DiameterErrorUnknownServingNode uint32 = 5423
	// DiameterAuthenticationDataUnavailable is Result-Code 4181
	DiameterAuthenticationDataUnavailable uint32 = 4181
	// DiameterErrorCancelSubscriptionPresent is Result-Code 4182
	DiameterErrorCancelSubscriptionPresent uint32 = 4182
)

// ULRFlags AVP
func ULRFlags(singleReg, s6as6d, skipSubsData, gprsSubsData, nodeType, initAttach, psLcsNotSupportedByUE, smsOnly bool) msg.Avp {
	a := msg.Avp{Code: uint32(1405), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	i := uint32(0)

	if singleReg {
		i = i | 0x00000001
	}
	if s6as6d {
		i = i | 0x00000002
	}
	if skipSubsData {
		i = i | 0x00000004
	}
	if gprsSubsData {
		i = i | 0x00000008
	}
	if nodeType {
		i = i | 0x00000010
	}
	if initAttach {
		i = i | 0x00000020
	}
	if psLcsNotSupportedByUE {
		i = i | 0x00000040
	}
	if smsOnly {
		i = i | 0x00000080
	}

	a.Encode(i)
	return a
}

// VisitedPLMNID AVP
func VisitedPLMNID(mcc, mnc string) msg.Avp {
	a := msg.Avp{Code: uint32(1407), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	b := make([]byte, 3)

	m := [3]int{0, 0, 0}
	if len(mcc) == 3 {
		for i := 0; i < 3; i++ {
			m[i], _ = strconv.Atoi(mcc[i : i+1])
		}
	}
	b[0] = (byte(m[1]) << 4) | byte(m[0])
	b[1] = byte(m[2])

	m = [3]int{0, 0, 0}
	if len(mnc) == 2 {
		for i := 0; i < 2; i++ {
			m[i], _ = strconv.Atoi(mnc[i : i+1])
		}
		m[2] = 0xf
	} else if len(mnc) == 3 {
		for i := 0; i < 3; i++ {
			m[i], _ = strconv.Atoi(mnc[i : i+1])
		}
	}
	b[1] = (byte(m[2]) << 4) | b[1]
	b[2] = (byte(m[1]) << 4) | byte(m[0])

	a.Encode(b)
	return a
}

// TerminalInformation AVP
func TerminalInformation(imei string, meid []byte, version string) msg.Avp {
	a := msg.Avp{Code: uint32(1401), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}

	var t []msg.Avp
	// IMEI
	if len(imei) != 0 {
		v := msg.Avp{Code: uint32(1402), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(imei)
		t = append(t, v)
	}
	// 3GPP2-MEID
	if meid != nil && len(meid) != 0 {
		v := msg.Avp{Code: uint32(1471), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(meid)
		t = append(t, v)
	}
	// Sofrware-Version
	if len(version) != 0 {
		v := msg.Avp{Code: uint32(1403), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(version)
		t = append(t, v)
	}

	a.Encode(msg.GroupedAVP(t))
	return a
}

// UESRVCCCapability AVP
func UESRVCCCapability(e msg.Enumerated) msg.Avp {
	a := msg.Avp{Code: uint32(1615), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(e)
	return a
}

const (
	// UeSrvccNotSupported is Enumerated value 0
	UeSrvccNotSupported msg.Enumerated = 0
	// UeSrvccSupported is Enumerated value 1
	UeSrvccSupported msg.Enumerated = 1
)

// UESRVCCCapabilityValue is definition of Enumerated value
var UESRVCCCapabilityValue = struct {
	UeSrvccNotSupported msg.Enumerated
	UeSrvccSupported    msg.Enumerated
}{0, 1}

// SGSNNumber AVP
type SGSNNumber string

// Encode return AVP struct of this value
func (v SGSNNumber) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(1489), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(string(v))
	return a
}

// HomogeneousSupportOfIMSVoiceOverPSSessions AVP
func HomogeneousSupportOfIMSVoiceOverPSSessions(e msg.Enumerated) msg.Avp {
	a := msg.Avp{Code: uint32(1493), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(e)
	return a
}

const (
	// NotSupported is Enumerated value 0
	NotSupported msg.Enumerated = 0
	// Supported is Enumerated value 1
	Supported msg.Enumerated = 1
)

// ContextIdentifire AVP
func ContextIdentifire(i uint32) msg.Avp {
	a := msg.Avp{Code: uint32(1423), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// ActiveAPN AVP
func ActiveAPN(id uint32) msg.Avp {
	a := msg.Avp{Code: uint32(1612), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []msg.Avp
	t = append(t, ContextIdentifire(id))
	/*
		<!--<avp name="Service-Selection"	value=""></avp>-->
		<!--<avp name="MIP6-Agent-Info"	value=""></avp>-->
		<!--<avp name="Visited-Network-Identifier"	value=""></avp>-->
		<!--<avp name="Specific-APN-Info"	value=""></avp>-->
	*/
	a.Encode(msg.GroupedAVP(t))
	return a
}

// EquivalentPLMNList AVP
func EquivalentPLMNList(plmns [][2]string) msg.Avp {
	a := msg.Avp{Code: uint32(1637), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	var t []msg.Avp
	for _, p := range plmns {
		t = append(t, VisitedPLMNID(p[0], p[1]))
	}
	a.Encode(msg.GroupedAVP(t))
	return a
}

// MMENumberForMTSMS AVP
type MMENumberForMTSMS []byte

// Encode return AVP struct of this value
func (v MMENumberForMTSMS) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(1645), VenID: v3gpp,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode([]byte(v))
	return a
}

// SMSRegisterRequest AVP
func SMSRegisterRequest(e msg.Enumerated) msg.Avp {
	a := msg.Avp{Code: uint32(1648), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(e)
	return a
}

const (
	// SmsRegistrationRequired is Enumerated value 0
	SmsRegistrationRequired msg.Enumerated = 0
	// SmsRegistrationNotPreferred is Enumerated value 1
	SmsRegistrationNotPreferred msg.Enumerated = 1
	// NotPreference is Enumerated value 2
	NotPreference msg.Enumerated = 2
)

// SGsMMEIdentity AVP
func SGsMMEIdentity(s string) msg.Avp {
	a := msg.Avp{Code: uint32(1664), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(s)
	return a
}

// CoupledNodeDiameterID AVP
func CoupledNodeDiameterID(i msg.DiameterIdentity) msg.Avp {
	a := msg.Avp{Code: uint32(1666), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// RequestedEUTRANAuthenticationInfo AVP
func RequestedEUTRANAuthenticationInfo(num uint32, resyncInfo []byte, immRespPref bool) msg.Avp {
	a := msg.Avp{Code: uint32(1408), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []msg.Avp
	// Number-Of-Requested-Vectors
	if num != 0 {
		v := msg.Avp{Code: uint32(1410), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(num)
		t = append(t, v)
	}
	// Re-Synchronization-Info
	if resyncInfo != nil {
		v := msg.Avp{Code: uint32(1411), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(resyncInfo)
		t = append(t, v)
	}
	// Immediate-Response-Preferred
	if immRespPref {
		v := msg.Avp{Code: uint32(1412), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(uint32(1))
		t = append(t, v)
	}
	a.Encode(msg.GroupedAVP(t))
	return a
}

// RequestedUTRANGERANAuthenticationInfo AVP
func RequestedUTRANGERANAuthenticationInfo(num uint32, resyncInfo []byte, immRespPref bool) msg.Avp {
	a := RequestedEUTRANAuthenticationInfo(num, resyncInfo, immRespPref)
	a.Code = uint32(1409)
	return a
}

// NORFlags AVP
func NORFlags(singleReg, sgsnRestrict, readySmSgsn, ueReachableMme, reserved, ueReachableSgsn, readySmMme, homogeneousSupportIMSVoPSSession, s6as6d, remMmeRegSm bool) msg.Avp {
	a := msg.Avp{Code: uint32(1443), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	i := uint32(0)

	if singleReg {
		i = i | 0x00000001
	}
	if sgsnRestrict {
		i = i | 0x00000002
	}
	if readySmSgsn {
		i = i | 0x00000004
	}
	if ueReachableMme {
		i = i | 0x00000008
	}
	// if reserved {
	// 	i = i | 0x00000010
	// }
	if ueReachableSgsn {
		i = i | 0x00000020
	}
	if readySmMme {
		i = i | 0x00000040
	}
	if homogeneousSupportIMSVoPSSession {
		i = i | 0x00000080
	}
	if s6as6d {
		i = i | 0x00000100
	}
	if remMmeRegSm {
		i = i | 0x00000200
	}

	a.Encode(i)
	return a
}
