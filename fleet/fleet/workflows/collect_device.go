package workflows

import (
	"fmt"

	"git.liero.se/opentelco/go-swpx/fleet/notification"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/fleetpb"
	"github.com/araddon/dateparse"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Get all devices that are flagged for collection and get basic info from them
// if the device cannot be collected, flag it as failed and move on
func CollectDeviceWorkflow(ctx workflow.Context, params *fleetpb.CollectDeviceParameters) (*devicepb.Device, error) {

	logger := workflow.GetLogger(ctx)

	// get the device from the database
	device, err := getDeviceById(ctx, params.DeviceId)
	if err != nil {
		return nil, err
	}

	discoverResponse, err := runDiscovery(ctx, device)
	if err != nil {
		if _, err := setDeviceUnreachable(ctx, device.Id); err != nil {
			return nil, err
		}
		if err := addEventCollectFailed(ctx, device.Id, "no respose from device"); err != nil {
			return nil, err
		}
		// add a notification about failure
		if _, err := addNotification(ctx, notification.NewDeviceCollectionFailed(device.Id, device.Hostname, err.Error())); err != nil {
			return nil, err
		}

		return nil, err
	}

	// success
	updateParams := &devicepb.UpdateParameters{
		Id:     device.Id,
		Status: &[]devicepb.Device_Status{devicepb.Device_DEVICE_STATUS_REACHABLE}[0],
	}

	// set the last seen as now
	updateParams.LastSeen = timestamppb.Now()

	if discoverResponse.NetworkElement.Sysname != "" {
		updateParams.Sysname = &discoverResponse.NetworkElement.Sysname
	}
	if discoverResponse.NetworkElement.Version != "" {
		updateParams.Version = &discoverResponse.NetworkElement.Version
	}
	// if the model is not known set it to the sysObjectId
	if updateParams.Model != nil {
		if discoverResponse.NetworkElement.SnmpObjectId != "" {
			updateParams.Model = &discoverResponse.NetworkElement.SnmpObjectId
		}
	}

	if discoverResponse.NetworkElement.Uptime != "" {
		ts, err := dateparse.ParseAny(discoverResponse.NetworkElement.Uptime)
		if err != nil {
			logger.Warn("could not parse uptime", "err", err)
		} else {
			updateParams.LastReboot = timestamppb.New(ts)
		}

	}

	device, err = updateDevice(ctx, updateParams)
	if err != nil {
		return nil, err
	}

	if err := addEventCollectSuccess(ctx, device.Id); err != nil {
		return nil, err
	}

	return device, nil
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
