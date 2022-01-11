package diameter

import (
	"bytes"
	"io"
	"net"
	"time"
)

// DialAndServe start diameter connection handling process as initiator.
func DialAndServe(c net.Conn) (err error) {
	state = closed
	return serve(c)
}

// ListenAndServe start diameter connection handling process as responder.
func ListenAndServe(c net.Conn) (err error) {
	state = waitCER
	return serve(c)
}

func serve(c net.Conn) error {
	if c == nil {
		return io.ErrUnexpectedEOF
	}
	conn = c
	if Peer.Host == "" {
		if names, err := net.LookupAddr(c.RemoteAddr().String()); err == nil {
			Peer.Host, Peer.Realm, _ = ResolveIdentiry(names[0])
		}
	}

	go func() {
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
	}()
	go func() {
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
				if f, ok := app.handlers[req.Code]; ok && f != nil {
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
	}()
	TraceState(shutdown.String(), state.String(), eventInit{}.String(), nil)

	if state != waitCER {
		notify <- eventConnect{}
	}
	for {
		event := <-notify
		old := state
		TraceState(old.String(), state.String(), event.String(), event.exec())

		if _, ok := event.(eventPeerDisc); ok {
			break
		}
	}

	return nil
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
func Send(m Message) (bool, []byte) {
	m.FlgR = true
	m.FlgE = false
	m.HbHID = nextHbH()

	if state != open {
		m = m.generateAnswerBy(UnableToDeliver)
	} else if _, ok := applications[m.AppID]; !ok && len(applications) != 0 {
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

	return m.FlgE, m.AVPs
}

// Handler handles Diameter message.
// Inputs are Retry flag and AVPs of Request, Outputs are Error flag and AVPs of Answer.
type Handler func(bool, []byte) (bool, []byte)

// Handle registers Diameter request handler for specified command.
func Handle(code, appID, venID uint32, h Handler) Handler {
	if _, ok := applications[appID]; !ok {
		applications[appID] = application{
			venID:    venID,
			handlers: make(map[uint32]func(bool, []byte) (bool, []byte))}
	}
	applications[appID].handlers[code] = h

	return func(r bool, avp []byte) (bool, []byte) {
		if Router {
			buf := new(bytes.Buffer)
			buf.Write(avp)
			SetRouteRecord(Local.Host).MarshalTo(buf)
			avp = buf.Bytes()
		}

		return Send(Message{
			FlgR: true, FlgP: true, FlgE: false, FlgT: r,
			Code: code, AppID: appID,
			HbHID: 0, EtEID: nextEtE(),
			AVPs: avp})
	}
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
