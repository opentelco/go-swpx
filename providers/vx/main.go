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

package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
	"go.vxfiber.dev/proto-go/inventory/device"
	"google.golang.org/grpc"

	"git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/proto/go/provider"
	"git.liero.se/opentelco/go-swpx/shared"
)

var VERSION *version.Version
var logger hclog.Logger

const (
	VERSION_BASE string = "1.0-beta"
)

var PROVIDER_NAME = "vx"

func init() {
	var err error
	if VERSION, err = version.NewVersion(VERSION_BASE); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger = hclog.New(&hclog.LoggerOptions{
		Name:  fmt.Sprintf("%s@%s", PROVIDER_NAME, VERSION.String()),
		Level: hclog.Debug,
		Color: hclog.AutoColor,
	})

}

// Provider is the implementation of the GRPC
type Provider struct {
	devClient device.ServiceClient
	logger    hclog.Logger
}

func (g *Provider) Version() (string, error) {
	return fmt.Sprintf("%s@%s", PROVIDER_NAME, VERSION.String()), nil
}

func (p *Provider) Name() (string, error) {
	return PROVIDER_NAME, nil
}
func (p *Provider) Setup(ctx context.Context, request *provider.SetupConfiguration) (*provider.SetupResponse, error) {

	return &provider.SetupResponse{}, nil
}

func (p *Provider) PostHandler(ctx context.Context, request *core.Response) (*core.Response, error) {
	p.logger.Named("post-handler").Debug("processing response", "changes", 0)
	return request, nil
}

func (p *Provider) PreHandler(ctx context.Context, request *core.Request) (*core.Request, error) {

	//  If s is not a valid textual representation of an IP address, ParseIP returns nil.
	isIp := net.ParseIP(request.Hostname)
	params := &device.GetParameters{}
	if isIp == nil {
		domain := parseDomain()
		p.logger.Info("appending domain to hostname", "hostname", request.Hostname, "domain", domain)
		params.Hostname = request.Hostname
		request.Hostname = fmt.Sprintf("%s%s", request.Hostname, domain)

	} else {
		params.Ip = isIp.String()
	}

	d, err := p.devClient.Get(ctx, &device.GetParameters{Hostname: request.Hostname})
	if err != nil || len(d.Devices) == 0 {
		p.logger.Warn("could not get OSS device", "hostname", params.Hostname, "error", err)
		return request, nil
	}
	host := d.Devices[0]
	switch strings.ToUpper(host.Vendor) {
	case "HUAWEI":
		p.logger.Debug("provider found device in oss, overwrite settings", "settings.resource_plugin", "vrp")
		request.Settings.ResourcePlugin = "vrp"
	case "CTC", "VXFIBER":
		p.logger.Debug("provider found device in oss, overwrite settings", "settings.resource_plugin", "ctc")
		request.Settings.ResourcePlugin = "ctc"
	}

	p.logger.Named("pre-handler").Debug("processing request in", "changes", 0)
	return request, nil
}

func (p *Provider) GetConfiguration(ctx context.Context) (*shared.Configuration, error) {
	return nil, nil
}

func main() {
	// when the config system in swpx is done this could be moved to that instead
	// for now this will be good enough
	grpcAddr := os.Getenv("OSS_INVENTORY_ADDR")
	if grpcAddr == "" {
		grpcAddr = "127.0.0.1:9001"
	}
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		logger.Error("could not connect to inventory GRPC")
		os.Exit(1)
	}
	devClient := device.NewServiceClient(conn)
	prov := &Provider{logger: logger, devClient: devClient}

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

func parseDomain() string {

	switch strings.ToUpper(os.Getenv("REGION")) {
	case "VX_SA1", "SA1":
		return ".joburg.net.venturanext.se"

	case "VX_UK1", "UK1":
		return ".net.uk1.vx.se"

	case "VX_DE1", "DE1":
		return ".net.de1.vx.se"

	case "VX_BE1", "BE1":
		return ".net.be1.vx.se"

	case "VX_AT1", "AT1":
		return ".net.at1.vx.se"

	case "VX_SE2", "SE2":
		return ".net.se2.vx.se"

	case "VX_SE1", "SE1":
		return ".net.se1.vx.se"

	default:
		return ""
	}

}
