package diameter

import (
	"bytes"
)

var sharedQ = make(chan Message, 65535)

func init() {
	for i := 0; i < 10; i++ {
		go func() {
			for req, ok := <-sharedQ; ok; req, ok = <-sharedQ {
				handleMsg(req)
			}
		}()
	}
}

func handleMsg(req Message) {
	var f Handler
	if app, ok := applications[req.AppID]; !ok {
		f = nil
	} else if f, ok = app.handlers[req.Code]; !ok {
		f = nil
	}
	if f == nil {
		ans := DefaultRxHandler(req)
		ans.FlgR = false
		ans.HbHID = req.HbHID
		ans.EtEID = req.EtEID
		req.notify <- eventSndMsg{ans, nil}
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

			req.notify <- eventSndMsg{Message{
				FlgR: false, FlgP: req.FlgP, FlgE: true, FlgT: false,
				Code: req.Code, AppID: req.AppID,
				HbHID: req.HbHID, EtEID: req.EtEID,
				AVPs: buf.Bytes()}, nil}
			return
		}
		avp = append(avp, a)
	}

	if req.FlgE, avp = f(req.FlgT, avp); avp != nil {
		buf := new(bytes.Buffer)
		for _, a := range avp {
			a.MarshalTo(buf)
		}

		req.notify <- eventSndMsg{Message{
			FlgR: false, FlgP: req.FlgP, FlgE: req.FlgE, FlgT: false,
			Code: req.Code, AppID: req.AppID,
			HbHID: req.HbHID, EtEID: req.EtEID,
			AVPs: buf.Bytes()}, nil}
	}
}
