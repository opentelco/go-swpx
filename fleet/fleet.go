package fleet

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/fleetpb"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/types/known/emptypb"
)

func New(device devicepb.DeviceServiceServer, config configurationpb.ConfigurationServiceServer, poller core.CoreServiceClient, logger hclog.Logger) fleetpb.FleetServiceServer {
	return &fleet{
		device: device,
		config: config,
		poller: poller,
		logger: logger.Named("fleet"),
	}
}

type fleet struct {
	logger hclog.Logger
	device devicepb.DeviceServiceServer
	config configurationpb.ConfigurationServiceServer
	poller core.CoreServiceClient

	fleetpb.UnimplementedFleetServiceServer
}

// CollectDevice collects information about the device from the network (with the help of the poller)
// and returns the device with the updated information
func (f *fleet) CollectDevice(ctx context.Context, params *fleetpb.CollectDeviceParameters) (*devicepb.Device, error) {
	return nil, ErrNotImplemented
}

// CollectConfig collects the running configuration from the device in the network (with the help of the poller) and
// returns the config as a string
func (f *fleet) CollectConfig(ctx context.Context, params *fleetpb.CollectConfigParameters) (*configurationpb.Configuration, error) {
	return nil, ErrNotImplemented
}

// DeleteDevice deletes the device, its configuration and all changes related to the device
func (f *fleet) DeleteDevice(ctx context.Context, params *devicepb.DeleteParameters) (*emptypb.Empty, error) {
	_, err := f.device.Delete(ctx, params)
	if err != nil {
		return nil, err
	}

	_, err = f.config.Delete(ctx, &configurationpb.DeleteParameters{DevceId: params.Id})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil

}
