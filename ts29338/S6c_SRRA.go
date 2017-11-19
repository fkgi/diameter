package ts29338

import (
	"bytes"
	"fmt"

	"github.com/fkgi/diameter"
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
           [ SC-Address ]
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
	OriginHost       diameter.Identity
	OriginRealm      diameter.Identity
	DestinationHost  diameter.Identity
	DestinationRealm diameter.Identity

	MSISDN teldata.E164
	teldata.IMSI
	// SMSMICorrelationID
	// []SupportedFeatures
	SCAddress teldata.E164
	MessageType
	SMRPSMEA sms.Address
	Flags    struct {
		GPRSIndicator bool
		SMRPPRI       bool
		SingleAttempt bool
	}
	RequiredInfo
	// []ProxyInfo
}

func (v SRR) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sOrigin-Host       =%s\n", diameter.Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm      =%s\n", diameter.Indent, v.OriginRealm)
	fmt.Fprintf(w, "%sDestination-Host  =%s\n", diameter.Indent, v.DestinationHost)
	fmt.Fprintf(w, "%sDestination-Realm =%s\n", diameter.Indent, v.DestinationRealm)

	fmt.Fprintf(w, "%sMSISDN            =%s\n", diameter.Indent, v.MSISDN)
	fmt.Fprintf(w, "%sIMSI              =%s\n", diameter.Indent, v.IMSI)
	fmt.Fprintf(w, "%sSCAddress         =%s\n", diameter.Indent, v.SCAddress)
	switch v.MessageType {
	case DeliverType:
		fmt.Fprintf(w, "%sSM-RP-MTI         =SM_DELIVER\n", diameter.Indent)
	case StatusReportType:
		fmt.Fprintf(w, "%sSM-RP-MTI         =SM_STATUS_REPORT\n", diameter.Indent)
	default:
		fmt.Fprintf(w, "%sSM-RP-MTI         =unknown\n", diameter.Indent)
	}
	fmt.Fprintf(w, "%sSM-RP-SMEA        =%s\n", diameter.Indent, v.SMRPSMEA)
	fmt.Fprintf(w, "%sGPRS support      =%t\n", diameter.Indent, v.Flags.GPRSIndicator)
	fmt.Fprintf(w, "%sSM-RP-PRI         =%t\n", diameter.Indent, v.Flags.SMRPPRI)
	fmt.Fprintf(w, "%sSingle-Attempt    =%t\n", diameter.Indent, v.Flags.SingleAttempt)

	switch v.RequiredInfo {
	case OnlyImsiRequested:
		fmt.Fprintf(w, "%sSM-Delivery-Not-Intended =ONLY_IMSI_REQUESTED\n", diameter.Indent)
	case OnlyMccMncRequested:
		fmt.Fprintf(w, "%sSM-Delivery-Not-Intended =ONLY_MCC_MNC_REQUESTED\n", diameter.Indent)
	}
	return w.String()
}

// ToRaw return diameter.RawMsg struct of this value
func (v SRR) ToRaw(s string) diameter.RawMsg {
	m := diameter.RawMsg{
		Ver:  diameter.DiaVer,
		FlgR: true, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388647, AppID: 16777312,
		AVP: make([]diameter.RawAVP, 0, 15)}

	m.AVP = append(m.AVP, setSessionID(s))
	m.AVP = append(m.AVP, setVendorSpecAppID(16777312))
	m.AVP = append(m.AVP, setAuthSessionState())

	m.AVP = append(m.AVP, setOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, setOriginRealm(v.OriginRealm))
	if len(v.DestinationHost) != 0 {
		m.AVP = append(m.AVP, setDestinationHost(v.DestinationHost))
	}
	m.AVP = append(m.AVP, setDestinationRealm(v.DestinationRealm))

	if v.MSISDN.Length() != 0 {
		m.AVP = append(m.AVP, setMSISDN(v.MSISDN))
	}
	if v.IMSI.Length() != 0 {
		m.AVP = append(m.AVP, setUserName(v.IMSI))
	}
	if v.SCAddress.Length() != 0 {
		m.AVP = append(m.AVP, setSCAddress(v.SCAddress))
	}
	if v.MessageType != UnknownType {
		m.AVP = append(m.AVP, setSMRPMTI(v.MessageType))
	}
	if v.SMRPSMEA.Addr != nil {
		m.AVP = append(m.AVP, setSMRPSMEA(v.SMRPSMEA))
	}
	if v.Flags.GPRSIndicator ||
		v.Flags.SingleAttempt ||
		v.Flags.SMRPPRI {
		m.AVP = append(m.AVP, setSRRFlags(
			v.Flags.GPRSIndicator,
			v.Flags.SingleAttempt,
			v.Flags.SMRPPRI))
	}
	if v.RequiredInfo != LocationRequested {
		m.AVP = append(m.AVP, setSMDeliveryNotIntended(v.RequiredInfo))
	}

	m.AVP = append(m.AVP, setRouteRecord(v.OriginHost))
	return m
}

