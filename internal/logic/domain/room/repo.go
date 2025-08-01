package room

import (
	"context"
	"errors"
	"fmt"
	"gim/internal/api/models"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"gim/pkg/logger"
	"gim/pkg/protocol/pb"
	"gim/pkg/rpc"
	"gim/pkg/util"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	"time"

	"gorm.io/gorm"
)

type chatRoomRepo struct{}

var ChatRoomRepo = new(chatRoomRepo)

type chatRoomMemberRepo struct{}

var ChatRoomMemberRepo = new(chatRoomMemberRepo)

const (
	TaskDailySignIn     = 1001 // 每日签到
	TaskSevenDaySignIn  = 1002 // 签到七天
	TaskFollowTwitter   = 1003 // 关注推特
	taskStatusKeyPrefix = "task_status"
	TaskDailySum        = 1005
)

// Add 添加聊天室
func (r *chatRoomRepo) Add(ctx context.Context, chatRoom *pb.ChatRoom) error {
	// 转换为数据库模型
	modelChatRoom := &models.ChatRoom{
		RoomID:         chatRoom.RoomId,
		Name:           chatRoom.Name,
		AvatarURL:      chatRoom.AvatarUrl,
		Introduction:   chatRoom.Introduction,
		CreatorId:      chatRoom.CreatorId,
		CreatorName:    chatRoom.CreatorName,
		OnlineCount:    chatRoom.OnlineCount,
		MemberCount:    chatRoom.MemberCount,
		MaxMemberCount: chatRoom.MaxMemberCount,
		Extra:          chatRoom.Extra,
		CreateTime:     time.UnixMilli(chatRoom.CreateTime), // 从毫秒级时间戳转换为 time.Time
		UpdateTime:     time.UnixMilli(chatRoom.UpdateTime), // 从毫秒级时间戳转换为 time.Time
	}
	return db.DB.Create(modelChatRoom).Error
}

// Get 获取聊天室信息
func (r *chatRoomRepo) Get(ctx context.Context, roomId int64) (*pb.ChatRoom, error) {
	var chatRoom models.ChatRoom
	err := db.DB.Where("room_id = ?", roomId).First(&chatRoom).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 转换为 proto 消息
	return &pb.ChatRoom{
		RoomId:         chatRoom.RoomID,
		Name:           chatRoom.Name,
		AvatarUrl:      chatRoom.AvatarURL,
		Introduction:   chatRoom.Introduction,
		OnlineCount:    chatRoom.OnlineCount,
		MemberCount:    chatRoom.MemberCount,
		MaxMemberCount: chatRoom.MaxMemberCount,
		Extra:          chatRoom.Extra,
		CreateTime:     chatRoom.CreateTime.Unix(), // 转换为 Unix 时间戳
		UpdateTime:     chatRoom.UpdateTime.Unix(), // 转换为 Unix 时间戳
	}, nil
}

// List 获取聊天室列表
func (r *chatRoomRepo) List(ctx context.Context, offset, limit int32) ([]*pb.ChatRoom, error) {
	var modelChatRooms []models.ChatRoom
	err := db.DB.Offset(int(offset)).Limit(int(limit)).Find(&modelChatRooms).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 转换为 proto 消息列表
	chatRooms := make([]*pb.ChatRoom, 0, len(modelChatRooms))
	for _, room := range modelChatRooms {
		chatRooms = append(chatRooms, &pb.ChatRoom{
			RoomId:         room.RoomID,
			Name:           room.Name,
			AvatarUrl:      room.AvatarURL,
			Introduction:   room.Introduction,
			OnlineCount:    room.OnlineCount,
			MemberCount:    room.MemberCount,
			MaxMemberCount: room.MaxMemberCount,
			Extra:          room.Extra,
			CreateTime:     room.CreateTime.Unix(),
			UpdateTime:     room.UpdateTime.Unix(),
			CreatorId:      room.CreatorId,   // 添加
			CreatorName:    room.CreatorName, // 添加
			Level:          calculateRoomLevel(room.MemberCount),
			RandomCount:    float32(room.RandomCount),
		})
	}
	return chatRooms, nil
}

// ChatRoomMessage 仓储接口
type chatRoomMessageRepo struct{}

var ChatRoomMessageRepo = new(chatRoomMessageRepo)

