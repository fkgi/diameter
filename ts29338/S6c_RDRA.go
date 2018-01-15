package ts29338

import (
	"bytes"
	"fmt"

	dia "github.com/fkgi/diameter"
	"github.com/fkgi/teldata"
)

/*
RDR is Report-SM-Delivery-Status-Request message.
 <RDR> ::= < Diameter Header: 8388649, REQ, PXY, 16777312 >
		   < Session-Id >
		   [ DRMP ] // not supported
		   [ Vendor-Specific-Application-Id ]
		   { Auth-Session-State }
		   { Origin-Host }
		   { Origin-Realm }
		   [ Destination-Host ]
		   { Destination-Realm }
		 * [ Supported-Features ] // not supported
		   { User-Identifier }
		   [ SMSMI-Correlation-ID ] // not supported
		   { SC-Address }
		   { SM-Delivery-Outcome }
		   [ RDR-Flags ]
		 * [ AVP ]
		 * [ Proxy-Info ] // not supported
		 * [ Route-Record ]
*/
type RDR struct {
	OriginHost       dia.Identity
	OriginRealm      dia.Identity
	DestinationHost  dia.Identity
	DestinationRealm dia.Identity

	MSISDN teldata.E164
	teldata.IMSI
	SCAddress teldata.E164

	DeliveryOutcome struct {
		MME struct {
			SMDeliveryCause
			AbsentDiag
		}
		MSC struct {
			SMDeliveryCause
			AbsentDiag
		}
		SGSN struct {
			SMDeliveryCause
			AbsentDiag
		}
	}
	Flags struct {
		SingleAttempt bool
	}

	// DRMP
	// SMSMICorrelationID
	// []SupportedFeatures
	// []ProxyInfo
}

func (v RDR) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sOrigin-Host       =%s\n", dia.Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm      =%s\n", dia.Indent, v.OriginRealm)
	fmt.Fprintf(w, "%sDestination-Host  =%s\n", dia.Indent, v.DestinationHost)
	fmt.Fprintf(w, "%sDestination-Realm =%s\n", dia.Indent, v.DestinationRealm)

	fmt.Fprintf(w, "%sMSISDN            =%s\n", dia.Indent, v.MSISDN)
	fmt.Fprintf(w, "%sIMSI              =%s\n", dia.Indent, v.IMSI)
	fmt.Fprintf(w, "%sSC Address        =%s\n", dia.Indent, v.SCAddress)

	switch v.DeliveryOutcome.MME.SMDeliveryCause {
	case UeMemoryCapacityExceeded:
		fmt.Fprintf(w, "%sMME Delivery Cause=UE memory capacity exceeded\n", dia.Indent)
	case AbsentUser:
		fmt.Fprintf(w, "%sMME Delivery Cause=Absent user\n", dia.Indent)
		fmt.Fprintf(w, "%s%sAbsent User Diagnostic =%s\n",
			dia.Indent, dia.Indent, v.DeliveryOutcome.MME.AbsentDiag)
	case SuccessfulTransfer:
		fmt.Fprintf(w, "%sMME Delivery Cause=Successful transferd\n", dia.Indent)
	}
	switch v.DeliveryOutcome.MSC.SMDeliveryCause {
	case UeMemoryCapacityExceeded:
		fmt.Fprintf(w, "%sMSC Delivery Cause=UE memory capacity exceeded\n", dia.Indent)
	case AbsentUser:
		fmt.Fprintf(w, "%sMSC Delivery Cause=Absent user\n", dia.Indent)
		fmt.Fprintf(w, "%s%sAbsent User Diagnostic =%s\n",
			dia.Indent, dia.Indent, v.DeliveryOutcome.MSC.AbsentDiag)
	case SuccessfulTransfer:
		fmt.Fprintf(w, "%sMSC Delivery Cause=Successful transferd\n", dia.Indent)
	}
	switch v.DeliveryOutcome.SGSN.SMDeliveryCause {
	case UeMemoryCapacityExceeded:
		fmt.Fprintf(w, "%sSGSN Delivery Cause=UE memory capacity exceeded\n", dia.Indent)
	case AbsentUser:
		fmt.Fprintf(w, "%sSGSN Delivery Cause=Absent user\n", dia.Indent)
		fmt.Fprintf(w, "%s%sAbsent User Diagnostic =%s\n",
			dia.Indent, dia.Indent, v.DeliveryOutcome.SGSN.AbsentDiag)
	case SuccessfulTransfer:
		fmt.Fprintf(w, "%sSGSN Delivery Cause=Successful transferd\n", dia.Indent)
	}

	fmt.Fprintf(w, "%sGPRS support      =%t\n", dia.Indent, v.Flags.SingleAttempt)

	return w.String()
}

