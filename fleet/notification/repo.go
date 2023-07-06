package notification

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/notificationpb"
)

type Repository interface {
	// GetByID returns a notification by its ID
	GetByID(ctx context.Context, id string) (*notificationpb.Notification, error)

	// List returns a list of notifications
	List(ctx context.Context, params *notificationpb.ListRequest) (*notificationpb.ListResponse, error)

	// Upsert creates a new notification or updates a notification
	Upsert(ctx context.Context, nonotification *notificationpb.Notification) (*notificationpb.Notification, error)

	// Mark one or more notifications as read
	MarkAsRead(ctx context.Context, params *notificationpb.MarkAsReadRequest) (*notificationpb.MarkAsReadResponse, error)

	// Delete deletes an existing notification
	Delete(ctx context.Context, id string) error
}
