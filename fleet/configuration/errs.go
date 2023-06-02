package configuration

import "errors"

var (
	ErrNotImplemented        = errors.New("not implemented")
	ErrInvalidArgumentID     = errors.New("invalid argument: ID")
	ErrConfigurationNotFound = errors.New("requested device configuration was not found")
	ErrConfigurationInvalid  = errors.New("device configuration is invalid")
)
