package main

import (
	"context"
	"strings"

	"git.liero.se/opentelco/go-dnc/models/pb/transportpb"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelementpb"
	"git.liero.se/opentelco/go-swpx/proto/go/resourcepb"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d *driver) Version() (string, error) {
	return VERSION.String(), nil
}

func (d *driver) Discover(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error) {
	errs := make([]*networkelementpb.TransientError, 0)
	var msgs []*transportpb.Message
	msgs = append(msgs, createTaskSystemInfo(req, &d.conf.Snmp))
	ne := &networkelementpb.Element{}
	ne.Hostname = req.Hostname

	for _, msg := range msgs {
		reply, err := d.dnc.Put(ctx, msg)
		if err != nil {
			d.logger.Error("could not complete discovery", "error", err.Error())
			return nil, err
		}

		switch task := reply.Task.Task.(type) {
		case *transportpb.Task_Snmpc:
			d.logger.Debug("the reply returns from dnc",
				"status", reply.Status.String(),
				"completed", reply.Completed.String(),
				"execution_time", reply.ExecutionTime.AsDuration().String(),
				"size", len(task.Snmpc.Metrics),
			)

			for _, m := range task.Snmpc.Metrics {
				switch {
				case strings.HasPrefix(m.Oid, oids.SysPrefix):
					parseSystemInformation(m, ne)
				}
			}
		}
	}
	ne.TransientErrors = &networkelementpb.TransientErrors{Errors: errs}
	return ne, nil
}

func (d *driver) TechnicalPortInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) BasicPortInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) AllPortInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Element, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) MapInterface(ctx context.Context, req *resourcepb.Request) (*resourcepb.PortIndex, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) MapEntityPhysical(ctx context.Context, req *resourcepb.Request) (*resourcepb.PortIndex, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) GetTransceiverInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Transceiver, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) GetAllTransceiverInformation(ctx context.Context, req *resourcepb.Request) (*networkelementpb.Transceivers, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) GetRunningConfig(ctx context.Context, req *resourcepb.GetRunningConfigParameters) (*resourcepb.GetRunningConfigResponse, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}
