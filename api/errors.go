package api

import (
	"net/http"

	"git.liero.se/opentelco/go-swpx/errors"
)

// ErrorStatusInvalidAddr is the status response when the app cannot do anything about the host
// the error can be traced back to bad dns settings
var ErrorStatusInvalidAddr = &ResponseStatus{
	AppErrorCode:    errors.ErrInvalidAddr,
	AppErrorMessage: "",

	Error:   true,
	Code:    http.StatusNotAcceptable,
	Type:    "failed",
	Message: "invalid or mistyped hostname/addr",
}

var ErrorStatusInvalidProvider = &ResponseStatus{
	AppErrorCode:    errors.ErrInvalidProvider,
	AppErrorMessage: "",

	Error:   true,
	Code:    http.StatusNotAcceptable,
	Type:    "failed",
	Message: "the selected provider does not exist",
}

var ErrorStatusRequestTimeout = &ResponseStatus{
	AppErrorCode:    errors.ErrTimeoutRequest,
	AppErrorMessage: "",

	Error:   true,
	Code:    http.StatusRequestTimeout,
	Type:    "failed",
	Message: "the request did not make it back in time",
}
