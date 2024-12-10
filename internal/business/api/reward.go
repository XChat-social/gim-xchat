package api

import (
	"context"
	app2 "gim/internal/business/domain/user/app"
	"gim/pkg/protocol/pb"
)

type RewardServer struct {
	pb.UnimplementedRewardServiceServer
}

func (s *RewardServer) DailySignIn(ctx context.Context, req *pb.DailySignInReq) (*pb.DailySignInResp, error) {
	rewardAmount, message, err := app2.UserApp.DailySignIn(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.DailySignInResp{
		RewardAmount: int32(rewardAmount),
		Message:      message,
	}, nil
}

func (s *RewardServer) ClaimSevenDayReward(ctx context.Context, req *pb.ClaimSevenDayRewardReq) (*pb.ClaimSevenDayRewardResp, error) {
	rewardAmount, message, err := app2.UserApp.ClaimSevenDayReward(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.ClaimSevenDayRewardResp{
		RewardAmount: int32(rewardAmount),
		Message:      message,
	}, nil
}
