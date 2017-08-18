package msg

// UnknownAVPTypeError is error of invalid AVP type
type UnknownAVPTypeError struct {
}

func (e *UnknownAVPTypeError) Error() string {
	return "unknow AVP data type"
}
