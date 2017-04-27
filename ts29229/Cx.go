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

// GetSupportedFeatures get AVP value
func GetSupportedFeatures(o msg.GroupedAVP) (SupportedFeatures, bool) {
	s := SupportedFeatures{}
	if a, ok := o.Get(628, 10415); ok {
		o = msg.GroupedAVP{}
		a.Decode(&o)
	} else {
		return s, false
	}
	if t, ok := msg.GetVendorID(o); ok {
		s.VendorID = t
	}
	if t, ok := GetFeatureListID(o); ok {
		s.FeatureListID = t
	}
	if t, ok := GetFeatureList(o); ok {
		s.FeatureList = t
	}
	return s, true
}

// FeatureListID AVP
type FeatureListID uint32

// Encode return AVP struct of this value
func (v FeatureListID) Encode() msg.Avp {
	a := msg.Avp{Code: 629, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetFeatureListID get AVP value
func GetFeatureListID(o msg.GroupedAVP) (FeatureListID, bool) {
	s := new(uint32)
	if a, ok := o.Get(629, 10415); ok {
		a.Decode(s)
	} else {
		return 0, false
	}
	return FeatureListID(*s), true
}

// FeatureList AVP
type FeatureList uint32

// Encode return AVP struct of this value
func (v FeatureList) Encode() msg.Avp {
	a := msg.Avp{Code: 630, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	a.Encode(uint32(v))
	return a
}

// GetFeatureList get AVP value
func GetFeatureList(o msg.GroupedAVP) (FeatureList, bool) {
	s := new(uint32)
	if a, ok := o.Get(630, 10415); ok {
		a.Decode(s)
	} else {
		return 0, false
	}
	return FeatureList(*s), true
}
