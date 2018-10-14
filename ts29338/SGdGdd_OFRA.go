package ts29338

import (
	"bytes"
	"fmt"

	dia "github.com/fkgi/diameter"
	"github.com/fkgi/sms"
	"github.com/fkgi/teldata"
)

/*
OFR is MO-Forward-ShortMessage-Request message.
 <OFR> ::= < Diameter Header: 8388645, REQ, PXY, 16777313 >
           < Session-Id >
           [ DRMP ]  // not supported
           [ Vendor-Specific-Application-Id ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ Destination-Host ]
           { Destination-Realm }
           { SC-Address }
           [ OFR-Flags ]
         * [ Supported-Features ]  // not supported
           { User-Identifier }
           { SM-RP-UI }
           [ SMSMI-Correlation-ID ]  // not supported
           [ SM-Delivery-Outcome ]  // not supported
         * [ AVP ]
         * [ Proxy-Info ]  // not supported
         * [ Route-Record ]
*/
type OFR struct {
	OriginHost       dia.Identity
	OriginRealm      dia.Identity
	DestinationHost  dia.Identity
	DestinationRealm dia.Identity

	SCAddress teldata.E164
	Flags     struct {
		S6aS6d bool
	}
	MSISDN teldata.E164
	teldata.IMSI
	SMSPDU sms.Submit

	// SupportedFeatures
	// SMSMICorrelationID
	// SMDeliveryOutcome
	// Proxy-Info
}

func (v OFR) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sOrigin-Host       =%s\n", dia.Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm      =%s\n", dia.Indent, v.OriginRealm)
	if len(v.DestinationHost) != 0 {
		fmt.Fprintf(w, "%sDestination-Host  =%s\n", dia.Indent, v.DestinationHost)
	} else {
		fmt.Fprintf(w, "%sDestination-Host  =not present\n", dia.Indent)
	}
	fmt.Fprintf(w, "%sDestination-Realm =%s\n", dia.Indent, v.DestinationRealm)

	fmt.Fprintf(w, "%sSC Address        =%s\n", dia.Indent, v.SCAddress)
	fmt.Fprintf(w, "%sSent from Gdd IF  =%t\n", dia.Indent, v.Flags.S6aS6d)
	if v.MSISDN.Length() == 0 && v.IMSI.Length() != 0 {
		fmt.Fprintf(w, "%sIMSI              =%s\n", dia.Indent, v.IMSI)
	} else {
		fmt.Fprintf(w, "%sMSISDN            =%s\n", dia.Indent, v.MSISDN)
	}
	fmt.Fprintf(w, "%sSMS Data Unit     =%s\n", dia.Indent, v.SMSPDU.String())
	return w.String()
}

// ToRaw return dia.RawMsg struct of this value
func (v OFR) ToRaw(s string) dia.RawMsg {
	m := dia.RawMsg{
		Ver:  dia.DiaVer,
		FlgR: true, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388645, AppID: 16777313,
		AVP: make([]dia.RawAVP, 0, 12)}

	m.AVP = append(m.AVP, dia.SetSessionID(s))
	m.AVP = append(m.AVP, dia.SetVendorSpecAppID(10415, m.AppID))
	m.AVP = append(m.AVP, dia.SetAuthSessionState(false))

	m.AVP = append(m.AVP, dia.SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, dia.SetOriginRealm(v.OriginRealm))
	if len(v.DestinationHost) != 0 {
		m.AVP = append(m.AVP, dia.SetDestinationHost(v.DestinationHost))
	}
	m.AVP = append(m.AVP, dia.SetDestinationRealm(v.DestinationRealm))

	m.AVP = append(m.AVP, setSCAddress(v.SCAddress))

	if v.Flags.S6aS6d {
		m.AVP = append(m.AVP, setOFRFlags(v.Flags.S6aS6d))
	}
	m.AVP = append(m.AVP, setUserIdentifier(v.IMSI, v.MSISDN))
	m.AVP = append(m.AVP, setSMRPUI(&v.SMSPDU))

	m.AVP = append(m.AVP, dia.SetRouteRecord(v.OriginHost))
	return m
}

// FromRaw make this value from dia.RawMsg struct
func (OFR) FromRaw(m dia.RawMsg) (dia.Request, string, error) {
	s := ""
	e := m.Validate(true, true, false, false)
	if e != nil {
		return nil, s, e
	}

	v := OFR{}
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

		case 3300:
			v.SCAddress, e = getSCAddress(a)
		case 3328:
			v.Flags.S6aS6d, e = getOFRFlags(a)
		case 3102:
			v.IMSI, v.MSISDN, e = getUserIdentifier(a)
		case 3301:
			v.SMSPDU, e = getSMRPUIasSubmit(a)
		}

		if e != nil {
			return nil, s, e
		}
	}

	if len(v.OriginHost) == 0 || len(v.OriginRealm) == 0 ||
		len(v.DestinationRealm) == 0 || v.SCAddress.Length() == 0 ||
		v.SMSPDU.DA.Addr == nil {
		e = dia.InvalidAVP(dia.DiameterMissingAvp)
	}

	if v.IMSI.Length() == 0 && v.MSISDN.Length() == 0 {
		e = dia.InvalidAVP(dia.DiameterMissingAvp)
	}
	return v, s, e
}

// Failed make error message for timeout
func (v OFR) Failed(c uint32) dia.Answer {
	return OFA{
		ResultCode:  c,
		OriginHost:  dia.Host,
		OriginRealm: dia.Realm}
}

