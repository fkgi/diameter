package ts29212

import "github.com/fkgi/diameter/msg"

// RATType AVP
type RATType msg.Enumerated

// Encode return AVP struct of this value
func (v RATType) Encode() msg.Avp {
	a := msg.Avp{Code: 1032, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(msg.Enumerated(v))
	return a
}

// GetRATType get AVP value
func GetRATType(o msg.GroupedAVP) (RATType, bool) {
	s := new(msg.Enumerated)
	if a, ok := o.Get(1032, 10415); ok {
		a.Decode(s)
	} else {
		return 0, false
	}
	return RATType(*s), true
}

const (
	// Wlan is Enumerated value 0
	Wlan msg.Enumerated = 0
	// Virtual is Enumerated value 1
	Virtual msg.Enumerated = 1
	// Utran is Enumerated value 1000
	Utran msg.Enumerated = 1000
	// GEran is Enumerated value 1001
	GEran msg.Enumerated = 1001
	// Gan is Enumerated value 1002
	Gan msg.Enumerated = 1002
	// HsdpaEvolution is Enumerated value 1003
	HsdpaEvolution msg.Enumerated = 1003
	// EUtran is Enumerated value 1004
	EUtran msg.Enumerated = 1004
	// Cdma20001x is Enumerated value 2000
	Cdma20001x msg.Enumerated = 2000
	// Hrpd is Enumerated value 2001
	Hrpd msg.Enumerated = 2001
	// Umb is Enumerated value 2002
	Umb msg.Enumerated = 2002
	// EHrpd is Enumerated value 2003
	EHrpd msg.Enumerated = 2003
)
