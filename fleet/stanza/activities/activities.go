package activities

import (
	"time"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"go.temporal.io/sdk/workflow"
)

type Activities struct {
	stanza stanzapb.StanzaServiceServer
}

func New(service stanzapb.StanzaServiceServer) *Activities {
	return &Activities{
		stanza: service,
	}
}

func ActivityOptionsNewNotification(ctx workflow.Context) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(),
		StartToCloseTimeout: time.Minute * 1,
		WaitForCancellation: false,
	})
}
