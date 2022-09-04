// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.17.3
// source: kmq/v1/options/mq.proto

package options

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

type MQ struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Subscribes []*MQ_Group `protobuf:"bytes,1,rep,name=subscribes,proto3" json:"subscribes,omitempty"`
}

func (x *MQ) Reset() {
	*x = MQ{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kmq_v1_options_mq_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MQ) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MQ) ProtoMessage() {}

func (x *MQ) ProtoReflect() protoreflect.Message {
	mi := &file_kmq_v1_options_mq_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MQ.ProtoReflect.Descriptor instead.
func (*MQ) Descriptor() ([]byte, []int) {
	return file_kmq_v1_options_mq_proto_rawDescGZIP(), []int{0}
}

func (x *MQ) GetSubscribes() []*MQ_Group {
	if x != nil {
		return x.Subscribes
	}
	return nil
}

type MQ_Group struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic   string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Channel string `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`
}

func (x *MQ_Group) Reset() {
	*x = MQ_Group{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kmq_v1_options_mq_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MQ_Group) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MQ_Group) ProtoMessage() {}

func (x *MQ_Group) ProtoReflect() protoreflect.Message {
	mi := &file_kmq_v1_options_mq_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MQ_Group.ProtoReflect.Descriptor instead.
func (*MQ_Group) Descriptor() ([]byte, []int) {
	return file_kmq_v1_options_mq_proto_rawDescGZIP(), []int{0, 0}
}

func (x *MQ_Group) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *MQ_Group) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

var File_kmq_v1_options_mq_proto protoreflect.FileDescriptor

var file_kmq_v1_options_mq_proto_rawDesc = []byte{
	0x0a, 0x17, 0x6b, 0x6d, 0x71, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2f, 0x6d, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x6b, 0x6d, 0x71, 0x2e, 0x76,
	0x31, 0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x77, 0x0a, 0x02, 0x4d, 0x51, 0x12,
	0x38, 0x0a, 0x0a, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6b, 0x6d, 0x71, 0x2e, 0x76, 0x31, 0x2e, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4d, 0x51, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x0a, 0x73,
	0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x73, 0x1a, 0x37, 0x0a, 0x05, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x42, 0x45, 0x5a, 0x43, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x43, 0x63, 0x68, 0x65, 0x65, 0x72, 0x73, 0x2f, 0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73, 0x2d,
	0x6d, 0x71, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6d, 0x71,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x3b, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_kmq_v1_options_mq_proto_rawDescOnce sync.Once
	file_kmq_v1_options_mq_proto_rawDescData = file_kmq_v1_options_mq_proto_rawDesc
)

func file_kmq_v1_options_mq_proto_rawDescGZIP() []byte {
	file_kmq_v1_options_mq_proto_rawDescOnce.Do(func() {
		file_kmq_v1_options_mq_proto_rawDescData = protoimpl.X.CompressGZIP(file_kmq_v1_options_mq_proto_rawDescData)
	})
	return file_kmq_v1_options_mq_proto_rawDescData
}

var file_kmq_v1_options_mq_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_kmq_v1_options_mq_proto_goTypes = []interface{}{
	(*MQ)(nil),       // 0: kmq.v1.options.MQ
	(*MQ_Group)(nil), // 1: kmq.v1.options.MQ.Group
}
var file_kmq_v1_options_mq_proto_depIdxs = []int32{
	1, // 0: kmq.v1.options.MQ.subscribes:type_name -> kmq.v1.options.MQ.Group
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_kmq_v1_options_mq_proto_init() }
func file_kmq_v1_options_mq_proto_init() {
	if File_kmq_v1_options_mq_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kmq_v1_options_mq_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MQ); i {
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
		file_kmq_v1_options_mq_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MQ_Group); i {
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
			RawDescriptor: file_kmq_v1_options_mq_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_kmq_v1_options_mq_proto_goTypes,
		DependencyIndexes: file_kmq_v1_options_mq_proto_depIdxs,
		MessageInfos:      file_kmq_v1_options_mq_proto_msgTypes,
	}.Build()
	File_kmq_v1_options_mq_proto = out.File
	file_kmq_v1_options_mq_proto_rawDesc = nil
	file_kmq_v1_options_mq_proto_goTypes = nil
	file_kmq_v1_options_mq_proto_depIdxs = nil
}
