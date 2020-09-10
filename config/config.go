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
