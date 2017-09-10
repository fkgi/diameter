package sock

import (
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

var timeoutMsg msg.ErrorMessage = "no response from peer node"

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

func (c *Conn) sendSysMsg(req msg.Message, nak msg.Message) error {
	req.HbHID = c.local.NextHbH()
	nak.HbHID = req.HbHID
	req.EtEID = c.local.NextEtE()
	nak.EtEID = req.EtEID
	c.sndstack[req.HbHID] = nil //make(chan msg.Message)

	c.setTransportDeadline()
	_, e := req.WriteTo(c.con)

	if e != nil {
		return e
	}
	c.wTimer = time.AfterFunc(c.peer.SndTimeout, func() {
		switch nak.Code {
		case 257:
			c.notify <- eventRcvCEA{nak}
		case 282:
			c.notify <- eventRcvDPA{nak}
		case 280:
			c.notify <- eventRcvDWA{nak}
		}
	})
	return nil
}

// Connect
type eventConnect struct{}

func (eventConnect) String() string {
	return "Connect"
}

func (v eventConnect) exec(c *Conn) error {
	if c.state != closed {
		return NotAcceptableEvent{event: v, state: c.state}
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
	// c.con.Close()
	//notify(&PurgeEvent{
	//	Tx: false, Req: false, Local: c.local, Peer: c.peer,
	//	Err: fmt.Errorf("no answer")})
	//})

	//notify(&CapabilityExchangeEvent{
	//	Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
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
		return NotAcceptableEvent{event: v, state: c.state}
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

	// notify(&WatchdogEvent{
	//	Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
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
		return NotAcceptableEvent{event: v, state: c.state}
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

	// notify(&PurgeEvent{
	// 	Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
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
