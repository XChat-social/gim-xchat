package room

import (
	"context"
	"gim/internal/api/models"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"gim/pkg/protocol/pb"
	"log"
	"time"

	"gorm.io/gorm"
)

type chatRoomRepo struct{}

var ChatRoomRepo = new(chatRoomRepo)

type chatRoomMemberRepo struct{}

var ChatRoomMemberRepo = new(chatRoomMemberRepo)

// Add 添加聊天室
func (r *chatRoomRepo) Add(ctx context.Context, chatRoom *pb.ChatRoom) error {
	// 转换为数据库模型
	modelChatRoom := &models.ChatRoom{
		RoomID:         chatRoom.RoomId,
		Name:           chatRoom.Name,
		AvatarURL:      chatRoom.AvatarUrl,
		Introduction:   chatRoom.Introduction,
		CreatorId:      chatRoom.CreatorId,
		OnlineCount:    chatRoom.OnlineCount,
		MemberCount:    chatRoom.MemberCount,
		MaxMemberCount: chatRoom.MaxMemberCount,
		Extra:          chatRoom.Extra,
		CreateTime:     time.UnixMilli(chatRoom.CreateTime), // 从毫秒级时间戳转换为 time.Time
		UpdateTime:     time.UnixMilli(chatRoom.UpdateTime), // 从毫秒级时间戳转换为 time.Time
	}
	return db.DB.Create(modelChatRoom).Error
}

// Get 获取聊天室信息
func (r *chatRoomRepo) Get(ctx context.Context, roomId int64) (*pb.ChatRoom, error) {
	var chatRoom models.ChatRoom
	err := db.DB.Where("room_id = ?", roomId).First(&chatRoom).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 转换为 proto 消息
	return &pb.ChatRoom{
		RoomId:         chatRoom.RoomID,
		Name:           chatRoom.Name,
		AvatarUrl:      chatRoom.AvatarURL,
		Introduction:   chatRoom.Introduction,
		OnlineCount:    chatRoom.OnlineCount,
		MemberCount:    chatRoom.MemberCount,
		MaxMemberCount: chatRoom.MaxMemberCount,
		Extra:          chatRoom.Extra,
		CreateTime:     chatRoom.CreateTime.Unix(), // 转换为 Unix 时间戳
		UpdateTime:     chatRoom.UpdateTime.Unix(), // 转换为 Unix 时间戳
	}, nil
}

// List 获取聊天室列表
func (r *chatRoomRepo) List(ctx context.Context, offset, limit int32) ([]*pb.ChatRoom, error) {
	var modelChatRooms []models.ChatRoom
	err := db.DB.Offset(int(offset)).Limit(int(limit)).Find(&modelChatRooms).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 转换为 proto 消息列表
	chatRooms := make([]*pb.ChatRoom, 0, len(modelChatRooms))
	for _, room := range modelChatRooms {
		chatRooms = append(chatRooms, &pb.ChatRoom{
			RoomId:         room.RoomID,
			Name:           room.Name,
			AvatarUrl:      room.AvatarURL,
			Introduction:   room.Introduction,
			OnlineCount:    room.OnlineCount,
			MemberCount:    room.MemberCount,
			MaxMemberCount: room.MaxMemberCount,
			Extra:          room.Extra,
			CreateTime:     room.CreateTime.Unix(),
			UpdateTime:     room.UpdateTime.Unix(),
			Level:          calculateRoomLevel(room.MemberCount),
		})
	}
	return chatRooms, nil
}

// ChatRoomMessage 仓储接口
type chatRoomMessageRepo struct{}

var ChatRoomMessageRepo = new(chatRoomMessageRepo)

// Add 添加聊天室消息
func (r *chatRoomMessageRepo) Add(ctx context.Context, message *pb.ChatRoomMessage) error {
	// 转换为数据库模型
	modelMessage := &models.ChatRoomMessage{
		RoomID:    uint64(message.RoomId),
		UserID:    uint64(message.UserId),
		RequestID: message.RequestId,
		Code:      int8(message.Code),
		Content:   message.Content,
		Seq:       uint64(message.Seq),
		SendTime:  time.Unix(message.SendTime, 0),
		Status:    int8(message.Status),
	}

	err := db.DB.Create(modelMessage).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// GetNextSeq 获取下一个消息序列号
func (r *chatRoomMessageRepo) GetNextSeq(ctx context.Context, roomId int64) (uint64, error) {
	var result struct {
		MaxSeq uint64 `gorm:"column:max_seq"`
	}
	err := db.DB.Model(&models.ChatRoomMessage{}).Where("room_id = ?", roomId).Select("COALESCE(MAX(seq), 0) as max_seq").Scan(&result).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return result.MaxSeq + 1, nil
}

// List 获取聊天室消息列表
func (r *chatRoomMessageRepo) List(ctx context.Context, roomId int64, offset, limit int32) ([]*pb.ChatRoomMessage, error) {
	var dbMessages []models.ChatRoomMessage
	err := db.DB.Where("room_id = ? AND status = 0", roomId).
		Order("seq DESC").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&dbMessages).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 转换为 proto 消息列表
	messages := make([]*pb.ChatRoomMessage, 0, len(dbMessages))
	for i := len(dbMessages) - 1; i >= 0; i-- {
		msg := dbMessages[i]
		messages = append(messages, &pb.ChatRoomMessage{
			Id:        int64(msg.ID),
			RoomId:    int64(msg.RoomID),
			UserId:    int64(msg.UserID),
			RequestId: msg.RequestID,
			Code:      int32(msg.Code),
			Content:   msg.Content,
			Seq:       int64(msg.Seq),
			SendTime:  msg.SendTime.Unix(),
			Status:    int32(msg.Status),
		})
	}
	return messages, nil
}

