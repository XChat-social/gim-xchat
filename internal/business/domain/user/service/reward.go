package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gim/internal/business/domain/user/model"
	"gim/internal/business/domain/user/repo"
	"gim/pkg/db"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"math/rand"
	"net/http"
	"time"
)

type rewardService struct{}

var RewardService = new(rewardService)

const (
	TaskDailySignIn     = 1001 // 每日签到
	TaskSevenDaySignIn  = 1002 // 签到七天
	TaskFollowTwitter   = 1003 // 关注推特
	taskStatusKeyPrefix = "task_status"
)

// DailySignIn 处理每日签到逻辑
func (*rewardService) DailySignIn(ctx context.Context, userID int64) (string, error) {
	// Redis Key
	key := fmt.Sprintf("signin:%d", userID)

	// 获取当天的位偏移
	offset := time.Now().UTC().Unix() / 86400 // 当前时间的天偏移量

	// 检查当天是否已签到
	signed, err := db.RedisCli.GetBit(key, offset).Result()
	if err != nil {
		return "", fmt.Errorf("failed to query redis: %w", err)
	}
	if signed == 1 {
		return "Already signed in today", nil
	}

	// 设置当天签到状态
	err = db.RedisCli.SetBit(key, offset, 1).Err()
	if err != nil {
		return "", fmt.Errorf("failed to set redis bit: %w", err)
	}

	// 构造任务状态 Redis Key
	taskKey := fmt.Sprintf("%s:%d:%d", taskStatusKeyPrefix, userID, TaskDailySignIn)

	// 将每日签到任务状态设置为待领取 (3)
	err = db.RedisCli.Set(taskKey, 3, 24*time.Hour).Err() // 设置过期时间为 1 天
	if err != nil {
		// 如果任务状态更新失败，手动回滚 Redis
		rollbackErr := db.RedisCli.SetBit(key, offset, 0).Err()
		if rollbackErr != nil {
			return "", fmt.Errorf("failed to set task status and rollback redis: %w", rollbackErr)
		}
		return "", fmt.Errorf("failed to set task status: %w", err)
	}

	return "Sign-in successful", nil
}

// FollowTwitter 关注推特并更新状态
func (s *rewardService) FollowTwitter(ctx context.Context, userId int64, officialTwitterID string) (int32, string, error) {
	user, err2 := repo.UserRepo.GetNew(userId)
	if err2 != nil {
		return 0, "", fmt.Errorf("failed to get user: %w", err2)
	}
	twitterId := user.TwitterID

	// 检查奖励状态
	if user.FollowReward == 1 {
		return 0, "Reward already claimed.", nil
	}

	// 获取 Twitter Access Token
	device, err := repo.AuthRepo.Get(userId, 0)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get Twitter access token: %w", err)
	}

	// 调用 API 创建关注关系
	isFollowing, err := followUser(device.AccessToken, twitterId, officialTwitterID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to follow official Twitter account: %w", err)
	}

	// 如果未成功创建关注关系
	if !isFollowing {
		return 0, "Failed to follow the official Twitter account.", nil
	}

	// 保存状态到 Redis，设置为待领取状态
	key := fmt.Sprintf("%s:%d:%d", taskStatusKeyPrefix, userId, TaskFollowTwitter)
	err = db.RedisCli.Set(key, 3, 24*time.Hour).Err() // 3 表示待领取状态，过期时间为 24 小时
	if err != nil {
		return 0, "Failed to save task status to Redis.", fmt.Errorf("failed to save task status to Redis: %w", err)
	}

	return 1, "follow successfully!", nil
}

