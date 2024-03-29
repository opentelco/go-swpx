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

package main

import (
	"context"
	"fmt"
	"log"

	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.opentelco.io/go-swpx/proto/go/providerpb"
	"go.opentelco.io/go-swpx/shared"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
)

var VERSION *version.Version
var logger hclog.Logger

const (
	VERSION_BASE  string = "1.0-beta"
	PROVIDER_NAME string = "default"
	WEIGHT        int64  = 1
)

func init() {
	var err error
	if VERSION, err = version.NewVersion(VERSION_BASE); err != nil {
		log.Fatal(err)
	}
}

// Here is a real implementation of Greeter
type Provider struct {
	logger hclog.Logger
}

func (g *Provider) Version() (string, error) {
	return fmt.Sprintf("%s@%s", PROVIDER_NAME, VERSION.String()), nil
}

func (p *Provider) Name() (string, error) {
	return PROVIDER_NAME, nil

}
func (p *Provider) ResolveSessionRequest(ctx context.Context, request *corepb.SessionRequest) (*corepb.SessionRequest, error) {
	return request, nil
}

func (p *Provider) ResolveResourcePlugin(ctx context.Context, request *corepb.SessionRequest) (*providerpb.ResolveResourcePluginResponse, error) {
	return &providerpb.ResolveResourcePluginResponse{}, nil
}

func (p *Provider) ProcessPollResponse(ctx context.Context, resp *corepb.PollResponse) (*corepb.PollResponse, error) {
	return resp, nil
}

// func (p *Provider)  ResolveSessionRequest(ctx context.Context, request *pb_corepb.SessionRequest) (*pb_corepb.SessionRequest, error) {return nil,nil}
// func (p *Provider)  PostHandler(ctx context.Context, response *corepb.Response) (*corepb.Response, error) {return nil,nil}

func main() {
	logger = hclog.New(&hclog.LoggerOptions{
		Name:  fmt.Sprintf("%s@%s", PROVIDER_NAME, VERSION.String()),
		Level: hclog.Trace,
	})

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			shared.PluginProviderKey: &shared.ProviderPlugin{
				Impl: &Provider{logger: logger},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
