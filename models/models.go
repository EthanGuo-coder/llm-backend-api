package models

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Conversation struct {
	ID          string    `json:"conversation_id"`
	Title       string    `json:"title"`
	Messages    []Message `json:"messages"`
	CreatedTime string    `json:"created_time"` // ISO 8601 格式的时间戳
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
