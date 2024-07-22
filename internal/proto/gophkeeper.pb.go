// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.6.1
// source: gophkeeper.proto

package proto

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type CheckEMailResponse_EMailStatus int32

const (
	CheckEMailResponse_BUSY      CheckEMailResponse_EMailStatus = 0
	CheckEMailResponse_AVAILABLE CheckEMailResponse_EMailStatus = 1
)

// Enum value maps for CheckEMailResponse_EMailStatus.
var (
	CheckEMailResponse_EMailStatus_name = map[int32]string{
		0: "BUSY",
		1: "AVAILABLE",
	}
	CheckEMailResponse_EMailStatus_value = map[string]int32{
		"BUSY":      0,
		"AVAILABLE": 1,
	}
)

func (x CheckEMailResponse_EMailStatus) Enum() *CheckEMailResponse_EMailStatus {
	p := new(CheckEMailResponse_EMailStatus)
	*p = x
	return p
}

func (x CheckEMailResponse_EMailStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CheckEMailResponse_EMailStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_gophkeeper_proto_enumTypes[0].Descriptor()
}

func (CheckEMailResponse_EMailStatus) Type() protoreflect.EnumType {
	return &file_gophkeeper_proto_enumTypes[0]
}

func (x CheckEMailResponse_EMailStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CheckEMailResponse_EMailStatus.Descriptor instead.
func (CheckEMailResponse_EMailStatus) EnumDescriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{7, 0}
}

type RegistrationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email    string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *RegistrationRequest) Reset() {
	*x = RegistrationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistrationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistrationRequest) ProtoMessage() {}

func (x *RegistrationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistrationRequest.ProtoReflect.Descriptor instead.
func (*RegistrationRequest) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{0}
}

func (x *RegistrationRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *RegistrationRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type RegistrationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId string `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
}

func (x *RegistrationResponse) Reset() {
	*x = RegistrationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistrationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistrationResponse) ProtoMessage() {}

func (x *RegistrationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistrationResponse.ProtoReflect.Descriptor instead.
func (*RegistrationResponse) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{1}
}

func (x *RegistrationResponse) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

type PassOTPRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId string `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	Password  string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *PassOTPRequest) Reset() {
	*x = PassOTPRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PassOTPRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PassOTPRequest) ProtoMessage() {}

func (x *PassOTPRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PassOTPRequest.ProtoReflect.Descriptor instead.
func (*PassOTPRequest) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{2}
}

func (x *PassOTPRequest) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *PassOTPRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type PassOTPResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId string `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
}

func (x *PassOTPResponse) Reset() {
	*x = PassOTPResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PassOTPResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PassOTPResponse) ProtoMessage() {}

func (x *PassOTPResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PassOTPResponse.ProtoReflect.Descriptor instead.
func (*PassOTPResponse) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{3}
}

func (x *PassOTPResponse) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

type MasterKeyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId          string `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	EncryptedMasterKey string `protobuf:"bytes,2,opt,name=encryptedMasterKey,proto3" json:"encryptedMasterKey,omitempty"`
	MasterKeyPassHint  string `protobuf:"bytes,3,opt,name=masterKeyPassHint,proto3" json:"masterKeyPassHint,omitempty"`
	HelloEncrypted     string `protobuf:"bytes,4,opt,name=helloEncrypted,proto3" json:"helloEncrypted,omitempty"`
}

func (x *MasterKeyRequest) Reset() {
	*x = MasterKeyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MasterKeyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MasterKeyRequest) ProtoMessage() {}

func (x *MasterKeyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MasterKeyRequest.ProtoReflect.Descriptor instead.
func (*MasterKeyRequest) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{4}
}

func (x *MasterKeyRequest) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *MasterKeyRequest) GetEncryptedMasterKey() string {
	if x != nil {
		return x.EncryptedMasterKey
	}
	return ""
}

func (x *MasterKeyRequest) GetMasterKeyPassHint() string {
	if x != nil {
		return x.MasterKeyPassHint
	}
	return ""
}

