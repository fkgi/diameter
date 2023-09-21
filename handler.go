package diameter

import (
	"bytes"
	"time"
)

// Handler handles Diameter message.
// Inputs are Retry flag and AVPs of Request. Outputs are Error flag and AVPs of Answer.
type Handler func(bool, []AVP) (bool, []AVP)

// Handle registers Diameter request handler for specified command.
// Input Handler is called when when request is from peer.
// Output Handler is used when send request to peer.
func Handle(code, appID, venID uint32, h Handler) Handler {
	if _, ok := applications[appID]; !ok {
		applications[appID] = application{
			venID:    venID,
			handlers: make(map[uint32]Handler)}
	}
	applications[appID].handlers[code] = h

	return func(r bool, avp []AVP) (bool, []AVP) {
		m := Message{
			FlgR: true, FlgP: true, FlgE: false, FlgT: r,
			Code: code, AppID: appID,
			HbHID: nextHbH(), EtEID: nextEtE()}
		m.setAVP(avp)
		m = send(m)

		var e error
		avp, e = m.getAVP()
		if e != nil {
			return true, []AVP{
				SetResultCode(UnableToDeliver),
				SetOriginHost(Local.Host),
				SetOriginRealm(Local.Realm)}
		}
		return m.FlgE, avp
	}
}

// DefaultRxHandler for receiving Diameter request message without Handler or ralay application.
var DefaultRxHandler func(Message) Message = func(m Message) Message {
	return m.generateAnswerBy(UnableToDeliver)
}

// DefaultTxHandler for sending Diameter request message without Handler or relay application.
func DefaultTxHandler(m Message) Message {
	if _, ok := applications[m.AppID]; !ok || len(applications) != 0 {
		return m.generateAnswerBy(UnableToDeliver)
	}

	m.HbHID = nextHbH()
	m.FlgR = true
	m.FlgP = true
	m.FlgE = false

	return send(m)
}

func send(m Message) Message {
	if Router {
		buf := bytes.NewBuffer(m.AVPs)
		SetRouteRecord(Local.Host).MarshalTo(buf)
		m.AVPs = buf.Bytes()
	}

	if state != open {
		m = m.generateAnswerBy(UnableToDeliver)
	} else {
		ch := make(chan Message)
		sndQueue[m.HbHID] = ch
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

	return m
}
