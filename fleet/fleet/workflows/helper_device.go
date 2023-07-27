package workflows

import (
	"fmt"

	"git.liero.se/opentelco/go-swpx/fleet/device"
	"git.liero.se/opentelco/go-swpx/fleet/fleet/activities"
	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/fleetpb"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

var devAct = &device.Activities{}

func setDeviceUnreachable(ctx workflow.Context, deviceId string) (*devicepb.Device, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           devicepb.TaskQueue_TASK_QUEUE_FLEET_DEVICE.String(),
		StartToCloseTimeout: sixty,
		WaitForCancellation: false,
	})
	var device devicepb.Device
	if err := workflow.ExecuteActivity(ctx, devAct.SetDeviceUnreachable, deviceId).Get(ctx, &device); err != nil {
		return nil, fmt.Errorf("failed to set device unreachable: %w", err)
	}
	return &device, nil
}

func getDeviceById(ctx workflow.Context, deviceId string) (*devicepb.Device, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           devicepb.TaskQueue_TASK_QUEUE_FLEET_DEVICE.String(),
		StartToCloseTimeout: sixty,
		WaitForCancellation: false,
	})
	var device devicepb.Device
	if err := workflow.ExecuteActivity(ctx, devAct.GetDeviceByID, deviceId).Get(ctx, &device); err != nil {
		return nil, fmt.Errorf("failed to collect device: %w", err)
	}
	return &device, nil

}

func runDiscovery(ctx workflow.Context, device *devicepb.Device) (*corepb.DiscoverResponse, error) {
	var target string
	if device.ManagementIp != "" {
		target = device.ManagementIp
	} else {
		target = device.Hostname
	}
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: thirty,
		TaskQueue:           fleetpb.TaskQueue_TASK_QUEUE_FLEET.String(),
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:        1,
			NonRetryableErrorTypes: []string{activities.ErrTypeDiscoveryFailed},
		},
	})
	var discoverResponse corepb.DiscoverResponse
	if err := workflow.ExecuteActivity(ctx, act.DiscoverWithPoller, &corepb.DiscoverRequest{
		Session: &corepb.SessionRequest{
			Hostname:      target,
			NetworkRegion: device.NetworkRegion,
		},
		Settings: &corepb.Settings{
			ResourcePlugin: "generic", // todo: make this configurable?
			RecreateIndex:  false,
			Timeout:        "15s",
			TqChannel:      corepb.Settings_CHANNEL_PRIMARY,
			Priority:       corepb.Settings_DEFAULT,
		},
	}).Get(ctx, &discoverResponse); err != nil {
		return nil, fmt.Errorf("failed to discover device: %w", err)
	}
	return &discoverResponse, nil
}

func updateDevice(ctx workflow.Context, params *devicepb.UpdateParameters) (*devicepb.Device, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           devicepb.TaskQueue_TASK_QUEUE_FLEET_DEVICE.String(),
		StartToCloseTimeout: thirty,
		WaitForCancellation: false,
	})
	var updatedDevice devicepb.Device
	if err := workflow.ExecuteActivity(ctx, devAct.UpdateDevice, params).Get(ctx, &updatedDevice); err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}
	return &updatedDevice, nil
}

func listDevices(ctx workflow.Context, params *devicepb.ListParameters) ([]*devicepb.Device, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           devicepb.TaskQueue_TASK_QUEUE_FLEET_DEVICE.String(),
		StartToCloseTimeout: sixty,
		WaitForCancellation: false,
	})
	var resp devicepb.ListResponse
	if err := workflow.ExecuteActivity(ctx, devAct.ListDevices, params).Get(ctx, &resp); err != nil {
		return nil, fmt.Errorf("failed to collect device: %w", err)
	}
	return resp.Devices, nil

}

func setScheduleLastRun(ctx workflow.Context, device *devicepb.Device, st devicepb.Device_Schedule_Type) (*devicepb.Device, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           devicepb.TaskQueue_TASK_QUEUE_FLEET_DEVICE.String(),
		StartToCloseTimeout: thirty,
		WaitForCancellation: false,
	})

	var resp devicepb.Device
	if err := workflow.ExecuteActivity(ctx, devAct.SetScheduleLastRun,
		device.Id,
		st,
		workflow.Now(ctx),
	).Get(ctx, &resp); err != nil {
		return nil, fmt.Errorf("failed to set schedule last run: %w", err)
	}
	return &resp, nil
}