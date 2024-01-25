package main

import (
	"context"
	"strings"

	"go.opentelco.io/go-dnc/models/pb/transportpb"
	"go.opentelco.io/go-swpx/proto/go/devicepb"
	"go.opentelco.io/go-swpx/proto/go/resourcepb"
	"go.opentelco.io/go-swpx/shared/oids"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d *driver) Version() (string, error) {
	return VERSION.String(), nil
}

func (d *driver) Discover(ctx context.Context, req *resourcepb.Request) (*devicepb.Device, error) {
	errs := make([]*devicepb.TransientError, 0)
	var msgs []*transportpb.Message
	msgs = append(msgs, createTaskSystemInfo(req, &d.conf.Snmp))
	ne := &devicepb.Device{}
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
	ne.TransientErrors = &devicepb.TransientErrors{Errors: errs}
	return ne, nil
}

func (d *driver) TechnicalPortInformation(ctx context.Context, req *resourcepb.Request) (*devicepb.Device, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) BasicPortInformation(ctx context.Context, req *resourcepb.Request) (*devicepb.Device, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) AllPortInformation(ctx context.Context, req *resourcepb.Request) (*devicepb.Device, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) MapInterface(ctx context.Context, req *resourcepb.Request) (*resourcepb.PortIndex, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) MapEntityPhysical(ctx context.Context, req *resourcepb.Request) (*resourcepb.PortIndex, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) GetTransceiverInformation(ctx context.Context, req *resourcepb.Request) (*devicepb.Transceiver, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) GetAllTransceiverInformation(ctx context.Context, req *resourcepb.Request) (*devicepb.Transceivers, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) GetRunningConfig(ctx context.Context, req *resourcepb.GetRunningConfigParameters) (*resourcepb.GetRunningConfigResponse, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}

func (d *driver) GetDeviceInformation(ctx context.Context, req *resourcepb.Request) (*devicepb.Device, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}
