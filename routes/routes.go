package routes

import "github.com/gin-gonic/gin"

// RegisterRoutes 统一注册所有模块路由
func RegisterRoutes(r *gin.Engine) {
	// 用户相关路由
	RegisterUserRoutes(r)

	// 会话相关路由
	RegisterConversationRoutes(r)

	// 聊天相关路由
	RegisterChatRoutes(r)
}
