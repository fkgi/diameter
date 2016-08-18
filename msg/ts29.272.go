package msg

import (
	"strconv"
)

const (
	DIAMETER_ERROR_USER_UNKNOWN               uint32 = 5001
	DIAMETER_ERROR_UNKNOWN_EPS_SUBSCRIPTION   uint32 = 5420
	DIAMETER_ERROR_RAT_NOT_ALLOWED            uint32 = 5421
	DIAMETER_ERROR_ROAMING_NOT_ALLOWED        uint32 = 5004
	DIAMETER_ERROR_EQUIPMENT_UNKNOWN          uint32 = 5422
	DIAMETER_ERROR_UNKOWN_SERVING_NODE        uint32 = 5423
	DIAMETER_AUTHENTICATION_DATA_UNAVAILABLE  uint32 = 4181
	DIAMETER_ERROR_CAMEL_SUBSCRIPTION_PRESENT uint32 = 4182
)

// ULRFlags AVP
func ULRFlags(singleReg, s6as6d, skipSubsData, gprsSubsData, nodeType, initAttach, psLcsNotSupportedByUE, smsOnly bool) Avp {
	a := Avp{Code: uint32(1405), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
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
func VisitedPLMNID(mcc, mnc string) Avp {
	a := Avp{Code: uint32(1407), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
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
func TerminalInformation(imei string, meid []byte, version string) Avp {
	a := Avp{Code: uint32(1401), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}

	var t []Avp
	// IMEI
	if len(imei) != 0 {
		v := Avp{Code: uint32(1402), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(imei)
		t = append(t, v)
	}
	// 3GPP2-MEID
	if meid != nil && len(meid) != 0 {
		v := Avp{Code: uint32(1471), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(meid)
		t = append(t, v)
	}
	// Sofrware-Version
	if len(version) != 0 {
		v := Avp{Code: uint32(1403), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(version)
		t = append(t, v)
	}

	a.Encode(t)
	return a
}

// UESRVCCCapability AVP
func UESRVCCCapability(e Enumerated) Avp {
	a := Avp{Code: uint32(1615), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(e)
	return a
}

const (
	UESRVCCCapability_UE_SRVCC_NOT_SUPPORTED Enumerated = 0
	UESRVCCCapability_UE_SRVCC_SUPPORTED     Enumerated = 1
)

// SGSNNumber AVP
func SGSNNumber(b []byte) Avp {
	a := Avp{Code: uint32(1489), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(b)
	return a
}

// HomogeneousSupportOfIMSVoiceOverPSSessions AVP
func HomogeneousSupportOfIMSVoiceOverPSSessions(e Enumerated) Avp {
	a := Avp{Code: uint32(1493), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(e)
	return a
}

const (
	HomogeneousSupportOfIMSVoiceOverPSSessions_NOT_SUPPORTED Enumerated = 0
	HomogeneousSupportOfIMSVoiceOverPSSessions_SUPPORTED     Enumerated = 1
)

// ContextIdentifire AVP
func ContextIdentifire(i uint32) Avp {
	a := Avp{Code: uint32(1423), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// ActiveAPN AVP
func ActiveAPN(id uint32) Avp {
	a := Avp{Code: uint32(1612), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []Avp
	t = append(t, ContextIdentifire(id))
	/*
		<!--<avp name="Service-Selection"	value=""></avp>-->
		<!--<avp name="MIP6-Agent-Info"	value=""></avp>-->
		<!--<avp name="Visited-Network-Identifier"	value=""></avp>-->
		<!--<avp name="Specific-APN-Info"	value=""></avp>-->
	*/
	a.Encode(t)
	return a
}

// EquivalentPLMNList AVP
func EquivalentPLMNList(plmns [][2]string) Avp {
	a := Avp{Code: uint32(1637), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	var t []Avp
	for _, p := range plmns {
		t = append(t, VisitedPLMNID(p[0], p[1]))
	}
	a.Encode(t)
	return a
}

// MMENumberForMTSMS AVP
func MMENumberForMTSMS(b []byte) Avp {
	a := Avp{Code: uint32(1645), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(b)
	return a
}

// SMSRegisterRequest AVP
func SMSRegisterRequest(e Enumerated) Avp {
	a := Avp{Code: uint32(1648), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(e)
	return a
}

const (
	SMSRegisterRequest_SMS_REGISTRATION_REQUIRED      Enumerated = 0
	SMSRegisterRequest_SMS_REGISTRATION_NOT_PREFERRED Enumerated = 1
	SMSRegisterRequest_NO_PREFERENCE                  Enumerated = 2
)

// SGsMMEIdentity AVP
func SGsMMEIdentity(s string) Avp {
	a := Avp{Code: uint32(1664), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(s)
	return a
}

// CoupledNodeDiameterID AVP
func CoupledNodeDiameterID(i DiameterIdentity) Avp {
	a := Avp{Code: uint32(1666), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// RequestedEUTRANAuthenticationInfo AVP
func RequestedEUTRANAuthenticationInfo(num uint32, resyncInfo []byte, immRespPref bool) Avp {
	a := Avp{Code: uint32(1408), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []Avp
	// Number-Of-Requested-Vectors
	if num != 0 {
		v := Avp{Code: uint32(1410), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(num)
		t = append(t, v)
	}
	// Re-Synchronization-Info
	if resyncInfo != nil {
		v := Avp{Code: uint32(1411), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(resyncInfo)
		t = append(t, v)
	}
	// Immediate-Response-Preferred
	if immRespPref {
		v := Avp{Code: uint32(1412), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.Encode(uint32(1))
		t = append(t, v)
	}
	a.Encode(t)
	return a
}

// RequestedUTRANGERANAuthenticationInfo AVP
func RequestedUTRANGERANAuthenticationInfo(num uint32, resyncInfo []byte, immRespPref bool) Avp {
	a := RequestedEUTRANAuthenticationInfo(num, resyncInfo, immRespPref)
	a.Code = uint32(1409)
	return a
}

// NORFlags AVP
func NORFlags(singleReg, sgsnRestrict, readySmSgsn, ueReachableMme, reserved, ueReachableSgsn, readySmMme, homogeneousSupportIMSVoPSSession, s6as6d, remMmeRegSm bool) Avp {
	a := Avp{Code: uint32(1443), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
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
