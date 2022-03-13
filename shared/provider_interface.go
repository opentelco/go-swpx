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

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"git.liero.se/opentelco/go-swpx/proto/go/core"
	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
	proto "git.liero.se/opentelco/go-swpx/proto/go/provider"
)

// NetworkElementPlugin is the interface that we're exposing as a plugin.
type Provider interface {
	Name() (string, error)
	Version() (string, error)

	// Process the request before it hits the main functions  in the Core
	PreHandler(ctx context.Context, request *pb_core.Request) (*pb_core.Request, error)

	// Process the network element after data has been collected
	PostHandler(ctx context.Context, response *core.Response) (*pb_core.Response, error)

	// Setup and initiate configuration and other things
	Setup(ctx context.Context, request *proto.SetupConfiguration) (*proto.SetupResponse, error)
}

// Here is an implementation that talks over GRPC
type ProviderGRPCClient struct {
	client proto.ProviderClient
}

func (p *ProviderGRPCClient) Name() (string, error) {
	resp, err := p.client.Name(context.Background(), &empty.Empty{})
	return resp.Name, err
}

func (p *ProviderGRPCClient) Version() (string, error) {
	resp, err := p.client.Version(context.Background(), &empty.Empty{})
	if err != nil {
		return resp.Version, err
	}
	return resp.Version, err
}

func (p *ProviderGRPCClient) PreHandler(ctx context.Context, req *pb_core.Request) (*pb_core.Request, error) {
	return p.client.PreHandler(ctx, req)
}

func (p *ProviderGRPCClient) PostHandler(ctx context.Context, resp *pb_core.Response) (*pb_core.Response, error) {
	return p.client.PostHandler(ctx, resp)
}

func (rpc *ProviderGRPCClient) Setup(ctx context.Context, r *proto.SetupConfiguration) (*proto.SetupResponse, error) {
	return rpc.client.Setup(ctx, r)
}

// func (p *ProviderGRPCClient)  PreHandler(request *core.Request) (*core.Request, error) {
// 	return p.client.PreHandler(context.Background(), request)
// }
//
// // Process the network element after data has been collected
// func (p *ProviderGRPCClient)  PostHandler(resp *pb_core.Response) (*pb_core.Response,error) {
// 	return p.client.PostHandler(context.Background(), resp)
// }

// ProviderGRPCServer is the RPC server that ProviderPRC talks to, conforming to the requirements of net/rpc
type ProviderGRPCServer struct {
	proto.UnimplementedProviderServer
	Impl Provider
}

func (rpc *ProviderGRPCServer) Name(ctx context.Context, _ *empty.Empty) (*proto.NameResponse, error) {
	res, err := rpc.Impl.Name()
	return &proto.NameResponse{Name: res}, err
}

func (rpc *ProviderGRPCServer) Version(ctx context.Context, _ *empty.Empty) (*proto.VersionResponse, error) {
	res, err := rpc.Impl.Version()
	return &proto.VersionResponse{Version: res}, err
}

func (rpc *ProviderGRPCServer) PreHandler(ctx context.Context, r *pb_core.Request) (*pb_core.Request, error) {
	return rpc.Impl.PreHandler(ctx, r)
}

func (rpc *ProviderGRPCServer) PostHandler(ctx context.Context, resp *pb_core.Response) (*pb_core.Response, error) {
	return rpc.Impl.PostHandler(ctx, resp)
}

func (rpc *ProviderGRPCServer) Setup(ctx context.Context, r *proto.SetupConfiguration) (*proto.SetupResponse, error) {
	return rpc.Impl.Setup(ctx, r)
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
