package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gim/internal/api/middleware"
	"gim/internal/api/models"
	"gim/pkg/db"
	"gim/pkg/logger"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// 任务类型常量
const (
	TaskTypeDaily        = 1 // 每日签到
	TaskTypeWeeklySignIn = 2 // 连续签到7天
	TaskTypeFollowX      = 3 // 关注官方X账号
	TaskTypeInvite       = 4 // 邀请好友
	TaskTypeShareX       = 5 // 分享到X
	TaskTypeCompleteKYC  = 6 // 完成KYC认证
)

const (
	TaskDailySignIn     = 1001 // 每日签到
	TaskSevenDaySignIn  = 1002 // 签到七天
	TaskFollowTwitter   = 1003 // 关注推特
	taskStatusKeyPrefix = "task_status"
)

// 任务状态常量
const (
	TaskStatusInProgress = 1 // 进行中
	TaskStatusCompleted  = 2 // 已完成
	TaskStatusClaimed    = 3 // 待领取
	TaskStatusExpired    = 4 // 已领取
)

// DailySignIn 每日签到
func (h *TaskHandler) DailySignIn(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Unauthorized"})
		return
	}

	log.Print("UserID:", userID)

	// 检查用户是否存在
	var user models.User
	result := h.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
		return
	}

	// Redis Key
	key := fmt.Sprintf("signin:%d", userID)
	offset := time.Now().UTC().Unix() / 86400
	signed, err := db.RedisCli.GetBit(key, offset).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query redis"})
		return
	}
	if signed == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Already signed in today"})
		return
	}

	// 设置当天签到状态
	err = db.RedisCli.SetBit(key, offset, 1).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to set redis bit"})
		return
	}

	// 构造任务状态 Redis Key
	taskKey := fmt.Sprintf("%s:%d:%d", taskStatusKeyPrefix, userID, TaskDailySignIn)

	// 将每日签到任务状态设置为待领取 (3)
	err = setWithMidnightExpire(taskKey, TaskStatusClaimed)
	//err = db.RedisCli.Set(taskKey, TaskStatusClaimed, 24*time.Hour).Err() // 设置过期时间为 1 天
	if err != nil {
		// 如果任务状态更新失败，手动回滚 Redis
		rollbackErr := db.RedisCli.SetBit(key, offset, 0).Err()
		if rollbackErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to set task status and rollback redis"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to set task status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Sign-in successful",
	})
}

