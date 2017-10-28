package diameter

import (
	"net"
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
	peer *Peer

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
	con, e := run(&p, c, closed)
	if e == nil {
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
			e = ConnectionRefused{}
		}
	}
	return con, e
}

// Accept new transport connection and return Conn
func Accept(p *Peer, c net.Conn) (*Conn, error) {
	return run(p, c, waitCER)
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

// Send send Diameter request
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

func (c *Conn) watchdog() {
	dwr := MakeDWR(c)
	req := dwr.ToRaw("")
	req.HbHID = nextHbH()
	req.EtEID = nextEtE()

	ch := make(chan RawMsg)
	c.sndstack[req.HbHID] = ch
	c.notify <- eventWatchdog{m: req}

	t := time.AfterFunc(c.peer.WDInterval, func() {
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

// PeerHost returns peer host name
func (c *Conn) PeerHost() Identity {
	return c.peer.Host
}

// PeerRealm returns peer realm name
func (c *Conn) PeerRealm() Identity {
	return c.peer.Realm
}

// PeerAddr returns transport connection of state machine
func (c *Conn) PeerAddr() net.Addr {
	return c.con.RemoteAddr()
}

// AuthApplication returns application ID of this connection
func (c *Conn) AuthApplication() map[uint32][]uint32 {
	return c.peer.AuthApps
}

func run(p *Peer, c net.Conn, s state) (*Conn, error) {
	if c == nil {
		return nil, ConnectionRefused{}
	}
	con := &Conn{
		peer:     p,
		notify:   make(chan stateEvent),
		state:    s,
		con:      c,
		sndstack: make(map[uint32]chan RawMsg, TxBuffer),
		rcvstack: make(chan RawMsg, RxBuffer)}
	go socketHandler(con)
	go eventHandler(con)
	for i := 0; i < Workers; i++ {
		go func() {
			for {
				m := <-con.rcvstack
				if m.Code == 0 {
					con.rcvstack <- m
					break
				}
				con.notify <- eventSndMsg{messageHandler(m)}
			}
		}()
	}

	return con, nil
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

func messageHandler(m RawMsg) RawMsg {
	var ans Answer
	var sid string
	if app, ok := supportedApps[m.AppID]; !ok {
		ans = GenericReq(m).Failed(DiameterApplicationUnsupported)
	} else if req, ok := app.req[m.Code]; !ok {
		ans = GenericReq(m).Failed(DiameterCommandUnspported)
	} else {
		var e error
		if req, sid, e = req.FromRaw(m); e == nil {
			ans = HandleMSG(req)
		}
		if ans == nil {
			ans = GenericReq(m).Failed(DiameterUnableToComply)
			// ToDo
			// invalid message handling
		}
		a := ans.ToRaw(sid)
		a.HbHID = m.HbHID
		a.EtEID = m.EtEID
		return a
	}

	if app, ok := supportedApps[0xffffffff]; !ok {
	} else if req, ok := app.req[0]; ok {
		req, sid, _ = req.FromRaw(m)
		ans = HandleMSG(req)
		if ans == nil {
			ans = GenericReq(m).Failed(DiameterUnableToComply)
			// ToDo
			// invalid message handling
		}
	}
	a := ans.ToRaw(sid)
	a.HbHID = m.HbHID
	a.EtEID = m.EtEID
	return a
}
