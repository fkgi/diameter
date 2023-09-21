package diameter

import "fmt"

// SetAuthAppID make Auth-Application-Id AVP
func SetAuthAppID(v uint32) (a AVP) {
	a = AVP{Code: 258, Mandatory: true}
	a.Encode(v)
	return
}

// GetAuthAppID read Auth-Application-Id AVP
func GetAuthAppID(a AVP) (v uint32, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

// SetVendorID make Vendor-Id AVP
func SetVendorID(v uint32) (a AVP) {
	a = AVP{Code: 266, Mandatory: true}
	a.Encode(v)
	return
}

// GetVendorID read Vendor-Id AVP
func GetVendorID(a AVP) (v uint32, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

// SetVendorSpecAppID make Vendor-Specific-Application-Id AVP
func SetVendorSpecAppID(vi, ai uint32) (a AVP) {
	a = AVP{Code: 260, Mandatory: true}
	a.Encode([]AVP{SetVendorID(vi), SetAuthAppID(ai)})
	return
}

// GetVendorSpecAppID read Vendor-Specific-Application-Id AVP
func GetVendorSpecAppID(a AVP) (vi, ai uint32, e error) {
	o := []AVP{}
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&o)
	}
	for _, a := range o {
		if a.VendorID != 0 {
			continue
		}
		switch a.Code {
		case 266:
			vi, e = GetVendorID(a)
		case 258:
			ai, e = GetAuthAppID(a)
		}
		if e != nil {
			break
		}
	}
	if e == nil && (vi == 0 || ai == 0) {
		e = InvalidAVP{Code: MissingAvp, AVP: a}
	}
	return
}

// SetSessionID make Session-ID AVP
func SetSessionID(v string) (a AVP) {
	a = AVP{Code: 263, Mandatory: true}
	a.Encode(v)
	return
}

