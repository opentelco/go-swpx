package core

import "time"

// Generated structs for serializing BSON when used in response cache

type BSONResponse struct {
	Response struct {
		Status Status `json:"status,omitempty" bson:"status,omitempty"`
		Data   Data   `json:"data,omitempty" bson:"data,omitempty"`
	}
	Hostname    string      `bson:"hostname"`
	Port        string      `bson:"port"`
	RequestType RequestType `bson:"request_type"`
	Timestamp   time.Time   `bson:"timestamp"`
}

type DhcpTable struct {
	IpAddress       string `json:"ip_address,omitempty" bson:"ipaddress,omitempty"`
	HardwareAddress string `json:"hardware_address,omitempty" bson:"hardwareaddress,omitempty"`
	Vendor          string `json:"vendor,omitempty" bson:"vendor,omitempty"`
	Vlan            int    `json:"vlan,omitempty" bson:"vlan,omitempty"`
	Timestamp       string `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

type LastChanged struct {
	Seconds int `json:"seconds,omitempty" bson:"seconds,omitempty"`
	Nanos   int `json:"nanos,omitempty" bson:"nanos,omitempty"`
}

type MacAddressTable struct {
	HardwareAddress string `json:"hardware_address,omitempty" bson:"hardwareaddress,omitempty"`
	Vlan            int    `json:"vlan,omitempty" bson:"vlan,omitempty"`
	Vendor          string `json:"vendor,omitempty" bson:"vendor,omitempty"`
	Timestamp       string `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

type Input struct {
	Bytes     int `json:"bytes,omitempty" bson:"bytes,omitempty"`
	Unicast   int `json:"unicast,omitempty" bson:"unicast,omitempty"`
	Broadcast int `json:"broadcast,omitempty" bson:"broadcast,omitempty"`
	Multicast int `json:"multicast,omitempty" bson:"multicast,omitempty"`
	Crcerrors int `json:"crcerrors,omitempty" bson:"crcerrors,omitempty"`
	Errors    int `json:"errors,omitempty" bson:"errors,omitempty"`
	Packets   int `json:"packets,omitempty" bson:"packets,omitempty"`
	Pauses    int `json:"pauses,omitempty" bson:"pauses,omitempty"`
}

type Output struct {
	Bytes     int `json:"bytes,omitempty" bson:"bytes,omitempty"`
	Unicast   int `json:"unicast,omitempty" bson:"unicast,omitempty"`
	Broadcast int `json:"broadcast,omitempty" bson:"broadcast,omitempty"`
	Multicast int `json:"multicast,omitempty" bson:"multicast,omitempty"`
	Errors    int `json:"errors,omitempty" bson:"errors,omitempty"`
	Packets   int `json:"packets,omitempty" bson:"packets,omitempty"`
	Pauses    int `json:"pauses,omitempty" bson:"pauses,omitempty"`
}

type Stats struct {
	Input  Input  `json:"input,omitempty" bson:"input,omitempty"`
	Output Output `json:"output,omitempty" bson:"output,omitempty"`
	Resets int    `json:"resets,omitempty" bson:"resets,omitempty"`
}

type Interface struct {
	Index             int               `json:"index,omitempty" bson:"index,omitempty"`
	Alias             string            `json:"alias,omitempty" bson:"alias,omitempty"`
	Description       string            `json:"description,omitempty" bson:"description,omitempty"`
	Type              int               `json:"type,omitempty" bson:"type,omitempty"`
	AdminStatus       int               `json:"admin_status,omitempty" bson:"adminstatus,omitempty"`
	OperationalStatus int               `json:"operational_status,omitempty" bson:"operationalstatus,omitempty"`
	LastChanged       LastChanged       `json:"last_changed,omitempty" bson:"lastchanged,omitempty"`
	ConnectorPresent  bool              `json:"connector_present,omitempty" bson:"connectorpresent,omitempty"`
	Speed             int               `json:"speed,omitempty" bson:"speed,omitempty"`
	Mtu               int               `json:"mtu,omitempty" bson:"mtu,omitempty"`
	Stats             Stats             `json:"stats,omitempty" bson:"stats,omitempty"`
	MacAddressTable   []MacAddressTable `json:"mac_address_table,omitempty" bson:"macaddresstable,omitempty"`
	DhcpTable         []DhcpTable       `json:"dhcp_table,omitempty" bson:"dhcptable,omitempty"`
	Config            string            `json:"config,omitempty" bson:"config,omitempty"`
	AggregatedId      string            `json:"aggregated_id,omitempty" bson:"aggregated_id,omitempty"`
	Duplex            string            `json:"duplex,omitempty" bson:"duplex,omitempty"`
	HwAddress         string            `json:"hw_address,omitempty" bson:"hwaddress,omitempty"`
	InterfaceStatus   int               `json:"interface_status,omitempty" bson:"interface_status,omitempty"`
	Neighbor          interface{}       `json:"neighbor,omitempty" bson:"neighbor,omitempty"`
	Transceiver       interface{}       `json:"transceiver,omitempty" bson:"transceiver,omitempty"`
}

type NetworkElement struct {
	Hostname             string      `json:"hostname,omitempty" bson:"hostname,omitempty"`
	Version              string      `json:"version,omitempty" bson:"version,omitempty"`
	Contact              string      `json:"contact,omitempty" bson:"contact,omitempty"`
	Sysname              string      `json:"sysname,omitempty" bson:"sysname,omitempty"`
	Location             string      `json:"location,omitempty" bson:"location,omitempty"`
	Interfaces           []Interface `json:"interfaces,omitempty" bson:"interfaces,omitempty"`
	AggregatedInterfaces interface{} `json:"aggregated_interfaces,omitempty" bson:"aggregated_interfaces,omitempty"`
	BridgeMacAddress     string      `json:"bridge_mac_address,omitempty" bson:"bridge_mac_address,omitempty"`
	Driver               string      `json:"driver,omitempty" bson:"driver,omitempty"`
	InterfaceIndex       int         `json:"interface_index,omitempty" bson:"interfaceindex,omitempty"`
	Modules              interface{} `json:"modules,omitempty" bson:"modules,omitempty"`
	SnmpObjectID         string      `json:"snmp_object_id,omitempty" bson:"snmp_object_id,omitempty"`
	SoftwareVersion      string      `json:"software_version,omitempty" bson:"software_version,omitempty"`
	Uptime               string      `json:"uptime,omitempty" bson:"uptime,omitempty"`
	Virtual              bool        `json:"virtual,omitempty" bson:"virtual,omitempty"`
}

type TransceiverStats struct {
	Current   float64 `json:"current,omitempty" bson:"current,omitempty"`
	Rx        float64 `json:"rx,omitempty" bson:"rx,omitempty"`
	Tx        float64 `json:"tx,omitempty" bson:"tx,omitempty"`
	Temp      int     `json:"temp,omitempty" bson:"temp,omitempty"`
	Voltage   float64 `json:"voltage,omitempty" bson:"voltage,omitempty"`
	Timestamp string  `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

type Transceiver struct {
	SerialNumber      string             `json:"serial_number,omitempty" bson:"serialnumber,omitempty"`
	PartNumber        string             `json:"part_number,omitempty" bson:"partnumber,omitempty"`
	TransceiverStats  []TransceiverStats `json:"stats,omitempty" bson:"stats,omitempty"`
	ConnectorType     string             `json:"connector_type,omitempty" bson:"connectortype,omitempty"`
	Ddm               bool               `json:"ddm,omitempty" bson:"ddm,omitempty"`
	ManufacturingDate string             `json:"manufacturing_date,omitempty" bson:"manufacturingdate,omitempty"`
	TransferDistance  string             `json:"transfer_distance,omitempty" bson:"transferdistance,omitempty"`
	Type              string             `json:"type,omitempty" bson:"type,omitempty"`
	Vendor            string             `json:"vendor,omitempty" bson:"vendor,omitempty"`
	Wavelength        string             `json:"wavelength,omitempty" bson:"wavelength,omitempty"`
}

type Data struct {
	NetworkElement  NetworkElement `json:"network_element,omitempty" bson:"network_element,omitempty"`
	PhysicalPort    string         `json:"physical_port,omitempty" bson:"physical_port,omitempty"`
	Transceiver     Transceiver    `json:"transceiver,omitempty" bson:"transceiver,omitempty"`
	Error           interface{}    `json:"error,omitempty" bson:"error,omitempty"`
	RequestObjectID string         `json:"request_object_id,omitempty" bson:"request_object_id,omitempty"`
}

type Status struct {
	Code    int    `json:"code,omitempty" bson:"code,omitempty"`
	Error   bool   `json:"error,omitempty" bson:"error,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
	Type    string `json:"type,omitempty" bson:"type,omitempty"`
}
