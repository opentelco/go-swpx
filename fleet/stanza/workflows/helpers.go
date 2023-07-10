package workflows

import (
	"fmt"
	"time"

	"git.liero.se/opentelco/go-swpx/fleet/device"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"go.temporal.io/sdk/workflow"
)

var devAct = device.Activities{}

func getDeviceById(ctx workflow.Context, deviceId string) (*devicepb.Device, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:           devicepb.TaskQueue_TASK_QUEUE_FLEET_DEVICE.String(),
		StartToCloseTimeout: time.Second * 60,
		WaitForCancellation: false,
	})
	var device devicepb.Device
	if err := workflow.ExecuteActivity(ctx, devAct.GetDeviceByID, deviceId).Get(ctx, &device); err != nil {
		return nil, fmt.Errorf("failed to collect device: %w", err)
	}
	return &device, nil

}
