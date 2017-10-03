package sock

import (
	"net"
	"time"

	"github.com/fkgi/diameter/msg"
)

// constant values
var (
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
}

// Dial make new Conn that use specified peernode and connection
func Dial(p *Peer, c net.Conn) (*Conn, error) {
	con := &Conn{
		peer:     p,
		notify:   make(chan stateEvent),
		state:    closed,
		con:      c,
		sndstack: make(map[uint32]chan msg.Answer)}
	con.run()

	ch := make(chan msg.Answer)
	cer := MakeCER(con)
	t := new(time.Timer)
	con.notify <- eventConnect{
		m: cer,
		c: ch}

	ack := <-ch
	t.Stop()

	if ack.Result() != msg.DiameterSuccess {
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
		sndstack: make(map[uint32]chan msg.Answer)}
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

// Send send Diameter request
func (c *Conn) Send(m msg.Request) msg.Answer {
	ch := make(chan msg.Answer)
	t := new(time.Timer)

	c.notify <- eventSndRequest{m: m, c: ch, t: t}

	a := <-ch
	t.Stop()

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
			} else if m.FlgR {
				c.notify <- eventRcvRequest{m}
			} else {
				c.notify <- eventRcvAnswer{m}
			}
		}
		c.notify <- eventPeerDisc{}
	}()
}

func (c *Conn) write(m msg.RawMsg) error {
	c.con.SetWriteDeadline(time.Now().Add(TransportTimeout))
	_, e := m.WriteTo(c.con)
	return e
}
