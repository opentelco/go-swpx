package repo

import (
	"context"
	"runtime"
	"testing"
	"time"

	"git.liero.se/opentelco/go-swpx/fleet/fleet/utils"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	codecs "github.com/amsokol/mongo-go-driver-protobuf"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/suite"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var db = "TEST_DEVICE"
var (
	testDeviceId1 = "test-device-1"
	testDeviceId2 = "test-device-2"
	testDeviceId3 = "test-device-3"
	testDeviceId4 = "test-device-4"
)

// Language: go
// test suite struct
type repoTestSuite struct {
	suite.Suite
	client *mongo.Client
	logger hclog.Logger
}

// setup test suite
func (s *repoTestSuite) SetupTest() {
	// do something
	s.client = s.setupDb()
	s.logger = hclog.New(hclog.DefaultOptions)
	s.populateCollections()
}

// tear down test suite
func (s *repoTestSuite) TearDownTest() {
	_ = s.client.Disconnect(context.Background())
}

// run suite
func Test_Suite(t *testing.T) {
	suite.Run(t, new(repoTestSuite))
}

func (s *repoTestSuite) setupDb() *mongo.Client {
	opts := &memongo.Options{
		MongoVersion: "5.0.0",
	}
	if runtime.GOARCH == "arm64" {
		if runtime.GOOS == "darwin" {
			// Only set the custom url as workaround for arm64 macs
			opts.DownloadURL = "https://fastdl.mongodb.org/osx/mongodb-macos-x86_64-5.0.5.tgz"
		}
	}
	opts.StartupTimeout = time.Second * 60
	server, err := memongo.StartWithOptions(opts)
	if err != nil {
		s.T().Fatal("Stopping! could not start memongo for tests; aborting!")
	}

	reg := codecs.Register(bson.NewRegistryBuilder()).Build()
	clientOpts := options.Client()
	clientOpts.Registry = reg

	// cmdMonitor := &event.CommandMonitor{
	// 	Started: func(_ context.Context, evt *event.CommandStartedEvent) {
	// 		log.Print(evt.Command)
	// 	},
	// }
	// clientOpts.SetMonitor(cmdMonitor)

	clientOpts.ApplyURI(server.URI())
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		s.T().Fatal("Stopping! could not create Mongo client for tests")
	}

	return client
}

func (s *repoTestSuite) populateCollections() {
	_, _ = s.client.
		Database(db).
		Collection(defaultDeviceCollection).InsertOne(
		context.Background(),
		&devicepb.Device{
			Id:                   testDeviceId1,
			Hostname:             "device-1",
			ManagementIp:         "192.168.1.101",
			Model:                "S5720",
			Version:              "",
			NetworkRegion:        "",
			PollerResourcePlugin: "",
			PollerProviderPlugin: "",
			State:                devicepb.Device_DEVICE_STATE_ACTIVE,
			Status:               devicepb.Device_DEVICE_STATUS_REACHABLE,
			Schedules: []*devicepb.Device_Schedule{
				{
					Interval: durationpb.New(time.Hour * 24),
					Type:     devicepb.Device_Schedule_COLLECT_CONFIG,
					LastRun:  timestamppb.New(time.Now().Add(-time.Hour * 2)),
					Active:   true,
				},
				{
					Interval: durationpb.New(time.Hour * 1),
					Type:     devicepb.Device_Schedule_COLLECT_DEVICE,
					LastRun:  timestamppb.New(time.Now().Add(-time.Hour * 1)),
					Active:   true,
				},
			},
		},
	)
	_, _ = s.client.
		Database(db).
		Collection(defaultDeviceCollection).InsertOne(
		context.Background(),
		&devicepb.Device{
			Id:                   testDeviceId2,
			Hostname:             "device-2",
			ManagementIp:         "192.168.1.102",
			Model:                "S5720",
			Version:              "",
			NetworkRegion:        "",
			PollerResourcePlugin: "",
			PollerProviderPlugin: "",
			State:                devicepb.Device_DEVICE_STATE_ACTIVE,
			Status:               devicepb.Device_DEVICE_STATUS_REACHABLE,
			Schedules: []*devicepb.Device_Schedule{
				{
					Interval: durationpb.New(time.Hour * 24),
					Type:     devicepb.Device_Schedule_COLLECT_CONFIG,
					LastRun:  timestamppb.New(time.Now().Add(-time.Hour * 2)),
					Active:   true,
				},
				{
					Interval: durationpb.New(time.Hour * 1),
					Type:     devicepb.Device_Schedule_COLLECT_DEVICE,
					LastRun:  timestamppb.New(time.Now().Add(-time.Hour * 1)),
					Active:   true,
				},
			},
		},
	)
	_, _ = s.client.
		Database(db).
		Collection(defaultDeviceCollection).InsertOne(
		context.Background(),
		&devicepb.Device{
			Id:                   testDeviceId3,
			Hostname:             "device-3",
			ManagementIp:         "192.168.1.103",
			Model:                "S5720",
			Version:              "",
			NetworkRegion:        "",
			PollerResourcePlugin: "",
			PollerProviderPlugin: "",
			State:                devicepb.Device_DEVICE_STATE_ACTIVE,
			Status:               devicepb.Device_DEVICE_STATUS_REACHABLE,
			Schedules: []*devicepb.Device_Schedule{
				{
					Interval: durationpb.New(time.Hour * 24),
					Type:     devicepb.Device_Schedule_COLLECT_CONFIG,
					LastRun:  timestamppb.New(time.Now().Add(-time.Hour * 2)),
					Active:   true,
				},
				{
					Interval: durationpb.New(time.Hour * 1),
					Type:     devicepb.Device_Schedule_COLLECT_DEVICE,
					LastRun:  timestamppb.New(time.Now().Add(-time.Hour * 1)),
					Active:   false,
				},
			},
		})
	_, _ = s.client.
		Database(db).
		Collection(defaultDeviceCollection).InsertOne(
		context.Background(),
		&devicepb.Device{
			Id:                   testDeviceId4,
			Hostname:             "device-4",
			ManagementIp:         "192.168.1.104",
			Model:                "S5720",
			Version:              "",
			NetworkRegion:        "",
			PollerResourcePlugin: "",
			PollerProviderPlugin: "",
			State:                devicepb.Device_DEVICE_STATE_ACTIVE,
			Status:               devicepb.Device_DEVICE_STATUS_REACHABLE,
			Schedules: []*devicepb.Device_Schedule{
				{
					Interval: durationpb.New(time.Hour * 24),
					Type:     devicepb.Device_Schedule_COLLECT_CONFIG,
					LastRun:  timestamppb.New(time.Now().Add(-time.Hour * 2)),
					Active:   false,
				},
				{
					Interval: durationpb.New(time.Hour * 1),
					Type:     devicepb.Device_Schedule_COLLECT_DEVICE,
					LastRun:  timestamppb.New(time.Now().Add(-time.Hour * 1)),
					Active:   false,
				},
			},
		},
	)
}

