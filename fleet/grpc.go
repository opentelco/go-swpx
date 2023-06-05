package fleet

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/fleetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcServer struct {
	grpc  *grpc.Server
	fleet fleetpb.FleetServiceServer

	fleetpb.UnimplementedFleetServiceServer
}

func NewGRPC(fleet fleetpb.FleetServiceServer, srv *grpc.Server) {
	impl := &grpcServer{fleet: fleet, grpc: srv}
	fleetpb.RegisterFleetServiceServer(srv, impl)
}

// CollectDevice collects information about the device from the network (with the help of the poller)
// and returns the device with the updated information
func (s *grpcServer) CollectDevice(ctx context.Context, params *fleetpb.CollectDeviceParameters) (*devicepb.Device, error) {
	res, err := s.fleet.CollectDevice(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil
}

// *** Configuration ***
// CollectConfig collects the running configuration from the device in the network (with the help of the poller) and
// returns the config as a string
func (s *grpcServer) CollectConfig(ctx context.Context, params *fleetpb.CollectConfigParameters) (*configurationpb.Configuration, error) {
	res, err := s.fleet.CollectConfig(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// DeleteDevice deletes the device, its configuration and all changes related to the device
func (f *grpcServer) DeleteDevice(ctx context.Context, params *devicepb.DeleteParameters) (*emptypb.Empty, error) {
	return f.fleet.DeleteDevice(ctx, params)
}

func (f *grpcServer) DiscoverDevice(ctx context.Context, params *devicepb.CreateParameters) (*devicepb.Device, error) {
	res, err := f.fleet.DiscoverDevice(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil
}

func errHandler(err error) error {
	if err == nil {
		return nil
	}
	switch err {
	case ErrNotImplemented:
		return grpc.Errorf(codes.Unimplemented, err.Error())
	default:
		return grpc.Errorf(codes.Internal, err.Error())
	}
}
