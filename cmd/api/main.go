package main

import (
	"crypto/tls"
	"flag"
	"gim/internal/api"
	"log"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// HTTP 服务地址
	httpServerEndpoint = flag.String("http-server-endpoint", ":8080", "HTTP server endpoint")
	// 数据库连接字符串
	dbConnStr = flag.String("db-conn-str", "xchat:6TsXay5!h.pMnm3@tcp(database-1.chw4qwku6qx0.eu-north-1.rds.amazonaws.com:3306)/xchat?charset=utf8&parseTime=true", "Database connection string")
	//dbConnStr = flag.String("db-conn-str", "root:root@tcp(127.0.0.1:3306)/xchat?charset=utf8&parseTime=true", "Database connection string")
	// Redis连接地址
	redisAddr = flag.String("redis-addr", "127.0.0.1:6379", "Redis server address")
	// Redis密码
	redisPassword = flag.String("redis-password", "", "Redis password")
	// Redis数据库
	redisDB = flag.Int("redis-db", 0, "Redis database")
)

func main() {
	flag.Parse()

	// 连接数据库
	db, err := gorm.Open(mysql.Open(*dbConnStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	log.Printf("成功连接到MySQL数据库: %s", *dbConnStr)

	// 连接Redis，启用TLS模式
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // 密码已经设置为空
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	// 测试Redis连接
	pong, err := rdb.Ping().Result()
	if err != nil {
		log.Fatalf("无法连接到Redis: %v", err)
	}
	log.Printf("成功连接到Redis服务器(TLS模式): %s, 响应: %s", rdb.Options().Addr, pong)

	// 创建API服务
	apiService := api.NewAPIService(db, rdb)

	// 初始化路由
	r := apiService.InitRouter()

	// 启动HTTP服务
	log.Printf("HTTP API服务启动在 %s...", *httpServerEndpoint)
	if err := r.Run(*httpServerEndpoint); err != nil {
		log.Fatalf("启动HTTP服务失败: %v", err)
	}
}
