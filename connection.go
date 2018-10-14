package diameter

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"
)

// constant values
var (
	TxBuffer = 65535
	RxBuffer = 65535

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

	Since        time.Time
	RxReq        uint64
	Reject       uint64
	Tx1xxx       uint64
	Tx2xxx       uint64
	Tx3xxx       uint64
	Tx4xxx       uint64
	Tx5xxx       uint64
	TxEtc        uint64
	TxReq        uint64
	TxReqFail    uint64
	TxReqTimeout uint64
	Rx1xxx       uint64
	Rx2xxx       uint64
	Rx3xxx       uint64
	Rx4xxx       uint64
	Rx5xxx       uint64
	RxEtc        uint64
}

func (c *Conn) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sPeer                =%s\n", Indent, c.Peer)
	fmt.Fprintf(w, "%sStatus              =%s\n", Indent, c.State())
	fmt.Fprintf(w, "%sUptime              =%s\n",
		Indent, time.Now().Sub(c.Since).String())
	fmt.Fprintf(w, "%sRx Request count    =%d\n", Indent, c.RxReq)
	fmt.Fprintf(w, "%s%sReject count  =%d\n", Indent, Indent, c.Reject)
	fmt.Fprintf(w, "%s%sTx 1xxx count =%d\n", Indent, Indent, c.Tx1xxx)
	fmt.Fprintf(w, "%s%sTx 2xxx count =%d\n", Indent, Indent, c.Tx2xxx)
	fmt.Fprintf(w, "%s%sTx 3xxx count =%d\n", Indent, Indent, c.Tx3xxx)
	fmt.Fprintf(w, "%s%sTx 4xxx count =%d\n", Indent, Indent, c.Tx4xxx)
	fmt.Fprintf(w, "%s%sTx 5xxx count =%d\n", Indent, Indent, c.Tx5xxx)
	fmt.Fprintf(w, "%s%sTx etc count  =%d\n", Indent, Indent, c.TxEtc)
	fmt.Fprintf(w, "%sTx Request count    =%d\n", Indent, c.TxReq)
	fmt.Fprintf(w, "%s%sFailed count  =%d\n", Indent, Indent, c.TxReqFail)
	fmt.Fprintf(w, "%s%sTimeout count =%d\n", Indent, Indent, c.TxReqTimeout)
	fmt.Fprintf(w, "%s%sRx 1xxx count =%d\n", Indent, Indent, c.Rx1xxx)
	fmt.Fprintf(w, "%s%sRx 2xxx count =%d\n", Indent, Indent, c.Rx2xxx)
	fmt.Fprintf(w, "%s%sRx 3xxx count =%d\n", Indent, Indent, c.Rx3xxx)
	fmt.Fprintf(w, "%s%sRx 4xxx count =%d\n", Indent, Indent, c.Rx4xxx)
	fmt.Fprintf(w, "%s%sRx 5xxx count =%d\n", Indent, Indent, c.Rx5xxx)
	fmt.Fprintf(w, "%s%sRx etc count  =%d\n", Indent, Indent, c.RxEtc)
	fmt.Fprintf(w, "%sRx queue length     =%d\n", Indent, c.RxQueue())
	fmt.Fprintf(w, "%sTx queue length     =%d\n", Indent, c.TxQueue())

	return w.String()
}

// RxQueue returns length of Rx queue
func (c *Conn) RxQueue() int {
	return len(c.rcvstack)
}

// TxQueue returns length of Tx queue
func (c *Conn) TxQueue() int {
	return len(c.sndstack)
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

	con := &Conn{
		Peer:     &p,
		notify:   make(chan stateEvent),
		state:    closed,
		con:      c,
		sndstack: make(map[uint32]chan RawMsg, TxBuffer),
		rcvstack: make(chan RawMsg, RxBuffer)}
	go socketHandler(con)
	Notify(StateUpdate{
		oldStat: shutdown, newStat: con.state,
		stateEvent: eventInit{}, conn: con, Err: nil})
	go eventHandler(con)

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
	con := &Conn{
		Peer:     p,
		notify:   make(chan stateEvent),
		state:    waitCER,
		con:      c,
		sndstack: make(map[uint32]chan RawMsg, TxBuffer),
		rcvstack: make(chan RawMsg, RxBuffer)}
	go socketHandler(con)

	Notify(StateUpdate{
		oldStat: shutdown, newStat: con.state,
		stateEvent: eventInit{}, conn: con, Err: nil})

	event := <-con.notify
	old := con.state
	e := event.exec(con)
	Notify(StateUpdate{
		oldStat: old, newStat: con.state,
		stateEvent: event, conn: con, Err: e})
	if e != nil {
		c.Close()
	}
	go eventHandler(con)

	return con, e
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
	for {
		event := <-c.notify
		old := c.state
		e := event.exec(c)

		Notify(StateUpdate{
			oldStat: old, newStat: c.state,
			stateEvent: event, conn: c, Err: e})

		if _, ok := event.(eventPeerDisc); ok {
			break
		}
	}
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
	} else if ans, ok := app.ans[a.Code]; !ok {
	} else if ack, _, e := ans.FromRaw(a); e == nil {
		return ack
	} else if avperr, ok := e.(InvalidAVP); ok {
		return m.Failed(uint32(avperr))
	} else {
		return m.Failed(DiameterUnableToComply)
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

	if req == nil {
		app, _ := supportedApps[0xffffffff]
		req, _ = app.req[0]
	}

	r, sid, e := req.FromRaw(m)
	f := func(ans Answer) {
		a := ans.ToRaw(sid)
		a.HbHID = m.HbHID
		a.EtEID = m.EtEID
		c.notify <- eventSndMsg{a}
	}
	if e != nil {
		if avperr, ok := e.(InvalidAVP); ok {
			f(req.Failed(uint32(avperr)))
		} else {
			f(req.Failed(DiameterUnableToComply))
		}
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
	if c == nil || c.state != open {
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

// State returns state machine state
func (c *Conn) State() string {
	return c.state.String()
}
