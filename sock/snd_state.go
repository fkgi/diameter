package sock

import (
	"fmt"
	"time"

	"github.com/fkgi/diameter/msg"
)

type state int

func (s state) String() string {
	switch s {
	case shutdown:
		return "shutdown"
	case closed:
		return "closed"
	case waitCER:
		return "waitCER"
	case waitCEA:
		return "waitCEA"
	case open:
		return "open"
	case closing:
		return "closing"
	}
	return "<nil>"
}

const (
	shutdown state = iota
	closed   state = iota
	waitCER  state = iota
	waitCEA  state = iota
	open     state = iota
	closing  state = iota
)

// NotAcceptableEvent is error
type NotAcceptableEvent struct {
	stateEvent
	state
}

func (e NotAcceptableEvent) Error() string {
	return fmt.Sprintf("Event %s is not acceptable in state %s",
		e.stateEvent, e.state)
}

// WatchdogExpired is error
type WatchdogExpired struct{}

func (e WatchdogExpired) Error() string {
	return "watchdog is expired"
}

// ConnectionRefused is error
type ConnectionRefused struct{}

func (e ConnectionRefused) Error() string {
	return "connection is refused"
}

type stateEvent interface {
	exec(p *Conn) error
	String() string
}

// Init
type eventInit struct{}

func (eventInit) String() string {
	return "Initialize"
}

func (v eventInit) exec(c *Conn) error {
	return NotAcceptableEvent{stateEvent: v, state: c.state}
}

// Connect
type eventConnect struct {
	m msg.RawMsg
}

func (eventConnect) String() string {
	return "Connect"
}

func (v eventConnect) exec(c *Conn) error {
	if c.state != closed {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	c.state = waitCEA

	e := c.write(v.m)
	Notify(CapabilityExchangeEvent{tx: true, req: true, conn: c, Err: e})
	if e == nil {
		c.sysTimer = time.AfterFunc(c.peer.SndTimeout, func() {
			cea, _ := msg.CER{}.FromRaw(v.m)
			nak := cea.(msg.CER).TimeoutMsg().ToRaw()
			nak.HbHID = v.m.HbHID
			nak.EtEID = v.m.EtEID
			c.notify <- eventRcvCEA{nak}
		})
	} else {
		c.con.Close()
	}
	return e
}

func watchdog(c *Conn) {
	dwr := MakeDWR(c)
	req := dwr.ToRaw()
	req.HbHID = NextHbH()
	req.EtEID = NextEtE()

	ch := make(chan msg.Answer)
	c.sndstack[req.HbHID] = ch

	c.notify <- eventWatchdog{}

	t := time.AfterFunc(c.peer.SndTimeout, func() {
		dwa := dwr.TimeoutMsg()
		nak := dwa.ToRaw()
		nak.HbHID = req.HbHID
		nak.EtEID = req.EtEID
		c.notify <- eventRcvDWA{nak}
	})
	<-ch
	t.Stop()
}

// Watchdog
type eventWatchdog struct{}

func (eventWatchdog) String() string {
	return "Watchdog"
}

func (v eventWatchdog) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	c.wdCounter++
	if c.wdCounter > c.peer.WDExpired {
		c.con.Close()
		return WatchdogExpired{}
	}

	dwr := MakeDWR(c)
	req := dwr.ToRaw()
	req.HbHID = NextHbH()
	req.EtEID = NextEtE()
	ch := make(chan msg.Answer)
	c.sndstack[req.HbHID] = ch

	e := c.write(req)

	Notify(WatchdogEvent{tx: true, req: true, conn: c, Err: e})
	if e == nil {
		go func() {
			t := time.AfterFunc(c.peer.SndTimeout, func() {
				dwa := dwr.TimeoutMsg()
				nak := dwa.ToRaw()
				nak.HbHID = req.HbHID
				nak.EtEID = req.EtEID
				c.notify <- eventRcvDWA{nak}
			})
			<-ch
			t.Stop()
		}()
	} else {
		c.con.Close()
	}
	return e
}

// Stop
type eventStop struct{}

func (eventStop) String() string {
	return "Stop"
}

func (v eventStop) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	c.state = closing
	c.sysTimer.Stop()

	dpr := MakeDPR(c)
	req := dpr.ToRaw()
	req.HbHID = NextHbH()
	req.EtEID = NextEtE()
	c.sndstack[req.HbHID] = nil //make(chan msg.Message)

	e := c.write(req)
	Notify(PurgeEvent{tx: true, req: true, conn: c, Err: e})
	if e == nil {
		c.sysTimer = time.AfterFunc(c.peer.SndTimeout, func() {
			dpa := dpr.TimeoutMsg()
			nak := dpa.ToRaw()
			nak.HbHID = req.HbHID
			nak.EtEID = req.EtEID
			c.notify <- eventRcvDPA{nak}
		})
	} else {
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

// Snd Request MSG
type eventSndRequest struct {
	m msg.RawMsg
}

func (eventSndRequest) String() string {
	return "Snd-REQ"
}

func (v eventSndRequest) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	e := c.write(v.m)
	Notify(MessageEvent{tx: true, req: true, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

// Snd Answer MSG
type eventSndAnswer struct {
	m msg.RawMsg
}

func (eventSndAnswer) String() string {
	return "Snd-ANS"
}

func (v eventSndAnswer) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	e := c.write(v.m)
	Notify(MessageEvent{tx: true, req: false, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}
