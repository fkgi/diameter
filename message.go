package diameter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

/*
Message is Diameter message.

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
type Message struct {
	FlgR  bool   // Request
	FlgP  bool   // Proxiable
	FlgE  bool   // Error
	FlgT  bool   // Potentially re-transmitted message
	Code  uint32 // Command-Code (24bit)
	AppID uint32 // Application-ID
	HbHID uint32 // Hop-by-Hop ID
	EtEID uint32 // End-to-End ID
	AVPs  []byte // Message body AVP binary data
}

func (m *Message) setAVP(avp []AVP) {
	buf := new(bytes.Buffer)
	for _, a := range avp {
		a.MarshalTo(buf)
	}
	m.AVPs = buf.Bytes()
}

func (m *Message) getAVP() ([]AVP, error) {
	avp := make([]AVP, 0, avpBufferSize)
	for rdr := bytes.NewReader(m.AVPs); rdr.Len() != 0; {
		a := AVP{}
		if e := a.wrapedUnmarshalFrom(rdr); e != nil {
			return nil, e
		}
		avp = append(avp, a)
	}
	return avp, nil
}

func (m Message) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "Flags        R=%t, P=%t, E=%t, T=%t\n",
		m.FlgR, m.FlgP, m.FlgE, m.FlgT)
	fmt.Fprintf(w, "Command-Code  =%d\n", m.Code)
	fmt.Fprintf(w, "Application-ID=%d\n", m.AppID)
	fmt.Fprintf(w, "Hop-by-Hop ID =%d\n", m.HbHID)
	fmt.Fprintf(w, "End-to-End ID =%d", m.EtEID)

	for rdr := bytes.NewReader(m.AVPs); rdr.Len() != 0; {
		a := AVP{}
		if e := a.UnmarshalFrom(rdr); e != nil {
			break
		}
		fmt.Fprintf(w, "\nAVP [%d]      =%x", a.Code, a.Data)
	}

	return w.String()
}

func (m Message) generateAnswerBy(result uint32) Message {
	buf := new(bytes.Buffer)
	SetResultCode(result).MarshalTo(buf)
	SetOriginHost(Local.Host).MarshalTo(buf)
	SetOriginRealm(Local.Realm).MarshalTo(buf)

	return Message{
		FlgR: false, FlgP: m.FlgP, FlgE: true, FlgT: false,
		Code: m.Code, AppID: m.AppID,
		HbHID: m.HbHID, EtEID: m.EtEID,
		AVPs: buf.Bytes()}
}

// MarshalTo write binary data to io.Writer
func (m Message) MarshalTo(w io.Writer) error {
	var b, buf bytes.Buffer

	b.WriteByte(1)

	binary.Write(&buf, binary.BigEndian, uint32(20+len(m.AVPs)))
	b.Write(buf.Bytes()[1:4])

	var flags byte
	if m.FlgR {
		flags |= 0x80
	}
	if m.FlgP {
		flags |= 0x40
	}
	if m.FlgE {
		flags |= 0x20
	}
	if m.FlgT {
		flags |= 0x10
	}
	b.WriteByte(flags)

	buf.Reset()
	binary.Write(&buf, binary.BigEndian, m.Code)
	b.Write(buf.Bytes()[1:4])

	binary.Write(&b, binary.BigEndian, m.AppID)
	binary.Write(&b, binary.BigEndian, m.HbHID)
	binary.Write(&b, binary.BigEndian, m.EtEID)

	b.Write(m.AVPs)

	_, err := b.WriteTo(w)
	return err
}

// UnmarshalFrom read binary data from io.Reader
func (m *Message) UnmarshalFrom(r io.Reader) error {
	b, err := readUntil(r, 20)
	if err != nil {
		return err
	}
	if b[0] != 1 {
		return InvalidMessage{
			Code: UnsupportedVersion}
	}

	b[0] = 0x00
	var lng uint32
	binary.Read(bytes.NewBuffer(b[0:4]), binary.BigEndian, &lng)

	m.FlgR = b[4]&0x80 == 0x80
	m.FlgP = b[4]&0x40 == 0x40
	m.FlgE = b[4]&0x20 == 0x20
	m.FlgT = b[4]&0x10 == 0x10

	b[4] = 0x00
	binary.Read(bytes.NewBuffer(b[4:8]), binary.BigEndian, &m.Code)

	binary.Read(bytes.NewBuffer(b[8:12]), binary.BigEndian, &m.AppID)
	binary.Read(bytes.NewBuffer(b[12:16]), binary.BigEndian, &m.HbHID)
	binary.Read(bytes.NewBuffer(b[16:20]), binary.BigEndian, &m.EtEID)

	if m.AVPs, err = readUntil(r, int(lng)-20); err != nil {
		return InvalidMessage{
			Code: InvalidMessageLength}
	}

	return nil
}

func readUntil(r io.Reader, l int) ([]byte, error) {
	buf := make([]byte, l)
	offset := 0

	for offset < l {
		i, err := r.Read(buf[offset:])
		offset += i
		if err != nil {
			return buf, err
		}
	}
	return buf, nil
}
