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
	"context"
	"encoding/json"
	"net"
	"net/http"
	
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"
	
	"git.liero.se/opentelco/go-swpx/core"
	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
)

// TechnicalInformationRequest is the request that holdes the TI request.
type TechnicalInformationRequest struct {
	Hostname      string          `json:"hostname"`
	Port          string          `json:"port"`
	Provider      string          `json:"provider"`
	Driver        string          `json:"driver"`
	Region        string          `json:"region"`
	RecreateIndex bool            `json:"recreate_index"`
	Timeout       TimeoutDuration `json:"timeout"`
	CacheTTL      TimeoutDuration `json:"cache_ttl"`
	ipAddr        []net.IP
	
	logger hclog.Logger
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
		r.logger.Error(err.Error())
		return err
	}

	for _, addr := range addrs {
		addr := net.ParseIP(addr)
		r.logger.Info("addr:", addr.String())
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
		r.logger.Error(err.Error())
		return core.NewError(err.Error(), core.ErrInvalidAddr)
	}
	return nil
}

type ServiceTechnicalInformation struct {
	*chi.Mux
	core *core.Core
	logger hclog.Logger
	storage  interface{}
}

func NewServiceTechnicalInformation(core *core.Core, logger hclog.Logger) *ServiceTechnicalInformation {
	h := &ServiceTechnicalInformation{
		Mux: chi.NewRouter(),
		
		core: core,
		logger: logger,
		
	}
	h.Post("/", h.GetTI)
	return h
}

// GetTI is the ti
func (s *ServiceTechnicalInformation) GetTI(w http.ResponseWriter, r *http.Request) {
	data := &TechnicalInformationRequest{
		logger: s.logger,
	}

	if err := render.Bind(r, data); err != nil {
		logger.Error(err.Error())
		render.JSON(w, r, NewResponse(ErrorStatusInvalidAddr, err))
		return
	}
	ti, _ := json.Marshal(data)
	s.logger.Info("TI:", string(ti))
	if err := data.Parse(); err != nil {
		render.JSON(w, r, NewResponse(ErrorStatusInvalidAddr, err))
		return
	}

	ctx, _ := context.WithTimeout(r.Context(), data.Timeout.Duration)
	req := &core.Request{
		Request: &pb_core.Request{
			ProviderPlugin:         data.Provider,
			ResourcePlugin:         data.Driver,
			RecreateIndex:          data.RecreateIndex,
			DisableDistributedLock: false,
			Timeout:                data.Timeout.String(),
			CacheTtl:               data.CacheTTL.String(),
			Hostname:               data.Hostname,
			Port:                   data.Port,
			Type: pb_core.Request_GET_TECHNICAL_INFO,
		},

		// Metadata
		Response: make(chan *pb_core.Response, 1),
		Context:  ctx,
	}
	
	if data.Port != "" {
		req.Type = pb_core.Request_GET_TECHNICAL_INFO_PORT
		// check response cache before sending request
	}
	
	// send the request
	
	resp, err := s.core.SendRequest(ctx, req)
	if err != nil {
		render.JSON(w, r, NewResponse(ResponseStatusNotFound, resp.Error))
	}
	wrappedResponse := NewResponse(ResponseStatusOK, resp)
	render.JSON(w, r, wrappedResponse)

	// handle it

}

