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
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"

	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.opentelco.io/go-swpx/proto/go/providerpb"
	"go.opentelco.io/go-swpx/shared"
)

var VERSION *version.Version
var logger hclog.Logger

const (
	VERSION_BASE string = "1.0-beta"
	SDD_VLAN     int64  = 296
)

var SDDNetwork *net.IPNet

var PROVIDER_NAME = "sait"

func init() {
	var err error
	if VERSION, err = version.NewVersion(VERSION_BASE); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// Provider is the implementation of the GRPC
type Provider struct {
	logger hclog.Logger
}

func (g *Provider) Version() (string, error) {
	return fmt.Sprintf("%s@%s", PROVIDER_NAME, VERSION.String()), nil
}

func (p *Provider) Name() (string, error) {
	return PROVIDER_NAME, nil
}

func (p *Provider) ProcessPollResponse(ctx context.Context, r *corepb.PollResponse) (*corepb.PollResponse, error) {
	changes := 0
	if r.Device == nil {
		p.logger.Warn("network element is empt ")
		return r, nil
	}

	// attach device
	// for ri, i := range r.Device.Interfaces {
	// 	for _, d := range i.DhcpTable {
	// 		if d.Vlan == SDD_VLAN {
	// 			p.logger.Debug("found SDD on interface", "interface", i.Description, "sdd-ipAddr", d.IpAddress, "sdd-mac", d.HardwareAddress)
	// 			r.Device.Interfaces[ri].ConnectedSdd = &networkelementpb.Element{BridgeMacAddress: d.HardwareAddress, Hostname: d.IpAddress}
	// 			changes++
	// 		}
	// 	}
	// }

	p.logger.Named("post-handler").Debug("processing response", "changes", changes)
	return r, nil
}

func (p *Provider) ResolveSessionRequest(ctx context.Context, request *corepb.SessionRequest) (*corepb.SessionRequest, error) {
	p.logger.Named("pre-handler").Debug("processing request")

	if request.AccessId != "" {
		access, ok := translationMap[request.AccessId]
		if ok {
			p.logger.Named("pre-handler").Info("found access ID on provider", "access_id", request.AccessId, "network_element", access.Device, "port", access.Interface)
			request.Hostname = access.Device
			request.Port = access.Interface

		} else {
			p.logger.Named("pre-handler").Warn("could not find access id on provider", "access_id", request.AccessId)
			return request, fmt.Errorf("access id was not found with selected provider")
		}
	}

	return request, nil
}
func (p *Provider) ResolveResourcePlugin(ctx context.Context, request *corepb.SessionRequest) (*providerpb.ResolveResourcePluginResponse, error) {
	resp := &providerpb.ResolveResourcePluginResponse{}

	if request.AccessId != "" {
		access, ok := translationMap[request.AccessId]
		if ok {
			resp.ResourcePlugin = access.ResourcePlugin
		}
	} else {
		var ip net.IP
		ip = net.ParseIP(request.Hostname)
		if ip != nil {
			addrs, err := net.LookupHost(request.Hostname)
			if err != nil || len(addrs) == 0 {
				p.logger.Error("could not find host", "host", request.Hostname, "error", err)
				return nil, fmt.Errorf("could not find host, %w", err)
			}
			ip = net.ParseIP(addrs[0])
		}
		if SDDNetwork.Contains(ip) {
			resp.ResourcePlugin = "raycore"
		}
	}
	return resp, nil
}

func main() {
	var err error
	logger = hclog.New(&hclog.LoggerOptions{
		Name:  fmt.Sprintf("%s@%s", PROVIDER_NAME, VERSION.String()),
		Level: hclog.Debug,
		Color: hclog.AutoColor,
	})

	_, SDDNetwork, err = net.ParseCIDR("192.168.112.0/23")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	prov := &Provider{logger: logger}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			shared.PluginProviderKey: &shared.ProviderPlugin{
				Impl: prov,
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
		Logger:     logger,
	})

}
