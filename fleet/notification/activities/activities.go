package activities

import (
	"context"
	"time"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/notificationpb"
	"go.temporal.io/sdk/workflow"
)

type Activities struct {
	notification notificationpb.NotificationServiceServer
}

func New(notification notificationpb.NotificationServiceServer) *Activities {
	return &Activities{
		notification: notification,
	}
}

func ActivityOptionsNewNotification(ctx workflow.Context) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           notificationpb.TaskQueue_TASK_QUEUE_FLEET_NOTIFICATIONS.String(),
		StartToCloseTimeout: time.Minute * 1,
		WaitForCancellation: false,
	})
}

func (a *Activities) NewNotification(ctx context.Context, params *notificationpb.CreateRequest) (*notificationpb.Notification, error) {
	not, err := a.notification.Create(ctx, params)
	if err != nil {
		return nil, err
	}
	return not, nil
}
