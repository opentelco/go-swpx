package configuration

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*configurationpb.Configuration, error)
	List(ctx context.Context, params *configurationpb.ListParameters) ([]*configurationpb.Configuration, error)
	Upsert(ctx context.Context, deviceConf *configurationpb.Configuration) (*configurationpb.Configuration, error)
	Delete(ctx context.Context, id string) error
}
