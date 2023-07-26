package workflows

import (
	"fmt"
	"regexp"
	"time"

	"git.liero.se/opentelco/go-swpx/fleet/device"
	"git.liero.se/opentelco/go-swpx/fleet/stanza/activities"
	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	stb "git.liero.se/opentelco/go-swpx/proto/go/stanzapb"
	"go.temporal.io/sdk/workflow"
)

var devAct = device.Activities{}
var stanAct = activities.Activities{}

func getDeviceById(ctx workflow.Context, deviceId string) (*devicepb.Device, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           devicepb.TaskQueue_TASK_QUEUE_FLEET_DEVICE.String(),
		StartToCloseTimeout: time.Second * 60,
		WaitForCancellation: false,
	})
	var device devicepb.Device
	if err := workflow.ExecuteActivity(ctx, devAct.GetDeviceByID, deviceId).Get(ctx, &device); err != nil {
		return nil, fmt.Errorf("failed to collect device: %w", err)
	}
	return &device, nil

}

// cloneStanza clones the stanza and returns a new stanza with the same content but tied to the device
func cloneStanza(ctx workflow.Context, deviceId string, stanzaId string) (*stanzapb.Stanza, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(),
		StartToCloseTimeout: time.Second * 60,
		WaitForCancellation: false,
	})
	var device devicepb.Device
	if err := workflow.ExecuteActivity(ctx, stanAct.Clone, deviceId, stanzaId).Get(ctx, &device); err != nil {
		return nil, fmt.Errorf("failed to collect device: %w", err)
	}

	return nil, nil

}

func applyStanza(ctx workflow.Context, device *devicepb.Device, stanza *stanzapb.Stanza) (*stanzapb.ApplyResponse, error) {
	ctx = activities.ActivityOptionsApply(ctx)

	var target string
	if device.ManagementIp != "" {
		target = device.ManagementIp
	} else {
		target = device.Hostname
	}

	var resp stanzapb.ApplyResponse

	// regex to split content for each new line
	reNewLine := regexp.MustCompile(`\r?\n`)
	lines := []*stb.ConfigurationLine{}

	for _, l := range reNewLine.Split(stanza.Content, -1) {
		lines = append(lines, &stb.ConfigurationLine{
			Content: l,
		})

	}
	params := &corepb.ConfigureStanzaRequest{
		Session: &corepb.SessionRequest{
			Hostname:      target,
			NetworkRegion: device.NetworkRegion,
		},
		Settings: &corepb.Settings{
			ProviderPlugin: []string{device.PollerProviderPlugin},
			ResourcePlugin: device.PollerResourcePlugin,
			Timeout:        "90s",
		},
		Stanza: lines,
	}

	if err := workflow.ExecuteActivity(ctx, stanAct.Apply, params).Get(ctx, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func getById(ctx workflow.Context, id string) (*stanzapb.Stanza, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(),
		StartToCloseTimeout: time.Second * 60,
		WaitForCancellation: false,
	})
	var stanza stanzapb.Stanza
	if err := workflow.ExecuteActivity(ctx, stanAct.GetById, id).Get(ctx, &stanza); err != nil {
		return nil, fmt.Errorf("failed to collect stanza: %w", err)
	}
	return &stanza, nil

}
