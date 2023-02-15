package main

import (
	"fmt"
	"strconv"
	"strings"

	"git.liero.se/opentelco/go-dnc/models/pb/metric"
	shared2 "git.liero.se/opentelco/go-dnc/models/pb/shared"
	"git.liero.se/opentelco/go-dnc/models/pb/snmpc"
	"git.liero.se/opentelco/go-dnc/models/pb/transport"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelement"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/resources"
	"git.liero.se/opentelco/go-swpx/shared"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createVRPTransceiverMsg(el *proto.NetworkElement, conf *shared.Configuration) *transport.Message {
	task := &snmpc.Task{
		Deadline: el.Conf.Request.Deadline,
		Config: &snmpc.Config{
			Community:          conf.SNMP.Community,
			DynamicRepititions: true,
			MaxIterations:      1,
			NonRepeaters:       0,
			Version:            snmpc.SnmpVersion(conf.SNMP.Version),
			Timeout:            durationpb.New(conf.SNMP.Timeout),
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
		Status: shared2.Status_NEW,
		// RequestDeadline: el.Conf.Request.Deadline,
		Created: timestamppb.Now(),
	}

	return message
}

func (d *VRPDriver) parseTransceiverMessage(task *snmpc.Task, startIndex int) *networkelement.Transceiver {
	tempInt, _ := strconv.Atoi(task.Metrics[startIndex+1].GetValue())
	voltInt, _ := strconv.Atoi(task.Metrics[startIndex+2].GetValue())
	curInt, _ := strconv.Atoi(task.Metrics[startIndex+3].GetValue())

	rxInt, _ := strconv.Atoi(task.Metrics[startIndex+4].GetValue())
	txInt, _ := strconv.Atoi(task.Metrics[startIndex+5].GetValue())

	// no transceiver available, return nil
	if tempInt == -255 || rxInt == -1 || txInt == -1 {
		return nil
	}

	val := &networkelement.Transceiver{
		SerialNumber: strings.Trim(task.Metrics[startIndex+0].GetValue(), " "),
		Stats: []*networkelement.TransceiverStatistics{
			{
				Temp:    float64(tempInt),
				Voltage: float64(voltInt) / 1000,
				Current: float64(curInt) / 1000,
				Rx:      resources.ConvertToDb(int64(rxInt)),
				Tx:      resources.ConvertToDb(int64(txInt)),
			},
		},
		PartNumber: task.Metrics[startIndex+6].GetValue(),
	}
	return val
}
