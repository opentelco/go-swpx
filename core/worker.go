package core

import (
	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
)

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
func (w *worker) start(done chan *worker) {
	for {
		select {
		case msg := <-w.messages:
			resp := &pb_core.Response{RequestAccessId: msg.AccessId}

			// do work with payload
			err := handle(msg.Context, msg, resp, handleMsg)

			if err != nil {
				resp.Error = &pb_core.Error{Message: err.Error()}
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
