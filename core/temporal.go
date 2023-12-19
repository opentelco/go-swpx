package core

import (
	"context"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	TemporalTaskQueue        = "swpx-core"
	TemporalClientContextKey = "temporal_client"
)

// createTemporalWorker creates a new temporal worker and returns it
// can run the w.Run() in a go routine or whatever
func createTemporalWorker(tc client.Client, c *Core) worker.Worker {

	w := worker.New(tc, TemporalTaskQueue, worker.Options{
		BackgroundActivityContext: context.WithValue(context.Background(), TemporalClientContextKey, tc),
	})
	w.RegisterWorkflowWithOptions(
		RunDiagnosticWorkflow,
		workflow.RegisterOptions{
			// Name: "",
		})

	// w.RegisterWorkflowWithOptions(
	// 	MutexWorkflow,
	// 	workflow.RegisterOptions{
	// 		Name: dncEnums.MutexWorkflowName,
	// 	})

	act := NewActivities(c)
	w.RegisterActivity(act)
	return w
}
