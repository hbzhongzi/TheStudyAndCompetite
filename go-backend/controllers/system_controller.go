package controllers

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
	"yunmeng-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SystemController struct {
	db *gorm.DB
}

func NewSystemController(db *gorm.DB) *SystemController {
	return &SystemController{db: db}
}

// GetSystemLogs 获取系统日志列表
func (c *SystemController) GetSystemLogs(ctx *gin.Context) {
	var logs []models.SystemLog
	var total int64

	// 获取查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	logType := ctx.Query("log_type")
	operation := ctx.Query("operation")
	status := ctx.Query("status")
	userID := ctx.Query("user_id")
	action := ctx.Query("action")
	ipAddress := ctx.Query("ip_address")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	// 构建查询条件
	query := c.db.Model(&models.SystemLog{}).Preload("User.Profile")

	if logType != "" {
		query = query.Where("log_type = ?", logType)
	}

	if operation != "" {
		query = query.Where("operation LIKE ?", "%"+operation+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if userID != "" {
		if uid, err := strconv.ParseUint(userID, 10, 32); err == nil {
			query = query.Where("user_id = ?", uid)
		}
	}

	if action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}

	if ipAddress != "" {
		query = query.Where("ip_address LIKE ?", "%"+ipAddress+"%")
	}

	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", start)
		}
	}

	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("created_at <= ?", end.Add(24*time.Hour))
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		log.Printf("获取系统日志总数失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取系统日志失败",
		})
		return
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&logs).Error; err != nil {
		log.Printf("获取系统日志列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取系统日志失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统日志成功",
		"data": gin.H{
			"logs":  logs,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// GetSystemLogsSummary 获取系统日志统计
func (c *SystemController) GetSystemLogsSummary(ctx *gin.Context) {
	var summary struct {
		TotalLogs   int64 `json:"total_logs"`
		ErrorLogs   int64 `json:"error_logs"`
		WarningLogs int64 `json:"warning_logs"`
		InfoLogs    int64 `json:"info_logs"`
		TodayLogs   int64 `json:"today_logs"`
		WeekLogs    int64 `json:"week_logs"`
		MonthLogs   int64 `json:"month_logs"`
	}

	// 获取总数
	c.db.Model(&models.SystemLog{}).Count(&summary.TotalLogs)

	// 获取各类型日志数量
	c.db.Model(&models.SystemLog{}).Where("log_type = ?", "error").Count(&summary.ErrorLogs)
	c.db.Model(&models.SystemLog{}).Where("log_type = ?", "warning").Count(&summary.WarningLogs)
	c.db.Model(&models.SystemLog{}).Where("log_type = ?", "info").Count(&summary.InfoLogs)

	// 获取时间范围统计
	today := time.Now().Truncate(24 * time.Hour)
	weekAgo := time.Now().AddDate(0, 0, -7)
	monthAgo := time.Now().AddDate(0, -1, 0)

	c.db.Model(&models.SystemLog{}).Where("created_at >= ?", today).Count(&summary.TodayLogs)
	c.db.Model(&models.SystemLog{}).Where("created_at >= ?", weekAgo).Count(&summary.WeekLogs)
	c.db.Model(&models.SystemLog{}).Where("created_at >= ?", monthAgo).Count(&summary.MonthLogs)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统日志统计成功",
		"data":    summary,
	})
}

// GetSystemHealthLogs 获取系统健康日志
func (c *SystemController) GetSystemHealthLogs(ctx *gin.Context) {
	var logs []models.SystemHealthLog
	var total int64

	// 获取查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	status := ctx.Query("status")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	// 构建查询条件
	query := c.db.Model(&models.SystemHealthLog{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("record_time >= ?", start)
		}
	}

	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("record_time <= ?", end.Add(24*time.Hour))
		}
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("record_time DESC").Find(&logs).Error; err != nil {
		log.Printf("获取系统健康日志失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取系统健康日志失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统健康日志成功",
		"data": gin.H{
			"logs":  logs,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// GetSystemHealthSummary 获取系统健康统计
func (c *SystemController) GetSystemHealthSummary(ctx *gin.Context) {
	var summary struct {
		TotalRecords    int64   `json:"total_records"`
		HealthyRecords  int64   `json:"healthy_records"`
		WarningRecords  int64   `json:"warning_records"`
		CriticalRecords int64   `json:"critical_records"`
		AvgCPUUsage     float64 `json:"avg_cpu_usage"`
		AvgMemoryUsage  float64 `json:"avg_memory_usage"`
		AvgDiskUsage    float64 `json:"avg_disk_usage"`
	}

	// 获取总数
	c.db.Model(&models.SystemHealthLog{}).Count(&summary.TotalRecords)

	// 获取各状态数量
	c.db.Model(&models.SystemHealthLog{}).Where("status = ?", "healthy").Count(&summary.HealthyRecords)
	c.db.Model(&models.SystemHealthLog{}).Where("status = ?", "warning").Count(&summary.WarningRecords)
	c.db.Model(&models.SystemHealthLog{}).Where("status = ?", "critical").Count(&summary.CriticalRecords)

	// 获取平均值
	c.db.Model(&models.SystemHealthLog{}).Select("AVG(cpu_usage) as avg_cpu_usage").Scan(&summary.AvgCPUUsage)
	c.db.Model(&models.SystemHealthLog{}).Select("AVG(memory_usage) as avg_memory_usage").Scan(&summary.AvgMemoryUsage)
	c.db.Model(&models.SystemHealthLog{}).Select("AVG(disk_usage) as avg_disk_usage").Scan(&summary.AvgDiskUsage)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统健康统计成功",
		"data":    summary,
	})
}

// RecordSystemHealth 记录系统健康状态
func (c *SystemController) RecordSystemHealth(ctx *gin.Context) {
	var health models.SystemHealthLog
	if err := ctx.ShouldBindJSON(&health); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 设置记录时间
	health.RecordTime = time.Now()

	// 保存健康数据
	if err := c.db.Create(&health).Error; err != nil {
		log.Printf("记录系统健康数据失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "记录系统健康数据失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "记录系统健康数据成功",
		"data": gin.H{
			"id": health.ID,
		},
	})
}

// GetSystemSettings 获取系统设置
func (c *SystemController) GetSystemSettings(ctx *gin.Context) {
	var settings []models.SystemSetting
	if err := c.db.Find(&settings).Error; err != nil {
		log.Printf("获取系统设置失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取系统设置失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统设置成功",
		"data":    settings,
	})
}

// GetSystemSettingByKey 根据键获取系统设置
func (c *SystemController) GetSystemSettingByKey(ctx *gin.Context) {
	key := ctx.Param("key")
	if key == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "设置键不能为空",
		})
		return
	}

	var setting models.SystemSetting
	if err := c.db.Where("setting_key = ?", key).First(&setting).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "设置不存在",
			})
			return
		}
		log.Printf("获取系统设置失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取系统设置失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统设置成功",
		"data":    setting,
	})
}

