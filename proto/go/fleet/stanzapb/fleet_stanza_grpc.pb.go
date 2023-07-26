// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.2
// source: fleet_stanza.proto

package stanzapb

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

// StanzaServiceClient is the client API for StanzaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StanzaServiceClient interface {
	// Create a new stanza and return it
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*Stanza, error)
	// Clone a stanza and attach it to a device, this is used on Apply to create a new stanza that is applied to a device
	Clone(ctx context.Context, in *CloneRequest, opts ...grpc.CallOption) (*Stanza, error)
	// Get a stanza by id and return it
	GetByID(ctx context.Context, in *GetByIDRequest, opts ...grpc.CallOption) (*Stanza, error)
	// List stanzas, if no filters are used the list will return all stanzas not applied to a device
	// basically returning the library of stanzas that can be applied to a device.
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	// Update a stanza that is not applied to a device yet. If the stanza is applied to a device, it cannot be updated.
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*Stanza, error)
	// Delete a stanza that is not yet applied, Delete from the stanza library
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Apply a stanza to a device, this will duplicate the stanza and apply it to the device and return the applied stanza
	// if the Apply fails the stanza will be reverted by using the revert_content in the stanza. If no revert_content is set
	// the stanza will not be reverted and apply will return an error.
	Apply(ctx context.Context, in *ApplyRequest, opts ...grpc.CallOption) (*ApplyResponse, error)
	// Revert a stanza that has been applied to a device, this will use the revert_content in the stanza to revert the configuration
	// written to the device. If no revert_content is set the stanza will not be reverted and revert will return an error and the stanza
	Revert(ctx context.Context, in *RevertRequest, opts ...grpc.CallOption) (*RevertResponse, error)
}

type stanzaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStanzaServiceClient(cc grpc.ClientConnInterface) StanzaServiceClient {
	return &stanzaServiceClient{cc}
}

