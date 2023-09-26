package diameter

import (
	"bytes"
	"net"
	"time"
)

type Connection struct {
	wdTimer *time.Timer // system message timer
	wdCount int         //= 0         // watchdog expired counter

	Host    Identity // Peer diameter hostname
	Realm   Identity // Peer diameter realm
	stateID uint32   // Peer diameter state ID

	conn   net.Conn        // Transport connection
	notify chan stateEvent //= make(chan stateEvent, 16)            // state change notification queue
	state  conState        //= closed                               // current state
}

func (c *Connection) DialAndServe(con net.Conn) (e error) {
	c.conn = con
	c.state = closed
	return c.serve()
}

func (c *Connection) ListenAndServe(con net.Conn) (e error) {
	c.conn = con
	c.state = waitCER
	return c.serve()
}

func (c *Connection) serve() error {
	c.notify = make(chan stateEvent, 16)
	go func() {
		for {
			m := Message{}

			// conn.SetReadDeadline(time.Time{})
			if err := m.UnmarshalFrom(c.conn); err != nil {
				c.notify <- eventPeerDisc{reason: err}
				break
			}

			if m.AppID == 0 && m.Code == 257 && m.FlgR {
				c.notify <- eventRcvCER{m}
			} else if m.AppID == 0 && m.Code == 257 && !m.FlgR {
				c.notify <- eventRcvCEA{m}
			} else if m.AppID == 0 && m.Code == 280 && m.FlgR {
				c.notify <- eventRcvDWR{m}
			} else if m.AppID == 0 && m.Code == 280 && !m.FlgR {
				c.notify <- eventRcvDWA{m}
			} else if m.AppID == 0 && m.Code == 282 && m.FlgR {
				c.notify <- eventRcvDPR{m}
			} else if m.AppID == 0 && m.Code == 282 && !m.FlgR {
				c.notify <- eventRcvDPA{m}
			} else if m.FlgR {
				c.notify <- eventRcvReq{m}
			} else {
				c.notify <- eventRcvAns{m}
			}
		}
	}()
	go func() {
	rxNewMsg:
		for req, ok := <-rcvQueue; ok; req, ok = <-rcvQueue {
			if app, ok := applications[req.AppID]; ok {
				if f, ok := app.handlers[req.Code]; ok && f != nil {
					avp := make([]AVP, 0, avpBufferSize)
					for rdr := bytes.NewReader(req.AVPs); rdr.Len() != 0; {
						a := AVP{}
						if e := a.UnmarshalFrom(rdr); e != nil {
							buf := new(bytes.Buffer)
							SetResultCode(InvalidAvpValue).MarshalTo(buf)
							SetOriginHost(Host).MarshalTo(buf)
							SetOriginRealm(Realm).MarshalTo(buf)
							c.notify <- eventSndMsg{Message{
								FlgR: false, FlgP: req.FlgP, FlgE: req.FlgE, FlgT: false,
								Code: req.Code, AppID: req.AppID,
								HbHID: req.HbHID, EtEID: req.EtEID,
								AVPs: buf.Bytes()}}
							continue rxNewMsg
						}
						avp = append(avp, a)
					}

					req.FlgE, avp = f(req.FlgT, avp)
					buf := new(bytes.Buffer)
					for _, a := range avp {
						a.MarshalTo(buf)
					}
					c.notify <- eventSndMsg{Message{
						FlgR: false, FlgP: req.FlgP, FlgE: req.FlgE, FlgT: false,
						Code: req.Code, AppID: req.AppID,
						HbHID: req.HbHID, EtEID: req.EtEID,
						AVPs: buf.Bytes()}}
					continue rxNewMsg
				}
			}
			c.notify <- eventSndMsg{DefaultRxHandler(req)}
		}
	}()

	TraceEvent(shutdown.String(), c.state.String(), eventInit{}.String(), nil)

	if c.state != waitCER {
		c.notify <- eventConnect{}
	}
	for {
		event := <-c.notify
		old := c.state
		err := event.exec(c)
		TraceEvent(old.String(), c.state.String(), event.String(), err)

		if _, ok := event.(eventPeerDisc); ok {
			break
		}
	}

	return nil
}

// Close Diameter connection and stop state machine.
func (c *Connection) Close(cause Enumerated) {
	if c.state == open {
		c.notify <- eventLock{}
		for len(rcvQueue) != 0 || len(sndQueue) != 0 {
			time.Sleep(time.Millisecond * 100)
		}
	}
	c.notify <- eventStop{cause}
}
