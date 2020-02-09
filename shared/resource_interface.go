package shared

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/opentelco/go-swpx/proto/networkelement"
	proto "github.com/opentelco/go-swpx/proto/resource"
	"google.golang.org/grpc"
)

// NetworkElementPlugin is the interface that we're exposing as a plugin.
type Resource interface {
	Version() (string, error)
	TechnicalPortInformation(context.Context, *proto.NetworkElement) (*networkelement.Element, error)
	MapInterface(context.Context, *proto.NetworkElement) (*proto.NetworkElementInterface, error)
}

// Here is an implementation that talks over RPC
type ResourceGRPCClient struct {
	client proto.ResourceClient
}

// MapInterface is the client implmentation for the plugin-reosource. Connects to the RPC
func (rpc *ResourceGRPCClient) MapInterface(ctx context.Context, proto *proto.NetworkElement) (*proto.NetworkElementInterface, error) {
	return rpc.client.MapInterface(ctx, proto)
}

// TechnicalPortInformation is the client implmentation
func (rpc *ResourceGRPCClient) TechnicalPortInformation(ctx context.Context, proto *proto.NetworkElement) (*networkelement.Element, error) {
	return rpc.client.TechnicalPortInformation(ctx, proto)
}

func (rpc *ResourceGRPCClient) Version() (string, error) {
	resp, err := rpc.client.Version(context.Background(), &proto.Empty{})
	if err != nil {
		return resp.Version, err
	}

	return resp.Version, err
}

// ResourceGRCServer is the server struct
type ResourceGRCServer struct {
	// *plugin.MuxBroker
	Impl Resource
}

// Version returns the current version
func (rpc *ResourceGRCServer) Version(ctx context.Context, _ *proto.Empty) (*proto.VersionResponse, error) {
	res, err := rpc.Impl.Version()
	return &proto.VersionResponse{Version: res}, err
}

// TechnicalPortInformation is a lazy interface to get all information needed for a technical info call.
func (rpc *ResourceGRCServer) TechnicalPortInformation(ctx context.Context, ne *proto.NetworkElement) (*networkelement.Element, error) {
	return rpc.Impl.TechnicalPortInformation(ctx, ne)
}

//MapInterface has the purppose to map interface name to a index by asking the device
func (rpc *ResourceGRCServer) MapInterface(ctx context.Context, ne *proto.NetworkElement) (*proto.NetworkElementInterface, error) {
	return rpc.Impl.MapInterface(ctx, ne)
}

// ResourcePlugin is the implementation of plugin.Plugin so we can serve/consume this
type ResourcePlugin struct {
	// implement the plugin interface
	plugin.Plugin

	// concrete implementiation for plugins written in go.
	Impl Resource
}

// GRPCServer Implements RCP interface
func (p *ResourcePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterResourceServer(s, &ResourceGRCServer{
		Impl: p.Impl,
	})

	return nil
}

// GRPCClient implements the grpc client
func (p *ResourcePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ResourceGRPCClient{
		client: proto.NewResourceClient(c),
	}, nil
}

// Secure that a plugin is implemented
var _ plugin.GRPCPlugin = &ResourcePlugin{}
