package ts29338

import (
	"bytes"
	"fmt"

	dia "github.com/fkgi/diameter"
	"github.com/fkgi/sms"
	"github.com/fkgi/teldata"
)

/*
SRR is Send-Routing-Info-For-SM-Request message.
 <SRR> ::= < Diameter Header: 8388647, REQ, PXY, 16777312 >
           < Session-Id >
		   [ DRMP ]  // not supported
           [ Vendor-Specific-Application-Id ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ Destination-Host ]
           { Destination-Realm }
           [ MSISDN ]
           [ User-Name ]
           [ SMSMI-Correlation-ID ] // not supported
         * [ Supported-Features ] // not supported
           { SC-Address } // mandatory in table
           [ SM-RP-MTI ]
           [ SM-RP-SMEA ]
           [ SRR-Flags ]
           [ SM-Delivery-Not-Intended ]
         * [ AVP ]
         * [ Proxy-Info ] // not supported
		 * [ Route-Record ]
IP-SM-GW and MSISDN-less SMS are not supported.
*/
type SRR struct {
	// DRMP
	OriginHost       dia.Identity
	OriginRealm      dia.Identity
	DestinationHost  dia.Identity
	DestinationRealm dia.Identity

	MSISDN teldata.E164
	teldata.IMSI
	SCAddress teldata.E164

	MTType
	SMEAddr sms.Address
	Flags   struct {
		GPRSSupport   bool
		Prioritized   bool
		SingleAttempt bool
	}
	RequiredInfo

	// SMSMICorrelationID
	// []SupportedFeatures
	// []ProxyInfo
}

func (v SRR) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sOrigin-Host       =%s\n", dia.Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm      =%s\n", dia.Indent, v.OriginRealm)
	if len(v.DestinationHost) != 0 {
		fmt.Fprintf(w, "%sDestination-Host  =%s\n", dia.Indent, v.DestinationHost)
	} else {
		fmt.Fprintf(w, "%sDestination-Host  =not present\n", dia.Indent)
	}
	fmt.Fprintf(w, "%sDestination-Realm =%s\n", dia.Indent, v.DestinationRealm)

	if v.MSISDN.Length() == 0 && v.IMSI.Length() != 0 {
		fmt.Fprintf(w, "%sIMSI              =%s\n", dia.Indent, v.IMSI)
	} else {
		fmt.Fprintf(w, "%sMSISDN            =%s\n", dia.Indent, v.MSISDN)
	}
	fmt.Fprintf(w, "%sSC Address        =%s\n", dia.Indent, v.SCAddress)

	switch v.MTType {
	case DeliverMT:
		fmt.Fprintf(w, "%sMT Message Type   =Deliver\n", dia.Indent)
	case StatusReportMT:
		fmt.Fprintf(w, "%sMT Message Type   =Status Report\n", dia.Indent)
	default:
		fmt.Fprintf(w, "%sMT Message Type   =Unknown\n", dia.Indent)
	}
	fmt.Fprintf(w, "%sSME Address       =%s\n", dia.Indent, v.SMEAddr)
	fmt.Fprintf(w, "%sGPRS support      =%t\n", dia.Indent, v.Flags.GPRSSupport)
	fmt.Fprintf(w, "%sPrioritized       =%t\n", dia.Indent, v.Flags.Prioritized)
	fmt.Fprintf(w, "%sSingle Attempt    =%t\n", dia.Indent, v.Flags.SingleAttempt)

	switch v.RequiredInfo {
	case OnlyImsiRequested:
		fmt.Fprintf(w, "%sRequired data     =IMSI only\n", dia.Indent)
	case OnlyMccMncRequested:
		fmt.Fprintf(w, "%sRequired data     =MCC and MNC only\n", dia.Indent)
	default:
		fmt.Fprintf(w, "%sRequired data     =complete data\n", dia.Indent)
	}
	return w.String()
}

