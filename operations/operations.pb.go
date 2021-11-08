// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: operations/operations.proto

package operations

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

type JoinMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cluster []string `protobuf:"bytes,1,rep,name=cluster,proto3" json:"cluster,omitempty"`
}

func (x *JoinMessage) Reset() {
	*x = JoinMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_operations_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinMessage) ProtoMessage() {}

func (x *JoinMessage) ProtoReflect() protoreflect.Message {
	mi := &file_operations_operations_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinMessage.ProtoReflect.Descriptor instead.
func (*JoinMessage) Descriptor() ([]byte, []int) {
	return file_operations_operations_proto_rawDescGZIP(), []int{0}
}

func (x *JoinMessage) GetCluster() []string {
	if x != nil {
		return x.Cluster
	}
	return nil
}

type Key struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *Key) Reset() {
	*x = Key{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_operations_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Key) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Key) ProtoMessage() {}

func (x *Key) ProtoReflect() protoreflect.Message {
	mi := &file_operations_operations_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Key.ProtoReflect.Descriptor instead.
func (*Key) Descriptor() ([]byte, []int) {
	return file_operations_operations_proto_rawDescGZIP(), []int{1}
}

func (x *Key) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

type Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value [][]byte `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
}

func (x *Value) Reset() {
	*x = Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_operations_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_operations_operations_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_operations_operations_proto_rawDescGZIP(), []int{2}
}

func (x *Value) GetValue() [][]byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type KeyValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value [][]byte `protobuf:"bytes,2,rep,name=value,proto3" json:"value,omitempty"`
}

func (x *KeyValue) Reset() {
	*x = KeyValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_operations_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValue) ProtoMessage() {}

func (x *KeyValue) ProtoReflect() protoreflect.Message {
	mi := &file_operations_operations_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValue.ProtoReflect.Descriptor instead.
func (*KeyValue) Descriptor() ([]byte, []int) {
	return file_operations_operations_proto_rawDescGZIP(), []int{3}
}

func (x *KeyValue) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *KeyValue) GetValue() [][]byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type KeyValueVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key     []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value   [][]byte `protobuf:"bytes,2,rep,name=value,proto3" json:"value,omitempty"`
	Version uint64   `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *KeyValueVersion) Reset() {
	*x = KeyValueVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_operations_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValueVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValueVersion) ProtoMessage() {}

func (x *KeyValueVersion) ProtoReflect() protoreflect.Message {
	mi := &file_operations_operations_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValueVersion.ProtoReflect.Descriptor instead.
func (*KeyValueVersion) Descriptor() ([]byte, []int) {
	return file_operations_operations_proto_rawDescGZIP(), []int{4}
}

func (x *KeyValueVersion) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *KeyValueVersion) GetValue() [][]byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *KeyValueVersion) GetVersion() uint64 {
	if x != nil {
		return x.Version
	}
	return 0
}

