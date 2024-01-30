package core

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelco.io/go-swpx/proto/go/analysispb"
	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflow/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var fingerprintSearchAttributeKey = "transaction_id"

func (c *Core) RunDiagnostic(ctx context.Context, request *corepb.RunDiagnosticRequest) (*corepb.RunDiagnosticResponse, error) {

	c.logger.Info("Diagnostic request received")

	opts := client.StartWorkflowOptions{
		TaskQueue: TemporalTaskQueue,
	}
	if request.Fingerprint != "" {
		opts.SearchAttributes = map[string]interface{}{fingerprintSearchAttributeKey: request.Fingerprint}
	}

	wr, err := c.tc.ExecuteWorkflow(ctx, opts, RunDiagnosticWorkflow, request)

	if err != nil {
		c.logger.Error("Error executing workflow", "error", err)
		return nil, err
	}

	return &corepb.RunDiagnosticResponse{Id: wr.GetID()}, nil
}

func (c *Core) RunQuickDiagnostic(ctx context.Context, request *corepb.RunDiagnosticRequest) (*corepb.RunDiagnosticResponse, error) {
	opts := client.StartWorkflowOptions{
		TaskQueue: TemporalTaskQueue,
	}
	if request.Fingerprint != "" {
		opts.SearchAttributes = map[string]interface{}{fingerprintSearchAttributeKey: request.Fingerprint}
	}

	wr, err := c.tc.ExecuteWorkflow(ctx, opts, RunQuickDiagnosticWorkflow, request)

	if err != nil {
		c.logger.Error("error executing workflow", "error", err)
		return nil, err
	}

	return &corepb.RunDiagnosticResponse{Id: wr.GetID()}, nil
}

func (c *Core) ListDiagnostics(ctx context.Context, req *corepb.ListDiagnosticsRequest) (*corepb.ListDiagnosticsResponse, error) {

	params := &workflowservice.ListWorkflowExecutionsRequest{
		Namespace: c.config.Temporal.Namespace,
		Query:     "transaction_id='" + req.Fingerprint + "'",
	}

	if req.Limit != 0 {
		params.PageSize = int32(req.Limit)
	}

	if req.Offset != 0 {
		// int64 to byte array
		params.NextPageToken = []byte(fmt.Sprintf("%d", req.Offset))
	}

	res, err := c.tc.ListWorkflow(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(res.Executions) == 0 {
		return &corepb.ListDiagnosticsResponse{
			Limit: req.Limit,
			Total: 0,
		}, nil
	}

	out := &corepb.ListDiagnosticsResponse{
		Diagnostics: []*analysispb.Report{},
		Offset:      0,
		Limit:       req.Limit,
		Total:       int64(len(res.Executions)),
	}

	for _, wfResult := range res.Executions {
		report, err := parseExecution(ctx, c.tc, wfResult)
		if err != nil {
			return nil, err
		}
		out.Diagnostics = append(out.Diagnostics, report)
	}

	return out, nil
}

func (c *Core) GetDiagnostic(ctx context.Context, request *corepb.GetDiagnosticRequest) (*analysispb.Report, error) {

	params := &workflowservice.ListWorkflowExecutionsRequest{
		Namespace: c.config.Temporal.Namespace,
		Query:     "WorkflowId='" + request.Id + "'",
	}

	res, err := c.tc.ListWorkflow(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(res.Executions) == 0 {
		return nil, fmt.Errorf("no workflow found")
	}

	wfResult := res.Executions[0]
	return parseExecution(ctx, c.tc, wfResult)

}

func parseExecution(ctx context.Context, tc client.Client, e *workflow.WorkflowExecutionInfo) (*analysispb.Report, error) {
	report := &analysispb.Report{}
	tid, ok := e.SearchAttributes.IndexedFields[fingerprintSearchAttributeKey]
	if ok {
		if tid.GetData() != nil {
			var fingerprint string
			err := json.Unmarshal(tid.GetData(), &fingerprint) // just to sanitize the data
			if err != nil {
				return nil, fmt.Errorf("could not marshal fingerprint: %w", err)
			}
			report.Fingerprint = &fingerprint
		}
	}

	switch e.Type.Name {
	case "RunDiagnosticWorkflow":
		report.Type = analysispb.Report_TYPE_DETAILED
	case "RunQuickDiagnosticWorkflow":
		report.Type = analysispb.Report_TYPE_QUICK
	default:
		report.Type = analysispb.Report_TYPE_NOT_SET
	}

	switch e.Status {
	case enums.WORKFLOW_EXECUTION_STATUS_RUNNING:
		report.Started = timestamppb.New(*e.StartTime)
		report.Status = analysispb.Report_STATUS_RUNNING

	case enums.WORKFLOW_EXECUTION_STATUS_COMPLETED:
		err := tc.GetWorkflow(ctx, e.Execution.WorkflowId, "").Get(ctx, &report)
		if err != nil {
			return nil, err
		}
		report.Started = timestamppb.New(*e.StartTime)
		report.Status = analysispb.Report_STATUS_COMPLETED
		report.Completed = timestamppb.New(*e.CloseTime)

	case enums.WORKFLOW_EXECUTION_STATUS_FAILED:
		err := tc.GetWorkflow(ctx, e.Execution.WorkflowId, "").Get(ctx, &report)
		if err != nil {
			msg := "the diagnostic failed, please try again or contact support for more information"
			report.Error = &msg
		}
		report.Started = timestamppb.New(*e.StartTime)
		report.Status = analysispb.Report_STATUS_FAILED
		report.Completed = timestamppb.New(*e.CloseTime)

	case enums.WORKFLOW_EXECUTION_STATUS_CANCELED:
		report.Started = timestamppb.New(*e.StartTime)
		report.Status = analysispb.Report_STATUS_CANCELED
		report.Completed = timestamppb.New(*e.CloseTime)

	case enums.WORKFLOW_EXECUTION_STATUS_TERMINATED:
		report.Started = timestamppb.New(*e.StartTime)
		report.Status = analysispb.Report_STATUS_TERMINATED
		report.Completed = timestamppb.New(*e.CloseTime)
	}

	return report, nil

}
