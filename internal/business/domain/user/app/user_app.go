package app

import (
	"context"
	"gim/internal/business/domain/user/service"
	"time"

	"gim/internal/business/domain/user/repo"
	"gim/pkg/protocol/pb"
)

type userApp struct{}

var UserApp = new(userApp)

func (*userApp) Get(ctx context.Context, userId int64) (*pb.User, error) {
	user, err := repo.UserRepo.Get(userId)
	return user.ToProto(), err
}

func (*userApp) GetNew(ctx context.Context, userId int64) (*pb.User, error) {
	user, err := repo.UserRepo.GetNew(userId)
	return user.ToProto(), err
}

func (*userApp) Update(ctx context.Context, userId int64, req *pb.UpdateUserReq) error {
	u, err := repo.UserRepo.Get(userId)
	if err != nil {
		return err
	}
	if u == nil {
		return nil
	}

	u.Nickname = req.Nickname
	u.Sex = req.Sex
	u.AvatarUrl = req.AvatarUrl
	u.Extra = req.Extra
	u.UpdateTime = time.Now()

	err = repo.UserRepo.Save(u)
	if err != nil {
		return err
	}
	return nil
}

func (*userApp) GetByIds(ctx context.Context, userIds []int64) (map[int64]*pb.User, error) {
	users, err := repo.UserRepo.GetByIds(userIds)
	if err != nil {
		return nil, err
	}

	pbUsers := make(map[int64]*pb.User, len(users))
	for i := range users {
		pbUsers[users[i].Id] = users[i].ToProto()
	}
	return pbUsers, nil
}

func (*userApp) Search(ctx context.Context, key string) ([]*pb.User, error) {
	users, err := repo.UserRepo.Search(key)
	if err != nil {
		return nil, err
	}

	pbUsers := make([]*pb.User, len(users))
	for i, v := range users {
		pbUsers[i] = v.ToProto()
	}
	return pbUsers, nil
}

func (*userApp) DailySignIn(ctx context.Context, userId int64) (string, error) {
	return service.RewardService.DailySignIn(ctx, userId)
}

func (a *userApp) FollowTwitter(ctx context.Context, userId int64, officialTwitterID string) (int32, string, error) {
	return service.RewardService.FollowTwitter(ctx, userId, officialTwitterID)
}

func (a *userApp) GetTaskStatus(ctx context.Context, userId int64, taskId int64) (int, error) {
	return service.RewardService.GetTaskStatus(ctx, userId, taskId)
}

func (a *userApp) ClaimTaskReward(ctx context.Context, userId int64, taskId int64) (string, error) {
	return service.RewardService.ClaimTaskReward(ctx, userId, taskId)
}
