package diameter

import "fmt"
import "strconv"

// UnknownAVPType is error of invalid AVP type
type UnknownAVPType struct {
}

func (e UnknownAVPType) Error() string {
	return "unknow AVP data type"
}

// UnknownApplicationID is error of invalid AVP type
type UnknownApplicationID struct {
}

func (e UnknownApplicationID) Error() string {
	return "unknow application id"
}

// UnknownCommand is error of invalid AVP type
type UnknownCommand struct {
}

func (e UnknownCommand) Error() string {
	return "unknow command code"
}

// InvalidMessage is error of invalid message
type InvalidMessage struct {
}

func (e InvalidMessage) Error() string {
	return "invalid message data"
}

// NoMandatoryAVP is error of invalid avp
type NoMandatoryAVP struct {
}

func (e NoMandatoryAVP) Error() string {
	return "mandatory AVP not found"
}

// InvalidAVP is error of invalid AVP value
type InvalidAVP uint32

func (e InvalidAVP) Error() string {
	switch uint32(e) {
	case DiameterInvalidAvpBits:
		return "invalid AVP Bits"
	case DiameterInvalidAvpValue:
		return "invalid AVP Value"
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
	return "Answer message with failure: code=" + strconv.FormatInt(int64(e.Answer.Result()), 10)
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
