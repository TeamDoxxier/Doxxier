// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.20.3
// source: direct_message.proto

package cmixx

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

type DirectMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PayloadType    uint32 `protobuf:"varint,1,opt,name=PayloadType,proto3" json:"PayloadType,omitempty"`
	DmToken        uint32 `protobuf:"varint,2,opt,name=DmToken,proto3" json:"DmToken,omitempty"`
	RoundId        uint64 `protobuf:"varint,3,opt,name=RoundId,proto3" json:"RoundId,omitempty"`
	Payload        []byte `protobuf:"bytes,4,opt,name=Payload,proto3" json:"Payload,omitempty"`
	LocalTimestamp int64  `protobuf:"varint,5,opt,name=LocalTimestamp,proto3" json:"LocalTimestamp,omitempty"`
	Nonce          []byte `protobuf:"bytes,6,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
}

func (x *DirectMessage) Reset() {
	*x = DirectMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_direct_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DirectMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DirectMessage) ProtoMessage() {}

func (x *DirectMessage) ProtoReflect() protoreflect.Message {
	mi := &file_direct_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DirectMessage.ProtoReflect.Descriptor instead.
func (*DirectMessage) Descriptor() ([]byte, []int) {
	return file_direct_message_proto_rawDescGZIP(), []int{0}
}

func (x *DirectMessage) GetPayloadType() uint32 {
	if x != nil {
		return x.PayloadType
	}
	return 0
}

func (x *DirectMessage) GetDmToken() uint32 {
	if x != nil {
		return x.DmToken
	}
	return 0
}

func (x *DirectMessage) GetRoundId() uint64 {
	if x != nil {
		return x.RoundId
	}
	return 0
}

func (x *DirectMessage) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *DirectMessage) GetLocalTimestamp() int64 {
	if x != nil {
		return x.LocalTimestamp
	}
	return 0
}

func (x *DirectMessage) GetNonce() []byte {
	if x != nil {
		return x.Nonce
	}
	return nil
}

var File_direct_message_proto protoreflect.FileDescriptor

var file_direct_message_proto_rawDesc = []byte{
	0x0a, 0x14, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x63, 0x6d, 0x69, 0x78, 0x78, 0x22, 0xbd, 0x01,
	0x0a, 0x0d, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x20, 0x0a, 0x0b, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x44, 0x6d, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x07, 0x44, 0x6d, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x52,
	0x6f, 0x75, 0x6e, 0x64, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x52, 0x6f,
	0x75, 0x6e, 0x64, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12,
	0x26, 0x0a, 0x0e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x4e, 0x6f, 0x6e, 0x63, 0x65,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x42, 0x20, 0x5a,
	0x1e, 0x64, 0x6f, 0x78, 0x78, 0x69, 0x65, 0x72, 0x2e, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x6c, 0x69,
	0x62, 0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x63, 0x6d, 0x69, 0x78, 0x78, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_direct_message_proto_rawDescOnce sync.Once
	file_direct_message_proto_rawDescData = file_direct_message_proto_rawDesc
)

func file_direct_message_proto_rawDescGZIP() []byte {
	file_direct_message_proto_rawDescOnce.Do(func() {
		file_direct_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_direct_message_proto_rawDescData)
	})
	return file_direct_message_proto_rawDescData
}

var file_direct_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_direct_message_proto_goTypes = []any{
	(*DirectMessage)(nil), // 0: cmixx.DirectMessage
}
var file_direct_message_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_direct_message_proto_init() }
func file_direct_message_proto_init() {
	if File_direct_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_direct_message_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*DirectMessage); i {
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
			RawDescriptor: file_direct_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_direct_message_proto_goTypes,
		DependencyIndexes: file_direct_message_proto_depIdxs,
		MessageInfos:      file_direct_message_proto_msgTypes,
	}.Build()
	File_direct_message_proto = out.File
	file_direct_message_proto_rawDesc = nil
	file_direct_message_proto_goTypes = nil
	file_direct_message_proto_depIdxs = nil
}
