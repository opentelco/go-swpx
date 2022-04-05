package main

import (
	"fmt"
	"strings"

	"git.liero.se/opentelco/go-dnc/models/pb/metric"
	shared2 "git.liero.se/opentelco/go-dnc/models/pb/shared"
	"git.liero.se/opentelco/go-dnc/models/pb/snmpc"
	"git.liero.se/opentelco/go-dnc/models/pb/transport"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelement"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createTaskGetPortStats(index int64, el *proto.NetworkElement, conf *shared.Configuration) *transport.Message {
	task := &snmpc.Task{
		Config: &snmpc.Config{
			Community:          conf.SNMP.Community,
			DynamicRepititions: false,
			NonRepeaters:       12,
			MaxIterations:      1,
			Version:            snmpc.SnmpVersion(conf.SNMP.Version),
			Timeout:            durationpb.New(conf.SNMP.Timeout),
			Retries:            conf.SNMP.Retries,
		},
		Type: snmpc.Type_GET,
		Oids: []*snmpc.Oid{

			// OUT
			{Oid: fmt.Sprintf(oids.IfOutErrorsF, index), Name: "ifOutErrors", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOutOctetsF, index), Name: "ifOutOctets", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOutUcastPktsF, index), Name: "ifOutUcastPkts", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOutBroadcastPktsF, index), Name: "ifOutBroadcastPkts", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOutMulticastPktsF, index), Name: "ifOutMulticastPkts", Type: metric.MetricType_INT},

			// In
			{Oid: fmt.Sprintf(oids.IfInErrorsF, index), Name: "ifInErrors", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfInOctetsF, index), Name: "ifInOctetsF", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfInUcastPktsF, index), Name: "ifInUcastPkts", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfInBroadcastPktsF, index), Name: "ifInBroadcastPkts", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfInMulticastPktsF, index), Name: "ifInMulticastPkts", Type: metric.MetricType_INT},

			// huawei spec.
			{Oid: fmt.Sprintf(oids.HuaIfEtherStatInCRCPktsF, index), Name: "HuaIfEtherStatInCRCPkts", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.HuaIfEtherStatInPausePktsF, index), Name: "HuaIfEtherStatInPausePkts", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.HuaIfEtherStatOutPausePktsF, index), Name: "HuaIfEtherStatOutPausePkts", Type: metric.MetricType_INT},
		},
	}

	// task.Parameters = params
	message := &transport.Message{
		Id:      ksuid.New().String(),
		Target:  el.Hostname,
		Type:    transport.Type_SNMP,
		Task:    &transport.Message_Snmpc{Snmpc: task},
		Status:  shared2.Status_NEW,
		Created: &timestamppb.Timestamp{},
	}
	return message
}

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

	case strings.HasPrefix(m.Oid, oids.IfAdminStatus):
		elementInterface.AdminStatus = networkelement.InterfaceStatus(m.GetIntValue())

	case strings.HasPrefix(m.Oid, oids.IfOperStatus):
		elementInterface.OperationalStatus = networkelement.InterfaceStatus(m.GetIntValue())

	}
}

func (v *VRPDriver) getHuaweiInterfaceStats(m *metric.Metric, elementInterface *networkelement.Interface) {
	switch {
	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatInCRCPkts):
		elementInterface.Stats.Input.CrcErrors = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatInPausePkts):
		elementInterface.Stats.Input.Pauses = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.HuaIfEthIfStatReset):
		elementInterface.Stats.Resets = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatOutPausePkts):
		elementInterface.Stats.Output.Pauses = m.GetIntValue()
	}
}
