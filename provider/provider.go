package provider

import (
	"fmt"
	"net"
	"time"

	"github.com/fkgi/diameter/msg"
)

// constant values
var (
	MsgStackLen      = 10000
	VendorID         = uint32(41102)
	ProductName      = "yagtagarasu"
	FirmwareRevision = uint32(170407001)
)

const (
	shutdown         = iota
	closed           = iota
	waitConnAck      = iota
	waitICEA         = iota
	waitConnAckElect = iota
	waitReturns      = iota
	rOpen            = iota
	iOpen            = iota
	closing          = iota

	// Tzero is zero duration
	Tzero = time.Duration(0)
)

// Provider is state machine of Diameter
type Provider struct {
	local *LocalNode
	peer  *PeerNode

	wT *time.Timer // watchdog timer
	wE int         // watchdog expired counter

	notify     chan stateEvent
	state      int
	icon, rcon net.Conn
	rcvstack   chan *msg.Message
	sndstack   map[uint32]chan *msg.Message

	cachemsg msg.Message
}

// Close stop state machine
func (p *Provider) Close(cause msg.Enumerated) {
	if p.state == rOpen || p.state == iOpen {
		p.notify <- eventStop{cause}
	}
}

// State returns status of state machine
func (p *Provider) State() string {
	switch p.state {
	case closed:
		return "Closed"
	case waitConnAck:
		return "Wait-Conn-Ack"
	case waitICEA:
		return "Wait-I-CEA"
	case waitConnAckElect:
		return "Wait-Conn-Ack-Elect"
	case waitReturns:
		return "Wait-Returns"
	case rOpen:
		return "R-Open"
	case iOpen:
		return "I-Open"
	case closing:
		return "Closing"
	}
	return "Shutdown"
}

// LocalHost returns local host name
func (p *Provider) LocalHost() msg.DiameterIdentity {
	return p.local.Host
}

// LocalRealm returns local realm name
func (p *Provider) LocalRealm() msg.DiameterIdentity {
	return p.local.Realm
}

// PeerHost returns peer host name
func (p *Provider) PeerHost() msg.DiameterIdentity {
	return p.peer.Host
}

// PeerRealm returns peer realm name
func (p *Provider) PeerRealm() msg.DiameterIdentity {
	return p.peer.Realm
}

// Connection returns transport connection of state machine
func (p *Provider) Connection() net.Conn {
	return p.activeConnection()
}

// Properties returns properties of this state machine
func (p *Provider) Properties() Properties {
	return p.peer.Properties
}

func (p *Provider) run() {
	if Notificator != nil {
		Notificator(&StateUpdate{
			OldState: "Shutdown", NewState: p.State(), Event: "Start",
			Local: p.local, Peer: p.peer, Err: nil})
	}

	for p.state != shutdown {
		event := <-p.notify
		old := p.State()
		e := event.exec(p)

		if Notificator != nil {
			Notificator(&StateUpdate{
				OldState: old, NewState: p.State(), Event: event.name(),
				Local: p.local, Peer: p.peer, Err: e})
		}
	}
}

