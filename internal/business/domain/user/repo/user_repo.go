package repo

import (
	"gim/internal/business/domain/user/model"
)

type userRepo struct{}

var UserRepo = new(userRepo)

// Get 获取单个用户
func (*userRepo) Get(userId int64) (*model.User, error) {
	user, err := UserCache.Get(userId)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	user, err = UserDao.Get(userId)
	if err != nil {
		return nil, err
	}

	if user != nil {
		err = UserCache.Set(*user)
		if err != nil {
			return nil, err
		}
	}
	return user, err
}

func (*userRepo) GetNew(userId int64) (*model.User, error) {
	user, err := UserDao.Get(userId)
	if err != nil {
		return nil, err
	}
	return user, err
}

// GetByPhoneNumber 通过手机号获取用户
func (*userRepo) GetByPhoneNumber(phoneNumber string) (*model.User, error) {
	return UserDao.GetByPhoneNumber(phoneNumber)
}

// GetByTwitterID 通过 TwitterID 获取用户
func (*userRepo) GetByTwitterID(twitterID string) (*model.User, error) {
	return UserDao.GetByTwitterID(twitterID)
}

// GetByIds 获取多个用户
func (*userRepo) GetByIds(userIds []int64) ([]model.User, error) {
	return UserDao.GetByIds(userIds)
}

// Search 搜索用户
func (*userRepo) Search(key string) ([]model.User, error) {
	return UserDao.Search(key)
}

// Save 保存用户
func (*userRepo) Save(user *model.User) error {
	userId := user.Id
	err := UserDao.Save(user)
	if err != nil {
		return err
	}

	if userId != 0 {
		err = UserCache.Del(user.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *userRepo) InviteCodeExists(code string) (bool, error) {
	return UserDao.InviteCodeExists(code)
}

func (*userRepo) GetUserByInviteCode(inviteCode string) (*model.User, error) {
	return UserDao.GetUserByInviteCode(inviteCode)
}

func (*userRepo) UpdateInviteCodeStatus(id int64, code string) error {
	return UserDao.UpdateInviteCodeStatus(id, code)
}

func (r *userRepo) GetByWalletAddress(address string) (*model.User, error) {
	return UserDao.GetByWalletAddress(address)
}