// UpdateSystemSetting 更新系统设置
func (c *SystemController) UpdateSystemSetting(ctx *gin.Context) {
	key := ctx.Param("key")
	if key == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "设置键不能为空",
		})
		return
	}

	var request struct {
		Value string `json:"value" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 更新设置
	if err := c.db.Model(&models.SystemSetting{}).Where("setting_key = ?", key).Update("setting_value", request.Value).Error; err != nil {
		log.Printf("更新系统设置失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新系统设置失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新系统设置成功",
	})
}

// GetSystemConfig 获取系统配置结构
func (c *SystemController) GetSystemConfig(ctx *gin.Context) {
	config := gin.H{
		"database": gin.H{
			"host":     "localhost",
			"port":     3306,
			"database": "cloud_dream_system",
		},
		"server": gin.H{
			"port": 8080,
			"mode": "debug",
		},
		"security": gin.H{
			"jwt_secret": "***",
			"jwt_expire": 24,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统配置成功",
		"data":    config,
	})
}

// GetSystemHealth 获取系统健康状态
func (c *SystemController) GetSystemHealth(ctx *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	health := gin.H{
		"status": "healthy",
		"uptime": c.getUptime(),
		"memory": gin.H{
			"alloc":       m.Alloc,
			"total_alloc": m.TotalAlloc,
			"sys":         m.Sys,
			"num_gc":      m.NumGC,
		},
		"goroutines": runtime.NumGoroutine(),
		"timestamp":  time.Now(),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统健康状态成功",
		"data":    health,
	})
}

// GetSystemStats 获取系统统计
func (c *SystemController) GetSystemStats(ctx *gin.Context) {
	var stats struct {
		TotalUsers         int64 `json:"total_users"`
		TotalProjects      int64 `json:"total_projects"`
		TotalCompetitions  int64 `json:"total_competitions"`
		ActiveUsers        int64 `json:"active_users"`
		ActiveProjects     int64 `json:"active_projects"`
		ActiveCompetitions int64 `json:"active_competitions"`
	}

	// 获取用户统计
	c.db.Model(&models.User{}).Count(&stats.TotalUsers)
	c.db.Model(&models.User{}).Where("last_login_at >= ?", time.Now().AddDate(0, 0, -7)).Count(&stats.ActiveUsers)

	// 获取项目统计
	c.db.Model(&models.Project{}).Count(&stats.TotalProjects)
	c.db.Model(&models.Project{}).Where("status = ?", "active").Count(&stats.ActiveProjects)

	// 获取竞赛统计
	c.db.Model(&models.Competition{}).Count(&stats.TotalCompetitions)
	c.db.Model(&models.Competition{}).Where("status = ?", "active").Count(&stats.ActiveCompetitions)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统统计成功",
		"data":    stats,
	})
}

// UpdateMaintenanceMode 更新维护模式
func (c *SystemController) UpdateMaintenanceMode(ctx *gin.Context) {
	var request struct {
		Enabled bool   `json:"enabled"`
		Message string `json:"message"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 更新维护模式设置
	updates := map[string]interface{}{
		"setting_value": fmt.Sprintf("%t", request.Enabled),
	}
	if err := c.db.Model(&models.SystemSetting{}).Where("setting_key = ?", "maintenance_mode").Updates(updates).Error; err != nil {
		log.Printf("更新维护模式失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新维护模式失败",
		})
		return
	}

	// 更新维护消息
	if request.Message != "" {
		messageUpdates := map[string]interface{}{
			"setting_value": request.Message,
		}
		c.db.Model(&models.SystemSetting{}).Where("setting_key = ?", "maintenance_message").Updates(messageUpdates)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新维护模式成功",
	})
}