// ToRaw return dia.RawMsg struct of this value
func (v SRR) ToRaw(s string) dia.RawMsg {
	m := dia.RawMsg{
		Ver:  dia.DiaVer,
		FlgR: true, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388647, AppID: 16777312,
		AVP: make([]dia.RawAVP, 0, 15)}

	m.AVP = append(m.AVP, dia.SetSessionID(s))
	m.AVP = append(m.AVP, dia.SetVendorSpecAppID(10415, m.AppID))
	m.AVP = append(m.AVP, dia.SetAuthSessionState(false))

	m.AVP = append(m.AVP, dia.SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, dia.SetOriginRealm(v.OriginRealm))
	if len(v.DestinationHost) != 0 {
		m.AVP = append(m.AVP, dia.SetDestinationHost(v.DestinationHost))
	}
	m.AVP = append(m.AVP, dia.SetDestinationRealm(v.DestinationRealm))

	if v.MSISDN.Length() != 0 {
		m.AVP = append(m.AVP, setMSISDN(v.MSISDN))
	} else if v.IMSI.Length() != 0 {
		m.AVP = append(m.AVP, setUserName(v.IMSI))
	} else {
		m.AVP = append(m.AVP, setMSISDN(v.MSISDN))
	}
	m.AVP = append(m.AVP, setSCAddress(v.SCAddress))

	if v.MTType != UnknownMT {
		m.AVP = append(m.AVP, setSMRPMTI(v.MTType))
	}
	if v.SMEAddr.Addr != nil {
		m.AVP = append(m.AVP, setSMRPSMEA(v.SMEAddr))
	}
	if v.Flags.GPRSSupport || v.Flags.SingleAttempt || v.Flags.Prioritized {
		m.AVP = append(m.AVP, setSRRFlags(
			v.Flags.GPRSSupport, v.Flags.SingleAttempt, v.Flags.Prioritized))
	}
	if v.RequiredInfo != LocationRequested {
		m.AVP = append(m.AVP, setSMDeliveryNotIntended(v.RequiredInfo))
	}

	m.AVP = append(m.AVP, dia.SetRouteRecord(v.OriginHost))
	return m
}

// FromRaw make this value from dia.RawMsg struct
func (SRR) FromRaw(m dia.RawMsg) (dia.Request, string, error) {
	s := ""
	e := m.Validate(true, true, false, false)
	if e != nil {
		return nil, s, e
	}

	v := SRR{}
	for _, a := range m.AVP {
		switch a.Code {
		case 263:
			s, e = dia.GetSessionID(a)
		case 264:
			v.OriginHost, e = dia.GetOriginHost(a)
		case 296:
			v.OriginRealm, e = dia.GetOriginRealm(a)
		case 293:
			v.DestinationHost, e = dia.GetDestinationHost(a)
		case 283:
			v.DestinationRealm, e = dia.GetDestinationRealm(a)

		case 701:
			v.MSISDN, e = getMSISDN(a)
		case 1:
			v.IMSI, e = getUserName(a)
		case 3300:
			v.SCAddress, e = getSCAddress(a)
		case 3308:
			v.MTType, e = getSMRPMTI(a)
		case 3309:
			v.SMEAddr, e = getSMRPSMEA(a)
		case 3310:
			v.Flags.GPRSSupport, v.Flags.SingleAttempt, v.Flags.Prioritized, e = getSRRFlags(a)
		case 3311:
			v.RequiredInfo, e = getSMDeliveryNotIntended(a)
		}

		if e != nil {
			return nil, s, e
		}
	}

	if len(v.OriginHost) == 0 || len(v.OriginRealm) == 0 ||
		len(v.DestinationRealm) == 0 || v.SCAddress.Length() == 0 {
		e = dia.InvalidAVP(dia.DiameterMissingAvp)
	} else if v.MSISDN.Length() == 0 && v.IMSI.Length() == 0 {
		e = dia.InvalidAVP(dia.DiameterMissingAvp)
	}
	return v, s, e
}

// Failed make error message for timeout
func (v SRR) Failed(c uint32) dia.Answer {
	return SRA{
		ResultCode:  c,
		OriginHost:  dia.Host,
		OriginRealm: dia.Realm}
}

/*
SRA is SendRoutingInfoForSMAnswer message.
 <SRA> ::= < Diameter Header: 8388647, PXY, 16777312 >
           < Session-Id >
		   [ DRMP ] // not supported
           [ Vendor-Specific-Application-Id ]
           [ Result-Code ]
           [ Experimental-Result ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ User-Name ]
         * [ Supported-Features ] // not supported
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
         * [ Proxy-Info ] // not supported
         * [ Route-Record ]
IP-SM-GW and MSISDN-less SMS are not supported.
*/
type SRA struct {
	// DRMP
	ResultCode  uint32
	OriginHost  dia.Identity
	OriginRealm dia.Identity

	teldata.IMSI
	// ExtID  string
	ServingNode [2]struct {
		NodeType
		Address teldata.E164
		Host    dia.Identity
		Realm   dia.Identity
		LMSI    uint32
	}

	MWDStat struct { // for Inform-SC
		MSISDN   teldata.E164
		NoSCAddr bool
		MNRF     bool
		MCEF     bool
		MNRG     bool
	}
	AbsentUserDiag struct { // for Inform-SC
		MME  AbsentDiag
		MSC  AbsentDiag
		SGSN AbsentDiag
	}

	FailedAVP []dia.RawAVP
	// []SupportedFeatures
	// []ProxyInfo
}

