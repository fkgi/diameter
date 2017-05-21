package connection

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/fkgi/diameter/msg"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// LocalNode is local node of Diameter
type LocalNode struct {
	Realm msg.DiameterIdentity
	Host  msg.DiameterIdentity
	Properties

	hbHID     chan uint32
	etEID     chan uint32
	sessionID chan uint32
}

// Properties is set of diameter properties
type Properties struct {
	Tw time.Duration // DWR send interval time
	Ew int           // watchdog expired count

	Ts time.Duration // transport packet send timeout

	Tp time.Duration // pending Diameter answer time
	Cp int           // retry Diameter request count

	Apps []AuthApplication
	// for Vendor-Specific-Application-Id,
	//     Auth-Application-Id, Supported-Vendor-Id AVP
}

// AuthApplication is set of authenticated application
type AuthApplication struct {
	VendorID msg.VendorID
	AppID    msg.AuthApplicationID
}

// InitIDs initiate each IDs
func (l *LocalNode) InitIDs() {
	l.hbHID = make(chan uint32, 1)
	l.hbHID <- rand.Uint32()

	l.etEID = make(chan uint32, 1)
	tmp := uint32(time.Now().Unix() ^ 0xFFF)
	tmp = (tmp << 20) | (rand.Uint32() ^ 0x000FFFFF)
	l.etEID <- tmp

	l.sessionID = make(chan uint32, 1)
	l.sessionID <- rand.Uint32()
}

// NextHbH make HbH ID
func (l *LocalNode) NextHbH() uint32 {
	ret := <-l.hbHID
	l.hbHID <- ret + 1
	return ret
}

// NextEtE make EtE ID
func (l *LocalNode) NextEtE() uint32 {
	ret := <-l.etEID
	l.etEID <- ret + 1
	return ret
}

// NextSession make session ID
func (l *LocalNode) NextSession() string {
	ret := <-l.sessionID
	l.sessionID <- ret + 1
	return fmt.Sprintf("%s;%d;%d;0",
		l.Host, time.Now().Unix()+2208988800, ret)
}

// Dial make new Connection that use specified peernode and connection
func (l *LocalNode) Dial(n *PeerNode, c net.Conn) *Connection {
	p := &Connection{
		notify:   make(chan stateEvent),
		state:    closed,
		con:      c,
		rcvstack: make(chan *msg.Message, MsgStackLen),
		sndstack: make(map[uint32]chan *msg.Message),
		openNtfy: make(chan bool, 1),
		local:    l,
		peer:     n}

	go p.run()
	r := p.makeCER(p.con)
	r.HbHID = p.local.NextHbH()

	p.notify <- eventConnect{r}

	return p
}

// Accept accept new transport connection and return Connection
func (l *LocalNode) Accept(c net.Conn) *Connection {
	p := &Connection{
		notify:   make(chan stateEvent),
		state:    closed,
		con:      c,
		rcvstack: make(chan *msg.Message, MsgStackLen),
		sndstack: make(map[uint32]chan *msg.Message),
		openNtfy: make(chan bool, 1),
		local:    l,
		peer:     &PeerNode{Properties: l.Properties}}

	go p.run()
	p.notify <- eventAccept{}

	return p
}

func (l *LocalNode) String() string {
	if l == nil {
		return "<nil>"
	}
	return string(l.Host)
}

// PeerNode is peer node of Diameter
type PeerNode struct {
	Realm msg.DiameterIdentity
	Host  msg.DiameterIdentity
	Properties
}

func (p *PeerNode) String() string {
	if p == nil {
		return "<nil>"
	}
	return string(p.Host)
}