func (c *stanzaServiceClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*Stanza, error) {
	out := new(Stanza)
	err := c.cc.Invoke(ctx, "/fleet.stanza.StanzaService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stanzaServiceClient) Clone(ctx context.Context, in *CloneRequest, opts ...grpc.CallOption) (*Stanza, error) {
	out := new(Stanza)
	err := c.cc.Invoke(ctx, "/fleet.stanza.StanzaService/Clone", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stanzaServiceClient) GetByID(ctx context.Context, in *GetByIDRequest, opts ...grpc.CallOption) (*Stanza, error) {
	out := new(Stanza)
	err := c.cc.Invoke(ctx, "/fleet.stanza.StanzaService/GetByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stanzaServiceClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/fleet.stanza.StanzaService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stanzaServiceClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*Stanza, error) {
	out := new(Stanza)
	err := c.cc.Invoke(ctx, "/fleet.stanza.StanzaService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stanzaServiceClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/fleet.stanza.StanzaService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stanzaServiceClient) Apply(ctx context.Context, in *ApplyRequest, opts ...grpc.CallOption) (*ApplyResponse, error) {
	out := new(ApplyResponse)
	err := c.cc.Invoke(ctx, "/fleet.stanza.StanzaService/Apply", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stanzaServiceClient) Revert(ctx context.Context, in *RevertRequest, opts ...grpc.CallOption) (*RevertResponse, error) {
	out := new(RevertResponse)
	err := c.cc.Invoke(ctx, "/fleet.stanza.StanzaService/Revert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StanzaServiceServer is the server API for StanzaService service.
// All implementations must embed UnimplementedStanzaServiceServer
// for forward compatibility
type StanzaServiceServer interface {
	// Create a new stanza and return it
	Create(context.Context, *CreateRequest) (*Stanza, error)
	// Clone a stanza and attach it to a device, this is used on Apply to create a new stanza that is applied to a device
	Clone(context.Context, *CloneRequest) (*Stanza, error)
	// Get a stanza by id and return it
	GetByID(context.Context, *GetByIDRequest) (*Stanza, error)
	// List stanzas, if no filters are used the list will return all stanzas not applied to a device
	// basically returning the library of stanzas that can be applied to a device.
	List(context.Context, *ListRequest) (*ListResponse, error)
	// Update a stanza that is not applied to a device yet. If the stanza is applied to a device, it cannot be updated.
	Update(context.Context, *UpdateRequest) (*Stanza, error)
	// Delete a stanza that is not yet applied, Delete from the stanza library
	Delete(context.Context, *DeleteRequest) (*emptypb.Empty, error)
	// Apply a stanza to a device, this will duplicate the stanza and apply it to the device and return the applied stanza
	// if the Apply fails the stanza will be reverted by using the revert_content in the stanza. If no revert_content is set
	// the stanza will not be reverted and apply will return an error.
	Apply(context.Context, *ApplyRequest) (*ApplyResponse, error)
	// Revert a stanza that has been applied to a device, this will use the revert_content in the stanza to revert the configuration
	// written to the device. If no revert_content is set the stanza will not be reverted and revert will return an error and the stanza
	Revert(context.Context, *RevertRequest) (*RevertResponse, error)
	mustEmbedUnimplementedStanzaServiceServer()
}

// UnimplementedStanzaServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStanzaServiceServer struct {
}

func (UnimplementedStanzaServiceServer) Create(context.Context, *CreateRequest) (*Stanza, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedStanzaServiceServer) Clone(context.Context, *CloneRequest) (*Stanza, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Clone not implemented")
}
func (UnimplementedStanzaServiceServer) GetByID(context.Context, *GetByIDRequest) (*Stanza, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}
func (UnimplementedStanzaServiceServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedStanzaServiceServer) Update(context.Context, *UpdateRequest) (*Stanza, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedStanzaServiceServer) Delete(context.Context, *DeleteRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedStanzaServiceServer) Apply(context.Context, *ApplyRequest) (*ApplyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Apply not implemented")
}
func (UnimplementedStanzaServiceServer) Revert(context.Context, *RevertRequest) (*RevertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Revert not implemented")
}
func (UnimplementedStanzaServiceServer) mustEmbedUnimplementedStanzaServiceServer() {}

// UnsafeStanzaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StanzaServiceServer will
// result in compilation errors.
type UnsafeStanzaServiceServer interface {
	mustEmbedUnimplementedStanzaServiceServer()
}

func RegisterStanzaServiceServer(s grpc.ServiceRegistrar, srv StanzaServiceServer) {
	s.RegisterService(&StanzaService_ServiceDesc, srv)
}

func _StanzaService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StanzaServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.stanza.StanzaService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StanzaServiceServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StanzaService_Clone_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StanzaServiceServer).Clone(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.stanza.StanzaService/Clone",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StanzaServiceServer).Clone(ctx, req.(*CloneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StanzaService_GetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StanzaServiceServer).GetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.stanza.StanzaService/GetByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StanzaServiceServer).GetByID(ctx, req.(*GetByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StanzaService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StanzaServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.stanza.StanzaService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StanzaServiceServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StanzaService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StanzaServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.stanza.StanzaService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StanzaServiceServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StanzaService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StanzaServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.stanza.StanzaService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StanzaServiceServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StanzaService_Apply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StanzaServiceServer).Apply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.stanza.StanzaService/Apply",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StanzaServiceServer).Apply(ctx, req.(*ApplyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StanzaService_Revert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StanzaServiceServer).Revert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fleet.stanza.StanzaService/Revert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StanzaServiceServer).Revert(ctx, req.(*RevertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// StanzaService_ServiceDesc is the grpc.ServiceDesc for StanzaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StanzaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fleet.stanza.StanzaService",
	HandlerType: (*StanzaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _StanzaService_Create_Handler,
		},
		{
			MethodName: "Clone",
			Handler:    _StanzaService_Clone_Handler,
		},
		{
			MethodName: "GetByID",
			Handler:    _StanzaService_GetByID_Handler,
		},
		{
			MethodName: "List",
			Handler:    _StanzaService_List_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _StanzaService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _StanzaService_Delete_Handler,
		},
		{
			MethodName: "Apply",
			Handler:    _StanzaService_Apply_Handler,
		},
		{
			MethodName: "Revert",
			Handler:    _StanzaService_Revert_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fleet_stanza.proto",
}