// FollowTwitter 关注推特并更新状态
func (h *TaskHandler) FollowTwitter(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Unauthorized"})
		return
	}

	var req struct {
		OfficialTwitterID string `json:"official_twitter_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid parameters"})
		return
	}

	// 检查用户是否存在
	var user models.User
	result := h.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
		return
	}

	// 检查奖励状态
	if user.FollowReward == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Reward already claimed"})
		return
	}

	// 从Redis中获取用户的Twitter访问令牌
	redisKey := fmt.Sprintf("twitter:access_token:%s", user.TwitterID)
	accessToken, err := db.RedisCli.Get(redisKey).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to get Twitter access token"})
		return
	}

	// 调用 API 创建关注关系
	isFollowing, err := h.followUser(accessToken, user.TwitterID, req.OfficialTwitterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to follow official Twitter account"})
		return
	}

	// 如果未成功创建关注关系
	if !isFollowing {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Failed to follow the official Twitter account"})
		return
	}

	// 保存状态到 Redis，设置为待领取状态
	key := fmt.Sprintf("%s:%d:%d", taskStatusKeyPrefix, userID, TaskFollowTwitter)
	//err = setWithMidnightExpire(key, TaskStatusClaimed)
	err = db.RedisCli.Set(key, TaskStatusClaimed, 24*time.Hour).Err() // 设置过期时间为 24 小时
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to save task status to Redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Follow successfully!",
	})
}

// GetTaskStatus 获取任务状态
func (h *TaskHandler) GetTaskStatus(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Unauthorized"})
		return
	}

	// 获取任务类型
	taskType := c.Query("task_type")
	if taskType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Task type is required"})
		return
	}

	var taskID int
	switch taskType {
	case "daily_signin":
		taskID = TaskDailySignIn
	case "seven_day_signin":
		taskID = TaskSevenDaySignIn
	case "follow_twitter":
		taskID = TaskFollowTwitter
	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid task type"})
		return
	}

	// 构造任务状态 Redis Key
	key := fmt.Sprintf("%s:%d:%d", taskStatusKeyPrefix, userID, taskID)

	// 从 Redis 获取任务状态
	statusStr, err := db.RedisCli.Get(key).Result()
	logger.Sugar.Infof("Get %s status: %s", key, statusStr)
	if err == redis.Nil {
		// Key 不存在，返回未完成状态
		status := TaskStatusInProgress

		// 如果是七天签到任务，检查是否满足条件
		if taskID == TaskSevenDaySignIn {
			consecutive, err := h.checkSevenDaySignIn(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to check seven-day sign-in"})
				return
			}
			if consecutive {
				// 更新任务状态为待领取并设置过期时间为 1 天
				err = db.RedisCli.Set(key, TaskStatusClaimed, 24*time.Hour).Err()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update task status"})
					return
				}
				status = TaskStatusClaimed // 状态为待领取
			}
		}

		// 如果是关注任务，检查数据库中的关注状态
		if taskID == TaskFollowTwitter {
			var user models.User
			result := h.DB.First(&user, userID)
			if result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
				return
			}

			// 检查奖励状态是否已领取
			if user.FollowReward == 1 {
				// 更新任务状态为已领取，并保存到 Redis
				err = db.RedisCli.Set(key, TaskStatusExpired, 24*time.Hour).Err() // 状态为已领取，过期时间为 1 天
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update task status"})
					return
				}
				status = TaskStatusExpired // 状态为已领取
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "Success",
			"status":  status,
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to get task status"})
		return
	}

	// 转换状态为整数
	var status int
	fmt.Sscanf(statusStr, "%d", &status)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"status":  status,
	})
}

// ClaimTaskReward 领取任务奖励
func (h *TaskHandler) ClaimTaskReward(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Unauthorized"})
		return
	}

	// 获取任务类型
	var req struct {
		TaskType string `json:"task_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid parameters"})
		return
	}

	var taskID int
	var rewardAmount int64
	switch req.TaskType {
	case "daily_signin":
		taskID = TaskDailySignIn
		rewardAmount = 10
	case "seven_day_signin":
		taskID = TaskSevenDaySignIn
		rewardAmount = int64(h.calculateRandomReward())
	case "follow_twitter":
		taskID = TaskFollowTwitter
		rewardAmount = 50
	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid task type"})
		return
	}

	// 构造任务状态 Redis Key
	key := fmt.Sprintf("%s:%d:%d", taskStatusKeyPrefix, userID, taskID)

	// 获取任务状态
	statusStr, err := db.RedisCli.Get(key).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Task not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to get task status"})
		return
	}

	// 转换状态为整数
	var status int
	fmt.Sscanf(statusStr, "%d", &status)

	// 检查状态是否可以领取
	if status != TaskStatusClaimed {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Task is not ready for claiming"})
		return
	}

	// 开始事务
	tx := h.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to begin transaction"})
		return
	}

	// 更新用户积分
	var updateFields map[string]interface{} = map[string]interface{}{
		"xpoint": gorm.Expr("xpoint + ?", rewardAmount),
	}

	// 如果是关注推特任务，更新关注奖励状态
	if taskID == TaskFollowTwitter {
		updateFields["follow_reward"] = 1
	}

	result := tx.Model(&models.User{}).Where("id = ?", userID).Updates(updateFields)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update user points"})
		return
	}

	// 记录积分变动日志
	var reason string
	switch taskID {
	case TaskDailySignIn:
		reason = "Daily sign-in reward"
	case TaskSevenDaySignIn:
		reason = "Seven-day sign-in reward"
	case TaskFollowTwitter:
		reason = "Twitter follow reward"
	}

	xpointLog := models.XPointLog{
		UserID:       uint64(userID),
		ChangeAmount: int(rewardAmount),
		Reason:       reason,
		CreateTime:   time.Now(),
	}

	result = tx.Create(&xpointLog)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create xpoint log"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Transaction failed"})
		return
	}

	// 更新任务状态为已领取
	err = db.RedisCli.Set(key, TaskStatusExpired, 0).Err()
	if err != nil {
		// 这里不回滚事务，因为奖励已经发放，只是状态更新失败
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "Reward claimed successfully, but status update failed",
			"reward":  rewardAmount,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": fmt.Sprintf("Task reward claimed successfully! +%d Xpoint!", rewardAmount),
		"reward":  rewardAmount,
	})
}

