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

package shared

import (
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"

	proto "git.liero.se/opentelco/go-swpx/proto/go/provider"
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
var DefaultCorePath = fmt.Sprintf(".%s", CoreName)

func init() {
	VERSION, _ = version.NewVersion(VersionString)
}

type DNCRequest struct {
	OriginalID string
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

type ConfigTelnet struct {
	Username            string        `json:"username" mapstructure:"username" yaml:"username" toml:"username"`
	Password            string        `json:"password" mapstructure:"password" yaml:"password" toml:"password"`
	Port                int32         `json:"port" mapstructure:"port" yaml:"port" toml:"port"`
	ScreenLength        string        `json:"screen_length" mapstructure:"screen_length" yaml:"screen_length" toml:"screen_length"`
	ScreenLengthCommand string        `json:"screen_length_command" mapstructure:"screen_length_command" yaml:"screen_length_command" toml:"screen_length_command"`
	RegexPrompt         string        `json:"default_prompt" mapstructure:"default_prompt" yaml:"default_prompt" toml:"default_prompt"`
	Errors              string        `json:"default_errors" mapstructure:"default_errors" yaml:"default_errors" toml:"default_errors"`
	TTL                 time.Duration `json:"ttl" mapstructure:"ttl" yaml:"ttl" toml:"ttl"`
	ReadDeadLine        time.Duration `json:"read_dead_line" mapstructure:"read_dead_line" yaml:"read_dead_line" toml:"read_dead_line"`
	WriteDeadLine       time.Duration `json:"write_dead_line" mapstructure:"write_dead_line" yaml:"write_dead_line" toml:"write_dead_line"`
}

type ConfigSSH struct {
	Username            string        `json:"username" mapstructure:"username" yaml:"username" toml:"username"`
	Password            string        `json:"password" mapstructure:"password" yaml:"password" toml:"password"`
	Port                int32         `json:"port" mapstructure:"port" yaml:"port" toml:"port"`
	ScreenLength        string        `json:"screen_length" mapstructure:"screen_length" yaml:"screen_length" toml:"screen_length"`
	ScreenLengthCommand string        `json:"screen_length_command" mapstructure:"screen_length_command" yaml:"screen_length_command" toml:"screen_length_command"`
	RegexPrompt         string        `json:"default_prompt" mapstructure:"default_prompt" yaml:"default_prompt" toml:"default_prompt"`
	Errors              string        `json:"default_errors" mapstructure:"default_errors" yaml:"default_errors" toml:"default_errors"`
	TTL                 time.Duration `json:"ttl" mapstructure:"ttl" yaml:"ttl" toml:"ttl"`
	ReadDeadLine        time.Duration `json:"read_dead_line" mapstructure:"read_dead_line" yaml:"read_dead_line" toml:"read_dead_line"`
	WriteDeadLine       time.Duration `json:"write_dead_line" mapstructure:"write_dead_line" yaml:"write_dead_line" toml:"write_dead_line"`
	SSHKeyPath          string        `json:"ssh_key_path" mapstructure:"ssh_key_path" yaml:"ssh_key_path" toml:"ssh_key_path"`
}

type ConfigSNMP struct {
	Community          string        `json:"community" toml:"community" yaml:"community"`
	Version            uint8         `json:"version" toml:"version" yaml:"version"`
	Port               uint16        `json:"port" toml:"port" yaml:"port"`
	Timeout            time.Duration `json:"timeout" toml:"timeout" yaml:"timeout"`
	Retries            int32         `json:"retries" toml:"retries" yaml:"retries"`
	DynamicRepetitions bool          `json:"dynamic_repetitions" toml:"dynamic_repetitions" yaml:"dynamic_repetitions"`
}

type ConfigMongo struct {
	Server         string `json:"server" toml:"server" yaml:"server"`
	Database       string `json:"database" toml:"database" yaml:"database"`
	Username       string `json:"username" toml:"username" yaml:"username"`
	Password       string `json:"password" toml:"password" yaml:"password"`
	Collection     string `json:"collection" toml:"collection" yaml:"collection"`
	TimeoutSeconds int    `json:"timeout_seconds" toml:"timeout_seconds" yaml:"timeout_seconds"`
}

type ConfigNATS struct {
	EventServers []string `json:"event_servers" toml:"event_servers" yaml:"event_servers" mapstructure:"event_servers"`
}

type Configuration struct {
	SNMP            ConfigSNMP   `json:"snmp" toml:"snmp" yaml:"snmp" mapstructure:"snmp"`
	InterfaceCache  ConfigMongo  `json:"interface_cache" toml:"interface_cache" yaml:"interface_cache" mapstructure:"interface_cache"`
	ResponseCache   ConfigMongo  `json:"response_cache" toml:"response_cache" yaml:"response_cache" mapstructure:"response_cache"`
	NATS            ConfigNATS   `json:"nats" toml:"nats" yaml:"nats" mapstructure:"nats"`
	Telnet          ConfigTelnet `json:"telnet" mapstructure:"telnet" yaml:"telnet" toml:"telnet"`
	Ssh             ConfigSSH    `json:"ssh" mapstructure:"ssh" yaml:"ssh" toml:"ssh"`
	DefaultProvider string       `json:"default_provider" mapstructure:"default_provider" yaml:"default_provider" toml:"default_provider"`
}

func GetConfig() *Configuration {
	conf := &Configuration{}
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/swpx/")
	viper.AddConfigPath("/etc/swpx/")

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

func Conf2proto(conf *Configuration) *proto.Configuration {
	return &proto.Configuration{
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
		Ssh: &proto.ConfigSSH{
			User:                conf.Ssh.Username,
			Password:            conf.Ssh.Password,
			Port:                conf.Ssh.Port,
			ScreenLength:        conf.Ssh.ScreenLength,
			ScreenLengthCommand: conf.Ssh.ScreenLengthCommand,
			RegexPrompt:         conf.Ssh.RegexPrompt,
			Errors:              conf.Ssh.Errors,
			TTL:                 uint64(conf.Ssh.TTL),
			ReadDeadLine:        uint64(conf.Ssh.ReadDeadLine),
			WriteDeadLine:       uint64(conf.Ssh.WriteDeadLine),
			SSHKeyPath:          conf.Ssh.SSHKeyPath,
		},
	}
}

func Proto2conf(protoConf *proto.Configuration) *Configuration {
	return &Configuration{
		SNMP: ConfigSNMP{
			Community:          protoConf.SNMP.Community,
			Version:            uint8(protoConf.SNMP.Version),
			Port:               uint16(protoConf.SNMP.Port),
			Timeout:            time.Duration(protoConf.SNMP.Timeout),
			Retries:            protoConf.SNMP.Retries,
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
		Ssh: ConfigSSH{
			Username:            protoConf.Ssh.User,
			Password:            protoConf.Ssh.Password,
			Port:                protoConf.Ssh.Port,
			ScreenLength:        protoConf.Ssh.ScreenLength,
			ScreenLengthCommand: protoConf.Ssh.ScreenLengthCommand,
			RegexPrompt:         protoConf.Ssh.RegexPrompt,
			Errors:              protoConf.Ssh.Errors,
			TTL:                 time.Duration(protoConf.Ssh.TTL),
			ReadDeadLine:        time.Duration(protoConf.Ssh.ReadDeadLine),
			WriteDeadLine:       time.Duration(protoConf.Ssh.WriteDeadLine),
			SSHKeyPath:          protoConf.Ssh.SSHKeyPath,
		},
	}
}
