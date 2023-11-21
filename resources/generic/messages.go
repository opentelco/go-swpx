package main

import (
	"time"

	"git.liero.se/opentelco/go-dnc/models/pb/metricpb"
	"git.liero.se/opentelco/go-dnc/models/pb/sharedpb"
	"git.liero.se/opentelco/go-dnc/models/pb/snmpcpb"
	"git.liero.se/opentelco/go-dnc/models/pb/transportpb"
	"git.liero.se/opentelco/go-swpx/config"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelementpb"
	"git.liero.se/opentelco/go-swpx/proto/go/resourcepb"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createTaskSystemInfo(req *resourcepb.Request, conf *config.Snmp) *transportpb.Message {
	task := &snmpcpb.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req.Timeout, defaultDeadlineTimeout))),
		Config: &snmpcpb.Config{
			Community:          conf.Community,
			DynamicRepititions: false,
			NonRepeaters:       12,
			MaxIterations:      1,
			Version:            snmpcpb.SnmpVersion(conf.Version),
			Timeout:            durationpb.New(conf.Timeout.AsDuration()),
			Retries:            int32(conf.Retries),
		},
		Type: snmpcpb.Type_GET,
		Oids: []*snmpcpb.Oid{

			// OUT
			{Oid: oids.SysDescr, Name: "sysDescr", Type: metricpb.MetricType_STRING},
			{Oid: oids.SysObjectID, Name: "sysObjectID", Type: metricpb.MetricType_STRING},
			{Oid: oids.SysUpTime, Name: "sysUpTime", Type: metricpb.MetricType_TIMETICKS},
			{Oid: oids.SysContact, Name: "sysContact", Type: metricpb.MetricType_STRING},
			{Oid: oids.SysName, Name: "sysName", Type: metricpb.MetricType_STRING},
			{Oid: oids.SysLocation, Name: "sysLocation", Type: metricpb.MetricType_STRING},
			// {Oid: oids.SysORLastChange, Name: "sysORLastChange", Type: metricpb.MetricType_TIMETICKS},
		},
	}

	// task.Parameters = params
	message := &transportpb.Message{
		Session: &transportpb.Session{
			NetworkRegion: req.NetworkRegion,
			Target:        req.Hostname,
			Port:          int32(conf.Port),

			Type: transportpb.Type_SNMP,
		},
		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transportpb.Task{
			Task: &transportpb.Task_Snmpc{Snmpc: task},
		},
		Status:   sharedpb.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req.Timeout, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}
	return message
}

func parseSystemInformation(m *metricpb.Metric, ne *networkelementpb.Element) {
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
