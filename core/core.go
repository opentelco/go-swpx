package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"sort"

	"git.liero.se/opentelco/go-swpx/shared"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
)

const (
	REQUEST_BUFFER_SIZE int    = 2000000
	MAX_REQUESTS               = 1000000
	WORKERS                    = 1
	APP_NAME            string = "swpx-core"
	RESOURCE            string = "resource"
	PLUGIN_PATH         string = "./plugins"
	RESOURCES           string = "resources"
	PROVIDERS           string = "providers"
	VERSION_STRING      string = "0.1-beta"

	// Name of queues SWPX Listens to.
	CommandQueue string = "opentelco.dnc.cmd"
	EventQueue   string = "opentelco.dnc.events"
)

var (
	// this contants should be moved to environment variables and arguments to CMD
	EVENT_SERVERS []string    = []string{"nats://localhost:14222", "nats://localhost:24222", "nats://localhost:34222"}
	TEST_CHAN     chan string = make(chan string, 0)
)

var (
	logger             hclog.Logger
	VERSION            *version.Version
	StopRequestHandler chan bool = make(chan bool)

	// Global request queue
	RequestQueue chan *Request = make(chan *Request, REQUEST_BUFFER_SIZE)

	// if swpx has an established InterfaceCache
	useCache bool

)

var Cache InterfaceCache

func init() {
	// Create an hclog.Logger
	VERSION, _ = version.NewVersion(VERSION_STRING)
	logger = hclog.New(&hclog.LoggerOptions{
		Name:   APP_NAME,
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

}

// LoadPlugins loads plugins in a given folder
func LoadPlugins(pluginPath string) (map[string]*plugin.Client, error) {
	loadedPlugins := make(map[string]*plugin.Client)
	plugs, err := ioutil.ReadDir(pluginPath)
	if err != nil {
		return loadedPlugins, err
	}
	for _, p := range plugs {
		if !p.IsDir() {
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
	Swarm *workerPool

	transport Transport
}

// start the core app
func (c *Core) Start() {
	c.Swarm.start(RequestQueue)

	// catch interrupt and kill all plugins
	csig := make(chan os.Signal, 1)
	signal.Notify(csig, os.Interrupt)
	go func() {
		for _ = range csig {
			plugin.CleanupClients()
			os.Exit(1)
		}
	}()

}

//
func CreateCore() *Core {
	var err error

	// create core
	core := &Core{
		Swarm: newWorkerPool(WORKERS, MAX_REQUESTS),
		//transport: Transport(t),
	}

	// load all provider and resource plugins (files)
	if availableProviders, err = LoadPlugins(path.Join(PLUGIN_PATH, PROVIDERS)); err != nil {
		logger.Error(err.Error())
	}
	// load resource plugins, vrp etc
	if availableResources, err = LoadPlugins(path.Join(PLUGIN_PATH, RESOURCES)); err != nil {
		logger.Error(err.Error())
	}

	// interate the resources and connect to the plugin
	// TODO move to a function instead?
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
			log.Printf("succesfully dispensed resource plugin (%s)", name)
			if resource, ok := raw.(shared.Resource); ok {
				v, err := resource.Version()
				resources[name] = resource
				logger.Error("something went wrong", "version", v, "error", err)
			} else {
				logger.Error(fmt.Sprintf("type assertions failed. %s plugin does not implement Plugin %T", name, raw))
				os.Exit(1)
			}

		} else {
			logger.Error(err.Error())
			log.Println("error trying to dispense resource or provider: '", err, "'")
		}

	}

	// iterate providers and connect to the plugin.
	// TODO move to a function instead?
	for name, p := range availableProviders {
		var raw interface{}
		var err error

		logger.Debug("connect to plugin %s", name)
		rpc, err := p.Client()
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		raw, err = rpc.Dispense("provider")
		if err == nil {
			provider := raw.(shared.Provider)

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
			rpc.Close()
			log.Println("error trying to dispense resource or provider: '", err, "'")
		}
	}

	// Sort the list of providers by their Weight()
	sortedProviders = providers.Slice()
	sort.Sort(byWeight(sortedProviders))
	for n, p := range sortedProviders {
		name, _ := p.Name()
		w, _ := p.Weight()
		println(n, name, w)
	}


	// setup mongodb cache
	mc, err := initMongoDB("mongodb://localhost:27017")
	if err != nil {
		logger.Warn("could not establish mongodb connection: %s",err.Error())
		useCache = false
		return core
	}
	Cache,err = NewCache(mc,logger)
	if err != nil {
		logger.Error("cannot set cache: %s",err.Error())
		useCache = false
		return core
	}

	useCache = true



	return core
}
