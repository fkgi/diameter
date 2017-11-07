package diameter

import (
	"net"
	"strings"
	"time"
)

// constant values
var (
	TxBuffer = 65535
	RxBuffer = 65535
	Workers  = 1

	VendorID         uint32 = 41102
	ProductName             = "yatagarasu"
	FirmwareRevision uint32 = 170819001
)

// Conn is state machine of Diameter
type Conn struct {
	*Peer

	wdTimer *time.Timer // system message timer
	wdCount int         // watchdog expired counter

	notify chan stateEvent
	state
	con      net.Conn
	sndstack map[uint32]chan RawMsg
	rcvstack chan RawMsg
}

// Dial make new Conn that use specified peernode and connection
func Dial(p Peer, c net.Conn, d time.Duration) (*Conn, error) {
	if c == nil {
		return nil, ConnectionRefused{}
	}
	if len(p.Host) == 0 {
		return nil, ConnectionRefused{}
	}
	if len(p.Realm) == 0 {
		var e error
		p.Realm, e = ParseIdentity(string(
			p.Host[strings.Index(string(p.Host), ".")+1:]))
		if e != nil {
			return nil, e
		}
	}
	if p.WDExpired == 0 {
		p.WDExpired = WDExpired
	}
	if p.WDInterval == 0 {
		p.WDInterval = WDInterval
	}

	con := run(&p, c, closed)
	cer := MakeCER(con)
	req := cer.ToRaw("")
	req.HbHID = nextHbH()
	req.EtEID = nextEtE()

	ch := make(chan RawMsg)
	con.sndstack[req.HbHID] = ch
	con.notify <- eventConnect{m: req}

	t := time.AfterFunc(d, func() {
		m := cer.Failed(DiameterTooBusy).ToRaw("")
		m.HbHID = req.HbHID
		m.EtEID = req.EtEID
		con.notify <- eventRcvCEA{m}
	})

	ack := <-ch
	t.Stop()

	if ack.Code == 0 {
		return nil, ConnectionRefused{}
	}
	return con, nil
}

// Accept new transport connection and return Conn
func Accept(p *Peer, c net.Conn) (*Conn, error) {
	if c == nil {
		return nil, ConnectionRefused{}
	}
	return run(p, c, waitCER), nil
}

// NewSession make new session
/*
func (c *Conn) NewSession() *Session {
	s := &Session{
		id: nextSession(),
		c:  c}
	return s
}
*/

// Send Diameter request
func (c *Conn) Send(m Request, d time.Duration) Answer {
	sid := nextSession()
	req := m.ToRaw(sid)
	req.HbHID = nextHbH()
	req.EtEID = nextEtE()

	ch := make(chan RawMsg)
	c.sndstack[req.HbHID] = ch
	c.notify <- eventSndMsg{m: req}

	t := time.AfterFunc(d, func() {
		m := m.Failed(DiameterTooBusy).ToRaw(sid)
		m.HbHID = req.HbHID
		m.EtEID = req.EtEID
		c.notify <- eventRcvMsg{m}
	})

	a := <-ch
	t.Stop()
	if a.Code == 0 {
		return m.Failed(DiameterUnableToDeliver)
	}

	if app, ok := supportedApps[a.AppID]; !ok {
		// ToDo
		// invalid application handling
	} else if ans, ok := app.ans[a.Code]; !ok {
		// ToDo
		// invalid code handling
	} else {
		ack, _, e := ans.FromRaw(a)
		if e != nil {
			// ToDo
			// invalid message handling
			ack = m.Failed(DiameterUnableToComply)
		}
		return ack
	}

	if app, ok := supportedApps[0xffffffff]; !ok {
	} else if ans, ok := app.ans[0]; ok {
		ack, _, _ := ans.FromRaw(a)
		return ack
	}

	return m.Failed(DiameterUnableToComply)
}

