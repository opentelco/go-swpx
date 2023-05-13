package api

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"git.liero.se/opentelco/go-swpx/core"
	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
)

type coreGrpcImpl struct {
	pb_core.UnimplementedCoreServer
	core   *core.Core
	grpc   *grpc.Server
	logger hclog.Logger
}

// Request to SWP-core
func (s *coreGrpcImpl) Poll(ctx context.Context, request *pb_core.Request) (*pb_core.Response, error) {

	start := time.Now()
	if request.Type == pb_core.Request_NOT_SET {
		request.Type = pb_core.Request_GET_TECHNICAL_INFO
	}

	if request.NetworkRegion == "" {
		return nil, status.Error(codes.InvalidArgument, "network_region is required")
	}

	resp, err := s.core.SendRequest(ctx, request)
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

func (s *coreGrpcImpl) Command(ctx context.Context, request *pb_core.CommandRequest) (*pb_core.CommandResponse, error) {
	panic("implement me")
}

func (s *coreGrpcImpl) Information(ctx context.Context, request *emptypb.Empty) (*pb_core.InformationResponse, error) {
	panic("implement me")
}

func NewCoreGrpcServer(core *core.Core, logger hclog.Logger) *coreGrpcImpl {

	grpcServer := grpc.NewServer()
	instance := &coreGrpcImpl{
		core:   core,
		grpc:   grpcServer,
		logger: logger,
	}

	pb_core.RegisterCoreServer(grpcServer, instance)

	return instance
}

func (s *coreGrpcImpl) ListenAndServe(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	s.logger.Info("starting grpc server", "addr", addr)
	err = s.grpc.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
