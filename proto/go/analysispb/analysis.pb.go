//
// File: analysis.proto
// Project: src
// File Created: Sunday, 14th February 2021 1:47:04 pm
// Author: Mathias Ehrlin (mathias.ehrlin@vx.se)
// -----
// Last Modified: Sunday, 14th February 2021 1:56:06 pm
// Modified By: Mathias Ehrlin (mathias.ehrlin@vx.se>)
// -----
// Copyright - 2021 VX Service Delivery AB
//
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// -----

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: analysis.proto

package analysispb

import (
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

type Analysis_Level int32

const (
	Analysis_NOT_SET Analysis_Level = 0
	Analysis_FAILURE Analysis_Level = 1
	Analysis_WARNING Analysis_Level = 2
	Analysis_OK      Analysis_Level = 3
)

// Enum value maps for Analysis_Level.
var (
	Analysis_Level_name = map[int32]string{
		0: "NOT_SET",
		1: "FAILURE",
		2: "WARNING",
		3: "OK",
	}
	Analysis_Level_value = map[string]int32{
		"NOT_SET": 0,
		"FAILURE": 1,
		"WARNING": 2,
		"OK":      3,
	}
)

func (x Analysis_Level) Enum() *Analysis_Level {
	p := new(Analysis_Level)
	*p = x
	return p
}

func (x Analysis_Level) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Analysis_Level) Descriptor() protoreflect.EnumDescriptor {
	return file_analysis_proto_enumTypes[0].Descriptor()
}

func (Analysis_Level) Type() protoreflect.EnumType {
	return &file_analysis_proto_enumTypes[0]
}

func (x Analysis_Level) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Analysis_Level.Descriptor instead.
func (Analysis_Level) EnumDescriptor() ([]byte, []int) {
	return file_analysis_proto_rawDescGZIP(), []int{0, 0}
}

type Analysis struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Level     Analysis_Level `protobuf:"varint,1,opt,name=level,proto3,enum=analysis.Analysis_Level" json:"level,omitempty"`
	Message   string         `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Value     string         `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	Threshold string         `protobuf:"bytes,4,opt,name=threshold,proto3" json:"threshold,omitempty"`
	Type      string         `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *Analysis) Reset() {
	*x = Analysis{}
	if protoimpl.UnsafeEnabled {
		mi := &file_analysis_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Analysis) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Analysis) ProtoMessage() {}

func (x *Analysis) ProtoReflect() protoreflect.Message {
	mi := &file_analysis_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Analysis.ProtoReflect.Descriptor instead.
func (*Analysis) Descriptor() ([]byte, []int) {
	return file_analysis_proto_rawDescGZIP(), []int{0}
}

func (x *Analysis) GetLevel() Analysis_Level {
	if x != nil {
		return x.Level
	}
	return Analysis_NOT_SET
}

func (x *Analysis) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Analysis) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Analysis) GetThreshold() string {
	if x != nil {
		return x.Threshold
	}
	return ""
}

func (x *Analysis) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

var File_analysis_proto protoreflect.FileDescriptor

var file_analysis_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x22, 0xd4, 0x01, 0x0a, 0x08, 0x41,
	0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x12, 0x2e, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69,
	0x73, 0x2e, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x2e, 0x4c, 0x65, 0x76, 0x65, 0x6c,
	0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x68, 0x72, 0x65, 0x73,
	0x68, 0x6f, 0x6c, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x68, 0x72, 0x65,
	0x73, 0x68, 0x6f, 0x6c, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x36, 0x0a, 0x05, 0x4c, 0x65, 0x76,
	0x65, 0x6c, 0x12, 0x0b, 0x0a, 0x07, 0x4e, 0x4f, 0x54, 0x5f, 0x53, 0x45, 0x54, 0x10, 0x00, 0x12,
	0x0b, 0x0a, 0x07, 0x46, 0x41, 0x49, 0x4c, 0x55, 0x52, 0x45, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07,
	0x57, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10,
	0x03, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x2e, 0x6c, 0x69, 0x65, 0x72, 0x6f, 0x2e, 0x73,
	0x65, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x74, 0x65, 0x6c, 0x63, 0x6f, 0x2f, 0x67, 0x6f, 0x2d, 0x73,
	0x77, 0x70, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x6e, 0x61,
	0x6c, 0x79, 0x73, 0x69, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_analysis_proto_rawDescOnce sync.Once
	file_analysis_proto_rawDescData = file_analysis_proto_rawDesc
)

func file_analysis_proto_rawDescGZIP() []byte {
	file_analysis_proto_rawDescOnce.Do(func() {
		file_analysis_proto_rawDescData = protoimpl.X.CompressGZIP(file_analysis_proto_rawDescData)
	})
	return file_analysis_proto_rawDescData
}

var file_analysis_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_analysis_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_analysis_proto_goTypes = []interface{}{
	(Analysis_Level)(0), // 0: analysis.Analysis.Level
	(*Analysis)(nil),    // 1: analysis.Analysis
}
var file_analysis_proto_depIdxs = []int32{
	0, // 0: analysis.Analysis.level:type_name -> analysis.Analysis.Level
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_analysis_proto_init() }
func file_analysis_proto_init() {
	if File_analysis_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_analysis_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Analysis); i {
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
			RawDescriptor: file_analysis_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_analysis_proto_goTypes,
		DependencyIndexes: file_analysis_proto_depIdxs,
		EnumInfos:         file_analysis_proto_enumTypes,
		MessageInfos:      file_analysis_proto_msgTypes,
	}.Build()
	File_analysis_proto = out.File
	file_analysis_proto_rawDesc = nil
	file_analysis_proto_goTypes = nil
	file_analysis_proto_depIdxs = nil
}
