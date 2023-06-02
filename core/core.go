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

package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path"

	"git.liero.se/opentelco/go-swpx/config"
	"git.liero.se/opentelco/go-swpx/shared"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
	"go.mongodb.org/mongo-driver/mongo"
	"go.temporal.io/sdk/client"
)

const (
	RequestBufferSize int = 2000000 // nolint

	PluginProviderStr string = "provider"
	PluginResourceStr string = "resource"

	PluginPath string = "./plugins"
	Resources  string = "resources"
	Providers  string = "providers"

	VersionString string = "0.1-beta"
)

var (
	VERSION *version.Version
)

func init() {
	// Create an hclog.Logger
	VERSION, _ = version.NewVersion(VersionString)

}

// core app
type Core struct {
	cacheEnabled bool

	pollResponseCache PollResponseCache
	interfaceCache    PollInterfaceCache

	resources resourceMap
	providers providerMap

	tc client.Client

	config *config.Configuration
	logger hclog.Logger
}

// Start the core application
func (c *Core) Start() error {
	c.logger.Debug("starting core")

	go func() {
		// catch interrupt and kill all plugins
		csig := make(chan os.Signal, 1)
		signal.Notify(csig, os.Interrupt)
		for range csig {
			plugin.CleanupClients()
		}
	}()

	return nil
}

// New creates a new SWPX Core Application
func New(conf *config.Configuration, mongoClient *mongo.Client, logger hclog.Logger) (*Core, error) {
	var err error
	var ctx = context.Background()

	core := &Core{
		resources: make(map[string]shared.Resource),
		providers: make(map[string]shared.Provider),
		config:    conf,
		logger:    logger,
	}

	if mongoClient == nil {
		core.cacheEnabled = false
	}

	var availableResources = make(map[string]*plugin.Client)
	var availableProviders = make(map[string]*plugin.Client)

	// load all provider and resource plugins (files)
	if availableProviders, err = core.LoadPlugins(path.Join(PluginPath, Providers), PluginProviderStr); err != nil {
		logger.Error("error getting available provider resources", "error", err)
	}
	// load resource plugins, vrp etc
	if availableResources, err = core.LoadPlugins(path.Join(PluginPath, Resources), PluginResourceStr); err != nil {
		logger.Error("error getting available resources resources", "error", err)
	}

	if core.config.Request.DefaultProvider != "" {
		if _, ok := availableProviders[core.config.Request.DefaultProvider]; !ok {
			logger.Warn("the selected provider was not found, falling back on no provider", "selected", core.config.Request.DefaultProvider)
		} else {
			logger.Info("selected default_provider found and loaded", "default_provider", core.config.Request.DefaultProvider)
		}
	}

	// load the resources
	if err := core.LoadResourcePlugins(availableResources); err != nil {
		return nil, err
	}
	// load the Providers
	if err := core.LoadProviderPlugins(availableProviders); err != nil {
		return nil, err
	}

	// setup mongodb client

	if core.interfaceCache, err = newInterfaceCache(ctx, mongoClient, conf.GetMongoByLabel(config.LabelMongoInterfaceCache), logger); err != nil {
		logger.Error("error creating cache", "error", err)
		logger.Info("no mongo connection established", "cache_enabled", false)
		core.cacheEnabled = false
		return core, nil
	}

	if core.pollResponseCache, err = newResponseCache(ctx, mongoClient, conf.GetMongoByLabel(config.LabelMongoResponseCache), logger); err != nil {
		logger.Error("cannot set response cache", "error", err)
		return core, nil
	}
	core.cacheEnabled = true

	return core, nil
}

// LoadResourcePlugins iterates the resources and connect to the plugin
func (c *Core) LoadResourcePlugins(availableResources map[string]*plugin.Client) error {
	for name, p := range availableResources {
		var raw interface{}
		var err error

		c.logger.Debug("connect to resource", "name", name)
		rrpc, err := p.Client()
		if err != nil {
			return fmt.Errorf("could not get resource client: %w", err)
		}

		raw, err = rrpc.Dispense("resource")
		if err == nil {
			c.logger.Info("dispense and load resource plugin", "plugin", name)

			if resource, ok := raw.(shared.Resource); ok {

				c.logger.Debug("call version on resource plugin", "plugin", name)
				v, err := resource.Version()
				c.resources[name] = resource
				if err != nil {
					return fmt.Errorf("could not get version for plugin: %w", err)
				}
				c.logger.Info("loaded resource plugin", "plugin", name, "version", v)
			} else {
				return fmt.Errorf("type assertions failed for plugin: %s", name)
			}
		} else {
			return fmt.Errorf("failed to dispense resource plugin: %w", err)
		}

	}
	return nil
}

// LoadProviderPlugins iterates providers and connect to the plugin.
func (c *Core) LoadProviderPlugins(availableProviders map[string]*plugin.Client) error {
	for name, p := range availableProviders {
		var raw interface{}
		var err error

		c.logger.Debug("connecting to plugin", "plugin", name, "version", p.NegotiatedVersion(), "protocol", p.Protocol())

		rpc, err := p.Client()
		if err != nil {
			return fmt.Errorf("could not get provider client: %w", err)
		}

		err = rpc.Ping()
		if err != nil {
			return fmt.Errorf("could not ping provider plugin: %w", err)
		}

		raw, err = rpc.Dispense("provider")
		if err == nil {
			provider, ok := raw.(shared.Provider)
			if !ok || provider == nil {
				return fmt.Errorf("failed to load provider plugin: %s", name, "provider", provider, "ok", ok)
			}

			// get information about the provider to use on request
			var (
				err error
				n   string // name
				w   int64  // weight
				v   string // version
			)

			if n, err = provider.Name(); err != nil {
				return fmt.Errorf("failed to get provider name: %w", err)
			}
			if v, err = provider.Version(); err != nil {
				return fmt.Errorf("failed to get provider version: %w", err)
			}

			c.providers[n] = provider

			c.logger.Info("loaded provider", "name", n, "version", v, "weight", w)

			continue
		} else {

			if err := rpc.Close(); err != nil {
				return fmt.Errorf("error trying to dispense provider: %w", err)
			}
		}
	}
	return nil
}

// LoadPlugins loads plugins in a given folder
func (c *Core) LoadPlugins(pluginPath string, pluginType string) (map[string]*plugin.Client, error) {
	loadedPlugins := make(map[string]*plugin.Client)
	c.logger.Debug("searching for plugins", "path", pluginPath)
	plugs, err := ioutil.ReadDir(pluginPath)
	if err != nil {
		return loadedPlugins, err
	}
	for _, p := range plugs {
		if !p.IsDir() {
			if c.config.BlacklistProvider.Has(p.Name()) {
				c.logger.Debug("plugin is blacklisted", "plugin", p.Name())
				continue
			}
			c.logger.Debug("found plugin", "type", pluginType, "plugin", p.Name())
			loadedPlugins[p.Name()] = plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig:  shared.Handshake,
				Plugins:          shared.PluginMap,
				Cmd:              exec.Command(path.Join(pluginPath, p.Name())),
				Managed:          true,
				AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC}, // Only allow GRPC
				Logger:           c.logger.Named(pluginType),
			})

		}
	}
	return loadedPlugins, nil
}

type Transport interface {
	Ping()
}
