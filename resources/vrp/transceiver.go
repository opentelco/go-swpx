package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.liero.se/opentelco/go-dnc/models/pb/metric"
	shared2 "git.liero.se/opentelco/go-dnc/models/pb/shared"
	"git.liero.se/opentelco/go-dnc/models/pb/snmpc"
	"git.liero.se/opentelco/go-dnc/models/pb/transport"
	"git.liero.se/opentelco/go-swpx/config"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelement"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createVRPTransceiverMsg(req *proto.Request, conf *config.ResourceVRP) *transport.Message {
	task := &snmpc.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req.Timeout, defaultDeadlineTimeout))),
		Config: &snmpc.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: true,
			MaxIterations:      1,
			NonRepeaters:       0,
			Version:            snmpc.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpc.Type_GET,
		Oids: []*snmpc.Oid{
			{Oid: fmt.Sprintf(oids.HuaIfVRPOpticalVendorSNF, req.PhysicalPortIndex), Name: "hwEntityOpticalVendorSn", Type: metric.MetricType_STRING},
			{Oid: fmt.Sprintf(oids.HuaIfVRPOpticalTemperatureF, req.PhysicalPortIndex), Name: "hwEntityOpticalTemperature", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.HuaIfVRPOpticalVoltageF, req.PhysicalPortIndex), Name: "hwEntityOpticalVoltage", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.HuaIfVRPOpticalBiasF, req.PhysicalPortIndex), Name: "hwEntityOpticalBiasCurrent", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.HuaIfVRPOpticalRxPowerF, req.PhysicalPortIndex), Name: "hwEntityOpticalRxPower", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.HuaIfVRPOpticalTxPowerF, req.PhysicalPortIndex), Name: "hwEntityOpticalTxPower", Type: metric.MetricType_INT},
			{Oid: fmt.Sprintf(oids.HuaIfVRPVendorPNF, req.PhysicalPortIndex), Name: "hwEntityOpticalVenderPn", Type: metric.MetricType_STRING},
		},
	}

	// task.Parameters = params
	message := &transport.Message{
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
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req.Timeout, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
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
		d.logger.Warn("could not parse transceiver (no transceiver)", "temp", tempInt, "rx", rxInt, "tx", txInt)
		return nil
	}

	val := &networkelement.Transceiver{
		SerialNumber: strings.Trim(task.Metrics[startIndex+0].GetValue(), " "),
		Stats: []*networkelement.TransceiverStatistics{
			{
				Temp:    float64(tempInt),
				Voltage: float64(voltInt) / 1000,
				Current: float64(curInt) / 1000,
				Rx:      convertToDb(int64(rxInt)),
				Tx:      convertToDb(int64(txInt)),
			},
		},
		PartNumber: task.Metrics[startIndex+6].GetValue(),
	}
	return val
}
