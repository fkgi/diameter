package diameter

import "fmt"

// UnknownAVPType is error of invalid AVP type
type UnknownAVPType struct{}

func (e UnknownAVPType) Error() string {
	return "unknow AVP data type"
}

// InvalidMessage is error of invalid message
type InvalidMessage uint32

func (e InvalidMessage) Error() string {
	switch uint32(e) {
	case DiameterUnsupportedVersion:
		return "unsupported verion"
	case DiameterInvalidHdrBits:
		return "invalid header bit"
	}
	return "invalid message"
}

// InvalidAVP is error of invalid AVP value
type InvalidAVP uint32

func (e InvalidAVP) Error() string {
	switch uint32(e) {
	case DiameterInvalidAvpBits:
		return "invalid AVP Bits"
	case DiameterInvalidAvpValue:
		return "invalid AVP Value"
	case DiameterMissingAvp:
		return "missing mandatory AVP"
	}
	return "invalid AVP"
}

// UnknownIDAnswer is error
type UnknownIDAnswer struct {
	RawMsg
}

func (e UnknownIDAnswer) Error() string {
	return "Unknown message recieved"
}

// FailureAnswer is error
type FailureAnswer struct {
	Answer
}

func (e FailureAnswer) Error() string {
	r := e.Answer.Result()
	if r > 10000 {
		return fmt.Sprintf("Answer message with failure: code=%d (vendor=%d)",
			r%10000, r/10000)
	}
	return fmt.Sprintf("Answer message with failure: code=%d", r)
}

// NotAcceptableEvent is error
type NotAcceptableEvent struct {
	stateEvent
	state
}

func (e NotAcceptableEvent) Error() string {
	return fmt.Sprintf("Event %s is not acceptable in state %s",
		e.stateEvent, e.state)
}

// WatchdogExpired is error
type WatchdogExpired struct{}

func (e WatchdogExpired) Error() string {
	return "watchdog is expired"
}

// ConnectionRefused is error
type ConnectionRefused struct{}

func (e ConnectionRefused) Error() string {
	return "connection is refused"
}
