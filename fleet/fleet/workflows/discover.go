package workflows

import (
	"errors"
	"fmt"
	"net"
	"time"

	"git.liero.se/opentelco/go-swpx/fleet/fleet/activities"
	"git.liero.se/opentelco/go-swpx/fleet/fleet/utils"
	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"github.com/araddon/dateparse"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var safeNow = func(ctx workflow.Context) time.Time {
	return workflow.Now(ctx)
}

var act = &activities.Activities{}

var (
	// dont overwrite theese variables, they are used as pointers in the workflows
	stateActive     = devicepb.Device_DEVICE_STATE_ACTIVE
	statusReachable = devicepb.Device_DEVICE_STATUS_REACHABLE
	statusUnreached = devicepb.Device_DEVICE_STATUS_UNREACHABLE
)

// DiscoverWorkflow is a workflow that discovers a device using the switch poller.
// If the requests fails or returns nil an non retryable error is returned.
func DiscoverWorkflow(ctx workflow.Context, params *devicepb.CreateParameters) (*devicepb.Device, error) {
	logger := workflow.GetLogger(ctx)

	if params == nil {
		return nil, errors.New("invalid parameters")
	}

	target, err := resolveTarget(params)
	if err != nil {
		return nil, fmt.Errorf("could not resolve target: %w", err)
	}
	logger.Info("resolved target", "target", target)

	// discover the device with the help of the poller
	// activity options
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: thirty,
		WaitForCancellation: false,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:        2,
			NonRetryableErrorTypes: []string{activities.ErrTypeDiscoveryFailed},
		},
	})
	var discoverResponse corepb.DiscoverResponse
	if err := workflow.ExecuteActivity(ctx, act.DiscoverWithPoller, &corepb.DiscoverRequest{
		Session: &corepb.SessionRequest{
			Hostname:      target,
			NetworkRegion: *params.NetworkRegion,
		},
		Settings: &corepb.Settings{
			ResourcePlugin: "generic",
			RecreateIndex:  false,
			Timeout:        "15s",
			TqChannel:      corepb.Settings_CHANNEL_PRIMARY,
			Priority:       corepb.Settings_DEFAULT,
		},
	}).Get(ctx, &discoverResponse); err != nil {
		return nil, fmt.Errorf("failed to discover device: %w", err)
	}

	// set the last seen as now
	params.LastSeen = timestamppb.New(safeNow(ctx))
	params.State = &stateActive
	params.Status = &statusReachable

	// set the device parameters
	if discoverResponse.NetworkElement.Sysname != "" {
		params.Sysname = &discoverResponse.NetworkElement.Sysname
		if net.ParseIP(target) != nil {
			params.Hostname = &discoverResponse.NetworkElement.Sysname
			params.ManagementIp = &target
		}
	}
	if discoverResponse.NetworkElement.Version != "" {
		params.Version = &discoverResponse.NetworkElement.Version

		if params.PollerResourcePlugin == nil || *params.PollerResourcePlugin == "" {
			plugin := utils.ParseVersionToResourcePlugin(discoverResponse.NetworkElement.Version)
			params.PollerResourcePlugin = &plugin
		}
	}

	// // if the model is not known set it to the sysObjectId
	if params.Model != nil {
		if discoverResponse.NetworkElement.SnmpObjectId != "" {
			params.Model = &discoverResponse.NetworkElement.SnmpObjectId
		}
	}

	if discoverResponse.NetworkElement.Uptime != "" {
		ts, err := dateparse.ParseAny(discoverResponse.NetworkElement.Uptime)
		if err != nil {
			logger.Warn("could not parse uptime", "err", err)
		} else {
			params.LastReboot = timestamppb.New(ts)
		}
	}

	// parameters are set, create the device in the database
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: thirty,
		WaitForCancellation: false,
	})
	var device devicepb.Device
	if err := workflow.ExecuteActivity(ctx, act.CreateDevice, params).Get(ctx, &device); err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	// create an event for the device creation
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: thirty,
		WaitForCancellation: false,
	})

	var event devicepb.Event
	// returns event but we are not interested in it
	err = workflow.ExecuteActivity(ctx, act.AddDeviceEvent, &devicepb.Event{
		DeviceId: device.Id,
		Type:     devicepb.Event_DEVICE,
		Message:  "device was created by discovery",
		Action:   devicepb.Event_CREATE,
		Outcome:  devicepb.Event_SUCCESS,
	}).Get(ctx, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to create device event: %w", err)
	}

	return &device, nil
}

func resolveTarget(params *devicepb.CreateParameters) (string, error) {
	// validate the request parameters
	if params.ManagementIp == nil && params.Hostname == nil {
		return "", fmt.Errorf("either management ip or hostname must be set")
	}

	// we have not decided how to handle the unique identifier yet
	// ip, hostname, serial number, mac address, uuid, region ?
	var target string

	if params.ManagementIp != nil {
		target = *params.ManagementIp
	}

	if target == "" {
		if params.Hostname != nil && *params.Hostname != "" {
			target = *params.Hostname
			//  if hostname is an ip, set it as management ip
			if net.ParseIP(target) != nil {
				params.ManagementIp = &target
			}
		}

	}

	if target == "" {
		return "", fmt.Errorf("could not determine target")
	}
	return target, nil
}
