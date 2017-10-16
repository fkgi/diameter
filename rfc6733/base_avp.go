package rfc6733

import (
	"net"

	"github.com/fkgi/diameter/msg"
)

func setHostIPAddress(v net.IP) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 257, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getHostIPAddress(a msg.RawAVP) (v net.IP, e error) {
	if e = a.Validate(0, 257, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setAuthAppID(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 258, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getAuthAppID(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(0, 258, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setVendorSpecAppID(vi, ai uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 260, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode([]msg.RawAVP{
		setVendorID(vi),
		setAuthAppID(ai)})
	return
}

func getVendorSpecAppID(a msg.RawAVP) (vi, ai uint32, e error) {
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
			vi, e = getVendorID(a)
		case 258:
			ai, e = getAuthAppID(a)
		}
	}
	if vi == 0 || ai == 0 {
		e = msg.NoMandatoryAVP{}
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

func setVendorID(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 266, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getVendorID(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(0, 266, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setFirmwareRevision(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 267, VenID: 0,
		FlgV: false, FlgM: false, FlgP: false}
	a.Encode(v)
	return
}

func getFirmwareRevision(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(0, 267, false, false, false); e == nil {
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

func setProductName(v string) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 269, VenID: 0,
		FlgV: false, FlgM: false, FlgP: false}
	a.Encode(v)
	return
}

func getProductName(a msg.RawAVP) (v string, e error) {
	if e = a.Validate(0, 269, false, false, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

const (
	// Rebooting is Enumerated value 0
	Rebooting msg.Enumerated = 0
	// Busy is Enumerated value 1
	Busy msg.Enumerated = 1
	// DoNotWantToTalkToYou is Enumerated value 2
	DoNotWantToTalkToYou msg.Enumerated = 2
)

func setDisconnectCause(v msg.Enumerated) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 273, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getDisconnectCause(a msg.RawAVP) (v msg.Enumerated, e error) {
	if e = a.Validate(0, 273, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	if v < 0 || v > 2 {
		e = msg.InvalidAVP{}
	}
	return
}

func setOriginStateID(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 278, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getOriginStateID(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(0, 278, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

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

func setSupportedVendorID(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 265, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return a
}

func getSupportedVendorID(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(0, 265, false, true, false); e == nil {
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

/*
// SessionID AVP
type SessionID string

// ToRaw return AVP struct of this value
func (v *SessionID) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 263, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(string(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *SessionID) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 263, false, true, false); e != nil {
		return
	}
	s := new(string)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = SessionID(*s)
	return
}

// AuthSessionState AVP
type AuthSessionState bool

const (
	// StateMaintained is value of AuthSessionState
	StateMaintained AuthSessionState = true
	// StateNotMaintained is value of AuthSessionState
	StateNotMaintained AuthSessionState = false
)

// ToRaw return AVP struct of this value
func (v *AuthSessionState) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 277, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		if *v {
			a.Encode(msg.Enumerated(0))
		} else {
			a.Encode(msg.Enumerated(1))
		}
	}
	return a
}

// FromRaw get AVP value
func (v *AuthSessionState) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 277, false, true, false); e != nil {
		return
	}
	s := new(msg.Enumerated)
	if e = a.Decode(s); e != nil {
		return
	}
	switch *s {
	case 0:
		*v = StateMaintained
	case 1:
		*v = StateNotMaintained
	default:
		e = msg.InvalidAVP{}
	}
	return
}

// UserName AVP
type UserName string

// ToRaw return AVP struct of this value
func (v *UserName) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 1, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(string(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *UserName) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 1, false, true, false); e != nil {
		return
	}
	s := new(string)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = UserName(*s)
	return
}

// DestinationHost AVP
type DestinationHost msg.DiameterIdentity

// ToRaw return AVP struct of this value
func (v *DestinationHost) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 293, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(msg.DiameterIdentity(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *DestinationHost) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 293, false, true, false); e != nil {
		return
	}
	s := new(msg.DiameterIdentity)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = DestinationHost(*s)
	return
}

// DestinationRealm AVP
type DestinationRealm msg.DiameterIdentity

// ToRaw return AVP struct of this value
func (v *DestinationRealm) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 283, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(msg.DiameterIdentity(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *DestinationRealm) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 283, false, true, false); e != nil {
		return
	}
	s := new(msg.DiameterIdentity)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = DestinationRealm(*s)
	return
}

// ExperimentalResultCode AVP
type ExperimentalResultCode uint32

// ToRaw return AVP struct of this value
func (v *ExperimentalResultCode) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 298, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(uint32(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *ExperimentalResultCode) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 298, false, true, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = ExperimentalResultCode(*s)
	return
}

// ExperimentalResult AVP
type ExperimentalResult struct {
	VendorID
	ExperimentalResultCode
}

// ToRaw return AVP struct of this value
func (v *ExperimentalResult) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 297, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		t := []msg.RawAVP{
			v.VendorID.ToRaw(),
			v.ExperimentalResultCode.ToRaw()}
		a.Encode(msg.GroupedAVP(t))
	}
	return a
}

// FromRaw get AVP value
func (v *ExperimentalResult) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 297, false, true, false); e != nil {
		return
	}
	o := msg.GroupedAVP{}
	if e = a.Decode(&o); e != nil {
		return
	}
	*v = ExperimentalResult{}
	for _, a := range o {
		if a.VenID != 0 {
			continue
		}
		switch a.Code {
		case 266:
			e = v.VendorID.FromRaw(a)
		case 298:
			e = v.ExperimentalResultCode.FromRaw(a)
		}
	}
	return
}

// RouteRecord AVP
type RouteRecord msg.DiameterIdentity

// ToRaw return AVP struct of this value
func (v *RouteRecord) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 282, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(msg.DiameterIdentity(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *RouteRecord) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 282, false, true, false); e != nil {
		return
	}
	s := new(msg.DiameterIdentity)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = RouteRecord(*s)
	return
}

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
