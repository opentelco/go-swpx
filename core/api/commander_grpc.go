package api

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"google.golang.org/grpc"
)

type commanderGrpcImpl struct {
	service corepb.CommanderServiceServer

	corepb.UnimplementedCommanderServiceServer
}

func NewCommanderGrpc(service corepb.CommanderServiceServer, srv *grpc.Server) {
	impl := &commanderGrpcImpl{
		service: service,
	}

	corepb.RegisterCommanderServiceServer(srv, impl)

}

func (s *commanderGrpcImpl) ConfigureStanza(ctx context.Context, req *corepb.ConfigureStanzaRequest) (*corepb.ConfigureStanzaResponse, error) {
	resp, err := s.service.ConfigureStanza(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
