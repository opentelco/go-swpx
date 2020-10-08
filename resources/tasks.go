/*
 * Copyright (c) 2020. Liero AB
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software
 * is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
 * IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
 * CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package resources

import (
	"fmt"
	"git.liero.se/opentelco/go-dnc/models/protobuf/metric"
	shared2 "git.liero.se/opentelco/go-dnc/models/protobuf/shared"
	"git.liero.se/opentelco/go-dnc/models/protobuf/snmpc"
	"git.liero.se/opentelco/go-dnc/models/protobuf/ssh"
	"git.liero.se/opentelco/go-dnc/models/protobuf/telnet"
	"git.liero.se/opentelco/go-dnc/models/protobuf/transport"
	proto "git.liero.se/opentelco/go-swpx/proto/resource"
	"git.liero.se/opentelco/go-swpx/shared"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/segmentio/ksuid"
)

func CreateDiscoveryMsg(el *proto.NetworkElement, conf shared.Configuration) *transport.Message {
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
func CreateSinglePortMsg(index int64, el *proto.NetworkElement, conf shared.Configuration) *transport.Message {
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

func CreateTaskSystemInfo(el *proto.NetworkElement, conf shared.Configuration) *transport.Message {
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

func CreateTaskGetPortStats(index int64, el *proto.NetworkElement, conf shared.Configuration) *transport.Message {
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

// CreateMsg uses the pbuf transport for DNC..
func CreateMsg(conf shared.Configuration) *transport.Message {
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

func CreatePortInformationMsg(el *proto.NetworkElement, conf shared.Configuration) *transport.Message {
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
		Type: snmpc.Type_WALK,
		Oids: []*snmpc.Oid{
			{Oid: oids.IfEntPhysicalName, Name: "ifPhysAddress", Type: metric.MetricType_STRING},
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

func CreateVRPTransceiverMsg(el *proto.NetworkElement, conf shared.Configuration) *transport.Message {
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
		Type: snmpc.Type_GET,
		Oids: []*snmpc.Oid{
			{Oid: fmt.Sprintf(oids.HuaIfVRPOpticalVendorSNF, el.PhysicalIndex), Name: "hwEntityOpticalVendorSn", Type: metric.MetricType_STRING},
			{Oid: fmt.Sprintf(oids.HuaIfVRPOpticalTemperatureF, el.PhysicalIndex), Name: "hwEntityOpticalTemperature", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.HuaIfVRPOpticalVoltageF, el.PhysicalIndex), Name: "hwEntityOpticalVoltage", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.HuaIfVRPOpticalBiasF, el.PhysicalIndex), Name: "hwEntityOpticalBiasCurrent", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.HuaIfVRPOpticalRxPowerF, el.PhysicalIndex), Name: "hwEntityOpticalRxPower", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.HuaIfVRPOpticalTxPowerF, el.PhysicalIndex), Name: "hwEntityOpticalTxPower", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.HuaIfVRPVendorPNF, el.PhysicalIndex), Name: "hwEntityOpticalVenderPn", Type: metric.MetricType_STRING},
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

func createTelnetInterfaceTask(el *proto.NetworkElement, conf shared.Configuration) *transport.Message {
	task := &telnet.Task{
		Type: telnet.Type_GET,
		Payload: []*telnet.Payload{
			{
				Command: fmt.Sprintf("display mac-address %s", el.Interface),
			},
			{
				Command: fmt.Sprintf("display dhcp snooping user-bind interface %s", el.Interface),
			},
			{
				Command: fmt.Sprintf("display current-configuration interface %s", el.Interface),
			},
			{
				Command: fmt.Sprintf("display current-configuration interface %s | include traffic-policy|shaping", el.Interface),
			},
			{
				Command: fmt.Sprintf("display traffic policy statistics interface %s inbound verbose classifier-base", el.Interface),
			},
			{
				Command: fmt.Sprintf("display qos queue statistics interface %s", el.Interface),
			},
		},
		Config: &telnet.Config{
			User:                conf.Telnet.Username,
			Password:            conf.Telnet.Password,
			Port:                conf.Telnet.Port,
			ScreenLength:        conf.Telnet.ScreenLength,
			ScreenLengthCommand: conf.Telnet.ScreenLengthCommand,
			RegexPrompt:         conf.Telnet.RegexPrompt,
			Ttl:                 &duration.Duration{Seconds: int64(conf.Telnet.TTL.Seconds())},
			ReadDeadLine:        &duration.Duration{Seconds: int64(conf.Telnet.ReadDeadLine.Seconds())},
			WriteDeadLine:       &duration.Duration{Seconds: int64(conf.Telnet.WriteDeadLine.Seconds())},
		},
		Host: el.Hostname,
	}

	message := &transport.Message{
		Id:      ksuid.New().String(),
		Target:  el.Hostname,
		Type:    transport.Type_TELNET,
		Task:    &transport.Message_Telnet{Telnet: task},
		Status:  shared2.Status_NEW,
		Created: &timestamp.Timestamp{},
	}
	return message

}

func CreateSSHInterfaceTask(el *proto.NetworkElement, conf shared.Configuration) *transport.Message {
	task := &ssh.Task{
		Type: ssh.Type_GET,
		Payload: []*ssh.Payload{
			{
				Command: fmt.Sprintf("display mac-address %s", el.Interface),
			},
			{
				Command: fmt.Sprintf("display dhcp snooping user-bind interface %s", el.Interface),
			},
			{
				Command: fmt.Sprintf("display current-configuration interface %s", el.Interface),
			},
			{
				Command: fmt.Sprintf("display current-configuration interface %s | include traffic-policy|shaping", el.Interface),
			},
			{
				Command: fmt.Sprintf("display traffic policy statistics interface %s inbound verbose classifier-base", el.Interface),
			},
			{
				Command: fmt.Sprintf("display qos queue statistics interface %s", el.Interface),
			},
		},
		Config: &ssh.Config{
			User:                conf.Ssh.Username,
			Password:            conf.Ssh.Password,
			Port:                conf.Ssh.Port,
			ScreenLength:        conf.Ssh.ScreenLength,
			ScreenLengthCommand: conf.Ssh.ScreenLengthCommand,
			RegexPrompt:         conf.Ssh.RegexPrompt,
			Ttl:                 &duration.Duration{Seconds: int64(conf.Ssh.TTL.Seconds())},
			ReadDeadLine:        &duration.Duration{Seconds: int64(conf.Ssh.ReadDeadLine.Seconds())},
			WriteDeadLine:       &duration.Duration{Seconds: int64(conf.Ssh.WriteDeadLine.Seconds())},
			SshKeyPath:          conf.Ssh.SSHKeyPath,
		},
		Host: el.Hostname,
	}

	message := &transport.Message{
		Id:      ksuid.New().String(),
		Target:  el.Hostname,
		Type:    transport.Type_SSH,
		Task:    &transport.Message_Ssh{Ssh: task},
		Status:  shared2.Status_NEW,
		Created: &timestamp.Timestamp{},
	}
	return message

}

func CreateAllPortsMsg(el *proto.NetworkElement, conf shared.Configuration) *transport.Message {
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
			{Oid: oids.IfDescr, Name: "ifDescr", Type: metric.MetricType_STRING},
			{Oid: oids.IfAlias, Name: "ifAlias", Type: metric.MetricType_STRING},
			{Oid: oids.IfType, Name: "ifType", Type: metric.MetricType_INT},
			{Oid: oids.IfMtu, Name: "ifMtu", Type: metric.MetricType_INT},
			{Oid: oids.IfPhysAddress, Name: "ifPhysAddress", Type: metric.MetricType_HWADDR},
			{Oid: oids.IfAdminStatus, Name: "ifAdminStatus", Type: metric.MetricType_INT},
			{Oid: oids.IfOperStatus, Name: "ifOperStatus", Type: metric.MetricType_INT},
			{Oid: oids.IfLastChange, Name: "ifLastChange", Type: metric.MetricType_TIMETICKS},

			{Oid: oids.IfHighSpeed, Name: "ifHighSpeed", Type: metric.MetricType_INT},
		},
	}

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

// Before the task is sent we need to set the MaxRepetitions to X
func createAllVRPTransceiverMsg(el *proto.NetworkElement, conf shared.Configuration) *transport.Message {
	task := &snmpc.Task{
		Config: &snmpc.Config{
			Community:          conf.SNMP.Community,
			DynamicRepititions: false,
			MaxIterations:      2,
			MaxRepetitions:     0,
			NonRepeaters:       0,
			Version:            snmpc.SnmpVersion(conf.SNMP.Version),
			Timeout:            ptypes.DurationProto(conf.SNMP.Timeout),
			Retries:            conf.SNMP.Retries,
		},
		Type: snmpc.Type_GET,
		Oids: []*snmpc.Oid{
			{Oid: oids.HuaIfVRPOpticalVendorSN, Name: "hwEntityOpticalVendorSn", Type: metric.MetricType_STRING},
			{Oid: oids.HuaIfVRPOpticalTemperature, Name: "hwEntityOpticalTemperature", Type: metric.MetricType_INT},
			{Oid: oids.HuaIfVRPOpticalVoltage, Name: "hwEntityOpticalVoltage", Type: metric.MetricType_INT},
			{Oid: oids.HuaIfVRPOpticalBias, Name: "hwEntityOpticalBiasCurrent", Type: metric.MetricType_INT},
			{Oid: oids.HuaIfVRPOpticalRxPower, Name: "hwEntityOpticalRxPower", Type: metric.MetricType_INT},
			{Oid: oids.HuaIfVRPOpticalTxPower, Name: "hwEntityOpticalTxPower", Type: metric.MetricType_INT},
			{Oid: oids.HuaIfVRPVendorPN, Name: "hwEntityOpticalVenderPn", Type: metric.MetricType_STRING},
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
