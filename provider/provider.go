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
	ProductName      = "mave"
	FirmwareRevision = uint32(160114001)
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

var stateStr = map[int]string{
	closed:           "Closed",
	waitConnAck:      "Wait-Conn-Ack",
	waitICEA:         "Wait-I-CEA",
	waitConnAckElect: "Wait-Conn-Ack-Elect",
	waitReturns:      "Wait-Returns",
	rOpen:            "R-Open",
	iOpen:            "I-Open",
	closing:          "Closing",
	shutdown:         "Shutdown",
}

// Provider is state machine of Diameter
type Provider struct {
	wT *time.Timer // watchdog timer
	wE int         // watchdog expired counter

	notify     chan stateEvent
	state      int
	icon, rcon *Connection
	rcvstack   chan *msg.Message
	sndstack   map[uint32]chan *msg.Message

	cachemsg msg.Message
}

// Open start state machine
func (p *Provider) Open() {
	if p.state == shutdown {
		p.state = closed
	}
}

// Close stop state machine
func (p *Provider) Close(cause msg.Enumerated) {
	if p.state == rOpen || p.state == iOpen {
		p.notify <- eventStop{cause}
	} else if p.state == closed {
		p.state = shutdown
	}
}

// State returns status of state machine
func (p *Provider) State() string {
	return stateStr[p.state]
}

func (p *Provider) run() {
	for {
		if event := <-p.notify; event == nil {
			break
		} else if e := event.exec(p); e != nil {
			monitor.Notify(monitor.Info,
				"not acceptable event for provider status",
				e.Error(), stateStr[p.state])
		}
	}
}

func (p *Provider) open() {
	monitor.Notify(monitor.Info, "start message handling process")

	p.resetWatchdog()

	for p.state == rOpen {
		m, e := p.rcon.Read(0)
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
		m, e := p.icon.Read(0)
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

	monitor.Notify(monitor.Info, "stop message handling process")
}

func (p *Provider) activeConnection() *Connection {
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
		if p.wE > c.Peer.Ew {
			p.notify <- eventStop{msg.Enumerated(0)}
		} else {
			r := c.makeDWR()
			monitor.Notify(monitor.Debug, "-> DWR")
			ap := sendReq(r, c, p)
			if ap == nil {
				monitor.Notify(monitor.Debug, "no DWA response")
				p.notify <- eventStop{msg.Enumerated(0)}
			} else {
				monitor.Notify(monitor.Debug, "<- DWA")
				p.resetWatchdog()
			}
		}
	}

	c := p.activeConnection()
	if c == nil {
		return
	}

	if p.wT != nil {
		p.wT.Reset(c.Peer.Tw)
	} else {
		p.wT = time.AfterFunc(c.Peer.Tw, f)
	}
}

// Send send Diameter request message
func (p *Provider) Send(r msg.Message) (a msg.Message, e error) {
	c := p.activeConnection()
	if c == nil {
		e = fmt.Errorf("connection is not open")
		return
	}
	r.EtEID = c.Local.NextEtE()

	for i := 0; i <= c.Peer.Cp; i++ {
		monitor.Notify(monitor.Debug, "-> REQ")
		ap := sendReq(r, c, p)
		if ap == nil {
			if i == c.Peer.Cp {
				monitor.Notify(monitor.Debug, "no response")
				e = fmt.Errorf("No answer")
			} else {
				monitor.Notify(monitor.Debug, "no response retry")
			}
			r.FlgT = true
		} else {
			monitor.Notify(monitor.Debug, "<- ANS")
			a = *ap
			e = nil
			break
		}
	}
	return
}

