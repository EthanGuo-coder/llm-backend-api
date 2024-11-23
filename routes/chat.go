package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EthanGuo-coder/llm-backend-api/middleware"
	"github.com/EthanGuo-coder/llm-backend-api/models"
	"github.com/EthanGuo-coder/llm-backend-api/services"
)

func RegisterChatRoutes(r *gin.Engine) {
	group := r.Group("/api/chat/:conversation_id")
	{
		group.POST("/", middleware.AuthMiddleware(), streamSendMessage) // 流式返回消息
	}
}

func streamSendMessage(c *gin.Context) {
	conversationID := c.Param("conversation_id")

	var req *models.AskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 流式处理消息并返回 SSE
	if err := services.StreamSendMessage(c, conversationID, req.ApiKey, req.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
