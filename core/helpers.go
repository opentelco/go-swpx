package core

import (
	"context"
	"strings"

	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.opentelco.io/go-swpx/shared"
)

func (c *Core) selectProviders(ctx context.Context, settings *corepb.Settings) ([]shared.Provider, error) {
	var selectedProviders []shared.Provider
	// check if a providers are selected in the request
	if len(settings.ProviderPlugin) > 0 {
		c.logger.Info("request has selected providers", "providers", strings.Join(settings.ProviderPlugin, ","))

		for _, provider := range settings.ProviderPlugin {
			c.logger.Debug("find provider in loaded plugins", "selected_provider", provider)
			if selectedProvider, ok := c.providers[provider]; ok {
				selectedProviders = append(selectedProviders, selectedProvider)
			}
		}

	} else {
		c.logger.Debug("request has no selected provider and default provider is set in config", "default_provider", c.config.Request.DefaultProvider)
		if provider, ok := c.providers[c.config.Request.DefaultProvider]; ok {
			selectedProviders = append(selectedProviders, provider)
		}
	}

	return selectedProviders, nil
}

func (c *Core) resolveResourcePlugin(ctx context.Context, sess *corepb.SessionRequest, selectedProviders []shared.Provider) (string, error) {
	for _, provider := range selectedProviders {
		rp, err := provider.ResolveResourcePlugin(ctx, sess)
		if err != nil {
			return "", err
		}
		if rp != nil {
			return rp.ResourcePlugin, nil
		}
	}
	return "", nil
}

func (c *Core) resolveSession(ctx context.Context, sess *corepb.SessionRequest, selectedProviders []shared.Provider) (*corepb.SessionRequest, error) {
	for _, provider := range selectedProviders {
		s, err := provider.ResolveSessionRequest(ctx, sess)
		if err != nil {
			return sess, err
		}
		if s != nil {
			return s, nil
		}
	}
	return sess, nil
}

// providerPostProcess post-process the response with the selected providers
// appends to the response in the input argument and returns an error if something went wrong
func providerPollPostProcess(ctx context.Context, selectedProviders []shared.Provider, response *corepb.PollResponse) error {
	// PostProcess the response with the selected Providers
	for _, selectedProvider := range selectedProviders {
		if selectedProvider != nil {
			nr, err := selectedProvider.ProcessPollResponse(ctx, response)
			if err != nil {
				return nil
			}
			response.NetworkElement = nr.NetworkElement
		}
	}
	return nil
}
