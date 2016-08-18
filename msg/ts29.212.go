package msg

// RATType AVP
func RATType(e Enumerated) Avp {
	a := Avp{Code: uint32(1032), FlgV: true, FlgM: true, FlgP: false, VenID: uint32(10415)}
	a.Encode(e)
	return a
}

const (
	RATType_WLAN            Enumerated = 0
	RATType_VIRTUAL         Enumerated = 1
	RATType_UTRAN           Enumerated = 1000
	RATType_GERAN           Enumerated = 1001
	RATType_GAN             Enumerated = 1002
	RATType_HSDPA_EVOLUTION Enumerated = 1003
	RATType_EUTRAN          Enumerated = 1004
	RATType_CDMA2000_1X     Enumerated = 2000
	RATType_HRPD            Enumerated = 2001
	RATType_UMB             Enumerated = 2002
	RATType_EHRPD           Enumerated = 2003
)
