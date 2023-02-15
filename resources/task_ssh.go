package resources

import (
	"fmt"

	shared2 "git.liero.se/opentelco/go-dnc/models/pb/shared"
	"git.liero.se/opentelco/go-dnc/models/pb/terminal"
	"git.liero.se/opentelco/go-dnc/models/pb/transport"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateBasicSshInterfaceTask(el *proto.NetworkElement, conf *shared.Configuration) *transport.Message {
	task := &terminal.Task{
		Deadline: el.Conf.Request.Deadline,
		Payload: []*terminal.Task_Payload{
			{
				Command: fmt.Sprintf("display mac-address %s", el.Interface),
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
			Port:   int32(el.Conf.Ssh.Port),
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
