package sock

import (
	"net"
	"strings"
	"time"

	"github.com/fkgi/diameter/msg"
)

// constant values
var (
	MsgStackLen                           = 10000
	VendorID         msg.VendorID         = 41102
	ProductName      msg.ProductName      = "yagtagarasu"
	FirmwareRevision msg.FirmwareRevision = 170819001
)

// Conn is state machine of Diameter
type Conn struct {
	local *Local
	peer  *Peer

	wTimer   *time.Timer // watchdog timer
	wCounter int         // watchdog expired counter

	notify   chan stateEvent
	state    int
	con      net.Conn
	sndstack map[uint32]chan msg.Message

	cachemsg msg.Message
}

func (c *Conn) setTransportDeadline() {
	c.con.SetWriteDeadline(time.Now().Add(TransportTimeout))
}

// Dial make new Conn that use specified peernode and connection
func Dial(l *Local, p *Peer) (*Conn, error) {
	dialer := net.Dialer{
		Timeout:   TransportTimeout,
		LocalAddr: l.Addr}
	c, e := dialer.Dial(p.Addr.Network(), p.Addr.String())
	if e != nil {
		return nil, e
	}

	con := &Conn{
		local:    l,
		peer:     p,
		notify:   make(chan stateEvent),
		state:    closed,
		con:      c,
		sndstack: make(map[uint32]chan msg.Message)}
	con.run()

	r := msg.CER{
		OriginHost:  msg.OriginHost(l.Host),
		OriginRealm: msg.OriginRealm(l.Realm),
		//HostIPAddress: make([]msg.HostIPAddress, 0),
		VendorID:    VendorID,
		ProductName: ProductName,
		// *OriginStateID:
		SupportedVendorID:           make([]msg.SupportedVendorID, 0),
		ApplicationID:               make([]msg.ApplicationID, 0),
		VendorSpecificApplicationID: make([]msg.VendorSpecificApplicationID, 0),
		FirmwareRevision:            &FirmwareRevision}

	switch l.Addr.Network() {
	case "tcp":
		s := l.Addr.String()
		s = s[:strings.LastIndex(s, ":")]
		r.HostIPAddress = []msg.HostIPAddress{msg.HostIPAddress(net.ParseIP(s))}
	case "sctp":
		s := l.Addr.String()
		s = s[:strings.LastIndex(s, ":")]
		r.HostIPAddress = []msg.HostIPAddress{}
		for _, i := range strings.Split(s, "/") {
			r.HostIPAddress = append(r.HostIPAddress, msg.HostIPAddress(net.ParseIP(i)))
		}
	}
	if l.StateID != 0 {
		r.OriginStateID = &l.StateID
	}
	for v, a := range l.AuthApps {
		if v != 0 {
			r.SupportedVendorID = append(
				r.SupportedVendorID, msg.SupportedVendorID(v))
			for _, i := range a {
				r.VendorSpecificApplicationID = append(
					r.VendorSpecificApplicationID,
					msg.VendorSpecificApplicationID{
						VendorID: v,
						App:      i})
			}
		} else {
			for _, i := range a {
				r.ApplicationID = append(r.ApplicationID, i)
			}
		}
	}
	m := r.Encode()
	m.HbHID = l.NextHbH()
	m.EtEID = l.NextEtE()

	con.notify <- eventConnect{m}

	return con, nil
}

// Accept new transport connection and return Conn
func Accept(l *Local, p *Peer) (*Conn, error) {
	lnr, e := net.Listen(l.Addr.Network(), l.Addr.String())
	if e != nil {
		return nil, e
	}
	c, e := lnr.Accept()
	lnr.Close()
	if e != nil {
		return nil, e
	}

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

// Close stop state machine
func (c *Conn) Close(cause msg.Enumerated, timeout int) {
	if c.state != open {
		return
	}

	r := msg.DPR{
		OriginHost:      msg.OriginHost(c.local.Host),
		OriginRealm:     msg.OriginRealm(c.local.Realm),
		DisconnectCause: msg.DisconnectCause(cause)}
	m := r.Encode()
	m.HbHID = c.local.NextHbH()
	m.EtEID = c.local.NextEtE()

	ch := make(chan msg.Message)
	c.sndstack[m.HbHID] = ch

	c.notify <- eventStop{m}

	t := time.AfterFunc(
		time.Second*time.Duration(timeout),
		func() {
			dwa := msg.DWA{
				ResultCode:  msg.DiameterSuccess,
				OriginHost:  msg.OriginHost(c.local.Host),
				OriginRealm: msg.OriginRealm(c.local.Realm),
			}
			if c.local.StateID != 0 {
				dwa.OriginStateID = &c.local.StateID
			}
			m := dwa.Encode()
			m.HbHID = c.local.NextHbH()
			m.EtEID = c.local.NextEtE()
			m.FlgE = true
			ch <- m
			//notify(&PurgeEvent{
			//	Tx: false, Req: false, Local: c.local, Peer: c.peer,
			//	Err: fmt.Errorf("no answer")})
		})
	a := <-ch
	t.Stop()
	delete(c.sndstack, m.HbHID)

	if avp, e := a.Decode(); e == nil {
		if result, ok := msg.GetResultCode(avp); !ok {
			if result == msg.DiameterSuccess {
				return
			}
		}
	}

	c.notify <- eventPeerDisc{}
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
func (c *Conn) Send(r msg.Message, d time.Duration) msg.Message {
	ch := make(chan msg.Message)
	c.sndstack[r.HbHID] = ch
	c.notify <- eventSndMsg{r}

	t := time.AfterFunc(d, func() {
		nack := msg.Message{
			Ver:  r.Ver,
			FlgR: false, FlgP: r.FlgP, FlgE: true, FlgT: false,
			HbHID: r.HbHID, EtEID: r.EtEID,
			Code: r.Code, AppID: r.AppID}
		host := msg.OriginHost(c.local.Host)
		realm := msg.OriginRealm(c.local.Realm)
		avps := []msg.Avp{
			msg.DiameterUnableToDeliver.Encode(),
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
		//notify(&StateUpdate{
		//	OldState: c.State(), NewState: c.State(), Event: "Start",
		//	Local: c.local, Peer: c.peer, Err: nil})

		for {
			event := <-c.notify
			//old = c.State()
			e := event.exec(c)

			//notify(&StateUpdate{
			//	OldState: old, NewState: c.State(), Event: event.name(),
			//	Local: c.local, Peer: c.peer, Err: e})

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

			c.wCounter = 0
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

/*
func (c *Conn) resetWatchdog() {
	f := func() {
		c.wT = nil

		c.wE++
		if c.wE > c.peer.Ew {
			c.Close(msg.Enumerated(0))
		} else {
			ch := make(chan *msg.Message)
			r := c.makeDWR()
			r.HbHID = c.local.NextHbH()
			c.sndstack[r.HbHID] = ch
			c.notify <- eventWatchdog{r}

			t := time.AfterFunc(c.peer.Tp, func() {
				ch <- nil
				notify(&WatchdogEvent{
					Tx: false, Req: false, Local: c.local, Peer: c.peer,
					Err: fmt.Errorf("no answer")})
			})
			ap := <-ch
			t.Stop()
			delete(c.sndstack, r.HbHID)
			if ap != nil {
				c.wE = 0
			}
		}
	}

	if c.wT != nil {
		c.wT.Reset(c.peer.Tw)
	} else {
		c.wT = time.AfterFunc(c.peer.Tw, f)
	}
}
*/
