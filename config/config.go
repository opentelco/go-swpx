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

type Config struct {
	Key       string
	Region    string
	Providers []ProviderConfig
	Drivers   []DriverConfig
}

type ProviderConfig struct {
	Weight int
}

type DriverConfig struct {
	OIDs          []string
	RegexMatchOID string
}

type MongoConfig struct {
	Addr string `hcl:"addr"`
	Port int `hcl:"port"`
	Database string `hcl:"database"`
	User string `hcl:"user,optional"`
	Password string `hcl:"password,optional`
}

type LoggerConf struct {
	Level string `hcl:"level,optional"`
	AsJson bool `hcl:"as_json,optional"`
}


type Configuration struct {
	Logger LoggerConf `hcl:"logger,block"`
	MongoConfig *MongoConfig `hcl:"mongo,block"`
	
	// Drivers and Resources
	Resources []*Resource `hcl:"resource,block"`
	Providers []*Provider `hcl:"provider,block"`
}


type Resource struct {
	Plugin string `hcl:",label"`
	Name string `hcl:"name"`
	Description string `hcl:"description,optional"`
	Version string `hcl:"version"`
	
}


type Provider struct {
	Plugin string `hcl:",label"`
	Name string `hcl:"name"`
	Description string `hcl:"description,optional"`
	Version string `hcl:"version"`
	
	
	
}