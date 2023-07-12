package activities

import (
	"context"
	"errors"
	"sync"

	"git.liero.se/opentelco/go-swpx/fleet/configuration"
	"git.liero.se/opentelco/go-swpx/fleet/device"
	"git.liero.se/opentelco/go-swpx/fleet/fleet/utils"
	"git.liero.se/opentelco/go-swpx/fleet/notification"
	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/fleetpb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/notificationpb"
	"github.com/hashicorp/go-hclog"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Activities struct {
	logger       hclog.Logger
	device       devicepb.DeviceServiceServer
	config       configurationpb.ConfigurationServiceServer
	poller       corepb.CoreServiceClient
	fleet        fleetpb.FleetServiceServer
	notification notificationpb.NotificationServiceServer
}

func New(device devicepb.DeviceServiceServer, config configurationpb.ConfigurationServiceServer, fleet fleetpb.FleetServiceServer, poller corepb.CoreServiceClient, notification notificationpb.NotificationServiceServer, logger hclog.Logger) *Activities {

	return &Activities{
		device:       device,
		config:       config,
		poller:       poller,
		fleet:        fleet,
		notification: notification,
		logger:       logger.Named("fleet"),
	}
}

// DiscoverWithPoller is an activity that discovers a device using the switch poller. If the requests fails or returns nil an non retryable error is returned.
func (a *Activities) DiscoverWithPoller(ctx context.Context, params *corepb.DiscoverRequest) (*corepb.DiscoverResponse, error) {
	d, err := a.poller.Discover(ctx, params)
	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError("could not complete discovery with poller", ErrTypeDiscoveryFailed, err)
	}
	if d == nil || d.NetworkElement == nil {
		return nil, temporal.NewNonRetryableApplicationError("no data from poller", ErrTypeDiscoveryFailed, errors.New("no discovery data from poller"))
	}

	return d, nil
}

func (a *Activities) CollectConfigWithPoller(ctx context.Context, params *corepb.CollectConfigRequest) (*corepb.CollectConfigResponse, error) {
	resp, err := a.poller.CollectConfig(ctx, params)
	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError("could not complete config collection with poller", ErrTypeConfigCollectionFailed, err)
	}
	if resp == nil {
		return nil, temporal.NewNonRetryableApplicationError("no data from poller", ErrTypeConfigCollectionFailed, errors.New("no config data from poller"))
	}
	return resp, nil

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

