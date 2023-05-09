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
	"github.com/araddon/dateparse"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createTaskGetPortStats(index int64, req *proto.Request, conf *config.ResourceVRP) *transport.Message {
	task := &snmpc.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpc.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: false,
			NonRepeaters:       13,
			MaxIterations:      1,
			Version:            snmpc.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
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

func (v *VRPDriver) getIfEntryInformation(m *metric.Metric, elementInterface *networkelement.Interface) {
	switch {
	case strings.HasPrefix(m.Oid, oids.IfOutOctets):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Output.Bytes = int64(i)

	case strings.HasPrefix(m.Oid, oids.IfInOctets):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Input.Bytes = int64(i)

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
		elementInterface.Type = networkelement.InterfaceType(i)

	case strings.HasPrefix(m.Oid, oids.IfMtu):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Mtu = int64(i)

	case strings.HasPrefix(m.Oid, oids.IfLastChange):
		elementInterface.LastChanged = timestamppb.New(dateparse.MustParse(m.GetValue()))

	case strings.HasPrefix(m.Oid, oids.IfPhysAddress):
		elementInterface.Hwaddress = m.GetValue()

	case strings.HasPrefix(m.Oid, oids.IfAdminStatus):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.AdminStatus = networkelement.InterfaceStatus(i)

	case strings.HasPrefix(m.Oid, oids.IfOperStatus):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.OperationalStatus = networkelement.InterfaceStatus(i)

	}
}

func (v *VRPDriver) getHuaweiInterfaceStats(m *metric.Metric, elementInterface *networkelement.Interface) {
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
