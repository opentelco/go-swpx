package fleet

import (
	"git.liero.se/opentelco/go-swpx/fleet/fleet/activities"
	"git.liero.se/opentelco/go-swpx/fleet/fleet/workflows"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/fleetpb"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func (f *fleet) newWorker() worker.Worker {
	w := worker.New(f.temporalClient, fleetpb.TaskQueue_TASK_QUEUE_FLEET.String(), worker.Options{})

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

	w.RegisterWorkflowWithOptions(
		workflows.CollectConfigScheduleWorkflow,
		workflow.RegisterOptions{
			Name: "fleet.schedule.collectConfig",
		})

	w.RegisterWorkflowWithOptions(
		workflows.CollectDeviceScheduleWorkflow,
		workflow.RegisterOptions{
			Name: "fleet.schedule.collectDevice",
		})

	act := activities.New(f.device, f.config, f, f.poller, f.notifications, f.logger)
	w.RegisterActivity(act)
	return w
}
