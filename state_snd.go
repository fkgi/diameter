package diameter

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

const busyTO = time.Millisecond * 100

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

	avps := []AVP{}
	avps = append(avps, SetOriginHost(Host))
	avps = append(avps, SetOriginRealm(Realm))

	h, _, _ := net.SplitHostPort(c.conn.LocalAddr().String())
	for _, h := range strings.Split(h, "/") {
		avps = append(avps, setHostIPAddress(net.ParseIP(h)))
	}

	avps = append(avps, SetVendorID(VendorID))
	avps = append(avps, setProductName(ProductName))
	if stateID != 0 {
		avps = append(avps, setOriginStateID(stateID))
	}
	if len(applications) == 0 {
		avps = append(avps, SetAuthAppID(0xffffffff))
	} else {
		vmap := make(map[uint32]interface{})
		for aid, app := range applications {
			if app.venID == 0 {
				avps = append(avps, SetAuthAppID(aid))
			} else if _, ok := vmap[app.venID]; ok {
				avps = append(avps, SetVendorSpecAppID(app.venID, aid))
			} else {
				avps = append(avps, setSupportedVendorID(app.venID))
				avps = append(avps, SetVendorSpecAppID(app.venID, aid))
				vmap[app.venID] = nil
			}
		}
	}
	// Inband-Security-Id
	// Acct-Application-Id
	avps = append(avps, setFirmwareRevision(FirmwareRev))

	buf := new(bytes.Buffer)
	for _, a := range overwriteAVP(avps, OverrideCER) {
		a.MarshalTo(buf)
	}

	cer := Message{
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 257, AppID: 0,
		HbHID: nextHbH(), EtEID: nextEtE(),
		AVPs:     buf.Bytes(),
		PeerName: c.Host, PeerRealm: c.Realm}

	c.sndQueue[cer.HbHID] = make(chan Message)
	c.wdTimer = time.AfterFunc(TransactionWait, func() {
		c.notify <- eventRcvCEA{cer.GenerateAnswerBy(UnableToDeliver)}
	})

	err := cer.MarshalTo(c.conn)
	if err != nil {
		err = TransportTxError{err: err}
		c.notify <- eventPeerDisc{reason: err}
	}

	if TraceTxMessage != nil {
		TraceTxMessage(c, cer, err)
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
		err := fmt.Errorf("watchdog is expired")
		c.notify <- eventPeerDisc{reason: err}
		return err
	}
	c.wdTimer.Stop()

	avps := []AVP{}
	avps = append(avps, SetOriginHost(Host))
	avps = append(avps, SetOriginRealm(Realm))
	if stateID != 0 {
		avps = append(avps, setOriginStateID(stateID))
	}

	buf := new(bytes.Buffer)
	for _, a := range overwriteAVP(avps, OverrideDWR) {
		a.MarshalTo(buf)
	}

	dwr := Message{
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 280, AppID: 0,
		HbHID: nextHbH(), EtEID: nextEtE(),
		AVPs:     buf.Bytes(),
		PeerName: c.Host, PeerRealm: c.Realm}

	c.sndQueue[dwr.HbHID] = make(chan Message)
	c.wdTimer = time.AfterFunc(TransactionWait, func() {
		c.notify <- eventRcvDWA{dwr.GenerateAnswerBy(UnableToDeliver)}
		c.notify <- eventWatchdog{}
	})

	c.conn.SetWriteDeadline(time.Now().Add(busyTO))
	err := dwr.MarshalTo(c.conn)
	c.conn.SetWriteDeadline(time.Time{})
	if err != nil {
		err = TransportTxError{err: err}
		c.notify <- eventPeerDisc{reason: err}
	}

	if TraceTxMessage != nil {
		TraceTxMessage(c, dwr, err)
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

func (v eventStop) String() string {
	return "Stop"
}

func (v eventStop) exec(c *Connection) error {
	if c.state != open && c.state != locked {
		err := notAcceptableEvent{e: v, s: c.state}
		c.notify <- eventPeerDisc{reason: err}
		return err
	}
	c.state = closing
	c.wdTimer.Stop()

	avps := []AVP{}
	avps = append(avps, SetOriginHost(Host))
	avps = append(avps, SetOriginRealm(Realm))
	avps = append(avps, setDisconnectCause(v.cause))

	buf := new(bytes.Buffer)
	for _, a := range overwriteAVP(avps, OverrideDPR) {
		a.MarshalTo(buf)
	}

	dpr := Message{
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 282, AppID: 0,
		HbHID: nextHbH(), EtEID: nextEtE(),
		AVPs:     buf.Bytes(),
		PeerName: c.Host, PeerRealm: c.Realm}

	c.sndQueue[dpr.HbHID] = make(chan Message)
	c.wdTimer = time.AfterFunc(TransactionWait, func() {
		c.notify <- eventRcvDPA{dpr.GenerateAnswerBy(UnableToDeliver)}
	})

	c.conn.SetWriteDeadline(time.Now().Add(busyTO))
	err := dpr.MarshalTo(c.conn)
	c.conn.SetWriteDeadline(time.Time{})
	if err != nil {
		err = TransportTxError{err: err}
		c.notify <- eventPeerDisc{reason: err}
	}

	if TraceTxMessage != nil {
		TraceTxMessage(c, dpr, err)
	}
	return err
}

// PeerDisc
type eventPeerDisc struct {
	reason error
}

func (v eventPeerDisc) String() string {
	if v.reason == nil {
		return "Peer-Disc(shutdown)"
	} else {
		return fmt.Sprintf("Peer-Disc(%s)", v.reason)
	}
}

func (v eventPeerDisc) exec(c *Connection) error {
	if c.state == closed {
		return v.reason
	}
	if c.state == closing {
		v.reason = nil
	}

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
	m  Message
	ch chan Message
	ts time.Time
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

	var err error
	if v.ch == nil {
		// Answer message handling
		if err = v.m.MarshalTo(c.conn); err != nil {
			err = TransportTxError{err: err}
		}
	} else if time.Since(v.ts) > TransactionWait {
		// timedout while waiting event queue
		v.ch <- v.m.GenerateAnswerBy(TooBusy)
		err = TransportTxError{
			err: errors.New("timeout while waiting event queue")}
	} else {
		// Request message handling
		c.sndQueue[v.m.HbHID] = v.ch
		c.conn.SetWriteDeadline(time.Now().Add(busyTO))
		err = v.m.MarshalTo(c.conn)
		c.conn.SetWriteDeadline(time.Time{})
		if err != nil {
			delete(c.sndQueue, v.m.HbHID)
			if e, ok := err.(*net.OpError); ok && e.Timeout() {
				v.ch <- v.m.GenerateAnswerBy(TooBusy)
			} else {
				v.ch <- v.m.GenerateAnswerBy(UnableToDeliver)
			}
			err = TransportTxError{err: err}
			//c.notify <- eventPeerDisc{reason: err}
		}
	}

	if TraceTxMessage != nil {
		TraceTxMessage(c, v.m, err)
	}
	return err
}
