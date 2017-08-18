package msg

import (
	"bytes"
	"encoding/binary"
	"io"
)

// Convert from bool list to byte array
func botob(fs ...bool) (b []byte) {
	b = make([]byte, len(fs)/8+1)
	for i, f := range fs {
		if f {
			b[i/8] |= (0x80 >> uint(i%8))
		}
	}
	return
}

// Convert from byte to bool list
func btobo(b []byte) (fs []bool) {
	fs = make([]bool, 8*len(b))
	for i := range fs {
		fs[i] = (b[i/8]>>uint(7-i%8))&0x01 == 0x01
	}
	return
}

// subfunction of SetData(v interface{})(e error)
func numConv(v interface{}) (b []byte, e error) {
	buf := new(bytes.Buffer)
	e = binary.Write(buf, binary.BigEndian, v)
	if e == nil {
		b = buf.Bytes()
	}
	return
}

func numRConv(b []byte, s int, v interface{}) (e error) {
	if len(b) != s {
		e = io.EOF
	} else {
		buf := bytes.NewReader(b)
		e = binary.Read(buf, binary.BigEndian, v)
	}
	return
}
