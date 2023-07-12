// notifications

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: fleet_notification.proto

package notificationpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TaskQueue int32

const (
	TaskQueue_TASK_QUEUE_FLEET_NOTIFICATIONS TaskQueue = 0
)

// Enum value maps for TaskQueue.
var (
	TaskQueue_name = map[int32]string{
		0: "TASK_QUEUE_FLEET_NOTIFICATIONS",
	}
	TaskQueue_value = map[string]int32{
		"TASK_QUEUE_FLEET_NOTIFICATIONS": 0,
	}
)

func (x TaskQueue) Enum() *TaskQueue {
	p := new(TaskQueue)
	*p = x
	return p
}

func (x TaskQueue) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TaskQueue) Descriptor() protoreflect.EnumDescriptor {
	return file_fleet_notification_proto_enumTypes[0].Descriptor()
}

func (TaskQueue) Type() protoreflect.EnumType {
	return &file_fleet_notification_proto_enumTypes[0]
}

func (x TaskQueue) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TaskQueue.Descriptor instead.
func (TaskQueue) EnumDescriptor() ([]byte, []int) {
	return file_fleet_notification_proto_rawDescGZIP(), []int{0}
}

type Notification_ResourceType int32

const (
	Notification_RESOURCE_TYPE_UNSPECIFIED Notification_ResourceType = 0
	Notification_RESOURCE_TYPE_DEVICE      Notification_ResourceType = 1
	Notification_RESOURCE_TYPE_CONFIG      Notification_ResourceType = 2
)

// Enum value maps for Notification_ResourceType.
var (
	Notification_ResourceType_name = map[int32]string{
		0: "RESOURCE_TYPE_UNSPECIFIED",
		1: "RESOURCE_TYPE_DEVICE",
		2: "RESOURCE_TYPE_CONFIG",
	}
	Notification_ResourceType_value = map[string]int32{
		"RESOURCE_TYPE_UNSPECIFIED": 0,
		"RESOURCE_TYPE_DEVICE":      1,
		"RESOURCE_TYPE_CONFIG":      2,
	}
)

func (x Notification_ResourceType) Enum() *Notification_ResourceType {
	p := new(Notification_ResourceType)
	*p = x
	return p
}

func (x Notification_ResourceType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Notification_ResourceType) Descriptor() protoreflect.EnumDescriptor {
	return file_fleet_notification_proto_enumTypes[1].Descriptor()
}

func (Notification_ResourceType) Type() protoreflect.EnumType {
	return &file_fleet_notification_proto_enumTypes[1]
}

func (x Notification_ResourceType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Notification_ResourceType.Descriptor instead.
func (Notification_ResourceType) EnumDescriptor() ([]byte, []int) {
	return file_fleet_notification_proto_rawDescGZIP(), []int{0, 0}
}

type ListRequest_Filter int32

const (
	ListRequest_NO_LIST_REQUEST_FILTER ListRequest_Filter = 0
	ListRequest_RESOURCE_TYPE_DEVICE   ListRequest_Filter = 1
	ListRequest_RESOURCE_TYPE_CONFIG   ListRequest_Filter = 2
	// include read notifications and unread notifications
	ListRequest_INCLUDE_READ ListRequest_Filter = 3
)

// Enum value maps for ListRequest_Filter.
var (
	ListRequest_Filter_name = map[int32]string{
		0: "NO_LIST_REQUEST_FILTER",
		1: "RESOURCE_TYPE_DEVICE",
		2: "RESOURCE_TYPE_CONFIG",
		3: "INCLUDE_READ",
	}
	ListRequest_Filter_value = map[string]int32{
		"NO_LIST_REQUEST_FILTER": 0,
		"RESOURCE_TYPE_DEVICE":   1,
		"RESOURCE_TYPE_CONFIG":   2,
		"INCLUDE_READ":           3,
	}
)

func (x ListRequest_Filter) Enum() *ListRequest_Filter {
	p := new(ListRequest_Filter)
	*p = x
	return p
}

func (x ListRequest_Filter) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ListRequest_Filter) Descriptor() protoreflect.EnumDescriptor {
	return file_fleet_notification_proto_enumTypes[2].Descriptor()
}

func (ListRequest_Filter) Type() protoreflect.EnumType {
	return &file_fleet_notification_proto_enumTypes[2]
}

