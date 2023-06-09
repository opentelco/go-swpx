package fleet

import (
	"context"
	"fmt"

	"git.liero.se/opentelco/go-swpx/fleet/configuration"
	"git.liero.se/opentelco/go-swpx/fleet/fleet/workflows"
	"git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/fleetpb"
	"github.com/hashicorp/go-hclog"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Create a new fleet service and start the Temporal worker for the fleet service
func New(device devicepb.DeviceServiceServer, config configurationpb.ConfigurationServiceServer, poller core.CoreServiceClient, tc client.Client, logger hclog.Logger) (fleetpb.FleetServiceServer, error) {
	f := &fleet{
		device:         device,
		config:         config,
		poller:         poller,
		temporalClient: tc,
		logger:         logger.Named("fleet"),
	}

	w := f.newWorker()
	if err := w.Start(); err != nil {
		return nil, err
	}

	return f, nil
}

type fleet struct {
	logger hclog.Logger
	device devicepb.DeviceServiceServer
	config configurationpb.ConfigurationServiceServer
	poller core.CoreServiceClient

	temporalClient client.Client

	fleetpb.UnimplementedFleetServiceServer
}

func (f *fleet) DiscoverDevice(ctx context.Context, params *fleetpb.DiscoverDeviceParameters) (*devicepb.Device, error) {

	wf, err := f.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			TaskQueue: TaskQueue,
		},
		workflows.DiscoverWorkflow,
		params.CreateDeviceParams,
	)
	if err != nil {
		return nil, fmt.Errorf("could not start discover workflow: %w", err)
	}
	if params.Blocking {
		var dev devicepb.Device
		if err := wf.Get(ctx, &dev); err != nil {
			return nil, fmt.Errorf("could not get discovery result: %w", err)
		}
		return &dev, nil
	}
	return &devicepb.Device{}, nil
}

// CollectDevice collects information about the device from the network (with the help of the poller)
// and returns the device with the updated information
func (f *fleet) CollectDevice(ctx context.Context, params *fleetpb.CollectDeviceParameters) (*devicepb.Device, error) {

	_, err := f.device.GetByID(ctx, &devicepb.GetByIDParameters{Id: params.DeviceId})
	if err != nil {
		return nil, err
	}
	wf, err := f.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			TaskQueue: TaskQueue,
		},
		workflows.CollectDeviceWorkflow,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("could not start collect workflow: %w", err)
	}
	if params.Blocking {
		var dev devicepb.Device
		if err := wf.Get(ctx, &dev); err != nil {
			return nil, fmt.Errorf("could not get collect result: %w", err)
		}
		return &dev, nil
	}
	return &devicepb.Device{}, nil
}

// CollectConfig collects the running configuration from the device in the network (with the help of the poller) and
// returns the config as a string
func (f *fleet) CollectConfig(ctx context.Context, params *fleetpb.CollectConfigParameters) (*configurationpb.Configuration, error) {

	dev, err := f.device.GetByID(ctx, &devicepb.GetByIDParameters{Id: params.DeviceId})
	if err != nil {
		return nil, err
	}
	if dev.PollerResourcePlugin == "" {
		return nil, fmt.Errorf("no poller plugin defined for device %s", dev.Id)
	}

	var target string
	if dev.ManagementIp != "" {
		target = dev.ManagementIp
	} else {
		target = dev.Hostname
	}

	res, err := f.poller.CollectConfig(ctx, &core.CollectConfigRequest{
		Settings: &core.Settings{
			ProviderPlugin: []string{dev.PollerProviderPlugin},
			ResourcePlugin: dev.PollerResourcePlugin,
			Timeout:        "60s",
			TqChannel:      core.Settings_CHANNEL_PRIMARY,
			Priority:       core.Settings_DEFAULT,
		},
		Session: &core.SessionRequest{
			Hostname:      target,
			NetworkRegion: dev.NetworkRegion,
		},
	})
	if err != nil {
		err = fmt.Errorf("could not collect config: %w", err)

		if _, err := f.device.AddEvent(ctx, &devicepb.Event{
			DeviceId: dev.Id,
			Type:     devicepb.Event_DEVICE,
			Message:  err.Error(),
			Action:   devicepb.Event_COLLECT_CONFIG,
			Outcome:  devicepb.Event_FAILURE,
		}); err != nil {
			f.logger.Warn("could not create event on device", "err", err)
		}

		return nil, fmt.Errorf("could not collect config: %w", err)
	}
	checksum, err := configuration.Hash(res.Config)
	if err != nil {
		return nil, fmt.Errorf("could not hash config: %w", err)
	}

	// determine if config has changed since last collection
	configResult, err := f.config.List(ctx, &configurationpb.ListParameters{
		DeviceId: dev.Id,
	})
	if err != nil {
		f.logger.Warn("could not list configs on device", "err", err)
	}
	var diffs string
	if len(configResult.Configurations) > 0 {
		f.logger.Debug("comparing config checksum", "old", configResult.Configurations[0].Checksum, "new", checksum)
		if configResult.Configurations[0].Checksum == checksum {
			if _, err := f.device.AddEvent(ctx, &devicepb.Event{
				DeviceId: dev.Id,
				Type:     devicepb.Event_DEVICE,
				Message:  "config has not changed since last collection, not storing",
				Action:   devicepb.Event_COLLECT_CONFIG,
				Outcome:  devicepb.Event_SUCCESS,
			}); err != nil {
				f.logger.Warn("could not create event on device", "err", err)
			}

			return configResult.Configurations[0], nil
		}

		changes, err := f.config.Diff(ctx, &configurationpb.DiffParameters{
			ConfigurationA:   configResult.Configurations[0].Configuration,
			ConfigurationAId: &configResult.Configurations[0].Id,
			ConfigurationB:   res.Config,
		})
		diffs = changes.Changes
		if err != nil {
			return nil, fmt.Errorf("could not diff configs: %w", err)
		}
	}

	config, err := f.config.Create(ctx, &configurationpb.CreateParameters{
		DeviceId:      dev.Id,
		Configuration: res.Config,
		Checksum:      checksum,
		Changes:       diffs,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create config: %w", err)
	}

	if _, err := f.device.AddEvent(ctx, &devicepb.Event{
		DeviceId: dev.Id,
		Type:     devicepb.Event_DEVICE,
		Message:  "collected configuration for device through the poller",
		Action:   devicepb.Event_COLLECT_CONFIG,
		Outcome:  devicepb.Event_SUCCESS,
	}); err != nil {
		f.logger.Warn("could not create event on device", "err", err)
	}
	return config, nil
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
