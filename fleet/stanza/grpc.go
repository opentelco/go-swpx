package stanza

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"github.com/gogo/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewGRPC(service stanzapb.StanzaServiceServer, server *grpc.Server) {
	g := &grpcImpl{
		stanzaService: service,
		grpc:          server,
	}
	stanzapb.RegisterStanzaServiceServer(server, g)
}

type grpcImpl struct {
	grpc *grpc.Server

	stanzaService stanzapb.StanzaServiceServer

	stanzapb.UnimplementedStanzaServiceServer
}

func (g *grpcImpl) GetByID(ctx context.Context, params *stanzapb.GetByIDRequest) (*stanzapb.Stanza, error) {
	res, err := g.stanzaService.GetByID(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}
func (g *grpcImpl) List(ctx context.Context, params *stanzapb.ListRequest) (*stanzapb.ListResponse, error) {
	res, err := g.stanzaService.List(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}
func (g *grpcImpl) Create(ctx context.Context, params *stanzapb.CreateRequest) (*stanzapb.Stanza, error) {
	res, err := g.stanzaService.Create(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}
func (g *grpcImpl) Delete(ctx context.Context, params *stanzapb.DeleteRequest) (*emptypb.Empty, error) {
	res, err := g.stanzaService.Delete(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil
}

func (g *grpcImpl) Apply(ctx context.Context, params *stanzapb.ApplyRequest) (*stanzapb.ApplyResponse, error) {
	res, err := g.stanzaService.Apply(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

func (g *grpcImpl) Revert(ctx context.Context, params *stanzapb.RevertRequest) (*stanzapb.RevertResponse, error) {
	res, err := g.stanzaService.Revert(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil
}

func errHandler(err error) error {
	if err == nil {
		return nil
	}
	switch err {
	case ErrNotImplemented:
		return status.Errorf(codes.Unimplemented, err.Error())
	default:
		return status.Errorf(codes.Internal, err.Error())
	}
}
