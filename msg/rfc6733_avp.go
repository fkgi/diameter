package msg

import "net"

// Get search AVP that have code and vendor-id.
func (g GroupedAVP) Get(c, v uint32) (*Avp, bool) {
	for _, a := range g {
		if a.Code == c && a.VenID == v {
			return &a, true
		}
	}
	return nil, false
}

// SessionID AVP
type SessionID string

// Encode return AVP struct of this value
func (v *SessionID) Encode() Avp {
	a := Avp{Code: 263, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(string(*v))
	return a
}

// Decode get AVP value
func (v *SessionID) Decode(a Avp) (e error) {
	if a.Code == 263 && a.VenID == 0 {
		s := new(string)
		if e = a.Decode(s); e == nil {
			*v = SessionID(*s)
		}
	} else {
		e = InvalidAVPError{}
	}
	return
}

// GetSessionID get AVP value
func GetSessionID(o GroupedAVP) (SessionID, bool) {
	a, ok := o.Get(263, 0)
	if !ok {
		return "", false
	}
	s := new(string)
	a.Decode(s)
	return SessionID(*s), true
}

// AuthSessionState AVP
type AuthSessionState bool

const (
	// StateMaintained is value of AuthSessionState
	StateMaintained AuthSessionState = true
	// StateNotMaintained is value of AuthSessionState
	StateNotMaintained AuthSessionState = false
)

// Encode return AVP struct of this value
func (v *AuthSessionState) Encode() Avp {
	a := Avp{Code: 277, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if *v {
		a.Encode(Enumerated(0))
	} else {
		a.Encode(Enumerated(1))
	}
	return a
}

// Decode get AVP value
func (v *AuthSessionState) Decode(a Avp) (e error) {
	if a.Code == 277 && a.VenID == 0 {
		s := new(Enumerated)
		if e = a.Decode(s); e == nil {
			switch *s {
			case 0:
				*v = StateMaintained
			case 1:
				*v = StateNotMaintained
			default:
				e = InvalidAVPError{}
			}
		}
	} else {
		e = InvalidAVPError{}
	}
	return
}

// GetAuthSessionState get AVP value
func GetAuthSessionState(o GroupedAVP) (AuthSessionState, bool) {
	a, ok := o.Get(277, 0)
	if !ok {
		return false, false
	}
	s := new(Enumerated)
	a.Decode(s)
	switch *s {
	case 0:
		return StateMaintained, true
	case 1:
		return StateNotMaintained, true
	}
	return false, false
}

// OriginHost AVP
type OriginHost DiameterIdentity

// Encode return AVP struct of this value
func (v *OriginHost) Encode() Avp {
	a := Avp{Code: 264, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(DiameterIdentity(*v))
	return a
}

// Decode get AVP value
func (v *OriginHost) Decode(a Avp) (e error) {
	if a.Code == 264 && a.VenID == 0 {
		s := new(DiameterIdentity)
		if e = a.Decode(s); e == nil {
			*v = OriginHost(*s)
		}
	} else {
		e = InvalidAVPError{}
	}
	return
}

// GetOriginHost get AVP value
func GetOriginHost(o GroupedAVP) (OriginHost, bool) {
	a, ok := o.Get(264, 0)
	if !ok {
		return "", false
	}
	s := new(DiameterIdentity)
	a.Decode(s)
	return OriginHost(*s), true
}

// OriginRealm AVP
type OriginRealm DiameterIdentity

// Encode return AVP struct of this value
func (v *OriginRealm) Encode() Avp {
	a := Avp{Code: 296, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(DiameterIdentity(*v))
	return a
}

// Decode get AVP value
func (v *OriginRealm) Decode(a Avp) (e error) {
	if a.Code == 296 && a.VenID == 0 {
		s := new(DiameterIdentity)
		if e = a.Decode(s); e == nil {
			*v = OriginRealm(*s)
		}
	} else {
		e = InvalidAVPError{}
	}
	return
}

// GetOriginRealm get AVP value
func GetOriginRealm(o GroupedAVP) (OriginRealm, bool) {
	a, ok := o.Get(296, 0)
	if !ok {
		return "", false
	}
	s := new(DiameterIdentity)
	a.Decode(s)
	return OriginRealm(*s), true
}

// DestinationHost AVP
type DestinationHost DiameterIdentity

// Encode return AVP struct of this value
func (v *DestinationHost) Encode() Avp {
	a := Avp{Code: 293, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(DiameterIdentity(*v))
	return a
}

// Decode get AVP value
func (v *DestinationHost) Decode(a Avp) (e error) {
	if a.Code == 293 && a.VenID == 0 {
		s := new(DiameterIdentity)
		if e = a.Decode(s); e == nil {
			*v = DestinationHost(*s)
		}
	} else {
		e = InvalidAVPError{}
	}
	return
}

// GetDestinationHost get AVP value
func GetDestinationHost(o GroupedAVP) (DestinationHost, bool) {
	a, ok := o.Get(293, 0)
	if !ok {
		return "", false
	}
	s := new(DiameterIdentity)
	a.Decode(s)
	return DestinationHost(*s), true
}

// DestinationRealm AVP
type DestinationRealm DiameterIdentity

// Encode return AVP struct of this value
func (v *DestinationRealm) Encode() Avp {
	a := Avp{Code: 283, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(DiameterIdentity(*v))
	return a
}

// Decode get AVP value
func (v *DestinationRealm) Decode(a Avp) (e error) {
	if a.Code == 283 && a.VenID == 0 {
		s := new(DiameterIdentity)
		if e = a.Decode(s); e == nil {
			*v = DestinationRealm(*s)
		}
	} else {
		e = InvalidAVPError{}
	}
	return
}

// GetDestinationRealm get AVP value
func GetDestinationRealm(o GroupedAVP) (DestinationRealm, bool) {
	a, ok := o.Get(283, 0)
	if !ok {
		return "", false
	}
	s := new(DiameterIdentity)
	a.Decode(s)
	return DestinationRealm(*s), true
}

// HostIPAddress AVP
type HostIPAddress net.IP

// Encode return AVP struct of this value
func (v HostIPAddress) Encode() Avp {
	a := Avp{Code: 257, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(net.IP(v))
	return a
}

// GetHostIPAddress get AVP value
func GetHostIPAddress(o GroupedAVP) (HostIPAddress, bool) {
	a, ok := o.Get(257, 0)
	if !ok {
		return nil, false
	}
	s := new(net.IP)
	a.Decode(s)
	return HostIPAddress(*s), true
}

// GetHostIPAddresses get AVP value
func GetHostIPAddresses(o GroupedAVP) (r []HostIPAddress) {
	for _, a := range o {
		if a.Code == 257 && a.VenID == 0 {
			s := new(net.IP)
			a.Decode(s)
			r = append(r, HostIPAddress(*s))
		}
	}
	return
}

// VendorID AVP
type VendorID uint32

// Encode return AVP struct of this value
func (v VendorID) Encode() Avp {
	a := Avp{Code: 266, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetVendorID get AVP value
func GetVendorID(o GroupedAVP) (VendorID, bool) {
	a, ok := o.Get(266, 0)
	if !ok {
		return 0, false
	}
	s := new(uint32)
	a.Decode(s)
	return VendorID(*s), true
}

// ProductName AVP
type ProductName string

// Encode return AVP struct of this value
func (v ProductName) Encode() Avp {
	a := Avp{Code: 269, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(string(v))
	return a
}

// GetProductName get AVP value
func GetProductName(o GroupedAVP) (ProductName, bool) {
	a, ok := o.Get(269, 0)
	if !ok {
		return "", false
	}
	s := new(string)
	a.Decode(s)
	return ProductName(*s), true
}

// ResultCode AVP
type ResultCode uint32

// Encode return AVP struct of this value
func (v ResultCode) Encode() Avp {
	a := Avp{Code: 268, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetResultCode get AVP value
func GetResultCode(o GroupedAVP) (ResultCode, bool) {
	a, ok := o.Get(268, 0)
	if !ok {
		return 0, false
	}
	s := new(uint32)
	a.Decode(s)
	return ResultCode(*s), true
}

// DisconnectCause AVP
type DisconnectCause Enumerated

// Encode return AVP struct of this value
func (v DisconnectCause) Encode() Avp {
	a := Avp{Code: 273, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(Enumerated(v))
	return a
}

// GetDisconnectCause get AVP value
func GetDisconnectCause(o GroupedAVP) (DisconnectCause, bool) {
	a, ok := o.Get(273, 0)
	if !ok {
		return 0, false
	}
	s := new(Enumerated)
	a.Decode(s)
	return DisconnectCause(*s), true
}

const (
	// Rebooting is Enumerated value 0
	Rebooting Enumerated = 0
	// Busy is Enumerated value 1
	Busy Enumerated = 1
	// DoNotWantToTalkToYou is Enumerated value 2
	DoNotWantToTalkToYou Enumerated = 2
)

// UserName AVP
type UserName string

// Encode return AVP struct of this value
func (v UserName) Encode() Avp {
	a := Avp{Code: 1, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(string(v))
	return a
}

// GetUserName get AVP value
func GetUserName(o GroupedAVP) (UserName, bool) {
	a, ok := o.Get(1, 0)
	if !ok {
		return "", false
	}
	s := new(string)
	a.Decode(s)
	return UserName(*s), true
}

// FirmwareRevision AVP
type FirmwareRevision uint32

// Encode return AVP struct of this value
func (v FirmwareRevision) Encode() Avp {
	a := Avp{Code: 267, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetFirmwareRevision get AVP value
func GetFirmwareRevision(o GroupedAVP) (FirmwareRevision, bool) {
	a, ok := o.Get(267, 0)
	if !ok {
		return 0, false
	}
	s := new(uint32)
	a.Decode(s)
	return FirmwareRevision(*s), true
}

// OriginStateID AVP
type OriginStateID uint32

// Encode return AVP struct of this value
func (v OriginStateID) Encode() Avp {
	a := Avp{Code: 278, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetOriginStateID get AVP value
func GetOriginStateID(o GroupedAVP) (OriginStateID, bool) {
	a, ok := o.Get(278, 0)
	if !ok {
		return 0, false
	}
	s := new(uint32)
	a.Decode(s)
	return OriginStateID(*s), true
}

// SupportedVendorID AVP
type SupportedVendorID uint32

// Encode return AVP struct of this value
func (v SupportedVendorID) Encode() Avp {
	a := Avp{Code: 265, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetSupportedVendorID get AVP value
func GetSupportedVendorID(o GroupedAVP) (SupportedVendorID, bool) {
	a, ok := o.Get(265, 0)
	if !ok {
		return 0, false
	}
	s := new(uint32)
	a.Decode(s)
	return SupportedVendorID(*s), true
}

// GetSupportedVendorIDs get AVP value
func GetSupportedVendorIDs(o GroupedAVP) (r []SupportedVendorID) {
	for _, a := range o {
		if a.Code == 265 && a.VenID == 0 {
			s := new(uint32)
			a.Decode(s)
			r = append(r, SupportedVendorID(*s))
		}
	}
	return
}

// ApplicationID is ID of diameter application
type ApplicationID interface {
	Equals(ApplicationID) bool
	Encode() Avp
}

func getApplicationID(o GroupedAVP) (id ApplicationID, ok bool) {
	id, ok = GetAuthApplicationID(o)
	if !ok {
		id, ok = GetAcctApplicationID(o)
	}
	return
}

// AuthApplicationID AVP
type AuthApplicationID uint32

// Equals compare Application ID
func (v AuthApplicationID) Equals(i ApplicationID) bool {
	if a, ok := i.(AuthApplicationID); ok {
		return a == v
	}
	return false
}

// Encode return AVP struct of this value
func (v AuthApplicationID) Encode() Avp {
	a := Avp{Code: 258, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetAuthApplicationID get AVP value
func GetAuthApplicationID(o GroupedAVP) (AuthApplicationID, bool) {
	a, ok := o.Get(258, 0)
	if !ok {
		return 0, false
	}
	s := new(uint32)
	a.Decode(s)
	return AuthApplicationID(*s), true
}

// GetAuthApplicationIDs get AVP value
func GetAuthApplicationIDs(o GroupedAVP) (r []AuthApplicationID) {
	for _, a := range o {
		if a.Code == 258 && a.VenID == 0 {
			s := new(uint32)
			a.Decode(s)
			r = append(r, AuthApplicationID(*s))
		}
	}
	return
}

// AcctApplicationID AVP
type AcctApplicationID uint32

// Equals compare Application ID
func (v AcctApplicationID) Equals(i ApplicationID) bool {
	if a, ok := i.(AcctApplicationID); ok {
		return a == v
	}
	return false
}

// Encode return AVP struct of this value
func (v AcctApplicationID) Encode() Avp {
	a := Avp{Code: 259, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetAcctApplicationID get AVP value
func GetAcctApplicationID(o GroupedAVP) (AcctApplicationID, bool) {
	a, ok := o.Get(259, 0)
	if !ok {
		return 0, false
	}
	s := new(uint32)
	a.Decode(s)
	return AcctApplicationID(*s), true
}

// GetAcctApplicationIDs get AVP value
func GetAcctApplicationIDs(o GroupedAVP) (r []AcctApplicationID) {
	for _, a := range o {
		if a.Code == 259 && a.VenID == 0 {
			s := new(uint32)
			a.Decode(s)
			r = append(r, AcctApplicationID(*s))
		}
	}
	return
}

// VendorSpecificApplicationID AVP
type VendorSpecificApplicationID struct {
	VendorID
	App ApplicationID
}

// Encode return AVP struct of this value
func (v VendorSpecificApplicationID) Encode() Avp {
	a := Avp{Code: 260, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(GroupedAVP([]Avp{
		v.VendorID.Encode(),
		v.App.Encode()}))
	return a
}

// GetVendorSpecificApplicationID get AVP value
func GetVendorSpecificApplicationID(o GroupedAVP) (VendorSpecificApplicationID, bool) {
	a, ok := o.Get(260, 0)
	if !ok {
		return VendorSpecificApplicationID{}, false
	}
	s := VendorSpecificApplicationID{}
	o = GroupedAVP{}
	a.Decode(&o)
	s.VendorID, _ = GetVendorID(o)
	s.App, _ = getApplicationID(o)
	return s, true
}

// GetVendorSpecificApplicationIDs get AVP value
func GetVendorSpecificApplicationIDs(o GroupedAVP) (r []VendorSpecificApplicationID) {
	for _, a := range o {
		if a.Code == 260 && a.VenID == 0 {
			s := VendorSpecificApplicationID{}
			o = GroupedAVP{}
			a.Decode(&o)
			s.VendorID, _ = GetVendorID(o)
			s.App, _ = getApplicationID(o)
			r = append(r, s)
		}
	}
	return
}

// ErrorMessage AVP
type ErrorMessage string

// Encode return AVP struct of this value
func (v ErrorMessage) Encode() Avp {
	a := Avp{Code: 281, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(string(v))
	return a
}

// GetErrorMessage get AVP value
func GetErrorMessage(o GroupedAVP) (ErrorMessage, bool) {
	a, ok := o.Get(281, 0)
	if !ok {
		return "", false
	}
	s := new(string)
	a.Decode(s)
	return ErrorMessage(*s), true
}

// FailedAVP AVP
type FailedAVP GroupedAVP

// Encode return AVP struct of this value
func (v FailedAVP) Encode() Avp {
	a := Avp{Code: 279, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(GroupedAVP(v))
	return a
}

// GetFailedAVP get AVP value
func GetFailedAVP(o GroupedAVP) (FailedAVP, bool) {
	a, ok := o.Get(279, 0)
	if !ok {
		return nil, false
	}
	s := new(GroupedAVP)
	a.Decode(s)
	return FailedAVP(*s), true
}

// ExperimentalResult AVP
type ExperimentalResult struct {
	VendorID
	Code uint32
}

// Encode return AVP struct of this value
func (v ExperimentalResult) Encode() Avp {
	t := make([]Avp, 2)
	t[0] = v.VendorID.Encode()
	t[1] = Avp{Code: 298, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	t[1].Encode(v.Code)

	a := Avp{Code: 297, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(GroupedAVP(t))
	return a
}

// GetExperimentalResult get AVP value
func GetExperimentalResult(o GroupedAVP) (ExperimentalResult, bool) {
	s := ExperimentalResult{}
	a, ok := o.Get(297, 0)
	if !ok {
		return s, false
	}
	o = GroupedAVP{}
	a.Decode(&o)
	if t, ok := GetVendorID(o); ok {
		s.VendorID = t
	}
	if t, ok := o.Get(298, 0); ok {
		t.Decode(&s.Code)
	}
	return s, true
}

// RouteRecord AVP
type RouteRecord DiameterIdentity

// Encode return AVP struct of this value
func (v RouteRecord) Encode() Avp {
	a := Avp{Code: 282, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(DiameterIdentity(v))
	return a
}

// GetRouteRecord get AVP value
func GetRouteRecord(o GroupedAVP) (RouteRecord, bool) {
	a, ok := o.Get(282, 0)
	if !ok {
		return "", false
	}
	s := new(DiameterIdentity)
	a.Decode(s)
	return RouteRecord(*s), true
}

// GetRouteRecords get AVP value
func GetRouteRecords(o GroupedAVP) (r []RouteRecord) {
	for _, a := range o {
		if a.Code == 282 && a.VenID == 0 {
			s := new(DiameterIdentity)
			a.Decode(s)
			r = append(r, RouteRecord(*s))
		}
	}
	return
}

// ProxyInfo AVP
type ProxyInfo struct {
	DiameterIdentity
	State string
}

// Encode return AVP struct of this value
func (v ProxyInfo) Encode() Avp {
	t := make([]Avp, 2)
	t[0] = Avp{Code: 280, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	t[0].Encode(v.DiameterIdentity)
	t[1] = Avp{Code: 33, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	t[1].Encode([]byte(v.State))

	a := Avp{Code: 284, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(GroupedAVP(t))
	return a
}

// GetProxyInfo get AVP value
func GetProxyInfo(o GroupedAVP) (ProxyInfo, bool) {
	s := ProxyInfo{}
	a, ok := o.Get(284, 0)
	if !ok {
		return s, false
	}
	o = GroupedAVP{}
	a.Decode(&o)
	if t, ok := o.Get(280, 0); ok {
		t.Decode(&s.DiameterIdentity)
	}
	if t, ok := o.Get(33, 0); ok {
		stat := new([]byte)
		t.Decode(stat)
		s.State = string(*stat)
	}
	return s, true
}

// GetProxyInfos get AVP value
func GetProxyInfos(o GroupedAVP) (r []ProxyInfo) {
	for _, a := range o {
		if a.Code == 284 && a.VenID == 0 {
			s := ProxyInfo{}
			o2 := GroupedAVP{}
			a.Decode(&o2)
			if t, ok := o2.Get(280, 0); ok {
				t.Decode(&s.DiameterIdentity)
			}
			if t, ok := o2.Get(33, 0); ok {
				stat := new([]byte)
				t.Decode(stat)
				s.State = string(*stat)
			}
			r = append(r, ProxyInfo(s))
		}
	}
	return
}