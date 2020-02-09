package api

import (
	"net/http"

	"github.com/opentelco/go-swpx/errors"
)

// Response is the main struct that returns to the client
type Response struct {
	Status *ResponseStatus `json:"status"`
	Data   interface{}     `json:"data"`
}

// ResponseStatus is the status of a Response.
type ResponseStatus struct {
	// AppErrorCode is used for internal
	AppErrorCode    errors.ErrorCode `json:"-"`
	AppErrorMessage string           `json:"-"`

	// Sent to the client
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message"`
}

// Render implements the chi Response return
func (rs *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// NewResponse creates and returnes a Response.
// if no response is passed as argument it will create a OK.
func NewResponse(status *ResponseStatus, payload interface{}) *Response {
	if payload != nil {
		if err, ok := payload.(errors.Error); ok {
			status.AppErrorMessage = err.Error()
			status.AppErrorCode = err.Code
			payload = nil
			status.Message = err.Error()
		}
	}

	return &Response{Status: status, Data: payload}
}

var (
	// ResponseStatusOK is  http.StatusOK
	ResponseStatusOK = &ResponseStatus{
		Error:   false,
		Code:    http.StatusOK,
		Type:    "success",
		Message: "success",
	}

	ResponseStatusNotFound = &ResponseStatus{
		Error:   true,
		Code:    http.StatusNotFound,
		Type:    "failure",
		Message: "could not find anything",
	}
)
