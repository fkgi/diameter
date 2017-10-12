package rfc6733

import (
	"net"

	"github.com/fkgi/diameter/msg"
)

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

// OriginHost AVP
type OriginHost msg.DiameterIdentity

// ToRaw return AVP struct of this value
func (v *OriginHost) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 264, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(msg.DiameterIdentity(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *OriginHost) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 264, false, true, false); e != nil {
		return
	}
	s := new(msg.DiameterIdentity)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = OriginHost(*s)
	return
}

// OriginRealm AVP
type OriginRealm msg.DiameterIdentity

// ToRaw return AVP struct of this value
func (v *OriginRealm) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 296, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(msg.DiameterIdentity(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *OriginRealm) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 296, false, true, false); e != nil {
		return
	}
	s := new(msg.DiameterIdentity)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = OriginRealm(*s)
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

// HostIPAddress AVP
type HostIPAddress net.IP

// ToRaw return AVP struct of this value
func (v *HostIPAddress) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 257, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(net.IP(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *HostIPAddress) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 257, false, true, false); e != nil {
		return
	}
	s := new(net.IP)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = HostIPAddress(*s)
	return
}

// VendorID AVP
type VendorID uint32

// ToRaw return AVP struct of this value
func (v *VendorID) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 266, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(uint32(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *VendorID) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 266, false, true, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = VendorID(*s)
	return
}

// ProductName AVP
type ProductName string

// ToRaw return AVP struct of this value
func (v *ProductName) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 269, VenID: 0,
		FlgV: false, FlgM: false, FlgP: false}
	if v != nil {
		a.Encode(string(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *ProductName) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 269, false, false, false); e != nil {
		return
	}
	s := new(string)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = ProductName(*s)
	return
}

// ResultCode AVP
type ResultCode uint32

// ToRaw return AVP struct of this value
func (v *ResultCode) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 268, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(uint32(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *ResultCode) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 268, false, true, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = ResultCode(*s)
	return
}

// DisconnectCause AVP
type DisconnectCause msg.Enumerated

const (
	// Rebooting is msg.Enumerated value 0
	Rebooting DisconnectCause = 0
	// Busy is msg.Enumerated value 1
	Busy DisconnectCause = 1
	// DoNotWantToTalkToYou is msg.Enumerated value 2
	DoNotWantToTalkToYou DisconnectCause = 2
)

// ToRaw return AVP struct of this value
func (v *DisconnectCause) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 273, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(msg.Enumerated(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *DisconnectCause) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 273, false, true, false); e != nil {
		return
	}
	s := new(msg.Enumerated)
	if e = a.Decode(s); e != nil {
		return
	}
	switch *s {
	case 0, 1, 2:
		*v = DisconnectCause(*s)
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

// FirmwareRevision AVP
type FirmwareRevision uint32

// ToRaw return AVP struct of this value
func (v *FirmwareRevision) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 267, VenID: 0,
		FlgV: false, FlgM: false, FlgP: false}
	if v != nil {
		a.Encode(uint32(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *FirmwareRevision) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 267, false, false, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = FirmwareRevision(*s)
	return
}

// OriginStateID AVP
type OriginStateID uint32

// ToRaw return AVP struct of this value
func (v *OriginStateID) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 278, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(uint32(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *OriginStateID) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 278, false, true, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = OriginStateID(*s)
	return
}

// SupportedVendorID AVP
type SupportedVendorID uint32

// ToRaw return AVP struct of this value
func (v *SupportedVendorID) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 265, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(uint32(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *SupportedVendorID) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 265, false, true, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = SupportedVendorID(*s)
	return
}

// AuthApplicationID AVP
type AuthApplicationID uint32

// ToRaw return AVP struct of this value
func (v *AuthApplicationID) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 258, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(uint32(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *AuthApplicationID) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 258, false, true, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = AuthApplicationID(*s)
	return
}

// VendorSpecificApplicationID AVP
type VendorSpecificApplicationID struct {
	VendorID
	AuthApplicationID
}

// ToRaw return AVP struct of this value
func (v *VendorSpecificApplicationID) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 260, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(msg.GroupedAVP([]msg.RawAVP{
			v.VendorID.ToRaw(),
			v.AuthApplicationID.ToRaw()}))
	}
	return a
}

// FromRaw get AVP value
func (v *VendorSpecificApplicationID) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 260, false, true, false); e != nil {
		return
	}
	o := msg.GroupedAVP{}
	if e = a.Decode(&o); e != nil {
		return
	}
	for _, a := range o {
		if a.VenID != 0 {
			continue
		}
		switch a.Code {
		case 266:
			e = v.VendorID.FromRaw(a)
		case 258:
			e = v.AuthApplicationID.FromRaw(a)
		case 259:
			e = msg.InvalidAVP{}
		}
	}
	return
}

// ErrorMessage AVP
type ErrorMessage string

// ToRaw return AVP struct of this value
func (v *ErrorMessage) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 281, VenID: 0,
		FlgV: false, FlgM: false, FlgP: false}
	if v != nil {
		a.Encode(string(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *ErrorMessage) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 281, false, false, false); e != nil {
		return
	}
	s := new(string)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = ErrorMessage(*s)
	return
}

// FailedAVP AVP
type FailedAVP msg.GroupedAVP

// ToRaw return AVP struct of this value
func (v *FailedAVP) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 279, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(msg.GroupedAVP(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *FailedAVP) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(0, 279, false, true, false); e != nil {
		return
	}
	s := new(msg.GroupedAVP)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = FailedAVP(*s)
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
