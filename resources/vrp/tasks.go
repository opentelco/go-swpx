package main

import (
	"git.liero.se/opentelco/go-dnc/models/pb/metric"
	shared2 "git.liero.se/opentelco/go-dnc/models/pb/shared"
	"git.liero.se/opentelco/go-dnc/models/pb/snmpc"
	"git.liero.se/opentelco/go-dnc/models/pb/transport"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Before the task is sent we need to set the MaxRepetitions to X
func CreateAllVRPTransceiverMsg(el *proto.NetworkElement, conf *shared.Configuration, maxRepetitions int32) *transport.Message {
	task := &snmpc.Task{
		Deadline: el.Conf.Request.Deadline,
		Config: &snmpc.Config{
			Community:          conf.SNMP.Community,
			DynamicRepititions: false, // FALSE because right now it will lookup the ifIndex to get repetitions which we cannot rely on
			MaxIterations:      3,
			MaxRepetitions:     maxRepetitions, // set this to the number of interfaces ( db.getCollection('interface_cache').find({"hostname": "172.16.56.21"}).count(); )
			NonRepeaters:       0,              // all oids should be repeated
			Version:            snmpc.SnmpVersion(conf.SNMP.Version),
			Timeout:            durationpb.New(conf.SNMP.Timeout),
			Retries:            conf.SNMP.Retries,
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
		Session: &transport.Session{
			Target: el.Hostname,
			Port:   int32(el.Conf.SNMP.Port),
			Source: "swpx",
			Type:   transport.Type_SNMP,
		},
		Id:   ksuid.New().String(),
		Type: transport.Type_SNMP,
		Task: &transport.Task{
			Task: &transport.Task_Snmpc{task},
		},
		Status:  shared2.Status_NEW,
		Created: timestamppb.Now(),
	}
	return message
}
