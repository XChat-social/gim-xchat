package service

import (
	"context"
	"encoding/hex"
	"gim/internal/business/domain/user/model"
	"gim/internal/business/domain/user/repo"
	"gim/pkg/gerrors"
	"gim/pkg/protocol/pb"
	"gim/pkg/rpc"
	"math/rand"
	"time"
)

type authService struct{}

var AuthService = new(authService)

// SignIn 登录
func (*authService) SignIn(ctx context.Context, phoneNumber, code string, deviceId int64) (bool, int64, string, error) {
	if !Verify(phoneNumber, code) {
		return false, 0, "", gerrors.ErrBadCode
	}

	user, err := repo.UserRepo.GetByPhoneNumber(phoneNumber)
	if err != nil {
		return false, 0, "", err
	}

	var isNew = false
	if user == nil {
		user = &model.User{
			PhoneNumber: phoneNumber,
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		}
		err := repo.UserRepo.Save(user)
		if err != nil {
			return false, 0, "", err
		}
		isNew = true
	}

	resp, err := rpc.GetLogicIntClient().GetDevice(ctx, &pb.GetDeviceReq{DeviceId: deviceId})
	if err != nil {
		return false, 0, "", err
	}

	// 方便测试
	token := "0"
	//token := util.RandString(40)
	err = repo.AuthRepo.Set(user.Id, resp.Device.DeviceId, model.Device{
		Type:   resp.Device.Type,
		Token:  token,
		Expire: time.Now().AddDate(0, 3, 0).Unix(),
	})
	if err != nil {
		return false, 0, "", err
	}

	return isNew, user.Id, token, nil
}

func Verify(phoneNumber, code string) bool {
	// 假装他成功了
	return true
}

// Auth 验证用户是否登录
func (*authService) Auth(ctx context.Context, userId, deviceId int64, token string) error {
	device, err := repo.AuthRepo.Get(userId, deviceId)
	if err != nil {
		return err
	}

	if device == nil {
		return gerrors.ErrUnauthorized
	}

	if device.Expire < time.Now().Unix() {
		return gerrors.ErrUnauthorized
	}

	if device.Token != token {
		return gerrors.ErrUnauthorized
	}
	return nil
}

// TwitterSignIn 实现 Twitter 登录逻辑
func (*authService) TwitterSignIn(ctx context.Context, twitterID, name, username, avatar, accessToken, walletAddress string) (bool, int64, string, error) {
	var user *model.User
	var addressUser *model.User
	var err error
	user, err = repo.UserRepo.GetByTwitterID(twitterID)
	addressUser, err = repo.UserRepo.GetByWalletAddress(walletAddress)

	if err != nil {
		return false, 0, "", err
	}

	if user == nil && addressUser != nil {
		addressUser.TwitterID = twitterID
		addressUser.Nickname = name
		addressUser.TwitterUsername = username
		addressUser.AvatarUrl = avatar
		addressUser.UpdateTime = time.Now()
		if err = repo.UserRepo.Update(addressUser); err != nil {
			return false, 0, "", err
		}
	}

	if user != nil && addressUser == nil {
		user.WalletAddress = walletAddress
		user.UpdateTime = time.Now()
		if err = repo.UserRepo.Update(user); err != nil {
			return false, 0, "", err
		}
	}

	var isNew = false
	if user == nil {
		// Step 2: 如果用户不存在，创建新用户
		isNew = true
		inviteCode, err := generateUniqueInviteCode()
		if err != nil {
			return false, 0, "", err
		}
		user = &model.User{
			TwitterID:       twitterID,
			Nickname:        name,
			TwitterUsername: username,
			AvatarUrl:       avatar,
			InviteCode:      inviteCode,
			CreateTime:      time.Now(),
			UpdateTime:      time.Now(),
		}
		if err := repo.UserRepo.Save(user); err != nil {
			return false, 0, "", err
		}
	}

	// Step 3: 生成 Token
	token := GenerateToken()

	// Step 4: 保存 Token 信息
	err = repo.AuthRepo.Set(user.Id, 0, model.Device{ // DeviceId 设为 0 表示无需设备信息
		Type:        0,
		Token:       token,
		AccessToken: accessToken,
		Expire:      time.Now().AddDate(0, 0, 1).Unix(),
	})
	if err != nil {
		return false, 0, "", err
	}

	return isNew, user.Id, token, nil
}

// WalletSignIn 通过钱包地址登录
func (s *authService) WalletSignIn(ctx context.Context, address string) (bool, int64, string, error) {
	// Step 1: 查找用户是否已经存在
	user, err := repo.UserRepo.GetByWalletAddress(address)
	if err != nil {
		return false, 0, "", err
	}

	var isNew bool
	if user == nil {
		// Step 2: 如果用户不存在，创建新用户
		isNew = true
		inviteCode, err := generateUniqueInviteCode()
		if err != nil {
			return false, 0, "", err
		}

		// 假设没有其他附加数据（比如昵称等），可以直接根据地址创建用户
		user = &model.User{
			WalletAddress: address,
			InviteCode:    inviteCode,
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		}

		// 保存新用户
		if err := repo.UserRepo.Save(user); err != nil {
			return false, 0, "", err
		}
	}

	// Step 3: 生成会话 Token
	token := GenerateToken()

	// Step 4: 保存 Token 信息
	err = repo.AuthRepo.Set(user.Id, 0, model.Device{
		Type:   0,
		Token:  token,
		Expire: time.Now().AddDate(0, 0, 1).Unix(),
	})
	if err != nil {
		return false, 0, "", err
	}

	// 返回用户是否为新用户、用户 ID 和生成的 Token
	return isNew, user.Id, token, nil
}

// generateUniqueInviteCode 生成唯一邀请码
func generateUniqueInviteCode() (string, error) {
	const codeLength = 8
	for {
		code := randomString(codeLength)
		exists, err := repo.UserRepo.InviteCodeExists(code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateToken 生成会话 Token
func GenerateToken() string {
	randomBytes := make([]byte, 16) // 16 字节随机数
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}
	randomString := hex.EncodeToString(randomBytes)
	return randomString
}
