package stanza

import (
	"git.liero.se/opentelco/go-swpx/fleet/stanza/activities"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"go.temporal.io/sdk/worker"
)

func (n *stanzaImpl) newWorker() worker.Worker {
	w := worker.New(n.temporalClient, stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(), worker.Options{})

	act := activities.New(n)
	w.RegisterActivity(act)
	return w
}
