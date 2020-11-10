package config

type ConfigConnection struct {
	Type                string `json:"type" hcl:",label"`
	User                string `json:"username" hcl:"username"`
	Password            string `json:"password" hcl:"password"`
	Port                int    `json:"port" hcl:"port"`
	ScreenLength        string `json:"screen_length" hcl:"screen_length,optional"`
	RegexPrompt         string `json:"default_prompt" hcl:"default_prompt"`
	Errors              string `json:"default_errors" hcl:"default_errors"`
	CacheTTL            string `json:"cache_ttl" hcl:"cache_ttl,optional"`
	ReadDeadLine        string `json:"read_dead_line" hcl:"read_dead_line,optional"`
	WriteDeadLine       string `json:"write_dead_line" hcl:"write_dead_line,optional"`
	SSHKeyPath          string `json:"ssh_key_path" hcl:"ssh_key_path,optional"`
}

