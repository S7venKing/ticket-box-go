package identityv1

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// User is the public representation of an account.
// Credentials (password, password hash) are never part of this message.
type User struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Email         string                 `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	FullName      string                 `protobuf:"bytes,3,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"`
	Phone         string                 `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty"`
	IsActive      bool                   `protobuf:"varint,5,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_identity_v1_identity_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

func (x *User) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *User) GetIsActive() bool {
	if x != nil {
		return x.IsActive
	}
	return false
}

func (x *User) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *User) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type CreateUserRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Email string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	// Plaintext password. Only ever hashed by the application layer; never logged.
	Password      string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	FullName      string `protobuf:"bytes,3,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"`
	Phone         string `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateUserRequest) Reset() {
	*x = CreateUserRequest{}
	mi := &file_identity_v1_identity_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserRequest) ProtoMessage() {}

func (x *CreateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUserRequest.ProtoReflect.Descriptor instead.
func (*CreateUserRequest) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{1}
}

func (x *CreateUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *CreateUserRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *CreateUserRequest) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

func (x *CreateUserRequest) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

type CreateUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *User                  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateUserResponse) Reset() {
	*x = CreateUserResponse{}
	mi := &file_identity_v1_identity_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserResponse) ProtoMessage() {}

func (x *CreateUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUserResponse.ProtoReflect.Descriptor instead.
func (*CreateUserResponse) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{2}
}

func (x *CreateUserResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type GetUserByIdRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserByIdRequest) Reset() {
	*x = GetUserByIdRequest{}
	mi := &file_identity_v1_identity_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserByIdRequest) ProtoMessage() {}

func (x *GetUserByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserByIdRequest.ProtoReflect.Descriptor instead.
func (*GetUserByIdRequest) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserByIdRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetUserByIdResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *User                  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserByIdResponse) Reset() {
	*x = GetUserByIdResponse{}
	mi := &file_identity_v1_identity_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserByIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserByIdResponse) ProtoMessage() {}

func (x *GetUserByIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserByIdResponse.ProtoReflect.Descriptor instead.
func (*GetUserByIdResponse) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{4}
}

func (x *GetUserByIdResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type GetUserByEmailRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserByEmailRequest) Reset() {
	*x = GetUserByEmailRequest{}
	mi := &file_identity_v1_identity_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserByEmailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserByEmailRequest) ProtoMessage() {}

func (x *GetUserByEmailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserByEmailRequest.ProtoReflect.Descriptor instead.
func (*GetUserByEmailRequest) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{5}
}

func (x *GetUserByEmailRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type GetUserByEmailResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *User                  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserByEmailResponse) Reset() {
	*x = GetUserByEmailResponse{}
	mi := &file_identity_v1_identity_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserByEmailResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserByEmailResponse) ProtoMessage() {}

func (x *GetUserByEmailResponse) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserByEmailResponse.ProtoReflect.Descriptor instead.
func (*GetUserByEmailResponse) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{6}
}

func (x *GetUserByEmailResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type AuthenticateUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AuthenticateUserRequest) Reset() {
	*x = AuthenticateUserRequest{}
	mi := &file_identity_v1_identity_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthenticateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateUserRequest) ProtoMessage() {}

func (x *AuthenticateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateUserRequest.ProtoReflect.Descriptor instead.
func (*AuthenticateUserRequest) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{7}
}

func (x *AuthenticateUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *AuthenticateUserRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type AuthenticateUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *User                  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	AccessToken   string                 `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AuthenticateUserResponse) Reset() {
	*x = AuthenticateUserResponse{}
	mi := &file_identity_v1_identity_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthenticateUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateUserResponse) ProtoMessage() {}

func (x *AuthenticateUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateUserResponse.ProtoReflect.Descriptor instead.
func (*AuthenticateUserResponse) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{8}
}

func (x *AuthenticateUserResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *AuthenticateUserResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

type UpdateUserProfileRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FullName      string                 `protobuf:"bytes,2,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"`
	Phone         string                 `protobuf:"bytes,3,opt,name=phone,proto3" json:"phone,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateUserProfileRequest) Reset() {
	*x = UpdateUserProfileRequest{}
	mi := &file_identity_v1_identity_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateUserProfileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserProfileRequest) ProtoMessage() {}

