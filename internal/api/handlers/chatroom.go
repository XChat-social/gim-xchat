package handlers

import (
	"fmt"
	"gim/internal/api/middleware"
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

// GetChatRoomMessages 获取聊天室消息历史
func (h *ChatRoomHandler) GetChatRoomMessages(c *gin.Context) {
	var req struct {
		RoomID   int64 `json:"roomId" binding:"required"`
		PageNo   int32 `json:"pageNo" binding:"required"`
		PageSize int32 `json:"pageSize" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数无效"})
		return
	}

	// 调用 gRPC 服务
	resp, err := rpc.GetLogicExtClient().GetChatRoomMessages(grpclib.NewContextFromGin(c), &pb.GetChatRoomMessagesReq{
		RoomId:     req.RoomID,
		PageNumber: req.PageNo,
		PageSize:   req.PageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 格式化消息数据
	messages := make([]gin.H, 0, len(resp.Messages))
	for _, msg := range resp.Messages {
		messages = append(messages, gin.H{
			"id":        msg.Id,
			"room_id":   msg.RoomId,
			"user_id":   msg.UserId,
			"content":   string(msg.Content),
			"seq":       msg.Seq,
			"send_time": msg.SendTime,
			"status":    msg.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"messages": messages,
			"total":    resp.Total,
			"page_no":  resp.PageNo,
			"pages":    resp.Pages,
		},
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
	var req pb.GetChatRoomsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}

	resp, err := rpc.GetLogicExtClient().GetChatRooms(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//

	c.JSON(http.StatusOK, resp)
}

// GetUserChatRooms 获取用户加入的聊天室列表
func (h *ChatRoomHandler) GetUserChatRooms(c *gin.Context) {
	var req pb.GetUserChatRoomsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}

	resp, err := rpc.GetLogicExtClient().GetUserChatRooms(grpclib.NewContextFromGin(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
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
			"user_id":    member.UserId,
			"nickname":   member.Nickname,
			"avatar_url": member.AvatarUrl,
			"join_time":  member.JoinTime,
			"status":     member.Status,
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
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Unauthorized"})
		return
	}

	var req struct {
		Content  string `json:"content" binding:"required"`
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
		Content:  []byte(req.Content),
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

func (h *ChatRoomHandler) CheckPermissionsByUserId(c *gin.Context) {
	roomID, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的聊天室ID"})
		return
	}

	// 获取当前用户ID
	//userID, exists := middleware.GetUserID(c)
	//if !exists {
	//	c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Unauthorized"})
	//	return
	//}

	fmt.Printf("进入prc111")
	// 调用gRPC服务发送消息
	resp, err := rpc.GetLogicExtClient().GetChatRoom(grpclib.NewContextFromGin(c), &pb.GetChatRoomReq{
		RoomId: roomID,
	})
	//resp, err := rpc.GetLogicExtClient().CheckPermissionsByUserId(grpclib.NewContextFromGin(c), &pb.CheckPermissionsByUserIdReq{
	//	UserId: userID,
	//	RoomId: roomID,
	//})
	fmt.Printf("出prc111")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"permissions": resp.Room.Level},
	})
}
