package config

type RequestConfig struct {
	// DefaultRequestTimeout is the default timeout to use when no timeout is specified in the request
	DefaultRequestTimeout Duration `hcl:"default_request_timeout"`

	// DefaultProvider is the default provider to use when no provider is specified in the request
	DefaultProvider string `hcl:"default_provider,optional"`

	// what cache TTL to use when no TTL is specified in the request
	DefaultCacheTTL Duration `hcl:"default_cache_ttl,optional"`
}
