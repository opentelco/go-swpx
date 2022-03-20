package core

import (
	"context"

	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
	"github.com/hashicorp/go-hclog"
)

// worker that does the work
type worker struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Pending        int    `json:"pending"`
	Executed       int    `json:"executed"`
	TimedOut       int    `json:"timed_out"`
	requests       chan *Request
	requestHandler RequestHandler
	logger         hclog.Logger
}

type workers []*worker

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

// start the worker and ready it to accept payloads
func (w *worker) start(done chan *worker) {
	w.logger.Debug("staring worker", "wid", w.Id)
	for req := range w.requests {
		resp := &pb_core.Response{
			// RequestAccessId: req.AccessId,
		}

		// do work with payload and return err if timeout or any other error
		err := w.handle(req.ctx, req, resp)

		if err != nil {
			resp.Error = &pb_core.Error{Message: err.Error()}
			w.TimedOut++
		}

		req.Response <- resp
		w.Executed++

		w.logger.Debug("put the response back in queue.")

		done <- w
	}
}

// handle is the worker handler that wraps the actual handler, aborts if timeout is reached
func (w *worker) handle(ctx context.Context, msg *Request, resp *pb_core.Response) error {

	// execute the passed fnc function
	c := make(chan error, 1)
	go func() {
		c <- w.requestHandler(ctx, msg, resp)
	}()

	select {
	case <-ctx.Done():
		w.logger.Error("timeout reached or context cancelled", "error", ctx.Err())
		return ctx.Err()
	case err := <-c:
		if err != nil {
			w.logger.Error("error in response", "error", err.Error())
		}
		return err
	}
}
