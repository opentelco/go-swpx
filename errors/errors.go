package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type ErrorCode int

// Wraps the Wrap function..
func Wrap(err error, msg string) error {
	return errors.Wrap(err, msg)
}

// Error is a local error type
type Error struct {
	Message    string
	Code       ErrorCode
	StatusCode uint
}

// Error implements the Error interface
func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// New creates a new error with a code
func New(msg string, code ErrorCode) Error {
	return Error{
		Message: msg,
		Code:    code,
	}
}
