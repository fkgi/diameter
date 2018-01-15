package ts29338

import (
	"bytes"
	"encoding/binary"
	"time"

	dia "github.com/fkgi/diameter"
	"github.com/fkgi/sms"
	"github.com/fkgi/teldata"
)

func setUserName(v teldata.IMSI) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 1, VenID: 0, FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v.String())
	return
}

func getUserName(a dia.RawAVP) (v teldata.IMSI, e error) {
	s := new(string)
	if a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil {
		v, e = teldata.ParseIMSI(*s)
	}
	return
}

func setMSISDN(v teldata.E164) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 701, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	a.Encode(v.Bytes())
	return
}

func getMSISDN(a dia.RawAVP) (v teldata.E164, e error) {
	s := new([]byte)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil {
		v, e = teldata.B2E164(*s)
	}
	return
}

// MTType indicate SM-RP-MTI
// SM-RP-MTI AVP contain the RP-Message Type Indicator of the Short Message.
type MTType int

const (
	// UnknownMT is no SM-RP-MTI
	UnknownMT MTType = iota
	// DeliverMT is SM_DELIVER
	DeliverMT
	// StatusReportMT is SM_STATUS_REPORT
	StatusReportMT
)

func setSMRPMTI(v MTType) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3308, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	switch v {
	case DeliverMT:
		a.Encode(dia.Enumerated(0))
	case StatusReportMT:
		a.Encode(dia.Enumerated(1))
	}
	return
}

func getSMRPMTI(a dia.RawAVP) (v MTType, e error) {
	s := new(dia.Enumerated)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e != nil {
	} else if *s == 0 {
		v = DeliverMT
	} else if *s == 1 {
		v = StatusReportMT
	} else {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpValue)
	}
	return
}

// SM-RP-SMEA AVP contain the RP-Originating SME-address of
// the Short Message Entity that has originated the SM.
// It shall be formatted according to the formatting rules of
// the address fields described in 3GPP TS 23.040.
func setSMRPSMEA(v sms.Address) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3309, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	_, b := v.Encode()
	a.Encode(b)
	return
}

func getSMRPSMEA(a dia.RawAVP) (v sms.Address, e error) {
	s := new([]byte)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil {
		v.Decode(byte(len(*s)-3)*2, *s)
	}
	return
}

// SRR-Flags AVP contain a bit mask.
// GPRS-Indicator shall be ture if the SMS-GMSC supports receiving
// of two serving nodes addresses from the HSS.
// SM-RP-PRI shall be true if the delivery of the short message shall
// be attempted when a service centre address is already contained
// in the Message Waiting Data file.
// Single-Attempt if true indicates that only one delivery attempt
// shall be performed for this particular SM.
func setSRRFlags(g, p, s bool) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3310, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	i := uint32(0)
	if g {
		i = i | 0x00000001
	}
	if p {
		i = i | 0x00000002
	}
	if s {
		i = i | 0x00000004
	}
	a.Encode(i)
	return
}

func getSRRFlags(a dia.RawAVP) (g, p, s bool, e error) {
	v := new(uint32)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(v); e == nil {
		g = (*v)&0x00000001 == 0x00000001
		p = (*v)&0x00000002 == 0x00000002
		s = (*v)&0x00000004 == 0x00000004
	}
	return
}

// RequiredInfo indicate SM-Delivery-Not-Intended AVP data
// SM-Delivery-Not-Intended AVP indicate by its presence
// that delivery of a short message is not intended.
type RequiredInfo int

const (
	// LocationRequested is no SM-Delivery-Not-Intended
	LocationRequested RequiredInfo = iota
	// OnlyImsiRequested is ONLY_IMSI_REQUESTED
	OnlyImsiRequested
	// OnlyMccMncRequested is ONLY_MCC_MNC_REQUESTED
	OnlyMccMncRequested
)

func setSMDeliveryNotIntended(v RequiredInfo) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3311, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	switch v {
	case OnlyImsiRequested:
		a.Encode(dia.Enumerated(0))
	case OnlyMccMncRequested:
		a.Encode(dia.Enumerated(1))
	}
	return
}

func getSMDeliveryNotIntended(a dia.RawAVP) (v RequiredInfo, e error) {
	s := new(dia.Enumerated)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e != nil {
	} else if *s == 0 {
		v = OnlyImsiRequested
	} else if *s == 1 {
		v = OnlyMccMncRequested
	} else {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpValue)
	}
	return
}

