package notification

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/notificationpb"
	"github.com/gogo/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func NewGRPC(notification notificationpb.NotificationServiceServer, serv *grpc.Server) {
	g := &grpcImpl{
		notification: notification,
		grpc:         serv,
	}
	notificationpb.RegisterNotificationServiceServer(serv, g)
}

type grpcImpl struct {
	grpc *grpc.Server

	notification notificationpb.NotificationServiceServer

	notificationpb.UnimplementedNotificationServiceServer
}

func (g *grpcImpl) GetByID(ctx context.Context, params *notificationpb.GetByIDRequest) (*notificationpb.Notification, error) {
	res, err := g.notification.GetByID(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}
func (g *grpcImpl) List(ctx context.Context, params *notificationpb.ListRequest) (*notificationpb.ListResponse, error) {
	res, err := g.notification.List(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}
func (g *grpcImpl) Create(ctx context.Context, params *notificationpb.CreateRequest) (*notificationpb.Notification, error) {
	res, err := g.notification.Create(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}
func (g *grpcImpl) Delete(ctx context.Context, params *notificationpb.DeleteRequest) (*notificationpb.DeleteResponse, error) {
	res, err := g.notification.Delete(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}
func (g *grpcImpl) MarkAsRead(ctx context.Context, params *notificationpb.MarkAsReadRequest) (*notificationpb.MarkAsReadResponse, error) {
	res, err := g.notification.MarkAsRead(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

func errHandler(err error) error {
	if err == nil {
		return nil
	}
	switch err {
	case ErrNotImplemented:
		return status.Errorf(codes.Unimplemented, err.Error())
	default:
		return status.Errorf(codes.Internal, err.Error())
	}
}
