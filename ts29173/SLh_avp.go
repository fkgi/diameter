package ts29173

import (
	"bytes"
	"encoding/binary"

	"github.com/fkgi/teldata"

	"github.com/fkgi/diameter/msg"
)

// LMSI AVP
type LMSI uint32

// ToRaw return AVP struct of this value
func (v *LMSI) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 2400, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.BigEndian, uint32(*v))
		a.Encode(buf.Bytes())
	}
	return a
}

// FromRaw get AVP value
func (v *LMSI) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 2400, true, true, false); e != nil {
		return
	}
	s := new([]byte)
	if e = a.Decode(s); e != nil || len(*s) != 4 {
		return
	}
	var i uint32
	binary.Read(bytes.NewBuffer(*s), binary.BigEndian, &i)
	*v = LMSI(i)
	return
}

// MMEName AVP
type MMEName msg.DiameterIdentity

type MSCNumber teldata.TBCD
type MMERealm msg.DiameterIdentity
type SGSNName msg.DiameterIdentity
type SGSNRealm msg.DiameterIdentity
