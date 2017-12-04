package diameter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

// Enumerated is Enumerated format AVP value
type Enumerated int32

/*
RawAVP is AVP data and header
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
type RawAVP struct {
	Code  uint32 // AVP Code
	FlgV  bool   // Vendor Specific AVP Flag
	FlgM  bool   // Mandatory AVP Flag
	FlgP  bool   // Protected AVP Flag
	VenID uint32 // Vendor-ID
	data  []byte // AVP Data
}

func (a RawAVP) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "%s%sAVP Code      =%d\n", Indent, Indent, a.Code)
	fmt.Fprintf(w, "%s%sFlags        V=%t, M=%t, P=%t\n",
		Indent, Indent, a.FlgV, a.FlgM, a.FlgP)
	if a.FlgV {
		fmt.Fprintf(w, "%s%sVendor-ID     =%d\n", Indent, Indent, a.VenID)
	}
	fmt.Fprintf(w, "%s%sData          =% x", Indent, Indent, a.data)
	return w.String()
}

// WriteTo wite binary data to io.Writer
func (a RawAVP) WriteTo(w io.Writer) (n int64, e error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, a.Code)
	b.Write(botob(a.FlgV, a.FlgM, a.FlgP))

	lng := uint32(8 + len(a.data))
	if a.FlgV {
		lng += 4
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, lng)
	b.Write(buf.Bytes()[1:4])

	if a.FlgV {
		binary.Write(b, binary.BigEndian, a.VenID)
	}
	b.Write(a.data)

	b.Write(make([]byte, (4-len(a.data)%4)%4))

	i, e := w.Write(b.Bytes())
	return int64(i), e
}

// ReadFrom read binary data from io.Reader
func (a *RawAVP) ReadFrom(r io.Reader) (n int64, e error) {
	buf, i, e := subread(r, 8)
	n += int64(i)
	if e != nil {
		return
	}

	binary.Read(bytes.NewBuffer(buf[0:4]), binary.BigEndian, &a.Code)

	flgs := btobo(buf[4:5])
	a.FlgV = flgs[0]
	a.FlgM = flgs[1]
	a.FlgP = flgs[2]

	buf[4] = 0x00
	var lng uint32
	binary.Read(bytes.NewBuffer(buf[4:8]), binary.BigEndian, &lng)
	l := lng - 8

	if a.FlgV {
		buf, i, e = subread(r, 4)
		n += int64(i)
		if e != nil {
			return
		}
		binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &a.VenID)
		l -= 4
	}

	a.data, i, e = subread(r, int(l))
	n += int64(i)
	if e != nil {
		return
	}

	_, i, e = subread(r, (4-int(lng%4))%4)
	n += int64(i)

	return
}

// Encode make AVP from primitive go value
func (a *RawAVP) Encode(d interface{}) (e error) {
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
		_, e = buf.Write([]byte(d))
	case URI:
		_, e = buf.Write([]byte(d.String()))
	case Enumerated:
		e = binary.Write(buf, binary.BigEndian, int32(d))
		//	case IPFilterRule:
		//		e = a.setIPFilterRuleData(d)
	case string:
		_, e = buf.Write([]byte(d))
	case []RawAVP:
		for _, avp := range d {
			_, e = avp.WriteTo(buf)
		}
	case []byte:
		_, e = buf.Write(d)
	case int32, int64, uint32, uint64, float32, float64:
		e = binary.Write(buf, binary.BigEndian, d)
	case nil:
	default:
		e = &UnknownAVPType{}
	}
	if e == nil {
		a.data = buf.Bytes()
	}
	return
}

// Decode make primitive go value from AVP
func (a RawAVP) Decode(d interface{}) (e error) {
	if a.data == nil {
		d = nil
		return
	}
	switch d := d.(type) {
	case *net.IP:
		if len(a.data) == 6 && a.data[0] == 0x00 && a.data[1] == 0x01 {
			*d = net.IP(a.data[2:6])
		} else if len(a.data) == 18 && a.data[0] == 0x00 && a.data[1] == 0x02 {
			*d = net.IP(a.data[2:18])
		} else {
			e = fmt.Errorf("invalid address family")
		}
	case *time.Time:
		if len(a.data) != 8 {
			e = io.EOF
		} else {
			buf := bytes.NewReader(a.data)
			var t uint64
			if e = binary.Read(buf, binary.BigEndian, &t); e != nil {
				*d = time.Unix(int64(t-2208988800), int64(0))
			}
		}
	case *Identity:
		*d, e = ParseIdentity(string(a.data))
	case *URI:
		*d, e = ParseURI(string(a.data))
	case *Enumerated:
		if len(a.data) != 4 {
			e = io.EOF
		} else {
			buf := bytes.NewReader(a.data)
			var t int32
			if e = binary.Read(buf, binary.BigEndian, &t); e != nil {
				*d = Enumerated(t)
			}
		}
		//	case *IPFilterRule:
		//		e = a.getIPFilterRuleData(d)
	case *string:
		*d = string(a.data)
	case *[]RawAVP:
		*d = make([]RawAVP, 0)
		for buf := bytes.NewReader(a.data); buf.Len() != 0; {
			avp := RawAVP{}
			if _, e = avp.ReadFrom(buf); e != nil {
				break
			}
			*d = append(*d, avp)
		}
	case *[]byte:
		b := make([]byte, len(a.data))
		copy(b, a.data)
		*d = b
	case *int32, *uint32, *float32:
		if len(a.data) != 4 {
			e = io.EOF
		} else {
			buf := bytes.NewReader(a.data)
			e = binary.Read(buf, binary.BigEndian, d)
		}
	case *int64, *uint64, *float64:
		if len(a.data) != 8 {
			e = io.EOF
		} else {
			buf := bytes.NewReader(a.data)
			e = binary.Read(buf, binary.BigEndian, d)
		}
	default:
		e = &UnknownAVPType{}
	}
	return
}

/*
// IPFilterRule is IP Filter Rule format AVP value
type IPFilterRule string

// IPFilterRule
func (a *RawAVP) setIPFilterRuleData(s IPFilterRule) (e error) {
	a.data = []byte(s)
	return
}

func (a RawAVP) getIPFilterRuleData(s *IPFilterRule) (e error) {
	*s = IPFilterRule(a.data)
	return
}
*/
