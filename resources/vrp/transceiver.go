package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
	"go.opentelco.io/go-dnc/models/pb/metricpb"
	"go.opentelco.io/go-dnc/models/pb/sharedpb"
	"go.opentelco.io/go-dnc/models/pb/snmpcpb"
	"go.opentelco.io/go-dnc/models/pb/transportpb"
	"go.opentelco.io/go-swpx/config"
	"go.opentelco.io/go-swpx/proto/go/devicepb"
	"go.opentelco.io/go-swpx/proto/go/resourcepb"
	"go.opentelco.io/go-swpx/shared/oids"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createVRPTransceiverMsg(req *resourcepb.Request, conf *config.ResourceVRP) *transportpb.Message {
	task := &snmpcpb.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req.Timeout, defaultDeadlineTimeout))),
		Config: &snmpcpb.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: true,
			MaxIterations:      1,
			NonRepeaters:       0,
			Version:            snmpcpb.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpcpb.Type_GET,
		Oids: []*snmpcpb.Oid{
			{
				Oid:  fmt.Sprintf(oids.HuaIfVRPOpticalVendorSNF, req.PhysicalPortIndex),
				Name: "hwEntityOpticalVendorSn",
				Type: metricpb.MetricType_STRING,
			},
			{
				Oid:  fmt.Sprintf(oids.HuaIfVRPOpticalTemperatureF, req.PhysicalPortIndex),
				Name: "hwEntityOpticalTemperature",
				Type: metricpb.MetricType_INT,
			},
			{
				Oid:  fmt.Sprintf(oids.HuaIfVRPOpticalVoltageF, req.PhysicalPortIndex),
				Name: "hwEntityOpticalVoltage",
				Type: metricpb.MetricType_INT,
			},
			{
				Oid:  fmt.Sprintf(oids.HuaIfVRPOpticalBiasF, req.PhysicalPortIndex),
				Name: "hwEntityOpticalBiasCurrent",
				Type: metricpb.MetricType_INT,
			},
			{
				Oid:  fmt.Sprintf(oids.HuaIfVRPOpticalRxPowerF, req.PhysicalPortIndex),
				Name: "hwEntityOpticalRxPower",
				Type: metricpb.MetricType_INT,
			},
			{
				Oid:  fmt.Sprintf(oids.HuaIfVRPOpticalTxPowerF, req.PhysicalPortIndex),
				Name: "hwEntityOpticalTxPower",
				Type: metricpb.MetricType_INT,
			},
			{
				Oid:  fmt.Sprintf(oids.HuaIfVRPVendorPNF, req.PhysicalPortIndex),
				Name: "hwEntityOpticalVenderPn",
				Type: metricpb.MetricType_STRING,
			},
		},
	}

	// task.Parameters = params
	message := &transportpb.Message{
		Session: &transportpb.Session{
			NetworkRegion: req.NetworkRegion,
			Target:        req.Hostname,
			Port:          int32(conf.Snmp.Port),
			Type:          transportpb.Type_SNMP,
		},
		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transportpb.Task{
			Task: &transportpb.Task_Snmpc{task},
		},
		Status:   sharedpb.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req.Timeout, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}

	return message
}

func (d *VRPDriver) parseTransceiverMessage(task *snmpcpb.Task, startIndex int) *devicepb.Transceiver {
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

	val := &devicepb.Transceiver{
		SerialNumber: strings.Trim(task.Metrics[startIndex+0].GetValue(), " "),
		Stats: &devicepb.Transceiver_Statistics{
			Temp:    float64(tempInt),
			Voltage: float64(voltInt) / 1000,
			Current: float64(curInt) / 1000,
			Rx:      convertToDb(int64(rxInt)),
			Tx:      convertToDb(int64(txInt)),
		},
		PartNumber: task.Metrics[startIndex+6].GetValue(),
	}
	return val
}
