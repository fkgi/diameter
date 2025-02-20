package diameter

import (
	"bytes"
	"errors"
	"time"
)

// Acceptable Application-ID and commands of the application.
// Empty map indicate that accept any application.
var applications = make(map[uint32]application)

type application struct {
	venID    uint32
	handlers map[uint32]Handler
}

// Handler handles Diameter message.
// Inputs are Retry flag and AVPs of Request. Outputs are Error flag and AVPs of Answer.
type Handler func(bool, []AVP) (bool, []AVP)

// Router select destination peer for specific message.
type Router func(Message) *Connection

// Handle registers Diameter request handler for specified command.
// Input Handler is called when when request is from peer.
// Output Handler is used when send request to peer.
func Handle(code, appID, venID uint32, h Handler, rt Router) Handler {
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

		var err error
		if rt == nil {
			err = errors.New("no route found")
		} else if c := rt(m); c == nil {
			err = errors.New("no route found")
		} else {
			m = c.send(m)
			avp, err = m.GetAVP()
		}

		if err != nil {
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
	return m.GenerateAnswerBy(UnableToDeliver)
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
				FlgR: false, FlgP: req.FlgP, FlgE: true, FlgT: false,
				Code: req.Code, AppID: req.AppID,
				HbHID: req.HbHID, EtEID: req.EtEID,
				AVPs: buf.Bytes()}
			return
		}
		avp = append(avp, a)
	}

	req.FlgE, avp = f(req.FlgT, avp)
	ans = Message{
		FlgR: false, FlgP: req.FlgP, FlgE: req.FlgE, FlgT: false,
		Code: req.Code, AppID: req.AppID,
		HbHID: req.HbHID, EtEID: req.EtEID}

	if avp != nil {
		buf := new(bytes.Buffer)
		for _, a := range avp {
			a.MarshalTo(buf)
		}
		ans.AVPs = buf.Bytes()
	}

	return
}

// DefaultTxHandler for sending Diameter request message without Handler or relay application.
func (c *Connection) DefaultTxHandler(m Message) Message {
	if c.state != open {
		return m.GenerateAnswerBy(UnableToDeliver)
	}
	if _, ok := c.commonApp[m.AppID]; !ok && len(c.commonApp) != 0 {
		return m.GenerateAnswerBy(UnableToDeliver)
	}

	m.HbHID = nextHbH()
	m.FlgR = true
	m.FlgE = false

	return c.send(m)
}

func (c *Connection) send(m Message) Message {
	if c.state != open {
		return m.GenerateAnswerBy(UnableToDeliver)
	}

	ch := make(chan Message)
	c.notify <- eventSndMsg{m, ch}

	t := time.AfterFunc(WDInterval, func() {
		c.notify <- eventRcvAns{m.GenerateAnswerBy(TooBusy)}
	})
	r, ok := <-ch
	t.Stop()

	if !ok {
		m = m.GenerateAnswerBy(UnableToDeliver)
	} else if m.Code != r.Code || m.AppID != r.AppID || m.EtEID != r.EtEID {
		m = m.GenerateAnswerBy(UnableToDeliver)
	} else {
		m = r
	}
	return m
}
