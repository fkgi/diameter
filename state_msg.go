package diameter

import (
	"bytes"
	"errors"
	"fmt"
)

type eventRcvReq struct {
	m Message
}

func (eventRcvReq) String() string {
	return "Rcv-REQ"
}

func (v eventRcvReq) exec(c *Connection) error {
	var err error
	if c.state == closed || c.state == waitCER || c.state == waitCEA {
		err = RejectRxMessage{
			State: c.state, ErrMsg: "Request Message is not acceptable"}
	}
	if TraceMessage != nil {
		TraceMessage(v.m, Rx, err)
	}
	if err != nil {
		return err
	}

	var ch chan Message
	// Auth-Session-State AVP=STATE_MAINTAINED
	if bytes.Contains(v.m.AVPs, []byte{
		0x00, 0x00, 0x01, 0x15, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x00}) {
		ch = c.rcvQueue
	} else {
		ch = sharedQ
	}
	v.m.notify = c.notify

	result := Success
	if c.state == locked {
		result = UnableToDeliver
	} else if len(ch) == cap(ch) {
		result = TooBusy
		err = errors.New("too busy, receive queue is full")
	} else if len(c.commonApp) == 0 {
		ch <- v.m
	} else if _, ok := c.commonApp[v.m.AppID]; ok {
		ch <- v.m
	} else {
		result = ApplicationUnsupported
		err = InvalidMessage{
			Code:   result,
			ErrMsg: fmt.Sprintf("unknown application %d", v.m.AppID)}
	}

	if c.wdCount == 0 {
		c.wdTimer.Stop()
		c.wdTimer.Reset(WDInterval)
	}
	if result != Success {
		ans := v.m.GenerateAnswerBy(result)
		if e := ans.MarshalTo(c.conn); e != nil {
			err = e
			c.notify <- eventPeerDisc{reason: err}
		}
		if TraceMessage != nil {
			TraceMessage(ans, Tx, err)
		}
	}

	return err
}

type eventRcvAns struct {
	m Message
}

func (eventRcvAns) String() string {
	return "Rcv-ANS"
}

func (v eventRcvAns) exec(c *Connection) (e error) {
	var err error

	if c.state != open && c.state != locked {
		err = InvalidMessage{
			Code:   UnableToComply,
			ErrMsg: "Answer Message is not acceptable in " + c.state.String() + " state"}
	} else if ch, ok := c.sndQueue[v.m.HbHID]; ok {
		ch <- v.m
	} else {
		err = InvalidMessage{
			Code:   UnableToComply,
			ErrMsg: "correlated request with the Hop-by-Hop ID not found"}
	}

	if TraceMessage != nil {
		TraceMessage(v.m, Rx, err)
	}
	if err == nil {
		delete(c.sndQueue, v.m.HbHID)

		if c.wdCount == 0 {
			c.wdTimer.Stop()
			c.wdTimer.Reset(WDInterval)
		}
	}

	return err
}
