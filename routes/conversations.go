package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EthanGuo-coder/llm-backend-api/services"
)

func RegisterConversationRoutes(r *gin.Engine) {
	group := r.Group("/api/conversations")
	{
		group.POST("/", createConversation)
		group.GET("/:conversation_id", getConversationHistory)
	}
}

func createConversation(c *gin.Context) {
	var req struct {
		Model  string `json:"model" binding:"required"`
		ApiKey string `json:"api_key" binding:"required"`
		Title  string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	conversation, err := services.CreateConversation(req.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, conversation)
}

func getConversationHistory(c *gin.Context) {
	conversationID := c.Param("conversation_id")
	history, err := services.GetConversationHistory(conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}
	c.JSON(http.StatusOK, history)
}
