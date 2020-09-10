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
	"log"
	"os"
	"time"

	"git.liero.se/opentelco/go-swpx/shared"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
	"github.com/spf13/viper"
)

var VERSION *version.Version
var logger hclog.Logger

const (
	VERSION_BASE string = "1.0-beta"
	//PROVIDER_NAME string = "zitius_provider"
	WEIGHT int64 = 23
)

var PROVIDER_NAME = "zitius_provider"

var conf shared.Configuration

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
	g.logger.Debug("message from provider, running version:", VERSION)
	return fmt.Sprintf("%s@%s", PROVIDER_NAME, VERSION.String()), nil
}

func (g *Provider) Weight() (int64, error) {
	// g.logger.Debug("return the weight", PROVIDER_NAME)
	return WEIGHT, nil
}

func (p *Provider) Lookup(id string) (string, error) {
	// duration := time.Duration(time.Millisecond * time.Duration(rand.Intn(1000-400)+400))
	// time.Sleep(duration)
	// p.logger.Debug("done", "execution_time", duration.String())
	return "nonon", nil
}

func (p *Provider) Match(id string) (bool, error) {
	return true, nil
}

func (p *Provider) Name() (string, error) {
	return PROVIDER_NAME, nil
}

func (p *Provider) GetConfiguration(ctx context.Context) (shared.Configuration, error) {
	return conf, nil
}

func loadConfig(logger hclog.Logger) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./" + PROVIDER_NAME)
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".providers/" + PROVIDER_NAME)
	viper.AddConfigPath("$HOME/." + PROVIDER_NAME)

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
		panic(fmt.Errorf("Unable to decode Config: %s \n", err))
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
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	loadConfig(logger)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			shared.PluginProviderKey: &shared.ProviderPlugin{Impl: &Provider{logger: logger}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
