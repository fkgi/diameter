package ts29336

import (
	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/ts29173"
	"github.com/fkgi/diameter/ts29329"
)

const v3gpp uint32 = 10415

// UserIdentifier AVP
type UserIdentifier struct {
	UserName msg.UserName
	MSISDN   ts29329.MSISDN
	ExtID    ExternalIdentifier
	LMSI     ts29173.LMSI
}

// Encode return AVP struct of this value
func (v UserIdentifier) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(3102), VenID: v3gpp,
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

// ExternalIdentifier AVP
type ExternalIdentifier string

// Encode return AVP struct of this value
func (v ExternalIdentifier) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(3111), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(string(v))
	return a
}
