package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
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

	// 配置 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 注册路由
	routes.RegisterRoutes(r)

	serverPort := config.AppConfig.Server.Port
	fmt.Printf("Server is running on port %d\n", serverPort)
	// 启动服务器
	r.Run(serverPort) // 启动在8080端口
}
