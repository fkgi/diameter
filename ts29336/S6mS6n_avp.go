package ts29336

import (
	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/rfc6733"
	"github.com/fkgi/diameter/ts29173"
	"github.com/fkgi/diameter/ts29329"
	"github.com/fkgi/teldata"
)

// UserIdentifier AVP
type UserIdentifier struct {
	rfc6733.UserName
	ts29329.MSISDN
	ExternalIdentifier
	ts29173.LMSI
}

// ToRaw return AVP struct of this value
func (v *UserIdentifier) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3102, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		t := make([]msg.RawAVP, 0, 4)
		if len(v.UserName) != 0 {
			t = append(t, v.UserName.ToRaw())
		}
		if v.MSISDN != nil || len(v.MSISDN) != 0 {
			t = append(t, v.MSISDN.ToRaw())
		}
		if len(v.ExternalIdentifier) != 0 {
			t = append(t, v.ExternalIdentifier.ToRaw())
		}
		if v.LMSI != 0 {
			t = append(t, v.LMSI.ToRaw())
		}
		a.Encode(msg.GroupedAVP(t))
	}
	return a
}

// FromRaw get AVP value
func (v *UserIdentifier) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3102, true, true, false); e != nil {
		return
	}
	o := msg.GroupedAVP{}
	if e = a.Decode(&o); e != nil {
		return
	}
	*v = UserIdentifier{}
	for _, a := range o {
		if a.Code == 1 && a.VenID == 0 {
			e = v.UserName.FromRaw(a)
		} else if a.Code == 701 && a.VenID == 10415 {
			e = v.MSISDN.FromRaw(a)
		} else if a.Code == 3111 && a.VenID == 10415 {
			e = v.ExternalIdentifier.FromRaw(a)
		} else if a.Code == 2400 && a.VenID == 10415 {
			e = v.LMSI.FromRaw(a)
		}
	}
	return
}

// ExternalIdentifier AVP
type ExternalIdentifier string

// ToRaw return AVP struct of this value
func (v *ExternalIdentifier) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3111, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(string(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *ExternalIdentifier) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3111, true, true, false); e != nil {
		return
	}
	s := new(string)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = ExternalIdentifier(*s)
	return
}

// IPSMGWNumber AVP
type IPSMGWNumber teldata.TBCD

// ToRaw return AVP struct of this value
func (v *IPSMGWNumber) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3100, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode([]byte(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *IPSMGWNumber) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3100, true, true, false); e != nil {
		return
	}
	s := new([]byte)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = IPSMGWNumber(*s)
	return
}

// IPSMGWName AVP
type IPSMGWName msg.DiameterIdentity

// ToRaw return AVP struct of this value
func (v *IPSMGWName) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3101, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(msg.DiameterIdentity(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *IPSMGWName) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3101, true, true, false); e != nil {
		return
	}
	s := new(msg.DiameterIdentity)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = IPSMGWName(*s)
	return
}

// IPSMGWRealm AVP
type IPSMGWRealm msg.DiameterIdentity

// ToRaw return AVP struct of this value
func (v *IPSMGWRealm) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3112, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(msg.DiameterIdentity(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *IPSMGWRealm) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3112, true, true, false); e != nil {
		return
	}
	s := new(msg.DiameterIdentity)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = IPSMGWRealm(*s)
	return
}
