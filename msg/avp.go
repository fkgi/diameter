package msg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// AVP value
type AVP interface {
	ToRaw() RawAVP
	FromRaw(RawAVP) (AVP, error)
}

// RawAVP is AVP data and header
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

// Validate header value
func (a RawAVP) Validate(i, c uint32, v, m, p bool) error {
	if a.VenID != i || a.Code != c {
		return InvalidAVP{}
	} else if a.FlgV != v || a.FlgM != m || a.FlgP != p {
		return InvalidAVP{}
	}
	return nil
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
	switch d := d.(type) {
	case nil:
		a.data = []byte{}
	case net.IP:
		e = a.setAddressData(d)
	case time.Time:
		e = a.setTimeData(d)
	case DiameterIdentity:
		e = a.setDiameterIdentityData(d)
	case DiameterURI:
		e = a.setDiameterURIData(d)
	case Enumerated:
		e = a.setEnumeratedData(d)
	case IPFilterRule:
		e = a.setIPFilterRuleData(d)
	case GroupedAVP:
		e = a.setGroupedData(d)
	case string:
		e = a.setUTF8StringData(d)
	case []byte:
		e = a.setOctetStringData(d)
	case int32:
		e = a.setInteger32Data(d)
	case int64:
		e = a.setInteger64Data(d)
	case uint32:
		e = a.setUnsigned32Data(d)
	case uint64:
		e = a.setUnsigned64Data(d)
	case float32:
		e = a.setFloat32Data(d)
	case float64:
		e = a.setFloat64Data(d)
	default:
		e = &UnknownAVPType{}
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
		e = a.getAddressData(d)
	case *time.Time:
		e = a.getTimeData(d)
	case *DiameterIdentity:
		e = a.getDiameterIdentityData(d)
	case *DiameterURI:
		e = a.getDiameterURIData(d)
	case *Enumerated:
		e = a.getEnumeratedData(d)
	case *IPFilterRule:
		e = a.getIPFilterRuleData(d)
	case *GroupedAVP:
		e = a.getGroupedData(d)
	case *string:
		e = a.getUTF8StringData(d)
	case *[]byte:
		e = a.getOctetStringData(d)
	case *int32:
		e = a.getInteger32Data(d)
	case *int64:
		e = a.getInteger64Data(d)
	case *uint32:
		e = a.getUnsigned32Data(d)
	case *uint64:
		e = a.getUnsigned64Data(d)
	case *float32:
		e = a.getFloat32Data(d)
	case *float64:
		e = a.getFloat64Data(d)
	default:
		e = &UnknownAVPType{}
	}
	return
}

// Basic Data Format AVP
// OctetString
func (a *RawAVP) setOctetStringData(d []byte) (e error) {
	a.data = make([]byte, len(d))
	copy(a.data, d)
	return
}

func (a RawAVP) getOctetStringData(d *[]byte) (e error) {
	*d = a.data
	return
}

// Integer32
func (a *RawAVP) setInteger32Data(d int32) (e error) {
	a.data, e = numConv(d)
	return
}

func (a RawAVP) getInteger32Data(d interface{}) (e error) {
	e = numRConv(a.data, 4, d)
	return
}

// Integer64
func (a *RawAVP) setInteger64Data(d int64) (e error) {
	a.data, e = numConv(d)
	return
}

func (a RawAVP) getInteger64Data(d *int64) (e error) {
	e = numRConv(a.data, 8, d)
	return
}

// Unsigned32
func (a *RawAVP) setUnsigned32Data(d uint32) (e error) {
	a.data, e = numConv(d)
	return
}

func (a RawAVP) getUnsigned32Data(d *uint32) (e error) {
	e = numRConv(a.data, 4, d)
	return
}

// Unsigned64
func (a *RawAVP) setUnsigned64Data(d uint64) (e error) {
	a.data, e = numConv(d)
	return
}

func (a RawAVP) getUnsigned64Data(d *uint64) (e error) {
	e = numRConv(a.data, 8, d)
	return
}

