package handlers

import (
	"fmt"
	"gim/internal/api/middleware"
	"gim/internal/api/models"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	DB  *gorm.DB
	RDB *redis.Client
	uh  *UserHandler
}

// CreateToken 创建Token
func (h *TokenHandler) CreateToken(c *gin.Context) {
	var req struct {
		TokenAddress string `json:"token_address" binding:"required"`
		TokenName    string `json:"token_name" binding:"required"`
		TokenSymbol  string `json:"token_symbol" binding:"required"`
		Decimals     int32  `json:"decimals" binding:"required"`
		TotalSupply  string `json:"total_supply" binding:"required"`
		ChainID      int32  `json:"chain_id" binding:"required"`
		IconURL      string `json:"icon_url" binding:"required"`
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

	var user models.User
	userResult := h.DB.First(&user, userID)
	if userResult.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
		return
	}

	// 创建Token
	token := models.Token{
		UserID:         userID,
		TokenAddress:   req.TokenAddress,
		TokenName:      req.TokenName,
		TokenSymbol:    req.TokenSymbol,
		Decimals:       int(req.Decimals),
		TotalSupply:    req.TotalSupply,
		IconUrl:        req.IconURL,
		CreatorAddress: user.WalletAddress,
		ChainID:        int(req.ChainID),
		Status:         1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	result := h.DB.Create(&token)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create token: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Failed to create token",
		"token":   token,
	})
}

// GetToken 获取Token信息
func (h *TokenHandler) GetToken(c *gin.Context) {
	tokenAddress := c.Param("tokenAddress")
	if tokenAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Token地址不能为空"})
		return
	}

	var token models.Token
	result := h.DB.Where("token_address = ?", tokenAddress).First(&token)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Token不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"token":   token,
	})
}

// UploadTokenIcon 上传Token图标
func (h *TokenHandler) UploadTokenIcon(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Failed to get file"})
		return
	}

	// 获取后缀，如 ".png"
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid file type"})
		return
	}

	// 生成唯一文件名
	filename := uuid.New().String() + ext

	// 保存路径（绝对或相对）
	savePath := filepath.Join("static/token-icons", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to save file"})
		return
	}

	// 构造访问 URL（假设你绑定了 /static 路由）
	fileURL := fmt.Sprintf("https://api.xchat.social/api/static/token-icons/%s", filename)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Upload success",
		"data": gin.H{
			"icon_url": fileURL,
		},
	})
}
