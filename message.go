package diameter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	// DiaVer is Diameter protocol Version
	DiaVer uint8 = 1
)

var (
	// Indent for String() output for RawMsg
	Indent = " | "
)

// Request is Diameter request
type Request interface {
	ToRaw(string) RawMsg                     // generate RawMsg with session-id
	FromRaw(RawMsg) (Request, string, error) // decode RawMsg and return session-id
	Failed(uint32) Answer
}

// Answer is Diameter answer
type Answer interface {
	ToRaw(string) RawMsg
	FromRaw(RawMsg) (Answer, string, error)
	Result() uint32
}

/*
RawMsg is Diameter message.
    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |    Version    |                 Message Length                |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   | Command Flags |                  Command Code                 |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                         Application-ID                        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                      Hop-by-Hop Identifier                    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                      End-to-End Identifier                    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |  AVPs ...
   +-+-+-+-+-+-+-+-+-+-+-+-+-
*/
type RawMsg struct {
	Ver   uint8  // Version = 1
	FlgR  bool   // Request
	FlgP  bool   // Proxiable
	FlgE  bool   // Error
	FlgT  bool   // Potentially re-transmitted message
	Code  uint32 // Command-Code (24bit)
	AppID uint32 // Application-ID
	HbHID uint32 // Hop-by-Hop ID
	EtEID uint32 // End-to-End ID
	AVP   []RawAVP
}

func (m RawMsg) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sVersion       =%d\n", Indent, m.Ver)
	fmt.Fprintf(w, "%sFlags        R=%t, P=%t, E=%t, T=%t\n",
		Indent, m.FlgR, m.FlgP, m.FlgE, m.FlgT)
	fmt.Fprintf(w, "%sCommand-Code  =%d\n", Indent, m.Code)
	fmt.Fprintf(w, "%sApplication-ID=%d\n", Indent, m.AppID)
	fmt.Fprintf(w, "%sHop-by-Hop ID =%d\n", Indent, m.HbHID)
	fmt.Fprintf(w, "%sEnd-to-End ID =%d", Indent, m.EtEID)
	for i, a := range m.AVP {
		fmt.Fprintf(w, "\n%sAVP [%d]\n%s", Indent, i, a)
	}

	return w.String()
}

// Validate header value
func (m RawMsg) Validate(i, c uint32, r, p, e, t bool) error {
	if m.AppID != i || m.Code != c {
		return InvalidMessage{}
	} else if m.FlgR != r {
		return InvalidMessage{}
	} else if r == true && m.FlgE == true {
		return InvalidMessage{}
	} else if r == false && m.FlgT == true {
		return InvalidMessage{}
	}
	return nil
}

// Clone make copy of this RawMsg
func (m RawMsg) Clone() RawMsg {
	avp := make([]RawAVP, len(m.AVP))
	for i := range m.AVP {
		avp[i] = RawAVP{
			Code:  m.AVP[i].Code,
			FlgV:  m.AVP[i].FlgV,
			FlgM:  m.AVP[i].FlgM,
			FlgP:  m.AVP[i].FlgP,
			VenID: m.AVP[i].VenID,
			data:  make([]byte, len(m.AVP[i].data))}
		copy(avp[i].data, m.AVP[i].data)
	}
	return RawMsg{
		Ver:   m.Ver,
		FlgR:  m.FlgR,
		FlgP:  m.FlgP,
		FlgE:  m.FlgE,
		FlgT:  m.FlgT,
		Code:  m.Code,
		AppID: m.AppID,
		HbHID: m.HbHID,
		EtEID: m.EtEID,
		AVP:   avp}
}

// WriteTo write binary data to io.Writer
func (m RawMsg) WriteTo(w io.Writer) (n int64, e error) {
	var b, dat bytes.Buffer

	for _, a := range m.AVP {
		if _, e = a.WriteTo(&dat); e != nil {
			return
		}
	}
	lng := uint32(20 + dat.Len())

	b.Write([]byte{byte(m.Ver)})

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, lng)
	b.Write(buf.Bytes()[1:4])

	b.Write(botob(m.FlgR, m.FlgP, m.FlgE, m.FlgT))

	buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, m.Code)
	b.Write(buf.Bytes()[1:4])

	binary.Write(&b, binary.BigEndian, m.AppID)
	binary.Write(&b, binary.BigEndian, m.HbHID)
	binary.Write(&b, binary.BigEndian, m.EtEID)

	b.Write(dat.Bytes())

	i, e := w.Write(b.Bytes())
	return int64(i), e
}

// ReadFrom read binary data from io.Reader
func (m *RawMsg) ReadFrom(r io.Reader) (n int64, e error) {
	buf, i, e := subread(r, 20)
	n += int64(i)
	if e != nil {
		return
	}
	m.Ver = buf[0]

	buf[0] = 0x00
	var lng uint32
	binary.Read(bytes.NewBuffer(buf[0:4]), binary.BigEndian, &lng)
	// l := m.leng - 20 + (4 - m.leng % 4) % 4

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

	buf, i, e = subread(r, int(lng)-20)
	n += int64(i)
	if e != nil {
		return
	}

	m.AVP = []RawAVP{}
	rdr := bytes.NewReader(buf)
	for rdr.Len() != 0 {
		a := RawAVP{}
		if _, e = a.ReadFrom(rdr); e != nil {
			return
		}
		m.AVP = append(m.AVP, a)
	}

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
