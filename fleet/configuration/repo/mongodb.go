package repo

import (
	"context"
	"time"

	"git.liero.se/opentelco/go-swpx/database"
	"git.liero.se/opentelco/go-swpx/fleet/configuration"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultDatabase                string = "swpx"
	defaultConfigurationCollection string = "configurations"
)

type repo struct {
	mc               *mongo.Client
	configCollection *mongo.Collection
	logger           hclog.Logger
}

func New(mc *mongo.Client, database string, logger hclog.Logger) (configuration.Repository, error) {
	if database == "" {
		database = defaultDatabase
	}

	return &repo{
		mc:               mc,
		configCollection: mc.Database(database).Collection(defaultConfigurationCollection),
		logger:           logger.Named("db"),
	}, nil
}

func (r *repo) GetByID(ctx context.Context, id string) (*configurationpb.Configuration, error) {
	// get a device configuration by its ID, this is used to get a specific device configuration
	filter := bson.M{"_id": id}
	var deviceConfiguration configurationpb.Configuration
	err := r.configCollection.FindOne(ctx, filter).Decode(&deviceConfiguration)
	if err != nil {
		return nil, configuration.ErrConfigurationNotFound
	}
	return &deviceConfiguration, nil
}

func (r *repo) List(ctx context.Context, params *configurationpb.ListParameters) ([]*configurationpb.Configuration, error) {
	// list device configurations, this is used to list all device configurations
	filter := bson.M{}
	if params.DeviceId != "" {
		filter["device_id"] = params.DeviceId
	}
	// order by created time (newest first)
	opts := options.Find().SetSort(bson.M{"created": -1})
	// set limit and offset
	if params.Limit != nil {
		opts.SetLimit(*params.Limit)
	}
	if params.Offset != nil {
		opts.SetSkip(*params.Offset)
	}

	var deviceConfigurations []*configurationpb.Configuration
	cur, err := r.configCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var deviceConfiguration configurationpb.Configuration
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

func (r *repo) Upsert(ctx context.Context, conf *configurationpb.Configuration) (*configurationpb.Configuration, error) {
	// upsert a device configuration into the device configuration collection and return the device configuration
	if conf.Id == "" { // if the device configuration does not have an ID, then generate one and set the created time
		conf.Id = database.NewID()
		conf.Created = timestamppb.New(time.Now())
	}

	_, err := r.configCollection.UpdateOne(ctx, bson.M{"_id": conf.Id}, bson.M{"$set": conf}, &options.UpdateOptions{Upsert: &[]bool{true}[0]})
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (r *repo) Delete(ctx context.Context, id string) error {

	// delete a device configuration from the device configuration collection
	_, err := r.configCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
