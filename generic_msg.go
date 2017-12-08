package diameter

import (
	"bytes"
	"fmt"
)

// GenericReq is generic format of diameter request
type GenericReq struct {
	FlgP     bool   // Proxiable
	FlgT     bool   // Potentially re-transmitted message
	Code     uint32 // Command-Code (24bit)
	VenID    uint32
	AppID    uint32 // Application-ID
	Stateful bool

	OriginHost       Identity
	OriginRealm      Identity
	DestinationHost  Identity
	DestinationRealm Identity

	AVP []RawAVP
}

func (v GenericReq) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sOrigin-Host       =%s\n", Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm      =%s\n", Indent, v.OriginRealm)
	fmt.Fprintf(w, "%sDestination-Host  =%s\n", Indent, v.DestinationHost)
	fmt.Fprintf(w, "%sDestination-Realm =%s\n", Indent, v.DestinationRealm)
	for i, avp := range v.AVP {
		fmt.Fprintf(w, "%sAVP[%d]    =\n%s", Indent, i, avp)
	}

	return w.String()
}

// ToRaw return RawMsg struct of this value
func (v GenericReq) ToRaw(s string) RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: true, FlgP: v.FlgP, FlgE: false, FlgT: v.FlgT,
		Code: v.Code, AppID: v.AppID,
		AVP: make([]RawAVP, 0, len(v.AVP)+6)}

	m.AVP = append(m.AVP, SetSessionID(s))
	m.AVP = append(m.AVP, SetVendorSpecAppID(v.VenID, v.AppID))
	m.AVP = append(m.AVP, SetAuthSessionState(v.Stateful))

	m.AVP = append(m.AVP, SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, SetOriginRealm(v.OriginRealm))
	if len(v.DestinationHost) != 0 {
		m.AVP = append(m.AVP, SetDestinationHost(v.DestinationHost))
	}
	m.AVP = append(m.AVP, SetDestinationRealm(v.DestinationRealm))

	for _, a := range v.AVP {
		a2 := RawAVP{
			FlgV: a.FlgV, FlgM: a.FlgM, FlgP: a.FlgP,
			Code: a.Code, VenID: a.VenID,
			data: make([]byte, len(a.data))}
		copy(a2.data, a.data)
		m.AVP = append(m.AVP, a2)
	}

	return m
}

// FromRaw make this value from RawMsg struct
func (GenericReq) FromRaw(m RawMsg) (Request, string, error) {
	var s string
	var e error
	v := GenericReq{
		FlgP: m.FlgP, FlgT: m.FlgT,
		Code: m.Code, AppID: m.AppID,
		AVP: make([]RawAVP, 0, len(m.AVP))}
	for _, a := range m.AVP {
		switch a.Code {
		case 263:
			s, e = GetSessionID(a)
		case 260:
			v.VenID, _, e = GetVendorSpecAppID(a)
		case 277:
			v.Stateful, e = GetAuthSessionState(a)
		case 264:
			v.OriginHost, e = GetOriginHost(a)
		case 296:
			v.OriginRealm, e = GetOriginRealm(a)
		case 293:
			v.DestinationHost, e = GetDestinationHost(a)
		case 283:
			v.DestinationRealm, e = GetDestinationRealm(a)
		default:
			a2 := RawAVP{
				FlgV: a.FlgV, FlgM: a.FlgM, FlgP: a.FlgP,
				Code: a.Code, VenID: a.VenID,
				data: make([]byte, len(a.data))}
			copy(a2.data, a.data)
			v.AVP = append(v.AVP, a2)
		}
		if e != nil {
			return nil, "", e
		}
	}

	return v, s, nil
}

// Failed make error message for timeout
func (v GenericReq) Failed(c uint32) Answer {
	return GenericAns{
		FlgP:        v.FlgP,
		Code:        v.Code,
		AppID:       v.AppID,
		Stateful:    v.Stateful,
		ResultCode:  c,
		OriginHost:  Host,
		OriginRealm: Realm}
}

// GenericAns is generic format of diameter request
type GenericAns struct {
	FlgP     bool   // Proxiable
	Code     uint32 // Command-Code (24bit)
	VenID    uint32
	AppID    uint32 // Application-ID
	Stateful bool

	ResultCode  uint32
	OriginHost  Identity
	OriginRealm Identity

	AVP []RawAVP
}

func (v GenericAns) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sResult-Code     =%d\n", Indent, v.ResultCode)
	fmt.Fprintf(w, "%sOrigin-Host     =%s\n", Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm    =%s\n", Indent, v.OriginRealm)
	for i, avp := range v.AVP {
		fmt.Fprintf(w, "%sAVP[%d]  =\n%s", Indent, i, avp)
	}

	return w.String()
}

// ToRaw return RawMsg struct of this value
func (v GenericAns) ToRaw(s string) RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: false, FlgP: v.FlgP, FlgE: false, FlgT: false,
		Code: v.Code, AppID: v.AppID,
		AVP: make([]RawAVP, 0, len(v.AVP)+6)}

	m.AVP = append(m.AVP, SetResultCode(v.ResultCode))
	m.AVP = append(m.AVP, SetSessionID(s))
	m.AVP = append(m.AVP, SetVendorSpecAppID(v.VenID, v.AppID))
	m.AVP = append(m.AVP, SetAuthSessionState(v.Stateful))

	m.AVP = append(m.AVP, SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, SetOriginRealm(v.OriginRealm))

	for _, a := range v.AVP {
		a2 := RawAVP{
			FlgV: a.FlgV, FlgM: a.FlgM, FlgP: a.FlgP,
			Code: a.Code, VenID: a.VenID,
			data: make([]byte, len(a.data))}
		copy(a2.data, a.data)
		m.AVP = append(m.AVP, a2)
	}

	return m
}

// FromRaw make this value from RawMsg struct
func (GenericAns) FromRaw(m RawMsg) (Answer, string, error) {
	var s string
	var e error
	v := GenericAns{
		FlgP: m.FlgP, Code: m.Code, AppID: m.AppID,
		AVP: make([]RawAVP, 0, len(m.AVP))}
	for _, a := range m.AVP {
		switch a.Code {
		case 268, 297:
			v.ResultCode, e = GetResultCode(a)
		case 263:
			s, e = GetSessionID(a)
		case 260:
			v.VenID, _, e = GetVendorSpecAppID(a)
		case 277:
			v.Stateful, e = GetAuthSessionState(a)
		case 264:
			v.OriginHost, e = GetOriginHost(a)
		case 296:
			v.OriginRealm, e = GetOriginRealm(a)
		default:
			a2 := RawAVP{
				FlgV: a.FlgV, FlgM: a.FlgM, FlgP: a.FlgP,
				Code: a.Code, VenID: a.VenID,
				data: make([]byte, len(a.data))}
			copy(a2.data, a.data)
			v.AVP = append(v.AVP, a2)
		}
		if e != nil {
			return nil, "", e
		}
	}

	return v, s, nil
}

// Result returns result-code
func (v GenericAns) Result() uint32 {
	return v.ResultCode
}
