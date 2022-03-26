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
	TechnicalPortInformation(context.Context, *proto.NetworkElement) (*networkelement.Element, error)

	// BasicPortInformation
	BasicPortInformation(context.Context, *proto.NetworkElement) (*networkelement.Element, error)

	// AllPortInformation
	AllPortInformation(context.Context, *proto.NetworkElement) (*networkelement.Element, error)

	// MapInterface Map interfaces (IF-MIB) from device with the swpx model
	MapInterface(context.Context, *proto.NetworkElement) (*proto.NetworkElementInterfaces, error)

	// MapEntityPhysical Map interfcaes from Envirnment MIB to the swpx model
	MapEntityPhysical(context.Context, *proto.NetworkElement) (*proto.NetworkElementInterfaces, error)

	// GetTransceiverInformation Get SFP (transceiver) information
	GetTransceiverInformation(ctx context.Context, ne *proto.NetworkElement) (*networkelement.Transceiver, error)

	// GetAllTransceiverInformation Maps transceivers to corresponding interfaces using physical port information in the wrapper
	GetAllTransceiverInformation(ctx context.Context, ne *proto.NetworkElementWrapper) (*networkelement.Element, error)

	// SetConfiguration sets config in the resource plugin
	SetConfiguration(ctx context.Context, conf *Configuration) error

	// GetConfiguration configuration from provider
	GetConfiguration(ctx context.Context) (*Configuration, error)
}

// Here is an implementation that talks over RPC
type ResourceGRPCClient struct {
	client proto.ResourceClient
	conf   *Configuration
}

// MapInterface is the client implementation for the plugin-resource. Connects to the RPC
func (rpc *ResourceGRPCClient) MapInterface(ctx context.Context, proto *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	return rpc.client.MapInterface(ctx, proto)
}

func (rpc *ResourceGRPCClient) AllPortInformation(ctx context.Context, proto *proto.NetworkElement) (*networkelement.Element, error) {
	return rpc.client.AllPortInformation(ctx, proto)
}

// TechnicalPortInformation is the client implementation
func (rpc *ResourceGRPCClient) TechnicalPortInformation(ctx context.Context, proto *proto.NetworkElement) (*networkelement.Element, error) {
	return rpc.client.TechnicalPortInformation(ctx, proto)
}

func (rpc *ResourceGRPCClient) BasicPortInformation(ctx context.Context, proto *proto.NetworkElement) (*networkelement.Element, error) {
	return rpc.client.BasicPortInformation(ctx, proto)
}

func (rpc *ResourceGRPCClient) MapEntityPhysical(ctx context.Context, proto *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	return rpc.client.MapEntityPhysical(ctx, proto)
}

func (rpc *ResourceGRPCClient) GetTransceiverInformation(ctx context.Context, proto *proto.NetworkElement) (*networkelement.Transceiver, error) {
	return rpc.client.GetTransceiverInformation(ctx, proto)
}

func (rpc *ResourceGRPCClient) GetConfiguration(ctx context.Context) (*Configuration, error) {
	return rpc.conf, nil
}

func (rpc *ResourceGRPCClient) SetConfiguration(ctx context.Context, conf *Configuration) error {
	rpc.conf = conf

	return nil
}

func (rpc *ResourceGRPCClient) GetAllTransceiverInformation(ctx context.Context, ne *proto.NetworkElementWrapper) (*networkelement.Element, error) {
	return rpc.client.GetAllTransceiverInformation(ctx, ne)
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
	conf *Configuration
}

// Version returns the current version
func (rpc *ResourceGRPCServer) Version(ctx context.Context, _ *emptypb.Empty) (*proto.VersionResponse, error) {
	res, err := rpc.Impl.Version()
	return &proto.VersionResponse{Version: res}, err
}

// TechnicalPortInformation is a lazy interface to get all information needed for a technical info call.
func (rpc *ResourceGRPCServer) TechnicalPortInformation(ctx context.Context, ne *proto.NetworkElement) (*networkelement.Element, error) {
	return rpc.Impl.TechnicalPortInformation(ctx, ne)
}

// BasicPortInformation is a lazy interface to get all information needed for a technical info call.
func (rpc *ResourceGRPCServer) BasicPortInformation(ctx context.Context, ne *proto.NetworkElement) (*networkelement.Element, error) {
	return rpc.Impl.BasicPortInformation(ctx, ne)
}

// AllPortInformation is a lazy interface to get all information needed for a technical info call.
func (rpc *ResourceGRPCServer) AllPortInformation(ctx context.Context, ne *proto.NetworkElement) (*networkelement.Element, error) {
	return rpc.Impl.AllPortInformation(ctx, ne)
}

//MapInterface has the purppose to map interface name to a index by asking the device
func (rpc *ResourceGRPCServer) MapInterface(ctx context.Context, ne *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	return rpc.Impl.MapInterface(ctx, ne)
}

func (rpc *ResourceGRPCServer) MapEntityPhysical(ctx context.Context, ne *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	return rpc.Impl.MapEntityPhysical(ctx, ne)
}

func (rpc *ResourceGRPCServer) GetConfiguration(ctx context.Context) (*Configuration, error) {
	return rpc.conf, nil
}

func (rpc *ResourceGRPCServer) GetTransceiverInformation(ctx context.Context, ne *proto.NetworkElement) (*networkelement.Transceiver, error) {
	return rpc.Impl.GetTransceiverInformation(ctx, ne)
}

func (rpc *ResourceGRPCServer) SetConfiguration(ctx context.Context, conf *Configuration) error {
	rpc.conf = conf

	return nil
}
func (rpc *ResourceGRPCServer) GetAllTransceiverInformation(ctx context.Context, wrapper *proto.NetworkElementWrapper) (*networkelement.Element, error) {
	return rpc.Impl.GetAllTransceiverInformation(ctx, wrapper)
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