// Serving-Node AVP shall contain the information about
// the network node serving the targeted SMS user.

// NodeType is tyoe of serving node address
type NodeType int

const (
	// NodeSGSN present SGSN
	NodeSGSN NodeType = iota
	// NodeMME present MME
	NodeMME
	// NodeMSC present MSC
	NodeMSC
	// NodeIPSMGW present IP-SM-GW
	// NodeIPSMGW
)

func setServingNode(t NodeType, d teldata.E164, n, r dia.Identity) dia.RawAVP {
	a := dia.RawAVP{Code: 2401, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	return setSN(a, t, d, n, r)
}

func setAdditionalServingNode(t NodeType, d teldata.E164, n, r dia.Identity) dia.RawAVP {
	a := dia.RawAVP{Code: 2406, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	return setSN(a, t, d, n, r)
}

func setSN(a dia.RawAVP, t NodeType, d teldata.E164, n, r dia.Identity) dia.RawAVP {
	v := make([]dia.RawAVP, 1, 3)
	switch t {
	case NodeSGSN:
		v[0] = dia.RawAVP{Code: 1489, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
		v[0].Encode(d.Bytes())
		if len(n) != 0 {
			nv := dia.RawAVP{Code: 2409, VenID: 10415, FlgV: true, FlgM: false, FlgP: false}
			nv.Encode(n)
			v = append(v, nv)
		}
		if len(r) != 0 {
			rv := dia.RawAVP{Code: 2410, VenID: 10415, FlgV: true, FlgM: false, FlgP: false}
			rv.Encode(r)
			v = append(v, rv)
		}
	case NodeMME:
		v[0] = dia.RawAVP{Code: 1645, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
		v[0].Encode(d.Bytes())
		if len(n) != 0 {
			nv := dia.RawAVP{Code: 2402, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
			nv.Encode(n)
			v = append(v, nv)
		}
		if len(r) != 0 {
			rv := dia.RawAVP{Code: 2408, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
			rv.Encode(r)
			v = append(v, rv)
		}
	case NodeMSC:
		v[0] = dia.RawAVP{Code: 2403, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
		v[0].Encode(d.Bytes())
		if len(n) != 0 {
			nv := dia.RawAVP{Code: 2402, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
			nv.Encode(n)
			v = append(v, nv)
		}
		if len(r) != 0 {
			rv := dia.RawAVP{Code: 2408, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
			rv.Encode(r)
			v = append(v, rv)
		}
	}

	a.Encode(v)
	return a
}

func getServingNode(a dia.RawAVP) (t NodeType, d teldata.E164, n, r dia.Identity, e error) {
	return getSN(a)
}

func getAdditionalServingNode(a dia.RawAVP) (t NodeType, d teldata.E164, n, r dia.Identity, e error) {
	return getSN(a)
}

func getSN(a dia.RawAVP) (t NodeType, d teldata.E164, n, r dia.Identity, e error) {
	o := []dia.RawAVP{}
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else {
		e = a.Decode(&o)
	}
	for _, a := range o {
		b := new([]byte)
		switch a.Code {
		case 1489:
			if !a.FlgV || !a.FlgM || a.FlgP {
				e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
			} else if e = a.Decode(b); e == nil {
				d, e = teldata.B2E164(*b)
				t = NodeSGSN
			}
		case 1645:
			if !a.FlgV || !a.FlgM || a.FlgP {
				e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
			} else if e = a.Decode(b); e == nil {
				d, e = teldata.B2E164(*b)
				t = NodeMME
			}
		case 2403:
			if !a.FlgV || !a.FlgM || a.FlgP {
				e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
			} else if e = a.Decode(b); e == nil {
				d, e = teldata.B2E164(*b)
				t = NodeMSC
			}
		}
	}
	switch t {
	case NodeSGSN:
		for _, a := range o {
			switch a.Code {
			case 2409:
				if !a.FlgV || !a.FlgM || a.FlgP {
					e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
				} else {
					e = a.Decode(&n)
				}
			case 2410:
				if !a.FlgV || !a.FlgM || a.FlgP {
					e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
				} else {
					e = a.Decode(&r)
				}
			}
		}
	case NodeMME, NodeMSC:
		for _, a := range o {
			switch a.Code {
			case 2402:
				if !a.FlgV || !a.FlgM || a.FlgP {
					e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
				} else {
					e = a.Decode(&n)
				}
			case 2408:
				if !a.FlgV || !a.FlgM || a.FlgP {
					e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
				} else {
					e = a.Decode(&r)
				}
			}
		}
	}
	return
}

func setLMSI(v uint32) dia.RawAVP {
	a := dia.RawAVP{Code: 2400, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, v)
	a.Encode(buf.Bytes())
	return a
}

func getLMSI(a dia.RawAVP) (v uint32, e error) {
	s := new([]byte)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil && len(*s) == 4 {
		binary.Read(bytes.NewBuffer(*s), binary.BigEndian, &v)
	}
	return
}

func setUserIdentifier(i teldata.IMSI, m teldata.E164) dia.RawAVP {
	a := dia.RawAVP{Code: 3102, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	t := make([]dia.RawAVP, 1)
	if m.Length() != 0 {
		t[0] = setMSISDN(m)
	} else if i.Length() != 0 {
		t[0] = setUserName(i)
	} else {
		t[0] = setMSISDN([]byte{})
	}
	a.Encode(t)
	return a
}

func getUserIdentifier(a dia.RawAVP) (i teldata.IMSI, m teldata.E164, e error) {
	o := []dia.RawAVP{}
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(&o); e == nil {
		for _, a := range o {
			switch a.Code {
			case 1:
				i, e = getUserName(a)
			case 701:
				m, e = getMSISDN(a)
			}
			if e != nil {
				return
			}
		}
	}
	return
}

// MWD-Status AVP contain a bit mask.
// SCAddrNotIncluded shall indicate the presence of the SC Address in the Message Waiting Data in the HSS.
// MNRF shall indicate that the MNRF flag is set in the HSS.
// MCEF shall indicate that the MCEF flag is set in the HSS.
// MNRG shall indicate that the MNRG flag is set in the HSS.
func setMWDStatus(sca, mnrf, mcef, mnrg bool) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3312, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}

	i := uint32(0)
	if sca {
		i = i | 0x00000001
	}
	if mnrf {
		i = i | 0x00000002
	}
	if mcef {
		i = i | 0x00000004
	}
	if mnrg {
		i = i | 0x00000008
	}
	a.Encode(i)
	return
}

func getMWDStatus(a dia.RawAVP) (sca, mnrf, mcef, mnrg bool, e error) {
	s := new(uint32)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil {
		sca = (*s)&0x00000001 == 0x00000001
		mnrf = (*s)&0x00000002 == 0x00000002
		mcef = (*s)&0x00000004 == 0x00000004
		mnrg = (*s)&0x00000008 == 0x00000008
	}
	return
}

// AbsentDiag indicate Absent-User-Diagnostic-SM
type AbsentDiag int

const (
	// NoAbsentDiag is "no diag data"
	NoAbsentDiag AbsentDiag = iota
	// NoPagingRespMSC is "no paging response via the MSC"
	NoPagingRespMSC
	// IMSIDetached is "IMSI detached"
	IMSIDetached
	// RoamingRestrict is "roaming restriction"
	RoamingRestrict
	// DeregisteredNonGPRS is "deregistered in the HLR for non GPRS"
	DeregisteredNonGPRS
	// PurgedNonGPRS is "MS purged for non GPRS"
	PurgedNonGPRS
	// NoPagingRespSGSN is "no paging response via the SGSN"
	NoPagingRespSGSN
	// GPRSDetached is "GPRS detached"
	GPRSDetached
	// DeregisteredGPRS is "deregistered in the HLR for GPRS"
	DeregisteredGPRS
	// PurgedGPRS is "MS purged for GPRS"
	PurgedGPRS
	// UnidentifiedSubsMSC is "Unidentified subscriber via the MSC"
	UnidentifiedSubsMSC
	// UnidentifiedSubsSGSN is "Unidentified subscriber via the SGSN"
	UnidentifiedSubsSGSN
	// DeregisteredIMS is "deregistered in the HSS/HLR for IMS"
	DeregisteredIMS
	// NoRespIPSMGW is "no response via the IP-SM-GW"
	NoRespIPSMGW
	// TempUnavailable is "the MS is temporarily unavailable"
	TempUnavailable
)

func (a AbsentDiag) String() string {
	switch a {
	case NoAbsentDiag:
		return "no diag data"
	case NoPagingRespMSC:
		return "no paging response via the MSC"
	case IMSIDetached:
		return "IMSI detached"
	case RoamingRestrict:
		return "roaming restriction"
	case DeregisteredNonGPRS:
		return "deregistered in the HLR for non GPRS"
	case PurgedNonGPRS:
		return "MS purged for non GPRS"
	case NoPagingRespSGSN:
		return "no paging response via the SGSN"
	case GPRSDetached:
		return "GPRS detached"
	case DeregisteredGPRS:
		return "deregistered in the HLR for GPRS"
	case PurgedGPRS:
		return "MS purged for GPRS"
	case UnidentifiedSubsMSC:
		return "Unidentified subscriber via the MSC"
	case UnidentifiedSubsSGSN:
		return "Unidentified subscriber via the SGSN"
	case DeregisteredIMS:
		return "deregistered in the HSS/HLR for IMS"
	case NoRespIPSMGW:
		return "no response via the IP-SM-GW"
	case TempUnavailable:
		return "the MS is temporarily unavailable"
	}
	return ""
}

// MMEAbsentUserDiagnosticSM AVP shall indicate the diagnostic
// explaining the absence of the user given by the MME.
func setMMEAbsentUserDiagnosticSM(v AbsentDiag) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3313, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	switch v {
	case NoPagingRespMSC:
		a.Encode(uint32(0))
	case IMSIDetached:
		a.Encode(uint32(1))
	case RoamingRestrict:
		a.Encode(uint32(2))
	case DeregisteredNonGPRS:
		a.Encode(uint32(3))
	case PurgedNonGPRS:
		a.Encode(uint32(4))
	case UnidentifiedSubsMSC:
		a.Encode(uint32(9))
	case TempUnavailable:
		a.Encode(uint32(13))
	}
	return
}

func getMMEAbsentUserDiagnosticSM(a dia.RawAVP) (v AbsentDiag, e error) {
	s := new(uint32)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil {
		switch *s {
		case 0:
			v = NoPagingRespMSC
		case 1:
			v = IMSIDetached
		case 2:
			v = RoamingRestrict
		case 3:
			v = DeregisteredNonGPRS
		case 4:
			v = PurgedNonGPRS
		case 9:
			v = UnidentifiedSubsMSC
		case 13:
			v = TempUnavailable
		default:
			e = dia.InvalidAVP(dia.DiameterInvalidAvpValue)
		}
	}
	return
}

// MSCAbsentUserDiagnosticSM AVP shall indicate the diagnostic
// explaining the absence of the user given by the MSC.
func setMSCAbsentUserDiagnosticSM(v AbsentDiag) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3314, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	switch v {
	case NoPagingRespMSC:
		a.Encode(uint32(0))
	case IMSIDetached:
		a.Encode(uint32(1))
	case RoamingRestrict:
		a.Encode(uint32(2))
	case DeregisteredNonGPRS:
		a.Encode(uint32(3))
	case PurgedNonGPRS:
		a.Encode(uint32(4))
	case UnidentifiedSubsMSC:
		a.Encode(uint32(9))
	case TempUnavailable:
		a.Encode(uint32(13))
	}
	return
}

func getMSCAbsentUserDiagnosticSM(a dia.RawAVP) (v AbsentDiag, e error) {
	s := new(uint32)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil {
		switch *s {
		case 0:
			v = NoPagingRespMSC
		case 1:
			v = IMSIDetached
		case 2:
			v = RoamingRestrict
		case 3:
			v = DeregisteredNonGPRS
		case 4:
			v = PurgedNonGPRS
		case 9:
			v = UnidentifiedSubsMSC
		case 13:
			v = TempUnavailable
		default:
			e = dia.InvalidAVP(dia.DiameterInvalidAvpValue)
		}
	}
	return
}

// SGSNAbsentUserDiagnosticSM AVP shall indicate the diagnostic
// explaining the absence of the user given by the SGSN.
func setSGSNAbsentUserDiagnosticSM(v AbsentDiag) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3315, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	switch v {
	case IMSIDetached:
		a.Encode(uint32(1))
	case RoamingRestrict:
		a.Encode(uint32(2))
	case NoPagingRespSGSN:
		a.Encode(uint32(5))
	case GPRSDetached:
		a.Encode(uint32(6))
	case DeregisteredGPRS:
		a.Encode(uint32(7))
	case PurgedGPRS:
		a.Encode(uint32(8))
	case UnidentifiedSubsSGSN:
		a.Encode(uint32(10))
	case TempUnavailable:
		a.Encode(uint32(13))
	}
	return
}

func getSGSNAbsentUserDiagnosticSM(a dia.RawAVP) (v AbsentDiag, e error) {
	s := new(uint32)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil {
		switch *s {
		case 1:
			v = IMSIDetached
		case 2:
			v = RoamingRestrict
		case 5:
			v = NoPagingRespSGSN
		case 6:
			v = GPRSDetached
		case 7:
			v = DeregisteredGPRS
		case 8:
			v = PurgedGPRS
		case 10:
			v = UnidentifiedSubsSGSN
		case 13:
			v = TempUnavailable
		default:
			e = dia.InvalidAVP(dia.DiameterInvalidAvpValue)
		}
	}
	return
}

// AbsentUserDiagnosticSM AVP shall indicate the diagnostic explaining the absence of the subscriber.
func setAbsentUserDiagnosticSM(v AbsentDiag) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3322, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	switch v {
	case NoPagingRespMSC:
		a.Encode(uint32(0))
	case IMSIDetached:
		a.Encode(uint32(1))
	case RoamingRestrict:
		a.Encode(uint32(2))
	case DeregisteredNonGPRS:
		a.Encode(uint32(3))
	case PurgedNonGPRS:
		a.Encode(uint32(4))
	case NoPagingRespSGSN:
		a.Encode(uint32(5))
	case GPRSDetached:
		a.Encode(uint32(6))
	case DeregisteredGPRS:
		a.Encode(uint32(7))
	case PurgedGPRS:
		a.Encode(uint32(8))
	case UnidentifiedSubsMSC:
		a.Encode(uint32(9))
	case UnidentifiedSubsSGSN:
		a.Encode(uint32(10))
	case DeregisteredIMS:
		a.Encode(uint32(11))
	case NoRespIPSMGW:
		a.Encode(uint32(12))
	case TempUnavailable:
		a.Encode(uint32(13))
	}
	return
}

func getAbsentUserDiagnosticSM(a dia.RawAVP) (v AbsentDiag, e error) {
	s := new(uint32)
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil {
		switch *s {
		case 0:
			v = NoPagingRespMSC
		case 1:
			v = IMSIDetached
		case 2:
			v = RoamingRestrict
		case 3:
			v = DeregisteredNonGPRS
		case 4:
			v = PurgedNonGPRS
		case 5:
			v = NoPagingRespSGSN
		case 6:
			v = GPRSDetached
		case 7:
			v = DeregisteredGPRS
		case 8:
			v = PurgedGPRS
		case 9:
			v = UnidentifiedSubsMSC
		case 10:
			v = UnidentifiedSubsSGSN
		case 11:
			v = DeregisteredIMS
		case 12:
			v = NoRespIPSMGW
		case 13:
			v = TempUnavailable
		default:
			e = dia.InvalidAVP(dia.DiameterInvalidAvpValue)
		}
	}
	return
}

// SMDeliveryCause AVP value
type SMDeliveryCause int

const (
	// NoOutcome indicate DeliveryOutcome AVP is not present
	NoOutcome SMDeliveryCause = iota
	// UeMemoryCapacityExceeded is Enumerated value 0
	UeMemoryCapacityExceeded
	// AbsentUser is Enumerated value 1
	AbsentUser
	// SuccessfulTransfer is Enumerated value 2
	SuccessfulTransfer
)

func setSMDeliveryOutcome(ec, cc, nc SMDeliveryCause, ed, cd, nd AbsentDiag) dia.RawAVP {
	a := dia.RawAVP{Code: 3316, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	t := make([]dia.RawAVP, 0, 2)
	if ec != NoOutcome {
		t = append(t, setNodeOutcome(3317, ec, ed))
	} else if cc != NoOutcome {
		t = append(t, setNodeOutcome(3318, cc, cd))
	}
	if nc != NoOutcome {
		t = append(t, setNodeOutcome(3319, nc, nd))
	}
	if len(t) == 0 {
		t = append(t, setNodeOutcome(3317, SuccessfulTransfer, NoAbsentDiag))
	}
	return a
}

func setNodeOutcome(code uint32, cause SMDeliveryCause, diag AbsentDiag) dia.RawAVP {
	a := dia.RawAVP{Code: code, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	t := make([]dia.RawAVP, 1, 2)
	t[0] = dia.RawAVP{Code: 3321, VenID: 10415, FlgV: true, FlgM: true, FlgP: false}
	switch cause {
	case UeMemoryCapacityExceeded:
		t[0].Encode(dia.Enumerated(0))
	case AbsentUser:
		t[0].Encode(dia.Enumerated(1))
		t = append(t, dia.RawAVP{Code: 3322, VenID: 10415, FlgV: true, FlgM: true, FlgP: false})
		t[1].Encode(uint32(diag))
	case SuccessfulTransfer:
		t[0].Encode(dia.Enumerated(2))
	}
	a.Encode(t)
	return a
}

func getSMDeliveryOutcome(a dia.RawAVP) (
	ec, cc, nc SMDeliveryCause, ed, cd, nd AbsentDiag, e error) {
	o := []dia.RawAVP{}
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(&o); e == nil {
		for _, a := range o {
			switch a.Code {
			case 3317:
				ec, ed, e = getNodeOutcome(a)
			case 3318:
				cc, cd, e = getNodeOutcome(a)
			case 3319:
				nc, nd, e = getNodeOutcome(a)
			}
			if e != nil {
				return
			}
		}
		if ec != NoOutcome && cc != NoOutcome {
			e = dia.InvalidAVP(dia.DiameterInvalidAvpValue)
		}
	}
	return
}

func getNodeOutcome(a dia.RawAVP) (c SMDeliveryCause, d AbsentDiag, e error) {
	o := []dia.RawAVP{}
	if !a.FlgV || !a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(&o); e == nil {
		for _, a := range o {
			switch a.Code {
			case 3321:
				var n dia.Enumerated
				if !a.FlgV || !a.FlgM || a.FlgP {
					e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
				} else if e = a.Decode(&n); e != nil {
				} else if n == 0 {
					c = UeMemoryCapacityExceeded
				} else if n == 1 {
					c = AbsentUser
				} else if n == 2 {
					c = SuccessfulTransfer
				}
			case 3322:
				var n uint32
				if !a.FlgV || !a.FlgM || a.FlgP {
					e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
				} else if e = a.Decode(&n); e == nil {
					d = AbsentDiag(n)
				}
			}
			if e != nil {
				return
			}
		}
	}
	return
}

// RDR-Flags AVP contain a bit mask.
// Single-Attempt-Delivery indicates that only one delivery attempt
// shall be performed for this particular SM.
func setRDRFlags(s bool) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3323, VenID: 10415, FlgV: true, FlgM: false, FlgP: false}
	i := uint32(0)
	if s {
		i = i | 0x00000001
	}
	a.Encode(i)
	return
}

func getRDRFlags(a dia.RawAVP) (s bool, e error) {
	v := new(uint32)
	if !a.FlgV || a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(v); e == nil {
		s = (*v)&0x00000001 == 0x00000001
	}
	return
}

// Maximum-UE-Availability-Time AVP shall contain the timestamp (in UTC)
// until which a UE using a power saving mechanism
// (such as extended idle mode DRX) is expected to be reachable for SM Delivery.
func setMaximumUEAvailabilityTime(v time.Time) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3329, VenID: 10415, FlgV: true, FlgM: false, FlgP: false}
	a.Encode(v)
	return a
}

func getMaximumUEAvailabilityTime(a dia.RawAVP) (v time.Time, e error) {
	if !a.FlgV || a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	}
	e = a.Decode(&v)
	return
}

// SM-SGMSC-Alert-Event AVP shall contain a bit mask.
// UE-Available-For-MT-SMS shall indicate that the UE is now available for MT SMS
// UE-Under-New-Serving-Node shall indicate that the UE has moved
// under the coverage of another MME or SGSN.
func setSMSGMSCAlertEvent(av, nn bool) (a dia.RawAVP) {
	a = dia.RawAVP{Code: 3333, VenID: 10415, FlgV: true, FlgM: false, FlgP: false}
	i := uint32(0)
	if av {
		i = i | 0x00000001
	}
	if nn {
		i = i | 0x00000002
	}
	a.Encode(i)
	return a
}

func getSMSGMSCAlertEvent(a dia.RawAVP) (av, nn bool, e error) {
	s := new(uint32)
	if !a.FlgV || a.FlgM || a.FlgP {
		e = dia.InvalidAVP(dia.DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e == nil {
		av = (*s)&0x00000001 == 0x00000001
		nn = (*s)&0x00000002 == 0x00000002
	}
	return
}
