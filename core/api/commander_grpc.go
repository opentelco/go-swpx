package api

import (
	"context"

	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.opentelco.io/go-swpx/proto/go/stanzapb"
	"google.golang.org/grpc"
)

type commanderGrpcImpl struct {
	service corepb.CommanderServer

	corepb.UnimplementedCommanderServer
}

func NewCommanderGrpc(service corepb.CommanderServer, srv *grpc.Server) {
	impl := &commanderGrpcImpl{
		service: service,
	}

	corepb.RegisterCommanderServer(srv, impl)

}

func (s *commanderGrpcImpl) ConfigureStanza(ctx context.Context, req *corepb.ConfigureStanzaRequest) (*stanzapb.ConfigureResponse, error) {
	resp, err := s.service.ConfigureStanza(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
