package room

import (
	"context"
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