// Add 添加聊天室消息
func (r *chatRoomMessageRepo) Add(ctx context.Context, message *pb.ChatRoomMessage) error {
	// 转换为数据库模型
	modelMessage := &models.ChatRoomMessage{
		RoomID:    uint64(message.RoomId),
		UserID:    uint64(message.UserId),
		RequestID: message.RequestId,
		Code:      int8(message.Code),
		Content:   message.Content,
		Seq:       uint64(message.Seq),
		SendTime:  time.Unix(message.SendTime, 0),
		Status:    int8(message.Status),
	}

	err := db.DB.Create(modelMessage).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// GetNextSeq 获取下一个消息序列号
func (r *chatRoomMessageRepo) GetNextSeq(roomId int64) (uint64, error) {
	var result struct {
		MaxSeq uint64 `gorm:"column:max_seq"`
	}
	err := db.DB.Model(&models.ChatRoomMessage{}).Where("room_id = ?", roomId).Select("COALESCE(MAX(seq), 0) as max_seq").Scan(&result).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return result.MaxSeq + 1, nil
}

// List 获取聊天室消息列表
func (r *chatRoomMessageRepo) List(ctx context.Context, roomId int64, offset, limit int32) ([]*pb.ChatRoomMessage, error) {
	var dbMessages []models.ChatRoomMessage
	err := db.DB.Where("room_id = ? AND status = 0", roomId).
		Order("seq DESC").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&dbMessages).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 提取所有消息的用户ID
	userIds := make(map[int64]int32)
	for _, msg := range dbMessages {
		userIds[int64(msg.UserID)] = 0 // 假设消息结构中SenderId字段表示发送者ID
	}

	// 获取所有用户信息并转换为map
	usersResp, err := rpc.GetBusinessIntClient().GetUsers(ctx, &pb.GetUsersReq{
		UserIds: userIds,
	})
	if err != nil {
		return nil, err
	}
	userMap := make(map[int64]*pb.User)
	for userId, user := range usersResp.Users {
		userMap[userId] = user
	}

	logger.Sugar.Info("users: %v", userMap)
	// 转换为 proto 消息列表
	messages := make([]*pb.ChatRoomMessage, 0, len(dbMessages))
	for i := len(dbMessages) - 1; i >= 0; i-- {
		msg := dbMessages[i]
		messages = append(messages, &pb.ChatRoomMessage{
			Id:           int64(msg.ID),
			RoomId:       int64(msg.RoomID),
			UserId:       int64(msg.UserID),
			RequestId:    msg.RequestID,
			LikeCount:    int64(msg.LikeCount),
			DislikeCount: int64(msg.DislikeCount),
			Code:         int32(msg.Code),
			Content:      msg.Content,
			Seq:          int64(msg.Seq),
			SendTime:     msg.SendTime.Unix(),
			Status:       int32(msg.Status),
			AvatarUrl:    userMap[int64(msg.UserID)].AvatarUrl, // 头像地址
			UserName:     userMap[int64(msg.UserID)].Nickname,
		})
	}
	return messages, nil
}

// Count 获取聊天室消息总数
func (r *chatRoomMessageRepo) Count(ctx context.Context, roomId int64) (int32, error) {
	var count int32
	err := db.DB.Model(&models.ChatRoomMessage{}).Where("room_id = ? AND status = 0", roomId).Count(&count).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return count, nil
}

// Count 获取聊天室总数
func (r *chatRoomRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := db.DB.Model(&models.ChatRoom{}).Count(&count).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return count, nil
}

// Add 添加聊天室成员
func (r *chatRoomMemberRepo) Add(ctx context.Context, member *pb.ChatRoomMember) error {
	// 转换为数据库模型
	dbMember := &models.ChatRoomMember{
		RoomID:    uint64(member.RoomId),
		UserID:    uint64(member.UserId),
		Nickname:  member.Nickname,
		AvatarURL: member.AvatarUrl,
		JoinTime:  time.UnixMilli(member.JoinTime),
		Status:    int8(member.Status),
	}
	return db.DB.Create(dbMember).Error
}

// Delete 删除聊天室成员
func (r *chatRoomMemberRepo) Delete(ctx context.Context, roomId, userId int64) error {
	return db.DB.Where("room_id = ? AND user_id = ?", roomId, userId).Delete(&models.ChatRoomMember{}).Error
}

