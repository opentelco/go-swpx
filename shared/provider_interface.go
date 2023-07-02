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

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"git.liero.se/opentelco/go-swpx/proto/go/providerpb"
)

// Provider is the interface that we're exposing as a Provider plugin.
type Provider interface {
	Name() (string, error)
	Version() (string, error)

	// Process the request before it hits the main functions  in the Core
	ResolveSessionRequest(ctx context.Context, request *corepb.SessionRequest) (*corepb.SessionRequest, error)

	// Process the network element after data has been collected
	ProcessPollResponse(ctx context.Context, resp *corepb.PollResponse) (*corepb.PollResponse, error)

	ResolveResourcePlugin(ctx context.Context, request *corepb.SessionRequest) (*providerpb.ResolveResourcePluginResponse, error)
}

// Here is an implementation that talks over GRPC
type ProviderGRPCClient struct {
	client providerpb.ProviderClient
}

func (p *ProviderGRPCClient) Name() (string, error) {
	resp, err := p.client.Name(context.Background(), &emptypb.Empty{})
	return resp.Name, err
}

func (p *ProviderGRPCClient) Version() (string, error) {
	resp, err := p.client.Version(context.Background(), &emptypb.Empty{})
	if err != nil {
		return resp.Version, err
	}
	return resp.Version, err
}

func (p *ProviderGRPCClient) ResolveSessionRequest(ctx context.Context, req *corepb.SessionRequest) (*corepb.SessionRequest, error) {
	return p.client.ResolveSessionRequest(ctx, req)
}

func (p *ProviderGRPCClient) ProcessPollResponse(ctx context.Context, resp *corepb.PollResponse) (*corepb.PollResponse, error) {
	return p.client.ProcessPollResponse(ctx, resp)
}

func (p *ProviderGRPCClient) ResolveResourcePlugin(ctx context.Context, req *corepb.SessionRequest) (*providerpb.ResolveResourcePluginResponse, error) {
	return p.client.ResolveResourcePlugin(ctx, req)
}

// ProviderGRPCServer is the RPC server that ProviderPRC talks to, conforming to the requirements of net/rpc
type ProviderGRPCServer struct {
	providerpb.UnimplementedProviderServer
	Impl Provider
}

func (rpc *ProviderGRPCServer) Name(ctx context.Context, _ *emptypb.Empty) (*providerpb.NameResponse, error) {
	res, err := rpc.Impl.Name()
	return &providerpb.NameResponse{Name: res}, err
}

func (rpc *ProviderGRPCServer) Version(ctx context.Context, _ *emptypb.Empty) (*providerpb.VersionResponse, error) {
	res, err := rpc.Impl.Version()
	return &providerpb.VersionResponse{Version: res}, err
}

func (rpc *ProviderGRPCServer) ResolveSessionRequest(ctx context.Context, request *corepb.SessionRequest) (*corepb.SessionRequest, error) {
	return rpc.Impl.ResolveSessionRequest(ctx, request)
}

func (rpc *ProviderGRPCServer) ProcessPollResponse(ctx context.Context, resp *corepb.PollResponse) (*corepb.PollResponse, error) {
	return rpc.Impl.ProcessPollResponse(ctx, resp)
}

func (rpc *ProviderGRPCServer) ResolveResourcePlugin(ctx context.Context, req *corepb.SessionRequest) (*providerpb.ResolveResourcePluginResponse, error) {
	return rpc.Impl.ResolveResourcePlugin(ctx, req)
}

type ProviderPlugin struct {
	// Implement the plugin interface
	plugin.Plugin

	// Concrete implementaion for plugins written in Go.
	Impl Provider
}

func (p *ProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	providerpb.RegisterProviderServer(s, &ProviderGRPCServer{Impl: p.Impl})
	return nil
}

func (p *ProviderPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ProviderGRPCClient{client: providerpb.NewProviderClient(c)}, nil
}
