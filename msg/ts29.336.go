package msg

// UserIdentifier AVP
func UserIdentifier(uname, msisdn, extid string, lmsi uint32) Avp {
	a := Avp{Code: uint32(3102), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []Avp
	if len(uname) != 0 {
		t = append(t, UserName(uname))
	}
	if len(msisdn) != 0 {
		t = append(t, MSISDN(msisdn))
	}
	if len(extid) != 0 {
		t = append(t, ExternalIdentifier(extid))
	}
	if lmsi != 0 {
		t = append(t, LMSI(lmsi))
	}
	a.Encode(t)
	return a
}

// ExternalIdentifier AVP
func ExternalIdentifier(extid string) Avp {
	a := Avp{Code: uint32(3111), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(extid)
	return a
}
