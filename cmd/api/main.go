package main

import (
	"crypto/tls"
	"flag"
	"gim/internal/api"
	"github.com/go-redis/redis"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// HTTP 服务地址
	httpServerEndpoint = flag.String("http-server-endpoint", ":8080", "HTTP server endpoint")

	// ✅ 生产数据库连接字符串
	dbConnStr = flag.String("db-conn-str", "xchat:6TsXay5!h.pMnm3@tcp(database-1.chw4qwku6qx0.eu-north-1.rds.amazonaws.com:3306)/xchat?charset=utf8&parseTime=true", "Database connection string")

	// ⬇️ 如需切换到本地数据库，取消注释下面一行
	// dbConnStr = flag.String("db-conn-str", "root:root@tcp(127.0.0.1:3306)/xchat?charset=utf8&parseTime=true", "Database connection string")

	// ✅ 生产 Redis 地址（Cluster + TLS）
	redisAddr = flag.String("redis-addr", "clustercfg.xchat-dev.y60xry.eun1.cache.amazonaws.com:6379", "Redis server address")

	// ⬇️ 如需切换到本地 Redis，取消注释下面一行
	// redisAddr = flag.String("redis-addr", "127.0.0.1:6379", "Redis server address")

	redisPassword = flag.String("redis-password", "", "Redis password")
)

func main() {
	flag.Parse()

	// 连接 MySQL 数据库
	db, err := gorm.Open(mysql.Open(*dbConnStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	log.Printf("✅ 成功连接到 MySQL 数据库: %s", *dbConnStr)

	// 连接 Redis（集群 + TLS）
	rdb := redis.NewClient(&redis.Options{
		Addr:     *redisAddr,
		Password: *redisPassword,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	// 测试 Redis 连接
	if pong, err := rdb.Ping().Result(); err != nil {
		log.Fatalf("❌ 无法连接到 Redis 集群: %v", err)
	} else {
		log.Printf("✅ 成功连接到 Redis 集群: %s，响应: %s", *redisAddr, pong)
	}

	// 创建 API 服务并启动 HTTP 服务器
	apiService := api.NewAPIService(db, rdb)
	r := apiService.InitRouter()

	log.Printf("🚀 HTTP API 服务启动在 %s...", *httpServerEndpoint)
	if err := r.Run(*httpServerEndpoint); err != nil {
		log.Fatalf("❌ 启动 HTTP 服务失败: %v", err)
	}
}
