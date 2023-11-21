package workflows

import (
	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"go.temporal.io/sdk/workflow"
)

func DiagnosticWorkflow(ctx workflow.Context, request *corepb.DiagnosticRequest) (*corepb.DiagnosticResponse, error) {

	logger := workflow.GetLogger(ctx)
	logger.Info("Diagnostic workflow started")
	defer logger.Info("Diagnostic workflow ended")

	return &corepb.DiagnosticResponse{}, nil
}
