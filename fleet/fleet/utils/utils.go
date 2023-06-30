package utils

import (
	"regexp"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
)

// Parse the content of the Version from SNMP to a useful resource plugin name
func ParseVersionToResourcePlugin(version string) string {
	// regexp match HUAWEI or VRP in version string
	if match, err := regexp.MatchString(`(?i)(huawei|vrp)`, version); err == nil && match {
		return "vrp"
	}

	return "generic"
}

func GetDeviceScheduleByType(dev *devicepb.Device, t devicepb.Device_Schedule_Type) *devicepb.Device_Schedule {
	for _, s := range dev.Schedules {
		if s.Type == t {
			return s
		}
	}
	return nil
}
