// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package resource

import (
	context "context"
	networkelement "git.liero.se/opentelco/go-swpx/proto/go/networkelement"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ResourceClient is the client API for Resource service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ResourceClient interface {
	// Get the version of the network element
	Version(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*VersionResponse, error)
	// Get technical information about a port
	TechnicalPortInformation(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*networkelement.Element, error)
	// Get technical information about all ports TODO: rename
	AllPortInformation(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*networkelement.Element, error)
	// Map the interfaces with ifIndex and description
	MapInterface(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*NetworkElementInterfaces, error)
	// Map the interace description and the environemnt index
	MapEntityPhysical(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*NetworkElementInterfaces, error)
	// Get transceiver information about a interface
	GetTransceiverInformation(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*networkelement.Transceiver, error)
	// Get transceiver information about all interfaces
	GetAllTransceiverInformation(ctx context.Context, in *NetworkElementWrapper, opts ...grpc.CallOption) (*networkelement.Element, error)
}

type resourceClient struct {
	cc grpc.ClientConnInterface
}

func NewResourceClient(cc grpc.ClientConnInterface) ResourceClient {
	return &resourceClient{cc}
}

func (c *resourceClient) Version(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*VersionResponse, error) {
	out := new(VersionResponse)
	err := c.cc.Invoke(ctx, "/resource.Resource/Version", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) TechnicalPortInformation(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*networkelement.Element, error) {
	out := new(networkelement.Element)
	err := c.cc.Invoke(ctx, "/resource.Resource/TechnicalPortInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) AllPortInformation(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*networkelement.Element, error) {
	out := new(networkelement.Element)
	err := c.cc.Invoke(ctx, "/resource.Resource/AllPortInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) MapInterface(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*NetworkElementInterfaces, error) {
	out := new(NetworkElementInterfaces)
	err := c.cc.Invoke(ctx, "/resource.Resource/MapInterface", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) MapEntityPhysical(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*NetworkElementInterfaces, error) {
	out := new(NetworkElementInterfaces)
	err := c.cc.Invoke(ctx, "/resource.Resource/MapEntityPhysical", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) GetTransceiverInformation(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*networkelement.Transceiver, error) {
	out := new(networkelement.Transceiver)
	err := c.cc.Invoke(ctx, "/resource.Resource/GetTransceiverInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) GetAllTransceiverInformation(ctx context.Context, in *NetworkElementWrapper, opts ...grpc.CallOption) (*networkelement.Element, error) {
	out := new(networkelement.Element)
	err := c.cc.Invoke(ctx, "/resource.Resource/GetAllTransceiverInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceServer is the server API for Resource service.
// All implementations must embed UnimplementedResourceServer
// for forward compatibility
type ResourceServer interface {
	// Get the version of the network element
	Version(context.Context, *emptypb.Empty) (*VersionResponse, error)
	// Get technical information about a port
	TechnicalPortInformation(context.Context, *NetworkElement) (*networkelement.Element, error)
	// Get technical information about all ports TODO: rename
	AllPortInformation(context.Context, *NetworkElement) (*networkelement.Element, error)
	// Map the interfaces with ifIndex and description
	MapInterface(context.Context, *NetworkElement) (*NetworkElementInterfaces, error)
	// Map the interace description and the environemnt index
	MapEntityPhysical(context.Context, *NetworkElement) (*NetworkElementInterfaces, error)
	// Get transceiver information about a interface
	GetTransceiverInformation(context.Context, *NetworkElement) (*networkelement.Transceiver, error)
	// Get transceiver information about all interfaces
	GetAllTransceiverInformation(context.Context, *NetworkElementWrapper) (*networkelement.Element, error)
	mustEmbedUnimplementedResourceServer()
}

// UnimplementedResourceServer must be embedded to have forward compatible implementations.
type UnimplementedResourceServer struct {
}

func (UnimplementedResourceServer) Version(context.Context, *emptypb.Empty) (*VersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Version not implemented")
}
func (UnimplementedResourceServer) TechnicalPortInformation(context.Context, *NetworkElement) (*networkelement.Element, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TechnicalPortInformation not implemented")
}
func (UnimplementedResourceServer) AllPortInformation(context.Context, *NetworkElement) (*networkelement.Element, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllPortInformation not implemented")
}
func (UnimplementedResourceServer) MapInterface(context.Context, *NetworkElement) (*NetworkElementInterfaces, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MapInterface not implemented")
}
func (UnimplementedResourceServer) MapEntityPhysical(context.Context, *NetworkElement) (*NetworkElementInterfaces, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MapEntityPhysical not implemented")
}
func (UnimplementedResourceServer) GetTransceiverInformation(context.Context, *NetworkElement) (*networkelement.Transceiver, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransceiverInformation not implemented")
}
func (UnimplementedResourceServer) GetAllTransceiverInformation(context.Context, *NetworkElementWrapper) (*networkelement.Element, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllTransceiverInformation not implemented")
}
func (UnimplementedResourceServer) mustEmbedUnimplementedResourceServer() {}

// UnsafeResourceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ResourceServer will
// result in compilation errors.
type UnsafeResourceServer interface {
	mustEmbedUnimplementedResourceServer()
}

func RegisterResourceServer(s grpc.ServiceRegistrar, srv ResourceServer) {
	s.RegisterService(&Resource_ServiceDesc, srv)
}

func _Resource_Version_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).Version(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/resource.Resource/Version",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).Version(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_TechnicalPortInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkElement)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).TechnicalPortInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/resource.Resource/TechnicalPortInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).TechnicalPortInformation(ctx, req.(*NetworkElement))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_AllPortInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkElement)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).AllPortInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/resource.Resource/AllPortInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).AllPortInformation(ctx, req.(*NetworkElement))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_MapInterface_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkElement)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).MapInterface(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/resource.Resource/MapInterface",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).MapInterface(ctx, req.(*NetworkElement))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_MapEntityPhysical_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkElement)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).MapEntityPhysical(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/resource.Resource/MapEntityPhysical",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).MapEntityPhysical(ctx, req.(*NetworkElement))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_GetTransceiverInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkElement)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).GetTransceiverInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/resource.Resource/GetTransceiverInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).GetTransceiverInformation(ctx, req.(*NetworkElement))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_GetAllTransceiverInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkElementWrapper)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).GetAllTransceiverInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/resource.Resource/GetAllTransceiverInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).GetAllTransceiverInformation(ctx, req.(*NetworkElementWrapper))
	}
	return interceptor(ctx, in, info, handler)
}

// Resource_ServiceDesc is the grpc.ServiceDesc for Resource service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Resource_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "resource.Resource",
	HandlerType: (*ResourceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Version",
			Handler:    _Resource_Version_Handler,
		},
		{
			MethodName: "TechnicalPortInformation",
			Handler:    _Resource_TechnicalPortInformation_Handler,
		},
		{
			MethodName: "AllPortInformation",
			Handler:    _Resource_AllPortInformation_Handler,
		},
		{
			MethodName: "MapInterface",
			Handler:    _Resource_MapInterface_Handler,
		},
		{
			MethodName: "MapEntityPhysical",
			Handler:    _Resource_MapEntityPhysical_Handler,
		},
		{
			MethodName: "GetTransceiverInformation",
			Handler:    _Resource_GetTransceiverInformation_Handler,
		},
		{
			MethodName: "GetAllTransceiverInformation",
			Handler:    _Resource_GetAllTransceiverInformation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "resource.proto",
}