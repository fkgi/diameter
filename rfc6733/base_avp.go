package rfc6733

import (
	"net"

	"github.com/fkgi/diameter/msg"
)

func setHostIPAddress(v net.IP) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 257, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getHostIPAddress(a msg.RawAVP) (v net.IP, e error) {
	if e = a.Validate(0, 257, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setAuthAppID(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 258, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getAuthAppID(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(0, 258, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setVendorSpecAppID(vi, ai uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 260, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode([]msg.RawAVP{
		setVendorID(vi),
		setAuthAppID(ai)})
	return
}

func getVendorSpecAppID(a msg.RawAVP) (vi, ai uint32, e error) {
	o := []msg.RawAVP{}
	if e = a.Validate(0, 260, false, true, false); e == nil {
		e = a.Decode(&o)
	}
	for _, a := range o {
		if a.VenID != 0 {
			continue
		}
		switch a.Code {
		case 266:
			vi, e = getVendorID(a)
		case 258:
			ai, e = getAuthAppID(a)
		}
	}
	if vi == 0 || ai == 0 {
		e = msg.NoMandatoryAVP{}
	}
	return
}

func setOriginHost(v msg.DiameterIdentity) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 264, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getOriginHost(a msg.RawAVP) (v msg.DiameterIdentity, e error) {
	if e = a.Validate(0, 264, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setVendorID(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 266, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getVendorID(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(0, 266, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setFirmwareRevision(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 267, VenID: 0,
		FlgV: false, FlgM: false, FlgP: false}
	a.Encode(v)
	return
}

func getFirmwareRevision(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(0, 267, false, false, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setResultCode(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 268, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getResultCode(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(0, 268, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setProductName(v string) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 269, VenID: 0,
		FlgV: false, FlgM: false, FlgP: false}
	a.Encode(v)
	return
}

func getProductName(a msg.RawAVP) (v string, e error) {
	if e = a.Validate(0, 269, false, false, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

const (
	// Rebooting is Enumerated value 0
	Rebooting msg.Enumerated = 0
	// Busy is Enumerated value 1
	Busy msg.Enumerated = 1
	// DoNotWantToTalkToYou is Enumerated value 2
	DoNotWantToTalkToYou msg.Enumerated = 2
)

func setDisconnectCause(v msg.Enumerated) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 273, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getDisconnectCause(a msg.RawAVP) (v msg.Enumerated, e error) {
	if e = a.Validate(0, 273, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	if v < 0 || v > 2 {
		e = msg.InvalidAVP{}
	}
	return
}

func setOriginStateID(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 278, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getOriginStateID(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(0, 278, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setFailedAVP(v []msg.RawAVP) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 279, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getFailedAVP(a msg.RawAVP) (v []msg.RawAVP, e error) {
	if e = a.Validate(0, 279, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setSupportedVendorID(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 265, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return a
}

func getSupportedVendorID(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(0, 265, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setErrorMessage(v string) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 281, VenID: 0,
		FlgV: false, FlgM: false, FlgP: false}
	a.Encode(v)
	return
}

func getErrorMessage(a msg.RawAVP) (v string, e error) {
	if e = a.Validate(0, 281, false, false, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

func setOriginRealm(v msg.DiameterIdentity) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 296, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getOriginRealm(a msg.RawAVP) (v msg.DiameterIdentity, e error) {
	if e = a.Validate(0, 296, false, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}