func (s *repoTestSuite) TestDeviceRepo_GetDevice() {

	repo, _ := New(s.client, db, s.logger)
	hasFiringSchedule := true
	t := devicepb.Device_Schedule_COLLECT_DEVICE
	res, err := repo.List(context.Background(), &devicepb.ListParameters{
		HasFiringSchedule: &hasFiringSchedule,
		ScheduleType:      &t,
	})
	s.NoError(err)
	s.Len(res, 2)

	t = devicepb.Device_Schedule_COLLECT_CONFIG
	res, err = repo.List(context.Background(), &devicepb.ListParameters{
		HasFiringSchedule: &hasFiringSchedule,
		ScheduleType:      &t,
	})
	s.NoError(err)
	s.Len(res, 3)

}

func (s *repoTestSuite) Test_SetSchedule() {
	repo, _ := New(s.client, db, s.logger)
	ctx := context.Background()

	d, err := repo.UpsertSchedule(ctx, testDeviceId1, &devicepb.Device_Schedule{
		Interval:    durationpb.New(time.Hour * 1),
		Type:        devicepb.Device_Schedule_COLLECT_DEVICE,
		LastRun:     timestamppb.New(time.Now().Add(-time.Hour * 1)),
		Active:      true,
		FailedCount: 3,
	})
	s.NoError(err)
	s.NotNil(d)
	s.Len(d.Schedules, 2)

	x := utils.GetDeviceScheduleByType(d, devicepb.Device_Schedule_COLLECT_DEVICE)
	s.Equal(x.FailedCount, int64(3))

}

// Test_Get_By_FiringSchedule tests the Get method with the HasFiringSchedule
// parameter set to true.
// The test creates 4 devices, 2 of them have a firing schedule
func (s *repoTestSuite) Test_Get_By_FiringSchedule() {
	repo, _ := New(s.client, db, s.logger)
	ctx := context.Background()

	// one device that was collected 3 hours ago with a schedule that runs every 24 hours
	// this should NOT be returned by the list method
	if _, err := repo.UpsertSchedule(ctx, testDeviceId1, &devicepb.Device_Schedule{
		Interval: durationpb.New(time.Hour * 24),
		Type:     devicepb.Device_Schedule_COLLECT_CONFIG,
		LastRun:  timestamppb.New(time.Now().Add(-time.Hour * 1)),
		Active:   true,
	}); err != nil {
		s.NoError(err)
	}

	// one device that was collected 3 hours ago with a schedule that runs every 1 hour
	// this should be returned by the list method
	if _, err := repo.UpsertSchedule(ctx, testDeviceId2, &devicepb.Device_Schedule{
		Interval: durationpb.New(time.Hour * 1),
		Type:     devicepb.Device_Schedule_COLLECT_CONFIG,
		LastRun:  timestamppb.New(time.Now().Add(-time.Hour * 3)),
		Active:   true,
	}); err != nil {
		s.NoError(err)
	}

	// one device that was collected 3 hours ago with a schedule that runs every 6 hours
	// but the schedule is not active
	// this should NOT be returned by the list method
	if _, err := repo.UpsertSchedule(ctx, testDeviceId3, &devicepb.Device_Schedule{
		Interval: durationpb.New(time.Hour * 6),
		Type:     devicepb.Device_Schedule_COLLECT_CONFIG,
		LastRun:  timestamppb.New(time.Now().Add(-time.Hour * 10)),
		Active:   false,
	}); err != nil {
		s.NoError(err)
	}

	// one device that was collected 10 hours ago with a schedule that runs every 8 hours
	// this should be returned by the list method
	if _, err := repo.UpsertSchedule(ctx, testDeviceId4, &devicepb.Device_Schedule{
		Interval: durationpb.New(time.Hour * 8),
		Type:     devicepb.Device_Schedule_COLLECT_CONFIG,
		LastRun:  timestamppb.New(time.Now().Add(-time.Hour * 10)),
		Active:   true,
	}); err != nil {
		s.NoError(err)
	}

	hasFiring := true
	t := devicepb.Device_Schedule_COLLECT_CONFIG

	d, err := repo.List(ctx, &devicepb.ListParameters{
		HasFiringSchedule: &hasFiring,
		ScheduleType:      &t,
	})

	s.NoError(err)
	s.NotNil(d)
	s.Len(d, 2)
	// print as json

}
