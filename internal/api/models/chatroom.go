package models

import (
	"time"

	"gorm.io/gorm"
)

// ChatRoom 聊天室模型
type ChatRoom struct {
	gorm.Model
	RoomID         int64  `gorm:"uniqueIndex;not null" json:"room_id"` // 聊天室ID
	Name           string `gorm:"size:50;not null" json:"name"`        // 聊天室名称
	AvatarURL      string `gorm:"size:255" json:"avatar_url"`          // 头像URL
	Introduction   string `gorm:"size:500" json:"introduction"`        // 简介
	MaxMemberCount int32  `gorm:"not null" json:"max_member_count"`    // 最大成员数
	Extra          string `gorm:"size:1024" json:"extra"`              // 附加字段
	CreatorID      int64  `gorm:"not null" json:"creator_id"`          // 创建者ID
}

// ChatRoomMember 聊天室成员模型
type ChatRoomMember struct {
	gorm.Model
	RoomID   int64     `gorm:"not null;index:idx_room_user" json:"room_id"` // 聊天室ID
	UserID   int64     `gorm:"not null;index:idx_room_user" json:"user_id"` // 用户ID
	JoinTime time.Time `gorm:"not null" json:"join_time"`                   // 加入时间
	Extra    string    `gorm:"size:1024" json:"extra"`                      // 附加字段
}
