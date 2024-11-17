package routes

import (
	"github.com/EthanGuo-coder/llm-backend-api/models"
	"log"
	"net/http"

	"github.com/EthanGuo-coder/llm-backend-api/services"

	"github.com/gin-gonic/gin"
)

func RegisterChatRoutes(r *gin.Engine) {
	group := r.Group("/api/conversations/:conversation_id/messages")
	{
		group.POST("/", streamSendMessage) // 流式返回消息
	}
}

func streamSendMessage(c *gin.Context) {
	conversationID := c.Param("conversation_id")

	var req *models.Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 流式处理消息并返回 SSE
	log.Println(conversationID, req.Model, req.ApiKey, req.Message)
	if err := services.StreamSendMessage(c, conversationID, req.Model, req.ApiKey, req.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
