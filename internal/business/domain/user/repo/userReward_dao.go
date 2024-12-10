package repo

import (
	"errors"
	"gim/internal/business/domain/user/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"github.com/jinzhu/gorm"
)

type userRewardDao struct {
}

var UserRewardDao = new(userRewardDao)

func (d *userRewardDao) Get(userId int64) (*model.UserReward, error) {
	var userReward = model.UserReward{UserID: userId}
	err := db.DB.First(&userReward).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gerrors.WrapError(err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &userReward, err
}

func (d *userRewardDao) Add(userReward model.UserReward) (int64, error) {
	err := db.DB.Create(&userReward).Error
	return userReward.ID, gerrors.WrapError(err)
}
