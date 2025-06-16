package room

import (
	"context"
	"gim/pkg/gerrors"
	"gim/pkg/grpclib"
	"gim/pkg/grpclib/picker"
	"gim/pkg/logger"
	"gim/pkg/mq"
	"gim/pkg/protocol/pb"
	"gim/pkg/rpc"
	"gim/pkg/util"
	"time"

	"google.golang.org/protobuf/proto"
)

type service struct{}

var Service = new(service)

func (s *service) Push(ctx context.Context, req *pb.PushRoomReq) error {
	seq, err := SeqRepo.GetNextSeq(req.RoomId)
	if err != nil {
		return err
	}

	msg := &pb.Message{
		Code:     req.Code,
		Content:  req.Content,
		Seq:      seq,
		SendTime: util.UnixMilliTime(time.Now()),
	}
	if req.IsPersist {
		err = s.AddMessage(req.RoomId, msg)
		if err != nil {
			return err
		}
	}

	pushRoomMsg := pb.PushRoomMsg{
		RoomId:  req.RoomId,
		Message: msg,
	}
	bytes, err := proto.Marshal(&pushRoomMsg)
	if err != nil {
		return gerrors.WrapError(err)
	}
	var topicName = mq.PushRoomTopic
	if req.IsPriority {
		topicName = mq.PushRoomPriorityTopic
	}
	err = mq.Publish(topicName, bytes)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddMessage(roomId int64, msg *pb.Message) error {
	err := MessageRepo.Add(roomId, msg)
	if err != nil {
		return err
	}
	return s.DelExpireMessage(roomId)
}

// DelExpireMessage 删除过期消息
func (s *service) DelExpireMessage(roomId int64) error {
	var (
		index int64 = 0
		stop  bool
		min   int64
		max   int64
	)

	for {
		msgs, err := MessageRepo.ListByIndex(roomId, index, index+20)
		if err != nil {
			return err
		}
		if len(msgs) == 0 {
			break
		}

		for _, v := range msgs {
			if v.SendTime > util.UnixMilliTime(time.Now().Add(-MessageExpireTime)) {
				stop = true
				break
			}

			if min == 0 {
				min = v.Seq
			}
			max = v.Seq
		}
		if stop {
			break
		}
	}

	return MessageRepo.DelBySeq(roomId, min, max)
}

// SubscribeRoom 订阅房间
func (s *service) SubscribeRoom(ctx context.Context, req *pb.SubscribeRoomReq) error {
	if req.Seq == 0 {
		return nil
	}

	messages, err := MessageRepo.List(req.RoomId, req.Seq)
	if err != nil {
		return err
	}

	for i := range messages {
		_, err := rpc.GetConnectIntClient().DeliverMessage(picker.ContextWithAddr(ctx, req.ConnAddr), &pb.DeliverMessageReq{
			DeviceId: req.DeviceId,
			Message:  messages[i],
		})
		if err != nil {
			logger.Sugar.Error(err)
		}
	}
	return nil
}

func (s *service) CreateChatRoom(ctx context.Context, req *pb.CreateChatRoomReq) (*pb.CreateChatRoomResp, error) {
	// 从context获取创建者ID
	creatorId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	// 生成聊天室ID
	roomId, err := util.GetNextSeq("chat_room")

	chatRoom := &pb.ChatRoom{
		RoomId:         roomId,
		Name:           req.Name,
		AvatarUrl:      req.AvatarUrl,
		Introduction:   req.Introduction,
		OnlineCount:    0, // 初始在线人数为0
		MemberCount:    1, // 初始成员数为1（创建者）
		MaxMemberCount: req.MaxMemberCount,
		Extra:          req.Extra,
		CreateTime:     util.UnixMilliTime(time.Now()),
		UpdateTime:     util.UnixMilliTime(time.Now()),
	}

	// 保存聊天室信息
	if err := ChatRoomRepo.Add(ctx, chatRoom); err != nil {
		return nil, err
	}

	// 获取创建者信息
	userInfo, err := rpc.GetBusinessIntClient().GetUser(ctx, &pb.GetUserReq{UserId: creatorId})
	if err != nil {
		return nil, err
	}

	// 创建者自动加入聊天室
	member := &pb.ChatRoomMember{
		UserId:    creatorId,
		Nickname:  userInfo.User.Nickname,
		AvatarUrl: userInfo.User.AvatarUrl,
		JoinTime:  util.UnixMilliTime(time.Now()),
		IsOnline:  true,
	}

	if err := ChatRoomMemberRepo.Add(ctx, member); err != nil {
		return nil, err
	}

	return &pb.CreateChatRoomResp{
		RoomId: roomId,
	}, nil
}

// GetChatRoom 获取聊天室信息
func (s *service) GetChatRoom(ctx context.Context, req *pb.GetChatRoomReq) (*pb.GetChatRoomResp, error) {
	// 获取聊天室信息
	room, err := ChatRoomRepo.Get(ctx, req.RoomId)
	if err != nil {
		return nil, err
	}

	return &pb.GetChatRoomResp{
		Room: room,
	}, nil
}

// GetChatRooms 获取聊天室列表
func (s *service) GetChatRooms(ctx context.Context, req *pb.GetChatRoomsReq) (*pb.GetChatRoomsResp, error) {
	// 获取聊天室总数
	total, err := ChatRoomRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// 计算分页参数
	offset := (req.PageNumber - 1) * req.PageSize
	limit := req.PageSize

	// 获取聊天室列表
	rooms, err := ChatRoomRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return &pb.GetChatRoomsResp{
		Rooms: rooms,
		Total: int32(total),
	}, nil
}

// JoinChatRoom 加入聊天室
func (s *service) JoinChatRoom(ctx context.Context, req *pb.JoinChatRoomReq) error {
	// 获取用户信息
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return err
	}

	// 创建成员信息
	member := &pb.ChatRoomMember{
		UserId:   userId,
		JoinTime: util.UnixMilliTime(time.Now()),
		IsOnline: true,
	}

	// 添加成员
	err = ChatRoomMemberRepo.Add(ctx, member)
	if err != nil {
		return err
	}

	return nil
}

