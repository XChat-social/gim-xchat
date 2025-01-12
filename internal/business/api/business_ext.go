package api

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	app2 "gim/internal/business/domain/user/app"
	"gim/pkg/db"
	"gim/pkg/grpclib"
	"gim/pkg/protocol/pb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/storyicon/sigverify"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type BusinessExtServer struct {
	pb.UnsafeBusinessExtServer
}

const (
	twitterAuthorizeURL = "https://twitter.com/i/oauth2/authorize"
	twitterTokenURL     = "https://api.twitter.com/2/oauth2/token"
	twitterUserInfoURL  = "https://api.twitter.com/2/users/me"
	clientID            = "YS01bVJhaXdIdEN4X3N5cjNlQzM6MTpjaQ"                 // 替换为实际的 Client ID
	clientSecret        = "JbYIsa77FIbRVc_ZY4238KaPV3Y-K-G5ylbOfuDHgVr8WtgEvO" // 替换为实际的 Client Secret
	redirectURI         = "https://api.xchat.social/api/twitter/signin"
	officialTwitterID   = "1837782128660017152" // 替换为 XChat 官方 Twitter 的实际 ID
	stateTTL            = 5 * time.Minute       // Redis state & verifier TTL
)

func (s *BusinessExtServer) SignIn(ctx context.Context, req *pb.SignInReq) (*pb.SignInResp, error) {
	isNew, userId, token, err := app2.AuthApp.SignIn(ctx, req.PhoneNumber, req.Code, req.DeviceId)
	if err != nil {
		return nil, err
	}
	return &pb.SignInResp{
		IsNew:  isNew,
		UserId: userId,
		Token:  token,
	}, nil
}

func (s *BusinessExtServer) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return &pb.GetUserResp{
			Code:    1,
			Message: "Failed to get user ID from context",
		}, err
	}
	user, err := app2.UserApp.GetNew(ctx, userId)
	if err != nil {
		return &pb.GetUserResp{
			Code:    1,
			Message: "Failed to get user information",
		}, err
	}
	return &pb.GetUserResp{
		Code:    0,
		Message: "Success",
		User:    user,
	}, nil
}

func (s *BusinessExtServer) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*emptypb.Empty, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}

	return new(emptypb.Empty), app2.UserApp.Update(ctx, userId, req)
}

func (s *BusinessExtServer) SearchUser(ctx context.Context, req *pb.SearchUserReq) (*pb.SearchUserResp, error) {
	users, err := app2.UserApp.Search(ctx, req.Key)
	return &pb.SearchUserResp{Users: users}, err
}

// GetTwitterAuthorizeURL 获取 Twitter 授权 URL
func (s *BusinessExtServer) GetTwitterAuthorizeURL(ctx context.Context, req *emptypb.Empty) (*pb.TwitterAuthorizeURLResp, error) {
	state := generateRandomState()
	codeVerifier := generateCodeVerifier()
	codeChallenge := generateCodeChallenge(codeVerifier)

	if err := saveToRedis(state, codeVerifier); err != nil {
		return &pb.TwitterAuthorizeURLResp{
			Code:    1,
			Message: "Failed to save state and code_verifier",
		}, err
	}

	authorizeURL := fmt.Sprintf(
		"%s?response_type=code&client_id=%s&redirect_uri=%s&scope=tweet.read users.read follows.read follows.write&state=%s&code_challenge=%s&code_challenge_method=S256",
		twitterAuthorizeURL, clientID, url.QueryEscape(redirectURI), state, codeChallenge,
	)

	return &pb.TwitterAuthorizeURLResp{
		Code:    0,
		Message: "Twitter authorize URL generated",
		Url:     authorizeURL,
	}, nil
}

