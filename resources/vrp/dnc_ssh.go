package main

import (
	"fmt"
	"time"

	"git.liero.se/opentelco/go-dnc/models/pb/sharedpb"
	"git.liero.se/opentelco/go-dnc/models/pb/terminalpb"
	"git.liero.se/opentelco/go-dnc/models/pb/transportpb"
	"git.liero.se/opentelco/go-swpx/config"
	"git.liero.se/opentelco/go-swpx/proto/go/resourcepb"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createBasicSSHInterfaceTask(req *resourcepb.Request, conf *config.ResourceVRP) *transportpb.Message {

	sshConf := conf.Transports.GetByLabel("ssh")
	task := &terminalpb.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req.Timeout, defaultDeadlineTimeout))),
		Payload: []*terminalpb.Task_Payload{
			{
				Command: fmt.Sprintf("display mac-address %s", req.Port),
			},
		},
		Config: &terminalpb.Config{
			User:                sshConf.User,
			Password:            sshConf.Password,
			RegexPrompt:         sshConf.RegexPrompt,
			ScreenLengthCommand: sshConf.ScreenLength,
			ReadDeadLine:        durationpb.New(sshConf.ReadDeadLine.AsDuration()),
			WriteDeadLine:       durationpb.New(sshConf.WriteDeadLine.AsDuration()),
			SshKeyPath:          sshConf.SSHKeyPath,
		},
	}

	message := &transportpb.Message{
		Session: &transportpb.Session{
			NetworkRegion: req.NetworkRegion,
			Target:        req.Hostname,
			Port:          int32(sshConf.Port),
			Type:          transportpb.Type_SSH,
		},

		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transportpb.Task{
			Task: &transportpb.Task_Terminal{Terminal: task},
		},
		Status:   sharedpb.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req.Timeout, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}
	return message
}

func createCollectConfigMsg(req *resourcepb.GetRunningConfigParameters, conf *config.ResourceVRP) *transportpb.Message {

	sshConf := conf.Transports.GetByLabel("ssh")
	task := &terminalpb.Task{
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req.Timeout, defaultDeadlineTimeout))),
		Payload: []*terminalpb.Task_Payload{
			{
				Command: "disp cur",
			},
		},
		Config: &terminalpb.Config{
			User:                sshConf.User,
			Password:            sshConf.Password,
			RegexPrompt:         sshConf.RegexPrompt,
			ScreenLengthCommand: sshConf.ScreenLength,
			ReadDeadLine:        durationpb.New(sshConf.ReadDeadLine.AsDuration()),
			WriteDeadLine:       durationpb.New(sshConf.WriteDeadLine.AsDuration()),
			SshKeyPath:          sshConf.SSHKeyPath,
		},
	}

	message := &transportpb.Message{
		Session: &transportpb.Session{
			NetworkRegion: req.NetworkRegion,
			Target:        req.Hostname,
			Port:          int32(sshConf.Port),
			Type:          transportpb.Type_SSH,
		},

		Id:     ksuid.New().String(),
		Source: VERSION.String(),
		Task: &transportpb.Task{
			Task: &transportpb.Task_Terminal{Terminal: task},
		},
		Status:   sharedpb.Status_NEW,
		Deadline: timestamppb.New(time.Now().Add(validateEOLTimeout(req.Timeout, defaultDeadlineTimeout))),
		Created:  timestamppb.Now(),
	}
	return message
}
