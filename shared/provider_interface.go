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
	"fmt"
	
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	
	"git.liero.se/opentelco/go-swpx/proto/go/core"
	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/proto/go/dnc"
	proto "git.liero.se/opentelco/go-swpx/proto/go/provider"
)

// NetworkElementPlugin is the interface that we're exposing as a plugin.
type Provider interface {
	Name() (string, error)
	Version() (string, error)

	Match(string) (bool, error)
	Lookup(string) (string, error)
	Weight() (int64, error)
	
	
	// Process the request before it hits the main functions  in the Core
	PreHandler(ctx context.Context, request *core.Request) (*core.Request, error)
	
	// Process the network element after data has been collected
	PostHandler(ctx context.Context,  response *core.Response) (*core.Response, error)
	

	GetConfiguration(ctx context.Context) (Configuration, error)
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

func (p *ProviderGRPCClient)  PreHandler(request *core.Request) (*core.Request, error) {
	return nil, nil
}

// Process the network element after data has been collected
func (p *ProviderGRPCClient)  PostHandler(resp *pb_core.Response) (*pb_core.Response,error) {
	return nil, nil

}

// GetConfiguration returns the configuration of the Provider.
func (p *ProviderGRPCClient) GetConfiguration(ctx context.Context) (Configuration, error) {
	resp, err := p.client.GetConfiguration(ctx, &proto.Empty{})
	if err != nil {
		return Configuration{}, err
	}

	if resp.Telnet == nil || resp.Ssh == nil || resp.SNMP == nil {
		return Configuration{}, fmt.Errorf("invalid config, restore to default")
	}

	return Proto2conf(resp), nil
}

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

func (rpc *ProviderGRPCServer) PreHandler(ctx context.Context, resp *pb_core.Request) (*pb_core.Request, error)  {
	return nil, nil
	
}

func (rpc *ProviderGRPCServer) PostHandler(ctx context.Context, resp *pb_core.Response)  (*pb_core.Response, error) {
	return nil,nil
}

func (rpc *ProviderGRPCServer) GetConfiguration(ctx context.Context,  _ *proto.Empty) (*proto.Configuration, error) {
	res, err := rpc.Impl.GetConfiguration(ctx)
	if err != nil {
		return nil, err
	}
	protoConf := Conf2proto(res)
	return &protoConf, err
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