func sendReq(r msg.Message, c *Connection, p *Provider) (a *msg.Message) {
	ch := make(chan *msg.Message)

	r.HbHID = c.Local.NextHbH()
	p.sndstack[r.HbHID] = ch

	p.notify <- eventSndMsg{r}
	t := time.AfterFunc(c.Peer.Tp, func() {
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
		monitor.Notify(monitor.Debug, "<- REQ")
		r = *rp
		ch = make(chan *msg.Message)
		go func() {
			if a := <-ch; a == nil {
				monitor.Notify(monitor.Debug, "not send response")
			} else {
				monitor.Notify(monitor.Debug, "-> ANS")
				a.HbHID = r.HbHID
				a.EtEID = r.EtEID
				p.notify <- eventSndMsg{*a}
			}
		}()
	}
	return
}

type stateEvent interface {
	exec(p *Provider) error
}

type eventStart struct {
	src, dst net.Addr
	l        *LocalNode
	p        *PeerNode
}

func (v eventStart) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event Start")

	switch p.state {
	case closed:
		// I-Snd-Conn-Req
		p.state = waitConnAck
		go func() {
			c, e := v.l.Connect(v.p, v.src, v.dst, v.p.Ts)
			p.icon = c
			if e == nil {
				p.notify <- eventIRcvConAck{}
			} else {
				p.notify <- eventIRcvConNack{}
			}
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

		e = fmt.Errorf("Start")
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
	c *Connection
}

func (v eventRConnCER) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Debug, "<- CER")
	monitor.Notify(monitor.Trace, "event R-Conn-CER")

	switch p.state {
	case closed:
		// R-Accept, Process-CER
		m, c := v.c.makeCEA(v.m)

		// R-Snd-CEA
		monitor.Notify(monitor.Debug, "-> CEA")
		e = v.c.Write(v.c.Peer.Ts, m)

		if c == 2001 && e == nil {
			p.rcon = v.c
			p.state = rOpen
			go p.open()
		} else {
			// Close
			monitor.Notify(monitor.Info, "DEA send failed")
			v.c.Close()
		}
	case waitConnAck:
		p.state = waitConnAckElect
		// R-Accept, Process-CER
		p.cachemsg, _ = v.c.makeCEA(v.m)
		p.rcon = v.c
	case waitICEA:
		p.state = waitReturns
		// R-Accept, Process-CER
		p.cachemsg, _ = v.c.makeCEA(v.m)
		p.rcon = v.c
		// Elect
		if msg.Compare(p.rcon.Local.Host, p.rcon.Peer.Host) > 0 {
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
		e = fmt.Errorf("R-Conn-CER")
	}
	return
}

type eventIRcvConAck struct{}

func (v eventIRcvConAck) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event I-Rcv-Con-Ack")

	switch p.state {
	case waitConnAck:
		p.state = waitICEA
		// I-Snd-CER
		r := p.icon.makeCER()
		r.HbHID = p.icon.Local.NextHbH()

		go func() {
			monitor.Notify(monitor.Debug, "-> CER")
			if e = p.icon.Write(p.icon.Peer.Ts, r); e != nil {
				p.notify <- eventIPeerDisc{}
			} else if a, e := p.icon.Read(0); e != nil {
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
		r := p.icon.makeCER()
		r.HbHID = p.icon.Local.NextHbH()

		go func() {
			monitor.Notify(monitor.Debug, "-> CER")
			if e = p.icon.Write(p.icon.Peer.Ts, r); e != nil {
				p.notify <- eventIPeerDisc{}
			} else if msg.Compare(p.rcon.Local.Host, p.rcon.Peer.Host) > 0 {
				// Elect
				p.notify <- eventWinElection{}
			} else if a, e := p.icon.Read(0); e != nil {
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

		e = fmt.Errorf("I-Rcv-Conn-Ack")
	}
	return
}

type eventIRcvConNack struct{}

func (v eventIRcvConNack) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event I-Rcv-Con-Nack")

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
			monitor.Notify(monitor.Debug, "-> CEA")
			if e = p.rcon.Write(p.rcon.Peer.Ts, p.cachemsg); e != nil {
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

		e = fmt.Errorf("I-Rcv-Conn-Nack")
	}
	return
}

type eventTimeout struct {
	s string
}

func (v eventTimeout) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event Timeout")

	switch p.state {
	case waitConnAck, waitICEA, waitConnAckElect, waitReturns, closing:
		// Error
		monitor.Notify(monitor.Major, "request timeout", v.s)
		p.state = closed
	default:
		// closed
		// r_open
		// i_open
		// shutdown

		e = fmt.Errorf("Timeout")
	}
	return
}

type eventRcvCER struct {
	m msg.Message
}

func (v eventRcvCER) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event Rcv-CER")
	monitor.Notify(monitor.Debug, "<- CER")

	e = fmt.Errorf("Rcv-CER")
	return
}

type eventRRcvCEA struct {
	m msg.Message
}

func (v eventRRcvCEA) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event R-Rcv-CEA")
	monitor.Notify(monitor.Debug, "<- CEA")

	e = fmt.Errorf("R-Rcv-CEA")
	return
}

