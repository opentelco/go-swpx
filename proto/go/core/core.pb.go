//
// Copyright (c) 2020. Liero AB
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
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: core.proto

package core

import (
	networkelement "git.liero.se/opentelco/go-swpx/proto/go/networkelement"
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

type Request_Type int32

const (
	Request_NOT_SET                 Request_Type = 0
	Request_GET_TECHNICAL_INFO      Request_Type = 1
	Request_GET_TECHNICAL_INFO_PORT Request_Type = 2
)

// Enum value maps for Request_Type.
var (
	Request_Type_name = map[int32]string{
		0: "NOT_SET",
		1: "GET_TECHNICAL_INFO",
		2: "GET_TECHNICAL_INFO_PORT",
	}
	Request_Type_value = map[string]int32{
		"NOT_SET":                 0,
		"GET_TECHNICAL_INFO":      1,
		"GET_TECHNICAL_INFO_PORT": 2,
	}
)

func (x Request_Type) Enum() *Request_Type {
	p := new(Request_Type)
	*p = x
	return p
}

func (x Request_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Request_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_core_proto_enumTypes[0].Descriptor()
}

func (Request_Type) Type() protoreflect.EnumType {
	return &file_core_proto_enumTypes[0]
}

func (x Request_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Request_Type.Descriptor instead.
func (Request_Type) EnumDescriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{1, 0}
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Code    int32  `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{0}
}

func (x *Error) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Error) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// settings for the request
	ProviderPlugin         string       `protobuf:"bytes,1,opt,name=provider_plugin,json=providerPlugin,proto3" json:"provider_plugin,omitempty"`
	ResourcePlugin         string       `protobuf:"bytes,2,opt,name=resource_plugin,json=resourcePlugin,proto3" json:"resource_plugin,omitempty"`
	RecreateIndex          bool         `protobuf:"varint,3,opt,name=recreate_index,json=recreateIndex,proto3" json:"recreate_index,omitempty"`
	DisableDistributedLock bool         `protobuf:"varint,4,opt,name=disable_distributed_lock,json=disableDistributedLock,proto3" json:"disable_distributed_lock,omitempty"`
	Timeout                string       `protobuf:"bytes,5,opt,name=timeout,proto3" json:"timeout,omitempty"`
	CacheTtl               string       `protobuf:"bytes,6,opt,name=cache_ttl,json=cacheTtl,proto3" json:"cache_ttl,omitempty"`
	Type                   Request_Type `protobuf:"varint,7,opt,name=type,proto3,enum=core.Request_Type" json:"type,omitempty"`
	// used to locate the resource
	AccessId string `protobuf:"bytes,8,opt,name=access_id,json=accessId,proto3" json:"access_id,omitempty"`
	Hostname string `protobuf:"bytes,9,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Port     string `protobuf:"bytes,10,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[1]
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
	return file_core_proto_rawDescGZIP(), []int{1}
}

func (x *Request) GetProviderPlugin() string {
	if x != nil {
		return x.ProviderPlugin
	}
	return ""
}

func (x *Request) GetResourcePlugin() string {
	if x != nil {
		return x.ResourcePlugin
	}
	return ""
}

func (x *Request) GetRecreateIndex() bool {
	if x != nil {
		return x.RecreateIndex
	}
	return false
}

func (x *Request) GetDisableDistributedLock() bool {
	if x != nil {
		return x.DisableDistributedLock
	}
	return false
}

func (x *Request) GetTimeout() string {
	if x != nil {
		return x.Timeout
	}
	return ""
}

func (x *Request) GetCacheTtl() string {
	if x != nil {
		return x.CacheTtl
	}
	return ""
}

func (x *Request) GetType() Request_Type {
	if x != nil {
		return x.Type
	}
	return Request_NOT_SET
}

func (x *Request) GetAccessId() string {
	if x != nil {
		return x.AccessId
	}
	return ""
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

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Request         *Request                `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
	NetworkElement  *networkelement.Element `protobuf:"bytes,2,opt,name=network_element,json=networkElement,proto3" json:"network_element,omitempty"`
	PhysicalPort    string                  `protobuf:"bytes,3,opt,name=physical_port,json=physicalPort,proto3" json:"physical_port,omitempty"`
	Error           *Error                  `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
	RequestAccessId string                  `protobuf:"bytes,5,opt,name=request_access_id,json=requestAccessId,proto3" json:"request_access_id,omitempty"`
	ExecutionTime   string                  `protobuf:"bytes,6,opt,name=execution_time,json=executionTime,proto3" json:"execution_time,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{2}
}

func (x *Response) GetRequest() *Request {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *Response) GetNetworkElement() *networkelement.Element {
	if x != nil {
		return x.NetworkElement
	}
	return nil
}

func (x *Response) GetPhysicalPort() string {
	if x != nil {
		return x.PhysicalPort
	}
	return ""
}

func (x *Response) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *Response) GetRequestAccessId() string {
	if x != nil {
		return x.RequestAccessId
	}
	return ""
}

func (x *Response) GetExecutionTime() string {
	if x != nil {
		return x.ExecutionTime
	}
	return ""
}

type CommandRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CommandRequest) Reset() {
	*x = CommandRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommandRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommandRequest) ProtoMessage() {}

func (x *CommandRequest) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommandRequest.ProtoReflect.Descriptor instead.
func (*CommandRequest) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{3}
}

type CommandResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CommandResponse) Reset() {
	*x = CommandResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommandResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommandResponse) ProtoMessage() {}

func (x *CommandResponse) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommandResponse.ProtoReflect.Descriptor instead.
func (*CommandResponse) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{4}
}

type InformationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InformationResponse) Reset() {
	*x = InformationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InformationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InformationResponse) ProtoMessage() {}

func (x *InformationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InformationResponse.ProtoReflect.Descriptor instead.
func (*InformationResponse) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{5}
}

type ProvideCPERequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ProvideCPERequest) Reset() {
	*x = ProvideCPERequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProvideCPERequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProvideCPERequest) ProtoMessage() {}

func (x *ProvideCPERequest) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProvideCPERequest.ProtoReflect.Descriptor instead.
func (*ProvideCPERequest) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{6}
}

type ProvideCPEResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ProvideCPEResponse) Reset() {
	*x = ProvideCPEResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProvideCPEResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProvideCPEResponse) ProtoMessage() {}

func (x *ProvideCPEResponse) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProvideCPEResponse.ProtoReflect.Descriptor instead.
func (*ProvideCPEResponse) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{7}
}

type ProvideAccessRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ProvideAccessRequest) Reset() {
	*x = ProvideAccessRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProvideAccessRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProvideAccessRequest) ProtoMessage() {}

func (x *ProvideAccessRequest) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProvideAccessRequest.ProtoReflect.Descriptor instead.
func (*ProvideAccessRequest) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{8}
}

type ProvideAccessResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ProvideAccessResponse) Reset() {
	*x = ProvideAccessResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProvideAccessResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProvideAccessResponse) ProtoMessage() {}

func (x *ProvideAccessResponse) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProvideAccessResponse.ProtoReflect.Descriptor instead.
func (*ProvideAccessResponse) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{9}
}

var File_core_proto protoreflect.FileDescriptor

var file_core_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x6f,
	0x72, 0x65, 0x1a, 0x15, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x65, 0x6c, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x35, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0xb2, 0x03,
	0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x70, 0x72, 0x6f,
	0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x50, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x12, 0x25, 0x0a, 0x0e, 0x72,
	0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0d, 0x72, 0x65, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x12, 0x38, 0x0a, 0x18, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x64, 0x69,
	0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x5f, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x16, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x44, 0x69, 0x73,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x4c, 0x6f, 0x63, 0x6b, 0x12, 0x18, 0x0a, 0x07,
	0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f,
	0x74, 0x74, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x63, 0x68, 0x65,
	0x54, 0x74, 0x6c, 0x12, 0x26, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x12, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x61,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x48, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x0b, 0x0a, 0x07, 0x4e, 0x4f, 0x54, 0x5f, 0x53, 0x45, 0x54, 0x10, 0x00, 0x12, 0x16, 0x0a,
	0x12, 0x47, 0x45, 0x54, 0x5f, 0x54, 0x45, 0x43, 0x48, 0x4e, 0x49, 0x43, 0x41, 0x4c, 0x5f, 0x49,
	0x4e, 0x46, 0x4f, 0x10, 0x01, 0x12, 0x1b, 0x0a, 0x17, 0x47, 0x45, 0x54, 0x5f, 0x54, 0x45, 0x43,
	0x48, 0x4e, 0x49, 0x43, 0x41, 0x4c, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x5f, 0x50, 0x4f, 0x52, 0x54,
	0x10, 0x02, 0x22, 0x90, 0x02, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x27, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0d, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52,
	0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x40, 0x0a, 0x0f, 0x6e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x5f, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x6c, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x2e, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0e, 0x6e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x68,
	0x79, 0x73, 0x69, 0x63, 0x61, 0x6c, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x70, 0x68, 0x79, 0x73, 0x69, 0x63, 0x61, 0x6c, 0x50, 0x6f, 0x72, 0x74, 0x12,
	0x21, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x12, 0x2a, 0x0a, 0x11, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x61, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x49, 0x64, 0x12, 0x25,
	0x0a, 0x0e, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f,
	0x6e, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x10, 0x0a, 0x0e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x11, 0x0a, 0x0f, 0x43, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x15, 0x0a, 0x13, 0x49, 0x6e,
	0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x13, 0x0a, 0x11, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x43, 0x50, 0x45, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x14, 0x0a, 0x12, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x43, 0x50, 0x45, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0x0a, 0x14,
	0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x17, 0x0a, 0x15, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x41,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xa7, 0x01,
	0x0a, 0x04, 0x43, 0x6f, 0x72, 0x65, 0x12, 0x25, 0x0a, 0x04, 0x50, 0x6f, 0x6c, 0x6c, 0x12, 0x0d,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a,
	0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x14, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x0b, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x19, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x86, 0x01, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x76,
	0x69, 0x64, 0x65, 0x12, 0x38, 0x0a, 0x03, 0x43, 0x50, 0x45, 0x12, 0x17, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x43, 0x50, 0x45, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x43, 0x50, 0x45, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a,
	0x06, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x1a, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x50,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x2e, 0x6c, 0x69, 0x65, 0x72, 0x6f, 0x2e, 0x73, 0x65,
	0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x74, 0x65, 0x6c, 0x63, 0x6f, 0x2f, 0x67, 0x6f, 0x2d, 0x73, 0x77,
	0x70, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_core_proto_rawDescOnce sync.Once
	file_core_proto_rawDescData = file_core_proto_rawDesc
)

func file_core_proto_rawDescGZIP() []byte {
	file_core_proto_rawDescOnce.Do(func() {
		file_core_proto_rawDescData = protoimpl.X.CompressGZIP(file_core_proto_rawDescData)
	})
	return file_core_proto_rawDescData
}

var file_core_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_core_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_core_proto_goTypes = []interface{}{
	(Request_Type)(0),              // 0: core.Request.Type
	(*Error)(nil),                  // 1: core.Error
	(*Request)(nil),                // 2: core.Request
	(*Response)(nil),               // 3: core.Response
	(*CommandRequest)(nil),         // 4: core.CommandRequest
	(*CommandResponse)(nil),        // 5: core.CommandResponse
	(*InformationResponse)(nil),    // 6: core.InformationResponse
	(*ProvideCPERequest)(nil),      // 7: core.ProvideCPERequest
	(*ProvideCPEResponse)(nil),     // 8: core.ProvideCPEResponse
	(*ProvideAccessRequest)(nil),   // 9: core.ProvideAccessRequest
	(*ProvideAccessResponse)(nil),  // 10: core.ProvideAccessResponse
	(*networkelement.Element)(nil), // 11: networkelement.Element
	(*emptypb.Empty)(nil),          // 12: google.protobuf.Empty
}
var file_core_proto_depIdxs = []int32{
	0,  // 0: core.Request.type:type_name -> core.Request.Type
	2,  // 1: core.Response.request:type_name -> core.Request
	11, // 2: core.Response.network_element:type_name -> networkelement.Element
	1,  // 3: core.Response.error:type_name -> core.Error
	2,  // 4: core.Core.Poll:input_type -> core.Request
	4,  // 5: core.Core.Command:input_type -> core.CommandRequest
	12, // 6: core.Core.Information:input_type -> google.protobuf.Empty
	7,  // 7: core.Provide.CPE:input_type -> core.ProvideCPERequest
	9,  // 8: core.Provide.Access:input_type -> core.ProvideAccessRequest
	3,  // 9: core.Core.Poll:output_type -> core.Response
	5,  // 10: core.Core.Command:output_type -> core.CommandResponse
	6,  // 11: core.Core.Information:output_type -> core.InformationResponse
	8,  // 12: core.Provide.CPE:output_type -> core.ProvideCPEResponse
	10, // 13: core.Provide.Access:output_type -> core.ProvideAccessResponse
	9,  // [9:14] is the sub-list for method output_type
	4,  // [4:9] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_core_proto_init() }
func file_core_proto_init() {
	if File_core_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_core_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
		file_core_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_core_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_core_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommandRequest); i {
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
		file_core_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommandResponse); i {
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
		file_core_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InformationResponse); i {
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
		file_core_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProvideCPERequest); i {
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
		file_core_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProvideCPEResponse); i {
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
		file_core_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProvideAccessRequest); i {
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
		file_core_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProvideAccessResponse); i {
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
			RawDescriptor: file_core_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_core_proto_goTypes,
		DependencyIndexes: file_core_proto_depIdxs,
		EnumInfos:         file_core_proto_enumTypes,
		MessageInfos:      file_core_proto_msgTypes,
	}.Build()
	File_core_proto = out.File
	file_core_proto_rawDesc = nil
	file_core_proto_goTypes = nil
	file_core_proto_depIdxs = nil
}
