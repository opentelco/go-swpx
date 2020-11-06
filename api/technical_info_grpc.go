package api

import (
	"context"
	"errors"
	"git.liero.se/opentelco/go-swpx/core"
	"git.liero.se/opentelco/go-swpx/proto/go/resource"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

type GRPCServer struct {
	requests chan *core.Request
	grpc     *grpc.Server
}

func (s *GRPCServer) TechnicalPortInformation(ctx context.Context, requestProto *resource.TechnicalInformationRequest) (*resource.TechnicalInformationResponse, error) {
	ctx, _ = context.WithTimeout(ctx, time.Duration(requestProto.Timeout)*time.Second)
	req := &core.Request{
		NetworkElement: requestProto.Hostname,
		Provider:       requestProto.Provider,
		Resource:       requestProto.Driver,
		DontUseIndex:   requestProto.RecreateIndex,

		Response: make(chan *resource.TechnicalInformationResponse, 1),
		Context:  ctx,
	}

	if requestProto.Port != "" {
		req.NetworkElementInterface = &requestProto.Port
		req.Type = core.GetTechnicalInformationPort
	} else {
		req.Type = core.GetTechnicalInformationElement
	}

	cachedResponse, err := core.ResponseCache.PopResponse(req.NetworkElement, *req.NetworkElementInterface, req.Type)
	if err != nil {
		logger.Error("error popping from cache: ", err.Error())
		return nil, err
	}

	if cachedResponse != nil {
		if time.Since(cachedResponse.Timestamp.AsTime()) < time.Duration(requestProto.CacheTtl)*time.Second {
			logger.Info("found response in cache")

			return cachedResponse.Response, nil
		}
		// if response is cached but ttl ran out, clear it from the cache
		if err := core.ResponseCache.Clear(req.NetworkElement, *req.NetworkElementInterface, req.Type); err != nil {
			logger.Error("error clearing cache:", err)
		}
	}

	s.requests <- req

	for {
		select {
		case resp := <-req.Response:
			if resp.Error != nil {
				return nil, errors.New(resp.Error.Message)
			}

			if err := core.ResponseCache.SetResponse(req.NetworkElement, *req.NetworkElementInterface, req.Type, resp); err != nil {
				logger.Error("error saving response to cache: ", err.Error())
				return nil, err
			}

			return resp, nil
		case <-req.Context.Done():
			logger.Info("timeout for request was hit")
			return nil, errors.New("timeout")
		}
	}
}

func NewGRPCServer(requests chan *core.Request) *GRPCServer {
	if requests == nil {
		log.Fatal("channel is nil, requests needs to be handled..")
	}

	logger = hclog.New(&hclog.LoggerOptions{
		Name:   APP_NAME,
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	grpcServer := grpc.NewServer()
	instance := &GRPCServer{requests, grpcServer}

	resource.RegisterTechnicalInformationServer(grpcServer, instance)

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
