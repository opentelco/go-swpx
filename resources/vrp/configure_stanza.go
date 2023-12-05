package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"git.liero.se/opentelco/go-dnc/client"
	"git.liero.se/opentelco/go-dnc/models/pb/terminalpb"
	"git.liero.se/opentelco/go-dnc/models/pb/transportpb"
	"git.liero.se/opentelco/go-swpx/proto/go/resourcepb"
	"git.liero.se/opentelco/go-swpx/proto/go/stanzapb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	commandEnableConfigure = "system-view"
	commandEndConfigure    = "quit"
)

func (d *VRPDriver) ConfigureStanza(ctx context.Context, req *resourcepb.ConfigureStanzaRequest) (*stanzapb.ConfigureResponse, error) {
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
	stanza := []*terminalpb.Task_Payload{{Command: commandEnableConfigure}}

	for _, s := range req.Stanza {
		stanza = append(stanza, &terminalpb.Task_Payload{
			Command:      s.Content,
			PromptErrors: []string{"Error: Unrecognized command found"},
		})
	}

	// we start at 0 so the last is linue number is the length of the stanza
	stanza = append(stanza, &terminalpb.Task_Payload{
		Command: commandEndConfigure,
	})

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
	response := &stanzapb.ConfigureResponse{
		StanzaResult: []*stanzapb.Result{},
	}

	for _, p := range pl {
		lineResult := &stanzapb.Result{
			Line:   p.Command,
			Status: stanzapb.Result_SUCCESS,
		}

		if p.Error != "" {
			lineResult.Status = stanzapb.Result_FAILED
			lineResult.Error = &p.Error
		}
		response.StanzaResult = append(response.StanzaResult, lineResult)
	}

	if res.Error != "" {
		d.logger.Error("error in stanza", "error", res.Error)
		return response, fmt.Errorf("error in stanza: %s", res.Error)
	}

	return response, nil
}

func prettyPrintJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
