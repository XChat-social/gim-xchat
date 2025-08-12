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
		api.GET("/users/xpoint/:userId", userHandler.GetUserXpoint)
		api.PUT("/users", middleware.Auth(rdb), userHandler.UpdateUser)
		api.GET("/users/search", userHandler.SearchUser)
		api.POST("/users/upload-avatar", middleware.Auth(rdb), tokenHandler.UploadTokenIcon)

		// 推特相关接口
		api.GET("/twitter/authorize-url", twitterHandler.GetTwitterAuthorizeURL)
		api.POST("/twitter/sign-in", twitterHandler.TwitterSignIn)
		api.POST("/twitter/follow", middleware.Auth(rdb), twitterHandler.FollowTwitter)
		api.GET("/twitter/users/search", twitterHandler.SearchTwitterUser)

		// 任务相关接口
		api.POST("/daily-sign-in", middleware.Auth(rdb), taskHandler.DailySignIn)
		api.GET("/task/status", middleware.Auth(rdb), taskHandler.GetTaskStatus)
		api.GET("/task/dailySum", middleware.Auth(rdb), taskHandler.GetDailySum)
		api.GET("/task/deleteKey", middleware.Auth(rdb), taskHandler.DeleteKey)
		api.POST("/task/claim", middleware.Auth(rdb), taskHandler.ClaimTaskReward)
		//api.GET("/task/seven-day-signin", middleware.Auth(rdb), taskHandler.CheckSevenDaySignIn)
		api.POST("/task/follow-twitter", middleware.Auth(rdb), taskHandler.FollowTwitter)
		//api.GET("/task/follow-twitter", middleware.Auth(rdb), taskHandler.CheckTwitterFollow)
		//api.PUT("/tasks/status", middleware.Auth(rdb), taskHandler.ModifyTaskStatus)

		// 邀请码相关接口
		api.POST("/invite-code", middleware.Auth(rdb), taskHandler.FillInviteCode)
		// 内测码相关接口
		api.POST("/gen-beta-code", middleware.Auth(rdb), taskHandler.GenBetaCode)
		api.POST("/redeem-beta-code", middleware.Auth(rdb), taskHandler.RedeemBetaCode)
		api.POST("/checkUserAuth", middleware.Auth(rdb), taskHandler.CheckUserAuth)

		// 钱包相关接口
		api.POST("/wallet/sign-in", userHandler.WalletSignIn)

		// Token相关接口
		api.POST("/tokens/create-token", middleware.Auth(rdb), tokenHandler.CreateToken)
		api.GET("/tokens/getTokenByUserId", middleware.Auth(rdb), tokenHandler.GetTokenByUserID)
		api.GET("/tokens/:tokenAddress", tokenHandler.GetToken)
		api.POST("/tokens/upload-icon", middleware.Auth(rdb), tokenHandler.UploadTokenIcon)
		api.GET("/tokens/getTokenHolders", middleware.Auth(rdb), tokenHandler.GetTokenHolders)

		// Token交易相关接口
		api.POST("/tokens/buy", middleware.Auth(rdb), tokenHandler.BuyToken)
		api.POST("/tokens/sell", middleware.Auth(rdb), tokenHandler.SellToken)
		api.GET("/tokens/holdings", middleware.Auth(rdb), tokenHandler.GetUserTokenHoldings)
	}

	// 处理Twitter回调
	r.GET("/twitter/signin", twitterHandler.HandleTwitterCallback)

	// 创建聊天室处理器
	chatRoomHandler := &handlers.ChatRoomHandler{DB: db, RDB: rdb}
	{
		// 聊天室相关接口
		chatroom := api.Group("/chatrooms", middleware.Auth(rdb))
		{
			chatroom.POST("", chatRoomHandler.CreateChatRoom)                                       // 创建聊天室
			chatroom.GET("", chatRoomHandler.GetChatRooms)                                          // 获取聊天室列表
			chatroom.GET("/user", chatRoomHandler.GetUserChatRooms)                                 // 获取用户加入的聊天室列表
			chatroom.GET("/:roomId", chatRoomHandler.GetChatRoom)                                   // 获取聊天室信息
			chatroom.POST("/:roomId/join", chatRoomHandler.JoinChatRoom)                            // 加入聊天室
			chatroom.POST("/:roomId/leave", chatRoomHandler.LeaveChatRoom)                          // 离开聊天室
			chatroom.GET("/:roomId/members", chatRoomHandler.GetChatRoomMembers)                    // 获取成员列表
			chatroom.POST("/:roomId/message", chatRoomHandler.SendMessage)                          // 发送聊天室消息
			chatroom.POST("/history", chatRoomHandler.GetChatRoomMessages)                          // 获取聊天室消息历史
			chatroom.GET("/:roomId/checkPermissions", chatRoomHandler.CheckPermissionsByUserId)     // 检查用户是否是聊天室成员
			chatroom.POST("/upload-icon", chatRoomHandler.UploadChatRoomIcon)                       // 上传聊天室图标
			chatroom.POST("/thumb", chatRoomHandler.ThumbMessage)                                   // like
			chatroom.GET("/getChatRoomsByParticipants", chatRoomHandler.GetChatRoomsByParticipants) // 获取聊天室列表根据总数排序
		}
	}
}
