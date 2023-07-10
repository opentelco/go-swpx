package device

import (
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"go.temporal.io/sdk/worker"
)

func (d *device) newWorker() worker.Worker {
	w := worker.New(d.temporalClient, devicepb.TaskQueue_TASK_QUEUE_FLEET_DEVICE.String(), worker.Options{})

	act := NewActivities(d)
	w.RegisterActivity(act)

	return w
}
