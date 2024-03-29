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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: plugin_resource.proto

package resourcepb

import (
	devicepb "go.opentelco.io/go-swpx/proto/go/devicepb"
	stanzapb "go.opentelco.io/go-swpx/proto/go/stanzapb"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type VersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty" bson:"version"`
}

func (x *VersionResponse) Reset() {
	*x = VersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_resource_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionResponse) ProtoMessage() {}

func (x *VersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_resource_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionResponse.ProtoReflect.Descriptor instead.
func (*VersionResponse) Descriptor() ([]byte, []int) {
	return file_plugin_resource_proto_rawDescGZIP(), []int{0}
}

func (x *VersionResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error   bool   `protobuf:"varint,1,opt,name=error,proto3" json:"error,omitempty" bson:"error"`
	Code    int32  `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty" bson:"code"`
	Type    string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty" bson:"type"`
	Message string `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty" bson:"message"`
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_resource_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_resource_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_plugin_resource_proto_rawDescGZIP(), []int{1}
}

func (x *Status) GetError() bool {
	if x != nil {
		return x.Error
	}
	return false
}

func (x *Status) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Status) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Status) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// PortIndexEntity is a entity that is used to map the port index to the description
// this is used both for the physical port and the logical port
type PortIndexEntity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index       int64  `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty" bson:"index"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty" bson:"description"`
	Alias       string `protobuf:"bytes,3,opt,name=alias,proto3" json:"alias,omitempty" bson:"alias"`
}

func (x *PortIndexEntity) Reset() {
	*x = PortIndexEntity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_resource_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PortIndexEntity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortIndexEntity) ProtoMessage() {}

func (x *PortIndexEntity) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_resource_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PortIndexEntity.ProtoReflect.Descriptor instead.
func (*PortIndexEntity) Descriptor() ([]byte, []int) {
	return file_plugin_resource_proto_rawDescGZIP(), []int{2}
}

func (x *PortIndexEntity) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *PortIndexEntity) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PortIndexEntity) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

// PortIndex is a map of the port index to the description
type PortIndex struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ports map[string]*PortIndexEntity `protobuf:"bytes,1,rep,name=ports,proto3" json:"ports,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3" bson:"ports" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *PortIndex) Reset() {
	*x = PortIndex{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_resource_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PortIndex) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortIndex) ProtoMessage() {}

func (x *PortIndex) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_resource_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PortIndex.ProtoReflect.Descriptor instead.
func (*PortIndex) Descriptor() ([]byte, []int) {
	return file_plugin_resource_proto_rawDescGZIP(), []int{3}
}

func (x *PortIndex) GetPorts() map[string]*PortIndexEntity {
	if x != nil {
		return x.Ports
	}
	return nil
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// hostname or ip address
	Hostname string `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty" bson:"hostname"`
	// port name (e.g. GigabitEthernet0/0/1 )
	Port string `protobuf:"bytes,2,opt,name=port,proto3" json:"port,omitempty" bson:"port"`
	// the number of interfaces discovered on the port
	// this is used to bulk get with by setting the "repition" to the number of interfaces
	NumInterfaces     int32 `protobuf:"varint,3,opt,name=num_interfaces,json=numInterfaces,proto3" json:"num_interfaces,omitempty" bson:"num_interfaces"`
	PhysicalPortIndex int64 `protobuf:"varint,4,opt,name=physical_port_index,json=physicalPortIndex,proto3" json:"physical_port_index,omitempty" bson:"physical_port_index"`
	LogicalPortIndex  int64 `protobuf:"varint,5,opt,name=logical_port_index,json=logicalPortIndex,proto3" json:"logical_port_index,omitempty" bson:"logical_port_index"`
	// should be a string we can parse to a duration
	// used to set the EOL timeout for requests
	Timeout string `protobuf:"bytes,6,opt,name=timeout,proto3" json:"timeout,omitempty" bson:"timeout"`
	// network regions passed down in from the SessionReqest.NetworkRegion
	NetworkRegion string `protobuf:"bytes,7,opt,name=network_region,json=networkRegion,proto3" json:"network_region,omitempty" bson:"network_region"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_resource_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_resource_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_plugin_resource_proto_rawDescGZIP(), []int{4}
}

func (x *Request) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *Request) GetPort() string {
	if x != nil {
		return x.Port
	}
	return ""
}

func (x *Request) GetNumInterfaces() int32 {
	if x != nil {
		return x.NumInterfaces
	}
	return 0
}

func (x *Request) GetPhysicalPortIndex() int64 {
	if x != nil {
		return x.PhysicalPortIndex
	}
	return 0
}

func (x *Request) GetLogicalPortIndex() int64 {
	if x != nil {
		return x.LogicalPortIndex
	}
	return 0
}

func (x *Request) GetTimeout() string {
	if x != nil {
		return x.Timeout
	}
	return ""
}

func (x *Request) GetNetworkRegion() string {
	if x != nil {
		return x.NetworkRegion
	}
	return ""
}

type GetRunningConfigParameters struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// hostname or ip address
	Hostname string `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty" bson:"hostname"`
	// should be a string we can parse to a duration
	// used to set the EOL timeout for requests
	Timeout string `protobuf:"bytes,2,opt,name=timeout,proto3" json:"timeout,omitempty" bson:"timeout"`
	// network regions passed down in from the SessionReqest.NetworkRegion
	NetworkRegion string `protobuf:"bytes,3,opt,name=network_region,json=networkRegion,proto3" json:"network_region,omitempty" bson:"network_region"`
}

