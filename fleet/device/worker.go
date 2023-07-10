package device

import (
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"go.temporal.io/sdk/worker"
)

func (d *device) newWorker() worker.Worker {
	w := worker.New(d.temporalClient, stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(), worker.Options{})

	act := NewActivities(d)
	w.RegisterActivity(act)

	return w
}
