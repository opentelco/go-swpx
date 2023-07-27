package stanza

import (
	"context"
	"fmt"

	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"github.com/hashicorp/go-hclog"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/types/known/emptypb"
)

func New(repo Repository, temporalClient client.Client, commanderClient corepb.CommanderServiceClient, logger hclog.Logger) (stanzapb.StanzaServiceServer, error) {
	n := &stanzaImpl{
		temporalClient:  temporalClient,
		repo:            repo,
		commanderClient: commanderClient,
		logger:          logger.Named("fleet.stanza"),
	}
	w := n.newWorker()
	err := w.Start()
	if err != nil {
		return nil, err
	}
	return n, nil
}

type stanzaImpl struct {
	repo            Repository
	logger          hclog.Logger
	temporalClient  client.Client
	commanderClient corepb.CommanderServiceClient

	stanzapb.UnimplementedStanzaServiceServer
}

func (s *stanzaImpl) GetByID(ctx context.Context, params *stanzapb.GetByIDRequest) (*stanzapb.Stanza, error) {
	if params.Id == "" {
		return nil, ErrStanzaNotFound
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
		Template:   params.Template,
		DeviceType: params.DeviceType,
	}
	if params.Description != nil {
		st.Description = *params.Description
	}
	if params.RevertTemplate != nil {
		st.RevertTemplate = *params.RevertTemplate
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

	wf, err := s.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			TaskQueue: stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(),
		},
		ApplyStanzaWorkflow,
		params.DeviceId,
		params.Id,
	)
	if err != nil {
		return nil, fmt.Errorf("could not start apply workflow: %w", err)
	}

	if params.Blocking {
		var result stanzapb.ApplyResponse
		if err := wf.Get(ctx, &result); err != nil {
			return nil, fmt.Errorf("could not get apply result: %w", err)
		}
		return &result, nil
	}
	return &stanzapb.ApplyResponse{}, nil
}

func (s *stanzaImpl) Revert(ctx context.Context, params *stanzapb.RevertRequest) (*stanzapb.RevertResponse, error) {
	return nil, ErrNotImplemented
}

// Validate the stanza and process the template to check if it is valid
// if the template cannot be processed an error is returned
func (s *stanzaImpl) Validate(ctx context.Context, params *stanzapb.ValidateRequest) (*stanzapb.ValidateResponse, error) {

	wf, err := s.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			TaskQueue: stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(),
		},
		ValidateWorkflow,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("could not start validate workflow: %w", err)
	}

	var result stanzapb.ValidateResponse
	err = wf.Get(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &result, nil
}

func (s *stanzaImpl) Attach(ctx context.Context, params *stanzapb.AttachRequest) (*stanzapb.Stanza, error) {
	stanza, err := s.repo.GetByID(ctx, params.Id)
	if err != nil {
		return nil, err
	}

	// create a cloned stanza and upsert it
	clonedStanza := &stanzapb.Stanza{
		Name:          stanza.Name,
		Description:   stanza.Description,
		Content:       params.Content,
		RevertContent: params.RevertContent,
		DeviceType:    stanza.DeviceType,
		DeviceId:      &params.DeviceId,
	}
	return s.repo.Upsert(ctx, clonedStanza)

}

func (s *stanzaImpl) SetApplied(ctx context.Context, params *stanzapb.SetAppliedRequest) (*stanzapb.Stanza, error) {
	stanza, err := s.repo.GetByID(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	stanza.AppliedAt = params.AppliedAt

	return s.repo.Upsert(ctx, stanza)

}
