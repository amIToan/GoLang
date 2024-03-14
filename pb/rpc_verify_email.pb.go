// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.3
// source: rpc_verify_email.proto

package pb

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

type VerifyEmailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	SecretCode string `protobuf:"bytes,2,opt,name=secret_code,json=secretCode,proto3" json:"secret_code,omitempty"`
}

func (x *VerifyEmailRequest) Reset() {
	*x = VerifyEmailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_verify_email_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyEmailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyEmailRequest) ProtoMessage() {}

func (x *VerifyEmailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_verify_email_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyEmailRequest.ProtoReflect.Descriptor instead.
func (*VerifyEmailRequest) Descriptor() ([]byte, []int) {
	return file_rpc_verify_email_proto_rawDescGZIP(), []int{0}
}

func (x *VerifyEmailRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *VerifyEmailRequest) GetSecretCode() string {
	if x != nil {
		return x.SecretCode
	}
	return ""
}

type VerifyEmailResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Mail *Mail `protobuf:"bytes,2,opt,name=mail,proto3" json:"mail,omitempty"`
}

func (x *VerifyEmailResponse) Reset() {
	*x = VerifyEmailResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_verify_email_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyEmailResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyEmailResponse) ProtoMessage() {}

func (x *VerifyEmailResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_verify_email_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyEmailResponse.ProtoReflect.Descriptor instead.
func (*VerifyEmailResponse) Descriptor() ([]byte, []int) {
	return file_rpc_verify_email_proto_rawDescGZIP(), []int{1}
}

func (x *VerifyEmailResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *VerifyEmailResponse) GetMail() *Mail {
	if x != nil {
		return x.Mail
	}
	return nil
}

var File_rpc_verify_email_proto protoreflect.FileDescriptor

var file_rpc_verify_email_proto_rawDesc = []byte{
	0x0a, 0x16, 0x72, 0x70, 0x63, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x5f, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x0a, 0x75, 0x73,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x45, 0x0a, 0x12, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x73,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x51, 0x0a, 0x13,
	0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x08, 0x2e, 0x70, 0x62, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x12, 0x1c, 0x0a, 0x04, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x08, 0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x04, 0x6d, 0x61, 0x69, 0x6c, 0x42,
	0x26, 0x5a, 0x24, 0x73, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74,
	0x65, 0x63, 0x68, 0x73, 0x63, 0x68, 0x6f, 0x6f, 0x6c, 0x2f, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65,
	0x62, 0x61, 0x6e, 0x6b, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_verify_email_proto_rawDescOnce sync.Once
	file_rpc_verify_email_proto_rawDescData = file_rpc_verify_email_proto_rawDesc
)

func file_rpc_verify_email_proto_rawDescGZIP() []byte {
	file_rpc_verify_email_proto_rawDescOnce.Do(func() {
		file_rpc_verify_email_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_verify_email_proto_rawDescData)
	})
	return file_rpc_verify_email_proto_rawDescData
}

var file_rpc_verify_email_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_verify_email_proto_goTypes = []interface{}{
	(*VerifyEmailRequest)(nil),  // 0: pb.VerifyEmailRequest
	(*VerifyEmailResponse)(nil), // 1: pb.VerifyEmailResponse
	(*User)(nil),                // 2: pb.User
	(*Mail)(nil),                // 3: pb.Mail
}
var file_rpc_verify_email_proto_depIdxs = []int32{
	2, // 0: pb.VerifyEmailResponse.user:type_name -> pb.User
	3, // 1: pb.VerifyEmailResponse.mail:type_name -> pb.Mail
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_rpc_verify_email_proto_init() }
func file_rpc_verify_email_proto_init() {
	if File_rpc_verify_email_proto != nil {
		return
	}
	file_user_proto_init()
	file_email_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_rpc_verify_email_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyEmailRequest); i {
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
		file_rpc_verify_email_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyEmailResponse); i {
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
			RawDescriptor: file_rpc_verify_email_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_verify_email_proto_goTypes,
		DependencyIndexes: file_rpc_verify_email_proto_depIdxs,
		MessageInfos:      file_rpc_verify_email_proto_msgTypes,
	}.Build()
	File_rpc_verify_email_proto = out.File
	file_rpc_verify_email_proto_rawDesc = nil
	file_rpc_verify_email_proto_goTypes = nil
	file_rpc_verify_email_proto_depIdxs = nil
}