package core

import (
	"context"
	"time"

	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.temporal.io/sdk/temporal"
)

type Activities struct {
	core *Core
}

func NewActivities(c *Core) *Activities {
	return &Activities{core: c}
}

func (a *Activities) CheckAvailability(ctx context.Context, request *corepb.SessionRequest) (*corepb.CheckAvailabilityResponse, error) {
	return &corepb.CheckAvailabilityResponse{Available: true, ResponseTime: 40}, nil
}

func (a *Activities) Poll(ctx context.Context, request *corepb.PollRequest) (*corepb.PollResponse, error) {
	start := time.Now()
	if request.Type == corepb.PollRequest_NOT_SET {
		request.Type = corepb.PollRequest_GET_TECHNICAL_INFO
	}

	if request.Session.NetworkRegion == "" {
		return nil, temporal.NewNonRetryableApplicationError("network_region is required", TemporalErrTypeRegionRequired, nil)
	}

	resp, err := a.core.PollNetworkElement(ctx, request)

	if err != nil {
		return nil, err
	}

	if resp == nil || resp.NetworkElement == nil {
		return nil, temporal.NewNonRetryableApplicationError("failed to get data", TemporalErrTypePollFailed, nil)
	}
	resp.ExecutionTime = time.Since(start).String()

	return resp, err
}

func (a *Activities) Test() error {
	return nil
}
