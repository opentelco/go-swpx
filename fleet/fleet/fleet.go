package fleet

import (
	"context"
	"fmt"

	"git.liero.se/opentelco/go-swpx/fleet/fleet/workflows"
	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/fleetpb"
	"github.com/hashicorp/go-hclog"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Create a new fleet service and start the Temporal worker for the fleet service
func New(device devicepb.DeviceServiceServer, config configurationpb.ConfigurationServiceServer, poller corepb.CoreServiceClient, tc client.Client, logger hclog.Logger) (fleetpb.FleetServiceServer, error) {
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

	if err := f.startSchedules(); err != nil {
		return nil, err
	}

	return f, nil
}

type fleet struct {
	logger hclog.Logger
	device devicepb.DeviceServiceServer
	config configurationpb.ConfigurationServiceServer
	poller corepb.CoreServiceClient

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
	wf, err := f.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			TaskQueue: TaskQueue,
		},
		workflows.CollectConfigWorkflow,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("could not start collect config workflow: %w", err)
	}
	if params.Blocking {
		var config configurationpb.Configuration
		if err := wf.Get(ctx, &config); err != nil {
			return nil, fmt.Errorf("could not get collect config result: %w", err)
		}
		return &config, nil
	}
	return &configurationpb.Configuration{}, nil

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

func (f *fleet) startSchedules() error {

	workflowOptions := client.StartWorkflowOptions{
		ID:           "schedule.collect-config",
		CronSchedule: "*/5 * * * *",
		TaskQueue:    TaskQueue,
	}

	wf, err := f.temporalClient.ExecuteWorkflow(
		context.Background(),
		workflowOptions,
		workflows.CollectConfigScheduleWorkflow,
	)
	if err != nil {
		return fmt.Errorf("could not start cron workflow CollectConfigScheduleWorkflow: %w", err)
	}

	f.logger.Info("cron workflow started", "workflowID", wf.GetID(), "runID", wf.GetRunID(), "workflow", "CollectConfigScheduleWorkflow")

	/// start cron workflow for collect device
	workflowOptions = client.StartWorkflowOptions{
		ID:           "schedule.collect-device",
		CronSchedule: "*/5 * * * *",
		TaskQueue:    TaskQueue,
	}

	wf, err = f.temporalClient.ExecuteWorkflow(
		context.Background(),
		workflowOptions,
		workflows.CollectDeviceScheduleWorkflow,
	)
	if err != nil {
		return fmt.Errorf("could not start cron workflow CollectDeviceScheduleWorkflow: %w", err)
	}

	f.logger.Info("cron workflow started", "workflowID", wf.GetID(), "runID", wf.GetRunID(), "workflow", "CollectDeviceScheduleWorkflow")

	return nil
}
