package models

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
	Message string `json:"message" binding:"required"`
}
