package fleet

import (
	"context"
	"errors"
	"fmt"
	"net"

	"git.liero.se/opentelco/go-swpx/fleet/configuration"
	"git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/fleetpb"
	"github.com/araddon/dateparse"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (f *fleet) DiscoverDevice(ctx context.Context, params *devicepb.CreateParameters) (*devicepb.Device, error) {

	var targetIsIp bool // determine if we have an IP or hostname
	var target string
	if params.ManagementIp != nil {
		res, err := f.device.List(ctx, &devicepb.ListParameters{Hostname: *params.Hostname})
		if err != nil {
			return nil, fmt.Errorf("could not discover and create device: %w", err)
		}
		if len(res.Devices) > 0 {
			return nil, fmt.Errorf("device with hostname %s already exists", *params.Hostname)
		}
		target = *params.ManagementIp

	}

	if params.Hostname != nil {
		res, err := f.device.List(ctx, &devicepb.ListParameters{Hostname: *params.Hostname})
		if err != nil {
			return nil, fmt.Errorf("could not discover and create device: %w", err)
		}
		if len(res.Devices) > 0 {
			return nil, fmt.Errorf("device with hostname %s already exists", *params.Hostname)
		}
		target = *params.Hostname
		if net.ParseIP(target) != nil {
			targetIsIp = true
		}
	}

	d, err := f.poller.Discover(ctx, &core.DiscoverRequest{
		Session: &core.SessionRequest{
			Hostname: target,
		},
		Settings: &core.Settings{
			ResourcePlugin: "generic",
			RecreateIndex:  false,
			Timeout:        "30s",
			TqChannel:      core.Settings_CHANNEL_PRIMARY,
			Priority:       core.Settings_DEFAULT,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("discoveyr failed: %w", err)
	}
	if d.NetworkElement == nil {
		return nil, errors.New("discoveyr failed: could not get any data from the poller")
	}

	// set the last seen as now
	params.LastSeen = timestamppb.Now()

	if d.NetworkElement.Sysname != "" {
		params.Sysname = &d.NetworkElement.Sysname
		if targetIsIp {
			params.Hostname = &d.NetworkElement.Sysname
			params.ManagementIp = &target
		}
	}
	if d.NetworkElement.Version != "" {
		params.Version = &d.NetworkElement.Version
	}

	// if the model is not known set it to the sysObjectId
	if params.Model != nil {
		if d.NetworkElement.SnmpObjectId != "" {
			params.Model = &d.NetworkElement.SnmpObjectId
		}
	}

	f.logger.Debug("device discovered", "device", d.NetworkElement)

	if d.NetworkElement.Uptime != "" {
		ts, err := dateparse.ParseAny(d.NetworkElement.Uptime)
		if err != nil {
			f.logger.Warn("could not parse uptime", "err", err)
		} else {
			params.LastReboot = timestamppb.New(ts)
		}

	}

	dev, err := f.device.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("could not create device: %w", err)
	}

	_, err = f.device.AddEvent(ctx, &devicepb.Event{
		DeviceId: dev.Id,
		Type:     devicepb.Event_DEVICE,
		Message:  "device was created by discovery",
		Action:   devicepb.Event_CREATE,
		Outcome:  devicepb.Event_SUCCESS,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create event on device: %w", err)
	}
	return dev, nil
}

// CollectDevice collects information about the device from the network (with the help of the poller)
// and returns the device with the updated information
func (f *fleet) CollectDevice(ctx context.Context, params *fleetpb.CollectDeviceParameters) (*devicepb.Device, error) {

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

	d, err := f.poller.Discover(ctx, &core.DiscoverRequest{
		Session: &core.SessionRequest{
			Hostname:      target,
			NetworkRegion: dev.NetworkRegion,
		},
		Settings: &core.Settings{
			// todo: change this to the poller plugin of the device when implemented in the plugin
			ResourcePlugin: "generic",
			Timeout:        "60s",
			TqChannel:      core.Settings_CHANNEL_PRIMARY,
			Priority:       core.Settings_DEFAULT,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("collect failed: %w", err)
	}
	if d.NetworkElement == nil {
		return nil, errors.New("collect failed: could not get any data from the poller")
	}
	updateParams := &devicepb.UpdateParameters{
		Id: dev.Id,
	}

	// set the last seen as now
	updateParams.LastSeen = timestamppb.Now()

	if d.NetworkElement.Sysname != "" {
		updateParams.Sysname = &d.NetworkElement.Sysname
	}
	if d.NetworkElement.Version != "" {
		updateParams.Version = &d.NetworkElement.Version
	}
	// if the model is not known set it to the sysObjectId
	if updateParams.Model != nil {
		if d.NetworkElement.SnmpObjectId != "" {
			updateParams.Model = &d.NetworkElement.SnmpObjectId
		}
	}

	f.logger.Debug("device collected", "device", d.NetworkElement)

	if d.NetworkElement.Uptime != "" {
		ts, err := dateparse.ParseAny(d.NetworkElement.Uptime)
		if err != nil {
			f.logger.Warn("could not parse uptime", "err", err)
		} else {
			updateParams.LastReboot = timestamppb.New(ts)
		}

	}

	dev, err = f.device.Update(ctx, updateParams)
	if err != nil {
		return nil, fmt.Errorf("could not collect device: %w", err)
	}

	_, err = f.device.AddEvent(ctx, &devicepb.Event{
		DeviceId: dev.Id,
		Type:     devicepb.Event_DEVICE,
		Message:  "device was updated by collection",
		Action:   devicepb.Event_CREATE,
		Outcome:  devicepb.Event_SUCCESS,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create event on device: %w", err)
	}
	return dev, nil
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
			ConfigurationA: configResult.Configurations[0].Configuration,
			ConfigurationB: res.Config,
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
