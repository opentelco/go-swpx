package stanza

import (
	"context"
	"time"

	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Activities struct {
	stanza    stanzapb.StanzaServiceServer
	commander corepb.CommanderServiceClient
}

func NewActivities(service stanzapb.StanzaServiceServer, commander corepb.CommanderServiceClient) *Activities {
	return &Activities{
		stanza:    service,
		commander: commander,
	}
}

func ActivityOptionsNewNotification(ctx workflow.Context) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(),
		StartToCloseTimeout: time.Minute * 1,
		WaitForCancellation: false,
	})
}

func ActivityOptionsApply(ctx workflow.Context) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           stanzapb.TaskQueue_TASK_QUEUE_FLEET_STANZA.String(),
		StartToCloseTimeout: time.Minute * 1,
		WaitForCancellation: false,
	})
}

func (a *Activities) GetById(ctx context.Context, id string) (*stanzapb.Stanza, error) {
	req := &stanzapb.GetByIDRequest{
		Id: id,
	}
	return a.stanza.GetByID(ctx, req)
}

func (a *Activities) Apply(ctx context.Context, req *corepb.ConfigureStanzaRequest) (*stanzapb.ApplyResponse, error) {

	resp, err := a.commander.ConfigureStanza(ctx, req)
	if err != nil {
		return nil, err
	}
	applyResponse := &stanzapb.ApplyResponse{
		Result: resp.StanzaResult,
	}

	return applyResponse, nil
}

func (a *Activities) Attach(ctx context.Context, req *stanzapb.AttachRequest) (*stanzapb.Stanza, error) {
	return a.stanza.Attach(ctx, req)
}

func (a *Activities) GenerateFromTemplate(ctx context.Context, st *stanzapb.Stanza, device *devicepb.Device) (*stanzapb.Stanza, error) {

	content, err := FromTemplate(st.Id, st.Template, device)
	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError("could not generate template", ErrTypeInvalidTemplate, err)
	}
	// set the content to the template if the content is empty
	if content == "" {
		content = st.Template
	}
	st.Content = content

	revertContent, err := FromTemplate(st.Id, st.RevertTemplate, device)
	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError("could not generate revert template", ErrTypeInvalidRevertTemplate, err)
	}

	if revertContent == "" {
		revertContent = st.RevertTemplate
	}

	st.RevertContent = revertContent

	return st, nil
}

func (a *Activities) SetApplied(ctx context.Context, req *stanzapb.SetAppliedRequest) (*stanzapb.Stanza, error) {
	return a.stanza.SetApplied(ctx, req)
}
