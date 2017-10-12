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

	c.con.SetWriteDeadline(time.Now().Add(TransportTimeout))
	_, e := v.m.WriteTo(c.con)
	Notify(CapabilityExchangeEvent{tx: true, req: true, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

// Watchdog
type eventWatchdog struct {
	m msg.RawMsg
}

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

	c.con.SetWriteDeadline(time.Now().Add(TransportTimeout))
	_, e := v.m.WriteTo(c.con)
	Notify(WatchdogEvent{tx: true, req: true, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

// Stop
type eventStop struct {
	m msg.RawMsg
}

func (eventStop) String() string {
	return "Stop"
}

func (v eventStop) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	c.state = closing
	c.sysTimer.Stop()

	c.con.SetWriteDeadline(time.Now().Add(TransportTimeout))
	_, e := v.m.WriteTo(c.con)
	Notify(PurgeEvent{tx: true, req: true, conn: c, Err: e})
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

	for _, ch := range c.sndstack {
		ch <- msg.RawMsg{}
	}
	c.rcvstack <- msg.RawMsg{}

	return nil
}

// Snd MSG
type eventSndMsg struct {
	m msg.RawMsg
}

func (eventSndMsg) String() string {
	return "Snd-REQ"
}

func (v eventSndMsg) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	c.con.SetWriteDeadline(time.Now().Add(TransportTimeout))
	_, e := v.m.WriteTo(c.con)
	Notify(MessageEvent{tx: true, req: v.m.FlgR, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}
