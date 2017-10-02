package sock

import (
	"time"

	"github.com/fkgi/diameter/msg"
)

// UnknownIDAnswer is error
type UnknownIDAnswer struct {
	msg.RawMsg
}

func (e UnknownIDAnswer) Error() string {
	return "Unknown message recieved"
}

// FailureAnswer is error
type FailureAnswer struct {
	msg.RawMsg
}

func (e FailureAnswer) Error() string {
	return "Answer message with failure"
}

// RcvCER
type eventRcvCER struct {
	m msg.RawMsg
}

func (eventRcvCER) String() string {
	return "Rcv-CER"
}

func (v eventRcvCER) exec(c *Conn) error {
	if c.state != waitCER {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	cer, e := msg.CER{}.FromRaw(v.m)
	Notify(CapabilityExchangeEvent{tx: false, req: true, conn: c, Err: e})

	if e != nil {
		// ToDo
		// make error answer for undecodable CER
		return e
	}

	cea := HandleCER(cer.(msg.CER), c)
	m := cea.ToRaw()
	m.HbHID = v.m.HbHID
	m.EtEID = v.m.EtEID
	if cea.ResultCode != msg.DiameterSuccess {
		m.FlgE = true
	}
	e = c.write(m)

	if e == nil && cea.ResultCode != msg.DiameterSuccess {
		e = FailureAnswer{m}
	}
	if e == nil {
		c.state = open
		c.sysTimer = time.AfterFunc(c.peer.WDInterval, func() {
			c.notify <- eventWatchdog{}
		})
	}

	Notify(CapabilityExchangeEvent{tx: true, req: false, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

// RcvCEA
type eventRcvCEA struct {
	m msg.RawMsg
}

func (eventRcvCEA) String() string {
	return "Rcv-CEA"
}

func (v eventRcvCEA) exec(c *Conn) error {
	if c.state != waitCEA {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}
	if _, ok := c.sndstack[v.m.HbHID]; !ok {
		return UnknownIDAnswer{v.m}
	}
	delete(c.sndstack, v.m.HbHID)
	c.sysTimer.Stop()

	cea, e := msg.CEA{}.FromRaw(v.m)
	if e == nil {
		HandleCEA(cea.(msg.CEA), c)
		if cea.Result() == msg.DiameterSuccess {
			c.state = open
			c.sysTimer = time.AfterFunc(c.peer.WDInterval, func() {
				c.notify <- eventWatchdog{}
			})
		} else {
			e = FailureAnswer{v.m}
		}
	}

	Notify(CapabilityExchangeEvent{tx: false, req: false, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

type eventRcvDWR struct {
	m msg.RawMsg
}

func (eventRcvDWR) String() string {
	return "Rcv-DWR"
}

func (v eventRcvDWR) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	dwr, e := msg.DWR{}.FromRaw(v.m)
	Notify(WatchdogEvent{tx: false, req: true, conn: c, Err: e})

	if e != nil {
		// ToDo
		// make error answer for undecodable CER
		return e
	}

	dwa := HandleDWR(dwr.(msg.DWR), c)
	m := dwa.ToRaw()
	m.HbHID = v.m.HbHID
	m.EtEID = v.m.EtEID
	if dwa.ResultCode != msg.DiameterSuccess {
		m.FlgE = true
	}
	e = c.write(m)

	if e == nil && dwa.ResultCode != msg.DiameterSuccess {
		e = FailureAnswer{m}
	}
	if e == nil {
		c.sysTimer.Reset(c.peer.WDInterval)
	}

	Notify(WatchdogEvent{tx: true, req: false, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

// RcvDWA
type eventRcvDWA struct {
	m msg.RawMsg
}

func (eventRcvDWA) String() string {
	return "Rcv-DWA"
}

func (v eventRcvDWA) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}
	if _, ok := c.sndstack[v.m.HbHID]; !ok {
		return UnknownIDAnswer{v.m}
	}

	dwa, e := msg.DWA{}.FromRaw(v.m)
	c.sndstack[v.m.HbHID] <- dwa
	delete(c.sndstack, v.m.HbHID)

	if e == nil {
		HandleDWA(dwa.(msg.DWA), c)
		if dwa.Result() == msg.DiameterSuccess {
			c.wdCounter = 0
			c.sysTimer.Reset(c.peer.WDInterval)
		} else {
			e = FailureAnswer{v.m}
		}
	}

	Notify(WatchdogEvent{tx: false, req: false, conn: c, Err: e})
	return e
}

type eventRcvDPR struct {
	m msg.RawMsg
}

func (eventRcvDPR) String() string {
	return "Rcv-DPR"
}

func (v eventRcvDPR) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	dpr, e := msg.DPR{}.FromRaw(v.m)
	Notify(PurgeEvent{tx: false, req: true, conn: c, Err: e})

	if e != nil {
		// ToDo
		// make error answer for undecodable CER
		return e
	}

	dpa := HandleDPR(dpr.(msg.DPR), c)
	m := dpa.ToRaw()
	m.HbHID = v.m.HbHID
	m.EtEID = v.m.EtEID
	if dpa.ResultCode != msg.DiameterSuccess {
		m.FlgE = true
	} else {
		c.state = closing
		c.sysTimer.Stop()
		c.sysTimer = time.AfterFunc(c.peer.SndTimeout, func() {
			c.con.Close()
		})
	}
	e = c.write(m)

	Notify(&PurgeEvent{tx: true, req: false, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

type eventRcvDPA struct {
	m msg.RawMsg
}

func (eventRcvDPA) String() string {
	return "Rcv-DPA"
}

func (v eventRcvDPA) exec(c *Conn) error {
	if c.state != closing {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}
	if _, ok := c.sndstack[v.m.HbHID]; !ok {
		return UnknownIDAnswer{v.m}
	}
	delete(c.sndstack, v.m.HbHID)
	c.sysTimer.Stop()

	dpa, e := msg.DPA{}.FromRaw(v.m)
	if e == nil {
		HandleDPA(dpa.(msg.DPA), c)
		if dpa.Result() != msg.DiameterSuccess {
			e = FailureAnswer{v.m}
		}
	}

	Notify(PurgeEvent{tx: false, req: false, conn: c, Err: e})
	c.con.Close()
	return e
}

type eventRcvRequest struct {
	m msg.RawMsg
}

func (eventRcvRequest) String() string {
	return "Rcv-REQ"
}

func (v eventRcvRequest) exec(c *Conn) (e error) {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	if app, ok := supportedApps[msg.AuthApplicationID(v.m.AppID)]; !ok {
		go func() {
			c.notify <- eventSndAnswer{MakeUnsupportedAnswer(v.m)}
		}()
		e = msg.UnknownApplicationID{}
	} else if req, ok := app.req[v.m.Code]; !ok {
		go func() {
			c.notify <- eventSndAnswer{MakeUnsupportedAnswer(v.m)}
		}()
		e = msg.InvalidMessage{}
	} else {
		go func() {
			if m, e := req.FromRaw(v.m); e != nil {
				// ToDo
				// invalid message handling
			} else if ans := HandleMSG(m); ans == nil {
				// ToDo
				// message handling failure handling
			} else {
				a := ans.ToRaw()
				a.HbHID = v.m.HbHID
				a.EtEID = v.m.EtEID
				c.notify <- eventSndAnswer{a}
			}
		}()
	}
	c.sysTimer.Reset(c.peer.WDInterval)

	Notify(MessageEvent{tx: false, req: true, conn: c, Err: e})
	return
}

type eventRcvAnswer struct {
	m msg.RawMsg
}

func (eventRcvAnswer) String() string {
	return "Rcv-ANS"
}

func (v eventRcvAnswer) exec(c *Conn) (e error) {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	if ch, ok := c.sndstack[v.m.HbHID]; ok {
		delete(c.sndstack, v.m.HbHID)

		if app, ok := supportedApps[msg.AuthApplicationID(v.m.AppID)]; !ok {
			// ToDo
			// invalid message handling
			e = msg.UnknownApplicationID{}
		} else if ans, ok := app.ans[v.m.Code]; !ok {
			// ToDo
			// invalid message handling
			e = msg.InvalidMessage{}
		} else if m, e2 := ans.FromRaw(v.m); e2 != nil {
			// ToDo
			// invalid message handling
			e = e2
		} else {
			ch <- m
		}
		c.sysTimer.Reset(c.peer.WDInterval)
	} else {
		e = UnknownIDAnswer{v.m}
	}

	Notify(MessageEvent{tx: false, req: false, conn: c, Err: e})
	return
}
