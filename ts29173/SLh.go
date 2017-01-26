package ts29173

import "github.com/fkgi/diameter/msg"

// LMSI AVP
func LMSI(id uint32) msg.Avp {
	a := msg.Avp{Code: uint32(2400), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(msg.ItoB(id))
	return a
}
