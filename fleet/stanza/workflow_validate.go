package stanza

import (
	"fmt"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"go.temporal.io/sdk/workflow"
)

func ValidateWorkflow(ctx workflow.Context, params *stanzapb.ValidateRequest) (*stanzapb.ValidateResponse, error) {
	currentWorkflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	logger := workflow.GetLogger(ctx)
	logger.Info("started", "currentWorkflowID", currentWorkflowID)

	stanza, err := getById(ctx, params.Id)
	if err != nil {
		return nil, fmt.Errorf("valditation failed, invalid stanza: %w", err)
	}

	if stanza.DeviceId != nil && *stanza.DeviceId != "" {
		if *stanza.DeviceId != params.DeviceId {
			return nil, fmt.Errorf("valditation failed, stanza already attached to device")
		}
	}

	device, err := getDeviceById(ctx, params.DeviceId)
	if err != nil {
		return nil, fmt.Errorf("valditation failed, invalid device: %w", err)
	}

	return generateFromTemplate(ctx, stanza, device)

}
