package handlers

import (
	"gim/internal/api/models"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	DB  *gorm.DB
	RDB *redis.Client
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
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 从请求头获取用户ID
	// 实际项目中应该从JWT token中提取
	userID := int64(12345)

	// 创建Token
	token := models.Token{
		UserID:       userID,
		TokenAddress: req.TokenAddress,
		TokenName:    req.TokenName,
		TokenSymbol:  req.TokenSymbol,
		Decimals:     req.Decimals,
		TotalSupply:  req.TotalSupply,
		CreatorAddr:  "", // 实际项目中应该设置为用户的钱包地址
		ChainID:      req.ChainID,
		Status:       1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result := h.DB.Create(&token)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
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
