package repo

import (
	"context"
	"fmt"
	"time"

	"git.liero.se/opentelco/go-swpx/database"
	"git.liero.se/opentelco/go-swpx/fleet/notification"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/notificationpb"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultDatabase               string = "swpx"
	defaultNotificationCollection string = "notifications"
)

type repo struct {
	mc                     *mongo.Client
	notificationCollection *mongo.Collection
	logger                 hclog.Logger
}

func New(mc *mongo.Client, database string, logger hclog.Logger) (notification.Repository, error) {
	if database == "" {
		database = defaultDatabase
	}

	return &repo{
		mc:                     mc,
		notificationCollection: mc.Database(database).Collection(defaultNotificationCollection),
		logger:                 logger.Named("db"),
	}, nil
}

func (r *repo) GetByID(ctx context.Context, id string) (*notificationpb.Notification, error) {
	// get a notification by its ID, this is used to get a specific notification
	if id == "" {
		return nil, notification.ErrNotificationNotFound
	}
	filter := bson.M{"_id": id}
	opts := options.FindOne()

	var nf notificationpb.Notification
	if err := r.notificationCollection.FindOne(ctx, filter, opts).Decode(&nf); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, notification.ErrNotificationNotFound
		}
		return nil, err
	}

	return &nf, nil
}

// List returns a list of notifications
func (r *repo) List(ctx context.Context, params *notificationpb.ListRequest) (*notificationpb.ListResponse, error) {

	var filter = bson.M{
		"read": false,
	}

	for _, f := range params.Filter {
		// all filters are OR filters
		switch f {
		case notificationpb.ListRequest_INCLUDE_READ:
			// filter by read true and false
			delete(filter, "read")

		case notificationpb.ListRequest_RESOURCE_TYPE_CONFIG:
			// filter by resource type
			filter["$or"] = []bson.M{
				{"resource_type": notificationpb.Notification_RESOURCE_TYPE_CONFIG},
			}

		case notificationpb.ListRequest_RESOURCE_TYPE_DEVICE:
			// filter by resource type
			filter["$or"] = []bson.M{
				{"resource_type": notificationpb.Notification_RESOURCE_TYPE_DEVICE},
			}

		}
	}
	if params.ResourceIds != nil {
		// filter by resource ID
		filter["resource_id"] = bson.M{"$in": params.ResourceIds}
	}

	if params.Ids != nil {
		// filter by notification ID
		filter["_id"] = bson.M{"$in": params.Ids}
	}

	opts := options.Find()
	if params.Limit != nil {
		opts.SetLimit(*params.Limit)
	}
	if params.Offset != nil {
		opts.SetSkip(*params.Offset)
	}

	cur, err := r.notificationCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var notifications []*notificationpb.Notification
	for cur.Next(ctx) {
		var nf notificationpb.Notification
		if err := cur.Decode(&nf); err != nil {
			return nil, err
		}
		notifications = append(notifications, &nf)
	}

	return &notificationpb.ListResponse{
		Notifications: notifications,
	}, nil

}

func (r *repo) Upsert(ctx context.Context, notification *notificationpb.Notification) (*notificationpb.Notification, error) {
	if notification.Id == "" { // if the notification does not have an ID, then generate one and set the created time
		notification.Id = database.NewID()
		notification.Timestamp = timestamppb.New(time.Now())
	}

	_, err := r.notificationCollection.UpdateOne(ctx, bson.M{"_id": notification.Id}, bson.M{"$set": notification},
		&options.UpdateOptions{Upsert: &[]bool{true}[0]})
	if err != nil {
		return nil, err
	}
	return notification, nil
}

// Mark one or more notifications as read
func (r *repo) MarkAsRead(ctx context.Context, params *notificationpb.MarkAsReadRequest) (*notificationpb.MarkAsReadResponse, error) {
	for i, id := range params.Ids {
		if id == "" {
			return nil, fmt.Errorf("invalid notification ID on place: %d: %w", i, notification.ErrInvalidArgument)
		}
	}

	// update all notifications with the given IDs
	filter := bson.M{"_id": bson.M{"$in": params.Ids}}
	update := bson.M{"$set": bson.M{"read": true}}
	_, err := r.notificationCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	res, err := r.List(ctx, &notificationpb.ListRequest{Ids: params.Ids})
	if err != nil {
		return nil, fmt.Errorf("could not get notiications after setting as read: %w", err)
	}

	return &notificationpb.MarkAsReadResponse{
		Notifications: res.Notifications,
	}, nil

}

// Delete deletes an existing notification
func (r *repo) Delete(ctx context.Context, id string) error {
	if id == "" {
		return notification.ErrInvalidArgument
	}

	filter := bson.M{"_id": id}
	_, err := r.notificationCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
