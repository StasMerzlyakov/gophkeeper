// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.6.1
// source: gophkeeper.proto

package proto

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	RegistrationService_CheckEMail_FullMethodName   = "/proto.RegistrationService/CheckEMail"
	RegistrationService_Registrate_FullMethodName   = "/proto.RegistrationService/Registrate"
	RegistrationService_PassOTP_FullMethodName      = "/proto.RegistrationService/PassOTP"
	RegistrationService_SetMasterKey_FullMethodName = "/proto.RegistrationService/SetMasterKey"
)

// RegistrationServiceClient is the client API for RegistrationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegistrationServiceClient interface {
	CheckEMail(ctx context.Context, in *CheckEMailRequest, opts ...grpc.CallOption) (*CheckEMailResponse, error)
	Registrate(ctx context.Context, in *RegistrationRequest, opts ...grpc.CallOption) (*RegistrationResponse, error)
	PassOTP(ctx context.Context, in *PassOTPRequest, opts ...grpc.CallOption) (*PassOTPResponse, error)
	SetMasterKey(ctx context.Context, in *MasterKeyRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type registrationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRegistrationServiceClient(cc grpc.ClientConnInterface) RegistrationServiceClient {
	return &registrationServiceClient{cc}
}

func (c *registrationServiceClient) CheckEMail(ctx context.Context, in *CheckEMailRequest, opts ...grpc.CallOption) (*CheckEMailResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckEMailResponse)
	err := c.cc.Invoke(ctx, RegistrationService_CheckEMail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registrationServiceClient) Registrate(ctx context.Context, in *RegistrationRequest, opts ...grpc.CallOption) (*RegistrationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegistrationResponse)
	err := c.cc.Invoke(ctx, RegistrationService_Registrate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registrationServiceClient) PassOTP(ctx context.Context, in *PassOTPRequest, opts ...grpc.CallOption) (*PassOTPResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PassOTPResponse)
	err := c.cc.Invoke(ctx, RegistrationService_PassOTP_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registrationServiceClient) SetMasterKey(ctx context.Context, in *MasterKeyRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, RegistrationService_SetMasterKey_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegistrationServiceServer is the server API for RegistrationService service.
// All implementations must embed UnimplementedRegistrationServiceServer
// for forward compatibility
type RegistrationServiceServer interface {
	CheckEMail(context.Context, *CheckEMailRequest) (*CheckEMailResponse, error)
	Registrate(context.Context, *RegistrationRequest) (*RegistrationResponse, error)
	PassOTP(context.Context, *PassOTPRequest) (*PassOTPResponse, error)
	SetMasterKey(context.Context, *MasterKeyRequest) (*empty.Empty, error)
	mustEmbedUnimplementedRegistrationServiceServer()
}

// UnimplementedRegistrationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRegistrationServiceServer struct {
}

func (UnimplementedRegistrationServiceServer) CheckEMail(context.Context, *CheckEMailRequest) (*CheckEMailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckEMail not implemented")
}
func (UnimplementedRegistrationServiceServer) Registrate(context.Context, *RegistrationRequest) (*RegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Registrate not implemented")
}
func (UnimplementedRegistrationServiceServer) PassOTP(context.Context, *PassOTPRequest) (*PassOTPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PassOTP not implemented")
}
func (UnimplementedRegistrationServiceServer) SetMasterKey(context.Context, *MasterKeyRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetMasterKey not implemented")
}
func (UnimplementedRegistrationServiceServer) mustEmbedUnimplementedRegistrationServiceServer() {}

// UnsafeRegistrationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegistrationServiceServer will
// result in compilation errors.
type UnsafeRegistrationServiceServer interface {
	mustEmbedUnimplementedRegistrationServiceServer()
}

func RegisterRegistrationServiceServer(s grpc.ServiceRegistrar, srv RegistrationServiceServer) {
	s.RegisterService(&RegistrationService_ServiceDesc, srv)
}

func _RegistrationService_CheckEMail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckEMailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationServiceServer).CheckEMail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistrationService_CheckEMail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrationServiceServer).CheckEMail(ctx, req.(*CheckEMailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistrationService_Registrate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationServiceServer).Registrate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistrationService_Registrate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrationServiceServer).Registrate(ctx, req.(*RegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistrationService_PassOTP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PassOTPRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationServiceServer).PassOTP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistrationService_PassOTP_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrationServiceServer).PassOTP(ctx, req.(*PassOTPRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistrationService_SetMasterKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MasterKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationServiceServer).SetMasterKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistrationService_SetMasterKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrationServiceServer).SetMasterKey(ctx, req.(*MasterKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegistrationService_ServiceDesc is the grpc.ServiceDesc for RegistrationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegistrationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.RegistrationService",
	HandlerType: (*RegistrationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckEMail",
			Handler:    _RegistrationService_CheckEMail_Handler,
		},
		{
			MethodName: "Registrate",
			Handler:    _RegistrationService_Registrate_Handler,
		},
		{
			MethodName: "PassOTP",
			Handler:    _RegistrationService_PassOTP_Handler,
		},
		{
			MethodName: "SetMasterKey",
			Handler:    _RegistrationService_SetMasterKey_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gophkeeper.proto",
}

const (
	AuthService_Login_FullMethodName   = "/proto.AuthService/Login"
	AuthService_PassOTP_FullMethodName = "/proto.AuthService/PassOTP"
)

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServiceClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	PassOTP(ctx context.Context, in *PassOTPRequest, opts ...grpc.CallOption) (*AuthResponse, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, AuthService_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) PassOTP(ctx context.Context, in *PassOTPRequest, opts ...grpc.CallOption) (*AuthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthResponse)
	err := c.cc.Invoke(ctx, AuthService_PassOTP_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility
type AuthServiceServer interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	PassOTP(context.Context, *PassOTPRequest) (*AuthResponse, error)
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (UnimplementedAuthServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthServiceServer) PassOTP(context.Context, *PassOTPRequest) (*AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PassOTP not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}

// UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServiceServer will
// result in compilation errors.
type UnsafeAuthServiceServer interface {
	mustEmbedUnimplementedAuthServiceServer()
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

func _AuthService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_PassOTP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PassOTPRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).PassOTP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_PassOTP_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).PassOTP(ctx, req.(*PassOTPRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _AuthService_Login_Handler,
		},
		{
			MethodName: "PassOTP",
			Handler:    _AuthService_PassOTP_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gophkeeper.proto",
}

const (
	Pinger_Ping_FullMethodName = "/proto.Pinger/Ping"
)

// PingerClient is the client API for Pinger service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PingerClient interface {
	Ping(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
}

type pingerClient struct {
	cc grpc.ClientConnInterface
}

func NewPingerClient(cc grpc.ClientConnInterface) PingerClient {
	return &pingerClient{cc}
}

func (c *pingerClient) Ping(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Pinger_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PingerServer is the server API for Pinger service.
// All implementations must embed UnimplementedPingerServer
// for forward compatibility
type PingerServer interface {
	Ping(context.Context, *empty.Empty) (*empty.Empty, error)
	mustEmbedUnimplementedPingerServer()
}

// UnimplementedPingerServer must be embedded to have forward compatible implementations.
type UnimplementedPingerServer struct {
}

func (UnimplementedPingerServer) Ping(context.Context, *empty.Empty) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedPingerServer) mustEmbedUnimplementedPingerServer() {}

// UnsafePingerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PingerServer will
// result in compilation errors.
type UnsafePingerServer interface {
	mustEmbedUnimplementedPingerServer()
}

func RegisterPingerServer(s grpc.ServiceRegistrar, srv PingerServer) {
	s.RegisterService(&Pinger_ServiceDesc, srv)
}

func _Pinger_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PingerServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Pinger_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PingerServer).Ping(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Pinger_ServiceDesc is the grpc.ServiceDesc for Pinger service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Pinger_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Pinger",
	HandlerType: (*PingerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Pinger_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gophkeeper.proto",
}

const (
	DataAccessor_Hello_FullMethodName                   = "/proto.DataAccessor/Hello"
	DataAccessor_GetBankCardList_FullMethodName         = "/proto.DataAccessor/GetBankCardList"
	DataAccessor_CreateBankCard_FullMethodName          = "/proto.DataAccessor/CreateBankCard"
	DataAccessor_DeleteBankCard_FullMethodName          = "/proto.DataAccessor/DeleteBankCard"
	DataAccessor_UpdateBankCard_FullMethodName          = "/proto.DataAccessor/UpdateBankCard"
	DataAccessor_GetUserPasswordDataList_FullMethodName = "/proto.DataAccessor/GetUserPasswordDataList"
	DataAccessor_CreateUserPasswordData_FullMethodName  = "/proto.DataAccessor/CreateUserPasswordData"
	DataAccessor_DeleteUserPasswordData_FullMethodName  = "/proto.DataAccessor/DeleteUserPasswordData"
	DataAccessor_UpdateUserPasswordData_FullMethodName  = "/proto.DataAccessor/UpdateUserPasswordData"
)

// DataAccessorClient is the client API for DataAccessor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataAccessorClient interface {
	Hello(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*HelloResponse, error)
	GetBankCardList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*BankCardListResponse, error)
	CreateBankCard(ctx context.Context, in *CreateBankCardRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	DeleteBankCard(ctx context.Context, in *DeleteBankCardRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	UpdateBankCard(ctx context.Context, in *UpdateBankCardRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	GetUserPasswordDataList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*UserPasswordDataResponse, error)
	CreateUserPasswordData(ctx context.Context, in *CreateUserPasswordDataRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	DeleteUserPasswordData(ctx context.Context, in *DeleteUserPasswordDataRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	UpdateUserPasswordData(ctx context.Context, in *UpdateUserPasswordDataRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type dataAccessorClient struct {
	cc grpc.ClientConnInterface
}

func NewDataAccessorClient(cc grpc.ClientConnInterface) DataAccessorClient {
	return &dataAccessorClient{cc}
}

func (c *dataAccessorClient) Hello(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*HelloResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, DataAccessor_Hello_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataAccessorClient) GetBankCardList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*BankCardListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BankCardListResponse)
	err := c.cc.Invoke(ctx, DataAccessor_GetBankCardList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataAccessorClient) CreateBankCard(ctx context.Context, in *CreateBankCardRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, DataAccessor_CreateBankCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataAccessorClient) DeleteBankCard(ctx context.Context, in *DeleteBankCardRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, DataAccessor_DeleteBankCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataAccessorClient) UpdateBankCard(ctx context.Context, in *UpdateBankCardRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, DataAccessor_UpdateBankCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataAccessorClient) GetUserPasswordDataList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*UserPasswordDataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserPasswordDataResponse)
	err := c.cc.Invoke(ctx, DataAccessor_GetUserPasswordDataList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataAccessorClient) CreateUserPasswordData(ctx context.Context, in *CreateUserPasswordDataRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, DataAccessor_CreateUserPasswordData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataAccessorClient) DeleteUserPasswordData(ctx context.Context, in *DeleteUserPasswordDataRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, DataAccessor_DeleteUserPasswordData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataAccessorClient) UpdateUserPasswordData(ctx context.Context, in *UpdateUserPasswordDataRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, DataAccessor_UpdateUserPasswordData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataAccessorServer is the server API for DataAccessor service.
// All implementations must embed UnimplementedDataAccessorServer
// for forward compatibility
type DataAccessorServer interface {
	Hello(context.Context, *empty.Empty) (*HelloResponse, error)
	GetBankCardList(context.Context, *empty.Empty) (*BankCardListResponse, error)
	CreateBankCard(context.Context, *CreateBankCardRequest) (*empty.Empty, error)
	DeleteBankCard(context.Context, *DeleteBankCardRequest) (*empty.Empty, error)
	UpdateBankCard(context.Context, *UpdateBankCardRequest) (*empty.Empty, error)
	GetUserPasswordDataList(context.Context, *empty.Empty) (*UserPasswordDataResponse, error)
	CreateUserPasswordData(context.Context, *CreateUserPasswordDataRequest) (*empty.Empty, error)
	DeleteUserPasswordData(context.Context, *DeleteUserPasswordDataRequest) (*empty.Empty, error)
	UpdateUserPasswordData(context.Context, *UpdateUserPasswordDataRequest) (*empty.Empty, error)
	mustEmbedUnimplementedDataAccessorServer()
}

// UnimplementedDataAccessorServer must be embedded to have forward compatible implementations.
type UnimplementedDataAccessorServer struct {
}

func (UnimplementedDataAccessorServer) Hello(context.Context, *empty.Empty) (*HelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}
func (UnimplementedDataAccessorServer) GetBankCardList(context.Context, *empty.Empty) (*BankCardListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBankCardList not implemented")
}
func (UnimplementedDataAccessorServer) CreateBankCard(context.Context, *CreateBankCardRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBankCard not implemented")
}
func (UnimplementedDataAccessorServer) DeleteBankCard(context.Context, *DeleteBankCardRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBankCard not implemented")
}
func (UnimplementedDataAccessorServer) UpdateBankCard(context.Context, *UpdateBankCardRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBankCard not implemented")
}
func (UnimplementedDataAccessorServer) GetUserPasswordDataList(context.Context, *empty.Empty) (*UserPasswordDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserPasswordDataList not implemented")
}
func (UnimplementedDataAccessorServer) CreateUserPasswordData(context.Context, *CreateUserPasswordDataRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUserPasswordData not implemented")
}
func (UnimplementedDataAccessorServer) DeleteUserPasswordData(context.Context, *DeleteUserPasswordDataRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserPasswordData not implemented")
}
func (UnimplementedDataAccessorServer) UpdateUserPasswordData(context.Context, *UpdateUserPasswordDataRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserPasswordData not implemented")
}
func (UnimplementedDataAccessorServer) mustEmbedUnimplementedDataAccessorServer() {}

// UnsafeDataAccessorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataAccessorServer will
// result in compilation errors.
type UnsafeDataAccessorServer interface {
	mustEmbedUnimplementedDataAccessorServer()
}

func RegisterDataAccessorServer(s grpc.ServiceRegistrar, srv DataAccessorServer) {
	s.RegisterService(&DataAccessor_ServiceDesc, srv)
}

func _DataAccessor_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataAccessorServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataAccessor_Hello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataAccessorServer).Hello(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataAccessor_GetBankCardList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataAccessorServer).GetBankCardList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataAccessor_GetBankCardList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataAccessorServer).GetBankCardList(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataAccessor_CreateBankCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBankCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataAccessorServer).CreateBankCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataAccessor_CreateBankCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataAccessorServer).CreateBankCard(ctx, req.(*CreateBankCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataAccessor_DeleteBankCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteBankCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataAccessorServer).DeleteBankCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataAccessor_DeleteBankCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataAccessorServer).DeleteBankCard(ctx, req.(*DeleteBankCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataAccessor_UpdateBankCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateBankCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataAccessorServer).UpdateBankCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataAccessor_UpdateBankCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataAccessorServer).UpdateBankCard(ctx, req.(*UpdateBankCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataAccessor_GetUserPasswordDataList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataAccessorServer).GetUserPasswordDataList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataAccessor_GetUserPasswordDataList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataAccessorServer).GetUserPasswordDataList(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataAccessor_CreateUserPasswordData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserPasswordDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataAccessorServer).CreateUserPasswordData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataAccessor_CreateUserPasswordData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataAccessorServer).CreateUserPasswordData(ctx, req.(*CreateUserPasswordDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataAccessor_DeleteUserPasswordData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserPasswordDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataAccessorServer).DeleteUserPasswordData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataAccessor_DeleteUserPasswordData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataAccessorServer).DeleteUserPasswordData(ctx, req.(*DeleteUserPasswordDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataAccessor_UpdateUserPasswordData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserPasswordDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataAccessorServer).UpdateUserPasswordData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataAccessor_UpdateUserPasswordData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataAccessorServer).UpdateUserPasswordData(ctx, req.(*UpdateUserPasswordDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DataAccessor_ServiceDesc is the grpc.ServiceDesc for DataAccessor service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataAccessor_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.DataAccessor",
	HandlerType: (*DataAccessorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hello",
			Handler:    _DataAccessor_Hello_Handler,
		},
		{
			MethodName: "GetBankCardList",
			Handler:    _DataAccessor_GetBankCardList_Handler,
		},
		{
			MethodName: "CreateBankCard",
			Handler:    _DataAccessor_CreateBankCard_Handler,
		},
		{
			MethodName: "DeleteBankCard",
			Handler:    _DataAccessor_DeleteBankCard_Handler,
		},
		{
			MethodName: "UpdateBankCard",
			Handler:    _DataAccessor_UpdateBankCard_Handler,
		},
		{
			MethodName: "GetUserPasswordDataList",
			Handler:    _DataAccessor_GetUserPasswordDataList_Handler,
		},
		{
			MethodName: "CreateUserPasswordData",
			Handler:    _DataAccessor_CreateUserPasswordData_Handler,
		},
		{
			MethodName: "DeleteUserPasswordData",
			Handler:    _DataAccessor_DeleteUserPasswordData_Handler,
		},
		{
			MethodName: "UpdateUserPasswordData",
			Handler:    _DataAccessor_UpdateUserPasswordData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gophkeeper.proto",
}
