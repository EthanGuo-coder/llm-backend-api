package models

// KnowledgeBase 知识库模型
type KnowledgeBase struct {
	ID             string `json:"kb_id"`
	Name           string `json:"kb_name"`
	EmbeddingModel string `json:"embedding_model,omitempty"`
}

// Document 文档模型
type Document struct {
	ID       string `json:"doc_id"`
	KBID     string `json:"kb_id"`
	Name     string `json:"doc_name"`
	FileType string `json:"file_type"`
}

// RetrieveResult 检索结果模型
type RetrieveResult struct {
	Content string  `json:"content"`
	Score   float32 `json:"score"`
	DocID   string  `json:"doc_id"`
	DocName string  `json:"doc_name"`
}

// CreateKnowledgeBaseRequest 创建知识库请求
type CreateKnowledgeBaseRequest struct {
	Name           string `json:"kb_name" binding:"required"`
	EmbeddingModel string `json:"embedding_model,omitempty"`
}

// CreateKnowledgeBaseResponse 创建知识库响应
type CreateKnowledgeBaseResponse struct {
	Success bool   `json:"success"`
	KBID    string `json:"kb_id"`
	Message string `json:"message"`
}

// ListKnowledgeBasesResponse 知识库列表响应
type ListKnowledgeBasesResponse struct {
	Success bool            `json:"success"`
	KBs     []KnowledgeBase `json:"kbs"`
	Message string          `json:"message"`
}

// UploadDocumentResponse 上传文档响应
type UploadDocumentResponse struct {
	Success bool   `json:"success"`
	DocID   string `json:"doc_id"`
	Message string `json:"message"`
}

// RetrieveRequest 检索请求
type RetrieveRequest struct {
	KBID  string `json:"kb_id" binding:"required"`
	Query string `json:"query" binding:"required"`
	TopK  int    `json:"top_k,omitempty"`
}

// RetrieveResponse 检索响应
type RetrieveResponse struct {
	Success bool             `json:"success"`
	Results []RetrieveResult `json:"results"`
	Message string           `json:"message"`
}

// ListDocumentsResponse 文档列表响应
type ListDocumentsResponse struct {
	Success bool       `json:"success"`
	Docs    []Document `json:"docs"`
	Message string     `json:"message"`
}

// DeleteKnowledgeBaseResponse 删除知识库响应
type DeleteKnowledgeBaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// DeleteDocumentResponse 删除文档响应
type DeleteDocumentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ListEmbeddingModelsResponse 嵌入模型列表响应
type ListEmbeddingModelsResponse struct {
	Success bool     `json:"success"`
	Models  []string `json:"models"`
	Message string   `json:"message"`
}

// RagChatRequest 基于知识库的对话请求
type RagChatRequest struct {
	ConversationID int64  `json:"conversation_id" binding:"required"`
	KBID           string `json:"kb_id" binding:"required"`
	Message        string `json:"message" binding:"required"`
	TopK           int    `json:"top_k,omitempty"`
}
