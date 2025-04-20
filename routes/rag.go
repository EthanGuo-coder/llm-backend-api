package routes

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/EthanGuo-coder/llm-backend-api/middleware"
	"github.com/EthanGuo-coder/llm-backend-api/models"
	"github.com/EthanGuo-coder/llm-backend-api/services"
	"github.com/EthanGuo-coder/llm-backend-api/utils"
)

// RegisterRagRoutes 注册 RAG 相关路由
func RegisterRagRoutes(r *gin.Engine) {
	group := r.Group("/api/rag")
	group.Use(middleware.AuthMiddleware())
	{
		// 知识库管理
		group.POST("/kb/create", createKnowledgeBase)
		group.GET("/kb/list", listKnowledgeBases)
		group.POST("/kb/delete", deleteKnowledgeBase)

		// 文档管理
		group.POST("/doc/upload", uploadDocument)
		group.GET("/doc/list", listDocuments)
		group.POST("/doc/delete", deleteDocument)

		// 检索功能
		group.POST("/retrieve", retrieveInfo)

		// 基于知识库的对话
		group.POST("/chat", ragChat)

		// 元数据
		group.GET("/models", listEmbeddingModels)
	}
}

// createKnowledgeBase 创建知识库
func createKnowledgeBase(c *gin.Context) {
	var req models.CreateKnowledgeBaseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 获取用户 ID
	userID := utils.GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 用户 ID 转换为字符串
	uid := strconv.FormatInt(userID, 10)

	// 获取 RAG 服务
	ragService, err := services.GetRAGService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 RAG 服务失败: " + err.Error()})
		return
	}

	// 调用 RAG 服务创建知识库
	resp, err := ragService.CreateKnowledgeBase(uid, req.Name, req.EmbeddingModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Message})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// listKnowledgeBases 获取用户的知识库列表
func listKnowledgeBases(c *gin.Context) {
	// 获取用户 ID
	userID := utils.GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 用户 ID 转换为字符串
	uid := strconv.FormatInt(userID, 10)

	// 获取 RAG 服务
	ragService, err := services.GetRAGService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 RAG 服务失败: " + err.Error()})
		return
	}

	// 调用 RAG 服务获取知识库列表
	resp, err := ragService.ListKnowledgeBases(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Message})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// deleteKnowledgeBase 删除知识库
func deleteKnowledgeBase(c *gin.Context) {
	var req struct {
		KBID string `json:"kb_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 获取用户 ID
	userID := utils.GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 用户 ID 转换为字符串
	uid := strconv.FormatInt(userID, 10)

	// 获取 RAG 服务
	ragService, err := services.GetRAGService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 RAG 服务失败: " + err.Error()})
		return
	}

	// 调用 RAG 服务删除知识库
	resp, err := ragService.DeleteKnowledgeBase(uid, req.KBID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Message})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// uploadDocument 上传文档
func uploadDocument(c *gin.Context) {
	// 获取用户 ID
	userID := utils.GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 用户 ID 转换为字符串
	uid := strconv.FormatInt(userID, 10)

	// 获取知识库 ID
	kbID := c.PostForm("kb_id")
	if kbID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少知识库 ID"})
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未上传文件"})
		return
	}
	defer file.Close()

	// 读取文件内容
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件失败"})
		return
	}

	// 获取文件名和类型
	fileName := header.Filename
	fileExt := strings.ToLower(strings.TrimPrefix(filepath.Ext(fileName), "."))

	// 获取 RAG 服务
	ragService, err := services.GetRAGService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 RAG 服务失败: " + err.Error()})
		return
	}

	// 调用 RAG 服务上传文档
	resp, err := ragService.UploadDocument(uid, kbID, fileName, fileContent, fileExt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Message})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// listDocuments 获取知识库中的文档列表
func listDocuments(c *gin.Context) {
	kbID := c.Query("kb_id")
	if kbID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少知识库 ID 参数"})
		return
	}

	// 获取 RAG 服务
	ragService, err := services.GetRAGService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 RAG 服务失败: " + err.Error()})
		return
	}

	// 调用 RAG 服务获取文档列表
	resp, err := ragService.ListDocuments(kbID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Message})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// deleteDocument 删除文档
func deleteDocument(c *gin.Context) {
	var req struct {
		KBID  string `json:"kb_id" binding:"required"`
		DocID string `json:"doc_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 获取用户 ID
	userID := utils.GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 用户 ID 转换为字符串
	uid := strconv.FormatInt(userID, 10)

	// 获取 RAG 服务
	ragService, err := services.GetRAGService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 RAG 服务失败: " + err.Error()})
		return
	}

	// 调用 RAG 服务删除文档
	resp, err := ragService.DeleteDocument(uid, req.KBID, req.DocID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Message})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// retrieveInfo 从知识库检索信息
func retrieveInfo(c *gin.Context) {
	var req models.RetrieveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 获取 RAG 服务
	ragService, err := services.GetRAGService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 RAG 服务失败: " + err.Error()})
		return
	}

	// 调用 RAG 服务检索信息
	resp, err := ragService.RetrieveInfo(req.KBID, req.Query, req.TopK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Message})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ragChat 基于知识库的对话
func ragChat(c *gin.Context) {
	var req models.RagChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 获取用户 ID
	userID := utils.GetUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取 RAG 服务
	ragService, err := services.GetRAGService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 RAG 服务失败: " + err.Error()})
		return
	}

	// 生成基于知识库的提示
	prompt, err := ragService.GenerateRagPrompt(req.KBID, req.Message, req.TopK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成提示失败: " + err.Error()})
		return
	}

	// 使用提示进行对话
	if err := services.StreamSendMessage(c, req.ConversationID, prompt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// listEmbeddingModels 获取支持的嵌入模型列表
func listEmbeddingModels(c *gin.Context) {
	// 获取 RAG 服务
	ragService, err := services.GetRAGService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 RAG 服务失败: " + err.Error()})
		return
	}

	// 调用 RAG 服务获取嵌入模型列表
	resp, err := ragService.ListEmbeddingModels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Message})
		return
	}

	c.JSON(http.StatusOK, resp)
}
