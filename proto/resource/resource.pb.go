// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.3
// source: git.liero.se/opentelco/go-swpx/proto/resource/resource.proto

package resource

import (
	context "context"
	networkelement "git.liero.se/opentelco/go-swpx/proto/networkelement"
	provider "git.liero.se/opentelco/go-swpx/proto/provider"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDescGZIP(), []int{0}
}

type VersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *VersionResponse) Reset() {
	*x = VersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionResponse) ProtoMessage() {}

func (x *VersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[1]
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
	return file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDescGZIP(), []int{1}
}

func (x *VersionResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type TechnicalPortInformationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response string `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *TechnicalPortInformationResponse) Reset() {
	*x = TechnicalPortInformationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TechnicalPortInformationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TechnicalPortInformationResponse) ProtoMessage() {}

func (x *TechnicalPortInformationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TechnicalPortInformationResponse.ProtoReflect.Descriptor instead.
func (*TechnicalPortInformationResponse) Descriptor() ([]byte, []int) {
	return file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDescGZIP(), []int{2}
}

func (x *TechnicalPortInformationResponse) GetResponse() string {
	if x != nil {
		return x.Response
	}
	return ""
}

// used in request
type NetworkElement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hostname       string                  `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Ip             string                  `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	Interface      string                  `protobuf:"bytes,3,opt,name=interface,proto3" json:"interface,omitempty"`
	InterfaceIndex int64                   `protobuf:"varint,4,opt,name=interface_index,json=interfaceIndex,proto3" json:"interface_index,omitempty"`
	PhysicalIndex  int64                   `protobuf:"varint,5,opt,name=physical_index,json=physicalIndex,proto3" json:"physical_index,omitempty"`
	Conf           *provider.Configuration `protobuf:"bytes,6,opt,name=conf,proto3" json:"conf,omitempty"`
}

func (x *NetworkElement) Reset() {
	*x = NetworkElement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetworkElement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkElement) ProtoMessage() {}

func (x *NetworkElement) ProtoReflect() protoreflect.Message {
	mi := &file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkElement.ProtoReflect.Descriptor instead.
func (*NetworkElement) Descriptor() ([]byte, []int) {
	return file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDescGZIP(), []int{3}
}

func (x *NetworkElement) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *NetworkElement) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *NetworkElement) GetInterface() string {
	if x != nil {
		return x.Interface
	}
	return ""
}

func (x *NetworkElement) GetInterfaceIndex() int64 {
	if x != nil {
		return x.InterfaceIndex
	}
	return 0
}

func (x *NetworkElement) GetPhysicalIndex() int64 {
	if x != nil {
		return x.PhysicalIndex
	}
	return 0
}

func (x *NetworkElement) GetConf() *provider.Configuration {
	if x != nil {
		return x.Conf
	}
	return nil
}

type NetworkElementInterface struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index       int64  `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Alias       string `protobuf:"bytes,3,opt,name=alias,proto3" json:"alias,omitempty"`
}

func (x *NetworkElementInterface) Reset() {
	*x = NetworkElementInterface{}
	if protoimpl.UnsafeEnabled {
		mi := &file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetworkElementInterface) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkElementInterface) ProtoMessage() {}

func (x *NetworkElementInterface) ProtoReflect() protoreflect.Message {
	mi := &file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkElementInterface.ProtoReflect.Descriptor instead.
func (*NetworkElementInterface) Descriptor() ([]byte, []int) {
	return file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDescGZIP(), []int{4}
}

func (x *NetworkElementInterface) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *NetworkElementInterface) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *NetworkElementInterface) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