// FillInviteCode 填写邀请码
func (h *TaskHandler) FillInviteCode(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Unauthorized"})
		return
	}

	var req struct {
		InviteCode string `json:"invite_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid parameters"})
		return
	}

	// 检查用户是否存在
	var user models.User
	result := h.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
		return
	}

	// 检查用户是否已填写过邀请码
	if user.InviterCode != "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "You have already used an invite code"})
		return
	}

	// 检查邀请码是否为用户自己
	if req.InviteCode == user.InviteCode {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Unable to redeem your own invitation code"})
		return
	}

	// 检查邀请码是否有效
	var inviterUser models.User
	result = h.DB.Where("invite_code = ?", req.InviteCode).First(&inviterUser)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Unable to redeem invitation code"})
		return
	}

	// 开始事务
	tx := h.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to begin transaction"})
		return
	}

	// 更新用户填写邀请码的状态
	result = tx.Model(&models.User{}).Where("id = ?", userID).Update("inviter_code", req.InviteCode)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update user invite code status"})
		return
	}

	// 给邀请人奖励，50 积分
	inviterRewardAmount := 50
	inviterReason := "Invite reward for user using invite code"

	// 更新邀请人的积分
	result = tx.Model(&models.User{}).Where("id = ?", inviterUser.ID).UpdateColumn("xpoint", gorm.Expr("xpoint + ?", inviterRewardAmount))
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update inviter's points"})
		return
	}

	// 记录邀请人积分变动日志
	inviterXpointLog := models.XPointLog{
		UserID:       uint64(inviterUser.ID),
		ChangeAmount: inviterRewardAmount,
		Reason:       inviterReason,
		CreateTime:   time.Now(),
	}

	result = tx.Create(&inviterXpointLog)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create inviter xpoint log"})
		return
	}

	// 给被邀请人奖励，10 积分
	inviteeRewardAmount := 10
	inviteeReason := "Reward for using invite code"

	// 更新被邀请人的积分
	result = tx.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("xpoint", gorm.Expr("xpoint + ?", inviteeRewardAmount))
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update invitee's points"})
		return
	}

	// 记录被邀请人积分变动日志
	inviteeXpointLog := models.XPointLog{
		UserID:       uint64(userID),
		ChangeAmount: inviteeRewardAmount,
		Reason:       inviteeReason,
		CreateTime:   time.Now(),
	}

	result = tx.Create(&inviteeXpointLog)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create invitee xpoint log"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Transaction failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": fmt.Sprintf("Your invitation code was confirmed! + %d Xpoint!", inviteeRewardAmount),
	})
}

// followUser 创建关注目标用户的关系
func (h *TaskHandler) followUser(accessToken, userId, targetId string) (bool, error) {
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

// calculateRandomReward 根据权重计算随机奖励
func (h *TaskHandler) calculateRandomReward() int {
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

// checkSevenDaySignIn 检查是否满足连续 7 天签到条件
func (h *TaskHandler) checkSevenDaySignIn(userId int64) (bool, error) {
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

// setWithMidnightExpire 设置24点过期的键值对
func setWithMidnightExpire(key string, value uint64) error {
	now := time.Now()
	loc := now.Location()

	// 计算次日00:00
	tomorrow := now.AddDate(0, 0, 1)
	endOfDay := time.Date(
		tomorrow.Year(),
		tomorrow.Month(),
		tomorrow.Day(),
		0, 0, 0, 0, loc,
	)

	// 计算时间差
	expireIn := endOfDay.Sub(now)

	// 处理时间穿越情况
	if expireIn < 0 {
		expireIn = 0 // 立即过期
	}

	return db.RedisCli.Set(key, value, expireIn).Err()
}
