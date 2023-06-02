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

	if params.PollerProvider == "" {
		params.PollerProvider = "default_provider"
	}

	if params.Hostname == "" {
		return nil, ErrHostnameRequired
	}

	device := &devicepb.Device{
		Hostname:             params.Hostname,
		Domain:               params.Domain,
		ManagementIp:         params.ManagementIp,
		SerialNumber:         params.SerialNumber,
		Model:                params.Model,
		Version:              params.Version,
		PollerResourcePlugin: params.PollerResourcePlugin,
		PollerProvider:       params.PollerProvider,
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

	deviceB := &devicepb.Device{
		Id: params.Id,
	}
	if params.Hostname != nil {
		deviceB.Hostname = *params.Hostname
	}
	if params.Domain != nil {
		deviceB.Domain = *params.Domain
	}
	if params.ManagementIp != nil {
		deviceB.ManagementIp = *params.ManagementIp
	}

	if params.SerialNumber != nil {
		deviceB.SerialNumber = *params.SerialNumber
	}
	if params.Model != nil {
		deviceB.Model = *params.Model
	}
	if params.Version != nil {
		deviceB.Version = *params.Version
	}
	if params.PollerResourcePlugin != nil {
		deviceB.PollerResourcePlugin = *params.PollerResourcePlugin
	}
	if params.PollerProvider != nil {
		deviceB.PollerProvider = *params.PollerProvider
	}

	changes := getChanges(deviceA, deviceB)
	if len(changes) > 0 {
		deviceB, err = d.repo.Upsert(ctx, deviceB)
		if err != nil {
			return nil, err
		}
		for _, change := range changes {
			_, err = d.repo.UpsertChange(ctx, change)
			if err != nil {
				return nil, err
			}
		}
		return deviceB, nil
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
	if params.Limit == 0 {
		params.Limit = 100
	}
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
