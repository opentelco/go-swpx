package activities

import (
	"context"
	"errors"

	"git.liero.se/opentelco/go-swpx/fleet/configuration"
	"git.liero.se/opentelco/go-swpx/fleet/device"
	"git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"github.com/hashicorp/go-hclog"
	"go.temporal.io/sdk/temporal"
)

type Activities struct {
	logger hclog.Logger
	device devicepb.DeviceServiceServer
	config configurationpb.ConfigurationServiceServer
	poller core.CoreServiceClient
}

func New(device devicepb.DeviceServiceServer, config configurationpb.ConfigurationServiceServer, poller core.CoreServiceClient, logger hclog.Logger) *Activities {

	return &Activities{
		device: device,
		config: config,
		poller: poller,
		logger: logger.Named("fleet"),
	}
}

// DiscoverWithPoller is an activity that discovers a device using the switch poller. If the requests fails or returns nil an non retryable error is returned.
func (a *Activities) DiscoverWithPoller(ctx context.Context, params *core.DiscoverRequest) (*core.DiscoverResponse, error) {
	d, err := a.poller.Discover(ctx, params)
	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError("could not complete discovery with poller", ErrTypeDiscoveryFailed, err)
	}
	if d == nil || d.NetworkElement == nil {
		return nil, temporal.NewNonRetryableApplicationError("no data from poller", ErrTypeDiscoveryFailed, errors.New("no discovery data from poller"))
	}

	return d, nil
}

func (a *Activities) CollectConfigWithPoller(ctx context.Context, params *core.CollectConfigRequest) (*core.CollectConfigResponse, error) {
	resp, err := a.poller.CollectConfig(ctx, params)
	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError("could not complete config collection with poller", ErrTypeConfigCollectionFailed, err)
	}
	if resp == nil {
		return nil, temporal.NewNonRetryableApplicationError("no data from poller", ErrTypeConfigCollectionFailed, errors.New("no config data from poller"))
	}
	return resp, nil

}

func (a *Activities) CreateDevice(ctx context.Context, params *devicepb.CreateParameters) (*devicepb.Device, error) {
	return a.device.Create(ctx, params)
}

func (a *Activities) AddDeviceEvent(ctx context.Context, event *devicepb.Event) (*devicepb.Event, error) {
	return a.device.AddEvent(ctx, event)
}

func (a *Activities) GetDeviceByID(ctx context.Context, id string) (*devicepb.Device, error) {
	dev, err := a.device.GetByID(ctx, &devicepb.GetByIDParameters{Id: id})
	if errors.Is(err, device.ErrDeviceNotFound) {
		return nil, temporal.NewNonRetryableApplicationError("device not found", ErrTypeDeviceNotFound, err)
	}
	return dev, err
}

func (a *Activities) SetDeviceUnreachable(ctx context.Context, id string) (*devicepb.Device, error) {
	dev, err := a.device.Update(ctx, &devicepb.UpdateParameters{
		Id:     id,
		Status: &[]devicepb.Device_Status{devicepb.Device_DEVICE_STATUS_UNREACHABLE}[0],
	})
	if errors.Is(err, device.ErrDeviceNotFound) {
		return nil, temporal.NewNonRetryableApplicationError("device not found", ErrTypeDeviceNotFound, err)
	}
	return dev, err
}

func (a *Activities) UpdateDevice(ctx context.Context, params *devicepb.UpdateParameters) (*devicepb.Device, error) {
	dev, err := a.device.Update(ctx, params)
	if errors.Is(err, device.ErrDeviceNotFound) {
		return nil, temporal.NewNonRetryableApplicationError("device not found", ErrTypeDeviceNotFound, err)
	}
	return dev, err
}

func (a *Activities) GetConfigurationByID(ctx context.Context, id string) (*configurationpb.Configuration, error) {
	c, err := a.config.GetByID(ctx, &configurationpb.GetByIDParameters{Id: id})
	if errors.Is(err, configuration.ErrConfigurationNotFound) {
		return nil, temporal.NewNonRetryableApplicationError("configuration not found", ErrTypeConfigNotFound, err)
	}
	return c, err
}

func (a *Activities) CreateConfiguration(ctx context.Context, params *configurationpb.CreateParameters) (*configurationpb.Configuration, error) {
	return a.config.Create(ctx, params)
}

func (a *Activities) ListConfigs(ctx context.Context, params *configurationpb.ListParameters) (*configurationpb.ListResponse, error) {
	return a.config.List(ctx, params)
}

func (a *Activities) DiffConfigs(ctx context.Context, params *configurationpb.DiffParameters) (*configurationpb.DiffResponse, error) {
	return a.config.Diff(ctx, params)
}
