package models

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Conversation struct {
	ID          string    `json:"conversation_id"`
	Title       string    `json:"title"`
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	CreatedTime int64     `json:"created_time"` // Unix 时间戳
}

// SSEDelta 定义结构体以匹配 JSON 数据格式
type SSEDelta struct {
	Content string `json:"content"`
}

type SSEChoice struct {
	Delta SSEDelta `json:"delta"`
}

type SSEResponse struct {
	Choices []SSEChoice `json:"choices"`
}

type AskReq struct {
	ApiKey  string `json:"api_key" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type ConversationReq struct {
	Model string `json:"model" binding:"required"`
	Title string `json:"title" binding:"required"`
}