type eventIRcvCEA struct {
	m msg.Message
}

func (v eventIRcvCEA) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event I-Rcv-CEA")
	monitor.Notify(monitor.Debug, "<- CEA")

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
			p.state = iOpen
			go p.open()
		} else {
			// Close
			monitor.Notify(monitor.Info, "CEA Nack received")
			p.icon.Close()
			p.icon = nil
			p.state = closed
		}
	case waitReturns:
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

		e = fmt.Errorf("I-Rcv-CEA")
	}
	return
}

type eventIRcvNonCEA struct {
	m msg.Message
}

func (v eventIRcvNonCEA) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event I-Rcv-Non-CEA")
	monitor.Notify(monitor.Debug, "<- ANS")

	switch p.state {
	case waitICEA:
		// Error
		p.icon.Close()
		p.icon = nil
		monitor.Notify(monitor.Major, "None CEA received")
		p.state = closed
	default:
		// closed
		// wait_conn_ack
		// wait_conn_ack_elect
		// wait_returns
		// r_open
		// i_open
		// closing
		// shutdown

		e = fmt.Errorf("I-Rcv-Non-CEA")
	}
	return
}

type eventRPeerDisc struct{}

func (v eventRPeerDisc) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event R-Peer-Disc")

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
		p.state = closed
	default:
		// closed
		// wait_conn_ack
		// wait_i_CEA
		// i_open
		// shutdown

		e = fmt.Errorf("R-Peer-Disc")
	}
	return
}

type eventIPeerDisc struct{}

func (v eventIPeerDisc) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event I-Peer-Disc")

	switch p.state {
	case waitICEA, iOpen, closing:
		// I-Disc
		p.icon.Close()
		p.icon = nil
		p.state = closed
	case waitReturns:
		// I-Disc
		p.icon.Close()
		p.icon = nil

		p.state = rOpen
		go p.open()
		// R-Snd-CEA
		go func() {
			monitor.Notify(monitor.Debug, "-> CEA")
			if e = p.rcon.Write(p.rcon.Peer.Ts, p.cachemsg); e != nil {
				p.notify <- eventRPeerDisc{}
			}
		}()
	default:
		// closed
		// wait_conn_ack
		// wait_conn_ack_elect
		// r_open
		// shutdown

		e = fmt.Errorf("I-Peer-Disc")
	}
	return
}

type eventRcvDWR struct {
	m msg.Message
}

