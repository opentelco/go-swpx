package config

type Snmp struct {
	StringVersion string `json:"stringVersion" hcl:",label"`

	Community string `json:"community" hcl:"community,optional"`

	Version uint8 `json:"version" hcl:"version,optional"`

	Port uint16 `json:"port" hcl:"port,optional"`

	Timeout Duration `json:"timeout" hcl:"timeout,optional"`
	Retries int      `json:"retries" hcl:"retries,optional"`

	MaxIterations  int `json:"max_iterations" hcl:"max_iterations,optional"`
	MaxRepetitions int `json:"max_repetitions" hcl:"max_repetitions,optional"`
	NonRepeaters   int `json:"non_repeaters" hcl:"non_repeaters,optional"`

	DynamicRepetitions bool `json:"dynamic_repetitions" hcl:"dynamic_repetitions,optional"`
}
