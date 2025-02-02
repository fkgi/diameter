package diameter

import (
	"bytes"
	"fmt"
	"time"
)

/*
Disconnect-Peer-Request message
 <DPR> ::= < Diameter Header: 282, REQ >
		   { Origin-Host }
		   { Origin-Realm }
		   { Disconnect-Cause }
		 * [ AVP ] // no any other AVP

Disconnect-Peer-Answer message
 <DPA> ::= < Diameter Header: 282 >
		   { Result-Code }
		   { Origin-Host }
		   { Origin-Realm }
		   [ Error-Message ] // ignored
		   [ Failed-AVP ]    // ignored
		 * [ AVP ]           // ignored
*/

type eventRcvDPR struct {
	m Message
}

func (eventRcvDPR) String() string {
	return "Rcv-DPR"
}

func (v eventRcvDPR) exec(c *Connection) error {
	var err error
	if c.state != open && c.state != locked {
		err = RejectRxMessage{
			State: c.state, ErrMsg: "DPR is not acceptable"}
	}
	if TraceMessage != nil {
		TraceMessage(v.m, Rx, err)
	}
	if err != nil {
		return err
	}

	// Notify(PurgeEvent{tx: false, req: true, conn: c, Err: e})
	var oHost Identity
	var oRealm Identity
	cause := Enumerated(-1)

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
		err = InvalidMessage{
			Code: result, ErrMsg: "DPR must not enable P and T flag"}
	} else if iavp, ok := err.(InvalidAVP); ok {
		result = iavp.Code
	} else if len(oHost) == 0 {
		result = MissingAvp
		err = InvalidAVP{Code: result, AVP: SetOriginHost("")}
	} else if len(oRealm) == 0 {
		result = MissingAvp
		err = InvalidAVP{Code: result, AVP: SetOriginRealm("")}
	} else if c.Host != "" && oHost != c.Host {
		result = UnknownPeer
		err = InvalidMessage{
			Code: result,
			ErrMsg: fmt.Sprintf(
				"peer host %s is not match with %s",
				oHost, c.Host)}
	} else if c.Realm != "" && oRealm != c.Realm {
		result = UnknownPeer
		err = InvalidMessage{
			Code: result,
			ErrMsg: fmt.Sprintf(
				"peer realm %s is not match with %s",
				oRealm, c.Realm)}
	} else if cause < 0 {
		result = MissingAvp
		err = InvalidAVP{Code: result, AVP: setDisconnectCause(Rebooting)}
	}

	buf := new(bytes.Buffer)
	SetResultCode(result).MarshalTo(buf)
	SetOriginHost(Host).MarshalTo(buf)
	SetOriginRealm(Realm).MarshalTo(buf)
	if iavp, ok := err.(InvalidAVP); ok {
		setFailedAVP([]AVP{iavp.AVP}).MarshalTo(buf)
	}

	dpa := Message{
		FlgR: false, FlgP: false, FlgE: result != Success, FlgT: false,
		Code: 282, AppID: 0,
		HbHID: v.m.HbHID, EtEID: v.m.EtEID,
		AVPs: buf.Bytes()}

	if e := dpa.MarshalTo(c.conn); e != nil {
		c.conn.Close()
		err = TransportTxError{err: e}
	} else if err == nil {
		c.state = closing
		c.wdTimer.Stop()
		c.wdTimer = time.AfterFunc(WDInterval, func() {
			c.conn.Close()
		})
	}

	if TraceMessage != nil {
		TraceMessage(dpa, Tx, err)
	}
	return err
}

type eventRcvDPA struct {
	m Message
}

func (eventRcvDPA) String() string {
	return "Rcv-DPA"
}

func (v eventRcvDPA) exec(c *Connection) error {
	var err error

	if v.m.FlgP {
		err = InvalidMessage{
			Code: InvalidHdrBits, ErrMsg: "DPA must not enable P flag"}
	} else if c.state != closing {
		err = InvalidMessage{
			Code:   UnableToComply,
			ErrMsg: "DPA is not acceptable in " + c.state.String() + " state"}
	} else if _, ok := c.sndQueue[v.m.HbHID]; !ok {
		err = InvalidMessage{
			Code:   UnableToComply,
			ErrMsg: "correlated request with the Hop-by-Hop ID not found"}
	}

	if err != nil {
		if TraceMessage != nil {
			TraceMessage(v.m, Rx, err)
		}
		return err
	}

	var result uint32
	var oHost Identity
	var oRealm Identity
	// var errorMsg string
	// var failedAVP []AVP

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
		err = InvalidMessage{
			Code:   InvalidHdrBits,
			ErrMsg: "error flag is true but success response code"}
	} else if result == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetResultCode(0)}
	} else if len(oHost) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetOriginHost("")}
	} else if len(oRealm) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetOriginRealm("")}
	} else if oHost != c.Host && oHost != Host {
		err = InvalidMessage{
			Code: UnknownPeer,
			ErrMsg: fmt.Sprintf(
				"peer host %s is not match with %s or %s",
				oHost, c.Host, Host)}
	} else if oRealm != c.Realm && oRealm != Realm {
		err = InvalidMessage{
			Code: UnknownPeer,
			ErrMsg: fmt.Sprintf(
				"peer realm %s is not match with %s or %s",
				oRealm, c.Realm, Host)}
	} else if result != Success {
		err = FailureAnswer{Code: result}
		delete(c.sndQueue, v.m.HbHID)
	} else if err != nil {
		// invalid AVP value
	} else {
		delete(c.sndQueue, v.m.HbHID)
		c.wdTimer.Stop()
		err = c.conn.Close()
	}

	if TraceMessage != nil {
		TraceMessage(v.m, Rx, err)
	}
	return err
}
