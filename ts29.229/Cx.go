package ts29229

import "github.com/fkgi/diameter/msg"

// SupportedFeatures AVP
func SupportedFeatures(vendorID, featureID, featureList uint32) msg.Avp {
	a := msg.Avp{Code: uint32(628), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	t := make([]msg.Avp, 3)

	// Vendor-Id
	t[0] = msg.VendorID(vendorID)
	// Feature-List-ID
	t[1] = FeatureListID(featureID)
	// Feature-List
	t[2] = FeatureList(featureList)

	a.Encode(t)
	return a
}

// FeatureListID AVP
func FeatureListID(i uint32) msg.Avp {
	a := msg.Avp{Code: uint32(629), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}

// FeatureList AVP
func FeatureList(i uint32) msg.Avp {
	a := msg.Avp{Code: uint32(630), FlgV: true, FlgM: false, FlgP: false, VenID: uint32(10415)}
	a.Encode(i)
	return a
}
