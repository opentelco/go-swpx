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

	// delete a specific change by its ID
	DeleteChangeByID(ctx context.Context, id string) error

	// delete all changes for a specific device (used when deleting a device)
	DeleteChangersByDeviceID(ctx context.Context, id string) error
}