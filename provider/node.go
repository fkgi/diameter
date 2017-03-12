package provider

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/extnet"
)

// LocalNode is local node of Diameter
type LocalNode struct {
	Realm msg.DiameterIdentity
	Host  msg.DiameterIdentity
	Addr  []net.IP // IP address for CER/A AVP

	hbHId chan uint32
	etEId chan uint32
}

// PeerNode is peer node of Diameter
type PeerNode struct {
	Realm msg.DiameterIdentity
	Host  msg.DiameterIdentity
	Addr  []net.IP // IP address for CER/A AVP

	Tw time.Duration // DWR send interval time
	Ew int           // watchdog expired count

	Ts time.Duration // transport packet send timeout

	Tp time.Duration // pending Diameter answer time
	Cp int           // retry Diameter request count

	SupportedApps [][2]uint32
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

// Connect is Low-level diameter connect
func (l *LocalNode) Connect(p *PeerNode, laddr, raddr net.Addr, s time.Duration) (c *Connection, e error) {
	if raddr == nil {
		e = fmt.Errorf("Remote address is nil")
	} else if laddr == nil {
		e = fmt.Errorf("Local address is nil")
	} else if p == nil {
		e = fmt.Errorf("Peer node is nil")
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

// Accept is Low-level diameter accept
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

func (l *LocalNode) String() string {
	if l == nil {
		return "<nil>"
	}
	return string(l.Host)
}

func (p *PeerNode) String() string {
	if p == nil {
		return "<nil>"
	}
	return string(p.Host)
}