type Ack struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *Ack) Reset() {
	*x = Ack{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_operations_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ack) ProtoMessage() {}

func (x *Ack) ProtoReflect() protoreflect.Message {
	mi := &file_operations_operations_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ack.ProtoReflect.Descriptor instead.
func (*Ack) Descriptor() ([]byte, []int) {
	return file_operations_operations_proto_rawDescGZIP(), []int{5}
}

func (x *Ack) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type KeyCost struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key  []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Cost uint64 `protobuf:"varint,2,opt,name=cost,proto3" json:"cost,omitempty"`
}

func (x *KeyCost) Reset() {
	*x = KeyCost{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_operations_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyCost) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyCost) ProtoMessage() {}

func (x *KeyCost) ProtoReflect() protoreflect.Message {
	mi := &file_operations_operations_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyCost.ProtoReflect.Descriptor instead.
func (*KeyCost) Descriptor() ([]byte, []int) {
	return file_operations_operations_proto_rawDescGZIP(), []int{6}
}

func (x *KeyCost) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *KeyCost) GetCost() uint64 {
	if x != nil {
		return x.Cost
	}
	return 0
}

type Outcome struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Out     bool     `protobuf:"varint,1,opt,name=out,proto3" json:"out,omitempty"`
	Value   [][]byte `protobuf:"bytes,2,rep,name=value,proto3" json:"value,omitempty"`
	Version uint64   `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *Outcome) Reset() {
	*x = Outcome{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_operations_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Outcome) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Outcome) ProtoMessage() {}

func (x *Outcome) ProtoReflect() protoreflect.Message {
	mi := &file_operations_operations_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Outcome.ProtoReflect.Descriptor instead.
func (*Outcome) Descriptor() ([]byte, []int) {
	return file_operations_operations_proto_rawDescGZIP(), []int{7}
}

func (x *Outcome) GetOut() bool {
	if x != nil {
		return x.Out
	}
	return false
}

func (x *Outcome) GetValue() [][]byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *Outcome) GetVersion() uint64 {
	if x != nil {
		return x.Version
	}
	return 0
}

type PingMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PingMessage) Reset() {
	*x = PingMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_operations_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingMessage) ProtoMessage() {}

func (x *PingMessage) ProtoReflect() protoreflect.Message {
	mi := &file_operations_operations_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingMessage.ProtoReflect.Descriptor instead.
func (*PingMessage) Descriptor() ([]byte, []int) {
	return file_operations_operations_proto_rawDescGZIP(), []int{8}
}

var File_operations_operations_proto protoreflect.FileDescriptor

var file_operations_operations_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x27, 0x0a, 0x0b, 0x4a, 0x6f, 0x69,
	0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x22, 0x17, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x1d, 0x0a, 0x05, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x32, 0x0a, 0x08, 0x4b, 0x65,
	0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x53,
	0x0a, 0x0f, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x22, 0x17, 0x0a, 0x03, 0x41, 0x63, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x2f, 0x0a, 0x07,
	0x4b, 0x65, 0x79, 0x43, 0x6f, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x73,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x63, 0x6f, 0x73, 0x74, 0x22, 0x4b, 0x0a,
	0x07, 0x4f, 0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6f, 0x75, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x6f, 0x75, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x0d, 0x0a, 0x0b, 0x50, 0x69,
	0x6e, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0xba, 0x05, 0x0a, 0x0a, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x2b, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12,
	0x0f, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4b, 0x65, 0x79,
	0x1a, 0x11, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x22, 0x00, 0x12, 0x2e, 0x0a, 0x03, 0x50, 0x75, 0x74, 0x12, 0x14, 0x2e, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x1a, 0x0f, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x41, 0x63, 0x6b, 0x22, 0x00, 0x12, 0x31, 0x0a, 0x06, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x12,
	0x14, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4b, 0x65, 0x79,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x0f, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x41, 0x63, 0x6b, 0x22, 0x00, 0x12, 0x29, 0x0a, 0x03, 0x44, 0x65, 0x6c, 0x12,
	0x0f, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4b, 0x65, 0x79,
	0x1a, 0x0f, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x41, 0x63,
	0x6b, 0x22, 0x00, 0x12, 0x33, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x12, 0x0f, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x4b, 0x65, 0x79, 0x1a, 0x11, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0b, 0x50, 0x75, 0x74, 0x49,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x12, 0x14, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x0f, 0x2e,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x41, 0x63, 0x6b, 0x22, 0x00,
	0x12, 0x39, 0x0a, 0x0e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x12, 0x14, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x0f, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x41, 0x63, 0x6b, 0x22, 0x00, 0x12, 0x31, 0x0a, 0x0b, 0x44,
	0x65, 0x6c, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x12, 0x0f, 0x2e, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4b, 0x65, 0x79, 0x1a, 0x0f, 0x2e, 0x6f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x41, 0x63, 0x6b, 0x22, 0x00, 0x12, 0x3b,
	0x0a, 0x09, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x1b, 0x2e, 0x6f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x0f, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x41, 0x63, 0x6b, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x12, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x73, 0x12, 0x0f, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4b,
	0x65, 0x79, 0x1a, 0x0f, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x41, 0x63, 0x6b, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x09, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x13, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x4b, 0x65, 0x79, 0x43, 0x6f, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4f, 0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65, 0x22, 0x00, 0x12, 0x32,
	0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x17, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a,
	0x0f, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x41, 0x63, 0x6b,
	0x22, 0x00, 0x12, 0x32, 0x0a, 0x04, 0x4a, 0x6f, 0x69, 0x6e, 0x12, 0x17, 0x2e, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4a, 0x6f, 0x69, 0x6e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x1a, 0x0f, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x41, 0x63, 0x6b, 0x22, 0x00, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x54, 0x68, 0x65, 0x74, 0x61, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x72,
	0x73, 0x2f, 0x53, 0x44, 0x43, 0x43, 0x2f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_operations_operations_proto_rawDescOnce sync.Once
	file_operations_operations_proto_rawDescData = file_operations_operations_proto_rawDesc
)

func file_operations_operations_proto_rawDescGZIP() []byte {
	file_operations_operations_proto_rawDescOnce.Do(func() {
		file_operations_operations_proto_rawDescData = protoimpl.X.CompressGZIP(file_operations_operations_proto_rawDescData)
	})
	return file_operations_operations_proto_rawDescData
}

var file_operations_operations_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_operations_operations_proto_goTypes = []interface{}{
	(*JoinMessage)(nil),     // 0: operations.JoinMessage
	(*Key)(nil),             // 1: operations.Key
	(*Value)(nil),           // 2: operations.Value
	(*KeyValue)(nil),        // 3: operations.KeyValue
	(*KeyValueVersion)(nil), // 4: operations.KeyValueVersion
	(*Ack)(nil),             // 5: operations.Ack
	(*KeyCost)(nil),         // 6: operations.KeyCost
	(*Outcome)(nil),         // 7: operations.Outcome
	(*PingMessage)(nil),     // 8: operations.PingMessage
}
var file_operations_operations_proto_depIdxs = []int32{
	1,  // 0: operations.Operations.Get:input_type -> operations.Key
	3,  // 1: operations.Operations.Put:input_type -> operations.KeyValue
	3,  // 2: operations.Operations.Append:input_type -> operations.KeyValue
	1,  // 3: operations.Operations.Del:input_type -> operations.Key
	1,  // 4: operations.Operations.GetInternal:input_type -> operations.Key
	3,  // 5: operations.Operations.PutInternal:input_type -> operations.KeyValue
	3,  // 6: operations.Operations.AppendInternal:input_type -> operations.KeyValue
	1,  // 7: operations.Operations.DelInternal:input_type -> operations.Key
	4,  // 8: operations.Operations.Replicate:input_type -> operations.KeyValueVersion
	1,  // 9: operations.Operations.DeleteFromReplicas:input_type -> operations.Key
	6,  // 10: operations.Operations.Migration:input_type -> operations.KeyCost
	8,  // 11: operations.Operations.Ping:input_type -> operations.PingMessage
	0,  // 12: operations.Operations.Join:input_type -> operations.JoinMessage
	2,  // 13: operations.Operations.Get:output_type -> operations.Value
	5,  // 14: operations.Operations.Put:output_type -> operations.Ack
	5,  // 15: operations.Operations.Append:output_type -> operations.Ack
	5,  // 16: operations.Operations.Del:output_type -> operations.Ack
	2,  // 17: operations.Operations.GetInternal:output_type -> operations.Value
	5,  // 18: operations.Operations.PutInternal:output_type -> operations.Ack
	5,  // 19: operations.Operations.AppendInternal:output_type -> operations.Ack
	5,  // 20: operations.Operations.DelInternal:output_type -> operations.Ack
	5,  // 21: operations.Operations.Replicate:output_type -> operations.Ack
	5,  // 22: operations.Operations.DeleteFromReplicas:output_type -> operations.Ack
	7,  // 23: operations.Operations.Migration:output_type -> operations.Outcome
	5,  // 24: operations.Operations.Ping:output_type -> operations.Ack
	5,  // 25: operations.Operations.Join:output_type -> operations.Ack
	13, // [13:26] is the sub-list for method output_type
	0,  // [0:13] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_operations_operations_proto_init() }
func file_operations_operations_proto_init() {
	if File_operations_operations_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_operations_operations_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinMessage); i {
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
		file_operations_operations_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Key); i {
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
		file_operations_operations_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Value); i {
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
		file_operations_operations_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValue); i {
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
		file_operations_operations_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValueVersion); i {
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
		file_operations_operations_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ack); i {
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
		file_operations_operations_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyCost); i {
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
		file_operations_operations_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Outcome); i {
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
		file_operations_operations_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingMessage); i {
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
			RawDescriptor: file_operations_operations_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_operations_operations_proto_goTypes,
		DependencyIndexes: file_operations_operations_proto_depIdxs,
		MessageInfos:      file_operations_operations_proto_msgTypes,
	}.Build()
	File_operations_operations_proto = out.File
	file_operations_operations_proto_rawDesc = nil
	file_operations_operations_proto_goTypes = nil
	file_operations_operations_proto_depIdxs = nil
}
