package diameter

import (
	"bufio"
	"errors"
	"net"
	"time"
)

var (
	WDInterval = time.Second * 30 // WDInterval is watchdog send interval time
	WDMaxSend  = 3                // WDMaxSend is watchdog expired count

	TransactionWait = time.Second * 30 // TransactionWait is wait time for answer message

	Host    Identity // Local diameter hostname
	Realm   Identity // Local diameter realm
	stateID uint32   // Local diameter state ID

	OverrideCER []AVP // Override CER by this AVPs
	OverrideCEA []AVP // Override CEA by this AVPs
	OverrideDWR []AVP // Override DWR by this AVPs
	OverrideDWA []AVP // Override DWA by this AVPs
	OverrideDPR []AVP // Override DPR by this AVPs
	OverrideDPA []AVP // Override DPA by this AVPs
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

func overwriteAVP(org []AVP, overwrite []AVP) []AVP {
	m := map[uint32][]AVP{}
	for _, a := range overwrite {
		if avps, ok := m[a.Code]; ok {
			m[a.Code] = append(avps, a)
		} else {
			m[a.Code] = []AVP{a}
		}
	}
	res := []AVP{}
	for _, a := range org {
		if avps, ok := m[a.Code]; ok {
			for _, a := range avps {
				if a.Data != nil {
					res = append(res, a)
				}
			}
			m[a.Code] = nil
		} else if !ok {
			res = append(res, a)
		}
	}
	for _, v := range m {
		for _, a := range v {
			res = append(res, a)
		}
	}
	return res
}

// DialAndServe start handling Diameter connection as initiator
func (c *Connection) DialAndServe(con net.Conn) (e error) {
	if c.conn != nil || c.state != closed {
		return errors.New("reusing connection is not acceptable")
	}
	c.conn = con
	c.state = closed
	return c.serve()
}

// ListenAndServe start handling Diameter connection as responder
func (c *Connection) ListenAndServe(con net.Conn) (e error) {
	if c.conn != nil || c.state != closed {
		return errors.New("reusing connection is not acceptable")
	}
	c.conn = con
	c.state = waitCER
	return c.serve()
}

func (c *Connection) serve() error {
	c.notify = make(chan stateEvent, 65535)
	c.sndQueue = make(map[uint32]chan Message, 65535)
	c.rcvQueue = make(chan Message, 1024)
	c.commonApp = make(map[uint32]application)

	go func() {
		// read transport socket
		r := bufio.NewReader(c.conn)
		for {
			m := Message{}
			if err := m.UnmarshalFrom(r); err != nil {
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
			handleMsg(req)
		}
	}()

	if TraceEvent != nil {
		TraceEvent(c, shutdown.String(), eventInit{}.String(), nil)
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
			TraceEvent(c, old.String(), event.String(), err)
		}

		if _, ok := event.(eventPeerDisc); ok {
			break
		}
	}

	var e error
	if old != closing {
		e = errors.New("connection aborted")
	}
	return e
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
