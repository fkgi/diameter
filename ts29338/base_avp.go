package ts29338

import (
	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/teldata"
)

func setUserName(v teldata.IMSI) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 1, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v.String())
	return
}

func getUserName(a msg.RawAVP) (v teldata.IMSI, e error) {
	s := new(string)
	if e = a.Validate(0, 1, false, true, false); e != nil {
	} else if e = a.Decode(s); e == nil {
		v, e = teldata.ParseIMSI(*s)
	}
	return
}

func setVendorSpecAppID(ai uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 260, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	v := []msg.RawAVP{
		msg.RawAVP{Code: 266, VenID: 0,
			FlgV: false, FlgM: true, FlgP: false},
		msg.RawAVP{Code: 258, VenID: 0,
			FlgV: false, FlgM: true, FlgP: false}}
	v[0].Encode(10415)
	v[1].Encode(ai)
	a.Encode(v)
	return
}

/*
func getVendorSpecAppID(a msg.RawAVP) (ai uint32, e error) {
	o := []msg.RawAVP{}
	if e = a.Validate(0, 260, false, true, false); e == nil {
		e = a.Decode(&o)
	}
	for _, a := range o {
		if a.VenID != 0 {
			continue
		}
		switch a.Code {
		case 266:
			if e = a.Validate(0, 266, false, true, false); e == nil {
				e = a.Decode(&vi)
			}
		case 258:
			if e = a.Validate(0, 258, false, true, false); e == nil {
				e = a.Decode(&ai)
			}
		}
	}
	if vi == 0 || ai == 0 {
		e = msg.NoMandatoryAVP{}
	}
	return
}
*/
func setSessionID(v string) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 263, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getSessionID(a msg.RawAVP) (v string, e error) {
	if e = a.Validate(0, 263, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setOriginHost(v msg.DiameterIdentity) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 264, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getOriginHost(a msg.RawAVP) (v msg.DiameterIdentity, e error) {
	if e = a.Validate(0, 264, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setResultCode(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 268, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getResultCode(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(0, 268, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setAuthSessionState() (a msg.RawAVP) {
	a = msg.RawAVP{Code: 277, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	// value is NO_STATE_MAINTAINED (1)
	a.Encode(msg.Enumerated(1))
	return
}

/*
func getAuthSessionState(a msg.RawAVP) (v bool, e error) {
	s := new(msg.Enumerated)
	if e = a.Validate(0, 277, false, true, false); e != nil {
	} else if e = a.Decode(s); e != nil {
		switch *s {
		case 0:
			v = StateMaintained
		case 1:
			v = StateNotMaintained
		default:
			e = msg.InvalidAVP{}
		}
	}
	return
}
*/

func setFailedAVP(v msg.GroupedAVP) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 279, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getFailedAVP(a msg.RawAVP) (v msg.GroupedAVP, e error) {
	if e = a.Validate(0, 279, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setErrorMessage(v string) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 281, VenID: 0,
		FlgV: false, FlgM: false, FlgP: false}
	a.Encode(v)
	return
}

func getErrorMessage(a msg.RawAVP) (v string, e error) {
	if e = a.Validate(0, 281, false, false, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setDestinationRealm(v msg.DiameterIdentity) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 283, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getDestinationRealm(a msg.RawAVP) (v msg.DiameterIdentity, e error) {
	if e = a.Validate(0, 283, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setDestinationHost(v msg.DiameterIdentity) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 293, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getDestinationHost(a msg.RawAVP) (v msg.DiameterIdentity, e error) {
	if e = a.Validate(0, 293, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setOriginRealm(v msg.DiameterIdentity) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 296, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getOriginRealm(a msg.RawAVP) (v msg.DiameterIdentity, e error) {
	if e = a.Validate(0, 296, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setExperimentalResult(i, c uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 297, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	v := []msg.RawAVP{
		msg.RawAVP{Code: 266, VenID: 0,
			FlgV: false, FlgM: true, FlgP: false},
		msg.RawAVP{Code: 298, VenID: 0,
			FlgV: false, FlgM: true, FlgP: false}}
	v[0].Encode(i)
	v[1].Encode(c)
	a.Encode(v)
	return
}

func getExperimentalResult(a msg.RawAVP) (i, c uint32, e error) {
	o := []msg.RawAVP{}
	if e = a.Validate(0, 297, false, true, false); e == nil {
		e = a.Decode(&o)
	}
	for _, a := range o {
		if a.VenID != 0 {
			continue
		}
		switch a.Code {
		case 266:
			if e = a.Validate(0, 266, false, true, false); e == nil {
				e = a.Decode(&i)
			}
		case 298:
			if e = a.Validate(0, 298, false, true, false); e == nil {
				e = a.Decode(&c)
			}
		}
	}
	if i == 0 || c == 0 {
		e = msg.NoMandatoryAVP{}
	}
	return
}

func setRouteRecord(v msg.DiameterIdentity) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 282, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getRouteRecord(a msg.RawAVP) (v msg.DiameterIdentity, e error) {
	if e = a.Validate(0, 282, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

/*
// ProxyHost AVP
type ProxyHost msg.DiameterIdentity

// ToRaw return AVP struct of this value
func (v *ProxyHost) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 280, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(msg.DiameterIdentity(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *ProxyHost) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 280, false, true, false); e != nil {
		return
	}
	s := new(msg.DiameterIdentity)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = ProxyHost(*s)
	return
}

// ProxyState AVP
type ProxyState []byte

// ToRaw return AVP struct of this value
func (v *ProxyState) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 33, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode([]byte(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *ProxyState) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 33, false, true, false); e != nil {
		return
	}
	s := new([]byte)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = ProxyState(*s)
	return
}

// ProxyInfo AVP
type ProxyInfo struct {
	ProxyHost
	ProxyState
}

// ToRaw return AVP struct of this value
func (v *ProxyInfo) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 284, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		t := []msg.RawAVP{
			v.ProxyHost.ToRaw(),
			v.ProxyState.ToRaw()}
		a.Encode(msg.GroupedAVP(t))
	}
	return a
}

// FromRaw get AVP value
func (v *ProxyInfo) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 284, false, true, false); e != nil {
		return
	}
	o := msg.GroupedAVP{}
	if e = a.Decode(&o); e != nil {
		return
	}
	*v = ProxyInfo{}
	for _, a := range o {
		if a.VenID != 0 {
			continue
		}
		switch a.Code {
		case 280:
			e = v.ProxyHost.FromRaw(a)
		case 33:
			e = v.ProxyState.FromRaw(a)
		}
	}
	return
}
*/
