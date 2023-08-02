package notification

import "errors"

var (
	ErrInvalidArgument      = errors.New("invalid argument")
	ErrNotificationNotFound = errors.New("notification not found")
	ErrNotImplemented       = errors.New("not implemented")

	ErrTypeFailedCreate = "FAILED_CREATE_NOTIFICATION"
)