// Recieve Diameter request
func (c *Conn) Recieve() (Request, func(Answer), error) {
	m := <-c.rcvstack
	if m.Code == 0 {
		c.rcvstack <- m
		return nil, nil, ConnectionRefused{}
	}

	var req Request

	if app, ok := supportedApps[m.AppID]; ok {
		req, _ = app.req[m.Code]
	}

	if req != nil {
	} else if app, ok := supportedApps[0xffffffff]; ok {
		req, _ = app.req[0]
	}
	if req == nil {
		return nil, nil, ConnectionRefused{}
	}

	r, sid, e := req.FromRaw(m)
	f := func(ans Answer) {
		a := ans.ToRaw(sid)
		a.HbHID = m.HbHID
		a.EtEID = m.EtEID
		c.notify <- eventSndMsg{a}
	}
	if e != nil {
		ans := req.Failed(DiameterUnableToComply)
		f(ans)
		return r, nil, e
	}
	return r, f, nil
}

func (c *Conn) watchdog() {
	dwr := MakeDWR(c)
	req := dwr.ToRaw("")
	req.HbHID = nextHbH()
	req.EtEID = nextEtE()

	ch := make(chan RawMsg)
	c.sndstack[req.HbHID] = ch
	c.notify <- eventWatchdog{m: req}

	t := time.AfterFunc(c.Peer.WDInterval, func() {
		m := dwr.Failed(DiameterTooBusy).ToRaw("")
		m.HbHID = req.HbHID
		m.EtEID = req.EtEID
		c.notify <- eventRcvDWA{m}
	})

	<-ch
	t.Stop()
}

// Close stop state machine
func (c *Conn) Close(d time.Duration) {
	if c.state != open {
		return
	}

	dpr := MakeDPR(c)
	req := dpr.ToRaw("")
	req.HbHID = nextHbH()
	req.EtEID = nextEtE()

	ch := make(chan RawMsg)
	c.sndstack[req.HbHID] = ch
	c.notify <- eventStop{m: req}

	t := time.AfterFunc(d, func() {
		m := dpr.Failed(DiameterTooBusy).ToRaw("")
		m.HbHID = req.HbHID
		m.EtEID = req.EtEID
		c.notify <- eventRcvDPA{m}
	})

	<-ch
	t.Stop()
}

// LocalAddr returns transport connection of state machine
func (c *Conn) LocalAddr() net.Addr {
	return c.con.LocalAddr()
}

// PeerAddr returns transport connection of state machine
func (c *Conn) PeerAddr() net.Addr {
	return c.con.RemoteAddr()
}

func run(p *Peer, c net.Conn, s state) *Conn {
	con := &Conn{
		Peer:     p,
		notify:   make(chan stateEvent),
		state:    s,
		con:      c,
		sndstack: make(map[uint32]chan RawMsg, TxBuffer),
		rcvstack: make(chan RawMsg, RxBuffer)}

	go socketHandler(con)
	go eventHandler(con)
	return con
}

func socketHandler(c *Conn) {
	for {
		m := RawMsg{}
		c.con.SetReadDeadline(time.Time{})
		if _, e := m.ReadFrom(c.con); e != nil {
			break
		}

		if m.AppID == 0 && m.Code == 257 && m.FlgR {
			c.notify <- eventRcvCER{m}
		} else if m.AppID == 0 && m.Code == 257 && !m.FlgR {
			c.notify <- eventRcvCEA{m}
		} else if m.AppID == 0 && m.Code == 280 && m.FlgR {
			c.notify <- eventRcvDWR{m}
		} else if m.AppID == 0 && m.Code == 280 && !m.FlgR {
			c.notify <- eventRcvDWA{m}
		} else if m.AppID == 0 && m.Code == 282 && m.FlgR {
			c.notify <- eventRcvDPR{m}
		} else if m.AppID == 0 && m.Code == 282 && !m.FlgR {
			c.notify <- eventRcvDPA{m}
		} else {
			c.notify <- eventRcvMsg{m}
		}
	}
	c.notify <- eventPeerDisc{}
}

func eventHandler(c *Conn) {
	old := shutdown
	Notify(StateUpdate{
		oldStat: old, newStat: c.state,
		stateEvent: eventInit{}, conn: c, Err: nil})

	for {
		event := <-c.notify
		old = c.state
		e := event.exec(c)

		Notify(StateUpdate{
			oldStat: old, newStat: c.state,
			stateEvent: event, conn: c, Err: e})

		if _, ok := event.(eventPeerDisc); ok {
			break
		}
	}
}