func (x ListRequest_Filter) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ListRequest_Filter.Descriptor instead.
func (ListRequest_Filter) EnumDescriptor() ([]byte, []int) {
	return file_fleet_notification_proto_rawDescGZIP(), []int{2, 0}
}

// Notification is a notification triggered by a fleet event
// (e.g. a new device is added to the fleet or a device is removed from the fleet,
// a device is not reachable, or a configuration not being successfully fetched)
type Notification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The ID of the notification
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" bson:"_id"`
	// The title of the notification
	// (e.g. "Device 1234 is not reachable")
	Title string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty" bson:"title"`
	// The ID of the device
	ResourceId string `protobuf:"bytes,3,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty" bson:"resource_id"`
	// the resource type of the notification (e.g. device, config)
	// used to identify the resource_id (e.g. device_id, config_id)
	ResourceType Notification_ResourceType `protobuf:"varint,4,opt,name=resource_type,json=resourceType,proto3,enum=notificationpb.Notification_ResourceType" json:"resource_type,omitempty" bson:"resource_type"`
	// The timestamp of the notification
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty" bson:"timestamp"`
	// The message of the notification e.g "Device 1234 is not reachable since 2020-01-01 12:00:00 UTC,
	// last seen 2020-01-01 11:00:00 UTC, last config update 2020-01-01 10:00:00 UTC"
	Message string `protobuf:"bytes,6,opt,name=message,proto3" json:"message,omitempty" bson:"message"`
	// used to mark the notification as read
	Read bool `protobuf:"varint,7,opt,name=read,proto3" json:"read,omitempty" bson:"read"`
}

func (x *Notification) Reset() {
	*x = Notification{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fleet_notification_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Notification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notification) ProtoMessage() {}

func (x *Notification) ProtoReflect() protoreflect.Message {
	mi := &file_fleet_notification_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Notification.ProtoReflect.Descriptor instead.
func (*Notification) Descriptor() ([]byte, []int) {
	return file_fleet_notification_proto_rawDescGZIP(), []int{0}
}

func (x *Notification) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Notification) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Notification) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *Notification) GetResourceType() Notification_ResourceType {
	if x != nil {
		return x.ResourceType
	}
	return Notification_RESOURCE_TYPE_UNSPECIFIED
}

func (x *Notification) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Notification) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Notification) GetRead() bool {
	if x != nil {
		return x.Read
	}
	return false
}

type GetByIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" bson:"_id"`
}

func (x *GetByIDRequest) Reset() {
	*x = GetByIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fleet_notification_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByIDRequest) ProtoMessage() {}

func (x *GetByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_fleet_notification_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByIDRequest.ProtoReflect.Descriptor instead.
func (*GetByIDRequest) Descriptor() ([]byte, []int) {
	return file_fleet_notification_proto_rawDescGZIP(), []int{1}
}

func (x *GetByIDRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResourceIds []string             `protobuf:"bytes,1,rep,name=resource_ids,json=resourceIds,proto3" json:"resource_ids,omitempty" bson:"resource_ids"`
	Ids         []string             `protobuf:"bytes,2,rep,name=ids,proto3" json:"ids,omitempty" bson:"ids"`
	Filter      []ListRequest_Filter `protobuf:"varint,3,rep,packed,name=filter,proto3,enum=notificationpb.ListRequest_Filter" json:"filter,omitempty" bson:"filter"`
	// how many notifications to return (default 10)
	Limit *int64 `protobuf:"varint,4,opt,name=limit,proto3,oneof" json:"limit,omitempty" bson:"limit"`
	// offset to start from (default 0)
	Offset *int64 `protobuf:"varint,5,opt,name=offset,proto3,oneof" json:"offset,omitempty" bson:"offset"`
}

func (x *ListRequest) Reset() {
	*x = ListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fleet_notification_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRequest) ProtoMessage() {}

func (x *ListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_fleet_notification_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRequest.ProtoReflect.Descriptor instead.
func (*ListRequest) Descriptor() ([]byte, []int) {
	return file_fleet_notification_proto_rawDescGZIP(), []int{2}
}

func (x *ListRequest) GetResourceIds() []string {
	if x != nil {
		return x.ResourceIds
	}
	return nil
}

func (x *ListRequest) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *ListRequest) GetFilter() []ListRequest_Filter {
	if x != nil {
		return x.Filter
	}
	return nil
}

func (x *ListRequest) GetLimit() int64 {
	if x != nil && x.Limit != nil {
		return *x.Limit
	}
	return 0
}

func (x *ListRequest) GetOffset() int64 {
	if x != nil && x.Offset != nil {
		return *x.Offset
	}
	return 0
}

type ListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Notifications []*Notification `protobuf:"bytes,1,rep,name=notifications,proto3" json:"notifications,omitempty" bson:"notifications"`
}

