package diameter

import (
	"errors"
	"net"
	"time"
)

var (
	WDInterval = time.Second * 30 // WDInterval is watchdog send interval time
	WDMaxSend  = 3                // WDMaxSend is watchdog expired count

	Host    Identity // Local diameter hostname
	Realm   Identity // Local diameter realm
	stateID uint32   // Local diameter state ID

	OverwriteAddr []net.IP // Overwrite IP addresses of local host in CER
)

// Connection of Diameter
type Connection struct {
	wdTimer *time.Timer // system message timer
	wdCount int         // watchdog expired counter

	Host    Identity // Peer diameter hostname
	Realm   Identity // Peer diameter realm
	stateID uint32   // Peer diameter state ID

	conn   net.Conn        // Transport connection
	notify chan stateEvent // state change notification queue
	state  conState        // current state

	sndQueue map[uint32]chan Message // Sending Request message queue
	rcvQueue chan Message            // Receiving Request message queue

	commonApp map[uint32]application
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
	c.sndQueue = make(map[uint32]chan Message, 65535)
	c.rcvQueue = make(chan Message, 65535)
	c.commonApp = make(map[uint32]application)

	go func() {
		// read transport socket
		for {
			m := Message{}
			if err := m.UnmarshalFrom(c.conn); err != nil {
				c.notify <- eventPeerDisc{reason: err}
				break
			}
			m.PeerName = c.Host
			m.PeerRealm = c.Realm

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
		for req, ok := <-c.rcvQueue; ok; req, ok = <-c.rcvQueue {
			ans, err := rxHandlerHelper(req)
			if err != nil {
				ans = DefaultRxHandler(req)
				ans.FlgR = false
				ans.HbHID = req.HbHID
				ans.EtEID = req.EtEID
			}
			if ans.AVPs != nil {
				c.notify <- eventSndMsg{ans}
			}
		}
	}()

	if TraceEvent != nil {
		TraceEvent(
			shutdown.String(), c.state.String(), eventInit{}.String(), nil)
	}

	if c.state != waitCER {
		c.notify <- eventConnect{}
	}

	var old conState
	for {
		event := <-c.notify
		old = c.state
		err := event.exec(c)
		if TraceEvent != nil {
			TraceEvent(old.String(), c.state.String(), event.String(), err)
		}

		if _, ok := event.(eventPeerDisc); ok {
			break
		}
	}

	if old != closing {
		if ConnectionAbortNotify != nil {
			ConnectionAbortNotify(c)
		}
		return errors.New("connection aborted")
	}

	if ConnectionDownNotify != nil {
		ConnectionDownNotify(c)
	}
	return nil
}

// Close Diameter connection and stop state machine.
func (c *Connection) Close(cause Enumerated) {
	if c.state == open {
		c.notify <- eventLock{}
		for len(c.rcvQueue) != 0 || len(c.sndQueue) != 0 {
			time.Sleep(time.Millisecond * 100)
		}
	}
	c.notify <- eventStop{cause}
}
