package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"git.liero.se/opentelco/go-dnc/client"
	"git.liero.se/opentelco/go-dnc/models/pb/terminalpb"
	"git.liero.se/opentelco/go-dnc/models/pb/transportpb"
	"git.liero.se/opentelco/go-swpx/proto/go/resourcepb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	commandEnableConfigure = "system-view"
	commandEndConfigure    = "quit"
)

func (d *VRPDriver) ConfigureStanza(ctx context.Context, req *resourcepb.ConfigureStanzaRequest) (*resourcepb.ConfigureStanzaResponse, error) {
	deadline := time.Now().Add(validateEOLTimeout(req.Timeout, defaultDeadlineTimeout))
	sshConf := d.conf.Transports.GetByLabel("ssh")
	msg, err := client.NewMessage(client.NewMessageParameters{
		Target:        req.Hostname,
		Port:          int32(sshConf.Port),
		Type:          transportpb.Type_SSH,
		NetworkRegion: req.NetworkRegion,
		Source:        VERSION.String(),
		Deadline:      deadline,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create dnc message: %w", err)
	}

	// setup the commands, first enable configure mode, then the stanza, then exit configure mode
	stanza := []*terminalpb.Task_Payload{{Command: commandEnableConfigure, LineNumber: 0}}

	for n, s := range req.Stanza {
		stanza = append(stanza, &terminalpb.Task_Payload{
			Command:      s,
			LineNumber:   int64(n + 1),
			PromptErrors: []string{"Error: Unrecognized command found"},
		})
	}

	// we start at 0 so the last is linue number is the length of the stanza
	stanza = append(stanza, &terminalpb.Task_Payload{
		Command:    commandEndConfigure,
		LineNumber: int64(len(req.Stanza))})

	task := &terminalpb.Task{
		Payload:  stanza,
		Deadline: timestamppb.New(deadline),
		Config: &terminalpb.Config{
			RegexPrompt:         sshConf.RegexPrompt,
			ScreenLengthCommand: sshConf.ScreenLength,
			ReadDeadLine:        durationpb.New(sshConf.ReadDeadLine.AsDuration()),
			WriteDeadLine:       durationpb.New(sshConf.WriteDeadLine.AsDuration()),
			SshKeyPath:          sshConf.SSHKeyPath,
		},
	}

	msg.Task = &transportpb.Task{
		Task: &transportpb.Task_Terminal{Terminal: task},
	}

	res, err := d.dnc.Put(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("could not complete dnc request: %w", err)
	}

	pl := res.Task.GetTerminal().Payload
	// sort pl by line number
	sort.Slice(pl, func(i, j int) bool {
		return pl[i].LineNumber < pl[j].LineNumber
	})
	for _, p := range pl {

		if p.Error != "" {
			d.logger.Error("error in stanza", "error", p.Error)
		}
	}

	if res.Error != "" {
		d.logger.Error("error in stanza", "error", res.Error)
		return nil, fmt.Errorf("error in stanza: %s", res.Error)
	}

	return &resourcepb.ConfigureStanzaResponse{}, nil
}