// FromRaw make this value from diameter.RawMsg struct
func (SRR) FromRaw(m diameter.RawMsg) (diameter.Request, string, error) {
	s := ""
	e := m.Validate(16777312, 8388647, true, true, false, false)
	if e != nil {
		return nil, s, e
	}

	v := SRR{}
	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 263 {
			s, e = getSessionID(a)
		} else if a.VenID == 0 && a.Code == 264 {
			v.OriginHost, e = getOriginHost(a)
		} else if a.VenID == 0 && a.Code == 296 {
			v.OriginRealm, e = getOriginRealm(a)
		} else if a.VenID == 0 && a.Code == 293 {
			v.DestinationHost, e = getDestinationHost(a)
		} else if a.VenID == 0 && a.Code == 283 {
			v.DestinationRealm, e = getDestinationRealm(a)

		} else if a.VenID == 10415 && a.Code == 701 {
			v.MSISDN, e = getMSISDN(a)
		} else if a.VenID == 0 && a.Code == 1 {
			v.IMSI, e = getUserName(a)
		} else if a.VenID == 10415 && a.Code == 3300 {
			v.SCAddress, e = getSCAddress(a)
		} else if a.VenID == 10415 && a.Code == 3308 {
			v.MessageType, e = getSMRPMTI(a)
		} else if a.VenID == 10415 && a.Code == 3309 {
			v.SMRPSMEA, e = getSMRPSMEA(a)
		} else if a.VenID == 10415 && a.Code == 3310 {
			v.Flags.GPRSIndicator,
				v.Flags.SingleAttempt,
				v.Flags.SMRPPRI, e = getSRRFlags(a)
		} else if a.VenID == 10415 && a.Code == 3311 {
			v.RequiredInfo, e = getSMDeliveryNotIntended(a)
		}

		if e != nil {
			return nil, s, e
		}
	}

	if len(v.OriginHost) == 0 ||
		len(v.OriginRealm) == 0 ||
		len(v.DestinationRealm) == 0 {
		e = diameter.NoMandatoryAVP{}
	}
	return v, s, e
}

// Failed make error message for timeout
func (v SRR) Failed(c uint32) diameter.Answer {
	return SRA{
		ResultCode:  c,
		OriginHost:  v.OriginHost,
		OriginRealm: v.OriginRealm}
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
	OriginHost  diameter.Identity
	OriginRealm diameter.Identity

	// []SupportedFeatures
	ServingNode struct {
		NodeType
		Address teldata.E164
		Name    diameter.Identity
		Realm   diameter.Identity
	}
	AdditionalServingNode struct {
		NodeType
		Address teldata.E164
		Name    diameter.Identity
		Realm   diameter.Identity
	}
	LMSI uint32
	User struct {
		teldata.IMSI
		MSISDN teldata.E164
		// ExtID  string
	}
	MWDStat struct {
		SCAddrNotIncluded bool
		MNRF              bool
		MCEF              bool
		MNRG              bool
	}

	MMEAbsentUserDiagnosticSM  uint32
	MSCAbsentUserDiagnosticSM  uint32
	SGSNAbsentUserDiagnosticSM uint32

	FailedAVP []diameter.RawAVP
	// []ProxyInfo
}

