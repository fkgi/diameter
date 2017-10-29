package diameter

// GenericReq is generic format of diameter request
type GenericReq struct {
	FlgP     bool   // Proxiable
	FlgT     bool   // Potentially re-transmitted message
	Code     uint32 // Command-Code (24bit)
	AppID    uint32 // Application-ID
	Stateful bool

	OriginHost       Identity
	OriginRealm      Identity
	DestinationHost  Identity
	DestinationRealm Identity

	AVP []RawAVP
}

// ToRaw return RawMsg struct of this value
func (v GenericReq) ToRaw(s string) RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: true, FlgP: v.FlgP, FlgE: false, FlgT: v.FlgT,
		Code: v.Code, AppID: v.AppID,
		AVP: make([]RawAVP, 0, len(v.AVP))}

	m.AVP = append(m.AVP, setSessionID(s))
	m.AVP = append(m.AVP, setVendorSpecAppID(v.AppID))
	m.AVP = append(m.AVP, setAuthSessionState(v.Stateful))

	m.AVP = append(m.AVP, setOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, setOriginRealm(v.OriginRealm))
	if len(v.DestinationHost) != 0 {
		m.AVP = append(m.AVP, setDestinationHost(v.DestinationHost))
	}
	m.AVP = append(m.AVP, setDestinationRealm(v.DestinationRealm))

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
		if a.VenID == 0 && a.Code == 263 {
			s, e = getSessionID(a)
		} else if a.VenID == 0 && a.Code == 260 {
			// skip Vendor-Specific-Application-ID
		} else if a.VenID == 0 && a.Code == 277 {
			v.Stateful, e = getAuthSessionState(a)
		} else if a.VenID == 0 && a.Code == 264 {
			v.OriginHost, e = getOriginHost(a)
		} else if a.VenID == 0 && a.Code == 296 {
			v.OriginRealm, e = getOriginRealm(a)
		} else if a.VenID == 0 && a.Code == 293 {
			v.DestinationHost, e = getDestinationHost(a)
		} else if a.VenID == 0 && a.Code == 283 {
			v.DestinationRealm, e = getDestinationRealm(a)
		} else {
			a2 := RawAVP{
				FlgV: a.FlgV, FlgM: a.FlgM, FlgP: a.FlgP,
				Code: a.Code, VenID: a.VenID,
				data: make([]byte, len(a.data))}
			copy(a2.data, a.data)
			v.AVP = append(v.AVP, a2)
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
	AppID    uint32 // Application-ID
	Stateful bool

	ResultCode  uint32
	OriginHost  Identity
	OriginRealm Identity

	AVP []RawAVP
}

// ToRaw return RawMsg struct of this value
func (v GenericAns) ToRaw(s string) RawMsg {
	return RawMsg(v).Clone()
}

// FromRaw make this value from RawMsg struct
func (GenericAns) FromRaw(m RawMsg) (Answer, string, error) {
	return GenericAns(m.Clone()), "", nil
}

// Result returns result-code
func (v GenericAns) Result() uint32 {
	for _, a := range v.AVP {
		if a.VenID != 0 || a.Code != 268 {
			continue
		} else if r, e := getResultCode(a); e == nil {
			return r
		}
	}
	return 0
}
