package msg

import (
	"net"
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

// SessionID AVP
func SessionID(s string) Avp {
	a := Avp{Code: uint32(263), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(s)
	return a
}

// AuthSessionState AVP
func AuthSessionState(e Enumerated) Avp {
	a := Avp{Code: uint32(277), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(e)
	return a
}

const (
	AuthSessionState_STATE_MAINTAINED    Enumerated = 0
	AuthSessionState_NO_STATE_MAINTAINED Enumerated = 1
)

// OriginHost AVP
func OriginHost(i DiameterIdentity) Avp {
	a := Avp{Code: uint32(264), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
}

// OriginRealm AVP
func OriginRealm(i DiameterIdentity) Avp {
	a := Avp{Code: uint32(296), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
}

// DestinationHost AVP
func DestinationHost(i DiameterIdentity) Avp {
	a := Avp{Code: uint32(293), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
}

// DestinationRealm AVP
func DestinationRealm(i DiameterIdentity) Avp {
	a := Avp{Code: uint32(283), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
}

// HostIPAddress AVP
func HostIPAddress(i net.IP) Avp {
	a := Avp{Code: uint32(257), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
}

// VendorID AVP
func VendorID(i uint32) Avp {
	a := Avp{Code: uint32(266), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
}

// ProductName AVP
func ProductName(i string) Avp {
	a := Avp{Code: uint32(269), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
}

// ResultCode AVP
func ResultCode(i uint32) Avp {
	a := Avp{Code: uint32(268), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
}

const (
	DIAMETER_MULTI_ROUND_AUTH          uint32 = 1001
	DIAMETER_SUCCESS                   uint32 = 2001
	DIAMETER_LIMITED_SUCCESS           uint32 = 2002
	DIAMETER_COMMAND_UNSUPPORTED       uint32 = 3001
	DIAMETER_UNABLE_TO_DELIVER         uint32 = 3002
	DIAMETER_REALM_NOT_SERVED          uint32 = 3003
	DIAMETER_TOO_BUSY                  uint32 = 3004
	DIAMETER_LOOP_DETECTED             uint32 = 3005
	DIAMETER_REDIRECT_INDICATION       uint32 = 3006
	DIAMETER_APPLICATION_UNSUPPORTED   uint32 = 3007
	DIAMETER_INVALID_HDR_BITS          uint32 = 3008
	DIAMETER_INVALID_AVP_BITS          uint32 = 3009
	DIAMETER_UNKNOWN_PEER              uint32 = 3010
	DIAMETER_AUTHENTICATION_REJECTED   uint32 = 4001
	DIAMETER_OUT_OF_SPACE              uint32 = 4002
	DIAMETER_ELECTION_LOST             uint32 = 4003
	DIAMETER_AVP_UNSUPPORTED           uint32 = 5001
	DIAMETER_UNKNOWN_SESSION_ID        uint32 = 5002
	DIAMETER_AUTHORIZATION_REJECTED    uint32 = 5003
	DIAMETER_INVALID_AVP_VALUE         uint32 = 5004
	DIAMETER_MISSING_AVP               uint32 = 5005
	DIAMETER_RESOURCES_EXCEEDED        uint32 = 5006
	DIAMETER_CONTRADICTING_AVPS        uint32 = 5007
	DIAMETER_AVP_NOT_ALLOWED           uint32 = 5008
	DIAMETER_AVP_OCCURS_TOO_MANY_TIMES uint32 = 5009
	DIAMETER_NO_COMMON_APPLICATION     uint32 = 5010
	DIAMETER_UNSUPPORTED_VERSION       uint32 = 5011
	DIAMETER_UNABLE_TO_COMPLY          uint32 = 5012
	DIAMETER_INVALID_BIT_IN_HEADER     uint32 = 5013
	DIAMETER_INVALID_AVP_LENGTH        uint32 = 5014
	DIAMETER_INVALID_MESSAGE_LENGTH    uint32 = 5015
	DIAMETER_INVALID_AVP_BIT_COMBO     uint32 = 5016
	DIAMETER_NO_COMMON_SECURITY        uint32 = 5017
)

// DisconnectCause AVP
func DisconnectCause(i Enumerated) Avp {
	a := Avp{Code: uint32(273), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
}

const (
	DisconnectCause_REBOOTING                  Enumerated = 0
	DisconnectCause_BUSY                       Enumerated = 1
	DisconnectCause_DO_NOT_WANT_TO_TALK_TO_YOU Enumerated = 2
)

// UserName AVP
func UserName(s string) Avp {
	a := Avp{Code: uint32(1), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(s)
	return a
}

// FirmwareRevision AVP
func FirmwareRevision(i uint32) Avp {
	a := Avp{Code: uint32(267), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
}

// SupportedVendorID AVP
func SupportedVendorID(i uint32) Avp {
	a := Avp{Code: uint32(265), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
}

// AuthApplicationID AVP
func AuthApplicationID(i uint32) Avp {
	a := Avp{Code: uint32(258), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
}

// AcctApplicationID AVP
func AcctApplicationID(i uint32) Avp {
	a := Avp{Code: uint32(259), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
}

// VendorSpecificApplicationID AVP
func VendorSpecificApplicationID(ven, app uint32, isAuth bool) Avp {
	t := make([]Avp, 2)
	t[0] = VendorID(ven)
	if isAuth {
		t[1] = AuthApplicationID(app)
	} else {
		t[1] = AcctApplicationID(app)
	}

	a := Avp{Code: uint32(260), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(t)
	return a
}

// ErrorMessage AVP
func ErrorMessage(s string) Avp {
	a := Avp{Code: uint32(281), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(s)
	return a
}

// FailedAVP AVP
func FailedAVP(f []Avp) Avp {
	a := Avp{Code: uint32(279), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(f)
	return a
}

// ExperimentalResult AVP
func ExperimentalResult(ven, code uint32) Avp {
	t := make([]Avp, 2)
	t[0] = VendorID(ven)
	t[1] = Avp{Code: uint32(298), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	t[1].Encode(code)

	a := Avp{Code: uint32(297), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(t)
	return a
}

// RouteRecord AVP
func RouteRecord(i DiameterIdentity) Avp {
	a := Avp{Code: uint32(282), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
}

// ProxyInfo AVP
func ProxyInfo(host DiameterIdentity, state []byte) Avp {
	t := make([]Avp, 2)
	t[0] = Avp{Code: uint32(280), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	t[0].Encode(host)
	t[1] = Avp{Code: uint32(33), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	t[1].Encode(state)

	a := Avp{Code: uint32(284), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(t)
	return a
}
