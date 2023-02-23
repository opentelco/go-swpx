package config

import (
	"testing"
)

func Test_Parse(t *testing.T) {

	bs, err := loadFile("config.hcl")
	if err != nil {
		t.Error(err)
	}

	cfg := &Configuration{}
	err = parseConfig(bs, "config.hcl", cfg)
	if err != nil {
		t.Error(err)
	}

	if cfg == nil {
		t.Error("no config?")
	}

	type ArrayItem struct {
		Field  string `hcl:"field-a"`
		Field2 string `hcl:"field-b"`
	}
	type InnerConfig struct {
		Field string      `hcl:"field"`
		Items []ArrayItem `hcl:"item,block"`
	}

}