// Count 获取聊天室成员总数
func (r *chatRoomMemberRepo) Count(ctx context.Context, roomId int64) (int64, error) {
	var count int64
	err := db.DB.Model(&models.ChatRoomMember{}).Where("room_id = ?", roomId).Count(&count).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return count, nil
}

// List 获取聊天室成员列表
func (r *chatRoomMemberRepo) List(ctx context.Context, roomId int64, offset, limit int32) ([]*pb.ChatRoomMember, error) {
	var dbMembers []models.ChatRoomMember
	err := db.DB.Where("room_id = ?", roomId).Offset(int(offset)).Limit(int(limit)).Find(&dbMembers).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 转换为 proto 消息列表
	members := make([]*pb.ChatRoomMember, 0, len(dbMembers))
	for _, member := range dbMembers {
		members = append(members, &pb.ChatRoomMember{
			RoomId:    int64(member.RoomID),
			UserId:    int64(member.UserID),
			Nickname:  member.Nickname,
			AvatarUrl: member.AvatarURL,
			JoinTime:  member.JoinTime.Unix(),
			Status:    int32(member.Status),
		})
	}
	return members, nil
}

// ListByUserId 获取用户加入的所有聊天室
func (r *chatRoomMemberRepo) ListByUserId(ctx context.Context, userId int64) ([]*pb.ChatRoomMember, error) {
	var dbMembers []models.ChatRoomMember
	err := db.DB.Where("user_id = ?", userId).Find(&dbMembers).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 转换为protobuf格式
	var members []*pb.ChatRoomMember
	for _, dbMember := range dbMembers {
		members = append(members, &pb.ChatRoomMember{
			RoomId:    int64(dbMember.RoomID),
			UserId:    int64(dbMember.UserID),
			Nickname:  dbMember.Nickname,
			AvatarUrl: dbMember.AvatarURL,
			JoinTime:  dbMember.JoinTime.Unix(),
			Status:    int32(dbMember.Status),
		})
	}
	return members, nil
}

// CountByUserId 获取用户加入的聊天室总数
func (r *chatRoomMemberRepo) CountByUserId(ctx context.Context, userId int64) (int64, error) {
	var count int64
	err := db.DB.Model(&models.ChatRoomMember{}).Where("user_id = ?", userId).Count(&count).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return count, nil
}

func (r *chatRoomMemberRepo) Get(ctx context.Context, roomId, userId int64) (*pb.ChatRoomMember, error) {
	var dbMember models.ChatRoomMember
	err := db.DB.Where("room_id = ? AND user_id = ?", roomId, userId).First(&dbMember).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	// 转换为protobuf格式
	member := &pb.ChatRoomMember{
		RoomId:    int64(dbMember.RoomID),
		UserId:    int64(dbMember.UserID),
		Nickname:  dbMember.Nickname,
		AvatarUrl: dbMember.AvatarURL,
		JoinTime:  dbMember.JoinTime.Unix(),
		Status:    int32(dbMember.Status),
	}
	return member, nil
}

func (r *chatRoomMemberRepo) AddUnreadCount(roomId int64, userId int64) error {
	return db.DB.Exec(
		"UPDATE chat_room_member SET unread_count = unread_count + ? WHERE room_id = ? AND status = 1 AND user_id != ?",
		1, roomId, userId,
	).Error
}

func (r *chatRoomMemberRepo) RefreshUnreadCount(roomId int64, userId int64) error {
	return db.DB.Model(&models.ChatRoomMember{}).Where("room_id = ? AND user_id = ?", roomId, userId).UpdateColumn("unread_count", 0).Error
}

// IncrOnlineCount 增加在线人数
func (r *chatRoomRepo) IncrOnlineCount(ctx context.Context, roomId int64) error {
	return db.DB.Model(&models.ChatRoom{}).Where("room_id = ?", roomId).UpdateColumn("online_count", gorm.Expr("online_count + ?", 1)).Error
}

// DecrOnlineCount 减少在线人数
func (r *chatRoomRepo) DecrOnlineCount(ctx context.Context, roomId int64) error {
	return db.DB.Model(&models.ChatRoom{}).Where("room_id = ?", roomId).UpdateColumn("online_count", gorm.Expr("online_count - ?", 1)).Error
}

