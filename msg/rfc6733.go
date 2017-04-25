package msg

import (
	"net"

	"github.com/fkgi/diameter/msg"
)

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
const iana uint32 = 0

// SessionID AVP
type SessionID string

// Encode return AVP struct of this value
func (v SessionID) Encode() Avp {
	a := Avp{Code: uint32(263), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(string(v))
	return a
}

// DecodeSessionID get AVP value
func DecodeSessionID(o GroupedAVP) (r []SessionID) {
	for _, a := range o {
		if a.Code == 263 && a.VenID == 0 {
			s := new(string)
			a.Decode(s)
			r = append(r, SessionID(*s))
		}
	}
	return
}

// AuthSessionState AVP (true=STATE_MAINTAINED / false=STATE_NOT_MAINTAINED)
type AuthSessionState bool

// Encode return AVP struct of this value
func (v AuthSessionState) Encode() Avp {
	a := Avp{Code: uint32(277), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	if v {
		a.Encode(Enumerated(0))
	} else {
		a.Encode(Enumerated(1))
	}
	return a
}

// DecodeAuthSessionState get AVP value
func DecodeAuthSessionState(o GroupedAVP) (r []AuthSessionState) {
	for _, a := range o {
		if a.Code == 277 && a.VenID == 0 {
			s := new(Enumerated)
			a.Decode(s)
			switch *s {
			case 0:
				r = append(r, AuthSessionState(true))
			case 1:
				r = append(r, AuthSessionState(false))
			}
		}
	}
	return
}

// OriginHost AVP
type OriginHost DiameterIdentity

// Encode return AVP struct of this value
func (v OriginHost) Encode() Avp {
	a := Avp{Code: uint32(264), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(DiameterIdentity(v))
	return a
}

// DecodeOriginHost get AVP value
func DecodeOriginHost(o GroupedAVP) (r []OriginHost) {
	for _, a := range o {
		if a.Code == 264 && a.VenID == 0 {
			s := new(DiameterIdentity)
			a.Decode(s)
			r = append(r, OriginHost(*s))
		}
	}
	return
}

// OriginRealm AVP
type OriginRealm DiameterIdentity

// Encode return AVP struct of this value
func (v OriginRealm) Encode() Avp {
	a := Avp{Code: uint32(296), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(DiameterIdentity(v))
	return a
}

// DecodeOriginRealm get AVP value
func DecodeOriginRealm(o GroupedAVP) (r []OriginRealm) {
	for _, a := range o {
		if a.Code == 296 && a.VenID == 0 {
			s := new(DiameterIdentity)
			a.Decode(s)
			r = append(r, OriginRealm(*s))
		}
	}
	return
}

// DestinationHost AVP
type DestinationHost DiameterIdentity

// Encode return AVP struct of this value
func (v DestinationHost) Encode() Avp {
	a := Avp{Code: uint32(293), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(DiameterIdentity(v))
	return a
}

// DecodeDestinationHost get AVP value
func DecodeDestinationHost(o GroupedAVP) (r []DestinationHost) {
	for _, a := range o {
		if a.Code == 293 && a.VenID == 0 {
			s := new(DiameterIdentity)
			a.Decode(s)
			r = append(r, DestinationHost(*s))
		}
	}
	return
}

// DestinationRealm AVP
type DestinationRealm DiameterIdentity

// Encode return AVP struct of this value
func (v DestinationRealm) Encode() Avp {
	a := Avp{Code: uint32(283), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(DiameterIdentity(v))
	return a
}

// DecodeDestinationRealm get AVP value
func DecodeDestinationRealm(o GroupedAVP) (r []DestinationRealm) {
	for _, a := range o {
		if a.Code == 283 && a.VenID == 0 {
			s := new(DiameterIdentity)
			a.Decode(s)
			r = append(r, DestinationRealm(*s))
		}
	}
	return
}

// HostIPAddress AVP
type HostIPAddress net.IP

// Encode return AVP struct of this value
func (v HostIPAddress) Encode() Avp {
	a := Avp{Code: uint32(257), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(net.IP(v))
	return a
}

// DecodeHostIPAddress get AVP value
func DecodeHostIPAddress(o GroupedAVP) (r []HostIPAddress) {
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
	a := Avp{Code: uint32(266), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// DecodeVendorID get AVP value
func DecodeVendorID(o GroupedAVP) (r []VendorID) {
	for _, a := range o {
		if a.Code == 266 && a.VenID == 0 {
			s := new(uint32)
			a.Decode(s)
			r = append(r, VendorID(*s))
		}
	}
	return
}

// ProductName AVP
type ProductName string

// Encode return AVP struct of this value
func (v ProductName) Encode() Avp {
	a := Avp{Code: uint32(269), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(string(v))
	return a
}

// DecodeProductName get AVP value
func DecodeProductName(o GroupedAVP) (r []ProductName) {
	for _, a := range o {
		if a.Code == 269 && a.VenID == 0 {
			s := new(string)
			a.Decode(s)
			r = append(r, ProductName(*s))
		}
	}
	return
}

// ResultCode AVP
type ResultCode uint32

// Encode return AVP struct of this value
func (v ResultCode) Encode() Avp {
	a := Avp{Code: uint32(268), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// DecodeResultCode get AVP value
func DecodeResultCode(o GroupedAVP) (r []ResultCode) {
	for _, a := range o {
		if a.Code == 268 && a.VenID == 0 {
			s := new(uint32)
			a.Decode(s)
			r = append(r, ResultCode(*s))
		}
	}
	return
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
	a := Avp{Code: uint32(273), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(Enumerated(v))
	return a
}

// DecodeDisconnectCause get AVP value
func DecodeDisconnectCause(o GroupedAVP) (r []DisconnectCause) {
	for _, a := range o {
		if a.Code == 273 && a.VenID == 0 {
			s := new(Enumerated)
			a.Decode(s)
			r = append(r, DisconnectCause(*s))
		}
	}
	return
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
	a := Avp{Code: uint32(1), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(string(v))
	return a
}

// DecodeUserName get AVP value
func DecodeUserName(o GroupedAVP) (r []UserName) {
	for _, a := range o {
		if a.Code == 1 && a.VenID == 0 {
			s := new(string)
			a.Decode(s)
			r = append(r, UserName(*s))
		}
	}
	return
}

// FirmwareRevision AVP
type FirmwareRevision uint32

// Encode return AVP struct of this value
func (v FirmwareRevision) Encode() Avp {
	a := Avp{Code: uint32(267), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// DecodeFirmwareRevision get AVP value
func DecodeFirmwareRevision(o GroupedAVP) (r []FirmwareRevision) {
	for _, a := range o {
		if a.Code == 267 && a.VenID == 0 {
			s := new(uint32)
			a.Decode(s)
			r = append(r, FirmwareRevision(*s))
		}
	}
	return
}

// SupportedVendorID AVP
type SupportedVendorID uint32

// Encode return AVP struct of this value
func (v SupportedVendorID) Encode() Avp {
	a := Avp{Code: uint32(265), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// DecodeSupportedVendorID get AVP value
func DecodeSupportedVendorID(o GroupedAVP) (r []SupportedVendorID) {
	for _, a := range o {
		if a.Code == 265 && a.VenID == 0 {
			s := new(uint32)
			a.Decode(s)
			r = append(r, SupportedVendorID(*s))
		}
	}
	return
}

// AuthApplicationID AVP
type AuthApplicationID uint32

// Encode return AVP struct of this value
func (v AuthApplicationID) Encode() Avp {
	a := Avp{Code: uint32(258), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// DecodeAuthApplicationID get AVP value
func DecodeAuthApplicationID(o GroupedAVP) (r []AuthApplicationID) {
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

// Encode return AVP struct of this value
func (v AcctApplicationID) Encode() Avp {
	a := Avp{Code: uint32(259), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// DecodeAcctApplicationID get AVP value
func DecodeAcctApplicationID(o GroupedAVP) (r []AcctApplicationID) {
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
	a := Avp{Code: uint32(260), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(GroupedAVP(t))
	return a
}

// DecodeVendorSpecificApplicationID get AVP value
func DecodeVendorSpecificApplicationID(o GroupedAVP) (r []VendorSpecificApplicationID) {
	for _, a := range o {
		if a.Code == 260 && a.VenID == 0 {
			s := VendorSpecificApplicationID{}
			o2 := new(GroupedAVP)
			a.Decode(o2)
			if t := DecodeVendorID(*o2); len(t) != 0 {
				s.VendorID = t[0]
			}
			if t := DecodeAuthApplicationID(*o2); len(t) != 0 {
				s.App = t[0]
			}
			if t := DecodeAcctApplicationID(*o2); len(t) != 0 {
				s.App = t[0]
			}
			r = append(r, s)
		}
	}
	return
}

// ErrorMessage AVP
type ErrorMessage string

// Encode return AVP struct of this value
func (v ErrorMessage) Encode() Avp {
	a := Avp{Code: uint32(281), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(string(v))
	return a
}

// DecodeErrorMessage get AVP value
func DecodeErrorMessage(o GroupedAVP) (r []ErrorMessage) {
	for _, a := range o {
		if a.Code == 281 && a.VenID == 0 {
			s := new(string)
			a.Decode(s)
			r = append(r, ErrorMessage(*s))
		}
	}
	return
}

// FailedAVP AVP
type FailedAVP GroupedAVP

// Encode return AVP struct of this value
func (v FailedAVP) Encode() Avp {
	a := Avp{Code: uint32(279), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(GroupedAVP(v))
	return a
}

// DecodeFailedAVP get AVP value
func DecodeFailedAVP(o GroupedAVP) (r []FailedAVP) {
	for _, a := range o {
		if a.Code == 279 && a.VenID == 0 {
			s := new(GroupedAVP)
			a.Decode(s)
			r = append(r, FailedAVP(*s))
		}
	}
	return
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
	t[1] = Avp{Code: uint32(298), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	t[1].Encode(v.Code)

	a := Avp{Code: uint32(297), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(GroupedAVP(t))
	return a
}

// DecodeExperimentalResult get AVP value
func DecodeExperimentalResult(o GroupedAVP) (r []ExperimentalResult) {
	for _, a := range o {
		if a.Code == 297 && a.VenID == 0 {
			s := ExperimentalResult{}
			o2 := new(GroupedAVP)
			a.Decode(o2)
			if t := DecodeVendorID(*o2); len(t) != 0 {
				s.VendorID = t[0]
			}
			for _, a := range *o2 {
				if a.Code == 298 && a.VenID == 0 {
					a.Decode(&s.Code)
					break
				}
			}
			r = append(r, s)
		}
	}
	return
}

// RouteRecord AVP
type RouteRecord DiameterIdentity

// Encode return AVP struct of this value
func (v RouteRecord) Encode() Avp {
	a := Avp{Code: uint32(282), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(DiameterIdentity(v))
	return a
}

// DecodeRouteRecord get AVP value
func DecodeRouteRecord(o GroupedAVP) (r []RouteRecord) {
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
	t[0] = Avp{Code: uint32(280), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	t[0].Encode(v.DiameterIdentity)
	t[1] = Avp{Code: uint32(33), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	t[1].Encode([]byte(v.State))

	a := Avp{Code: uint32(284), VenID: iana,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(msg.GroupedAVP(t))
	return a
}

// DecodeProxyInfo get AVP value
func DecodeProxyInfo(o GroupedAVP) (r []ProxyInfo) {
	for _, a := range o {
		if a.Code == 284 && a.VenID == 0 {
			s := ProxyInfo{}
			o2 := new(GroupedAVP)
			stat := new([]byte)
			a.Decode(o2)
			for _, a := range *o2 {
				if a.Code == 280 && a.VenID == 0 {
					a.Decode(&s.DiameterIdentity)
				}
				if a.Code == 33 && a.VenID == 0 {
					a.Decode(stat)
					s.State = string(*stat)
				}
			}
			r = append(r, s)
		}
	}
	return
}
