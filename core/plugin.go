package core

import (
	"sync"

	"git.liero.se/opentelco/go-swpx/shared"
	"github.com/hashicorp/go-plugin"
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
type byWeight []shared.Provider

func (p byWeight) Less(i int, j int) bool {
	pi, _ := p[i].Weight()
	pj, _ := p[j].Weight()
	return pi < pj
}
func (p byWeight) Len() int          { return len(p) }
func (p byWeight) Swap(i int, j int) { p[i], p[j] = p[j], p[i] }

var (
	availableResources map[string]*plugin.Client = make(map[string]*plugin.Client)
	availableProviders map[string]*plugin.Client = make(map[string]*plugin.Client)

	resources       resourceMap = make(map[string]shared.Resource)
	providers       providerMap = make(map[string]shared.Provider)
	sortedProviders []shared.Provider
)

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
