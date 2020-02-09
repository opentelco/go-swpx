package shared

import (
	"context"

	"github.com/hashicorp/go-plugin"
	dnc "github.com/opentelco/go-swpx/proto/dnc"
	proto "github.com/opentelco/go-swpx/proto/provider"
	"google.golang.org/grpc"
)

// NetworkElementPlugin is the interface that we're exposing as a plugin.
type Provider interface {
	Name() (string, error)
	Version() (string, error)

	Match(string) (bool, error)
	Lookup(string) (string, error)
	Weight() (int64, error)
}

// Here is an implementation that talks over GRPC
type ProviderGRPCClient struct{ client proto.ProviderClient }

func (p *ProviderGRPCClient) Name() (string, error) {
	resp, err := p.client.Name(context.Background(), &proto.Empty{})
	return resp.Name, err
}

func (p *ProviderGRPCClient) Version() (string, error) {
	resp, err := p.client.Version(context.Background(), &proto.Empty{})
	if err != nil {
		return resp.Version, err
	}

	return resp.Version, err
}

// Match is used to match the ID with the provider.
func (p *ProviderGRPCClient) Match(id string) (bool, error) {
	resp, err := p.client.Match(context.Background(), &proto.MatchRequest{OriginId: id})
	return resp.Match, err
}

// Lookup returns a Request after doing a lookup against the provider.
func (p *ProviderGRPCClient) Lookup(id string) (string, error) {
	res, err := p.client.Lookup(context.Background(), &dnc.LookupRequest{OriginId: id})

	return res.OriginId, err
}

// Weight returns the weight of the Provider. Lower to higher.
func (p *ProviderGRPCClient) Weight() (int64, error) {
	resp, err := p.client.Weight(context.Background(), &proto.Empty{})
	return resp.GetWeight(), err

}

//
// ProviderGRPCServer is the RPC server that ProviderPRC talks to, conforming to the requirements of net/rpc
type ProviderGRPCServer struct {
	Impl Provider
}

func (rpc *ProviderGRPCServer) Name(ctx context.Context, _ *proto.Empty) (*proto.NameResponse, error) {
	res, err := rpc.Impl.Name()
	return &proto.NameResponse{Name: res}, err
}

func (rpc *ProviderGRPCServer) Version(ctx context.Context, _ *proto.Empty) (*proto.VersionResponse, error) {
	res, err := rpc.Impl.Version()
	return &proto.VersionResponse{Version: res}, err
}

func (rpc *ProviderGRPCServer) Lookup(ctx context.Context, req *dnc.LookupRequest) (*dnc.LookupResponse, error) {
	res, err := rpc.Impl.Lookup(req.OriginId)
	return &dnc.LookupResponse{
		OriginId: res,
	}, err
}

func (rpc *ProviderGRPCServer) Match(ctx context.Context, req *proto.MatchRequest) (*proto.MatchResponse, error) {
	res, err := rpc.Impl.Match(req.OriginId)
	return &proto.MatchResponse{Match: res}, err
}

func (rpc *ProviderGRPCServer) Weight(ctx context.Context, _ *proto.Empty) (*proto.WeightResponse, error) {
	res, err := rpc.Impl.Weight()
	return &proto.WeightResponse{Weight: res}, err
}

type ProviderPlugin struct {
	// Implement the plugin interface
	plugin.Plugin

	// Concrete implementaion for plugins written in Go.
	Impl Provider
}

func (p *ProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterProviderServer(s, &ProviderGRPCServer{Impl: p.Impl})
	return nil
}

func (p *ProviderPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ProviderGRPCClient{client: proto.NewProviderClient(c)}, nil
}
