package diameter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"time"
)

const (
	VendorID    uint32 = 41102         // VendorID of this code
	ProductName        = "round-robin" // ProductName of this code
	FirmwareRev uint32 = 230619001     // FirmwareRev of this code

	avpBufferSize = 10
)

var (
	hbhID     = make(chan uint32, 1) // Hop-by-Hop ID source
	eteID     = make(chan uint32, 1) // End-to-End ID source
	sessionID = make(chan uint32, 1) // Session-ID source
)

func init() {
	ut := time.Now().Unix()
	hbhID <- rand.Uint32()
	eteID <- (uint32(ut^0xFFF) << 20) | (rand.Uint32() ^ 0xFFFFF)
	sessionID <- rand.Uint32()
	stateID = uint32(ut)
}

type application struct {
	venID    uint32
	handlers map[uint32]Handler
}

func nextHbH() uint32 {
	ret := <-hbhID
	hbhID <- ret + 1
	return ret
}

func nextEtE() uint32 {
	ret := <-eteID
	eteID <- ret + 1
	return ret
}

// NextSession generate new session ID data
func NextSession(h string) string {
	ret := <-sessionID
	sessionID <- ret + 1
	if h == "" {
		h = Host.String()
	}
	return fmt.Sprintf("%s;%d;%d;0", h, time.Now().Unix()+2208988800, ret)
}

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

	PeerName  Identity // peer node that send this message
	PeerRealm Identity
}

func (m *Message) SetAVP(avp []AVP) {
	buf := new(bytes.Buffer)
	for _, a := range avp {
		a.MarshalTo(buf)
	}
	m.AVPs = buf.Bytes()
}

func (m *Message) GetAVP() ([]AVP, error) {
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

func (m Message) GenerateAnswerBy(result uint32) Message {
	buf := new(bytes.Buffer)
	SetResultCode(result).MarshalTo(buf)
	for rdr := bytes.NewReader(m.AVPs); rdr.Len() != 0; {
		a := AVP{}
		if e := a.UnmarshalFrom(rdr); e != nil {
			continue
		}
		if a.VendorID != 0 {
			continue
		}
		switch a.Code {
		case 277:
			a.MarshalTo(buf)
		case 263:
			a.MarshalTo(buf)
		}
	}
	SetOriginHost(Host).MarshalTo(buf)
	SetOriginRealm(Realm).MarshalTo(buf)

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
