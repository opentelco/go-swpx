package device

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
)

type Repository interface {
	// Device
	GetByID(ctx context.Context, id string) (*devicepb.Device, error)
	List(ctx context.Context, parms *devicepb.ListParameters) ([]*devicepb.Device, error)
	Upsert(ctx context.Context, dev *devicepb.Device) (*devicepb.Device, error)
	Delete(ctx context.Context, id string) error

	UpsertChange(ctx context.Context, change *devicepb.Change) (*devicepb.Change, error)

	GetChangeByID(ctx context.Context, id string) (*devicepb.Change, error)
	ListChanges(ctx context.Context, params *devicepb.ListChangesParameters) ([]*devicepb.Change, error)

	AddEvent(ctx context.Context, event *devicepb.Event) (*devicepb.Event, error)
	GetEventByID(ctx context.Context, id string) (*devicepb.Event, error)
	ListEvents(ctx context.Context, params *devicepb.ListEventsParameters) ([]*devicepb.Event, error)

	// delete a specific change by its ID
	DeleteChangeByID(ctx context.Context, id string) error

	// delete all changes for a specific device (used when deleting a device)
	DeleteChangersByDeviceID(ctx context.Context, id string) error

	// Upsert a schedule on a device
	UpsertSchedule(ctx context.Context, deviceId string, schedule *devicepb.Device_Schedule) (*devicepb.Device, error)
}
