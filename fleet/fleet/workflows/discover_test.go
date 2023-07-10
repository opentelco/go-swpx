package workflows

import (
	"fmt"
	"time"

	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelementpb"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (t *unitTestSuite) TestDiscoverWorkflowParamsHost() {
	testTime := time.Date(2023, 06, 9, 0, 0, 0, 0, time.UTC)
	safeNow = func(ctx workflow.Context) time.Time {
		return testTime
	}
	env := t.NewTestWorkflowEnvironment()

	env.RegisterWorkflow(DiscoverWorkflow)
	params := &devicepb.CreateParameters{
		Hostname:      &[]string{"host-a1"}[0],
		NetworkRegion: &[]string{"ABC"}[0],
	}

	discoverWithPollerParams := &corepb.DiscoverRequest{
		Session: &corepb.SessionRequest{
			Hostname: "host-a1",
		},
		Settings: &corepb.Settings{
			ResourcePlugin: "generic",
			RecreateIndex:  false,
			Timeout:        "30s",
			TqChannel:      corepb.Settings_CHANNEL_PRIMARY,
			Priority:       corepb.Settings_DEFAULT,
		},
	}

	discoverWithPollerResp := &corepb.DiscoverResponse{
		NetworkElement: &networkelementpb.Element{
			Sysname:      "host-a1",
			Version:      "1.0.0",
			SnmpObjectId: "1.0.0.0.232.23132.2.0",
			Uptime:       "2023-06-09T00:00:00Z",
		},
	}
	env.OnActivity(testAct.DiscoverWithPoller, mock.Anything, discoverWithPollerParams).Return(discoverWithPollerResp, nil).Once()

	createDeviceParams := &devicepb.CreateParameters{
		Hostname: &[]string{"host-a1"}[0],
		Sysname:  &discoverWithPollerResp.NetworkElement.Sysname,
		// Model:         &discoverWithPollerResp.NetworkElementpb.SnmpObjectId,
		Version:       &discoverWithPollerResp.NetworkElement.Version,
		NetworkRegion: params.NetworkRegion,
		LastSeen:      timestamppb.New(testTime),
		LastReboot:    timestamppb.New(testTime),
		State:         &stateActive,
		Status:        &statusReachable,
	}
	returnDevice := &devicepb.Device{Id: "1234"}
	env.OnActivity(testDevAct.CreateDevice, mock.Anything, createDeviceParams).Return(returnDevice, nil).Once()
	env.OnActivity(testDevAct.AddDeviceEvent, mock.Anything, &devicepb.Event{
		DeviceId: returnDevice.Id,
		Type:     devicepb.Event_DEVICE,
		Message:  "device was created by discovery",
		Action:   devicepb.Event_CREATE,
		Outcome:  devicepb.Event_SUCCESS,
	}).Return(nil, nil).Once()
	env.ExecuteWorkflow(DiscoverWorkflow, params)

	t.True(env.IsWorkflowCompleted())
	err := env.GetWorkflowError()

	fmt.Println(err)
	env.AssertExpectations(t.T())

}
