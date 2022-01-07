package diameter

import (
	"bytes"
	"time"
)

/*
Disconnect-Peer-Request message
 <DPR>  ::= < Diameter Header: 282, REQ >
			{ Origin-Host }
			{ Origin-Realm }
			{ Disconnect-Cause }
		  * [ AVP ]

Disconnect-Peer-Answer message
 <DPA>  ::= < Diameter Header: 282 >
			{ Result-Code }
			{ Origin-Host }
			{ Origin-Realm }
			[ Error-Message ]
			[ Failed-AVP ]
		  * [ AVP ]
*/
type eventRcvDPR struct {
	m Message
}

func (eventRcvDPR) String() string {
	return "Rcv-DPR"
}

func (v eventRcvDPR) exec() error {
	RxReq++
	if state != open && state != locked {
		RejectReq++
		return notAcceptableEvent{e: v, s: state}
	}

	TraceMessage(v.m, Rx, nil)

	// Notify(PurgeEvent{tx: false, req: true, conn: c, Err: e})
	var oHost Identity
	var oRealm Identity
	cause := Enumerated(-1)
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
		case 273:
			if cause != -1 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				cause, err = getDisconnectCause(a)
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
		err = InvalidMessage(result)
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
		err = InvalidMessage(result)
	} else if oRealm != Peer.Realm {
		result = UnknownPeer
		err = InvalidMessage(result)
	} else if cause < 0 {
		result = MissingAvp
		err = InvalidAVP{Code: result, AVP: setDisconnectCause(Rebooting)}
	}

	buf := new(bytes.Buffer)
	SetResultCode(result).MarshalTo(buf)
	SetOriginHost(Local.Host).MarshalTo(buf)
	SetOriginRealm(Local.Realm).MarshalTo(buf)
	if iavp, ok := err.(InvalidAVP); ok {
		setFailedAVP([]AVP{iavp.AVP}).MarshalTo(buf)
	}

	dpa := Message{
		FlgR: false, FlgP: false, FlgE: result != Success, FlgT: false,
		Code: 282, AppID: 0,
		HbHID: v.m.HbHID, EtEID: v.m.EtEID,
		AVPs: buf.Bytes()}

	conn.SetWriteDeadline(time.Now().Add(TxTimeout))
	if e := dpa.MarshalTo(conn); e != nil {
		TxAnsFail++
		conn.Close()
		err = e
	} else if err == nil {
		countTxCode(result)
		state = closing
		wdTimer.Stop()
		wdTimer = time.AfterFunc(WDInterval, func() {
			conn.Close()
		})
	}

	TraceMessage(dpa, Tx, err)
	return err
}

type eventRcvDPA struct {
	m Message
}

func (eventRcvDPA) String() string {
	return "Rcv-DPA"
}

func (v eventRcvDPA) exec() error {
	// verify diameter header
	if v.m.FlgP {
		InvalidAns++
		return InvalidMessage(InvalidHdrBits)
	}
	if state != closing {
		InvalidAns++
		return notAcceptableEvent{e: v, s: state}
	}
	if _, ok := sndStack[v.m.HbHID]; !ok {
		InvalidAns++
		return unknownAnswer(v.m.HbHID)
	}

	var result uint32
	var oHost Identity
	var oRealm Identity
	// var errorMsg string
	// var failedAVP []AVP
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
		err = InvalidMessage(InvalidHdrBits)
	} else if err != nil {
		// invalid AVP value
	} else if result == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetResultCode(0)}
	} else if len(oHost) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetOriginHost("")}
	} else if len(oRealm) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetOriginRealm("")}
	} else if oHost != Peer.Host && oHost != Local.Host {
		err = InvalidMessage(UnknownPeer)
	} else if oRealm != Peer.Realm && oRealm != Local.Realm {
		err = InvalidMessage(UnknownPeer)
	} else if result != Success {
		err = FailureAnswer{Code: result}
		delete(sndStack, v.m.HbHID)
	} else {
		delete(sndStack, v.m.HbHID)
		wdTimer.Stop()
		conn.Close()
	}
	countRxCode(result)

	TraceMessage(v.m, Rx, err)
	return err
}
