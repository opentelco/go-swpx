package workflows

import "go.temporal.io/sdk/workflow"

// Get all devices that are flagged for config collection and collect their configs
// if the device cannot be collected, flag it as failed and move on
func CollectConfigWorkflow(ctx workflow.Context) error {
	return nil
}