func (x *MasterKeyRequest) GetHelloEncrypted() string {
	if x != nil {
		return x.HelloEncrypted
	}
	return ""
}

type MasterKeyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MasterKeyResponse) Reset() {
	*x = MasterKeyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MasterKeyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MasterKeyResponse) ProtoMessage() {}

func (x *MasterKeyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MasterKeyResponse.ProtoReflect.Descriptor instead.
func (*MasterKeyResponse) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{5}
}

type CheckEMailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *CheckEMailRequest) Reset() {
	*x = CheckEMailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckEMailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckEMailRequest) ProtoMessage() {}

func (x *CheckEMailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckEMailRequest.ProtoReflect.Descriptor instead.
func (*CheckEMailRequest) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{6}
}

func (x *CheckEMailRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type CheckEMailResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status CheckEMailResponse_EMailStatus `protobuf:"varint,1,opt,name=status,proto3,enum=proto.CheckEMailResponse_EMailStatus" json:"status,omitempty"`
}

func (x *CheckEMailResponse) Reset() {
	*x = CheckEMailResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckEMailResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckEMailResponse) ProtoMessage() {}

func (x *CheckEMailResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckEMailResponse.ProtoReflect.Descriptor instead.
func (*CheckEMailResponse) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{7}
}

func (x *CheckEMailResponse) GetStatus() CheckEMailResponse_EMailStatus {
	if x != nil {
		return x.Status
	}
	return CheckEMailResponse_BUSY
}

type LoginRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email    string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{8}
}

