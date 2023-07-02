package api

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"git.liero.se/opentelco/go-swpx/core"
	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
)

type coreGrpcImpl struct {
	core   *core.Core
	grpc   *grpc.Server
	logger hclog.Logger

	corepb.UnimplementedCoreServiceServer
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

// Request to SWP-core
func (s *coreGrpcImpl) Poll(ctx context.Context, request *corepb.PollRequest) (*corepb.PollResponse, error) {

	start := time.Now()
	if request.Type == corepb.PollRequest_NOT_SET {
		request.Type = corepb.PollRequest_GET_TECHNICAL_INFO
	}

	if request.Session.NetworkRegion == "" {
		return nil, status.Error(codes.InvalidArgument, "network_region is required")
	}

	resp, err := s.core.PollNetworkElement(ctx, request)
	if err != nil {
		return nil, err
	}

	if resp == nil || resp.NetworkElement == nil {
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

func (s *coreGrpcImpl) Command(ctx context.Context, request *corepb.CommandRequest) (*corepb.CommandResponse, error) {
	panic("implement me")
}

func (s *coreGrpcImpl) Information(ctx context.Context, request *emptypb.Empty) (*corepb.InformationResponse, error) {
	panic("implement me")
}

func NewGrpc(core *core.Core, srv *grpc.Server, logger hclog.Logger) {
	instance := &coreGrpcImpl{
		core:   core,
		logger: logger,
	}
	corepb.RegisterCoreServiceServer(srv, instance)
}
