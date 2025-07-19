package models

import (
	"time"
)

// User 表示用户信息
type User struct {
	ID              uint64    `json:"id" gorm:"primaryKey;column:id"`                  // 自增主键
	PhoneNumber     string    `json:"phone_number" gorm:"column:phone_number"`         // 手机号
	Nickname        string    `json:"nickname" gorm:"column:nickname"`                 // 昵称
	Sex             int8      `json:"sex" gorm:"column:sex"`                           // 性别：0 未知，1 男，2 女
	AvatarURL       string    `json:"avatar_url" gorm:"column:avatar_url"`             // 用户头像链接
	Extra           string    `json:"extra" gorm:"column:extra"`                       // 附加属性（可存储 JSON 等扩展信息）
	TwitterID       string    `json:"twitter_id" gorm:"column:twitter_id"`             // Twitter ID（可为空）
	TwitterUsername string    `json:"twitter_username" gorm:"column:twitter_username"` // Twitter 用户名
	CreateTime      time.Time `json:"create_time" gorm:"column:create_time"`           // 创建时间
	UpdateTime      time.Time `json:"update_time" gorm:"column:update_time"`           // 更新时间
	XPoint          float64   `json:"xpoint" gorm:"column:xpoint"`                     // 当前积分
	FollowReward    int       `json:"follow_reward" gorm:"column:follow_reward"`       // 推特关注奖励领取状态：false=未领取，true=已领取
	InviterCode     string    `json:"inviter_code" gorm:"column:inviter_code"`         // 填写的邀请码
	InviteCode      string    `json:"invite_code" gorm:"column:invite_code"`           // 用户的邀请码
	WalletAddress   string    `json:"wallet_address" gorm:"column:wallet_address"`     // 钱包地址
}

// TableName 设置表名
func (User) TableName() string {
	return "user"
}

// Token 表示用户创建的 Token 信息
type Token struct {
	ID             int64     `json:"id" gorm:"primaryKey;column:id"`                // 主键ID
	UserID         int64     `json:"user_id" gorm:"column:user_id"`                 // 创建者用户ID
	TokenAddress   string    `json:"token_address" gorm:"column:token_address"`     // Token 合约地址
	TokenName      string    `json:"token_name" gorm:"column:token_name"`           // Token 名称
	TokenSymbol    string    `json:"token_symbol" gorm:"column:token_symbol"`       // Token 符号
	Decimals       int       `json:"decimals" gorm:"column:decimals"`               // Token 精度（默认 18）
	TotalSupply    string    `json:"total_supply" gorm:"column:total_supply"`       // Token 发行总量
	CreatorAddress string    `json:"creator_address" gorm:"column:creator_address"` // 创建者钱包地址
	ChainID        int       `json:"chain_id" gorm:"column:chain_id"`               // 链 ID（如：1=Ethereum, 56=BSC）
	IconUrl        string    `json:"icon_url" gorm:"column:icon_url"`
	Status         int8      `json:"status" gorm:"column:status"`         // 状态：1=正常，0=禁用
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"` // 创建时间
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at"` // 更新时间
}

// TableName 设置表名
func (Token) TableName() string {
	return "token"
}

// XPointLog 积分变动日志
type XPointLog struct {
	ID           uint64    `json:"id" gorm:"primaryKey;column:id"`                       // 日志ID
	UserID       uint64    `json:"user_id" gorm:"column:user_id"`                        // 用户ID
	ChangeAmount int       `json:"change_amount" gorm:"column:change_amount"`            // 积分变化值
	Reason       string    `json:"reason" gorm:"column:reason"`                          // 变动原因（如每日签到、连续签到奖励等）
	CreateTime   time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"` // 变动时间
}

// TableName 指定表名
func (XPointLog) TableName() string {
	return "xpoint_log"
}

// TokenHolding Token持有记录
type TokenHolding struct {
	ID           int64     `json:"id" gorm:"primaryKey;column:id"`            // 主键ID
	UserID       int64     `json:"user_id" gorm:"column:user_id"`             // 用户ID
	TokenAddress string    `json:"token_address" gorm:"column:token_address"` // Token地址
	Amount       float64   `json:"amount" gorm:"column:amount"`               // 持有数量
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`       // 创建时间
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`       // 更新时间
}

// TableName 设置表名
func (TokenHolding) TableName() string {
	return "token_holdings"
}

// TokenTransaction Token交易记录
type TokenTransaction struct {
	ID              int64     `json:"id" gorm:"primaryKey;column:id"`                  // 主键ID
	UserID          int64     `json:"user_id" gorm:"column:user_id"`                   // 用户ID
	TokenAddress    string    `json:"token_address" gorm:"column:token_address"`       // Token地址
	TransactionType string    `json:"transaction_type" gorm:"column:transaction_type"` // 交易类型：buy, sell
	Amount          float64   `json:"amount" gorm:"column:amount"`                     // 交易数量
	Price           float64   `json:"price" gorm:"column:price"`                       // 单价
	TotalValue      float64   `json:"total_value" gorm:"column:total_value"`           // 总价值
	TransactionHash string    `json:"transaction_hash" gorm:"column:transaction_hash"` // 交易哈希
	Status          string    `json:"status" gorm:"column:status"`                     // 交易状态：pending, completed, failed
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at"`             // 创建时间
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:updated_at"`             // 更新时间
}

// TableName 设置表名
func (TokenTransaction) TableName() string {
	return "token_transactions"
}
