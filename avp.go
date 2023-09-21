package diameter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

// Enumerated is Enumerated format AVP value.
type Enumerated int32

// IPFilterRule is IP Filter Rule format AVP value.
type IPFilterRule string

/*
AVP data and header

	 0                   1                   2                   3
	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|                           AVP Code                            |
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|V M P r r r r r|                  AVP Length                   |
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|                        Vendor-ID (opt)                        |
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|    Data ...
	+-+-+-+-+-+-+-+-+
*/
type AVP struct {
	Code      uint32 // AVP Code
	VendorID  uint32 // Vendor-ID
	Mandatory bool   // Mandatory AVP Flag
	// Protected bool // Protected AVP Flag
	Data []byte // AVP Data
}

func (a AVP) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "AVP Code      =%d\n", a.Code)
	fmt.Fprintf(w, "Flags        V=%t, M=%t, P=%t\n",
		a.VendorID != 0, a.Mandatory, false)
	if a.VendorID != 0 {
		fmt.Fprintf(w, "Vendor-ID     =%d\n", a.VendorID)
	}
	fmt.Fprintf(w, "Data          =% x", a.Data)
	return w.String()
}

// MarshalTo wite binary data to io.Writer
func (a AVP) MarshalTo(w io.Writer) error {
	var b bytes.Buffer

	binary.Write(&b, binary.BigEndian, a.Code)

	var flags byte
	if a.VendorID != 0 {
		flags |= 128
	}
	if a.Mandatory {
		flags |= 64
	}
	b.WriteByte(flags)

	lng := uint32(8 + len(a.Data))
	if a.VendorID != 0 {
		lng += 4
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, lng)
	b.Write(buf.Bytes()[1:4])

	if a.VendorID != 0 {
		binary.Write(&b, binary.BigEndian, a.VendorID)
	}
	b.Write(a.Data)

	b.Write(make([]byte, (4-len(a.Data)%4)%4))

	_, err := b.WriteTo(w)
	return err
}

func (a *AVP) wrapedUnmarshalFrom(r io.Reader) error {
	err := a.UnmarshalFrom(r)
	if err != nil {
		if _, ok := err.(InvalidAVP); !ok {
			err = InvalidAVP{Code: InvalidAvpValue, AVP: *a, E: err}
		}
	}
	return err
}

// UnmarshalFrom read binary data from io.Reader
func (a *AVP) UnmarshalFrom(r io.Reader) error {
	buf, err := readUntil(r, 8)
	if err != nil {
		return err
	}

	binary.Read(bytes.NewBuffer(buf[0:4]), binary.BigEndian, &a.Code)

	vid := buf[4]&128 == 128
	a.Mandatory = buf[4]&64 == 64
	pbit := buf[4]&32 == 32

	buf[4] = 0x00
	var lng uint32
	binary.Read(bytes.NewBuffer(buf[4:8]), binary.BigEndian, &lng)
	l := lng - 8

	if vid {
		if buf, err = readUntil(r, 4); err != nil {
			return err
		}
		binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &a.VendorID)
		l -= 4
	}

	if a.Data, err = readUntil(r, int(l)); err != nil {
		return err
	}
	if _, err = readUntil(r, (4-int(lng%4))%4); err != nil {
		return err
	}
	if pbit {
		return InvalidAVP{
			Code: InvalidAvpBits, AVP: *a,
			E: fmt.Errorf("p bit is not supported")}
	}
	return nil
}

