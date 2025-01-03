// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: protobuf/project.proto

package billing

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type Tariff struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID    string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	SSD   int64  `protobuf:"varint,3,opt,name=SSD,proto3" json:"SSD,omitempty"`
	CPU   int64  `protobuf:"varint,4,opt,name=CPU,proto3" json:"CPU,omitempty"`
	RAM   int64  `protobuf:"varint,5,opt,name=RAM,proto3" json:"RAM,omitempty"`
	Price int64  `protobuf:"varint,6,opt,name=Price,proto3" json:"Price,omitempty"`
}

func (x *Tariff) Reset() {
	*x = Tariff{}
	mi := &file_protobuf_billing_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Tariff) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tariff) ProtoMessage() {}

func (x *Tariff) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_billing_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tariff.ProtoReflect.Descriptor instead.
func (*Tariff) Descriptor() ([]byte, []int) {
	return file_protobuf_billing_proto_rawDescGZIP(), []int{0}
}

func (x *Tariff) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Tariff) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Tariff) GetSSD() int64 {
	if x != nil {
		return x.SSD
	}
	return 0
}

func (x *Tariff) GetCPU() int64 {
	if x != nil {
		return x.CPU
	}
	return 0
}

func (x *Tariff) GetRAM() int64 {
	if x != nil {
		return x.RAM
	}
	return 0
}

func (x *Tariff) GetPrice() int64 {
	if x != nil {
		return x.Price
	}
	return 0
}

type GetTariffsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tariffs []*Tariff `protobuf:"bytes,1,rep,name=tariffs,proto3" json:"tariffs,omitempty"`
}

func (x *GetTariffsResponse) Reset() {
	*x = GetTariffsResponse{}
	mi := &file_protobuf_billing_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTariffsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTariffsResponse) ProtoMessage() {}

func (x *GetTariffsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_billing_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTariffsResponse.ProtoReflect.Descriptor instead.
func (*GetTariffsResponse) Descriptor() ([]byte, []int) {
	return file_protobuf_billing_proto_rawDescGZIP(), []int{1}
}

func (x *GetTariffsResponse) GetTariffs() []*Tariff {
	if x != nil {
		return x.Tariffs
	}
	return nil
}

var File_protobuf_billing_proto protoreflect.FileDescriptor

var file_protobuf_billing_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x62, 0x69, 0x6c, 0x6c, 0x69,
	0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x78, 0x0a, 0x06, 0x54, 0x61, 0x72, 0x69,
	0x66, 0x66, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x53, 0x53, 0x44, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x03, 0x53, 0x53, 0x44, 0x12, 0x10, 0x0a, 0x03, 0x43, 0x50, 0x55, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x43, 0x50, 0x55, 0x12, 0x10, 0x0a, 0x03, 0x52, 0x41,
	0x4d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x52, 0x41, 0x4d, 0x12, 0x14, 0x0a, 0x05,
	0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x50, 0x72, 0x69,
	0x63, 0x65, 0x22, 0x3b, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x54, 0x61, 0x72, 0x69, 0x66, 0x66, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x07, 0x74, 0x61, 0x72, 0x69,
	0x66, 0x66, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x54, 0x61, 0x72, 0x69, 0x66, 0x66, 0x52, 0x07, 0x74, 0x61, 0x72, 0x69, 0x66, 0x66, 0x73, 0x32,
	0x69, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5a,
	0x0a, 0x0a, 0x47, 0x65, 0x74, 0x54, 0x61, 0x72, 0x69, 0x66, 0x66, 0x73, 0x12, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x61,
	0x72, 0x69, 0x66, 0x66, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1b, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x15, 0x12, 0x13, 0x2f, 0x76, 0x31, 0x2f, 0x62, 0x69, 0x6c, 0x6c, 0x69,
	0x6e, 0x67, 0x2f, 0x74, 0x61, 0x72, 0x69, 0x66, 0x66, 0x73, 0x42, 0x11, 0x5a, 0x0f, 0x70, 0x6b,
	0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x69, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protobuf_billing_proto_rawDescOnce sync.Once
	file_protobuf_billing_proto_rawDescData = file_protobuf_billing_proto_rawDesc
)

func file_protobuf_billing_proto_rawDescGZIP() []byte {
	file_protobuf_billing_proto_rawDescOnce.Do(func() {
		file_protobuf_billing_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobuf_billing_proto_rawDescData)
	})
	return file_protobuf_billing_proto_rawDescData
}

var file_protobuf_billing_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_protobuf_billing_proto_goTypes = []any{
	(*Tariff)(nil),             // 0: api.Tariff
	(*GetTariffsResponse)(nil), // 1: api.GetTariffsResponse
	(*emptypb.Empty)(nil),      // 2: google.protobuf.Empty
}
var file_protobuf_billing_proto_depIdxs = []int32{
	0, // 0: api.GetTariffsResponse.tariffs:type_name -> api.Tariff
	2, // 1: api.UserService.GetTariffs:input_type -> google.protobuf.Empty
	1, // 2: api.UserService.GetTariffs:output_type -> api.GetTariffsResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_protobuf_billing_proto_init() }
func file_protobuf_billing_proto_init() {
	if File_protobuf_billing_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protobuf_billing_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protobuf_billing_proto_goTypes,
		DependencyIndexes: file_protobuf_billing_proto_depIdxs,
		MessageInfos:      file_protobuf_billing_proto_msgTypes,
	}.Build()
	File_protobuf_billing_proto = out.File
	file_protobuf_billing_proto_rawDesc = nil
	file_protobuf_billing_proto_goTypes = nil
	file_protobuf_billing_proto_depIdxs = nil
}
