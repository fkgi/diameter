package diameter

import "net"

func setHostIPAddress(v net.IP) (a AVP) {
	a = AVP{Code: 257, Mandatory: true}
	a.Encode(v)
	return
}

func getHostIPAddress(a AVP) (v net.IP, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

func setFirmwareRevision(v uint32) (a AVP) {
	a = AVP{Code: 267}
	a.Encode(v)
	return
}

func getFirmwareRevision(a AVP) (v uint32, e error) {
	if a.VendorID != 0 || a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

func setProductName(v string) (a AVP) {
	a = AVP{Code: 269}
	a.Encode(v)
	return
}

func getProductName(a AVP) (v string, e error) {
	if a.VendorID != 0 || a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

const (
	// Rebooting is Enumerated value 0
	Rebooting Enumerated = 0
	// Busy is Enumerated value 1
	Busy Enumerated = 1
	// DoNotWantToTalkToYou is Enumerated value 2
	DoNotWantToTalkToYou Enumerated = 2
)

func setDisconnectCause(v Enumerated) (a AVP) {
	a = AVP{Code: 273, Mandatory: true}
	a.Encode(v)
	return
}

func getDisconnectCause(a AVP) (v Enumerated, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	if v < 0 || v > 2 {
		e = InvalidAVP{Code: InvalidAvpValue, AVP: a}
	}
	return
}

func setOriginStateID(v uint32) (a AVP) {
	a = AVP{Code: 278, Mandatory: true}
	a.Encode(v)
	return
}

func getOriginStateID(a AVP) (v uint32, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

func setFailedAVP(v []AVP) (a AVP) {
	a = AVP{Code: 279, Mandatory: true}
	a.Encode(v)
	return
}

func getFailedAVP(a AVP) (v []AVP, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

func setSupportedVendorID(v uint32) (a AVP) {
	a = AVP{Code: 265, Mandatory: true}
	a.Encode(v)
	return a
}

func getSupportedVendorID(a AVP) (v uint32, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

func setInbandSecurityID(v uint32) (a AVP) {
	a = AVP{Code: 299, Mandatory: true}
	a.Encode(v)
	return
}

func getInbandSecurityID(a AVP) (v uint32, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}

func setErrorMessage(v string) (a AVP) {
	a = AVP{Code: 281}
	a.Encode(v)
	return
}

func getErrorMessage(a AVP) (v string, e error) {
	if a.VendorID != 0 || !a.Mandatory {
		e = InvalidAVP{Code: InvalidAvpBits, AVP: a}
	} else {
		e = a.wrapedDecode(&v)
	}
	return
}
