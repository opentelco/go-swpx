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

func ErrInvalidArgument(message string) *ResponseStatus {
	return &ResponseStatus{
		AppErrorCode:    core.ErrInvalidArgument,
		AppErrorMessage: message,

		Error:   true,
		Code:    http.StatusBadRequest,
		Type:    "invalid argument",
		Message: message,
	}
}

// ErrorStatusInvalidAddr is the status response when the app cannot do anything about the host
// the error can be traced back to bad dns settings
var ErrorStatusInvalidAddr = &ResponseStatus{
	AppErrorCode:    core.ErrInvalidAddr,
	AppErrorMessage: "",

	Error:   true,
	Code:    http.StatusBadRequest,
	Type:    "failed",
	Message: "invalid or mistyped hostname/addr",
}

var ErrorStatusInvalidProvider = &ResponseStatus{
	AppErrorCode:    core.ErrInvalidProvider,
	AppErrorMessage: "",

	Error:   true,
	Code:    http.StatusBadRequest,
	Type:    "failed",
	Message: "the selected provider does not exist",
}

var ErrorStatusRequestTimeout = &ResponseStatus{
	AppErrorCode:    core.ErrTimeoutRequest,
	AppErrorMessage: "",

	Error:   true,
	Code:    http.StatusRequestTimeout,
	Type:    "failed",
	Message: "the request did not make it back in time",
}
