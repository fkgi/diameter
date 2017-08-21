package sock

import (
	"fmt"

	"github.com/fkgi/diameter/msg"
)

var (
	stateMap = map[int]string{
		shutdown: "shutdown",
		closed:   "closed",
		waitCER:  "waitCER",
		waitCEA:  "waitCEA",
		open:     "open",
		closing:  "closing",
	}
)

const (
	shutdown = iota
	closed   = iota
	waitCER  = iota
	waitCEA  = iota
	open     = iota
	closing  = iota
)

// NotAcceptableEvent is error
type NotAcceptableEvent struct {
	event stateEvent
	state int
}

func (e NotAcceptableEvent) Error() string {
	return "Event " + e.event.String() +
		" is not acceptable in state " + stateMap[e.state]
}

type stateEvent interface {
	exec(p *Conn) error
	String() string
}

// Connect
type eventConnect struct {
	m msg.Message
}

func (eventConnect) String() string {
	return "Connect"
}

func (v eventConnect) exec(c *Conn) error {
	if c.state != closed {
		return NotAcceptableEvent{
			event: v,
			state: c.state}
	}

	c.state = waitCEA
	c.setTransportDeadline()
	_, e := v.m.WriteTo(c.con)

	//notify(&CapabilityExchangeEvent{
	//	Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

// Watchdog
type eventWatchdog struct {
	m msg.Message
}

func (eventWatchdog) String() string {
	return "Watchdog"
}

func (v eventWatchdog) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{
			event: v,
			state: c.state}
	}

	c.setTransportDeadline()
	_, e := v.m.WriteTo(c.con)

	// notify(&WatchdogEvent{
	//	Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

// Stop
type eventStop struct {
	m msg.Message
}

func (eventStop) String() string {
	return "Stop"
}

func (v eventStop) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{
			event: v,
			state: c.state}
	}

	c.state = closing
	if c.wTimer != nil {
		c.wTimer.Stop()
	}
	c.setTransportDeadline()
	_, e := v.m.WriteTo(c.con)

	// notify(&PurgeEvent{
	// 	Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

// PeerDisc
type eventPeerDisc struct {
	Err error
}

func (eventPeerDisc) String() string {
	return "Peer-Disc"
}

func (v eventPeerDisc) exec(c *Conn) error {
	c.con.Close()
	c.state = closed

	// notify(&DisconnectEvent{
	// 	Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
	return nil
}

// RcvCER
type eventRcvCER struct {
	m msg.Message
}

func (eventRcvCER) String() string {
	return "Rcv-CER"
}

func (v eventRcvCER) exec(c *Conn) error {
	if c.state != waitCER {
		return NotAcceptableEvent{
			event: v,
			state: c.state}
	}

	p := &Peer{}
	avp, e := v.m.Decode()
	if e == nil {
		if t, ok := msg.GetOriginHost(avp); ok {
			p.Host = msg.DiameterIdentity(t)
		}
		if t, ok := msg.GetOriginRealm(avp); ok {
			p.Realm = msg.DiameterIdentity(t)
		}
		getProvidedAuthApp(p, avp)
	}

	//notify(&CapabilityExchangeEvent{
	//	Tx: false, Req: true, Local: p.local, Peer: peer, Err: e})
	if e != nil {
		return e
	}

	a, code := c.makeCEA(v.m, p)
	c.setTransportDeadline()
	_, e = a.WriteTo(c.con)

	if e == nil {
		if code != msg.DiameterSuccess {
			e = fmt.Errorf("close with error response %d", code)
			c.con.Close()
		} else {
			c.state = open
			c.resetWatchdog()
		}
	}
	//notify(&CapabilityExchangeEvent{
	//	Tx: true, Req: false, Local: p.local, Peer: p.peer})
	return e
}

func getProvidedAuthApp(p *Peer, avp msg.GroupedAVP) {
	p.AuthApps = map[msg.VendorID][]msg.ApplicationID{}

	for _, vid := range msg.GetSupportedVendorIDs(avp) {
		p.AuthApps[msg.VendorID(vid)] = make([]msg.ApplicationID, 0)
	}
	for _, vsa := range msg.GetVendorSpecificApplicationIDs(avp) {
		if _, ok := p.AuthApps[vsa.VendorID]; !ok {
			p.AuthApps[vsa.VendorID] = make([]msg.ApplicationID, 0)
		}
		p.AuthApps[vsa.VendorID] = append(p.AuthApps[vsa.VendorID], vsa.App)
	}
	p.AuthApps[0] = make([]msg.ApplicationID, 0)
	for _, aid := range msg.GetAuthApplicationIDs(avp) {
		p.AuthApps[0] = append(p.AuthApps[0], aid)
	}
}

// RcvCEA
type eventRcvCEA struct {
	m msg.Message
}

