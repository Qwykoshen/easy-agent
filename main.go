package main

import (
	"easy-agent/config"
	"easy-agent/controllers"
	"easy-agent/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	config.InitDB()
	// 初始化DeepSeek配置
	config.InitDeepSeekConfig()
	config.InitBaidu()
	// 初始化DeepSeek客户端
	controllers.InitDeepSeekClient()

	// 创建Gin引擎
	r := gin.Default()

	// API路由组
	api := r.Group("/api/v1")
	{
		// 公开路由
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(middleware.JWTAuth())
		{
			// AI相关路由
			auth.POST("/chat", controllers.Chat)
			// 图像生成路由
			auth.POST("/generate-image", controllers.GenerateImage)
			// 语音生成路由
			auth.POST("/generate-voice", controllers.GenerateVoice)
		}
	}

	// 启动服务器
	r.Run(":8080")
}
