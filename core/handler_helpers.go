package core

import (
	"context"
	"strings"

	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/shared"
)

func (c *Core) selectProviders(ctx context.Context, request *pb_core.Request) ([]shared.Provider, error) {
	var selectedProviders []shared.Provider
	// check if a providers are selected in the request
	if len(request.Settings.ProviderPlugin) > 0 {
		c.logger.Info("request has selected providers", "providers", strings.Join(request.Settings.ProviderPlugin, ","))

		for _, provider := range request.Settings.ProviderPlugin {
			c.logger.Debug("find provider in loaded plugins", "selected_provider", provider)
			if selectedProvider, ok := c.providers[provider]; ok {
				selectedProviders = append(selectedProviders, selectedProvider)
			}
		}

	} else {
		c.logger.Info("request has selected provider and default provider is set in config", "default_provider", c.config.Request.DefaultProvider)
		if provider, ok := c.providers[c.config.Request.DefaultProvider]; ok {
			selectedProviders = append(selectedProviders, provider)
		}
	}

	return selectedProviders, nil
}

// providerPreProccessing pre-process the request with the selected providers
func (c *Core) providerPreProccessing(ctx context.Context, request *pb_core.Request, selectedProviders []shared.Provider) (*pb_core.Request, error) {
	// check if a providers are selected in the request
	for _, provider := range selectedProviders {
		var err error
		pname, _ := provider.Name()
		c.logger.Info("pre-process request with provider", "provider", pname)
		request, err = provider.PreHandler(ctx, request)
		if err != nil {
			return nil, err
		}
	}

	return request, nil
}

// providerPostProcess post-process the response with the selected providers
// appends to the response in the input argument and returns an error if something went wrong
func providerPostProcess(ctx context.Context, selectedProviders []shared.Provider, response *pb_core.Response) error {
	// PostProcess the response with the selected Providers
	for _, selectedProvider := range selectedProviders {
		if selectedProvider != nil {
			nr, err := selectedProvider.PostHandler(ctx, response)
			if err != nil {
				return nil
			}
			response.NetworkElement = nr.NetworkElement
		}
	}
	return nil
}