func (x *GetRunningConfigParameters) Reset() {
	*x = GetRunningConfigParameters{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_resource_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRunningConfigParameters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRunningConfigParameters) ProtoMessage() {}

func (x *GetRunningConfigParameters) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_resource_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRunningConfigParameters.ProtoReflect.Descriptor instead.
func (*GetRunningConfigParameters) Descriptor() ([]byte, []int) {
	return file_plugin_resource_proto_rawDescGZIP(), []int{5}
}

func (x *GetRunningConfigParameters) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *GetRunningConfigParameters) GetTimeout() string {
	if x != nil {
		return x.Timeout
	}
	return ""
}

func (x *GetRunningConfigParameters) GetNetworkRegion() string {
	if x != nil {
		return x.NetworkRegion
	}
	return ""
}

type GetRunningConfigResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Config string `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty" bson:"config"`
}

func (x *GetRunningConfigResponse) Reset() {
	*x = GetRunningConfigResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_resource_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRunningConfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRunningConfigResponse) ProtoMessage() {}

func (x *GetRunningConfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_resource_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRunningConfigResponse.ProtoReflect.Descriptor instead.
func (*GetRunningConfigResponse) Descriptor() ([]byte, []int) {
	return file_plugin_resource_proto_rawDescGZIP(), []int{6}
}

func (x *GetRunningConfigResponse) GetConfig() string {
	if x != nil {
		return x.Config
	}
	return ""
}

type ConfigureStanzaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// hostname or ip address
	Hostname string `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty" bson:"hostname"`
	// should be a string we can parse to a duration
	// used to set the EOL timeout for requests
	Timeout string `protobuf:"bytes,2,opt,name=timeout,proto3" json:"timeout,omitempty" bson:"timeout"`
	// network regions passed down in from the SessionReqest.NetworkRegion
	NetworkRegion string `protobuf:"bytes,3,opt,name=network_region,json=networkRegion,proto3" json:"network_region,omitempty" bson:"network_region"`
	// the configuration to send to the device, each line is a string in the array
	Stanza []*stanzapb.ConfigurationLine `protobuf:"bytes,4,rep,name=stanza,proto3" json:"stanza,omitempty" bson:"stanza"`
}

func (x *ConfigureStanzaRequest) Reset() {
	*x = ConfigureStanzaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_resource_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigureStanzaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigureStanzaRequest) ProtoMessage() {}

func (x *ConfigureStanzaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_resource_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigureStanzaRequest.ProtoReflect.Descriptor instead.
func (*ConfigureStanzaRequest) Descriptor() ([]byte, []int) {
	return file_plugin_resource_proto_rawDescGZIP(), []int{7}
}

func (x *ConfigureStanzaRequest) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *ConfigureStanzaRequest) GetTimeout() string {
	if x != nil {
		return x.Timeout
	}
	return ""
}

func (x *ConfigureStanzaRequest) GetNetworkRegion() string {
	if x != nil {
		return x.NetworkRegion
	}
	return ""
}

func (x *ConfigureStanzaRequest) GetStanza() []*stanzapb.ConfigurationLine {
	if x != nil {
		return x.Stanza
	}
	return nil
}

var File_plugin_resource_proto protoreflect.FileDescriptor

