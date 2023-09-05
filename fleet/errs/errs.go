package errs

import "errors"

type ErrorCode uint

const (
	ErrCodeNotFound ErrorCode = iota
	ErrCodeConflict
)

func New(msg string, code ErrorCode) error {
	return errors.New(msg)
}
