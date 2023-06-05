package device

import (
	"context"
	"encoding/json"
	"fmt"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	pollerDefaultProvider = "default"
	pollerDefaultResource = "generic"
)

var pj = protojson.MarshalOptions{
	Multiline:       false,
	AllowPartial:    false,
	UseProtoNames:   true,
	UseEnumNumbers:  false,
	EmitUnpopulated: true,
}

func New(repo Repository, logger hclog.Logger) devicepb.DeviceServiceServer {
	return &device{
		repo:   repo,
		logger: logger.Named("fleet"),
	}
}

type device struct {
	repo   Repository
	logger hclog.Logger

	devicepb.UnimplementedDeviceServiceServer
}

// *** Device ***
// Get a device by its ID, this is used to get a specific device
func (d *device) GetByID(ctx context.Context, params *devicepb.GetByIDParameters) (*devicepb.Device, error) {
	if params.Id == "" {
		return nil, ErrDeviceNotFound
	}
	return d.repo.GetByID(ctx, params.Id)
}

// Get a device by its hostname, managment ip or serial number etc (used to search for a device)
func (d *device) List(ctx context.Context, params *devicepb.ListParameters) (*devicepb.ListResponse, error) {
	devices, err := d.repo.List(ctx, params)
	return &devicepb.ListResponse{
		Devices: devices,
	}, err

}

// Create a device in the fleet
func (d *device) Create(ctx context.Context, params *devicepb.CreateParameters) (*devicepb.Device, error) {

	device := &devicepb.Device{}
	if params.Hostname == nil && params.ManagementIp == nil {
		return nil, ErrHostnameOrManagementIpRequired
	}

	if params.Hostname != nil {
		device.Hostname = *params.Hostname
	}

	if params.Domain != nil {
		device.Domain = *params.Domain
	}

	if params.ManagementIp != nil {
		device.ManagementIp = *params.ManagementIp
	}

	if params.SerialNumber != nil {
		device.SerialNumber = *params.SerialNumber
	}

	if params.Model != nil {
		device.Model = *params.Model
	}

	if params.Version != nil {
		device.Version = *params.Version
	}

	if params.PollerResourcePlugin != nil {
		device.PollerResourcePlugin = *params.PollerResourcePlugin
	} else {
		device.PollerResourcePlugin = pollerDefaultResource
	}

	if params.PollerProviderPlugin != nil {
		device.PollerProviderPlugin = *params.PollerProviderPlugin
	} else {
		device.PollerProviderPlugin = pollerDefaultProvider
	}

	if params.LastReboot != nil {
		device.LastReboot = params.LastReboot
	}
	if params.LastSeen != nil {
		device.LastSeen = params.LastSeen
	}

	if params.NetworkRegion != nil {
		device.NetworkRegion = *params.NetworkRegion
	}

	return d.repo.Upsert(ctx, device)
}

// Update a device in the fleet (this is used to update the device with new information)
func (d *device) Update(ctx context.Context, params *devicepb.UpdateParameters) (*devicepb.Device, error) {
	if params.Id == "" {
		return nil, ErrDeviceNotFound
	}
	deviceA, err := d.repo.GetByID(ctx, params.Id)
	if err != nil {
		return nil, err
	}

	x := proto.Clone(deviceA)
	var deviceToUpdate = x.(*devicepb.Device)

	if params.Hostname != nil {
		deviceToUpdate.Hostname = *params.Hostname
	}
	if params.Domain != nil {
		deviceToUpdate.Domain = *params.Domain
	}
	if params.ManagementIp != nil {
		deviceToUpdate.ManagementIp = *params.ManagementIp
	}

	if params.SerialNumber != nil {
		deviceToUpdate.SerialNumber = *params.SerialNumber
	}
	if params.Model != nil {
		deviceToUpdate.Model = *params.Model
	}
	if params.Version != nil {
		deviceToUpdate.Version = *params.Version
	}
	if params.PollerResourcePlugin != nil {
		deviceToUpdate.PollerResourcePlugin = *params.PollerResourcePlugin
	}
	if params.PollerProviderPlugin != nil {
		deviceToUpdate.PollerProviderPlugin = *params.PollerProviderPlugin
	}
	if params.LastReboot != nil {
		deviceToUpdate.LastReboot = params.LastReboot
	}
	if params.LastSeen != nil {
		deviceToUpdate.LastSeen = params.LastSeen
	}

	if params.NetworkRegion != nil {
		deviceToUpdate.NetworkRegion = *params.NetworkRegion
	}

	changes := getChanges(deviceA, deviceToUpdate)
	if len(changes) > 0 {

		deviceToUpdate, err = d.repo.Upsert(ctx, deviceToUpdate)

		if err != nil {
			return nil, err
		}
		for _, change := range changes {
			_, err = d.repo.UpsertChange(ctx, change)
			if err != nil {
				return nil, err
			}
		}
		return deviceToUpdate, nil
	} else {
		return deviceA, nil
	}

}

// Delete a device from the fleet (mark the device as deleted)
func (d *device) Delete(ctx context.Context, params *devicepb.DeleteParameters) (*emptypb.Empty, error) {
	err := d.repo.Delete(ctx, params.Id)
	if err != nil {
		return nil, err
	}

	err = d.repo.DeleteChangersByDeviceID(ctx, params.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (d *device) ListChanges(ctx context.Context, params *devicepb.ListChangesParameters) (*devicepb.ListChangesResponse, error) {
	c, err := d.repo.ListChanges(ctx, params)
	if err != nil {
		return nil, err
	}
	return &devicepb.ListChangesResponse{
		Changes: c,
	}, nil

}
func (d *device) GetChangeByID(ctx context.Context, params *devicepb.GetChangeByIDParameters) (*devicepb.Change, error) {
	c, err := d.repo.GetChangeByID(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// compare chagnes between two devices and return a list of changes
func getChanges(a, b *devicepb.Device) []*devicepb.Change {
	skipFields := []string{"id", "created", "updated"}
	amap := protoToMap(a)
	bmap := protoToMap(b)
	changes := make([]*devicepb.Change, 0)
	for k, v := range amap {
		if inStringArray(k, skipFields) {
			continue
		}
		if v != bmap[k] {
			changes = append(changes, &devicepb.Change{
				DeviceId: a.Id,
				Field:    k,
				OldValue: toString(v),
				NewValue: toString(bmap[k]),
				Created:  timestamppb.Now(),
			})
		}
	}
	return changes
}

// generic any type to string
func toString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

// convert from proto message to map[string]interface{} by marshalling to json and then unmarshalling to map
func protoToMap(m proto.Message) map[string]interface{} {
	var result map[string]interface{}
	b, _ := pj.Marshal(m)
	json.Unmarshal(b, &result)
	return result
}

// helper function to check if a string is in a string array
func inStringArray(s string, a []string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}
	return false
}

func (d *device) AddEvent(ctx context.Context, event *devicepb.Event) (*devicepb.Event, error) {
	if event.DeviceId == "" {
		return nil, ErrDeviceNotFound
	}

	return d.repo.AddEvent(ctx, event)
}

func (d *device) GetEventByID(ctx context.Context, params *devicepb.GetEventByIDParameters) (*devicepb.Event, error) {
	return d.repo.GetEventByID(ctx, params.Id)
}

// returns a list of events (default 100)
func (d *device) ListEvents(ctx context.Context, params *devicepb.ListEventsParameters) (*devicepb.ListEventsResponse, error) {

	events, err := d.repo.ListEvents(ctx, params)
	if err != nil {
		return nil, err
	}
	return &devicepb.ListEventsResponse{
		Events: events,
	}, nil
}
