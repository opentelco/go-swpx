package core

import (
	"context"
	"fmt"

	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.opentelco.io/go-swpx/proto/go/resourcepb"
	"go.opentelco.io/go-swpx/shared"
)

func (c *Core) CollectConfig(ctx context.Context, request *corepb.CollectConfigRequest) (*corepb.CollectConfigResponse, error) {

	selectedProviders, err := c.selectProviders(ctx, request.Settings)
	if err != nil {
		return nil, err
	}

	if len(selectedProviders) == 0 {

		return nil, NewError("no providers selected or selected not found", ErrCodeInvalidProvider)
	}

	var plugin shared.Resource
	plugin, err = c.resolveResourcePlugin(ctx, request.Session, request.Settings, selectedProviders)
	if err != nil {
		return nil, fmt.Errorf("could not resolve resource plugin: %w", err)
	}

	request.Session, err = c.resolveSession(ctx, request.Session, selectedProviders)
	if err != nil {
		return nil, fmt.Errorf("could not resolve resource session request: %w", err)
	}

	c.logger.Debug("request processed",
		"hostname", request.Session.Hostname,
	)

	// select resource-plugin to send the requests to
	c.logger.Info("selected resource plugin", "plugin", request.Settings.ResourcePlugin)

	resp, err := plugin.GetRunningConfig(ctx,
		&resourcepb.GetRunningConfigParameters{
			Hostname:      request.Session.Hostname,
			Timeout:       request.Settings.Timeout,
			NetworkRegion: request.Session.NetworkRegion,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to get running config: %w", err)
	}

	// compare with stored config from database
	return &corepb.CollectConfigResponse{
		Config:  resp.Config,
		Changes: []*corepb.ConfigChange{},
	}, nil

}
