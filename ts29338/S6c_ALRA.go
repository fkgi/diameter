package ts29338

import (
	"bytes"
	"fmt"
	"time"

	dia "github.com/fkgi/diameter"
	"github.com/fkgi/teldata"
)

/*
ALR is AlertServiceCentreRequest message.
 <ALR> ::= < Diameter Header: 8388648, REQ, PXY, 16777312 >
		   < Session-Id >
		   [ DRMP ] // not supported
		   [ Vendor-Specific-Application-Id ]
		   { Auth-Session-State }
		   { Origin-Host }
		   { Origin-Realm }
		   [ Destination-Host ]
		   { Destination-Realm }
		   { SC-Address }
		   { User-Identifier }
		   [ SMSMI-Correlation-ID ] // not supported
		   [ Maximum-UE-Availability-Time ]
		   [ SMS-GMSC-Alert-Event ]
		   [ Serving-Node ]
		 * [ Supported-Features ] // not supported
		 * [ AVP ]
		 * [ Proxy-Info ] // not supported
		 * [ Route-Record ]
*/
type ALR struct {
	OriginHost       dia.Identity
	OriginRealm      dia.Identity
	DestinationHost  dia.Identity
	DestinationRealm dia.Identity

	MSISDN teldata.E164
	teldata.IMSI
	SCAddress teldata.E164

	AvailTime time.Time
	Flags     struct {
		AvailForMT   bool
		UnderNewNode bool
	}
}

func (v ALR) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sOrigin-Host       =%s\n", dia.Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm      =%s\n", dia.Indent, v.OriginRealm)
	fmt.Fprintf(w, "%sDestination-Host  =%s\n", dia.Indent, v.DestinationHost)
	fmt.Fprintf(w, "%sDestination-Realm =%s\n", dia.Indent, v.DestinationRealm)

	fmt.Fprintf(w, "%sMSISDN            =%s\n", dia.Indent, v.MSISDN)
	fmt.Fprintf(w, "%sIMSI              =%s\n", dia.Indent, v.IMSI)
	fmt.Fprintf(w, "%sSC Address        =%s\n", dia.Indent, v.SCAddress)
	fmt.Fprintf(w, "%sMax UE Avail Time =%s\n", dia.Indent, v.AvailTime)

	fmt.Fprintf(w, "%sAvail for MT-SMS  =%t\n", dia.Indent, v.Flags.AvailForMT)
	fmt.Fprintf(w, "%sUnder new SrvNode =%t\n", dia.Indent, v.Flags.UnderNewNode)

	return w.String()
}

// ToRaw return dia.RawMsg struct of this value
func (v ALR) ToRaw(s string) dia.RawMsg {
	m := dia.RawMsg{
		Ver:  dia.DiaVer,
		FlgR: true, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388648, AppID: 16777312,
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

	m.AVP = append(m.AVP, setSCAddress(v.SCAddress))
	m.AVP = append(m.AVP, setUserIdentifier(v.IMSI, v.MSISDN))
	if !v.AvailTime.IsZero() {
		m.AVP = append(m.AVP, setMaximumUEAvailabilityTime(v.AvailTime))
	}
	if v.Flags.AvailForMT || v.Flags.UnderNewNode {
		m.AVP = append(m.AVP, setSMSGMSCAlertEvent(
			v.Flags.AvailForMT, v.Flags.UnderNewNode))
	}

	m.AVP = append(m.AVP, dia.SetRouteRecord(v.OriginHost))
	return m
}

// FromRaw make this value from dia.RawMsg struct
func (ALR) FromRaw(m dia.RawMsg) (dia.Request, string, error) {
	s := ""
	e := m.Validate(true, true, false, false)
	if e != nil {
		return nil, s, e
	}

	v := ALR{}
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
		case 3329:
			v.AvailTime, e = getMaximumUEAvailabilityTime(a)
		case 3333:
			v.Flags.AvailForMT, v.Flags.UnderNewNode, e = getSMSGMSCAlertEvent(a)
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
func (v ALR) Failed(c uint32) dia.Answer {
	return ALA{
		ResultCode:  c,
		OriginHost:  dia.Host,
		OriginRealm: dia.Realm}
}

/*
ALA is AlertServiceCentreAnswer message.
 <ALA> ::= < Diameter Header: 8388648, PXY, 16777312 >
		   < Session-Id >
		   [ DRMP ] // not supported
		   [ Vendor-Specific-Application-Id ]
		   [ Result-Code ]
		   [ Experimental-Result ]
		   { Auth-Session-State }
		   { Origin-Host }
		   { Origin-Realm }
		 * [ Supported-Features ] // not supported
		 * [ AVP ]
		 * [ Failed-AVP ]
		 * [ Proxy-Info ] // not supported
		 * [ Route-Record ]
*/
type ALA struct {
	ResultCode  uint32
	OriginHost  dia.Identity
	OriginRealm dia.Identity

	FailedAVP []dia.RawAVP
}

func (v ALA) String() string {
	w := new(bytes.Buffer)

	if v.ResultCode > 10000 {
		fmt.Fprintf(w, "%sExp-Result-Code   =%d:%d\n", dia.Indent, v.ResultCode/10000, v.ResultCode%10000)
	} else {
		fmt.Fprintf(w, "%sResult-Code       =%d\n", dia.Indent, v.ResultCode)
	}
	fmt.Fprintf(w, "%sOrigin-Host       =%s\n", dia.Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm      =%s\n", dia.Indent, v.OriginRealm)

	return w.String()
}

// ToRaw return dia.RawMsg struct of this value
func (v ALA) ToRaw(s string) dia.RawMsg {
	m := dia.RawMsg{
		Ver:  dia.DiaVer,
		FlgR: false, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388648, AppID: 16777312,
		AVP: make([]dia.RawAVP, 0, 10)}

	m.AVP = append(m.AVP, dia.SetResultCode(v.ResultCode))
	m.AVP = append(m.AVP, dia.SetSessionID(s))
	m.AVP = append(m.AVP, dia.SetVendorSpecAppID(10415, m.AppID))
	m.AVP = append(m.AVP, dia.SetAuthSessionState(false))
	m.AVP = append(m.AVP, dia.SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, dia.SetOriginRealm(v.OriginRealm))

	if v.ResultCode != dia.DiameterSuccess && len(v.FailedAVP) != 0 {
		m.AVP = append(m.AVP, dia.SetFailedAVP(v.FailedAVP))
	}
	return m
}

// FromRaw make this value from dia.RawMsg struct
func (ALA) FromRaw(m dia.RawMsg) (dia.Answer, string, error) {
	s := ""
	e := m.Validate(false, true, false, false)
	if e != nil {
		return nil, s, e
	}

	v := ALA{}
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

	if v.ResultCode == 0 ||
		len(v.OriginHost) == 0 || len(v.OriginRealm) == 0 {
		e = dia.InvalidAVP(dia.DiameterMissingAvp)
	}
	return v, s, e
}

// Result returns result-code
func (v ALA) Result() uint32 {
	return v.ResultCode
}