var file_plugin_resource_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x1a, 0x0c, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0c, 0x73, 0x74, 0x61, 0x6e, 0x7a, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2b, 0x0a, 0x0f, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x60, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x5f, 0x0a, 0x0f, 0x50, 0x6f, 0x72,
	0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x22, 0x96, 0x01, 0x0a, 0x09, 0x50,
	0x6f, 0x72, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x34, 0x0a, 0x05, 0x70, 0x6f, 0x72, 0x74,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x50, 0x6f, 0x72,
	0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x1a, 0x53,
	0x0a, 0x0a, 0x50, 0x6f, 0x72, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2f,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x22, 0xff, 0x01, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70,
	0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12,
	0x25, 0x0a, 0x0e, 0x6e, 0x75, 0x6d, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x6e, 0x75, 0x6d, 0x49, 0x6e, 0x74, 0x65,
	0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x12, 0x2e, 0x0a, 0x13, 0x70, 0x68, 0x79, 0x73, 0x69, 0x63,
	0x61, 0x6c, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x11, 0x70, 0x68, 0x79, 0x73, 0x69, 0x63, 0x61, 0x6c, 0x50, 0x6f, 0x72,
	0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x2c, 0x0a, 0x12, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x61,
	0x6c, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x10, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x61, 0x6c, 0x50, 0x6f, 0x72, 0x74, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x25,
	0x0a, 0x0e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x52,
	0x65, 0x67, 0x69, 0x6f, 0x6e, 0x22, 0x79, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x52, 0x75, 0x6e, 0x6e,
	0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74,
	0x65, 0x72, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x6e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x22, 0x32, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x22, 0xa8, 0x01, 0x0a, 0x16, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x65, 0x53, 0x74, 0x61, 0x6e, 0x7a, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x74,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x69,
	0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x5f, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x31, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x6e, 0x7a, 0x61, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x73,
	0x74, 0x61, 0x6e, 0x7a, 0x61, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x6e, 0x65, 0x52, 0x06, 0x73, 0x74, 0x61, 0x6e, 0x7a, 0x61, 0x32,
	0x96, 0x06, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x3c, 0x0a, 0x07,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x19, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x08, 0x44, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x36, 0x0a, 0x0c, 0x4d, 0x61, 0x70,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x12, 0x3b, 0x0a, 0x11, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x50, 0x68,
	0x79, 0x73, 0x69, 0x63, 0x61, 0x6c, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x39,
	0x0a, 0x14, 0x42, 0x61, 0x73, 0x69, 0x63, 0x50, 0x6f, 0x72, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x72,
	0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3d, 0x0a, 0x18, 0x54, 0x65, 0x63,
	0x68, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x50, 0x6f, 0x72, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x37, 0x0a, 0x12, 0x41, 0x6c, 0x6c, 0x50,
	0x6f, 0x72, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x11,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0e, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x39, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e,
	0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x19,
	0x47, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x49, 0x6e,
	0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x72, 0x12, 0x47, 0x0a, 0x1c, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x73, 0x12, 0x5c, 0x0a, 0x10, 0x47, 0x65,
	0x74, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x24,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x75, 0x6e,
	0x6e, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65,
	0x74, 0x65, 0x72, 0x73, 0x1a, 0x22, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e,
	0x47, 0x65, 0x74, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4e, 0x0a, 0x0f, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x75, 0x72, 0x65, 0x53, 0x74, 0x61, 0x6e, 0x7a, 0x61, 0x12, 0x20, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65,
	0x53, 0x74, 0x61, 0x6e, 0x7a, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x73, 0x74, 0x61, 0x6e, 0x7a, 0x61, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x6f, 0x2e, 0x6f,
	0x70, 0x65, 0x6e, 0x74, 0x65, 0x6c, 0x63, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x67, 0x6f, 0x2d, 0x73,
	0x77, 0x70, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_plugin_resource_proto_rawDescOnce sync.Once
	file_plugin_resource_proto_rawDescData = file_plugin_resource_proto_rawDesc
)

func file_plugin_resource_proto_rawDescGZIP() []byte {
	file_plugin_resource_proto_rawDescOnce.Do(func() {
		file_plugin_resource_proto_rawDescData = protoimpl.X.CompressGZIP(file_plugin_resource_proto_rawDescData)
	})
	return file_plugin_resource_proto_rawDescData
}

