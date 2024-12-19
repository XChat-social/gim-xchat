package model

import (
	"gim/pkg/protocol/pb"
	"time"
)

// User 账户
type User struct {
	Id              int64     // 用户id
	PhoneNumber     string    // 手机号
	Nickname        string    // 昵称
	Sex             int32     // 性别，1:男；2:女
	AvatarUrl       string    // 用户头像
	Extra           string    // 附加属性
	Xpoint          int32     // 用户当前积分
	CreateTime      time.Time // 创建时间
	UpdateTime      time.Time // 更新时间
	TwitterID       string    // 推特ID
	TwitterUsername string    // 推特用户名
	FollowReward    int32     // 推特关注状态领取: 0 未领取, 1 已领取
	InviteCode      string    // 邀请码
	InviterCode     string    // 邀请人邀请码
}

func (u *User) ToProto() *pb.User {
	if u == nil {
		return nil
	}

	return &pb.User{
		UserId:          u.Id,
		Nickname:        u.Nickname,
		Sex:             u.Sex,
		AvatarUrl:       u.AvatarUrl,
		Extra:           u.Extra,
		CreateTime:      u.CreateTime.Unix(),
		UpdateTime:      u.UpdateTime.Unix(),
		TwitterId:       u.TwitterID,
		TwitterUsername: u.TwitterUsername,
		Xpoint:          u.Xpoint,
		FollowReward:    u.FollowReward,
		InviteCode:      u.InviteCode,
	}
}
