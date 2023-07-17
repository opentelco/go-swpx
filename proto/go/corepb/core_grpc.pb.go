// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.2
// source: core.proto

package corepb

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

// CoreServiceClient is the client API for CoreService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CoreServiceClient interface {
	// Discover is used to get basic information about an network element
	Discover(ctx context.Context, in *DiscoverRequest, opts ...grpc.CallOption) (*DiscoverResponse, error)
	// SWP Polling call to get technical Information and other information about a network element
	// the request is sent to the correct poller based on the network_region of the request
	// the type of the request is used to determine what information to collect from the network element
	Poll(ctx context.Context, in *PollRequest, opts ...grpc.CallOption) (*PollResponse, error)
	// Information returns infomration about the switch poller. loaded resources plugins and provider plugins
	Information(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*InformationResponse, error)
	// CollectConfig collects the configuration of a network element check for any changes between the stored config and the
	// collected one. Returs a list of changes and the config collected from the network element
	CollectConfig(ctx context.Context, in *CollectConfigRequest, opts ...grpc.CallOption) (*CollectConfigResponse, error)
}

type coreServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCoreServiceClient(cc grpc.ClientConnInterface) CoreServiceClient {
	return &coreServiceClient{cc}
}

func (c *coreServiceClient) Discover(ctx context.Context, in *DiscoverRequest, opts ...grpc.CallOption) (*DiscoverResponse, error) {
	out := new(DiscoverResponse)
	err := c.cc.Invoke(ctx, "/core.CoreService/Discover", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coreServiceClient) Poll(ctx context.Context, in *PollRequest, opts ...grpc.CallOption) (*PollResponse, error) {
	out := new(PollResponse)
	err := c.cc.Invoke(ctx, "/core.CoreService/Poll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coreServiceClient) Information(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*InformationResponse, error) {
	out := new(InformationResponse)
	err := c.cc.Invoke(ctx, "/core.CoreService/Information", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coreServiceClient) CollectConfig(ctx context.Context, in *CollectConfigRequest, opts ...grpc.CallOption) (*CollectConfigResponse, error) {
	out := new(CollectConfigResponse)
	err := c.cc.Invoke(ctx, "/core.CoreService/CollectConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CoreServiceServer is the server API for CoreService service.
// All implementations must embed UnimplementedCoreServiceServer
// for forward compatibility
type CoreServiceServer interface {
	// Discover is used to get basic information about an network element
	Discover(context.Context, *DiscoverRequest) (*DiscoverResponse, error)
	// SWP Polling call to get technical Information and other information about a network element
	// the request is sent to the correct poller based on the network_region of the request
	// the type of the request is used to determine what information to collect from the network element
	Poll(context.Context, *PollRequest) (*PollResponse, error)
	// Information returns infomration about the switch poller. loaded resources plugins and provider plugins
	Information(context.Context, *emptypb.Empty) (*InformationResponse, error)
	// CollectConfig collects the configuration of a network element check for any changes between the stored config and the
	// collected one. Returs a list of changes and the config collected from the network element
	CollectConfig(context.Context, *CollectConfigRequest) (*CollectConfigResponse, error)
	mustEmbedUnimplementedCoreServiceServer()
}

// UnimplementedCoreServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCoreServiceServer struct {
}

func (UnimplementedCoreServiceServer) Discover(context.Context, *DiscoverRequest) (*DiscoverResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Discover not implemented")
}
func (UnimplementedCoreServiceServer) Poll(context.Context, *PollRequest) (*PollResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Poll not implemented")
}
func (UnimplementedCoreServiceServer) Information(context.Context, *emptypb.Empty) (*InformationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Information not implemented")
}
func (UnimplementedCoreServiceServer) CollectConfig(context.Context, *CollectConfigRequest) (*CollectConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectConfig not implemented")
}
func (UnimplementedCoreServiceServer) mustEmbedUnimplementedCoreServiceServer() {}

// UnsafeCoreServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CoreServiceServer will
// result in compilation errors.
type UnsafeCoreServiceServer interface {
	mustEmbedUnimplementedCoreServiceServer()
}

func RegisterCoreServiceServer(s grpc.ServiceRegistrar, srv CoreServiceServer) {
	s.RegisterService(&CoreService_ServiceDesc, srv)
}

func _CoreService_Discover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiscoverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoreServiceServer).Discover(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.CoreService/Discover",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoreServiceServer).Discover(ctx, req.(*DiscoverRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CoreService_Poll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PollRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoreServiceServer).Poll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.CoreService/Poll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoreServiceServer).Poll(ctx, req.(*PollRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CoreService_Information_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoreServiceServer).Information(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.CoreService/Information",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoreServiceServer).Information(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CoreService_CollectConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoreServiceServer).CollectConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.CoreService/CollectConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoreServiceServer).CollectConfig(ctx, req.(*CollectConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CoreService_ServiceDesc is the grpc.ServiceDesc for CoreService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CoreService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "core.CoreService",
	HandlerType: (*CoreServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Discover",
			Handler:    _CoreService_Discover_Handler,
		},
		{
			MethodName: "Poll",
			Handler:    _CoreService_Poll_Handler,
		},
		{
			MethodName: "Information",
			Handler:    _CoreService_Information_Handler,
		},
		{
			MethodName: "CollectConfig",
			Handler:    _CoreService_CollectConfig_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "core.proto",
}

// CommanderServiceClient is the client API for CommanderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommanderServiceClient interface {
	// configure a configuration stanza on a network element
	ConfigureStanza(ctx context.Context, in *ConfigureStanzaRequest, opts ...grpc.CallOption) (*ConfigureStanzaResponse, error)
	TogglePort(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error)
	EnablePortFlapDetection(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error)
}

type commanderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCommanderServiceClient(cc grpc.ClientConnInterface) CommanderServiceClient {
	return &commanderServiceClient{cc}
}

func (c *commanderServiceClient) ConfigureStanza(ctx context.Context, in *ConfigureStanzaRequest, opts ...grpc.CallOption) (*ConfigureStanzaResponse, error) {
	out := new(ConfigureStanzaResponse)
	err := c.cc.Invoke(ctx, "/core.CommanderService/ConfigureStanza", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commanderServiceClient) TogglePort(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error) {
	out := new(CommandResponse)
	err := c.cc.Invoke(ctx, "/core.CommanderService/TogglePort", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commanderServiceClient) EnablePortFlapDetection(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (*CommandResponse, error) {
	out := new(CommandResponse)
	err := c.cc.Invoke(ctx, "/core.CommanderService/EnablePortFlapDetection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommanderServiceServer is the server API for CommanderService service.
// All implementations must embed UnimplementedCommanderServiceServer
// for forward compatibility
type CommanderServiceServer interface {
	// configure a configuration stanza on a network element
	ConfigureStanza(context.Context, *ConfigureStanzaRequest) (*ConfigureStanzaResponse, error)
	TogglePort(context.Context, *CommandRequest) (*CommandResponse, error)
	EnablePortFlapDetection(context.Context, *CommandRequest) (*CommandResponse, error)
	mustEmbedUnimplementedCommanderServiceServer()
}

// UnimplementedCommanderServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCommanderServiceServer struct {
}

func (UnimplementedCommanderServiceServer) ConfigureStanza(context.Context, *ConfigureStanzaRequest) (*ConfigureStanzaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigureStanza not implemented")
}
func (UnimplementedCommanderServiceServer) TogglePort(context.Context, *CommandRequest) (*CommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TogglePort not implemented")
}
func (UnimplementedCommanderServiceServer) EnablePortFlapDetection(context.Context, *CommandRequest) (*CommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnablePortFlapDetection not implemented")
}
func (UnimplementedCommanderServiceServer) mustEmbedUnimplementedCommanderServiceServer() {}

// UnsafeCommanderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommanderServiceServer will
// result in compilation errors.
type UnsafeCommanderServiceServer interface {
	mustEmbedUnimplementedCommanderServiceServer()
}

func RegisterCommanderServiceServer(s grpc.ServiceRegistrar, srv CommanderServiceServer) {
	s.RegisterService(&CommanderService_ServiceDesc, srv)
}

func _CommanderService_ConfigureStanza_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigureStanzaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommanderServiceServer).ConfigureStanza(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.CommanderService/ConfigureStanza",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommanderServiceServer).ConfigureStanza(ctx, req.(*ConfigureStanzaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommanderService_TogglePort_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommanderServiceServer).TogglePort(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.CommanderService/TogglePort",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommanderServiceServer).TogglePort(ctx, req.(*CommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommanderService_EnablePortFlapDetection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommanderServiceServer).EnablePortFlapDetection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.CommanderService/EnablePortFlapDetection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommanderServiceServer).EnablePortFlapDetection(ctx, req.(*CommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CommanderService_ServiceDesc is the grpc.ServiceDesc for CommanderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CommanderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "core.CommanderService",
	HandlerType: (*CommanderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConfigureStanza",
			Handler:    _CommanderService_ConfigureStanza_Handler,
		},
		{
			MethodName: "TogglePort",
			Handler:    _CommanderService_TogglePort_Handler,
		},
		{
			MethodName: "EnablePortFlapDetection",
			Handler:    _CommanderService_EnablePortFlapDetection_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "core.proto",
}

// ProviderServiceClient is the client API for ProviderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProviderServiceClient interface {
	// Ask provider to return a valid CPE for a access
	CPE(ctx context.Context, in *ProvideCPERequest, opts ...grpc.CallOption) (*ProvideCPEResponse, error)
	// Ask a provider to return information about a selected access
	Access(ctx context.Context, in *ProvideAccessRequest, opts ...grpc.CallOption) (*ProvideAccessResponse, error)
}

type providerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProviderServiceClient(cc grpc.ClientConnInterface) ProviderServiceClient {
	return &providerServiceClient{cc}
}

func (c *providerServiceClient) CPE(ctx context.Context, in *ProvideCPERequest, opts ...grpc.CallOption) (*ProvideCPEResponse, error) {
	out := new(ProvideCPEResponse)
	err := c.cc.Invoke(ctx, "/core.ProviderService/CPE", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerServiceClient) Access(ctx context.Context, in *ProvideAccessRequest, opts ...grpc.CallOption) (*ProvideAccessResponse, error) {
	out := new(ProvideAccessResponse)
	err := c.cc.Invoke(ctx, "/core.ProviderService/Access", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProviderServiceServer is the server API for ProviderService service.
// All implementations must embed UnimplementedProviderServiceServer
// for forward compatibility
type ProviderServiceServer interface {
	// Ask provider to return a valid CPE for a access
	CPE(context.Context, *ProvideCPERequest) (*ProvideCPEResponse, error)
	// Ask a provider to return information about a selected access
	Access(context.Context, *ProvideAccessRequest) (*ProvideAccessResponse, error)
	mustEmbedUnimplementedProviderServiceServer()
}

// UnimplementedProviderServiceServer must be embedded to have forward compatible implementations.
type UnimplementedProviderServiceServer struct {
}

func (UnimplementedProviderServiceServer) CPE(context.Context, *ProvideCPERequest) (*ProvideCPEResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CPE not implemented")
}
func (UnimplementedProviderServiceServer) Access(context.Context, *ProvideAccessRequest) (*ProvideAccessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Access not implemented")
}
func (UnimplementedProviderServiceServer) mustEmbedUnimplementedProviderServiceServer() {}

// UnsafeProviderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProviderServiceServer will
// result in compilation errors.
type UnsafeProviderServiceServer interface {
	mustEmbedUnimplementedProviderServiceServer()
}

func RegisterProviderServiceServer(s grpc.ServiceRegistrar, srv ProviderServiceServer) {
	s.RegisterService(&ProviderService_ServiceDesc, srv)
}

func _ProviderService_CPE_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProvideCPERequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderServiceServer).CPE(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.ProviderService/CPE",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderServiceServer).CPE(ctx, req.(*ProvideCPERequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProviderService_Access_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProvideAccessRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderServiceServer).Access(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.ProviderService/Access",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderServiceServer).Access(ctx, req.(*ProvideAccessRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProviderService_ServiceDesc is the grpc.ServiceDesc for ProviderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProviderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "core.ProviderService",
	HandlerType: (*ProviderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CPE",
			Handler:    _ProviderService_CPE_Handler,
		},
		{
			MethodName: "Access",
			Handler:    _ProviderService_Access_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "core.proto",
}
