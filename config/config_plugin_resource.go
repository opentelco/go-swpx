package config

type ResourceExample struct {
	Description string `hcl:"description,optional"`
	Version     string `hcl:"version"`
}

type DNC struct {
	Addr string `hcl:"addr"`
}

type ResourceGeneric struct {
	DNC        DNC        `hcl:"dnc,block"`
	Snmp       Snmp       `hcl:"snmp,block"`
	Transports Transports `hcl:"transport,block"`
	Logger     Logger     `hcl:"logger,block"`
}

// ResourceVRP is the configuration for the VRP resource plugin and it comes bundled with the switchpoller
// it is the config that fills upp the "hcl.Body" block in the Resource struct
type ResourceVRP struct {
	DNC        DNC        `hcl:"dnc,block"`
	Snmp       Snmp       `hcl:"snmp,block"`
	Transports Transports `hcl:"transport,block"`
	Logger     Logger     `hcl:"logger,block"`
}

// ResourceCTC is the configuration for the CTC resource plugin and it comes bundled with the switchpoller
// it is the config that fills upp the "hcl.Body" block in the Resource struct
type ResourceCTC struct {
	DNC        DNC        `hcl:"dnc,block"`
	Snmp       Snmp       `hcl:"snmp,block"`
	Transports Transports `hcl:"transport,block"`
	Logger     Logger     `hcl:"logger,block"`
}
