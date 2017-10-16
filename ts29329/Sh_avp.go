package ts29329

import (
	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/teldata"
)

// MSISDN AVP
type MSISDN teldata.TBCD

// ToRaw return AVP struct of this value
func (v *MSISDN) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 701, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode([]byte(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *MSISDN) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 701, true, true, false); e != nil {
		return
	}
	s := new([]byte)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = MSISDN(*s)
	return
}

// UserData AVP
type UserData []byte

// ToRaw return AVP struct of this value
func (v *UserData) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 702, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode([]byte(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *UserData) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 702, true, true, false); e != nil {
		return
	}
	s := new([]byte)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = UserData(*s)
	return
}
