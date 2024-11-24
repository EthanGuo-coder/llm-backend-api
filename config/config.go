package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/EthanGuo-coder/llm-backend-api/models"
)

var AppConfig *models.Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	viper.SetConfigName("config") // 配置文件名称
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(configPath)

	// 读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// 解析配置到结构体
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("failed to parse configuration file: %w", err)
	}

	return nil
}
