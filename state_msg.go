package diameter

import (
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
	c.RxReq++
	if c.state == closed || c.state == waitCER || c.state == waitCEA {
		c.DscardReq++
		return notAcceptableEvent{e: v, s: c.state}
	}
	if TraceMessage != nil {
		TraceMessage(v.m, Rx, nil)
	}

	result := Success
	var err error
	if c.state == locked {
		result = UnableToDeliver
	} else if len(c.rcvQueue) == cap(c.rcvQueue) {
		result = TooBusy
		err = errors.New("too busy, receive queue is full")
	} else if len(c.commonApp) == 0 {
		c.rcvQueue <- v.m
	} else if _, ok := c.commonApp[v.m.AppID]; ok {
		c.rcvQueue <- v.m
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
			c.conn.Close()
			err = e
		} else {
			if TraceMessage != nil {
				TraceMessage(ans, Tx, err)
			}
			c.countTxCode(result)
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
	if c.state != open && c.state != locked {
		c.InvalidAns++
		return notAcceptableEvent{e: v, s: c.state}
	}

	if TraceMessage != nil {
		TraceMessage(v.m, Rx, nil)
	}

	if ch, ok := c.sndQueue[v.m.HbHID]; ok {
		ch <- v.m
	} else {
		c.InvalidAns++
		return unknownAnswer(v.m.HbHID)
	}
	delete(c.sndQueue, v.m.HbHID)

	if c.wdCount == 0 {
		c.wdTimer.Stop()
		c.wdTimer.Reset(WDInterval)
	}

	return nil
}
