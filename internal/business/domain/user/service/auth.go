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
func (*authService) TwitterSignIn(ctx context.Context, twitterID, name, username, avatar, accessToken, walletAddress string) (bool, int64, string, int64, error) {
	// 初始化返回值
	var (
		isNew      bool
		signStatus int64 // 0表示登录 1表示绑定
		user       *model.User
		err        error
	)

	// 查询 Twitter 用户
	user, err = repo.UserRepo.GetByTwitterID(twitterID)
	if err != nil {
		return false, 0, "", signStatus, err
	}

	// 如果提供了钱包地址，查询钱包用户
	var addressUser *model.User
	if walletAddress != "" {
		signStatus = 1
		addressUser, err = repo.UserRepo.GetByWalletAddress(walletAddress)
		if err != nil {
			return false, 0, "", signStatus, err
		}
	}

	// 处理用户冲突情况
	if user != nil && addressUser != nil && user.Id != addressUser.Id {
		token := GenerateToken()
		if err = saveAuthToken(user.Id, token, accessToken); err != nil {
			return false, 0, "", signStatus, err
		}
		return false, 0, "", signStatus, gerrors.ErrUserAlreadyExists
	}

	// 处理用户绑定情况
	if err = handleUserBinding(user, addressUser, twitterID, name, username, avatar, walletAddress); err != nil {
		return false, 0, "", signStatus, err
	}

	// 创建新用户
	if user == nil && addressUser == nil {
		user, err = createNewUser(twitterID, name, username, avatar, walletAddress)
		if err != nil {
			return false, 0, "", signStatus, err
		}
		isNew = true
	}

	// 生成并保存 Token
	token := GenerateToken()
	if err = saveAuthToken(user.Id, token, accessToken); err != nil {
		return false, 0, "", signStatus, err
	}

	return isNew, user.Id, token, signStatus, nil
}

// 辅助函数：保存认证 Token
func saveAuthToken(userId int64, token, accessToken string) error {
	return repo.AuthRepo.Set(userId, 0, model.Device{
		Type:        0,
		Token:       token,
		AccessToken: accessToken,
		Expire:      time.Now().AddDate(0, 0, 1).Unix(),
	})
}

// 辅助函数：处理用户绑定
func handleUserBinding(user, addressUser *model.User, twitterID, name, username, avatar, walletAddress string) error {
	if user == nil && addressUser != nil {
		// 绑定 Twitter 信息到钱包用户
		addressUser.TwitterID = twitterID
		addressUser.Nickname = name
		addressUser.TwitterUsername = username
		addressUser.AvatarUrl = avatar
		addressUser.UpdateTime = time.Now()
		return repo.UserRepo.Update(addressUser)
	}
	if user != nil && addressUser == nil && walletAddress != "" {
		// 绑定钱包地址到 Twitter 用户
		user.WalletAddress = walletAddress
		user.UpdateTime = time.Now()
		return repo.UserRepo.Update(user)
	}
	return nil
}

// 辅助函数：创建新用户
func createNewUser(twitterID, name, username, avatar, walletAddress string) (*model.User, error) {
	inviteCode, err := generateUniqueInviteCode()
	if err != nil {
		return nil, err
	}

	user := &model.User{
		TwitterID:       twitterID,
		Nickname:        name,
		TwitterUsername: username,
		AvatarUrl:       avatar,
		WalletAddress:   walletAddress,
		InviteCode:      inviteCode,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}

	if err := repo.UserRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

// WalletSignIn 通过钱包地址登录
func (s *authService) WalletSignIn(ctx context.Context, address string) (bool, int64, string, error) {
	// 查找用户是否已经存在
	user, err := repo.UserRepo.GetByWalletAddress(address)
	if err != nil {
		return false, 0, "", err
	}

	// 如果用户不存在，创建新用户
	var isNew bool
	if user == nil {
		user, err = createWalletUser(address)
		if err != nil {
			return false, 0, "", err
		}
		isNew = true
	}

	// 生成并保存 Token
	token := GenerateToken()
	if err = saveAuthToken(user.Id, token, ""); err != nil {
		return false, 0, "", err
	}

	return isNew, user.Id, token, nil
}

// 辅助函数：创建钱包用户
func createWalletUser(address string) (*model.User, error) {
	inviteCode, err := generateUniqueInviteCode()
	if err != nil {
		return nil, err
	}

	user := &model.User{
		WalletAddress: address,
		InviteCode:    inviteCode,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}

	if err := repo.UserRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
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
