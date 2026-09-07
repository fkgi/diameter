package diameter

import (
	"bytes"
	"math/rand/v2"
	"time"
)

const (
	minWorkers = 100
	maxWorkers = 20000 - minWorkers
)

var (
	sharedQ       = make(chan Message, maxWorkers+minWorkers)
	activeWorkers = make(chan int, 1)
)

func init() {
	activeWorkers <- 0
	for range minWorkers {
		go func() {
			for req, ok := <-sharedQ; ok; req, ok = <-sharedQ {
				handleMsg(req)
			}
		}()
	}
	go func() {
		act := true
		for act {
			acl := len(sharedQ)
			if acl < minWorkers/2 {
				time.Sleep(time.Millisecond * 10)
				continue
			}

			a := <-activeWorkers
			if a+acl > maxWorkers {
				acl = maxWorkers - a
			}
			a += acl
			activeWorkers <- a

			for range acl {
				go func() {
					for c := 0; c < 500; c++ {
						if len(sharedQ) < minWorkers/2 {
							time.Sleep(time.Millisecond * time.Duration(8+rand.IntN(4)))
						} else if req, ok := <-sharedQ; !ok {
							act = false
							break
						} else {
							handleMsg(req)
							c = 0
						}
					}
					activeWorkers <- (<-activeWorkers - 1)
				}()
			}
			time.Sleep(time.Millisecond * 10)
		}
	}()
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
		req.notify <- eventSndMsg{ans, nil, time.Now()}
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
				AVPs: buf.Bytes()}, nil, time.Now()}
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
			AVPs: buf.Bytes()}, nil, time.Now()}
	}
}
