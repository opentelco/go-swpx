package workflows

import (
	"fmt"
	"time"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/fleetpb"
	"go.temporal.io/sdk/workflow"
)

func CollectConfigScheduleWorkflow(ctx workflow.Context) error {
	if err := runScheduledConfigCollection(ctx); err != nil {
		return fmt.Errorf("failed to run scheduled config collection: %w", err)
	}
	return nil
}

func runScheduledConfigCollection(ctx workflow.Context) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           fleetpb.TaskQueue_TASK_QUEUE_FLEET.String(),
		StartToCloseTimeout: time.Hour * 12,
		WaitForCancellation: false,
		HeartbeatTimeout:    time.Minute * 20,
	})
	if err := workflow.ExecuteActivity(ctx, act.CollectConfigsFromDevices).Get(ctx, nil); err != nil {
		return fmt.Errorf("failed to collect configs from devices: %w", err)
	}
	return nil
}
