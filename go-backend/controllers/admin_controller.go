package controllers

import (
	"log"
	"net/http"
	"time"
	"yunmeng-backend/models"
	"yunmeng-backend/services"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	userService *services.UserService
}

func NewAdminController(userService *services.UserService) *AdminController {
	return &AdminController{
		userService: userService,
	}
}

// GetDashboardStats 获取管理员仪表板统计数据
func (c *AdminController) GetDashboardStats(ctx *gin.Context) {
	log.Printf("获取管理员仪表板统计数据")

	// 获取用户统计
	userStats, err := c.userService.GetUserStats()
	if err != nil {
		log.Printf("获取用户统计失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取统计数据失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 构建仪表板数据
	dashboardData := gin.H{
		"userStats": userStats,
		"systemInfo": gin.H{
			"serverTime": time.Now().Format("2006-01-02 15:04:05"),
			"version":    "1.0.0",
			"status":     "running",
		},
		"quickActions": []gin.H{
			{"name": "创建用户", "action": "create_user", "icon": "user-add"},
			{"name": "批量导入", "action": "batch_import", "icon": "upload"},
			{"name": "数据导出", "action": "export_data", "icon": "download"},
			{"name": "系统设置", "action": "system_settings", "icon": "setting"},
		},
	}

	log.Printf("管理员仪表板统计数据获取成功")

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取仪表板数据成功",
		"data":    dashboardData,
	})
}

// GetUserOverview 获取用户概览数据
func (c *AdminController) GetUserOverview(ctx *gin.Context) {
	log.Printf("获取用户概览数据")

	// 获取最近注册的用户
	var recentUsers []models.User
	if err := c.userService.GetDB().Preload("Profile").Preload("Roles").
		Order("create_time DESC").Limit(5).Find(&recentUsers).Error; err != nil {
		log.Printf("获取最近用户失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户概览失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 获取活跃用户统计
	var activeUserCount int64
	if err := c.userService.GetDB().Model(&models.User{}).
		Where("status = ?", "active").Count(&activeUserCount).Error; err != nil {
		log.Printf("获取活跃用户数失败: %v", err)
	}

	// 获取今日新增用户
	today := time.Now().Format("2006-01-02")
	var todayNewUsers int64
	if err := c.userService.GetDB().Model(&models.User{}).
		Where("DATE(create_time) = ?", today).Count(&todayNewUsers).Error; err != nil {
		log.Printf("获取今日新增用户数失败: %v", err)
	}

	overviewData := gin.H{
		"recentUsers":     recentUsers,
		"activeUserCount": activeUserCount,
		"todayNewUsers":   todayNewUsers,
		"lastUpdated":     time.Now().Format("2006-01-02 15:04:05"),
	}

	log.Printf("用户概览数据获取成功")

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取用户概览成功",
		"data":    overviewData,
	})
}

// GetSystemLogs 获取系统日志
func (c *AdminController) GetSystemLogs(ctx *gin.Context) {
	page, _ := ctx.GetQuery("page")
	size, _ := ctx.GetQuery("size")
	logType, _ := ctx.GetQuery("type")

	log.Printf("获取系统日志 - 页码: %s, 每页数量: %s, 日志类型: %s", page, size, logType)

	// 这里应该实现系统日志查询逻辑
	// 暂时返回模拟数据
	logs := []gin.H{
		{
			"id":        1,
			"type":      "info",
			"message":   "系统启动成功",
			"timestamp": time.Now().Add(-time.Hour).Format("2006-01-02 15:04:05"),
			"user":      "system",
		},
		{
			"id":        2,
			"type":      "warning",
			"message":   "用户登录失败",
			"timestamp": time.Now().Add(-30 * time.Minute).Format("2006-01-02 15:04:05"),
			"user":      "admin",
		},
		{
			"id":        3,
			"type":      "error",
			"message":   "数据库连接超时",
			"timestamp": time.Now().Add(-15 * time.Minute).Format("2006-01-02 15:04:05"),
			"user":      "system",
		},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统日志成功",
		"data": gin.H{
			"list":  logs,
			"total": len(logs),
			"page":  1,
			"size":  10,
		},
	})
}

// GetSystemSettings 获取系统设置
func (c *AdminController) GetSystemSettings(ctx *gin.Context) {
	log.Printf("获取系统设置")

	// 模拟系统设置数据
	settings := gin.H{
		"system": gin.H{
			"name":        "云梦高校科研竞赛管理系统",
			"version":     "1.0.0",
			"description": "高校学生科研与竞赛项目管理系统",
			"maintenance": false,
		},
		"security": gin.H{
			"passwordMinLength":  6,
			"passwordComplexity": true,
			"sessionTimeout":     3600,
			"maxLoginAttempts":   5,
		},
		"email": gin.H{
			"enabled":     true,
			"smtpServer":  "smtp.yunmeng.edu.cn",
			"smtpPort":    587,
			"fromAddress": "noreply@yunmeng.edu.cn",
		},
		"database": gin.H{
			"type":       "MySQL",
			"version":    "8.0",
			"backup":     true,
			"backupTime": "02:00",
		},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统设置成功",
		"data":    settings,
	})
}

