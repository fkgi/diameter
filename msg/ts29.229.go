package msg

// SupportedFeatures AVP
func SupportedFeatures(vendorID, featureID, featureList uint32) Avp {
	a := Avp{Code: uint32(628), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	t := make([]Avp, 3)

	// Vendor-Id
	t[0] = VendorID(vendorID)
	// Feature-List-ID
	t[1] = Avp{Code: uint32(629), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	t[1].Encode(featureID)
	// Feature-List
	t[2] = Avp{Code: uint32(630), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	t[2].Encode(featureList)

	a.Encode(t)
	return a
}
