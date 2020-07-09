package core

import (
	"container/heap"
	"context"
	"fmt" // "github.com/davecgh/go-spew/spew"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"

	"github.com/segmentio/ksuid"

	"git.liero.se/opentelco/go-swpx/errors"
	proto "git.liero.se/opentelco/go-swpx/proto/resource"
	"git.liero.se/opentelco/go-swpx/shared"
)

var start time.Time
var mongoClient *mongo.Client
var mongoCache *mongo.Collection
var useCache bool

func init() {
	start = time.Now()

	initMongoCache()
}

func initMongoCache() {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	//TODO parametrize the URI (eg. read from config file)
	var err error
	if mongoClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017")); err != nil {
		logger.Error("Error initializing Mongo client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = mongoClient.Connect(ctx); err != nil {
		logger.Error("Error connecting Mongo client:", err)
	}

	// Check the connection

	if err = mongoClient.Ping(context.TODO(), nil); err != nil {
		logger.Error("Can't ping Mongo client:", err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = mongoClient.Ping(ctx, readpref.Primary())

	useCache = true
}

// NetInterface which is cached in MongoDB
type NetInterface struct {
	Index       int64  `bson:"index,omitempty"`
	Description string `bson:"description,omitempty"`
	Alias       string `bson:"alias,omitempty"`
}

// worker that does the work
type worker struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Pending  int    `json:"pending"`
	Executed int    `json:"executed"`
	TimedOut int    `json:"timed_out"`
	messages chan *Request
}

type workers []*worker

// implementation
type workerPool struct {
	pool     workers
	done     chan *worker
	response chan *Response
	index    int
	ops      int
}

// create a new workerPool of workers
func newWorkerPool(nWorker int, nRequester int) *workerPool {
	logger.Debug("create new workerPool", "workers", nWorker, "buffer", nRequester)
	done := make(chan *worker, nWorker)
	// create balancer
	b := &workerPool{
		pool:     make(workers, 0, nWorker),
		done:     done,
		response: make(chan *Response, nRequester),
		index:    0,
		ops:      0,
	}
	// Start the workers
	for i := 0; i < nWorker; i++ {
		logger.Debug("staring worker:", "wid", i+1)
		w := &worker{
			Id:       i + 1,
			Name:     ksuid.New().String(),
			Pending:  0,
			Executed: 0,
			TimedOut: 0,
			messages: make(chan *Request, nRequester),
		}

		// push the worker to the heap pool
		heap.Push(&b.pool, w)
		go w.start(b.done, b.response)
	}
	return b
}

// start the workerPool and listen for new events on the channel
func (b *workerPool) start(requestChan chan *Request) {
	// print status every 30 secs. not for production
	// TODO prometheus endpoint.. ?
	go func() {
		for {
			b.print()
			time.Sleep(30 * time.Second)
		}
	}()

	// start going through the incoming channel
	go func() {
		for {
			select {
			case req := <-requestChan:
				// dispatch the work to a plugin
				b.dispatch(req)
			case w := <-b.done:
				// complete the job
				b.completed(w)
			}

		}
	}()
}

// prints stats to the console (dev)
func (b *workerPool) print() {
	sum := 0
	sumsq := 0
	for _, w := range b.pool {
		fmt.Printf("%3d ", w.Pending)
		sum += w.Pending
		sumsq += w.Pending * w.Pending
	}
	avg := float64(sum) / float64(len(b.pool))
	variance := float64(sumsq)/float64(len(b.pool)) - avg*avg
	fmt.Printf("   %3.2f %3.2f   ops: %10d     elapsed/s: %8.2f  avg/s: %.1f\n", avg, variance, b.ops, time.Since(start).Seconds(), float64(b.ops)/float64(time.Since(start).Seconds()))
}

// format stats from the workerPool of workers
func (b *workerPool) printer() string {
	sum := 0
	sumsq := 0

	r := new(wresp)
	for _, w := range b.pool {
		r.Workers = append(r.Workers, *w)
		sum += w.Pending
		sumsq += w.Pending * w.Pending
	}

	avg := float64(sum) / float64(len(b.pool))
	vari := float64(sumsq)/float64(len(b.pool)) - avg*avg
	r.Average = avg
	r.Variance = vari
	r.Elapsed = time.Since(start).Seconds()
	r.OpsSecond = float64(b.ops) / float64(time.Since(start).Seconds())

	return fmt.Sprintf("%v", r)
}

// dispatches a request to a worker and takes a request as argument.
func (b *workerPool) dispatch(request *Request) {
	if false {
		logger.Error("dispatch: if false?")
		w := b.pool[b.index]
		w.messages <- request
		w.Pending++
		b.index++
		if b.index >= len(b.pool) {
			b.index = 0
		}
		return
	}
	// take a woprker form the pool
	w := heap.Pop(&b.pool).(*worker)
	// give it a paylod
	w.messages <- request
	// add to pending
	w.Pending++
	// push it back to the pool
	heap.Push(&b.pool, w)
}

func (b *workerPool) completed(w *worker) {
	if false {
		logger.Error("completed: if false?")
		w.Pending--
		b.ops++
		return
	}
	w.Pending--
	b.ops++
	heap.Remove(&b.pool, w.Id)
	heap.Push(&b.pool, w)
}

// Heap implementation for worker-pool
func (p workers) Len() int           { return len(p) }
func (p workers) Less(i, j int) bool { return p[i].Pending < p[j].Pending }

func (p *workers) Swap(i, j int) {
	a := *p
	a[i], a[j] = a[j], a[i]
	a[i].Id = i
	a[j].Id = j
}

func (p *workers) Push(x interface{}) {
	a := *p
	n := len(a)
	a = a[0 : n+1]
	w := x.(*worker)
	a[n] = w
	w.Id = n
	*p = a
}

func (p *workers) Pop() interface{} {
	a := *p
	*p = a[0 : len(a)-1]
	w := a[len(a)-1]
	w.Id = -1 // for safety
	return w
}

// worker stats
type wresp struct {
	Workers    []worker `json:"workers"`
	Variance   float64  `json:"variance"`
	Average    float64  `json:"average"`
	Operations int      `json:"total_operations"`
	Failed     int      `json:"failed_operations"`
	Elapsed    float64  `json:"elapsed"`
	OpsSecond  float64  `json:"operations_second"`
}

// start the worker and ready it to accept payloads
func (w *worker) start(done chan *worker, res chan<- *Response) {
	for {
		select {
		case msg := <-w.messages:
			resp := &Response{RequestObjectID: msg.ObjectID}

			// do work with payload
			err := handle(msg.Context, msg, resp, handleMsg)

			if err != nil {
				resp.Error = err
				w.TimedOut++
			}

			msg.Response <- resp
			w.Executed++
			// response is not nandled this way right now. May never be.

			logger.Debug("response back in queue.")
			done <- w
		}
	}
}

// TODO this just runs some functions.. not a real implementation
func providerFunc(provider shared.Provider, msg *Request) {
	name, err := provider.Name()
	if err != nil {
		logger.Debug("getting provider name failed", "error", err)
	}
	ver, err := provider.Version()
	if err != nil {
		logger.Debug("getting provider name failed", "error", err)
	}
	weight, err := provider.Weight()
	if err != nil {
		logger.Debug("getting provider name failed", "error", err)
	}

	l, err := provider.Lookup(msg.ObjectID)
	if err != nil {
		log.Println(err)
	}

	logger.Debug("this is coming back from the plugin", "id", l)
	logger.Debug("data from provider plugin", "provider", name, "version", ver, "weight", weight)
}

func handle(ctx context.Context, msg *Request, resp *Response, f func(msg *Request, resp *Response) error) error {
	c := make(chan error, 1)
	go func() { c <- f(msg, resp) }()
	select {
	case <-ctx.Done():
		log.Println("got a timeout, letrs go")
		return ctx.Err()
	case err := <-c:
		if err != nil {
			log.Println("err: ", err)
		}
		return err
	}
	return nil
}

func handleMsg(msg *Request, resp *Response) error {
	logger.Debug("worker has payload")
	log.Printf("the user has sent in %s as provider", msg.Provider)
	var conf shared.Configuration

	// check if a provider is selected in the request
	if msg.Provider != "" {
		provider := providers[msg.Provider]
		if provider == nil {
			resp.Error = errors.New("the provider is missing/does not exist", errors.ErrInvalidProvider)
			return resp.Error
		}
		// run some provider funcs
		providerFunc(provider, msg)
		conf, _ = provider.GetConfiguration(msg.Context)

	} else {
		// no provider selected, walk all providers
		for pname, provider := range providers {
			logger.Debug("parsing provider", "provider", pname)
			providerFunc(provider, msg)
		}
	}

	// select resource-plugin to send the requests to
	plugin := resources[msg.Resource]
	if plugin == nil {
		logger.Error("selected driver is not a installed resource-driver-plugin", "selected-driver", msg.Resource)
		resp.Error = errors.New("selected driver is missing/does not exist", errors.ErrInvalidResource)
		return nil
	}

	plugin.SetConfiguration(msg.Context, conf)

	// implementation of different messages that SWP-X can handle right now
	// TODO is this the best way to to this.. ?
	switch msg.Type {
	case GetTechnicalInformationElement:
		return handleGetTechnicalInformationElement(msg, resp, plugin, conf)
	case GetTechnicalInformationPort:
		return handleGetTechnicalInformationPort(msg, resp, plugin, conf)
	}

	return nil
}

// handleGetTechnicalInformationElement gets full information of a Element
func handleGetTechnicalInformationElement(msg *Request, resp *Response, plugin shared.Resource, conf shared.Configuration) error {
	ver, err := plugin.Version()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	logger.Info("calling version ok", "version", ver)

	return nil
}

// handleGetTechnicalInformationPort gets information related to the selected interface
func handleGetTechnicalInformationPort(msg *Request, resp *Response, plugin shared.Resource, conf shared.Configuration) error {
	protConf := shared.Conf2proto(conf)
	req := &proto.NetworkElement{
		Hostname:  msg.NetworkElement,
		Interface: *msg.NetworkElementInterface,
		Conf:      &protConf,
	}

	var iface *proto.NetworkElementInterface
	var cacheResult bson.M
	var err error

	if useCache {
		mongoCache = mongoClient.Database("test").Collection("myCollection")
		err = mongoCache.FindOne(context.Background(), bson.M{"host": req.Hostname, "interface": req.Interface}).Decode(&cacheResult)
	}

	if err == mongo.ErrNoDocuments || !useCache {
		if iface, err = plugin.MapInterface(msg.Context, req); err != nil {
			logger.Error("error running map interface", "err", err.Error())
			resp.Error = errors.New(err.Error(), errors.ErrInvalidPort)
			return err
		}

		// save in cache upon success
		if useCache {
			if err = cache(req, iface); err != nil {
				return err
			}
		}

	} else if err != nil {
		logger.Error("Error fetching from cache:", err)
		return err
	}

	var ok bool
	if _, ok = cacheResult["index"].(int64); ok {
		iface = &proto.NetworkElementInterface{}
		iface.Index = cacheResult["index"].(int64)
	}

	//if the return is 0 something went wrong
	if iface.Index == 0 {
		logger.Error("error running map interface", "err", "index is zero")
		resp.Error = errors.New("interface index returned zero", errors.ErrInvalidPort)
		return err
	}

	logger.Info("got back info from MapInterface", "index", iface.Index)

	req.InterfaceIndex = iface.Index

	ti, err := plugin.TechnicalPortInformation(msg.Context, req)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	logger.Info("calling technical info ok ", "result", ti)
	resp.NetworkElement = ti

	return nil

}

func cache(req *proto.NetworkElement, iface *proto.NetworkElementInterface) error {
	_, err := mongoCache.InsertOne(
		context.Background(),
		bson.M{"host": req.Hostname, "interface": req.Interface, "index": iface.Index},
	)
	if err != nil {
		logger.Error("Error saving info in cache: ", err)
	}

	return err
}
