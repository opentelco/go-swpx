package device

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"github.com/gogo/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcServer struct {
	grpc   *grpc.Server
	device devicepb.DeviceServiceServer

	devicepb.UnimplementedDeviceServiceServer
}

func NewGRPC(device devicepb.DeviceServiceServer, srv *grpc.Server) {
	impl := &grpcServer{device: device, grpc: srv}
	devicepb.RegisterDeviceServiceServer(srv, impl)
}

func (s *grpcServer) GetByID(ctx context.Context, params *devicepb.GetByIDParameters) (*devicepb.Device, error) {
	res, err := s.device.GetByID(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// Get a device by its hostname, managment ip or serial number etc (used to search for a device)
func (s *grpcServer) List(ctx context.Context, params *devicepb.ListParameters) (*devicepb.ListResponse, error) {
	res, err := s.device.List(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// Create a device in the device
func (s *grpcServer) Create(ctx context.Context, params *devicepb.CreateParameters) (*devicepb.Device, error) {
	res, err := s.device.Create(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil
}

// Update a device in the device (this is used to update the device with new information)
func (s *grpcServer) Update(ctx context.Context, params *devicepb.UpdateParameters) (*devicepb.Device, error) {
	res, err := s.device.Update(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// Delete a device from the device (mark the device as deleted)
func (s *grpcServer) Delete(ctx context.Context, params *devicepb.DeleteParameters) (*emptypb.Empty, error) {
	res, err := s.device.Delete(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil
}

func (s *grpcServer) ListChanges(ctx context.Context, params *devicepb.ListChangesParameters) (*devicepb.ListChangesResponse, error) {
	res, err := s.device.ListChanges(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

func (s *grpcServer) GetChangeByID(ctx context.Context, params *devicepb.GetChangeByIDParameters) (*devicepb.Change, error) {
	res, err := s.device.GetChangeByID(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

func (d *grpcServer) GetEventByID(ctx context.Context, params *devicepb.GetEventByIDParameters) (*devicepb.Event, error) {
	res, err := d.device.GetEventByID(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// returns a list of events (default 100)
func (d *grpcServer) ListEvents(ctx context.Context, params *devicepb.ListEventsParameters) (*devicepb.ListEventsResponse, error) {
	res, err := d.device.ListEvents(ctx, params)
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
	case ErrDeviceNotFound:
		return status.Errorf(codes.NotFound, err.Error())
	case ErrDeviceAlreadyExists:
		return status.Errorf(codes.AlreadyExists, err.Error())
	case ErrNotImplemented:
		return status.Errorf(codes.Unimplemented, err.Error())
	default:
		return status.Errorf(codes.Internal, err.Error())
	}
}
