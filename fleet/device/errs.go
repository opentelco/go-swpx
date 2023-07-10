package device

import "errors"

var (
	ErrDeviceNotFound                          = errors.New("requested device was not found")
	ErrDeviceConfigurationNotFound             = errors.New("requested device configuration was not found")
	ErrNotImplemented                          = errors.New("not implemented")
	ErrHostnameRequired                        = errors.New("hostname is required")
	ErrDeviceAlreadyExists                     = errors.New("device already exists")
	ErrDeviceConfigurationInvalid              = errors.New("device configuration is invalid")
	ErrDeviceInvalid                           = errors.New("device is invalid")
	ErrHostnameOrManagementIpRequired          = errors.New("hostname or management ip is required")
	ErrInvalidArgumentScheduleNotSet           = errors.New("invalid argument: schedule not set")
	ErrInvalidArgumentScheduleTypeNotSet       = errors.New("invalid argument: schedule type not set")
	ErrInvalidArgumentScheduleIntervalTooShort = errors.New("invalid argument: schedule interval too short")
)

var (
	ErrTypeDeviceNotFound = "DEVICE_NOT_FOUND"
)
