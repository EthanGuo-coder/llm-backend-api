package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/EthanGuo-coder/llm-backend-api/config"
	"github.com/EthanGuo-coder/llm-backend-api/routes"
	"github.com/EthanGuo-coder/llm-backend-api/storage"
)

func main() {
	// 加载配置文件
	if err := config.LoadConfig("."); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	// 初始化 Redis
	if err := storage.InitializeRedis(); err != nil {
		log.Fatalf("Error initializing Redis: %v", err)
	}
	// 初始化 SQLite
	if err := storage.InitializeSQLite(); err != nil {
		log.Fatalf("Error initializing SQLite: %v", err)
	}

	r := gin.Default()
	r.RedirectTrailingSlash = true

	// 注册路由
	routes.RegisterRoutes(r)

	serverPort := config.AppConfig.Server.Port
	fmt.Printf("Server is running on port %d\n", serverPort)
	// 启动服务器
	r.Run(serverPort) // 启动在8080端口
}
