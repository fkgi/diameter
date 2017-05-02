package msg

import "net"

/*
const (
	DONT_CACHE                 Enumerated = 0
	ALL_SESSION                Enumerated = 1
	ALL_REALM                  Enumerated = 2
	REALM_AND_APPLICATION      Enumerated = 3
	ALL_APPLICATION            Enumerated = 4
	ALL_HOST                   Enumerated = 5
	ALL_USER                   Enumerated = 6
	AUTHENTICATE_ONLY          Enumerated = 1
	//	AUTHORIZE_ONLY Enumerated = 2
	//	AUTHORIZE_AUTHENTICATE Enumerated = 3
	//	AUTHORIZE_ONLY Enumerated = 0
	//	AUTHORIZE_AUTHENTICATE Enumerated = 1
	REFUSE_SERVICE          Enumerated = 0
	TRY_AGAIN               Enumerated = 1
	ALLOW_SERVICE           Enumerated = 2
	TRY_AGAIN_ALLOW_SERVICE Enumerated = 3
	EVENT_RECORD            Enumerated = 1
	START_RECORD            Enumerated = 2
	INTERIM_RECORD          Enumerated = 3
	STOP_RECORD             Enumerated = 4
	DELIVER_AND_GRANT       Enumerated = 1
	GRANT_AND_STORE         Enumerated = 2
	GRANT_AND_LOSE          Enumerated = 3
)
*/

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
func (v SessionID) Encode() Avp {
	a := Avp{Code: 263, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(string(v))
	return a
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

// AuthSessionState AVP (true=STATE_MAINTAINED / false=STATE_NOT_MAINTAINED)
type AuthSessionState bool

// StateMaintained is value of AuthSessionState
const StateMaintained bool = true

// StateNotMaintained is value of AuthSessionState
const StateNotMaintained bool = false

// Encode return AVP struct of this value
func (v AuthSessionState) Encode() Avp {
	a := Avp{Code: 277, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v {
		a.Encode(Enumerated(0))
	} else {
		a.Encode(Enumerated(1))
	}
	return a
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
		return AuthSessionState(true), true
	case 1:
		return AuthSessionState(false), true
	}
	return false, false
}

// OriginHost AVP
type OriginHost DiameterIdentity

// Encode return AVP struct of this value
func (v OriginHost) Encode() Avp {
	a := Avp{Code: 264, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(DiameterIdentity(v))
	return a
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
func (v OriginRealm) Encode() Avp {
	a := Avp{Code: 296, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(DiameterIdentity(v))
	return a
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
func (v DestinationHost) Encode() Avp {
	a := Avp{Code: 293, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(DiameterIdentity(v))
	return a
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
func (v DestinationRealm) Encode() Avp {
	a := Avp{Code: 283, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(DiameterIdentity(v))
	return a
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

const (
	// DiameterMultiRoundAuth is Result-Code 1001
	DiameterMultiRoundAuth ResultCode = 1001

	// DiameterSuccess is Result-Code 2001
	DiameterSuccess ResultCode = 2001
	// DiameterLimitedSuccess is Result-Code 2002
	DiameterLimitedSuccess ResultCode = 2002

	// DiameterCommandUnspported is Result-Code 3001
	DiameterCommandUnspported ResultCode = 3001
	// DiameterUnableToDeliver is Result-Code 3002
	DiameterUnableToDeliver ResultCode = 3002
	// DiameterRealmNotServed is Result-Code 3003
	DiameterRealmNotServed ResultCode = 3003
	// DiameterTooBusy is Result-Code 3004
	DiameterTooBusy ResultCode = 3004
	// DiameterLoopDetected is Result-Code 3005
	DiameterLoopDetected ResultCode = 3005
	// DiameterRedirectIndication is Result-Code 3006
	DiameterRedirectIndication ResultCode = 3006
	// DiameterApplicationUnsupported is Result-Code 3007
	DiameterApplicationUnsupported ResultCode = 3007
	// DiameterInvalidHdrBits is Result-Code 3008
	DiameterInvalidHdrBits ResultCode = 3008
	// DiameterInvalidAvpBits is Result-Code 3009
	DiameterInvalidAvpBits ResultCode = 3009
	// DiameterUnknownPeer is Result-Code 3010
	DiameterUnknownPeer ResultCode = 3010

	// DiameterAuthenticationRejected is Result-Code 4001
	DiameterAuthenticationRejected ResultCode = 4001
	// DiameterOutOfSpace is Result-Code 4002
	DiameterOutOfSpace ResultCode = 4002
	// DiameterElectionLost is Result-Code 4003
	DiameterElectionLost ResultCode = 4003

	// DiameterAvpUnsupported is Result-Code 5001
	DiameterAvpUnsupported ResultCode = 5001
	// DiameterUnknownSessionID is Result-Code 5002
	DiameterUnknownSessionID ResultCode = 5002
	// DiameterAuthorizationRejected is Result-Code 5003
	DiameterAuthorizationRejected ResultCode = 5003
	// DiameterInvalidAvpValue is Result-Code 5004
	DiameterInvalidAvpValue ResultCode = 5004
	// DiameterMissingAvp is Result-Code 5005
	DiameterMissingAvp ResultCode = 5005
	// DiameterResourcesExceeded is Result-Code 5006
	DiameterResourcesExceeded ResultCode = 5006
	// DiameterContradictingAvps is Result-Code 5007
	DiameterContradictingAvps ResultCode = 5007
	//DiameterAvpNotAllowed is Result-Code 5008
	DiameterAvpNotAllowed ResultCode = 5008
	// DiameterAvpOccursTooManyTimes is Result-Code 5009
	DiameterAvpOccursTooManyTimes ResultCode = 5009
	// DiameterNoCommonApplication is Result-Code 5010
	DiameterNoCommonApplication ResultCode = 5010
	// DiameterUnsupportedVersion is Result-Code 5011
	DiameterUnsupportedVersion ResultCode = 5011
	// DiameterUnableToComply is Result-Code 5012
	DiameterUnableToComply ResultCode = 5012
	// DiameterInvalidBitInHeader is Result-Code 5013
	DiameterInvalidBitInHeader ResultCode = 5013
	// DiameterInvalidAvpLength is Result-Code 5014
	DiameterInvalidAvpLength ResultCode = 5014
	// DiameterInvalidMessageLength is Result-Code 5015
	DiameterInvalidMessageLength ResultCode = 5015
	// DiameterInvalidAvpBitCombo is Result-Code 5016
	DiameterInvalidAvpBitCombo ResultCode = 5016
	// DiameterNoCommonSecurity is Result-Code 5017
	DiameterNoCommonSecurity ResultCode = 5017
)

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

// AuthApplicationID AVP
type AuthApplicationID uint32

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

// AcctApplicationID AVP
type AcctApplicationID uint32

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

// VendorSpecificApplicationID AVP
type VendorSpecificApplicationID struct {
	VendorID
	App interface{}
}

// Encode return AVP struct of this value
func (v VendorSpecificApplicationID) Encode() Avp {
	t := make([]Avp, 2)
	t[0] = v.VendorID.Encode()
	switch d := v.App.(type) {
	case AuthApplicationID:
		t[1] = d.Encode()
	case AcctApplicationID:
		t[1] = d.Encode()
	}
	a := Avp{Code: 260, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(GroupedAVP(t))
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
	if t, ok := GetVendorID(o); ok {
		s.VendorID = t
	}
	if t, ok := GetAuthApplicationID(o); ok {
		s.App = t
	}
	if t, ok := GetAcctApplicationID(o); ok {
		s.App = t
	}
	return s, true
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
