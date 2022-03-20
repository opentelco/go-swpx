package tests

import (
	"fmt"
	"log"
	"testing"

	"git.liero.se/opentelco/go-swpx/config"
)

func Test_Parse(t *testing.T) {

	cfg := config.Configuration{}

	bs, err := config.LoadFile("config.hcl")
	if err != nil {
		t.Fatal(err)
	}
	config.ParseConfig(bs, "config.hcl", &cfg)

	db := cfg.Poller.GetCacheConfig(config.ResponseCacheType)
	if db == nil {
		log.Fatal("could not find cache type")
	}
	fmt.Println(db)

	db = cfg.Poller.GetCacheConfig(config.InterfaceCacheType)
	if db == nil {
		log.Fatal("could not find cache type")
	}
	fmt.Println(db)
}