func (p *Provider) open() {
	if Notificator != nil {
		Notificator(&ConnectionStateChange{
			Open: true, Local: p.local, Peer: p.peer})
	}
	p.resetWatchdog()

	for p.state == rOpen {
		m := msg.Message{}
		p.rcon.SetReadDeadline(time.Time{})
		_, e := m.ReadFrom(p.rcon)

		if Notificator != nil {
			Notificator(&MessageTransfer{
				Tx: false, Local: p.local, Peer: p.peer, Err: e, dump: m.PrintStack})
		}

		if e != nil {
			// Disconnected
			p.notify <- eventRPeerDisc{}
			break
		}

		p.wE = 0
		p.resetWatchdog()
		if m.AppID == 0 && m.Code == 257 && m.FlgR {
			p.notify <- eventRcvCER{m}
		} else if m.AppID == 0 && m.Code == 257 && !m.FlgR {
			p.notify <- eventRRcvCEA{m}
		} else if m.AppID == 0 && m.Code == 280 && m.FlgR {
			p.notify <- eventRcvDWR{m}
		} else if m.AppID == 0 && m.Code == 280 && !m.FlgR {
			p.notify <- eventRcvDWA{m}
		} else if m.AppID == 0 && m.Code == 282 && m.FlgR {
			p.notify <- eventRcvDPR{m}
		} else if m.AppID == 0 && m.Code == 282 && !m.FlgR {
			p.notify <- eventRRcvDPA{m}
		} else {
			p.notify <- eventRcvMsg{m}
		}
	}

	for p.state == iOpen {
		m := msg.Message{}
		p.icon.SetReadDeadline(time.Time{})
		_, e := m.ReadFrom(p.icon)

		if Notificator != nil {
			Notificator(&MessageTransfer{
				Tx: false, Local: p.local, Peer: p.peer, Err: e, dump: m.PrintStack})
		}

		if e != nil {
			// Disconnected
			p.notify <- eventIPeerDisc{}
			break
		}

		p.wE = 0
		p.resetWatchdog()
		if m.AppID == 0 && m.Code == 257 && m.FlgR {
			p.notify <- eventRcvCER{m}
		} else if m.AppID == 0 && m.Code == 257 && !m.FlgR {
			p.notify <- eventIRcvCEA{m}
		} else if m.AppID == 0 && m.Code == 280 && m.FlgR {
			p.notify <- eventRcvDWR{m}
		} else if m.AppID == 0 && m.Code == 280 && !m.FlgR {
			p.notify <- eventRcvDWA{m}
		} else if m.AppID == 0 && m.Code == 282 && m.FlgR {
			p.notify <- eventRcvDPR{m}
		} else if m.AppID == 0 && m.Code == 282 && !m.FlgR {
			p.notify <- eventIRcvDPA{m}
		} else {
			p.notify <- eventRcvMsg{m}
		}
	}

	if p.wT != nil {
		p.wT.Stop()
	}
	if Notificator != nil {
		Notificator(&ConnectionStateChange{
			Open: false, Local: p.local, Peer: p.peer})
	}
}

func (p *Provider) activeConnection() net.Conn {
	if p.state == iOpen {
		return p.icon
	} else if p.state == rOpen {
		return p.rcon
	}
	return nil
}

