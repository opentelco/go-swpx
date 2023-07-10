package stanza

import (
	"context"
	"fmt"

	"git.liero.se/opentelco/go-swpx/fleet/stanza/workflows"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"github.com/hashicorp/go-hclog"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/types/known/emptypb"
)

const TaskQueue = "STANZA"

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

func (s *stanzaImpl) GetByID(ctx context.Context, params *stanzapb.GetByIDRequest) (*stanzapb.Stanza, error) {
	if params.Id == "" {
		return nil, ErrNotificationNotFound
	}
	return s.repo.GetByID(ctx, params.Id)
}

func (s *stanzaImpl) List(ctx context.Context, params *stanzapb.ListRequest) (*stanzapb.ListResponse, error) {
	res, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *stanzaImpl) Create(ctx context.Context, params *stanzapb.CreateRequest) (*stanzapb.Stanza, error) {
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

	return s.repo.Upsert(ctx, st)
}

func (s *stanzaImpl) Delete(ctx context.Context, params *stanzapb.DeleteRequest) (*emptypb.Empty, error) {
	err := s.repo.Delete(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *stanzaImpl) Apply(ctx context.Context, params *stanzapb.ApplyRequest) (*stanzapb.ApplyResponse, error) {

	stanza, err := s.GetByID(ctx, &stanzapb.GetByIDRequest{Id: params.Id})
	if err != nil {
		return nil, fmt.Errorf("could not get stanza: %w", err)
	}

	wf, err := s.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			TaskQueue: TaskQueue,
		},
		workflows.ApplyStanzaWorkflow,
		params.DeviceId,
		stanza,
	)
	if err != nil {
		return nil, fmt.Errorf("could not start apply workflow: %w", err)
	}

	if params.Blocking {
		if err := wf.Get(ctx, nil); err != nil {
			return nil, fmt.Errorf("could not get apply result: %w", err)
		}
		return &stanzapb.ApplyResponse{}, nil
	}
	return &stanzapb.ApplyResponse{}, nil
}

func (s *stanzaImpl) Revert(ctx context.Context, params *stanzapb.RevertRequest) (*stanzapb.RevertResponse, error) {
	return nil, ErrNotImplemented
}
