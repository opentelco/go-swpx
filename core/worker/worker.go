package worker

import (
	"git.liero.se/opentelco/go-swpx/core"
	"git.liero.se/opentelco/go-swpx/core/activities"
	"git.liero.se/opentelco/go-swpx/core/workflows"
	"github.com/hashicorp/go-hclog"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func New(tc client.Client, taskQueue string, c *core.Core, logger hclog.Logger) worker.Worker {

	w := worker.New(tc, taskQueue, worker.Options{})

	// register workflows
	w.RegisterWorkflowWithOptions(
		workflows.CollectConfigWorkflow,
		workflow.RegisterOptions{
			Name: "swpx.collectConfig",
		})

	// register activities
	act := activities.New(c)
	w.RegisterActivity(act)

	return w
}
