package fleet

import (
	"git.liero.se/opentelco/go-swpx/fleet/fleet/activities"
	"git.liero.se/opentelco/go-swpx/fleet/fleet/workflows"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const TaskQueue = "FLEET"

func (f *fleet) newWorker() worker.Worker {
	w := worker.New(f.temporalClient, TaskQueue, worker.Options{})

	w.RegisterWorkflowWithOptions(
		workflows.CollectConfigWorkflow,
		workflow.RegisterOptions{
			Name: "fleet.collectConfig",
		})

	w.RegisterWorkflowWithOptions(
		workflows.DiscoverWorkflow,
		workflow.RegisterOptions{
			Name: "fleet.discoverDevice",
		})
	w.RegisterWorkflowWithOptions(
		workflows.CollectDeviceWorkflow,
		workflow.RegisterOptions{
			Name: "fleet.collectDevice",
		})

	act := activities.New(f.device, f.config, f.poller, f.logger)
	w.RegisterActivity(act)
	return w
}
