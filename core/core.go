package core

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"sort"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
	"github.com/opentelco/go-swpx/core/transport/nats"
	"github.com/opentelco/go-swpx/shared"
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

// this contants should be moved to environment variables and arguments to CMD
var (
	EVENT_SERVERS []string    = []string{"nats://localhost:14222", "nats://localhost:24222", "nats://localhost:34222"}
	TEST_CHAN     chan string = make(chan string, 0)
)

var (
	logger             hclog.Logger
	VERSION            *version.Version
	StopRequestHandler chan bool     = make(chan bool)
	RequestQueue       chan *Request = make(chan *Request, REQUEST_BUFFER_SIZE)
)

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

// Transport is the interface that talks down to the DNC
type Transport interface {
	// Subscribe(queue string, requestChannel chan interface{}, responseChannel chan interface{}) error
	Ping()
	// Subscribe(subject string, queue string, channel chan interface{}) error
}

type Core struct {
	Swarm *swarm

	transport Transport
}

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

func CreateCore() *Core {
	var err error

	t, err := nats.New(EVENT_SERVERS)

	_, err = t.BindRecvQueueChan("vrp-driver", "dispatchers", TEST_CHAN)
	go func() {
		for {
			select {
			case s := <-TEST_CHAN:
				log.Println("got this string:", s)
			}
		}
	}()
	core := &Core{
		Swarm:     newSwarm(WORKERS, MAX_REQUESTS),
		transport: Transport(t),
	}

	// Create swarm, starts when Core is started.

	// load all provider annd resource plugins (files)
	if availableProviders, err = LoadPlugins(path.Join(PLUGIN_PATH, PROVIDERS)); err != nil {
		logger.Error(err.Error())
	}
	if availableResources, err = LoadPlugins(path.Join(PLUGIN_PATH, RESOURCES)); err != nil {
		logger.Error(err.Error())
	}

	// interate the resources and connect.
	for name, p := range availableResources {
		var raw interface{}
		var err error

		logger.Debug("connect to plugin", "name", name)
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
				logger.Error("Type assertions failed. %s plugin does not implement Plugin", name)
				os.Exit(1)
			}

		} else {
			logger.Error(err.Error())
			log.Println("error trying to dispense resource or provider: '", err, "'")
		}

	}
	// iterate providers and connect to the plugin.
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

	// Sort the list of providers by Weight
	sortedProviders = providers.Slice()
	sort.Sort(byWeight(sortedProviders))
	for n, p := range sortedProviders {
		name, _ := p.Name()
		w, _ := p.Weight()
		println(n, name, w)
	}

	// start the handler for Requests
	// HandleRequests()

	// start the swarm to handle incoming requests.

	// go func() {
	// 	var i int
	// 	for {
	// 		i++
	// 		fmt.Printf("----------%d-------------", i)
	// 		for n := 0; n < 10000; n++ {
	// 			r := &Request{
	// 				ObjectID: fmt.Sprintf("ZIT%d", rand.Int()),
	// 			}
	// 			// Add request to queue
	// 			RequestQueue <- r

	// 		}
	// 		time.Sleep(time.Second * 30)
	// 	}
	// }()
	return core
}

// HandleRequests handels the incoming jobs
// func HandleRequests() error {
// 	logger.Debug("request handler started")
// 	go func() {
// 		for {
// 			select {
// 			case r := <-RequestQueue:
// 				tx := &Response{}
// 				log.Printf("the user has sent in %s as provider", r.Provider)
// 				if r.Provider != "" {
// 					provider := providers[r.Provider]

// 					if provider == nil {
// 						tx.Error = errors.New("the selected provider does not exist")
// 						r.Response <- tx
// 						break
// 					}
// 				}

// 				for _, provider := range sortedProviders {
// 					logger.Debug("parsing provider", "provider", provider.Name)

// 					name, err := provider.Name()
// 					if err != nil {
// 						logger.Debug("getting provider name failed", "error", err)
// 					}
// 					ver, err := provider.Version()
// 					if err != nil {
// 						logger.Debug("getting provider name failed", "error", err)
// 					}
// 					weight, err := provider.Weight()
// 					if err != nil {
// 						logger.Debug("getting provider name failed", "error", err)
// 					}

// 					l, err := provider.Lookup(r.ObjectID)
// 					if err != nil {
// 						log.Println(err)
// 					}

// 					// m, err := provider.Match(strconv.Itoa(rand.Int()))
// 					// if err != nil {
// 					// 	log.Println(err)
// 					// }
// 					logger.Debug("this is coming back from the plugin", "id", l)
// 					logger.Debug("data from provider plugin", "provider", name, "version", ver, "weight", weight)
// 				}
// 				logger.Debug("handler received a new request", "objectID", r.ObjectID)
// 			case <-StopRequestHandler:
// 				logger.Debug("go kill signal from something, exit.")
// 				break
// 			}
// 		}
// 	}()
// 	return nil
// }
