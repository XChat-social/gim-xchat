package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"gim/pkg/grpclib/picker"
	_ "gim/pkg/grpclib/resolver/addrs"
	"gim/pkg/logger"
	"gim/pkg/protocol/pb"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

type defaultBuilder struct{}

func (*defaultBuilder) Build() Configuration {
	logger.Level = zap.DebugLevel
	logger.Target = logger.Console

	// 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:      "clustercfg.xchat-dev.y60xry.eun1.cache.amazonaws.com:6379", // 替换为你的 Redis 地址
		Password:  "",                                                          // 如果没有密码留空
		TLSConfig: &tls.Config{},                                               // 配置 TLS 支持
	})

	// 测试 Redis 连接
	if _, err := rdb.Ping().Result(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	return Configuration{
		MySQL:                "xchat:6TsXay5!h.pMnm3@tcp(database-1.chw4qwku6qx0.eu-north-1.rds.amazonaws.com:3306)/xchat?charset=utf8&parseTime=true",
		RedisHost:            "clustercfg.xchat-dev.y60xry.eun1.cache.amazonaws.com:6379", // Redis 地址
		RedisPassword:        "",
		PushRoomSubscribeNum: 100,
		PushAllSubscribeNum:  100,

		ConnectLocalAddr:     "127.0.0.1:8000",
		ConnectRPCListenAddr: ":8000",
		ConnectTCPListenAddr: ":8001",
		ConnectWSListenAddr:  ":8002",

		LogicRPCListenAddr:    ":8010",
		BusinessRPCListenAddr: ":8020",
		FileHTTPListenAddr:    "8030",

		ConnectIntClientBuilder: func() pb.ConnectIntClient {
			conn, err := grpc.DialContext(context.TODO(), "addrs:///127.0.0.1:8000", grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
				grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, picker.AddrPickerName)))
			if err != nil {
				panic(err)
			}
			return pb.NewConnectIntClient(conn)
		},
		LogicIntClientBuilder: func() pb.LogicIntClient {
			conn, err := grpc.DialContext(context.TODO(), "addrs:///127.0.0.1:8010", grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
				grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
			if err != nil {
				panic(err)
			}
			return pb.NewLogicIntClient(conn)
		},
		BusinessIntClientBuilder: func() pb.BusinessIntClient {
			conn, err := grpc.DialContext(context.TODO(), "addrs:///127.0.0.1:8020", grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
				grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
			if err != nil {
				panic(err)
			}
			return pb.NewBusinessIntClient(conn)
		},
	}
}
