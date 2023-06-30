package workflows

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

func CollectDeviceScheduleWorkflow(ctx workflow.Context) error {
	if err := runScheduledDeviceCollection(ctx); err != nil {
		return fmt.Errorf("failed to run scheduled config collection: %w", err)
	}
	return nil
}

func runScheduledDeviceCollection(ctx workflow.Context) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 12,
		WaitForCancellation: false,
		HeartbeatTimeout:    time.Minute * 20,
	})
	if err := workflow.ExecuteActivity(ctx, act.CollectDevices).Get(ctx, nil); err != nil {
		return fmt.Errorf("failed to collect devices: %w", err)
	}
	return nil
}