// Count 获取聊天室消息总数
func (r *chatRoomMessageRepo) Count(ctx context.Context, roomId int64) (int32, error) {
	var count int32
	err := db.DB.Model(&models.ChatRoomMessage{}).Where("room_id = ? AND status = 0", roomId).Count(&count).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return count, nil
}

// Count 获取聊天室总数
func (r *chatRoomRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := db.DB.Model(&models.ChatRoom{}).Count(&count).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return count, nil
}

// Add 添加聊天室成员
func (r *chatRoomMemberRepo) Add(ctx context.Context, member *pb.ChatRoomMember) error {
	// 转换为数据库模型
	dbMember := &models.ChatRoomMember{
		RoomID:    uint64(member.RoomId),
		UserID:    uint64(member.UserId),
		Nickname:  member.Nickname,
		AvatarURL: member.AvatarUrl,
		JoinTime:  time.UnixMilli(member.JoinTime),
		Status:    int8(member.Status),
	}
	return db.DB.Create(dbMember).Error
}

// Delete 删除聊天室成员
func (r *chatRoomMemberRepo) Delete(ctx context.Context, roomId, userId int64) error {
	return db.DB.Where("room_id = ? AND user_id = ?", roomId, userId).Delete(&models.ChatRoomMember{}).Error
}

// Count 获取聊天室成员总数
func (r *chatRoomMemberRepo) Count(ctx context.Context, roomId int64) (int64, error) {
	var count int64
	err := db.DB.Model(&models.ChatRoomMember{}).Where("room_id = ?", roomId).Count(&count).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return count, nil
}

// List 获取聊天室成员列表
func (r *chatRoomMemberRepo) List(ctx context.Context, roomId int64, offset, limit int32) ([]*pb.ChatRoomMember, error) {
	var dbMembers []models.ChatRoomMember
	err := db.DB.Where("room_id = ?", roomId).Offset(int(offset)).Limit(int(limit)).Find(&dbMembers).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 转换为 proto 消息列表
	members := make([]*pb.ChatRoomMember, 0, len(dbMembers))
	for _, member := range dbMembers {
		members = append(members, &pb.ChatRoomMember{
			RoomId:    int64(member.RoomID),
			UserId:    int64(member.UserID),
			Nickname:  member.Nickname,
			AvatarUrl: member.AvatarURL,
			JoinTime:  member.JoinTime.Unix(),
			Status:    int32(member.Status),
		})
	}
	return members, nil
}

// ListByUserId 获取用户加入的所有聊天室
func (r *chatRoomMemberRepo) ListByUserId(ctx context.Context, userId int64) ([]*pb.ChatRoomMember, error) {
	var dbMembers []models.ChatRoomMember
	err := db.DB.Where("user_id = ?", userId).Find(&dbMembers).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 转换为protobuf格式
	var members []*pb.ChatRoomMember
	for _, dbMember := range dbMembers {
		members = append(members, &pb.ChatRoomMember{
			RoomId:    int64(dbMember.RoomID),
			UserId:    int64(dbMember.UserID),
			Nickname:  dbMember.Nickname,
			AvatarUrl: dbMember.AvatarURL,
			JoinTime:  dbMember.JoinTime.Unix(),
			Status:    int32(dbMember.Status),
		})
	}
	return members, nil
}

// CountByUserId 获取用户加入的聊天室总数
func (r *chatRoomMemberRepo) CountByUserId(ctx context.Context, userId int64) (int64, error) {
	var count int64
	err := db.DB.Model(&models.ChatRoomMember{}).Where("user_id = ?", userId).Count(&count).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return count, nil
}

