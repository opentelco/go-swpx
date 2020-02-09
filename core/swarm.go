package core

import (
	"container/heap"
	"context"
	"fmt" // "github.com/davecgh/go-spew/spew"
	"log"
	"time"

	"github.com/opentelco/go-dnc/models/transport"
	"github.com/opentelco/go-swpx/errors"
	proto "github.com/opentelco/go-swpx/proto/resource"
	"github.com/opentelco/go-swpx/shared"
)

var start time.Time

func init() {
	start = time.Now()
}

type worker struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Pending  int    `json:"pending"`
	Executed int    `json:"executed"`
	TimedOut int    `json:"timed_out"`
	messages chan *Request
}

type workerPool []*worker

//
// Balancer
//

type swarm struct {
	pool     workerPool
	done     chan *worker
	response chan *Response
	index    int
	ops      int
}

func newSwarm(nWorker int, nRequester int) *swarm {
	logger.Debug("create new swarm", "workers", nWorker, "buffer", nRequester)
	done := make(chan *worker, nWorker)
	// create balancer
	b := &swarm{
		pool:     make(workerPool, 0, nWorker),
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
			Name:     transport.NewID().String(),
			Pending:  0,
			Executed: 0,
			TimedOut: 0,
			messages: make(chan *Request, nRequester),
		}

		heap.Push(&b.pool, w)
		go w.start(b.done, b.response)
	}
	return b
}

func (b *swarm) start(j chan *Request) {
	// print status every 30 secs. not for production
	go func() {
		for {
			b.print()
			time.Sleep(30 * time.Second)
		}
	}()

	// start going through the channel
	go func() {
		for {
			select {
			case req := <-j:
				b.dispatch(req)
			case w := <-b.done:
				b.completed(w)
			}

		}
	}()
}

// prints stats to the console
func (b *swarm) print() {
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

// format stats fror the swarm of workers
func (b *swarm) printer() string {
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

// dispatches a request to a worker and
// takes a request as argument.
func (b *swarm) dispatch(request *Request) {
	if false {
		w := b.pool[b.index]
		w.messages <- request
		w.Pending++
		b.index++
		if b.index >= len(b.pool) {
			b.index = 0
		}
		return
	}
	w := heap.Pop(&b.pool).(*worker)
	w.messages <- request
	w.Pending++ //
	heap.Push(&b.pool, w)
}

func (b *swarm) completed(w *worker) {
	if false {
		w.Pending--
		b.ops++
		return
	}
	w.Pending--
	b.ops++
	heap.Remove(&b.pool, w.Id)
	heap.Push(&b.pool, w)
}

// Heap to get worker
//

func (p workerPool) Len() int           { return len(p) }
func (p workerPool) Less(i, j int) bool { return p[i].Pending < p[j].Pending }

func (p *workerPool) Swap(i, j int) {
	a := *p
	a[i], a[j] = a[j], a[i]
	a[i].Id = i
	a[j].Id = j
}

func (p *workerPool) Push(x interface{}) {
	a := *p
	n := len(a)
	a = a[0 : n+1]
	w := x.(*worker)
	a[n] = w
	w.Id = n
	*p = a
}

func (p *workerPool) Pop() interface{} {
	a := *p
	*p = a[0 : len(a)-1]
	w := a[len(a)-1]
	w.Id = -1 // for safety
	return w
}

//
// worker
//

type wresp struct {
	Workers    []worker `json:"workers"`
	Variance   float64  `json:"variance"`
	Average    float64  `json:"average"`
	Operations int      `json:"total_operations"`
	Failed     int      `json:"failed_operations"`
	Elapsed    float64  `json:"elapsed"`
	OpsSecond  float64  `json:"operations_second"`
}

func (w *worker) start(done chan *worker, res chan<- *Response) {
	for {
		select {
		case msg := <-w.messages:
			resp := &Response{RequestObjectID: msg.ObjectID}

			// do work with payload
			err := handle(msg.Context, msg, resp, func(msg *Request, resp *Response) error {
				logger.Debug("worker has payload")
				log.Printf("the user has sent in %s as provider", msg.Provider)
				// check if a provider is selected in the request
				if msg.Provider != "" {
					provider := providers[msg.Provider]
					if provider == nil {
						resp.Error = errors.New("the provider is missing/does not exist", errors.ErrInvalidProvider)
						return resp.Error
					}
					providerFunc(provider, msg)

				} else {
					// no provider selected, walk all providers
					for pname, provider := range providers {
						logger.Debug("parsing provider", "provider", pname)
						providerFunc(provider, msg)
					}
				}

				// select resource
				res := resources[msg.Resource]
				if res == nil {
					logger.Error("selected driver is not a installed resource-driver-plugin", "selected-driver", msg.Resource)
					resp.Error = errors.New("selected driver is missing/does not exist", errors.ErrInvalidResource)
					return nil
				}

				switch msg.Type {
				// Get full information of a Element
				case GetTechnicalInformationElement:
					ver, err := res.Version()
					if err != nil {
						logger.Error(err.Error())
						return err
					}
					logger.Info("calling version ok", "version", ver)

					// Get information related to the selected interface
				case GetTechnicalInformationPort:
					req := &proto.NetworkElement{
						Hostname:  msg.NetworkElement,
						Interface: *msg.NetworkElementInterface,
					}

					iface, err := res.MapInterface(msg.Context, req)
					if err != nil {
						logger.Error("error running map interrace", "err", err.Error())
						resp.Error = errors.New(err.Error(), errors.ErrInvalidPort)
						return err
					}
					// if the return is 0 somethng went wrong
					if iface.Index == 0 {
						logger.Error("error running map interrace", "err", "index is zero")
						resp.Error = errors.New("interface index returned zero", errors.ErrInvalidPort)
						return err
					}

					logger.Info("got back info from MapInterface", "index", iface.Index)

					req.InterfaceIndex = iface.Index

					ti, err := res.TechnicalPortInformation(msg.Context, req)
					if err != nil {
						logger.Error(err.Error())
						return err
					}
					logger.Info("calling technical info ok ", "result", ti)
					resp.NetworkElement = ti
					return nil

				}

				return nil
			})

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

	// m, err := provider.Match(strconv.Itoa(rand.Int()))
	// if err != nil {
	// 	log.Println(err)
	// }
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
