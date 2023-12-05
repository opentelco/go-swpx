package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.liero.se/opentelco/go-dnc/models/pb/metricpb"
	"git.liero.se/opentelco/go-dnc/models/pb/sharedpb"
	"git.liero.se/opentelco/go-dnc/models/pb/snmpcpb"
	"git.liero.se/opentelco/go-dnc/models/pb/transportpb"
	"git.liero.se/opentelco/go-swpx/config"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelementpb"
	"git.liero.se/opentelco/go-swpx/proto/go/resourcepb"
	"git.liero.se/opentelco/go-swpx/shared"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"github.com/araddon/dateparse"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createTaskGetPortStats(index int64, req *resourcepb.Request, conf *config.ResourceCTC) *transportpb.Message {
	task := &snmpcpb.Task{
		Deadline: timestamppb.New(time.Now().Add(shared.ValidateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpcpb.Config{
			Community:          conf.Snmp.Community,
			DynamicRepititions: false,
			NonRepeaters:       10,
			MaxIterations:      1,
			Version:            snmpcpb.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpcpb.Type_GET,
		Oids: []*snmpcpb.Oid{

			// OUT
			{Oid: fmt.Sprintf(oids.IfOutErrorsF, index), Name: "ifOutErrors", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOutOctetsF, index), Name: "ifOutOctets", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOutUcastPktsF, index), Name: "ifOutUcastPkts", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOutBroadcastPktsF, index), Name: "ifOutBroadcastPkts", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfOutMulticastPktsF, index), Name: "ifOutMulticastPkts", Type: metricpb.MetricType_INT},

			// In
			{Oid: fmt.Sprintf(oids.IfInErrorsF, index), Name: "ifInErrors", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfInOctetsF, index), Name: "ifInOctetsF", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfInUcastPktsF, index), Name: "ifInUcastPkts", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfInBroadcastPktsF, index), Name: "ifInBroadcastPkts", Type: metricpb.MetricType_INT},
			{Oid: fmt.Sprintf(oids.IfInMulticastPktsF, index), Name: "ifInMulticastPkts", Type: metricpb.MetricType_INT},
		},
	}

	message := &transportpb.Message{
		Session: &transportpb.Session{
			NetworkRegion: req.NetworkRegion,
			Target:        req.Hostname,
			Type:          transportpb.Type_SNMP,
		},
		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transportpb.Task{
			Task: &transportpb.Task_Snmpc{task},
		},
		Status:   sharedpb.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(shared.ValidateEOLTimeout(req, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}

	return message
}

func (d *driver) getIfEntryInformation(m *metricpb.Metric, elementInterface *networkelementpb.Port) {
	switch {
	case strings.HasPrefix(m.Oid, oids.IfOutOctets):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Output.Bits = int64(i) * 8 // mulitply to get bits

	case strings.HasPrefix(m.Oid, oids.IfInOctets):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Input.Bits = int64(i) * 8 // mulitply to get bits

	// --
	case strings.HasPrefix(m.Oid, oids.IfHCOutOctets):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Output.Bits = int64(i) * 8 // mulitply to get bits

	case strings.HasPrefix(m.Oid, oids.IfHCInOctets):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Stats.Input.Bits = int64(i) * 8 // mulitply to get bits

	// --
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
		elementInterface.Type = networkelementpb.Port_Type(i)

	case strings.HasPrefix(m.Oid, oids.IfMtu):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.Mtu = int64(i)

	case strings.HasPrefix(m.Oid, oids.IfLastChange):
		elementInterface.LastChanged = timestamppb.New(dateparse.MustParse(m.GetValue()))

	case strings.HasPrefix(m.Oid, oids.IfPhysAddress):
		elementInterface.Hwaddress = m.GetValue()

	case strings.HasPrefix(m.Oid, oids.IfAdminStatus):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.AdminStatus = networkelementpb.Port_Status(i)

	case strings.HasPrefix(m.Oid, oids.IfOperStatus):
		i, _ := strconv.Atoi(m.GetValue())
		elementInterface.OperationalStatus = networkelementpb.Port_Status(i)

	}
}
