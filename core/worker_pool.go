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
	"container/heap"
	"context"
	"fmt"

	// "github.com/davecgh/go-spew/spew"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/segmentio/ksuid"

	"git.liero.se/opentelco/go-swpx/proto/go/core"
	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
)

var start time.Time

func init() {
	start = time.Now()
}

// implementation
type workerPool struct {
	pool     workers
	done     chan *worker
	response chan *pb_core.Response
	index    int
	ops      int
	handler  RequestHandler
	logger   hclog.Logger
}

// create a new workerPool of workers
func newWorkerPool(nWorker int, nRequester int, logger hclog.Logger) *workerPool {
	logger = logger.Named("workerPool")

	logger.Debug("create new workerPool", "workers", nWorker, "buffer", nRequester)
	done := make(chan *worker, nWorker)
	// create balancer
	b := &workerPool{
		pool:     make(workers, 0, nWorker),
		done:     done,
		response: make(chan *pb_core.Response, nRequester),
		index:    0,
		ops:      0,
		logger:   logger,
	}
	b.SetHandler(b._defaultRequestHandler) // set default handler
	// Start the workers
	for i := 0; i < nWorker; i++ {

		w := &worker{
			Id:             i + 1,
			Name:           ksuid.New().String(),
			Pending:        0,
			Executed:       0,
			TimedOut:       0,
			requests:       make(chan *Request, nRequester),
			requestHandler: b.handler,
			logger:         logger.Named(fmt.Sprintf("worker-%d", i+1)),
		}

		// push the worker to the heap pool
		heap.Push(&b.pool, w)
		go w.start(b.done)
	}
	return b
}

// start the workerPool and listen for new events on the channel
func (b *workerPool) start(requestChan chan *Request) {
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
func (p *workerPool) print() {
	sum := 0
	sumsq := 0
	for _, w := range p.pool {
		// fmt.Printf("%3d ", w.Pending)
		sum += w.Pending
		sumsq += w.Pending * w.Pending
	}
	avg := float64(sum) / float64(len(p.pool))
	variance := float64(sumsq)/float64(len(p.pool)) - avg*avg
	p.logger.Info("statistics", "avg", avg, "variance", variance, "ops", p.ops, "elapsed", time.Since(start).Seconds(), "avgs", float64(p.ops)/float64(time.Since(start).Seconds()))
}

// dispatches a request to a worker and takes a request as argument.
func (p *workerPool) dispatch(request *Request) {
	if false {
		p.logger.Error("dispatch: if false?")
		w := p.pool[p.index]
		w.requests <- request
		w.Pending++
		p.index++
		if p.index >= len(p.pool) {
			p.index = 0
		}
		return
	}
	// take a woprker form the pool
	w := heap.Pop(&p.pool).(*worker)
	// give it a paylod
	w.requests <- request
	// add to pending
	w.Pending++
	// push it back to the pool
	heap.Push(&p.pool, w)
}

func (p *workerPool) completed(w *worker) {
	if false {
		w.logger.Error("completed: if false?")
		w.Pending--
		p.ops++
		return
	}
	w.Pending--
	p.ops++
	heap.Remove(&p.pool, w.Id)
	heap.Push(&p.pool, w)
}

// _defaultRequestHandler is set when the pool is created as a default handler
// a method to use logging in the handler
func (p *workerPool) _defaultRequestHandler(ctx context.Context, msg *Request, resp *core.Response) error {
	p.logger.Warn("default request handler for pool in use", "request_hostname", msg.Hostname)
	return nil
}

// SetHandler sets the Request Handler for the worker pool
// enables the core to inject a handler after initiated the Pool
func (p *workerPool) SetHandler(handler RequestHandler) {
	p.handler = handler
}
