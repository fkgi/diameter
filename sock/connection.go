package sock

import (
	"net"
	"time"

	"github.com/fkgi/diameter/msg"
)

// constant values
var (
	MsgStackLen                           = 10000
	VendorID         msg.VendorID         = 41102
	ProductName      msg.ProductName      = "yatagarasu"
	FirmwareRevision msg.FirmwareRevision = 170819001
)

// Conn is state machine of Diameter
type Conn struct {
	local *Local
	peer  *Peer

	sysTimer  *time.Timer // system message timer
	wdCounter int         // watchdog expired counter

	notify chan stateEvent
	state
	con      net.Conn
	sndstack map[uint32]chan msg.Message
}

// Dial make new Conn that use specified peernode and connection
func (l *Local) Dial(p *Peer, c net.Conn) (*Conn, error) {
	con := &Conn{
		local:    l,
		peer:     p,
		notify:   make(chan stateEvent),
		state:    closed,
		con:      c,
		sndstack: make(map[uint32]chan msg.Message)}
	con.run()

	con.notify <- eventConnect{}

	return con, nil
}

// Accept new transport connection and return Conn
func (l *Local) Accept(p *Peer, c net.Conn) (*Conn, error) {
	con := &Conn{
		local:    l,
		peer:     p,
		notify:   make(chan stateEvent),
		state:    waitCER,
		con:      c,
		sndstack: make(map[uint32]chan msg.Message)}
	con.run()

	return con, nil
}

// State return status of connection
func (c *Conn) State() string {
	return c.state.String()
}

// Close stop state machine
func (c *Conn) Close(timeout time.Duration) {
	if c.state != open {
		return
	}

	c.notify <- eventStop{}
}

// LocalHost returns local host name
func (c *Conn) LocalHost() msg.DiameterIdentity {
	return c.local.Host
}

// LocalRealm returns local realm name
func (c *Conn) LocalRealm() msg.DiameterIdentity {
	return c.local.Realm
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
func (c *Conn) AuthApplication() map[msg.VendorID][]msg.ApplicationID {
	return c.peer.AuthApps
}

// Send send Diameter request
func (c *Conn) Send(r msg.Message) msg.Message {
	ch := make(chan msg.Message)
	r.HbHID = c.local.NextHbH()
	c.sndstack[r.HbHID] = ch

	c.notify <- eventSndMsg{r}

	t := time.AfterFunc(c.peer.SndTimeout,
		func() {
			nack := msg.Message{
				Ver:  r.Ver,
				FlgR: false, FlgP: r.FlgP, FlgE: true, FlgT: false,
				HbHID: r.HbHID, EtEID: r.EtEID,
				Code: r.Code, AppID: r.AppID}
			host := msg.OriginHost(c.local.Host)
			realm := msg.OriginRealm(c.local.Realm)
			state := msg.AuthSessionState(msg.StateNotMaintained)
			avps := []msg.Avp{
				msg.DiameterUnableToDeliver.Encode(),
				state.Encode(),
				host.Encode(),
				realm.Encode()}
			nack.Encode(avps)

			ch <- nack
		})

	a := <-ch
	t.Stop()
	delete(c.sndstack, r.HbHID)

	return a
}

func (c *Conn) run() {
	go func() {
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
	}()
	go func() {
		for {
			m := msg.Message{}
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
	}()
}

func (c *Conn) sendSysMsg(req, nak msg.Message) error {
	req.HbHID = c.local.NextHbH()
	nak.HbHID = req.HbHID
	req.EtEID = c.local.NextEtE()
	nak.EtEID = req.EtEID
	c.sndstack[req.HbHID] = nil //make(chan msg.Message)

	if e := c.write(req); e != nil {
		return e
	}

	c.sysTimer = time.AfterFunc(c.peer.SndTimeout, func() {
		switch nak.Code {
		case 257:
			c.notify <- eventRcvCEA{nak}
		case 282:
			c.notify <- eventRcvDPA{nak}
		}
	})
	return nil
}

func (c *Conn) write(m msg.Message) error {
	c.con.SetWriteDeadline(time.Now().Add(TransportTimeout))
	_, e := m.WriteTo(c.con)
	return e
}
