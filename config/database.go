package config

import (
	"fmt"
	"log"

	"easy-agent/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	// 数据库连接配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",       // 用户名
		"root",       // 密码
		"localhost",  // 主机
		13306,        // 端口
		"easy_agent", // 数据库名
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表
	db.AutoMigrate(&models.User{})

	DB = db
}
