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

package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/gorilla/context"

	"git.liero.se/opentelco/go-swpx/core"

	"github.com/go-chi/chi"
)

const APP_NAME = "go-swpx"

var (
	logger hclog.Logger
	router *chi.Mux
)

type TimeoutDuration struct {
	time.Duration
}

func (d *TimeoutDuration) UnmarshalJSON(b []byte) (err error) {
	d.Duration, err = time.ParseDuration(strings.Trim(string(b), `"`))
	return
}

func (d TimeoutDuration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

type Server struct {
	*chi.Mux

	core   *core.Core
	logger hclog.Logger
}

func (s *Server) ListenAndServe(host string) error {
	s.logger.Info(fmt.Sprintf("Listen on %s\n", host))
	return http.ListenAndServe(host, context.ClearHandler(s))
}

// NewServer creates and returns a server.
func NewServer(core *core.Core, logger hclog.Logger) *Server {
	if logger != nil {
		logger = hclog.New(&hclog.LoggerOptions{
			Name:   APP_NAME,
			Output: os.Stdout,
			Level:  hclog.Debug,
		})
	}

	srv := &Server{
		Mux:    chi.NewRouter(),
		core:   core,
		logger: logger,
	}

	srv.Use(middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		render.SetContentType(render.ContentTypeJSON),
		middleware.Heartbeat("/ping"))

	// Mount the default path
	srv.Get("/", func(w http.ResponseWriter, r *http.Request) {
		Render(w, r, NewResponse(ResponseStatusNotFound, nil))
	})

	srv.Route("/v1", func(r chi.Router) {
		r.Mount("/poll", NewPollService(core, logger))
	})
	return srv
}

// Request is the interface for every request.
type Request interface {
	Parse()
}

// Service is any service of the API
type Service interface {
	Validate() error
	Parse() error
	Routes() chi.Router
}

// Routers is the available routers for the API
func Routers() *chi.Mux {
	r := chi.NewRouter()
	return r
}
