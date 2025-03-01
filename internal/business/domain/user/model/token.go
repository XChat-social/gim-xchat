package model

import (
	"gim/pkg/protocol/pb"
	"time"
)

type Token struct {
	ID             int64     // 主键ID
	UserID         int64     // 创建者用户ID
	TokenAddress   string    // token合约地址
	TokenName      string    // token名称
	TokenSymbol    string    // token符号
	Decimals       int32     // token精度
	TotalSupply    string    // 发行总量
	CreatorAddress string    // 创建者钱包地址
	ChainID        int32     // 链ID
	Status         int32     // 状态 1:正常 0:禁用
	CreatedAt      time.Time // 创建时间
	UpdatedAt      time.Time // 更新时间
}

func (t Token) ToProto() *pb.Token {
	return &pb.Token{
		Id:           t.ID,
		UserId:       t.UserID,
		TokenAddress: t.TokenAddress,
		TokenName:    t.TokenName,
		TokenSymbol:  t.TokenSymbol,
		Decimals:     t.Decimals,
		TotalSupply:  t.TotalSupply,
		CreatorAddr:  t.CreatorAddress,
		ChainId:      t.ChainID,
		Status:       t.Status,
		CreatedAt:    t.CreatedAt.Unix(),
		UpdatedAt:    t.UpdatedAt.Unix(),
	}
}
