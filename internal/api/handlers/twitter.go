package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gim/internal/api/middleware"
	"gim/internal/api/models"
	"gim/internal/business/domain/user/model"
	"gim/pkg/db"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type TwitterHandler struct {
	DB  *gorm.DB
	RDB *redis.Client
}

const (
	twitterAuthorizeURL = "https://twitter.com/i/oauth2/authorize"
	twitterTokenURL     = "https://api.twitter.com/2/oauth2/token"
	twitterUserInfoURL  = "https://api.twitter.com/2/users/me"
	clientID            = "YS01bVJhaXdIdEN4X3N5cjNlQzM6MTpjaQ"                 // Twitter Client ID
	clientSecret        = "JbYIsa77FIbRVc_ZY4238KaPV3Y-K-G5ylbOfuDHgVr8WtgEvO" // Twitter Client Secret
	redirectURI         = "https://api.xchat.social/api/twitter/signin"
	officialTwitterID   = "1837782128660017152" // XChat 官方 Twitter ID
	stateTTL            = 5 * time.Minute       // Redis state & verifier TTL
)

// GetTwitterAuthorizeURL 获取推特授权URL
func (h *TwitterHandler) GetTwitterAuthorizeURL(c *gin.Context) {
	walletAddress := c.Query("wallet_address")
	if walletAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Wallet address cannot be empty"})
		return
	}

	codeVerifier := generateCodeVerifier()
	codeChallenge := generateCodeChallenge(codeVerifier)
	state := fmt.Sprintf("%s:%s", generateRandomState(), walletAddress)

	err := h.RDB.Set(fmt.Sprintf("twitter:state:%s", state), codeVerifier, stateTTL).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to save state"})
		return
	}

	authorizeURL := fmt.Sprintf(
		"%s?response_type=code&client_id=%s&redirect_uri=%s&scope=tweet.read users.read follows.read follows.write&state=%s&code_challenge=%s&code_challenge_method=S256",
		twitterAuthorizeURL, clientID, url.QueryEscape(redirectURI), state, codeChallenge,
	)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"url":     authorizeURL,
	})
}

// HandleTwitterCallback 处理Twitter回调
func (h *TwitterHandler) HandleTwitterCallback(c *gin.Context) {
	authorizationCode := c.Query("code")
	if authorizationCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Authorization code is required"})
		return
	}

	state := c.Query("state")
	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "State is required"})
		return
	}

	stateParts := strings.Split(state, ":")
	if len(stateParts) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid state format"})
		return
	}
	walletAddress := stateParts[1]

	codeVerifier, err := h.RDB.Get(fmt.Sprintf("twitter:state:%s", state)).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid or expired state"})
		return
	}

	tempToken := randomString(32)

	tokenData := fmt.Sprintf("%s:%s:%s", authorizationCode, codeVerifier, walletAddress)
	err = h.RDB.Set(fmt.Sprintf("twitter:temp_token:%s", tempToken), tokenData, 5*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to save temporary token"})
		return
	}

	redirectURL := fmt.Sprintf(
		"http://localhost:3000?redirect=redirectx&temp_token=%s",
		url.QueryEscape(tempToken),
	)

	c.Redirect(http.StatusFound, redirectURL)
}

