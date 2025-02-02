package diameter

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"
)

type conState int

func (s conState) String() string {
	switch s {
	case closed:
		return "closed"
	case waitCER:
		return "waitCER"
	case waitCEA:
		return "waitCEA"
	case open:
		return "open"
	case locked:
		return "locked"
	case closing:
		return "closing"
	case shutdown:
		return "shutdown"
	}
	return "<nil>"
}

const (
	closed conState = iota
	waitCER
	waitCEA
	open
	locked
	closing
	shutdown
)

type stateEvent interface {
	exec(*Connection) error
	fmt.Stringer
}

// Init
type eventInit struct{}

func (eventInit) String() string {
	return "Initialize"
}

func (v eventInit) exec(c *Connection) error {
	return notAcceptableEvent{e: v, s: c.state}
}

// Connect
type eventConnect struct{}

func (eventConnect) String() string {
	return "Connect"
}

func (v eventConnect) exec(c *Connection) error {
	if c.state != closed {
		return notAcceptableEvent{e: v, s: c.state}
	}
	c.state = waitCEA

	buf := new(bytes.Buffer)
	SetOriginHost(Host).MarshalTo(buf)
	SetOriginRealm(Realm).MarshalTo(buf)

	if len(OverwriteAddr) != 0 {
		for _, h := range OverwriteAddr {
			setHostIPAddress(h).MarshalTo(buf)
		}
	} else {
		h, _, _ := net.SplitHostPort(c.conn.LocalAddr().String())
		for _, h := range strings.Split(h, "/") {
			setHostIPAddress(net.ParseIP(h)).MarshalTo(buf)
		}
	}

	SetVendorID(VendorID).MarshalTo(buf)
	setProductName(ProductName).MarshalTo(buf)
	if stateID != 0 {
		setOriginStateID(stateID).MarshalTo(buf)
	}
	if len(applications) == 0 {
		SetAuthAppID(0xffffffff).MarshalTo(buf)
	} else {
		vmap := make(map[uint32]interface{})
		for aid, app := range applications {
			if app.venID == 0 {
				SetAuthAppID(aid).MarshalTo(buf)
			} else if _, ok := vmap[app.venID]; ok {
				SetVendorSpecAppID(app.venID, aid).MarshalTo(buf)
			} else {
				setSupportedVendorID(app.venID).MarshalTo(buf)
				SetVendorSpecAppID(app.venID, aid).MarshalTo(buf)
				vmap[app.venID] = nil
			}
		}
	}
	// Inband-Security-Id
	// Acct-Application-Id
	setFirmwareRevision(FirmwareRev).MarshalTo(buf)

	cer := Message{
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 257, AppID: 0,
		HbHID: nextHbH(), EtEID: nextEtE(),
		AVPs:     buf.Bytes(),
		PeerName: c.Host, PeerRealm: c.Realm}

	c.sndQueue[cer.HbHID] = make(chan Message)
	c.wdTimer = time.AfterFunc(WDInterval, func() {
		c.notify <- eventRcvCEA{cer.GenerateAnswerBy(UnableToDeliver)}
	})

	err := cer.MarshalTo(c.conn)
	if err != nil {
		c.conn.Close()
		err = TransportTxError{err: err}
	}

	if TraceMessage != nil {
		TraceMessage(cer, Tx, err)
	}
	return err
}

// Watchdog
type eventWatchdog struct{}

func (eventWatchdog) String() string {
	return "Watchdog"
}

func (v eventWatchdog) exec(c *Connection) error {
	if c.state != open && c.state != locked {
		return notAcceptableEvent{e: v, s: c.state}
	}

	c.wdCount++
	if c.wdCount > WDMaxSend {
		c.conn.Close()
		return fmt.Errorf("watchdog is expired")
	}
	c.wdTimer.Stop()

	buf := new(bytes.Buffer)
	SetOriginHost(Host).MarshalTo(buf)
	SetOriginRealm(Realm).MarshalTo(buf)
	if stateID != 0 {
		setOriginStateID(stateID).MarshalTo(buf)
	}

	dwr := Message{
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 280, AppID: 0,
		HbHID: nextHbH(), EtEID: nextEtE(),
		AVPs:     buf.Bytes(),
		PeerName: c.Host, PeerRealm: c.Realm}

	c.sndQueue[dwr.HbHID] = make(chan Message)
	c.wdTimer = time.AfterFunc(WDInterval, func() {
		c.notify <- eventRcvDWA{dwr.GenerateAnswerBy(UnableToDeliver)}
		c.notify <- eventWatchdog{}
	})

	err := dwr.MarshalTo(c.conn)
	if err != nil {
		c.conn.Close()
		err = TransportTxError{err: err}
	}

	if TraceMessage != nil {
		TraceMessage(dwr, Tx, err)
	}
	return err
}

// Lock
type eventLock struct{}

func (eventLock) String() string {
	return "Lock"
}

func (v eventLock) exec(c *Connection) error {
	if c.state != open {
		return notAcceptableEvent{e: v, s: c.state}
	}
	c.state = locked
	return nil
}

// Stop
type eventStop struct {
	cause Enumerated
}

func (eventStop) String() string {
	return "Stop"
}

func (v eventStop) exec(c *Connection) error {
	if c.state != open && c.state != locked {
		c.conn.Close()
		return notAcceptableEvent{e: v, s: c.state}
	}
	c.state = closing
	c.wdTimer.Stop()

	buf := new(bytes.Buffer)
	SetOriginHost(Host).MarshalTo(buf)
	SetOriginRealm(Realm).MarshalTo(buf)
	setDisconnectCause(v.cause).MarshalTo(buf)

	dpr := Message{
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 282, AppID: 0,
		HbHID: nextHbH(), EtEID: nextEtE(),
		AVPs:     buf.Bytes(),
		PeerName: c.Host, PeerRealm: c.Realm}

	c.sndQueue[dpr.HbHID] = make(chan Message)
	c.wdTimer = time.AfterFunc(WDInterval, func() {
		c.notify <- eventRcvDPA{dpr.GenerateAnswerBy(UnableToDeliver)}
	})

	err := dpr.MarshalTo(c.conn)
	if err != nil {
		c.conn.Close()
		err = TransportTxError{err: err}
	}

	if TraceMessage != nil {
		TraceMessage(dpr, Tx, err)
	}
	return err
}

// PeerDisc
type eventPeerDisc struct {
	reason error
}

func (eventPeerDisc) String() string {
	return "Peer-Disc"
}

func (v eventPeerDisc) exec(c *Connection) error {
	c.conn.Close()
	c.state = closed

	for _, ch := range c.sndQueue {
		close(ch)
	}
	close(c.rcvQueue)

	return v.reason
}

// Snd MSG
type eventSndMsg struct {
	m Message
}

func (eventSndMsg) String() string {
	return "Snd-MSG"
}

func (v eventSndMsg) exec(c *Connection) error {
	if c.state != open && c.state != locked {
		return notAcceptableEvent{e: v, s: c.state}
	}

	v.m.PeerName = c.Host
	v.m.PeerRealm = c.Realm

	err := v.m.MarshalTo(c.conn)
	if err != nil {
		c.conn.Close()
		err = TransportTxError{err: err}
	}

	if TraceMessage != nil {
		TraceMessage(v.m, Tx, err)
	}
	return err
}
