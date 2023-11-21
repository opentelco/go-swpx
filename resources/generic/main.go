package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"git.liero.se/opentelco/go-dnc/client"
	"git.liero.se/opentelco/go-swpx/config"
	"git.liero.se/opentelco/go-swpx/proto/go/resourcepb"
	"git.liero.se/opentelco/go-swpx/proto/go/stanzapb"
	"git.liero.se/opentelco/go-swpx/shared"
	"github.com/gogo/status"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
	"google.golang.org/grpc/codes"
)

var VERSION *version.Version

var logger hclog.Logger

const (
	VersionBase            = "1.0-beta"
	DriverName             = "generic"
	defaultDeadlineTimeout = 90 * time.Second
)

func init() {
	var err error
	if VERSION, err = version.NewVersion(VersionBase); err != nil {
		log.Fatal(err)
	}
}

type driver struct {
	logger hclog.Logger
	dnc    client.Client
	conf   *config.ResourceGeneric
}

func main() {
	logger = hclog.New(&hclog.LoggerOptions{
		Name:       fmt.Sprintf("%s@%s", DriverName, VERSION.String()),
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	var appConf config.ResourceGeneric
	configPath := os.Getenv("GENERIC_CONFIG_FILE")
	if configPath == "" {
		configPath = "generic.hcl"
	}
	err := config.LoadConfig(configPath, &appConf)
	if err != nil {
		logger.Error("failed to loadd.config", "error", err)
		os.Exit(1)
	}

	logger.Info("connecting to DNC", "address", appConf.DNC.Addr)
	dncClient, err := client.NewGRPC(appConf.DNC.Addr)
	if err != nil {
		logger.Error("failed to create DNC Client", "error", err)
		os.Exit(1)
	}
	driver := &driver{
		logger: logger,
		dnc:    dncClient,
		conf:   &appConf,
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			shared.PluginResourceKey: &shared.ResourcePlugin{Impl: driver},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

func validateEOLTimeout(reqTimeout string, defaultDuration time.Duration) time.Duration {
	dur, _ := time.ParseDuration(reqTimeout)

	if dur.Seconds() == 0 {
		dur = defaultDuration

	}

	return dur
}

func (d *driver) ConfigureStanza(ctx context.Context, req *resourcepb.ConfigureStanzaRequest) (*stanzapb.ConfigureResponse, error) {
	return nil, status.New(codes.Unimplemented, "not implemented").Err()
}
