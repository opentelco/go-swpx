package notification

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/notificationpb"
	"github.com/hashicorp/go-hclog"
	"go.temporal.io/sdk/client"
)

func New(repo Repository, temporalClient client.Client, logger hclog.Logger) (notificationpb.NotificationServiceServer, error) {
	n := &notificationImpl{
		repo:   repo,
		logger: logger.Named("fleet-notitification"),
	}
	w := n.newWorker()
	err := w.Start()
	if err != nil {
		return nil, err
	}
	return n, nil
}

type notificationImpl struct {
	repo           Repository
	logger         hclog.Logger
	temporalClient client.Client

	notificationpb.UnimplementedNotificationServiceServer
}

func (n *notificationImpl) GetByID(ctx context.Context, params *notificationpb.GetByIDRequest) (*notificationpb.Notification, error) {
	if params.Id == "" {
		return nil, ErrNotificationNotFound
	}
	return n.repo.GetByID(ctx, params.Id)
}

func (n *notificationImpl) List(ctx context.Context, params *notificationpb.ListRequest) (*notificationpb.ListResponse, error) {
	res, err := n.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *notificationImpl) Create(ctx context.Context, params *notificationpb.CreateRequest) (*notificationpb.Notification, error) {
	notif := &notificationpb.Notification{
		Title:        params.Title,
		ResourceId:   params.ResourceId,
		ResourceType: params.ResourceType,
	}
	if params.Message != nil {
		notif.Message = *params.Message
	}
	return n.repo.Upsert(ctx, notif)
}

func (n *notificationImpl) Delete(ctx context.Context, params *notificationpb.DeleteRequest) (*notificationpb.DeleteResponse, error) {
	err := n.repo.Delete(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	return &notificationpb.DeleteResponse{}, nil
}

func (n *notificationImpl) MarkAsRead(ctx context.Context, params *notificationpb.MarkAsReadRequest) (*notificationpb.MarkAsReadResponse, error) {
	res, err := n.repo.MarkAsRead(ctx, params)
	if err != nil {
		return nil, err
	}
	return res, nil
}
