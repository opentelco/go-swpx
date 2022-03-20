package config

import (
	"github.com/hashicorp/hcl/v2"
)

type Resource struct {
	Plugin      string `hcl:",label"`
	Description string `hcl:"description,optional"`
	Version     string `hcl:"version"`

	Config hcl.Body `hcl:",remain"`
}