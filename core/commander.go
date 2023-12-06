package core

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.opentelco.io/go-swpx/proto/go/resourcepb"
	"go.opentelco.io/go-swpx/proto/go/stanzapb"
)

type commander struct {
	service corepb.CommanderServer

	core *Core

	logger hclog.Logger

	corepb.UnimplementedCommanderServer
}

func NewCommander(core *Core, logger hclog.Logger) (corepb.CommanderServer, error) {
	return &commander{
		core:   core,
		logger: logger,
	}, nil
}

func (c *commander) ConfigureStanza(ctx context.Context, req *corepb.ConfigureStanzaRequest) (*stanzapb.ConfigureResponse, error) {
	selectedProviders, err := c.core.selectProviders(ctx, req.Settings)
	if err != nil {
		return nil, err
	}

	if req.Settings.ResourcePlugin == "" {
		req.Settings.ResourcePlugin, err = c.core.resolveResourcePlugin(ctx, req.Session, selectedProviders)
		if err != nil {
			return nil, fmt.Errorf("ConfigureStanza: could not resolve resource plugin: %w", err)
		}
	}

	req.Session, err = c.core.resolveSession(ctx, req.Session, selectedProviders)
	if err != nil {
		return nil, fmt.Errorf("ConfigureStanza: could not resolve resource session request: %w", err)
	}

	c.logger.Debug("ConfigureStanza: request processed",
		"hostname", req.Session.Hostname,
		"resource-plugin", req.Settings.ResourcePlugin,
	)

	// select resource-plugin to send the requests to
	c.logger.Info("ConfigureStanza: selected resource plugin", "plugin", req.Settings.ResourcePlugin)

	plugin := c.core.resources[req.Settings.ResourcePlugin]
	if plugin == nil {
		return nil, NewError("ConfigureStanza: selected resource plugin is missing/does not exist", ErrCodeInvalidResource)
	}

	res, err := plugin.ConfigureStanza(ctx, &resourcepb.ConfigureStanzaRequest{
		Hostname:      req.Session.Hostname,
		Timeout:       req.Settings.Timeout,
		NetworkRegion: req.Session.NetworkRegion,
		Stanza:        req.Stanza,
	})
	if err != nil {
		return nil, fmt.Errorf("ConfigureStanza: failed to configure stanza: %w", err)
	}

	c.logger.Debug("ConfigureStanza: stanza configured", "result", res)
	return res, nil
}
