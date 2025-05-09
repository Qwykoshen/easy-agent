package config

import (
	"os"
)

var (
	DeepSeekAPIKey string
)

// InitDeepSeekConfig initializes the DeepSeek API configuration
func InitDeepSeekConfig() {
	// 从环境变量获取API密钥
	DeepSeekAPIKey = os.Getenv("DEEPSEEK_API_KEY")
	if DeepSeekAPIKey == "" {
		// 如果环境变量未设置，可以设置一个默认值或从配置文件读取
		// TODO: 实现从配置文件读取API密钥的逻辑
		panic("DEEPSEEK_API_KEY environment variable is not set")
	}
}