func (x *UpdateUserProfileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserProfileRequest.ProtoReflect.Descriptor instead.
func (*UpdateUserProfileRequest) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{9}
}

func (x *UpdateUserProfileRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateUserProfileRequest) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

func (x *UpdateUserProfileRequest) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

type UpdateUserProfileResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *User                  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateUserProfileResponse) Reset() {
	*x = UpdateUserProfileResponse{}
	mi := &file_identity_v1_identity_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateUserProfileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserProfileResponse) ProtoMessage() {}

func (x *UpdateUserProfileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserProfileResponse.ProtoReflect.Descriptor instead.
func (*UpdateUserProfileResponse) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{10}
}

func (x *UpdateUserProfileResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type ActivateUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ActivateUserRequest) Reset() {
	*x = ActivateUserRequest{}
	mi := &file_identity_v1_identity_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ActivateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivateUserRequest) ProtoMessage() {}

func (x *ActivateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivateUserRequest.ProtoReflect.Descriptor instead.
func (*ActivateUserRequest) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{11}
}

func (x *ActivateUserRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ActivateUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *User                  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ActivateUserResponse) Reset() {
	*x = ActivateUserResponse{}
	mi := &file_identity_v1_identity_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ActivateUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivateUserResponse) ProtoMessage() {}

func (x *ActivateUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivateUserResponse.ProtoReflect.Descriptor instead.
func (*ActivateUserResponse) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{12}
}

func (x *ActivateUserResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type DeactivateUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeactivateUserRequest) Reset() {
	*x = DeactivateUserRequest{}
	mi := &file_identity_v1_identity_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeactivateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeactivateUserRequest) ProtoMessage() {}

func (x *DeactivateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeactivateUserRequest.ProtoReflect.Descriptor instead.
func (*DeactivateUserRequest) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{13}
}

func (x *DeactivateUserRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeactivateUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *User                  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeactivateUserResponse) Reset() {
	*x = DeactivateUserResponse{}
	mi := &file_identity_v1_identity_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeactivateUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeactivateUserResponse) ProtoMessage() {}