func (x *LoginRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type LoginResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId string `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{9}
}

func (x *LoginResponse) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

type AuthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *AuthResponse) Reset() {
	*x = AuthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthResponse) ProtoMessage() {}

func (x *AuthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthResponse.ProtoReflect.Descriptor instead.
func (*AuthResponse) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{10}
}

func (x *AuthResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type PingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{11}
}

type PingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{12}
}

type HelloRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HelloRequest) Reset() {
	*x = HelloRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloRequest) ProtoMessage() {}

func (x *HelloRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloRequest.ProtoReflect.Descriptor instead.
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{13}
}

type HelloResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HelloEncrypted string `protobuf:"bytes,1,opt,name=helloEncrypted,proto3" json:"helloEncrypted,omitempty"`
}

func (x *HelloResponse) Reset() {
	*x = HelloResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gophkeeper_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloResponse) ProtoMessage() {}

func (x *HelloResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloResponse.ProtoReflect.Descriptor instead.
func (*HelloResponse) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{14}
}

func (x *HelloResponse) GetHelloEncrypted() string {
	if x != nil {
		return x.HelloEncrypted
	}
	return ""
}

var File_gophkeeper_proto protoreflect.FileDescriptor

var file_gophkeeper_proto_rawDesc = []byte{
	0x0a, 0x10, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x50, 0x0a, 0x13, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x05, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x60,
	0x01, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x22, 0x35, 0x0a, 0x14, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x4b, 0x0a, 0x0e, 0x50,
	0x61, 0x73, 0x73, 0x4f, 0x54, 0x50, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x30, 0x0a, 0x0f, 0x50, 0x61, 0x73, 0x73,
	0x4f, 0x54, 0x50, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0xb7, 0x01, 0x0a, 0x10, 0x4d,
	0x61, 0x73, 0x74, 0x65, 0x72, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x2e,
	0x0a, 0x12, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x4d, 0x61, 0x73, 0x74, 0x65,
	0x72, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x65, 0x6e, 0x63, 0x72,
	0x79, 0x70, 0x74, 0x65, 0x64, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4b, 0x65, 0x79, 0x12, 0x2c,
	0x0a, 0x11, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x73, 0x73, 0x48,
	0x69, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x6d, 0x61, 0x73, 0x74, 0x65,
	0x72, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x73, 0x73, 0x48, 0x69, 0x6e, 0x74, 0x12, 0x26, 0x0a, 0x0e,
	0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x45, 0x6e, 0x63, 0x72, 0x79,
	0x70, 0x74, 0x65, 0x64, 0x22, 0x13, 0x0a, 0x11, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4b, 0x65,
	0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x32, 0x0a, 0x11, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x45, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d,
	0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x72, 0x02, 0x60, 0x01, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x7b, 0x0a,
	0x12, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x45, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x45, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x45,
	0x4d, 0x61, 0x69, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x22, 0x26, 0x0a, 0x0b, 0x45, 0x4d, 0x61, 0x69, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x08, 0x0a, 0x04, 0x42, 0x55, 0x53, 0x59, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x41,
	0x56, 0x41, 0x49, 0x4c, 0x41, 0x42, 0x4c, 0x45, 0x10, 0x01, 0x22, 0x49, 0x0a, 0x0c, 0x4c, 0x6f,
	0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x05, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02,
	0x60, 0x01, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x2e, 0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x24, 0x0a, 0x0c, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x0d, 0x0a, 0x0b, 0x50,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x0e, 0x0a, 0x0c, 0x50, 0x69,
	0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x0e, 0x0a, 0x0c, 0x48, 0x65,
	0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x37, 0x0a, 0x0d, 0x48, 0x65,
	0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x68,
	0x65, 0x6c, 0x6c, 0x6f, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70,
	0x74, 0x65, 0x64, 0x32, 0x9c, 0x02, 0x0a, 0x13, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x41, 0x0a, 0x0a, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x45, 0x4d, 0x61, 0x69, 0x6c, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x45, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x45, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45,
	0x0a, 0x0a, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x12, 0x1a, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x07, 0x50, 0x61, 0x73, 0x73, 0x4f, 0x54, 0x50,
	0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x4f, 0x54, 0x50,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x50, 0x61, 0x73, 0x73, 0x4f, 0x54, 0x50, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x41, 0x0a, 0x0c, 0x53, 0x65, 0x74, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4b, 0x65, 0x79, 0x12,
	0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4b, 0x65,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x32, 0x78, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x32, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x07, 0x50, 0x61, 0x73, 0x73, 0x4f, 0x54, 0x50,
	0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x4f, 0x54, 0x50,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x39, 0x0a, 0x06,
	0x50, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x12, 0x2f, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x12,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x42, 0x0a, 0x0c, 0x44, 0x61, 0x74, 0x61, 0x41,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x32, 0x0a, 0x05, 0x48, 0x65, 0x6c, 0x6c, 0x6f,
	0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x65,
	0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x10, 0x5a, 0x0e, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gophkeeper_proto_rawDescOnce sync.Once
	file_gophkeeper_proto_rawDescData = file_gophkeeper_proto_rawDesc
)

func file_gophkeeper_proto_rawDescGZIP() []byte {
	file_gophkeeper_proto_rawDescOnce.Do(func() {
		file_gophkeeper_proto_rawDescData = protoimpl.X.CompressGZIP(file_gophkeeper_proto_rawDescData)
	})
	return file_gophkeeper_proto_rawDescData
}

var file_gophkeeper_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_gophkeeper_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_gophkeeper_proto_goTypes = []any{
	(CheckEMailResponse_EMailStatus)(0), // 0: proto.CheckEMailResponse.EMailStatus
	(*RegistrationRequest)(nil),         // 1: proto.RegistrationRequest
	(*RegistrationResponse)(nil),        // 2: proto.RegistrationResponse
	(*PassOTPRequest)(nil),              // 3: proto.PassOTPRequest
	(*PassOTPResponse)(nil),             // 4: proto.PassOTPResponse
	(*MasterKeyRequest)(nil),            // 5: proto.MasterKeyRequest
	(*MasterKeyResponse)(nil),           // 6: proto.MasterKeyResponse
	(*CheckEMailRequest)(nil),           // 7: proto.CheckEMailRequest
	(*CheckEMailResponse)(nil),          // 8: proto.CheckEMailResponse
	(*LoginRequest)(nil),                // 9: proto.LoginRequest
	(*LoginResponse)(nil),               // 10: proto.LoginResponse
	(*AuthResponse)(nil),                // 11: proto.AuthResponse
	(*PingRequest)(nil),                 // 12: proto.PingRequest
	(*PingResponse)(nil),                // 13: proto.PingResponse
	(*HelloRequest)(nil),                // 14: proto.HelloRequest
	(*HelloResponse)(nil),               // 15: proto.HelloResponse
}
var file_gophkeeper_proto_depIdxs = []int32{
	0,  // 0: proto.CheckEMailResponse.status:type_name -> proto.CheckEMailResponse.EMailStatus
	7,  // 1: proto.RegistrationService.CheckEMail:input_type -> proto.CheckEMailRequest
	1,  // 2: proto.RegistrationService.Registrate:input_type -> proto.RegistrationRequest
	3,  // 3: proto.RegistrationService.PassOTP:input_type -> proto.PassOTPRequest
	5,  // 4: proto.RegistrationService.SetMasterKey:input_type -> proto.MasterKeyRequest
	9,  // 5: proto.AuthService.Login:input_type -> proto.LoginRequest
	3,  // 6: proto.AuthService.PassOTP:input_type -> proto.PassOTPRequest
	12, // 7: proto.Pinger.Ping:input_type -> proto.PingRequest
	14, // 8: proto.DataAccessor.Hello:input_type -> proto.HelloRequest
	8,  // 9: proto.RegistrationService.CheckEMail:output_type -> proto.CheckEMailResponse
	2,  // 10: proto.RegistrationService.Registrate:output_type -> proto.RegistrationResponse
	4,  // 11: proto.RegistrationService.PassOTP:output_type -> proto.PassOTPResponse
	6,  // 12: proto.RegistrationService.SetMasterKey:output_type -> proto.MasterKeyResponse
	10, // 13: proto.AuthService.Login:output_type -> proto.LoginResponse
	11, // 14: proto.AuthService.PassOTP:output_type -> proto.AuthResponse
	13, // 15: proto.Pinger.Ping:output_type -> proto.PingResponse
	15, // 16: proto.DataAccessor.Hello:output_type -> proto.HelloResponse
	9,  // [9:17] is the sub-list for method output_type
	1,  // [1:9] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_gophkeeper_proto_init() }
func file_gophkeeper_proto_init() {
	if File_gophkeeper_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gophkeeper_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*RegistrationRequest); i {
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
		file_gophkeeper_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*RegistrationResponse); i {
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
		file_gophkeeper_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*PassOTPRequest); i {
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
		file_gophkeeper_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*PassOTPResponse); i {
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
		file_gophkeeper_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*MasterKeyRequest); i {
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
		file_gophkeeper_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*MasterKeyResponse); i {
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
		file_gophkeeper_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*CheckEMailRequest); i {
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
		file_gophkeeper_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*CheckEMailResponse); i {
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
		file_gophkeeper_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*LoginRequest); i {
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
		file_gophkeeper_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*LoginResponse); i {
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
		file_gophkeeper_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*AuthResponse); i {
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
		file_gophkeeper_proto_msgTypes[11].Exporter = func(v any, i int) any {
			switch v := v.(*PingRequest); i {
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
		file_gophkeeper_proto_msgTypes[12].Exporter = func(v any, i int) any {
			switch v := v.(*PingResponse); i {
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
		file_gophkeeper_proto_msgTypes[13].Exporter = func(v any, i int) any {
			switch v := v.(*HelloRequest); i {
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
		file_gophkeeper_proto_msgTypes[14].Exporter = func(v any, i int) any {
			switch v := v.(*HelloResponse); i {
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
			RawDescriptor: file_gophkeeper_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   4,
		},
		GoTypes:           file_gophkeeper_proto_goTypes,
		DependencyIndexes: file_gophkeeper_proto_depIdxs,
		EnumInfos:         file_gophkeeper_proto_enumTypes,
		MessageInfos:      file_gophkeeper_proto_msgTypes,
	}.Build()
	File_gophkeeper_proto = out.File
	file_gophkeeper_proto_rawDesc = nil
	file_gophkeeper_proto_goTypes = nil
	file_gophkeeper_proto_depIdxs = nil
}