// Float32
func (a *RawAVP) setFloat32Data(d float32) (e error) {
	a.data, e = numConv(d)
	return
}

func (a RawAVP) getFloat32Data(d *float32) (e error) {
	e = numRConv(a.data, 4, d)
	return
}

// Float64
func (a *RawAVP) setFloat64Data(d float64) (e error) {
	a.data, e = numConv(d)
	return
}

func (a RawAVP) getFloat64Data(d *float64) (e error) {
	e = numRConv(a.data, 8, d)
	return
}

// GroupedAVP is Grouped format AVP value
type GroupedAVP []RawAVP

// Grouped
func (a *RawAVP) setGroupedData(d GroupedAVP) (e error) {
	buf := new(bytes.Buffer)
	for _, avp := range d {
		_, e = avp.WriteTo(buf)
	}
	if e == nil {
		a.data = buf.Bytes()
	}
	return
}

func (a RawAVP) getGroupedData(d *GroupedAVP) (e error) {
	*d = make([]RawAVP, 0)
	for buf := bytes.NewReader(a.data); buf.Len() != 0; {
		avp := RawAVP{}
		if _, e = avp.ReadFrom(buf); e != nil {
			break
		}
		*d = append(*d, avp)
	}
	return
}

// Common Derived AVP Data Formats
// Address
func (a *RawAVP) setAddressData(ip net.IP) (e error) {
	if ip.To4() != nil {
		ip = ip.To4()
		a.data = make([]byte, 6)
		a.data[0] = 0x00
		a.data[1] = 0x01
		for i := 2; i < 6; i++ {
			a.data[i] = ip[i-2]
		}
	} else if ip.To16() != nil {
		ip = ip.To16()
		a.data = make([]byte, 18)
		a.data[0] = 0x00
		a.data[1] = 0x02
		for i := 2; i < 18; i++ {
			a.data[i] = ip[i-2]
		}
	} else {
		e = fmt.Errorf("invalid net.IP struct")
	}
	return
}

func (a RawAVP) getAddressData(ip *net.IP) (e error) {
	if len(a.data) == 6 && a.data[0] == 0x00 && a.data[1] == 0x01 {
		*ip = net.IP(a.data[2:6])
	} else if len(a.data) == 18 && a.data[0] == 0x00 && a.data[1] == 0x02 {
		*ip = net.IP(a.data[2:18])
	} else {
		e = fmt.Errorf("invalid address family")
	}
	return
}

// Time
func (a *RawAVP) setTimeData(t time.Time) (e error) {
	a.data, e = numConv(int64(t.Unix() + 2208988800))
	return
}

func (a RawAVP) getTimeData(t *time.Time) (e error) {
	var d uint64
	if e = numRConv(a.data, 8, &d); e != nil {
		*t = time.Unix(int64(d-2208988800), int64(0))
	}
	return
}

// UTF8String
func (a *RawAVP) setUTF8StringData(s string) (e error) {
	a.data = []byte(strings.TrimSpace(s))
	return
}

func (a RawAVP) getUTF8StringData(s *string) (e error) {
	*s = string(a.data)
	return
}

// DiameterIdentity
func (a *RawAVP) setDiameterIdentityData(u DiameterIdentity) (e error) {
	a.data = []byte(u)
	return
}

func (a RawAVP) getDiameterIdentityData(u *DiameterIdentity) (e error) {
	*u, e = ParseDiameterIdentity(string(a.data))
	return
}

// DiameterURI
func (a *RawAVP) setDiameterURIData(u DiameterURI) (e error) {
	a.data = []byte(u.String())
	return
}

func (a RawAVP) getDiameterURIData(u *DiameterURI) (e error) {
	*u, e = ParseDiameterURI(string(a.data))
	return
}

// Enumerated is Enumerated format AVP value
type Enumerated int32

// Enumerated
func (a *RawAVP) setEnumeratedData(n Enumerated) (e error) {
	a.data, e = numConv(int32(n))
	return
}

func (a RawAVP) getEnumeratedData(n *Enumerated) (e error) {
	e = numRConv(a.data, 4, n)
	return
}

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
