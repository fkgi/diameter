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

func (v eventRcvReq) exec() error {
	RxReq++
	if state == closed || state == waitCER || state == waitCEA {
		RejectReq++
		return notAcceptableEvent{e: v, s: state}
	}
	TraceMessage(v.m, Rx, nil)

	result := Success
	var err error
	if state == locked {
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

	if wdCount == 0 {
		wdTimer.Stop()
		wdTimer.Reset(WDInterval)
	}
	if result != Success {
		ans := v.m.generateAnswerBy(result)
		if e := ans.MarshalTo(conn); e != nil {
			TxAnsFail++
			conn.Close()
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

func (v eventRcvAns) exec() (e error) {
	if state != open && state != locked {
		InvalidAns++
		return notAcceptableEvent{e: v, s: state}
	}

	TraceMessage(v.m, Rx, nil)

	if ch, ok := sndQueue[v.m.HbHID]; ok {
		ch <- v.m
	} else {
		InvalidAns++
		return unknownAnswer(v.m.HbHID)
	}
	delete(sndQueue, v.m.HbHID)

	if wdCount == 0 {
		wdTimer.Stop()
		wdTimer.Reset(WDInterval)
	}

	return nil
}