func (s *rewardService) GetTaskStatus(ctx context.Context, userId int64, taskId int64) (int, error) {
	// 构造 Redis Key
	key := fmt.Sprintf("%s:%d:%d", taskStatusKeyPrefix, userId, taskId)

	// 从 Redis 获取任务状态
	statusStr, err := db.RedisCli.Get(key).Result()
	if errors.Is(err, redis.Nil) {
		// Key 不存在，返回未完成状态
		status := 1

		// 如果是七天签到任务，检查是否满足条件
		if taskId == TaskSevenDaySignIn {
			consecutive, err := checkSevenDaySignIn(userId)
			if err != nil {
				return 0, fmt.Errorf("failed to check seven-day sign-in: %w", err)
			}
			if consecutive {
				// 更新任务状态为待领取并设置过期时间为 1 天
				err = db.RedisCli.Set(key, 3, 24*time.Hour).Err()
				if err != nil {
					return 0, fmt.Errorf("failed to update task status for seven-day sign-in: %w", err)
				}
				status = 3 // 状态为待领取
			}
		}

		// 如果是关注任务，检查数据库中的关注状态
		if taskId == TaskFollowTwitter {
			user, err := repo.UserRepo.GetNew(userId)
			if err != nil {
				return 0, fmt.Errorf("failed to get user: %w", err)
			}

			// 检查奖励状态是否已领取
			if user.FollowReward == 1 {
				// 更新任务状态为已领取，并保存到 Redis
				err = db.RedisCli.Set(key, 4, 24*time.Hour).Err() // 状态为已领取，过期时间为 1 天
				if err != nil {
					return 0, fmt.Errorf("failed to update task status for follow reward: %w", err)
				}
				status = 4 // 状态为已领取
			}
		}

		return status, nil
	} else if err != nil {
		return 0, fmt.Errorf("failed to get task status: %w", err)
	}

	// 转换状态为整数
	var status int
	fmt.Sscanf(statusStr, "%d", &status)
	return status, nil
}

// ClaimTaskReward 领取任务奖励
func (s *rewardService) ClaimTaskReward(ctx context.Context, userId int64, taskId int64) (string, error) {
	key := fmt.Sprintf("%s:%d:%d", taskStatusKeyPrefix, userId, taskId)
	statusStr, err := db.RedisCli.Get(key).Result()
	if errors.Is(err, redis.Nil) {
		return "Task not found", nil
	} else if err != nil {
		return "Failed to get task status", fmt.Errorf("failed to get task status: %w", err)
	}

	// 转换状态为整数
	var status int
	fmt.Sscanf(statusStr, "%d", &status)

	// 检查状态是否可以领取
	if status != 3 {
		return "Task is not ready for claiming", nil
	}

	// 更新状态为已领取
	err = db.RedisCli.Set(key, 4, time.Hour*24).Err()
	if err != nil {
		return "Failed to update task status", fmt.Errorf("failed to update task status: %w", err)
	}

	// 初始化奖励逻辑
	var changeAmount int
	var reason string
	rewardFields := map[string]interface{}{}

	switch taskId {
	case TaskDailySignIn:
		changeAmount = 1
		reason = "Daily sign-in reward"
	case TaskSevenDaySignIn:
		changeAmount = calculateRandomReward()
		reason = "Seven-day sign-in reward"
	case TaskFollowTwitter:
		changeAmount = 50
		reason = "Twitter follow reward"
		rewardFields = map[string]interface{}{
			"follow_reward": 1,
		}
	default:
		return "Unknown task ID", fmt.Errorf("unknown task ID: %d", taskId)
	}

	// 更新积分、奖励状态并记录日志
	err = updateUserXPoint(ctx, userId, changeAmount, reason, rewardFields)
	if err != nil {
		return "Failed to process task reward", fmt.Errorf("failed to process task reward: %w", err)
	}

	// 动态生成返回消息
	successMessage := fmt.Sprintf("Task reward claimed successfully! +%d Xpoint!", changeAmount)
	return successMessage, nil
}

