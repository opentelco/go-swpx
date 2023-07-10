package stanza

import (
	"git.liero.se/opentelco/go-swpx/fleet/stanza/activities"
	"git.liero.se/opentelco/go-swpx/fleet/stanza/workflows"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func (n *stanzaImpl) newWorker() worker.Worker {
	w := worker.New(n.temporalClient, stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(), worker.Options{})
	w.RegisterWorkflowWithOptions(
		workflows.ApplyStanzaWorkflow,
		workflow.RegisterOptions{
			Name: "fleet.stanza.apply",
		})

	act := activities.New(n)

	w.RegisterActivity(act)
	return w
}
