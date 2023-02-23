package config

type ProviderExample struct {
	Description string `hcl:"description,optional"`
	Version     string `hcl:"version"`
}

type ProviderVX struct {
	Description string `hcl:"description,optional"`
	Version     string `hcl:"version"`
}
