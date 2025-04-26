package models

import (
	"time"
)

// User 用户模型
type User struct {
	UserID          int64     `json:"user_id" gorm:"primaryKey;column:user_id"`
	Nickname        string    `json:"nickname" gorm:"column:nickname"`
	Sex             int32     `json:"sex" gorm:"column:sex"`
	AvatarURL       string    `json:"avatar_url" gorm:"column:avatar_url"`
	Extra           string    `json:"extra" gorm:"column:extra"`
	CreateTime      time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime      time.Time `json:"update_time" gorm:"column:update_time"`
	TwitterID       string    `json:"twitter_id" gorm:"column:twitter_id"`
	TwitterUsername string    `json:"twitter_username" gorm:"column:twitter_username"`
	XPoint          int32     `json:"xpoint" gorm:"column:xpoint"`
	FollowReward    int32     `json:"follow_reward" gorm:"column:follow_reward"`
	InviteCode      string    `json:"invite_code" gorm:"column:invite_code"`
	InviterCode     string    `json:"inviter_code" gorm:"column:inviter_code"`
	WalletAddress   string    `json:"wallet_address" gorm:"column:wallet_address"`
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}

// Token 代币模型
type Token struct {
	ID           int64     `json:"id" gorm:"primaryKey;column:id"`
	UserID       int64     `json:"user_id" gorm:"column:user_id"`
	TokenAddress string    `json:"token_address" gorm:"column:token_address"`
	TokenName    string    `json:"token_name" gorm:"column:token_name"`
	TokenSymbol  string    `json:"token_symbol" gorm:"column:token_symbol"`
	Decimals     int32     `json:"decimals" gorm:"column:decimals"`
	TotalSupply  string    `json:"total_supply" gorm:"column:total_supply"`
	CreatorAddr  string    `json:"creator_addr" gorm:"column:creator_addr"`
	ChainID      int32     `json:"chain_id" gorm:"column:chain_id"`
	Status       int32     `json:"status" gorm:"column:status"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName 设置表名
func (Token) TableName() string {
	return "tokens"
}

// Task 任务模型
type Task struct {
	ID         int64     `json:"id" gorm:"primaryKey;column:id"`
	UserID     int64     `json:"user_id" gorm:"column:user_id"`
	TaskType   int32     `json:"task_type" gorm:"column:task_type"`
	Status     int32     `json:"status" gorm:"column:status"`
	Reward     int32     `json:"reward" gorm:"column:reward"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"`
}

// TableName 设置表名
func (Task) TableName() string {
	return "tasks"
}
