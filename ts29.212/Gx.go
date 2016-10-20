package ts29212

import "github.com/fkgi/diameter/msg"

// RATType AVP
func RATType(e msg.Enumerated) msg.Avp {
	a := msg.Avp{Code: uint32(1032), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(e)
	return a
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
