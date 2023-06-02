package configuration

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"github.com/gogo/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcServer struct {
	grpc *grpc.Server
	c    configurationpb.ConfigurationServiceServer

	configurationpb.UnimplementedConfigurationServiceServer
}

func NewGRPC(c configurationpb.ConfigurationServiceServer, srv *grpc.Server) {
	impl := &grpcServer{c: c, grpc: srv}
	configurationpb.RegisterConfigurationServiceServer(srv, impl)
}

// Get a device configuration by its ID, this is used to get a specific device configuration
func (s *grpcServer) GetByID(ctx context.Context, params *configurationpb.GetByIDParameters) (*configurationpb.Configuration, error) {
	res, err := s.c.GetByID(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// Compare compares the configuration of a device with the configuration in the database and returns the changes
// if no specific configuration is specified the latest configuration is used to compare with
func (s *grpcServer) Compare(ctx context.Context, params *configurationpb.CompareParameters) (*configurationpb.CompareResponse, error) {
	res, err := s.c.Compare(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// List lists all configurations for a device
func (s *grpcServer) List(ctx context.Context, params *configurationpb.ListParameters) (*configurationpb.ListResponse, error) {
	res, err := s.c.List(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// Create a device configuration in the device (this is used to store the configuration of a device)
func (s *grpcServer) Create(ctx context.Context, params *configurationpb.CreateParameters) (*configurationpb.Configuration, error) {
	res, err := s.c.Create(ctx, params)
	if err != nil {
		return nil, errHandler(err)
	}
	return res, nil

}

// Delete a device configuration from the device (removes the configuration from the database)
func (s *grpcServer) Delete(ctx context.Context, params *configurationpb.DeleteParameters) (*emptypb.Empty, error) {
	res, err := s.c.Delete(ctx, params)
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
	case ErrConfigurationNotFound:
		return status.Errorf(codes.NotFound, err.Error())
	case ErrConfigurationInvalid:
		return status.Errorf(codes.InvalidArgument, err.Error())
	case ErrConfigurationInvalid:
		return status.Errorf(codes.InvalidArgument, err.Error())
	case ErrNotImplemented:
		return status.Errorf(codes.Unimplemented, err.Error())
	default:
		return status.Errorf(codes.Internal, err.Error())
	}
}
