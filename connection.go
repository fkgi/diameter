package diameter

import (
	"bytes"
	"net"
	"time"
)

// DialAndServe start diameter connection handling process as initiator.
func DialAndServe(localFqdn, peerFqdn string, sctp bool) (err error) {
	Local.Host, Local.Realm, err = ResolveIdentiry(localFqdn)
	if err != nil {
		return
	}
	Peer.Host, Peer.Realm, err = ResolveIdentiry(peerFqdn)
	if err != nil {
		return
	}

	if sctp {

	} else {
		var la *net.TCPAddr
		la, err = net.ResolveTCPAddr("tcp", localFqdn)
		if err != nil {
			return
		}
		conn, err = (&net.Dialer{LocalAddr: la}).Dial("tcp", peerFqdn)
		if err != nil {
			return
		}
	}

	go socketHandler()
	go messageHandler()
	TraceState(shutdown.String(), state.String(), eventInit{}.String(), nil)

	notify <- eventConnect{}
	for {
		event := <-notify
		old := state
		TraceState(old.String(), state.String(), event.String(), event.exec())

		if _, ok := event.(eventPeerDisc); ok {
			break
		}
	}

	return err
}

// ListenAndServe start diameter connection handling process as responder.
func ListenAndServe(fqdn string, sctp bool) (err error) {
	Local.Host, Local.Realm, err = ResolveIdentiry(fqdn)
	if err != nil {
		return
	}

	if sctp {

	} else {
		var l net.Listener
		l, err = net.Listen("tcp", fqdn)
		if err != nil {
			return
		}

		t := time.AfterFunc(WDInterval, func() {
			l.Close()
		})
		conn, err = l.Accept()
		t.Stop()
		if err != nil {
			return
		}
	}

	state = waitCER
	go socketHandler()
	go messageHandler()
	TraceState(shutdown.String(), state.String(), eventInit{}.String(), nil)

	for {
		event := <-notify
		old := state
		TraceState(old.String(), state.String(), event.String(), event.exec())

		if _, ok := event.(eventPeerDisc); ok {
			break
		}
	}

	return err
}

func socketHandler() {
	for {
		m := Message{}
		conn.SetReadDeadline(time.Time{})
		if err := m.UnmarshalFrom(conn); err != nil {
			break
		}

		if m.AppID == 0 && m.Code == 257 && m.FlgR {
			notify <- eventRcvCER{m}
		} else if m.AppID == 0 && m.Code == 257 && !m.FlgR {
			notify <- eventRcvCEA{m}
		} else if m.AppID == 0 && m.Code == 280 && m.FlgR {
			notify <- eventRcvDWR{m}
		} else if m.AppID == 0 && m.Code == 280 && !m.FlgR {
			notify <- eventRcvDWA{m}
		} else if m.AppID == 0 && m.Code == 282 && m.FlgR {
			notify <- eventRcvDPR{m}
		} else if m.AppID == 0 && m.Code == 282 && !m.FlgR {
			notify <- eventRcvDPA{m}
		} else if m.FlgR {
			notify <- eventRcvReq{m}
		} else {
			notify <- eventRcvAns{m}
		}
	}
	notify <- eventPeerDisc{}
}

func messageHandler() {
	for req, ok := <-rcvStack; ok; req, ok = <-rcvStack {
		/*
			avps := make([]AVP, 0, 10)
			var err error
			for rdr := bytes.NewReader(req.AVPs); rdr.Len() != 0; {
				a := AVP{}
				if err = a.wrapedUnmarshalFrom(rdr); err != nil {
					break
				}
				avps = append(avps, a)
			}
			if iavp, ok := err.(InvalidAVP); ok {
				notify <- eventSndMsg{req.generateAnswerBy(iavp.Code)}
			} else {
		*/
		var avps []byte
		var flgE bool

		if app, ok := applications[req.AppID]; ok {
			if f, ok := app.handlers[req.Code]; ok {
				flgE, avps = f(req.FlgT, req.AVPs)
			}
		}
		if avps == nil {
			flgE, avps = DefaultHandler(req)
		}

		notify <- eventSndMsg{Message{
			FlgR: false, FlgP: req.FlgP, FlgE: flgE, FlgT: false,
			Code: req.Code, AppID: req.AppID,
			HbHID: req.HbHID, EtEID: req.EtEID,
			AVPs: avps}}
		/*
			}
		*/
	}
}

// message -> error, avps
var DefaultHandler func(Message) (bool, []byte)

func init() {
	DefaultHandler = func(_ Message) (bool, []byte) {
		buf := new(bytes.Buffer)
		SetResultCode(UnableToDeliver).MarshalTo(buf)
		SetOriginHost(Local.Host).MarshalTo(buf)
		SetOriginRealm(Local.Realm).MarshalTo(buf)
		return true, buf.Bytes()
	}
}

// Send Diameter request
func Send(code, appID, msgID uint32, retry bool, avp []byte) []byte {
	buf := new(bytes.Buffer)
	buf.Write(avp)
	if Router {
		SetRouteRecord(Local.Host).MarshalTo(buf)
	}

	m := Message{
		FlgR: true, FlgP: true, FlgE: false, FlgT: retry,
		Code: code, AppID: appID,
		HbHID: nextHbH(), EtEID: msgID,
		AVPs: buf.Bytes()}

	if state != open {
		m = m.generateAnswerBy(UnableToDeliver)
	} else if _, ok := applications[appID]; !ok && len(applications) != 0 {
		m = m.generateAnswerBy(UnableToDeliver)
	} else {
		ch := make(chan Message)
		sndStack[m.HbHID] = ch
		notify <- eventSndMsg{m}

		t := time.AfterFunc(WDInterval, func() {
			notify <- eventRcvAns{m.generateAnswerBy(TooBusy)}
		})
		r, ok := <-ch
		t.Stop()

		if !ok {
			m = m.generateAnswerBy(UnableToDeliver)
		} else if m.Code != r.Code || m.AppID != r.AppID || m.EtEID != r.EtEID {
			m = m.generateAnswerBy(UnableToDeliver)
		} else {
			m = r
		}
	}

	return m.AVPs
}

// Handle Diameter request
func Handle(code, appID, venID uint32, handler func(bool, []byte) (bool, []byte)) {
	if _, ok := applications[appID]; !ok {
		applications[appID] = application{
			venID:    venID,
			handlers: make(map[uint32]func(bool, []byte) (bool, []byte))}
	}
	applications[appID].handlers[code] = handler
}

// Close stop state machine
func Close(cause Enumerated) {
	if state == open {
		notify <- eventLock{}
		for len(rcvStack) != 0 || len(sndStack) != 0 {
			time.Sleep(time.Millisecond * 100)
		}
	}
	notify <- eventStop{cause}
}
