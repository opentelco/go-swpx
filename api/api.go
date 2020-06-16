package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"git.liero.se/opentelco/go-swpx/core"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/gorilla/context"

	"github.com/go-chi/chi"
)

var router *chi.Mux

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

	requests chan *core.Request
}

func (s *Server) ListenAndServe(host string) error {
	log.Printf("Listen on %s\n", host)
	return http.ListenAndServe(host, context.ClearHandler(s))
}

// New creates and returnes a server.
func New(requestQueue chan *core.Request) *Server {
	if requestQueue == nil {
		log.Fatal("channel is nil, requests needs to be handled..")
	}
	s := &Server{
		Mux:      chi.NewRouter(),
		requests: requestQueue,
	}
	s.Use(middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		render.SetContentType(render.ContentTypeJSON),
		middleware.Heartbeat("/ping"))

	// Mount the default path
	s.Get("/", func(w http.ResponseWriter, r *http.Request) {
		Render(w, r, NewResponse(ResponseStatusNotFound, nil))
	})

	s.Route("/v1", func(r chi.Router) {
		r.Mount("/ti", NewServiceTechnicalInformation(s.requests))
	})
	return s
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