func (x *ListResponse) Reset() {
	*x = ListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fleet_notification_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListResponse) ProtoMessage() {}

func (x *ListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_fleet_notification_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListResponse.ProtoReflect.Descriptor instead.
func (*ListResponse) Descriptor() ([]byte, []int) {
	return file_fleet_notification_proto_rawDescGZIP(), []int{3}
}

func (x *ListResponse) GetNotifications() []*Notification {
	if x != nil {
		return x.Notifications
	}
	return nil
}

type CreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The title of the notification
	// (e.g. "Device 1234 is not reachable")
	Title string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty" bson:"title"`
	// The ID of the device
	ResourceId string `protobuf:"bytes,3,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty" bson:"resource_id"`
	// the resource type of the notification (e.g. device, config)
	// used to identify the resource_id (e.g. device_id, config_id)
	ResourceType Notification_ResourceType `protobuf:"varint,4,opt,name=resource_type,json=resourceType,proto3,enum=notificationpb.Notification_ResourceType" json:"resource_type,omitempty" bson:"resource_type"`
	// The message of the notification e.g "Device 1234 is not reachable since 2020-01-01 12:00:00 UTC,
	// last seen 2020-01-01 11:00:00 UTC, last config update 2020-01-01 10:00:00 UTC"
	Message *string `protobuf:"bytes,6,opt,name=message,proto3,oneof" json:"message,omitempty" bson:"message"`
}

func (x *CreateRequest) Reset() {
	*x = CreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fleet_notification_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRequest) ProtoMessage() {}

func (x *CreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_fleet_notification_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRequest.ProtoReflect.Descriptor instead.
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return file_fleet_notification_proto_rawDescGZIP(), []int{4}
}

func (x *CreateRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *CreateRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *CreateRequest) GetResourceType() Notification_ResourceType {
	if x != nil {
		return x.ResourceType
	}
	return Notification_RESOURCE_TYPE_UNSPECIFIED
}

func (x *CreateRequest) GetMessage() string {
	if x != nil && x.Message != nil {
		return *x.Message
	}
	return ""
}

