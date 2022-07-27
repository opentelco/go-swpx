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
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"

	"git.liero.se/opentelco/go-swpx/shared"
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

	// Global request queue
	RequestQueue = make(chan *Request, RequestBufferSize)
)

func init() {
	// Create an hclog.Logger
	VERSION, _ = version.NewVersion(VersionString)

}

// core app
type Core struct {
	// workers and queues
	swarm *workerPool

	cacheEnabled   bool
	responseCache  ResponseCache
	interfaceCache InterfaceCache

	config *shared.Configuration
	logger hclog.Logger
}

// Start the core application
func (c *Core) Start() error {
	c.swarm.start(RequestQueue)

	go func() {
		// catch interrupt and kill all plugins
		csig := make(chan os.Signal, 1)
		signal.Notify(csig, os.Interrupt)
		for range csig {
			plugin.CleanupClients()
			// TODO need to close swarm c.swarm.Close() ?
		}
	}()

	return nil
}

// New creates a new SWPX Core Application
func New(logger hclog.Logger) (*Core, error) {
	var err error

	// create core
	conf := shared.GetConfig()

	swarm := newWorkerPool(conf.PollerWorkers, conf.MaxPollerRequests, logger)

	core := &Core{
		swarm:  swarm,
		logger: logger,
	}

	logger.Info("setting core requestHandler for pool")
	swarm.SetHandler(core.RequestHandler)

	core.config = conf

	// load all provider and resource plugins (files)
	if availableProviders, err = core.LoadPlugins(path.Join(PluginPath, Providers), PluginProviderStr); err != nil {
		logger.Error("error getting available provider resources", "error", err)
	}
	// load resource plugins, vrp etc
	if availableResources, err = core.LoadPlugins(path.Join(PluginPath, Resources), PluginResourceStr); err != nil {
		logger.Error("error getting available resources resources", "error", err)
	}

	if core.config.DefaultProvider != "" {
		if _, ok := availableProviders[core.config.DefaultProvider]; !ok {
			logger.Warn("the selected provider was not found, falling back on no provider", "default_provider", core.config.DefaultProvider)
		} else {
			logger.Info("selected default_provider found and loaded", "default_provider", core.config.DefaultProvider)
		}
	}

	// load the resources
	if err := core.LoadResourcePlugins(); err != nil {
		return nil, err
	}
	// load the Providers
	if err := core.LoadProviderPlugins(); err != nil {
		return nil, err
	}

	// setup mongodb cache
	mongoClient, err := initMongoDb(conf.InterfaceCache, logger.Named("mongodb"))
	if err != nil {
		logger.Warn("could not establish mongodb connection", "error", err)
		logger.Info("no mongo connection established", "cache_enabled", false)
		core.cacheEnabled = false
		return core, nil
	}

	if core.interfaceCache, err = newInterfaceCache(mongoClient, logger, conf.InterfaceCache); err != nil {
		logger.Error("error creating cache", "error", err)
		logger.Info("no mongo connection established", "cache_enabled", false)
		core.cacheEnabled = false
		return core, nil
	}

	if core.responseCache, err = newResponseCache(mongoClient, logger, conf.ResponseCache); err != nil {
		logger.Error("cannot set response cache", "error", err)
		return core, nil
	}
	core.cacheEnabled = true

	return core, nil
}

// LoadResourcePlugins iterates the resources and connect to the plugin
func (c *Core) LoadResourcePlugins() error {
	for name, p := range availableResources {
		var raw interface{}
		var err error

		c.logger.Debug("connect to resource", "name", name)
		rrpc, err := p.Client()
		if err != nil {
			c.logger.Error("could not return resource client", "error", err.Error())
			continue
		}
		raw, err = rrpc.Dispense("resource")
		if err == nil {

			c.logger.Info("successfully dispensed resource plugin", "plugin", name)
			if resource, ok := raw.(shared.Resource); ok {
				_, err := resource.Version()
				resources[name] = resource
				if err != nil {
					c.logger.Error("could not get version for plugin", "plugin", name, "error", err.Error())
				}
			} else {
				c.logger.Error("type assertions failed for plugin", "plugin", name)
				os.Exit(1)
			}
		} else {
			c.logger.Error("failed to dispense resource or provider", "error", err.Error())
		}

	}
	return nil
}

// LoadProviderPlugins iterates providers and connect to the plugin.
func (c *Core) LoadProviderPlugins() error {
	for name, p := range availableProviders {
		var raw interface{}
		var err error

		c.logger.Debug("connecting to plugin", "plugin", name, "version", p.NegotiatedVersion(), "protocol", p.Protocol())

		rpc, err := p.Client()
		if err != nil {
			c.logger.Error(err.Error())
			continue
		}
		err = rpc.Ping()
		if err != nil {
			c.logger.Error("could nog ping", "error", err)
		}

		raw, err = rpc.Dispense("provider")
		if err == nil {
			provider, ok := raw.(shared.Provider)
			if !ok || provider == nil {
				c.logger.Error("failed to load provider_plugin", "plugin", name, "provider", provider, "ok", ok)
				continue
			}

			// get information about the provider to use on request
			var (
				err error
				n   string // name
				w   int64  // weight
				v   string // version
			)

			if n, err = provider.Name(); err != nil {
				log.Fatal("could not load provider")
			}
			if v, err = provider.Version(); err != nil {
				log.Fatal("could not load provider")
			}

			providers[n] = provider

			c.logger.Debug("loaded provider", "name", n, "version", v, "weight", w)

			continue
		} else {
			c.logger.Error(err.Error())

			rpcErr := rpc.Close()
			if rpcErr != nil {
				c.logger.Error("error trying to dispense resource or provider", "error", rpcErr.Error())
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
