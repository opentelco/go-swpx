package activities

import (
	"context"
	"time"

	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"go.temporal.io/sdk/workflow"
)

type Activities struct {
	stanza    stanzapb.StanzaServiceServer
	commander corepb.CommanderServiceClient
}

func New(service stanzapb.StanzaServiceServer, commander corepb.CommanderServiceClient) *Activities {
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
	_ = resp

	return nil, nil
}

func (a *Activities) Clone(ctx context.Context, stanzaId string) (*stanzapb.Stanza, error) {

	req := &stanzapb.GetByIDRequest{
		Id: stanzaId,
	}
	stanza, err := a.stanza.GetByID(ctx, req)
	if err != nil {
		return nil, err
	}

	a.stanza.Create(ctx, &stanzapb.CreateRequest{
		Name:          stanza.Name,
		Description:   &stanza.Description,
		Content:       stanza.Content,
		RevertContent: &stanza.RevertContent,
		DeviceType:    *&stanza.DeviceType,
	})

	return nil, nil
}
