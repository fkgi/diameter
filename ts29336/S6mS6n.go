package ts29336

import (
	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/ts29173"
	"github.com/fkgi/diameter/ts29329"
)

// UserIdentifier AVP
type UserIdentifier struct {
	UserName msg.UserName
	MSISDN   ts29329.MSISDN
	ExtID    ExternalIdentifier
	LMSI     ts29173.LMSI
}

// Encode return AVP struct of this value
func (v UserIdentifier) Encode() msg.Avp {
	a := msg.Avp{Code: 3102, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	var t []msg.Avp
	if len(v.UserName) != 0 {
		t = append(t, v.UserName.Encode())
	}
	if len(v.MSISDN) != 0 {
		t = append(t, v.MSISDN.Encode())
	}
	if len(v.ExtID) != 0 {
		t = append(t, v.ExtID.Encode())
	}
	if v.LMSI != 0 {
		t = append(t, v.LMSI.Encode())
	}
	a.Encode(msg.GroupedAVP(t))
	return a
}

// GetUserIdentifier get AVP value
func GetUserIdentifier(o msg.GroupedAVP) (UserIdentifier, bool) {
	s := UserIdentifier{}
	if a, ok := o.Get(3102, 10415); ok {
		o = msg.GroupedAVP{}
		a.Decode(&o)
	} else {
		return s, false
	}
	if t, ok := msg.GetUserName(o); ok {
		s.UserName = t
	}
	if t, ok := ts29329.GetMSISDN(o); ok {
		s.MSISDN = t
	}
	if t, ok := GetExternalIdentifier(o); ok {
		s.ExtID = t
	}
	if t, ok := ts29173.GetLMSI(o); ok {
		s.LMSI = t
	}
	return s, true
}

// ExternalIdentifier AVP
type ExternalIdentifier string

// Encode return AVP struct of this value
func (v ExternalIdentifier) Encode() msg.Avp {
	a := msg.Avp{Code: 3111, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(string(v))
	return a
}

// GetExternalIdentifier get AVP value
func GetExternalIdentifier(o msg.GroupedAVP) (ExternalIdentifier, bool) {
	s := new(string)
	if a, ok := o.Get(3111, 10415); ok {
		a.Decode(s)
	} else {
		return "", false
	}
	return ExternalIdentifier(*s), true
}
