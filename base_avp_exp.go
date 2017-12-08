package diameter

// SetVendorSpecAppID make Vendor-Specific-Application-Id AVP
func SetVendorSpecAppID(vi, ai uint32) (a RawAVP) {
	a = RawAVP{Code: 260, VenID: 0, FlgV: false, FlgM: true, FlgP: false}
	a.Encode([]RawAVP{setVendorID(vi), setAuthAppID(ai)})
	return
}

// GetVendorSpecAppID read Vendor-Specific-Application-Id AVP
func GetVendorSpecAppID(a RawAVP) (vi, ai uint32, e error) {
	o := []RawAVP{}
	if a.FlgV || !a.FlgM || a.FlgP {
		e = InvalidAVP(DiameterInvalidAvpBits)
	} else {
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
		e = InvalidAVP(DiameterMissingAvp)
	}
	return
}

// SetSessionID make Session-ID AVP
func SetSessionID(v string) (a RawAVP) {
	a = RawAVP{Code: 263, VenID: 0, FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

// GetSessionID read Session-ID AVP
func GetSessionID(a RawAVP) (v string, e error) {
	if a.FlgV || !a.FlgM || a.FlgP {
		e = InvalidAVP(DiameterInvalidAvpBits)
	} else {
		e = a.Decode(&v)
	}
	return
}

// SetOriginHost make Origin-Host AVP
func SetOriginHost(v Identity) (a RawAVP) {
	a = RawAVP{Code: 264, VenID: 0, FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

// GetOriginHost read Origin-Host AVP
func GetOriginHost(a RawAVP) (v Identity, e error) {
	if a.FlgV || !a.FlgM || a.FlgP {
		e = InvalidAVP(DiameterInvalidAvpBits)
	} else {
		e = a.Decode(&v)
	}
	return
}

// SetResultCode make Result-Code AVP
func SetResultCode(c uint32) (a RawAVP) {
	if c < 10000 {
		a = RawAVP{Code: 268, VenID: 0, FlgV: false, FlgM: true, FlgP: false}
		a.Encode(c)
		return
	}
	a = RawAVP{Code: 297, VenID: 0, FlgV: false, FlgM: true, FlgP: false}
	v := []RawAVP{
		RawAVP{Code: 266, VenID: 0, FlgV: false, FlgM: true, FlgP: false},
		RawAVP{Code: 298, VenID: 0, FlgV: false, FlgM: true, FlgP: false}}
	v[0].Encode(c / 10000)
	v[1].Encode(c % 10000)
	a.Encode(v)
	return
}

// GetResultCode read Result-Code AVP
func GetResultCode(a RawAVP) (c uint32, e error) {
	if a.FlgV || !a.FlgM || a.FlgP {
		e = InvalidAVP(DiameterInvalidAvpBits)
	} else if a.Code == 268 {
		e = a.Decode(&c)
		return
	} else if a.Code == 297 {
		o := []RawAVP{}
		e = a.Decode(&o)
		for _, a := range o {
			switch a.Code {
			case 266:
				if a.FlgV || !a.FlgM || a.FlgP {
					e = InvalidAVP(DiameterInvalidAvpBits)
				} else {
					var i uint32
					e = a.Decode(&i)
					c += i * 10000
				}
			case 298:
				if a.FlgV || !a.FlgM || a.FlgP {
					e = InvalidAVP(DiameterInvalidAvpBits)
				} else {
					var r uint32
					e = a.Decode(&r)
					c += r
				}
			}
		}
		if c < 10000 || c%10000 == 0 {
			e = InvalidAVP(DiameterMissingAvp)
		}
	}
	return
}

// SetAuthSessionState make Auth-Session-State AVP
func SetAuthSessionState(v bool) (a RawAVP) {
	a = RawAVP{Code: 277, VenID: 0, FlgV: false, FlgM: true, FlgP: false}
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
func GetAuthSessionState(a RawAVP) (v bool, e error) {
	s := new(Enumerated)
	if a.FlgV || !a.FlgM || a.FlgP {
		e = InvalidAVP(DiameterInvalidAvpBits)
	} else if e = a.Decode(s); e != nil {
		switch *s {
		case 0:
			v = true
		case 1:
			v = false
		default:
			e = InvalidAVP(DiameterInvalidAvpValue)
		}
	}
	return
}

// SetFailedAVP make Failed-AVP AVP
func SetFailedAVP(v []RawAVP) (a RawAVP) {
	a = RawAVP{Code: 279, VenID: 0, FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

// GetFailedAVP read Failed-AVP AVP
func GetFailedAVP(a RawAVP) (v []RawAVP, e error) {
	if a.FlgV || !a.FlgM || a.FlgP {
		e = InvalidAVP(DiameterInvalidAvpBits)
	} else {
		e = a.Decode(&v)
	}
	return
}

// SetRouteRecord make Route-Record AVP
func SetRouteRecord(v Identity) (a RawAVP) {
	a = RawAVP{Code: 282, VenID: 0, FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

// GetRouteRecord read Route-Record AVP
func GetRouteRecord(a RawAVP) (v Identity, e error) {
	if a.FlgV || !a.FlgM || a.FlgP {
		e = InvalidAVP(DiameterInvalidAvpBits)
	} else {
		e = a.Decode(&v)
	}
	return
}

// SetDestinationRealm make Destination-Realm AVP
func SetDestinationRealm(v Identity) (a RawAVP) {
	a = RawAVP{Code: 283, VenID: 0, FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

// GetDestinationRealm read Destination-Realm AVP
func GetDestinationRealm(a RawAVP) (v Identity, e error) {
	if a.FlgV || !a.FlgM || a.FlgP {
		e = InvalidAVP(DiameterInvalidAvpBits)
	} else {
		e = a.Decode(&v)
	}
	return
}

// SetDestinationHost make Destination-Host AVP
func SetDestinationHost(v Identity) (a RawAVP) {
	a = RawAVP{Code: 293, VenID: 0, FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

// GetDestinationHost read Destination-Host AVP
func GetDestinationHost(a RawAVP) (v Identity, e error) {
	if a.FlgV || !a.FlgM || a.FlgP {
		e = InvalidAVP(DiameterInvalidAvpBits)
	} else {
		e = a.Decode(&v)
	}
	return
}

// SetOriginRealm make Origin-Realm AVP
func SetOriginRealm(v Identity) (a RawAVP) {
	a = RawAVP{Code: 296, VenID: 0, FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

// GetOriginRealm read Origin-Realm AVP
func GetOriginRealm(a RawAVP) (v Identity, e error) {
	if a.FlgV || !a.FlgM || a.FlgP {
		e = InvalidAVP(DiameterInvalidAvpBits)
	} else {
		e = a.Decode(&v)
	}
	return
}

/*
// ProxyHost AVP
type ProxyHost diameter.Identity

// ToRaw return AVP struct of this value
func (v *ProxyHost) ToRaw() diameter.RawAVP {
	a := diameter.RawAVP{Code: 280, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(diameter.Identity(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *ProxyHost) FromRaw(a diameter.RawAVP) (e error) {
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
func (v *ProxyState) ToRaw() diameter.RawAVP {
	a := diameter.RawAVP{Code: 33, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode([]byte(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *ProxyState) FromRaw(a diameter.RawAVP) (e error) {
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
func (v *ProxyInfo) ToRaw() diameter.RawAVP {
	a := diameter.RawAVP{Code: 284, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	if v != nil {
		t := []diameter.RawAVP{
			v.ProxyHost.ToRaw(),
			v.ProxyState.ToRaw()}
		a.Encode(diameter.GroupedAVP(t))
	}
	return a
}

// FromRaw get AVP value
func (v *ProxyInfo) FromRaw(a diameter.RawAVP) (e error) {
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
