package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"

	"git.liero.se/opentelco/go-swpx/fleet/graph/mappers"
	"git.liero.se/opentelco/go-swpx/fleet/graph/model"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
)

// MarkNotificationsAsRead is the resolver for the markNotificationsAsRead field.
func (r *mutationResolver) MarkNotificationsAsRead(ctx context.Context, input model.MarkNotificationsAsReadParams) ([]*model.Notification, error) {
	panic(fmt.Errorf("not implemented: MarkNotificationsAsRead - markNotificationsAsRead"))
}

// Device is the resolver for the device field.
func (r *queryResolver) Device(ctx context.Context, id string) (*model.Device, error) {
	dev, err := r.devices.GetByID(ctx, &devicepb.GetByIDParameters{Id: id})
	if err != nil {
		return nil, err
	}
	return mappers.ToDevice(dev).ToGQL(), nil
}

// Devices is the resolver for the devices field.
func (r *queryResolver) Devices(ctx context.Context, params *model.ListDevicesParams) (*model.ListDeviceResponse, error) {
	res, err := r.devices.List(ctx, params.ToProto())
	if err != nil {
		return nil, err
	}

	return mappers.ToDeviceListResponse(res).ToGQL(), nil
}

// Notifications is the resolver for the notifications field.
func (r *queryResolver) Notifications(ctx context.Context, params *model.ListNotificationsParams) (*model.ListNotificationsResponse, error) {
	resp, err := r.notifications.List(ctx, params.ToProto())
	if err != nil {
		return nil, err
	}
	return mappers.ToListNotificationResponse(resp).ToGQL(), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
