package fleet

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcServer struct {
	grpc  *grpc.Server
	fleet fleetpb.FleetServer

	fleetpb.UnimplementedFleetServer
}

func NewGRPC(fleet fleetpb.FleetServer, srv *grpc.Server) {
	impl := &grpcServer{fleet: fleet, grpc: srv}
	fleetpb.RegisterFleetServer(srv, impl)
}

func (s *grpcServer) GetDeviceByID(ctx context.Context, params *fleetpb.GetDeviceByIDParameters) (*fleetpb.Device, error) {
	res, err := s.fleet.GetDeviceByID(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// CollectDevice collects information about the device from the network (with the help of the poller)
// and returns the device with the updated information
func (s *grpcServer) CollectDevice(ctx context.Context, params *fleetpb.CollectDeviceParameters) (*fleetpb.Device, error) {
	res, err := s.fleet.CollectDevice(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil
}

// Get a device by its hostname, managment ip or serial number etc (used to search for a device)
func (s *grpcServer) ListDevices(ctx context.Context, params *fleetpb.ListDevicesParameters) (*fleetpb.ListDevicesResponse, error) {
	res, err := s.fleet.ListDevices(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// Create a device in the fleet
func (s *grpcServer) CreateDevice(ctx context.Context, params *fleetpb.CreateDeviceParameters) (*fleetpb.Device, error) {
	res, err := s.fleet.CreateDevice(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil
}

// Update a device in the fleet (this is used to update the device with new information)
func (s *grpcServer) UpdateDevice(ctx context.Context, params *fleetpb.UpdateDeviceParameters) (*fleetpb.Device, error) {
	res, err := s.fleet.UpdateDevice(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// Delete a device from the fleet (mark the device as deleted)
func (s *grpcServer) DeleteDevice(ctx context.Context, params *fleetpb.DeleteDeviceParameters) (*emptypb.Empty, error) {
	res, err := s.fleet.DeleteDevice(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil
}

// *** Configuration ***
// CollectConfig collects the running configuration from the device in the network (with the help of the poller) and
// returns the config as a string
func (s *grpcServer) CollectConfig(ctx context.Context, params *fleetpb.CollectConfigParameters) (*fleetpb.DeviceConfiguration, error) {
	res, err := s.fleet.CollectConfig(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// Get a device configuration by its ID, this is used to get a specific device configuration
func (s *grpcServer) GetDeviceConfigurationByID(ctx context.Context, params *fleetpb.GetDeviceConfigurationByIDParameters) (*fleetpb.DeviceConfiguration, error) {
	res, err := s.fleet.GetDeviceConfigurationByID(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// CompareDeviceConfiguration compares the configuration of a device with the configuration in the database and returns the changes
// if no specific configuration is specified the latest configuration is used to compare with
func (s *grpcServer) CompareDeviceConfiguration(ctx context.Context, params *fleetpb.CompareDeviceConfigurationParameters) (*fleetpb.CompareDeviceConfigurationResponse, error) {
	res, err := s.fleet.CompareDeviceConfiguration(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// ListDeviceConfigurations lists all configurations for a device
func (s *grpcServer) ListDeviceConfigurations(ctx context.Context, params *fleetpb.ListDeviceConfigurationsParameters) (*fleetpb.ListDeviceConfigurationsResponse, error) {
	res, err := s.fleet.ListDeviceConfigurations(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// Create a device configuration in the fleet (this is used to store the configuration of a device)
func (s *grpcServer) CreateDeviceConfiguration(ctx context.Context, params *fleetpb.CreateDeviceConfigurationParameters) (*fleetpb.DeviceConfiguration, error) {
	res, err := s.fleet.CreateDeviceConfiguration(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// Delete a device configuration from the fleet (removes the configuration from the database)
func (s *grpcServer) DeleteDeviceConfiguration(ctx context.Context, params *fleetpb.DeleteDeviceConfigurationParameters) (*emptypb.Empty, error) {
	res, err := s.fleet.DeleteDeviceConfiguration(ctx, params)
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
	case ErrDeviceConfigurationNotFound:
		return grpc.Errorf(codes.NotFound, err.Error())
	case ErrDeviceNotFound:
		return grpc.Errorf(codes.NotFound, err.Error())
	case ErrDeviceAlreadyExists:
		return grpc.Errorf(codes.AlreadyExists, err.Error())
	case ErrDeviceConfigurationInvalid:
		return grpc.Errorf(codes.InvalidArgument, err.Error())
	case ErrDeviceInvalid:
		return grpc.Errorf(codes.InvalidArgument, err.Error())
	case ErrDeviceConfigurationInvalid:
		return grpc.Errorf(codes.InvalidArgument, err.Error())
	case ErrNotImplemented:
		return grpc.Errorf(codes.Unimplemented, err.Error())
	default:
		return grpc.Errorf(codes.Internal, err.Error())
	}
}
