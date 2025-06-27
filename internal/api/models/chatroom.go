package models

import (
	"time"
)

// ChatRoom 聊天室模型
type ChatRoom struct {
	ID             uint      `gorm:"primarykey" json:"id"`                                  // 自增主键
	RoomID         int64     `gorm:"uniqueIndex;not null" json:"room_id"`                   // 聊天室ID
	Name           string    `gorm:"size:50;not null" json:"name"`                          // 聊天室名称
	AvatarURL      string    `gorm:"size:255;not null" json:"avatar_url"`                   // 聊天室头像
	Introduction   string    `gorm:"size:255;not null" json:"introduction"`                 // 聊天室简介
	CreatorId      int64     `gorm:"column:creator_id;not null"`                            // 创建者ID
	UserNum        int32     `gorm:"default:0;not null" json:"user_num"`                    // 当前人数
	OnlineCount    int32     `gorm:"default:0;not null" json:"online_count"`                // 在线人数
	MemberCount    int32     `gorm:"default:0;not null" json:"member_count"`                // 成员总数
	MaxMemberCount int32     `gorm:"default:1000;not null" json:"max_member_count"`         // 最大成员数量
	Status         int8      `gorm:"default:1;not null" json:"status"`                      // 状态 1:正常 2:已关闭
	Extra          string    `gorm:"size:1024;default:'';not null" json:"extra"`            // 附加属性
	CreateTime     time.Time `gorm:"default:CURRENT_TIMESTAMP;not null" json:"create_time"` // 创建时间
	UpdateTime     time.Time `gorm:"default:CURRENT_TIMESTAMP;not null" json:"update_time"` // 更新时间
}

// ChatRoomMember 聊天室成员模型
type ChatRoomMember struct {
	ID         uint      `gorm:"primarykey" json:"id"`                                                // 自增主键
	RoomID     uint64    `gorm:"not null;uniqueIndex:idx_room_user" json:"room_id"`                   // 聊天室ID
	UserID     uint64    `gorm:"not null;uniqueIndex:idx_room_user;index:idx_user_id" json:"user_id"` // 用户ID
	Nickname   string    `gorm:"size:50;not null" json:"nickname"`                                    // 昵称
	AvatarURL  string    `gorm:"size:255;not null" json:"avatar_url"`                                 // 头像
	JoinTime   time.Time `gorm:"default:CURRENT_TIMESTAMP;not null" json:"join_time"`                 // 加入时间
	Extra      string    `gorm:"size:1024;default:'';not null" json:"extra"`                          // 附加属性
	Status     int8      `gorm:"default:1;not null" json:"status"`                                    // 状态 1:正常 2:已退出
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP;not null" json:"create_time"`               // 创建时间
	UpdateTime time.Time `gorm:"default:CURRENT_TIMESTAMP;not null" json:"update_time"`               // 更新时间
}

// TableName 指定表名
func (ChatRoomMember) TableName() string {
	return "chat_room_member"
}
