package diameter

import (
	"errors"
	"time"
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
	} else if len(rcvStack) == cap(rcvStack) {
		result = TooBusy
		err = errors.New("too busy, recieve stack is full")
	} else if len(applications) == 0 {
		rcvStack <- v.m
	} else if _, ok := applications[v.m.AppID]; ok {
		rcvStack <- v.m
	} else {
		result = ApplicationUnsupported
		err = InvalidMessage(result)
	}

	if wdCount == 0 {
		wdTimer.Stop()
		wdTimer.Reset(WDInterval)
	}
	if result != Success {
		ans := v.m.generateAnswerBy(result)
		conn.SetWriteDeadline(time.Now().Add(TxTimeout))
		if e := ans.MarshalTo(conn); e != nil {
			TxAnsFail++
			conn.Close()
			err = e
		} else {
			TraceMessage(ans, Tx, err)
			countTxCode(result)
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

	if ch, ok := sndStack[v.m.HbHID]; ok {
		ch <- v.m
	} else {
		InvalidAns++
		return unknownAnswer(v.m.HbHID)
	}
	delete(sndStack, v.m.HbHID)

	if wdCount == 0 {
		wdTimer.Stop()
		wdTimer.Reset(WDInterval)
	}

	return nil
}
