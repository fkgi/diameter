package diameter

import (
	"time"
)

// RcvCER
type eventRcvCER struct {
	m RawMsg
}

func (eventRcvCER) String() string {
	return "Rcv-CER"
}

func (v eventRcvCER) exec(c *Conn) error {
	if c.state != waitCER {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	cer, _, e := CER{}.FromRaw(v.m)
	Notify(CapabilityExchangeEvent{tx: false, req: true, conn: c, Err: e})

	if e != nil {
		// ToDo
		// make error answer for undecodable CER
		c.con.Close()
		return e
	}

	cea := HandleCER(cer.(CER), c)
	m := cea.ToRaw("")
	m.HbHID = v.m.HbHID
	m.EtEID = v.m.EtEID
	if cea.ResultCode != DiameterSuccess {
		m.FlgE = true
	}
	c.con.SetWriteDeadline(time.Now().Add(TransportTimeout))
	_, e = m.WriteTo(c.con)

	if e == nil && cea.ResultCode != DiameterSuccess {
		e = FailureAnswer{m}
	}
	if e == nil {
		c.state = open
		c.wdTimer = time.AfterFunc(c.Peer.WDInterval, c.watchdog)
	}

	Notify(CapabilityExchangeEvent{tx: true, req: false, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

// RcvCEA
type eventRcvCEA struct {
	m RawMsg
}

func (eventRcvCEA) String() string {
	return "Rcv-CEA"
}

func (v eventRcvCEA) exec(c *Conn) error {
	if c.state != waitCEA {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}
	ch, ok := c.sndstack[v.m.HbHID]
	if !ok {
		return UnknownIDAnswer{v.m}
	}
	delete(c.sndstack, v.m.HbHID)

	cea, _, e := CEA{}.FromRaw(v.m)
	if e == nil {
		HandleCEA(cea.(CEA), c)
		if cea.Result() == uint32(DiameterSuccess) {
			c.state = open
			c.wdTimer = time.AfterFunc(c.Peer.WDInterval, c.watchdog)
		} else {
			e = FailureAnswer{v.m}
		}
	}

	Notify(CapabilityExchangeEvent{tx: false, req: false, conn: c, Err: e})
	if e != nil {
		c.con.Close()
		v.m = RawMsg{}
	}
	ch <- v.m
	return e
}

type eventRcvDWR struct {
	m RawMsg
}

func (eventRcvDWR) String() string {
	return "Rcv-DWR"
}

func (v eventRcvDWR) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	dwr, _, e := DWR{}.FromRaw(v.m)
	Notify(WatchdogEvent{tx: false, req: true, conn: c, Err: e})

	if e != nil {
		// ToDo
		// make error answer for undecodable CER
		return e
	}

	dwa := HandleDWR(dwr.(DWR), c)
	m := dwa.ToRaw("")
	m.HbHID = v.m.HbHID
	m.EtEID = v.m.EtEID
	if dwa.ResultCode != DiameterSuccess {
		m.FlgE = true
	}
	c.con.SetWriteDeadline(time.Now().Add(TransportTimeout))
	_, e = m.WriteTo(c.con)

	if e == nil && dwa.ResultCode != DiameterSuccess {
		e = FailureAnswer{m}
	}
	if e == nil {
		c.wdTimer.Reset(c.Peer.WDInterval)
	}

	Notify(WatchdogEvent{tx: true, req: false, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

// RcvDWA
type eventRcvDWA struct {
	m RawMsg
}

func (eventRcvDWA) String() string {
	return "Rcv-DWA"
}

func (v eventRcvDWA) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}
	ch, ok := c.sndstack[v.m.HbHID]
	if !ok {
		return UnknownIDAnswer{v.m}
	}
	delete(c.sndstack, v.m.HbHID)

	dwa, _, e := DWA{}.FromRaw(v.m)
	if e == nil {
		HandleDWA(dwa.(DWA), c)
		if dwa.Result() == uint32(DiameterSuccess) {
			c.wdCount = 0
			c.wdTimer.Reset(c.Peer.WDInterval)
		} else {
			e = FailureAnswer{v.m}
		}
	}

	Notify(WatchdogEvent{tx: false, req: false, conn: c, Err: e})
	if e != nil {
		v.m = RawMsg{}
	}
	ch <- v.m
	return e
}

type eventRcvDPR struct {
	m RawMsg
}

func (eventRcvDPR) String() string {
	return "Rcv-DPR"
}

func (v eventRcvDPR) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	dpr, _, e := DPR{}.FromRaw(v.m)
	Notify(PurgeEvent{tx: false, req: true, conn: c, Err: e})

	if e != nil {
		// ToDo
		// make error answer for undecodable CER
		return e
	}

	dpa := HandleDPR(dpr.(DPR), c)
	m := dpa.ToRaw("")
	m.HbHID = v.m.HbHID
	m.EtEID = v.m.EtEID
	if dpa.ResultCode != DiameterSuccess {
		m.FlgE = true
	} else {
		c.state = closing
		c.wdTimer.Stop()
		c.wdTimer = time.AfterFunc(TransportTimeout, func() {
			c.con.Close()
		})
	}
	c.con.SetWriteDeadline(time.Now().Add(TransportTimeout))
	_, e = m.WriteTo(c.con)

	Notify(&PurgeEvent{tx: true, req: false, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

type eventRcvDPA struct {
	m RawMsg
}

func (eventRcvDPA) String() string {
	return "Rcv-DPA"
}

func (v eventRcvDPA) exec(c *Conn) error {
	if c.state != closing {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}
	ch, ok := c.sndstack[v.m.HbHID]
	if !ok {
		return UnknownIDAnswer{v.m}
	}
	delete(c.sndstack, v.m.HbHID)

	dpa, _, e := DPA{}.FromRaw(v.m)
	if e == nil {
		HandleDPA(dpa.(DPA), c)
		if dpa.Result() != uint32(DiameterSuccess) {
			e = FailureAnswer{v.m}
		}
	}

	Notify(PurgeEvent{tx: false, req: false, conn: c, Err: e})
	c.con.Close()
	if e != nil {
		v.m = RawMsg{}
	}
	ch <- v.m
	return e
}

type eventRcvMsg struct {
	m RawMsg
}

func (eventRcvMsg) String() string {
	return "Rcv-MSG"
}

func (v eventRcvMsg) exec(c *Conn) (e error) {
	if c.state != open {
		return NotAcceptableEvent{stateEvent: v, state: c.state}
	}

	if v.m.FlgR {
		var cause uint32

		if app, ok := supportedApps[v.m.AppID]; !ok {
			cause = DiameterApplicationUnsupported
		} else if _, ok = app.req[v.m.Code]; !ok {
			cause = DiameterCommandUnspported
		}

		if cause == 0 {
		} else if app, ok := supportedApps[0xffffffff]; !ok {
		} else if _, ok = app.req[0]; ok {
			cause = 0
		}

		if cause == 0 {
			req, sid, _ := GenericReq{}.FromRaw(v.m)
			a := req.Failed(cause).ToRaw(sid)
			a.HbHID = v.m.HbHID
			a.EtEID = v.m.EtEID
			c.con.SetWriteDeadline(time.Now().Add(TransportTimeout))
			_, e = a.WriteTo(c.con)
		} else {
			c.rcvstack <- v.m
		}
	} else {
		ch, ok := c.sndstack[v.m.HbHID]
		if !ok {
			return
		}
		delete(c.sndstack, v.m.HbHID)
		ch <- v.m
	}
	c.wdTimer.Reset(c.Peer.WDInterval)

	Notify(MessageEvent{tx: false, req: true, conn: c, Err: e})
	if e != nil {
		c.con.Close()
	}
	return
}
