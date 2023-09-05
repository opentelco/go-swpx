package mappers

import (
	"git.liero.se/opentelco/go-swpx/fleet/graph/model"
	"git.liero.se/opentelco/go-swpx/fleet/internal"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
)

type Devices []*devicepb.Device

type Device struct {
	*devicepb.Device
}

func (d Device) ToGQL() *model.Device {

	return &model.Device{
		ID:                   d.Id,
		Hostname:             &d.Hostname,
		Domain:               &d.Domain,
		ManagementIP:         &d.ManagementIp,
		SerialNumber:         &d.SerialNumber,
		Model:                &d.Model,
		Version:              &d.Version,
		NetworkRegion:        &d.NetworkRegion,
		PollerResourcePlugin: &d.PollerProviderPlugin,
		PollerProviderPlugin: &d.PollerProviderPlugin,
		State:                DeviceState(d.State).ToGQL(),
		Status:               DeviceStatus(d.Status).ToGQL(),
		Events:               &model.EventConnection{},
		Changes:              &model.ChangeConnection{},
		LastSeen:             internal.TimestampPBToTimePointer(d.LastSeen),
		CreatedAt:            d.Created.AsTime(),
		UpdatedAt:            d.Updated.AsTime(),
		LastReboot:           internal.TimestampPBToTimePointer(d.LastReboot),
	}
}

func (d Devices) ToGQL() []*model.Device {
	res := make([]*model.Device, len(d))
	for i, v := range d {
		res[i] = Device{v}.ToGQL()
	}
	return res
}

type DeviceListParameters struct{ *devicepb.ListParameters }

func (p DeviceListParameters) ToGQL() *model.ListDevicesParams {
	return &model.ListDevicesParams{
		Search:       p.Search,
		Hostname:     p.Hostname,
		ManagementIP: p.ManagementIp,
		Limit:        internal.PointerInt64ToPointerInt(p.Limit),
		Offset:       internal.PointerInt64ToPointerInt(p.Offset),
	}
}

func ToDeviceListResponse(r *devicepb.ListResponse) *DeviceListResponse {
	return &DeviceListResponse{r}
}

func ToDevice(r *devicepb.Device) *Device {
	return &Device{r}
}

type DeviceListResponse struct{ *devicepb.ListResponse }

func (r DeviceListResponse) ToGQL() *model.ListDeviceResponse {
	return &model.ListDeviceResponse{
		Devices:  Devices(r.Devices).ToGQL(),
		PageInfo: ToPageInfo(r.PageInfo).ToGQL(),
	}
}

type DeviceStatus devicepb.Device_Status

func (s DeviceStatus) ToGQL() model.DeviceStatus {
	switch devicepb.Device_Status(s) {
	case devicepb.Device_DEVICE_STATUS_NEW:
		return model.DeviceStatusNew
	case devicepb.Device_DEVICE_STATUS_REACHABLE:
		return model.DeviceStatusReachable
	case devicepb.Device_DEVICE_STATUS_UNREACHABLE:
		return model.DeviceStatusUnreachable
	default:
		return model.DeviceStatusNotSet
	}

}

type DeviceState devicepb.Device_State

func (s DeviceState) ToGQL() model.DeviceState {
	switch devicepb.Device_State(s) {
	case devicepb.Device_DEVICE_STATE_NEW:
		return model.DeviceStateNew
	case devicepb.Device_DEVICE_STATE_ACTIVE:
		return model.DeviceStateActive
	case devicepb.Device_DEVICE_STATE_INACTIVE:
		return model.DeviceStateInactive
	case devicepb.Device_DEVICE_STATE_ROUGE:
		return model.DeviceStateRouge
	default:
		return model.DeviceStateNotSet
	}

}

type DeviceEvents []*devicepb.Event
type DeviceEvent struct {
	*devicepb.Event
}

type DeviceEventType devicepb.Event_Type

