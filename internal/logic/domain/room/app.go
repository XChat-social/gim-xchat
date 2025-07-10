package room

import (
	"context"
	"fmt"
	"gim/pkg/protocol/pb"
)

type app struct{}

var App = new(app)

// Push 推送房间消息
func (s *app) Push(ctx context.Context, req *pb.PushRoomReq) error {
	return Service.Push(ctx, req)
}

// SubscribeRoom 订阅房间
func (s *app) SubscribeRoom(ctx context.Context, req *pb.SubscribeRoomReq) error {
	return Service.SubscribeRoom(ctx, req)
}

// CreateChatRoom 创建聊天室
func (s *app) CreateChatRoom(ctx context.Context, req *pb.CreateChatRoomReq) (*pb.CreateChatRoomResp, error) {
	return Service.CreateChatRoom(ctx, req)
}

// GetChatRoom 获取聊天室信息
func (s *app) GetChatRoom(ctx context.Context, req *pb.GetChatRoomReq) (*pb.GetChatRoomResp, error) {
	return Service.GetChatRoom(ctx, req)
}

// GetChatRooms 获取聊天室列表
func (s *app) GetChatRooms(ctx context.Context, req *pb.GetChatRoomsReq) (*pb.GetChatRoomsResp, error) {
	return Service.GetChatRooms(ctx, req)
}

// GetUserChatRooms 获取用户加入的聊天室列表
func (s *app) GetUserChatRooms(ctx context.Context, userId int64, req *pb.GetUserChatRoomsReq) (*pb.GetUserChatRoomsResp, error) {
	return Service.GetUserChatRooms(ctx, userId, req)
}

// JoinChatRoom 加入聊天室
func (s *app) JoinChatRoom(ctx context.Context, req *pb.JoinChatRoomReq) error {
	return Service.JoinChatRoom(ctx, req)
}

// LeaveChatRoom 离开聊天室
func (s *app) LeaveChatRoom(ctx context.Context, req *pb.LeaveChatRoomReq) error {
	return Service.LeaveChatRoom(ctx, req)
}

// GetChatRoomMembers 获取聊天室成员列表
func (s *app) GetChatRoomMembers(ctx context.Context, req *pb.GetChatRoomMembersReq) (*pb.GetChatRoomMembersResp, error) {
	return Service.GetChatRoomMembers(ctx, req)
}

// SendChatRoomMessage 发送聊天室消息
func (s *app) SendChatRoomMessage(ctx context.Context, req *pb.SendChatRoomMessageReq) (*pb.SendChatRoomMessageResp, error) {
	return Service.SendChatRoomMessage(ctx, req)
}

// GetChatRoomMessages 获取聊天室消息历史
func (a *app) GetChatRoomMessages(ctx context.Context, req *pb.GetChatRoomMessagesReq) (*pb.GetChatRoomMessagesResp, error) {
	return Service.GetChatRoomMessages(ctx, req)
}

// GetUserCreatedChatRooms 获取用户创建的聊天室列表
func (s *app) GetUserCreatedChatRooms(ctx context.Context, userId int64) ([]*pb.ChatRoom, error) {
	return Service.GetUserCreatedChatRooms(ctx, userId)
}

// CheckPermissionsByUserId 根据用户ID检查权限
func (s *app) CheckPermissionsByUserId(ctx context.Context, req *pb.CheckPermissionsByUserIdReq) (*pb.CheckPermissionsByUserIdResp, error) {
	fmt.Printf("进入prc1")
	return Service.CheckPermissionsByUserId(ctx, req)
}
