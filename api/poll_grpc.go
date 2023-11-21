package api

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"git.liero.se/opentelco/go-swpx/config"
	"git.liero.se/opentelco/go-swpx/core"
	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelement"
)

type coreGrpcImpl struct {
	pb_core.UnimplementedCoreServer

	conf   *config.Configuration
	core   *core.Core
	grpc   *grpc.Server
	logger hclog.Logger
}

func NewCoreGrpcServer(core *core.Core, conf *config.Configuration, logger hclog.Logger) *coreGrpcImpl {

	grpcServer := grpc.NewServer()
	instance := &coreGrpcImpl{
		core:   core,
		conf:   conf,
		grpc:   grpcServer,
		logger: logger,
	}

	if len(conf.Request.OverrideOkList) > 0 {
		logger.Warn("[override] automated OK enabled", "hostnames (prefix)", conf.Request.OverrideOkList)
	}

	pb_core.RegisterCoreServer(grpcServer, instance)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	return instance
}

// Request to SWP-core
func (s *coreGrpcImpl) Poll(ctx context.Context, request *pb_core.Request) (*pb_core.Response, error) {

	start := time.Now()
	if request.Type == pb_core.Request_NOT_SET {
		request.Type = pb_core.Request_GET_TECHNICAL_INFO
	}

	// if the switch is in the list, return OK.
	// noc has probably tried to migrate a dead switch
	if In(request.Hostname, s.conf.Request.OverrideOkList...) {
		s.logger.Warn("[override] automated OK", "hostname", request.Hostname, "port", request.Port, "type", request.Type.String())
		resp := &pb_core.Response{
			NetworkElement: &networkelement.Element{
				Hostname: request.Hostname,
				Interfaces: []*networkelement.Interface{
					{
						Alias:       request.Port,
						Description: request.Port,
					},
				},
			},
			Error:           nil,
			RequestAccessId: "",
			ExecutionTime:   "0s",
		}
		return resp, nil
	}

	resp, err := s.core.SendRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	if resp == nil || resp.NetworkElement == nil {
		return nil, status.Error(codes.NotFound, "no data from poller")
	}
	resp.ExecutionTime = time.Since(start).String()

	return resp, nil
}

// Helper to get a .In behaviour
func In(hostname string, list ...string) bool {
	for _, item := range list {
		if strings.HasPrefix(hostname, item) {
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
