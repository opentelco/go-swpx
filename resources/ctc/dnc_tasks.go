package main

import (
	"fmt"
	"time"

	"git.liero.se/opentelco/go-dnc/models/pb/metric"
	shared2 "git.liero.se/opentelco/go-dnc/models/pb/shared"
	"git.liero.se/opentelco/go-dnc/models/pb/snmpc"
	"git.liero.se/opentelco/go-dnc/models/pb/terminal"
	"git.liero.se/opentelco/go-dnc/models/pb/transport"
	"git.liero.se/opentelco/go-swpx/config"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createCTCSSHInterfaceTask(req *proto.Request, conf *config.ResourceCTC) *transport.Message {
	sshConf := conf.Transports.GetByLabel("ssh")

	task := &terminal.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Payload: []*terminal.Task_Payload{
			{
				Command: fmt.Sprintf("show mac address-table dynamic interface  %s", req.Port),
			},
		},
		Config: &terminal.Config{
			User:                sshConf.User,
			Password:            sshConf.Password,
			RegexPrompt:         sshConf.RegexPrompt,
			ScreenLengthCommand: sshConf.ScreenLength, // must be sent down to the terminal
			ReadDeadLine:        durationpb.New(sshConf.ReadDeadLine.AsDuration()),
			WriteDeadLine:       durationpb.New(sshConf.WriteDeadLine.AsDuration()),
			SshKeyPath:          sshConf.SSHKeyPath,
		},
	}

	message := &transport.Message{
		NetworkRegion: req.NetworkRegion,
		Session: &transport.Session{
			Target: req.Hostname,
			Port:   int32(sshConf.Port),
			Type:   transport.Type_SSH,
		},
		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transport.Task{
			Task: &transport.Task_Terminal{task},
		},
		Status:   shared2.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}
	return message

}

func createCTCDiscoveryMsg(req *proto.Request, conf *config.ResourceCTC) *transport.Message {
	task := &snmpc.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Config: &snmpc.Config{
			Community:          conf.Snmp.Community,
			MaxRepetitions:     0,
			DynamicRepititions: true,
			MaxIterations:      5,
			NonRepeaters:       0,
			Version:            snmpc.SnmpVersion(conf.Snmp.Version),
			Timeout:            durationpb.New(conf.Snmp.Timeout.AsDuration()),
			Retries:            int32(conf.Snmp.Retries),
		},
		Type: snmpc.Type_BULKWALK,
		Oids: []*snmpc.Oid{
			{Oid: oids.IfAlias, Name: "ifAlias", Type: metric.MetricType_STRING},
			{Oid: oids.IfIndex, Name: "ifIndex", Type: metric.MetricType_INT},
			{Oid: oids.IfXEntry, Name: "ifXEntry", Type: metric.MetricType_STRING},
		},
	}

	message := &transport.Message{
		NetworkRegion: req.NetworkRegion,
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
