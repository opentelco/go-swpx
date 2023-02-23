package shared

import (
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
)

const (
	PluginVersion    uint   = 1
	VersionString    string = "1.0-beta"
	CoreName         string = "swpx"
	MagicCookieKey   string = "SWPX_PLUGIN_COOKE"
	MagicCookieValue string = "IVbRejBEZ6EnzwoR4wUz"

	PluginResourceKey string = "resource"
	PluginProviderKey string = "provider"
)

var VERSION *version.Version

func init() {
	VERSION, _ = version.NewVersion(VersionString)
}

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  PluginVersion,
	MagicCookieKey:   MagicCookieKey,
	MagicCookieValue: MagicCookieValue,
}

// pluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	PluginResourceKey: &ResourcePlugin{},
	PluginProviderKey: &ProviderPlugin{},
}
