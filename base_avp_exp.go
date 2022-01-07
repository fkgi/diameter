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
