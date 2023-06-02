package cmd

import (
	"time"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/fleetpb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func prettyPrintJSON(m proto.Message) string {
	json := protojson.MarshalOptions{
		Multiline:       true,
		Indent:          "  ",
		AllowPartial:    false,
		UseProtoNames:   false,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
	}
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

func getDeviceClient(cmd *cobra.Command) (devicepb.DeviceServiceClient, error) {

	addr, _ := cmd.Flags().GetString("fleet-addr")
	conn, err := grpc.Dial(addr, grpc.WithTimeout(5*time.Second), grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return devicepb.NewDeviceServiceClient(conn), nil
}

func getFleetClient(cmd *cobra.Command) (fleetpb.FleetServiceClient, error) {

	addr, _ := cmd.Flags().GetString("fleet-addr")
	conn, err := grpc.Dial(addr, grpc.WithTimeout(5*time.Second), grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return fleetpb.NewFleetServiceClient(conn), nil
}
