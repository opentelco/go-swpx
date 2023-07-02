package core

import (
	"context"
	"fmt"

	"git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/proto/go/resource"
)

func (c *Core) CollectConfig(ctx context.Context, request *core.CollectConfigRequest) (*core.CollectConfigResponse, error) {

	selectedProviders, err := c.selectProviders(ctx, request.Settings)
	if err != nil {
		return nil, err
	}

	if len(selectedProviders) == 0 {

		return nil, NewError("no providers selected or selected not found", ErrCodeInvalidProvider)
	}

	// pre-process the request with the selected providers
	err = c.providerGenericPreProccessing(ctx, request.Session, request.Settings, selectedProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to pre-proccess request: %w", err)
	}

	c.logger.Debug("request processed",
		"hostname", request.Session.Hostname,
		"resource-plugin", request.Settings.ResourcePlugin,
	)

	// select resource-plugin to send the requests to
	c.logger.Info("selected resource plugin", "plugin", request.Settings.ResourcePlugin)

	plugin := c.resources[request.Settings.ResourcePlugin]
	if plugin == nil {
		return nil, NewError("selected resource plugin is missing/does not exist", ErrCodeInvalidResource)
	}

	resp, err := plugin.GetRunningConfig(ctx, &resource.GetRunningConfigParameters{
		Hostname: request.Session.Hostname,
		Timeout:  request.Session.Hostname,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get running config: %w", err)
	}

	// compare with stored config from database
	return &core.CollectConfigResponse{
		Config:  resp.Config,
		Changes: []*core.ConfigChange{},
	}, nil

}
