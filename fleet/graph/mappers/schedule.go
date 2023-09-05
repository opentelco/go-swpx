package mappers

import (
	"git.liero.se/opentelco/go-swpx/fleet/graph/model"
	"git.liero.se/opentelco/go-swpx/fleet/internal"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/schedulepb"
)

type ScheduleType schedulepb.Schedule_Type

func (t ScheduleType) ToGQL() model.ScheduleType {
	switch schedulepb.Schedule_Type(t) {
	case schedulepb.Schedule_COLLECT_CONFIG:
		return model.ScheduleTypeConfig
	case schedulepb.Schedule_COLLECT_DEVICE:
		return model.ScheduleTypeDevice
	default:
		return model.ScheduleTypeNotSet
	}

}

type Schedule struct{ *schedulepb.Schedule }

func (s Schedule) ToGQL() *model.DeviceSchedule {
	if s.Schedule == nil {
		return nil
	}

	return &model.DeviceSchedule{
		Interval:    s.Interval.AsDuration(),
		Type:        ScheduleType(s.Type).ToGQL(),
		LastRun:     internal.TimestampPBToTimePointer(s.LastRun),
		Active:      s.Active,
		FailedCount: int(s.FailedCount),
	}
}
