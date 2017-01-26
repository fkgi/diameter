package ts29336

import (
	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/ts29173"
	"github.com/fkgi/diameter/ts29329"
)

// UserIdentifier AVP
func UserIdentifier(uname, msisdn, extid string, lmsi uint32) msg.Avp {
	a := msg.Avp{Code: uint32(3102), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	var t []msg.Avp
	if len(uname) != 0 {
		t = append(t, msg.UserName(uname))
	}
	if len(msisdn) != 0 {
		t = append(t, ts29329.MSISDN(msisdn))
	}
	if len(extid) != 0 {
		t = append(t, ExternalIdentifier(extid))
	}
	if lmsi != 0 {
		t = append(t, ts29173.LMSI(lmsi))
	}
	a.Encode(t)
	return a
}

// ExternalIdentifier AVP
func ExternalIdentifier(extid string) msg.Avp {
	a := msg.Avp{Code: uint32(3111), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(extid)
	return a
}
