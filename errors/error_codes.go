package errors

const (
	// ErrInvalidAddr if the Request has a invalid hostname
	ErrInvalidAddr     ErrorCode = 1010
	ErrInvalidProvider           = 1011
	ErrInvalidResource           = 1012

	// ErrInvalidPort is returned if the request has a invalid port
	ErrInvalidPort    = 1020
	ErrTimeoutRequest = 2010

	// ErrDickPick is a ha.ha. error..
	ErrDickPick = 9010
)
