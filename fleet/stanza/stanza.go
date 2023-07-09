package stanza

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"github.com/hashicorp/go-hclog"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/types/known/emptypb"
)

func New(repo Repository, temporalClient client.Client, logger hclog.Logger) (stanzapb.StanzaServiceServer, error) {
	n := &stanzaImpl{
		temporalClient: temporalClient,
		repo:           repo,
		logger:         logger.Named("fleet-notitification"),
	}
	w := n.newWorker()
	err := w.Start()
	if err != nil {
		return nil, err
	}
	return n, nil
}

type stanzaImpl struct {
	repo           Repository
	logger         hclog.Logger
	temporalClient client.Client

	stanzapb.UnimplementedStanzaServiceServer
}

func (n *stanzaImpl) GetByID(ctx context.Context, params *stanzapb.GetByIDRequest) (*stanzapb.Stanza, error) {
	if params.Id == "" {
		return nil, ErrNotificationNotFound
	}
	return n.repo.GetByID(ctx, params.Id)
}

func (n *stanzaImpl) List(ctx context.Context, params *stanzapb.ListRequest) (*stanzapb.ListResponse, error) {
	res, err := n.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *stanzaImpl) Create(ctx context.Context, params *stanzapb.CreateRequest) (*stanzapb.Stanza, error) {
	st := &stanzapb.Stanza{
		Name:       params.Name,
		Content:    params.Content,
		DeviceType: params.DeviceType,
	}
	if params.Description != nil {
		st.Description = *params.Description
	}
	if params.RevertContent != nil {
		st.RevertContent = *params.RevertContent
	}

	return n.repo.Upsert(ctx, st)
}

func (n *stanzaImpl) Delete(ctx context.Context, params *stanzapb.DeleteRequest) (*emptypb.Empty, error) {
	err := n.repo.Delete(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
