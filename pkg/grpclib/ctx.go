package grpclib

import (
	"context"
	"gim/pkg/gerrors"
	"gim/pkg/logger"
	"strconv"

	"google.golang.org/grpc/metadata"
)

const (
	CtxUserId       = "user_id"
	CtxDeviceId     = "device_id"
	CtxToken        = "token"
	CtxRequestId    = "request_id"
	CtxUserIdDash   = "user-id"   // 支持的 user-id
	CtxDeviceIdDash = "device-id" // 支持的 device-id
)

func ContextWithRequestId(ctx context.Context, requestId int64) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(CtxRequestId, strconv.FormatInt(requestId, 10)))
}

func Get(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values, ok := md[key]
	if !ok || len(values) == 0 {
		return ""
	}
	return values[0]
}

// GetCtxRequestId 获取ctx的app_id
func GetCtxRequestId(ctx context.Context) int64 {
	requestIdStr := Get(ctx, CtxRequestId)
	requestId, err := strconv.ParseInt(requestIdStr, 10, 64)
	if err != nil {
		return 0
	}
	return requestId
}

// GetCtxData 获取ctx的用户数据，依次返回user_id, device_id
func GetCtxData(ctx context.Context) (int64, int64, error) {
	var (
		userId   int64
		deviceId int64
		err      error
	)
	token := Get(ctx, CtxToken)
	userId1 := Get(ctx, CtxUserId)
	deviceId1 := Get(ctx, CtxDeviceId)

	// Logging for debugging purposes
	logger.Sugar.Info("token:", token)
	logger.Sugar.Info("userId1:", userId1)
	logger.Sugar.Info("deviceId1:", deviceId1)

	// 优先检查 user_id
	userIdStr := Get(ctx, CtxUserId)
	if userIdStr == "" {
		// 如果没有找到 user_id，尝试 user-id
		userIdStr = Get(ctx, CtxUserIdDash)
	}
	logger.Sugar.Info("userIdStr:", userIdStr)
	userId, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, 0, gerrors.ErrUnauthorized
	}

	// 获取 device_id
	deviceIdStr := Get(ctx, CtxDeviceId)
	if deviceIdStr == "" {
		// 如果没有找到 device_id，尝试 device-id
		deviceIdStr = Get(ctx, CtxDeviceIdDash)
	}
	deviceId, err = strconv.ParseInt(deviceIdStr, 10, 64)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, 0, gerrors.ErrUnauthorized
	}
	return userId, deviceId, nil
}

// GetCtxToken 获取ctx的token
func GetCtxToken(ctx context.Context) string {
	return Get(ctx, CtxToken)
}

// NewAndCopyRequestId 创建一个context,并且复制RequestId
func NewAndCopyRequestId(ctx context.Context) context.Context {
	newCtx := context.TODO()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return newCtx
	}

	requestIds, ok := md[CtxRequestId]
	if !ok && len(requestIds) == 0 {
		return newCtx
	}
	return metadata.NewOutgoingContext(newCtx, metadata.Pairs(CtxRequestId, requestIds[0]))
}
