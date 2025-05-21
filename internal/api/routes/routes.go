package routes

import (
	"gim/internal/api/handlers"
	"gim/internal/api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// SetupRoutes 配置所有API路由
func SetupRoutes(r *gin.Engine, db *gorm.DB, rdb *redis.Client) {
	// 静态资源映射（icon访问路径）
	r.Static("/static", "./static")

	// 创建处理器
	userHandler := &handlers.UserHandler{DB: db, RDB: rdb}
	twitterHandler := &handlers.TwitterHandler{DB: db, RDB: rdb}
	taskHandler := &handlers.TaskHandler{DB: db, RDB: rdb}
	tokenHandler := &handlers.TokenHandler{DB: db, RDB: rdb}

	// 添加CORS中间件
	r.Use(middleware.CORS())

	// API路由组
	api := r.Group("/api/v1")
	{
		// 用户相关接口
		api.POST("/sign-in", userHandler.SignIn)
		api.GET("/users/:userId", userHandler.GetUser)
		api.PUT("/users", middleware.Auth(rdb), userHandler.UpdateUser)
		api.GET("/users/search", userHandler.SearchUser)

		// 推特相关接口
		api.GET("/twitter/authorize-url", twitterHandler.GetTwitterAuthorizeURL)
		api.POST("/twitter/sign-in", twitterHandler.TwitterSignIn)
		api.POST("/twitter/follow", middleware.Auth(rdb), twitterHandler.FollowTwitter)
		api.GET("/twitter/users/search", twitterHandler.SearchTwitterUser)

		// 任务相关接口
		api.POST("/daily-sign-in", middleware.Auth(rdb), taskHandler.DailySignIn)
		api.GET("/task/status", middleware.Auth(rdb), taskHandler.GetTaskStatus)
		api.POST("/task/claim", middleware.Auth(rdb), taskHandler.ClaimTaskReward)
		//api.GET("/task/seven-day-signin", middleware.Auth(rdb), taskHandler.CheckSevenDaySignIn)
		api.POST("/task/follow-twitter", middleware.Auth(rdb), taskHandler.FollowTwitter)
		//api.GET("/task/follow-twitter", middleware.Auth(rdb), taskHandler.CheckTwitterFollow)
		//api.PUT("/tasks/status", middleware.Auth(rdb), taskHandler.ModifyTaskStatus)

		// 邀请码相关接口
		api.POST("/invite-code", middleware.Auth(rdb), taskHandler.FillInviteCode)

		// 钱包相关接口
		api.POST("/wallet/sign-in", userHandler.WalletSignIn)

		// Token相关接口
		api.POST("/tokens/create-token", middleware.Auth(rdb), tokenHandler.CreateToken)
		api.GET("/tokens/:tokenAddress", tokenHandler.GetToken)
		api.POST("/tokens/upload-icon", middleware.Auth(rdb), tokenHandler.UploadTokenIcon)
	}

	// 处理Twitter回调
	r.GET("/twitter/signin", twitterHandler.HandleTwitterCallback)
}
