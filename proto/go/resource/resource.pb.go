// Code generated by protoc-gen-go. DO NOT EDIT.
// source: git.liero.se/opentelco/go-swpx/proto/src/resource.proto

package resource

import (
	context "context"
	fmt "fmt"
	networkelement "git.liero.se/opentelco/go-swpx/proto/go/networkelement"
	provider "git.liero.se/opentelco/go-swpx/proto/go/provider"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ff568e5765ad370, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type VersionResponse struct {
	Version              string   `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VersionResponse) Reset()         { *m = VersionResponse{} }
func (m *VersionResponse) String() string { return proto.CompactTextString(m) }
func (*VersionResponse) ProtoMessage()    {}
func (*VersionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ff568e5765ad370, []int{1}
}

func (m *VersionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VersionResponse.Unmarshal(m, b)
}
func (m *VersionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VersionResponse.Marshal(b, m, deterministic)
}
func (m *VersionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VersionResponse.Merge(m, src)
}
func (m *VersionResponse) XXX_Size() int {
	return xxx_messageInfo_VersionResponse.Size(m)
}
func (m *VersionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VersionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VersionResponse proto.InternalMessageInfo

func (m *VersionResponse) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

type Status struct {
	Error                bool     `protobuf:"varint,1,opt,name=error,proto3" json:"error,omitempty"`
	Code                 int32    `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	Type                 string   `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Message              string   `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ff568e5765ad370, []int{2}
}

func (m *Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Status.Unmarshal(m, b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Status.Marshal(b, m, deterministic)
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return xxx_messageInfo_Status.Size(m)
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetError() bool {
	if m != nil {
		return m.Error
	}
	return false
}

func (m *Status) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Status) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Status) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

// used in resource request
type NetworkElement struct {
	Hostname             string                  `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Ip                   string                  `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	Interface            string                  `protobuf:"bytes,3,opt,name=interface,proto3" json:"interface,omitempty"`
	InterfaceIndex       int64                   `protobuf:"varint,4,opt,name=interface_index,json=interfaceIndex,proto3" json:"interface_index,omitempty"`
	PhysicalIndex        int64                   `protobuf:"varint,5,opt,name=physical_index,json=physicalIndex,proto3" json:"physical_index,omitempty"`
	Conf                 *provider.Configuration `protobuf:"bytes,6,opt,name=conf,proto3" json:"conf,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *NetworkElement) Reset()         { *m = NetworkElement{} }
func (m *NetworkElement) String() string { return proto.CompactTextString(m) }
func (*NetworkElement) ProtoMessage()    {}
func (*NetworkElement) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ff568e5765ad370, []int{3}
}

func (m *NetworkElement) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkElement.Unmarshal(m, b)
}
func (m *NetworkElement) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkElement.Marshal(b, m, deterministic)
}
func (m *NetworkElement) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkElement.Merge(m, src)
}
func (m *NetworkElement) XXX_Size() int {
	return xxx_messageInfo_NetworkElement.Size(m)
}
func (m *NetworkElement) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkElement.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkElement proto.InternalMessageInfo

func (m *NetworkElement) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *NetworkElement) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *NetworkElement) GetInterface() string {
	if m != nil {
		return m.Interface
	}
	return ""
}

func (m *NetworkElement) GetInterfaceIndex() int64 {
	if m != nil {
		return m.InterfaceIndex
	}
	return 0
}

func (m *NetworkElement) GetPhysicalIndex() int64 {
	if m != nil {
		return m.PhysicalIndex
	}
	return 0
}

func (m *NetworkElement) GetConf() *provider.Configuration {
	if m != nil {
		return m.Conf
	}
	return nil
}

type NetworkElementInterface struct {
	Index                int64    `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Description          string   `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Alias                string   `protobuf:"bytes,3,opt,name=alias,proto3" json:"alias,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetworkElementInterface) Reset()         { *m = NetworkElementInterface{} }
func (m *NetworkElementInterface) String() string { return proto.CompactTextString(m) }
func (*NetworkElementInterface) ProtoMessage()    {}
func (*NetworkElementInterface) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ff568e5765ad370, []int{4}
}

func (m *NetworkElementInterface) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkElementInterface.Unmarshal(m, b)
}
func (m *NetworkElementInterface) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkElementInterface.Marshal(b, m, deterministic)
}
func (m *NetworkElementInterface) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkElementInterface.Merge(m, src)
}
func (m *NetworkElementInterface) XXX_Size() int {
	return xxx_messageInfo_NetworkElementInterface.Size(m)
}
func (m *NetworkElementInterface) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkElementInterface.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkElementInterface proto.InternalMessageInfo

