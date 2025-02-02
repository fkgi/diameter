package diameter

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

/*
Device-Watchdog-Request message
 <DWR> ::= < Diameter Header: 280, REQ >
		   { Origin-Host }
		   { Origin-Realm }
		   [ Origin-State-Id ]
		 * [ AVP ] // ignored

Device-Watchdog-Answer message
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

func (v eventRcvDWR) exec(c *Connection) error {
	var err error
	if c.state != open && c.state != locked {
		err = RejectRxMessage{
			State: c.state, ErrMsg: "DWR is not acceptable"}
	}
	if TraceMessage != nil {
		TraceMessage(v.m, Rx, err)
	}
	if err != nil {
		return err
	}

	var oHost Identity
	var oRealm Identity
	var oState uint32

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
	} else if oHost != c.Host {
		result = UnknownPeer
		err = InvalidMessage{
			Code: result,
			ErrMsg: fmt.Sprintf(
				"peer host %s is not match with %s",
				oHost, c.Host)}
	} else if oRealm != c.Realm {
		result = UnknownPeer
		err = InvalidMessage{
			Code: result,
			ErrMsg: fmt.Sprintf(
				"peer realm %s is not match with %s",
				oRealm, c.Realm)}
	} else if oState != 0 && oState != c.stateID {
		result = UnknownPeer
		err = InvalidMessage{
			Code: result,
			ErrMsg: fmt.Sprintf("peer state %d is not match with %d",
				oState, c.stateID)}
	}

	buf := new(bytes.Buffer)
	SetResultCode(result).MarshalTo(buf)
	SetOriginHost(Host).MarshalTo(buf)
	SetOriginRealm(Realm).MarshalTo(buf)
	if stateID != 0 {
		setOriginStateID(stateID).MarshalTo(buf)
	}
	if iavp, ok := err.(InvalidAVP); ok {
		setFailedAVP([]AVP{iavp.AVP}).MarshalTo(buf)
	}

	dwa := Message{
		FlgR: false, FlgP: false, FlgE: result != Success, FlgT: false,
		Code: 280, AppID: 0,
		HbHID: v.m.HbHID, EtEID: v.m.EtEID,
		AVPs: buf.Bytes()}

	if e := dwa.MarshalTo(c.conn); e != nil {
		c.conn.Close()
		err = TransportTxError{err: e}
	} else if err == nil && c.wdCount == 0 {
		c.wdTimer.Stop()
		c.wdTimer.Reset(WDInterval)
	}

	if TraceMessage != nil {
		TraceMessage(dwa, Tx, err)
	}
	return err
}

// RcvDWA
type eventRcvDWA struct {
	m Message
}

func (eventRcvDWA) String() string {
	return "Rcv-DWA"
}

func (v eventRcvDWA) exec(c *Connection) error {
	var err error

	if v.m.FlgP {
		err = InvalidMessage{
			Code: InvalidHdrBits, ErrMsg: "DWA must not enable P flag"}
	} else if c.state != open {
		err = InvalidMessage{
			Code:   UnableToComply,
			ErrMsg: "DWA is not acceptable in " + c.state.String() + " state"}
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
	var oState uint32

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
		case 281:
			//	errorMsg, e = getErrorMessage(a)
		case 279:
			// failedAVP, e = getFailedAVP(a)
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
	} else if oState != 0 && oState != c.stateID {
		err = errors.New("peer may be abruptly restarted")
		c.Close(Rebooting)
	} else if result == Success && err != nil {
		// invalid AVP value
	} else {
		if result == Success {
			c.wdCount = 0
		} else {
			err = FailureAnswer{Code: result}
		}
		delete(c.sndQueue, v.m.HbHID)
		c.wdTimer.Stop()
		c.wdTimer = time.AfterFunc(WDInterval, func() {
			c.notify <- eventWatchdog{}
		})
	}

	if TraceMessage != nil {
		TraceMessage(v.m, Rx, err)
	}
	return err
}