func (a *Activities) CollectConfigsFromDevices(ctx context.Context) error {
	var hasFiringSchedule = true
	var scheduleType = devicepb.Device_Schedule_COLLECT_DEVICE
	var limit int64 = 100
	var keepGoing = true
	result := make(chan struct{})

	a.logger.Debug("collecting configs from devices...")
	// heartbeat for the activity
	go func() {
		devicesCollected := 0
		for {
			select {
			case <-result:
				devicesCollected += 1
				activity.RecordHeartbeat(ctx, devicesCollected)
				a.logger.Debug("heartbeat", "devicesConfigCollected", devicesCollected)
			case <-ctx.Done():
				return
			}
		}
	}()

	for keepGoing {

		res, err := a.device.List(ctx, &devicepb.ListParameters{
			HasFiringSchedule: &hasFiringSchedule,
			ScheduleType:      &scheduleType,
			Limit:             &limit,
		})
		if err != nil {
			return err
		}
		a.logger.Debug("got devices to collect configs for", "count", len(res.Devices))
		// stop if we have no devices
		if len(res.Devices) == 0 {
			keepGoing = false
		}

		wg := sync.WaitGroup{}
		for _, d := range res.Devices {
			wg.Add(1)

			go func(d *devicepb.Device) {
				defer wg.Done()

				a.logger.Debug("collecting config for device", "device", d.Id)
				s := utils.GetDeviceScheduleByType(d, devicepb.Device_Schedule_COLLECT_DEVICE)
				_, err := a.fleet.CollectConfig(ctx, &fleetpb.CollectConfigParameters{
					DeviceId: d.Id,
					Blocking: true,
				})
				if err != nil {
					a.logger.Warn("could not collect config", "device", d.Id, "error", err)
					s.FailedCount += 1
				} else {
					s.FailedCount = 0
				}

				s.LastRun = timestamppb.Now()

				// if we have too many failures, deactivate the schedule
				if s.FailedCount >= device.MaxScheduleFailures {
					s.Active = false
					s.FailedCount = 0

					// update schedule
					if _, err := a.device.SetSchedule(ctx, &devicepb.SetScheduleParameters{
						DeviceId: d.Id,
						Schedule: s,
					}); err != nil {
						a.logger.Error("could not deactivate schedule", "error", err)
					}

					if _, err := a.device.AddEvent(ctx, &devicepb.Event{
						DeviceId: d.Id,
						Type:     devicepb.Event_DEVICE,
						Message:  "device schedule deactivated due to too many failures",
						Action:   devicepb.Event_COLLECT_CONFIG,
						Outcome:  devicepb.Event_FAILURE,
					}); err != nil {
						a.logger.Error("could not add event about deactivate device schedule", "error", err)
					}

					if _, err := a.notification.Create(ctx, notification.NewDeviceConfigScheduleCancelled(d.Id, d.Hostname)); err != nil {
						a.logger.Error("could not create notification", "error", err)
					}

				} else { // no failures, update schedule and continue

					// update schedule
					s.LastRun = timestamppb.Now()
					if _, err := a.device.SetSchedule(ctx, &devicepb.SetScheduleParameters{
						DeviceId: d.Id,
						Schedule: s,
					}); err != nil {
						a.logger.Error("could not update schedule (collect config)", "error", err)
					}
				}

				result <- struct{}{}

			}(d)

		}

		wg.Wait()
	}

	return nil

}

// CollectDevices is used in the schedule workflow to collect information from devices
// that have a firing schedule of type COLLECT_DEVICE
func (a *Activities) CollectDevices(ctx context.Context) error {
	var hasFiringSchedule = true
	var scheduleType = devicepb.Device_Schedule_COLLECT_DEVICE
	var limit int64 = 100
	var keepGoing = true
	result := make(chan struct{})

	// heartbeat for the activity
	go func() {
		devicesCollected := 0
		for {
			select {
			case <-result:
				devicesCollected += 1
				activity.RecordHeartbeat(ctx, devicesCollected)
				a.logger.Debug("heartbeat", "devicesCollected", devicesCollected)
			case <-ctx.Done():
				return
			}
		}
	}()

	for keepGoing {

		res, err := a.device.List(ctx, &devicepb.ListParameters{
			HasFiringSchedule: &hasFiringSchedule,
			ScheduleType:      &scheduleType,
			Limit:             &limit,
		})
		if err != nil {
			return err
		}
		a.logger.Debug("devices marked for collection by schedule", "count", len(res.Devices))
		// stop if we have no devices
		if len(res.Devices) == 0 {
			keepGoing = false
			a.logger.Debug("no devices marked for collection by schedule")
			return nil
		}

		wg := sync.WaitGroup{}
		for _, d := range res.Devices {
			wg.Add(1)

			go func(d *devicepb.Device) {
				defer wg.Done()
				a.logger.Debug("collect device", "device", d.Id, "hostname", d.Hostname)
				s := utils.GetDeviceScheduleByType(d, devicepb.Device_Schedule_COLLECT_DEVICE)
				_, err := a.fleet.CollectDevice(ctx, &fleetpb.CollectDeviceParameters{
					DeviceId: d.Id,
					Blocking: true,
				})

				if err != nil {
					a.logger.Warn("could not collect device", "device", d.Id, "error", err)
					s.FailedCount += 1
				} else {
					s.FailedCount = 0
				}

				s.LastRun = timestamppb.Now()
				if _, err := a.device.SetSchedule(ctx, &devicepb.SetScheduleParameters{
					DeviceId: d.Id,
					Schedule: s,
				}); err != nil {
					a.logger.Error("could not update schedule (collect device)", "error", err)
				}
				result <- struct{}{}

			}(d)

		}

		wg.Wait()
	}

	return nil

}