func (m *NetworkElementInterface) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *NetworkElementInterface) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *NetworkElementInterface) GetAlias() string {
	if m != nil {
		return m.Alias
	}
	return ""
}

// map port to entire interface for faster retrieval
type NetworkElementInterfaces struct {
	Interfaces           map[string]*NetworkElementInterface `protobuf:"bytes,1,rep,name=interfaces,proto3" json:"interfaces,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                            `json:"-"`
	XXX_unrecognized     []byte                              `json:"-"`
	XXX_sizecache        int32                               `json:"-"`
}

func (m *NetworkElementInterfaces) Reset()         { *m = NetworkElementInterfaces{} }
func (m *NetworkElementInterfaces) String() string { return proto.CompactTextString(m) }
func (*NetworkElementInterfaces) ProtoMessage()    {}
func (*NetworkElementInterfaces) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ff568e5765ad370, []int{5}
}

func (m *NetworkElementInterfaces) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkElementInterfaces.Unmarshal(m, b)
}
func (m *NetworkElementInterfaces) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkElementInterfaces.Marshal(b, m, deterministic)
}
func (m *NetworkElementInterfaces) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkElementInterfaces.Merge(m, src)
}
func (m *NetworkElementInterfaces) XXX_Size() int {
	return xxx_messageInfo_NetworkElementInterfaces.Size(m)
}
func (m *NetworkElementInterfaces) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkElementInterfaces.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkElementInterfaces proto.InternalMessageInfo

func (m *NetworkElementInterfaces) GetInterfaces() map[string]*NetworkElementInterface {
	if m != nil {
		return m.Interfaces
	}
	return nil
}

type Transceivers struct {
	Transceivers         map[int32]*networkelement.Transceiver `protobuf:"bytes,1,rep,name=transceivers,proto3" json:"transceivers,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                              `json:"-"`
	XXX_unrecognized     []byte                                `json:"-"`
	XXX_sizecache        int32                                 `json:"-"`
}

func (m *Transceivers) Reset()         { *m = Transceivers{} }
func (m *Transceivers) String() string { return proto.CompactTextString(m) }
func (*Transceivers) ProtoMessage()    {}
func (*Transceivers) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ff568e5765ad370, []int{6}
}

func (m *Transceivers) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transceivers.Unmarshal(m, b)
}
func (m *Transceivers) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transceivers.Marshal(b, m, deterministic)
}
func (m *Transceivers) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transceivers.Merge(m, src)
}
func (m *Transceivers) XXX_Size() int {
	return xxx_messageInfo_Transceivers.Size(m)
}
func (m *Transceivers) XXX_DiscardUnknown() {
	xxx_messageInfo_Transceivers.DiscardUnknown(m)
}

var xxx_messageInfo_Transceivers proto.InternalMessageInfo

func (m *Transceivers) GetTransceivers() map[int32]*networkelement.Transceiver {
	if m != nil {
		return m.Transceivers
	}
	return nil
}

