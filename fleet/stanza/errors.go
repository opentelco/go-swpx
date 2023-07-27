package stanza

import "errors"

var (
	ErrTypeFailedCreate          = "FAILED_CREATE_STANZA"
	ErrTypeInvalidTemplate       = "INVALID_TEMPLATE"
	ErrTypeInvalidRevertTemplate = "INVALID_REVERT_TEMPLATE"
)

var (
	ErrStanzaNotFound  = errors.New("stanza not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotImplemented  = errors.New("not implemented")
)
