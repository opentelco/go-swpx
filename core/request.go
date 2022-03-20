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
	"context"
	"fmt"
	"time"

	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
)

// Request is the internal representation of an incoming request
// it is passed between the api and the core
type Request struct {
	*pb_core.Request

	Response chan *pb_core.Response

	ctx context.Context
}

func NewRequest(ctx context.Context, request *pb_core.Request) *Request {
	return &Request{
		Request: request,
		// Metadata
		Response: make(chan *pb_core.Response, 1),
		ctx:      ctx,
	}
}

// GetCacheTTL is a helper function
func (r *Request) GetCacheTTL() time.Duration {
	ttl, err := time.ParseDuration(r.Settings.CacheTtl)
	if err != nil {
		return 0
	}
	return ttl
}

// SendRequest to CORE
func (c *Core) SendRequest(ctx context.Context, request *Request) (*pb_core.Response, error) {

	if !request.Settings.RecreateIndex && request.GetCacheTTL() != 0 {
		cr, err := c.responseCache.Pop(ctx, request.Hostname, request.Port, request.AccessId, request.Type)
		if err != nil {
			c.logger.Warn("could not pop from cache", "error", err)
		}
		// if a cached response exists
		if cr != nil {
			c.logger.Debug("found cached item", "age", time.Since(cr.Timestamp))
			if time.Since(cr.Timestamp) < request.GetCacheTTL() {
				c.logger.Info("found response in cache")
				return cr.Response, nil
			}

		}
	}

	RequestQueue <- request
	// cache is not set
	for {
		select {
		case resp := <-request.Response:
			if resp == nil {
				return nil, fmt.Errorf("received emptry response for: %s", request.Type.String())
			}
			if resp.Error != nil && resp.Error.Message != "" {
				return nil, fmt.Errorf("received error in response, (%d): %s", resp.Error.Code, resp.Error.Message)
			}

			if err := c.responseCache.Upsert(ctx, request.Hostname, request.Port, request.AccessId, request.Type, resp); err != nil {
				c.logger.Error("error saving response to cache", "error", err.Error())
			}

			return resp, nil

		case <-request.ctx.Done():
			c.logger.Error("timeout for request was hit")
			return nil, fmt.Errorf("timeout for request reached")

		}
	}
}
