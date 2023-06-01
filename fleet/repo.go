package fleet

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleetpb"
)

type Repository interface {
	// Device
	GetDeviceByID(ctx context.Context, id string) (*fleetpb.Device, error)
	ListDevices(ctx context.Context, parms *fleetpb.ListDevicesParameters) ([]*fleetpb.Device, error)
	UpsertDevice(ctx context.Context, dev *fleetpb.Device) (*fleetpb.Device, error)
	DeleteDevice(ctx context.Context, id string) error

	// DeviceConfiguration
	GetDeviceConfigurationByID(ctx context.Context, id string) (*fleetpb.DeviceConfiguration, error)
	ListDeviceConfiguration(ctx context.Context, params *fleetpb.ListDeviceConfigurationsParameters) ([]*fleetpb.DeviceConfiguration, error)
	UpsertDeviceConfiguration(ctx context.Context, deviceConf *fleetpb.DeviceConfiguration) (*fleetpb.DeviceConfiguration, error)
	DeleteDeviceConfiguration(ctx context.Context, id string) error
}