// TwitterSignIn 实现推特登录
func (s *BusinessExtServer) TwitterSignIn(ctx context.Context, req *pb.TwitterSignInReq) (*pb.TwitterSignInResp, error) {
	if req.AuthorizationCode == "" || req.State == "" {
		return &pb.TwitterSignInResp{
			Code:    1,
			Message: "Authorization code and state are required",
		}, errors.New("invalid parameters")
	}

	codeVerifier, err := getFromRedis(req.State)
	if err != nil {
		return &pb.TwitterSignInResp{
			Code:    1,
			Message: "Invalid or expired state",
		}, err
	}
	defer deleteFromRedis(req.State)

	accessToken, err := exchangeCodeForToken(req.AuthorizationCode, codeVerifier)
	if err != nil {
		return &pb.TwitterSignInResp{
			Code:    1,
			Message: "Failed to exchange code for token",
		}, err
	}

	twitterUser, err := getTwitterUserInfo(accessToken)
	if err != nil {
		return &pb.TwitterSignInResp{
			Code:    1,
			Message: "Failed to fetch Twitter user info",
		}, err
	}

	isNew, userId, token, err := app2.AuthApp.TwitterSignIn(ctx, twitterUser.ID, twitterUser.Name, twitterUser.Username, twitterUser.Avatar, accessToken, req.WalletAddress)
	if err != nil {
		return &pb.TwitterSignInResp{
			Code:    1,
			Message: "Failed to sign in user",
		}, err
	}

	getNew, err := app2.UserApp.GetNew(ctx, userId)
	if err != nil {
		return &pb.TwitterSignInResp{
			Code:    1,
			Message: "Failed to get user information",
		}, err
	}

	return &pb.TwitterSignInResp{
		Code:    0,
		Message: "Twitter sign-in successful",
		IsNew:   isNew,
		UserId:  userId,
		Token:   token,
		UserInfo: &pb.User{
			UserId:          userId,
			Nickname:        twitterUser.Name,
			AvatarUrl:       twitterUser.Avatar,
			TwitterId:       twitterUser.ID,
			TwitterUsername: twitterUser.Username,
			InviteCode:      getNew.InviteCode,
		},
	}, nil
}

// WalletSignIn 实现钱包登录
func (s *BusinessExtServer) WalletSignIn(ctx context.Context, req *pb.WalletSignInReq) (*pb.WalletSignInResp, error) {
	if req.WalletAddress == "" || req.Signature == "" {
		return &pb.WalletSignInResp{
			Code:    1,
			Message: "Address and signature are required",
		}, errors.New("invalid parameters")
	}

	// 验证签名
	valid, err := verifySignature(req.WalletAddress, req.Message, req.Signature)
	if err != nil || !valid {
		return &pb.WalletSignInResp{
			Code:    1,
			Message: "Invalid signature",
		}, err
	}

	// 查找或创建用户
	isNew, userId, token, err := app2.AuthApp.WalletSignIn(ctx, req.WalletAddress)
	if err != nil {
		return &pb.WalletSignInResp{
			Code:    1,
			Message: "Failed to sign in user",
		}, err
	}

	getNew, err := app2.UserApp.GetNew(ctx, userId)
	if err != nil {
		return &pb.WalletSignInResp{
			Code:    1,
			Message: "Failed to get user information",
		}, err
	}

	return &pb.WalletSignInResp{
		Code:    0,
		Message: "Wallet sign-in successful",
		IsNew:   isNew,
		UserId:  userId,
		Token:   token,
		UserInfo: &pb.User{
			UserId:        userId,
			WalletAddress: req.WalletAddress,
			InviteCode:    getNew.InviteCode,
		},
	}, nil
}

// verifySignature 验证签名
func verifySignature(address, message, signature string) (bool, error) {
	valid, err := sigverify.VerifyEllipticCurveHexSignatureEx(
		common.HexToAddress(address),
		[]byte(message),
		signature,
	)
	if err != nil {
		return false, fmt.Errorf("签名验证失败: %v", err)
	}
	return valid, nil
}

// DailySignIn 每日签到
func (s *BusinessExtServer) DailySignIn(ctx context.Context, req *emptypb.Empty) (*pb.DailySignInResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return &pb.DailySignInResp{
			Code:    1,
			Message: "Failed to get user ID from context",
		}, err
	}
	message, err := app2.UserApp.DailySignIn(ctx, userId)
	if err != nil {
		return &pb.DailySignInResp{
			Code:    1,
			Message: message,
		}, err
	}

	return &pb.DailySignInResp{
		Code:    0,
		Message: message,
	}, nil
}

// FollowTwitter 关注推特
func (s *BusinessExtServer) FollowTwitter(ctx context.Context, req *emptypb.Empty) (*pb.TwitterFollowResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return &pb.TwitterFollowResp{
			Code:    1,
			Message: "Failed to get user ID from context",
		}, err
	}
	success, message, err := app2.UserApp.FollowTwitter(ctx, userId, officialTwitterID)
	if err != nil {
		return &pb.TwitterFollowResp{
			Code:    1,
			Message: message,
			Success: 0,
		}, err
	}
	return &pb.TwitterFollowResp{
		Code:    0,
		Message: message,
		Success: success,
	}, nil
}

