package config

// Transport is config for Telnet and SSH
type Transport struct {
	Label string `json:"type" hcl:",label"`

	// user and password to use, best is to let the DNC have this or the credentials will be sent over the wire
	User     string `json:"username" hcl:"username,optional"`
	Password string `json:"password" hcl:"password,optional"`

	// default poport to use (e.g. 22)
	Port int `json:"port" hcl:"port,optional"`

	// default command to set the screen lenght (e.g. "terminal length 0")
	ScreenLength string `json:"screen_length" hcl:"screen_length"`

	// default regex to find the prompt (e.g. ".*#")
	RegexPrompt string `json:"default_prompt" hcl:"default_prompt"`

	// default regex to find the error (e.g. "error|invalid|denied|fail")
	Errors string `json:"default_errors" hcl:"default_errors"`

	// default value for how long to read the prompt before giving up (e.g. "10s")
	ReadDeadLine Duration `json:"read_dead_line" hcl:"read_dead_line"`

	// default value for how long to wait to write to the device before giving up (e.g. "10s")
	WriteDeadLine Duration `json:"write_dead_line" hcl:"write_dead_line"`

	// where to find the SSH key
	SSHKeyPath string `json:"ssh_key_path" hcl:"ssh_key_path,optional"`
}

type Transports []*Transport

func (t Transports) GetByLabel(label string) *Transport {
	for _, c := range t {
		if c.Label == label {
			return c
		}
	}
	return nil
}
