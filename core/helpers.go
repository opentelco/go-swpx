package core

import (
	"context"
	"strings"

	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/shared"
)

func (c *Core) selectProviders(ctx context.Context, settings *pb_core.Settings) ([]shared.Provider, error) {
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

// providerPreProccessing pre-process the request with the selected providers
// writes any changes to the settings/session parameters (as they are pointers)
func (c *Core) providerGenericPreProccessing(ctx context.Context, session *pb_core.SessionRequest, settings *pb_core.Settings, selectedProviders []shared.Provider) error {
	// check if a providers are selected in the request
	for _, provider := range selectedProviders {
		var err error
		pname, _ := provider.Name()
		c.logger.Info("pre-process request", "provider-plugin", pname)

		if settings.ResourcePlugin == "" {
			rp, err := provider.ResolveResourcePlugin(ctx, session)
			if err != nil {
				return err
			}
			if rp != nil {
				settings.ResourcePlugin = rp.ResourcePlugin
			}
		}

		s, err := provider.ResolveSessionRequest(ctx, session)
		if err != nil {
			return err
		} else {
			session = s
			c.logger.Debug("session resolved", "session", session)
		}

	}

	return nil
}

// providerPostProcess post-process the response with the selected providers
// appends to the response in the input argument and returns an error if something went wrong
func providerPollPostProcess(ctx context.Context, selectedProviders []shared.Provider, response *pb_core.PollResponse) error {
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
