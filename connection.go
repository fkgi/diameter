package diameter

import (
	"errors"
	"net"
	"time"
)

type Connection struct {
	wdTimer *time.Timer // system message timer
	wdCount int         // watchdog expired counter

	Host    Identity // Peer diameter hostname
	Realm   Identity // Peer diameter realm
	stateID uint32   // Peer diameter state ID

	conn   net.Conn        // Transport connection
	notify chan stateEvent // state change notification queue
	state  conState        // current state
}

func (c *Connection) DialAndServe(con net.Conn) (e error) {
	if c.conn != nil || c.state != closed {
		return errors.New("reusing connection is not acceptable")
	}
	c.conn = con
	return c.serve()
}

func (c *Connection) ListenAndServe(con net.Conn) (e error) {
	if c.conn != nil || c.state != closed {
		return errors.New("reusing connection is not acceptable")
	}
	c.conn = con
	c.state = waitCER
	return c.serve()
}

func (c *Connection) serve() error {
	c.notify = make(chan stateEvent, 16)
	go func() {
		// read transport socket
		for {
			m := Message{}
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
		// handle Rx Diameter message
		for req, ok := <-rcvQueue; ok; req, ok = <-rcvQueue {
			ans, err := rxHandlerHelper(req)
			if err != nil {
				ans = DefaultRxHandler(req)
				ans.FlgR = false
				ans.HbHID = req.HbHID
				ans.EtEID = req.EtEID
			}
			c.notify <- eventSndMsg{ans}
		}
	}()

	if TraceEvent != nil {
		TraceEvent(
			shutdown.String(), c.state.String(), eventInit{}.String(), nil)
	}

	if c.state != waitCER {
		c.notify <- eventConnect{}
	}
	for {
		event := <-c.notify
		old := c.state
		err := event.exec(c)
		if TraceEvent != nil {
			TraceEvent(old.String(), c.state.String(), event.String(), err)
		}

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
