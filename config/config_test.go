package config

import (
	"fmt"
	"testing"
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

	fmt.Println(cfg.Poller.Transports)

	type ArrayItem struct {
		Field  string `hcl:"field-a"`
		Field2 string `hcl:"field-b"`
	}
	type InnerConfig struct {
		Field string      `hcl:"field"`
		Items []ArrayItem `hcl:"item,block"`
	}
	v := InnerConfig{}
	pv, err := cfg.GetProviderConfig("vx", &v)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(pv, v)

	y := InnerConfig{}
	pv, err = cfg.GetProviderConfig("sait", &y)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pv, y.Items)

	fmt.Println("field value: ", v.Field)

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

	pv, err := cfg.GetResourceConfig("vrp", nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pv)

}

func Test_CTC(t *testing.T) {

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

	pv, err := cfg.GetResourceConfig("ctc", nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pv)

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

	fmt.Println(cfg.Poller.NATS.GetServers())

}
