package models

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Conversation struct {
	ID          string    `json:"conversation_id"`
	Title       string    `json:"title"`
	Model       string    `json:"model"`
	ApiKey      string    `json:"api_key"` // 新增字段
	Messages    []Message `json:"messages"`
	CreatedTime int64     `json:"created_time"` // Unix 时间戳
}

type ConversationSummary struct {
	ID          string `json:"conversation_id"`
	Title       string `json:"title"`
	CreatedTime int64  `json:"created_time"`
}

type ConversationHistory struct {
	ID       string    `json:"conversation_id"`
	Title    string    `json:"title"`
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ConversationReq struct {
	Model string `json:"model" binding:"required"`
	Title string `json:"title" binding:"required"`
}

type CreateConversationReq struct {
	Model  string `json:"model" binding:"required"`
	Title  string `json:"title" binding:"required"`
	ApiKey string `json:"api_key" binding:"required"`
}
