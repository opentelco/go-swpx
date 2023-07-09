/*
 * Copyright (c) 2023. Liero AB
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software
 * is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
 * IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
 * CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package shared

import (
	"context"

	"git.liero.se/opentelco/go-swpx/config"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelementpb"
	"git.liero.se/opentelco/go-swpx/proto/go/resourcepb"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// NetworkElementPlugin is the interface that we're exposing as a plugin.

type Resource interface {
	// return the version of the plugin
	Version() (string, error)

	// discover the network element, simple lookup on the device
	Discover(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error)

	// TechnicalPortInformation Gets all the technical information for a Port
	// from interface name/descr a SNMP index must be found. This functions helps to solve this problem
	TechnicalPortInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error)

	// BasicPortInformation
	BasicPortInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error)

	// AllPortInformation
	AllPortInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error)

	// MapInterface Map interfaces (IF-MIB) from device with the swpx model
	MapInterface(ctx context.Context, req *resourcepb.Request) (*resourcepb.PortIndex, error)

	// MapEntityPhysical Map interfcaes from Envirnment MIB to the swpx model
	MapEntityPhysical(ctx context.Context, req *resourcepb.Request) (*resourcepb.PortIndex, error)

	// GetTransceiverInformation Get SFP (transceiver) information
	GetTransceiverInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Transceiver, error)

	// GetAllTransceiverInformation Maps transceivers to corresponding interfaces using physical port information in the wrapper
	GetAllTransceiverInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Transceivers, error)

	GetRunningConfig(ctx context.Context, req *resourcepb.GetRunningConfigParameters) (*resourcepb.GetRunningConfigResponse, error)
}

// Here is an implementation that talks over RPC
type ResourceGRPCClient struct {
	client resourcepb.ResourceClient
	conf   *config.Configuration
}

func (rpc *ResourceGRPCClient) Discover(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error) {
	return rpc.client.Discover(ctx, req)
}

// MapInterface is the client implementation for the plugin-resource. Connects to the RPC
func (rpc *ResourceGRPCClient) MapInterface(ctx context.Context, req *resourcepb.Request) (*resourcepb.PortIndex, error) {
	return rpc.client.MapInterface(ctx, req)
}

func (rpc *ResourceGRPCClient) MapEntityPhysical(ctx context.Context, req *resourcepb.Request) (*resourcepb.PortIndex, error) {
	return rpc.client.MapEntityPhysical(ctx, req)
}

func (rpc *ResourceGRPCClient) AllPortInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error) {
	return rpc.client.AllPortInformation(ctx, req)
}

// TechnicalPortInformation is the client implementation
func (rpc *ResourceGRPCClient) TechnicalPortInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error) {
	return rpc.client.TechnicalPortInformation(ctx, req)
}

func (rpc *ResourceGRPCClient) BasicPortInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error) {
	return rpc.client.BasicPortInformation(ctx, req)
}

func (rpc *ResourceGRPCClient) GetTransceiverInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Transceiver, error) {
	return rpc.client.GetTransceiverInformation(ctx, req)
}

func (rpc *ResourceGRPCClient) GetAllTransceiverInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Transceivers, error) {
	return rpc.client.GetAllTransceiverInformation(ctx, req)
}

func (rpc *ResourceGRPCClient) GetRunningConfig(ctx context.Context, req *resourcepb.GetRunningConfigParameters) (*resourcepb.GetRunningConfigResponse, error) {
	return rpc.client.GetRunningConfig(ctx, req)
}

func (rpc *ResourceGRPCClient) Version() (string, error) {
	resp, err := rpc.client.Version(context.Background(), &emptypb.Empty{})
	if err != nil {
		return "", err
	}

	return resp.Version, err
}

// ResourceGRPCServer is the server struct
type ResourceGRPCServer struct {
	resourcepb.UnimplementedResourceServer
	// *plugin.MuxBroker
	Impl Resource
	conf *config.Configuration
}

// Version returns the current version
func (rpc *ResourceGRPCServer) Version(ctx context.Context, _ *emptypb.Empty) (*resourcepb.VersionResponse, error) {
	res, err := rpc.Impl.Version()
	return &resourcepb.VersionResponse{Version: res}, err
}

func (rpc *ResourceGRPCServer) Discover(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error) {
	return rpc.Impl.Discover(ctx, req)
}

// MapInterface has the purppose to map interface name to a index by asking the device
func (rpc *ResourceGRPCServer) MapInterface(ctx context.Context, req *resourcepb.Request) (*resourcepb.PortIndex, error) {
	return rpc.Impl.MapInterface(ctx, req)
}

func (rpc *ResourceGRPCServer) MapEntityPhysical(ctx context.Context, req *resourcepb.Request) (*resourcepb.PortIndex, error) {
	return rpc.Impl.MapEntityPhysical(ctx, req)
}

// TechnicalPortInformation is a lazy interface to get all information needed for a technical info call.
func (rpc *ResourceGRPCServer) TechnicalPortInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error) {
	return rpc.Impl.TechnicalPortInformation(ctx, req)
}

// BasicPortInformation is a lazy interface to get all information needed for a technical info call.
func (rpc *ResourceGRPCServer) BasicPortInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error) {
	return rpc.Impl.BasicPortInformation(ctx, req)
}

// AllPortInformation is a lazy interface to get all information needed for a technical info call.
func (rpc *ResourceGRPCServer) AllPortInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error) {
	return rpc.Impl.AllPortInformation(ctx, req)
}

func (rpc *ResourceGRPCServer) GetTransceiverInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Transceiver, error) {
	return rpc.Impl.GetTransceiverInformation(ctx, req)
}

func (rpc *ResourceGRPCServer) GetAllTransceiverInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Transceivers, error) {
	return rpc.Impl.GetAllTransceiverInformation(ctx, req)
}

func (rpc *ResourceGRPCServer) GetRunningConfig(ctx context.Context, req *resourcepb.GetRunningConfigParameters) (*resourcepb.GetRunningConfigResponse, error) {
	return rpc.Impl.GetRunningConfig(ctx, req)
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
	resourcepb.RegisterResourceServer(s, &ResourceGRPCServer{
		Impl: p.Impl,
	})

	return nil
}

// GRPCClient implements the grpc client
func (p *ResourcePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ResourceGRPCClient{
		client: resourcepb.NewResourceClient(c),
	}, nil
}

// Secure that a plugin is implemented
var _ plugin.GRPCPlugin = &ResourcePlugin{}
