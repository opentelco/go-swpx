// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: fleet_device.proto

package devicepb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DeviceServiceClient is the client API for DeviceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeviceServiceClient interface {
	// *** Device ***
	// Get a device by its ID, this is used to get a specific device
	GetByID(ctx context.Context, in *GetByIDParameters, opts ...grpc.CallOption) (*Device, error)
	// Get a device by its hostname, managment ip or serial number etc (used to search for a device)
	List(ctx context.Context, in *ListParameters, opts ...grpc.CallOption) (*ListResponse, error)
	// Create a device in the fleet
	// note: if device needs to be discovered use the FleetService instead
	// Creating a device will append the default schedles to the device.
	// - CollectDevice every hour
	// - CollectConfig every 24 hours
	Create(ctx context.Context, in *CreateParameters, opts ...grpc.CallOption) (*Device, error)
	// Update a device in the fleet (this is used to update the device with new information)
	Update(ctx context.Context, in *UpdateParameters, opts ...grpc.CallOption) (*Device, error)
	// Delete a device from the fleet and all changes for the device
	// To purge a device use the Delete in the FleetService instead (as it also also deletes the configuration)
	Delete(ctx context.Context, in *DeleteParameters, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Get changes for a device, changes are created when a device is updated
	GetChangeByID(ctx context.Context, in *GetChangeByIDParameters, opts ...grpc.CallOption) (*Change, error)
	// returns a list of changes (default 100)
	ListChanges(ctx context.Context, in *ListChangesParameters, opts ...grpc.CallOption) (*ListChangesResponse, error)
	// add an event
	AddEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*Event, error)
	// Get changes for a device, changes are created when a device is updated
	GetEventByID(ctx context.Context, in *GetEventByIDParameters, opts ...grpc.CallOption) (*Event, error)
	// returns a list of events (default 100)
	ListEvents(ctx context.Context, in *ListEventsParameters, opts ...grpc.CallOption) (*ListEventsResponse, error)
}

type deviceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeviceServiceClient(cc grpc.ClientConnInterface) DeviceServiceClient {
	return &deviceServiceClient{cc}
}

