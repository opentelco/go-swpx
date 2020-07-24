package core

import (
	"container/heap"
	"context"
	"fmt" // "github.com/davecgh/go-spew/spew"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/ksuid"

	"git.liero.se/opentelco/go-swpx/errors"
	proto "git.liero.se/opentelco/go-swpx/proto/resource"
	"git.liero.se/opentelco/go-swpx/shared"
)

var start time.Time

func init() {
	start = time.Now()
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
		logger.Error(err.Error())
	}

	logger.Debug("this is coming back from the plugin", "id", l)
	logger.Debug("data from provider plugin", "provider", name, "version", ver, "weight", weight)
}

func handle(ctx context.Context, msg *Request, resp *Response, f func(msg *Request, resp *Response) error) error {
	c := make(chan error, 1)
	go func() { c <- f(msg, resp) }()
	select {
	case <-ctx.Done():
		logger.Error("got a timeout, letrs go")
		return ctx.Err()
	case err := <-c:
		if err != nil {
			logger.Error("err: ", err.Error())
		}
		return err
	}
	return nil
}

func handleMsg(msg *Request, resp *Response) error {
	logger.Debug("worker has payload")
	logger.Info("the user has sent in %s as provider", msg.Provider)
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

	iface := &proto.NetworkElementInterface{}
	var cachedInterface *CachedInterface
	var cachedPhysPort *CachedPhysicalPortInformation
	var err error

	if useCache {
		logger.Debug("cache enabled, pop object from cache")
		cachedInterface, err = Cache.Pop(req.Hostname, req.Interface)
		if cachedInterface != nil {
			iface.Index = cachedInterface.InterfaceIndex
			resp.Transceiver = cachedInterface.TransceiverInformation
		}

		cachedPhysPort, err = PhysicalPortCache.PopPhysical(msg.Provider, msg.Resource)
		if cachedPhysPort != nil {
			findPhysicalPort(cachedPhysPort.PhysicalPortInformation, req, resp)
		}
	}

	// did not find cached item or cached is disabled
	if cachedInterface == nil || !useCache {
		var physPortResponse *proto.PhysicalPortinformationResponse
		if physPortResponse, err = plugin.GetPhysicalPort(msg.Context, req); err != nil {
			logger.Error("error running getphysport", "err", err.Error())
			resp.Error = errors.New(err.Error(), errors.ErrInvalidPort)
			return err
		}

		findPhysicalPort(physPortResponse.Data, req, resp)

		transceiver, err := plugin.GetVRPTransceiverInformation(msg.Context, req)
		resp.Transceiver = transceiver
		if iface, err = plugin.MapInterface(msg.Context, req); err != nil {
			logger.Error("error running map interface", "err", err.Error())
			resp.Error = errors.New(err.Error(), errors.ErrInvalidPort)
			return err
		}

		// save in cache upon success (if enabled)
		if useCache {
			if _, err = Cache.Set(req, iface, physPortResponse, transceiver); err != nil {
				return err
			}
			if _, err = PhysicalPortCache.SetPhysical(msg.Provider, msg.Resource, physPortResponse); err != nil {
				return err
			}

		}

	} else if err != nil {
		logger.Error("error fetching from cache:", err)
		return err
	}

	//if the return is 0 something went wrong
	if iface.Index == 0 {
		logger.Error("error running map interface", "err", "index is zero")
		resp.Error = errors.New("interface index returned zero", errors.ErrInvalidPort)
		return err
	}

	logger.Info("found index for selected interface", "index", iface.Index)

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

func findPhysicalPort(data []*proto.PhysicalPortInformation, req *proto.NetworkElement, resp *Response) {
	for _, element := range data {
		if element.Value == req.Interface {

			resp.PhysicalPort = element

			fields := strings.Split(element.Oid, ".")
			index, err := strconv.Atoi(fields[len(fields)-1])
			if err != nil {
				logger.Error("can't convert phys.port to int: ", err)
				return
			}
			req.PhysicalIndex = int64(index)
		}
	}
}
