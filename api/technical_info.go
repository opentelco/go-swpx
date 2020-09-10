package api

import (
	"context"
	"encoding/json"
	"git.liero.se/opentelco/go-swpx/core"
	"git.liero.se/opentelco/go-swpx/errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net"
	"net/http"
	"time"
)

// TechnicalInformationRequest is the request that holdes the TI request.
type TechnicalInformationRequest struct {
	Hostname     string          `json:"hostname"`
	Port         string          `json:"port"`
	Provider     string          `json:"provider"`
	Driver       string          `json:"driver"`
	Region       string          `json:"region"`
	DontUseIndex bool            `json:"dont_use_index"`
	Timeout      TimeoutDuration `json:"timeout"`
	TTL          TimeoutDuration `json:"ttl"`
	ipAddr       []net.IP
}

func (req *TechnicalInformationRequest) Bind(r *http.Request) error {
	return nil
}

func (r *TechnicalInformationRequest) parseDriver() error {
	return nil
}

func (r *TechnicalInformationRequest) parseAddr() error {
	// Parse hostname/ip for host
	addrs, err := net.LookupHost(r.Hostname)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	for _, addr := range addrs {
		addr := net.ParseIP(addr)
		logger.Info("addr:", addr.String())
		if addr == nil {
			r.ipAddr = append(r.ipAddr, addr)
		} else {
		}
	}

	return nil
}

// Parse the incoming request
func (r *TechnicalInformationRequest) Parse() error {
	if err := r.parseAddr(); err != nil {
		logger.Error(err.Error())
		return errors.New(err.Error(), errors.ErrInvalidAddr)
	}
	return nil
}

type ServiceTechnicalInformation struct {
	*chi.Mux

	requests chan *core.Request
	storage  interface{}
}

func NewServiceTechnicalInformation(requestChan chan *core.Request) *ServiceTechnicalInformation {
	h := &ServiceTechnicalInformation{
		Mux: chi.NewRouter(),

		requests: requestChan,
	}
	h.Post("/", h.GetTI)
	return h
}

// GetTI is the ti
func (s *ServiceTechnicalInformation) GetTI(w http.ResponseWriter, r *http.Request) {
	data := &TechnicalInformationRequest{}

	if err := render.Bind(r, data); err != nil {
		logger.Error(err.Error())
		render.JSON(w, r, NewResponse(ErrorStatusInvalidAddr, err))
		return
	}
	ti, _ := json.Marshal(data)
	logger.Info("TI:", string(ti))
	if err := data.Parse(); err != nil {
		render.JSON(w, r, NewResponse(ErrorStatusInvalidAddr, err))
		return
	}

	ctx, _ := context.WithTimeout(r.Context(), data.Timeout.Duration)
	req := &core.Request{
		NetworkElement: data.Hostname,
		Provider:       data.Provider,
		Resource:       data.Driver,
		DontUseIndex:   data.DontUseIndex,

		// Metadata
		Response: make(chan *core.Response, 1),
		Context:  ctx,
	}

	if data.Port != "" {
		req.NetworkElementInterface = &data.Port
		req.Type = core.GetTechnicalInformationPort
	} else {
		req.Type = core.GetTechnicalInformationElement
	}

	// check response cache before sending request
	cachedResponse, err := core.ResponseCache.PopResponse(req.NetworkElement, *req.NetworkElementInterface, req.Type)
	if err != nil {
		logger.Error("error popping from cache: ", err.Error())
		render.JSON(w, r, NewResponse(ResponseStatusError, err.Error()))
		return
	}

	if cachedResponse != nil {
		if time.Since(cachedResponse.Timestamp) < data.TTL.Duration {
			logger.Info("found response in cache")
			render.JSON(w, r, NewResponse(ResponseStatusOK, cachedResponse.Response.Data))
			return
		}
		// if response is cached but ttl ran out, clear it from the cache
		if err := core.ResponseCache.Clear(req.NetworkElement, *req.NetworkElementInterface, req.Type); err != nil {
			logger.Error("error clearing cache:", err)
		}
	}

	// send the request
	s.requests <- req

	// handle it
	for {
		select {
		case resp := <-req.Response:
			if resp.Error != nil {
				render.JSON(w, r, NewResponse(ResponseStatusNotFound, resp.Error))
				return
			}

			wrappedResponse := NewResponse(ResponseStatusOK, resp)
			if err = core.ResponseCache.SetResponse(req.NetworkElement, *req.NetworkElementInterface, req.Type, wrappedResponse); err != nil {
				logger.Error("error saving response to cache: ", err.Error())
			}

			render.JSON(w, r, wrappedResponse)
			return
		case <-req.Context.Done():
			logger.Info("timeout for request was hit")

			render.JSON(w, r, NewResponse(ErrorStatusRequestTimeout, nil))
			return
		}

	}
}
