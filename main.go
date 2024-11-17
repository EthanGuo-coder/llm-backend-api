package main

import (
	"github.com/gin-gonic/gin"

	"github.com/EthanGuo-coder/llm-backend-api/routes"
	"github.com/EthanGuo-coder/llm-backend-api/storage"
)

func main() {
	// 初始化 Redis 存储
	storage.InitializeRedis("localhost:6379", "", 0)

	r := gin.Default()
	r.RedirectTrailingSlash = false

	routes.RegisterConversationRoutes(r)
	routes.RegisterChatRoutes(r)

	// 定义 POST 路由
	r.POST("/api/chat", streamChat)

	// 启动服务器
	r.Run(":8080") // 启动在8080端口
}
