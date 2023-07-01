package main

import (
	"fmt"
	"time"

	"git.liero.se/opentelco/go-dnc/models/pb/metric"
	"git.liero.se/opentelco/go-dnc/models/pb/shared"
	shared2 "git.liero.se/opentelco/go-dnc/models/pb/shared"
	"git.liero.se/opentelco/go-dnc/models/pb/snmpc"
	"git.liero.se/opentelco/go-dnc/models/pb/transport"
	"git.liero.se/opentelco/go-swpx/config"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createPhysicalPortIndex(request *proto.Request, conf *config.ResourceVRP) *transport.Message {
	task := &snmpc.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(request, defaultDeadlineTimeout))),
		Config: &snmpc.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: true,
			MaxIterations:      1,
			NonRepeaters:       0,
			Version:            snmpc.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpc.Type_BULKWALK,
		Oids: []*snmpc.Oid{
			{Oid: oids.IfEntPhysicalName, Name: "ifPhysAddress", Type: metric.MetricType_STRING},
		},
	}

	// task.Parameters = params
	message := &transport.Message{
		NetworkRegion: request.NetworkRegion,
		Session: &transport.Session{
			Target: request.Hostname,
			Port:   int32(conf.Snmp.Port),
			Type:   transport.Type_SNMP,
		},
		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transport.Task{
			Task: &transport.Task_Snmpc{task},
		},
		Status:   shared.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(request, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}
	return message
}

// createAllVRPTransceiverMsg creates a message for all transceivers on the device with
// the help of bulk walk
func createAllVRPTransceiverMsg(request *proto.Request, conf *config.ResourceVRP, maxRepetitions int32) *transport.Message {
	task := &snmpc.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(request, defaultDeadlineTimeout))),
		Config: &snmpc.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: false, // FALSE because right now it will lookup the ifIndex to get repetitions which we cannot rely on
			MaxIterations:      3,
			MaxRepetitions:     maxRepetitions, // set this to the number of interfaces
			NonRepeaters:       0,              // all oids should be repeated
			Version:            snmpc.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpc.Type_BULKWALK,
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

	message := &transport.Message{
		NetworkRegion: request.NetworkRegion,
		Session: &transport.Session{
			Target: request.Hostname,
			Port:   int32(conf.Snmp.Port),
			Type:   transport.Type_SNMP,
		},
		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transport.Task{
			Task: &transport.Task_Snmpc{task},
		},
		Status:   shared.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(request, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}
	return message
}

func createLogicalPortIndex(req *proto.Request, conf *config.ResourceVRP) *transport.Message {
	task := &snmpc.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpc.Config{
			Community:          conf.Snmp.Community,
			MaxRepetitions:     72,
			DynamicRepititions: true,
			MaxIterations:      2,
			NonRepeaters:       0,
			Version:            snmpc.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpc.Type_BULKWALK,
		Oids: []*snmpc.Oid{
			{Oid: oids.IfIndex, Name: "ifIndex", Type: metric.MetricType_INT},
			{Oid: oids.IfDescr, Name: "ifDescr", Type: metric.MetricType_STRING},
			{Oid: oids.IfAlias, Name: "ifAlias", Type: metric.MetricType_STRING},
		},
	}

	// task.Parameters = params
	message := &transport.Message{
		NetworkRegion: req.NetworkRegion,
		Session: &transport.Session{
			Target: req.Hostname,
			Port:   int32(conf.Snmp.Port),
			Type:   transport.Type_SNMP,
		},
		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transport.Task{
			Task: &transport.Task_Snmpc{task},
		},
		Status:   shared.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}
	return message
}

func createTaskSystemInfo(req *proto.Request, conf *config.ResourceVRP) *transport.Message {
	task := &snmpc.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpc.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: false,
			NonRepeaters:       12,
			MaxIterations:      1,
			Version:            snmpc.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
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
		NetworkRegion: req.NetworkRegion,
		Session: &transport.Session{
			Target: req.Hostname,
			Port:   int32(conf.Snmp.Port),
			Type:   transport.Type_SNMP,
		},
		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transport.Task{
			Task: &transport.Task_Snmpc{task},
		},
		Status:   shared.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}
	return message
}

func createAllPortsMsg(req *proto.Request, conf *config.ResourceVRP) *transport.Message {
	task := &snmpc.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpc.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: true,
			MaxIterations:      1,
			NonRepeaters:       0,
			Version:            snmpc.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpc.Type_BULKWALK,
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
		NetworkRegion: req.NetworkRegion,
		Session: &transport.Session{
			Target: req.Hostname,
			Port:   int32(conf.Snmp.Port),
			Type:   transport.Type_SNMP,
		},
		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transport.Task{
			Task: &transport.Task_Snmpc{task},
		},
		Status:   shared2.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}
	return message
}

func createSinglePortMsg(index int64, req *proto.Request, conf *config.ResourceVRP) *transport.Message {
	task := &snmpc.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpc.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: false,
			NonRepeaters:       12,
			MaxIterations:      1,
			Version:            snmpc.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
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

	message := &transport.Message{
		NetworkRegion: req.NetworkRegion,
		Session: &transport.Session{
			Target: req.Hostname,
			Port:   int32(conf.Snmp.Port),
			Type:   transport.Type_SNMP,
		},
		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transport.Task{
			Task: &transport.Task_Snmpc{task},
		},
		Status:   shared2.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}
	return message
}
func createSinglePortMsgShort(index int64, req *proto.Request, conf *config.ResourceCTC) *transport.Message {
	task := &snmpc.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpc.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: false,
			NonRepeaters:       5,
			MaxIterations:      1,
			Version:            snmpc.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpc.Type_GET,
		Oids: []*snmpc.Oid{
			{Oid: fmt.Sprintf(oids.IfDescrF, index), Name: "ifDescr", Type: metric.MetricType_STRING},
			{Oid: fmt.Sprintf(oids.IfAliasF, index), Name: "ifAlias", Type: metric.MetricType_STRING},
			{Oid: fmt.Sprintf(oids.IfAdminStatusF, index), Name: "ifAdminStatus", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOperStatusF, index), Name: "ifOperStatus", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfHighSpeedF, index), Name: "ifHighSpeed", Type: metric.MetricType_INT},
		},
	}

	message := &transport.Message{
		NetworkRegion: req.NetworkRegion,
		Session: &transport.Session{
			Target: req.Hostname,
			Port:   int32(conf.Snmp.Port),
			Type:   transport.Type_SNMP,
		},
		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transport.Task{
			Task: &transport.Task_Snmpc{task},
		},
		Status:   shared2.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}
	return message
}

// CreateMsg uses the pbuf transport for DNC..
func createMsg(req *proto.Request, conf *config.ResourceCTC) *transport.Message {

	task := &snmpc.Task{
		// Deadline: timestamppb.New(dl),
		Config: &snmpc.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: true,
			MaxIterations:      200,
			Version:            snmpc.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpc.Type_BULKWALK,
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

	message := &transport.Message{
		NetworkRegion: req.NetworkRegion,
		Session: &transport.Session{
			Target: req.Hostname,
			Type:   transport.Type_SNMP,
		},
		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transport.Task{
			Task: &transport.Task_Snmpc{task},
		},
		Status:   shared2.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}
	return message
}
