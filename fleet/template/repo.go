package template

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/templatepb"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*templatepb.Template, error)
	List(ctx context.Context, params *templatepb.ListParameters) ([]*templatepb.Template, error)
	Upsert(ctx context.Context, deviceConf *templatepb.Template) (*templatepb.Template, error)
	Delete(ctx context.Context, id string) error
}