func (v SRA) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sResult-Code       =%d\n", diameter.Indent, v.ResultCode)
	fmt.Fprintf(w, "%sOrigin-Host       =%s\n", diameter.Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm      =%s\n", diameter.Indent, v.OriginRealm)

	switch v.ServingNode.NodeType {
	case NodeSGSN:
		fmt.Fprintf(w, "%sServing-Node(SGSN)\n", diameter.Indent)
	case NodeMME:
		fmt.Fprintf(w, "%sServing-Node(MME)\n", diameter.Indent)
	case NodeMSC:
		fmt.Fprintf(w, "%sServing-Node(MSC)\n", diameter.Indent)
	}
	fmt.Fprintf(w, "%s%sAddress =%s\n", diameter.Indent, diameter.Indent, v.ServingNode.Address)
	fmt.Fprintf(w, "%s%sHost    =%s\n", diameter.Indent, diameter.Indent, v.ServingNode.Name)
	fmt.Fprintf(w, "%s%sRealm   =%s\n", diameter.Indent, diameter.Indent, v.ServingNode.Realm)
	switch v.AdditionalServingNode.NodeType {
	case NodeSGSN:
		fmt.Fprintf(w, "%sServing-Node(SGSN)\n", diameter.Indent)
	case NodeMME:
		fmt.Fprintf(w, "%sServing-Node(MME)\n", diameter.Indent)
	case NodeMSC:
		fmt.Fprintf(w, "%sServing-Node(MSC)\n", diameter.Indent)
	}
	fmt.Fprintf(w, "%s%sAddress =%s\n", diameter.Indent, diameter.Indent, v.AdditionalServingNode.Address)
	fmt.Fprintf(w, "%s%sHost    =%s\n", diameter.Indent, diameter.Indent, v.AdditionalServingNode.Name)
	fmt.Fprintf(w, "%s%sRealm   =%s\n", diameter.Indent, diameter.Indent, v.AdditionalServingNode.Realm)
	fmt.Fprintf(w, "%sLMSI              =%x\n", diameter.Indent, v.LMSI)
	fmt.Fprintf(w, "%sIMSI              =%s\n", diameter.Indent, v.User.IMSI)
	fmt.Fprintf(w, "%sMSISDN            =%s\n", diameter.Indent, v.User.MSISDN)

	fmt.Fprintf(w, "%sSCAddrNotIncluded =%t\n", diameter.Indent, v.MWDStat.SCAddrNotIncluded)
	fmt.Fprintf(w, "%sMNRF              =%t\n", diameter.Indent, v.MWDStat.MNRF)
	fmt.Fprintf(w, "%sMCEF              =%t\n", diameter.Indent, v.MWDStat.MCEF)
	fmt.Fprintf(w, "%sMNRG              =%t\n", diameter.Indent, v.MWDStat.MNRG)

	return w.String()
}

// ToRaw return diameter.RawMsg struct of this value
func (v SRA) ToRaw(s string) diameter.RawMsg {
	m := diameter.RawMsg{
		Ver:  diameter.DiaVer,
		FlgR: false, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388647, AppID: 16777312,
		AVP: make([]diameter.RawAVP, 0, 20)}

	m.AVP = append(m.AVP, setSessionID(s))
	m.AVP = append(m.AVP, setVendorSpecAppID(16777312))
	if v.ResultCode >= 5000 && v.ResultCode <= 5999 {
		m.AVP = append(m.AVP, setExperimentalResult(10415, v.ResultCode))
	} else {
		m.AVP = append(m.AVP, setResultCode(v.ResultCode))
	}
	m.AVP = append(m.AVP, setAuthSessionState())
	m.AVP = append(m.AVP, setOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, setOriginRealm(v.OriginRealm))

	if v.User.IMSI.Length() != 0 {
		m.AVP = append(m.AVP, setUserName(v.User.IMSI))
	}
	if v.ServingNode.Address.Length() != 0 {
		m.AVP = append(m.AVP, setServingNode(
			v.ServingNode.NodeType,
			v.ServingNode.Address,
			v.ServingNode.Name,
			v.ServingNode.Realm))
	}
	if v.AdditionalServingNode.Address.Length() != 0 {
		m.AVP = append(m.AVP, setAdditionalServingNode(
			v.AdditionalServingNode.NodeType,
			v.AdditionalServingNode.Address,
			v.AdditionalServingNode.Name,
			v.AdditionalServingNode.Realm))
	}
	if v.LMSI != 0 {
		m.AVP = append(m.AVP, setLMSI(v.LMSI))
	}
	if v.User.MSISDN.Length() != 0 {
		m.AVP = append(m.AVP, setUserIdentifier(v.User.MSISDN))
	}
	if v.MWDStat.SCAddrNotIncluded ||
		v.MWDStat.MNRF ||
		v.MWDStat.MCEF ||
		v.MWDStat.MNRG {
		m.AVP = append(m.AVP, setMWDStatus(
			v.MWDStat.SCAddrNotIncluded,
			v.MWDStat.MNRF,
			v.MWDStat.MCEF,
			v.MWDStat.MNRG))
	}
	if v.MMEAbsentUserDiagnosticSM != 0 {
		m.AVP = append(m.AVP,
			setMMEAbsentUserDiagnosticSM(v.MMEAbsentUserDiagnosticSM))
	}
	if v.MSCAbsentUserDiagnosticSM != 0 {
		m.AVP = append(m.AVP,
			setMSCAbsentUserDiagnosticSM(v.MSCAbsentUserDiagnosticSM))
	}
	if v.SGSNAbsentUserDiagnosticSM != 0 {
		m.AVP = append(m.AVP,
			setSGSNAbsentUserDiagnosticSM(v.SGSNAbsentUserDiagnosticSM))
	}
	if len(v.FailedAVP) != 0 {
		m.AVP = append(m.AVP, setFailedAVP(v.FailedAVP))
	}
	return m
}

