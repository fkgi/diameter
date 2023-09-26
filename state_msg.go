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
	RxReq++
	if c.state == closed || c.state == waitCER || c.state == waitCEA {
		RejectReq++
		return notAcceptableEvent{e: v, s: c.state}
	}
	TraceMessage(v.m, Rx, nil)

	result := Success
	var err error
	if c.state == locked {
		result = UnableToDeliver
	} else if len(rcvQueue) == cap(rcvQueue) {
		result = TooBusy
		err = errors.New("too busy, receive queue is full")
	} else if len(applications) == 0 {
		rcvQueue <- v.m
	} else if _, ok := applications[v.m.AppID]; ok {
		rcvQueue <- v.m
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
		ans := v.m.generateAnswerBy(result)
		if e := ans.MarshalTo(c.conn); e != nil {
			TxAnsFail++
			c.conn.Close()
			err = e
		} else {
			TraceMessage(ans, Tx, err)
			CountTxCode(result)
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
		InvalidAns++
		return notAcceptableEvent{e: v, s: c.state}
	}

	TraceMessage(v.m, Rx, nil)

	if ch, ok := sndQueue[v.m.HbHID]; ok {
		ch <- v.m
	} else {
		InvalidAns++
		return unknownAnswer(v.m.HbHID)
	}
	delete(sndQueue, v.m.HbHID)

	if c.wdCount == 0 {
		c.wdTimer.Stop()
		c.wdTimer.Reset(WDInterval)
	}

	return nil
}