// ToRaw return dia.RawMsg struct of this value
func (v RDR) ToRaw(s string) dia.RawMsg {
	m := dia.RawMsg{
		Ver:  dia.DiaVer,
		FlgR: true, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388649, AppID: 16777312,
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

	m.AVP = append(m.AVP, setUserIdentifier(v.IMSI, v.MSISDN))
	m.AVP = append(m.AVP, setSCAddress(v.SCAddress))

	m.AVP = append(m.AVP, setSMDeliveryOutcome(
		v.DeliveryOutcome.MME.SMDeliveryCause,
		v.DeliveryOutcome.MSC.SMDeliveryCause,
		v.DeliveryOutcome.SGSN.SMDeliveryCause,
		v.DeliveryOutcome.MME.AbsentDiag,
		v.DeliveryOutcome.MSC.AbsentDiag,
		v.DeliveryOutcome.SGSN.AbsentDiag))

	if v.Flags.SingleAttempt {
		m.AVP = append(m.AVP, setRDRFlags(v.Flags.SingleAttempt))
	}

	m.AVP = append(m.AVP, dia.SetRouteRecord(v.OriginHost))
	return m
}

// FromRaw make this value from dia.RawMsg struct
func (RDR) FromRaw(m dia.RawMsg) (dia.Request, string, error) {
	s := ""
	e := m.Validate(true, true, false, false)
	if e != nil {
		return nil, s, e
	}

	v := RDR{}
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

		case 3102:
			v.IMSI, v.MSISDN, e = getUserIdentifier(a)
		case 3300:
			v.SCAddress, e = getSCAddress(a)
		case 3316:
			v.DeliveryOutcome.MME.SMDeliveryCause,
				v.DeliveryOutcome.MSC.SMDeliveryCause,
				v.DeliveryOutcome.SGSN.SMDeliveryCause,
				v.DeliveryOutcome.MME.AbsentDiag,
				v.DeliveryOutcome.MSC.AbsentDiag,
				v.DeliveryOutcome.SGSN.AbsentDiag,
				e = getSMDeliveryOutcome(a)
		case 3323:
			v.Flags.SingleAttempt, e = getRDRFlags(a)
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
	} else if v.DeliveryOutcome.MME.SMDeliveryCause == NoOutcome &&
		v.DeliveryOutcome.MSC.SMDeliveryCause == NoOutcome &&
		v.DeliveryOutcome.SGSN.SMDeliveryCause == NoOutcome {
		e = dia.InvalidAVP(dia.DiameterMissingAvp)
	}
	return v, s, e
}

// Failed make error message for timeout
func (v RDR) Failed(c uint32) dia.Answer {
	return SRA{
		ResultCode:  c,
		OriginHost:  dia.Host,
		OriginRealm: dia.Realm}
}

/*
RDA is ReportSMDeliveryStatusAnswer message.
 <RDA> ::= < Diameter Header: 8388649, PXY, 16777312 >
		   < Session-Id >
		   [ DRMP ] // not supported
		   [ Vendor-Specific-Application-Id ]
		   [ Result-Code ]
		   [ Experimental-Result ]
		   { Auth-Session-State }
		   { Origin-Host }
		   { Origin-Realm }
		 * [ Supported-Features ] // not supported
		   [ User-Identifier ]
		 * [ AVP ]
		 * [ Failed-AVP ]
		 * [ Proxy-Info ] // not supported
		 * [ Route-Record ]
*/
type RDA struct {
	// DRMP
	ResultCode  uint32
	OriginHost  dia.Identity
	OriginRealm dia.Identity

	MSISDN teldata.E164

	FailedAVP []dia.RawAVP
}

func (v RDA) String() string {
	w := new(bytes.Buffer)

	if v.ResultCode > 10000 {
		fmt.Fprintf(w, "%sExp-Result-Code   =%d:%d\n", dia.Indent, v.ResultCode/10000, v.ResultCode%10000)
	} else {
		fmt.Fprintf(w, "%sResult-Code       =%d\n", dia.Indent, v.ResultCode)
	}
	fmt.Fprintf(w, "%sOrigin-Host       =%s\n", dia.Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm      =%s\n", dia.Indent, v.OriginRealm)

	if v.MSISDN.Length() != 0 {
		fmt.Fprintf(w, "%sMSISDN            =%s\n", dia.Indent, v.MSISDN)
	}

	return w.String()
}

// ToRaw return dia.RawMsg struct of this value
func (v RDA) ToRaw(s string) dia.RawMsg {
	m := dia.RawMsg{
		Ver:  dia.DiaVer,
		FlgR: false, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388649, AppID: 16777312,
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
	if v.MSISDN.Length() != 0 {
		m.AVP = append(m.AVP, setUserIdentifier("", v.MSISDN))
	}
	return m
}

// FromRaw make this value from dia.RawMsg struct
func (RDA) FromRaw(m dia.RawMsg) (dia.Answer, string, error) {
	s := ""
	e := m.Validate(false, true, false, false)
	if e != nil {
		return nil, s, e
	}

	v := RDA{}
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

		case 3102:
			_, v.MSISDN, e = getUserIdentifier(a)
		}
		if e != nil {
			return nil, s, e
		}
	}

	if v.ResultCode == 0 ||
		len(v.OriginHost) == 0 || len(v.OriginRealm) == 0 {
		e = dia.InvalidAVP(dia.DiameterMissingAvp)
	}
	return v, s, e
}

// Result returns result-code
func (v RDA) Result() uint32 {
	return v.ResultCode
}
