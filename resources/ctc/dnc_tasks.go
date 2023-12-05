package main

import (
	"fmt"
	"time"

	"go.opentelco.io/go-dnc/models/pb/terminalpb"
	"go.opentelco.io/go-dnc/models/pb/transportpb"
	"go.opentelco.io/go-swpx/config"
	"go.opentelco.io/go-swpx/proto/go/resourcepb"
	"go.opentelco.io/go-swpx/shared"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createCTCSSHInterfaceTask(req *resourcepb.Request, conf *config.ResourceCTC) *transportpb.Message {
	sshConf := conf.Transports.GetByLabel("ssh")

	task := &terminalpb.Task{
		Deadline: timestamppb.New(time.Now().Add(shared.ValidateEOLTimeout(req, defaultDeadlineTimeout))),
		Payload: []*terminalpb.Task_Payload{
			{
				Command: fmt.Sprintf("show mac address-table dynamic interface  %s", req.Port),
			},
		},
		Config: &terminalpb.Config{
			User:                sshConf.User,
			Password:            sshConf.Password,
			RegexPrompt:         sshConf.RegexPrompt,
			ScreenLengthCommand: sshConf.ScreenLength, // must be sent down to the terminal
			ReadDeadLine:        durationpb.New(sshConf.ReadDeadLine.AsDuration()),
			WriteDeadLine:       durationpb.New(sshConf.WriteDeadLine.AsDuration()),
			SshKeyPath:          sshConf.SSHKeyPath,
		},
	}

	message := createSSHMessage(req, sshConf)
	message.Task.Task = &transportpb.Task_Terminal{Terminal: task}
	return message
}
