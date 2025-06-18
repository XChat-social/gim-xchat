package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v4"
)

// JWTSecret JWT密钥
const JWTSecret = "XChatSecret_2024_a7b9c3d5e8f2g4h6j8k0m1n3p5q7r9t2v4w6y8z0"

// UserClaims 用户JWT声明
type UserClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// CORS 处理跨域请求
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Auth 认证中间件（完成版）
func Auth(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Missing or invalid token"})
			c.Abort()
			return
		}

		// 解析 token 字符串
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析 JWT
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid token"})
			c.Abort()
			return
		}

		// 提取用户ID
		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid token claims"})
			c.Abort()
			return
		}

		// 将认证信息放入上下文
		c.Set("user_id", claims.UserID) // 改为 user_id
		c.Set("device_id", 1)           // 设置一个默认的 device_id
		c.Set("token", tokenString)     // 保存原始 token

		c.Next()
	}
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	return userID.(int64), true
}

// GenerateToken 生成JWT token
func GenerateToken(userID int64) (string, error) {
	claims := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

// InvalidateToken 使token失效（登出）
func InvalidateToken(rdb *redis.Client, tokenString string) error {
	// 解析token以获取过期时间
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return fmt.Errorf("无法解析token声明")
	}

	// 计算token剩余有效期
	expiresAt := claims.ExpiresAt.Time
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil // token已过期，无需加入黑名单
	}

	// 将token加入Redis黑名单，过期时间与token相同
	return rdb.Set(
		fmt.Sprintf("token:blacklist:%s", tokenString),
		"1",
		ttl,
	).Err()
}
