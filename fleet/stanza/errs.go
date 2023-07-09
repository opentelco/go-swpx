package stanza

import "errors"

var (
	ErrInvalidArgument      = errors.New("invalid argument")
	ErrNotificationNotFound = errors.New("stanza not found")
	ErrNotImplemented       = errors.New("not implemented")
)
