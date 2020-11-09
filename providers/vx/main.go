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
	"time"
	
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
	"github.com/spf13/viper"
	
	"git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/shared"
)

var VERSION *version.Version
var logger hclog.Logger

const (
	VERSION_BASE string = "1.0-beta"
)

var PROVIDER_NAME = "vx"

var conf *shared.Configuration

func init() {
	// var err error
	// if VERSION, err = version.NewVersion(VERSION_BASE); err != nil {
	// 	log.Fatal(err)
	// }
}

// Provider is the implementation of the GRPC
type Provider struct {
	logger hclog.Logger
}

func (g *Provider) Version() (string, error) {
	g.logger.Debug("message from provider, running version:", VERSION)
	return fmt.Sprintf("%s@%s", PROVIDER_NAME, VERSION.String()), nil
}

func (p *Provider) Name() (string, error) {
	return PROVIDER_NAME, nil
}

func (p *Provider) PostHandler(ctx context.Context, request *core.Response) (*core.Response, error) {
	p.logger.Named("post-handler").Debug("processing response", "changes", 0)
	return request, nil
}

func (p *Provider) PreHandler(ctx context.Context, response *core.Request) (*core.Request, error) {
	p.logger.Named("pre-handler").Debug("processing request in",  "changes", 0)
	return response, nil
}

func (p *Provider) GetConfiguration(ctx context.Context) (*shared.Configuration, error) {
	return conf, nil
}

func loadConfig(logger hclog.Logger) {

	defaultSnmpConf := shared.ConfigSNMP{
		Community: "semipublic",
		Retries:   3,
		Version:   2,
		Timeout:   time.Second * 5,
	}

	viper.SetDefault("snmp", defaultSnmpConf)
	err := viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("unable to decode Config: %s \n", err))
	}
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   shared.MagicCookieKey,
	MagicCookieValue: shared.MagicCookieValue,
	
}

func main() {
	logger = hclog.New(&hclog.LoggerOptions{
		Name:       fmt.Sprintf("%s@%s", PROVIDER_NAME, VERSION.String()),
		Level:      hclog.Debug,
		Color: hclog.AutoColor,
	})

	// loadConfig(logger)
	
	prov := &Provider{logger: logger}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: map[string]plugin.Plugin{
			shared.PluginProviderKey: &shared.ProviderPlugin{
				Impl:   prov,
			},
		},
		GRPCServer:       plugin.DefaultGRPCServer,
		// Logger:           logger,
	})
	
}
