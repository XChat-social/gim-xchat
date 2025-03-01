package app

import (
	"context"
	"gim/internal/business/domain/user/model"
	"gim/internal/business/domain/user/repo"
	"gim/pkg/gerrors"
	"gim/pkg/protocol/pb"
)

type TokenAppInterface interface {
	Create(ctx context.Context, token *pb.Token) (*pb.Token, error)
	Get(ctx context.Context, tokenAddress string) (*pb.Token, error)
}

type tokenApp struct{}

var TokenApp TokenAppInterface = new(tokenApp)

// Create 创建token
func (a *tokenApp) Create(ctx context.Context, token *pb.Token) (*pb.Token, error) {
	if token.TokenAddress == "" {
		return nil, gerrors.ErrBadRequest
	}

	tokenModel := &model.Token{
		TokenAddress:   token.TokenAddress,
		TokenName:      token.TokenName,
		TokenSymbol:    token.TokenSymbol,
		Decimals:       token.Decimals,
		TotalSupply:    token.TotalSupply,
		CreatorAddress: token.CreatorAddr,
		ChainID:        token.ChainId,
		Status:         token.Status,
	}
	if err := repo.TokenRepo.Create(tokenModel); err != nil {
		return nil, err
	}
	return tokenModel.ToProto(), nil
}

// Get 获取token信息
func (a *tokenApp) Get(ctx context.Context, tokenAddress string) (*pb.Token, error) {
	token, err := repo.TokenRepo.GetByAddress(tokenAddress)
	if err != nil {
		return nil, err
	}
	return token.ToProto(), nil
}

// convertToProtoToken 将model转换为proto
func convertToProtoToken(token *model.Token) *pb.Token {
	return &pb.Token{
		Id:           token.ID,
		UserId:       token.UserID,
		TokenAddress: token.TokenAddress,
		TokenName:    token.TokenName,
		TokenSymbol:  token.TokenSymbol,
		Decimals:     token.Decimals,
		TotalSupply:  token.TotalSupply,
		CreatorAddr:  token.CreatorAddress,
		ChainId:      token.ChainID,
		Status:       token.Status,
		CreatedAt:    token.CreatedAt.Unix(),
		UpdatedAt:    token.UpdatedAt.Unix(),
	}
}
