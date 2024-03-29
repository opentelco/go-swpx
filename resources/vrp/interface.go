package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
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

func createTaskGetPortStats(index int64, req *resourcepb.Request, conf *config.ResourceVRP) *transportpb.Message {
	task := &snmpcpb.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req.Timeout, defaultDeadlineTimeout))),
		Config: &snmpcpb.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: false,
			NonRepeaters:       13,
			MaxIterations:      1,
			Version:            snmpcpb.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpcpb.Type_GET,
		Oids: []*snmpcpb.Oid{

			// OUT
			{Oid: fmt.Sprintf(oids.IfOutErrorsF, index), Name: "ifOutErrors", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfHCOutOctetsF, index), Name: "ifHCOutOctets", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOutUcastPktsF, index), Name: "ifOutUcastPkts", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOutBroadcastPktsF, index), Name: "ifOutBroadcastPkts", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOutMulticastPktsF, index), Name: "ifOutMulticastPkts", Type: metricpb.MetricType_INT},

			// In
			{Oid: fmt.Sprintf(oids.IfInErrorsF, index), Name: "ifInErrors", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfHCInOctetsF, index), Name: "ifHCInOctets", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfInUcastPktsF, index), Name: "ifInUcastPkts", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfInBroadcastPktsF, index), Name: "ifInBroadcastPkts", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfInMulticastPktsF, index), Name: "ifInMulticastPkts", Type: metricpb.MetricType_INT},

			// huawei spec.
			{Oid: fmt.Sprintf(oids.HuaIfEtherStatInCRCPktsF, index), Name: "HuaIfEtherStatInCRCPkts", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.HuaIfEtherStatInPausePktsF, index), Name: "HuaIfEtherStatInPausePkts", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.HuaIfEtherStatOutPausePktsF, index), Name: "HuaIfEtherStatOutPausePkts", Type: metricpb.MetricType_INT},
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

func (v *VRPDriver) getIfEntryInformation(m *metricpb.Metric, elementInterface *devicepb.Port) {
	switch {
	case strings.HasPrefix(m.Oid, oids.IfOutOctets):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Output.Bits = int64(i) * 8 // mulitply to get bits

	case strings.HasPrefix(m.Oid, oids.IfInOctets):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Input.Bits = int64(i) * 8 // mulitply to get bits

	case strings.HasPrefix(m.Oid, oids.IfHCOutOctets):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Output.Bits = int64(i) * 8 // mulitply to get bits

	case strings.HasPrefix(m.Oid, oids.IfHCInOctets):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Input.Bits = int64(i) * 8 // mulitply to get bits

	case strings.HasPrefix(m.Oid, oids.IfInUcastPkts):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Input.Unicast = int64(i)

	case strings.HasPrefix(m.Oid, oids.IfInErrors):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Input.Errors = int64(i)

	case strings.HasPrefix(m.Oid, oids.IfOutErrors):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Output.Errors = int64(i)

	case strings.HasPrefix(m.Oid, oids.IfDescr):
		elementInterface.Description = m.GetValue()

	case strings.HasPrefix(m.Oid, oids.IfType):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Type = devicepb.Port_Type(i)

	case strings.HasPrefix(m.Oid, oids.IfMtu):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Mtu = int64(i)

	case strings.HasPrefix(m.Oid, oids.IfLastChange):
		elementInterface.LastChanged = timestamppb.New(dateparse.MustParse(m.GetValue()))

	case strings.HasPrefix(m.Oid, oids.IfPhysAddress):
		elementInterface.MacAddress = m.GetValue()

	case strings.HasPrefix(m.Oid, oids.IfAdminStatus):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.AdminStatus = devicepb.Port_Status(i)

	case strings.HasPrefix(m.Oid, oids.IfOperStatus):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.OperationalStatus = devicepb.Port_Status(i)

	}
}

func (v *VRPDriver) getHuaweiInterfaceStats(m *metricpb.Metric, elementInterface *devicepb.Port) {
	switch {
	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatInCRCPkts):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Input.CrcErrors = int64(i)

	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatInPausePkts):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Input.Pauses = int64(i)

	case strings.HasPrefix(m.Oid, oids.HuaIfEthIfStatReset):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Resets = int64(i)

	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatOutPausePkts):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Output.Pauses = int64(i)
	}
}
