package shared

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/networkelement"
	proto "git.liero.se/opentelco/go-swpx/proto/resource"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// NetworkElementPlugin is the interface that we're exposing as a plugin.

type Resource interface {
	Version() (string, error)
	// Gets all the technical information for a Port
	TechnicalPortInformation(context.Context, *proto.NetworkElement) (*networkelement.Element, error) // From interface name/descr a SNMP index must be found. This functions helps to solve this problem

	// TODO should return a slice of *proto.NetworkElementInterface so we can cache all results
	MapInterface(context.Context, *proto.NetworkElement) (*proto.NetworkElementInterfaces, error)
	MapEntityPhysical(context.Context, *proto.NetworkElement) (*proto.NetworkElementInterfaces, error)
	GetTransceiverInformation(ctx context.Context, ne *proto.NetworkElement) (*networkelement.Transceiver, error)

	SetConfiguration(ctx context.Context, conf Configuration) error
	GetConfiguration(ctx context.Context) (Configuration, error)
}

// Here is an implementation that talks over RPC
type ResourceGRPCClient struct {
	client proto.ResourceClient
	conf   Configuration
}

// MapInterface is the client implementation for the plugin-resource. Connects to the RPC
func (rpc *ResourceGRPCClient) MapInterface(ctx context.Context, proto *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	return rpc.client.MapInterface(ctx, proto)
}

// TechnicalPortInformation is the client implementation
func (rpc *ResourceGRPCClient) TechnicalPortInformation(ctx context.Context, proto *proto.NetworkElement) (*networkelement.Element, error) {
	return rpc.client.TechnicalPortInformation(ctx, proto)
}

func (rpc *ResourceGRPCClient) MapEntityPhysical(ctx context.Context, proto *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	return rpc.client.MapEntityPhysical(ctx, proto)
}

func (rpc *ResourceGRPCClient) GetTransceiverInformation(ctx context.Context, proto *proto.NetworkElement) (*networkelement.Transceiver, error) {
	return rpc.client.GetTransceiverInformation(ctx, proto)
}

func (rpc *ResourceGRPCClient) GetConfiguration(ctx context.Context) (Configuration, error) {
	return rpc.conf, nil
}

func (rpc *ResourceGRPCClient) SetConfiguration(ctx context.Context, conf Configuration) error {
	rpc.conf = conf

	return nil
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
	conf Configuration
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
func (rpc *ResourceGRCServer) MapInterface(ctx context.Context, ne *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	return rpc.Impl.MapInterface(ctx, ne)
}

func (rpc *ResourceGRCServer) MapEntityPhysical(ctx context.Context, ne *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	return rpc.Impl.MapEntityPhysical(ctx, ne)
}

func (rpc *ResourceGRCServer) GetConfiguration(ctx context.Context) (Configuration, error) {
	return rpc.conf, nil
}

func (rpc *ResourceGRCServer) GetTransceiverInformation(ctx context.Context, ne *proto.NetworkElement) (*networkelement.Transceiver, error) {
	return rpc.Impl.GetTransceiverInformation(ctx, ne)
}

func (rpc *ResourceGRCServer) SetConfiguration(ctx context.Context, conf Configuration) error {
	rpc.conf = conf

	return nil
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
