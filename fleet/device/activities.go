package device

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"go.temporal.io/sdk/temporal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Activities struct {
	device devicepb.DeviceServiceServer
}

func NewActivities(device devicepb.DeviceServiceServer) *Activities {
	return &Activities{
		device: device,
	}
}

func (a *Activities) GetByID(ctx context.Context, deviceId string) (*devicepb.Device, error) {
	return a.device.GetByID(ctx, &devicepb.GetByIDParameters{Id: deviceId})
}

func (a *Activities) CreateDevice(ctx context.Context, params *devicepb.CreateParameters) (*devicepb.Device, error) {
	return a.device.Create(ctx, params)
}

func (a *Activities) AddDeviceEvent(ctx context.Context, event *devicepb.Event) (*devicepb.Event, error) {
	return a.device.AddEvent(ctx, event)
}

func (a *Activities) GetDeviceByID(ctx context.Context, id string) (*devicepb.Device, error) {
	dev, err := a.device.GetByID(ctx, &devicepb.GetByIDParameters{Id: id})
	if errors.Is(err, ErrDeviceNotFound) {
		return nil, temporal.NewNonRetryableApplicationError("device not found", ErrTypeDeviceNotFound, err)
	}
	return dev, err
}

func (a *Activities) SetDeviceUnreachable(ctx context.Context, id string) (*devicepb.Device, error) {
	dev, err := a.device.Update(ctx, &devicepb.UpdateParameters{
		Id:     id,
		Status: &[]devicepb.Device_Status{devicepb.Device_DEVICE_STATUS_UNREACHABLE}[0],
	})
	if errors.Is(err, ErrDeviceNotFound) {
		return nil, temporal.NewNonRetryableApplicationError("device not found", ErrTypeDeviceNotFound, err)
	}
	return dev, err
}

func (a *Activities) UpdateDevice(ctx context.Context, params *devicepb.UpdateParameters) (*devicepb.Device, error) {
	dev, err := a.device.Update(ctx, params)
	if errors.Is(err, ErrDeviceNotFound) {
		return nil, temporal.NewNonRetryableApplicationError("device not found", ErrTypeDeviceNotFound, err)
	}
	return dev, err
}

func (a *Activities) ListDevices(ctx context.Context, params *devicepb.ListParameters) (*devicepb.ListResponse, error) {
	return a.device.List(ctx, params)
}

func (a *Activities) SetScheduleLastRun(ctx context.Context, id string, scheduleType devicepb.Device_Schedule_Type, lastRun time.Time) (*devicepb.Device, error) {
	d, err := a.device.GetByID(ctx, &devicepb.GetByIDParameters{Id: id})
	if err != nil {
		return nil, err
	}

	schedule := getDeviceScheduleByType(d.Schedules, scheduleType)
	if schedule == nil {
		return nil, fmt.Errorf("schedule %s not found", scheduleType.String())
	}
	schedule.LastRun = timestamppb.New(lastRun)
	params := &devicepb.SetScheduleParameters{
		DeviceId: id,
		Schedule: schedule,
	}

	return a.device.SetSchedule(ctx, params)
}
