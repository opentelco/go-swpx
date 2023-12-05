package api

import (
	"context"

	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.opentelco.io/go-swpx/proto/go/stanzapb"
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

func (s *commanderGrpcImpl) ConfigureStanza(ctx context.Context, req *corepb.ConfigureStanzaRequest) (*stanzapb.ConfigureResponse, error) {
	resp, err := s.service.ConfigureStanza(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
