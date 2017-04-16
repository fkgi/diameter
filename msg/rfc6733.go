package msg

import (
	"fmt"
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

func Set(o []Avp, d interface{}) ([]Avp, error) {
	var a Avp
	if d == nil {
		return o, fmt.Errorf("nil AVP data")
	}
	switch d := d.(type) {
	case SessionID:
		a = d.avp()
	}
	return append(o, a), nil
}

func Get(o []Avp, d interface{}) error {
	switch d := d.(type) {
	case *SessionID:
		d = d.avp()
	}
	return nil
}

// SessionID AVP
type SessionID string

func (v SessionID) avp() Avp {
	a := Avp{Code: uint32(263), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(string(v))
	return a
}

// AuthSessionState AVP (true=STATE_MAINTAINED / false=STATE_NOT_MAINTAINED)
func AuthSessionState(b bool) Avp {
	a := Avp{Code: uint32(277), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	if b {
		a.Encode(Enumerated(0))
	} else {
		a.Encode(Enumerated(1))
	}
	return a
}

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
	// DiameterMultiRoundAuth is Result-Code 1001
	DiameterMultiRoundAuth uint32 = 1001

	// DiameterSuccess is Result-Code 2001
	DiameterSuccess uint32 = 2001
	// DiameterLimitedSuccess is Result-Code 2002
	DiameterLimitedSuccess uint32 = 2002

	// DiameterCommandUnspported is Result-Code 3001
	DiameterCommandUnspported uint32 = 3001
	// DiameterUnableToDeliver is Result-Code 3002
	DiameterUnableToDeliver uint32 = 3002
	// DiameterRealmNotServed is Result-Code 3003
	DiameterRealmNotServed uint32 = 3003
	// DiameterTooBusy is Result-Code 3004
	DiameterTooBusy uint32 = 3004
	// DiameterLoopDetected is Result-Code 3005
	DiameterLoopDetected uint32 = 3005
	// DiameterRedirectIndication is Result-Code 3006
	DiameterRedirectIndication uint32 = 3006
	// DiameterApplicationUnsupported is Result-Code 3007
	DiameterApplicationUnsupported uint32 = 3007
	// DiameterInvalidHdrBits is Result-Code 3008
	DiameterInvalidHdrBits uint32 = 3008
	// DiameterInvalidAvpBits is Result-Code 3009
	DiameterInvalidAvpBits uint32 = 3009
	// DiameterUnknownPeer is Result-Code 3010
	DiameterUnknownPeer uint32 = 3010

	// DiameterAuthenticationRejected is Result-Code 4001
	DiameterAuthenticationRejected uint32 = 4001
	// DiameterOutOfSpace is Result-Code 4002
	DiameterOutOfSpace uint32 = 4002
	// DiameterElectionLost is Result-Code 4003
	DiameterElectionLost uint32 = 4003

	// DiameterAvpUnsupported is Result-Code 5001
	DiameterAvpUnsupported uint32 = 5001
	// DiameterUnknownSessionID is Result-Code 5002
	DiameterUnknownSessionID uint32 = 5002
	// DiameterAuthorizationRejected is Result-Code 5003
	DiameterAuthorizationRejected uint32 = 5003
	// DiameterInvalidAvpValue is Result-Code 5004
	DiameterInvalidAvpValue uint32 = 5004
	// DiameterMissingAvp is Result-Code 5005
	DiameterMissingAvp uint32 = 5005
	// DiameterResourcesExceeded is Result-Code 5006
	DiameterResourcesExceeded uint32 = 5006
	// DiameterContradictingAvps is Result-Code 5007
	DiameterContradictingAvps uint32 = 5007
	//DiameterAvpNotAllowed is Result-Code 5008
	DiameterAvpNotAllowed uint32 = 5008
	// DiameterAvpOccursTooManyTimes is Result-Code 5009
	DiameterAvpOccursTooManyTimes uint32 = 5009
	// DiameterNoCommonApplication is Result-Code 5010
	DiameterNoCommonApplication uint32 = 5010
	// DiameterUnsupportedVersion is Result-Code 5011
	DiameterUnsupportedVersion uint32 = 5011
	// DiameterUnableToComply is Result-Code 5012
	DiameterUnableToComply uint32 = 5012
	// DiameterInvalidBitInHeader is Result-Code 5013
	DiameterInvalidBitInHeader uint32 = 5013
	// DiameterInvalidAvpLength is Result-Code 5014
	DiameterInvalidAvpLength uint32 = 5014
	// DiameterInvalidMessageLength is Result-Code 5015
	DiameterInvalidMessageLength uint32 = 5015
	// DiameterInvalidAvpBitCombo is Result-Code 5016
	DiameterInvalidAvpBitCombo uint32 = 5016
	// DiameterNoCommonSecurity is Result-Code 5017
	DiameterNoCommonSecurity uint32 = 5017
)

// DisconnectCause AVP
func DisconnectCause(i Enumerated) Avp {
	a := Avp{Code: uint32(273), FlgV: false, FlgM: true, FlgP: false, VenID: uint32(0)}
	a.Encode(i)
	return a
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
