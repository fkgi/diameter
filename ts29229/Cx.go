package ts29229

import "github.com/fkgi/diameter/msg"

const v3gpp uint32 = 10415

// SupportedFeatures AVP
type SupportedFeatures struct {
	VendorID      msg.VendorID
	FeatureListID FeatureListID
	FeatureList   FeatureList
}

// Encode return AVP struct of this value
func (v SupportedFeatures) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(628), VenID: v3gpp,
		FlgV: true, FlgM: true, FlgP: false}
	t := make([]msg.Avp, 3)

	// Vendor-Id
	t[0] = v.VendorID.Encode()
	// Feature-List-ID
	t[1] = v.FeatureListID.Encode()
	// Feature-List
	t[2] = v.FeatureList.Encode()

	a.Encode(msg.GroupedAVP(t))
	return a
}

// FeatureListID AVP
type FeatureListID uint32

// Encode return AVP struct of this value
func (v FeatureListID) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(629), VenID: v3gpp,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// FeatureList AVP
type FeatureList uint32

// Encode return AVP struct of this value
func (v FeatureList) Encode() msg.Avp {
	a := msg.Avp{Code: uint32(630), VenID: v3gpp,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode(uint32(v))
	return a
}