// GetSessionID read Session-ID AVP
func GetSessionID(a AVP) (v string, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

// SetOriginHost make Origin-Host AVP
func SetOriginHost(v Identity) (a AVP) {
	a = AVP{Code: 264, Mandatory: true}
	a.Encode(v)
	return
}

// GetOriginHost read Origin-Host AVP
func GetOriginHost(a AVP) (v Identity, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

const (
	MultiRoundAuth uint32 = 1001 // MultiRoundAuth is Result-Code 1001

	Success        uint32 = 2001 // Success is Result-Code 2001
	LimitedSuccess uint32 = 2002 // LimitedSuccess is Result-Code 2002

	CommandUnspported      uint32 = 3001 // CommandUnspported is Result-Code 3001
	UnableToDeliver        uint32 = 3002 // UnableToDeliver is Result-Code 3002
	RealmNotServed         uint32 = 3003 // RealmNotServed is Result-Code 3003
	TooBusy                uint32 = 3004 // TooBusy is Result-Code 3004
	LoopDetected           uint32 = 3005 // LoopDetected is Result-Code 3005
	RedirectIndication     uint32 = 3006 // RedirectIndication is Result-Code 3006
	ApplicationUnsupported uint32 = 3007 // ApplicationUnsupported is Result-Code 3007
	InvalidHdrBits         uint32 = 3008 // InvalidHdrBits is Result-Code 3008
	InvalidAvpBits         uint32 = 3009 // InvalidAvpBits is Result-Code 3009
	UnknownPeer            uint32 = 3010 // UnknownPeer is Result-Code 3010

	AuthenticationRejected uint32 = 4001 // AuthenticationRejected is Result-Code 4001
	OutOfSpace             uint32 = 4002 // OutOfSpace is Result-Code 4002
	ElectionLost           uint32 = 4003 // ElectionLost is Result-Code 4003

	AvpUnsupported        uint32 = 5001 // AvpUnsupported is Result-Code 5001
	UnknownSessionID      uint32 = 5002 // UnknownSessionID is Result-Code 5002
	AuthorizationRejected uint32 = 5003 // AuthorizationRejected is Result-Code 5003
	InvalidAvpValue       uint32 = 5004 // InvalidAvpValue is Result-Code 5004
	MissingAvp            uint32 = 5005 // MissingAvp is Result-Code 5005
	ResourcesExceeded     uint32 = 5006 // ResourcesExceeded is Result-Code 5006
	ContradictingAvps     uint32 = 5007 // ContradictingAvps is Result-Code 5007
	AvpNotAllowed         uint32 = 5008 // AvpNotAllowed is Result-Code 5008
	AvpOccursTooManyTimes uint32 = 5009 // AvpOccursTooManyTimes is Result-Code 5009
	NoCommonApplication   uint32 = 5010 // NoCommonApplication is Result-Code 5010
	UnsupportedVersion    uint32 = 5011 // UnsupportedVersion is Result-Code 5011
	UnableToComply        uint32 = 5012 // UnableToComply is Result-Code 5012
	InvalidBitInHeader    uint32 = 5013 // InvalidBitInHeader is Result-Code 5013
	InvalidAvpLength      uint32 = 5014 // InvalidAvpLength is Result-Code 5014
	InvalidMessageLength  uint32 = 5015 // InvalidMessageLength is Result-Code 5015
	InvalidAvpBitCombo    uint32 = 5016 // InvalidAvpBitCombo is Result-Code 5016
	NoCommonSecurity      uint32 = 5017 // NoCommonSecurity is Result-Code 5017
)

// SetResultCode make Result-Code AVP
func SetResultCode(c uint32) (a AVP) {
	if c < 10000 {
		a = AVP{Code: 268, Mandatory: true}
		a.Encode(c)
		return
	}
	a = AVP{Code: 297, Mandatory: true}
	v := []AVP{{Code: 266, Mandatory: true}, {Code: 298, Mandatory: true}}
	v[0].Encode(c / 10000)
	v[1].Encode(c % 10000)
	a.Encode(v)
	return
}

// GetResultCode read Result-Code AVP
func GetResultCode(a AVP) (c uint32, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else if a.Code == 268 {
		e = a.wrapedDecode(&c)
	} else if a.Code == 297 {
		o := []AVP{}
		if e = a.wrapedDecode(&o); e == nil {
			for _, a := range o {
				switch a.Code {
				case 266:
					if a.VendorID != 0 || !a.Mandatory {
						e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
					} else {
						var i uint32
						e = a.Decode(&i)
						c += i * 10000
					}
				case 298:
					if a.VendorID != 0 || !a.Mandatory {
						e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
					} else {
						var r uint32
						e = a.Decode(&r)
						c += r
					}
				}
				if e != nil {
					break
				}
			}
			if e == nil && c < 10000 {
				e = fmt.Errorf("AVP 266 not found")
				e = InvalidAVP{Code: MissingAvp, AVP: a, E: e}
			} else if e == nil && c%10000 == 0 {
				e = fmt.Errorf("AVP 298 not found")
				e = InvalidAVP{Code: MissingAvp, AVP: a, E: e}
			}
		}
	}
	return
}

// SetAuthSessionState make Auth-Session-State AVP
func SetAuthSessionState(v bool) (a AVP) {
	a = AVP{Code: 277, Mandatory: true}
	if v {
		// value is STATE_MAINTAINED
		a.Encode(Enumerated(0))
	} else {
		// value is NO_STATE_MAINTAINED
		a.Encode(Enumerated(1))
	}
	return
}

// GetAuthSessionState read Auth-Session-State AVP
func GetAuthSessionState(a AVP) (v bool, e error) {
	s := new(Enumerated)
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else if e = a.wrapedDecode(s); e != nil {
		switch *s {
		case 0:
			v = true
		case 1:
			v = false
		default:
			e = InvalidAVP{Code: InvalidAvpValue, AVP: a}
		}
	}
	return
}

// SetFailedAVP make Failed-AVP AVP
func SetFailedAVP(v []AVP) (a AVP) {
	a = AVP{Code: 279, Mandatory: true}
	a.Encode(v)
	return
}

// GetFailedAVP read Failed-AVP AVP
func GetFailedAVP(a AVP) (v []AVP, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

// SetRouteRecord make Route-Record AVP
func SetRouteRecord(v Identity) (a AVP) {
	a = AVP{Code: 282, Mandatory: true}
	a.Encode(v)
	return
}

// GetRouteRecord read Route-Record AVP
func GetRouteRecord(a AVP) (v Identity, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

// SetDestinationRealm make Destination-Realm AVP
func SetDestinationRealm(v Identity) (a AVP) {
	a = AVP{Code: 283, Mandatory: true}
	a.Encode(v)
	return
}

// GetDestinationRealm read Destination-Realm AVP
func GetDestinationRealm(a AVP) (v Identity, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

// SetDestinationHost make Destination-Host AVP
func SetDestinationHost(v Identity) (a AVP) {
	a = AVP{Code: 293, Mandatory: true}
	a.Encode(v)
	return
}

// GetDestinationHost read Destination-Host AVP
func GetDestinationHost(a AVP) (v Identity, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

// SetOriginRealm make Origin-Realm AVP
func SetOriginRealm(v Identity) (a AVP) {
	a = AVP{Code: 296, Mandatory: true}
	a.Encode(v)
	return
}

// GetOriginRealm read Origin-Realm AVP
func GetOriginRealm(a AVP) (v Identity, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

/*
// ProxyHost AVP
type ProxyHost diameter.Identity

// ToRaw return AVP struct of this value
func (v *ProxyHost) ToRaw() diameter.AVP {
	a := diameter.AVP{Code: 280, VenID: 0,
		FlgV: false, Mandatory: true, FlgP: false}
	if v != nil {
		a.Encode(diameter.Identity(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *ProxyHost) FromRaw(a diameter.AVP) (e error) {
	if e = a.Validate(0, 280, false, true, false); e != nil {
		return
	}
	s := new(diameter.Identity)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = ProxyHost(*s)
	return
}

// ProxyState AVP
type ProxyState []byte

// ToRaw return AVP struct of this value
func (v *ProxyState) ToRaw() diameter.AVP {
	a := diameter.AVP{Code: 33, VenID: 0,
		FlgV: false, Mandatory: true, FlgP: false}
	if v != nil {
		a.Encode([]byte(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *ProxyState) FromRaw(a diameter.AVP) (e error) {
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
func (v *ProxyInfo) ToRaw() diameter.AVP {
	a := diameter.AVP{Code: 284, VenID: 0,
		FlgV: false, Mandatory: true, FlgP: false}
	if v != nil {
		t := []diameter.AVP{
			v.ProxyHost.ToRaw(),
			v.ProxyState.ToRaw()}
		a.Encode(diameter.GroupedAVP(t))
	}
	return a
}

// FromRaw get AVP value
func (v *ProxyInfo) FromRaw(a diameter.AVP) (e error) {
	if e = a.Validate(0, 284, false, true, false); e != nil {
		return
	}
	o := diameter.GroupedAVP{}
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
