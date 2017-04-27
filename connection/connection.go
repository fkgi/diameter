package connection

import (
	"fmt"
	"net"
	"time"

	"github.com/fkgi/diameter/msg"
)

// constant values
var (
	MsgStackLen                           = 10000
	VendorID         msg.VendorID         = 41102
	ProductName      msg.ProductName      = "yagtagarasu"
	FirmwareRevision msg.FirmwareRevision = 170407001
)

// Connection is state machine of Diameter
type Connection struct {
	local *LocalNode
	peer  *PeerNode

	wT *time.Timer // watchdog timer
	wE int         // watchdog expired counter

	notify   chan stateEvent
	state    int
	con      net.Conn
	rcvstack chan *msg.Message
	sndstack map[uint32]chan *msg.Message

	cachemsg msg.Message
}

// Close stop state machine
func (p *Connection) Close(cause msg.Enumerated) {
	if p.state != open {
		return
	}

	ch := make(chan *msg.Message)
	r := p.makeDPR(cause)
	r.HbHID = p.local.NextHbH()
	p.sndstack[r.HbHID] = ch

	p.notify <- eventStop{r}

	t := time.AfterFunc(p.peer.Tp, func() {
		ch <- nil
		notify(&PurgeEvent{
			Tx: false, Req: false, Local: p.local, Peer: p.peer,
			Err: fmt.Errorf("no answer")})
	})
	ap := <-ch
	t.Stop()
	delete(p.sndstack, r.HbHID)

	if ap != nil {
		p.notify <- eventPeerDisc{nil}
	}
}

// LocalHost returns local host name
func (p *Connection) LocalHost() msg.DiameterIdentity {
	return p.local.Host
}

// LocalRealm returns local realm name
func (p *Connection) LocalRealm() msg.DiameterIdentity {
	return p.local.Realm
}

// PeerHost returns peer host name
func (p *Connection) PeerHost() msg.DiameterIdentity {
	return p.peer.Host
}

// PeerRealm returns peer realm name
func (p *Connection) PeerRealm() msg.DiameterIdentity {
	return p.peer.Realm
}

// Connection returns transport connection of state machine
func (p *Connection) Connection() net.Conn {
	return p.con
}

// Properties returns properties of this state machine
func (p *Connection) Properties() Properties {
	return p.peer.Properties
}

// Send send Diameter request message
func (p *Connection) Send(r msg.Message) (a msg.Message, e error) {
	r.EtEID = p.local.NextEtE()

	for i := 0; i <= p.peer.Cp; i++ {
		a = p.Transmit(r)

		res := msg.DiameterUnableToComply
		if avp, e := a.Decode(); e == nil {
			if t, ok := msg.GetResultCode(avp); ok {
				res = t
			}
		}

		if res < 3000 || res > 3999 || i >= p.peer.Cp {
			break
		}
		r.FlgT = true
	}
	return
}

// Transmit send Diameter request message
func (p *Connection) Transmit(r msg.Message) (a msg.Message) {
	r.HbHID = p.local.NextHbH()

	ch := make(chan *msg.Message)
	p.sndstack[r.HbHID] = ch
	p.notify <- eventSndMsg{r}

	t := time.AfterFunc(p.peer.Tp, func() {
		a := p.makeUnableToDeliver(r)
		ch <- &a
	})

	a = *<-ch
	t.Stop()
	delete(p.sndstack, r.HbHID)

	return
}

// Recieve recieve Diameter request message
func (p *Connection) Recieve() (r msg.Message, ch chan *msg.Message, e error) {
	rp := <-p.rcvstack
	if rp == nil {
		e = fmt.Errorf("closed")
		p.rcvstack <- nil
		return
	}
	r = *rp
	ch = make(chan *msg.Message)

	go func() {
		if a := <-ch; a == nil {
			e = fmt.Errorf("request is discarded")
			notify(&MessageEvent{
				Tx: true, Req: false, Local: p.local, Peer: p.peer, Err: e})
		} else {
			a.HbHID = r.HbHID
			a.EtEID = r.EtEID
			p.notify <- eventSndMsg{*a}
		}
	}()
	return
}

const (
	shutdown = iota
	closed   = iota
	waitCER  = iota
	waitCEA  = iota
	open     = iota
	closing  = iota
)

// State returns status of state machine
func (p *Connection) State() string {
	switch p.state {
	case closed:
		return "Closed"
	case waitCER:
		return "WaitCER"
	case waitCEA:
		return "WaitCEA"
	case open:
		return "Open"
	case closing:
		return "Closing"
	}
	return "Shutdown"
}

func (p *Connection) run() {
	go func() {
		for {
			m := msg.Message{}
			p.con.SetReadDeadline(time.Time{})
			if _, e := m.ReadFrom(p.con); e != nil {
				p.notify <- eventPeerDisc{e}
				break
			}

			p.wE = 0
			if m.AppID == 0 && m.Code == 257 && m.FlgR {
				p.notify <- eventRcvCER{m}
			} else if m.AppID == 0 && m.Code == 257 && !m.FlgR {
				p.notify <- eventRcvCEA{m}
			} else if m.AppID == 0 && m.Code == 280 && m.FlgR {
				p.notify <- eventRcvDWR{m}
			} else if m.AppID == 0 && m.Code == 280 && !m.FlgR {
				p.notify <- eventRcvDWA{m}
			} else if m.AppID == 0 && m.Code == 282 && m.FlgR {
				p.notify <- eventRcvDPR{m}
			} else if m.AppID == 0 && m.Code == 282 && !m.FlgR {
				p.notify <- eventRcvDPA{m}
			} else {
				p.notify <- eventRcvMsg{m}
			}
		}
		if p.wT != nil {
			p.wT.Stop()
		}
	}()

	old := "Shutdown"
	notify(&StateUpdate{
		OldState: old, NewState: p.State(), Event: "Start",
		Local: p.local, Peer: p.peer, Err: nil})

	for p.state != shutdown {
		event := <-p.notify
		old = p.State()
		e := event.exec(p)

		notify(&StateUpdate{
			OldState: old, NewState: p.State(), Event: event.name(),
			Local: p.local, Peer: p.peer, Err: e})
	}
}

func (p *Connection) resetWatchdog() {
	f := func() {
		p.wT = nil

		p.wE++
		if p.wE > p.peer.Ew {
			p.Close(msg.Enumerated(0))
		} else {
			ch := make(chan *msg.Message)
			r := p.makeDWR()
			r.HbHID = p.local.NextHbH()
			p.sndstack[r.HbHID] = ch
			p.notify <- eventWatchdog{r}

			t := time.AfterFunc(p.peer.Tp, func() {
				ch <- nil
				notify(&WatchdogEvent{
					Tx: false, Req: false, Local: p.local, Peer: p.peer,
					Err: fmt.Errorf("no answer")})
			})
			ap := <-ch
			t.Stop()
			delete(p.sndstack, r.HbHID)
			if ap != nil {
				p.wE = 0
			}
		}
	}

	if p.wT != nil {
		p.wT.Reset(p.peer.Tw)
	} else {
		p.wT = time.AfterFunc(p.peer.Tw, f)
	}
}
