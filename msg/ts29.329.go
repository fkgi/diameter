package msg

// MSISDN AVP
func MSISDN(msisdn string) Avp {
	a := Avp{Code: uint32(701), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(stotbcd(msisdn))
	return a
}

// UserData AVP
func UserData(data []byte) Avp {
	a := Avp{Code: uint32(702), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(data)
	return a
}
