package repo

import (
	"context"
	"time"

	"git.liero.se/opentelco/go-swpx/fleet"
	"git.liero.se/opentelco/go-swpx/proto/go/fleetpb"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultDatabase                      string = "swpx"
	defaultDeviceCollection              string = "devices"
	defaultDeviceConfigurationCollection string = "device_configurations"
)

type repo struct {
	mc *mongo.Client

	deviceCollection       *mongo.Collection
	deviceConfigCollection *mongo.Collection

	logger hclog.Logger
}

func New(mc *mongo.Client, database string, logger hclog.Logger) (fleet.Repository, error) {
	if database == "" {
		database = defaultDatabase
	}

	return &repo{
		mc:                     mc,
		deviceCollection:       mc.Database(database).Collection(defaultDeviceCollection),
		deviceConfigCollection: mc.Database(database).Collection(defaultDeviceConfigurationCollection),
		logger:                 logger.Named("db"),
	}, nil
}

func (r *repo) GetDeviceByID(ctx context.Context, id string) (*fleetpb.Device, error) {
	// get a device by its ID, this is used to get a specific device
	filter := bson.M{"_id": id}
	var device fleetpb.Device
	err := r.deviceCollection.FindOne(ctx, filter).Decode(&device)
	if err != nil {
		// if errors is mongo.ErrNoDocuments then return nil, nil

		return nil, fleet.ErrDeviceNotFound
	}

	return &device, nil
}

func (r *repo) ListDevices(ctx context.Context, params *fleetpb.ListDevicesParameters) ([]*fleetpb.Device, error) {
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

	var devices []*fleetpb.Device
	cur, err := r.deviceCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var device fleetpb.Device
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

// UpsertDevice updates or inserts a device into the database
func (r *repo) UpsertDevice(ctx context.Context, dev *fleetpb.Device) (*fleetpb.Device, error) {

	if dev.Id == "" { // if the device does not have an ID, then generate one and set the created time
		dev.Id = fleet.NewID()
		dev.Created = timestamppb.New(time.Now())
	}
	dev.Updated = timestamppb.New(time.Now())

	// upsert the device into the device collection and return the device
	_, err := r.deviceCollection.UpdateOne(ctx, bson.M{"_id": dev.Id}, bson.M{"$set": dev}, &options.UpdateOptions{Upsert: &[]bool{true}[0]})
	if err != nil {
		return nil, err
	}

	return dev, nil
}

func (r *repo) DeleteDevice(ctx context.Context, id string) error {
	// delete a device from the device collection
	_, err := r.deviceCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) GetDeviceConfigurationByID(ctx context.Context, id string) (*fleetpb.DeviceConfiguration, error) {
	// get a device configuration by its ID, this is used to get a specific device configuration
	filter := bson.M{"_id": id}
	var deviceConfiguration fleetpb.DeviceConfiguration
	err := r.deviceConfigCollection.FindOne(ctx, filter).Decode(&deviceConfiguration)
	if err != nil {
		return nil, fleet.ErrDeviceConfigurationNotFound
	}
	return &deviceConfiguration, nil
}

func (r *repo) ListDeviceConfiguration(ctx context.Context, params *fleetpb.ListDeviceConfigurationsParameters) ([]*fleetpb.DeviceConfiguration, error) {
	// list device configurations, this is used to list all device configurations
	filter := bson.M{}
	if params.DeviceId != "" {
		filter["device_id"] = params.DeviceId
	}

	var deviceConfigurations []*fleetpb.DeviceConfiguration
	cur, err := r.deviceConfigCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var deviceConfiguration fleetpb.DeviceConfiguration
		err := cur.Decode(&deviceConfiguration)
		if err != nil {
			return nil, err
		}
		deviceConfigurations = append(deviceConfigurations, &deviceConfiguration)
	}
	if err := cur.Err(); err != nil {

		return nil, err
	}
	return deviceConfigurations, nil
}

func (r *repo) UpsertDeviceConfiguration(ctx context.Context, conf *fleetpb.DeviceConfiguration) (*fleetpb.DeviceConfiguration, error) {
	// upsert a device configuration into the device configuration collection and return the device configuration
	if conf.Id == "" { // if the device configuration does not have an ID, then generate one and set the created time
		conf.Id = fleet.NewID()
		conf.Created = timestamppb.New(time.Now())
	}

	_, err := r.deviceConfigCollection.UpdateOne(ctx, bson.M{"_id": conf.Id}, bson.M{"$set": conf}, &options.UpdateOptions{Upsert: &[]bool{true}[0]})
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (r *repo) DeleteDeviceConfiguration(ctx context.Context, id string) error {

	// delete a device configuration from the device configuration collection
	_, err := r.deviceConfigCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
