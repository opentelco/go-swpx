package main

import (
	"strings"

	"git.liero.se/opentelco/go-dnc/models/protobuf/metric"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelement"
	"git.liero.se/opentelco/go-swpx/shared/oids"
)

func (v *VRPDriver) getIfEntryInformation(m *metric.Metric, elementInterface *networkelement.Interface) {
	switch {
	case strings.HasPrefix(m.Oid, oids.IfOutOctets):
		elementInterface.Stats.Output.Bytes = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInOctets):
		elementInterface.Stats.Input.Bytes = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInUcastPkts):
		elementInterface.Stats.Input.Unicast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInErrors):
		elementInterface.Stats.Input.Errors = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfOutErrors):
		elementInterface.Stats.Output.Errors = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfDescr):
		elementInterface.Description = m.GetStringValue()

	case strings.HasPrefix(m.Oid, oids.IfType):
		elementInterface.Type = networkelement.InterfaceType(m.GetIntValue())

	case strings.HasPrefix(m.Oid, oids.IfMtu):
		elementInterface.Mtu = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfLastChange):
		elementInterface.LastChanged = m.GetTimestampValue()

	case strings.HasPrefix(m.Oid, oids.IfPhysAddress):
		elementInterface.Hwaddress = m.GetStringValue()

	case strings.HasPrefix(m.Oid, oids.IfOperStatus):
		elementInterface.AdminStatus = networkelement.InterfaceStatus(m.GetIntValue())

	case strings.HasPrefix(m.Oid, oids.IfAdminStatus):
		elementInterface.OperationalStatus = networkelement.InterfaceStatus(m.GetIntValue())

	}
}
