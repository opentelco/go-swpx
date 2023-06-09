package utils

import "regexp"

// Parse the content of the Version from SNMP to a useful resource plugin name
func ParseVersionToResourcePlugin(version string) string {
	// regexp match HUAWEI or VRP in version string
	if match, err := regexp.MatchString(`(?i)(huawei|vrp)`, version); err == nil && match {
		return "vrp"
	}

	return "generic"
}