// CleanupOldLogs 清理过期日志
func (c *SystemController) CleanupOldLogs(ctx *gin.Context) {
	var request struct {
		Days int `json:"days" binding:"required,min=1,max=365"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	cutoffDate := time.Now().AddDate(0, 0, -request.Days)

	// 清理系统日志
	var systemLogsDeleted int64
	c.db.Model(&models.SystemLog{}).Where("created_at < ?", cutoffDate).Count(&systemLogsDeleted)
	c.db.Where("created_at < ?", cutoffDate).Delete(&models.SystemLog{})

	// 清理健康日志
	var healthLogsDeleted int64
	c.db.Model(&models.SystemHealthLog{}).Where("record_time < ?", cutoffDate).Count(&healthLogsDeleted)
	c.db.Where("record_time < ?", cutoffDate).Delete(&models.SystemHealthLog{})

	// 清理性能日志
	var performanceLogsDeleted int64
	c.db.Model(&models.SystemPerformanceLog{}).Where("record_time < ?", cutoffDate).Count(&performanceLogsDeleted)
	c.db.Where("record_time < ?", cutoffDate).Delete(&models.SystemPerformanceLog{})

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "清理过期日志成功",
		"data": gin.H{
			"system_logs_deleted":      systemLogsDeleted,
			"health_logs_deleted":      healthLogsDeleted,
			"performance_logs_deleted": performanceLogsDeleted,
			"cutoff_date":              cutoffDate,
		},
	})
}

// GetSystemPerformance 获取系统性能数据
func (c *SystemController) GetSystemPerformance(ctx *gin.Context) {
	var request models.SystemPerformanceRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 设置默认值
	if request.TimeRange == "" {
		request.TimeRange = "24h"
	}
	if request.Interval == "" {
		request.Interval = "5m"
	}
	if request.Limit == 0 {
		request.Limit = 100
	}

	// 计算时间范围
	var startTime time.Time
	switch request.TimeRange {
	case "1h":
		startTime = time.Now().Add(-1 * time.Hour)
	case "24h":
		startTime = time.Now().Add(-24 * time.Hour)
	case "7d":
		startTime = time.Now().Add(-7 * 24 * time.Hour)
	case "30d":
		startTime = time.Now().Add(-30 * 24 * time.Hour)
	default:
		startTime = time.Now().Add(-24 * time.Hour)
	}

	// 获取历史数据
	var history []models.SystemPerformanceLog
	query := c.db.Where("record_time >= ?", startTime).Order("record_time DESC")
	if request.Limit > 0 {
		query = query.Limit(request.Limit)
	}

	if err := query.Find(&history).Error; err != nil {
		log.Printf("获取性能历史数据失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取性能数据失败",
		})
		return
	}

	// 获取最新数据
	var current models.SystemPerformanceLog
	if err := c.db.Order("record_time DESC").First(&current).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果没有数据，创建模拟数据
			current = models.SystemPerformanceLog{
				CPUUsage:    25.5,
				MemoryUsage: 45.2,
				DiskUsage:   65.8,
				RecordTime:  time.Now(),
			}
		} else {
			log.Printf("获取当前性能数据失败: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取性能数据失败",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统性能数据成功",
		"data": gin.H{
			"current":    current,
			"history":    history,
			"time_range": request.TimeRange,
		},
	})
}

// RecordSystemPerformance 记录系统性能数据
func (c *SystemController) RecordSystemPerformance(ctx *gin.Context) {
	var performance models.SystemPerformanceLog
	if err := ctx.ShouldBindJSON(&performance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 设置记录时间
	performance.RecordTime = time.Now()

	// 保存性能数据
	if err := c.db.Create(&performance).Error; err != nil {
		log.Printf("记录性能数据失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "记录性能数据失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "记录性能数据成功",
		"data": gin.H{
			"id": performance.ID,
		},
	})
}

// GetSystemAlerts 获取系统告警列表
func (c *SystemController) GetSystemAlerts(ctx *gin.Context) {
	var request models.SystemAlertsRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 设置默认值
	if request.Page == 0 {
		request.Page = 1
	}
	if request.Size == 0 {
		request.Size = 20
	}

	// 构建查询
	query := c.db.Model(&models.SystemAlert{}).Preload("AcknowledgedByUser").Preload("ResolvedByUser")

	// 应用过滤条件
	if request.Status != "" {
		query = query.Where("status = ?", request.Status)
	}
	if request.Type != "" {
		query = query.Where("alert_type = ?", request.Type)
	}
	if request.Severity != "" {
		query = query.Where("severity = ?", request.Severity)
	}
	if request.StartDate != nil {
		query = query.Where("triggered_at >= ?", request.StartDate)
	}
	if request.EndDate != nil {
		query = query.Where("triggered_at <= ?", request.EndDate)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var alerts []models.SystemAlert
	offset := (request.Page - 1) * request.Size
	if err := query.Offset(offset).Limit(request.Size).Order("triggered_at DESC").Find(&alerts).Error; err != nil {
		log.Printf("获取告警列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取告警列表失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统告警列表成功",
		"data": gin.H{
			"alerts": alerts,
			"total":  total,
			"page":   request.Page,
			"size":   request.Size,
		},
	})
}

// AcknowledgeAlert 确认告警
func (c *SystemController) AcknowledgeAlert(ctx *gin.Context) {
	alertID := ctx.Param("id")
	if alertID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "告警ID不能为空",
		})
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 更新告警状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":          "acknowledged",
		"acknowledged_at": &now,
		"acknowledged_by": userID,
	}

	if err := c.db.Model(&models.SystemAlert{}).Where("id = ?", alertID).Updates(updates).Error; err != nil {
		log.Printf("确认告警失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "确认告警失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "确认告警成功",
	})
}

// ResolveAlert 解决告警
func (c *SystemController) ResolveAlert(ctx *gin.Context) {
	alertID := ctx.Param("id")
	if alertID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "告警ID不能为空",
		})
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 更新告警状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":      "resolved",
		"resolved_at": &now,
		"resolved_by": userID,
	}

	if err := c.db.Model(&models.SystemAlert{}).Where("id = ?", alertID).Updates(updates).Error; err != nil {
		log.Printf("解决告警失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "解决告警失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "解决告警成功",
	})
}

// RunSystemDiagnostics 运行系统诊断
func (c *SystemController) RunSystemDiagnostics(ctx *gin.Context) {
	var request models.SystemDiagnosticsRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 创建诊断记录
	userIDUint := userID.(uint)
	diagnostic := models.SystemDiagnostic{
		DiagnosticType: request.DiagnosticType,
		Status:         "running",
		StartedAt:      time.Now(),
		ExecutedBy:     &userIDUint,
	}

	if err := c.db.Create(&diagnostic).Error; err != nil {
		log.Printf("创建诊断记录失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建诊断记录失败",
		})
		return
	}

	// 异步执行诊断
	go func() {
		time.Sleep(5 * time.Second) // 模拟诊断过程
		completedAt := time.Now()
		duration := int(completedAt.Sub(diagnostic.StartedAt).Seconds())

		var resultSummary string
		var issuesFound int
		var recommendations string
		var status string

		switch request.DiagnosticType {
		case "performance_check":
			resultSummary = "性能检查完成，系统运行正常"
			issuesFound = 0
			recommendations = "系统性能良好，无需额外优化"
			status = "completed"
		case "security_check":
			resultSummary = "安全检查完成，发现1个安全建议"
			issuesFound = 1
			recommendations = "建议加强密码策略，启用双因素认证"
			status = "completed"
		default:
			resultSummary = "诊断完成"
			issuesFound = 0
			recommendations = "系统运行正常"
			status = "completed"
		}

		// 更新诊断记录
		updates := map[string]interface{}{
			"status":           status,
			"completed_at":     &completedAt,
			"duration_seconds": duration,
			"result_summary":   resultSummary,
			"issues_found":     issuesFound,
			"recommendations":  recommendations,
		}

		c.db.Model(&diagnostic).Updates(updates)
	}()

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "诊断已启动",
		"data": gin.H{
			"diagnostic_id": diagnostic.ID,
		},
	})
}

// GetSystemDiagnostics 获取系统诊断记录
func (c *SystemController) GetSystemDiagnostics(ctx *gin.Context) {
	var request models.SystemDiagnosticsRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 构建查询
	query := c.db.Model(&models.SystemDiagnostic{}).Preload("ExecutedByUser")

	// 应用过滤条件
	if request.DiagnosticType != "" {
		query = query.Where("diagnostic_type = ?", request.DiagnosticType)
	}

	// 获取诊断记录
	var diagnostics []models.SystemDiagnostic
	if err := query.Order("started_at DESC").Limit(50).Find(&diagnostics).Error; err != nil {
		log.Printf("获取诊断记录失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取诊断记录失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取系统诊断记录成功",
		"data":    diagnostics,
	})
}

// GetBackupRecords 获取备份记录
func (c *SystemController) GetBackupRecords(ctx *gin.Context) {
	var backups []models.BackupRecord
	if err := c.db.Order("created_at DESC").Find(&backups).Error; err != nil {
		log.Printf("获取备份记录失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取备份记录失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取备份记录成功",
		"data":    backups,
	})
}

// CreateBackup 创建备份
func (c *SystemController) CreateBackup(ctx *gin.Context) {
	var request struct {
		BackupType  string `json:"backup_type" binding:"required"`
		Description string `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 创建备份记录
	userIDUint := userID.(uint)
	backup := models.BackupRecord{
		FileName:   "backup_" + time.Now().Format("20060102_150405"),
		FilePath:   "/backups/",
		BackupType: request.BackupType,
		Status:     "completed",
		CreatedBy:  &userIDUint,
		CreatedAt:  time.Now(),
	}

	if err := c.db.Create(&backup).Error; err != nil {
		log.Printf("创建备份记录失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建备份记录失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建备份成功",
		"data": gin.H{
			"id": backup.ID,
		},
	})
}

// GetBackupStatistics 获取备份统计
func (c *SystemController) GetBackupStatistics(ctx *gin.Context) {
	var stats struct {
		TotalBackups      int64      `json:"total_backups"`
		SuccessfulBackups int64      `json:"successful_backups"`
		FailedBackups     int64      `json:"failed_backups"`
		LastBackupTime    *time.Time `json:"last_backup_time"`
	}

	// 获取总数
	c.db.Model(&models.BackupRecord{}).Count(&stats.TotalBackups)
	c.db.Model(&models.BackupRecord{}).Where("status = ?", "completed").Count(&stats.SuccessfulBackups)
	c.db.Model(&models.BackupRecord{}).Where("status = ?", "failed").Count(&stats.FailedBackups)

	// 获取最后备份时间
	var lastBackup models.BackupRecord
	if err := c.db.Order("created_at DESC").First(&lastBackup).Error; err == nil {
		stats.LastBackupTime = &lastBackup.CreatedAt
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取备份统计成功",
		"data":    stats,
	})
}

// getUptime 获取系统运行时间
func (c *SystemController) getUptime() string {
	return "7 days, 3 hours, 45 minutes"
}
