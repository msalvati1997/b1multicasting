// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: registry.proto

package protoregistry

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

type MulticastType int32

const (
	MulticastType_BMULTICAST   MulticastType = 0
	MulticastType_TOCMULTICAST MulticastType = 1
	MulticastType_TODMULTICAST MulticastType = 2
	MulticastType_COMULTICAST  MulticastType = 3
)

// Enum value maps for MulticastType.
var (
	MulticastType_name = map[int32]string{
		0: "BMULTICAST",
		1: "TOCMULTICAST",
		2: "TODMULTICAST",
		3: "COMULTICAST",
	}
	MulticastType_value = map[string]int32{
		"BMULTICAST":   0,
		"TOCMULTICAST": 1,
		"TODMULTICAST": 2,
		"COMULTICAST":  3,
	}
)

func (x MulticastType) Enum() *MulticastType {
	p := new(MulticastType)
	*p = x
	return p
}

func (x MulticastType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MulticastType) Descriptor() protoreflect.EnumDescriptor {
	return file_registry_proto_enumTypes[0].Descriptor()
}

func (MulticastType) Type() protoreflect.EnumType {
	return &file_registry_proto_enumTypes[0]
}

func (x MulticastType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MulticastType.Descriptor instead.
func (MulticastType) EnumDescriptor() ([]byte, []int) {
	return file_registry_proto_rawDescGZIP(), []int{0}
}

type Status int32

const (
	Status_OPENING  Status = 0
	Status_STARTING Status = 1
	Status_ACTIVE   Status = 2
	Status_CLOSING  Status = 3
	Status_CLOSED   Status = 4
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0: "OPENING",
		1: "STARTING",
		2: "ACTIVE",
		3: "CLOSING",
		4: "CLOSED",
	}
	Status_value = map[string]int32{
		"OPENING":  0,
		"STARTING": 1,
		"ACTIVE":   2,
		"CLOSING":  3,
		"CLOSED":   4,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_registry_proto_enumTypes[1].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_registry_proto_enumTypes[1]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_registry_proto_rawDescGZIP(), []int{1}
}

type Rinfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MulticastId   string        `protobuf:"bytes,1,opt,name=multicastId,proto3" json:"multicastId,omitempty"`
	MulticastType MulticastType `protobuf:"varint,2,opt,name=multicastType,proto3,enum=registryservice.MulticastType" json:"multicastType,omitempty"`
	ClientPort    uint32        `protobuf:"varint,3,opt,name=clientPort,proto3" json:"clientPort,omitempty"`
}

func (x *Rinfo) Reset() {
	*x = Rinfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rinfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rinfo) ProtoMessage() {}

func (x *Rinfo) ProtoReflect() protoreflect.Message {
	mi := &file_registry_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rinfo.ProtoReflect.Descriptor instead.
func (*Rinfo) Descriptor() ([]byte, []int) {
	return file_registry_proto_rawDescGZIP(), []int{0}
}

func (x *Rinfo) GetMulticastId() string {
	if x != nil {
		return x.MulticastId
	}
	return ""
}

func (x *Rinfo) GetMulticastType() MulticastType {
	if x != nil {
		return x.MulticastType
	}
	return MulticastType_BMULTICAST
}

func (x *Rinfo) GetClientPort() uint32 {
	if x != nil {
		return x.ClientPort
	}
	return 0
}

type Ranswer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId  string  `protobuf:"bytes,1,opt,name=clientId,proto3" json:"clientId,omitempty"`
	GroupInfo *MGroup `protobuf:"bytes,2,opt,name=groupInfo,proto3" json:"groupInfo,omitempty"`
}

func (x *Ranswer) Reset() {
	*x = Ranswer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ranswer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ranswer) ProtoMessage() {}

func (x *Ranswer) ProtoReflect() protoreflect.Message {
	mi := &file_registry_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ranswer.ProtoReflect.Descriptor instead.
func (*Ranswer) Descriptor() ([]byte, []int) {
	return file_registry_proto_rawDescGZIP(), []int{1}
}

func (x *Ranswer) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *Ranswer) GetGroupInfo() *MGroup {
	if x != nil {
		return x.GroupInfo
	}
	return nil
}

type MGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MulticastId   string                 `protobuf:"bytes,1,opt,name=multicastId,proto3" json:"multicastId,omitempty"`
	MulticastType MulticastType          `protobuf:"varint,2,opt,name=multicastType,proto3,enum=registryservice.MulticastType" json:"multicastType,omitempty"`
	Status        Status                 `protobuf:"varint,3,opt,name=status,proto3,enum=registryservice.Status" json:"status,omitempty"`
	ReadyMembers  uint64                 `protobuf:"varint,4,opt,name=readyMembers,proto3" json:"readyMembers,omitempty"`
	Members       map[string]*MemberInfo `protobuf:"bytes,5,rep,name=members,proto3" json:"members,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *MGroup) Reset() {
	*x = MGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MGroup) ProtoMessage() {}

func (x *MGroup) ProtoReflect() protoreflect.Message {
	mi := &file_registry_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MGroup.ProtoReflect.Descriptor instead.
func (*MGroup) Descriptor() ([]byte, []int) {
	return file_registry_proto_rawDescGZIP(), []int{2}
}

func (x *MGroup) GetMulticastId() string {
	if x != nil {
		return x.MulticastId
	}
	return ""
}

func (x *MGroup) GetMulticastType() MulticastType {
	if x != nil {
		return x.MulticastType
	}
	return MulticastType_BMULTICAST
}

func (x *MGroup) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_OPENING
}

func (x *MGroup) GetReadyMembers() uint64 {
	if x != nil {
		return x.ReadyMembers
	}
	return 0
}

func (x *MGroup) GetMembers() map[string]*MemberInfo {
	if x != nil {
		return x.Members
	}
	return nil
}

type MulticastId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MulticastId string `protobuf:"bytes,1,opt,name=multicastId,proto3" json:"multicastId,omitempty"`
}

func (x *MulticastId) Reset() {
	*x = MulticastId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MulticastId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MulticastId) ProtoMessage() {}

func (x *MulticastId) ProtoReflect() protoreflect.Message {
	mi := &file_registry_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MulticastId.ProtoReflect.Descriptor instead.
func (*MulticastId) Descriptor() ([]byte, []int) {
	return file_registry_proto_rawDescGZIP(), []int{3}
}

func (x *MulticastId) GetMulticastId() string {
	if x != nil {
		return x.MulticastId
	}
	return ""
}

type RequestData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MulticastId string `protobuf:"bytes,1,opt,name=multicastId,proto3" json:"multicastId,omitempty"`
	MId         string `protobuf:"bytes,2,opt,name=mId,proto3" json:"mId,omitempty"`
}

func (x *RequestData) Reset() {
	*x = RequestData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestData) ProtoMessage() {}

func (x *RequestData) ProtoReflect() protoreflect.Message {
	mi := &file_registry_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestData.ProtoReflect.Descriptor instead.
func (*RequestData) Descriptor() ([]byte, []int) {
	return file_registry_proto_rawDescGZIP(), []int{4}
}

func (x *RequestData) GetMulticastId() string {
	if x != nil {
		return x.MulticastId
	}
	return ""
}

func (x *RequestData) GetMId() string {
	if x != nil {
		return x.MId
	}
	return ""
}

type MemberInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Ready   bool   `protobuf:"varint,3,opt,name=ready,proto3" json:"ready,omitempty"`
}

func (x *MemberInfo) Reset() {
	*x = MemberInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberInfo) ProtoMessage() {}

func (x *MemberInfo) ProtoReflect() protoreflect.Message {
	mi := &file_registry_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberInfo.ProtoReflect.Descriptor instead.
func (*MemberInfo) Descriptor() ([]byte, []int) {
	return file_registry_proto_rawDescGZIP(), []int{5}
}

func (x *MemberInfo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MemberInfo) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *MemberInfo) GetReady() bool {
	if x != nil {
		return x.Ready
	}
	return false
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_registry_proto_msgTypes[6]
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
	return file_registry_proto_rawDescGZIP(), []int{6}
}

var File_registry_proto protoreflect.FileDescriptor

var file_registry_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x22, 0x8f, 0x01, 0x0a, 0x05, 0x52, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x20, 0x0a, 0x0b, 0x6d,
	0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x49, 0x64, 0x12, 0x44, 0x0a,
	0x0d, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x0d, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x50, 0x6f, 0x72,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x50,
	0x6f, 0x72, 0x74, 0x22, 0x5c, 0x0a, 0x07, 0x52, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x12, 0x1a,
	0x0a, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x35, 0x0a, 0x09, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x4d, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x09, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x66,
	0x6f, 0x22, 0xde, 0x02, 0x0a, 0x06, 0x4d, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x20, 0x0a, 0x0b,
	0x6d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x49, 0x64, 0x12, 0x44,
	0x0a, 0x0d, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0d, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x2f, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x65, 0x61, 0x64, 0x79, 0x4d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x72, 0x65, 0x61,
	0x64, 0x79, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12, 0x3e, 0x0a, 0x07, 0x6d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x1a, 0x57, 0x0a, 0x0c, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x31, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x2f, 0x0a, 0x0b, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x49,
	0x64, 0x12, 0x20, 0x0a, 0x0b, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73,
	0x74, 0x49, 0x64, 0x22, 0x41, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x20, 0x0a, 0x0b, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61,
	0x73, 0x74, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6d, 0x49, 0x64, 0x22, 0x4c, 0x0a, 0x0a, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72,
	0x65, 0x61, 0x64, 0x79, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x2a, 0x54, 0x0a,
	0x0d, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x63, 0x61, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e,
	0x0a, 0x0a, 0x42, 0x4d, 0x55, 0x4c, 0x54, 0x49, 0x43, 0x41, 0x53, 0x54, 0x10, 0x00, 0x12, 0x10,
	0x0a, 0x0c, 0x54, 0x4f, 0x43, 0x4d, 0x55, 0x4c, 0x54, 0x49, 0x43, 0x41, 0x53, 0x54, 0x10, 0x01,
	0x12, 0x10, 0x0a, 0x0c, 0x54, 0x4f, 0x44, 0x4d, 0x55, 0x4c, 0x54, 0x49, 0x43, 0x41, 0x53, 0x54,
	0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x4f, 0x4d, 0x55, 0x4c, 0x54, 0x49, 0x43, 0x41, 0x53,
	0x54, 0x10, 0x03, 0x2a, 0x48, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a,
	0x07, 0x4f, 0x50, 0x45, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x53, 0x54,
	0x41, 0x52, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x43, 0x54, 0x49,
	0x56, 0x45, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4c, 0x4f, 0x53, 0x49, 0x4e, 0x47, 0x10,
	0x03, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x4c, 0x4f, 0x53, 0x45, 0x44, 0x10, 0x04, 0x32, 0xe0, 0x02,
	0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x12, 0x3e, 0x0a, 0x08, 0x72, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x69, 0x6e, 0x66, 0x6f, 0x1a, 0x18,
	0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x52, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x22, 0x00, 0x12, 0x45, 0x0a, 0x0a, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x1c, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x79, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x17, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x22,
	0x00, 0x12, 0x40, 0x0a, 0x05, 0x72, 0x65, 0x61, 0x64, 0x79, 0x12, 0x1c, 0x2e, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x17, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x79, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x22, 0x00, 0x12, 0x45, 0x0a, 0x0a, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x12, 0x1c, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x1a,
	0x17, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x4d, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x09, 0x67, 0x65,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1c, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x72, 0x79, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x63,
	0x61, 0x73, 0x74, 0x49, 0x64, 0x1a, 0x17, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x00,
	0x42, 0x11, 0x5a, 0x0f, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_registry_proto_rawDescOnce sync.Once
	file_registry_proto_rawDescData = file_registry_proto_rawDesc
)

func file_registry_proto_rawDescGZIP() []byte {
	file_registry_proto_rawDescOnce.Do(func() {
		file_registry_proto_rawDescData = protoimpl.X.CompressGZIP(file_registry_proto_rawDescData)
	})
	return file_registry_proto_rawDescData
}

var file_registry_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_registry_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_registry_proto_goTypes = []interface{}{
	(MulticastType)(0),  // 0: registryservice.MulticastType
	(Status)(0),         // 1: registryservice.Status
	(*Rinfo)(nil),       // 2: registryservice.Rinfo
	(*Ranswer)(nil),     // 3: registryservice.Ranswer
	(*MGroup)(nil),      // 4: registryservice.MGroup
	(*MulticastId)(nil), // 5: registryservice.MulticastId
	(*RequestData)(nil), // 6: registryservice.RequestData
	(*MemberInfo)(nil),  // 7: registryservice.MemberInfo
	(*Empty)(nil),       // 8: registryservice.Empty
	nil,                 // 9: registryservice.MGroup.MembersEntry
}
var file_registry_proto_depIdxs = []int32{
	0,  // 0: registryservice.Rinfo.multicastType:type_name -> registryservice.MulticastType
	4,  // 1: registryservice.Ranswer.groupInfo:type_name -> registryservice.MGroup
	0,  // 2: registryservice.MGroup.multicastType:type_name -> registryservice.MulticastType
	1,  // 3: registryservice.MGroup.status:type_name -> registryservice.Status
	9,  // 4: registryservice.MGroup.members:type_name -> registryservice.MGroup.MembersEntry
	7,  // 5: registryservice.MGroup.MembersEntry.value:type_name -> registryservice.MemberInfo
	2,  // 6: registryservice.Registry.register:input_type -> registryservice.Rinfo
	6,  // 7: registryservice.Registry.startGroup:input_type -> registryservice.RequestData
	6,  // 8: registryservice.Registry.ready:input_type -> registryservice.RequestData
	6,  // 9: registryservice.Registry.closeGroup:input_type -> registryservice.RequestData
	5,  // 10: registryservice.Registry.getStatus:input_type -> registryservice.MulticastId
	3,  // 11: registryservice.Registry.register:output_type -> registryservice.Ranswer
	4,  // 12: registryservice.Registry.startGroup:output_type -> registryservice.MGroup
	4,  // 13: registryservice.Registry.ready:output_type -> registryservice.MGroup
	4,  // 14: registryservice.Registry.closeGroup:output_type -> registryservice.MGroup
	4,  // 15: registryservice.Registry.getStatus:output_type -> registryservice.MGroup
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_registry_proto_init() }
func file_registry_proto_init() {
	if File_registry_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_registry_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rinfo); i {
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
		file_registry_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ranswer); i {
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
		file_registry_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MGroup); i {
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
		file_registry_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MulticastId); i {
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
		file_registry_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestData); i {
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
		file_registry_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemberInfo); i {
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
		file_registry_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_registry_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_registry_proto_goTypes,
		DependencyIndexes: file_registry_proto_depIdxs,
		EnumInfos:         file_registry_proto_enumTypes,
		MessageInfos:      file_registry_proto_msgTypes,
	}.Build()
	File_registry_proto = out.File
	file_registry_proto_rawDesc = nil
	file_registry_proto_goTypes = nil
	file_registry_proto_depIdxs = nil
}