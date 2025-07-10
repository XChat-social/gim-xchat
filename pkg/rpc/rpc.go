package rpc

import (
	"context"
	"fmt"
	"gim/config"
	"gim/pkg/protocol/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	connectIntClient  pb.ConnectIntClient
	logicIntClient    pb.LogicIntClient
	businessIntClient pb.BusinessIntClient
	logicExtClient    pb.LogicExtClient
)

func GetConnectIntClient() pb.ConnectIntClient {
	if connectIntClient == nil {
		connectIntClient = config.Config.ConnectIntClientBuilder()
	}
	return connectIntClient
}

func GetLogicIntClient() pb.LogicIntClient {
	if logicIntClient == nil {
		logicIntClient = config.Config.LogicIntClientBuilder()
	}
	return logicIntClient
}

func GetBusinessIntClient() pb.BusinessIntClient {
	if businessIntClient == nil {
		businessIntClient = config.Config.BusinessIntClientBuilder()
	}
	return businessIntClient
}

func GetLogicExtClient() pb.LogicExtClient {
	if logicExtClient == nil {
		// 添加认证信息拦截器
		conn, err := grpc.Dial("127.0.0.1:8010",
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				// 从 gin.Context 中获取认证信息
				if md, ok := metadata.FromIncomingContext(ctx); ok {
					// 创建新的 metadata
					outCtx := metadata.NewOutgoingContext(ctx, md)
					// 使用新的 context 调用
					return invoker(outCtx, method, req, reply, cc, opts...)
				}
				return invoker(ctx, method, req, reply, cc, opts...)
			}),
		)
		if err != nil {
			panic(err)
		}
		logicExtClient = pb.NewLogicExtClient(conn)
	}
	fmt.Printf("获取client")
	return logicExtClient
}

func GetSender(deviceID, userID int64) (*pb.Sender, error) {
	user, err := GetBusinessIntClient().GetUser(context.TODO(), &pb.GetUserReq{UserId: userID})
	if err != nil {
		return nil, err
	}
	return &pb.Sender{
		UserId:    userID,
		DeviceId:  deviceID,
		AvatarUrl: user.User.AvatarUrl,
		Nickname:  user.User.Nickname,
		Extra:     user.User.Extra,
	}, nil
}
