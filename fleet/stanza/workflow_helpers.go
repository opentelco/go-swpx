package stanza

import (
	"fmt"
	"regexp"
	"time"

	"git.liero.se/opentelco/go-swpx/fleet/device"
	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	stb "git.liero.se/opentelco/go-swpx/proto/go/stanzapb"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var devAct = device.Activities{}
var stanAct = Activities{}

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

// attachStanza clones the stanza and returns a new stanza with the same content but tied to the device
func attachStanza(ctx workflow.Context, deviceId string, stanzaId string, validationResult *stanzapb.ValidateResponse) (*stanzapb.Stanza, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(),
		StartToCloseTimeout: time.Second * 60,
		WaitForCancellation: false,
	})

	params := &stanzapb.AttachRequest{
		Id:            stanzaId,
		DeviceId:      deviceId,
		Content:       validationResult.Content,
		RevertContent: validationResult.RevertContent,
	}

	var stanza stanzapb.Stanza
	if err := workflow.ExecuteActivity(ctx, stanAct.Attach, params).Get(ctx, &stanza); err != nil {
		return nil, fmt.Errorf("failed to attach stanza: %w", err)
	}

	return &stanza, nil

}

func applyStanza(ctx workflow.Context, device *devicepb.Device, stanza *stanzapb.Stanza) (*stanzapb.ApplyResponse, error) {
	ctx = ActivityOptionsApply(ctx)

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

func process(ctx workflow.Context, stanza *stanzapb.Stanza, device *devicepb.Device) (*stanzapb.ValidateResponse, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(),
		StartToCloseTimeout: time.Second * 60,
		WaitForCancellation: false,
	})
	var resp stanzapb.ValidateResponse
	if err := workflow.ExecuteActivity(ctx, stanAct.GenerateFromTemplate, stanza, device).Get(ctx, &resp); err != nil {
		return nil, fmt.Errorf("failed to validate stanza: %w", err)
	}
	return &resp, nil

}

func generateFromTemplate(ctx workflow.Context, stanza *stanzapb.Stanza, device *devicepb.Device) (*stanzapb.ValidateResponse, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(),
		StartToCloseTimeout: time.Second * 60,
		WaitForCancellation: false,
	})

	var result stanzapb.Stanza
	if err := workflow.ExecuteActivity(ctx, stanAct.GenerateFromTemplate, stanza, device).Get(ctx, &result); err != nil {
		workflow.GetLogger(ctx).Error("failed to generate stanza content from template", "error", err)
	}

	return &stanzapb.ValidateResponse{
		StanzaId:      result.Id,
		Content:       result.Content,
		RevertContent: result.RevertContent,
	}, nil
}

func setApplied(ctx workflow.Context, stanzaId string) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(),
		StartToCloseTimeout: time.Second * 60,
		WaitForCancellation: false,
	})

	appliedAt := timestamppb.New(workflow.Now(ctx))
	if err := workflow.ExecuteActivity(ctx, stanAct.SetApplied, &stanzapb.SetAppliedRequest{Id: stanzaId, AppliedAt: appliedAt}).Get(ctx, nil); err != nil {
		return fmt.Errorf("failed to set stanza as applied: %w", err)
	}

	return nil

}
