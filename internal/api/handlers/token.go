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
	"strconv"
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
	// 开始数据库事务
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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
	var holding models.TokenHolding
	parseFloat, _ := strconv.ParseFloat(req.TotalSupply, 64)
	holding = models.TokenHolding{
		UserID:       userID,
		TokenAddress: req.TokenAddress,
		Amount:       parseFloat,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := tx.Create(&holding).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create holding record"})
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
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Token address cannot be empty"})
		return
	}

	var token models.Token
	result := h.DB.Where("token_address = ?", tokenAddress).First(&token)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Token not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"token":   token,
	})
}

// GetTokenByUserID 根据用户ID获取Token信息
func (h *TokenHandler) GetTokenByUserID(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "User ID cannot be empty"})
		return
	}

	var tokens []models.Token
	result := h.DB.Where("user_id = ?", userID).Find(&tokens)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Database error"})
		return
	}

	if len(tokens) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "No tokens found for this user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"tokens":  tokens,
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

// BuyToken 购买Token
func (h *TokenHandler) BuyToken(c *gin.Context) {
	var req struct {
		TokenAddress string  `json:"token_address" binding:"required"`
		Amount       float64 `json:"amount" binding:"required,gt=0"`
		Price        float64 `json:"price" binding:"required,gt=0"`
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

	// 检查Token是否存在
	var token models.Token
	if err := h.DB.Where("token_address = ?", req.TokenAddress).First(&token).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Token not found"})
		return
	}

	// 开始数据库事务
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建交易记录
	transaction := models.TokenTransaction{
		UserID:          userID,
		TokenAddress:    req.TokenAddress,
		TransactionType: "buy",
		Amount:          req.Amount,
		Price:           req.Price,
		TotalValue:      req.Amount * req.Price,
		Status:          "completed",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create transaction record"})
		return
	}

	// 更新或创建Token持有记录
	var holding models.TokenHolding
	result := tx.Where("user_id = ? AND token_address = ?", userID, req.TokenAddress).First(&holding)
	if result.Error == gorm.ErrRecordNotFound {
		// 创建新的持有记录
		holding = models.TokenHolding{
			UserID:       userID,
			TokenAddress: req.TokenAddress,
			Amount:       req.Amount,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := tx.Create(&holding).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create holding record"})
			return
		}
	} else if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query holding record"})
		return
	} else {
		// 更新现有持有记录
		holding.Amount += req.Amount
		holding.UpdatedAt = time.Now()
		if err := tx.Save(&holding).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update holding record"})
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Transaction commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Purchase successful",
		"data": gin.H{
			"transaction_id": transaction.ID,
			"amount":         req.Amount,
			"total_value":    transaction.TotalValue,
			"new_balance":    holding.Amount,
		},
	})
}

// SellToken 出售Token
func (h *TokenHandler) SellToken(c *gin.Context) {
	var req struct {
		TokenAddress string  `json:"token_address" binding:"required"`
		Amount       float64 `json:"amount" binding:"required,gt=0"`
		Price        float64 `json:"price" binding:"required,gt=0"`
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

	// 检查Token是否存在
	var token models.Token
	if err := h.DB.Where("token_address = ?", req.TokenAddress).First(&token).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Token not found"})
		return
	}

	// 检查用户是否有足够的Token余额
	var holding models.TokenHolding
	if err := h.DB.Where("user_id = ? AND token_address = ?", userID, req.TokenAddress).First(&holding).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Token not held"})
		return
	}

	if holding.Amount < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Insufficient token balance"})
		return
	}

	// 开始数据库事务
	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建交易记录
	transaction := models.TokenTransaction{
		UserID:          userID,
		TokenAddress:    req.TokenAddress,
		TransactionType: "sell",
		Amount:          req.Amount,
		Price:           req.Price,
		TotalValue:      req.Amount * req.Price,
		Status:          "completed",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create transaction record"})
		return
	}

	// 更新Token持有记录
	holding.Amount -= req.Amount
	holding.UpdatedAt = time.Now()
	if err := tx.Save(&holding).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update holding record"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Transaction commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Sale successful",
		"data": gin.H{
			"transaction_id": transaction.ID,
			"amount":         req.Amount,
			"total_value":    transaction.TotalValue,
			"new_balance":    holding.Amount,
		},
	})
}

// GetUserTokenHoldings 获取用户Token持有记录
func (h *TokenHandler) GetUserTokenHoldings(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Unauthorized"})
		return
	}

	var holdings []models.TokenHolding
	if err := h.DB.Where("user_id = ?", userID).Find(&holdings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query holding records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Query successful",
		"data":    holdings,
	})
}
