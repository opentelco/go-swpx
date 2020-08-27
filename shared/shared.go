package shared

import (
	"fmt"
	"github.com/spf13/viper"
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
	Username            string        `mapstructure:"username" yaml:"username" toml:"username"`
	Password            string        `mapstructure:"password" yaml:"password" toml:"password"`
	Port                int32         `mapstructure:"port" yaml:"port" toml:"port"`
	ScreenLength        string        `mapstructure:"screen_length" yaml:"screen_length" toml:"screen_length"`
	ScreenLengthCommand string        `mapstructure:"screen_length_command" yaml:"screen_length_command" toml:"screen_length_command"`
	RegexPrompt         string        `mapstructure:"default_prompt" yaml:"default_prompt" toml:"default_prompt"`
	Errors              string        `mapstructure:"default_errors" yaml:"default_errors" toml:"default_errors"`
	TTL                 time.Duration `mapstructure:"ttl" yaml:"ttl" toml:"ttl"`
	ReadDeadLine        time.Duration `mapstructure:"read_dead_line" yaml:"read_dead_line" toml:"read_dead_line"`
	WriteDeadLine       time.Duration `mapstructure:"write_dead_line" yaml:"write_dead_line" toml:"write_dead_line"`
}

type ConfigSNMP struct {
	Community          string        `json:"community" toml:"community" yaml:"community"`
	Version            uint8         `json:"version" toml:"version" yaml:"version"`
	Port               uint16        `json:"port" toml:"port" yaml:"port"`
	Timeout            time.Duration `json:"timeout" toml:"timeout" yaml:"timeout"`
	Retries            int32         `json:"retries" toml:"retries" yaml:"retries"`
	DynamicRepetitions bool          `json:"dynamic_repetitions" yaml:"dynamic_repetitions" toml:"dynamic_repetitions" yaml:"dynamic_repetitions"`
}

type ConfigMongo struct {
	Server         string `json:"server" toml:"server" yaml:"server"`
	Database       string `json:"database" toml:"database" yaml:"database"`
	Collection     string `json:"collection" toml:"collection" yaml:"collection"`
	TimeoutSeconds int    `json:"timeout_seconds" toml:"timeout_seconds" yaml:"timeout_seconds"`
}

type ConfigNATS struct {
	EventServers []string `json:"event_servers" toml:"event_servers" yaml:"event_servers" mapstructure:"event_servers"`
}

type Configuration struct {
	SNMP   ConfigSNMP   `json:"snmp" toml:"snmp" yaml:"snmp"`
	Telnet ConfigTelnet `json:"telnet" toml:"telnet" yaml:"telnet"`
	Mongo  ConfigMongo  `json:"mongo" toml:"mongo" yaml:"mongo"`
	NATS   ConfigNATS   `json:"nats" toml:"nats" yaml:"nats"`
}

func GetConfig() Configuration {
	conf := Configuration{}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}

	return conf
}

func Conf2proto(conf Configuration) proto.Configuration {
	return proto.Configuration{
		SNMP: &proto.ConfigSNMP{
			Community:          conf.SNMP.Community,
			Version:            uint32(conf.SNMP.Version),
			Port:               uint32(conf.SNMP.Port),
			Timeout:            uint64(conf.SNMP.Timeout),
			Retries:            int32(conf.SNMP.Retries),
			DynamicRepetitions: conf.SNMP.DynamicRepetitions,
		},
		Telnet: &proto.ConfigTelnet{
			User:                conf.Telnet.Username,
			Password:            conf.Telnet.Password,
			Port:                conf.Telnet.Port,
			ScreenLength:        conf.Telnet.ScreenLength,
			ScreenLengthCommand: conf.Telnet.ScreenLengthCommand,
			RegexPrompt:         conf.Telnet.RegexPrompt,
			Errors:              conf.Telnet.Errors,
			TTL:                 uint64(conf.Telnet.TTL),
			ReadDeadLine:        uint64(conf.Telnet.ReadDeadLine),
			WriteDeadLine:       uint64(conf.Telnet.WriteDeadLine),
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
			DynamicRepetitions: protoConf.SNMP.DynamicRepetitions,
		},
		Telnet: ConfigTelnet{
			Username:            protoConf.Telnet.User,
			Password:            protoConf.Telnet.Password,
			Port:                protoConf.Telnet.Port,
			ScreenLength:        protoConf.Telnet.ScreenLength,
			ScreenLengthCommand: protoConf.Telnet.ScreenLengthCommand,
			RegexPrompt:         protoConf.Telnet.RegexPrompt,
			Errors:              protoConf.Telnet.Errors,
			TTL:                 time.Duration(protoConf.Telnet.TTL),
			ReadDeadLine:        time.Duration(protoConf.Telnet.ReadDeadLine),
			WriteDeadLine:       time.Duration(protoConf.Telnet.WriteDeadLine),
		},
	}
}
