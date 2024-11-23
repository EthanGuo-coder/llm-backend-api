package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EthanGuo-coder/llm-backend-api/middleware"
	"github.com/EthanGuo-coder/llm-backend-api/models"
	"github.com/EthanGuo-coder/llm-backend-api/services"
	"github.com/EthanGuo-coder/llm-backend-api/utils"
)

func RegisterConversationRoutes(r *gin.Engine) {
	group := r.Group("/api/conversations")
	{
		group.POST("/create", createConversation)
		group.GET("/history/:conversation_id", getConversationHistory)

		group.GET("/list", middleware.AuthMiddleware(), getUserConversations)
		group.POST("/del/:conversation_id", middleware.AuthMiddleware(), deleteConversation)
	}
}

func createConversation(c *gin.Context) {
	var req models.ConversationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 从上下文获取 userID
	userID := utils.GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 创建新会话
	conversation, err := services.CreateConversation(userID, req.Title, req.Model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

func getUserConversations(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	conversations, err := services.GetUserConversations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conversations)
}

func deleteConversation(c *gin.Context) {
	conversationID := c.Param("conversation_id")
	userID := utils.GetUserIDFromContext(c)

	if err := services.DeleteUserConversation(userID, conversationID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversation deleted successfully"})
}
