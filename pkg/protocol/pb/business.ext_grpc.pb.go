// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.0
// source: business.ext.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	BusinessExt_SignIn_FullMethodName                 = "/pb.BusinessExt/SignIn"
	BusinessExt_GetUser_FullMethodName                = "/pb.BusinessExt/GetUser"
	BusinessExt_UpdateUser_FullMethodName             = "/pb.BusinessExt/UpdateUser"
	BusinessExt_SearchUser_FullMethodName             = "/pb.BusinessExt/SearchUser"
	BusinessExt_GetTwitterAuthorizeURL_FullMethodName = "/pb.BusinessExt/GetTwitterAuthorizeURL"
	BusinessExt_TwitterSignIn_FullMethodName          = "/pb.BusinessExt/TwitterSignIn"
	BusinessExt_DailySignIn_FullMethodName            = "/pb.BusinessExt/DailySignIn"
	BusinessExt_FollowTwitter_FullMethodName          = "/pb.BusinessExt/FollowTwitter"
	BusinessExt_GetTaskStatus_FullMethodName          = "/pb.BusinessExt/GetTaskStatus"
	BusinessExt_ClaimTaskReward_FullMethodName        = "/pb.BusinessExt/ClaimTaskReward"
	BusinessExt_ModifyTaskStatus_FullMethodName       = "/pb.BusinessExt/ModifyTaskStatus"
)

// BusinessExtClient is the client API for BusinessExt service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BusinessExtClient interface {
	// 登录
	SignIn(ctx context.Context, in *SignInReq, opts ...grpc.CallOption) (*SignInResp, error)
	// 获取用户信息
	GetUser(ctx context.Context, in *GetUserReq, opts ...grpc.CallOption) (*GetUserResp, error)
	// 更新用户信息
	UpdateUser(ctx context.Context, in *UpdateUserReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 搜索用户(这里简单数据库实现，生产环境建议使用ES)
	SearchUser(ctx context.Context, in *SearchUserReq, opts ...grpc.CallOption) (*SearchUserResp, error)
	// 获取推特授权 URL
	GetTwitterAuthorizeURL(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TwitterAuthorizeURLResp, error)
	// 推特登录
	TwitterSignIn(ctx context.Context, in *TwitterSignInReq, opts ...grpc.CallOption) (*TwitterSignInResp, error)
	// 每日签到接口
	DailySignIn(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*DailySignInResp, error)
	// 关注推特
	FollowTwitter(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TwitterFollowResp, error)
	// 查询任务状态
	GetTaskStatus(ctx context.Context, in *GetTaskStatusReq, opts ...grpc.CallOption) (*GetTaskStatusResp, error)
	// 领取任务奖励
	ClaimTaskReward(ctx context.Context, in *ClaimTaskRewardReq, opts ...grpc.CallOption) (*ClaimTaskRewardResp, error)
	ModifyTaskStatus(ctx context.Context, in *ModifyTaskStatusReq, opts ...grpc.CallOption) (*ModifyTaskStatusResp, error)
}

type businessExtClient struct {
	cc grpc.ClientConnInterface
}

func NewBusinessExtClient(cc grpc.ClientConnInterface) BusinessExtClient {
	return &businessExtClient{cc}
}