func (t DeviceEventType) ToGQL() model.DeviceEventType {
	switch devicepb.Event_Type(t) {
	case devicepb.Event_TYPE_NOT_SET:
		return model.DeviceEventTypeNotSet
	case devicepb.Event_DEVICE:
		return model.DeviceEventTypeDevice
	case devicepb.Event_CONFIGURATION:
		return model.DeviceEventTypeConfiguration
	default:
		return model.DeviceEventTypeNotSet
	}
}

type DeviceEventAction devicepb.Event_Action

func (a DeviceEventAction) ToGQL() model.DeviceEventAction {
	switch devicepb.Event_Action(a) {
	case devicepb.Event_ACTION_NOT_SET:
		return model.DeviceEventActionNotSet
	case devicepb.Event_CREATE:
		return model.DeviceEventActionCreate
	case devicepb.Event_UPDATE:
		return model.DeviceEventActionUpdate
	case devicepb.Event_COLLECT_CONFIG:
		return model.DeviceEventActionCollectConfig
	case devicepb.Event_COLLECT_DEVICE:
		return model.DeviceEventActionCollectDevice
	default:
		return model.DeviceEventActionNotSet
	}
}

type DeviceEventOutcome devicepb.Event_Outcome

func (o DeviceEventOutcome) ToGQL() model.DeviceEventOutcome {
	switch devicepb.Event_Outcome(o) {
	case devicepb.Event_OUTCOME_NOT_SET:
		return model.DeviceEventOutcomeNotSet
	case devicepb.Event_SUCCESS:
		return model.DeviceEventOutcomeSuccess
	case devicepb.Event_FAILURE:
		return model.DeviceEventOutcomeFailure
	default:
		return model.DeviceEventOutcomeNotSet
	}
}

func (e DeviceEvent) ToGQL() *model.DeviceEvent {
	return &model.DeviceEvent{
		ID:        e.Id,
		DeviceID:  e.DeviceId,
		Type:      DeviceEventType(e.Type).ToGQL(),
		Message:   e.Message,
		Action:    DeviceEventAction(e.Action).ToGQL(),
		Outcome:   DeviceEventOutcome(e.Outcome).ToGQL(),
		CreatedAt: e.Created.AsTime(),
	}
}

func (e DeviceEvents) ToGQL() []*model.DeviceEvent {
	res := make([]*model.DeviceEvent, len(e))
	for i, v := range e {
		res[i] = DeviceEvent{v}.ToGQL()
	}
	return res
}

type DeviceChanges []*devicepb.Change
type DeviceChange struct {
	*devicepb.Change
}

func (c DeviceChange) ToGQL() *model.DeviceChange {
	return &model.DeviceChange{
		ID:        c.Id,
		Field:     c.Field,
		DeviceID:  c.DeviceId,
		OldValue:  c.OldValue,
		NewValue:  c.NewValue,
		CreatedAt: c.Created.AsTime(),
	}
}

func (c DeviceChanges) ToGQL() []*model.DeviceChange {
	res := make([]*model.DeviceChange, len(c))
	for i, v := range c {
		res[i] = DeviceChange{v}.ToGQL()
	}
	return res
}

type ListDeviceChangesResponse struct{ *devicepb.ListChangesResponse }

func (r ListDeviceChangesResponse) ToGQL() *model.ListDeviceChangesResponse {
	return &model.ListDeviceChangesResponse{
		Changes:  DeviceChanges(r.Changes).ToGQL(),
		PageInfo: ToPageInfo(r.PageInfo).ToGQL(),
	}
}

func ToListDeviceChangesResponse(r *devicepb.ListChangesResponse) *ListDeviceChangesResponse {
	return &ListDeviceChangesResponse{r}
}

type ListDeviceEventsResponse struct{ *devicepb.ListEventsResponse }

func (r ListDeviceEventsResponse) ToGQL() *model.ListDeviceEventsResponse {
	return &model.ListDeviceEventsResponse{
		Events:   DeviceEvents(r.Events).ToGQL(),
		PageInfo: ToPageInfo(r.PageInfo).ToGQL(),
	}
}

func ToListDeviceEventsResponse(r *devicepb.ListEventsResponse) *ListDeviceEventsResponse {
	return &ListDeviceEventsResponse{r}
}