// ListByCreatorId 获取用户创建的聊天室列表
func (r *chatRoomRepo) ListByCreatorId(ctx context.Context, creatorId int64) ([]*pb.ChatRoom, error) {
	var modelChatRooms []models.ChatRoom
	err := db.DB.Where("creator_id = ?", creatorId).Find(&modelChatRooms).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 转换为 proto 消息列表
	chatRooms := make([]*pb.ChatRoom, 0, len(modelChatRooms))
	for _, room := range modelChatRooms {
		chatRooms = append(chatRooms, &pb.ChatRoom{
			RoomId:         room.RoomID,
			Name:           room.Name,
			AvatarUrl:      room.AvatarURL,
			Introduction:   room.Introduction,
			OnlineCount:    room.OnlineCount,
			MemberCount:    room.MemberCount,
			MaxMemberCount: room.MaxMemberCount,
			Extra:          room.Extra,
			CreateTime:     room.CreateTime.Unix(),
			UpdateTime:     room.UpdateTime.Unix(),
		})
	}
	return chatRooms, nil
}

// ListByUserId 获取用户加入的聊天室列表
func (r *chatRoomRepo) ListByUserId(ctx context.Context, userId int64, offset, limit int32) ([]*pb.ChatRoomAndUnreadCount, error) {
	var modelChatRooms []models.ChatRoomAndUnreadCount

	err := db.DB.Table("chat_room").
		Select("chat_room.*, chat_room_member.unread_count as unread_count").
		Joins("JOIN chat_room_member ON chat_room.room_id = chat_room_member.room_id").
		Where("chat_room_member.user_id = ? and chat_room_member.status = 1", userId).
		Offset(int(offset)).Limit(int(limit)).
		Find(&modelChatRooms).Error

	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	// 提取roomIDs
	roomIDs := make([]int64, 0, len(modelChatRooms))
	for _, room := range modelChatRooms {
		roomIDs = append(roomIDs, room.RoomID)
	}

	// 查询每个聊天室的最新消息
	messageMap := make(map[int64]string)
	if len(roomIDs) > 0 {
		var messages []models.ChatRoomMessage
		err := db.DB.Raw(
			" SELECT * FROM (SELECT id, content, room_id, seq, ROW_NUMBER()OVER(PARTITION BY room_id ORDER BY seq DESC) AS rowNumber FROM chat_room_message where room_id IN (?)) t WHERE rowNumber = 1",
			roomIDs).
			Scan(&messages).Error

		if err == nil {
			for _, msg := range messages {
				if _, exists := messageMap[int64(msg.RoomID)]; !exists {
					messageMap[int64(msg.RoomID)] = string(msg.Content)
				}
			}
		}
	}

	// 构造返回结果
	chatRooms := make([]*pb.ChatRoomAndUnreadCount, 0, len(modelChatRooms))
	for _, room := range modelChatRooms {
		chatRooms = append(chatRooms, &pb.ChatRoomAndUnreadCount{
			RoomId:         room.RoomID,
			Name:           room.Name,
			AvatarUrl:      room.AvatarURL,
			Introduction:   room.Introduction,
			CreatorId:      room.CreatorId,
			OnlineCount:    room.OnlineCount,
			MemberCount:    room.MemberCount,
			MaxMemberCount: room.MaxMemberCount,
			Extra:          room.Extra,
			CreateTime:     room.CreateTime.Unix(),
			UpdateTime:     room.UpdateTime.Unix(),
			Level:          calculateRoomLevel(room.MemberCount),
			UnreadCount:    room.UnreadCount,
			LastMessage:    messageMap[room.RoomID], // 从map中获取
		})
	}

	return chatRooms, nil
}

