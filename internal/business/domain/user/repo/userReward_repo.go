package repo

import "gim/internal/business/domain/user/model"

type userRewardRepo struct{}

var UserReward = new(userRewardRepo)

func (u *userRewardRepo) Get(userId int64) (*model.UserReward, error) {
	return UserRewardDao.Get(userId)
}

func (u *userRewardRepo) Add(userReward model.UserReward) (int64, error) {
	return UserRewardDao.Add(userReward)
}
