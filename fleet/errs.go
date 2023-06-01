package fleet

import "errors"

var (
	ErrDeviceNotFound              = errors.New("requested device was not found")
	ErrDeviceConfigurationNotFound = errors.New("requested device configuration was not found")
	ErrNotImplemented              = errors.New("not implemented")
	ErrHostnameRequired            = errors.New("hostname is required")
	ErrDeviceAlreadyExists         = errors.New("device already exists")
	ErrDeviceConfigurationInvalid  = errors.New("device configuration is invalid")
	ErrDeviceInvalid               = errors.New("device is invalid")
)
