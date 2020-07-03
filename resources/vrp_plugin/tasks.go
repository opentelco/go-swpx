package main

import (
	"fmt"
	"git.liero.se/opentelco/go-dnc/models/protobuf/metric"
	shared2 "git.liero.se/opentelco/go-dnc/models/protobuf/shared"
	"git.liero.se/opentelco/go-dnc/models/protobuf/snmpc"
	"git.liero.se/opentelco/go-dnc/models/protobuf/transport"
	"git.liero.se/opentelco/go-swpx/shared"

	proto "git.liero.se/opentelco/go-swpx/proto/resource"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/segmentio/ksuid"
)

func createDiscoveryMsg(el *proto.NetworkElement, conf shared.Configuration) *transport.Message {
	task := &snmpc.Task{
		Config: &snmpc.Config{
			Community:          conf.SNMP.Community,
			DynamicRepititions: true,
			MaxIterations:      1,
			NonRepeaters:       0,
			Version:            snmpc.SnmpVersion(conf.SNMP.Version),
			Timeout:            ptypes.DurationProto(conf.SNMP.Timeout),
			Retries:            conf.SNMP.Retries,
		},
		Type: snmpc.Type_BULKGET,
		Oids: []*snmpc.Oid{
			{Oid: oids.IfIndex, Name: "ifIndex", Type: metric.MetricType_INT},
			{Oid: oids.IfDescr, Name: "ifDescr", Type: metric.MetricType_STRING},
			{Oid: oids.IfAlias, Name: "ifAlias", Type: metric.MetricType_STRING},
		},
	}

	// task.Parameters = params
	message := &transport.Message{
		Id:      ksuid.New().String(),
		Target:  el.Hostname,
		Type:    transport.Type_SNMP,
		Task:    &transport.Message_Snmpc{Snmpc: task},
		Status:  shared2.Status_NEW,
		Created: &timestamp.Timestamp{},
	}
	return message
}

// createMsg uses the pbuf transport for DNC..
func createSinglePortMsg(index int64, el *proto.NetworkElement, conf shared.Configuration) *transport.Message {
	task := &snmpc.Task{
		Config: &snmpc.Config{
			Community:          conf.SNMP.Community,
			DynamicRepititions: false,
			NonRepeaters:       12,
			MaxIterations:      1,
			Version:            snmpc.SnmpVersion(conf.SNMP.Version),
			Timeout:            ptypes.DurationProto(conf.SNMP.Timeout),
			Retries:            conf.SNMP.Retries,
		},
		Type: snmpc.Type_GET,
		Oids: []*snmpc.Oid{
			{Oid: fmt.Sprintf(oids.IfDescrF, index), Name: "ifDescr", Type: metric.MetricType_STRING},
			{Oid: fmt.Sprintf(oids.IfAliasF, index), Name: "ifAlias", Type: metric.MetricType_STRING},
			{Oid: fmt.Sprintf(oids.IfTypeF, index), Name: "ifType", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfMtuF, index), Name: "ifMtu", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfPhysAddressF, index), Name: "ifPhysAddress", Type: metric.MetricType_HWADDR},
			{Oid: fmt.Sprintf(oids.IfAdminStatusF, index), Name: "ifAdminStatus", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOperStatusF, index), Name: "ifOperStatus", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfLastChangeF, index), Name: "ifLastChange", Type: metric.MetricType_TIMETICKS},

			{Oid: fmt.Sprintf(oids.IfHighSpeedF, index), Name: "ifHighSpeed", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfConnectorPresentF, index), Name: "ifConnectorPresent", Type: metric.MetricType_BOOL},
		},
	}

	// task.Parameters = params
	message := &transport.Message{
		Id:      ksuid.New().String(),
		Target:  el.Hostname,
		Type:    transport.Type_SNMP,
		Task:    &transport.Message_Snmpc{Snmpc: task},
		Status:  shared2.Status_NEW,
		Created: &timestamp.Timestamp{},
	}
	return message
}

