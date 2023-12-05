/*
 * Copyright (c) 2023. Liero AB
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
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hashicorp/go-hclog"

	"go.opentelco.io/go-swpx/core"
	"go.opentelco.io/go-swpx/proto/go/corepb"
)

var (
	ErrInvalidRequest = fmt.Errorf("invalid request body")
)

// Poll is the request that holds the TI request.
type Poll struct {
	AccessId string `json:"access_id"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`

	NetworkRegion string `json:"network_region"`

	Provider []string `json:"provider"`
	Driver   string   `json:"driver"` // optional, need to be able to set with provider

	Region        string          `json:"region"`
	RecreateIndex bool            `json:"recreate_index"`
	Type          string          `json:"type"`
	Timeout       TimeoutDuration `json:"timeout"`
	CacheTTL      TimeoutDuration `json:"cache_ttl"`

	logger hclog.Logger
}

func (req *Poll) Bind(r *http.Request) error {
	return nil
}

// Parse the incoming request
func (r *Poll) Parse() error {

	if r.AccessId == "" && r.Hostname == "" {
		return fmt.Errorf("access_id and hostname cannot both be empty: %w", ErrInvalidRequest)
	}

	return nil
}

type PollService struct {
	*chi.Mux
	core   *core.Core
	logger hclog.Logger
}

func NewPollService(core *core.Core, logger hclog.Logger) *PollService {
	h := &PollService{
		Mux: chi.NewRouter(),

		core:   core,
		logger: logger,
	}
	h.Post("/", h.Poll)
	return h
}

// Poll is the ti
func (s *PollService) Poll(w http.ResponseWriter, r *http.Request) {
	data := &Poll{
		logger: s.logger,
	}

	if err := render.Bind(r, data); err != nil {
		s.logger.Error(err.Error())
		render.JSON(w, r, NewResponse(ErrorStatusInvalidAddr, err))
		return
	}
	if err := data.Parse(); err != nil {
		render.JSON(w, r, NewResponse(ErrorStatusInvalidAddr, err))
		return
	}

	// set the Type
	t := corepb.PollRequest_Type_value[data.Type]
	pbType := corepb.PollRequest_Type(t)

	if pbType == corepb.PollRequest_NOT_SET {
		pbType = corepb.PollRequest_GET_TECHNICAL_INFO
	}

	ctx, cancel := context.WithTimeout(r.Context(), data.Timeout.Duration)
	defer cancel()

	if data.NetworkRegion == "" {
		render.JSON(w, r, NewResponse(ErrInvalidArgument("network_region empty/nil"), fmt.Errorf("network_region cannot be empty")))
		return
	}

	request := &corepb.PollRequest{
		Settings: &corepb.Settings{
			ProviderPlugin: data.Provider,
			ResourcePlugin: data.Driver,
			RecreateIndex:  data.RecreateIndex,
			Timeout:        data.Timeout.String(),
			CacheTtl:       data.CacheTTL.String(),
		},
		Session: &corepb.SessionRequest{
			NetworkRegion: data.NetworkRegion,
			AccessId:      data.AccessId, // if set Hostname and port might be overwritten by the provider plugin.PreHandler()
			Hostname:      data.Hostname,
			Port:          data.Port,
		},

		Type: pbType,
	}

	// send the request
	resp, err := s.core.PollNetworkElement(ctx, request)
	if err != nil {
		render.JSON(w, r, NewResponse(ResponseStatusNotFound, err))
		return
	}

	wrappedResponse := NewResponse(ResponseStatusOK, resp)
	render.JSON(w, r, wrappedResponse)

}
