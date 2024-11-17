package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EthanGuo-coder/llm-backend-api/services"
)

func RegisterChatRoutes(r *gin.Engine) {
	group := r.Group("/api/conversations/:conversation_id/messages")
	{
		group.POST("/", sendMessage)
	}
}

func sendMessage(c *gin.Context) {
	conversationID := c.Param("conversation_id")
	var req struct {
		Model   string `json:"model" binding:"required"`
		ApiKey  string `json:"api_key" binding:"required"`
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	response, err := services.SendMessage(conversationID, req.Model, req.ApiKey, req.Message)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}
	c.JSON(http.StatusOK, response)
}
