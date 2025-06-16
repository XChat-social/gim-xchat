package rpc

import (
	"context"
	"gim/config"
	"gim/pkg/protocol/pb"
	"google.golang.org/grpc"
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
		conn, err := grpc.Dial("127.0.0.1:8010", grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		logicExtClient = pb.NewLogicExtClient(conn)
	}
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
