// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: internal/proto/authgophkeeper.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AuthGophkeeperClient is the client API for AuthGophkeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthGophkeeperClient interface {
	UserRegister(ctx context.Context, in *UserRegisterRequest, opts ...grpc.CallOption) (*UserRegisterResponse, error)
	UserLogin(ctx context.Context, in *UserLoginRequest, opts ...grpc.CallOption) (*UserLoginResponse, error)
	Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error)
}

type authGophkeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthGophkeeperClient(cc grpc.ClientConnInterface) AuthGophkeeperClient {
	return &authGophkeeperClient{cc}
}

func (c *authGophkeeperClient) UserRegister(ctx context.Context, in *UserRegisterRequest, opts ...grpc.CallOption) (*UserRegisterResponse, error) {
	out := new(UserRegisterResponse)
	err := c.cc.Invoke(ctx, "/authgophkeeper.AuthGophkeeper/UserRegister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authGophkeeperClient) UserLogin(ctx context.Context, in *UserLoginRequest, opts ...grpc.CallOption) (*UserLoginResponse, error) {
	out := new(UserLoginResponse)
	err := c.cc.Invoke(ctx, "/authgophkeeper.AuthGophkeeper/UserLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authGophkeeperClient) Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error) {
	out := new(RefreshResponse)
	err := c.cc.Invoke(ctx, "/authgophkeeper.AuthGophkeeper/Refresh", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthGophkeeperServer is the server API for AuthGophkeeper service.
// All implementations must embed UnimplementedAuthGophkeeperServer
// for forward compatibility
type AuthGophkeeperServer interface {
	UserRegister(context.Context, *UserRegisterRequest) (*UserRegisterResponse, error)
	UserLogin(context.Context, *UserLoginRequest) (*UserLoginResponse, error)
	Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error)
	mustEmbedUnimplementedAuthGophkeeperServer()
}

// UnimplementedAuthGophkeeperServer must be embedded to have forward compatible implementations.
type UnimplementedAuthGophkeeperServer struct {
}

func (UnimplementedAuthGophkeeperServer) UserRegister(context.Context, *UserRegisterRequest) (*UserRegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserRegister not implemented")
}
func (UnimplementedAuthGophkeeperServer) UserLogin(context.Context, *UserLoginRequest) (*UserLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserLogin not implemented")
}
func (UnimplementedAuthGophkeeperServer) Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}
func (UnimplementedAuthGophkeeperServer) mustEmbedUnimplementedAuthGophkeeperServer() {}

// UnsafeAuthGophkeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthGophkeeperServer will
// result in compilation errors.
type UnsafeAuthGophkeeperServer interface {
	mustEmbedUnimplementedAuthGophkeeperServer()
}

func RegisterAuthGophkeeperServer(s grpc.ServiceRegistrar, srv AuthGophkeeperServer) {
	s.RegisterService(&AuthGophkeeper_ServiceDesc, srv)
}

func _AuthGophkeeper_UserRegister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthGophkeeperServer).UserRegister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authgophkeeper.AuthGophkeeper/UserRegister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthGophkeeperServer).UserRegister(ctx, req.(*UserRegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthGophkeeper_UserLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthGophkeeperServer).UserLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authgophkeeper.AuthGophkeeper/UserLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthGophkeeperServer).UserLogin(ctx, req.(*UserLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthGophkeeper_Refresh_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthGophkeeperServer).Refresh(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authgophkeeper.AuthGophkeeper/Refresh",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthGophkeeperServer).Refresh(ctx, req.(*RefreshRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthGophkeeper_ServiceDesc is the grpc.ServiceDesc for AuthGophkeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthGophkeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "authgophkeeper.AuthGophkeeper",
	HandlerType: (*AuthGophkeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserRegister",
			Handler:    _AuthGophkeeper_UserRegister_Handler,
		},
		{
			MethodName: "UserLogin",
			Handler:    _AuthGophkeeper_UserLogin_Handler,
		},
		{
			MethodName: "Refresh",
			Handler:    _AuthGophkeeper_Refresh_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/authgophkeeper.proto",
}
