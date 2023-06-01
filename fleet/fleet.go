package fleet

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleetpb"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func New(repo Repository, logger hclog.Logger) fleetpb.FleetServer {
	return &fleet{
		repo:   repo,
		logger: logger.Named("fleet"),
	}
}

type fleet struct {
	repo   Repository
	logger hclog.Logger

	fleetpb.UnimplementedFleetServer
}

// *** Device ***
// Get a device by its ID, this is used to get a specific device
func (f *fleet) GetDeviceByID(ctx context.Context, params *fleetpb.GetDeviceByIDParameters) (*fleetpb.Device, error) {
	if params.Id == "" {
		return nil, ErrDeviceNotFound
	}
	return f.repo.GetDeviceByID(ctx, params.Id)
}

// CollectDevice collects information about the device from the network (with the help of the poller)
// and returns the device with the updated information
func (f *fleet) CollectDevice(ctx context.Context, params *fleetpb.CollectDeviceParameters) (*fleetpb.Device, error) {
	return nil, ErrNotImplemented
}

// Get a device by its hostname, managment ip or serial number etc (used to search for a device)
func (f *fleet) ListDevices(ctx context.Context, params *fleetpb.ListDevicesParameters) (*fleetpb.ListDevicesResponse, error) {
	devices, err := f.repo.ListDevices(ctx, params)
	return &fleetpb.ListDevicesResponse{
		Devices: devices,
	}, err

}

// Create a device in the fleet
func (f *fleet) CreateDevice(ctx context.Context, params *fleetpb.CreateDeviceParameters) (*fleetpb.Device, error) {

	if params.PollerProvider == "" {
		params.PollerProvider = "default_provider"
	}

	if params.Hostname == "" {
		return nil, ErrHostnameRequired
	}

	device := &fleetpb.Device{
		Hostname:             params.Hostname,
		Domain:               params.Domain,
		ManagementIp:         params.ManagementIp,
		SerialNumber:         params.SerialNumber,
		Model:                params.Model,
		Version:              params.Version,
		PollerResourcePlugin: params.PollerResourcePlugin,
		PollerProvider:       params.PollerProvider,
	}

	return f.repo.UpsertDevice(ctx, device)
}

// Update a device in the fleet (this is used to update the device with new information)
func (f *fleet) UpdateDevice(ctx context.Context, params *fleetpb.UpdateDeviceParameters) (*fleetpb.Device, error) {
	return nil, nil
}

// Delete a device from the fleet (mark the device as deleted)
func (f *fleet) DeleteDevice(ctx context.Context, params *fleetpb.DeleteDeviceParameters) (*emptypb.Empty, error) {
	return nil, nil
}

// *** Configuration ***
// CollectConfig collects the running configuration from the device in the network (with the help of the poller) and
// returns the config as a string
func (f *fleet) CollectConfig(ctx context.Context, params *fleetpb.CollectConfigParameters) (*fleetpb.DeviceConfiguration, error) {
	return nil, nil
}

// Get a device configuration by its ID, this is used to get a specific device configuration
func (f *fleet) GetDeviceConfigurationByID(ctx context.Context, params *fleetpb.GetDeviceConfigurationByIDParameters) (*fleetpb.DeviceConfiguration, error) {
	return nil, nil
}

// CompareDeviceConfiguration compares the configuration of a device with the configuration in the database and returns the changes
// if no specific configuration is specified the latest configuration is used to compare with
func (f *fleet) CompareDeviceConfiguration(ctx context.Context, params *fleetpb.CompareDeviceConfigurationParameters) (*fleetpb.CompareDeviceConfigurationResponse, error) {
	return nil, nil
}

// ListDeviceConfigurations lists all configurations for a device
func (f *fleet) ListDeviceConfigurations(ctx context.Context, params *fleetpb.ListDeviceConfigurationsParameters) (*fleetpb.ListDeviceConfigurationsResponse, error) {
	return nil, nil
}

// Create a device configuration in the fleet (this is used to store the configuration of a device)
func (f *fleet) CreateDeviceConfiguration(ctx context.Context, params *fleetpb.CreateDeviceConfigurationParameters) (*fleetpb.DeviceConfiguration, error) {
	return nil, nil
}

// Delete a device configuration from the fleet (removes the configuration from the database)
func (f *fleet) DeleteDeviceConfiguration(ctx context.Context, params *fleetpb.DeleteDeviceConfigurationParameters) (*emptypb.Empty, error) {
	return nil, nil
}

// used by the repo to generate a new ID for a device or configuration
func NewID() string {
	guid, _ := uuid.GenerateUUID()
	return guid
}