var file_plugin_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_plugin_resource_proto_goTypes = []interface{}{
	(*VersionResponse)(nil),            // 0: resource.VersionResponse
	(*Status)(nil),                     // 1: resource.Status
	(*PortIndexEntity)(nil),            // 2: resource.PortIndexEntity
	(*PortIndex)(nil),                  // 3: resource.PortIndex
	(*Request)(nil),                    // 4: resource.Request
	(*GetRunningConfigParameters)(nil), // 5: resource.GetRunningConfigParameters
	(*GetRunningConfigResponse)(nil),   // 6: resource.GetRunningConfigResponse
	(*ConfigureStanzaRequest)(nil),     // 7: resource.ConfigureStanzaRequest
	nil,                                // 8: resource.PortIndex.PortsEntry
	(*stanzapb.ConfigurationLine)(nil), // 9: stanza.ConfigurationLine
	(*emptypb.Empty)(nil),              // 10: google.protobuf.Empty
	(*devicepb.Device)(nil),            // 11: device.Device
	(*devicepb.Transceiver)(nil),       // 12: device.Transceiver
	(*devicepb.Transceivers)(nil),      // 13: device.Transceivers
	(*stanzapb.ConfigureResponse)(nil), // 14: stanza.ConfigureResponse
}
var file_plugin_resource_proto_depIdxs = []int32{
	8,  // 0: resource.PortIndex.ports:type_name -> resource.PortIndex.PortsEntry
	9,  // 1: resource.ConfigureStanzaRequest.stanza:type_name -> stanza.ConfigurationLine
	2,  // 2: resource.PortIndex.PortsEntry.value:type_name -> resource.PortIndexEntity
	10, // 3: resource.Resource.Version:input_type -> google.protobuf.Empty
	4,  // 4: resource.Resource.Discover:input_type -> resource.Request
	4,  // 5: resource.Resource.MapInterface:input_type -> resource.Request
	4,  // 6: resource.Resource.MapEntityPhysical:input_type -> resource.Request
	4,  // 7: resource.Resource.BasicPortInformation:input_type -> resource.Request
	4,  // 8: resource.Resource.TechnicalPortInformation:input_type -> resource.Request
	4,  // 9: resource.Resource.AllPortInformation:input_type -> resource.Request
	4,  // 10: resource.Resource.GetDeviceInformation:input_type -> resource.Request
	4,  // 11: resource.Resource.GetTransceiverInformation:input_type -> resource.Request
	4,  // 12: resource.Resource.GetAllTransceiverInformation:input_type -> resource.Request
	5,  // 13: resource.Resource.GetRunningConfig:input_type -> resource.GetRunningConfigParameters
	7,  // 14: resource.Resource.ConfigureStanza:input_type -> resource.ConfigureStanzaRequest
	0,  // 15: resource.Resource.Version:output_type -> resource.VersionResponse
	11, // 16: resource.Resource.Discover:output_type -> device.Device
	3,  // 17: resource.Resource.MapInterface:output_type -> resource.PortIndex
	3,  // 18: resource.Resource.MapEntityPhysical:output_type -> resource.PortIndex
	11, // 19: resource.Resource.BasicPortInformation:output_type -> device.Device
	11, // 20: resource.Resource.TechnicalPortInformation:output_type -> device.Device
	11, // 21: resource.Resource.AllPortInformation:output_type -> device.Device
	11, // 22: resource.Resource.GetDeviceInformation:output_type -> device.Device
	12, // 23: resource.Resource.GetTransceiverInformation:output_type -> device.Transceiver
	13, // 24: resource.Resource.GetAllTransceiverInformation:output_type -> device.Transceivers
	6,  // 25: resource.Resource.GetRunningConfig:output_type -> resource.GetRunningConfigResponse
	14, // 26: resource.Resource.ConfigureStanza:output_type -> stanza.ConfigureResponse
	15, // [15:27] is the sub-list for method output_type
	3,  // [3:15] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_plugin_resource_proto_init() }
func file_plugin_resource_proto_init() {
	if File_plugin_resource_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_plugin_resource_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VersionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plugin_resource_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Status); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plugin_resource_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PortIndexEntity); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plugin_resource_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PortIndex); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plugin_resource_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plugin_resource_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRunningConfigParameters); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plugin_resource_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRunningConfigResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_plugin_resource_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigureStanzaRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_plugin_resource_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_plugin_resource_proto_goTypes,
		DependencyIndexes: file_plugin_resource_proto_depIdxs,
		MessageInfos:      file_plugin_resource_proto_msgTypes,
	}.Build()
	File_plugin_resource_proto = out.File
	file_plugin_resource_proto_rawDesc = nil
	file_plugin_resource_proto_goTypes = nil
	file_plugin_resource_proto_depIdxs = nil
}
