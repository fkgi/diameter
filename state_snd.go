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
	exec() error
	fmt.Stringer
}

// Init
type eventInit struct{}

func (eventInit) String() string {
	return "Initialize"
}

func (v eventInit) exec() error {
	return notAcceptableEvent{e: v, s: state}
}

// Connect
type eventConnect struct{}

func (eventConnect) String() string {
	return "Connect"
}

func (v eventConnect) exec() error {
	if state != closed {
		return notAcceptableEvent{e: v, s: state}
	}
	state = waitCEA

	buf := new(bytes.Buffer)
	SetOriginHost(Local.Host).MarshalTo(buf)
	SetOriginRealm(Local.Realm).MarshalTo(buf)

	if len(OverwriteAddr) != 0 {
		for _, h := range OverwriteAddr {
			setHostIPAddress(h).MarshalTo(buf)
		}
	} else {
		h, _, _ := net.SplitHostPort(conn.LocalAddr().String())
		for _, h := range strings.Split(h, "/") {
			setHostIPAddress(net.ParseIP(h)).MarshalTo(buf)
		}
	}

	SetVendorID(VendorID).MarshalTo(buf)
	setProductName(ProductName).MarshalTo(buf)
	if Local.state != 0 {
		setOriginStateID(Local.state).MarshalTo(buf)
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
		AVPs: buf.Bytes()}

	TxReq++
	sndQueue[cer.HbHID] = make(chan Message)

	wdTimer = time.AfterFunc(WDInterval, func() {
		notify <- eventRcvCEA{cer.generateAnswerBy(UnableToDeliver)}
	})

	err := cer.MarshalTo(conn)
	if err != nil {
		conn.Close()
	}

	TraceMessage(cer, Tx, err)
	return err
}

// Watchdog
type eventWatchdog struct{}

func (eventWatchdog) String() string {
	return "Watchdog"
}

func (v eventWatchdog) exec() error {
	if state != open && state != locked {
		return notAcceptableEvent{e: v, s: state}
	}

	wdCount++
	if wdCount > WDMaxSend {
		conn.Close()
		return fmt.Errorf("watchdog is expired")
	}
	wdTimer.Stop()

	buf := new(bytes.Buffer)
	SetOriginHost(Local.Host).MarshalTo(buf)
	SetOriginRealm(Local.Realm).MarshalTo(buf)
	setOriginStateID(Local.state).MarshalTo(buf)

	dwr := Message{
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 280, AppID: 0,
		HbHID: nextHbH(), EtEID: nextEtE(),
		AVPs: buf.Bytes()}

	TxReq++
	sndQueue[dwr.HbHID] = make(chan Message)

	wdTimer = time.AfterFunc(WDInterval, func() {
		notify <- eventRcvDWA{dwr.generateAnswerBy(UnableToDeliver)}
		notify <- eventWatchdog{}
	})

	err := dwr.MarshalTo(conn)
	if err != nil {
		conn.Close()
	}

	TraceMessage(dwr, Tx, err)
	return err
}

// Lock
type eventLock struct{}

func (eventLock) String() string {
	return "Lock"
}

func (v eventLock) exec() error {
	if state != open {
		return notAcceptableEvent{e: v, s: state}
	}
	state = locked
	return nil
}

// Stop
type eventStop struct {
	cause Enumerated
}

func (eventStop) String() string {
	return "Stop"
}

func (v eventStop) exec() error {
	if state != open && state != locked {
		conn.Close()
		return notAcceptableEvent{e: v, s: state}
	}
	state = closing
	wdTimer.Stop()

	buf := new(bytes.Buffer)
	SetOriginHost(Local.Host).MarshalTo(buf)
	SetOriginRealm(Local.Realm).MarshalTo(buf)
	setDisconnectCause(v.cause).MarshalTo(buf)

	dpr := Message{
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 282, AppID: 0,
		HbHID: nextHbH(), EtEID: nextEtE(),
		AVPs: buf.Bytes()}

	TxReq++
	sndQueue[dpr.HbHID] = make(chan Message)

	wdTimer = time.AfterFunc(WDInterval, func() {
		notify <- eventRcvDPA{dpr.generateAnswerBy(UnableToDeliver)}
	})

	err := dpr.MarshalTo(conn)
	if err != nil {
		conn.Close()
	}

	TraceMessage(dpr, Tx, err)
	return err
}

// PeerDisc
type eventPeerDisc struct {
	reason error
}

func (eventPeerDisc) String() string {
	return "Peer-Disc"
}

func (v eventPeerDisc) exec() error {
	conn.Close()
	state = closed

	for _, ch := range sndQueue {
		close(ch)
	}
	close(rcvQueue)

	return v.reason
}

// Snd MSG
type eventSndMsg struct {
	m Message
}

func (eventSndMsg) String() string {
	return "Snd-MSG"
}

func (v eventSndMsg) exec() error {
	if state != open && state != locked {
		return notAcceptableEvent{e: v, s: state}
	}

	TxReq++
	err := v.m.MarshalTo(conn)
	if err != nil {
		conn.Close()
	}

	TraceMessage(v.m, Tx, err)
	return err
}
