package handlers

import (
	"fmt"
	"gim/internal/api/middleware"
	"gim/internal/api/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis"
	"github.com/storyicon/sigverify"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// SignIn 处理登录请求
func (h *UserHandler) SignIn(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		Code        string `json:"code" binding:"required"`
		DeviceID    int64  `json:"device_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid parameters"})
		return
	}

	isNew := false
	userId := int64(12345)
	token := "mock_token_" + time.Now().Format("20060102150405")

	c.JSON(http.StatusOK, gin.H{
		"is_new":  isNew,
		"user_id": userId,
		"token":   token,
	})
}

// GetUser 获取用户信息
func (h *UserHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid user ID format"})
		return
	}

	var user models.User
	result := h.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"user":    user,
	})
}

// UpdateUser 更新用户信息
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req struct {
		Nickname    string `json:"nickname"`
		Sex         int32  `json:"sex"`
		AvatarURL   string `json:"avatar_url"`
		Extra       string `json:"extra"`
		Description string `json:"description"`
		Email       string `json:"email"`
		Birthday    string `json:"birthday"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid parameters"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Unauthorized"})
		return
	}

	updates := map[string]interface{}{
		"update_time": time.Now(),
	}

	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Sex != 0 {
		updates["sex"] = req.Sex
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}
	if req.Extra != "" {
		updates["extra"] = req.Extra
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Birthday != "" {
		updates["birthday"] = req.Birthday
	}

	if len(updates) <= 1 {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "No fields to update"})
		return
	}

	result := h.DB.Model(&models.User{}).Where("user_id = ?", userID).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Update failed: " + result.Error.Error()})
		return
	}

	var user models.User
	h.DB.First(&user, userID)

	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"message":   "Update successful",
		"user_info": user,
	})
}

// SearchUser 搜索用户
func (h *UserHandler) SearchUser(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Search keyword cannot be empty"})
		return
	}

	var users []models.User
	result := h.DB.Where("nickname LIKE ? OR twitter_username LIKE ?", "%"+key+"%", "%"+key+"%").Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Search failed: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// WalletSignIn 钱包登录
func (h *UserHandler) WalletSignIn(c *gin.Context) {
	var req struct {
		WalletAddress string `json:"wallet_address" binding:"required"`
		Signature     string `json:"signature" binding:"required"`
		Message       string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid parameters"})
		return
	}

	if !common.IsHexAddress(req.WalletAddress) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid wallet address format"})
		return
	}

	verified, err := verifySignature(req.Message, req.Signature, req.WalletAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Signature verification failed: " + err.Error()})
		return
	}
	if !verified {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Signature verification failed"})
		return
	}

	var user models.User
	result := h.DB.Where("wallet_address = ?", req.WalletAddress).First(&user)
	isNew := result.Error != nil

	if isNew {
		user = models.User{
			WalletAddress: req.WalletAddress,
			InviteCode:    randomString(8),
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		}
		result = h.DB.Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create user: " + result.Error.Error()})
			return
		}
	}

	token, err := middleware.GenerateToken(int64(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to generate token: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"message":   "Success",
		"is_new":    isNew,
		"user_id":   user.ID,
		"token":     token,
		"user_info": user,
	})
}

// verifySignature 验证签名
func verifySignature(address, message, signature string) (bool, error) {
	valid, err := sigverify.VerifyEllipticCurveHexSignatureEx(
		common.HexToAddress(address),
		[]byte(message),
		signature,
	)
	if err != nil {
		return false, fmt.Errorf("signature verification failed: %v", err)
	}
	return valid, nil
}
