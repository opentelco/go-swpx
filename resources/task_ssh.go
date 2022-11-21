package resources

import (
	"fmt"

	shared2 "git.liero.se/opentelco/go-dnc/models/pb/shared"
	"git.liero.se/opentelco/go-dnc/models/pb/ssh"
	"git.liero.se/opentelco/go-dnc/models/pb/transport"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared"
	"github.com/segmentio/ksuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateBasicSshInterfaceTask(el *proto.NetworkElement, conf *shared.Configuration) *transport.Message {
	task := &ssh.Task{
		Type: ssh.Type_GET,
		Payload: []*ssh.Payload{
			{
				Command: fmt.Sprintf("display mac-address %s", el.Interface),
			},
		},
		Config: &ssh.Config{
			User:                conf.Ssh.Username,
			Password:            conf.Ssh.Password,
			Port:                conf.Ssh.Port,
			ScreenLength:        conf.Ssh.ScreenLength,
			ScreenLengthCommand: conf.Ssh.ScreenLengthCommand,
			RegexPrompt:         conf.Ssh.RegexPrompt,
			Ttl:                 &durationpb.Duration{Seconds: int64(conf.Ssh.TTL.Seconds())},
			ReadDeadLine:        &durationpb.Duration{Seconds: int64(conf.Ssh.ReadDeadLine.Seconds())},
			WriteDeadLine:       &durationpb.Duration{Seconds: int64(conf.Ssh.WriteDeadLine.Seconds())},
		},
		Host: el.Hostname,
	}

	message := &transport.Message{
		Id:      ksuid.New().String(),
		Target:  el.Hostname,
		Type:    transport.Type_SSH,
		Task:    &transport.Message_Ssh{Ssh: task},
		Status:  shared2.Status_NEW,
		Created: timestamppb.Now(),
	}
	return message
}
