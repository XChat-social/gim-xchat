package model

import "time"

// XPointLog 定义积分变动日志实体
type XPointLog struct {
	ID           uint64    `gorm:"primary_key;auto_increment" json:"id"`     // 日志ID，自增主键
	UserID       uint64    `gorm:"not null" json:"user_id"`                  // 用户ID
	ChangeAmount int       `gorm:"not null" json:"change_amount"`            // 积分变化值
	Reason       string    `gorm:"type:varchar(255);not null" json:"reason"` // 变动原因
	CreateTime   time.Time `gorm:"autoCreateTime" json:"create_time"`        // 创建时间
}

// TableName 指定数据库表名
func (XPointLog) TableName() string {
	return "xpoint_log"
}
