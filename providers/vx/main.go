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
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
	"go.vxfiber.dev/proto-go/inventory/device"
	"go.vxfiber.dev/proto-go/region"
	"go.vxfiber.dev/vx-bouncer/iam/go-sdk"
	"go.vxfiber.dev/vx-bouncer/iam/go-sdk/appauth"
	"go.vxfiber.dev/vx-bouncer/iam/go-sdk/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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
	deviceClients map[string]device.ServiceClient
	logger        hclog.Logger
	appToken      string
	appRegion     region.Region
}

func (g *Provider) Version() (string, error) {
	return fmt.Sprintf("%s@%s", PROVIDER_NAME, VERSION.String()), nil
}

func (p *Provider) Name() (string, error) {
	return PROVIDER_NAME, nil
}

func (p *Provider) ProcessPollResponse(ctx context.Context, request *core.PollResponse) (*core.PollResponse, error) {
	p.logger.Named("post.ProcessPollResponse").Debug("processing response", "changes", 0)
	return request, nil
}

func (p *Provider) ResolveSessionRequest(ctx context.Context, request *core.SessionRequest) (*core.SessionRequest, error) {
	countChanges := 0
	//  If s is not a valid textual representation of an IP address, ParseIP returns nil.

	isIp := net.ParseIP(request.Hostname)
	region := p.parseRegion(request.NetworkRegion)
	if region == nil {
		return request, fmt.Errorf("could not parse region from network region '%s'", request.NetworkRegion)
	}

	if isIp == nil {
		p.logger.Info("appending domain to hostname", "hostname", request.Hostname, "domain", region.domain)
		request.Hostname = fmt.Sprintf("%s%s", request.Hostname, region.domain)
	}

	p.logger.Named("pre.ResolveSessionRequest").Debug("processing request in", "changes", countChanges)
	return request, nil

}

func (p *Provider) ResolveResourcePlugin(ctx context.Context, request *core.SessionRequest) (*provider.ResolveResourcePluginResponse, error) {
	ctx = sdk.WithToken(ctx, p.appToken)
	countChanges := 0
	//  If s is not a valid textual representation of an IP address, ParseIP returns nil.

	region := p.parseRegion(request.NetworkRegion)
	if region == nil {
		return nil, fmt.Errorf("could not parse region from network region '%s'", request.NetworkRegion)
	}

	if region.deviceClient == nil {
		return nil, errors.New("provider has not been able to connect to the requested deviceClient to do lookups in the OSS for the region")
	}

	params := &device.GetParameters{
		// The inventory does not have fully qualified domain names for devices
		// so we need to strip the domain from the hostname before we do the lookup
		Hostname: HostFromFQDN(request.Hostname),
	}
	d, err := region.deviceClient.Get(ctx, params)
	if err != nil {
		p.logger.Warn("could not get OSS device", "hostname", params.Hostname, "error", err, "region", region.region)
		return nil, err
	}

	if len(d.Devices) == 0 {
		p.logger.Warn("could not find device in OSS", "hostname", params.Hostname, "region", region.region)
		return nil, fmt.Errorf("could not find device in provider inventory system")
	}

	resp := &provider.ResolveResourcePluginResponse{}

	host := d.Devices[0]
	switch strings.ToUpper(host.Vendor) {
	case "HUAWEI":
		p.logger.Debug("provider found device in oss, overwrite settings", "settings.resource_plugin", "vrp")
		resp.ResourcePlugin = "vrp"
		countChanges++
	case "CTC", "VXFIBER":
		p.logger.Debug("provider found device in oss, overwrite settings", "settings.resource_plugin", "ctc")
		resp.ResourcePlugin = "ctc"
		countChanges++
	}
	p.logger.Named("pre.ResolveResourcePlugin").Debug("processing request in", "changes", countChanges)
	return resp, nil
}

func setupEnv() error {

	if os.Getenv("APP_AUTH_ROLE_NAME") == "" {
		os.Setenv("APP_AUTH_ROLE_NAME", "role-test-fiberop-swpx")
	}
	if os.Getenv("APP_AUTH_MOUNT_PATH") == "" {
		os.Setenv("APP_AUTH_MOUNT_PATH", "fo-approle")
	}
	if os.Getenv("APP_AUTH_REGION") == "" {
		os.Setenv("APP_AUTH_REGION", "VX_SE2")
	}

	if os.Getenv("VAULT_ADDR") == "" {
		os.Setenv("VAULT_ADDR", "https://vault.vxfiber.dev")
	}
	if os.Getenv("VAULT_TOKEN") == "" {
		return errors.New("cannot run without vault token, VAULT_TOKEN is not set")
	}
	if os.Getenv("DISABLE_SECURITY") == "" {
		os.Setenv("DISABLE_SECURITY", "false")
	}
	return nil
}