func (x *DeactivateUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_identity_v1_identity_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeactivateUserResponse.ProtoReflect.Descriptor instead.
func (*DeactivateUserResponse) Descriptor() ([]byte, []int) {
	return file_identity_v1_identity_proto_rawDescGZIP(), []int{14}
}

func (x *DeactivateUserResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

var File_identity_v1_identity_proto protoreflect.FileDescriptor

const file_identity_v1_identity_proto_rawDesc = "" +
	"\n" +
	"\x1aidentity/v1/identity.proto\x12\videntity.v1\x1a\x1fgoogle/protobuf/timestamp.proto\"\xf2\x01\n" +
	"\x04User\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x14\n" +
	"\x05email\x18\x02 \x01(\tR\x05email\x12\x1b\n" +
	"\tfull_name\x18\x03 \x01(\tR\bfullName\x12\x14\n" +
	"\x05phone\x18\x04 \x01(\tR\x05phone\x12\x1b\n" +
	"\tis_active\x18\x05 \x01(\bR\bisActive\x129\n" +
	"\n" +
	"created_at\x18\x06 \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x129\n" +
	"\n" +
	"updated_at\x18\a \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAt\"x\n" +
	"\x11CreateUserRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\x12\x1b\n" +
	"\tfull_name\x18\x03 \x01(\tR\bfullName\x12\x14\n" +
	"\x05phone\x18\x04 \x01(\tR\x05phone\";\n" +
	"\x12CreateUserResponse\x12%\n" +
	"\x04user\x18\x01 \x01(\v2\x11.identity.v1.UserR\x04user\"$\n" +
	"\x12GetUserByIdRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"<\n" +
	"\x13GetUserByIdResponse\x12%\n" +
	"\x04user\x18\x01 \x01(\v2\x11.identity.v1.UserR\x04user\"-\n" +
	"\x15GetUserByEmailRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\"?\n" +
	"\x16GetUserByEmailResponse\x12%\n" +
	"\x04user\x18\x01 \x01(\v2\x11.identity.v1.UserR\x04user\"K\n" +
	"\x17AuthenticateUserRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\"d\n" +
	"\x18AuthenticateUserResponse\x12%\n" +
	"\x04user\x18\x01 \x01(\v2\x11.identity.v1.UserR\x04user\x12!\n" +
	"\faccess_token\x18\x02 \x01(\tR\vaccessToken\"]\n" +
	"\x18UpdateUserProfileRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x1b\n" +
	"\tfull_name\x18\x02 \x01(\tR\bfullName\x12\x14\n" +
	"\x05phone\x18\x03 \x01(\tR\x05phone\"B\n" +
	"\x19UpdateUserProfileResponse\x12%\n" +
	"\x04user\x18\x01 \x01(\v2\x11.identity.v1.UserR\x04user\"%\n" +
	"\x13ActivateUserRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"=\n" +
	"\x14ActivateUserResponse\x12%\n" +
	"\x04user\x18\x01 \x01(\v2\x11.identity.v1.UserR\x04user\"'\n" +
	"\x15DeactivateUserRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"?\n" +
	"\x16DeactivateUserResponse\x12%\n" +
	"\x04user\x18\x01 \x01(\v2\x11.identity.v1.UserR\x04user2\x82\x05\n" +
	"\x0fIdentityService\x12M\n" +
	"\n" +
	"CreateUser\x12\x1e.identity.v1.CreateUserRequest\x1a\x1f.identity.v1.CreateUserResponse\x12P\n" +
	"\vGetUserById\x12\x1f.identity.v1.GetUserByIdRequest\x1a .identity.v1.GetUserByIdResponse\x12Y\n" +
	"\x0eGetUserByEmail\x12\".identity.v1.GetUserByEmailRequest\x1a#.identity.v1.GetUserByEmailResponse\x12_\n" +
	"\x10AuthenticateUser\x12$.identity.v1.AuthenticateUserRequest\x1a%.identity.v1.AuthenticateUserResponse\x12b\n" +
	"\x11UpdateUserProfile\x12%.identity.v1.UpdateUserProfileRequest\x1a&.identity.v1.UpdateUserProfileResponse\x12S\n" +
	"\fActivateUser\x12 .identity.v1.ActivateUserRequest\x1a!.identity.v1.ActivateUserResponse\x12Y\n" +
	"\x0eDeactivateUser\x12\".identity.v1.DeactivateUserRequest\x1a#.identity.v1.DeactivateUserResponseB?Z=github.com/S7venKing/ticket-box-go/gen/identity/v1;identityv1b\x06proto3"

var (
	file_identity_v1_identity_proto_rawDescOnce sync.Once
	file_identity_v1_identity_proto_rawDescData []byte
)

func file_identity_v1_identity_proto_rawDescGZIP() []byte {
	file_identity_v1_identity_proto_rawDescOnce.Do(func() {
		file_identity_v1_identity_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_identity_v1_identity_proto_rawDesc), len(file_identity_v1_identity_proto_rawDesc)))
	})
	return file_identity_v1_identity_proto_rawDescData
}

