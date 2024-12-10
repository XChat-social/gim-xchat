package service

import (
	"context"
	"fmt"
	"gim/internal/business/domain/user/model"
	"gim/pkg/db"
	"github.com/jinzhu/gorm"
	"math/rand"
	"time"
)

type rewardService struct{}

var RewardService = new(rewardService)

// DailySignIn 处理每日签到逻辑
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

	// 更新积分到数据库并记录日志
	err = updateXPoint(ctx, userID, 1, "Daily sign-in")
	if err != nil {
		// 如果数据库操作失败，手动回滚 Redis
		rollbackErr := db.RedisCli.SetBit(key, offset, 0).Err()
		if rollbackErr != nil {
			return 0, "", fmt.Errorf("failed to update xpoint and rollback redis: %w", rollbackErr)
		}
		return 0, "", fmt.Errorf("failed to update xpoint: %w", err)
	}

	return 1, "Sign-in successful", nil
}

// ClaimSevenDayReward 处理连续 7 天签到奖励逻辑
func (*rewardService) ClaimSevenDayReward(ctx context.Context, userID int64) (int, string, error) {
	// Redis Key
	key := fmt.Sprintf("signin:%d", userID)

	// 获取今天的位偏移
	todayOffset := time.Now().UTC().Unix() / 86400

	// 检查最近 7 天的签到状态
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

	// 更新积分到数据库并记录日志
	err = updateXPoint(ctx, userID, rewardAmount, "Seven-day consecutive sign-in")
	if err != nil {
		return 0, "", fmt.Errorf("failed to update xpoint: %w", err)
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

// updateXPoint 更新用户积分并记录日志
func updateXPoint(ctx context.Context, userID int64, changeAmount int, reason string) error {
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
