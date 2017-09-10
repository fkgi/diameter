package sock

import (
	"fmt"
	"time"

	"github.com/fkgi/diameter/msg"
)

// UnknownIDAnswer is error
type UnknownIDAnswer struct {
	msg.Message
}

func (e UnknownIDAnswer) Error() string {
	return "Unknown message recieved"
}

// FailureAnswer is error
type FailureAnswer struct {
	msg.Message
}

func (e FailureAnswer) Error() string {
	return "Answer message with failure recieved"
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
					c.notify <- eventWatchdog{}
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
	if _, ok := c.sndstack[v.m.HbHID]; !ok {
		return UnknownIDAnswer{v.m}
	}
	delete(c.sndstack, v.m.HbHID)
	c.wTimer.Stop()

	cea := msg.CEA{}
	e := cea.Decode(v.m)
	if e == nil {
		if cea.ResultCode == msg.DiameterSuccess {
			HandleCEA(cea, c)
			c.state = open
			c.wTimer = time.AfterFunc(c.peer.WDInterval, func() {
				c.notify <- eventWatchdog{}
			})
		} else {
			e = FailureAnswer{v.m}
		}
	}

	//notify(&CapabilityExchangeEvent{
	//	Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
	if e != nil {
		c.con.Close()
	}
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
	if _, ok := c.sndstack[v.m.HbHID]; !ok {
		return UnknownIDAnswer{v.m}
	}
	delete(c.sndstack, v.m.HbHID)
	c.wTimer.Stop()

	dwa := msg.DWA{}
	e := dwa.Decode(v.m)
	if e == nil {
		if dwa.ResultCode == msg.DiameterSuccess {
			HandleDWA(dwa, c)
			c.wCounter = 0
		} else {
			e = FailureAnswer{v.m}
		}
	}
	//notify(&WatchdogEvent{
	//	Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
	if c.wCounter > c.peer.WDExpired {
		c.con.Close()
	} else {
		c.wTimer = time.AfterFunc(c.peer.WDInterval, func() {
			c.notify <- eventWatchdog{}
		})
	}
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
	if _, ok := c.sndstack[v.m.HbHID]; !ok {
		return UnknownIDAnswer{v.m}
	}
	delete(c.sndstack, v.m.HbHID)
	c.wTimer.Stop()

	dpa := msg.DPA{}
	e := dpa.Decode(v.m)
	if e == nil {
		if dpa.ResultCode == msg.DiameterSuccess {
			HandleDPA(dpa, c)
		} else {
			e = FailureAnswer{v.m}
		}
	}
	//notify(&PurgeEvent{
	//	Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
	c.con.Close()
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
		return NotAcceptableEvent{event: v, state: c.state}
	}
	c.wTimer.Stop()

	if v.m.FlgR {
		HandleMSG(v.m, c)
		c.wCounter = 0
		c.wTimer = time.AfterFunc(c.peer.WDInterval, func() {
			c.notify <- eventWatchdog{}
		})
	} else if ch, ok := c.sndstack[v.m.HbHID]; ok {
		delete(c.sndstack, v.m.HbHID)
		ch <- v.m
		c.wCounter = 0
		c.wTimer = time.AfterFunc(c.peer.WDInterval, func() {
			c.notify <- eventWatchdog{}
		})
	} else {
		return UnknownIDAnswer{v.m}
	}

	//notify(&MessageEvent{
	//	Tx: false, Req: v.m.FlgR, Local: p.local, Peer: p.peer, Err: e})
	return
}
