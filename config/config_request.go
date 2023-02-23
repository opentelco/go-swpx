package config

import "time"

type RequestConfig struct {
	// DefaultRequestTimeout is the default timeout to use when no timeout is specified in the request
	DefaultRequestTimeout Duration `hcl:"default_request_timeout"`

	// DefaultTaskQueuePrefix is the default prefix to use when no prefix is specified in the request
	DefaultTaskQueuePrefix string `hcl:"default_task_queue_prefix"`

	// DefaultProvider is the default provider to use when no provider is specified in the request
	DefaultProvider string `hcl:"default_provider,optional"`

	// what cache TTL to use when no TTL is specified in the request
	DefaultCacheTTL Duration `hcl:"default_cache_ttl,optional"`
}

type Duration string

func (d Duration) AsDuration() time.Duration {
	if d == "" {
		return 0
	}
	duration, err := time.ParseDuration(string(d))
	if err != nil {
		return 0
	}
	return duration
}
