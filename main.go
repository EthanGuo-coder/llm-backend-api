package main

import (
	"github.com/gin-gonic/gin"

	"github.com/EthanGuo-coder/llm-backend-api/routes"
	"github.com/EthanGuo-coder/llm-backend-api/storage"
)

func main() {
	// 初始化 SQLite 数据库
	dbPath := "llm_backend.db"
	storage.InitializeSQLite(dbPath)

	// 初始化 Redis
	redisAddr := "localhost:6379"
	redisPassword := ""
	redisDB := 0
	storage.InitializeRedis(redisAddr, redisPassword, redisDB)

	// 初始化 Redis 存储
	storage.InitializeRedis("localhost:6379", "", 0)

	r := gin.Default()
	r.RedirectTrailingSlash = true

	routes.RegisterConversationRoutes(r)
	routes.RegisterChatRoutes(r)
	routes.RegisterUserRoutes(r)

	// 启动服务器
	r.Run(":8080") // 启动在8080端口
}