var file_identity_v1_identity_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_identity_v1_identity_proto_goTypes = []any{
	(*User)(nil),                      // 0: identity.v1.User
	(*CreateUserRequest)(nil),         // 1: identity.v1.CreateUserRequest
	(*CreateUserResponse)(nil),        // 2: identity.v1.CreateUserResponse
	(*GetUserByIdRequest)(nil),        // 3: identity.v1.GetUserByIdRequest
	(*GetUserByIdResponse)(nil),       // 4: identity.v1.GetUserByIdResponse
	(*GetUserByEmailRequest)(nil),     // 5: identity.v1.GetUserByEmailRequest
	(*GetUserByEmailResponse)(nil),    // 6: identity.v1.GetUserByEmailResponse
	(*AuthenticateUserRequest)(nil),   // 7: identity.v1.AuthenticateUserRequest
	(*AuthenticateUserResponse)(nil),  // 8: identity.v1.AuthenticateUserResponse
	(*UpdateUserProfileRequest)(nil),  // 9: identity.v1.UpdateUserProfileRequest
	(*UpdateUserProfileResponse)(nil), // 10: identity.v1.UpdateUserProfileResponse
	(*ActivateUserRequest)(nil),       // 11: identity.v1.ActivateUserRequest
	(*ActivateUserResponse)(nil),      // 12: identity.v1.ActivateUserResponse
	(*DeactivateUserRequest)(nil),     // 13: identity.v1.DeactivateUserRequest
	(*DeactivateUserResponse)(nil),    // 14: identity.v1.DeactivateUserResponse
	(*timestamppb.Timestamp)(nil),     // 15: google.protobuf.Timestamp
}
var file_identity_v1_identity_proto_depIdxs = []int32{
	15, // 0: identity.v1.User.created_at:type_name -> google.protobuf.Timestamp
	15, // 1: identity.v1.User.updated_at:type_name -> google.protobuf.Timestamp
	0,  // 2: identity.v1.CreateUserResponse.user:type_name -> identity.v1.User
	0,  // 3: identity.v1.GetUserByIdResponse.user:type_name -> identity.v1.User
	0,  // 4: identity.v1.GetUserByEmailResponse.user:type_name -> identity.v1.User
	0,  // 5: identity.v1.AuthenticateUserResponse.user:type_name -> identity.v1.User
	0,  // 6: identity.v1.UpdateUserProfileResponse.user:type_name -> identity.v1.User
	0,  // 7: identity.v1.ActivateUserResponse.user:type_name -> identity.v1.User
	0,  // 8: identity.v1.DeactivateUserResponse.user:type_name -> identity.v1.User
	1,  // 9: identity.v1.IdentityService.CreateUser:input_type -> identity.v1.CreateUserRequest
	3,  // 10: identity.v1.IdentityService.GetUserById:input_type -> identity.v1.GetUserByIdRequest
	5,  // 11: identity.v1.IdentityService.GetUserByEmail:input_type -> identity.v1.GetUserByEmailRequest
	7,  // 12: identity.v1.IdentityService.AuthenticateUser:input_type -> identity.v1.AuthenticateUserRequest
	9,  // 13: identity.v1.IdentityService.UpdateUserProfile:input_type -> identity.v1.UpdateUserProfileRequest
	11, // 14: identity.v1.IdentityService.ActivateUser:input_type -> identity.v1.ActivateUserRequest
	13, // 15: identity.v1.IdentityService.DeactivateUser:input_type -> identity.v1.DeactivateUserRequest
	2,  // 16: identity.v1.IdentityService.CreateUser:output_type -> identity.v1.CreateUserResponse
	4,  // 17: identity.v1.IdentityService.GetUserById:output_type -> identity.v1.GetUserByIdResponse
	6,  // 18: identity.v1.IdentityService.GetUserByEmail:output_type -> identity.v1.GetUserByEmailResponse
	8,  // 19: identity.v1.IdentityService.AuthenticateUser:output_type -> identity.v1.AuthenticateUserResponse
	10, // 20: identity.v1.IdentityService.UpdateUserProfile:output_type -> identity.v1.UpdateUserProfileResponse
	12, // 21: identity.v1.IdentityService.ActivateUser:output_type -> identity.v1.ActivateUserResponse
	14, // 22: identity.v1.IdentityService.DeactivateUser:output_type -> identity.v1.DeactivateUserResponse
	16, // [16:23] is the sub-list for method output_type
	9,  // [9:16] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_identity_v1_identity_proto_init() }
func file_identity_v1_identity_proto_init() {
	if File_identity_v1_identity_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_identity_v1_identity_proto_rawDesc), len(file_identity_v1_identity_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_identity_v1_identity_proto_goTypes,
		DependencyIndexes: file_identity_v1_identity_proto_depIdxs,
		MessageInfos:      file_identity_v1_identity_proto_msgTypes,
	}.Build()
	File_identity_v1_identity_proto = out.File
	file_identity_v1_identity_proto_goTypes = nil
	file_identity_v1_identity_proto_depIdxs = nil
}
