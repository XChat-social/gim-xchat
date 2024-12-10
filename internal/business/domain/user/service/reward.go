package service

import (
	"context"
	"fmt"
	"gim/internal/business/domain/user/model"
	"gim/internal/business/domain/user/repo"
	"gim/pkg/db"
	"math/rand"
	"time"
)

type rewardService struct{}

var RewardService = new(rewardService)

func (*rewardService) DailySignIn(ctx context.Context, userID int64) (int, string, error) {
	// Redis Key
	key := fmt.Sprintf("signin:%d", userID)

	// 获取当天的位偏移
	offset := time.Now().UTC().Unix() / 86400 // 当前时间的天偏移量

	// 检查当天是否已签到
	signed, err := db.RedisCli.GetBit(key, offset).Result()
	if err != nil {
		return 0, "", fmt.Errorf("failed to query redis: %w", err)
	}
	if signed == 1 {
		return 0, "Already signed in today", nil
	}

	// 设置签到状态
	err = db.RedisCli.SetBit(key, offset, 1).Err()
	if err != nil {
		return 0, "", fmt.Errorf("failed to set redis bit: %w", err)
	}

	// 奖励 +1 XPoint
	return 1, "Sign-in successful", nil
}

func (*rewardService) ClaimSevenDayReward(ctx context.Context, userID int64) (int, string, error) {
	// Redis Key
	key := fmt.Sprintf("signin:%d", userID)

	// 获取今天的位偏移
	todayOffset := time.Now().UTC().Unix() / 86400

	// 检查最近7天的签到状态
	consecutive := true
	for i := 0; i < 7; i++ {
		bit, err := db.RedisCli.GetBit(key, todayOffset-int64(i)).Result()
		if err != nil || bit == 0 {
			consecutive = false
			break
		}
	}
	if !consecutive {
		return 0, "Not enough consecutive signin days", nil
	}

	// 检查是否已领取奖励
	rewardKey := fmt.Sprintf("reward_claimed:%d", userID)
	claimed, err := db.RedisCli.Get(rewardKey).Result()
	if err == nil && claimed == "true" {
		return 0, "Reward already claimed", nil
	}

	// 计算随机奖励
	rewardAmount := calculateRandomReward()

	// 持久化奖励信息到数据库
	_, err = repo.UserRewardDao.Add(model.UserReward{
		UserID:       userID,
		RewardType:   2,
		RewardAmount: int32(rewardAmount),
		RewardDate:   time.Now(),
	})
	if err != nil {
		return 0, "", fmt.Errorf("failed to add reward: %w", err)
	}

	// 标记奖励已领取
	db.RedisCli.Set(rewardKey, "true", 7*24*time.Hour) // 设置 7 天的过期时间

	return rewardAmount, "Reward claimed successfully", nil
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
