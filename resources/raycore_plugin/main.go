/*
 * Copyright (c) 2020. Liero AB
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

package main

import (
	"fmt"
	"log"
	"os"

	"git.liero.se/opentelco/go-swpx/shared"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
)

var VERSION *version.Version

const (
	VERSION_BASE string = "1.0-beta"
	DRIVER_NAME  string = "raycore-driver"
)

func init() {
	var err error
	if VERSION, err = version.NewVersion(VERSION_BASE); err != nil {
		log.Fatal(err)
	}
}

// Here is a real implementation of Driver
type RaycoreDriver struct {
	logger hclog.Logger
}

func (g *RaycoreDriver) Version() (string, error) {
	// g.logger.Debug("message from resource-driver running at version:", VERSION)
	return fmt.Sprintf("%s@%s", DRIVER_NAME, VERSION.String()), nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   shared.MagicCookieKey,
	MagicCookieValue: shared.MagicCookieValue,
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       fmt.Sprintf("%s@%s", DRIVER_NAME, VERSION.String()),
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	resource := &RaycoreDriver{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"resource": &shared.ResourceDriverPlugin{Impl: resource},
	}

	logger.Debug("message from resource-driver", VERSION.String())

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
