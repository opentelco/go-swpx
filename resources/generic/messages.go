package main

import (
	"time"

	"git.liero.se/opentelco/go-dnc/models/pb/metric"
	"git.liero.se/opentelco/go-dnc/models/pb/shared"
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

func createTaskSystemInfo(req *proto.Request, conf *config.Snmp) *transport.Message {
	task := &snmpc.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req.Timeout, defaultDeadlineTimeout))),
		Config: &snmpc.Config{
			Community:          conf.Community,
			DynamicRepititions: false,
			NonRepeaters:       12,
			MaxIterations:      1,
			Version:            snmpc.SnmpVersion(conf.Version),
			Timeout:            durationpb.New(conf.Timeout.AsDuration()),
			Retries:            int32(conf.Retries),
		},
		Type: snmpc.Type_GET,
		Oids: []*snmpc.Oid{

			// OUT
			{Oid: oids.SysDescr, Name: "sysDescr", Type: metric.MetricType_STRING},
			{Oid: oids.SysObjectID, Name: "sysObjectID", Type: metric.MetricType_STRING},
			{Oid: oids.SysUpTime, Name: "sysUpTime", Type: metric.MetricType_TIMETICKS},
			{Oid: oids.SysContact, Name: "sysContact", Type: metric.MetricType_STRING},
			{Oid: oids.SysName, Name: "sysName", Type: metric.MetricType_STRING},
			{Oid: oids.SysLocation, Name: "sysLocation", Type: metric.MetricType_STRING},
			// {Oid: oids.SysORLastChange, Name: "sysORLastChange", Type: metric.MetricType_TIMETICKS},
		},
	}

	// task.Parameters = params
	message := &transport.Message{
		Session: &transport.Session{
			Target: req.Hostname,
			Port:   int32(conf.Port),

			Type: transport.Type_SNMP,
		},
		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transport.Task{
			Task: &transport.Task_Snmpc{Snmpc: task},
		},
		Status:   shared.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req.Timeout, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}
	return message
}

func parseSystemInformation(m *metric.Metric, ne *networkelement.Element) {
	switch m.Oid {
	case oids.SysContact:
		ne.Contact = m.GetValue()
	case oids.SysDescr:
		ne.Version = m.GetValue()
	case oids.SysLocation:
		ne.Location = m.GetValue()
	case oids.SysName:
		ne.Sysname = m.GetValue()

	// case oids.SysORLastChange:
	// 	ne.LastChanged = m.GetValue()

	case oids.SysObjectID:
		ne.SnmpObjectId = m.GetValue()

	case oids.SysUpTime:
		ne.Uptime = m.GetValue()
	}
}
