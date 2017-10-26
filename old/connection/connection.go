package connection

import (
	"fmt"
	"net"
	"time"

	"github.com/fkgi/diameter/msg"
)

// constant values
var (
	MsgStackLen                           = 10000
	VendorID         msg.VendorID         = 41102
	ProductName      msg.ProductName      = "yagtagarasu"
	FirmwareRevision msg.FirmwareRevision = 170407001
)

// Connection is state machine of Diameter
type Connection struct {
	local *LocalNode
	peer  *PeerNode

	wT *time.Timer // watchdog timer
	wE int         // watchdog expired counter

	notify   chan stateEvent
	state    int
	con      net.Conn
	rcvstack chan *msg.Message
	sndstack map[uint32]chan *msg.Message

	openNtfy chan bool

	cachemsg msg.Message
}

// Close stop state machine
func (c *Connection) Close(cause msg.Enumerated) {
	if c.state != open {
		return
	}

	ch := make(chan *msg.Message)
	r := c.makeDPR(cause)
	r.HbHID = c.local.NextHbH()
	c.sndstack[r.HbHID] = ch

	c.notify <- eventStop{r}

	t := time.AfterFunc(c.peer.Tp, func() {
		ch <- nil
		notify(&PurgeEvent{
			Tx: false, Req: false, Local: c.local, Peer: c.peer,
			Err: fmt.Errorf("no answer")})
	})
	ap := <-ch
	t.Stop()
	delete(c.sndstack, r.HbHID)

	if ap != nil {
		c.notify <- eventPeerDisc{nil}
	}
}

// LocalHost returns local host name
func (c *Connection) LocalHost() msg.DiameterIdentity {
	return c.local.Host
}

// LocalRealm returns local realm name
func (c *Connection) LocalRealm() msg.DiameterIdentity {
	return c.local.Realm
}

// PeerHost returns peer host name
func (c *Connection) PeerHost() msg.DiameterIdentity {
	return c.peer.Host
}

// PeerRealm returns peer realm name
func (c *Connection) PeerRealm() msg.DiameterIdentity {
	return c.peer.Realm
}

// Connection returns transport connection of state machine
func (c *Connection) Connection() net.Conn {
	return c.con
}

// Properties returns properties of this state machine
func (c *Connection) Properties() Properties {
	return c.peer.Properties
}

// Send send Diameter request message
func (c *Connection) Send(r msg.Message) (a msg.Message, e error) {
	r.EtEID = c.local.NextEtE()

	for i := 0; i <= c.peer.Cp; i++ {
		a = c.Transmit(r)

		if avp, e := a.Decode(); e == nil {
			res, _ := msg.GetResultCode(avp)
			if res != msg.DiameterUnableToDeliver {
				break
			}
			ori, _ := msg.GetOriginHost(avp)
			if ori != msg.OriginHost(c.local.Host) {
				break
			}
		}
		if i >= c.peer.Cp {
			break
		}
		r.FlgT = true
	}
	return
}

// Transmit send Diameter request message
func (c *Connection) Transmit(r msg.Message) (a msg.Message) {
	r.HbHID = c.local.NextHbH()

	ch := make(chan *msg.Message)
	c.sndstack[r.HbHID] = ch
	c.notify <- eventSndMsg{r}

	t := time.AfterFunc(c.peer.Tp, func() {
		a := c.makeUnableToDeliver(r)
		ch <- &a
	})

	a = *<-ch
	t.Stop()
	delete(c.sndstack, r.HbHID)

	return
}

// Recieve recieve Diameter request message
func (c *Connection) Recieve() (r msg.Message, ch chan *msg.Message, e error) {
	rp := <-c.rcvstack
	if rp == nil {
		e = fmt.Errorf("closed")
		c.rcvstack <- nil
		return
	}
	r = *rp
	ch = make(chan *msg.Message)

	go func() {
		if a := <-ch; a == nil {
			e = fmt.Errorf("request is discarded")
			notify(&MessageEvent{
				Tx: true, Req: false, Local: c.local, Peer: c.peer, Err: e})
		} else {
			a.HbHID = r.HbHID
			a.EtEID = r.EtEID
			c.notify <- eventSndMsg{*a}
		}
	}()
	return
}

// WaitOpen wait for open state of this connection
func (c *Connection) WaitOpen() bool {
	return <-c.openNtfy
}

const (
	shutdown = iota
	closed   = iota
	waitCER  = iota
	waitCEA  = iota
	open     = iota
	closing  = iota
)

// State returns status of state machine
func (c *Connection) State() string {
	switch c.state {
	case closed:
		return "Closed"
	case waitCER:
		return "WaitCER"
	case waitCEA:
		return "WaitCEA"
	case open:
		return "Open"
	case closing:
		return "Closing"
	}
	return "Shutdown"
}

func (c *Connection) run() {
	go func() {
		for {
			m := msg.Message{}
			c.con.SetReadDeadline(time.Time{})
			if _, e := m.ReadFrom(c.con); e != nil {
				c.notify <- eventPeerDisc{e}
				break
			}

			c.wE = 0
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
		if c.wT != nil {
			c.wT.Stop()
		}
	}()

	old := "Shutdown"
	notify(&StateUpdate{
		OldState: old, NewState: c.State(), Event: "Start",
		Local: c.local, Peer: c.peer, Err: nil})

	for c.state != shutdown {
		event := <-c.notify
		old = c.State()
		e := event.exec(c)

		notify(&StateUpdate{
			OldState: old, NewState: c.State(), Event: event.name(),
			Local: c.local, Peer: c.peer, Err: e})
	}
}

func (c *Connection) resetWatchdog() {
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
