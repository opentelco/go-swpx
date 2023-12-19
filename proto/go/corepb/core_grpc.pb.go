// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: core.proto

package corepb

import (
	context "context"
	analysispb "go.opentelco.io/go-swpx/proto/go/analysispb"
	stanzapb "go.opentelco.io/go-swpx/proto/go/stanzapb"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PollerClient is the client API for Poller service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PollerClient interface {
	// Discover is used to get basic information about an network element, used to make a quick check of the device
	// with a generic request
	Discover(ctx context.Context, in *DiscoverRequest, opts ...grpc.CallOption) (*DiscoverResponse, error)
	// CheckAvailability is used to check if a network element is available and responding to requests
	// this does not imply that the network element is working correctly or that it is configured correctly but
	// that it is responding to requests and that the poller can connect to it over SNMP/ICMP
	// the availability also verifys checking that hostname is resolvable (if hostname is used in the request)
	CheckAvailability(ctx context.Context, in *SessionRequest, opts ...grpc.CallOption) (*CheckAvailabilityResponse, error)
	// RunDiagnostic is used to run a diagnostic on a network element or a specific port on the network element
	// It will collect data from the network element and then wait for a period of time and collect data again
	// and return the difference between the two collections of data to the client. The data will also be analyzed
	// by the poller and Report of the diagnostic will be returned to the client.
	// the diagnostic will be run the number of times specified in the request and the time between each poll is 10 seconds.
	// connecting to a device can take up to one minute depending on the device and protocol used so a standard diagnostic
	// will take aproximately 1 minute to complete.
	RunDiagnostic(ctx context.Context, in *RunDiagnosticRequest, opts ...grpc.CallOption) (*RunDiagnosticResponse, error)
	// GetDiagnostic returns the report of a diagnostic that has been run on a network element or a specific port on the network element
	GetDiagnostic(ctx context.Context, in *GetDiagnosticRequest, opts ...grpc.CallOption) (*analysispb.Report, error)
	// ListDiagnostics returns a list of diagnostics that has been run on a network element or a specific port on the network element
	ListDiagnostics(ctx context.Context, in *ListDiagnosticsRequest, opts ...grpc.CallOption) (*ListDiagnosticsResponse, error)
	// GetDeviceInformation returns the technical information about a device
	// port etc is not considered in this request
	CollectDeviceInformation(ctx context.Context, in *CollectDeviceInformationRequest, opts ...grpc.CallOption) (*DeviceInformationResponse, error)
	// get basic information about a device
	// port etc is not considered in this request
	CollectBasicDeviceInformation(ctx context.Context, in *CollectBasicDeviceInformationRequest, opts ...grpc.CallOption) (*DeviceInformationResponse, error)
	// PortInformation returns information about a port on a device
	CollectPortInformation(ctx context.Context, in *CollectPortInformationRequest, opts ...grpc.CallOption) (*PortInformationResponse, error)
	// Get all basic information about a port on a device
	CollectBasicPortInformation(ctx context.Context, in *CollectBasicPortInformationRequest, opts ...grpc.CallOption) (*PortInformationResponse, error)
	// CollectConfig collects the configuration of a network element check for any changes between the stored config and the
	// collected one. Returs a list of changes and the config collected from the network element
	CollectConfig(ctx context.Context, in *CollectConfigRequest, opts ...grpc.CallOption) (*CollectConfigResponse, error)
	// deprecated, use specific RPC:s instead
	// SWP Polling call to get technical Information and other information about a network element
	// the request is sent to the correct poller based on the network_region of the request
	// the type of the request is used to determine what information to collect from the network element
	Poll(ctx context.Context, in *PollRequest, opts ...grpc.CallOption) (*PollResponse, error)
}

type pollerClient struct {
	cc grpc.ClientConnInterface
}

func NewPollerClient(cc grpc.ClientConnInterface) PollerClient {
	return &pollerClient{cc}
}

func (c *pollerClient) Discover(ctx context.Context, in *DiscoverRequest, opts ...grpc.CallOption) (*DiscoverResponse, error) {
	out := new(DiscoverResponse)
	err := c.cc.Invoke(ctx, "/core.Poller/Discover", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pollerClient) CheckAvailability(ctx context.Context, in *SessionRequest, opts ...grpc.CallOption) (*CheckAvailabilityResponse, error) {
	out := new(CheckAvailabilityResponse)
	err := c.cc.Invoke(ctx, "/core.Poller/CheckAvailability", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pollerClient) RunDiagnostic(ctx context.Context, in *RunDiagnosticRequest, opts ...grpc.CallOption) (*RunDiagnosticResponse, error) {
	out := new(RunDiagnosticResponse)
	err := c.cc.Invoke(ctx, "/core.Poller/RunDiagnostic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pollerClient) GetDiagnostic(ctx context.Context, in *GetDiagnosticRequest, opts ...grpc.CallOption) (*analysispb.Report, error) {
	out := new(analysispb.Report)
	err := c.cc.Invoke(ctx, "/core.Poller/GetDiagnostic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pollerClient) ListDiagnostics(ctx context.Context, in *ListDiagnosticsRequest, opts ...grpc.CallOption) (*ListDiagnosticsResponse, error) {
	out := new(ListDiagnosticsResponse)
	err := c.cc.Invoke(ctx, "/core.Poller/ListDiagnostics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pollerClient) CollectDeviceInformation(ctx context.Context, in *CollectDeviceInformationRequest, opts ...grpc.CallOption) (*DeviceInformationResponse, error) {
	out := new(DeviceInformationResponse)
	err := c.cc.Invoke(ctx, "/core.Poller/CollectDeviceInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pollerClient) CollectBasicDeviceInformation(ctx context.Context, in *CollectBasicDeviceInformationRequest, opts ...grpc.CallOption) (*DeviceInformationResponse, error) {
	out := new(DeviceInformationResponse)
	err := c.cc.Invoke(ctx, "/core.Poller/CollectBasicDeviceInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pollerClient) CollectPortInformation(ctx context.Context, in *CollectPortInformationRequest, opts ...grpc.CallOption) (*PortInformationResponse, error) {
	out := new(PortInformationResponse)
	err := c.cc.Invoke(ctx, "/core.Poller/CollectPortInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pollerClient) CollectBasicPortInformation(ctx context.Context, in *CollectBasicPortInformationRequest, opts ...grpc.CallOption) (*PortInformationResponse, error) {
	out := new(PortInformationResponse)
	err := c.cc.Invoke(ctx, "/core.Poller/CollectBasicPortInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pollerClient) CollectConfig(ctx context.Context, in *CollectConfigRequest, opts ...grpc.CallOption) (*CollectConfigResponse, error) {
	out := new(CollectConfigResponse)
	err := c.cc.Invoke(ctx, "/core.Poller/CollectConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pollerClient) Poll(ctx context.Context, in *PollRequest, opts ...grpc.CallOption) (*PollResponse, error) {
	out := new(PollResponse)
	err := c.cc.Invoke(ctx, "/core.Poller/Poll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PollerServer is the server API for Poller service.
// All implementations must embed UnimplementedPollerServer
// for forward compatibility
type PollerServer interface {
	// Discover is used to get basic information about an network element, used to make a quick check of the device
	// with a generic request
	Discover(context.Context, *DiscoverRequest) (*DiscoverResponse, error)
	// CheckAvailability is used to check if a network element is available and responding to requests
	// this does not imply that the network element is working correctly or that it is configured correctly but
	// that it is responding to requests and that the poller can connect to it over SNMP/ICMP
	// the availability also verifys checking that hostname is resolvable (if hostname is used in the request)
	CheckAvailability(context.Context, *SessionRequest) (*CheckAvailabilityResponse, error)
	// RunDiagnostic is used to run a diagnostic on a network element or a specific port on the network element
	// It will collect data from the network element and then wait for a period of time and collect data again
	// and return the difference between the two collections of data to the client. The data will also be analyzed
	// by the poller and Report of the diagnostic will be returned to the client.
	// the diagnostic will be run the number of times specified in the request and the time between each poll is 10 seconds.
	// connecting to a device can take up to one minute depending on the device and protocol used so a standard diagnostic
	// will take aproximately 1 minute to complete.
	RunDiagnostic(context.Context, *RunDiagnosticRequest) (*RunDiagnosticResponse, error)
	// GetDiagnostic returns the report of a diagnostic that has been run on a network element or a specific port on the network element
	GetDiagnostic(context.Context, *GetDiagnosticRequest) (*analysispb.Report, error)
	// ListDiagnostics returns a list of diagnostics that has been run on a network element or a specific port on the network element
	ListDiagnostics(context.Context, *ListDiagnosticsRequest) (*ListDiagnosticsResponse, error)
	// GetDeviceInformation returns the technical information about a device
	// port etc is not considered in this request
	CollectDeviceInformation(context.Context, *CollectDeviceInformationRequest) (*DeviceInformationResponse, error)
	// get basic information about a device
	// port etc is not considered in this request
	CollectBasicDeviceInformation(context.Context, *CollectBasicDeviceInformationRequest) (*DeviceInformationResponse, error)
	// PortInformation returns information about a port on a device
	CollectPortInformation(context.Context, *CollectPortInformationRequest) (*PortInformationResponse, error)
	// Get all basic information about a port on a device
	CollectBasicPortInformation(context.Context, *CollectBasicPortInformationRequest) (*PortInformationResponse, error)
	// CollectConfig collects the configuration of a network element check for any changes between the stored config and the
	// collected one. Returs a list of changes and the config collected from the network element
	CollectConfig(context.Context, *CollectConfigRequest) (*CollectConfigResponse, error)
	// deprecated, use specific RPC:s instead
	// SWP Polling call to get technical Information and other information about a network element
	// the request is sent to the correct poller based on the network_region of the request
	// the type of the request is used to determine what information to collect from the network element
	Poll(context.Context, *PollRequest) (*PollResponse, error)
	mustEmbedUnimplementedPollerServer()
}

// UnimplementedPollerServer must be embedded to have forward compatible implementations.
type UnimplementedPollerServer struct {
}

func (UnimplementedPollerServer) Discover(context.Context, *DiscoverRequest) (*DiscoverResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Discover not implemented")
}
func (UnimplementedPollerServer) CheckAvailability(context.Context, *SessionRequest) (*CheckAvailabilityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckAvailability not implemented")
}
func (UnimplementedPollerServer) RunDiagnostic(context.Context, *RunDiagnosticRequest) (*RunDiagnosticResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunDiagnostic not implemented")
}
func (UnimplementedPollerServer) GetDiagnostic(context.Context, *GetDiagnosticRequest) (*analysispb.Report, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDiagnostic not implemented")
}
func (UnimplementedPollerServer) ListDiagnostics(context.Context, *ListDiagnosticsRequest) (*ListDiagnosticsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDiagnostics not implemented")
}
func (UnimplementedPollerServer) CollectDeviceInformation(context.Context, *CollectDeviceInformationRequest) (*DeviceInformationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectDeviceInformation not implemented")
}
func (UnimplementedPollerServer) CollectBasicDeviceInformation(context.Context, *CollectBasicDeviceInformationRequest) (*DeviceInformationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectBasicDeviceInformation not implemented")
}
func (UnimplementedPollerServer) CollectPortInformation(context.Context, *CollectPortInformationRequest) (*PortInformationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectPortInformation not implemented")
}
func (UnimplementedPollerServer) CollectBasicPortInformation(context.Context, *CollectBasicPortInformationRequest) (*PortInformationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectBasicPortInformation not implemented")
}
func (UnimplementedPollerServer) CollectConfig(context.Context, *CollectConfigRequest) (*CollectConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectConfig not implemented")
}
func (UnimplementedPollerServer) Poll(context.Context, *PollRequest) (*PollResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Poll not implemented")
}
func (UnimplementedPollerServer) mustEmbedUnimplementedPollerServer() {}

// UnsafePollerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PollerServer will
// result in compilation errors.
type UnsafePollerServer interface {
	mustEmbedUnimplementedPollerServer()
}

func RegisterPollerServer(s grpc.ServiceRegistrar, srv PollerServer) {
	s.RegisterService(&Poller_ServiceDesc, srv)
}

func _Poller_Discover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiscoverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PollerServer).Discover(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Poller/Discover",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PollerServer).Discover(ctx, req.(*DiscoverRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Poller_CheckAvailability_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PollerServer).CheckAvailability(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Poller/CheckAvailability",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PollerServer).CheckAvailability(ctx, req.(*SessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Poller_RunDiagnostic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunDiagnosticRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PollerServer).RunDiagnostic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Poller/RunDiagnostic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PollerServer).RunDiagnostic(ctx, req.(*RunDiagnosticRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Poller_GetDiagnostic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDiagnosticRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PollerServer).GetDiagnostic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Poller/GetDiagnostic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PollerServer).GetDiagnostic(ctx, req.(*GetDiagnosticRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Poller_ListDiagnostics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDiagnosticsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PollerServer).ListDiagnostics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Poller/ListDiagnostics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PollerServer).ListDiagnostics(ctx, req.(*ListDiagnosticsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Poller_CollectDeviceInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectDeviceInformationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PollerServer).CollectDeviceInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Poller/CollectDeviceInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PollerServer).CollectDeviceInformation(ctx, req.(*CollectDeviceInformationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Poller_CollectBasicDeviceInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectBasicDeviceInformationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PollerServer).CollectBasicDeviceInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Poller/CollectBasicDeviceInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PollerServer).CollectBasicDeviceInformation(ctx, req.(*CollectBasicDeviceInformationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Poller_CollectPortInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectPortInformationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PollerServer).CollectPortInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Poller/CollectPortInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PollerServer).CollectPortInformation(ctx, req.(*CollectPortInformationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Poller_CollectBasicPortInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectBasicPortInformationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PollerServer).CollectBasicPortInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Poller/CollectBasicPortInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PollerServer).CollectBasicPortInformation(ctx, req.(*CollectBasicPortInformationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Poller_CollectConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PollerServer).CollectConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Poller/CollectConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PollerServer).CollectConfig(ctx, req.(*CollectConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Poller_Poll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PollRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PollerServer).Poll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Poller/Poll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PollerServer).Poll(ctx, req.(*PollRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Poller_ServiceDesc is the grpc.ServiceDesc for Poller service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Poller_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "core.Poller",
	HandlerType: (*PollerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Discover",
			Handler:    _Poller_Discover_Handler,
		},
		{
			MethodName: "CheckAvailability",
			Handler:    _Poller_CheckAvailability_Handler,
		},
		{
			MethodName: "RunDiagnostic",
			Handler:    _Poller_RunDiagnostic_Handler,
		},
		{
			MethodName: "GetDiagnostic",
			Handler:    _Poller_GetDiagnostic_Handler,
		},
		{
			MethodName: "ListDiagnostics",
			Handler:    _Poller_ListDiagnostics_Handler,
		},
		{
			MethodName: "CollectDeviceInformation",
			Handler:    _Poller_CollectDeviceInformation_Handler,
		},
		{
			MethodName: "CollectBasicDeviceInformation",
			Handler:    _Poller_CollectBasicDeviceInformation_Handler,
		},
		{
			MethodName: "CollectPortInformation",
			Handler:    _Poller_CollectPortInformation_Handler,
		},
		{
			MethodName: "CollectBasicPortInformation",
			Handler:    _Poller_CollectBasicPortInformation_Handler,
		},
		{
			MethodName: "CollectConfig",
			Handler:    _Poller_CollectConfig_Handler,
		},
		{
			MethodName: "Poll",
			Handler:    _Poller_Poll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "core.proto",
}

// CommanderClient is the client API for Commander service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommanderClient interface {
	// configure a configuration stanza on a network element
	ConfigureStanza(ctx context.Context, in *ConfigureStanzaRequest, opts ...grpc.CallOption) (*stanzapb.ConfigureResponse, error)
}

type commanderClient struct {
	cc grpc.ClientConnInterface
}

func NewCommanderClient(cc grpc.ClientConnInterface) CommanderClient {
	return &commanderClient{cc}
}

func (c *commanderClient) ConfigureStanza(ctx context.Context, in *ConfigureStanzaRequest, opts ...grpc.CallOption) (*stanzapb.ConfigureResponse, error) {
	out := new(stanzapb.ConfigureResponse)
	err := c.cc.Invoke(ctx, "/core.Commander/ConfigureStanza", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommanderServer is the server API for Commander service.
// All implementations must embed UnimplementedCommanderServer
// for forward compatibility
type CommanderServer interface {
	// configure a configuration stanza on a network element
	ConfigureStanza(context.Context, *ConfigureStanzaRequest) (*stanzapb.ConfigureResponse, error)
	mustEmbedUnimplementedCommanderServer()
}

// UnimplementedCommanderServer must be embedded to have forward compatible implementations.
type UnimplementedCommanderServer struct {
}

func (UnimplementedCommanderServer) ConfigureStanza(context.Context, *ConfigureStanzaRequest) (*stanzapb.ConfigureResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigureStanza not implemented")
}
func (UnimplementedCommanderServer) mustEmbedUnimplementedCommanderServer() {}

// UnsafeCommanderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommanderServer will
// result in compilation errors.
type UnsafeCommanderServer interface {
	mustEmbedUnimplementedCommanderServer()
}

func RegisterCommanderServer(s grpc.ServiceRegistrar, srv CommanderServer) {
	s.RegisterService(&Commander_ServiceDesc, srv)
}

func _Commander_ConfigureStanza_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigureStanzaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommanderServer).ConfigureStanza(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Commander/ConfigureStanza",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommanderServer).ConfigureStanza(ctx, req.(*ConfigureStanzaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Commander_ServiceDesc is the grpc.ServiceDesc for Commander service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Commander_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "core.Commander",
	HandlerType: (*CommanderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConfigureStanza",
			Handler:    _Commander_ConfigureStanza_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "core.proto",
}

// ProviderClient is the client API for Provider service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProviderClient interface {
	// Ask provider to return a valid CPE for a access
	CPE(ctx context.Context, in *ProvideCPERequest, opts ...grpc.CallOption) (*ProvideCPEResponse, error)
	// Ask a provider to return information about a selected access
	Access(ctx context.Context, in *ProvideAccessRequest, opts ...grpc.CallOption) (*ProvideAccessResponse, error)
}

type providerClient struct {
	cc grpc.ClientConnInterface
}

func NewProviderClient(cc grpc.ClientConnInterface) ProviderClient {
	return &providerClient{cc}
}

func (c *providerClient) CPE(ctx context.Context, in *ProvideCPERequest, opts ...grpc.CallOption) (*ProvideCPEResponse, error) {
	out := new(ProvideCPEResponse)
	err := c.cc.Invoke(ctx, "/core.Provider/CPE", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerClient) Access(ctx context.Context, in *ProvideAccessRequest, opts ...grpc.CallOption) (*ProvideAccessResponse, error) {
	out := new(ProvideAccessResponse)
	err := c.cc.Invoke(ctx, "/core.Provider/Access", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProviderServer is the server API for Provider service.
// All implementations must embed UnimplementedProviderServer
// for forward compatibility
type ProviderServer interface {
	// Ask provider to return a valid CPE for a access
	CPE(context.Context, *ProvideCPERequest) (*ProvideCPEResponse, error)
	// Ask a provider to return information about a selected access
	Access(context.Context, *ProvideAccessRequest) (*ProvideAccessResponse, error)
	mustEmbedUnimplementedProviderServer()
}

// UnimplementedProviderServer must be embedded to have forward compatible implementations.
type UnimplementedProviderServer struct {
}

func (UnimplementedProviderServer) CPE(context.Context, *ProvideCPERequest) (*ProvideCPEResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CPE not implemented")
}
func (UnimplementedProviderServer) Access(context.Context, *ProvideAccessRequest) (*ProvideAccessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Access not implemented")
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

func _Provider_CPE_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProvideCPERequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderServer).CPE(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Provider/CPE",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderServer).CPE(ctx, req.(*ProvideCPERequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Provider_Access_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProvideAccessRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderServer).Access(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Provider/Access",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderServer).Access(ctx, req.(*ProvideAccessRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Provider_ServiceDesc is the grpc.ServiceDesc for Provider service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Provider_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "core.Provider",
	HandlerType: (*ProviderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CPE",
			Handler:    _Provider_CPE_Handler,
		},
		{
			MethodName: "Access",
			Handler:    _Provider_Access_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "core.proto",
}
