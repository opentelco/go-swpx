package workflows

import (
	"fmt"

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
	}
	if discoverResponse == nil || discoverResponse.NetworkElement == nil {
		if _, err := setDeviceUnreachable(ctx, device.Id); err != nil {
			return nil, err
		}

		if err := addEventCollectFailed(ctx, device.Id, "no data from device"); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("failed to collect device: %w", err)
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

	if _, err := updateDevice(ctx, updateParams); err != nil {
		return nil, err
	}

	if err := addEventCollectSuccess(ctx, device.Id); err != nil {
		return nil, err
	}

	return nil, nil
}
