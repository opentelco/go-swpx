package repo

import (
	"context"
	"time"

	"git.liero.se/opentelco/go-swpx/database"
	"git.liero.se/opentelco/go-swpx/fleet/device"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultDatabase          string = "swpx"
	defaultDeviceCollection  string = "devices"
	defaultChangesCollection string = "device_changes"
)

type repo struct {
	mc *mongo.Client

	deviceCollection  *mongo.Collection
	changesCollection *mongo.Collection

	logger hclog.Logger
}

func New(mc *mongo.Client, database string, logger hclog.Logger) (device.Repository, error) {
	if database == "" {
		database = defaultDatabase
	}

	return &repo{
		mc:                mc,
		deviceCollection:  mc.Database(database).Collection(defaultDeviceCollection),
		changesCollection: mc.Database(database).Collection(defaultChangesCollection),
		logger:            logger.Named("db"),
	}, nil
}

func (r *repo) GetByID(ctx context.Context, id string) (*devicepb.Device, error) {
	// get a device by its ID, this is used to get a specific device
	filter := bson.M{"_id": id}
	var dv devicepb.Device
	err := r.deviceCollection.FindOne(ctx, filter).Decode(&dv)
	if err != nil {
		// if errors is mongo.ErrNoDocuments then return nil, nil

		return nil, device.ErrDeviceNotFound
	}

	return &dv, nil
}

func (r *repo) List(ctx context.Context, params *devicepb.ListParameters) ([]*devicepb.Device, error) {
	filter := bson.M{}
	if params.Hostname != "" {
		filter["hostname"] = params.Hostname
	}
	if params.ManagementIp != "" {
		filter["management_ip"] = params.ManagementIp
	}

	if params.Search != "" {
		filter["$or"] = bson.A{
			bson.M{"hostname": bson.M{"$regex": params.Search}},
			bson.M{"management_ip": bson.M{"$regex": params.Search}},
		}
	}

	var devices []*devicepb.Device
	cur, err := r.deviceCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var device devicepb.Device
		err := cur.Decode(&device)
		if err != nil {
			return nil, err
		}
		devices = append(devices, &device)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return devices, nil

}

// Upsert updates or inserts a device into the database
func (r *repo) Upsert(ctx context.Context, dev *devicepb.Device) (*devicepb.Device, error) {

	if dev.Id == "" { // if the device does not have an ID, then generate one and set the created time
		dev.Id = database.NewID()
		dev.Created = timestamppb.New(time.Now())
	}
	dev.Updated = timestamppb.New(time.Now())

	// upsert the device into the device collection and return the device
	_, err := r.deviceCollection.UpdateOne(ctx, bson.M{"_id": dev.Id}, bson.M{"$set": dev},
		&options.UpdateOptions{Upsert: &[]bool{true}[0]})
	if err != nil {
		return nil, err
	}

	return dev, nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	// delete a device from the device collection
	_, err := r.deviceCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) UpsertChange(ctx context.Context, change *devicepb.Change) (*devicepb.Change, error) {
	// add a device change to the changes collection
	if change.Id == "" {
		change.Id = database.NewID()
		change.Created = timestamppb.New(time.Now())
	}

	_, err := r.changesCollection.InsertOne(ctx, change)
	if err != nil {
		return nil, err
	}

	return r.GetChangeByID(ctx, change.Id)
}

func (r *repo) GetChangeByID(ctx context.Context, id string) (*devicepb.Change, error) {
	// get a device change by its ID
	filter := bson.M{"_id": id}
	var change devicepb.Change
	err := r.changesCollection.FindOne(ctx, filter).Decode(&change)
	if err != nil {
		return nil, err
	}

	return &change, nil
}

func (r *repo) ListChanges(ctx context.Context, params *devicepb.ListChangesParameters) ([]*devicepb.Change, error) {
	// list changes for a device
	filter := bson.M{}
	if params.DeviceId != "" {
		filter["device_id"] = params.DeviceId
	}
	if params.Limit > 0 {
		filter["limit"] = params.Limit
	}
	if params.Offset > 0 {
		filter["offset"] = params.Offset
	}

	var changes []*devicepb.Change
	cur, err := r.changesCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {

		var change devicepb.Change
		err := cur.Decode(&change)
		if err != nil {
			return nil, err
		}
		changes = append(changes, &change)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return changes, nil
}

func (r *repo) DeleteChangeByID(ctx context.Context, id string) error {
	// delete a device change by its ID
	_, err := r.changesCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
func (r *repo) DeleteChangersByDeviceID(ctx context.Context, id string) error {
	// delete all device changes for a device
	_, err := r.changesCollection.DeleteMany(ctx, bson.M{"device_id": id})
	if err != nil {
		return err
	}

	return nil
}
