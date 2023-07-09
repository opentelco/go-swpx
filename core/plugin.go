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
	"sync"

	"git.liero.se/opentelco/go-swpx/shared"
)

type providerMap map[string]shared.Provider

func (p providerMap) Slice() (ps []shared.Provider) {
	for _, v := range p {
		ps = append(ps, v)
	}
	return
}

// sorters
type resourceMap map[string]shared.Resource

var ()

// ResourcePlugins is the container for resource plugins
type ResourcePlugins struct {
	Plugins map[string]shared.Resource
	*sync.Mutex
}

// ProviderPlugins is the container for provider plugins
type ProviderPlugins struct {
	Plugins map[string]shared.Provider
	*sync.Mutex
}

// NewProviderPlugins creates a new ProviderPlugins struct and initates the plugin-map
func NewProviderPlugins() *ProviderPlugins {
	return &ProviderPlugins{
		Plugins: make(map[string]shared.Provider),
	}
}

// NewResourcePlugins creates a new ResourcePlugins sruct and initates the plugin-map
func NewResourcePlugins() *ResourcePlugins {
	return &ResourcePlugins{
		Plugins: make(map[string]shared.Resource),
	}
}
