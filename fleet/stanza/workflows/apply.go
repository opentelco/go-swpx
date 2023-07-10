package workflows

import (
	"time"

	"git.liero.se/opentelco/go-dnc/client"
	"git.liero.se/opentelco/go-dnc/models/pb/sharedpb"
	"git.liero.se/opentelco/go-dnc/models/pb/transportpb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"go.temporal.io/sdk/workflow"
)

func ApplyStanzaWorkflow(ctx workflow.Context, deviceId string, stanza *stanzapb.Stanza) error {
	currentWorkflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	logger := workflow.GetLogger(ctx)

	device, err := getDeviceById(ctx, deviceId)
	if err != nil {
		return err
	}

	resourceID := client.NewResourceID(device.ManagementIp, 22, transportpb.Type_SSH, sharedpb.Channel_CONFIGURE, device.NetworkRegion)

	logger.Info("started", "currentWorkflowID", currentWorkflowID, "resourceID", resourceID)

	mutex := client.NewMutex(currentWorkflowID)
	unlockFunc, err := mutex.Lock(ctx, resourceID, 10*time.Minute)
	if err != nil {
		return err
	}
	defer unlockFunc()
	logger.Info("resource locked")

	// emulate long running process
	logger.Info("critical operation started")
	_ = workflow.Sleep(ctx, 10*time.Second)
	logger.Info("critical operation finished")

	logger.Info("finished")
	return nil
}
