package stanza

import (
	"time"

	"git.liero.se/opentelco/go-dnc/client"
	"git.liero.se/opentelco/go-dnc/models/pb/sharedpb"
	"git.liero.se/opentelco/go-dnc/models/pb/transportpb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"go.temporal.io/sdk/workflow"
)

func ApplyStanzaWorkflow(ctx workflow.Context, deviceId string, stanzaId string) (*stanzapb.ApplyResponse, error) {
	currentWorkflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	// logger := workflow.GetLogger(ctx)

	stanza, err := getById(ctx, stanzaId)
	if err != nil {
		return nil, err
	}

	device, err := getDeviceById(ctx, deviceId)
	if err != nil {
		return nil, err
	}

	var validation stanzapb.ValidateResponse
	if err := workflow.ExecuteChildWorkflow(ctx, ValidateWorkflow, &stanzapb.ValidateRequest{Id: stanzaId, DeviceId: deviceId}).Get(ctx, &validation); err != nil {
		return nil, err
	}

	stanza.Content = validation.Content
	stanza.RevertContent = validation.RevertContent
	resourceID := client.NewResourceID(device.ManagementIp, 22, transportpb.Type_SSH, sharedpb.Channel_CONFIGURE, device.NetworkRegion)

	mutex := client.NewMutex(currentWorkflowID)
	unlockFunc, err := mutex.Lock(ctx, resourceID, 10*time.Minute)
	if err != nil {
		return nil, err
	}
	defer unlockFunc()

	// attach stanza to device (create new stanza version and attach it to device)
	stanza, err = attachStanza(ctx, device.Id, stanza.Id, &validation)
	if err != nil {
		return nil, err
	}

	resp, err := applyStanza(ctx, device, stanza)
	if err != nil {
		return nil, err
	}
	resp.Stanza = stanza

	if err := setApplied(ctx, stanza.Id); err != nil {
		return nil, err
	}

	return resp, nil
}
