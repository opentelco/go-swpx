package api

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.opentelco.io/go-swpx/core"
	"go.opentelco.io/go-swpx/proto/go/analysispb"
	"go.opentelco.io/go-swpx/proto/go/corepb"
)

type coreGrpcImpl struct {
	core   *core.Core
	grpc   *grpc.Server
	logger hclog.Logger

	corepb.UnimplementedPollerServer
}

var automatedOkList = []string{
	"mulbarton-migration-a1",
	"mulbarton-migration-a2",
	"mulbarton-migration-a3",
	"only-for-migration-a1",
	"only-for-migration-a2",
	"only-for-migration-a3",
	"only-for-migration-a4",
	"only-for-migration-a5",
	"only-for-migration-a6",
	"only-for-migration-a7",
	"only-for-migration-a8",
	"only-for-migration-a9",
	"only-for-migration-a10",
}

func (s *coreGrpcImpl) Discover(ctx context.Context, request *corepb.DiscoverRequest) (*corepb.DiscoverResponse, error) {
	return s.core.Discover(ctx, request)
}

func (s *coreGrpcImpl) RunDiagnostic(ctx context.Context, request *corepb.RunDiagnosticRequest) (*corepb.RunDiagnosticResponse, error) {
	return s.core.RunDiagnostic(ctx, request)
}

func (s *coreGrpcImpl) RunQuickDiagnostic(ctx context.Context, request *corepb.RunDiagnosticRequest) (*corepb.RunDiagnosticResponse, error) {
	return s.core.RunQuickDiagnostic(ctx, request)
}

func (s *coreGrpcImpl) GetDiagnostic(ctx context.Context, request *corepb.GetDiagnosticRequest) (*analysispb.Report, error) {
	return s.core.GetDiagnostic(ctx, request)
}

func (s *coreGrpcImpl) ListDiagnostics(ctx context.Context, req *corepb.ListDiagnosticsRequest) (*corepb.ListDiagnosticsResponse, error) {

	return s.core.ListDiagnostics(ctx, req)
}

func (s *coreGrpcImpl) CollectDeviceInformation(context.Context, *corepb.CollectDeviceInformationRequest) (*corepb.DeviceInformationResponse, error) {
	return nil, nil
}

func (s *coreGrpcImpl) CollectBasicDeviceInformation(ctx context.Context, req *corepb.CollectBasicDeviceInformationRequest) (*corepb.DeviceInformationResponse, error) {
	return s.core.CollectBasicDeviceInformation(ctx, req)
}

// Request to SWP-core
func (s *coreGrpcImpl) Poll(ctx context.Context, request *corepb.PollRequest) (*corepb.PollResponse, error) {

	start := time.Now()
	if request.Type == corepb.PollRequest_NOT_SET {
		request.Type = corepb.PollRequest_GET_TECHNICAL_INFO
	}

	if request.Session.NetworkRegion == "" {
		return nil, status.Error(codes.InvalidArgument, "network_region is required")
	}

	resp, err := s.core.PollDevice(ctx, request)
	if err != nil {
		return nil, err
	}

	if resp == nil || resp.Device == nil {
		return nil, status.Error(codes.NotFound, "response is empty, no data from go-dnc")
	}
	resp.ExecutionTime = time.Since(start).String()

	return resp, nil
}

// Helper to get a .In behaviour
func In(hostname string, list ...string) bool {
	for _, item := range list {
		if item == hostname {
			return true
		}
	}
	return false
}

func (s *coreGrpcImpl) CollectConfig(ctx context.Context, request *corepb.CollectConfigRequest) (*corepb.CollectConfigResponse, error) {
	return s.core.CollectConfig(ctx, request)
}

func NewGrpc(core *core.Core, srv *grpc.Server, logger hclog.Logger) {
	instance := &coreGrpcImpl{
		core:   core,
		logger: logger,
	}
	corepb.RegisterPollerServer(srv, instance)
}
