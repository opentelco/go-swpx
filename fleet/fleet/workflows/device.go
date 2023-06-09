package workflows

import (
	"fmt"

	"git.liero.se/opentelco/go-swpx/fleet/fleet/activities"
	"git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func setDeviceUnreachable(ctx workflow.Context, deviceId string) (*devicepb.Device, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: sixty,
		WaitForCancellation: false,
	})
	var device devicepb.Device
	if err := workflow.ExecuteActivity(ctx, act.SetDeviceUnreachable, deviceId).Get(ctx, &device); err != nil {
		return nil, fmt.Errorf("failed to set device unreachable: %w", err)
	}
	return &device, nil
}

func getDeviceById(ctx workflow.Context, deviceId string) (*devicepb.Device, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: sixty,
		WaitForCancellation: false,
	})
	var device devicepb.Device
	if err := workflow.ExecuteActivity(ctx, act.GetDeviceByID, deviceId).Get(ctx, &device); err != nil {
		return nil, fmt.Errorf("failed to collect device: %w", err)
	}
	return &device, nil

}

func addEventCollectFailed(ctx workflow.Context, deviceId string, reason string) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: thirty,
		WaitForCancellation: false,
	})
	evt := devicepb.Event{
		DeviceId: deviceId,
		Type:     devicepb.Event_DEVICE,
		Message:  fmt.Sprintf("failed to collect device: %s", reason),
		Action:   devicepb.Event_COLLECT_DEVICE,
		Outcome:  devicepb.Event_FAILURE,
	}
	if err := workflow.ExecuteActivity(ctx, act.AddDeviceEvent, &evt).Get(ctx, &evt); err != nil {
		return fmt.Errorf("failed to add 'collection failed' event to device: %w", err)
	}
	return nil
}

func addEventCollectSuccess(ctx workflow.Context, deviceId string) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: thirty,
		WaitForCancellation: false,
	})
	evt := devicepb.Event{
		DeviceId: deviceId,
		Type:     devicepb.Event_DEVICE,
		Message:  "device collected successfully",
		Action:   devicepb.Event_COLLECT_DEVICE,
		Outcome:  devicepb.Event_SUCCESS,
	}
	if err := workflow.ExecuteActivity(ctx, act.AddDeviceEvent, &evt).Get(ctx, &evt); err != nil {
		return fmt.Errorf("failed to add 'collection success' event to device: %w", err)
	}
	return nil
}

func runDiscovery(ctx workflow.Context, device *devicepb.Device) (*core.DiscoverResponse, error) {
	var target string
	if device.ManagementIp != "" {
		target = device.ManagementIp
	} else {
		target = device.Hostname
	}
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: thirty,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:        2,
			NonRetryableErrorTypes: []string{activities.ErrTypeDiscoveryFailed},
		},
	})
	var discoverResponse core.DiscoverResponse
	if err := workflow.ExecuteActivity(ctx, act.DiscoverWithPoller, &core.DiscoverRequest{
		Session: &core.SessionRequest{
			Hostname: target,
		},
		Settings: &core.Settings{
			ResourcePlugin: "generic", // todo: make this configurable?
			RecreateIndex:  false,
			Timeout:        "15s",
			TqChannel:      core.Settings_CHANNEL_PRIMARY,
			Priority:       core.Settings_DEFAULT,
		},
	}).Get(ctx, &discoverResponse); err != nil {
		return nil, fmt.Errorf("failed to discover device: %w", err)
	}
	return &discoverResponse, nil
}

func updateDevice(ctx workflow.Context, params *devicepb.UpdateParameters) (*devicepb.Device, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: thirty,
		WaitForCancellation: false,
	})
	var updatedDevice devicepb.Device
	if err := workflow.ExecuteActivity(ctx, act.UpdateDevice, params).Get(ctx, &updatedDevice); err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}
	return &updatedDevice, nil
}
