package connection

import (
	"fmt"
	"time"

	"github.com/fkgi/diameter/msg"
)

type stateEvent interface {
	exec(p *Connection) error
	name() string
}

type eventConnect struct{}

func (eventConnect) name() string {
	return "Connect"
}

func (eventConnect) exec(p *Connection) error {
	if p.state != closed {
		return fmt.Errorf("not acceptable event")
	}
	r := p.makeCER(p.con)
	r.HbHID = p.local.NextHbH()
	p.state = waitCEA

	p.con.SetWriteDeadline(time.Now().Add(p.peer.Ts))
	_, e := r.WriteTo(p.con)

	if Notificator != nil {
		Notificator(&CapabilityExchangeEvent{
			Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
	}
	if e != nil {
		p.con.Close()
	}
	return nil
}

type eventAccept struct{}

func (e eventAccept) name() string {
	return "Accept"
}

func (eventAccept) exec(p *Connection) error {
	if p.state != closed {
		return fmt.Errorf("not acceptable event")
	}
	p.state = waitCER
	return nil
}

type eventWatchdog struct {
	m msg.Message
}

func (eventWatchdog) name() string {
	return "Watchdog"
}

func (v eventWatchdog) exec(p *Connection) (e error) {
	if p.state != open {
		return fmt.Errorf("not acceptable event")
	}

	p.con.SetWriteDeadline(time.Now().Add(p.peer.Ts))
	_, e = v.m.WriteTo(p.con)

	if Notificator != nil {
		Notificator(&WatchdogEvent{
			Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
	}
	if e != nil {
		p.con.Close()
	}
	return
}

type eventStop struct {
	m msg.Message
}

func (eventStop) name() string {
	return "Stop"
}

func (v eventStop) exec(p *Connection) (e error) {
	p.state = closing

	p.con.SetWriteDeadline(time.Now().Add(p.peer.Ts))
	_, e = v.m.WriteTo(p.con)
	if Notificator != nil {
		Notificator(&PurgeEvent{
			Tx: true, Req: true, Local: p.local, Peer: p.peer, Err: e})
	}
	if e != nil {
		p.con.Close()
	}
	return
}

type eventPeerDisc struct {
	Err error
}

func (v eventPeerDisc) name() string {
	return "Peer-Disc"
}

func (v eventPeerDisc) exec(p *Connection) (e error) {
	p.con.Close()
	p.state = shutdown
	return
}

type eventRcvCER struct {
	m msg.Message
}

func (eventRcvCER) name() string {
	return "Rcv-CER"
}

func (v eventRcvCER) exec(p *Connection) (e error) {
	if p.state != waitCER {
		e = fmt.Errorf("not acceptable message")
	}
	if Notificator != nil {
		Notificator(&CapabilityExchangeEvent{
			Tx: false, Req: true, Local: p.local, Peer: p.peer, Err: e})
	}
	if e != nil {
		return
	}

	if avp, e := v.m.Decode(); e == nil {
		for _, a := range avp {
			if a.Code == uint32(264) && a.VenID == uint32(0) {
				a.Decode(&(p.peer.Host))
			}
			if a.Code == uint32(296) && a.VenID == uint32(0) {
				a.Decode(&(p.peer.Realm))
			}
		}
	}
	a, code := p.makeCEA(v.m, p.con)
	p.con.SetWriteDeadline(time.Now().Add(p.peer.Ts))
	_, e = a.WriteTo(p.con)

	if e == nil {
		if code != 2001 {
			e = fmt.Errorf("close with error response %d", code)
			p.con.Close()
		} else {
			p.state = open
			p.resetWatchdog()
		}
	}
	if Notificator != nil {
		Notificator(&CapabilityExchangeEvent{
			Tx: true, Req: false, Local: p.local, Peer: p.peer})
	}
	return
}

type eventRcvCEA struct {
	m msg.Message
}

func (eventRcvCEA) name() string {
	return "Rcv-CEA"
}

func (v eventRcvCEA) exec(p *Connection) (e error) {
	if p.state != waitCEA {
		e = fmt.Errorf("not acceptable message")
	} else {
		c := new(uint32)
		if avp, e := v.m.Decode(); e == nil {
			for _, a := range avp {
				if a.Code == uint32(268) && a.VenID == uint32(0) {
					a.Decode(c)
					break
				}
			}
		}
		if c != nil && *c == 2001 {
			p.state = open
			p.resetWatchdog()
		} else {
			e = fmt.Errorf("CEA Nack received")
			p.con.Close()
		}
	}
	if Notificator != nil {
		Notificator(&CapabilityExchangeEvent{
			Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
	}
	return
}

type eventRcvDWR struct {
	m msg.Message
}

func (eventRcvDWR) name() string {
	return "Rcv-DWR"
}

func (v eventRcvDWR) exec(p *Connection) (e error) {
	if p.state != open {
		e = fmt.Errorf("not acceptable message")
	}
	if Notificator != nil {
		Notificator(&WatchdogEvent{
			Tx: false, Req: true, Local: p.local, Peer: p.peer, Err: e})
	}
	if e != nil {
		return
	}

	a, _ := p.makeDWA(v.m)
	p.con.SetWriteDeadline(time.Now().Add(p.peer.Ts))
	_, e = a.WriteTo(p.con)
	if Notificator != nil {
		Notificator(&WatchdogEvent{
			Tx: true, Req: false, Local: p.local, Peer: p.peer, Err: e})
	}
	if e != nil {
		p.con.Close()
		p.state = shutdown
	} else {
		p.resetWatchdog()
	}
	return
}

type eventRcvDWA struct {
	m msg.Message
}

func (eventRcvDWA) name() string {
	return "Rcv-DWA"
}

func (v eventRcvDWA) exec(p *Connection) (e error) {
	if p.state != open {
		e = fmt.Errorf("not acceptable message")
	} else if ch, ok := p.sndstack[v.m.HbHID]; ok {
		delete(p.sndstack, v.m.HbHID)
		ch <- &v.m
		p.resetWatchdog()
	} else {
		e = fmt.Errorf("unknown DWA recieved")
	}
	if Notificator != nil {
		Notificator(&WatchdogEvent{
			Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
	}
	return
}

type eventRcvDPR struct {
	m msg.Message
}

func (eventRcvDPR) name() string {
	return "Rcv-DPR"
}

func (v eventRcvDPR) exec(p *Connection) (e error) {
	if p.state != open {
		e = fmt.Errorf("not acceptable message")
	}
	if Notificator != nil {
		Notificator(&PurgeEvent{
			Tx: false, Req: true, Local: p.local, Peer: p.peer, Err: e})
	}
	if e != nil {
		return
	}

	p.state = closing
	a, _ := p.makeDPA(v.m)
	p.con.SetWriteDeadline(time.Now().Add(p.peer.Ts))
	_, e = a.WriteTo(p.con)
	if Notificator != nil {
		Notificator(&PurgeEvent{
			Tx: true, Req: false, Local: p.local, Peer: p.peer, Err: e})
	}
	if e != nil {
		p.con.Close()
	}
	return
}

type eventRcvDPA struct {
	m msg.Message
}

func (eventRcvDPA) name() string {
	return "Rcv-DPA"
}

func (v eventRcvDPA) exec(p *Connection) (e error) {
	println("dpa")
	if p.state != closing {
		e = fmt.Errorf("not acceptable message")
	} else if ch, ok := p.sndstack[v.m.HbHID]; ok {

		delete(p.sndstack, v.m.HbHID)
		ch <- &v.m
		p.con.Close()
	} else {
		e = fmt.Errorf("unknown DPA recieved")
	}
	if Notificator != nil {
		Notificator(&PurgeEvent{
			Tx: false, Req: false, Local: p.local, Peer: p.peer, Err: e})
	}
	return
}

type eventSndMsg struct {
	m msg.Message
}

func (eventSndMsg) name() string {
	return "Snd-MSG"
}

func (v eventSndMsg) exec(p *Connection) error {
	if p.state != open {
		return fmt.Errorf("not acceptable event")
	}

	p.con.SetWriteDeadline(time.Now().Add(p.peer.Ts))
	_, e := v.m.WriteTo(p.con)
	if Notificator != nil {
		Notificator(&MessageEvent{
			Tx: true, Req: v.m.FlgR, Local: p.local, Peer: p.peer, Err: e})
	}
	if e != nil {
		p.con.Close()
	}
	return e
}

type eventRcvMsg struct {
	m msg.Message
}

func (eventRcvMsg) name() string {
	return "Rcv-MSG"
}

func (v eventRcvMsg) exec(p *Connection) (e error) {
	if p.state != open {
		return fmt.Errorf("not acceptable event")
	}

	if v.m.FlgR {
		p.rcvstack <- &v.m
		p.resetWatchdog()
	} else if ch, ok := p.sndstack[v.m.HbHID]; ok {
		ch <- &v.m
		p.resetWatchdog()
	} else {
		e = fmt.Errorf("unknown answer message received")
	}

	if Notificator != nil {
		Notificator(&MessageEvent{
			Tx: false, Req: v.m.FlgR, Local: p.local, Peer: p.peer, Err: e})
	}
	return
}
