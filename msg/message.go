package msg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	// DiaVer is Diameter protocol Version
	DiaVer = uint8(1)
)

var (
	// Indent for String() output for message
	Indent = " | "
)

// Message is Diameter message
type Message struct {
	Ver   uint8  // Version = 1
	leng  uint32 // Message Length (24bit)
	FlgR  bool   // Request
	FlgP  bool   // Proxiable
	FlgE  bool   // Error
	FlgT  bool   // Potentially re-transmitted message
	Code  uint32 // Command-Code (24bit)
	AppID uint32 // Application-ID
	HbHID uint32 // Hop-by-Hop ID
	EtEID uint32 // End-to-End ID
	data  []byte
}

// PrintStack show message parameter
func (m Message) PrintStack(w io.Writer) {
	fmt.Fprintf(w, "Version       =%d\n", m.Ver)
	fmt.Fprintf(w, "Message Length=%d\n", m.leng)
	fmt.Fprintf(w, "Flags        R=%t, P=%t, E=%t, T=%t\n", m.FlgR, m.FlgP, m.FlgE, m.FlgT)
	fmt.Fprintf(w, "Command-Code  =%d\n", m.Code)
	fmt.Fprintf(w, "Application-ID=%d\n", m.AppID)
	fmt.Fprintf(w, "Hop-by-Hop ID =%d\n", m.HbHID)
	fmt.Fprintf(w, "End-to-End ID =%d\n", m.EtEID)

	if avp, e := m.Decode(); e == nil {
		for i, a := range avp {
			fmt.Fprintf(w, "AVP [%d]\n", i)
			a.PrintStack(w)
		}
	}
}

func (m Message) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "%sVersion       =%d\n", Indent, m.Ver)
	fmt.Fprintf(w, "%sMessage Length=%d\n", Indent, m.leng)
	fmt.Fprintf(w, "%sFlags        R=%t, P=%t, E=%t, T=%t\n",
		Indent, m.FlgR, m.FlgP, m.FlgE, m.FlgT)
	fmt.Fprintf(w, "%sCommand-Code  =%d\n", Indent, m.Code)
	fmt.Fprintf(w, "%sApplication-ID=%d\n", Indent, m.AppID)
	fmt.Fprintf(w, "%sHop-by-Hop ID =%d\n", Indent, m.HbHID)
	fmt.Fprintf(w, "%sEnd-to-End ID =%d", Indent, m.EtEID)

	if avp, e := m.Decode(); e == nil {
		for i, a := range avp {
			fmt.Fprintf(w, "\n%sAVP [%d]\n%s", Indent, i, a)
		}
	}

	return w.String()
}

// WriteTo write binary data to io.Writer
func (m Message) WriteTo(w io.Writer) (n int64, e error) {
	var b bytes.Buffer
	m.leng = uint32(20 + len(m.data))

	i := 0
	if i, e = b.Write([]byte{byte(m.Ver)}); e != nil {
		return
	}
	n += int64(i)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, m.leng)
	if i, e = b.Write(buf.Bytes()[1:4]); e != nil {
		return
	}
	n += int64(i)
	if i, e = b.Write(botob(m.FlgR, m.FlgP, m.FlgE, m.FlgT)); e != nil {
		return
	}
	n += int64(i)
	buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, m.Code)
	if i, e = b.Write(buf.Bytes()[1:4]); e != nil {
		return
	}
	n += int64(i)
	if e = binary.Write(&b, binary.BigEndian, m.AppID); e != nil {
		return
	}
	n += 4
	if e = binary.Write(&b, binary.BigEndian, m.HbHID); e != nil {
		return
	}
	n += 4
	if e = binary.Write(&b, binary.BigEndian, m.EtEID); e != nil {
		return
	}
	n += 4
	if i, e = b.Write(m.data); e != nil {
		return
	}
	n += int64(i)
	i, e = w.Write(b.Bytes())
	n += int64(i)

	return
}

// ReadFrom read binary data from io.Reader
func (m *Message) ReadFrom(r io.Reader) (n int64, e error) {
	i := 0
	buf, i, e := subread(r, 20)
	n += int64(i)
	if e != nil {
		return
	}
	m.Ver = buf[0]

	buf[0] = 0x00
	binary.Read(bytes.NewBuffer(buf[0:4]), binary.BigEndian, &m.leng)

	flgs := btobo(buf[4:5])
	m.FlgR = flgs[0]
	m.FlgP = flgs[1]
	m.FlgE = flgs[2]
	m.FlgT = flgs[3]

	buf[4] = 0x00
	binary.Read(bytes.NewBuffer(buf[4:8]), binary.BigEndian, &m.Code)

	binary.Read(bytes.NewBuffer(buf[8:12]), binary.BigEndian, &m.AppID)
	binary.Read(bytes.NewBuffer(buf[12:16]), binary.BigEndian, &m.HbHID)
	binary.Read(bytes.NewBuffer(buf[16:20]), binary.BigEndian, &m.EtEID)

	m.data, i, e = subread(r, int(m.leng)-20)
	n += int64(i)

	return
}

func subread(r io.Reader, l int) (buf []byte, o int, e error) {
	buf = make([]byte, l)
	i := 0
	for o < l {
		i, e = r.Read(buf[o:])
		o += i
		if e != nil {
			return
		}
	}
	return
}

// Encode convert AVP data to binary and set
func (m *Message) Encode(avp GroupedAVP) (e error) {
	buf := new(bytes.Buffer)
	for _, a := range avp {
		if _, e = a.WriteTo(buf); e != nil {
			return
		}
	}
	m.data = buf.Bytes()
	m.leng = uint32(20 + len(m.data))
	return
}

// Decode get and convert binary to AVP data
func (m Message) Decode() (avp GroupedAVP, e error) {
	avp = make([]Avp, 0)
	// l := m.leng - 20 + (4 - m.leng % 4) % 4

	buf := bytes.NewReader(m.data)
	for buf.Len() != 0 {
		a := Avp{}
		if _, e = a.ReadFrom(buf); e != nil {
			return
		}
		avp = append(avp, a)
	}
	return
}
