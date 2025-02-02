package diameter

import "fmt"

// InvalidMessage is error of invalid message
type InvalidMessage struct {
	Code   uint32
	ErrMsg string
}

func (e InvalidMessage) Error() string {
	switch uint32(e.Code) {
	case UnknownSessionID:
		return "unknowns session ID, " + e.ErrMsg
	case UnsupportedVersion:
		return "unsupported verion, " + e.ErrMsg
	case InvalidHdrBits:
		return "invalid header bit, " + e.ErrMsg
	case ApplicationUnsupported:
		return "unsupported application, " + e.ErrMsg
	case UnknownPeer:
		return "unknown peer, " + e.ErrMsg
	case UnableToComply:
		return "unable to comply, " + e.ErrMsg
	}
	return fmt.Sprintf("generic invalid message (code=%d), %s", e.Code, e.ErrMsg)
}

// InvalidAVP is error of invalid AVP value
type InvalidAVP struct {
	Code uint32
	AVP
	E error
}

func (e InvalidAVP) Error() string {
	err := ""
	if e.E != nil {
		err = ", error=" + e.E.Error()
	}
	switch e.Code {
	case InvalidAvpBits:
		return fmt.Sprintf("invalid AVP Bits (code=%d)%v",
			e.AVP.Code, err)
	case InvalidAvpValue:
		return fmt.Sprintf("invalid AVP Value (code=%d, value=%v)%v",
			e.AVP.Code, e.AVP.Data, err)
	case MissingAvp:
		return fmt.Sprintf("missing mandatory AVP (code=%d)%v",
			e.AVP.Code, err)
	case AvpOccursTooManyTimes:
		return fmt.Sprintf("AVP occures too many time (code=%d)%v",
			e.AVP.Code, err)
	case AvpUnsupported:
		return fmt.Sprintf("unsupported AVP with mandatory flag (code=%d)%v",
			e.AVP.Code, err)
	}
	return fmt.Sprintf("generic invalid AVP (code=%d)%v", e.AVP.Code, err)
}

// FailureAnswer is error response from peer.
type FailureAnswer struct {
	Code   uint32
	VenID  uint32
	ErrMsg string
	Avps   []AVP
}

func (e FailureAnswer) Error() string {
	acodes := ""
	for _, a := range e.Avps {
		acodes = fmt.Sprintf("%s/%d", acodes, a.Code)
	}
	return fmt.Sprintf("error answer: code=%d (vendor=%d), message=%s, avp=%s",
		e.Code, e.VenID, e.ErrMsg, acodes)
}

type RejectRxMessage struct {
	State  conState
	ErrMsg string
}

func (err RejectRxMessage) Error() string {
	return fmt.Sprintf("Rx message is rejected in state %s: %s",
		err.State, err.ErrMsg)
}

type TransportTxError struct {
	err error
}

func (err TransportTxError) Error() string {
	return fmt.Sprintf("failed to send data on lower layer: %s", err.err)
}

type notAcceptableEvent struct {
	e stateEvent
	s conState
}

func (err notAcceptableEvent) Error() string {
	return fmt.Sprintf("event %s is not acceptable in state %s", err.e, err.s)
}
