package main

import (
	"gim/config"
	"gim/internal/business/api"
	"gim/pkg/interceptor"
	"gim/pkg/logger"
	"gim/pkg/protocol/pb"
	"gim/pkg/urlwhitelist"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// 初始化 gRPC 服务
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.NewInterceptor("business_interceptor", urlwhitelist.Business)))

	// 注册服务
	pb.RegisterBusinessIntServer(server, &api.BusinessIntServer{})
	pb.RegisterBusinessExtServer(server, &api.BusinessExtServer{})

	// gRPC-Web 包装器，添加 CORS 支持
	grpcWebServer := grpcweb.WrapServer(server,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false), // 允许所有请求
		grpcweb.WithAllowedRequestHeaders([]string{"*"}),  // 允许所有头
	)

	// 启动 HTTP 服务以支持 gRPC-Web
	go func() {
		httpServer := &http.Server{
			Addr: ":8081", // 例如 ":8081"
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if grpcWebServer.IsGrpcWebRequest(r) || grpcWebServer.IsAcceptableGrpcCorsRequest(r) {
					grpcWebServer.ServeHTTP(w, r)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}),
		}

		logger.Logger.Info("gRPC-Web 服务已经开启", zap.String("addr", ":8081"))
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Logger.Error("gRPC-Web 服务启动失败", zap.Error(err))
		}
	}()

	// 启动 gRPC 服务
	go func() {
		listen, err := net.Listen("tcp", config.Config.BusinessRPCListenAddr)
		if err != nil {
			logger.Logger.Fatal("gRPC 服务监听失败", zap.Error(err))
		}

		logger.Logger.Info("gRPC 服务已经开启", zap.String("addr", config.Config.BusinessRPCListenAddr))
		if err := server.Serve(listen); err != nil {
			logger.Logger.Error("gRPC 服务启动失败", zap.Error(err))
		}
	}()

	// 监听服务关闭信号，服务平滑重启
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	s := <-c
	logger.Logger.Info("接收到信号，准备关闭服务", zap.Any("signal", s))
	server.GracefulStop()
}
