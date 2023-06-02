package config

type Configuration struct {
	// HttpAddr is the address to listen on for HTTP requests
	HttpAddr string `hcl:"http_addr"`

	// GrpcAddr is the address to listen on for GRPC requests
	GrpcAddr string `hcl:"grpc_addr"`

	Logger Logger `hcl:"logger,block"`

	BlacklistProvider ListStr `hcl:"blacklist_provider,optional"`

	MongoServer *MongoDb      `hcl:"mongodb,block"`
	MongoCaches []*MongoCache `hcl:"mongodb-cache,block"`
	Request     RequestConfig `hcl:"request,block"`
	Temporal    Temporal      `hcl:"temporal,block"`
}

type ListStr []string

func (l ListStr) Has(s string) bool {
	for _, v := range l {
		if v == s {
			return true
		}
	}
	return false
}

func (cfg Configuration) GetMongoByLabel(label string) *MongoCache {
	for _, c := range cfg.MongoCaches {
		if c.Label == label {
			return c
		}
	}
	return nil
}

type MongoDb struct {
	User     string   `hcl:"user,optional"`
	Password string   `hcl:"password,optional"`
	Timeout  Duration `hcl:"timeout" json:"timeout"` // Parse timeout as Duration (from string)
	Addr     string   `hcl:"addr"`
	Port     int      `hcl:"port"`

	// default one if not override in cache configs
	Database string `hcl:"database"`
}

type MongoCache struct {
	Label      string `hcl:",label"`
	Database   string `hcl:"database"`
	Collection string `hcl:"collection"`
}

type Logger struct {
	Level  string `hcl:"level,optional"`
	AsJson bool   `hcl:"as_json,optional"`
}

// Temporal config stanza
type Temporal struct {
	Address   string `hcl:"addr"`
	Namespace string `hcl:"namespace"`
	TaskQueue string `hcl:"task_queue"`
}
