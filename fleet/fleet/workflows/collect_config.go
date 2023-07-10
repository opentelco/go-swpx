package workflows

import (
	"fmt"
	"time"

	"git.liero.se/opentelco/go-swpx/fleet/configuration"
	"git.liero.se/opentelco/go-swpx/fleet/fleet/activities"
	"git.liero.se/opentelco/go-swpx/fleet/notification"
	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/fleetpb"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// Get all devices that are flagged for config collection and collect their configs
// if the device cannot be collected, flag it as failed and move on
func CollectConfigWorkflow(ctx workflow.Context, params *fleetpb.CollectConfigParameters) (*configurationpb.Configuration, error) {

	logger := workflow.GetLogger(ctx)

	// get the device from the database
	device, err := getDeviceById(ctx, params.DeviceId)
	if err != nil {
		return nil, err
	}

	collectedConfig, err := runConfigCollection(ctx, device)
	if err != nil {
		if _, err := setDeviceUnreachable(ctx, device.Id); err != nil {
			return nil, err
		}
		if err := addEventCollectConfigFailed(ctx, device.Id, "no respose from device"); err != nil {
			return nil, err
		}

		if _, err := addNotification(ctx,
			notification.NewDeviceConfigurationCollectionFailed(device.Id, device.Hostname, err.Error())); err != nil {
			return nil, err
		}

		return nil, err
	}

	checksum, err := configuration.Hash(collectedConfig.Config)
	if err != nil {
		return nil, fmt.Errorf("could not hash config: %w", err)
	}

	listOfConfigs, err := listConfigs(ctx, &configurationpb.ListParameters{DeviceId: device.Id})
	if err != nil {
		return nil, err
	}
	var diffs string
	if len(listOfConfigs.Configurations) > 0 {
		logger.Debug("comparing config checksum", "old", listOfConfigs.Configurations[0].Checksum, "new", checksum)
		if listOfConfigs.Configurations[0].Checksum == checksum {

			// add event no changes, not storing config
			if err := addEventCollectConfigNoChanges(ctx, device.Id); err != nil {
				return nil, err
			}

			return listOfConfigs.Configurations[0], nil
		}

		changes, err := diffConfigs(ctx, &configurationpb.DiffParameters{
			ConfigurationA:   listOfConfigs.Configurations[0].Configuration,
			ConfigurationAId: &listOfConfigs.Configurations[0].Id,
			ConfigurationB:   collectedConfig.Config,
		})
		if err != nil {
			return nil, fmt.Errorf("could not diff configs: %w", err)
		}
		diffs = changes.Changes
	}
	storedConfig, err := storeConfig(ctx, &configurationpb.CreateParameters{
		DeviceId:      device.Id,
		Configuration: collectedConfig.Config,
		Checksum:      checksum,
		Changes:       diffs,
	})
	if err != nil {
		return nil, err
	}

	// add event config collected
	if err := addEventCollectConfigSuccess(ctx, device.Id); err != nil {
		return nil, err
	}

	return storedConfig, nil
}

func runConfigCollection(ctx workflow.Context, device *devicepb.Device) (*corepb.CollectConfigResponse, error) {
	var target string
	if device.ManagementIp != "" {
		target = device.ManagementIp
	} else {
		target = device.Hostname
	}
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 2,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:        1,
			NonRetryableErrorTypes: []string{activities.ErrTypeDiscoveryFailed},
		},
	})
	var configResponse corepb.CollectConfigResponse
	if err := workflow.ExecuteActivity(ctx, act.CollectConfigWithPoller, &corepb.CollectConfigRequest{
		Session: &corepb.SessionRequest{
			Hostname:      target,
			NetworkRegion: device.NetworkRegion,
		},
		Settings: &corepb.Settings{
			ResourcePlugin: device.PollerResourcePlugin,
			ProviderPlugin: []string{device.PollerProviderPlugin},
			RecreateIndex:  false,
			Timeout:        "90s",
			TqChannel:      corepb.Settings_CHANNEL_PRIMARY,
			Priority:       corepb.Settings_DEFAULT,
		},
	}).Get(ctx, &configResponse); err != nil {
		return nil, fmt.Errorf("failed to discover device: %w", err)
	}
	return &configResponse, nil

}

func addEventCollectConfigFailed(ctx workflow.Context, deviceId string, reason string) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: thirty,
		WaitForCancellation: false,
	})
	evt := devicepb.Event{
		DeviceId: deviceId,
		Type:     devicepb.Event_CONFIGURATION,
		Message:  fmt.Sprintf("failed to collect configuration for device: %s", reason),
		Action:   devicepb.Event_COLLECT_CONFIG,
		Outcome:  devicepb.Event_FAILURE,
	}
	if err := workflow.ExecuteActivity(ctx, devAct.AddDeviceEvent, &evt).Get(ctx, &evt); err != nil {
		return fmt.Errorf("failed to add 'collect configuration failed' event to device: %w", err)
	}
	return nil
}

func addEventCollectConfigSuccess(ctx workflow.Context, deviceId string) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: thirty,
		WaitForCancellation: false,
	})
	evt := devicepb.Event{
		DeviceId: deviceId,
		Type:     devicepb.Event_CONFIGURATION,
		Message:  "device cofiguration collected successfully",
		Action:   devicepb.Event_COLLECT_CONFIG,
		Outcome:  devicepb.Event_SUCCESS,
	}
	if err := workflow.ExecuteActivity(ctx, devAct.AddDeviceEvent, &evt).Get(ctx, &evt); err != nil {
		return fmt.Errorf("failed to add 'collect configuration success' event to device: %w", err)
	}
	return nil
}

func addEventCollectConfigNoChanges(ctx workflow.Context, deviceId string) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: thirty,
		WaitForCancellation: false,
	})
	evt := &devicepb.Event{
		DeviceId: deviceId,
		Type:     devicepb.Event_DEVICE,
		Message:  "config has not changed since last collection, not storing",
		Action:   devicepb.Event_COLLECT_CONFIG,
		Outcome:  devicepb.Event_SUCCESS,
	}
	if err := workflow.ExecuteActivity(ctx, devAct.AddDeviceEvent, &evt).Get(ctx, &evt); err != nil {
		return fmt.Errorf("failed to add 'collect configuration no changes' event to device: %w", err)
	}
	return nil
}

func listConfigs(ctx workflow.Context, params *configurationpb.ListParameters) (*configurationpb.ListResponse, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: thirty,
		WaitForCancellation: false,
	})
	var resp configurationpb.ListResponse
	if err := workflow.ExecuteActivity(ctx, act.ListConfigs, params).Get(ctx, &resp); err != nil {
		return nil, fmt.Errorf("failed to list configuration: %w", err)
	}
	return &resp, nil
}

func diffConfigs(ctx workflow.Context, params *configurationpb.DiffParameters) (*configurationpb.DiffResponse, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: thirty,
		WaitForCancellation: false,
	})
	var resp configurationpb.DiffResponse
	if err := workflow.ExecuteActivity(ctx, act.DiffConfigs, params).Get(ctx, &resp); err != nil {
		return nil, fmt.Errorf("failed to diff configurations: %w", err)
	}
	return &resp, nil
}

func storeConfig(ctx workflow.Context, params *configurationpb.CreateParameters) (*configurationpb.Configuration, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: thirty,
		WaitForCancellation: false,
	})
	var resp configurationpb.Configuration
	if err := workflow.ExecuteActivity(ctx, act.CreateConfiguration, params).Get(ctx, &resp); err != nil {
		return nil, fmt.Errorf("failed to store configuration: %w", err)
	}
	return &resp, nil
}
