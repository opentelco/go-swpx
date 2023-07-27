package repo

import (
	"context"
	"time"

	"git.liero.se/opentelco/go-swpx/database"
	"git.liero.se/opentelco/go-swpx/fleet/stanza"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultDatabase           string = "swpx"
	defaultTemplateCollection string = "stanzas"
)

type repo struct {
	mc               *mongo.Client
	stanzaCollection *mongo.Collection
	logger           hclog.Logger
}

func New(mc *mongo.Client, database string, logger hclog.Logger) (stanza.Repository, error) {
	if database == "" {
		database = defaultDatabase
	}

	return &repo{
		mc:               mc,
		stanzaCollection: mc.Database(database).Collection(defaultTemplateCollection),
		logger:           logger.Named("db"),
	}, nil
}

func (r *repo) GetByID(ctx context.Context, id string) (*stanzapb.Stanza, error) {
	// get a stanza by its ID, this is used to get a specific device configuration
	filter := bson.M{"_id": id}
	var st stanzapb.Stanza
	err := r.stanzaCollection.FindOne(ctx, filter).Decode(&st)
	if err != nil {
		return nil, stanza.ErrStanzaNotFound
	}
	return &st, nil
}

func (r *repo) List(ctx context.Context, params *stanzapb.ListRequest) (*stanzapb.ListResponse, error) {
	// list stanzas, default filters will list all stanzas not applied to a device (device_id = "")
	filter := bson.M{
		"device_id": "",
	}

	for _, v := range params.Filters {
		switch v {
		case stanzapb.ListRequest_FILTER_APPLIED:
			delete(filter, "device_id")
			// set applied_at to not null
			filter["applied_at"] = bson.M{"$ne": nil}

		case stanzapb.ListRequest_FILTER_NOT_APPLIED:
			delete(filter, "device_id")
			// set applied_at to null
			filter["applied_at"] = nil

		}
	}

	if params.DeviceId != nil {
		filter["device_id"] = params.DeviceId
	}

	if params.DeviceType != nil {
		filter["device_type"] = params.DeviceType
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

	var stanzas []*stanzapb.Stanza
	cur, err := r.stanzaCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var stanza stanzapb.Stanza
		err := cur.Decode(&stanza)
		if err != nil {
			return nil, err
		}
		stanzas = append(stanzas, &stanza)
	}
	if err := cur.Err(); err != nil {

		return nil, err
	}

	resp := &stanzapb.ListResponse{
		Stanzas: stanzas,
		Total:   int64(len(stanzas)),
	}
	return resp, nil
}

func (r *repo) Upsert(ctx context.Context, stanza *stanzapb.Stanza) (*stanzapb.Stanza, error) {
	// upsert a device configuration into the device configuration collection and return the device configuration
	if stanza.Id == "" { // if the device configuration does not have an ID, then generate one and set the created time
		stanza.Id = database.NewID()
		stanza.CreatedAt = timestamppb.New(time.Now())
	}

	stanza.UpdatedAt = timestamppb.New(time.Now())

	_, err := r.stanzaCollection.UpdateOne(ctx, bson.M{"_id": stanza.Id}, bson.M{"$set": stanza}, &options.UpdateOptions{Upsert: &[]bool{true}[0]})
	if err != nil {
		return nil, err
	}
	return stanza, nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	// delete a device configuration from the device configuration collection
	_, err := r.stanzaCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