// TwitterSignIn 推特登录
func (h *TwitterHandler) TwitterSignIn(c *gin.Context) {
	var req struct {
		TempToken string `json:"temp_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid parameters"})
		return
	}

	tokenData, err := h.RDB.Get(fmt.Sprintf("twitter:temp_token:%s", req.TempToken)).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid or expired temporary token"})
		return
	}

	parts := strings.Split(tokenData, ":")
	if len(parts) != 3 {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Invalid temporary token data format"})
		return
	}

	authorizationCode := parts[0]
	codeVerifier := parts[1]
	walletAddress := parts[2]

	accessToken, err := exchangeCodeForToken(authorizationCode, codeVerifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to get Twitter access token: " + err.Error()})
		return
	}

	twitterUser, err := getTwitterUserInfo(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to get Twitter user info: " + err.Error()})
		return
	}

	// 将accessToken存入Redis，设置过期时间为7天
	redisKey := fmt.Sprintf("twitter:access_token:%s", twitterUser.ID)
	err = h.RDB.Set(redisKey, accessToken, 7*24*time.Hour).Err()
	if err != nil {
		// 记录错误但不中断流程
		fmt.Printf("Failed to save Twitter access token to Redis: %v\n", err)
	}

	var user models.User
	result := h.DB.Where("twitter_id = ?", twitterUser.ID).First(&user)
	isNew := result.Error != nil

	inviteCode := randomString(8)

	if isNew {
		user = models.User{
			Nickname:        twitterUser.Name,
			AvatarURL:       twitterUser.Avatar,
			TwitterID:       twitterUser.ID,
			TwitterUsername: twitterUser.Username,
			InviteCode:      inviteCode,
			WalletAddress:   walletAddress,
			CreateTime:      time.Now(),
			UpdateTime:      time.Now(),
		}
		result = h.DB.Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create user: " + result.Error.Error()})
			return
		}
	} else {
		updates := map[string]interface{}{
			"nickname":         twitterUser.Name,
			"twitter_username": twitterUser.Username,
			"wallet_address":   walletAddress,
			"update_time":      time.Now(),
		}
		result = h.DB.Model(&user).Updates(updates)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update user info: " + result.Error.Error()})
			return
		}
	}

	token, err := middleware.GenerateToken(int64(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to generate token: " + err.Error()})
		return
	}

	h.RDB.Del(fmt.Sprintf("twitter:temp_token:%s", req.TempToken))

	c.JSON(http.StatusOK, gin.H{
		"code":        200,
		"message":     "Success",
		"is_new":      isNew,
		"user_id":     user.ID,
		"token":       token,
		"user_info":   user,
		"err_message": "",
	})
}

// FollowTwitter 关注Twitter
func (h *TwitterHandler) FollowTwitter(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Unauthorized"})
		return
	}

	// 查询用户信息
	var user models.User
	result := h.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
		return
	}

	// 检查用户是否有Twitter账号
	if user.TwitterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "User has no Twitter account"})
		return
	}

	// 获取Twitter访问令牌
	accessToken, err := h.getTwitterAccessToken(user.TwitterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to get Twitter access token: " + err.Error()})
		return
	}

	// 调用 API 创建关注关系
	_, errFollowing := h.followUser(accessToken, user.TwitterID, officialTwitterID)
	if errFollowing != nil {
		// 如果包含具体 HTTP 状态码，如 429 或 401
		if strings.Contains(errFollowing.Error(), "HTTP 429") {
			c.JSON(http.StatusTooManyRequests, gin.H{"code": 429, "message": errFollowing.Error()})
		} else if strings.Contains(errFollowing.Error(), "HTTP 401") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": errFollowing.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": errFollowing.Error()})
		}
		return
	}

	// 保存状态到 Redis，设置为待领取状态
	key := fmt.Sprintf("%s:%d:%d", taskStatusKeyPrefix, userID, TaskFollowTwitter)
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

// 获取Twitter访问令牌
func (h *TwitterHandler) getTwitterAccessToken(twitterID string) (string, error) {
	// 从Redis中获取用户的Twitter访问令牌
	redisKey := fmt.Sprintf("twitter:access_token:%s", twitterID)
	accessToken, err := h.RDB.Get(redisKey).Result()

	// 如果Redis中有有效的令牌，直接返回
	if err == nil && accessToken != "" {
		return accessToken, nil
	}

	// 如果Redis中没有，则使用应用级别的访问令牌
	data := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	resp, err := http.PostForm(twitterTokenURL, data)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get access token: status code %d", resp.StatusCode)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// 将新获取的令牌存入Redis，设置适当的过期时间
	// 通常应用级别的令牌有效期较长，这里设置为1小时
	h.RDB.Set(redisKey, tokenResp.AccessToken, 1*time.Hour)

	return tokenResp.AccessToken, nil
}

// 关注Twitter用户
// followUser 创建关注目标用户的关系
func (h *TwitterHandler) followUser(accessToken, userId, targetId string) (bool, error) {
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

// SearchTwitterUser 搜索Twitter用户
func (h *TwitterHandler) SearchTwitterUser(c *gin.Context) {
	twitterUsername := c.Query("twitter_username")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	pageNumStr := c.DefaultQuery("page_num", "1")

	pageSize, _ := strconv.Atoi(pageSizeStr)
	pageNum, _ := strconv.Atoi(pageNumStr)

	var users []model.User
	offset := (pageNum - 1) * pageSize
	result := h.DB.Where("twitter_username LIKE ?", "%"+twitterUsername+"%").
		Offset(offset).Limit(pageSize).Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Search failed: " + result.Error.Error()})
		return
	}

	var total int64
	h.DB.Model(&model.User{}).Where("twitter_username LIKE ?", "%"+twitterUsername+"%").Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"users":   users,
		"total":   total,
	})
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
