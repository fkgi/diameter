package diameter

const (
	// MultiRoundAuth is Result-Code 1001
	MultiRoundAuth uint32 = 1001

	// Success is Result-Code 2001
	Success uint32 = 2001
	// LimitedSuccess is Result-Code 2002
	LimitedSuccess uint32 = 2002

	// CommandUnspported is Result-Code 3001
	CommandUnspported uint32 = 3001
	// UnableToDeliver is Result-Code 3002
	UnableToDeliver uint32 = 3002
	// RealmNotServed is Result-Code 3003
	RealmNotServed uint32 = 3003
	// TooBusy is Result-Code 3004
	TooBusy uint32 = 3004
	// LoopDetected is Result-Code 3005
	LoopDetected uint32 = 3005
	// RedirectIndication is Result-Code 3006
	RedirectIndication uint32 = 3006
	// ApplicationUnsupported is Result-Code 3007
	ApplicationUnsupported uint32 = 3007
	// InvalidHdrBits is Result-Code 3008
	InvalidHdrBits uint32 = 3008
	// InvalidAvpBits is Result-Code 3009
	InvalidAvpBits uint32 = 3009
	// UnknownPeer is Result-Code 3010
	UnknownPeer uint32 = 3010

	// AuthenticationRejected is Result-Code 4001
	AuthenticationRejected uint32 = 4001
	// OutOfSpace is Result-Code 4002
	OutOfSpace uint32 = 4002
	// ElectionLost is Result-Code 4003
	ElectionLost uint32 = 4003

	// AvpUnsupported is Result-Code 5001
	AvpUnsupported uint32 = 5001
	// UnknownSessionID is Result-Code 5002
	UnknownSessionID uint32 = 5002
	// AuthorizationRejected is Result-Code 5003
	AuthorizationRejected uint32 = 5003
	// InvalidAvpValue is Result-Code 5004
	InvalidAvpValue uint32 = 5004
	// MissingAvp is Result-Code 5005
	MissingAvp uint32 = 5005
	// ResourcesExceeded is Result-Code 5006
	ResourcesExceeded uint32 = 5006
	// ContradictingAvps is Result-Code 5007
	ContradictingAvps uint32 = 5007
	// AvpNotAllowed is Result-Code 5008
	AvpNotAllowed uint32 = 5008
	// AvpOccursTooManyTimes is Result-Code 5009
	AvpOccursTooManyTimes uint32 = 5009
	// NoCommonApplication is Result-Code 5010
	NoCommonApplication uint32 = 5010
	// UnsupportedVersion is Result-Code 5011
	UnsupportedVersion uint32 = 5011
	// UnableToComply is Result-Code 5012
	UnableToComply uint32 = 5012
	// InvalidBitInHeader is Result-Code 5013
	InvalidBitInHeader uint32 = 5013
	// InvalidAvpLength is Result-Code 5014
	InvalidAvpLength uint32 = 5014
	// InvalidMessageLength is Result-Code 5015
	InvalidMessageLength uint32 = 5015
	// InvalidAvpBitCombo is Result-Code 5016
	InvalidAvpBitCombo uint32 = 5016
	// NoCommonSecurity is Result-Code 5017
	NoCommonSecurity uint32 = 5017
)
