package sock

import (
	"net"
	"time"

	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/rfc6733"
)

// constant values
var (
	TxBuffer                = 65535
	RxBuffer                = 65535
	VendorID         uint32 = 41102
	ProductName             = "yatagarasu"
	FirmwareRevision uint32 = 170819001
)

// Conn is state machine of Diameter
type Conn struct {
	peer *Peer

	sysTimer  *time.Timer // system message timer
	wdCounter int         // watchdog expired counter

	notify chan stateEvent
	state
	con      net.Conn
	sndstack map[uint32]chan msg.RawMsg
	rcvstack chan msg.RawMsg
}

// Dial make new Conn that use specified peernode and connection
func Dial(p *Peer, c net.Conn) (*Conn, error) {
	con := &Conn{
		peer:     p,
		notify:   make(chan stateEvent),
		state:    closed,
		con:      c,
		sndstack: make(map[uint32]chan msg.RawMsg, TxBuffer),
		rcvstack: make(chan msg.RawMsg, RxBuffer)}
	go socketHandler(con)
	go eventHandler(con)
	go messageHandler(con)

	cer := MakeCER(con)
	req := cer.ToRaw()
	req.HbHID = nextHbH()
	req.EtEID = nextEtE()

	ch := make(chan msg.RawMsg)
	con.sndstack[req.HbHID] = ch
	con.notify <- eventConnect{m: req}

	t := time.AfterFunc(p.SndTimeout, func() {
		m := cer.Failed(
			uint32(rfc6733.DiameterTooBusy),
			"no response from peer node").ToRaw()
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
	con := &Conn{
		peer:     p,
		notify:   make(chan stateEvent),
		state:    waitCER,
		con:      c,
		sndstack: make(map[uint32]chan msg.RawMsg, TxBuffer),
		rcvstack: make(chan msg.RawMsg, RxBuffer)}
	go socketHandler(con)
	go eventHandler(con)
	go messageHandler(con)

	return con, nil
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
func (c *Conn) Send(m msg.Request) msg.Answer {
	req := m.ToRaw()
	req.HbHID = nextHbH()
	req.EtEID = nextEtE()

	ch := make(chan msg.RawMsg)
	c.sndstack[req.HbHID] = ch
	c.notify <- eventSndMsg{m: req}

	t := time.AfterFunc(c.peer.SndTimeout, func() {
		m := m.Failed(
			uint32(rfc6733.DiameterTooBusy),
			"no response from peer node").ToRaw()
		m.HbHID = req.HbHID
		m.EtEID = req.EtEID
		c.notify <- eventRcvMsg{m}
	})

	a := <-ch
	t.Stop()
	if a.Code == 0 {
		return m.Failed(
			uint32(rfc6733.DiameterUnableToDeliver),
			"failed to send")
	}

	app, ok := supportedApps[a.AppID]
	if !ok {
		return m.Failed(
			uint32(rfc6733.DiameterUnableToComply),
			"invalid Application-ID answer")
	}
	ans, ok := app.ans[a.Code]
	if !ok {
		return m.Failed(
			uint32(rfc6733.DiameterUnableToComply),
			"invalid Command-Code answer")
	}
	ack, e := ans.FromRaw(a)
	if e != nil {
		return m.Failed(
			uint32(rfc6733.DiameterUnableToComply),
			"invalid data answer")
		// ToDo
		// invalid message handling
	}
	return ack
}

func (c *Conn) watchdog() {
	dwr := MakeDWR(c)
	req := dwr.ToRaw()
	req.HbHID = nextHbH()
	req.EtEID = nextEtE()

	ch := make(chan msg.RawMsg)
	c.sndstack[req.HbHID] = ch
	c.notify <- eventWatchdog{m: req}

	t := time.AfterFunc(c.peer.SndTimeout, func() {
		m := dwr.Failed(
			uint32(rfc6733.DiameterTooBusy),
			"no response from peer node").ToRaw()
		m.HbHID = req.HbHID
		m.EtEID = req.EtEID
		c.notify <- eventRcvDWA{m}
	})

	<-ch
	t.Stop()
}

// Close stop state machine
func (c *Conn) Close(timeout time.Duration) {
	if c.state != open {
		return
	}

	dpr := MakeDPR(c)
	req := dpr.ToRaw()
	req.HbHID = nextHbH()
	req.EtEID = nextEtE()

	ch := make(chan msg.RawMsg)
	c.sndstack[req.HbHID] = ch
	c.notify <- eventStop{m: req}

	t := time.AfterFunc(c.peer.SndTimeout, func() {
		m := dpr.Failed(
			uint32(rfc6733.DiameterTooBusy),
			"no response from peer node").ToRaw()
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
func (c *Conn) PeerHost() msg.DiameterIdentity {
	return c.peer.Host
}

// PeerRealm returns peer realm name
func (c *Conn) PeerRealm() msg.DiameterIdentity {
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

func socketHandler(c *Conn) {
	for {
		m := msg.RawMsg{}
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

func messageHandler(c *Conn) {
	for {
		m := <-c.rcvstack
		if m.Code == 0 {
			break
		}
		if app, ok := supportedApps[m.AppID]; !ok {
			c.notify <- eventSndMsg{
				makeErrorMsg(m, rfc6733.DiameterApplicationUnsupported)}
		} else if req, ok := app.req[m.Code]; !ok {
			c.notify <- eventSndMsg{
				makeErrorMsg(m, rfc6733.DiameterCommandUnspported)}
		} else if r, e := req.FromRaw(m); e != nil {
			// ToDo
			// invalid message handling
		} else if ans := HandleMSG(r); ans == nil {
			// ToDo
			// message handling failure handling
		} else {
			a := ans.ToRaw()
			a.HbHID = m.HbHID
			a.EtEID = m.EtEID
			c.notify <- eventSndMsg{a}
		}
	}
}

func makeErrorMsg(m msg.RawMsg, c uint32) msg.RawMsg {
	r := msg.RawMsg{}
	r.Ver = m.Ver
	r.FlgP = m.FlgP
	r.Code = m.Code
	r.AppID = m.AppID
	r.HbHID = m.HbHID
	r.EtEID = m.EtEID

	host := msg.RawAVP{Code: 264, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	host.Encode(Host)
	realm := msg.RawAVP{Code: 296, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	realm.Encode(Realm)
	rcode := msg.RawAVP{Code: 268, VenID: 0,
		FlgV: false, FlgM: true, FlgP: false}
	rcode.Encode(c)

	r.AVP = []msg.RawAVP{rcode, host, realm}

	for _, a := range m.AVP {
		if e := a.Validate(0, 263, false, true, false); e == nil {
			r.AVP = append(r.AVP, a)
		}
	}
	return r
}
