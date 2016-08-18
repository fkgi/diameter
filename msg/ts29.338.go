package msg

// SCAddress AVP
func SCAddress(msisdn string) Avp {
	a := Avp{Code: uint32(3300), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(msisdn)
	return a
}

// OFRFlags AVP
func OFRFlags(s6as6d bool) Avp {
	a := Avp{Code: uint32(3328), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	i := uint32(0)

	if s6as6d {
		i = i | 0x00000001
	}

	a.Encode(i)
	return a
}

// SMRPUI AVP
func SMRPUI(s []byte) Avp {
	a := Avp{Code: uint32(3301), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(s)
	return a
}

/*
func NewAVP_SMSMICorrelationID() Avp {
	a := Avp{Code: uint32(3324), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	t := make([]Avp, 0)
	t = append(t, NewAVP_HSSID(id))
	t = append(t, NewAVP_OriginatingSIPURI(id))
	t = append(t, NewAVP_DestinationSIPURI(id))

	a.Encode(t)
	return a
}

func NewAVP_SMDeliveryOutcome() Avp {
	a := Avp{Code: uint32(3316), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	t := make([]Avp, 0)
	t = append(t, NewAVP_MMESMDeliveryOutcome(id))
	t = append(t, NewAVP_MSCSMDeliveryOutcome(id))
	t = append(t, NewAVP_SGSNSMDeliveryOutcome(id))
	t = append(t, NewAVP_IPSMGWSMDeliveryOutcome(id))

	a.Encode(t)
	return a
}
*/