func (p *Provider) resetWatchdog() {
	f := func() {
		p.wT = nil
		c := p.activeConnection()
		if c == nil {
			return
		}

		p.wE++
		if p.wE > p.peer.Ew {
			p.notify <- eventStop{msg.Enumerated(0)}
		} else {
			ch := make(chan *msg.Message)
			r := p.makeDWR()
			r.HbHID = p.local.NextHbH()
			p.sndstack[r.HbHID] = ch

			c.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
			_, e := r.WriteTo(c)

			if Notificator != nil {
				Notificator(&WatchdogEvent{
					Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
			}

			if e != nil {
				switch p.state {
				case rOpen:
					p.notify <- eventRPeerDisc{}
				case iOpen:
					p.notify <- eventIPeerDisc{}
				}
				return
			}

			t := time.AfterFunc(p.peer.Properties.Tp, func() {
				ch <- nil
			})
			ap := <-ch
			t.Stop()
			delete(p.sndstack, r.HbHID)

			if ap == nil {
				if Notificator != nil {
					Notificator(&WatchdogEvent{
						Tx: false, Req: false, Local: p.local, Peer: p.peer,
						Err: fmt.Errorf("no answer")})
				}
			} else {
				p.resetWatchdog()
			}
		}
	}

	if c := p.activeConnection(); c != nil {
		if p.wT != nil {
			p.wT.Reset(p.peer.Tw)
		} else {
			p.wT = time.AfterFunc(p.peer.Tw, f)
		}
	}
}

// Send send Diameter request message
func (p *Provider) Send(r msg.Message) (a msg.Message, e error) {
	c := p.activeConnection()
	if c == nil {
		e = fmt.Errorf("connection is not open")
		return
	}
	r.EtEID = p.local.NextEtE()

	for i := 0; i <= p.peer.Cp; i++ {
		ap := sendReq(r, c, p)
		if ap == nil {
			if i == p.peer.Cp {
				e = fmt.Errorf("No answer, retry expired")
			} else {
				e = fmt.Errorf("No answer, retry x%d", i+1)
			}
			if Notificator != nil {
				Notificator(&MessageEvent{
					Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
			}
			r.FlgT = true
		} else {
			a = *ap
			e = nil
			break
		}
	}
	return
}

func sendReq(r msg.Message, c net.Conn, p *Provider) (a *msg.Message) {
	ch := make(chan *msg.Message)

	r.HbHID = p.local.NextHbH()
	p.sndstack[r.HbHID] = ch

	p.notify <- eventSndMsg{r}
	t := time.AfterFunc(p.peer.Tp, func() {
		ch <- nil
	})

	a = <-ch
	delete(p.sndstack, r.HbHID)
	t.Stop()
	return
}

// Recieve recieve Diameter request message
func (p *Provider) Recieve() (r msg.Message, ch chan *msg.Message, e error) {
	if p.state != rOpen && p.state != iOpen {
		e = fmt.Errorf("connection is not open")
	} else if rp := <-p.rcvstack; rp == nil {
		e = fmt.Errorf("Peer Node closed")
		p.rcvstack <- nil
	} else {
		r = *rp
		ch = make(chan *msg.Message)
		go func() {
			if a := <-ch; a != nil {
				a.HbHID = r.HbHID
				a.EtEID = r.EtEID
				p.notify <- eventSndMsg{*a}
			} else if Notificator != nil {
				Notificator(&MessageEvent{
					Tx: true, Req: false, Local: p.local, Peer: p.peer,
					Err: fmt.Errorf("No answer")})
			}
		}()
	}
	return
}

type stateEvent interface {
	exec(p *Provider) error
	name() string
}

type eventStart struct {
	con net.Conn
}

func (v eventStart) name() string {
	return "Start"
}

func (v eventStart) exec(p *Provider) (e error) {
	switch p.state {
	case closed:
		// I-Snd-Conn-Req
		p.state = waitConnAck
		go func() {
			p.icon = v.con
			p.notify <- eventIRcvConAck{}
		}()
	default:
		// wait_conn_ack
		// wait_i_CEA
		// wait_conn_ack_elect
		// wait_returns
		// r_open
		// i_open
		// closing
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

/*
The election is performed on the responder.
The responder compares the Origin-Host received in the CER with its own Origin-Host as two streams of octets.
If the local Origin-Host lexicographically succeeds the received Origin-Host, a Win-Election event is issued locally.
Diameter identities are in ASCII form; therefore, the lexical comparison is consistent with DNS case insensitivity,
 where octets that fall in the ASCII range 'a' through 'z' MUST compare equally to their uppercase counterparts between 'A' and 'Z'.
See Appendix D for interactions between the Diameter protocol and Internationalized Domain Name (IDNs).

The winner of the election MUST close the connection it initiated.
Historically, maintaining the responder side of a connection was more efficient than maintaining the initiator side.
However, current practices makes this distinction irrelevant.
*/
type eventRConnCER struct {
	m msg.Message
	c net.Conn //*Connection
}

func (v eventRConnCER) name() string {
	return "R-Conn-CER"
}

func (v eventRConnCER) exec(p *Provider) (e error) {
	if Notificator != nil {
		Notificator(&CapabilityExchangeEvent{
			Tx: false, Req: true, Local: p.local, Peer: p.peer})
	}

	switch p.state {
	case closed:
		// R-Accept, Process-CER
		m, code := p.makeCEA(v.m, v.c)

		// R-Snd-CEA
		v.c.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
		_, e := m.WriteTo(v.c)

		if e != nil {
			v.c.Close()
			p.state = shutdown
		} else if code != 2001 {
			e = fmt.Errorf("close with error response %d", code)
			v.c.Close()
			p.state = shutdown
		} else {
			p.rcon = v.c
			p.state = rOpen
			go p.open()
		}

		if Notificator != nil {
			Notificator(&CapabilityExchangeEvent{
				Tx: true, Req: false, Local: p.local, Peer: p.peer})
		}
	case waitConnAck:
		p.state = waitConnAckElect
		// R-Accept, Process-CER
		p.cachemsg, _ = p.makeCEA(v.m, v.c)
		p.rcon = v.c
	case waitICEA:
		p.state = waitReturns
		// R-Accept, Process-CER
		p.cachemsg, _ = p.makeCEA(v.m, v.c)
		p.rcon = v.c
		// Elect
		if msg.Compare(p.local.Host, p.peer.Host) > 0 {
			p.notify <- eventWinElection{}
		}
	default:
		// wait_conn_ack_elect
		// wait_returns
		// r_open
		// i_open
		// closing
		// shutdown

		// R-Reject
		v.c.Close()
		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventIRcvConAck struct{}

func (v eventIRcvConAck) name() string {
	return "I-Rcv-Con-Ack"
}

func (v eventIRcvConAck) exec(p *Provider) (e error) {
	switch p.state {
	case waitConnAck:
		p.state = waitICEA
		// I-Snd-CER
		r := p.makeCER(p.icon)
		r.HbHID = p.local.NextHbH()

		go func() {
			p.icon.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
			_, e := r.WriteTo(p.icon)

			if Notificator != nil {
				Notificator(&CapabilityExchangeEvent{
					Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
			}
			if e != nil {
				p.notify <- eventIPeerDisc{}
				return
			}

			a := msg.Message{}
			p.icon.SetReadDeadline(time.Time{})
			_, e = a.ReadFrom(p.icon)
			if e != nil {
				if Notificator != nil {
					Notificator(&CapabilityExchangeEvent{
						Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
				}
				p.notify <- eventIPeerDisc{}
			} else if r.HbHID == a.HbHID && a.AppID == 0 && a.Code == 257 && !a.FlgR {
				p.notify <- eventIRcvCEA{a}
			} else {
				p.notify <- eventIRcvNonCEA{a}
			}
		}()
	case waitConnAckElect:
		p.state = waitReturns
		// I-Snd-CER
		r := p.makeCER(p.icon)
		r.HbHID = p.local.NextHbH()

		go func() {
			p.icon.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
			_, e := r.WriteTo(p.icon)
			if Notificator != nil {
				Notificator(&CapabilityExchangeEvent{
					Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
			}
			if e != nil {
				p.notify <- eventIPeerDisc{}
				return
			}
			if msg.Compare(p.local.Host, p.peer.Host) > 0 {
				// Elect
				p.notify <- eventWinElection{}
				return
			}

			a := msg.Message{}
			p.icon.SetReadDeadline(time.Time{})
			_, e = a.ReadFrom(p.icon)
			if e != nil {
				if Notificator != nil {
					Notificator(&CapabilityExchangeEvent{
						Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
				}
				p.notify <- eventIPeerDisc{}
			} else if r.HbHID == a.HbHID && a.AppID == 0 && a.Code == 257 && a.FlgR {
				p.notify <- eventIRcvCEA{a}
			} else {
				p.notify <- eventIRcvNonCEA{a}
			}
		}()
	default:
		// closed
		// wait_i_CEA
		// wait_returns
		// r_open
		// i_open
		// closing
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventIRcvConNack struct{}

func (v eventIRcvConNack) name() string {
	return "I-Rcv-Con-Nack"
}

func (v eventIRcvConNack) exec(p *Provider) (e error) {
	switch p.state {
	case waitConnAck:
		// Cleanup
		if p.icon != nil {
			p.icon.Close()
		}
		p.state = closed
	case waitConnAckElect:
		if p.icon != nil {
			p.icon.Close()
		}

		p.state = rOpen
		go p.open()
		// R-Snd-CEA
		go func() {
			p.rcon.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
			_, e := p.cachemsg.WriteTo(p.rcon)
			if Notificator != nil {
				Notificator(&CapabilityExchangeEvent{
					Tx: true, Req: false, Local: p.local, Peer: p.peer, Err: e})
			}
			if e != nil {
				p.notify <- eventRPeerDisc{}
			}
		}()
	default:
		// closed
		// wait_i_CEA
		// wait_returns
		// r_open
		// i_open
		// closing
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventTimeout struct {
	s string
}

func (v eventTimeout) name() string {
	return "Timeout"
}

func (v eventTimeout) exec(p *Provider) (e error) {
	switch p.state {
	case waitConnAck, waitICEA, waitConnAckElect, waitReturns, closing:
		// Error
		// monitor.Notificator(monitor.Major, "request timeout", v.s)
		p.state = closed
	default:
		// closed
		// r_open
		// i_open
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventRcvCER struct {
	m msg.Message
}

func (v eventRcvCER) name() string {
	return "Rcv-CER"
}

func (v eventRcvCER) exec(p *Provider) (e error) {
	e = fmt.Errorf("not acceptable event")

	if Notificator != nil {
		Notificator(&CapabilityExchangeEvent{
			Tx: false, Req: true, Local: p.local, Peer: p.peer, Err: e})
	}

	return
}

type eventRRcvCEA struct {
	m msg.Message
}

func (v eventRRcvCEA) name() string {
	return "R-Rcv-CEA"
}

func (v eventRRcvCEA) exec(p *Provider) (e error) {
	e = fmt.Errorf("not acceptable event")
	if Notificator != nil {
		Notificator(&CapabilityExchangeEvent{
			Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
	}

	return
}

type eventIRcvCEA struct {
	m msg.Message
}

func (v eventIRcvCEA) name() string {
	return "I-Rcv-CEA"
}

func (v eventIRcvCEA) exec(p *Provider) (e error) {

	switch p.state {
	case waitICEA:
		// Process-CEA
		c := new(uint32)
		if avp, e := v.m.Decode(); e == nil {
			for _, a := range avp {
				if a.Code == uint32(268) && a.VenID == uint32(0) {
					a.Decode(c)
					break
				}
			}
		}
		if c != nil && *c == 2001 {
			if Notificator != nil {
				Notificator(&CapabilityExchangeEvent{
					Tx: false, Req: false, Local: p.local, Peer: p.peer})
			}
			p.state = iOpen
			go p.open()
		} else {
			// Close
			if Notificator != nil {
				Notificator(&CapabilityExchangeEvent{
					Tx: false, Req: false, Local: p.local, Peer: p.peer,
					Err: fmt.Errorf("CEA Nack received")})
			}
			p.icon.Close()
			p.icon = nil
			p.state = shutdown
		}
	case waitReturns:
		if Notificator != nil {
			Notificator(&CapabilityExchangeEvent{
				Tx: false, Req: false, Local: p.local, Peer: p.peer})
		}

		// R-Disc
		p.rcon.Close()
		p.rcon = nil
		p.state = iOpen
		go p.open()
	default:
		// closed
		// wait_conn_ack
		// wait_conn_ack_elect
		// r_open
		// i_open
		// closing
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventIRcvNonCEA struct {
	m msg.Message
}

func (v eventIRcvNonCEA) name() string {
	return "I-Rcv-Non-CEA"
}

func (v eventIRcvNonCEA) exec(p *Provider) (e error) {

	switch p.state {
	case waitICEA:
		if Notificator != nil {
			Notificator(&CapabilityExchangeEvent{
				Tx: false, Req: false, Local: p.local, Peer: p.peer,
				Err: fmt.Errorf("None CEA received")})
		}
		// Error
		p.icon.Close()
		p.icon = nil
		p.state = shutdown
	default:
		// closed
		// wait_conn_ack
		// wait_conn_ack_elect
		// wait_returns
		// r_open
		// i_open
		// closing
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventRPeerDisc struct {
	Err error
}

func (v eventRPeerDisc) name() string {
	return "R-Peer-Disc"
}

func (v eventRPeerDisc) exec(p *Provider) (e error) {
	switch p.state {
	case waitConnAckElect:
		// R-Disc
		p.rcon.Close()
		p.rcon = nil
		p.state = waitConnAck
	case waitReturns:
		// R-Disc
		p.rcon.Close()
		p.rcon = nil
		p.state = waitICEA
	case rOpen, closing:
		// R-Disc
		p.rcon.Close()
		p.rcon = nil
		p.state = shutdown
	default:
		// closed
		// wait_conn_ack
		// wait_i_CEA
		// i_open
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventIPeerDisc struct {
	Err error
}

func (v eventIPeerDisc) name() string {
	return "I-Peer-Disc"
}

func (v eventIPeerDisc) exec(p *Provider) (e error) {
	switch p.state {
	case waitICEA, iOpen, closing:
		// I-Disc
		p.icon.Close()
		p.icon = nil
		p.state = shutdown
	case waitReturns:
		// I-Disc
		p.icon.Close()
		p.icon = nil

		p.state = rOpen
		go p.open()
		// R-Snd-CEA
		go func() {
			p.rcon.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
			_, e := p.cachemsg.WriteTo(p.rcon)
			if Notificator != nil {
				Notificator(&CapabilityExchangeEvent{
					Tx: true, Req: false, Local: p.local, Peer: p.peer, Err: e})
			}
			if e != nil {
				p.notify <- eventRPeerDisc{}
			}
		}()
	default:
		// closed
		// wait_conn_ack
		// wait_conn_ack_elect
		// r_open
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventRcvDWR struct {
	m msg.Message
}

func (v eventRcvDWR) name() string {
	return "Rcv-DWR"
}

func (v eventRcvDWR) exec(p *Provider) (e error) {
	if Notificator != nil {
		Notificator(&WatchdogEvent{
			Tx: false, Req: true, Local: p.local, Peer: p.peer})
	}

	switch p.state {
	case rOpen:
		// Process-DWR
		a, _ := p.makeDWA(v.m)
		// R-Snd-DWA
		go func() {
			p.rcon.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
			_, e := a.WriteTo(p.rcon)
			if Notificator != nil {
				Notificator(&WatchdogEvent{
					Tx: true, Req: false, Local: p.local, Peer: p.peer, Err: e})
			}
			if e != nil {
				p.notify <- eventRPeerDisc{}
			}
		}()
	case iOpen:
		// Process-DWR
		a, _ := p.makeDWA(v.m)
		// I-Snd-DWA
		go func() {
			p.icon.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
			_, e := a.WriteTo(p.icon)
			if Notificator != nil {
				Notificator(&WatchdogEvent{
					Tx: true, Req: false, Local: p.local, Peer: p.peer, Err: e})
			}
			if e != nil {
				p.notify <- eventIPeerDisc{}
			}
		}()
	default:
		// closed
		// wait_conn_ack
		// wait_i_CEA
		// wait_conn_ack_elect
		// wait_returns
		// closing
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventRcvDWA struct {
	m msg.Message
}

func (v eventRcvDWA) name() string {
	return "Rcv-DWA"
}

func (v eventRcvDWA) exec(p *Provider) (e error) {
	switch p.state {
	case rOpen, iOpen:
		// Process-DWA
		if ch, ok := p.sndstack[v.m.HbHID]; ok {
			if Notificator != nil {
				Notificator(&WatchdogEvent{
					Tx: false, Req: false, Local: p.local, Peer: p.peer})
			}
			// readed message is stacked answer
			delete(p.sndstack, v.m.HbHID)
			ch <- &v.m
		} else {
			if Notificator != nil {
				Notificator(&WatchdogEvent{
					Tx: false, Req: false, Local: p.local, Peer: p.peer,
					Err: fmt.Errorf("unknown DWA recieved")})
			}
		}
	default:
		// closed
		// wait_conn_ack
		// wait_i_CEA
		// wait_conn_ack_elect
		// wait_returns
		// i_open
		// closing
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventStop struct {
	c msg.Enumerated
}

func (v eventStop) name() string {
	return "Stop"
}

func (v eventStop) exec(p *Provider) (e error) {
	switch p.state {
	case rOpen:
		// R-Snd-DPR
		r := p.makeDPR(v.c)
		r.HbHID = p.local.NextHbH()

		p.state = closing
		go func() {
			p.rcon.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
			_, e := r.WriteTo(p.rcon)
			if Notificator != nil {
				Notificator(&PurgeEvent{
					Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
			}
			if e != nil {
				p.notify <- eventRPeerDisc{}
			}
		}()
	case iOpen:
		// I-Snd-DPR
		r := p.makeDPR(v.c)
		r.HbHID = p.local.NextHbH()

		p.state = closing
		go func() {
			p.icon.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
			_, e := r.WriteTo(p.icon)
			if Notificator != nil {
				Notificator(&PurgeEvent{
					Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
			}
			if e != nil {
				p.notify <- eventIPeerDisc{}
			}
		}()
	default:
		// closed
		// wait_conn_ack
		// wait_i_CEA
		// wait_conn_ack_elect
		// wait_returns
		// r_open
		// i_open
		// closing
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventRcvDPR struct {
	m msg.Message
}

func (v eventRcvDPR) name() string {
	return "Rcv-DPR"
}

func (v eventRcvDPR) exec(p *Provider) (e error) {
	if Notificator != nil {
		Notificator(&PurgeEvent{
			Tx: false, Req: true, Local: p.local, Peer: p.peer})
	}

	switch p.state {
	case rOpen:
		// R-Snd-DPA
		a, _ := p.makeDPA(v.m)
		p.state = closing
		go func() {
			p.rcon.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
			_, e := a.WriteTo(p.rcon)
			if Notificator != nil {
				Notificator(&PurgeEvent{
					Tx: true, Req: false, Local: p.local, Peer: p.peer, Err: e})
			}
			if e != nil {
				p.notify <- eventRPeerDisc{}
			}
		}()
	case iOpen:
		// I-Snd-DPA
		a, _ := p.makeDPA(v.m)
		p.state = closing
		go func() {
			p.icon.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
			_, e := a.WriteTo(p.icon)
			if Notificator != nil {
				Notificator(&PurgeEvent{
					Tx: true, Req: false, Local: p.local, Peer: p.peer, Err: e})
			}
			if e != nil {
				p.notify <- eventIPeerDisc{}
			}
		}()
	default:
		// closed
		// wait_conn_ack
		// wait_i_CEA
		// wait_conn_ack_elect
		// wait_returns
		// closing
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventRRcvDPA struct {
	m msg.Message
}

func (v eventRRcvDPA) name() string {
	return "R-Rcv-DPA"
}

func (v eventRRcvDPA) exec(p *Provider) (e error) {
	if Notificator != nil {
		Notificator(&PurgeEvent{
			Tx: false, Req: false, Local: p.local, Peer: p.peer})
	}

	switch p.state {
	case closing:
		// R-Disc
		p.rcon.Close()
		p.rcon = nil
		p.state = shutdown
	default:
		// closed
		// wait_conn_ack
		// wait_i_CEA
		// wait_conn_ack_elect
		// wait_returns
		// r_open
		// i_open
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventIRcvDPA struct {
	m msg.Message
}

func (v eventIRcvDPA) name() string {
	return "I-Rcv-DPA"
}

func (v eventIRcvDPA) exec(p *Provider) (e error) {
	if Notificator != nil {
		Notificator(&PurgeEvent{
			Tx: false, Req: false, Local: p.local, Peer: p.peer})
	}

	switch p.state {
	case closing:
		// I-Disc
		p.icon.Close()
		p.icon = nil
		p.state = shutdown
	default:
		// closed
		// wait_conn_ack
		// wait_i_CEA
		// wait_conn_ack_elect
		// wait_returns
		// r_open
		// i_open
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventWinElection struct{}

func (v eventWinElection) name() string {
	return "Win-Election"
}

func (v eventWinElection) exec(p *Provider) (e error) {
	switch p.state {
	case waitReturns:
		// I-Disc
		p.icon.Close()
		p.icon = nil

		p.state = rOpen
		// R-Snd-CEA
		go func() {
			p.rcon.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
			_, e := p.cachemsg.WriteTo(p.rcon)
			if Notificator != nil {
				Notificator(&CapabilityExchangeEvent{
					Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
			}
			if e != nil {
				p.notify <- eventRPeerDisc{}
			}
		}()
	default:
		// closed
		// wait_conn_ack
		// wait_i_CEA
		// wait_conn_ack_elect
		// r_open
		// i_open
		// closing
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventSndMsg struct {
	m msg.Message
}

func (v eventSndMsg) name() string {
	return "Send-Message"
}

func (v eventSndMsg) exec(p *Provider) (e error) {
	switch p.state {
	case rOpen:
		go func() {
			p.rcon.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
			_, e := v.m.WriteTo(p.rcon)
			if Notificator != nil {
				Notificator(&MessageEvent{
					Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
			}
			if e != nil {
				p.notify <- eventRPeerDisc{}
			}
		}()
	case iOpen:
		go func() {
			p.icon.SetWriteDeadline(time.Now().Add(p.peer.Properties.Ts))
			_, e := v.m.WriteTo(p.icon)
			if Notificator != nil {
				Notificator(&MessageEvent{
					Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
			}
			if e != nil {
				p.notify <- eventIPeerDisc{}
			}
		}()
	default:
		// closed
		// wait_conn_ack
		// wait_i_CEA
		// wait_conn_ack_elect
		// wait_returns
		// closing
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}

type eventRcvMsg struct {
	m msg.Message
}

func (v eventRcvMsg) name() string {
	return "Rcv-Message"
}

func (v eventRcvMsg) exec(p *Provider) (e error) {
	switch p.state {
	case rOpen, iOpen:
		// Process
		if v.m.FlgR {
			// readed message is request
			if Notificator != nil {
				Notificator(&MessageEvent{
					Tx: false, Req: true, Local: p.local, Peer: p.peer})
			}
			p.rcvstack <- &v.m
		} else if ch, ok := p.sndstack[v.m.HbHID]; ok {
			// readed message is stacked answer
			if Notificator != nil {
				Notificator(&MessageEvent{
					Tx: false, Req: false, Local: p.local, Peer: p.peer})
			}
			ch <- &v.m
		} else {
			if Notificator != nil {
				Notificator(&MessageEvent{
					Tx: false, Req: false, Local: p.local, Peer: p.peer,
					Err: fmt.Errorf("unknown answer message recieved")})
			}
		}
	default:
		// closed
		// wait_conn_ack
		// wait_i_CEA
		// wait_conn_ack_elect
		// wait_returns
		// closing
		// shutdown

		e = fmt.Errorf("not acceptable event")
	}
	return
}
