package main

import (
	"fmt"
	"time"

	"go.opentelco.io/go-dnc/models/pb/metricpb"
	"go.opentelco.io/go-dnc/models/pb/snmpcpb"
	"go.opentelco.io/go-dnc/models/pb/transportpb"
	"go.opentelco.io/go-swpx/config"
	"go.opentelco.io/go-swpx/proto/go/resourcepb"
	"go.opentelco.io/go-swpx/shared"
	"go.opentelco.io/go-swpx/shared/oids"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createPhysicalPortIndex(req *resourcepb.Request, conf *config.ResourceCTC) *transportpb.Message {
	task := &snmpcpb.Task{
		Deadline: timestamppb.New(time.Now().Add(shared.ValidateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpcpb.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: true,
			MaxIterations:      1,
			NonRepeaters:       0,
			Version:            snmpcpb.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpcpb.Type_BULKWALK,
		Oids: []*snmpcpb.Oid{
			{Oid: oids.IfEntPhysicalName, Name: "ifPhysAddress", Type: metricpb.MetricType_STRING},
		},
	}

	message := createSnmpMessage(req, conf)
	message.Task.Task = &transportpb.Task_Snmpc{Snmpc: task}
	return message
}

// createAllVRPTransceiverMsg creates a message for all transceivers on the device with
// the help of bulk walk
func createAllVRPTransceiverMsg(req *resourcepb.Request, conf *config.ResourceCTC, maxRepetitions int32) *transportpb.Message {
	task := &snmpcpb.Task{
		Deadline: timestamppb.New(time.Now().Add(shared.ValidateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpcpb.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: false, // FALSE because right now it will lookup the ifIndex to get repetitions which we cannot rely on
			MaxIterations:      3,
			MaxRepetitions:     maxRepetitions, // set this to the number of interfaces
			NonRepeaters:       0,              // all oids should be repeated
			Version:            snmpcpb.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpcpb.Type_BULKWALK,
		Oids: []*snmpcpb.Oid{
			{Oid: oids.HuaIfVRPOpticalVendorSN, Name: "hwEntityOpticalVendorSn", Type: metricpb.MetricType_STRING},
			{Oid: oids.HuaIfVRPOpticalTemperature, Name: "hwEntityOpticalTemperature", Type: metricpb.MetricType_INT},
			{Oid: oids.HuaIfVRPOpticalVoltage, Name: "hwEntityOpticalVoltage", Type: metricpb.MetricType_INT},
			{Oid: oids.HuaIfVRPOpticalBias, Name: "hwEntityOpticalBiasCurrent", Type: metricpb.MetricType_INT},
			{Oid: oids.HuaIfVRPOpticalRxPower, Name: "hwEntityOpticalRxPower", Type: metricpb.MetricType_INT},
			{Oid: oids.HuaIfVRPOpticalTxPower, Name: "hwEntityOpticalTxPower", Type: metricpb.MetricType_INT},
			{Oid: oids.HuaIfVRPVendorPN, Name: "hwEntityOpticalVenderPn", Type: metricpb.MetricType_STRING},
		},
	}

	message := createSnmpMessage(req, conf)
	message.Task.Task = &transportpb.Task_Snmpc{Snmpc: task}
	return message
}

func createLogicalPortIndex(req *resourcepb.Request, conf *config.ResourceCTC) *transportpb.Message {
	task := &snmpcpb.Task{
		Deadline: timestamppb.New(time.Now().Add(shared.ValidateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpcpb.Config{
			Community:          conf.Snmp.Community,
			MaxRepetitions:     72,
			DynamicRepititions: true,
			MaxIterations:      2,
			NonRepeaters:       0,
			Version:            snmpcpb.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpcpb.Type_BULKWALK,
		Oids: []*snmpcpb.Oid{
			{Oid: oids.IfIndex, Name: "ifIndex", Type: metricpb.MetricType_INT},
			{Oid: oids.IfDescr, Name: "ifDescr", Type: metricpb.MetricType_STRING},
			{Oid: oids.IfAlias, Name: "ifAlias", Type: metricpb.MetricType_STRING},
		},
	}

	message := createSnmpMessage(req, conf)
	message.Task.Task = &transportpb.Task_Snmpc{Snmpc: task}
	return message
}

func createTaskSystemInfo(req *resourcepb.Request, conf *config.ResourceCTC) *transportpb.Message {
	task := &snmpcpb.Task{
		Deadline: timestamppb.New(time.Now().Add(shared.ValidateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpcpb.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: false,
			NonRepeaters:       12,
			MaxIterations:      1,
			Version:            snmpcpb.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpcpb.Type_GET,
		Oids: []*snmpcpb.Oid{

			// OUT
			{Oid: oids.SysDescr, Name: "sysDescr", Type: metricpb.MetricType_STRING},
			{Oid: oids.SysObjectID, Name: "sysObjectID", Type: metricpb.MetricType_STRING},
			{Oid: oids.SysUpTime, Name: "sysUpTime", Type: metricpb.MetricType_STRING},
			{Oid: oids.SysContact, Name: "sysContact", Type: metricpb.MetricType_STRING},
			{Oid: oids.SysName, Name: "sysName", Type: metricpb.MetricType_STRING},
			{Oid: oids.SysLocation, Name: "sysLocation", Type: metricpb.MetricType_STRING},
			{Oid: oids.SysORLastChange, Name: "sysORLastChange", Type: metricpb.MetricType_INT},
		},
	}

	message := createSnmpMessage(req, conf)
	message.Task.Task = &transportpb.Task_Snmpc{Snmpc: task}
	return message
}

func createAllPortsMsg(req *resourcepb.Request, conf *config.ResourceCTC) *transportpb.Message {
	task := &snmpcpb.Task{
		Deadline: timestamppb.New(time.Now().Add(shared.ValidateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpcpb.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: true,
			MaxIterations:      1,
			NonRepeaters:       0,
			Version:            snmpcpb.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpcpb.Type_BULKWALK,
		Oids: []*snmpcpb.Oid{
			{Oid: oids.IfDescr, Name: "ifDescr", Type: metricpb.MetricType_STRING},
			{Oid: oids.IfAlias, Name: "ifAlias", Type: metricpb.MetricType_STRING},
			{Oid: oids.IfType, Name: "ifType", Type: metricpb.MetricType_INT},
			{Oid: oids.IfMtu, Name: "ifMtu", Type: metricpb.MetricType_INT},
			{Oid: oids.IfPhysAddress, Name: "ifPhysAddress", Type: metricpb.MetricType_HWADDR},
			{Oid: oids.IfAdminStatus, Name: "ifAdminStatus", Type: metricpb.MetricType_INT},
			{Oid: oids.IfOperStatus, Name: "ifOperStatus", Type: metricpb.MetricType_INT},
			{Oid: oids.IfLastChange, Name: "ifLastChange", Type: metricpb.MetricType_TIMETICKS},

			{Oid: oids.IfHighSpeed, Name: "ifHighSpeed", Type: metricpb.MetricType_INT},
		},
	}

	message := createSnmpMessage(req, conf)
	message.Task.Task = &transportpb.Task_Snmpc{Snmpc: task}
	return message
}

func createSinglePortMsg(index int64, req *resourcepb.Request, conf *config.ResourceCTC) *transportpb.Message {
	task := &snmpcpb.Task{
		Deadline: timestamppb.New(time.Now().Add(shared.ValidateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpcpb.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: false,
			NonRepeaters:       12,
			MaxIterations:      1,
			Version:            snmpcpb.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpcpb.Type_GET,
		Oids: []*snmpcpb.Oid{
			{Oid: fmt.Sprintf(oids.IfDescrF, index), Name: "ifDescr", Type: metricpb.MetricType_STRING},
			{Oid: fmt.Sprintf(oids.IfAliasF, index), Name: "ifAlias", Type: metricpb.MetricType_STRING},
			{Oid: fmt.Sprintf(oids.IfTypeF, index), Name: "ifType", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfMtuF, index), Name: "ifMtu", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfPhysAddressF, index), Name: "ifPhysAddress", Type: metricpb.MetricType_HWADDR},
			{Oid: fmt.Sprintf(oids.IfAdminStatusF, index), Name: "ifAdminStatus", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOperStatusF, index), Name: "ifOperStatus", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfLastChangeF, index), Name: "ifLastChange", Type: metricpb.MetricType_TIMETICKS},
			{Oid: fmt.Sprintf(oids.IfHighSpeedF, index), Name: "ifHighSpeed", Type: metricpb.MetricType_INT},
		},
	}

	message := createSnmpMessage(req, conf)
	message.Task.Task = &transportpb.Task_Snmpc{Snmpc: task}
	return message
}
func createSinglePortMsgShort(index int64, req *resourcepb.Request, conf *config.ResourceCTC) *transportpb.Message {
	task := &snmpcpb.Task{
		Deadline: timestamppb.New(time.Now().Add(shared.ValidateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpcpb.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: false,
			NonRepeaters:       5,
			MaxIterations:      1,
			Version:            snmpcpb.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpcpb.Type_GET,
		Oids: []*snmpcpb.Oid{
			{Oid: fmt.Sprintf(oids.IfDescrF, index), Name: "ifDescr", Type: metricpb.MetricType_STRING},
			{Oid: fmt.Sprintf(oids.IfAliasF, index), Name: "ifAlias", Type: metricpb.MetricType_STRING},
			{Oid: fmt.Sprintf(oids.IfAdminStatusF, index), Name: "ifAdminStatus", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOperStatusF, index), Name: "ifOperStatus", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfHighSpeedF, index), Name: "ifHighSpeed", Type: metricpb.MetricType_INT},
		},
	}

	message := createSnmpMessage(req, conf)
	message.Task.Task = &transportpb.Task_Snmpc{Snmpc: task}
	return message
}

// CreateMsg uses the pbuf transportpb for DNC..
func createMsg(req *resourcepb.Request, conf *config.ResourceCTC) *transportpb.Message {

	task := &snmpcpb.Task{
		// Deadline: timestamppb.New(dl),
		Config: &snmpcpb.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: true,
			MaxIterations:      200,
			Version:            snmpcpb.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpcpb.Type_BULKWALK,
		Oids: []*snmpcpb.Oid{

			{Oid: oids.IfIndex, Name: "ifIndex", Type: metricpb.MetricType_INT},
			{Oid: oids.IfDescr, Name: "ifDescr", Type: metricpb.MetricType_STRING},
			{Oid: oids.IfType, Name: "ifType", Type: metricpb.MetricType_INT},
			{Oid: oids.IfMtu, Name: "ifMtu", Type: metricpb.MetricType_INT},
			{Oid: oids.IfPhysAddress, Name: "ifPhysAddress", Type: metricpb.MetricType_HWADDR},
			{Oid: oids.IfAdminStatus, Name: "ifAdminStatus", Type: metricpb.MetricType_INT},
			{Oid: oids.IfOperStatus, Name: "ifOperStatus", Type: metricpb.MetricType_INT},
			{Oid: oids.IfLastChange, Name: "ifLastChange", Type: metricpb.MetricType_TIMETICKS},
			{Oid: oids.IfInErrors, Name: "ifInErrors", Type: metricpb.MetricType_INT},
			{Oid: oids.IfOutErrors, Name: "ifOutErrors", Type: metricpb.MetricType_INT},

			{Oid: oids.IfHighSpeed, Name: "ifHighSpeed", Type: metricpb.MetricType_INT},
			{Oid: oids.IfHCOutUcastPkts, Name: "ifHCOutUcastPkts", Type: metricpb.MetricType_INT},
		},
	}

	message := createSnmpMessage(req, conf)
	message.Task.Task = &transportpb.Task_Snmpc{Snmpc: task}
	return message
}

func createCTCDiscoveryMsg(req *resourcepb.Request, conf *config.ResourceCTC) *transportpb.Message {
	task := &snmpcpb.Task{
		Deadline: timestamppb.New(time.Now().Add(shared.ValidateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpcpb.Config{
			Community:          conf.Snmp.Community,
			MaxRepetitions:     0,
			DynamicRepititions: true,
			MaxIterations:      5,
			NonRepeaters:       0,
			Version:            snmpcpb.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpcpb.Type_BULKWALK,
		Oids: []*snmpcpb.Oid{
			{Oid: oids.IfAlias, Name: "ifAlias", Type: metricpb.MetricType_STRING},
			{Oid: oids.IfIndex, Name: "ifIndex", Type: metricpb.MetricType_INT},
			{Oid: oids.IfXEntry, Name: "ifXEntry", Type: metricpb.MetricType_STRING},
		},
	}

	message := createSnmpMessage(req, conf)
	message.Task.Task = &transportpb.Task_Snmpc{Snmpc: task}
	return message
}
