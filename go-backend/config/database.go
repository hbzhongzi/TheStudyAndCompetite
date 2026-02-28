package config

import (
	"fmt"
	"log"
	"os"

	"yunmeng-backend/models"
	"yunmeng-backend/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Charset  string
}

// NewDatabaseConfig 创建数据库配置
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "127.0.0.1"),
		Port:     getEnv("DB_PORT", "3306"),
		Username: getEnv("DB_USERNAME", "root"),
		Password: getEnv("DB_PASSWORD", "root"),
		Database: getEnv("DB_DATABASE", "cloud_dream_system"),
		Charset:  getEnv("DB_CHARSET", "utf8mb4"),
	}
}

// ConnectDatabase 连接数据库
func ConnectDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
	)

	// 配置GORM日志
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 获取底层的sql.DB对象
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100) // 最大打开连接数

	log.Println("数据库连接成功")
	return db, nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate(db *gorm.DB) error {
	log.Println("开始自动迁移数据库表...")

	// 自动迁移所有模型
	err := db.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.Role{},
		&models.UserRole{},
		&models.LoginLog{},
		&models.Project{},
		&models.ProjectType{},
		&models.StudentTeacher{},
	)

	if err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	log.Println("数据库表迁移完成")
	return nil
}

// InitDefaultData 初始化默认数据
func InitDefaultData(db *gorm.DB) error {
	log.Println("开始初始化默认数据...")

	// 检查是否已有角色数据
	var count int64
	db.Model(&models.Role{}).Count(&count)
	if count > 0 {
		log.Println("角色数据已存在，跳过角色初始化")
	} else {
		// 创建默认角色
		roles := []models.Role{
			{
				RoleKey:     "admin",
				RoleName:    "系统管理员",
				Description: "拥有系统所有权限",
			},
			{
				RoleKey:     "teacher",
				RoleName:    "教师",
				Description: "可以管理项目和指导学生",
			},
			{
				RoleKey:     "student",
				RoleName:    "学生",
				Description: "可以参与项目和竞赛",
			},
		}

		for _, role := range roles {
			if err := db.Create(&role).Error; err != nil {
				return fmt.Errorf("创建角色失败: %v", err)
			}
		}
		log.Println("默认角色数据初始化完成")
	}

	// 检查是否已有管理员用户
	var adminUser models.User
	err := db.Where("username = ?", "admin").First(&adminUser).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		log.Println("创建默认管理员用户...")

		// 创建管理员用户
		hashedPassword, err := utils.HashPassword("123456")
		if err != nil {
			return fmt.Errorf("密码加密失败: %v", err)
		}

		adminUser = models.User{
			Username: "admin",
			Password: hashedPassword,
			Email:    "admin@yunmeng.edu.cn",
			Status:   "active",
		}

		if err := db.Create(&adminUser).Error; err != nil {
			return fmt.Errorf("创建管理员用户失败: %v", err)
		}

		// 创建用户档案
		adminProfile := models.UserProfile{
			UserID:     adminUser.ID,
			RealName:   "系统管理员",
			Phone:      "",
			Department: "系统管理部",
		}

		if err := db.Create(&adminProfile).Error; err != nil {
			return fmt.Errorf("创建管理员档案失败: %v", err)
		}

		// 获取管理员角色
		var adminRole models.Role
		if err := db.Where("role_key = ?", "admin").First(&adminRole).Error; err != nil {
			return fmt.Errorf("获取管理员角色失败: %v", err)
		}

		// 关联用户和角色
		userRole := models.UserRole{
			UserID: adminUser.ID,
			RoleID: adminRole.ID,
		}

		if err := db.Create(&userRole).Error; err != nil {
			return fmt.Errorf("关联管理员角色失败: %v", err)
		}

		log.Println("默认管理员用户创建完成")
		log.Println("管理员账号: admin")
		log.Println("管理员密码: 123456")
	} else if err != nil {
		return fmt.Errorf("检查管理员用户失败: %v", err)
	} else {
		log.Println("管理员用户已存在，跳过创建")
	}

	log.Println("默认数据初始化完成")
	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
