package diameter

import (
	"bytes"
	"errors"
	"time"
)

// Handler handles Diameter message.
// Inputs are Retry flag and AVPs of Request. Outputs are Error flag and AVPs of Answer.
type Handler func(bool, []AVP) (bool, []AVP)

// Handle registers Diameter request handler for specified command.
// Input Handler is called when when request is from peer.
// Output Handler is used when send request to peer.
func Handle(code, appID, venID uint32, h Handler, s func() *Connection) Handler {
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
		m.SetAVP(avp)
		m = s().send(m)

		var e error
		avp, e = m.GetAVP()
		if e != nil {
			return true, []AVP{
				SetResultCode(UnableToDeliver),
				SetOriginHost(Host),
				SetOriginRealm(Realm)}
		}
		return m.FlgE, avp
	}
}

// DefaultRxHandler for receiving Diameter request message without Handler or ralay application.
var DefaultRxHandler func(Message) Message = func(m Message) Message {
	return m.generateAnswerBy(UnableToDeliver)
}

func rxHandlerHelper(req Message) (ans Message, err error) {
	app, ok := applications[req.AppID]
	if !ok {
		err = errors.New("application not registered")
		return
	}
	f, ok := app.handlers[req.Code]
	if !ok || f == nil {
		err = errors.New("command not registered")
		return
	}

	avp := make([]AVP, 0, avpBufferSize)
	for rdr := bytes.NewReader(req.AVPs); rdr.Len() != 0; {
		a := AVP{}
		if e := a.UnmarshalFrom(rdr); e != nil {
			buf := new(bytes.Buffer)
			SetResultCode(InvalidAvpValue).MarshalTo(buf)
			SetOriginHost(Host).MarshalTo(buf)
			SetOriginRealm(Realm).MarshalTo(buf)

			ans = Message{
				FlgR: false, FlgP: req.FlgP, FlgE: req.FlgE, FlgT: false,
				Code: req.Code, AppID: req.AppID,
				HbHID: req.HbHID, EtEID: req.EtEID,
				AVPs: buf.Bytes()}
			return
		}
		avp = append(avp, a)
	}

	req.FlgE, avp = f(req.FlgT, avp)
	buf := new(bytes.Buffer)
	for _, a := range avp {
		a.MarshalTo(buf)
	}

	ans = Message{
		FlgR: false, FlgP: req.FlgP, FlgE: req.FlgE, FlgT: false,
		Code: req.Code, AppID: req.AppID,
		HbHID: req.HbHID, EtEID: req.EtEID,
		AVPs: buf.Bytes()}
	return
}

// DefaultTxHandler for sending Diameter request message without Handler or relay application.
func (c *Connection) DefaultTxHandler(m Message) Message {
	if _, ok := applications[m.AppID]; !ok && len(applications) != 0 {
		return m.generateAnswerBy(UnableToDeliver)
	}

	buf := bytes.NewBuffer(m.AVPs)
	SetRouteRecord(Host).MarshalTo(buf)
	m.AVPs = buf.Bytes()

	m.HbHID = nextHbH()
	m.FlgR = true
	m.FlgE = false

	return c.send(m)
}

func (c *Connection) send(m Message) Message {
	if c.state != open {
		return m.generateAnswerBy(UnableToDeliver)
	}

	ch := make(chan Message)
	sndQueue[m.HbHID] = ch
	c.notify <- eventSndMsg{m}

	t := time.AfterFunc(WDInterval, func() {
		c.notify <- eventRcvAns{m.generateAnswerBy(TooBusy)}
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
	return m
}
