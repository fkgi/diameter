package ts29329

import "github.com/fkgi/diameter/msg"

// MSISDN AVP
func MSISDN(msisdn string) msg.Avp {
	a := msg.Avp{Code: uint32(701), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(msg.StoTbcd(msisdn))
	return a
}

// UserData AVP
func UserData(data []byte) msg.Avp {
	a := msg.Avp{Code: uint32(702), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(data)
	return a
}
