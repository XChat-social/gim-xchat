package handlers

import (
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

	// 使用新的 context 调用 gRPC
	resp, err := rpc.GetLogicExtClient().CreateChatRoom(grpclib.NewContextFromGin(c), &pb.CreateChatRoomReq{
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

	resp, err := rpc.GetLogicExtClient().GetChatRoom(grpclib.NewContextFromGin(c), &pb.GetChatRoomReq{
		RoomId: roomID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"room_id":          resp.Room.RoomId,
			"name":             resp.Room.Name,
			"avatar_url":       resp.Room.AvatarUrl,
			"introduction":     resp.Room.Introduction,
			"max_member_count": resp.Room.MaxMemberCount,
			"online_count":     resp.Room.OnlineCount,
			"member_count":     resp.Room.MemberCount,
			"extra":            resp.Room.Extra,
			"create_time":      resp.Room.CreateTime,
		},
	})
}

// GetChatRooms 获取聊天室列表
func (h *ChatRoomHandler) GetChatRooms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	resp, err := rpc.GetLogicExtClient().GetChatRooms(grpclib.NewContextFromGin(c), &pb.GetChatRoomsReq{
		PageNumber: int32(page),     // 错误：应该是 page_number
		PageSize:   int32(pageSize), // 正确
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	var roomList []gin.H
	for _, room := range resp.Rooms {
		roomList = append(roomList, gin.H{
			"room_id":          room.RoomId,
			"name":             room.Name,
			"avatar_url":       room.AvatarUrl,
			"introduction":     room.Introduction,
			"max_member_count": room.MaxMemberCount,
			"online_count":     room.OnlineCount,
			"member_count":     room.MemberCount,
			"extra":            room.Extra,
			"create_time":      room.CreateTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"rooms": roomList,
			"total": resp.Total,
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

	// 只接收第二个返回值（error）
	_, err = rpc.GetLogicExtClient().JoinChatRoom(grpclib.NewContextFromGin(c), &pb.JoinChatRoomReq{
		RoomId: roomID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
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

	_, err = rpc.GetLogicExtClient().LeaveChatRoom(grpclib.NewContextFromGin(c), &pb.LeaveChatRoomReq{
		RoomId: roomID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
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

	resp, err := rpc.GetLogicExtClient().GetChatRoomMembers(grpclib.NewContextFromGin(c), &pb.GetChatRoomMembersReq{
		RoomId:     roomID,
		PageSize:   int32(pageSize),
		PageNumber: int32(page),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	var memberList []gin.H
	for _, member := range resp.Members {
		memberList = append(memberList, gin.H{
			"user_id":   member.UserId,
			"join_time": member.JoinTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"members": memberList,
			"total":   resp.Total,
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

	// 调用gRPC服务发送消息
	resp, err := rpc.GetLogicExtClient().SendChatRoomMessage(grpclib.NewContextFromGin(c), &pb.SendChatRoomMessageReq{
		RoomId:   roomID,
		UserId:   userID,
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