func createTaskSystemInfo(el *proto.NetworkElement, conf shared.Configuration) *transport.Message {
	task := &snmpc.Task{
		Config: &snmpc.Config{
			Community:          conf.SNMP.Community,
			DynamicRepititions: false,
			NonRepeaters:       12,
			MaxIterations:      1,
			Version:            snmpc.SnmpVersion(conf.SNMP.Version),
			Timeout:            ptypes.DurationProto(conf.SNMP.Timeout),
			Retries:            conf.SNMP.Retries,
		},
		Type: snmpc.Type_GET,
		Oids: []*snmpc.Oid{

			// OUT
			{Oid: oids.SysDescr, Name: "sysDescr", Type: metric.MetricType_STRING},
			{Oid: oids.SysObjectID, Name: "sysObjectID", Type: metric.MetricType_STRING},
			{Oid: oids.SysUpTime, Name: "sysUpTime", Type: metric.MetricType_STRING},
			{Oid: oids.SysContact, Name: "sysContact", Type: metric.MetricType_STRING},
			{Oid: oids.SysName, Name: "sysName", Type: metric.MetricType_STRING},
			{Oid: oids.SysLocation, Name: "sysLocation", Type: metric.MetricType_STRING},
			{Oid: oids.SysORLastChange, Name: "sysORLastChange", Type: metric.MetricType_INT},
		},
	}

	// task.Parameters = params
	message := &transport.Message{
		Id:      ksuid.New().String(),
		Target:  el.Hostname,
		Type:    transport.Type_SNMP,
		Task:    &transport.Message_Snmpc{Snmpc: task},
		Status:  shared2.Status_NEW,
		Created: &timestamp.Timestamp{},
	}
	return message
}

func createTaskGetPortStats(index int64, el *proto.NetworkElement, conf shared.Configuration) *transport.Message {
	task := &snmpc.Task{
		Config: &snmpc.Config{
			Community:          conf.SNMP.Community,
			DynamicRepititions: false,
			NonRepeaters:       12,
			MaxIterations:      1,
			Version:            snmpc.SnmpVersion(conf.SNMP.Version),
			Timeout:            ptypes.DurationProto(conf.SNMP.Timeout),
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
		Created: &timestamp.Timestamp{},
	}
	return message
}

// createMsg uses the pbuf transport for DNC..
func createMsg(conf shared.Configuration) *transport.Message {
	task := &snmpc.Task{
		Config: &snmpc.Config{
			Community:          conf.SNMP.Community,
			DynamicRepititions: true,
			MaxIterations:      200,
			Version:            snmpc.SnmpVersion(conf.SNMP.Version),
			Timeout:            ptypes.DurationProto(conf.SNMP.Timeout),
			Retries:            conf.SNMP.Retries,
		},
		Type: snmpc.Type_BULKGET,
		Oids: []*snmpc.Oid{

			{Oid: oids.IfIndex, Name: "ifIndex", Type: metric.MetricType_INT},
			{Oid: oids.IfDescr, Name: "ifDescr", Type: metric.MetricType_STRING},
			{Oid: oids.IfType, Name: "ifType", Type: metric.MetricType_INT},
			{Oid: oids.IfMtu, Name: "ifMtu", Type: metric.MetricType_INT},
			{Oid: oids.IfPhysAddress, Name: "ifPhysAddress", Type: metric.MetricType_HWADDR},
			{Oid: oids.IfAdminStatus, Name: "ifAdminStatus", Type: metric.MetricType_INT},
			{Oid: oids.IfOperStatus, Name: "ifOperStatus", Type: metric.MetricType_INT},
			{Oid: oids.IfLastChange, Name: "ifLastChange", Type: metric.MetricType_TIMETICKS},
			{Oid: oids.IfInErrors, Name: "ifInErrors", Type: metric.MetricType_INT},
			{Oid: oids.IfOutErrors, Name: "ifOutErrors", Type: metric.MetricType_INT},

			{Oid: oids.IfHighSpeed, Name: "ifHighSpeed", Type: metric.MetricType_INT},
			{Oid: oids.IfConnectorPresent, Name: "ifConnectorPresent", Type: metric.MetricType_BOOL},
			{Oid: oids.IfHCOutUcastPkts, Name: "ifHCOutUcastPkts", Type: metric.MetricType_INT},
		},
	}

	// task.Parameters = params
	message := &transport.Message{
		Id:      ksuid.New().String(),
		Target:  "10.9.0.10",
		Type:    transport.Type_SNMP,
		Task:    &transport.Message_Snmpc{Snmpc: task},
		Status:  shared2.Status_NEW,
		Created: &timestamp.Timestamp{},
	}
	return message
}
