package config

import (
	"os"
)

var (
	BaiduAPIKey    string
	BaiduSecretKey string
	AccessToken    string
)

func InitBaidu() {
	BaiduAPIKey = os.Getenv("BAIDU_API_KEY")
	BaiduSecretKey = os.Getenv("BAIDU_SECRET_KEY")
	AccessToken = os.Getenv("BAIDU_ACCESS_TOKEN")
}
