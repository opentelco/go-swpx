package stanza

import (
	"context"

	"git.liero.se/opentelco/go-dnc/models/pb/dispatcherpb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func (n *stanzaImpl) newWorker() worker.Worker {
	w := worker.New(n.temporalClient, stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(), worker.Options{
		BackgroundActivityContext: context.WithValue(context.Background(), dispatcherpb.ContextKey_CLIENT, n.temporalClient),
	})

	w.RegisterWorkflowWithOptions(
		ApplyStanzaWorkflow,
		workflow.RegisterOptions{
			Name: "fleet.stanza.apply",
		})

	w.RegisterWorkflowWithOptions(
		ValidateWorkflow,
		workflow.RegisterOptions{
			Name: "fleet.stanza.validate",
		})

	act := NewActivities(n, n.commanderClient)

	w.RegisterActivity(act)
	return w
}
