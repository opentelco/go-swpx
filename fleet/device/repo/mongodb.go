package repo

import (
	"context"
	"fmt"
	"log"
	"time"

	"git.liero.se/opentelco/go-swpx/database"
	"git.liero.se/opentelco/go-swpx/fleet/device"
	"git.liero.se/opentelco/go-swpx/fleet/fleet/utils"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/commonpb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultDatabase          string = "swpx"
	defaultDeviceCollection  string = "devices"
	defaultChangesCollection string = "device_changes"
	defaultEventsCollection  string = "device_events"
)

type repo struct {
	mc *mongo.Client

	deviceCollection  *mongo.Collection
	eventsCollection  *mongo.Collection
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
		eventsCollection:  mc.Database(database).Collection(defaultEventsCollection),
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

func (r *repo) ListChanges(ctx context.Context, params *devicepb.ListChangesParameters) (*devicepb.ListChangesResponse, error) {

	// list changes for a device
	filter := bson.M{}
	if params.DeviceId != "" {
		filter["device_id"] = params.DeviceId
	}

	options := options.Find()
	options.Limit = params.Limit
	options.Skip = params.Offset

	var changes []*devicepb.Change
	cur, err := r.changesCollection.Find(ctx, filter, options)
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

	totalCount, err := r.changesCollection.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	pageInfo := &commonpb.PageInfo{
		Count:  int64(len(changes)),
		Offset: params.Offset,
		Limit:  params.Limit,
		Total:  totalCount,
	}

	return &devicepb.ListChangesResponse{
		Changes:  changes,
		PageInfo: pageInfo,
	}, nil

}

func (r *repo) DeleteChangeByID(ctx context.Context, id string) error {
	// delete a device change by its ID
	_, err := r.changesCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
func (r *repo) DeleteChangesByDeviceID(ctx context.Context, id string) error {
	// delete all device changes for a device
	_, err := r.changesCollection.DeleteMany(ctx, bson.M{"device_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) AddEvent(ctx context.Context, event *devicepb.Event) (*devicepb.Event, error) {
	// add a event to the events collection and return it
	event.Id = database.NewID()
	event.Created = timestamppb.New(time.Now())

	_, err := r.eventsCollection.InsertOne(ctx, event)
	if err != nil {
		return nil, err
	}

	return r.GetEventByID(ctx, event.Id)
}

func (r *repo) GetEventByID(ctx context.Context, id string) (*devicepb.Event, error) {
	// get a event by its ID
	filter := bson.M{"_id": id}
	var event devicepb.Event
	err := r.eventsCollection.FindOne(ctx, filter).Decode(&event)
	if err != nil {
		return nil, err
	}

	return &event, nil

}
func (r *repo) ListEvents(ctx context.Context, params *devicepb.ListEventsParameters) (*devicepb.ListEventsResponse, error) {
	// list events for a device
	filter := bson.M{}
	if params.DeviceId != "" {
		filter["device_id"] = params.DeviceId
	}
	options := options.Find()
	options.Limit = params.Limit
	options.Skip = params.Offset

	var events []*devicepb.Event
	cur, err := r.eventsCollection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {

		var event devicepb.Event
		err := cur.Decode(&event)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	totalCount, err := r.eventsCollection.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	pageInfo := &commonpb.PageInfo{
		Count:  int64(len(events)),
		Offset: params.Offset,
		Limit:  params.Limit,
		Total:  totalCount,
	}

	return &devicepb.ListEventsResponse{
		Events:   events,
		PageInfo: pageInfo,
	}, nil

}

func (r *repo) UpsertSchedule(ctx context.Context, deviceId string, schedule *devicepb.Device_Schedule) (*devicepb.Device, error) {
	// upsert a schedule for a device by deviceId unique by schedule.Typeand return the schedule
	filter := bson.M{"_id": deviceId}
	update := bson.M{"$set": bson.M{"schedules.$[elem]": schedule}}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.type": int(schedule.Type.Number())}, // Filter to match the item by Type
		},
	})
	_, err := r.deviceCollection.FindOneAndUpdate(context.Background(), filter, update, opts).DecodeBytes()
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, deviceId)
}

func GetEnumValues[T interface {
	Number() protoreflect.EnumNumber
}](enums []T) []int {
	values := make([]int, len(enums))
	for i, item := range enums {
		values[i] = int(item.Number())
	}
	return values
}

func (r *repo) List(ctx context.Context, params *devicepb.ListParameters) (*devicepb.ListResponse, error) {
	filter := bson.M{}
	if params.Hostname != nil {
		filter["hostname"] = *params.Hostname
	}
	if params.ManagementIp != nil {
		filter["management_ip"] = *params.ManagementIp
	}

	if params.Search != nil {
		filter["$or"] = bson.A{
			bson.M{"hostname": bson.M{"$regex": params.Search}},
			bson.M{"management_ip": bson.M{"$regex": params.Search}},
		}
	}

	options := options.Find()
	if params.HasFiringSchedule != nil {
		if params.ScheduleType == nil {
			return nil, fmt.Errorf("schedule type is required when filtering by firing schedule")
		}
		filter["schedules"] = bson.M{
			"$elemMatch": bson.M{
				"active": true,
				"type":   protoreflect.EnumNumber(*params.ScheduleType),
			},
		}

		// sort by schedule last_run ascending
		options.SetSort(bson.M{"schedules.last_run": 1})
	}

	options.Limit = params.Limit
	options.Skip = params.Offset

	var devices []*devicepb.Device
	cur, err := r.deviceCollection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var device devicepb.Device
		err := cur.Decode(&device)
		if err != nil {
			return nil, err
		}
		// check if device has a firing schedule and if it has been run
		// this could be done by some amazing mongo query but I'm not sure how and after 1 hour i gave up
		// so we just do it here in code for now
		if params.HasFiringSchedule != nil {
			s := utils.GetDeviceScheduleByType(&device, *params.ScheduleType)
			if s.LastRun == nil {
				devices = append(devices, &device)
			} else {
				// check if now is after last_run + interval and if so add to devices list
				if now.After(s.LastRun.AsTime().Add(time.Duration(s.Interval.Seconds) * time.Second)) {
					devices = append(devices, &device)
				}
			}
		} else {
			devices = append(devices, &device)
		}
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	totalCount, err := r.deviceCollection.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	pageInfo := &commonpb.PageInfo{
		Count:  int64(len(devices)),
		Offset: params.Offset,
		Limit:  params.Limit,
		Total:  totalCount,
	}

	return &devicepb.ListResponse{
		Devices:  devices,
		PageInfo: pageInfo,
	}, nil

}
