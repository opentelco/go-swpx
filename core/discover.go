package core

import (
	"context"
	"fmt"

	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.opentelco.io/go-swpx/proto/go/resourcepb"
)

// Discover will do basic discovery of a device and return basic info about it
// This can be used with any resource plugin but is not guaranteed to work with all.
// prefered way is to use the generic plugin to for all discovery. For best result no
// providers should be selected either as it could lead to unexpected results.
func (c *Core) Discover(ctx context.Context, request *corepb.DiscoverRequest) (*corepb.DiscoverResponse, error) {

	selectedProviders, err := c.selectProviders(ctx, request.Settings)
	if err != nil {
		return nil, err
	}

	if request.Settings.ResourcePlugin == "" {
		request.Settings.ResourcePlugin, err = c.resolveResourcePlugin(ctx, request.Session, selectedProviders)
		if err != nil {
			return nil, fmt.Errorf("could not resolve resource plugin: %w", err)
		}
	}

	request.Session, err = c.resolveSession(ctx, request.Session, selectedProviders)
	if err != nil {
		return nil, fmt.Errorf("could not resolve resource session request: %w", err)
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

	resp, err := plugin.Discover(ctx, &resourcepb.Request{
		NetworkRegion: request.Session.NetworkRegion,
		Hostname:      request.Session.Hostname,
		Timeout:       request.Settings.Timeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get running config: %w", err)
	}

	// compare with stored config from database
	return &corepb.DiscoverResponse{
		Device: resp,
	}, nil

}
