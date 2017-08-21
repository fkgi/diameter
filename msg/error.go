package msg

// UnknownAVPTypeError is error of invalid AVP type
type UnknownAVPTypeError struct {
}

func (e UnknownAVPTypeError) Error() string {
	return "unknow AVP data type"
}

// InvalidMessageError is error of invalid message
type InvalidMessageError struct {
}

func (e InvalidMessageError) Error() string {
	return "invalid message data"
}

// NoMandatoryAVPError is error of invalid avp
type NoMandatoryAVPError struct {
}

func (e NoMandatoryAVPError) Error() string {
	return "mandatory AVP not found"
}

// InvalidAVPError is error of invalid AVP value
type InvalidAVPError struct {
}

func (e InvalidAVPError) Error() string {
	return "invalid AVP data"
}
