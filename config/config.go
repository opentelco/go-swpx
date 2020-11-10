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

package config

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"time"
	
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)
const (
	DefaultMongoTimeout = time.Second * 5
)

var (
	ErrResourceConfigurationMissing = fmt.Errorf("resource plugin was not found or loaded")
	ErrProviderConfigurationMissing = fmt.Errorf("provider plugin was not found or loaded")
)

type Configuration struct {
	Logger      LoggerConf   `hcl:"logger,block"`
	MongoConfig *MongoConfig `hcl:"mongo,block"`

	NATS *NATSConfig `hcl:"nats,block"`
	// Drivers and Resources
	Resources []*Resource `hcl:"resource,block"`
	Providers []*Provider `hcl:"provider,block"`
}

func (c *Configuration) GetProviderConfig(plug string) (*Provider, error) {
	if c.Providers == nil || len(c.Providers) == 0 {
		return nil, ErrProviderConfigurationMissing
	}

	for _, p := range c.Providers {
		if p.Plugin == plug {
			return p, nil
		}
	}
	return nil, ErrProviderConfigurationMissing
}

func (c *Configuration) GetResourceConfig(plug string) (*Resource, error) {
	for _, r := range c.Resources {
		if r.Plugin == plug {
			return r, nil
		}
	}
	return nil, ErrResourceConfigurationMissing
}

type MongoConfig struct {
	Database   string `hcl:"database"`
	User       string `hcl:"user,optional"`
	Password   string `hcl:"password,optional`
	// Parse timeout as Duration (from string)
	Timeout string `hcl:"timeout" json:"timeout"`
	
	Servers    []*MongoServerEntry `hcl:"server,block"`
}

type MongoServerEntry struct {
	Addr    string `hcl:"addr"`
	Port    int    `hcl:"port"`
	Replica string `hcl:"replica,optional"`
}

func (c MongoConfig) GetServers() ([]string, error) {
	servers := make([]string,len(c.Servers))
	
	for _, s := range c.Servers {
		servers = append(servers, net.JoinHostPort(s.Addr, strconv.Itoa(s.Port)))
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("no mongo servers configured")
	}
	
	return servers,nil
}
func (c MongoConfig) GetTimeout() (time.Duration) {
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

func LoadFile(fname string) ([]byte, error) {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func ParseConfig(src []byte, filename string, cfg interface{}) error {
	var diags hcl.Diagnostics
	if src == nil || len(src) == 0 {
		return fmt.Errorf("no byte array supplied")
	}

	file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return fmt.Errorf("config parse: %w", diags)
	}

	diags = gohcl.DecodeBody(file.Body, nil, cfg)
	if diags.HasErrors() {
		return fmt.Errorf("config parse: %w", diags)
	}

	return nil
}
