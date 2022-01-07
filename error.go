package diameter

import "fmt"

// InvalidMessage is error of invalid message
type InvalidMessage uint32

func (e InvalidMessage) Error() string {
	switch uint32(e) {
	case UnknownSessionID:
		return "unknowns session ID"
	case UnsupportedVersion:
		return "unsupported verion"
	case InvalidHdrBits:
		return "invalid header bit"
	case ApplicationUnsupported:
		return "unsupported application"
	case UnknownPeer:
		return "unknown peer"
	}
	return fmt.Sprintf("generic invalid message (code=%d)", e)
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
	}
	return "generic invalid AVP"
}

// FailureAnswer is error
type FailureAnswer struct {
	Code  uint32
	VenID uint32
}

func (e FailureAnswer) Error() string {
	if e.VenID != 0 {
		return fmt.Sprintf("Answer message with failure: code=%d (vendor=%d)",
			e.Code, e.VenID)
	}
	return fmt.Sprintf("Answer message with failure: code=%d", e.Code)
}

// ConnectionRefused is error
type ConnectionRefused struct{}

func (e ConnectionRefused) Error() string {
	return "connection is refused"
}

type notAcceptableEvent struct {
	e stateEvent
	s conState
}

func (err notAcceptableEvent) Error() string {
	return fmt.Sprintf("event %s is not acceptable in state %s", err.e, err.s)
}

type unknownAnswer uint32

func (e unknownAnswer) Error() string {
	return fmt.Sprintf("Unknown hop-by-hop ID %d answer", e)
}
