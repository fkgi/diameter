package diameter

import (
	"time"
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
	closed
	waitCER
	waitCEA
	open
	closing
)

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
	m RawMsg
}

func (eventConnect) String() string {
	return "Connect"
}

func (v eventConnect) exec(c *Conn) error {
	if c.state != closed {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}
	c.state = waitCEA

	c.TxReq++
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
	m RawMsg
}

func (eventWatchdog) String() string {
	return "Watchdog"
}

func (v eventWatchdog) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	c.wdCount++
	if c.wdCount > c.Peer.WDExpired {
		c.con.Close()
		return WatchdogExpired{}
	}

	c.TxReq++
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
	m RawMsg
}

func (eventStop) String() string {
	return "Stop"
}

func (v eventStop) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	c.state = closing
	c.wdTimer.Stop()
	c.Since = time.Time{}

	c.TxReq++
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
	c.Since = time.Time{}

	for _, ch := range c.sndstack {
		ch <- RawMsg{}
	}
	c.rcvstack <- RawMsg{}

	return nil
}

// Snd MSG
type eventSndMsg struct {
	m RawMsg
}

func (eventSndMsg) String() string {
	return "Snd-MSG"
}

func (v eventSndMsg) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	c.TxReq++
	c.con.SetWriteDeadline(time.Now().Add(TransportTimeout))
	_, e := v.m.WriteTo(c.con)
	Notify(MessageEvent{tx: true, req: v.m.FlgR, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}