// GetTaskStatus 获取任务状态
func (s *BusinessExtServer) GetTaskStatus(ctx context.Context, req *pb.GetTaskStatusReq) (*pb.GetTaskStatusResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return &pb.GetTaskStatusResp{
			Code:    1,
			Message: "Failed to get user ID from context",
		}, err
	}
	if req.TaskId == 0 {
		return &pb.GetTaskStatusResp{
			Code:    1,
			Message: "Task ID is required",
		}, errors.New("invalid parameters")
	}
	status, err := app2.UserApp.GetTaskStatus(ctx, userId, req.TaskId)
	if err != nil {
		return &pb.GetTaskStatusResp{
			Code:    1,
			Message: "Failed to get task status",
		}, err
	}
	return &pb.GetTaskStatusResp{
		Code:    0,
		Message: "Success",
		Status:  int32(status),
	}, nil
}

// ModifyTaskStatus 用于修改指定任务的状态，主要用于测试
func (s *BusinessExtServer) ModifyTaskStatus(ctx context.Context, req *pb.ModifyTaskStatusReq) (*pb.ModifyTaskStatusResp, error) {
	key := fmt.Sprintf("%s:%d:%d", "task_status", req.UserId, req.TaskId)

	// 将任务状态写入 Redis
	err := db.RedisCli.Set(key, req.Status, 24*time.Hour).Err() // 设置状态，并且过期时间为 24 小时
	if err != nil {
		return &pb.ModifyTaskStatusResp{
			Code:    1, // 错误码
			Message: fmt.Sprintf("Failed to modify task status: %v", err),
		}, fmt.Errorf("failed to modify task status: %w", err)
	}

	// 返回成功响应
	return &pb.ModifyTaskStatusResp{
		Code:    0, // 成功
		Message: fmt.Sprintf("Task %d for user %d status updated to %d", req.TaskId, req.UserId, req.Status),
	}, nil
}

// ClaimTaskReward 领取任务奖励
func (s *BusinessExtServer) ClaimTaskReward(ctx context.Context, req *pb.ClaimTaskRewardReq) (*pb.ClaimTaskRewardResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return &pb.ClaimTaskRewardResp{
			Code:    1,
			Message: "Failed to get user ID from context",
		}, err
	}
	if req.TaskId == 0 {
		return &pb.ClaimTaskRewardResp{
			Code:    1,
			Message: "Task ID is required",
		}, errors.New("invalid parameters")
	}
	message, err := app2.UserApp.ClaimTaskReward(ctx, userId, req.TaskId)
	if err != nil {
		return &pb.ClaimTaskRewardResp{
			Code:    1,
			Message: message,
		}, err
	}
	return &pb.ClaimTaskRewardResp{
		Code:    0,
		Message: message,
	}, nil
}

// FillInviteCode 填写邀请码
func (s *BusinessExtServer) FillInviteCode(ctx context.Context, req *pb.FillInviteCodeReq) (*pb.FillInviteCodeResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return &pb.FillInviteCodeResp{
			Code:    1,
			Message: "Failed to get user ID from context",
		}, err
	}
	message, err := app2.UserApp.FillInviteCode(ctx, userId, req.InviteCode)
	if err != nil {
		return &pb.FillInviteCodeResp{
			Code:    1,
			Message: message,
		}, err
	}
	return &pb.FillInviteCodeResp{
		Code:    0,
		Message: message,
	}, nil
}

// exchangeCodeForToken 用授权码换取 Access Token
func exchangeCodeForToken(code, codeVerifier string) (string, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"client_id":     {clientID},
		"code_verifier": {codeVerifier},
	}
	resp, err := http.PostForm(twitterTokenURL, data)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errorResp)
		return "", fmt.Errorf("failed to get access token: %s (%s)", errorResp.Error, errorResp.ErrorDescription)
	}
	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	return tokenResp.AccessToken, nil
}

// 获取 Twitter 用户信息
func getTwitterUserInfo(accessToken string) (*struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Avatar   string `json:"profile_image_url"`
}, error) {
	req, err := http.NewRequest("GET", twitterUserInfoURL+"?user.fields=profile_image_url", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %d", resp.StatusCode)
	}

	var apiResponse struct {
		Data struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
			Avatar   string `json:"profile_image_url"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &apiResponse.Data, nil
}

// Redis 操作简化
func saveToRedis(state, codeVerifier string) error {
	return db.RedisCli.Set(fmt.Sprintf("twitter:state:%s", state), codeVerifier, stateTTL).Err()
}

func getFromRedis(state string) (string, error) {
	return db.RedisCli.Get(fmt.Sprintf("twitter:state:%s", state)).Result()
}

func deleteFromRedis(state string) {
	_ = db.RedisCli.Del(fmt.Sprintf("twitter:state:%s", state)).Err()
}

// 工具函数
func generateRandomState() string {
	return randomString(16)
}

func generateCodeVerifier() string {
	return randomString(43)
}

func generateCodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func randomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
