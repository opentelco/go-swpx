package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/opentelco/go-swpx/shared"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
)

var VERSION *version.Version
var logger hclog.Logger

const (
	VERSION_BASE  string = "1.0-beta"
	PROVIDER_NAME string = "ssab_provider"
	WEIGHT        int64  = 30
)

func init() {
	var err error
	if VERSION, err = version.NewVersion(VERSION_BASE); err != nil {
		log.Fatal(err)
	}
}

// Here is a real implementation of Greeter
type Provider struct {
	logger hclog.Logger
}

func (g *Provider) Version() (string, error) {
	g.logger.Debug("message from provider, running version:", VERSION)
	return fmt.Sprintf("%s@%s", PROVIDER_NAME, VERSION.String()), nil
}

func (g *Provider) Weight() (int64, error) {
	// g.logger.Debug("return the weight", PROVIDER_NAME)
	return WEIGHT, nil
}

func (p *Provider) Lookup(id string) (string, error) {
	duration := time.Duration(time.Millisecond * time.Duration(rand.Intn(1000-400)+400))
	time.Sleep(duration)
	p.logger.Debug("done", "execution_time", duration.String())
	return "nonon", nil
}

func (p *Provider) Match(id string) (bool, error) {
	return true, nil
}

func (p *Provider) Name() (string, error) {
	return PROVIDER_NAME, nil

}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   shared.MAGIC_COOKIE_KEY,
	MagicCookieValue: shared.MAGIC_COOKIE_VALUE,
}

func main() {
	logger = hclog.New(&hclog.LoggerOptions{
		Name:       fmt.Sprintf("%s@%s", PROVIDER_NAME, VERSION.String()),
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			shared.PLUGIN_PROVIDER_KEY: &shared.ProviderPlugin{Impl: &Provider{logger: logger}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
