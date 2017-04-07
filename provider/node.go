package provider

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/fkgi/diameter/msg"
)

// LocalNode is local node of Diameter
type LocalNode struct {
	Realm msg.DiameterIdentity
	Host  msg.DiameterIdentity
	Properties

	hbHId chan uint32
	etEId chan uint32
}

// Properties is set of diameter properties
type Properties struct {
	Tw time.Duration // DWR send interval time
	Ew int           // watchdog expired count

	Ts time.Duration // transport packet send timeout

	Tp time.Duration // pending Diameter answer time
	Cp int           // retry Diameter request count

	Apps [][2]uint32
	// for Vendor-Specific-Application-Id,
	//     Auth-Application-Id, Supported-Vendor-Id AVP
}

// InitIDs initiate each IDs
func (l *LocalNode) InitIDs() {
	l.hbHId = make(chan uint32, 1)
	l.hbHId <- rand.Uint32()

	l.etEId = make(chan uint32, 1)
	tmp := uint32(time.Now().Unix() ^ 0xFFF)
	tmp = (tmp << 20) | (rand.Uint32() ^ 0x000FFFFF)
	l.etEId <- tmp
}

// NextHbH make HbH ID
func (l *LocalNode) NextHbH() uint32 {
	ret := <-l.hbHId
	l.hbHId <- ret + 1
	return ret
}

// NextEtE make EtE ID
func (l *LocalNode) NextEtE() uint32 {
	ret := <-l.etEId
	l.etEId <- ret + 1
	return ret
}

// Dial make new provider that use specified peernode and connection
func (l *LocalNode) Dial(n *PeerNode, c net.Conn) *Provider {
	p := &Provider{
		notify:   make(chan stateEvent),
		state:    closed,
		rcvstack: make(chan *msg.Message, MsgStackLen),
		sndstack: make(map[uint32]chan *msg.Message),
		local:    l,
		peer:     n}

	go p.run()
	p.notify <- eventStart{c}

	return p
}

// Accept accept new transport connection and return provider
func (l *LocalNode) Accept(c net.Conn) *Provider {
	p := &Provider{
		notify:   make(chan stateEvent),
		state:    shutdown,
		rcvstack: make(chan *msg.Message, MsgStackLen),
		sndstack: make(map[uint32]chan *msg.Message),
		local:    l,
		peer:     &PeerNode{Properties: l.Properties}}

	m := msg.Message{}
	c.SetReadDeadline(time.Time{})
	_, e := m.ReadFrom(c)

	if e != nil {
		c.Close()
	} else if !(m.AppID == 0 && m.Code == 257 && m.FlgR) {
		e = fmt.Errorf("initial message is not CER")
		c.Close()
	}

	if e == nil {
		if avp, e := m.Decode(); e == nil {
			for _, a := range avp {
				if a.Code == uint32(264) && a.VenID == uint32(0) {
					a.Decode(&(p.peer.Host))
				}
				if a.Code == uint32(296) && a.VenID == uint32(0) {
					a.Decode(&(p.peer.Realm))
				}
			}
		}
		p.state = closed

		go p.run()
		p.notify <- eventRConnCER{m, c}
	} else if Notificator != nil {
		Notificator(&CapabilityExchangeEvent{
			Tx: false, Req: true,
			Local: p.local, Peer: p.peer, Err: e})
	}

	return p
}

// Connect is Low-level diameter connect    laddr, raddr net.Addr, s time.Duration
/*
func (l *LocalNode) Connect(p *PeerNode, con net.Conn) (c *Connection, e error) {
	if p == nil {
		e = fmt.Errorf("Peer node is nil")
	}
	c = &Connection{p, l, con}
			} else if laddr.Network() == "sctp" {
				if a, ok := laddr.(*extnet.SCTPAddr); !ok {
					e = fmt.Errorf("address type mismatch")
				} else {
					dialer := extnet.SCTPDialer{
						InitTimeout: s,
						PPID:        46,
						Unordered:   true,
						LocalAddr:   a}
					var con net.Conn
					if con, e = dialer.Dial(raddr.Network(), raddr.String()); e == nil {
						c = &Connection{p, l, con}
					}
				}

			} else if laddr.Network() == "tcp" {
				dialer := net.Dialer{
					Timeout:   s,
					LocalAddr: laddr}

				var con net.Conn
				if con, e = dialer.Dial(raddr.Network(), raddr.String()); e == nil {
					c = &Connection{p, l, con}
				}
			} else {
				e = fmt.Errorf("invalid address type")
			}
		// output logs
		if Notificator != nil {
			Notificator(&TransportStateChange{
				Open: true, Local: l, Peer: p, LAddr: laddr, PAddr: raddr, Err: e})
		}
	return
}
*/

// Accept is Low-level diameter accept
/*
func (l *LocalNode) Accept(lnr net.Listener) (c *Connection, e error) {
	if lnr == nil {
		e = fmt.Errorf("Local listener is nil")
	} else {
		var con net.Conn
		if con, e = lnr.Accept(); e == nil {
			c = &Connection{nil, l, con}
		}
	}

	// output logs
	if Notificator != nil {
		var pa net.Addr
		if e == nil {
			pa = c.conn.RemoteAddr()
		}
		Notificator(&TransportStateChange{
			Open: true, Local: l, Peer: nil, LAddr: lnr.Addr(), PAddr: pa, Err: e})
	}
	return
}
*/
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
