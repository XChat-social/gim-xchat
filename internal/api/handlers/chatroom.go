package handlers

import (
	"gim/internal/api/models"
	"gim/pkg/grpclib"
	"gim/pkg/protocol/pb"
	"gim/pkg/rpc"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type ChatRoomHandler struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// CreateChatRoom 创建聊天室
func (h *ChatRoomHandler) CreateChatRoom(c *gin.Context) {
	var req struct {
		Name           string `json:"name" binding:"required"`
		AvatarURL      string `json:"avatar_url"`
		Introduction   string `json:"introduction"`
		Extra          string `json:"extra"`
		MaxMemberCount int32  `json:"max_member_count" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数无效"})
		return
	}

	// 获取当前用户ID
	userID, _, err := grpclib.GetCtxData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	// 调用gRPC服务创建聊天室
	resp, err := rpc.GetLogicExtClient().CreateChatRoom(c.Request.Context(), &pb.CreateChatRoomReq{
		Name:           req.Name,
		AvatarUrl:      req.AvatarURL,
		Introduction:   req.Introduction,
		Extra:          req.Extra,
		MaxMemberCount: req.MaxMemberCount,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 创建聊天室记录
	chatRoom := &models.ChatRoom{
		RoomID:         resp.RoomId,
		Name:           req.Name,
		AvatarURL:      req.AvatarURL,
		Introduction:   req.Introduction,
		MaxMemberCount: req.MaxMemberCount,
		Extra:          req.Extra,
		CreatorID:      userID,
	}

	if err := h.DB.Create(chatRoom).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建聊天室失败"})
		return
	}

	// 创建者自动加入聊天室
	member := &models.ChatRoomMember{
		RoomID:   resp.RoomId,
		UserID:   userID,
		JoinTime: time.Now(),
	}

	if err := h.DB.Create(member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "加入聊天室失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"room_id": resp.RoomId},
	})
}

// GetChatRoom 获取聊天室信息
func (h *ChatRoomHandler) GetChatRoom(c *gin.Context) {
	roomID, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的聊天室ID"})
		return
	}

	var chatRoom models.ChatRoom
	if err := h.DB.Where("room_id = ?", roomID).First(&chatRoom).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "聊天室不存在"})
		return
	}

	// 获取在线人数
	onlineCount := 0 // TODO: 从Redis获取在线人数

	// 获取成员总数
	var memberCount int64
	h.DB.Model(&models.ChatRoomMember{}).Where("room_id = ?", roomID).Count(&memberCount)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"room_id":          chatRoom.RoomID,
			"name":             chatRoom.Name,
			"avatar_url":       chatRoom.AvatarURL,
			"introduction":     chatRoom.Introduction,
			"max_member_count": chatRoom.MaxMemberCount,
			"online_count":     onlineCount,
			"member_count":     memberCount,
			"extra":            chatRoom.Extra,
			"creator_id":       chatRoom.CreatorID,
			"create_time":      chatRoom.CreatedAt.Unix(),
		},
	})
}

// GetChatRooms 获取聊天室列表
func (h *ChatRoomHandler) GetChatRooms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var chatRooms []models.ChatRoom
	var total int64

	h.DB.Model(&models.ChatRoom{}).Count(&total)
	h.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&chatRooms)

	var roomList []gin.H
	for _, room := range chatRooms {
		// 获取在线人数
		onlineCount := 0 // TODO: 从Redis获取在线人数

		// 获取成员总数
		var memberCount int64
		h.DB.Model(&models.ChatRoomMember{}).Where("room_id = ?", room.RoomID).Count(&memberCount)

		roomList = append(roomList, gin.H{
			"room_id":          room.RoomID,
			"name":             room.Name,
			"avatar_url":       room.AvatarURL,
			"introduction":     room.Introduction,
			"max_member_count": room.MaxMemberCount,
			"online_count":     onlineCount,
			"member_count":     memberCount,
			"extra":            room.Extra,
			"creator_id":       room.CreatorID,
			"create_time":      room.CreatedAt.Unix(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"rooms": roomList,
			"total": total,
		},
	})
}

// JoinChatRoom 加入聊天室
func (h *ChatRoomHandler) JoinChatRoom(c *gin.Context) {
	roomID, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的聊天室ID"})
		return
	}

	userID, _, err := grpclib.GetCtxData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	// 检查聊天室是否存在
	var chatRoom models.ChatRoom
	if err := h.DB.Where("room_id = ?", roomID).First(&chatRoom).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "聊天室不存在"})
		return
	}

	// 检查是否已经是成员
	var count int64
	h.DB.Model(&models.ChatRoomMember{}).Where("room_id = ? AND user_id = ?", roomID, userID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "已经是聊天室成员"})
		return
	}

	// 检查成员数量是否达到上限
	h.DB.Model(&models.ChatRoomMember{}).Where("room_id = ?", roomID).Count(&count)
	if count >= int64(chatRoom.MaxMemberCount) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "聊天室已满"})
		return
	}

	// 加入聊天室
	member := &models.ChatRoomMember{
		RoomID:   roomID,
		UserID:   userID,
		JoinTime: time.Now(),
	}

	if err := h.DB.Create(member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "加入聊天室失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200})
}

// LeaveChatRoom 离开聊天室
func (h *ChatRoomHandler) LeaveChatRoom(c *gin.Context) {
	roomID, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的聊天室ID"})
		return
	}

	userID, _, err := grpclib.GetCtxData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 400, "message": "获取用户信息失败"})
		return
	}

	// 删除成员记录
	result := h.DB.Where("room_id = ? AND user_id = ?", roomID, userID).Delete(&models.ChatRoomMember{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "离开聊天室失败"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不是聊天室成员"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200})
}

// GetChatRoomMembers 获取聊天室成员列表
func (h *ChatRoomHandler) GetChatRoomMembers(c *gin.Context) {
	roomID, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的聊天室ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var members []models.ChatRoomMember
	var total int64

	h.DB.Model(&models.ChatRoomMember{}).Where("room_id = ?", roomID).Count(&total)
	h.DB.Where("room_id = ?", roomID).Offset((page - 1) * pageSize).Limit(pageSize).Find(&members)

	var memberList []gin.H
	for _, member := range members {
		// TODO: 获取用户信息和在线状态
		memberList = append(memberList, gin.H{
			"user_id":   member.UserID,
			"join_time": member.JoinTime.Unix(),
			"extra":     member.Extra,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"members": memberList,
			"total":   total,
		},
	})
}

// SendMessage 发送聊天室消息
func (h *ChatRoomHandler) SendMessage(c *gin.Context) {
	roomID, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的聊天室ID"})
		return
	}

	// 获取当前用户ID
	userID, _, err := grpclib.GetCtxData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	// 检查用户是否是聊天室成员
	var count int64
	h.DB.Model(&models.ChatRoomMember{}).Where("room_id = ? AND user_id = ?", roomID, userID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "不是聊天室成员"})
		return
	}

	var req struct {
		Content  []byte `json:"content" binding:"required"`
		SendTime int64  `json:"send_time,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数无效"})
		return
	}

	// 如果没有指定发送时间，使用当前时间
	if req.SendTime == 0 {
		req.SendTime = time.Now().Unix()
	}

	// 调用gRPC服务发送消息x
	resp, err := rpc.GetLogicExtClient().SendChatRoomMessage(c.Request.Context(), &pb.SendChatRoomMessageReq{
		RoomId:   roomID,
		UserId:   userID, // 添加发送者ID
		Content:  req.Content,
		SendTime: req.SendTime,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"seq": resp.Seq},
	})
}
