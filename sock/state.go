package sock

import (
	"fmt"
	"time"

	"github.com/fkgi/diameter/msg"
)

var (
	stateMap = map[int]string{
		shutdown: "shutdown",
		closed:   "closed",
		waitCER:  "waitCER",
		waitCEA:  "waitCEA",
		open:     "open",
		closing:  "closing",
	}
)

const (
	shutdown = iota
	closed   = iota
	waitCER  = iota
	waitCEA  = iota
	open     = iota
	closing  = iota
)

// NotAcceptableEvent is error
type NotAcceptableEvent struct {
	event stateEvent
	state int
}

func (e NotAcceptableEvent) Error() string {
	return "Event " + e.event.String() +
		" is not acceptable in state " + stateMap[e.state]
}

type stateEvent interface {
	exec(p *Conn) error
	String() string
}

// Connect
type eventConnect struct {
	m msg.Message
}

func (eventConnect) String() string {
	return "Connect"
}

func (v eventConnect) exec(c *Conn) error {
	if c.state != closed {
		return NotAcceptableEvent{event: v, state: c.state}
	}

	c.state = waitCEA
	c.setTransportDeadline()
	_, e := v.m.WriteTo(c.con)

	//notify(&CapabilityExchangeEvent{
	//	Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

// Watchdog
type eventWatchdog struct {
	m msg.Message
}

func (eventWatchdog) String() string {
	return "Watchdog"
}

func (v eventWatchdog) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{event: v, state: c.state}
	}

	c.wCounter++
	if c.wCounter > c.peer.WDExpired {
		c.con.Close()
		return nil
	}

	v.m.HbHID = c.local.NextHbH()
	v.m.EtEID = c.local.NextEtE()

	// ch := make(chan msg.Message)
	// c.sndstack[v.m.HbHID] = ch

	c.setTransportDeadline()
	_, e := v.m.WriteTo(c.con)

	// notify(&WatchdogEvent{
	//	Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
	if e != nil {
		c.con.Close()
	} else {
		c.wTimer = time.AfterFunc(c.peer.WDInterval, func() {
			dwr := MakeDWR(c)
			c.notify <- eventWatchdog{dwr.Encode()}
		})
	}
	return e
}

// Stop
type eventStop struct {
	m msg.Message
}

func (eventStop) String() string {
	return "Stop"
}

func (v eventStop) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{event: v, state: c.state}
	}

	c.state = closing
	c.wTimer.Stop()

	c.setTransportDeadline()
	_, e := v.m.WriteTo(c.con)

	// notify(&PurgeEvent{
	// 	Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

// PeerDisc
type eventPeerDisc struct{}

func (eventPeerDisc) String() string {
	return "Peer-Disc"
}

func (v eventPeerDisc) exec(c *Conn) error {
	c.con.Close()
	c.state = closed

	// notify(&DisconnectEvent{
	// 	Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
	return nil
}

// RcvCER
type eventRcvCER struct {
	m msg.Message
}

func (eventRcvCER) String() string {
	return "Rcv-CER"
}

func (v eventRcvCER) exec(c *Conn) error {
	if c.state != waitCER {
		return NotAcceptableEvent{event: v, state: c.state}
	}

	cer := msg.CER{}
	e := cer.Decode(v.m)
	//notify(&CapabilityExchangeEvent{
	//	Tx: false, Req: true, Local: p.local, Peer: peer, Err: e})

	if e == nil {
		cea := HandleCER(cer, c)

		m := cea.Encode()
		m.HbHID = c.local.NextHbH()
		m.EtEID = c.local.NextEtE()
		if cea.ResultCode != msg.DiameterSuccess {
			m.FlgE = true
		}

		c.setTransportDeadline()
		_, e = m.WriteTo(c.con)

		if e == nil {
			if cea.ResultCode != msg.DiameterSuccess {
				e = fmt.Errorf("close with error response %d", cea.ResultCode)
				c.con.Close()
			} else {
				c.state = open
				c.wTimer = time.AfterFunc(c.peer.WDInterval, func() {
					dwr := MakeDWR(c)
					c.notify <- eventWatchdog{dwr.Encode()}
				})
			}
		}
	}
	//notify(&CapabilityExchangeEvent{
	//	Tx: true, Req: false, Local: p.local, Peer: p.peer})
	return e
}

// RcvCEA
type eventRcvCEA struct {
	m msg.Message
}

func (eventRcvCEA) String() string {
	return "Rcv-CEA"
}

func (v eventRcvCEA) exec(c *Conn) error {
	if c.state != waitCEA {
		return NotAcceptableEvent{event: v, state: c.state}
	}

	cea := msg.CEA{}
	e := cea.Decode(v.m)

	if e == nil {
		HandleCEA(cea, c)

		if cea.ResultCode == msg.DiameterSuccess {
			c.state = open
			c.wTimer = time.AfterFunc(c.peer.WDInterval, func() {
				dwr := MakeDWR(c)
				c.notify <- eventWatchdog{dwr.Encode()}
			})
		} else {
			e = fmt.Errorf("CEA Nack received")
			c.con.Close()
		}
	}
	//notify(&CapabilityExchangeEvent{
	//	Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
	return e
}

