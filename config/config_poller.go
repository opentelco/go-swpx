package config

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type Poller struct {
	Logger LoggerConf `hcl:"logger,block"`

	MongoCaches []*MongoConfig `hcl:"mongodb,block"`

	NATS       *NATSConfig `hcl:"nats,block"`
	Snmp       Snmp        `hcl:"snmp,block"`
	Transports []Transport `hcl:"transport,block"`
}

const (
	ResponseCacheType  string = "responseCache"
	InterfaceCacheType string = "interfaceCache"
)

func (p Poller) GetCacheConfig(cacheType string) *MongoConfig {
	for _, c := range p.MongoCaches {
		if c.CacheType == cacheType {
			return c
		}
	}
	return nil
}

type MongoConfig struct {
	CacheType  string              `hcl:",label"`
	Database   string              `hcl:"database"`
	Collection string              `hcl:"collection"`
	User       string              `hcl:"user,optional"`
	Password   string              `hcl:"password,optional"`
	Timeout    string              `hcl:"timeout" json:"timeout"` // Parse timeout as Duration (from string)
	Servers    []*MongoServerEntry `hcl:"server,block"`
}

type MongoServerEntry struct {
	Addr    string `hcl:"addr"`
	Port    int    `hcl:"port"`
	Replica string `hcl:"replica,optional"`
}

func (c MongoConfig) GetServers() ([]string, error) {
	servers := make([]string, len(c.Servers))

	for _, s := range c.Servers {
		servers = append(servers, net.JoinHostPort(s.Addr, strconv.Itoa(s.Port)))
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("no mongo servers configured")
	}

	return servers, nil
}
func (c MongoConfig) GetTimeout() time.Duration {
	d, err := time.ParseDuration(c.Timeout)
	if err != nil {
		return DefaultMongoTimeout
	}
	return d
}

type LoggerConf struct {
	Level  string `hcl:"level,optional"`
	AsJson bool   `hcl:"as_json,optional"`
}
