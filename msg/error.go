package msg

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
type InvalidAVP struct {
}

func (e InvalidAVP) Error() string {
	return "invalid AVP data"
}
