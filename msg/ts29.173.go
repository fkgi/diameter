package msg

// LMSI AVP
func LMSI(id uint32) Avp {
	a := Avp{Code: uint32(2400), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(itob(id))
	return a
}
