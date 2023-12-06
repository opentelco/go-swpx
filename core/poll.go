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

package core

import (
	"context"
	"fmt"
	"time"

	"go.opentelco.io/go-swpx/proto/go/corepb"
)

// PollDevice sends the request!
func (c *Core) PollDevice(ctx context.Context, request *corepb.PollRequest) (*corepb.PollResponse, error) {

	c.logger.Debug("polling network element",
		"hostname", request.Session.Hostname,
		"port", request.Session.Port,
		"accessId", request.Session.AccessId,
		"type", request.Type,
		"region", request.Session.NetworkRegion,
		"recreateIndex", request.Settings.RecreateIndex,
		"cacheTTL", request.Settings.CacheTtl,
		"timeout", request.Settings.Timeout,
	)
	cacheTTLduration, _ := time.ParseDuration(request.Settings.CacheTtl)

	if !request.Settings.RecreateIndex && cacheTTLduration != 0 {
		cr, err := c.pollResponseCache.Pop(ctx, request.Session.Hostname, request.Session.Port, request.Session.AccessId, request.Type)
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

	timeoutDur, _ := time.ParseDuration(request.Settings.Timeout)
	if timeoutDur == 0 {
		c.logger.Debug("using default timeout, since none was specified", "timeout", c.config.Request.DefaultRequestTimeout.AsDuration())
		timeoutDur = c.config.Request.DefaultRequestTimeout.AsDuration()
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, timeoutDur)
	defer cancel()

	var resp *corepb.PollResponse
	respCh := make(chan *corepb.PollResponse)
	errCh := make(chan error)
	go func() {
		resp, err := c.doPollRequest(ctx, request)
		if err != nil {
			errCh <- err
		}
		respCh <- resp
	}()

	select {
	case <-ctxTimeout.Done():
		return nil, fmt.Errorf("timeout reached for requet: %w", ctxTimeout.Err())
	case err := <-errCh:
		return nil, fmt.Errorf("could not complete request: %w", err)

	case resp = <-respCh:
		if err := c.pollResponseCache.Upsert(ctx, request.Session.Hostname, request.Session.Port, request.Session.AccessId, request.Type, resp); err != nil {
			c.logger.Warn("error saving response to cache", "error", err.Error())
		}
		return resp, nil
	}

}