func (v eventRcvDWR) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event Rcv-DWR")
	monitor.Notify(monitor.Debug, "<- DWR")

	switch p.state {
	case rOpen:
		// Process-DWR
		a, _ := p.rcon.makeDWA(v.m)
		// R-Snd-DWA
		go func() {
			monitor.Notify(monitor.Debug, "-> DWA")
			if e = p.rcon.Write(p.rcon.Peer.Ts, a); e != nil {
				p.notify <- eventRPeerDisc{}
			}
		}()
	case iOpen:
		// Process-DWR
		a, _ := p.icon.makeDWA(v.m)
		// I-Snd-DWA
		go func() {
			monitor.Notify(monitor.Debug, "-> DWA")
			if e = p.icon.Write(p.icon.Peer.Ts, a); e != nil {
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

		e = fmt.Errorf("Rcv-DWR")
	}
	return
}

type eventRcvDWA struct {
	m msg.Message
}

func (v eventRcvDWA) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event Rcv-DWA")

	switch p.state {
	case rOpen, iOpen:
		// Process-DWA
		if ch, ok := p.sndstack[v.m.HbHID]; ok {
			// readed message is stacked answer
			delete(p.sndstack, v.m.HbHID)
			ch <- &v.m
		} else {
			monitor.Notify(monitor.Debug, "unknown DWA recieved")
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

		e = fmt.Errorf("Rcv-DWA")
	}
	return
}

type eventStop struct {
	c msg.Enumerated
}

func (v eventStop) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event Stop")

	switch p.state {
	case rOpen:
		// R-Snd-DPR
		r := p.rcon.makeDPR(v.c)
		r.HbHID = p.rcon.Local.NextHbH()

		p.state = closing
		go func() {
			monitor.Notify(monitor.Debug, "-> DPR")
			if e = p.rcon.Write(p.rcon.Peer.Ts, r); e != nil {
				p.notify <- eventRPeerDisc{}
			}
		}()
	case iOpen:
		// I-Snd-DPR
		r := p.icon.makeDPR(v.c)
		r.HbHID = p.icon.Local.NextHbH()

		p.state = closing
		go func() {
			monitor.Notify(monitor.Debug, "-> DPR")
			if e = p.icon.Write(p.icon.Peer.Ts, r); e != nil {
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

		e = fmt.Errorf("Stop")
	}
	return
}

type eventRcvDPR struct {
	m msg.Message
}

func (v eventRcvDPR) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event Rcv-DPR")
	monitor.Notify(monitor.Debug, "<- DPR")

	switch p.state {
	case rOpen:
		// R-Snd-DPA
		a, _ := p.rcon.makeDPA(v.m)
		p.state = closing
		go func() {
			monitor.Notify(monitor.Debug, "-> DPA")
			if e = p.rcon.Write(p.rcon.Peer.Ts, a); e != nil {
				p.notify <- eventRPeerDisc{}
			}
		}()
	case iOpen:
		// I-Snd-DPA
		a, _ := p.icon.makeDPA(v.m)
		p.state = closing
		go func() {
			monitor.Notify(monitor.Debug, "-> DPA")
			if e = p.icon.Write(p.icon.Peer.Ts, a); e != nil {
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

		e = fmt.Errorf("Rcv-DPR")
	}
	return
}

type eventRRcvDPA struct {
	m msg.Message
}

func (v eventRRcvDPA) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event R-Rcv-DPA")
	monitor.Notify(monitor.Debug, "<- DPA")

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

		e = fmt.Errorf("R-Rcv-DPA")
	}
	return
}

type eventIRcvDPA struct {
	m msg.Message
}

func (v eventIRcvDPA) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event I-Rcv-DPA")
	monitor.Notify(monitor.Debug, "<- DPA")

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

		e = fmt.Errorf("I-Rcv-DPA")
	}
	return
}

type eventWinElection struct{}

func (v eventWinElection) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event Win-Election")

	switch p.state {
	case waitReturns:
		// I-Disc
		p.icon.Close()
		p.icon = nil

		p.state = rOpen
		// R-Snd-CEA
		go func() {
			monitor.Notify(monitor.Debug, "-> CER")
			if e = p.rcon.Write(p.rcon.Peer.Ts, p.cachemsg); e != nil {
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

		e = fmt.Errorf("Win-Election")
	}
	return
}

type eventSndMsg struct {
	m msg.Message
}

func (v eventSndMsg) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event Send-Message")

	switch p.state {
	case rOpen:
		go func() {
			if e = p.rcon.Write(p.rcon.Peer.Ts, v.m); e != nil {
				p.notify <- eventRPeerDisc{}
			}
		}()
	case iOpen:
		go func() {
			if e = p.icon.Write(p.icon.Peer.Ts, v.m); e != nil {
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

		e = fmt.Errorf("Send-Message")
	}
	return
}

type eventRcvMsg struct {
	m msg.Message
}

func (v eventRcvMsg) exec(p *Provider) (e error) {
	monitor.Notify(monitor.Trace, "event Rcv-Message")

	switch p.state {
	case rOpen, iOpen:
		// Process
		if v.m.FlgR {
			// readed message is request
			p.rcvstack <- &v.m
		} else if ch, ok := p.sndstack[v.m.HbHID]; ok {
			// readed message is stacked answer
			ch <- &v.m
		} else {
			monitor.Notify(monitor.Debug, "unknown answer message recieved")
		}
	default:
		// closed
		// wait_conn_ack
		// wait_i_CEA
		// wait_conn_ack_elect
		// wait_returns
		// closing
		// shutdown

		e = fmt.Errorf("Rcv-Message")
	}
	return
}
