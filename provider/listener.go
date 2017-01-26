package provider

import (
	"fmt"
	"net"

	"github.com/fkgi/diameter/msg"
)

// Listener is Diameter server listener
type Listener struct {
	local *LocalNode
	provs map[*PeerNode]*Provider
}

// Listen make Listener with LocalNode
func Listen(n *LocalNode) (l *Listener) {
	//	n.InitIDs()
	l = &Listener{}
	l.local = n
	l.provs = make(map[*PeerNode]*Provider)
	return
}

// Bind bind network listener to Listener
func (l *Listener) Bind(lnr net.Listener) {
	for {
		c, e := l.local.Accept(lnr)
		// output logs
		if Notificator != nil {
			Notificator(&TransportBind{
				Local: l.local, LAddr: lnr.Addr(), Err: e})
		}
		if e != nil {
			break
		} else {
			go l.bindProvider(c)
		}
	}
}

func (l *Listener) bindProvider(c *Connection) {
	// R-Accept
	m, e := c.Read(0)
	if e == nil && !(m.AppID == 0 && m.Code == 257 && m.FlgR) {
		e = fmt.Errorf("not CER")
	}
	if e != nil {
		// output logs
		if Notificator != nil {
			Notificator(&CapabilityExchangeEvent{
				Tx: false, Req: true, Local: l.local, Peer: nil, Err: e})
		}
		c.Close()
		return
	}

	avp, e := m.Decode()
	var h, r msg.DiameterIdentity
	for _, a := range avp {
		if a.Code == uint32(264) && a.VenID == uint32(0) {
			a.Decode(&h)
		}
		if a.Code == uint32(296) && a.VenID == uint32(0) {
			a.Decode(&r)
		}
	}

	for k, v := range l.provs {
		if k.Host == h && k.Realm == r {
			c.Peer = k
			v.Notify <- eventRConnCER{m, c}
			return
		}
	}

	if Notificator != nil {
		e = fmt.Errorf("CER from unknown peer")
		Notificator(&CapabilityExchangeEvent{
			Tx: false, Req: true, Local: l.local, Peer: nil, Err: e})
	}
	c.Close()
}

// AddPeer add new PeerNode to Listener
func (l *Listener) AddPeer(n *PeerNode) (p *Provider) {
	p = &Provider{}

	p.Notify = make(chan stateEvent)

	p.state = shutdown

	p.rcvstack = make(chan *msg.Message, MsgStackLen)
	p.sndstack = make(map[uint32]chan *msg.Message)

	go p.run()

	l.provs[n] = p
	return
}

// Dial add new PeerNode to Listener
func (l *Listener) Dial(n *PeerNode, laddr, raddr net.Addr) (e error) {
	p := l.provs[n]
	if p == nil {
		return
	}
	if p.state == shutdown {
		p.state = closed
	}
	p.Notify <- eventStart{laddr, raddr, l.local, n}
	return
}
