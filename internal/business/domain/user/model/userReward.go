package model

import "time"

// UserReward 用户奖励实体
type UserReward struct {
	ID           int64     // 自增主键
	UserID       int64     // 用户ID
	RewardType   int32     // 奖励类型，1:每日签到；2:连续7天
	RewardAmount int32     // 奖励数量
	RewardDate   time.Time // 奖励日期
	CreateTime   time.Time // 创建时间
	UpdateTime   time.Time // 更新时间
}
