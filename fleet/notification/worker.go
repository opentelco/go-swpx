package notification

import (
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/notificationpb"
	"go.temporal.io/sdk/worker"
)

func (n *notificationImpl) newWorker() worker.Worker {
	w := worker.New(n.temporalClient, notificationpb.TaskQueue_TASK_QUEUE_FLEET_NOTIFICATIONS.String(), worker.Options{})

	act := NewActivities(n)
	w.RegisterActivity(act)
	return w
}
