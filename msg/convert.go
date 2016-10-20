package msg

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// ItoB convert from uint32 to byte array
func ItoB(i uint32) (b []byte) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, i)
	return buf.Bytes()
}

// Convert from byte array to uint32
func btoi(b []byte) (i uint32) {
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, &i)
	return
}

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

// StoTbcd convert from string to TBCD string byte array
func StoTbcd(s string) []byte {
	if len(s)%2 != 0 {
		s = s + " "
	}
	r := make([]byte, len(s)/2)
	for i, c := range s {
		v := byte(0x0f)
		switch c {
		case '0':
			v = 0x00
		case '1':
			v = 0x01
		case '2':
			v = 0x02
		case '3':
			v = 0x03
		case '4':
			v = 0x04
		case '5':
			v = 0x05
		case '6':
			v = 0x06
		case '7':
			v = 0x07
		case '8':
			v = 0x08
		case '9':
			v = 0x09
		case '*':
			v = 0x0a
		case '#':
			v = 0x0b
		case 'a', 'A':
			v = 0x0c
		case 'b', 'B':
			v = 0x0d
		case 'c', 'C':
			v = 0x0e
		default:
			v = 0x0f
		}
		if i%2 == 1 {
			v = v << 4
		}
		r[i/2] = r[i/2] | v
	}
	return r
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
		e = fmt.Errorf("invalid data size")
	} else {
		buf := bytes.NewReader(b)
		e = binary.Read(buf, binary.BigEndian, v)
	}
	return
}
