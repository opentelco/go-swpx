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
	"sort"
	
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
	
	"git.liero.se/opentelco/go-swpx/shared"
)

const (
	RequestBufferSize int    = 2000000
	MaxRequests              = 1000000
	WORKERS                  = 1
	AppName           string = "swpx-core"
	PluginPath        string = "./plugins"
	Resources         string = "resources"
	Providers         string = "providers"
	VersionString     string = "0.1-beta"
)

var (
	logger  hclog.Logger
	VERSION *version.Version

	// Global request queue
	RequestQueue = make(chan *Request, RequestBufferSize)

	useCache       bool
	InterfaceCache *cache
	ResponseCache  *cache
)

func init() {
	// Create an hclog.Logger
	VERSION, _ = version.NewVersion(VersionString)

}

// LoadPlugins loads plugins in a given folder
func LoadPlugins(pluginPath string) (map[string]*plugin.Client, error) {
	loadedPlugins := make(map[string]*plugin.Client)
	logger.Debug("searching for plugins", "path", pluginPath)
	plugs, err := ioutil.ReadDir(pluginPath)
	if err != nil {
		return loadedPlugins, err
	}
	for _, p := range plugs {
		if !p.IsDir() {
			logger.Debug("found plugin", "plugin",p.Name())
			loadedPlugins[p.Name()] = plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig:  shared.Handshake,
				Plugins:          shared.PluginMap,
				Cmd:              exec.Command(path.Join(pluginPath, p.Name())),
				Managed:          true,
				AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC}, // Only allow GRPC
				Logger:           logger,
			})
			defer loadedPlugins[p.Name()].Kill()

		}
	}
	return loadedPlugins, nil
}

// TODO Transport is the interface that talks down to the DNC
// TODO not in use
type Transport interface {
	// Subscribe(queue string, requestChannel chan interface{}, responseChannel chan interface{}) error
	Ping()
	// Subscribe(subject string, queue string, channel chan interface{}) error
}

// core app
type Core struct {
	// workers and queues
	swarm *workerPool

	transport Transport
	logger hclog.Logger
}

// start the core app
func (c *Core) Start() error{
	c.swarm.start(RequestQueue)

	// catch interrupt and kill all plugins
	csig := make(chan os.Signal, 1)
	signal.Notify(csig, os.Interrupt)
	for range csig {
		
		plugin.CleanupClients()
		// TODO need to close swarm c.swarm.Close() ?
		return nil
		
	}
	
	return nil
}

//
func New(log hclog.Logger) (*Core,error) {
	var err error

	if log != nil {
		logger = log
	}
	
	// create core
	core := &Core{
		swarm: newWorkerPool(WORKERS, MaxRequests),
		//transport: Transport(t),
	}
	conf := shared.GetConfig()

	// load all provider and resource plugins (files)
	if availableProviders, err = LoadPlugins(path.Join(PluginPath, Providers)); err != nil {
		logger.Error("error getting available provider resources", "error", err)
	}
	// load resource plugins, vrp etc
	if availableResources, err = LoadPlugins(path.Join(PluginPath, Resources)); err != nil {
		logger.Error("error getting available resources resources", "error", err)
	}

	loadResources()
	loadProviders()

	// Sort the list of providers by their Weight()
	sortedProviders = providers.Slice()
	sort.Sort(byWeight(sortedProviders))
	for n, p := range sortedProviders {
		name, _ := p.Name()
		w, _ := p.Weight()
		println(n, name, w)
	}

	// setup mongodb cache
	mongoClient, err := initMongoDB(conf.InterfaceCache)
	if err != nil {
		logger.Warn("could not establish mongodb connection","error", err)
		logger.Info("no mongo connection established","cache_enabled", false)
		useCache = false
		return core, nil
	}
	if InterfaceCache, err = NewCache(mongoClient, logger, conf.InterfaceCache); err != nil {
		logger.Error("error creating cache", "error", err)
		logger.Info("no mongo connection established","cache_enabled", false)
		useCache = false
		return core,nil
	}

	if ResponseCache, err = NewCache(mongoClient, logger, conf.ResponseCache); err != nil {
		logger.Error("cannot set response cache", "error", err)
		return core, nil
	}

	useCache = true

	return core,nil
}

// iterate the resources and connect to the plugin
func loadResources() {
	for name, p := range availableResources {
		var raw interface{}
		var err error

		logger.Debug("connect to resource", "name", name)
		rrpc, err := p.Client()
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		raw, err = rrpc.Dispense("resource")
		if err == nil {
			logger.Info("succesfully dispensed resource plugin","plugin", name)
			if resource, ok := raw.(shared.Resource); ok {
				_, err := resource.Version()
				resources[name] = resource
				if err != nil {
					logger.Error("could not get version for plugin", "plugin", name, "error", err.Error())
					}
			} else {
				logger.Error("type assertions failed for plugin", "plugin", name)
				os.Exit(1)
			}
		} else {
			logger.Error("failed to dispense resource or provider", "error", err.Error())
		}

	}
}

// iterate providers and connect to the plugin.
func loadProviders() {
	for name, p := range availableProviders {
		var raw interface{}
		var err error

		logger.Debug("connecting to plugin","plugin", name)
		rpc, err := p.Client()
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		raw, err = rpc.Dispense("provider")
		if err == nil {
			provider, ok := raw.(shared.Provider)
			if !ok || provider == nil {
				logger.Error("failed to load provider_plugin", "plugin", name)
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
			if w, err = provider.Weight(); err != nil {
				log.Fatal("could not load provider")
			}

			providers[n] = provider

			logger.Debug("loaded provider", "name", n, "version", v, "weight", w)

			continue
		} else {
			logger.Error(err.Error())

			rpcErr := rpc.Close()
			if rpcErr != nil {
				logger.Error("error trying to dispense resource or provider", rpcErr.Error(), "'")
			}
		}
	}
}
