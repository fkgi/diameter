package sock

import (
	"net"
	"time"

	"github.com/fkgi/diameter/msg"
)

// constant values
var (
	ReadBuffer                            = 65535
	VendorID         msg.VendorID         = 41102
	ProductName      msg.ProductName      = "yatagarasu"
	FirmwareRevision msg.FirmwareRevision = 170819001
)

// Conn is state machine of Diameter
type Conn struct {
	peer *Peer

	sysTimer  *time.Timer // system message timer
	wdCounter int         // watchdog expired counter

	notify chan stateEvent
	state
	con      net.Conn
	sndstack map[uint32]chan msg.Answer
	rcvstack chan msg.RawMsg
}

// Dial make new Conn that use specified peernode and connection
func Dial(p *Peer, c net.Conn) (*Conn, error) {
	con := &Conn{
		peer:     p,
		notify:   make(chan stateEvent),
		state:    closed,
		con:      c,
		sndstack: make(map[uint32]chan msg.Answer),
		rcvstack: make(chan msg.RawMsg, ReadBuffer)}
	go socketHandler(con)
	go eventHandler(con)
	go messageHandler(con)

	cer := MakeCER(con)
	req := cer.ToRaw()
	req.HbHID = NextHbH()
	req.EtEID = NextEtE()

	ch := make(chan msg.Answer)
	con.sndstack[req.HbHID] = ch
	con.notify <- eventConnect{m: req}

	t := time.AfterFunc(p.SndTimeout, func() {
		m := cer.Failed(
			msg.DiameterTooBusy,
			"no response from peer node").ToRaw()
		m.HbHID = req.HbHID
		m.EtEID = req.EtEID
		con.notify <- eventRcvCEA{m}
	})

	ack := <-ch
	t.Stop()

	if ack == nil || ack.Result() != msg.DiameterSuccess {
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
		sndstack: make(map[uint32]chan msg.Answer),
		rcvstack: make(chan msg.RawMsg, ReadBuffer)}
	go socketHandler(con)
	go eventHandler(con)
	go messageHandler(con)

	return con, nil
}

// State return status of connection
func (c *Conn) State() string {
	return c.state.String()
}

// Send send Diameter request
func (c *Conn) Send(m msg.Request) msg.Answer {
	req := m.ToRaw()
	req.HbHID = NextHbH()
	req.EtEID = NextEtE()

	ch := make(chan msg.Answer)
	c.sndstack[req.HbHID] = ch
	c.notify <- eventSndMsg{m: req}

	t := time.AfterFunc(c.peer.SndTimeout, func() {
		m := m.Failed(
			msg.DiameterTooBusy,
			"no response from peer node").ToRaw()
		m.HbHID = req.HbHID
		m.EtEID = req.EtEID
		c.notify <- eventRcvMsg{m}
	})

	a := <-ch
	t.Stop()
	if a == nil {
		a = m.Failed(
			msg.DiameterUnableToDeliver,
			"failed to send")
	}

	return a
}

func (c *Conn) watchdog() {
	dwr := MakeDWR(c)
	req := dwr.ToRaw()
	req.HbHID = NextHbH()
	req.EtEID = NextEtE()

	ch := make(chan msg.Answer)
	c.sndstack[req.HbHID] = ch
	c.notify <- eventWatchdog{m: req}

	t := time.AfterFunc(c.peer.SndTimeout, func() {
		m := dwr.Failed(
			msg.DiameterTooBusy,
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
	req.HbHID = NextHbH()
	req.EtEID = NextEtE()

	ch := make(chan msg.Answer)
	c.sndstack[req.HbHID] = ch
	c.notify <- eventStop{m: req}

	t := time.AfterFunc(c.peer.SndTimeout, func() {
		m := dpr.Failed(
			msg.DiameterTooBusy,
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
func (c *Conn) AuthApplication() map[msg.VendorID][]msg.AuthApplicationID {
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
		if m.FlgR {
			requestHandler(c, m)
		} else {
			answerHandler(c, m)
		}
	}
}

func requestHandler(c *Conn, m msg.RawMsg) {
	if app, ok := supportedApps[msg.AuthApplicationID(m.AppID)]; !ok {
		c.notify <- eventSndMsg{
			makeErrorMsg(m, msg.DiameterApplicationUnsupported)}
	} else if req, ok := app.req[m.Code]; !ok {
		c.notify <- eventSndMsg{
			makeErrorMsg(m, msg.DiameterCommandUnspported)}
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

func answerHandler(c *Conn, m msg.RawMsg) {
	ch, ok := c.sndstack[m.HbHID]
	if !ok {
		return
	}
	delete(c.sndstack, m.HbHID)

	if app, ok := supportedApps[msg.AuthApplicationID(m.AppID)]; !ok {
		// ToDo
		// invalid message handling
	} else if ans, ok := app.ans[m.Code]; !ok {
		// ToDo
		// invalid message handling
	} else if ack, e := ans.FromRaw(m); e != nil {
		// ToDo
		// invalid message handling
	} else {
		ch <- ack
	}
}

func makeErrorMsg(m msg.RawMsg, c msg.ResultCode) msg.RawMsg {
	r := msg.RawMsg{}
	r.Ver = m.Ver
	r.FlgP = m.FlgP
	r.Code = m.Code
	r.AppID = m.AppID
	r.HbHID = m.HbHID
	r.EtEID = m.EtEID

	host := msg.OriginHost(Host)
	realm := msg.OriginRealm(Realm)
	r.AVP = []msg.RawAVP{
		c.ToRaw(),
		host.ToRaw(),
		realm.ToRaw()}

	for _, a := range m.AVP {
		if e := a.Validate(0, 263, false, true, false); e == nil {
			r.AVP = append(r.AVP, a)
		}
	}
	return r
}