var vxTestRegions map[string]string = map[string]string{
	"VX_SE1": "localhost:9001",
	"VX_SE2": "localhost:9002",
}

var vxRegions map[string]string = map[string]string{
	"VX_SA1": "localhost:9003",
	"VX_DE1": "localhost:9004",
	"VX_AT1": "localhost:9005",
	"VX_UK1": "localhost:9006",
}

func main() {

	logger = hclog.New(&hclog.LoggerOptions{
		Name:        PROVIDER_NAME,
		Level:       hclog.Debug,
		JSONFormat:  true,
		DisableTime: true,
	})
	if err := setupEnv(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// setup security for the application
	var appToken string
	var appRegion region.Region
	if os.Getenv("DISABLE_SECURITY") == "true" {
		logger.Warn("RUNNING IN INSECURE MODE, app auth is disabled, enforcing of endpoints is disabled")
	} else {

		appCfg := appauth.Enable()
		aa, err := appauth.New(appCfg, logger.Named("app-auth"))
		if err != nil {
			logger.Error("could not enable app auth", "error", err)
			os.Exit(1)
		}

		if err := aa.Authenticate(); err != nil {
			logger.Error("could not authenticate app", "error", err)
			os.Exit(1)

		}

		appToken = aa.GetToken()
		appRegion = *aa.GetRegion()

		logger.Info("app auth enabled", "region", aa.GetRegion())
	}

	_ = appRegion
	_ = appToken

	// when the config system in swpx is done this could be moved to that instead
	// for now this will be good enough
	grpcAddr := os.Getenv("OSS_INVENTORY_ADDR")
	if grpcAddr == "" {
		grpcAddr = "127.0.0.1:9001"
	}

	devClients := map[string]device.ServiceClient{}
	var vxr map[string]string
	switch appRegion {
	case region.Region_VX_SE1, region.Region_VX_SE2:
		vxr = vxTestRegions
	default:
		vxr = vxRegions
	}
	for k, v := range vxr {
		conn, err := grpc.Dial(
			v,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			middleware.WithGrpcClientAuthInterceptor(),
			grpc.WithTimeout(5*time.Second),
			grpc.WithBlock(),
		)
		if err != nil {
			logger.Error("could not connect to inventory GRPC server", "error", err, "region", k, "address", v)
			continue
		}
		devClients[k] = device.NewServiceClient(conn)
	}

	for k := range devClients {
		logger.Info("connected to inventory GRPC server", "region", k)
	}

	prov := &Provider{logger: logger, deviceClients: devClients, appToken: appToken, appRegion: appRegion}

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

type networkRegion struct {
	region       string
	domain       string
	deviceClient device.ServiceClient
}

func (p *Provider) parseRegion(msgRegion string) *networkRegion {

	switch strings.ToUpper(msgRegion) {
	case "VX_SA1", "SA1":
		return &networkRegion{
			region:       "VX_SA1",
			domain:       ".joburg.net.venturanext.se",
			deviceClient: p.deviceClients["VX_SA1"],
		}
	case "VX_UK1", "UK1":
		return &networkRegion{
			region:       "VX_UK1",
			domain:       ".net.uk1.vx.se",
			deviceClient: p.deviceClients["VX_UK1"],
		}

	case "VX_DE1", "DE1":
		return &networkRegion{
			region:       "VX_DE1",
			domain:       ".net.de1.vx.se",
			deviceClient: p.deviceClients["VX_DE1"],
		}

	case "VX_AT1", "AT1":
		return &networkRegion{
			region:       "VX_AT1",
			domain:       ".net.at1.vx.se",
			deviceClient: p.deviceClients["VX_AT1"],
		}

	case "VX_SE2", "SE2":
		return &networkRegion{
			region:       "VX_SE2",
			domain:       ".net.se2.vx.se",
			deviceClient: p.deviceClients["VX_SE2"],
		}

	case "VX_SE1", "SE1":
		return &networkRegion{
			region:       "VX_SE1",
			domain:       ".net.se1.vx.se",
			deviceClient: p.deviceClients["VX_SE1"],
		}

	default:
		return nil
	}

}