/*
OFA is MO-Forward-Short-Message-Answer message.
 <OFA> ::= < Diameter Header: 8388645, PXY, 16777313 >
           < Session-Id >
           [ DRMP ] // not supported
           [ Vendor-Specific-Application-Id ]
           [ Result-Code ]
           [ Experimental-Result ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
         * [ Supported-Features ] // not supported
           [ SM-Delivery-Failure-Cause ]
		   [ SM-RP-UI ]
		   [ External-Identifier ] // not supported
         * [ AVP ]
         * [ Failed-AVP ]
         * [ Proxy-Info ] // not supported
         * [ Route-Record ]
*/
type OFA struct {
	ResultCode  uint32
	OriginHost  dia.Identity
	OriginRealm dia.Identity

	SMSPDU sms.SubmitReport
	DeliveryFailureCause

	FailedAVP []dia.RawAVP

	// SupportedFeatures
	// ExternalIdentifier
	// Proxy-Info
}

func (v OFA) String() string {
	w := new(bytes.Buffer)

	if v.ResultCode > 10000 {
		fmt.Fprintf(w, "%sExp-Result-Code   =%d:%d\n", dia.Indent, v.ResultCode/10000, v.ResultCode%10000)
	} else {
		fmt.Fprintf(w, "%sResult-Code       =%d\n", dia.Indent, v.ResultCode)
	}
	fmt.Fprintf(w, "%sOrigin-Host       =%s\n", dia.Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm      =%s\n", dia.Indent, v.OriginRealm)

	if v.ResultCode == dia.DiameterSuccess {
		fmt.Fprintf(w, "%sSMS Data Unit     =%s\n", dia.Indent, v.SMSPDU.String())
	}

	if v.ResultCode == DiameterErrorSmDeliveryFailure {
		switch v.DeliveryFailureCause {
		case CauseMemoryCapacityExceeded:
			fmt.Fprintf(w, "%sFailure Cause     =MEMORY_CAPACITY_EXCEEDED\n", dia.Indent)
		case CauseEquipmentProtocolError:
			fmt.Fprintf(w, "%sFailure Cause     =EQUIPMENT_PROTOCOL_ERROR\n", dia.Indent)
			fmt.Fprintf(w, "%sSMS Data Unit     =%s\n", dia.Indent, v.SMSPDU.String())
		case CauseEquipmentNotSMEquipped:
			fmt.Fprintf(w, "%sFailure Cause     =EQUIPMENT_NOT_SM-EQUIPPED\n", dia.Indent)
		case CauseUnknownServiceCenter:
			fmt.Fprintf(w, "%sFailure Cause     =UNKNOWN_SERVICE_CENTRE\n", dia.Indent)
		case CauseSCCongestion:
			fmt.Fprintf(w, "%sFailure Cause     =SC-CONGESTION\n", dia.Indent)
		case CauseInvalidSMEAddress:
			fmt.Fprintf(w, "%sFailure Cause     =INVALID_SME-ADDRESS\n", dia.Indent)
		case CauseUserNotSCUser:
			fmt.Fprintf(w, "%sFailure Cause     =USER_NOT_SC-USER\n", dia.Indent)
		}
	}
	return w.String()
}

// ToRaw return dia.RawMsg struct of this value
func (v OFA) ToRaw(s string) dia.RawMsg {
	m := dia.RawMsg{
		Ver:  dia.DiaVer,
		FlgR: false, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388645, AppID: 16777313,
		AVP: make([]dia.RawAVP, 0, 20)}

	m.AVP = append(m.AVP, dia.SetResultCode(v.ResultCode))
	m.AVP = append(m.AVP, dia.SetSessionID(s))
	m.AVP = append(m.AVP, dia.SetVendorSpecAppID(10415, m.AppID))

	m.AVP = append(m.AVP, dia.SetAuthSessionState(false))
	m.AVP = append(m.AVP, dia.SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, dia.SetOriginRealm(v.OriginRealm))

	switch v.ResultCode {
	case dia.DiameterSuccess:
		m.AVP = append(m.AVP, setSMRPUI(&v.SMSPDU))
	case DiameterErrorSmDeliveryFailure:
		m.AVP = append(m.AVP, setSMSubmissionFailureCause(
			v.DeliveryFailureCause, v.SMSPDU))
	}
	return m
}

// FromRaw make this value from dia.RawMsg struct
func (OFA) FromRaw(m dia.RawMsg) (dia.Answer, string, error) {
	s := ""
	e := m.Validate(false, true, false, false)
	if e != nil {
		return nil, s, e
	}

	v := OFA{}
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
		}
		if e != nil {
			return nil, s, e
		}
	}
	switch v.ResultCode {
	case dia.DiameterSuccess:
		for _, a := range m.AVP {
			switch a.Code {
			case 3301:
				v.SMSPDU, e = getSMRPUIasSubmitReport(a)
			}
			if e != nil {
				return nil, s, e
			}
		}
	case DiameterErrorSmDeliveryFailure:
		for _, a := range m.AVP {
			switch a.Code {
			case 3303:
				v.DeliveryFailureCause, v.SMSPDU, e = getSMSubmissionFailureCause(a)
			}
			if e != nil {
				return nil, s, e
			}
		}
	}
	if v.ResultCode == 0 ||
		len(v.OriginHost) == 0 || len(v.OriginRealm) == 0 {
		e = dia.InvalidAVP(dia.DiameterMissingAvp)
	}
	return v, s, e
}

// Result returns result-code
func (v OFA) Result() uint32 {
	return v.ResultCode
}