// Encode make AVP from primitive go value
func (a *AVP) Encode(d interface{}) (e error) {
	buf := new(bytes.Buffer)

	switch d := d.(type) {
	case net.IP:
		buf.WriteByte(0x00)
		if d.To4() != nil {
			d = d.To4()
			buf.WriteByte(0x01)
			for i := 0; i < 4; i++ {
				buf.WriteByte(d[i])
			}
		} else if d.To16() != nil {
			d = d.To16()
			buf.WriteByte(0x02)
			for i := 0; i < 16; i++ {
				buf.WriteByte(d[i])
			}
		} else {
			e = fmt.Errorf("invalid net.IP struct")
		}
	case time.Time:
		e = binary.Write(buf, binary.BigEndian, int64(d.Unix()+2208988800))
	case Identity:
		buf.Write([]byte(d))
	case URI:
		buf.Write([]byte(d.String()))
	case Enumerated:
		e = binary.Write(buf, binary.BigEndian, int32(d))
	case IPFilterRule:
		buf.WriteString(string(d))
	case string:
		buf.WriteString(d)
	case []AVP:
		for _, avp := range d {
			if e = avp.MarshalTo(buf); e != nil {
				break
			}
		}
	case []byte:
		buf.Write(d)
	case int32, int64, uint32, uint64, float32, float64:
		e = binary.Write(buf, binary.BigEndian, d)
	case nil:
	default:
		e = fmt.Errorf("unacceptable type value for AVP")
	}

	if e == nil {
		a.Data = buf.Bytes()
	}
	return
}

func (a AVP) wrapedDecode(d interface{}) error {
	err := a.Decode(d)
	if err == io.EOF {
		err = InvalidAVP{Code: InvalidAvpLength, AVP: a}
	} else if err != nil {
		err = InvalidAVP{Code: InvalidAvpValue, AVP: a, E: err}
	}
	return err
}

// Decode make primitive go value from AVP
func (a AVP) Decode(d interface{}) (e error) {
	if a.Data == nil {
		d = nil
		return
	}

	switch d := d.(type) {
	case *net.IP:
		if a.Data[0] != 0x00 || a.Data[1] < 0x01 || a.Data[1] > 0x02 {
			e = fmt.Errorf("invalid address family")
		} else if len(a.Data) == 6 && a.Data[1] == 0x01 {
			*d = net.IP(a.Data[2:6])
		} else if len(a.Data) == 18 && a.Data[1] == 0x02 {
			*d = net.IP(a.Data[2:18])
		} else {
			e = io.EOF
		}
	case *time.Time:
		if len(a.Data) != 8 {
			e = io.EOF
		} else {
			buf := bytes.NewReader(a.Data)
			var t uint64
			if e = binary.Read(buf, binary.BigEndian, &t); e == nil {
				*d = time.Unix(int64(t-2208988800), int64(0))
			}
		}
	case *Identity:
		*d, e = ParseIdentity(string(a.Data))
	case *URI:
		*d, e = ParseURI(string(a.Data))
	case *Enumerated:
		if len(a.Data) != 4 {
			e = io.EOF
		} else {
			buf := bytes.NewReader(a.Data)
			var t Enumerated
			if e = binary.Read(buf, binary.BigEndian, &t); e == nil {
				*d = t
			}
		}
	case *IPFilterRule:
		*d = IPFilterRule(a.Data)
	case *string:
		*d = string(a.Data)
	case *[]AVP:
		*d = make([]AVP, 0)
		for buf := bytes.NewReader(a.Data); buf.Len() != 0; {
			avp := AVP{}
			if e = avp.UnmarshalFrom(buf); e != nil {
				break
			}
			*d = append(*d, avp)
		}
	case *[]byte:
		b := make([]byte, len(a.Data))
		copy(b, a.Data)
		*d = b
	case *int32, *uint32, *float32:
		if len(a.Data) != 4 {
			e = io.EOF
		} else {
			buf := bytes.NewReader(a.Data)
			e = binary.Read(buf, binary.BigEndian, d)
		}
	case *int64, *uint64, *float64:
		if len(a.Data) != 8 {
			e = io.EOF
		} else {
			buf := bytes.NewReader(a.Data)
			e = binary.Read(buf, binary.BigEndian, d)
		}
	default:
		e = fmt.Errorf("unacceptable type value for AVP")
	}

	return
}

func SetGenericAVP(c, v uint32, m bool, b interface{}) AVP {
	a := AVP{Code: c, VendorID: v, Mandatory: m}
	a.Encode(b)
	return a
}