// LeaveChatRoom 离开聊天室
func (s *service) LeaveChatRoom(ctx context.Context, req *pb.LeaveChatRoomReq) error {
	// 获取用户信息
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return err
	}

	// 从聊天室移除成员
	err = ChatRoomMemberRepo.Delete(ctx, req.RoomId, userId)
	if err != nil {
		return err
	}

	return nil
}

// GetChatRoomMembers 获取聊天室成员列表
func (s *service) GetChatRoomMembers(ctx context.Context, req *pb.GetChatRoomMembersReq) (*pb.GetChatRoomMembersResp, error) {
	// 获取成员总数
	total, err := ChatRoomMemberRepo.Count(ctx, req.RoomId)
	if err != nil {
		return nil, err
	}

	// 计算分页参数
	offset := (req.PageNumber - 1) * req.PageSize
	limit := req.PageSize

	// 获取成员列表
	members, err := ChatRoomMemberRepo.List(ctx, req.RoomId, offset, limit)
	if err != nil {
		return nil, err
	}

	return &pb.GetChatRoomMembersResp{
		Members: members,
		Total:   int32(total),
	}, nil
}

// SendChatRoomMessage 发送聊天室消息
func (s *service) SendChatRoomMessage(ctx context.Context, req *pb.SendChatRoomMessageReq) (*pb.SendChatRoomMessageResp, error) {
	// 获取用户ID
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	// 验证用户是否在聊天室中
	member, err := ChatRoomMemberRepo.Get(ctx, req.RoomId, userId)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, gerrors.ErrNotInChatRoom
	}

	// 设置发送时间
	sendTime := req.SendTime
	if sendTime == 0 {
		sendTime = util.UnixMilliTime(time.Now())
	}

	// 构造消息内容
	messageData := &pb.ChatRoomMessageData{
		SenderId:   userId,
		SenderName: member.Nickname,
		Content:    req.Content,
	}

	contentBytes, err := proto.Marshal(messageData)
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 获取消息序列号
	seq, err := SeqRepo.GetNextSeq(req.RoomId)
	if err != nil {
		return nil, err
	}

	// 构造推送请求
	pushReq := &pb.PushRoomReq{
		RoomId:     req.RoomId,
		Code:       1, // 普通消息类型
		Content:    contentBytes,
		SendTime:   sendTime,
		IsPersist:  true,  // 持久化消息
		IsPriority: false, // 普通优先级
	}

	// 推送消息
	err = s.Push(ctx, pushReq)
	if err != nil {
		return nil, err
	}

	return &pb.SendChatRoomMessageResp{
		Seq: seq,
	}, nil
}
