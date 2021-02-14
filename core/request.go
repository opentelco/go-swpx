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

type _Metadata struct{}

// TODO: Change Request to interface so we can handl different types of requests?
// - ProvideRequest
//    Provide.CPE()
//    Provide.XXX_()
// - Core.Poll?
//    - PollRequest with different Types?
//       - GET_TECHNICAL_INFO
//       - GET_TECHNICAL_INFO_PORT
//       - GET_MAC_TABLE...
type _Request interface {
	// Get how long the requset has left to live
	TTL() time.Time

	// Get the value of Timeout
	Timeout() time.Duration

	// Get the raw underlaying request
	Raw() interface{}
}

// Request is the internal representation of a incoming request
// it is passed between the api and the core
type Request struct {
	*pb_core.Request
	// metadata to handle the request

	Response chan *pb_core.Response
	Context  context.Context
}

// GetCacheTTL is a helper function
func (r *Request) GetCacheTTL() time.Duration {
	ttl, err := time.ParseDuration(r.CacheTtl)
	if err != nil {
		return 0
	}
	return ttl
}

// SendRequest to CORE
func (c *Core) SendRequest(ctx context.Context, request *Request) (*pb_core.Response, error) {

	// TODO: move to handler, provider calls should not care about this?
	// if recreate is set no use to get the cache
	if !request.RecreateIndex && request.GetCacheTTL() != 0 {
		cr, err := CacheResponse.Pop(ctx, request.Hostname, request.Port, request.Type)
		if err != nil {
			c.logger.Warn("could not pop from cache", "error", err)
		}
		// if a cached response exists
		if cr != nil {
			if time.Since(cr.Timestamp.AsTime()) < request.GetCacheTTL() {
				c.logger.Info("found response in cache")
				return cr.Response, nil
			}
			// if response is cached but ttl ran out, clear it from the cache
			if err := CacheResponse.Clear(context.TODO(), request.Hostname, request.Port, request.Type); err != nil {
				logger.Error("error clearing cache:", err)
			}
		}
	}

	RequestQueue <- request
	// cache is not set
	for {
		select {
		case resp := <-request.Response:
			if err := CacheResponse.Upsert(context.TODO(), request.Hostname, request.Port, request.Type, resp); err != nil {
				logger.Error("error saving response to cache: ", err.Error())
			}

			return resp, nil
		case <-request.Context.Done():
			c.logger.Error("timeout for request was hit")
			return nil, fmt.Errorf("timeout for request reached")

		}
	}
}
