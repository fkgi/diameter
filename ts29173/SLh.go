package ts29173

import "github.com/fkgi/diameter/msg"

const v3gpp uint32 = 10415

// LMSI AVP
type LMSI uint32

// Encode return AVP struct of this value
func (v LMSI) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(2400), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(msg.ItoB(uint32(v)))
	return a
}
