package ts29329

import "github.com/fkgi/diameter/msg"

const v3gpp uint32 = 10415

// MSISDN AVP
type MSISDN string

// Encode return AVP struct of this value
func (v MSISDN) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(701), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(msg.StoTbcd(string(v)))
	return a
}

// UserData AVP
type UserData []byte

// Encode return AVP struct of this value
func (v UserData) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(702), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode([]byte(v))
	return a
}
