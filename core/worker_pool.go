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
	"fmt" // "github.com/davecgh/go-spew/spew"
	"time"
	
	"github.com/segmentio/ksuid"
	
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
}

// create a new workerPool of workers
func newWorkerPool(nWorker int, nRequester int) *workerPool {
	logger.Debug("create new workerPool", "workers", nWorker, "buffer", nRequester)
	done := make(chan *worker, nWorker)
	// create balancer
	b := &workerPool{
		pool:     make(workers, 0, nWorker),
		done:     done,
		response: make(chan *pb_core.Response, nRequester),
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
func (b *workerPool) print() {
	sum := 0
	sumsq := 0
	for _, w := range b.pool {
		// fmt.Printf("%3d ", w.Pending)
		sum += w.Pending
		sumsq += w.Pending * w.Pending
	}
	avg := float64(sum) / float64(len(b.pool))
	variance := float64(sumsq)/float64(len(b.pool)) - avg*avg
	logger.Info("statistics", "avg", avg, "variance", variance, "ops", b.ops, "elapsed",time.Since(start).Seconds(),  "avgs", float64(b.ops)/float64(time.Since(start).Seconds()))
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



