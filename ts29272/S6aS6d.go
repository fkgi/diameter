package ts29272

import (
	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/teldata"
)

/*
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
func ULRFlags(singleReg, s6as6d, skipSubsData, gprsSubsData, nodeType, initAttach, psLcsNotSupportedByUE, smsOnly bool) msg.RawAVP {
	a := msg.RawAVP{Code: uint32(1405), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
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

	a.ToRaw(i)
	return a
}

// VisitedPLMNID AVP
func VisitedPLMNID(mcc, mnc string) msg.RawAVP {
	a := msg.RawAVP{Code: uint32(1407), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
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

	a.ToRaw(b)
	return a
}

// TerminalInformation AVP
func TerminalInformation(imei string, meid []byte, version string) msg.RawAVP {
	a := msg.RawAVP{Code: uint32(1401), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}

	var t []msg.RawAVP
	// IMEI
	if len(imei) != 0 {
		v := msg.RawAVP{Code: uint32(1402), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.ToRaw(imei)
		t = append(t, v)
	}
	// 3GPP2-MEID
	if meid != nil && len(meid) != 0 {
		v := msg.RawAVP{Code: uint32(1471), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.ToRaw(meid)
		t = append(t, v)
	}
	// Sofrware-Version
	if len(version) != 0 {
		v := msg.RawAVP{Code: uint32(1403), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.ToRaw(version)
		t = append(t, v)
	}

	a.ToRaw(msg.GroupedAVP(t))
	return a
}

// UESRVCCCapability AVP
func UESRVCCCapability(e msg.Enumerated) msg.RawAVP {
	a := msg.RawAVP{Code: uint32(1615), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.ToRaw(e)
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

*/

// SGSNNumber AVP
type SGSNNumber teldata.TBCD

// ToRaw return AVP struct of this value
func (v *SGSNNumber) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 1489, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(teldata.TBCD(*v).String())
	}
	return a
}

// FromRaw get AVP value
func (v *SGSNNumber) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 1489, true, true, false); e != nil {
		return
	}
	s := new(string)
	if e = a.Decode(s); e != nil {
		return
	}
	if t, e := teldata.ParseTBCD(*s); e != nil {
		*v = SGSNNumber(t)
	}
	return
}

/*
// HomogeneousSupportOfIMSVoiceOverPSSessions AVP
func HomogeneousSupportOfIMSVoiceOverPSSessions(e msg.Enumerated) msg.RawAVP {
	a := msg.RawAVP{Code: uint32(1493), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.ToRaw(e)
	return a
}

const (
	// NotSupported is Enumerated value 0
	NotSupported msg.Enumerated = 0
	// Supported is Enumerated value 1
	Supported msg.Enumerated = 1
)

// ContextIdentifire AVP
func ContextIdentifire(i uint32) msg.RawAVP {
	a := msg.RawAVP{Code: uint32(1423), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.ToRaw(i)
	return a
}

// ActiveAPN AVP
func ActiveAPN(id uint32) msg.RawAVP {
	a := msg.RawAVP{Code: uint32(1612), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []msg.RawAVP
	t = append(t, ContextIdentifire(id))
	/*
		<!--<avp name="Service-Selection"	value=""></avp>-->
		<!--<avp name="MIP6-Agent-Info"	value=""></avp>-->
		<!--<avp name="Visited-Network-Identifier"	value=""></avp>-->
		<!--<avp name="Specific-APN-Info"	value=""></avp>-->
*/
/*
	a.ToRaw(msg.GroupedAVP(t))
	return a
}

// EquivalentPLMNList AVP
func EquivalentPLMNList(plmns [][2]string) msg.RawAVP {
	a := msg.RawAVP{Code: uint32(1637), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	var t []msg.RawAVP
	for _, p := range plmns {
		t = append(t, VisitedPLMNID(p[0], p[1]))
	}
	a.ToRaw(msg.GroupedAVP(t))
	return a
}
*/

// MMENumberForMTSMS AVP
type MMENumberForMTSMS teldata.TBCD

// ToRaw return AVP struct of this value
func (v *MMENumberForMTSMS) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 1645, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	if v != nil {
		a.Encode(teldata.TBCD(*v).String())
	}
	return a
}

// FromRaw get AVP value
func (v *MMENumberForMTSMS) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 1645, true, false, false); e != nil {
		return
	}
	s := new(string)
	if e = a.Decode(s); e != nil {
		return
	}
	if t, e := teldata.ParseTBCD(*s); e != nil {
		*v = MMENumberForMTSMS(t)
	}
	return
}

/*
const (
	// SmsRegistrationRequired is Enumerated value 0
	SmsRegistrationRequired msg.Enumerated = 0
	// SmsRegistrationNotPreferred is Enumerated value 1
	SmsRegistrationNotPreferred msg.Enumerated = 1
	// NotPreference is Enumerated value 2
	NotPreference msg.Enumerated = 2
)

// SGsMMEIdentity AVP
func SGsMMEIdentity(s string) msg.RawAVP {
	a := msg.RawAVP{Code: uint32(1664), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.ToRaw(s)
	return a
}

// CoupledNodeDiameterID AVP
func CoupledNodeDiameterID(i msg.DiameterIdentity) msg.RawAVP {
	a := msg.RawAVP{Code: uint32(1666), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.ToRaw(i)
	return a
}

// RequestedEUTRANAuthenticationInfo AVP
func RequestedEUTRANAuthenticationInfo(num uint32, resyncInfo []byte, immRespPref bool) msg.RawAVP {
	a := msg.RawAVP{Code: uint32(1408), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []msg.RawAVP
	// Number-Of-Requested-Vectors
	if num != 0 {
		v := msg.RawAVP{Code: uint32(1410), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.ToRaw(num)
		t = append(t, v)
	}
	// Re-Synchronization-Info
	if resyncInfo != nil {
		v := msg.RawAVP{Code: uint32(1411), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.ToRaw(resyncInfo)
		t = append(t, v)
	}
	// Immediate-Response-Preferred
	if immRespPref {
		v := msg.RawAVP{Code: uint32(1412), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
		v.ToRaw(uint32(1))
		t = append(t, v)
	}
	a.ToRaw(msg.GroupedAVP(t))
	return a
}

// RequestedUTRANGERANAuthenticationInfo AVP
func RequestedUTRANGERANAuthenticationInfo(num uint32, resyncInfo []byte, immRespPref bool) msg.RawAVP {
	a := RequestedEUTRANAuthenticationInfo(num, resyncInfo, immRespPref)
	a.Code = uint32(1409)
	return a
}

// NORFlags AVP
func NORFlags(singleReg, sgsnRestrict, readySmSgsn, ueReachableMme, reserved, ueReachableSgsn, readySmMme, homogeneousSupportIMSVoPSSession, s6as6d, remMmeRegSm bool) msg.RawAVP {
	a := msg.RawAVP{Code: uint32(1443), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
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

	a.ToRaw(i)
	return a
}
*/
