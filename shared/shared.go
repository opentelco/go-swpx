package shared

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"

	proto "git.liero.se/opentelco/go-swpx/proto/provider"
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

type ConfigTelnet struct {
	User          string        `json:"username" yaml:"username" toml:"username"`
	Password      string        `json:"password" yaml:"password" toml:"password"`
	Port          int32         `json:"port" yaml:"port" toml:"port"`
	ScreenLength  string        `json:"screen_length" yaml:"screen_length" toml:"screen_length"`
	RegexPrompt   string        `json:"default_prompt" yaml:"default_prompt" toml:"default_prompt"`
	Errors        string        `json:"default_errors" yaml:"default_errors" toml:"default_errors"`
	TTL           time.Duration `json:"ttl" yaml:"ttl" toml:"ttl"`
	ReadDeadLine  time.Duration `json:"read_dead_line" yaml:"read_dead_line" toml:"read_dead_line"`
	WriteDeadLine time.Duration `json:"write_dead_line" yaml:"write_dead_line" toml:"write_dead_line"`
}

type ConfigSNMP struct {
	Community          string        `json:"community" toml:"community" yaml:"community"`
	Version            uint8         `json:"version" toml:"version" yaml:"version"`
	Port               uint16        `json:"port" toml:"port" yaml:"port"`
	Timeout            time.Duration `json:"timeout" toml:"timeout" yaml:"timeout"`
	Retries            int32         `json:"retries" toml:"retries" yaml:"retries"`
	MaxIterations      int32         `json:"max_iterations" yaml:"max_iterations" toml:"max_iterations" yaml:"max_iterations"`
	MaxRepetitions     int32         `json:"max_repetitions" yaml:"max_repetitions" toml:"max_repetitions" yaml:"max_repetitions"`
	NonRepeaters       int32         `json:"non_repeaters" yaml:"non_repeaters" toml:"non_repeaters" yaml:"non_repeaters"`
	DynamicRepetitions bool          `json:"dynamic_repetitions" yaml:"dynamic_repetitions" toml:"dynamic_repetitions" yaml:"dynamic_repetitions"`
}

type Configuration struct {
	SNMP   ConfigSNMP   `json:"community" toml:"snmp" yaml:"snmp"`
	Telnet ConfigTelnet `json:"telnet" toml:"telnet" yaml:"telnet"`
}

func Conf2proto(conf Configuration) proto.Configuration {
	return proto.Configuration{
		SNMP: &proto.ConfigSNMP{
			Community:          conf.SNMP.Community,
			Version:            uint32(conf.SNMP.Version),
			Port:               uint32(conf.SNMP.Port),
			Timeout:            uint64(conf.SNMP.Timeout),
			Retries:            int32(conf.SNMP.Retries),
			MaxIterations:      int32(conf.SNMP.MaxIterations),
			MaxRepetitions:     conf.SNMP.MaxRepetitions,
			NonRepeaters:       conf.SNMP.NonRepeaters,
			DynamicRepetitions: conf.SNMP.DynamicRepetitions,
		},
		Telnet: &proto.ConfigTelnet{
			User:          conf.Telnet.User,
			Password:      conf.Telnet.Password,
			Port:          conf.Telnet.Port,
			ScreenLength:  conf.Telnet.ScreenLength,
			RegexPrompt:   conf.Telnet.RegexPrompt,
			Errors:        conf.Telnet.Errors,
			TTL:           uint64(conf.Telnet.TTL),
			ReadDeadLine:  uint64(conf.Telnet.ReadDeadLine),
			WriteDeadLine: uint64(conf.Telnet.WriteDeadLine),
		},
	}
}

func Proto2conf(protoConf proto.Configuration) Configuration {
	return Configuration{
		SNMP: ConfigSNMP{
			Community:          protoConf.SNMP.Community,
			Version:            uint8(protoConf.SNMP.Version),
			Port:               uint16(protoConf.SNMP.Port),
			Timeout:            time.Duration(protoConf.SNMP.Timeout),
			Retries:            int32(protoConf.SNMP.Retries),
			MaxIterations:      int32(protoConf.SNMP.MaxIterations),
			MaxRepetitions:     protoConf.SNMP.MaxRepetitions,
			NonRepeaters:       protoConf.SNMP.NonRepeaters,
			DynamicRepetitions: protoConf.SNMP.DynamicRepetitions,
		},
		Telnet: ConfigTelnet{
			User:          protoConf.Telnet.User,
			Password:      protoConf.Telnet.Password,
			Port:          protoConf.Telnet.Port,
			ScreenLength:  protoConf.Telnet.ScreenLength,
			RegexPrompt:   protoConf.Telnet.RegexPrompt,
			Errors:        protoConf.Telnet.Errors,
			TTL:           time.Duration(protoConf.Telnet.TTL),
			ReadDeadLine:  time.Duration(protoConf.Telnet.ReadDeadLine),
			WriteDeadLine: time.Duration(protoConf.Telnet.WriteDeadLine),
		},
	}
}