func (r *chatRoomRepo) CheckPermissionsByUserId(ctx context.Context, req *pb.CheckPermissionsByUserIdReq) (*pb.CheckPermissionsByUserIdResp, error) {
	// 校验当前人为创建的聊天室则返回true
	//creator, _ := checkCreator(req.RoomId, req.UserId)
	//if creator != (models.ChatRoom{}) {
	//	return &pb.CheckPermissionsByUserIdResp{
	//		HasPermission: true,
	//	}, nil
	//}

	exists, _ := ChatRoomRepo.CheckUserIsHolding(req.UserId, req.RoomId)

	if exists {
		// 校验当前用户是否在聊天室
		chatMember, _ := checkMember(req.RoomId, req.UserId)

		if chatMember == (models.ChatRoomMember{}) {
			// 获取用户信息
			user, _ := getUser(req.UserId)
			if user != (models.User{}) {
				logger.Sugar.Info("用户加入聊天室 - 用户ID: %d, 昵称: %s, 头像: %s, 房间ID: %d", user.ID, user.Nickname, user.AvatarURL, req.RoomId)
				// 创建成员信息
				member := &pb.ChatRoomMember{
					UserId:    int64(user.ID),
					Nickname:  user.Nickname,
					AvatarUrl: user.AvatarURL,
					RoomId:    req.RoomId,
					Status:    1,
					JoinTime:  util.UnixMilliTime(time.Now()),
				}

				// 添加成员
				err := ChatRoomMemberRepo.Add(ctx, member)
				if err != nil {
					return nil, err
				}
			}
		} else {
			logger.Sugar.Info("更新用户状态：%v", chatMember)
			err := updateChatMemberStatus(chatMember.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	return &pb.CheckPermissionsByUserIdResp{
		HasPermission: exists,
	}, nil
}

func (r *chatRoomRepo) CheckUserIsHolding(userId int64, roomId int64) (bool, error) {
	var count int64
	query := db.DB.Table("chat_room").
		Joins("INNER JOIN token ON chat_room.creator_id = token.user_id").
		Joins("INNER JOIN token_holdings ON token.token_address = token_holdings.token_address").
		Where("token_holdings.user_id = ? AND chat_room.room_id = ? AND amount > 0", userId, roomId)

	result := query.Count(&count)

	if result.Error != nil {
		// 明确处理记录不存在的情况
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		// 其他数据库错误
		logger.Sugar.Info("数据库查询错误: %v", result.Error)
		return false, nil
	}

	exists := count > 0
	logger.Sugar.Info("权限检查结果 - 用户ID: %d, 房间ID: %d, 权限: %v", userId, roomId, exists)
	return exists, nil
}

func (r *chatRoomRepo) CheckThumbed(userId int64, messageId int64) (bool, error) {
	var count int64
	query := db.DB.Table("chat_room_message").
		Joins("INNER JOIN chat_room_message_thumbs ON chat_room_message.id = chat_room_message_thumbs.message_id").
		Where("chat_room_message_thumbs.user_id = ? AND chat_room_message_thumbs.message_id = ? ", userId, messageId)
	result := query.Count(&count)

	if result.Error != nil {
		// 明确处理记录不存在的情况
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		// 其他数据库错误
		logger.Sugar.Info("数据库查询错误: %v", result.Error)
		return false, nil
	}

	exists := count > 0
	logger.Sugar.Infow("权限检查结果",
		"userId", userId,
		"messageId", messageId,
		"hasPermission", exists,
	)
	return exists, nil
}

func (r *chatRoomRepo) ThumbAndXpoint(userId int64, messageId int64, isLike bool, roomId int64, messageUserId int64) error {
	// 开启事务
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	flag := 1
	if isLike {
		flag = 0
	}
	//thumb
	thumb := &models.ChatRoomMessageThumbs{
		UserID:     uint64(userId),
		MessageID:  messageId,
		Thumb:      int8(flag),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if err := tx.Create(&thumb).Error; err != nil {
		tx.Rollback()
		return err
	}
	if isLike {
		if err := tx.Exec(
			"UPDATE chat_room_message SET like_count = like_count + 1 WHERE id = ?",
			messageId,
		).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err := tx.Exec(
			"UPDATE chat_room_message SET dislike_count = dislike_count + 1 WHERE id = ?",
			messageId,
		).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 获取积分权重
	roomInfo := &models.ChatRoom{}
	if err := tx.Model(&models.ChatRoom{}).Where("room_id = ?", roomId).First(&roomInfo).Error; err != nil {
		tx.Rollback()
		return err
	}
	xPoint := roomInfo.RandomCount * (float64(calculateRoomLevel(roomInfo.MemberCount))*0.1 + 1)
	logger.Sugar.Info("xPoint: %f", xPoint)

	// 查询当前人的积分
	xPointInfo := &models.User{}
	if err := tx.Where("id = ?", messageUserId).Select("xpoint").First(&xPointInfo).Error; err != nil {
		tx.Rollback()
		return err
	}
	point := xPointInfo.XPoint

	// 更新当前人的积分【like时新增，dislike时减少】
	if isLike {
		if err := tx.Model(&models.User{}).Where("id = ?", messageUserId).Update("xpoint", xPoint+point).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 增加redis值
		err := updateRedisSum(messageUserId, roomId, xPoint)
		if err != nil {
			return err
		}

	} else {
		if point-xPoint > 0 {
			if point >= xPoint {
				if err := tx.Model(&models.User{}).Where("id = ?", messageUserId).Update("xpoint", point-xPoint).Error; err != nil {
					tx.Rollback()
					return err
				}
			}

		} else {
			if err := tx.Model(&models.User{}).Where("id = ?", messageUserId).Update("xpoint", 0).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		// 减少redis值
		err := updateRedisSum(messageUserId, roomId, -xPoint)
		if err != nil {
			return err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func updateRedisSum(messageUserId int64, roomId int64, xPoint float64) error {
	// 构造任务统计 Redis Key
	sumKey := fmt.Sprintf("%s:%d:%d:%d", taskStatusKeyPrefix, messageUserId, roomId, TaskDailySum)
	// 查询当前用户是否已存在每日统计
	redisRes, err := db.RedisCli.Get(sumKey).Result()

	// string 转 decimal
	changePoint := decimal.NewFromFloat(xPoint)

	if errors.Is(err, redis.Nil) {
		logger.Sugar.Info("daily not found")
		err = setWithMidnightExpire(sumKey, changePoint.Round(2))
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	dailySum, err := decimal.NewFromString(redisRes)
	if err != nil {
		return err
	}
	logger.Sugar.Info("dailySum: %f", dailySum)
	// 累加
	err = setWithMidnightExpire(sumKey, dailySum.Add(changePoint).Round(2))
	if err != nil {
		return err
	}
	return nil
}

func getUser(userId int64) (models.User, error) {
	var user models.User
	result := db.DB.First(&user, userId)

	if result.Error != nil {
		// 处理记录不存在的情况
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 返回空用户结构体和自定义错误
			return models.User{}, nil
		}
		// 其他数据库错误
		return models.User{}, gerrors.WrapError(result.Error)
	}
	return user, nil
}

func checkMember(roomId int64, userId int64) (models.ChatRoomMember, error) {
	var chatRoomMember models.ChatRoomMember
	result := db.DB.Table("chat_room_member").
		Where("room_id = ? AND user_id = ?", roomId, userId).
		First(&chatRoomMember)
	if result.Error != nil {
		// 处理记录不存在的情况
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 返回空用户结构体和自定义错误
			return models.ChatRoomMember{}, nil
		}
		// 其他数据库错误
		return models.ChatRoomMember{}, gerrors.WrapError(result.Error)
	}
	return chatRoomMember, nil
}

func checkCreator(roomId int64, userId int64) (models.ChatRoom, error) {
	var chatRoom models.ChatRoom
	result := db.DB.Table("chat_room").
		Where("room_id = ? AND creator_id = ?", roomId, userId).
		First(&chatRoom)
	if result.Error != nil {
		// 处理记录不存在的情况
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 返回空用户结构体和自定义错误
			return models.ChatRoom{}, nil
		}
		// 其他数据库错误
		return models.ChatRoom{}, gerrors.WrapError(result.Error)
	}
	return chatRoom, nil
}

// 更新聊天室状态
func updateChatMemberStatus(memberId uint) error {
	return db.DB.Model(&models.ChatRoomMember{}).Where("id = ?", memberId).UpdateColumn("status", 1).Error
}

func setWithMidnightExpire(key string, value decimal.Decimal) error {
	now := time.Now()
	loc := now.Location()

	// 计算次日00:00
	tomorrow := now.AddDate(0, 0, 1)
	endOfDay := time.Date(
		tomorrow.Year(),
		tomorrow.Month(),
		tomorrow.Day(),
		0, 0, 0, 0, loc,
	)

	// 计算时间差
	expireIn := endOfDay.Sub(now)

	// 处理时间穿越情况
	if expireIn < 0 {
		expireIn = 0 // 立即过期
	}
	logger.Sugar.Info("value: %f", value)
	return db.RedisCli.Set(key, value.String(), expireIn).Err()
}