// map port to entire interface for faster retrieval
type NetworkElementInterfaces struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Interfaces map[string]*NetworkElementInterface `protobuf:"bytes,1,rep,name=interfaces,proto3" json:"interfaces,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *NetworkElementInterfaces) Reset() {
	*x = NetworkElementInterfaces{}
	if protoimpl.UnsafeEnabled {
		mi := &file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetworkElementInterfaces) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkElementInterfaces) ProtoMessage() {}

func (x *NetworkElementInterfaces) ProtoReflect() protoreflect.Message {
	mi := &file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkElementInterfaces.ProtoReflect.Descriptor instead.
func (*NetworkElementInterfaces) Descriptor() ([]byte, []int) {
	return file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDescGZIP(), []int{5}
}

func (x *NetworkElementInterfaces) GetInterfaces() map[string]*NetworkElementInterface {
	if x != nil {
		return x.Interfaces
	}
	return nil
}

var File_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto protoreflect.FileDescriptor

var file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDesc = []byte{
	0x0a, 0x3c, 0x67, 0x69, 0x74, 0x2e, 0x6c, 0x69, 0x65, 0x72, 0x6f, 0x2e, 0x73, 0x65, 0x2f, 0x6f,
	0x70, 0x65, 0x6e, 0x74, 0x65, 0x6c, 0x63, 0x6f, 0x2f, 0x67, 0x6f, 0x2d, 0x73, 0x77, 0x70, 0x78,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x1a, 0x1c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x2f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x2b, 0x0a, 0x0f, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x3e, 0x0a, 0x20, 0x54, 0x65, 0x63, 0x68, 0x6e, 0x69, 0x63,
	0x61, 0x6c, 0x50, 0x6f, 0x72, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xd7, 0x01, 0x0a, 0x0e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x70, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61,
	0x63, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x5f,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x25, 0x0a, 0x0e, 0x70,
	0x68, 0x79, 0x73, 0x69, 0x63, 0x61, 0x6c, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0d, 0x70, 0x68, 0x79, 0x73, 0x69, 0x63, 0x61, 0x6c, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x12, 0x2b, 0x0a, 0x04, 0x63, 0x6f, 0x6e, 0x66, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x63, 0x6f, 0x6e, 0x66, 0x22,
	0x67, 0x0a, 0x17, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x22, 0xd0, 0x01, 0x0a, 0x18, 0x4e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x66, 0x61, 0x63, 0x65, 0x73, 0x12, 0x52, 0x0a, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61,
	0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x45, 0x6c, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x2e, 0x49, 0x6e,
	0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x1a, 0x60, 0x0a, 0x0f, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x37,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x85, 0x03, 0x0a, 0x08,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x35, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x0f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x19, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x4d, 0x0a, 0x18, 0x54, 0x65, 0x63, 0x68, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x50, 0x6f, 0x72, 0x74,
	0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x45, 0x6c,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x1a, 0x17, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x65,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x4c,
	0x0a, 0x0c, 0x4d, 0x61, 0x70, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x12, 0x18,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x1a, 0x22, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x45, 0x6c, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x12, 0x51, 0x0a, 0x11,
	0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x50, 0x68, 0x79, 0x73, 0x69, 0x63, 0x61,
	0x6c, 0x12, 0x18, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x1a, 0x22, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x45, 0x6c,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x12,
	0x52, 0x0a, 0x19, 0x47, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x72, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x45,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x1a, 0x1b, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x72, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x2e, 0x6c, 0x69, 0x65, 0x72, 0x6f,
	0x2e, 0x73, 0x65, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x74, 0x65, 0x6c, 0x63, 0x6f, 0x2f, 0x67, 0x6f,
	0x2d, 0x73, 0x77, 0x70, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDescOnce sync.Once
	file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDescData = file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDesc
)

func file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDescGZIP() []byte {
	file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDescOnce.Do(func() {
		file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDescData = protoimpl.X.CompressGZIP(file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDescData)
	})
	return file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDescData
}

var file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_goTypes = []interface{}{
	(*Empty)(nil),                            // 0: resource.Empty
	(*VersionResponse)(nil),                  // 1: resource.VersionResponse
	(*TechnicalPortInformationResponse)(nil), // 2: resource.TechnicalPortInformationResponse
	(*NetworkElement)(nil),                   // 3: resource.NetworkElement
	(*NetworkElementInterface)(nil),          // 4: resource.NetworkElementInterface
	(*NetworkElementInterfaces)(nil),         // 5: resource.NetworkElementInterfaces
	nil,                                      // 6: resource.NetworkElementInterfaces.InterfacesEntry
	(*provider.Configuration)(nil),           // 7: provider.Configuration
	(*networkelement.Element)(nil),           // 8: networkelement.Element
	(*networkelement.Transceiver)(nil),       // 9: networkelement.Transceiver
}
var file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_depIdxs = []int32{
	7, // 0: resource.NetworkElement.conf:type_name -> provider.Configuration
	6, // 1: resource.NetworkElementInterfaces.interfaces:type_name -> resource.NetworkElementInterfaces.InterfacesEntry
	4, // 2: resource.NetworkElementInterfaces.InterfacesEntry.value:type_name -> resource.NetworkElementInterface
	0, // 3: resource.Resource.Version:input_type -> resource.Empty
	3, // 4: resource.Resource.TechnicalPortInformation:input_type -> resource.NetworkElement
	3, // 5: resource.Resource.MapInterface:input_type -> resource.NetworkElement
	3, // 6: resource.Resource.MapEntityPhysical:input_type -> resource.NetworkElement
	3, // 7: resource.Resource.GetTransceiverInformation:input_type -> resource.NetworkElement
	1, // 8: resource.Resource.Version:output_type -> resource.VersionResponse
	8, // 9: resource.Resource.TechnicalPortInformation:output_type -> networkelement.Element
	5, // 10: resource.Resource.MapInterface:output_type -> resource.NetworkElementInterfaces
	5, // 11: resource.Resource.MapEntityPhysical:output_type -> resource.NetworkElementInterfaces
	9, // 12: resource.Resource.GetTransceiverInformation:output_type -> networkelement.Transceiver
	8, // [8:13] is the sub-list for method output_type
	3, // [3:8] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_init() }
func file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_init() {
	if File_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TechnicalPortInformationResponse); i {
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
		file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetworkElement); i {
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
		file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetworkElementInterface); i {
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
		file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetworkElementInterfaces); i {
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
			RawDescriptor: file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_goTypes,
		DependencyIndexes: file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_depIdxs,
		MessageInfos:      file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_msgTypes,
	}.Build()
	File_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto = out.File
	file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_rawDesc = nil
	file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_goTypes = nil
	file_git_liero_se_opentelco_go_swpx_proto_resource_resource_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ResourceClient is the client API for Resource service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ResourceClient interface {
	Version(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*VersionResponse, error)
	TechnicalPortInformation(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*networkelement.Element, error)
	MapInterface(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*NetworkElementInterfaces, error)
	MapEntityPhysical(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*NetworkElementInterfaces, error)
	GetTransceiverInformation(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*networkelement.Transceiver, error)
}

type resourceClient struct {
	cc grpc.ClientConnInterface
}

func NewResourceClient(cc grpc.ClientConnInterface) ResourceClient {
	return &resourceClient{cc}
}

func (c *resourceClient) Version(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*VersionResponse, error) {
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

// ResourceServer is the server API for Resource service.
type ResourceServer interface {
	Version(context.Context, *Empty) (*VersionResponse, error)
	TechnicalPortInformation(context.Context, *NetworkElement) (*networkelement.Element, error)
	MapInterface(context.Context, *NetworkElement) (*NetworkElementInterfaces, error)
	MapEntityPhysical(context.Context, *NetworkElement) (*NetworkElementInterfaces, error)
	GetTransceiverInformation(context.Context, *NetworkElement) (*networkelement.Transceiver, error)
}

// UnimplementedResourceServer can be embedded to have forward compatible implementations.
type UnimplementedResourceServer struct {
}

func (*UnimplementedResourceServer) Version(context.Context, *Empty) (*VersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Version not implemented")
}
func (*UnimplementedResourceServer) TechnicalPortInformation(context.Context, *NetworkElement) (*networkelement.Element, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TechnicalPortInformation not implemented")
}
func (*UnimplementedResourceServer) MapInterface(context.Context, *NetworkElement) (*NetworkElementInterfaces, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MapInterface not implemented")
}
func (*UnimplementedResourceServer) MapEntityPhysical(context.Context, *NetworkElement) (*NetworkElementInterfaces, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MapEntityPhysical not implemented")
}
func (*UnimplementedResourceServer) GetTransceiverInformation(context.Context, *NetworkElement) (*networkelement.Transceiver, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransceiverInformation not implemented")
}

func RegisterResourceServer(s *grpc.Server, srv ResourceServer) {
	s.RegisterService(&_Resource_serviceDesc, srv)
}

func _Resource_Version_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
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
		return srv.(ResourceServer).Version(ctx, req.(*Empty))
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

var _Resource_serviceDesc = grpc.ServiceDesc{
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
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "git.liero.se/opentelco/go-swpx/proto/resource/resource.proto",
}
