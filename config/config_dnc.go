package config

import (
	"fmt"
	"strings"
)

type TestVRP struct {
	NATS *NATSConfig `hcl:"nats,block"`
}

type NATSConfig struct {
	Username string `hcl:"username,optional"`
	Password string `hcl:"password,optional"`
	NatsServers []*NatsServerEntry `hcl:"server,block"`
}

func (n NATSConfig) GetServers() (string,error) {
	servers := make([]string,0)
	
	for _, s := range n.NatsServers {
		servers = append(servers, fmt.Sprintf("nats://%s:%d", s.Addr, s.Port))
	}
	if len(servers) == 0 {
		return "", fmt.Errorf("no mongo servers configured")
	}
	return strings.Join(servers,","),nil
}

type NatsServerEntry struct {
	Addr string `hcl:"addr"`
	Port int    `hcl:"port"`
}