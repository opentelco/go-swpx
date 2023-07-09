package stanza

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*stanzapb.Stanza, error)
	List(ctx context.Context, params *stanzapb.ListRequest) (*stanzapb.ListResponse, error)
	Upsert(ctx context.Context, deviceConf *stanzapb.Stanza) (*stanzapb.Stanza, error)
	Delete(ctx context.Context, id string) error
}