// FillInviteCode 填写邀请码
func (s *rewardService) FillInviteCode(ctx context.Context, userId int64, inviteCode string) (string, error) {
	// 检查用户是否已填写过邀请码
	user, err := repo.UserRepo.GetNew(userId)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if user.InviterCode != "" {
		return "You have already used an invite code", nil
	}

	// 检查邀请码是否为用户自己
	if inviteCode == user.InviteCode {
		return "", fmt.Errorf("unable to redeem your own invitation code")
	}

	// 检查邀请码是否有效
	inviterUser, err := repo.UserRepo.GetUserByInviteCode(inviteCode)
	if err != nil || inviterUser == nil {
		return "", fmt.Errorf("unable to redeem invitation code")
	}

	// 更新用户填写邀请码的状态
	err = repo.UserRepo.UpdateInviteCodeStatus(userId, inviteCode)
	if err != nil {
		return "", fmt.Errorf("failed to update user invite code status: %w", err)
	}

	// 给邀请人奖励，50 积分
	rewardAmount := 50
	reason := "Invite reward for user using invite code"

	// 更新邀请人的积分
	err = updateUserXPoint(ctx, inviterUser.Id, rewardAmount, reason, nil)
	if err != nil {
		return "", fmt.Errorf("failed to update inviter's points: %w", err)
	}

	// 给被邀请人奖励，10 积分
	inviteeRewardAmount := 10
	inviteeReason := "Reward for using invite code"

	// 更新被邀请人的积分
	err = updateUserXPoint(ctx, userId, inviteeRewardAmount, inviteeReason, nil)
	if err != nil {
		return "", fmt.Errorf("failed to update invitee's points: %w", err)
	}

	// 只返回被邀请人积分的奖励信息
	return fmt.Sprintf("Your invitation code was confirmed! + %d Xpoint!", inviteeRewardAmount), nil
}

// followUser 创建关注目标用户的关系
func followUser(accessToken, userId, targetId string) (bool, error) {
	// 构造 API 请求 URL
	url := fmt.Sprintf("https://api.twitter.com/2/users/%s/following", userId)

	// 构造请求体
	payload := map[string]string{
		"target_user_id": targetId,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to encode payload: %w", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResp struct {
			Title  string `json:"title"`
			Detail string `json:"detail"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errorResp)
		return false, fmt.Errorf("failed to follow user: %s (%s)", errorResp.Title, errorResp.Detail)
	}

	// 成功创建关注关系
	return true, nil
}

// 根据权重计算随机奖励
func calculateRandomReward() int {
	rewardWeights := []struct {
		Amount int
		Weight int
	}{
		{10, 30},
		{20, 25},
		{30, 20},
		{40, 15},
		{50, 10},
	}

	r := rand.Intn(100) + 1
	sum := 0
	for _, rw := range rewardWeights {
		sum += rw.Weight
		if r <= sum {
			return rw.Amount
		}
	}
	return 10 // 默认最小奖励
}

// updateUserXPoint 更新用户积分、奖励状态并记录日志
func updateUserXPoint(ctx context.Context, userID int64, changeAmount int, reason string, rewardFields map[string]interface{}) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新用户积分
	err := tx.Model(&model.User{}).Where("id = ?", userID).Update("xpoint", gorm.Expr("xpoint + ?", changeAmount)).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update xpoint in user table: %w", err)
	}

	// 更新奖励状态或其他字段
	if len(rewardFields) > 0 {
		err = tx.Model(&model.User{}).Where("id = ?", userID).Updates(rewardFields).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update reward fields in user table: %w", err)
		}
	}

	// 记录积分变动日志
	log := model.XPointLog{
		UserID:       uint64(userID),
		ChangeAmount: changeAmount,
		Reason:       reason,
		CreateTime:   time.Now(),
	}
	if err := tx.Create(&log).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert xpoint log: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// checkSevenDaySignIn 检查是否满足连续 7 天签到条件
func checkSevenDaySignIn(userId int64) (bool, error) {
	key := fmt.Sprintf("signin:%d", userId)
	todayOffset := time.Now().UTC().Unix() / 86400
	for i := 0; i < 7; i++ {
		bit, err := db.RedisCli.GetBit(key, todayOffset-int64(i)).Result()
		if err != nil {
			return false, fmt.Errorf("failed to get bit from Redis: %w", err)
		}
		if bit == 0 {
			return false, nil
		}
	}
	return true, nil
}
