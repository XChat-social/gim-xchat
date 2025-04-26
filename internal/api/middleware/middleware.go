package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
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

// Auth 认证中间件
func Auth(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权：缺少Authorization头"})
			c.Abort()
			return
		}

		// 检查Authorization格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权：Authorization格式错误"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析JWT token
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权：" + err.Error()})
			c.Abort()
			return
		}

		// 验证token是否有效
		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
			// 检查token是否在Redis黑名单中（用于登出）
			_, err := rdb.Get(fmt.Sprintf("token:blacklist:%s", tokenString)).Result()
			if err == nil {
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权：token已失效"})
				c.Abort()
				return
			}

			// 将用户ID存储在上下文中
			c.Set("userID", claims.UserID)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权：无效的token"})
			c.Abort()
			return
		}
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
