/*
 * Copyright (c) 2020. Liero AB
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
	"git.liero.se/opentelco/go-swpx/proto/go/networkelement"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// NetworkElementPlugin is the interface that we're exposing as a plugin.

type Resource interface {
	Version() (string, error)

	// TechnicalPortInformation Gets all the technical information for a Port
	// from interface name/descr a SNMP index must be found. This functions helps to solve this problem
	TechnicalPortInformation(ctx context.Context, req *proto.Request) (*networkelement.Element, error)

	// BasicPortInformation
	BasicPortInformation(ctx context.Context, req *proto.Request) (*networkelement.Element, error)

	// AllPortInformation
	AllPortInformation(ctx context.Context, req *proto.Request) (*networkelement.Element, error)

	// MapInterface Map interfaces (IF-MIB) from device with the swpx model
	MapInterface(ctx context.Context, req *proto.Request) (*proto.PortIndex, error)

	// MapEntityPhysical Map interfcaes from Envirnment MIB to the swpx model
	MapEntityPhysical(ctx context.Context, req *proto.Request) (*proto.PortIndex, error)

	// GetTransceiverInformation Get SFP (transceiver) information
	GetTransceiverInformation(ctx context.Context, req *proto.Request) (*networkelement.Transceiver, error)

	// GetAllTransceiverInformation Maps transceivers to corresponding interfaces using physical port information in the wrapper
	GetAllTransceiverInformation(ctx context.Context, req *proto.Request) (*networkelement.Transceivers, error)

	GetRunningConfig(ctx context.Context, req *proto.GetRunningConfigParameters) (*proto.GetRunningConfigResponse, error)
}

// Here is an implementation that talks over RPC
type ResourceGRPCClient struct {
	client proto.ResourceClient
	conf   *config.Configuration
}

// MapInterface is the client implementation for the plugin-resource. Connects to the RPC
func (rpc *ResourceGRPCClient) MapInterface(ctx context.Context, req *proto.Request) (*proto.PortIndex, error) {
	return rpc.client.MapInterface(ctx, req)
}

func (rpc *ResourceGRPCClient) MapEntityPhysical(ctx context.Context, req *proto.Request) (*proto.PortIndex, error) {
	return rpc.client.MapEntityPhysical(ctx, req)
}

func (rpc *ResourceGRPCClient) AllPortInformation(ctx context.Context, req *proto.Request) (*networkelement.Element, error) {
	return rpc.client.AllPortInformation(ctx, req)
}

// TechnicalPortInformation is the client implementation
func (rpc *ResourceGRPCClient) TechnicalPortInformation(ctx context.Context, req *proto.Request) (*networkelement.Element, error) {
	return rpc.client.TechnicalPortInformation(ctx, req)
}

func (rpc *ResourceGRPCClient) BasicPortInformation(ctx context.Context, req *proto.Request) (*networkelement.Element, error) {
	return rpc.client.BasicPortInformation(ctx, req)
}

func (rpc *ResourceGRPCClient) GetTransceiverInformation(ctx context.Context, req *proto.Request) (*networkelement.Transceiver, error) {
	return rpc.client.GetTransceiverInformation(ctx, req)
}

func (rpc *ResourceGRPCClient) GetAllTransceiverInformation(ctx context.Context, req *proto.Request) (*networkelement.Transceivers, error) {
	return rpc.client.GetAllTransceiverInformation(ctx, req)
}

func (rpc *ResourceGRPCClient) GetRunningConfig(ctx context.Context, req *proto.GetRunningConfigParameters) (*proto.GetRunningConfigResponse, error) {
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
	proto.UnimplementedResourceServer
	// *plugin.MuxBroker
	Impl Resource
	conf *config.Configuration
}

// Version returns the current version
func (rpc *ResourceGRPCServer) Version(ctx context.Context, _ *emptypb.Empty) (*proto.VersionResponse, error) {
	res, err := rpc.Impl.Version()
	return &proto.VersionResponse{Version: res}, err
}

// MapInterface has the purppose to map interface name to a index by asking the device
func (rpc *ResourceGRPCServer) MapInterface(ctx context.Context, req *proto.Request) (*proto.PortIndex, error) {
	return rpc.Impl.MapInterface(ctx, req)
}

func (rpc *ResourceGRPCServer) MapEntityPhysical(ctx context.Context, req *proto.Request) (*proto.PortIndex, error) {
	return rpc.Impl.MapEntityPhysical(ctx, req)
}

// TechnicalPortInformation is a lazy interface to get all information needed for a technical info call.
func (rpc *ResourceGRPCServer) TechnicalPortInformation(ctx context.Context, req *proto.Request) (*networkelement.Element, error) {
	return rpc.Impl.TechnicalPortInformation(ctx, req)
}

// BasicPortInformation is a lazy interface to get all information needed for a technical info call.
func (rpc *ResourceGRPCServer) BasicPortInformation(ctx context.Context, req *proto.Request) (*networkelement.Element, error) {
	return rpc.Impl.BasicPortInformation(ctx, req)
}

// AllPortInformation is a lazy interface to get all information needed for a technical info call.
func (rpc *ResourceGRPCServer) AllPortInformation(ctx context.Context, req *proto.Request) (*networkelement.Element, error) {
	return rpc.Impl.AllPortInformation(ctx, req)
}

func (rpc *ResourceGRPCServer) GetTransceiverInformation(ctx context.Context, req *proto.Request) (*networkelement.Transceiver, error) {
	return rpc.Impl.GetTransceiverInformation(ctx, req)
}

func (rpc *ResourceGRPCServer) GetAllTransceiverInformation(ctx context.Context, req *proto.Request) (*networkelement.Transceivers, error) {
	return rpc.Impl.GetAllTransceiverInformation(ctx, req)
}

func (rpc *ResourceGRPCServer) GetRunningConfig(ctx context.Context, req *proto.GetRunningConfigParameters) (*proto.GetRunningConfigResponse, error) {
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
	proto.RegisterResourceServer(s, &ResourceGRPCServer{
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