type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" bson:"_id"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fleet_notification_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_fleet_notification_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_fleet_notification_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" bson:"_id"`
}

func (x *DeleteResponse) Reset() {
	*x = DeleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fleet_notification_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteResponse) ProtoMessage() {}

func (x *DeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_fleet_notification_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteResponse.ProtoReflect.Descriptor instead.
func (*DeleteResponse) Descriptor() ([]byte, []int) {
	return file_fleet_notification_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type MarkAsReadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids []string `protobuf:"bytes,1,rep,name=ids,proto3" json:"ids,omitempty" bson:"ids"`
}

func (x *MarkAsReadRequest) Reset() {
	*x = MarkAsReadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fleet_notification_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarkAsReadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarkAsReadRequest) ProtoMessage() {}

func (x *MarkAsReadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_fleet_notification_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarkAsReadRequest.ProtoReflect.Descriptor instead.
func (*MarkAsReadRequest) Descriptor() ([]byte, []int) {
	return file_fleet_notification_proto_rawDescGZIP(), []int{7}
}

func (x *MarkAsReadRequest) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

type MarkAsReadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Notifications []*Notification `protobuf:"bytes,1,rep,name=notifications,proto3" json:"notifications,omitempty" bson:"notifications"`
}

func (x *MarkAsReadResponse) Reset() {
	*x = MarkAsReadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fleet_notification_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarkAsReadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarkAsReadResponse) ProtoMessage() {}

func (x *MarkAsReadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_fleet_notification_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarkAsReadResponse.ProtoReflect.Descriptor instead.
func (*MarkAsReadResponse) Descriptor() ([]byte, []int) {
	return file_fleet_notification_proto_rawDescGZIP(), []int{8}
}

func (x *MarkAsReadResponse) GetNotifications() []*Notification {
	if x != nil {
		return x.Notifications
	}
	return nil
}

var File_fleet_notification_proto protoreflect.FileDescriptor

var file_fleet_notification_proto_rawDesc = []byte{
	0x0a, 0x18, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x5f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf0, 0x02, 0x0a, 0x0c,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x49, 0x64, 0x12, 0x4e, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x29, 0x2e, 0x6e, 0x6f, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x65, 0x61, 0x64, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x72, 0x65, 0x61, 0x64, 0x22, 0x61, 0x0a, 0x0c, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x19, 0x52,
	0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x52, 0x45,
	0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x45, 0x56, 0x49,
	0x43, 0x45, 0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x10, 0x02, 0x22, 0x20,
	0x0a, 0x0e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x22, 0xb7, 0x02, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x49, 0x64, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x3a, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x2e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65,
	0x72, 0x12, 0x19, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x48, 0x00, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06,
	0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x48, 0x01, 0x52, 0x06,
	0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x88, 0x01, 0x01, 0x22, 0x6a, 0x0a, 0x06, 0x46, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x16, 0x4e, 0x4f, 0x5f, 0x4c, 0x49, 0x53, 0x54, 0x5f, 0x52,
	0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x5f, 0x46, 0x49, 0x4c, 0x54, 0x45, 0x52, 0x10, 0x00, 0x12,
	0x18, 0x0a, 0x14, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x44, 0x45, 0x56, 0x49, 0x43, 0x45, 0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x52, 0x45, 0x53,
	0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x49,
	0x47, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x49, 0x4e, 0x43, 0x4c, 0x55, 0x44, 0x45, 0x5f, 0x52,
	0x45, 0x41, 0x44, 0x10, 0x03, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x42,
	0x09, 0x0a, 0x07, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x52, 0x0a, 0x0c, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x0d, 0x6e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x70, 0x62, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0d, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0xc1,
	0x01, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x4e, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x29,
	0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0x1f, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x20, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x25, 0x0a, 0x11, 0x4d, 0x61, 0x72, 0x6b, 0x41, 0x73, 0x52,
	0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x69, 0x64, 0x73, 0x22, 0x58, 0x0a, 0x12,
	0x4d, 0x61, 0x72, 0x6b, 0x41, 0x73, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x42, 0x0a, 0x0d, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2a, 0x2f, 0x0a, 0x09, 0x54, 0x61, 0x73, 0x6b, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x12, 0x22, 0x0a, 0x1e, 0x54, 0x41, 0x53, 0x4b, 0x5f, 0x51, 0x55, 0x45, 0x55,
	0x45, 0x5f, 0x46, 0x4c, 0x45, 0x45, 0x54, 0x5f, 0x4e, 0x4f, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41,
	0x54, 0x49, 0x4f, 0x4e, 0x53, 0x10, 0x00, 0x32, 0x90, 0x03, 0x0a, 0x13, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x49, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x44, 0x12, 0x1e, 0x2e, 0x6e, 0x6f, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x42,
	0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x6e, 0x6f, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x43, 0x0a, 0x04, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x1b, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x70, 0x62, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1c, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x47, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x49, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x12, 0x1d, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1e, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x0a, 0x4d, 0x61, 0x72, 0x6b, 0x41, 0x73, 0x52, 0x65, 0x61,
	0x64, 0x12, 0x21, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x70, 0x62, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x41, 0x73, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x41, 0x73, 0x52, 0x65, 0x61, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69,
	0x74, 0x2e, 0x6c, 0x69, 0x65, 0x72, 0x6f, 0x2e, 0x73, 0x65, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x74,
	0x65, 0x6c, 0x63, 0x6f, 0x2f, 0x67, 0x6f, 0x2d, 0x73, 0x77, 0x70, 0x78, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2f, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_fleet_notification_proto_rawDescOnce sync.Once
	file_fleet_notification_proto_rawDescData = file_fleet_notification_proto_rawDesc
)

func file_fleet_notification_proto_rawDescGZIP() []byte {
	file_fleet_notification_proto_rawDescOnce.Do(func() {
		file_fleet_notification_proto_rawDescData = protoimpl.X.CompressGZIP(file_fleet_notification_proto_rawDescData)
	})
	return file_fleet_notification_proto_rawDescData
}

var file_fleet_notification_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_fleet_notification_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_fleet_notification_proto_goTypes = []interface{}{
	(TaskQueue)(0),                 // 0: notificationpb.TaskQueue
	(Notification_ResourceType)(0), // 1: notificationpb.Notification.ResourceType
	(ListRequest_Filter)(0),        // 2: notificationpb.ListRequest.Filter
	(*Notification)(nil),           // 3: notificationpb.Notification
	(*GetByIDRequest)(nil),         // 4: notificationpb.GetByIDRequest
	(*ListRequest)(nil),            // 5: notificationpb.ListRequest
	(*ListResponse)(nil),           // 6: notificationpb.ListResponse
	(*CreateRequest)(nil),          // 7: notificationpb.CreateRequest
	(*DeleteRequest)(nil),          // 8: notificationpb.DeleteRequest
	(*DeleteResponse)(nil),         // 9: notificationpb.DeleteResponse
	(*MarkAsReadRequest)(nil),      // 10: notificationpb.MarkAsReadRequest
	(*MarkAsReadResponse)(nil),     // 11: notificationpb.MarkAsReadResponse
	(*timestamppb.Timestamp)(nil),  // 12: google.protobuf.Timestamp
}
var file_fleet_notification_proto_depIdxs = []int32{
	1,  // 0: notificationpb.Notification.resource_type:type_name -> notificationpb.Notification.ResourceType
	12, // 1: notificationpb.Notification.timestamp:type_name -> google.protobuf.Timestamp
	2,  // 2: notificationpb.ListRequest.filter:type_name -> notificationpb.ListRequest.Filter
	3,  // 3: notificationpb.ListResponse.notifications:type_name -> notificationpb.Notification
	1,  // 4: notificationpb.CreateRequest.resource_type:type_name -> notificationpb.Notification.ResourceType
	3,  // 5: notificationpb.MarkAsReadResponse.notifications:type_name -> notificationpb.Notification
	4,  // 6: notificationpb.NotificationService.GetByID:input_type -> notificationpb.GetByIDRequest
	5,  // 7: notificationpb.NotificationService.List:input_type -> notificationpb.ListRequest
	7,  // 8: notificationpb.NotificationService.Create:input_type -> notificationpb.CreateRequest
	8,  // 9: notificationpb.NotificationService.Delete:input_type -> notificationpb.DeleteRequest
	10, // 10: notificationpb.NotificationService.MarkAsRead:input_type -> notificationpb.MarkAsReadRequest
	3,  // 11: notificationpb.NotificationService.GetByID:output_type -> notificationpb.Notification
	6,  // 12: notificationpb.NotificationService.List:output_type -> notificationpb.ListResponse
	3,  // 13: notificationpb.NotificationService.Create:output_type -> notificationpb.Notification
	9,  // 14: notificationpb.NotificationService.Delete:output_type -> notificationpb.DeleteResponse
	11, // 15: notificationpb.NotificationService.MarkAsRead:output_type -> notificationpb.MarkAsReadResponse
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_fleet_notification_proto_init() }
func file_fleet_notification_proto_init() {
	if File_fleet_notification_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fleet_notification_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Notification); i {
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
		file_fleet_notification_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByIDRequest); i {
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
		file_fleet_notification_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRequest); i {
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
		file_fleet_notification_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListResponse); i {
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
		file_fleet_notification_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRequest); i {
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
		file_fleet_notification_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRequest); i {
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
		file_fleet_notification_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteResponse); i {
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
		file_fleet_notification_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarkAsReadRequest); i {
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
		file_fleet_notification_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarkAsReadResponse); i {
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
	file_fleet_notification_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_fleet_notification_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_fleet_notification_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_fleet_notification_proto_goTypes,
		DependencyIndexes: file_fleet_notification_proto_depIdxs,
		EnumInfos:         file_fleet_notification_proto_enumTypes,
		MessageInfos:      file_fleet_notification_proto_msgTypes,
	}.Build()
	File_fleet_notification_proto = out.File
	file_fleet_notification_proto_rawDesc = nil
	file_fleet_notification_proto_goTypes = nil
	file_fleet_notification_proto_depIdxs = nil
}
