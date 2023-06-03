package core

import (
	"context"
	"fmt"

	"git.liero.se/opentelco/go-swpx/proto/go/core"
	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/proto/go/resource"
)

// Discover will do basic discovery of a device and return basic info about it
// This can be used with any resource plugin but is not guaranteed to work with all.
// prefered way is to use the generic plugin to for all discovery. For best result no
// providers should be selected either as it could lead to unexpected results.
func (c *Core) Discover(ctx context.Context, request *pb_core.DiscoverRequest) (*pb_core.DiscoverResponse, error) {

	selectedProviders, err := c.selectProviders(ctx, request.Settings)
	if err != nil {
		return nil, err
	}

	// pre-process the request with the selected providers
	err = c.providerGenericPreProccessing(ctx, request.Session, request.Settings, selectedProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to pre-proccess request: %w", err)
	}

	// select resource-plugin to send the requests to
	c.logger.Info("selected resource plugin", "plugin", request.Settings.ResourcePlugin)

	plugin := c.resources[request.Settings.ResourcePlugin]
	if plugin == nil {
		return nil, NewError("selected resource plugin is missing/does not exist", ErrCodeInvalidResource)
	}

	resp, err := plugin.Discover(ctx, &resource.Request{
		Hostname: request.Session.Hostname,
		Timeout:  request.Session.Hostname,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get running config: %w", err)
	}

	// compare with stored config from database
	return &core.DiscoverResponse{
		NetworkElement: resp,
	}, nil

}
