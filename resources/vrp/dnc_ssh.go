package main

import (
	"fmt"
	"time"

	shared2 "git.liero.se/opentelco/go-dnc/models/pb/shared"
	"git.liero.se/opentelco/go-dnc/models/pb/terminal"
	"git.liero.se/opentelco/go-dnc/models/pb/transport"
	"git.liero.se/opentelco/go-swpx/config"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createBasicSSHInterfaceTask(req *proto.Request, conf *config.ResourceVRP) *transport.Message {

	sshConf := conf.Transports.GetByLabel("ssh")
	task := &terminal.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req, defaultDeadlineTimeout))),
		Payload: []*terminal.Task_Payload{
			{
				Command: fmt.Sprintf("display mac-address %s", req.Port),
			},
		},
		Config: &terminal.Config{
			User:                sshConf.User,
			Password:            sshConf.Password,
			RegexPrompt:         sshConf.RegexPrompt,
			ScreenLengthCommand: sshConf.ScreenLength,
			ReadDeadLine:        durationpb.New(sshConf.ReadDeadLine.AsDuration()),
			WriteDeadLine:       durationpb.New(sshConf.WriteDeadLine.AsDuration()),
			SshKeyPath:          sshConf.SSHKeyPath,
		},
	}

	message := &transport.Message{
		Session: &transport.Session{
			Target: req.Hostname,
			Port:   int32(sshConf.Port),
			Source: "swpx",
			Type:   transport.Type_SSH,
		},
		Id:   ksuid.New().String(),
		Type: transport.Type_SSH,
		Task: &transport.Task{
			Task: &transport.Task_Terminal{task},
		},
		Status:  shared2.Status_NEW,
		Created: timestamppb.Now(),
	}
	return message
}
