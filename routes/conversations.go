package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/EthanGuo-coder/llm-backend-api/middleware"
	"github.com/EthanGuo-coder/llm-backend-api/models"
	"github.com/EthanGuo-coder/llm-backend-api/services"
	"github.com/EthanGuo-coder/llm-backend-api/utils"
)

func RegisterConversationRoutes(r *gin.Engine) {
	group := r.Group("/api/conversations")
	{
		group.POST("/create", middleware.AuthMiddleware(), createConversation)                      // 创建新会话
		group.GET("/history/:conversation_id", middleware.AuthMiddleware(), getConversationHistory) // 用户单会话对话记录

		group.GET("/list", middleware.AuthMiddleware(), getUserConversations)                // 用户会话列表
		group.POST("/del/:conversation_id", middleware.AuthMiddleware(), deleteConversation) // 删除用户会话（某一个）
	}
}

func createConversation(c *gin.Context) {
	var req *models.CreateConversationReq

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
	conversation, err := services.CreateConversation(userID, req.Title, req.Model, req.ApiKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conversation)
}

func getConversationHistory(c *gin.Context) {
	conversationIDStr := c.Param("conversation_id")
	// 将字符串转换为 int64
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		// 处理转换错误（例如：参数不是有效的数字）
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}
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
	conversationIDStr := c.Param("conversation_id")
	// 将字符串转换为 int64
	conversationID, err := strconv.ParseInt(conversationIDStr, 10, 64)
	if err != nil {
		// 处理转换错误（例如：参数不是有效的数字）
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}
	userID := utils.GetUserIDFromContext(c)

	if err := services.DeleteUserConversation(userID, conversationID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversation deleted successfully"})
}