type eventRcvDWR struct {
	m msg.Message
}

func (eventRcvDWR) String() string {
	return "Rcv-DWR"
}

func (v eventRcvDWR) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{event: v, state: c.state}
	}

	dwr := msg.DWR{}
	e := dwr.Decode(v.m)
	//notify(&WatchdogEvent{
	//	Tx: false, Req: true, Local: p.local, Peer: p.peer, Err: e})

	if e == nil {
		dwa := HandleDWR(dwr, c)

		m := dwa.Encode()
		m.HbHID = c.local.NextHbH()
		m.EtEID = c.local.NextEtE()
		if dwa.ResultCode != msg.DiameterSuccess {
			m.FlgE = true
		}

		c.setTransportDeadline()
		_, e = m.WriteTo(c.con)
		// notify(&WatchdogEvent{
		//	Tx: true, Req: false, Local: p.local, Peer: p.peer, Err: e})
	}
	if e != nil {
		c.con.Close()
		c.state = shutdown
	} else {
		c.wCounter = 0
		c.wTimer.Reset(c.peer.WDInterval)
	}
	return e
}

// RcvDWA
type eventRcvDWA struct {
	m msg.Message
}

func (eventRcvDWA) String() string {
	return "Rcv-DWA"
}

func (v eventRcvDWA) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{event: v, state: c.state}
	}

	dwa := msg.DWA{}
	e := dwa.Decode(v.m)

	if e == nil {
		HandleDWA(dwa, c)
		if ch, ok := c.sndstack[v.m.HbHID]; ok {
			delete(c.sndstack, v.m.HbHID)
			ch <- v.m

			c.wCounter = 0
			c.wTimer.Reset(c.peer.WDInterval)
		} else {
			e = fmt.Errorf("unknown DWA recieved")
		}
	}
	//notify(&WatchdogEvent{
	//	Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
	return e
}

type eventRcvDPR struct {
	m msg.Message
}

func (eventRcvDPR) String() string {
	return "Rcv-DPR"
}

func (v eventRcvDPR) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{event: v, state: c.state}
	}

	dpr := msg.DPR{}
	e := dpr.Decode(v.m)
	//notify(&PurgeEvent{
	//	Tx: false, Req: true, Local: p.local, Peer: p.peer, Err: e})

	if e == nil {
		dpa := HandleDPR(dpr, c)
		m := dpa.Encode()
		m.HbHID = c.local.NextHbH()
		m.EtEID = c.local.NextEtE()
		if dpa.ResultCode != msg.DiameterSuccess {
			m.FlgE = true
		} else {
			c.state = closing
			c.wTimer.Stop()
		}

		c.setTransportDeadline()
		_, e = m.WriteTo(c.con)
		//notify(&PurgeEvent{
		//	Tx: true, Req: false, Local: p.local, Peer: p.peer, Err: e})
		if e != nil {
			c.con.Close()
		}
	}
	return e
}

type eventRcvDPA struct {
	m msg.Message
}

func (eventRcvDPA) String() string {
	return "Rcv-DPA"
}

func (v eventRcvDPA) exec(c *Conn) error {
	if c.state != closing {
		return NotAcceptableEvent{event: v, state: c.state}
	}

	dpa := msg.DPA{}
	e := dpa.Decode(v.m)

	if e == nil {
		HandleDPA(dpa, c)
		if ch, ok := c.sndstack[v.m.HbHID]; ok {
			ch <- v.m
			// p.con.Close()
		} else {
			e = fmt.Errorf("unknown DPA recieved")
		}
	}
	//notify(&PurgeEvent{
	//	Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
	return e
}

type eventSndMsg struct {
	m msg.Message
}

func (eventSndMsg) String() string {
	return "Snd-MSG"
}

func (v eventSndMsg) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{event: v, state: c.state}
	}

	c.setTransportDeadline()
	_, e := v.m.WriteTo(c.con)
	//notify(&MessageEvent{
	//	Tx: true, Req: v.m.FlgR, Local: p.local, Peer: p.peer, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

type eventRcvMsg struct {
	m msg.Message
}

func (eventRcvMsg) String() string {
	return "Rcv-MSG"
}

func (v eventRcvMsg) exec(c *Conn) (e error) {
	if c.state != open {
		return NotAcceptableEvent{
			event: v,
			state: c.state}
	}

	if v.m.FlgR {
		HandleMSG(v.m, c)
		c.wCounter = 0
		c.wTimer.Reset(c.peer.WDInterval)
	} else if ch, ok := c.sndstack[v.m.HbHID]; ok {
		ch <- v.m
		c.wCounter = 0
		c.wTimer.Reset(c.peer.WDInterval)
	} else {
		e = fmt.Errorf("unknown answer message received")
	}

	//notify(&MessageEvent{
	//	Tx: false, Req: v.m.FlgR, Local: p.local, Peer: p.peer, Err: e})
	return
}