func (c *deviceServiceClient) GetByID(ctx context.Context, in *GetByIDParameters, opts ...grpc.CallOption) (*Device, error) {
	out := new(Device)
	err := c.cc.Invoke(ctx, "/fleet.device.DeviceService/GetByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) List(ctx context.Context, in *ListParameters, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/fleet.device.DeviceService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) Create(ctx context.Context, in *CreateParameters, opts ...grpc.CallOption) (*Device, error) {
	out := new(Device)
	err := c.cc.Invoke(ctx, "/fleet.device.DeviceService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) Update(ctx context.Context, in *UpdateParameters, opts ...grpc.CallOption) (*Device, error) {
	out := new(Device)
	err := c.cc.Invoke(ctx, "/fleet.device.DeviceService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) Delete(ctx context.Context, in *DeleteParameters, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/fleet.device.DeviceService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) GetChangeByID(ctx context.Context, in *GetChangeByIDParameters, opts ...grpc.CallOption) (*Change, error) {
	out := new(Change)
	err := c.cc.Invoke(ctx, "/fleet.device.DeviceService/GetChangeByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) ListChanges(ctx context.Context, in *ListChangesParameters, opts ...grpc.CallOption) (*ListChangesResponse, error) {
	out := new(ListChangesResponse)
	err := c.cc.Invoke(ctx, "/fleet.device.DeviceService/ListChanges", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) AddEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/fleet.device.DeviceService/AddEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) GetEventByID(ctx context.Context, in *GetEventByIDParameters, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/fleet.device.DeviceService/GetEventByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) ListEvents(ctx context.Context, in *ListEventsParameters, opts ...grpc.CallOption) (*ListEventsResponse, error) {
	out := new(ListEventsResponse)
	err := c.cc.Invoke(ctx, "/fleet.device.DeviceService/ListEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceServiceServer is the server API for DeviceService service.
// All implementations must embed UnimplementedDeviceServiceServer
// for forward compatibility
type DeviceServiceServer interface {
	// *** Device ***
	// Get a device by its ID, this is used to get a specific device
	GetByID(context.Context, *GetByIDParameters) (*Device, error)
	// Get a device by its hostname, managment ip or serial number etc (used to search for a device)
	List(context.Context, *ListParameters) (*ListResponse, error)
	// Create a device in the fleet
	// note: if device needs to be discovered use the FleetService instead
	// Creating a device will append the default schedles to the device.
	// - CollectDevice every hour
	// - CollectConfig every 24 hours
	Create(context.Context, *CreateParameters) (*Device, error)
	// Update a device in the fleet (this is used to update the device with new information)
	Update(context.Context, *UpdateParameters) (*Device, error)
	// Delete a device from the fleet and all changes for the device
	// To purge a device use the Delete in the FleetService instead (as it also also deletes the configuration)
	Delete(context.Context, *DeleteParameters) (*emptypb.Empty, error)
	// Get changes for a device, changes are created when a device is updated
	GetChangeByID(context.Context, *GetChangeByIDParameters) (*Change, error)
	// returns a list of changes (default 100)
	ListChanges(context.Context, *ListChangesParameters) (*ListChangesResponse, error)
	// add an event
	AddEvent(context.Context, *Event) (*Event, error)
	// Get changes for a device, changes are created when a device is updated
	GetEventByID(context.Context, *GetEventByIDParameters) (*Event, error)
	// returns a list of events (default 100)
	ListEvents(context.Context, *ListEventsParameters) (*ListEventsResponse, error)
	mustEmbedUnimplementedDeviceServiceServer()
}

// UnimplementedDeviceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDeviceServiceServer struct {
}

func (UnimplementedDeviceServiceServer) GetByID(context.Context, *GetByIDParameters) (*Device, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}
func (UnimplementedDeviceServiceServer) List(context.Context, *ListParameters) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedDeviceServiceServer) Create(context.Context, *CreateParameters) (*Device, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedDeviceServiceServer) Update(context.Context, *UpdateParameters) (*Device, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedDeviceServiceServer) Delete(context.Context, *DeleteParameters) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedDeviceServiceServer) GetChangeByID(context.Context, *GetChangeByIDParameters) (*Change, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChangeByID not implemented")
}
func (UnimplementedDeviceServiceServer) ListChanges(context.Context, *ListChangesParameters) (*ListChangesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListChanges not implemented")
}
func (UnimplementedDeviceServiceServer) AddEvent(context.Context, *Event) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddEvent not implemented")
}
func (UnimplementedDeviceServiceServer) GetEventByID(context.Context, *GetEventByIDParameters) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventByID not implemented")
}
func (UnimplementedDeviceServiceServer) ListEvents(context.Context, *ListEventsParameters) (*ListEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEvents not implemented")
}
func (UnimplementedDeviceServiceServer) mustEmbedUnimplementedDeviceServiceServer() {}

// UnsafeDeviceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeviceServiceServer will
// result in compilation errors.
type UnsafeDeviceServiceServer interface {
	mustEmbedUnimplementedDeviceServiceServer()
}

func RegisterDeviceServiceServer(s grpc.ServiceRegistrar, srv DeviceServiceServer) {
	s.RegisterService(&DeviceService_ServiceDesc, srv)
}

func _DeviceService_GetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIDParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).GetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.device.DeviceService/GetByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).GetByID(ctx, req.(*GetByIDParameters))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.device.DeviceService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).List(ctx, req.(*ListParameters))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.device.DeviceService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).Create(ctx, req.(*CreateParameters))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.device.DeviceService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).Update(ctx, req.(*UpdateParameters))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.device.DeviceService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).Delete(ctx, req.(*DeleteParameters))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_GetChangeByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChangeByIDParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).GetChangeByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.device.DeviceService/GetChangeByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).GetChangeByID(ctx, req.(*GetChangeByIDParameters))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_ListChanges_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListChangesParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).ListChanges(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.device.DeviceService/ListChanges",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).ListChanges(ctx, req.(*ListChangesParameters))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_AddEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).AddEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.device.DeviceService/AddEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).AddEvent(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_GetEventByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventByIDParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).GetEventByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.device.DeviceService/GetEventByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).GetEventByID(ctx, req.(*GetEventByIDParameters))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_ListEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEventsParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).ListEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.device.DeviceService/ListEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).ListEvents(ctx, req.(*ListEventsParameters))
	}
	return interceptor(ctx, in, info, handler)
}

// DeviceService_ServiceDesc is the grpc.ServiceDesc for DeviceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeviceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fleet.device.DeviceService",
	HandlerType: (*DeviceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetByID",
			Handler:    _DeviceService_GetByID_Handler,
		},
		{
			MethodName: "List",
			Handler:    _DeviceService_List_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _DeviceService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _DeviceService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _DeviceService_Delete_Handler,
		},
		{
			MethodName: "GetChangeByID",
			Handler:    _DeviceService_GetChangeByID_Handler,
		},
		{
			MethodName: "ListChanges",
			Handler:    _DeviceService_ListChanges_Handler,
		},
		{
			MethodName: "AddEvent",
			Handler:    _DeviceService_AddEvent_Handler,
		},
		{
			MethodName: "GetEventByID",
			Handler:    _DeviceService_GetEventByID_Handler,
		},
		{
			MethodName: "ListEvents",
			Handler:    _DeviceService_ListEvents_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fleet_device.proto",
}
