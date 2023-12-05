package tests

import (
	"fmt"
	"log"
	"testing"

	"go.opentelco.io/go-swpx/config"
)

func Test_Parse(t *testing.T) {

	cfg := config.Configuration{}
	err := config.LoadConfig("config.hcl", &cfg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cfg)

	db := cfg.GetMongoByLabel(config.LabelMongoInterfaceCache)
	if db == nil {
		log.Fatal("could not find cache type")
	}
	fmt.Println(db)

	db = cfg.GetMongoByLabel(config.LabelMongoResponseCache)
	if db == nil {
		log.Fatal("could not find cache type")
	}
	fmt.Println(db)

}
