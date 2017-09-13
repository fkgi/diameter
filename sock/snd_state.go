package sock

import (
	"fmt"

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

var timeoutMsg msg.ErrorMessage = "no response from peer node"

// NotAcceptableEvent is error
type NotAcceptableEvent struct {
	stateEvent
	state
}

func (e NotAcceptableEvent) Error() string {
	return fmt.Sprintf("Event %s is not acceptable in state %s",
		e.stateEvent, e.state)
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
type eventConnect struct{}

func (eventConnect) String() string {
	return "Connect"
}

func (v eventConnect) exec(c *Conn) error {
	if c.state != closed {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	c.state = waitCEA

	req := MakeCER(c)
	nak := msg.CEA{
		ResultCode:                  msg.DiameterUnableToDeliver,
		OriginHost:                  req.OriginHost,
		OriginRealm:                 req.OriginRealm,
		HostIPAddress:               req.HostIPAddress,
		VendorID:                    req.VendorID,
		ProductName:                 req.ProductName,
		OriginStateID:               req.OriginStateID,
		ErrorMessage:                &timeoutMsg,
		SupportedVendorID:           req.SupportedVendorID,
		ApplicationID:               req.ApplicationID,
		VendorSpecificApplicationID: req.VendorSpecificApplicationID,
		FirmwareRevision:            req.FirmwareRevision}
	e := c.sendSysMsg(req.Encode(), nak.Encode())

	Notify(CapabilityExchangeEvent{tx: true, req: true, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
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

	c.wCounter++
	if c.wCounter > c.peer.WDExpired {
		c.con.Close()
		return nil
	}

	req := MakeDWR(c)
	nak := msg.DWA{
		ResultCode:    msg.DiameterUnableToDeliver,
		OriginHost:    req.OriginHost,
		OriginRealm:   req.OriginRealm,
		ErrorMessage:  &timeoutMsg,
		OriginStateID: req.OriginStateID}
	e := c.sendSysMsg(req.Encode(), nak.Encode())

	Notify(WatchdogEvent{tx: true, req: true, conn: c, Err: e})
	if e != nil {
		c.con.Close()
		//} else {
		//	c.wTimer = time.AfterFunc(c.peer.WDInterval, func() {
		//		c.notify <- eventWatchdog{}
		//	})
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
	c.wTimer.Stop()

	req := MakeDPR(c)
	nak := msg.DPA{
		ResultCode:   msg.DiameterUnableToDeliver,
		OriginHost:   req.OriginHost,
		OriginRealm:  req.OriginRealm,
		ErrorMessage: &timeoutMsg}
	e := c.sendSysMsg(req.Encode(), nak.Encode())

	Notify(PurgeEvent{tx: true, req: true, conn: c, Err: e})
	if e != nil {
		c.con.Close()
		//	} else {
		//		c.wTimer = time.AfterFunc(
		//			c.peer.SndTimeout,
		//			func() {
		//				delete(c.sndstack, m.HbHID)
		//				c.con.Close()
		//notify(&PurgeEvent{
		//	Tx: false, Req: false, Local: c.local, Peer: c.peer,
		//	Err: fmt.Errorf("no answer")})
		//			})
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

// Snd MSG
type eventSndMsg struct {
	m msg.Message
}

func (eventSndMsg) String() string {
	return "Snd-MSG"
}

func (v eventSndMsg) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	c.setTransportDeadline()
	_, e := v.m.WriteTo(c.con)

	Notify(MessageEvent{tx: true, req: v.m.FlgR, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}
