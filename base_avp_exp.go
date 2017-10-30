package diameter

// SetVendorSpecAppID make Vendor-Specific-Application-Id AVP
func SetVendorSpecAppID(vi, ai uint32) (a RawAVP) {
	a = RawAVP{Code: 260, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode([]RawAVP{
		setVendorID(vi),
		setAuthAppID(ai)})
	return
}

// GetVendorSpecAppID read Vendor-Specific-Application-Id AVP
func GetVendorSpecAppID(a RawAVP) (vi, ai uint32, e error) {
	o := []RawAVP{}
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
		e = NoMandatoryAVP{}
	}
	return
}

// SetSessionID make Session-ID AVP
func SetSessionID(v string) (a RawAVP) {
	a = RawAVP{Code: 263, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

// GetSessionID read Session-ID AVP
func GetSessionID(a RawAVP) (v string, e error) {
	if e = a.Validate(0, 263, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

// SetOriginHost make Origin-Host AVP
func SetOriginHost(v Identity) (a RawAVP) {
	a = RawAVP{Code: 264, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

// GetOriginHost read Origin-Host AVP
func GetOriginHost(a RawAVP) (v Identity, e error) {
	if e = a.Validate(0, 264, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

// SetResultCode make Result-Code AVP
func SetResultCode(v uint32) (a RawAVP) {
	a = RawAVP{Code: 268, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

// GetResultCode read Result-Code AVP
func GetResultCode(a RawAVP) (v uint32, e error) {
	if e = a.Validate(0, 268, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

// SetAuthSessionState make Auth-Session-State AVP
func SetAuthSessionState(v bool) (a RawAVP) {
	a = RawAVP{Code: 277, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
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
	if e = a.Validate(0, 277, false, true, false); e != nil {
	} else if e = a.Decode(s); e != nil {
		switch *s {
		case 0:
			v = true
		case 1:
			v = false
		default:
			e = InvalidAVP{}
		}
	}
	return
}

// SetDestinationRealm make Destination-Realm AVP
func SetDestinationRealm(v Identity) (a RawAVP) {
	a = RawAVP{Code: 283, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

// GetDestinationRealm read Destination-Realm AVP
func GetDestinationRealm(a RawAVP) (v Identity, e error) {
	if e = a.Validate(0, 283, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

// SetDestinationHost make Destination-Host AVP
func SetDestinationHost(v Identity) (a RawAVP) {
	a = RawAVP{Code: 293, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

// GetDestinationHost read Destination-Host AVP
func GetDestinationHost(a RawAVP) (v Identity, e error) {
	if e = a.Validate(0, 293, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

// SetOriginRealm make Origin-Realm AVP
func SetOriginRealm(v Identity) (a RawAVP) {
	a = RawAVP{Code: 296, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

// GetOriginRealm read Origin-Realm AVP
func GetOriginRealm(a RawAVP) (v Identity, e error) {
	if e = a.Validate(0, 296, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}