func (c *businessExtClient) SignIn(ctx context.Context, in *SignInReq, opts ...grpc.CallOption) (*SignInResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SignInResp)
	err := c.cc.Invoke(ctx, BusinessExt_SignIn_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *businessExtClient) GetUser(ctx context.Context, in *GetUserReq, opts ...grpc.CallOption) (*GetUserResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserResp)
	err := c.cc.Invoke(ctx, BusinessExt_GetUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *businessExtClient) UpdateUser(ctx context.Context, in *UpdateUserReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, BusinessExt_UpdateUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *businessExtClient) SearchUser(ctx context.Context, in *SearchUserReq, opts ...grpc.CallOption) (*SearchUserResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchUserResp)
	err := c.cc.Invoke(ctx, BusinessExt_SearchUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *businessExtClient) GetTwitterAuthorizeURL(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TwitterAuthorizeURLResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TwitterAuthorizeURLResp)
	err := c.cc.Invoke(ctx, BusinessExt_GetTwitterAuthorizeURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *businessExtClient) TwitterSignIn(ctx context.Context, in *TwitterSignInReq, opts ...grpc.CallOption) (*TwitterSignInResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TwitterSignInResp)
	err := c.cc.Invoke(ctx, BusinessExt_TwitterSignIn_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *businessExtClient) DailySignIn(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*DailySignInResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DailySignInResp)
	err := c.cc.Invoke(ctx, BusinessExt_DailySignIn_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *businessExtClient) FollowTwitter(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TwitterFollowResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TwitterFollowResp)
	err := c.cc.Invoke(ctx, BusinessExt_FollowTwitter_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *businessExtClient) GetTaskStatus(ctx context.Context, in *GetTaskStatusReq, opts ...grpc.CallOption) (*GetTaskStatusResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTaskStatusResp)
	err := c.cc.Invoke(ctx, BusinessExt_GetTaskStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *businessExtClient) ClaimTaskReward(ctx context.Context, in *ClaimTaskRewardReq, opts ...grpc.CallOption) (*ClaimTaskRewardResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClaimTaskRewardResp)
	err := c.cc.Invoke(ctx, BusinessExt_ClaimTaskReward_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *businessExtClient) ModifyTaskStatus(ctx context.Context, in *ModifyTaskStatusReq, opts ...grpc.CallOption) (*ModifyTaskStatusResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ModifyTaskStatusResp)
	err := c.cc.Invoke(ctx, BusinessExt_ModifyTaskStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BusinessExtServer is the server API for BusinessExt service.
// All implementations must embed UnimplementedBusinessExtServer
// for forward compatibility.
type BusinessExtServer interface {
	// 登录
	SignIn(context.Context, *SignInReq) (*SignInResp, error)
	// 获取用户信息
	GetUser(context.Context, *GetUserReq) (*GetUserResp, error)
	// 更新用户信息
	UpdateUser(context.Context, *UpdateUserReq) (*emptypb.Empty, error)
	// 搜索用户(这里简单数据库实现，生产环境建议使用ES)
	SearchUser(context.Context, *SearchUserReq) (*SearchUserResp, error)
	// 获取推特授权 URL
	GetTwitterAuthorizeURL(context.Context, *emptypb.Empty) (*TwitterAuthorizeURLResp, error)
	// 推特登录
	TwitterSignIn(context.Context, *TwitterSignInReq) (*TwitterSignInResp, error)
	// 每日签到接口
	DailySignIn(context.Context, *emptypb.Empty) (*DailySignInResp, error)
	// 关注推特
	FollowTwitter(context.Context, *emptypb.Empty) (*TwitterFollowResp, error)
	// 查询任务状态
	GetTaskStatus(context.Context, *GetTaskStatusReq) (*GetTaskStatusResp, error)
	// 领取任务奖励
	ClaimTaskReward(context.Context, *ClaimTaskRewardReq) (*ClaimTaskRewardResp, error)
	ModifyTaskStatus(context.Context, *ModifyTaskStatusReq) (*ModifyTaskStatusResp, error)
	mustEmbedUnimplementedBusinessExtServer()
}

// UnimplementedBusinessExtServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBusinessExtServer struct{}

func (UnimplementedBusinessExtServer) SignIn(context.Context, *SignInReq) (*SignInResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignIn not implemented")
}
func (UnimplementedBusinessExtServer) GetUser(context.Context, *GetUserReq) (*GetUserResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedBusinessExtServer) UpdateUser(context.Context, *UpdateUserReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedBusinessExtServer) SearchUser(context.Context, *SearchUserReq) (*SearchUserResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchUser not implemented")
}
func (UnimplementedBusinessExtServer) GetTwitterAuthorizeURL(context.Context, *emptypb.Empty) (*TwitterAuthorizeURLResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTwitterAuthorizeURL not implemented")
}
func (UnimplementedBusinessExtServer) TwitterSignIn(context.Context, *TwitterSignInReq) (*TwitterSignInResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TwitterSignIn not implemented")
}
func (UnimplementedBusinessExtServer) DailySignIn(context.Context, *emptypb.Empty) (*DailySignInResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DailySignIn not implemented")
}
func (UnimplementedBusinessExtServer) FollowTwitter(context.Context, *emptypb.Empty) (*TwitterFollowResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FollowTwitter not implemented")
}
func (UnimplementedBusinessExtServer) GetTaskStatus(context.Context, *GetTaskStatusReq) (*GetTaskStatusResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTaskStatus not implemented")
}
func (UnimplementedBusinessExtServer) ClaimTaskReward(context.Context, *ClaimTaskRewardReq) (*ClaimTaskRewardResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClaimTaskReward not implemented")
}
func (UnimplementedBusinessExtServer) ModifyTaskStatus(context.Context, *ModifyTaskStatusReq) (*ModifyTaskStatusResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyTaskStatus not implemented")
}
func (UnimplementedBusinessExtServer) mustEmbedUnimplementedBusinessExtServer() {}
func (UnimplementedBusinessExtServer) testEmbeddedByValue()                     {}

// UnsafeBusinessExtServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BusinessExtServer will
// result in compilation errors.
type UnsafeBusinessExtServer interface {
	mustEmbedUnimplementedBusinessExtServer()
}

func RegisterBusinessExtServer(s grpc.ServiceRegistrar, srv BusinessExtServer) {
	// If the following call pancis, it indicates UnimplementedBusinessExtServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BusinessExt_ServiceDesc, srv)
}

func _BusinessExt_SignIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignInReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessExtServer).SignIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BusinessExt_SignIn_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessExtServer).SignIn(ctx, req.(*SignInReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BusinessExt_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessExtServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BusinessExt_GetUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessExtServer).GetUser(ctx, req.(*GetUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BusinessExt_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessExtServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BusinessExt_UpdateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessExtServer).UpdateUser(ctx, req.(*UpdateUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BusinessExt_SearchUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessExtServer).SearchUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BusinessExt_SearchUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessExtServer).SearchUser(ctx, req.(*SearchUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BusinessExt_GetTwitterAuthorizeURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessExtServer).GetTwitterAuthorizeURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BusinessExt_GetTwitterAuthorizeURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessExtServer).GetTwitterAuthorizeURL(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BusinessExt_TwitterSignIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TwitterSignInReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessExtServer).TwitterSignIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BusinessExt_TwitterSignIn_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessExtServer).TwitterSignIn(ctx, req.(*TwitterSignInReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BusinessExt_DailySignIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessExtServer).DailySignIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BusinessExt_DailySignIn_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessExtServer).DailySignIn(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BusinessExt_FollowTwitter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessExtServer).FollowTwitter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BusinessExt_FollowTwitter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessExtServer).FollowTwitter(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BusinessExt_GetTaskStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTaskStatusReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessExtServer).GetTaskStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BusinessExt_GetTaskStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessExtServer).GetTaskStatus(ctx, req.(*GetTaskStatusReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BusinessExt_ClaimTaskReward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClaimTaskRewardReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessExtServer).ClaimTaskReward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BusinessExt_ClaimTaskReward_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessExtServer).ClaimTaskReward(ctx, req.(*ClaimTaskRewardReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BusinessExt_ModifyTaskStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyTaskStatusReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusinessExtServer).ModifyTaskStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BusinessExt_ModifyTaskStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusinessExtServer).ModifyTaskStatus(ctx, req.(*ModifyTaskStatusReq))
	}
	return interceptor(ctx, in, info, handler)
}

