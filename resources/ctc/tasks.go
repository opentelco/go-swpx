package main

import (
	"fmt"

	"git.liero.se/opentelco/go-dnc/models/pb/metric"
	shared2 "git.liero.se/opentelco/go-dnc/models/pb/shared"
	"git.liero.se/opentelco/go-dnc/models/pb/snmpc"
	"git.liero.se/opentelco/go-dnc/models/pb/terminal"
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
	task := &terminal.Task{
		Deadline: el.Conf.Request.Deadline,
		Payload: []*terminal.Task_Payload{
			{
				Command: fmt.Sprintf("show mac address-table dynamic interface  %s", el.Interface),
			},
		},
		Config: &terminal.Config{
			User:                conf.Ssh.Username,
			Password:            conf.Ssh.Password,
			RegexPrompt:         conf.Ssh.RegexPrompt,
			ScreenLengthCommand: conf.Ssh.ScreenLengthCommand,
			ReadDeadLine:        &durationpb.Duration{Seconds: int64(conf.Ssh.ReadDeadLine.Seconds())},
			WriteDeadLine:       &durationpb.Duration{Seconds: int64(conf.Ssh.WriteDeadLine.Seconds())},
			SshKeyPath:          conf.Ssh.SSHKeyPath,
		},
	}

	message := &transport.Message{
		Session: &transport.Session{
			Target: el.Hostname,
			Port:   int32(el.Conf.Telnet.Port),
			Source: "swpx",
			Type:   transport.Type_TELNET,
		},
		Id:   ksuid.New().String(),
		Type: transport.Type_TELNET,
		Task: &transport.Task{
			Task: &transport.Task_Terminal{task},
		},
		Status:  shared2.Status_NEW,
		Created: timestamppb.Now(),
		// RequestDeadline: el.Conf.Request.Deadline,
	}
	return message

}

func CreateCTCDiscoveryMsg(el *proto.NetworkElement, conf *shared.Configuration) *transport.Message {
	task := &snmpc.Task{
		Deadline: el.Conf.Request.Deadline,
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

	message := &transport.Message{
		Session: &transport.Session{
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
