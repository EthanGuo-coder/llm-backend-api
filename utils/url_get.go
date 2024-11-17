package utils

import (
	"fmt"
	"strings"

	"github.com/EthanGuo-coder/llm-backend-api/constant"
)

var URLMapping = map[string]string{
	"gpt": constant.GPTBaseURL,
	"glm": constant.GLMBaseURL,
}

// GetBaseURL 根据关键字模糊匹配并返回对应的 BaseURL
func GetBaseURL(model string) (string, error) {
	// 转换关键字为小写，确保匹配不区分大小写
	keyword := strings.ToLower(model)

	// 遍历映射表，进行模糊匹配（基于前缀）
	for prefix, url := range URLMapping {
		if strings.HasPrefix(keyword, prefix) {
			return url, nil
		}
	}

	return "", fmt.Errorf("unsupported keyword: %s", keyword)
}
