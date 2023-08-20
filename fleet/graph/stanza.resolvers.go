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

// Device is the resolver for the device field.
func (r *stanzaResolver) Device(ctx context.Context, obj *model.Stanza) (*model.Device, error) {
	if obj == nil {
		return nil, fmt.Errorf("stanza is nil, cannot get device for nil stanza")
	}

	if obj.Device == nil {
		return nil, nil
	}

	if obj.Device.ID == "" {
		return &model.Device{}, nil
	}
	res, err := r.devices.GetByID(ctx, &devicepb.GetByIDParameters{Id: obj.Device.ID})
	if err != nil {
		return nil, err
	}
	return mappers.ToDevice(res).ToGQL(), nil
}

// Stanza returns StanzaResolver implementation.
func (r *Resolver) Stanza() StanzaResolver { return &stanzaResolver{r} }

type stanzaResolver struct{ *Resolver }