func (r *chatRoomMemberRepo) Get(ctx context.Context, roomId, userId int64) (*pb.ChatRoomMember, error) {
	var dbMember models.ChatRoomMember
	err := db.DB.Where("room_id = ? AND user_id = ?", roomId, userId).First(&dbMember).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	// 转换为protobuf格式
	member := &pb.ChatRoomMember{
		RoomId:    int64(dbMember.RoomID),
		UserId:    int64(dbMember.UserID),
		Nickname:  dbMember.Nickname,
		AvatarUrl: dbMember.AvatarURL,
		JoinTime:  dbMember.JoinTime.Unix(),
		Status:    int32(dbMember.Status),
	}
	return member, nil
}

// IncrOnlineCount 增加在线人数
func (r *chatRoomRepo) IncrOnlineCount(ctx context.Context, roomId int64) error {
	return db.DB.Model(&models.ChatRoom{}).Where("room_id = ?", roomId).UpdateColumn("online_count", gorm.Expr("online_count + ?", 1)).Error
}

// DecrOnlineCount 减少在线人数
func (r *chatRoomRepo) DecrOnlineCount(ctx context.Context, roomId int64) error {
	return db.DB.Model(&models.ChatRoom{}).Where("room_id = ?", roomId).UpdateColumn("online_count", gorm.Expr("online_count - ?", 1)).Error
}

// ListByCreatorId 获取用户创建的聊天室列表
func (r *chatRoomRepo) ListByCreatorId(ctx context.Context, creatorId int64) ([]*pb.ChatRoom, error) {
	var modelChatRooms []models.ChatRoom
	err := db.DB.Where("creator_id = ?", creatorId).Find(&modelChatRooms).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 转换为 proto 消息列表
	chatRooms := make([]*pb.ChatRoom, 0, len(modelChatRooms))
	for _, room := range modelChatRooms {
		chatRooms = append(chatRooms, &pb.ChatRoom{
			RoomId:         room.RoomID,
			Name:           room.Name,
			AvatarUrl:      room.AvatarURL,
			Introduction:   room.Introduction,
			OnlineCount:    room.OnlineCount,
			MemberCount:    room.MemberCount,
			MaxMemberCount: room.MaxMemberCount,
			Extra:          room.Extra,
			CreateTime:     room.CreateTime.Unix(),
			UpdateTime:     room.UpdateTime.Unix(),
		})
	}
	return chatRooms, nil
}

// ListByUserId 获取用户加入的聊天室列表
func (r *chatRoomRepo) ListByUserId(ctx context.Context, userId int64, offset, limit int32) ([]*pb.ChatRoom, error) {
	var modelChatRooms []models.ChatRoom
	err := db.DB.Table("chat_room").
		Joins("JOIN chat_room_member ON chat_room.room_id = chat_room_member.room_id").
		Where("chat_room_member.user_id = ?", userId).
		Offset(int(offset)).Limit(int(limit)).
		Find(&modelChatRooms).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 转换为 proto 消息列表
	chatRooms := make([]*pb.ChatRoom, 0, len(modelChatRooms))
	for _, room := range modelChatRooms {
		chatRooms = append(chatRooms, &pb.ChatRoom{
			RoomId:         room.RoomID,
			Name:           room.Name,
			AvatarUrl:      room.AvatarURL,
			Introduction:   room.Introduction,
			CreatorId:      room.CreatorId,
			OnlineCount:    room.OnlineCount,
			MemberCount:    room.MemberCount,
			MaxMemberCount: room.MaxMemberCount,
			Extra:          room.Extra,
			CreateTime:     room.CreateTime.Unix(),
			UpdateTime:     room.UpdateTime.Unix(),
			Level:          calculateRoomLevel(room.MemberCount),
		})
	}
	return chatRooms, nil
}

func (r *chatRoomRepo) CheckPermissionsByUserId(ctx context.Context, req *pb.CheckPermissionsByUserIdReq) (*pb.CheckPermissionsByUserIdResp, error) {
	var tokenHoldings models.TokenHolding
	result := db.DB.Table("chat_room").
		Joins("INNER JOIN token ON chat_room.creator_id = token.user_id").
		Joins("INNER JOIN token_holdings ON token.token_address = token_holdings.token_address").
		Where("token_holdings.user_id = ? AND chat_room.room_id = ?", req.UserId, req.RoomId).
		Limit(1).
		Find(&tokenHoldings) // 只查询常量值1

	exists := result.RowsAffected > 0
	permission := &pb.CheckPermissionsByUserIdResp{
		HasPermission: exists,
	}
	log.Printf("----------- %v", permission)
	if result.Error != nil {
		// 处理数据库错误
		return permission, gerrors.WrapError(result.Error)
	}
	log.Printf("----------- %v", permission)
	return permission, nil
}
