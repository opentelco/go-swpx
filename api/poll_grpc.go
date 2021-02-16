package api

import (
	"context"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"

	"git.liero.se/opentelco/go-swpx/core"
	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
)

type GRPCServer struct {
	core   *core.Core
	grpc   *grpc.Server
	logger hclog.Logger
}

//  Request to SWP-core
func (s *GRPCServer) Poll(ctx context.Context, request *pb_core.Request) (*pb_core.Response, error) {

	request.Type = pb_core.Request_GET_TECHNICAL_INFO

	req := &core.Request{
		Request: request,
		// Metadata
		Response: make(chan *pb_core.Response, 1),
		Context:  ctx,
	}
	resp, err := s.core.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *GRPCServer) Command(ctx context.Context, request *pb_core.CommandRequest) (*pb_core.CommandResponse, error) {
	panic("implement me")
}

func (s *GRPCServer) Information(ctx context.Context, request *empty.Empty) (*pb_core.InformationResponse, error) {
	panic("implement me")
}

// func (s *GRPCServer) TechnicalPortInformation(ctx context.Context, requestProto *pb_core.Request) (*pb_core.Response, error) {
// ctx, _ = context.WithTimeout(ctx, time.Duration(requestProto.Timeout)*time.Second)
//
// req := &core.Request{
// 	Request: requestProto,
// 	Response: make(chan *pb_core.Response, 1),
// 	Context:  ctx,
// }
//
// if requestProto.Port != "" {
// 	req.NetworkElementInterface = &requestProto.Port
// 	req.Type = core.GetTechnicalInformationPort
// } else {
// 	req.Type = core.GetTechnicalInformationElement
// }
//
// cachedResponse, err := core.responseCache.PopResponse(req.NetworkElement, *req.NetworkElementInterface, req.Type)
// if err != nil {
// 	logger.Error("error popping from cache: ", err.Error())
// 	return nil, err
// }
//
// if cachedResponse != nil {
// 	if time.Since(cachedResponse.Timestamp.AsTime()) < time.Duration(requestProto.CacheTtl)*time.Second {
// 		logger.Info("found response in cache")
//
// 		return cachedResponse.Response, nil
// 	}
// 	// if response is cached but ttl ran out, clear it from the cache
// 	if err := core.responseCache.Clear(req.NetworkElement, *req.NetworkElementInterface, req.Type); err != nil {
// 		logger.Error("error clearing cache:", err)
// 	}
// }
//
// s.requests <- req
//
// for {
// 	select {
// 	case resp := <-req.Response:
// 		if resp.Error != nil {
// 			return nil, errors.New(resp.Error.Message)
// 		}
//
// 		if err := core.responseCache.SetResponse(req.Hostname, *req.NetworkElementInterface, req.Type, resp); err != nil {
// 			logger.Error("error saving response to cache: ", err.Error())
// 			return nil, err
// 		}
//
// 		return resp, nil
// 	case <-req.Context.Done():
// 		logger.Info("timeout for request was hit")
// 		return nil, errors.New("timeout")
// 	}
// }
// }

func NewGRPCServer(core *core.Core, logger hclog.Logger) *GRPCServer {

	grpcServer := grpc.NewServer()
	instance := &GRPCServer{core, grpcServer, logger}

	pb_core.RegisterCoreServer(grpcServer, instance)

	return instance
}

func (s GRPCServer) ListenAndServe(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	err = s.grpc.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
