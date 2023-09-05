// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: fleet.proto

package fleetpb

import (
	context "context"
	configurationpb "git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	devicepb "git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FleetServiceClient is the client API for FleetService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FleetServiceClient interface {
	// DiscoverDevice discovers a device in the network, creates the device with the information provided
	// by the poller in the discovery. (e.g. sysname, ip address, mac address, etc)
	// NOTE: hostname OR management_ip must be set
	// Creating a device will append the default schedles to the device.
	// - CollectDevice every hour
	// - CollectConfig every 24 hours
	DiscoverDevice(ctx context.Context, in *DiscoverDeviceParameters, opts ...grpc.CallOption) (*devicepb.Device, error)
	// CollectDevice collects information about the device from the network (with the help of the poller)
	// and returns the device with the updated information
	CollectDevice(ctx context.Context, in *CollectDeviceParameters, opts ...grpc.CallOption) (*devicepb.Device, error)
	// CollectConfig collects the running configuration from the device in the network (with the help of the poller) and
	// returns the config as a string
	CollectConfig(ctx context.Context, in *CollectConfigParameters, opts ...grpc.CallOption) (*configurationpb.Configuration, error)
	// Delete a device, chagnes and all stored configuration
	DeleteDevice(ctx context.Context, in *devicepb.DeleteParameters, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type fleetServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFleetServiceClient(cc grpc.ClientConnInterface) FleetServiceClient {
	return &fleetServiceClient{cc}
}

func (c *fleetServiceClient) DiscoverDevice(ctx context.Context, in *DiscoverDeviceParameters, opts ...grpc.CallOption) (*devicepb.Device, error) {
	out := new(devicepb.Device)
	err := c.cc.Invoke(ctx, "/fleet.FleetService/DiscoverDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fleetServiceClient) CollectDevice(ctx context.Context, in *CollectDeviceParameters, opts ...grpc.CallOption) (*devicepb.Device, error) {
	out := new(devicepb.Device)
	err := c.cc.Invoke(ctx, "/fleet.FleetService/CollectDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fleetServiceClient) CollectConfig(ctx context.Context, in *CollectConfigParameters, opts ...grpc.CallOption) (*configurationpb.Configuration, error) {
	out := new(configurationpb.Configuration)
	err := c.cc.Invoke(ctx, "/fleet.FleetService/CollectConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fleetServiceClient) DeleteDevice(ctx context.Context, in *devicepb.DeleteParameters, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/fleet.FleetService/DeleteDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FleetServiceServer is the server API for FleetService service.
// All implementations must embed UnimplementedFleetServiceServer
// for forward compatibility
type FleetServiceServer interface {
	// DiscoverDevice discovers a device in the network, creates the device with the information provided
	// by the poller in the discovery. (e.g. sysname, ip address, mac address, etc)
	// NOTE: hostname OR management_ip must be set
	// Creating a device will append the default schedles to the device.
	// - CollectDevice every hour
	// - CollectConfig every 24 hours
	DiscoverDevice(context.Context, *DiscoverDeviceParameters) (*devicepb.Device, error)
	// CollectDevice collects information about the device from the network (with the help of the poller)
	// and returns the device with the updated information
	CollectDevice(context.Context, *CollectDeviceParameters) (*devicepb.Device, error)
	// CollectConfig collects the running configuration from the device in the network (with the help of the poller) and
	// returns the config as a string
	CollectConfig(context.Context, *CollectConfigParameters) (*configurationpb.Configuration, error)
	// Delete a device, chagnes and all stored configuration
	DeleteDevice(context.Context, *devicepb.DeleteParameters) (*emptypb.Empty, error)
	mustEmbedUnimplementedFleetServiceServer()
}

// UnimplementedFleetServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFleetServiceServer struct {
}

func (UnimplementedFleetServiceServer) DiscoverDevice(context.Context, *DiscoverDeviceParameters) (*devicepb.Device, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DiscoverDevice not implemented")
}
func (UnimplementedFleetServiceServer) CollectDevice(context.Context, *CollectDeviceParameters) (*devicepb.Device, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectDevice not implemented")
}
func (UnimplementedFleetServiceServer) CollectConfig(context.Context, *CollectConfigParameters) (*configurationpb.Configuration, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectConfig not implemented")
}
func (UnimplementedFleetServiceServer) DeleteDevice(context.Context, *devicepb.DeleteParameters) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDevice not implemented")
}
func (UnimplementedFleetServiceServer) mustEmbedUnimplementedFleetServiceServer() {}

// UnsafeFleetServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FleetServiceServer will
// result in compilation errors.
type UnsafeFleetServiceServer interface {
	mustEmbedUnimplementedFleetServiceServer()
}

func RegisterFleetServiceServer(s grpc.ServiceRegistrar, srv FleetServiceServer) {
	s.RegisterService(&FleetService_ServiceDesc, srv)
}

func _FleetService_DiscoverDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiscoverDeviceParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FleetServiceServer).DiscoverDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.FleetService/DiscoverDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FleetServiceServer).DiscoverDevice(ctx, req.(*DiscoverDeviceParameters))
	}
	return interceptor(ctx, in, info, handler)
}

func _FleetService_CollectDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectDeviceParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FleetServiceServer).CollectDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.FleetService/CollectDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FleetServiceServer).CollectDevice(ctx, req.(*CollectDeviceParameters))
	}
	return interceptor(ctx, in, info, handler)
}

func _FleetService_CollectConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectConfigParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FleetServiceServer).CollectConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.FleetService/CollectConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FleetServiceServer).CollectConfig(ctx, req.(*CollectConfigParameters))
	}
	return interceptor(ctx, in, info, handler)
}

func _FleetService_DeleteDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(devicepb.DeleteParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FleetServiceServer).DeleteDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.FleetService/DeleteDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FleetServiceServer).DeleteDevice(ctx, req.(*devicepb.DeleteParameters))
	}
	return interceptor(ctx, in, info, handler)
}

// FleetService_ServiceDesc is the grpc.ServiceDesc for FleetService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FleetService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fleet.FleetService",
	HandlerType: (*FleetServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DiscoverDevice",
			Handler:    _FleetService_DiscoverDevice_Handler,
		},
		{
			MethodName: "CollectDevice",
			Handler:    _FleetService_CollectDevice_Handler,
		},
		{
			MethodName: "CollectConfig",
			Handler:    _FleetService_CollectConfig_Handler,
		},
		{
			MethodName: "DeleteDevice",
			Handler:    _FleetService_DeleteDevice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fleet.proto",
}
