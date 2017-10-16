package sock

import (
	"time"

	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/rfc6733"
)

// Session is diameter message session
type Session struct {
	id string
	c  *Conn
}

// Send send Diameter request
func (s *Session) Send(m msg.Request) msg.Answer {
	req := m.ToRaw()
	req.HbHID = nextHbH()
	req.EtEID = nextEtE()
	// req.AVP = append(req.AVP, s.id.ToRaw())

	ch := make(chan msg.RawMsg)
	s.c.sndstack[req.HbHID] = ch
	s.c.notify <- eventSndMsg{m: req}

	t := time.AfterFunc(s.c.peer.SndTimeout, func() {
		m := m.Failed(
			uint32(rfc6733.DiameterTooBusy),
			"no response from peer node").ToRaw()
		m.HbHID = req.HbHID
		m.EtEID = req.EtEID
		s.c.notify <- eventRcvMsg{m}
	})

	a := <-ch
	t.Stop()
	if a.Code == 0 {
		return m.Failed(
			uint32(rfc6733.DiameterUnableToDeliver),
			"failed to send")
	}

	app, ok := supportedApps[a.AppID]
	if !ok {
		return m.Failed(
			uint32(rfc6733.DiameterUnableToComply),
			"invalid Application-ID answer")
	}
	ans, ok := app.ans[a.Code]
	if !ok {
		return m.Failed(
			uint32(rfc6733.DiameterUnableToComply),
			"invalid Command-Code answer")
	}
	ack, e := ans.FromRaw(a)
	if e != nil {
		return m.Failed(
			uint32(rfc6733.DiameterUnableToComply),
			"invalid data answer")
		// ToDo
		// invalid message handling
	}
	return ack
}
