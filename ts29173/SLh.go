package ts29173

import (
	"bytes"
	"encoding/binary"

	"github.com/fkgi/diameter/msg"
)

// LMSI AVP
type LMSI uint32

// Encode return AVP struct of this value
func (v LMSI) Encode() msg.Avp {
	a := msg.Avp{Code: 2400, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(v))
	a.Encode(buf.Bytes())
	return a
}

// GetLMSI get AVP value
func GetLMSI(o msg.GroupedAVP) (LMSI, bool) {
	s := new([]byte)
	if a, ok := o.Get(2400, 10415); ok {
		a.Decode(s)
	} else {
		return 0, false
	}
	var i uint32
	binary.Read(bytes.NewBuffer(*s), binary.BigEndian, &i)
	return LMSI(i), true
}
