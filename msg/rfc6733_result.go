package msg

const (
	// DiameterMultiRoundAuth is Result-Code 1001
	DiameterMultiRoundAuth ResultCode = 1001

	// DiameterSuccess is Result-Code 2001
	DiameterSuccess ResultCode = 2001
	// DiameterLimitedSuccess is Result-Code 2002
	DiameterLimitedSuccess ResultCode = 2002

	// DiameterCommandUnspported is Result-Code 3001
	DiameterCommandUnspported ResultCode = 3001
	// DiameterUnableToDeliver is Result-Code 3002
	DiameterUnableToDeliver ResultCode = 3002
	// DiameterRealmNotServed is Result-Code 3003
	DiameterRealmNotServed ResultCode = 3003
	// DiameterTooBusy is Result-Code 3004
	DiameterTooBusy ResultCode = 3004
	// DiameterLoopDetected is Result-Code 3005
	DiameterLoopDetected ResultCode = 3005
	// DiameterRedirectIndication is Result-Code 3006
	DiameterRedirectIndication ResultCode = 3006
	// DiameterApplicationUnsupported is Result-Code 3007
	DiameterApplicationUnsupported ResultCode = 3007
	// DiameterInvalidHdrBits is Result-Code 3008
	DiameterInvalidHdrBits ResultCode = 3008
	// DiameterInvalidAvpBits is Result-Code 3009
	DiameterInvalidAvpBits ResultCode = 3009
	// DiameterUnknownPeer is Result-Code 3010
	DiameterUnknownPeer ResultCode = 3010

	// DiameterAuthenticationRejected is Result-Code 4001
	DiameterAuthenticationRejected ResultCode = 4001
	// DiameterOutOfSpace is Result-Code 4002
	DiameterOutOfSpace ResultCode = 4002
	// DiameterElectionLost is Result-Code 4003
	DiameterElectionLost ResultCode = 4003

	// DiameterAvpUnsupported is Result-Code 5001
	DiameterAvpUnsupported ResultCode = 5001
	// DiameterUnknownSessionID is Result-Code 5002
	DiameterUnknownSessionID ResultCode = 5002
	// DiameterAuthorizationRejected is Result-Code 5003
	DiameterAuthorizationRejected ResultCode = 5003
	// DiameterInvalidAvpValue is Result-Code 5004
	DiameterInvalidAvpValue ResultCode = 5004
	// DiameterMissingAvp is Result-Code 5005
	DiameterMissingAvp ResultCode = 5005
	// DiameterResourcesExceeded is Result-Code 5006
	DiameterResourcesExceeded ResultCode = 5006
	// DiameterContradictingAvps is Result-Code 5007
	DiameterContradictingAvps ResultCode = 5007
	//DiameterAvpNotAllowed is Result-Code 5008
	DiameterAvpNotAllowed ResultCode = 5008
	// DiameterAvpOccursTooManyTimes is Result-Code 5009
	DiameterAvpOccursTooManyTimes ResultCode = 5009
	// DiameterNoCommonApplication is Result-Code 5010
	DiameterNoCommonApplication ResultCode = 5010
	// DiameterUnsupportedVersion is Result-Code 5011
	DiameterUnsupportedVersion ResultCode = 5011
	// DiameterUnableToComply is Result-Code 5012
	DiameterUnableToComply ResultCode = 5012
	// DiameterInvalidBitInHeader is Result-Code 5013
	DiameterInvalidBitInHeader ResultCode = 5013
	// DiameterInvalidAvpLength is Result-Code 5014
	DiameterInvalidAvpLength ResultCode = 5014
	// DiameterInvalidMessageLength is Result-Code 5015
	DiameterInvalidMessageLength ResultCode = 5015
	// DiameterInvalidAvpBitCombo is Result-Code 5016
	DiameterInvalidAvpBitCombo ResultCode = 5016
	// DiameterNoCommonSecurity is Result-Code 5017
	DiameterNoCommonSecurity ResultCode = 5017
)
