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
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: fleet.proto

package fleetpb

import (
	configurationpb "git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	devicepb "git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
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

type CollectDeviceParameters struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// the id of the device to collect information for
	DeviceId string `protobuf:"bytes,1,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty" bson:"device_id"`
	// blocking is used to indicate if the call should block until the device is collected
	Blocking bool `protobuf:"varint,2,opt,name=blocking,proto3" json:"blocking,omitempty" bson:"blocking"`
}

func (x *CollectDeviceParameters) Reset() {
	*x = CollectDeviceParameters{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fleet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CollectDeviceParameters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CollectDeviceParameters) ProtoMessage() {}

func (x *CollectDeviceParameters) ProtoReflect() protoreflect.Message {
	mi := &file_fleet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CollectDeviceParameters.ProtoReflect.Descriptor instead.
func (*CollectDeviceParameters) Descriptor() ([]byte, []int) {
	return file_fleet_proto_rawDescGZIP(), []int{0}
}

func (x *CollectDeviceParameters) GetDeviceId() string {
	if x != nil {
		return x.DeviceId
	}
	return ""
}

func (x *CollectDeviceParameters) GetBlocking() bool {
	if x != nil {
		return x.Blocking
	}
	return false
}

type CollectConfigParameters struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// the id of the device to collect the configuration for
	DeviceId string `protobuf:"bytes,1,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty" bson:"device_id"`
	// blocking is used to indicate if the call should block until the configuration is collected
	Blocking bool `protobuf:"varint,2,opt,name=blocking,proto3" json:"blocking,omitempty" bson:"blocking"`
}

func (x *CollectConfigParameters) Reset() {
	*x = CollectConfigParameters{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fleet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CollectConfigParameters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CollectConfigParameters) ProtoMessage() {}

func (x *CollectConfigParameters) ProtoReflect() protoreflect.Message {
	mi := &file_fleet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CollectConfigParameters.ProtoReflect.Descriptor instead.
func (*CollectConfigParameters) Descriptor() ([]byte, []int) {
	return file_fleet_proto_rawDescGZIP(), []int{1}
}

func (x *CollectConfigParameters) GetDeviceId() string {
	if x != nil {
		return x.DeviceId
	}
	return ""
}

func (x *CollectConfigParameters) GetBlocking() bool {
	if x != nil {
		return x.Blocking
	}
	return false
}

type DiscoverDeviceParameters struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CreateDeviceParams *devicepb.CreateParameters `protobuf:"bytes,1,opt,name=create_device_params,json=createDeviceParams,proto3" json:"create_device_params,omitempty" bson:"create_device_params"`
	// blocking is used to indicate if the call should block until the device is discovered
	Blocking bool `protobuf:"varint,2,opt,name=blocking,proto3" json:"blocking,omitempty" bson:"blocking"`
}

func (x *DiscoverDeviceParameters) Reset() {
	*x = DiscoverDeviceParameters{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fleet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscoverDeviceParameters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoverDeviceParameters) ProtoMessage() {}

func (x *DiscoverDeviceParameters) ProtoReflect() protoreflect.Message {
	mi := &file_fleet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoverDeviceParameters.ProtoReflect.Descriptor instead.
func (*DiscoverDeviceParameters) Descriptor() ([]byte, []int) {
	return file_fleet_proto_rawDescGZIP(), []int{2}
}

func (x *DiscoverDeviceParameters) GetCreateDeviceParams() *devicepb.CreateParameters {
	if x != nil {
		return x.CreateDeviceParams
	}
	return nil
}

func (x *DiscoverDeviceParameters) GetBlocking() bool {
	if x != nil {
		return x.Blocking
	}
	return false
}

var File_fleet_proto protoreflect.FileDescriptor

var file_fleet_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x66,
	0x6c, 0x65, 0x65, 0x74, 0x1a, 0x12, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x5f, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x5f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x52, 0x0a, 0x17, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x69, 0x6e, 0x67, 0x22, 0x52, 0x0a, 0x17, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12,
	0x1b, 0x0a, 0x09, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x22, 0x88, 0x01, 0x0a, 0x18, 0x44, 0x69, 0x73,
	0x63, 0x6f, 0x76, 0x65, 0x72, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x50, 0x0a, 0x14, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74,
	0x65, 0x72, 0x73, 0x52, 0x12, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x69, 0x6e, 0x67, 0x32, 0xc3, 0x02, 0x0a, 0x0c, 0x46, 0x6c, 0x65, 0x65, 0x74, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x49, 0x0a, 0x0e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1f, 0x2e, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x44,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x1a, 0x14, 0x2e, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x22, 0x00, 0x12,
	0x47, 0x0a, 0x0d, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x1e, 0x2e, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73,
	0x1a, 0x14, 0x2e, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x0d, 0x43, 0x6f, 0x6c, 0x6c,
	0x65, 0x63, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1e, 0x2e, 0x66, 0x6c, 0x65, 0x65,
	0x74, 0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x1a, 0x22, 0x2e, 0x66, 0x6c, 0x65, 0x65,
	0x74, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12,
	0x48, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x1e, 0x2e, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x37, 0x5a, 0x35, 0x67, 0x69, 0x74,
	0x2e, 0x6c, 0x69, 0x65, 0x72, 0x6f, 0x2e, 0x73, 0x65, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x74, 0x65,
	0x6c, 0x63, 0x6f, 0x2f, 0x67, 0x6f, 0x2d, 0x73, 0x77, 0x70, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2f, 0x66, 0x6c, 0x65, 0x65, 0x74,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fleet_proto_rawDescOnce sync.Once
	file_fleet_proto_rawDescData = file_fleet_proto_rawDesc
)

func file_fleet_proto_rawDescGZIP() []byte {
	file_fleet_proto_rawDescOnce.Do(func() {
		file_fleet_proto_rawDescData = protoimpl.X.CompressGZIP(file_fleet_proto_rawDescData)
	})
	return file_fleet_proto_rawDescData
}

var file_fleet_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_fleet_proto_goTypes = []interface{}{
	(*CollectDeviceParameters)(nil),       // 0: fleet.CollectDeviceParameters
	(*CollectConfigParameters)(nil),       // 1: fleet.CollectConfigParameters
	(*DiscoverDeviceParameters)(nil),      // 2: fleet.DiscoverDeviceParameters
	(*devicepb.CreateParameters)(nil),     // 3: fleet.device.CreateParameters
	(*devicepb.DeleteParameters)(nil),     // 4: fleet.device.DeleteParameters
	(*devicepb.Device)(nil),               // 5: fleet.device.Device
	(*configurationpb.Configuration)(nil), // 6: fleet.configuration.Configuration
	(*emptypb.Empty)(nil),                 // 7: google.protobuf.Empty
}
var file_fleet_proto_depIdxs = []int32{
	3, // 0: fleet.DiscoverDeviceParameters.create_device_params:type_name -> fleet.device.CreateParameters
	2, // 1: fleet.FleetService.DiscoverDevice:input_type -> fleet.DiscoverDeviceParameters
	0, // 2: fleet.FleetService.CollectDevice:input_type -> fleet.CollectDeviceParameters
	1, // 3: fleet.FleetService.CollectConfig:input_type -> fleet.CollectConfigParameters
	4, // 4: fleet.FleetService.DeleteDevice:input_type -> fleet.device.DeleteParameters
	5, // 5: fleet.FleetService.DiscoverDevice:output_type -> fleet.device.Device
	5, // 6: fleet.FleetService.CollectDevice:output_type -> fleet.device.Device
	6, // 7: fleet.FleetService.CollectConfig:output_type -> fleet.configuration.Configuration
	7, // 8: fleet.FleetService.DeleteDevice:output_type -> google.protobuf.Empty
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_fleet_proto_init() }
func file_fleet_proto_init() {
	if File_fleet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fleet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CollectDeviceParameters); i {
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
		file_fleet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CollectConfigParameters); i {
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
		file_fleet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscoverDeviceParameters); i {
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
			RawDescriptor: file_fleet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_fleet_proto_goTypes,
		DependencyIndexes: file_fleet_proto_depIdxs,
		MessageInfos:      file_fleet_proto_msgTypes,
	}.Build()
	File_fleet_proto = out.File
	file_fleet_proto_rawDesc = nil
	file_fleet_proto_goTypes = nil
	file_fleet_proto_depIdxs = nil
}
