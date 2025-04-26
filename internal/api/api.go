package api

import (
	"gim/internal/api/middleware"
	"gim/internal/api/routes"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// APIService API服务结构体
type APIService struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// NewAPIService 创建新的API服务
func NewAPIService(db *gorm.DB, rdb *redis.Client) *APIService {
	return &APIService{
		DB:  db,
		RDB: rdb,
	}
}

// InitRouter 初始化路由
func (s *APIService) InitRouter() *gin.Engine {
	r := gin.Default()

	// 添加CORS中间件，使用middleware包中的CORS函数
	r.Use(middleware.CORS())

	// 注册路由
	routes.SetupRoutes(r, s.DB, s.RDB)

	return r
}
