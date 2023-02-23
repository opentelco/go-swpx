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
	"time"

	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
)

// SendRequest sends the request!
func (c *Core) SendRequest(ctx context.Context, request *pb_core.Request) (*pb_core.Response, error) {

	cacheTTLduration, _ := time.ParseDuration(request.Settings.CacheTtl)

	if !request.Settings.RecreateIndex && cacheTTLduration != 0 {
		cr, err := c.responseCache.Pop(ctx, request.Hostname, request.Port, request.AccessId, request.Type)
		if err != nil {
			c.logger.Warn("could not pop from cache", "error", err)
		}
		// if a cached response exists
		if cr != nil {
			c.logger.Debug("found cached item", "age", time.Since(cr.Timestamp))
			if time.Since(cr.Timestamp) < cacheTTLduration {
				c.logger.Info("found response in cache")
				return cr.Response, nil
			}

		}
	}

	resp, err := c.doRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	if err := c.responseCache.Upsert(ctx, request.Hostname, request.Port, request.AccessId, request.Type, resp); err != nil {
		c.logger.Error("error saving response to cache", "error", err.Error())
	}
	return resp, nil

}