func (v SRA) String() string {
	w := new(bytes.Buffer)

	if v.ResultCode > 10000 {
		fmt.Fprintf(w, "%sExp-Result-Code   =%d:%d\n", dia.Indent, v.ResultCode/10000, v.ResultCode%10000)
	} else {
		fmt.Fprintf(w, "%sResult-Code       =%d\n", dia.Indent, v.ResultCode)
	}
	fmt.Fprintf(w, "%sOrigin-Host       =%s\n", dia.Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm      =%s\n", dia.Indent, v.OriginRealm)

	fmt.Fprintf(w, "%sIMSI              =%s\n", dia.Indent, v.IMSI)
	for i, sn := range v.ServingNode {
		if sn.Address.Length() == 0 {
			continue
		}
		switch sn.NodeType {
		case NodeSGSN:
			fmt.Fprintf(w, "%sServing-Node#%d(SGSN)\n", dia.Indent, i+1)
		case NodeMME:
			fmt.Fprintf(w, "%sServing-Node#%d(MME)\n", dia.Indent, i+1)
		case NodeMSC:
			fmt.Fprintf(w, "%sServing-Node#%d(MSC)\n", dia.Indent, i+1)
		}
		fmt.Fprintf(w, "%s%sAddress =%s\n", dia.Indent, dia.Indent, sn.Address)
		fmt.Fprintf(w, "%s%sHost    =%s\n", dia.Indent, dia.Indent, sn.Host)
		fmt.Fprintf(w, "%s%sRealm   =%s\n", dia.Indent, dia.Indent, sn.Realm)
		fmt.Fprintf(w, "%s%sLMSI    =%x\n", dia.Indent, dia.Indent, sn.LMSI)
	}

	if v.MWDStat.MSISDN != nil || v.MWDStat.NoSCAddr || v.MWDStat.MNRF || v.MWDStat.MCEF || v.MWDStat.MNRG {
		fmt.Fprintf(w, "%sMWD information for Inform-SC\n", dia.Indent)
		fmt.Fprintf(w, "%s%sMSISDN in MWD     =%s\n", dia.Indent, dia.Indent, v.MWDStat.MSISDN)
		fmt.Fprintf(w, "%s%sno SCAddr in MWD  =%t\n", dia.Indent, dia.Indent, v.MWDStat.NoSCAddr)
		fmt.Fprintf(w, "%s%sMNRF              =%t\n", dia.Indent, dia.Indent, v.MWDStat.MNRF)
		fmt.Fprintf(w, "%s%sMCEF              =%t\n", dia.Indent, dia.Indent, v.MWDStat.MCEF)
		fmt.Fprintf(w, "%s%sMNRG              =%t\n", dia.Indent, dia.Indent, v.MWDStat.MNRG)
	}

	if v.AbsentUserDiag.MME != NoAbsentDiag || v.AbsentUserDiag.MSC != NoAbsentDiag || v.AbsentUserDiag.SGSN != NoAbsentDiag {
		fmt.Fprintf(w, "%sAbsent User Diagnostics for SM\n", dia.Indent)
		fmt.Fprintf(w, "%s%sMME  =%d\n", dia.Indent, dia.Indent, v.AbsentUserDiag.MME)
		fmt.Fprintf(w, "%s%sMSC  =%d\n", dia.Indent, dia.Indent, v.AbsentUserDiag.MSC)
		fmt.Fprintf(w, "%s%sSGSN =%d\n", dia.Indent, dia.Indent, v.AbsentUserDiag.SGSN)
	}
	return w.String()
}

// ToRaw return dia.RawMsg struct of this value
func (v SRA) ToRaw(s string) dia.RawMsg {
	m := dia.RawMsg{
		Ver:  dia.DiaVer,
		FlgR: false, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388647, AppID: 16777312,
		AVP: make([]dia.RawAVP, 0, 20)}

	m.AVP = append(m.AVP, dia.SetResultCode(v.ResultCode))
	m.AVP = append(m.AVP, dia.SetSessionID(s))
	m.AVP = append(m.AVP, dia.SetVendorSpecAppID(10415, m.AppID))
	m.AVP = append(m.AVP, dia.SetAuthSessionState(false))
	m.AVP = append(m.AVP, dia.SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, dia.SetOriginRealm(v.OriginRealm))

	if v.ResultCode != dia.DiameterSuccess {
		if len(v.FailedAVP) != 0 {
			m.AVP = append(m.AVP, dia.SetFailedAVP(v.FailedAVP))
		}
		return m
	}

	m.AVP = append(m.AVP, setUserName(v.IMSI))

	if v.ServingNode[0].Address.Length() != 0 {
		m.AVP = append(m.AVP, setServingNode(
			v.ServingNode[0].NodeType,
			v.ServingNode[0].Address,
			v.ServingNode[0].Host,
			v.ServingNode[0].Realm))
		if v.ServingNode[0].NodeType == NodeMSC && v.ServingNode[0].LMSI != 0 {
			m.AVP = append(m.AVP, setLMSI(v.ServingNode[0].LMSI))
		}

		if v.ServingNode[1].Address.Length() != 0 &&
			v.ServingNode[0].NodeType != v.ServingNode[1].NodeType {
			m.AVP = append(m.AVP, setAdditionalServingNode(
				v.ServingNode[1].NodeType,
				v.ServingNode[1].Address,
				v.ServingNode[1].Host,
				v.ServingNode[1].Realm))
			if v.ServingNode[1].NodeType == NodeMSC && v.ServingNode[1].LMSI != 0 {
				m.AVP = append(m.AVP, setLMSI(v.ServingNode[1].LMSI))
			}
		}
	}
	if v.MWDStat.MSISDN.Length() != 0 {
		m.AVP = append(m.AVP, setUserIdentifier("", v.MWDStat.MSISDN))
	}
	if v.MWDStat.NoSCAddr || v.MWDStat.MNRF || v.MWDStat.MCEF || v.MWDStat.MNRG {
		m.AVP = append(m.AVP, setMWDStatus(
			v.MWDStat.NoSCAddr, v.MWDStat.MNRF, v.MWDStat.MCEF, v.MWDStat.MNRG))
	}
	if v.AbsentUserDiag.MME != NoAbsentDiag {
		m.AVP = append(m.AVP, setMMEAbsentUserDiagnosticSM(v.AbsentUserDiag.MME))
	}
	if v.AbsentUserDiag.MSC != NoAbsentDiag {
		m.AVP = append(m.AVP, setMSCAbsentUserDiagnosticSM(v.AbsentUserDiag.MSC))
	}
	if v.AbsentUserDiag.SGSN != NoAbsentDiag {
		m.AVP = append(m.AVP, setSGSNAbsentUserDiagnosticSM(v.AbsentUserDiag.SGSN))
	}
	return m
}

// FromRaw make this value from dia.RawMsg struct
func (SRA) FromRaw(m dia.RawMsg) (dia.Answer, string, error) {
	s := ""
	e := m.Validate(false, true, false, false)
	if e != nil {
		return nil, s, e
	}

	v := SRA{}
	lmsi := uint32(0)
	for _, a := range m.AVP {
		switch a.Code {
		case 263:
			s, e = dia.GetSessionID(a)
		case 268, 297:
			v.ResultCode, e = dia.GetResultCode(a)
		case 264:
			v.OriginHost, e = dia.GetOriginHost(a)
		case 296:
			v.OriginRealm, e = dia.GetOriginRealm(a)

		case 1:
			v.IMSI, e = getUserName(a)
		case 3102:
			_, v.MWDStat.MSISDN, e = getUserIdentifier(a)
		case 2401:
			v.ServingNode[0].NodeType, v.ServingNode[0].Address,
				v.ServingNode[0].Host, v.ServingNode[0].Realm, e = getServingNode(a)
		case 2406:
			v.ServingNode[1].NodeType, v.ServingNode[1].Address,
				v.ServingNode[1].Host, v.ServingNode[1].Realm, e = getAdditionalServingNode(a)
		case 2400:
			lmsi, e = getLMSI(a)
		case 3312:
			v.MWDStat.NoSCAddr, v.MWDStat.MNRF, v.MWDStat.MCEF, v.MWDStat.MNRG, e = getMWDStatus(a)

		case 3313:
			v.AbsentUserDiag.MME, e = getMMEAbsentUserDiagnosticSM(a)
		case 3314:
			v.AbsentUserDiag.MSC, e = getMSCAbsentUserDiagnosticSM(a)
		case 3315:
			v.AbsentUserDiag.SGSN, e = getSGSNAbsentUserDiagnosticSM(a)
		}
		if e != nil {
			return nil, s, e
		}
	}

	if lmsi != 0 {
		if v.ServingNode[0].NodeType == NodeMSC {
			v.ServingNode[0].LMSI = lmsi
		}
		if v.ServingNode[1].NodeType == NodeMSC {
			v.ServingNode[1].LMSI = lmsi
		}
	}

	if v.ResultCode == 0 ||
		len(v.OriginHost) == 0 || len(v.OriginRealm) == 0 {
		e = dia.InvalidAVP(dia.DiameterMissingAvp)
	}
	if v.ResultCode == dia.DiameterSuccess {
		if v.IMSI.Length() == 0 {
			e = dia.InvalidAVP(dia.DiameterMissingAvp)
		} else if v.ServingNode[0].Address.Length() == 0 {
			e = dia.InvalidAVP(dia.DiameterMissingAvp)
		} else if v.ServingNode[1].Address.Length() != 0 &&
			v.ServingNode[0].NodeType == v.ServingNode[1].NodeType {
			e = dia.InvalidAVP(dia.DiameterInvalidAvpValue)
		}
	}
	return v, s, e
}

// Result returns result-code
func (v SRA) Result() uint32 {
	return v.ResultCode
}
