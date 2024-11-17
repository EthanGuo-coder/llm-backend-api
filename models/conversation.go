package models

type Conversation struct {
	ID          string    `json:"conversation_id"`
	Title       string    `json:"title"`
	Messages    []Message `json:"messages"`
	CreatedTime string    `json:"created_time"` // ISO 8601 格式的时间戳
}
