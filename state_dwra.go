package diameter

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

/*
DeviceWatchdogRequest message
 <DWR> ::= < Diameter Header: 280, REQ >
		   { Origin-Host }
		   { Origin-Realm }
		   [ Origin-State-Id ]
		 * [ AVP ] // ignored

Device-Watchdo-gAnswer message
 <DWA> ::= < Diameter Header: 280 >
		   { Result-Code }
		   { Origin-Host }
		   { Origin-Realm }
		   [ Error-Message ] // ignored
		   [ Failed-AVP ]    // ignored
		   [ Origin-State-Id ]
		 * [ AVP ]           // ignored
*/

type eventRcvDWR struct {
	m Message
}

func (eventRcvDWR) String() string {
	return "Rcv-DWR"
}

func (v eventRcvDWR) exec() error {
	RxReq++
	if state != open && state != locked {
		RejectReq++
		return notAcceptableEvent{e: v, s: state}
	}
	TraceMessage(v.m, Rx, nil)

	var oHost Identity
	var oRealm Identity
	var oState uint32
	var err error

	for rdr := bytes.NewReader(v.m.AVPs); rdr.Len() != 0; {
		a := AVP{}
		if err = a.wrapedUnmarshalFrom(rdr); err != nil {
			break
		}
		if a.VendorID != 0 {
			if a.Mandatory {
				err = InvalidAVP{Code: AvpUnsupported, AVP: a}
				break
			}
			continue
		}

		switch a.Code {
		case 264:
			if len(oHost) != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				oHost, err = GetOriginHost(a)
			}
		case 296:
			if len(oRealm) != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				oRealm, err = GetOriginRealm(a)
			}
		case 278:
			if oState != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				oState, err = getOriginStateID(a)
			}
		default:
			if a.Mandatory {
				err = InvalidAVP{Code: AvpUnsupported, AVP: a}
			}
		}

		if err != nil {
			break
		}
	}

	result := Success
	if v.m.FlgP || v.m.FlgT {
		result = InvalidHdrBits
		err = InvalidMessage{
			Code: result, ErrMsg: "DWR must not enable P and T flag"}
	} else if iavp, ok := err.(InvalidAVP); ok {
		result = iavp.Code
	} else if len(oHost) == 0 {
		result = MissingAvp
		err = InvalidAVP{Code: result, AVP: SetOriginHost("")}
	} else if len(oRealm) == 0 {
		result = MissingAvp
		err = InvalidAVP{Code: result, AVP: SetOriginRealm("")}
	} else if oHost != Peer.Host {
		result = UnknownPeer
		err = InvalidMessage{
			Code: result,
			ErrMsg: fmt.Sprintf(
				"peer host %s is not match with %s",
				oHost, Peer.Host)}
	} else if oRealm != Peer.Realm {
		result = UnknownPeer
		err = InvalidMessage{
			Code: result,
			ErrMsg: fmt.Sprintf(
				"peer realm %s is not match with %s",
				oRealm, Peer.Realm)}
	} else if oState != 0 && oState != Peer.state {
		result = UnknownPeer
		err = InvalidMessage{
			Code:   result,
			ErrMsg: fmt.Sprintf("peer state %d is not match with %d", oState, Peer.state)}
	}

	buf := new(bytes.Buffer)
	SetResultCode(result).MarshalTo(buf)
	SetOriginHost(Local.Host).MarshalTo(buf)
	SetOriginRealm(Local.Realm).MarshalTo(buf)
	if Local.state != 0 {
		setOriginStateID(Local.state).MarshalTo(buf)
	}
	if iavp, ok := err.(InvalidAVP); ok {
		setFailedAVP([]AVP{iavp.AVP}).MarshalTo(buf)
	}

	dwa := Message{
		FlgR: false, FlgP: false, FlgE: result != Success, FlgT: false,
		Code: 280, AppID: 0,
		HbHID: v.m.HbHID, EtEID: v.m.EtEID,
		AVPs: buf.Bytes()}

	if e := dwa.MarshalTo(conn); e != nil {
		TxAnsFail++
		conn.Close()
		err = e
	} else if err == nil && wdCount == 0 {
		CountTxCode(result)
		wdTimer.Stop()
		wdTimer.Reset(WDInterval)
	}

	TraceMessage(dwa, Tx, err)
	return err
}

// RcvDWA
type eventRcvDWA struct {
	m Message
}

func (eventRcvDWA) String() string {
	return "Rcv-DWA"
}

func (v eventRcvDWA) exec() error {
	// verify diameter header
	if v.m.FlgP {
		InvalidAns++
		return InvalidMessage{
			Code: InvalidHdrBits, ErrMsg: "DWA must not enable P flag"}
	}
	if state != open {
		InvalidAns++
		return notAcceptableEvent{e: v, s: state}
	}
	if _, ok := sndQueue[v.m.HbHID]; !ok {
		InvalidAns++
		return unknownAnswer(v.m.HbHID)
	}

	var result uint32
	var oHost Identity
	var oRealm Identity
	// var errorMsg string
	// var failedAVP []AVP
	var oState uint32
	var err error

	for rdr := bytes.NewReader(v.m.AVPs); rdr.Len() != 0; {
		a := AVP{}
		if err = a.wrapedUnmarshalFrom(rdr); err != nil {
			break
		}
		if a.VendorID != 0 {
			if a.Mandatory {
				err = InvalidAVP{Code: AvpUnsupported, AVP: a}
				break
			}
			continue
		}

		switch a.Code {
		case 268:
			if result != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				result, err = GetResultCode(a)
			}
		case 264:
			if len(oHost) != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				oHost, err = GetOriginHost(a)
			}
		case 296:
			if len(oRealm) != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				oRealm, err = GetOriginRealm(a)
			}
		case 278:
			if oState != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				oState, err = getOriginStateID(a)
			}
			// case 281:
			//	errorMsg, e = getErrorMessage(a)
		case 279:
			//failedAVP, e = getFailedAVP(a)
		default:
			if a.Mandatory {
				err = InvalidAVP{Code: AvpUnsupported, AVP: a}
			}
		}

		if err != nil {
			break
		}
	}

	if v.m.FlgE && result == Success {
		err = InvalidMessage{
			Code:   InvalidHdrBits,
			ErrMsg: "error flag is true but success response code"}
	} else if err != nil {
		// invalid AVP value
	} else if result == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetResultCode(0)}
	} else if len(oHost) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetOriginHost("")}
	} else if len(oRealm) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetOriginRealm("")}
	} else if oHost != Peer.Host && oHost != Local.Host {
		err = InvalidMessage{
			Code: UnknownPeer,
			ErrMsg: fmt.Sprintf(
				"peer host %s is not match with %s or %s",
				oHost, Peer.Host, Local.Host)}
	} else if oRealm != Peer.Realm && oRealm != Local.Realm {
		err = InvalidMessage{
			Code: UnknownPeer,
			ErrMsg: fmt.Sprintf(
				"peer realm %s is not match with %s or %s",
				oRealm, Peer.Realm, Local.Host)}
	} else if oState != 0 && oState != Peer.state {
		err = errors.New("peer may be abruptly restarted")
		Close(Rebooting)
	} else {
		if result == Success {
			wdCount = 0
		} else {
			err = FailureAnswer{Code: result}
		}
		delete(sndQueue, v.m.HbHID)
		wdTimer.Stop()
		wdTimer = time.AfterFunc(WDInterval, func() {
			notify <- eventWatchdog{}
		})
	}
	CountRxCode(result)

	TraceMessage(v.m, Rx, err)
	return err
}
