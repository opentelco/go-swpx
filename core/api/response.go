/*
 * Copyright (c) 2020. Liero AB
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software
 * is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
 * IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
 * CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package api

import (
	"net/http"

	"git.liero.se/opentelco/go-swpx/core"
)

// Response is the main struct that returns to the client
type Response struct {
	Status *ResponseStatus `json:"status" bson:"status"`
	Data   interface{}     `json:"data" bson:"data"`
}

// ResponseStatus is the status of a Response.
type ResponseStatus struct {
	// AppErrorCode is used for internal
	AppErrorCode    core.ErrorCode `json:"-" bson:"-"`
	AppErrorMessage string         `json:"-" bson:"-"`

	// Sent to the client
	Error   bool   `json:"error" bson:"error"`
	Code    int    `json:"code" bson:"code"`
	Type    string `json:"type,omitempty" bson:"type,omitempty"`
	Message string `json:"message" bson:"message"`
}

// Render implements the chi Response return
func (rs *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// NewResponse creates and returnes a Response.
// if no response is passed as argument it will create a OK.
func NewResponse(status *ResponseStatus, payload interface{}) *Response {
	switch err := payload.(type) {
	case core.Error:
		status.AppErrorMessage = err.Error()
		status.AppErrorCode = err.Code
		payload = nil
		status.Message = err.Error()
		return &Response{Status: status, Data: payload}
	case error:
		return &Response{Status: &ResponseStatus{
			Error:   true,
			Code:    http.StatusInternalServerError,
			Type:    "internal-error",
			Message: err.Error(),
		}, Data: nil}

	default:
		return &Response{Status: status, Data: payload}

	}
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
	ResponseStatusError = &ResponseStatus{
		Error:   true,
		Code:    http.StatusInternalServerError,
		Type:    "failure",
		Message: "internal server error",
	}

	ResponseStatusNothingFound = &ResponseStatus{
		Error:   true,
		Code:    http.StatusNoContent,
		Type:    "error",
		Message: "could not get any data from the poller",
	}
)