type NetworkElementWrapper struct {
	Element              *NetworkElement           `protobuf:"bytes,1,opt,name=element,proto3" json:"element,omitempty"`
	NumInterfaces        int32                     `protobuf:"varint,2,opt,name=numInterfaces,proto3" json:"numInterfaces,omitempty"`
	FullElement          *networkelement.Element   `protobuf:"bytes,3,opt,name=fullElement,proto3" json:"fullElement,omitempty"`
	PhysInterfaces       *NetworkElementInterfaces `protobuf:"bytes,4,opt,name=physInterfaces,proto3" json:"physInterfaces,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *NetworkElementWrapper) Reset()         { *m = NetworkElementWrapper{} }
func (m *NetworkElementWrapper) String() string { return proto.CompactTextString(m) }
func (*NetworkElementWrapper) ProtoMessage()    {}
func (*NetworkElementWrapper) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ff568e5765ad370, []int{7}
}

func (m *NetworkElementWrapper) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkElementWrapper.Unmarshal(m, b)
}
func (m *NetworkElementWrapper) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkElementWrapper.Marshal(b, m, deterministic)
}
func (m *NetworkElementWrapper) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkElementWrapper.Merge(m, src)
}
func (m *NetworkElementWrapper) XXX_Size() int {
	return xxx_messageInfo_NetworkElementWrapper.Size(m)
}
func (m *NetworkElementWrapper) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkElementWrapper.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkElementWrapper proto.InternalMessageInfo

func (m *NetworkElementWrapper) GetElement() *NetworkElement {
	if m != nil {
		return m.Element
	}
	return nil
}

func (m *NetworkElementWrapper) GetNumInterfaces() int32 {
	if m != nil {
		return m.NumInterfaces
	}
	return 0
}

func (m *NetworkElementWrapper) GetFullElement() *networkelement.Element {
	if m != nil {
		return m.FullElement
	}
	return nil
}

func (m *NetworkElementWrapper) GetPhysInterfaces() *NetworkElementInterfaces {
	if m != nil {
		return m.PhysInterfaces
	}
	return nil
}

func init() {
	proto.RegisterType((*Empty)(nil), "resource.Empty")
	proto.RegisterType((*VersionResponse)(nil), "resource.VersionResponse")
	proto.RegisterType((*Status)(nil), "resource.Status")
	proto.RegisterType((*NetworkElement)(nil), "resource.NetworkElement")
	proto.RegisterType((*NetworkElementInterface)(nil), "resource.NetworkElementInterface")
	proto.RegisterType((*NetworkElementInterfaces)(nil), "resource.NetworkElementInterfaces")
	proto.RegisterMapType((map[string]*NetworkElementInterface)(nil), "resource.NetworkElementInterfaces.InterfacesEntry")
	proto.RegisterType((*Transceivers)(nil), "resource.Transceivers")
	proto.RegisterMapType((map[int32]*networkelement.Transceiver)(nil), "resource.Transceivers.TransceiversEntry")
	proto.RegisterType((*NetworkElementWrapper)(nil), "resource.NetworkElementWrapper")
}

func init() {
	proto.RegisterFile("git.liero.se/opentelco/go-swpx/proto/src/resource.proto", fileDescriptor_8ff568e5765ad370)
}

var fileDescriptor_8ff568e5765ad370 = []byte{
	// 710 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0xdd, 0x6e, 0xd3, 0x4c,
	0x10, 0x95, 0xf3, 0xd3, 0xa4, 0x93, 0x36, 0xf9, 0xba, 0x6a, 0x55, 0x37, 0x5f, 0x25, 0x82, 0x05,
	0x22, 0x52, 0x45, 0x02, 0x46, 0xa8, 0xc0, 0x5d, 0x41, 0x51, 0x55, 0xd4, 0xa2, 0x62, 0x2a, 0x40,
	0x08, 0xa9, 0x35, 0xce, 0x24, 0xb5, 0xea, 0xec, 0x5a, 0xbb, 0x9b, 0xb6, 0x79, 0x23, 0x9e, 0x81,
	0xa7, 0xe0, 0x8e, 0x57, 0xe1, 0x12, 0x79, 0x77, 0xe3, 0x38, 0x41, 0xa6, 0x95, 0xe0, 0x6e, 0x67,
	0x72, 0xe6, 0x9c, 0x3d, 0x93, 0x99, 0x35, 0xec, 0x0e, 0x43, 0xd9, 0x89, 0x42, 0xe4, 0xac, 0x23,
	0xb0, 0xcb, 0x62, 0xa4, 0x12, 0xa3, 0x80, 0x75, 0x87, 0xec, 0xa1, 0xb8, 0x8a, 0xaf, 0xbb, 0x31,
	0x67, 0x92, 0x75, 0x05, 0x0f, 0xba, 0x1c, 0x05, 0x1b, 0xf3, 0x00, 0x3b, 0x2a, 0x45, 0xaa, 0xd3,
	0xb8, 0xb9, 0x41, 0x51, 0x5e, 0x31, 0x7e, 0x71, 0x8a, 0x11, 0x8e, 0x90, 0x4a, 0x0d, 0x68, 0xd6,
	0x63, 0xce, 0x2e, 0xc3, 0x3e, 0x72, 0x1d, 0x3b, 0x15, 0x28, 0xf7, 0x46, 0xb1, 0x9c, 0x38, 0x3b,
	0xd0, 0x78, 0x8f, 0x5c, 0x84, 0x8c, 0x7a, 0x28, 0x62, 0x46, 0x05, 0x12, 0x1b, 0x2a, 0x97, 0x3a,
	0x65, 0x5b, 0x2d, 0xab, 0xbd, 0xec, 0x4d, 0x43, 0xe7, 0x0c, 0x96, 0xde, 0x49, 0x5f, 0x8e, 0x05,
	0x59, 0x87, 0x32, 0x72, 0xce, 0xb8, 0x42, 0x54, 0x3d, 0x1d, 0x10, 0x02, 0xa5, 0x80, 0xf5, 0xd1,
	0x2e, 0xb4, 0xac, 0x76, 0xd9, 0x53, 0xe7, 0x24, 0x27, 0x27, 0x31, 0xda, 0x45, 0x45, 0xa5, 0xce,
	0x89, 0xc2, 0x08, 0x85, 0xf0, 0x87, 0x68, 0x97, 0xb4, 0x82, 0x09, 0x9d, 0x1f, 0x16, 0xd4, 0xdf,
	0x68, 0x07, 0x3d, 0x6d, 0x80, 0x34, 0xa1, 0x7a, 0xce, 0x84, 0xa4, 0xfe, 0x08, 0xcd, 0x7d, 0xd2,
	0x98, 0xd4, 0xa1, 0x10, 0xc6, 0x4a, 0x6e, 0xd9, 0x2b, 0x84, 0x31, 0xd9, 0x86, 0xe5, 0x90, 0x4a,
	0xe4, 0x03, 0x3f, 0x98, 0x2a, 0xce, 0x12, 0xe4, 0x01, 0x34, 0xd2, 0xe0, 0x34, 0xa4, 0x7d, 0xbc,
	0x56, 0xf2, 0x45, 0xaf, 0x9e, 0xa6, 0x0f, 0x92, 0x2c, 0xb9, 0x0f, 0xf5, 0xf8, 0x7c, 0x22, 0xc2,
	0xc0, 0x8f, 0x0c, 0xae, 0xac, 0x70, 0xab, 0xd3, 0xac, 0x86, 0xed, 0x24, 0x76, 0xe9, 0xc0, 0x5e,
	0x6a, 0x59, 0xed, 0x9a, 0xbb, 0xd9, 0x49, 0x7b, 0xfc, 0x8a, 0xd1, 0x41, 0x38, 0x1c, 0x73, 0x5f,
	0x26, 0x7d, 0x55, 0x20, 0x67, 0x08, 0x9b, 0xf3, 0xc6, 0x0e, 0xd2, 0x7b, 0xad, 0x43, 0x59, 0xab,
	0x58, 0x4a, 0x45, 0x07, 0xa4, 0x05, 0xb5, 0x3e, 0x8a, 0x80, 0x87, 0x71, 0xc2, 0x62, 0x4c, 0x66,
	0x53, 0x49, 0x9d, 0x1f, 0x85, 0xbe, 0x30, 0x4e, 0x75, 0xe0, 0x7c, 0xb7, 0xc0, 0xce, 0x51, 0x12,
	0xc4, 0x03, 0x48, 0xbd, 0x0a, 0xdb, 0x6a, 0x15, 0xdb, 0x35, 0xd7, 0xed, 0xa4, 0xd3, 0x94, 0x57,
	0xd7, 0x99, 0x1d, 0x7b, 0x54, 0xf2, 0x89, 0x97, 0x61, 0x69, 0x9e, 0x41, 0x63, 0xe1, 0x67, 0xf2,
	0x1f, 0x14, 0x2f, 0x70, 0x62, 0xfe, 0xae, 0xe4, 0x48, 0x76, 0xa1, 0x7c, 0xe9, 0x47, 0x63, 0x3d,
	0x1b, 0x35, 0xf7, 0xee, 0x8d, 0x9a, 0x9e, 0xc6, 0xbf, 0x28, 0x3c, 0xb3, 0x9c, 0x6f, 0x16, 0xac,
	0x9c, 0x70, 0x9f, 0x8a, 0x00, 0xc3, 0x64, 0x16, 0xc9, 0x21, 0xac, 0xc8, 0x4c, 0x6c, 0x8c, 0xb4,
	0x67, 0xa4, 0x59, 0xf4, 0x5c, 0xa0, 0xaf, 0x3f, 0x57, 0xdd, 0xfc, 0x0c, 0x6b, 0xbf, 0x41, 0xb2,
	0x16, 0xca, 0xda, 0xc2, 0xe3, 0x79, 0x0b, 0xff, 0x77, 0xcc, 0xaa, 0x4d, 0x37, 0x2d, 0xc3, 0x91,
	0xbd, 0xfc, 0x4f, 0x0b, 0x36, 0xe6, 0x3d, 0x7e, 0xe0, 0x7e, 0x1c, 0x23, 0x27, 0x2e, 0x54, 0x4c,
	0xad, 0x92, 0xa9, 0xb9, 0x76, 0x5e, 0x57, 0xbc, 0x29, 0x90, 0xdc, 0x83, 0x55, 0x3a, 0x1e, 0xcd,
	0xfa, 0x6d, 0x76, 0x6d, 0x3e, 0x49, 0x9e, 0x43, 0x6d, 0x30, 0x8e, 0x22, 0x53, 0xad, 0xe6, 0x23,
	0x19, 0xd0, 0x85, 0x0b, 0x4f, 0xc9, 0xb3, 0x58, 0xf2, 0x5a, 0xcf, 0x7e, 0x46, 0xa1, 0xa4, 0xaa,
	0x9d, 0x9b, 0xa7, 0xc4, 0x5b, 0xa8, 0x74, 0xbf, 0x96, 0xa0, 0xea, 0x99, 0x2a, 0xf2, 0x14, 0x2a,
	0xe6, 0xa5, 0x21, 0x8d, 0x19, 0x97, 0x7a, 0x85, 0x9a, 0x5b, 0xb3, 0xc4, 0xe2, 0x6b, 0x74, 0x04,
	0xf6, 0x09, 0x06, 0xe7, 0x34, 0x59, 0xbb, 0x63, 0xc6, 0xe5, 0x01, 0x1d, 0x30, 0x3e, 0x52, 0x9b,
	0x45, 0x72, 0xfb, 0xd5, 0xcc, 0xf3, 0x4a, 0xf6, 0x81, 0xec, 0x45, 0xff, 0x82, 0xe8, 0x10, 0x56,
	0x8e, 0xfc, 0x78, 0xb6, 0xc4, 0xf9, 0x14, 0xb7, 0xe8, 0x1c, 0x79, 0x0b, 0x6b, 0x47, 0x7e, 0xdc,
	0xa3, 0x32, 0x94, 0x93, 0x63, 0xf3, 0xc8, 0xfc, 0x25, 0xa5, 0x07, 0x5b, 0xfb, 0x28, 0x33, 0x43,
	0x79, 0x3b, 0xc3, 0x7f, 0x1a, 0x6b, 0xf2, 0x11, 0xb6, 0xf7, 0x51, 0xee, 0x45, 0x51, 0x0e, 0xed,
	0x9d, 0x3c, 0x5a, 0x33, 0xf2, 0xb9, 0xed, 0x7c, 0xe9, 0x7e, 0x7a, 0x74, 0xab, 0x8f, 0xdf, 0x90,
	0xa5, 0xdf, 0xbe, 0x2f, 0x4b, 0x2a, 0xf5, 0xe4, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8f, 0x0c,
	0x18, 0xb7, 0x37, 0x07, 0x00, 0x00,
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
	AllPortInformation(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*networkelement.Element, error)
	MapInterface(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*NetworkElementInterfaces, error)
	MapEntityPhysical(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*NetworkElementInterfaces, error)
	GetTransceiverInformation(ctx context.Context, in *NetworkElement, opts ...grpc.CallOption) (*networkelement.Transceiver, error)
	GetAllTransceiverInformation(ctx context.Context, in *NetworkElementWrapper, opts ...grpc.CallOption) (*networkelement.Element, error)
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
type ResourceServer interface {
	Version(context.Context, *Empty) (*VersionResponse, error)
	TechnicalPortInformation(context.Context, *NetworkElement) (*networkelement.Element, error)
	AllPortInformation(context.Context, *NetworkElement) (*networkelement.Element, error)
	MapInterface(context.Context, *NetworkElement) (*NetworkElementInterfaces, error)
	MapEntityPhysical(context.Context, *NetworkElement) (*NetworkElementInterfaces, error)
	GetTransceiverInformation(context.Context, *NetworkElement) (*networkelement.Transceiver, error)
	GetAllTransceiverInformation(context.Context, *NetworkElementWrapper) (*networkelement.Element, error)
}

// UnimplementedResourceServer can be embedded to have forward compatible implementations.
type UnimplementedResourceServer struct {
}

func (*UnimplementedResourceServer) Version(ctx context.Context, req *Empty) (*VersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Version not implemented")
}
func (*UnimplementedResourceServer) TechnicalPortInformation(ctx context.Context, req *NetworkElement) (*networkelement.Element, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TechnicalPortInformation not implemented")
}
func (*UnimplementedResourceServer) AllPortInformation(ctx context.Context, req *NetworkElement) (*networkelement.Element, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllPortInformation not implemented")
}
func (*UnimplementedResourceServer) MapInterface(ctx context.Context, req *NetworkElement) (*NetworkElementInterfaces, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MapInterface not implemented")
}
func (*UnimplementedResourceServer) MapEntityPhysical(ctx context.Context, req *NetworkElement) (*NetworkElementInterfaces, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MapEntityPhysical not implemented")
}
func (*UnimplementedResourceServer) GetTransceiverInformation(ctx context.Context, req *NetworkElement) (*networkelement.Transceiver, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransceiverInformation not implemented")
}
func (*UnimplementedResourceServer) GetAllTransceiverInformation(ctx context.Context, req *NetworkElementWrapper) (*networkelement.Element, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllTransceiverInformation not implemented")
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
	Metadata: "git.liero.se/opentelco/go-swpx/proto/src/resource.proto",
}