package main

import (
	"log"
	"os"

	"yunmeng-backend/config"
	"yunmeng-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("云梦高校科研竞赛管理系统后端服务启动中...")

	// 加载数据库配置
	dbConfig := config.NewDatabaseConfig()
	log.Printf("数据库配置: %s:%s/%s", dbConfig.Host, dbConfig.Port, dbConfig.Database)

	// 连接数据库
	db, err := config.ConnectDatabase(dbConfig)
	if err != nil {
		log.Fatal("数据库连接失败: ", err)
	}

	// 注意：如果遇到数据库迁移错误，请先手动执行 sql/migrate_existing.sql 脚本
	// 然后再取消下面的注释来启用自动迁移
	/*
		if err := config.AutoMigrate(db); err != nil {
			log.Fatal("数据库迁移失败: ", err)
		}
	*/

	// 初始化默认数据
	if err := config.InitDefaultData(db); err != nil {
		log.Fatal("初始化默认数据失败: ", err)
	}

	// 创建Gin实例
	r := gin.Default()

	// 添加CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))

	// 注册路由
	routes.RegisterRoutes(r, db)

	// 获取端口配置
	port := getEnv("PORT", "8080")
	log.Printf("服务器启动中，监听端口: %s", port)
	log.Println("前端地址: http://localhost:5173")
	log.Printf("后端API地址: http://localhost:%s/api", port)
	log.Println("注意：如果遇到数据库问题，请先执行 sql/migrate_existing.sql 脚本")

	// 启动服务
	if err := r.Run(":" + port); err != nil {
		log.Fatal("服务器启动失败: ", err)
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
