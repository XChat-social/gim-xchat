package room

import (
	"context"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"gim/pkg/protocol/pb"

	"gorm.io/gorm"
)

type chatRoomRepo struct{}

var ChatRoomRepo = new(chatRoomRepo)

type chatRoomMemberRepo struct{}

var ChatRoomMemberRepo = new(chatRoomMemberRepo)

// Add 添加聊天室
func (r *chatRoomRepo) Add(ctx context.Context, chatRoom *pb.ChatRoom) error {
	return db.DB.Create(chatRoom).Error
}

// Get 获取聊天室信息
func (r *chatRoomRepo) Get(ctx context.Context, roomId int64) (*pb.ChatRoom, error) {
	var chatRoom pb.ChatRoom
	err := db.DB.Where("room_id = ?", roomId).First(&chatRoom).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return &chatRoom, nil
}

// List 获取聊天室列表
func (r *chatRoomRepo) List(ctx context.Context, offset, limit int32) ([]*pb.ChatRoom, error) {
	var chatRooms []*pb.ChatRoom
	err := db.DB.Offset(int(offset)).Limit(int(limit)).Find(&chatRooms).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return chatRooms, nil
}

// Count 获取聊天室总数
func (r *chatRoomRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := db.DB.Model(&pb.ChatRoom{}).Count(&count).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return count, nil
}

// Add 添加聊天室成员
func (r *chatRoomMemberRepo) Add(ctx context.Context, member *pb.ChatRoomMember) error {
	return db.DB.Create(member).Error
}

// Delete 删除聊天室成员
func (r *chatRoomMemberRepo) Delete(ctx context.Context, roomId, userId int64) error {
	return db.DB.Where("room_id = ? AND user_id = ?", roomId, userId).Delete(&pb.ChatRoomMember{}).Error
}

// Count 获取聊天室成员总数
func (r *chatRoomMemberRepo) Count(ctx context.Context, roomId int64) (int64, error) {
	var count int64
	err := db.DB.Model(&pb.ChatRoomMember{}).Where("room_id = ?", roomId).Count(&count).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return count, nil
}

// List 获取聊天室成员列表
func (r *chatRoomMemberRepo) List(ctx context.Context, roomId int64, offset, limit int32) ([]*pb.ChatRoomMember, error) {
	var members []*pb.ChatRoomMember
	err := db.DB.Where("room_id = ?", roomId).Offset(int(offset)).Limit(int(limit)).Find(&members).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return members, nil
}

// ListByUserId 获取用户加入的所有聊天室
func (r *chatRoomMemberRepo) ListByUserId(ctx context.Context, userId int64) ([]*pb.ChatRoomMember, error) {
	var members []*pb.ChatRoomMember
	err := db.DB.Where("user_id = ?", userId).Find(&members).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return members, nil
}

func (r *chatRoomMemberRepo) Get(ctx context.Context, roomId, userId int64) (*pb.ChatRoomMember, error) {
	var member pb.ChatRoomMember
	err := db.DB.Where("room_id = ? AND user_id = ?", roomId, userId).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

// IncrOnlineCount 增加聊天室在线人数
func (r *chatRoomRepo) IncrOnlineCount(ctx context.Context, roomId int64) error {
	return db.DB.Model(&pb.ChatRoom{}).Where("room_id = ?", roomId).UpdateColumn("online_count", gorm.Expr("online_count + ?", 1)).Error
}

// DecrOnlineCount 减少聊天室在线人数
func (r *chatRoomRepo) DecrOnlineCount(ctx context.Context, roomId int64) error {
	return db.DB.Model(&pb.ChatRoom{}).Where("room_id = ?", roomId).UpdateColumn("online_count", gorm.Expr("online_count - ?", 1)).Error
}
