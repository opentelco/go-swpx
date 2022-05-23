package main

import (
	"fmt"

	"git.liero.se/opentelco/go-dnc/models/pb/metric"
	shared2 "git.liero.se/opentelco/go-dnc/models/pb/shared"
	"git.liero.se/opentelco/go-dnc/models/pb/snmpc"
	"git.liero.se/opentelco/go-dnc/models/pb/telnet"
	"git.liero.se/opentelco/go-dnc/models/pb/transport"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	PromptStringRegex   string = "([<\\[]|)(\\S+)([>#\\]])$"
	ScreenLengthCommand string = "terminal length 0"
)

func CreateCTCTelnetInterfaceTask(el *proto.NetworkElement, conf *shared.Configuration) *transport.Message {
	task := &telnet.Task{
		Type: telnet.Type_GET,
		Payload: []*telnet.Payload{
			{
				Command: fmt.Sprintf("show mac address-table dynamic interface  %s", el.Interface),
			},
		},
		Config: &telnet.Config{
			User:                conf.Telnet.Username,
			Password:            conf.Telnet.Password,
			Port:                conf.Telnet.Port,
			ScreenLength:        ScreenLengthCommand,
			ScreenLengthCommand: ScreenLengthCommand,
			RegexPrompt:         PromptStringRegex,
			Ttl:                 &durationpb.Duration{Seconds: int64(conf.Telnet.TTL.Seconds())},
			ReadDeadLine:        &durationpb.Duration{Seconds: 10},
			WriteDeadLine:       &durationpb.Duration{Seconds: 10},
		},
		Host: el.Hostname,
	}

	message := &transport.Message{
		Id:      ksuid.New().String(),
		Target:  el.Hostname,
		Type:    transport.Type_TELNET,
		Task:    &transport.Message_Telnet{Telnet: task},
		Status:  shared2.Status_NEW,
		Created: timestamppb.Now(),
	}
	return message

}

func CreateCTCDiscoveryMsg(el *proto.NetworkElement, conf *shared.Configuration) *transport.Message {
	task := &snmpc.Task{
		Config: &snmpc.Config{
			Community:          conf.SNMP.Community,
			MaxRepetitions:     0,
			DynamicRepititions: true,
			MaxIterations:      5,
			NonRepeaters:       0,
			Version:            snmpc.SnmpVersion(conf.SNMP.Version),
			Timeout:            durationpb.New(conf.SNMP.Timeout),
			Retries:            conf.SNMP.Retries,
		},
		Type: snmpc.Type_BULKWALK,
		Oids: []*snmpc.Oid{
			{Oid: oids.IfAlias, Name: "ifAlias", Type: metric.MetricType_STRING},
			{Oid: oids.IfIndex, Name: "ifIndex", Type: metric.MetricType_INT},
			{Oid: oids.IfXEntry, Name: "ifXEntry", Type: metric.MetricType_STRING},
		},
	}

	// task.Parameters = params
	message := &transport.Message{
		Id:      ksuid.New().String(),
		Target:  el.Hostname,
		Type:    transport.Type_SNMP,
		Task:    &transport.Message_Snmpc{Snmpc: task},
		Status:  shared2.Status_NEW,
		Created: timestamppb.Now(),
	}
	return message
}
