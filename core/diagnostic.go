package core

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"go.temporal.io/sdk/client"
)

func (c *Core) Diagnostic(ctx context.Context, request *corepb.DiagnosticRequest) (*corepb.DiagnosticResponse, error) {

	c.logger.Info("Diagnostic request received")

	wr, err := c.tc.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		TaskQueue: TemporalTaskQueue,
	}, DiagnosticWorkflow, request)

	if err != nil {
		c.logger.Error("Error executing workflow", "error", err)
		return nil, err
	}

	var resp corepb.DiagnosticResponse
	if err := wr.Get(ctx, &resp); err != nil {
		c.logger.Error("Error getting workflow result", "error", err)
		return nil, err
	}

	return &resp, nil
}