func (eventRcvCEA) String() string {
	return "Rcv-CEA"
}

func (v eventRcvCEA) exec(c *Conn) (e error) {
	if c.state != waitCEA {
		return NotAcceptableEvent{
			event: v,
			state: c.state}
	}

	var r msg.ResultCode
	if avp, e := v.m.Decode(); e == nil {
		if t, ok := msg.GetResultCode(avp); ok {
			r = t
		}
		getProvidedAuthApp(c.peer, avp)
	}
	if r == msg.DiameterSuccess {
		c.state = open
		c.resetWatchdog()
	} else {
		e = fmt.Errorf("CEA Nack received")
		c.con.Close()
	}
	//notify(&CapabilityExchangeEvent{
	//	Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
	return
}

type eventRcvDWR struct {
	m msg.Message
}

func (eventRcvDWR) String() string {
	return "Rcv-DWR"
}

func (v eventRcvDWR) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{
			event: v,
			state: c.state}
	}
	//notify(&WatchdogEvent{
	//	Tx: false, Req: true, Local: p.local, Peer: p.peer, Err: e})

	a, _ := c.makeDWA(v.m)
	c.setTransportDeadline()
	_, e := a.WriteTo(c.con)
	// notify(&WatchdogEvent{
	//	Tx: true, Req: false, Local: p.local, Peer: p.peer, Err: e})
	if e != nil {
		c.con.Close()
		c.state = shutdown
	} else {
		c.resetWatchdog()
	}
	return e
}

// RcvDWA
type eventRcvDWA struct {
	m msg.Message
}

func (eventRcvDWA) String() string {
	return "Rcv-DWA"
}

func (v eventRcvDWA) exec(c *Conn) (e error) {
	if c.state != open {
		return NotAcceptableEvent{
			event: v,
			state: c.state}
	}

	if ch, ok := c.sndstack[v.m.HbHID]; ok {
		delete(c.sndstack, v.m.HbHID)
		ch <- v.m
		c.resetWatchdog()
	} else {
		e = fmt.Errorf("unknown DWA recieved")
	}
	//notify(&WatchdogEvent{
	//	Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
	return
}

type eventRcvDPR struct {
	m msg.Message
}

func (eventRcvDPR) String() string {
	return "Rcv-DPR"
}

func (v eventRcvDPR) exec(c *Conn) (e error) {
	if c.state != open {
		return NotAcceptableEvent{
			event: v,
			state: c.state}
	}

	//notify(&PurgeEvent{
	//	Tx: false, Req: true, Local: p.local, Peer: p.peer, Err: e})

	c.state = closing
	a, _ := c.makeDPA(v.m)
	c.setTransportDeadline()
	_, e = a.WriteTo(c.con)
	//notify(&PurgeEvent{
	//	Tx: true, Req: false, Local: p.local, Peer: p.peer, Err: e})
	if e != nil {
		c.con.Close()
	}
	return
}

type eventRcvDPA struct {
	m msg.Message
}

func (eventRcvDPA) String() string {
	return "Rcv-DPA"
}

func (v eventRcvDPA) exec(c *Conn) (e error) {
	if c.state != closing {
		return NotAcceptableEvent{
			event: v,
			state: c.state}
	}

	if ch, ok := c.sndstack[v.m.HbHID]; ok {
		ch <- v.m
		// p.con.Close()
	} else {
		e = fmt.Errorf("unknown DPA recieved")
	}
	//notify(&PurgeEvent{
	//	Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
	return
}

type eventSndMsg struct {
	m msg.Message
}

func (eventSndMsg) String() string {
	return "Snd-MSG"
}

func (v eventSndMsg) exec(c *Conn) error {
	if c.state != open {
		return NotAcceptableEvent{
			event: v,
			state: c.state}
	}

	c.setTransportDeadline()
	_, e := v.m.WriteTo(c.con)
	//notify(&MessageEvent{
	//	Tx: true, Req: v.m.FlgR, Local: p.local, Peer: p.peer, Err: e})
	if e != nil {
		c.con.Close()
	}
	return e
}

type eventRcvMsg struct {
	m msg.Message
}

func (eventRcvMsg) String() string {
	return "Rcv-MSG"
}

func (v eventRcvMsg) exec(c *Conn) (e error) {
	if c.state != open {
		return NotAcceptableEvent{
			event: v,
			state: c.state}
	}

	if v.m.FlgR {
		c.rcvstack <- &v.m
		c.resetWatchdog()
	} else if ch, ok := c.sndstack[v.m.HbHID]; ok {
		ch <- &v.m
		c.resetWatchdog()
	} else {
		e = fmt.Errorf("unknown answer message received")
	}

	//notify(&MessageEvent{
	//	Tx: false, Req: v.m.FlgR, Local: p.local, Peer: p.peer, Err: e})
	return
}
