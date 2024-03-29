//
// Copyright (c) 2023. Liero AB
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the Software
// is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
// CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: plugin_provider.proto

package providerpb

import (
	context "context"
	corepb "go.opentelco.io/go-swpx/proto/go/corepb"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Provider_Name_FullMethodName                  = "/provider.Provider/Name"
	Provider_Version_FullMethodName               = "/provider.Provider/Version"
	Provider_ResolveSessionRequest_FullMethodName = "/provider.Provider/ResolveSessionRequest"
	Provider_ResolveResourcePlugin_FullMethodName = "/provider.Provider/ResolveResourcePlugin"
	Provider_ProcessPollResponse_FullMethodName   = "/provider.Provider/ProcessPollResponse"
)

// ProviderClient is the client API for Provider service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProviderClient interface {
	Name(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*NameResponse, error)
	Version(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*VersionResponse, error)
	// PRE.1 Always called first in the chain of Provider RPCs
	// Resolve any hostname and or port in the session request
	ResolveSessionRequest(ctx context.Context, in *corepb.SessionRequest, opts ...grpc.CallOption) (*corepb.SessionRequest, error)
	// PRE.2 Called second in the chain of Provider RPCs
	// From the resolved session request, resolve the resource plugin to be used
	// This is only called if the settings.resource_plugin is empty
	ResolveResourcePlugin(ctx context.Context, in *corepb.SessionRequest, opts ...grpc.CallOption) (*ResolveResourcePluginResponse, error)
	// POST.1 Called in the end after returning any response to the client
	// Process the Poll response with the provider's own logic
	ProcessPollResponse(ctx context.Context, in *corepb.PollResponse, opts ...grpc.CallOption) (*corepb.PollResponse, error)
}

type providerClient struct {
	cc grpc.ClientConnInterface
}

func NewProviderClient(cc grpc.ClientConnInterface) ProviderClient {
	return &providerClient{cc}
}

func (c *providerClient) Name(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*NameResponse, error) {
	out := new(NameResponse)
	err := c.cc.Invoke(ctx, Provider_Name_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerClient) Version(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*VersionResponse, error) {
	out := new(VersionResponse)
	err := c.cc.Invoke(ctx, Provider_Version_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerClient) ResolveSessionRequest(ctx context.Context, in *corepb.SessionRequest, opts ...grpc.CallOption) (*corepb.SessionRequest, error) {
	out := new(corepb.SessionRequest)
	err := c.cc.Invoke(ctx, Provider_ResolveSessionRequest_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerClient) ResolveResourcePlugin(ctx context.Context, in *corepb.SessionRequest, opts ...grpc.CallOption) (*ResolveResourcePluginResponse, error) {
	out := new(ResolveResourcePluginResponse)
	err := c.cc.Invoke(ctx, Provider_ResolveResourcePlugin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerClient) ProcessPollResponse(ctx context.Context, in *corepb.PollResponse, opts ...grpc.CallOption) (*corepb.PollResponse, error) {
	out := new(corepb.PollResponse)
	err := c.cc.Invoke(ctx, Provider_ProcessPollResponse_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProviderServer is the server API for Provider service.
// All implementations must embed UnimplementedProviderServer
// for forward compatibility
type ProviderServer interface {
	Name(context.Context, *emptypb.Empty) (*NameResponse, error)
	Version(context.Context, *emptypb.Empty) (*VersionResponse, error)
	// PRE.1 Always called first in the chain of Provider RPCs
	// Resolve any hostname and or port in the session request
	ResolveSessionRequest(context.Context, *corepb.SessionRequest) (*corepb.SessionRequest, error)
	// PRE.2 Called second in the chain of Provider RPCs
	// From the resolved session request, resolve the resource plugin to be used
	// This is only called if the settings.resource_plugin is empty
	ResolveResourcePlugin(context.Context, *corepb.SessionRequest) (*ResolveResourcePluginResponse, error)
	// POST.1 Called in the end after returning any response to the client
	// Process the Poll response with the provider's own logic
	ProcessPollResponse(context.Context, *corepb.PollResponse) (*corepb.PollResponse, error)
	mustEmbedUnimplementedProviderServer()
}

// UnimplementedProviderServer must be embedded to have forward compatible implementations.
type UnimplementedProviderServer struct {
}

func (UnimplementedProviderServer) Name(context.Context, *emptypb.Empty) (*NameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Name not implemented")
}
func (UnimplementedProviderServer) Version(context.Context, *emptypb.Empty) (*VersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Version not implemented")
}
func (UnimplementedProviderServer) ResolveSessionRequest(context.Context, *corepb.SessionRequest) (*corepb.SessionRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResolveSessionRequest not implemented")
}
func (UnimplementedProviderServer) ResolveResourcePlugin(context.Context, *corepb.SessionRequest) (*ResolveResourcePluginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResolveResourcePlugin not implemented")
}
func (UnimplementedProviderServer) ProcessPollResponse(context.Context, *corepb.PollResponse) (*corepb.PollResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessPollResponse not implemented")
}
func (UnimplementedProviderServer) mustEmbedUnimplementedProviderServer() {}

// UnsafeProviderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProviderServer will
// result in compilation errors.
type UnsafeProviderServer interface {
	mustEmbedUnimplementedProviderServer()
}

func RegisterProviderServer(s grpc.ServiceRegistrar, srv ProviderServer) {
	s.RegisterService(&Provider_ServiceDesc, srv)
}

func _Provider_Name_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderServer).Name(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Provider_Name_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderServer).Name(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Provider_Version_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderServer).Version(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Provider_Version_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderServer).Version(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Provider_ResolveSessionRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(corepb.SessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderServer).ResolveSessionRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Provider_ResolveSessionRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderServer).ResolveSessionRequest(ctx, req.(*corepb.SessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Provider_ResolveResourcePlugin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(corepb.SessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderServer).ResolveResourcePlugin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Provider_ResolveResourcePlugin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderServer).ResolveResourcePlugin(ctx, req.(*corepb.SessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Provider_ProcessPollResponse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(corepb.PollResponse)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderServer).ProcessPollResponse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Provider_ProcessPollResponse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderServer).ProcessPollResponse(ctx, req.(*corepb.PollResponse))
	}
	return interceptor(ctx, in, info, handler)
}

// Provider_ServiceDesc is the grpc.ServiceDesc for Provider service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Provider_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "provider.Provider",
	HandlerType: (*ProviderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Name",
			Handler:    _Provider_Name_Handler,
		},
		{
			MethodName: "Version",
			Handler:    _Provider_Version_Handler,
		},
		{
			MethodName: "ResolveSessionRequest",
			Handler:    _Provider_ResolveSessionRequest_Handler,
		},
		{
			MethodName: "ResolveResourcePlugin",
			Handler:    _Provider_ResolveResourcePlugin_Handler,
		},
		{
			MethodName: "ProcessPollResponse",
			Handler:    _Provider_ProcessPollResponse_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugin_provider.proto",
}
