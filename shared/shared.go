package shared

import (
	"fmt"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
)

const (
	PLUGIN_VERSION     uint   = 1
	VERSION_STRING     string = "1.0-beta"
	CORE_NAME          string = "swpx"
	MAGIC_COOKIE_KEY   string = "SWPX_PLUGIN_COOKE"
	MAGIC_COOKIE_VALUE string = "IVbRejBEZ6EnzwoR4wUz"

	PLUGIN_RESOURCE_KEY string = "resource"
	PLUGIN_PROVIDER_KEY string = "provider"
)

var VERSION *version.Version
var DEFAULT_CORE_PATH string = fmt.Sprintf(".%s", CORE_NAME)

func init() {
	VERSION, _ = version.NewVersion(VERSION_STRING)
}

type DNCRequest struct {
	OriginalID string
}

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  PLUGIN_VERSION,
	MagicCookieKey:   MAGIC_COOKIE_KEY,
	MagicCookieValue: MAGIC_COOKIE_VALUE,
}

// pluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	PLUGIN_RESOURCE_KEY: &ResourcePlugin{},
	PLUGIN_PROVIDER_KEY: &ProviderPlugin{},
}