// FromRaw make this value from diameter.RawMsg struct
func (SRA) FromRaw(m diameter.RawMsg) (diameter.Answer, string, error) {
	s := ""
	e := m.Validate(16777312, 8388647, false, true, false, false)
	if e != nil {
		return nil, s, e
	}

	v := SRA{}
	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 263 {
			s, e = getSessionID(a)
		} else if a.VenID == 0 && a.Code == 268 {
			v.ResultCode, e = getResultCode(a)
		} else if a.VenID == 0 && a.Code == 297 {
			_, v.ResultCode, e = getExperimentalResult(a)
		} else if a.VenID == 0 && a.Code == 264 {
			v.OriginHost, e = getOriginHost(a)
		} else if a.VenID == 0 && a.Code == 296 {
			v.OriginRealm, e = getOriginRealm(a)

		} else if a.VenID == 0 && a.Code == 1 {
			v.User.IMSI, e = getUserName(a)
		} else if a.VenID == 10415 && a.Code == 3102 {
			v.User.MSISDN, e = getUserIdentifier(a)
		} else if a.VenID == 10415 && a.Code == 2401 {
			v.ServingNode.NodeType,
				v.ServingNode.Address,
				v.ServingNode.Name,
				v.ServingNode.Realm, e = getServingNode(a)
		} else if a.VenID == 10415 && a.Code == 2406 {
			v.AdditionalServingNode.NodeType,
				v.AdditionalServingNode.Address,
				v.AdditionalServingNode.Name,
				v.AdditionalServingNode.Realm, e = getAdditionalServingNode(a)
		} else if a.VenID == 10415 && a.Code == 2400 {
			v.LMSI, e = getLMSI(a)
		} else if a.VenID == 10415 && a.Code == 3312 {
			v.MWDStat.SCAddrNotIncluded,
				v.MWDStat.MNRF,
				v.MWDStat.MCEF,
				v.MWDStat.MNRG, e = getMWDStatus(a)

		} else if a.VenID == 10415 && a.Code == 3313 {
			v.MMEAbsentUserDiagnosticSM, e = getMMEAbsentUserDiagnosticSM(a)
		} else if a.VenID == 10415 && a.Code == 3314 {
			v.MSCAbsentUserDiagnosticSM, e = getMSCAbsentUserDiagnosticSM(a)
		} else if a.VenID == 10415 && a.Code == 3315 {
			v.SGSNAbsentUserDiagnosticSM, e = getSGSNAbsentUserDiagnosticSM(a)
		}
		if e != nil {
			return nil, s, e
		}
	}

	if len(v.OriginHost) == 0 ||
		len(v.OriginRealm) == 0 {
		e = diameter.NoMandatoryAVP{}
	}
	return v, s, e
}

// Result returns result-code
func (v SRA) Result() uint32 {
	return v.ResultCode
}

/*
AlertServiceCentreRequest is ALR message.
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
*/
/*
 AlertServiceCentreAnswer is ALA message.
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
ReportSMDeliveryStatusRequest is RDR message.
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
*/
/*
ReportSMDeliveryStatusAnswer is RDA message.
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