// BusinessExt_ServiceDesc is the grpc.ServiceDesc for BusinessExt service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BusinessExt_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.BusinessExt",
	HandlerType: (*BusinessExtServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignIn",
			Handler:    _BusinessExt_SignIn_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _BusinessExt_GetUser_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _BusinessExt_UpdateUser_Handler,
		},
		{
			MethodName: "SearchUser",
			Handler:    _BusinessExt_SearchUser_Handler,
		},
		{
			MethodName: "GetTwitterAuthorizeURL",
			Handler:    _BusinessExt_GetTwitterAuthorizeURL_Handler,
		},
		{
			MethodName: "TwitterSignIn",
			Handler:    _BusinessExt_TwitterSignIn_Handler,
		},
		{
			MethodName: "DailySignIn",
			Handler:    _BusinessExt_DailySignIn_Handler,
		},
		{
			MethodName: "FollowTwitter",
			Handler:    _BusinessExt_FollowTwitter_Handler,
		},
		{
			MethodName: "GetTaskStatus",
			Handler:    _BusinessExt_GetTaskStatus_Handler,
		},
		{
			MethodName: "ClaimTaskReward",
			Handler:    _BusinessExt_ClaimTaskReward_Handler,
		},
		{
			MethodName: "ModifyTaskStatus",
			Handler:    _BusinessExt_ModifyTaskStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "business.ext.proto",
}
