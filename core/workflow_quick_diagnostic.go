package core

import (
	"fmt"
	"time"

	"go.opentelco.io/go-swpx/proto/go/analysispb"
	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func RunQuickDiagnosticWorkflow(ctx workflow.Context, request *corepb.RunDiagnosticRequest) (*analysispb.Report, error) {

	logger := workflow.GetLogger(ctx)
	logger.Info("Diagnostic workflow started")
	defer logger.Info("Diagnostic workflow ended")

	var availabilityResp corepb.CheckAvailabilityResponse
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			TaskQueue:           TemporalTaskQueue,
			StartToCloseTimeout: time.Duration(60 * time.Second),
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 2,
			},
		}),
		coreActivities.CheckAvailability,
		request.Session,
	).Get(ctx, &availabilityResp); err != nil {
		return nil, err
	}

	var pollResponse corepb.PollResponse
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			TaskQueue:           TemporalTaskQueue,
			StartToCloseTimeout: time.Duration(80 * time.Second),
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 1,
			},
		}),
		coreActivities.Poll,
		&corepb.PollRequest{
			Session:  request.Session,
			Settings: request.Settings,
			Type:     corepb.PollRequest_GET_BASIC_INFO,
		},
	).Get(ctx, &pollResponse); err != nil {
		return nil, err
	}

	report := &analysispb.Report{}

	if pollResponse.Device == nil {
		return nil, fmt.Errorf("diagnosis report was missing from device")
	}

	var responses []*corepb.PollResponse
	responses = append(responses, &pollResponse)

	// analyze the link
	r, err := analyzeLink(request.Session.Port, responses)
	if err != nil {
		return nil, err
	}
	report.Analysis = append(report.Analysis, r...)

	// analyze the errors
	r, err = analyzeErrors(request.Session.Port, responses)
	if err != nil {
		return nil, err
	}
	report.Analysis = append(report.Analysis, r...)

	// analyze the traffic
	r, err = analyzeTraffic(request.Session.Port, responses)
	if err != nil {
		return nil, err
	}
	report.Analysis = append(report.Analysis, r...)

	// analyze the transceiver
	r, err = analyzeTransceiver(request.Session.Port, responses)
	if err != nil {
		return nil, err
	}
	report.Analysis = append(report.Analysis, r...)

	return report, nil
}
