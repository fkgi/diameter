package diameter

// GenericReq is generic format of diameter request
type GenericReq RawMsg

// ToRaw return RawMsg struct of this value
func (v GenericReq) ToRaw(s string) RawMsg {
	return RawMsg(v).Clone()
}

// FromRaw make this value from RawMsg struct
func (GenericReq) FromRaw(m RawMsg) (Request, string, error) {
	return GenericReq(m.Clone()), "", nil
}

// Failed make error message for timeout
func (v GenericReq) Failed(c uint32) Answer {
	sid := RawAVP{}
	for _, a := range v.AVP {
		if e := a.Validate(0, 263, false, true, false); e == nil {
			sid = RawAVP{
				Code:  a.Code,
				FlgV:  a.FlgV,
				FlgM:  a.FlgM,
				FlgP:  a.FlgP,
				VenID: a.VenID,
				data:  make([]byte, len(a.data))}
			copy(sid.data, a.data)
		}
	}
	return GenericAns{
		Ver:   v.Ver,
		FlgR:  false,
		FlgP:  v.FlgP,
		FlgE:  true,
		FlgT:  false,
		Code:  v.Code,
		AppID: v.AppID,
		HbHID: v.HbHID,
		EtEID: v.EtEID,
		AVP: []RawAVP{
			setResultCode(c),
			setOriginHost(Host),
			setOriginRealm(Realm),
			sid}}
}

// GenericAns is generic format of diameter request
type GenericAns RawMsg

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
