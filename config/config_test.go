package config

import (
	"encoding/json"
	"fmt"
	"testing"
	
	"github.com/hashicorp/hcl/v2/gohcl"
)

func Test_Parse(t *testing.T) {
	
	bs, err := LoadFile("config.hcl")
	if err != nil {
		t.Error(err)
	}
	
	cfg := &Configuration{}
	err = ParseConfig(bs, "config.hcl", cfg)
	if err != nil {
		t.Error(err)
	}
	
	if cfg == nil {
		t.Error("no config?")
	}
	
	pv, err := cfg.GetProviderConfig("vx")
	if err != nil {
		t.Error(err)
	}
	
	fmt.Println(pv)
	
}
func Test_VRP(t *testing.T) {
	
	bs, err := LoadFile("config.hcl")
	if err != nil {
		t.Error(err)
	}
	
	cfg := &Configuration{}
	err = ParseConfig(bs, "config.hcl", cfg)
	if err != nil {
		t.Error(err)
	}
	
	if cfg == nil {
		t.Error("no config?")
	}
	
	pv, err := cfg.GetResourceConfig("vrp")
	if err != nil {
		t.Error(err)
	}
	
	n := &TestVRP{}
	diags := gohcl.DecodeBody(pv.Body, nil, n)
	if diags.HasErrors() {
		t.Error(err)
	}
	
	js, _ := json.MarshalIndent(n, "", "  ")
	fmt.Println(string(js))
	
}


func Test_Nats(t *testing.T) {
	
	bs, err := LoadFile("config.hcl")
	if err != nil {
		t.Error(err)
	}
	
	cfg := &Configuration{}
	err = ParseConfig(bs, "config.hcl", cfg)
	if err != nil {
		t.Error(err)
	}
	
	if cfg == nil {
		t.Error("no config?")
	}
	
	
	fmt.Println(cfg.NATS.GetServers())
	
	
	
}