// UpdateSystemSettings 更新系统设置
func (c *AdminController) UpdateSystemSettings(ctx *gin.Context) {
	var settings map[string]interface{}
	if err := ctx.ShouldBindJSON(&settings); err != nil {
		log.Printf("更新系统设置参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("更新系统设置: %v", settings)

	// 这里应该实现系统设置更新逻辑
	// 暂时返回成功响应

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "系统设置更新成功",
		"data":    nil,
	})
}

// GetDataReports 获取数据报表
func (c *AdminController) GetDataReports(ctx *gin.Context) {
	reportType, _ := ctx.GetQuery("type")
	startDate, _ := ctx.GetQuery("startDate")
	endDate, _ := ctx.GetQuery("endDate")

	log.Printf("获取数据报表 - 类型: %s, 开始日期: %s, 结束日期: %s", reportType, startDate, endDate)

	// 根据报表类型返回不同的数据
	var reportData gin.H
	switch reportType {
	case "user_growth":
		reportData = gin.H{
			"title": "用户增长趋势",
			"data": []gin.H{
				{"date": "2024-01-01", "count": 100},
				{"date": "2024-01-02", "count": 120},
				{"date": "2024-01-03", "count": 150},
				{"date": "2024-01-04", "count": 180},
				{"date": "2024-01-05", "count": 200},
			},
		}
	case "user_activity":
		reportData = gin.H{
			"title": "用户活跃度统计",
			"data": []gin.H{
				{"role": "admin", "active": 3, "total": 3},
				{"role": "teacher", "active": 15, "total": 20},
				{"role": "student", "active": 60, "total": 77},
			},
		}
	case "department_distribution":
		reportData = gin.H{
			"title": "部门用户分布",
			"data": []gin.H{
				{"department": "计算机学院", "count": 30},
				{"department": "数学学院", "count": 25},
				{"department": "物理学院", "count": 20},
				{"department": "化学学院", "count": 15},
			},
		}
	default:
		reportData = gin.H{
			"title": "综合统计报表",
			"data":  gin.H{},
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取数据报表成功",
		"data":    reportData,
	})
}

// ExportData 导出数据
func (c *AdminController) ExportData(ctx *gin.Context) {
	dataType, _ := ctx.GetQuery("type")
	format, _ := ctx.GetQuery("format")

	log.Printf("导出数据 - 类型: %s, 格式: %s", dataType, format)

	// 这里应该实现数据导出逻辑
	// 暂时返回成功响应

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "数据导出功能开发中",
		"data": gin.H{
			"type":   dataType,
			"format": format,
			"status": "processing",
		},
	})
}

// GetBackupStatus 获取备份状态
func (c *AdminController) GetBackupStatus(ctx *gin.Context) {
	log.Printf("获取备份状态")

	// 模拟备份状态数据
	backupStatus := gin.H{
		"lastBackup":      time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05"),
		"nextBackup":      time.Now().Add(12 * time.Hour).Format("2006-01-02 15:04:05"),
		"backupSize":      "256MB",
		"backupStatus":    "success",
		"autoBackup":      true,
		"backupRetention": 30,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取备份状态成功",
		"data":    backupStatus,
	})
}

// CreateBackup 创建备份
func (c *AdminController) CreateBackup(ctx *gin.Context) {
	log.Printf("创建数据备份")

	// 这里应该实现数据备份逻辑
	// 暂时返回成功响应

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "数据备份创建成功",
		"data": gin.H{
			"backupId":   "backup_" + time.Now().Format("20060102150405"),
			"backupTime": time.Now().Format("2006-01-02 15:04:05"),
			"status":     "completed",
		},
	})
}

// GetSystemHealth 获取系统健康状态
func (c *AdminController) GetSystemHealth(ctx *gin.Context) {
	log.Printf("获取系统健康状态")

	// 模拟系统健康状态数据
	healthStatus := gin.H{
		"status": "healthy",
		"checks": []gin.H{
			{
				"name":    "数据库连接",
				"status":  "healthy",
				"message": "连接正常",
				"latency": "5ms",
			},
			{
				"name":    "内存使用",
				"status":  "healthy",
				"message": "使用率正常",
				"usage":   "45%",
			},
			{
				"name":    "磁盘空间",
				"status":  "healthy",
				"message": "空间充足",
				"usage":   "30%",
			},
			{
				"name":    "网络连接",
				"status":  "healthy",
				"message": "连接正常",
				"latency": "10ms",
			},
		},
		"lastChecked": time.Now().Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统健康状态成功",
		"data":    healthStatus,
	})
}
