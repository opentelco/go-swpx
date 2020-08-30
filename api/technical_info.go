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
	"git.liero.se/opentelco/go-swpx/core"
	"git.liero.se/opentelco/go-swpx/errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net"
	"net/http"
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

	ipAddr []net.IP
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
			render.JSON(w, r, NewResponse(ResponseStatusOK, resp))
			return
		case <-req.Context.Done():
			logger.Info("timeout for request was hit")

			render.JSON(w, r, NewResponse(ErrorStatusRequestTimeout, nil))
			return
		}

	}

	render.JSON(w, r, NewResponse(ResponseStatusOK, ""))

}
