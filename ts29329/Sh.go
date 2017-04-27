package ts29329

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fkgi/diameter/msg"
)

// MSISDN AVP
type MSISDN []byte

// ParseMSISDN create TBCD value from string
func ParseMSISDN(s string) (MSISDN, error) {
	if strings.ContainsRune(s, '\x00') {
		return nil, fmt.Errorf("invalid charactor")
	} else if len(s)%2 != 0 {
		s = s + "\x00"
	}

	r := make([]byte, len(s)/2)
	for i, c := range s {
		var v byte
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
		case '\x00':
			v = 0x0f
		default:
			return r, fmt.Errorf("invalid charactor %c", c)
		}
		if i%2 == 1 {
			v = v << 4
		}
		r[i/2] = r[i/2] | v
	}
	return r, nil
}

// String return string value of the TBCD digit
func (v MSISDN) String() string {
	var b bytes.Buffer
	so := [2]byte{}
	for _, c := range v {
		so[0] = c & 0x0f
		so[1] = (c & 0xf0) >> 4
		for _, s := range so {
			switch s {
			case 0x00:
				b.WriteRune('0')
			case 0x01:
				b.WriteRune('1')
			case 0x02:
				b.WriteRune('2')
			case 0x03:
				b.WriteRune('3')
			case 0x04:
				b.WriteRune('4')
			case 0x05:
				b.WriteRune('5')
			case 0x06:
				b.WriteRune('6')
			case 0x07:
				b.WriteRune('7')
			case 0x08:
				b.WriteRune('8')
			case 0x09:
				b.WriteRune('9')
			case 0x0a:
				b.WriteRune('*')
			case 0x0b:
				b.WriteRune('#')
			case 0x0c:
				b.WriteRune('a')
			case 0x0d:
				b.WriteRune('b')
			case 0x0e:
				b.WriteRune('c')
			case 0x0f:
			}
		}
	}
	return b.String()
}

// Encode return AVP struct of this value
func (v MSISDN) Encode() msg.Avp {
	a := msg.Avp{Code: 701, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode([]byte(v))
	return a
}

// GetMSISDN get AVP value
func GetMSISDN(o msg.GroupedAVP) (MSISDN, bool) {
	s := new([]byte)
	if a, ok := o.Get(701, 10415); ok {
		a.Decode(s)
	} else {
		return nil, false
	}
	return MSISDN(*s), true
}

// UserData AVP
type UserData []byte

// Encode return AVP struct of this value
func (v UserData) Encode() msg.Avp {
	a := msg.Avp{Code: 702, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode([]byte(v))
	return a
}

// GetUserData get AVP value
func GetUserData(o msg.GroupedAVP) (UserData, bool) {
	s := new([]byte)
	if a, ok := o.Get(702, 10415); ok {
		a.Decode(s)
	} else {
		return nil, false
	}
	return UserData(*s), true
}
