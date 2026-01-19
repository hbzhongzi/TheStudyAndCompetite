package models

import (
	"time"

	"gorm.io/datatypes"
)

// SystemLog 系统日志表
type SystemLog struct {
	ID         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	LogType    string     `json:"log_type" gorm:"type:enum('info','warning','error','debug','security');default:info;comment:日志类型"`
	Operation  string     `json:"operation" gorm:"type:varchar(100);comment:操作名称"`
	Status     string     `json:"status" gorm:"type:enum('success','failed','pending');default:success;comment:执行状态"`
	UserID     *uint      `json:"user_id" gorm:"comment:操作用户ID"`
	Action     string     `json:"action" gorm:"type:varchar(255);not null;comment:执行动作"`
	Details    string     `json:"details" gorm:"type:text;comment:详细内容"`
	IPAddress  string     `json:"ip_address" gorm:"type:varchar(50);comment:IP地址"`
	UserAgent  string     `json:"user_agent" gorm:"type:text;comment:用户代理"`
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	ExpireTime *time.Time `json:"expire_time" gorm:"comment:日志过期时间"`

	// 关联关系
	User *User `json:"user" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (SystemLog) TableName() string {
	return "system_logs"
}

// SystemHealthLog 系统健康监控日志表
type SystemHealthLog struct {
	ID                uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CPUUsage          *float64  `json:"cpu_usage" gorm:"type:decimal(5,2);comment:CPU使用率(%)"`
	MemoryUsage       *float64  `json:"memory_usage" gorm:"type:decimal(5,2);comment:内存使用率(%)"`
	DiskUsage         *float64  `json:"disk_usage" gorm:"type:decimal(5,2);comment:磁盘使用率(%)"`
	DBStatus          string    `json:"db_status" gorm:"type:enum('healthy','warning','error','offline');default:healthy;comment:数据库状态"`
	DBConnectionCount int       `json:"db_connection_count" gorm:"default:0;comment:数据库连接数"`
	ActiveUsers       int       `json:"active_users" gorm:"default:0;comment:活跃用户数"`
	RequestCount      int       `json:"request_count" gorm:"default:0;comment:请求数量"`
	ErrorCount        int       `json:"error_count" gorm:"default:0;comment:错误数量"`
	ResponseTimeAvg   *float64  `json:"response_time_avg" gorm:"type:decimal(10,3);comment:平均响应时间(ms)"`
	RecordTime        time.Time `json:"record_time" gorm:"autoCreateTime;comment:记录时间"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
}

// TableName 指定表名
func (SystemHealthLog) TableName() string {
	return "system_health_logs"
}

// SystemSetting 系统配置表
type SystemSetting struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SettingKey   string    `json:"setting_key" gorm:"type:varchar(100);unique;not null;comment:配置键"`
	SettingValue string    `json:"setting_value" gorm:"type:text;comment:配置值"`
	Description  string    `json:"description" gorm:"type:varchar(255);comment:配置描述"`
	Category     string    `json:"category" gorm:"type:varchar(50);default:general;comment:配置分类"`
	IsPublic     bool      `json:"is_public" gorm:"default:false;comment:是否公开"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	UpdateTime   time.Time `json:"update_time" gorm:"autoUpdateTime;comment:更新时间"`
	UpdatedBy    *uint     `json:"updated_by" gorm:"comment:修改人ID"`

	// 关联关系
	UpdatedByUser *User `json:"updated_by_user" gorm:"foreignKey:UpdatedBy"`
}

// TableName 指定表名
func (SystemSetting) TableName() string {
	return "system_settings"
}

// BackupRecord 备份记录表
type BackupRecord struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	FileName     string     `json:"file_name" gorm:"type:varchar(255);not null;comment:文件名"`
	FilePath     string     `json:"file_path" gorm:"type:varchar(500);not null;comment:文件路径"`
	FileSize     *int64     `json:"file_size" gorm:"comment:文件大小"`
	BackupType   string     `json:"backup_type" gorm:"type:enum('full','incremental','manual');default:manual;comment:备份类型"`
	Status       string     `json:"status" gorm:"type:enum('pending','in_progress','success','failed','cancelled');default:pending;comment:备份状态"`
	ErrorMessage string     `json:"error_message" gorm:"type:text;comment:错误信息"`
	CreatedBy    *uint      `json:"created_by" gorm:"comment:创建者ID"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	StartedAt    *time.Time `json:"started_at" gorm:"comment:开始时间"`
	CompletedAt  *time.Time `json:"completed_at" gorm:"comment:完成时间"`

	// 关联关系
	CreatedByUser *User `json:"created_by_user" gorm:"foreignKey:CreatedBy"`
}

// TableName 指定表名
func (BackupRecord) TableName() string {
	return "backup_records"
}

// SystemStats 系统统计信息
type SystemStats struct {
	TotalUsers        int64 `json:"total_users"`
	ActiveUsers       int64 `json:"active_users"`
	TotalProjects     int64 `json:"total_projects"`
	TotalCompetitions int64 `json:"total_competitions"`
	Logs24h           int64 `json:"logs_24h"`
	Backups7d         int64 `json:"backups_7d"`
	HealthLogs24h     int64 `json:"health_logs_24h"`
	ErrorLogs24h      int64 `json:"error_logs_24h"`
}

// SystemLogSummary 系统日志统计
type SystemLogSummary struct {
	LogDate      string `json:"log_date"`
	LogType      string `json:"log_type"`
	Operation    string `json:"operation"`
	Status       string `json:"status"`
	ActionCount  int64  `json:"action_count"`
	UniqueUsers  int64  `json:"unique_users"`
	UniqueIPs    int64  `json:"unique_ips"`
	FailedCount  int64  `json:"failed_count"`
	SuccessCount int64  `json:"success_count"`
}

// SystemHealthSummary 系统健康统计
type SystemHealthSummary struct {
	HealthDate      string  `json:"health_date"`
	AvgCPUUsage     float64 `json:"avg_cpu_usage"`
	AvgMemoryUsage  float64 `json:"avg_memory_usage"`
	AvgDiskUsage    float64 `json:"avg_disk_usage"`
	AvgResponseTime float64 `json:"avg_response_time"`
	MaxActiveUsers  int     `json:"max_active_users"`
	TotalRequests   int     `json:"total_requests"`
	TotalErrors     int     `json:"total_errors"`
	DBErrorCount    int     `json:"db_error_count"`
	DBWarningCount  int     `json:"db_warning_count"`
}

// BackupStatistics 备份统计信息
type BackupStatistics struct {
	BackupType         string     `json:"backup_type"`
	Status             string     `json:"status"`
	Count              int64      `json:"count"`
	AvgDurationSeconds float64    `json:"avg_duration_seconds"`
	TotalSizeBytes     int64      `json:"total_size_bytes"`
	LastBackup         *time.Time `json:"last_backup"`
}

// SystemSettingRequest 系统配置请求
type SystemSettingRequest struct {
	SettingKey   string `json:"setting_key" binding:"required"`
	SettingValue string `json:"setting_value"`
	Description  string `json:"description"`
	Category     string `json:"category"`
	IsPublic     bool   `json:"is_public"`
}

// SystemSettingUpdateRequest 系统配置更新请求
type SystemSettingUpdateRequest struct {
	SettingValue string `json:"setting_value"`
	Description  string `json:"description"`
	Category     string `json:"category"`
	IsPublic     bool   `json:"is_public"`
}

// BackupRecordRequest 备份记录请求
type BackupRecordRequest struct {
	FileName   string `json:"file_name" binding:"required"`
	FilePath   string `json:"file_path" binding:"required"`
	BackupType string `json:"backup_type" binding:"omitempty,oneof=full incremental manual"`
}

// SystemLogRequest 系统日志请求
type SystemLogRequest struct {
	LogType    string     `json:"log_type" binding:"omitempty,oneof=info warning error debug security"`
	Operation  string     `json:"operation"`
	Status     string     `json:"status" binding:"omitempty,oneof=success failed pending"`
	Action     string     `json:"action" binding:"required"`
	Details    string     `json:"details"`
	IPAddress  string     `json:"ip_address"`
	UserAgent  string     `json:"user_agent"`
	ExpireTime *time.Time `json:"expire_time"`
}

// SystemHealthLogRequest 系统健康日志请求
type SystemHealthLogRequest struct {
	CPUUsage          *float64 `json:"cpu_usage"`
	MemoryUsage       *float64 `json:"memory_usage"`
	DiskUsage         *float64 `json:"disk_usage"`
	DBStatus          string   `json:"db_status" binding:"omitempty,oneof=healthy warning error offline"`
	DBConnectionCount int      `json:"db_connection_count"`
	ActiveUsers       int      `json:"active_users"`
	RequestCount      int      `json:"request_count"`
	ErrorCount        int      `json:"error_count"`
	ResponseTimeAvg   *float64 `json:"response_time_avg"`
}

// SystemLogResponse 系统日志响应
type SystemLogResponse struct {
	ID         uint       `json:"id"`
	LogType    string     `json:"log_type"`
	Operation  string     `json:"operation"`
	Status     string     `json:"status"`
	UserID     *uint      `json:"user_id"`
	Action     string     `json:"action"`
	Details    string     `json:"details"`
	IPAddress  string     `json:"ip_address"`
	UserAgent  string     `json:"user_agent"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpireTime *time.Time `json:"expire_time"`

	// 关联数据
	User *CompetitionUserResponse `json:"user"`
}

// SystemHealthLogResponse 系统健康日志响应
type SystemHealthLogResponse struct {
	ID                uint      `json:"id"`
	CPUUsage          *float64  `json:"cpu_usage"`
	MemoryUsage       *float64  `json:"memory_usage"`
	DiskUsage         *float64  `json:"disk_usage"`
	DBStatus          string    `json:"db_status"`
	DBConnectionCount int       `json:"db_connection_count"`
	ActiveUsers       int       `json:"active_users"`
	RequestCount      int       `json:"request_count"`
	ErrorCount        int       `json:"error_count"`
	ResponseTimeAvg   *float64  `json:"response_time_avg"`
	RecordTime        time.Time `json:"record_time"`
	CreatedAt         time.Time `json:"created_at"`
}

// SystemSettingResponse 系统配置响应
type SystemSettingResponse struct {
	ID           uint      `json:"id"`
	SettingKey   string    `json:"setting_key"`
	SettingValue string    `json:"setting_value"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	IsPublic     bool      `json:"is_public"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdateTime   time.Time `json:"update_time"`
	UpdatedBy    *uint     `json:"updated_by"`

	// 关联数据
	UpdatedByUser *CompetitionUserResponse `json:"updated_by_user"`
}

// BackupRecordResponse 备份记录响应
type BackupRecordResponse struct {
	ID           uint       `json:"id"`
	FileName     string     `json:"file_name"`
	FilePath     string     `json:"file_path"`
	FileSize     *int64     `json:"file_size"`
	BackupType   string     `json:"backup_type"`
	Status       string     `json:"status"`
	ErrorMessage string     `json:"error_message"`
	CreatedBy    *uint      `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	StartedAt    *time.Time `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`

	// 关联数据
	CreatedByUser *CompetitionUserResponse `json:"created_by_user"`
}

// SystemHealthResponse 系统健康检查响应
type SystemHealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Uptime    string                 `json:"uptime"`
	Database  DatabaseHealthInfo     `json:"database"`
	Services  map[string]ServiceInfo `json:"services"`
	Metrics   SystemMetrics          `json:"metrics"`
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	CPUUsage        float64 `json:"cpu_usage"`
	MemoryUsage     float64 `json:"memory_usage"`
	DiskUsage       float64 `json:"disk_usage"`
	ActiveUsers     int     `json:"active_users"`
	RequestCount    int     `json:"request_count"`
	ErrorCount      int     `json:"error_count"`
	ResponseTimeAvg float64 `json:"response_time_avg"`
}

// DatabaseHealthInfo 数据库健康信息
type DatabaseHealthInfo struct {
	Status      string `json:"status"`
	Version     string `json:"version"`
	Connections int    `json:"connections"`
	Size        string `json:"size"`
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Status       string    `json:"status"`
	ResponseTime int64     `json:"response_time_ms"`
	LastCheck    time.Time `json:"last_check"`
}

// SystemMaintenanceRequest 系统维护请求
type SystemMaintenanceRequest struct {
	MaintenanceMode bool   `json:"maintenance_mode"`
	Message         string `json:"message"`
}

// SystemLogFilter 系统日志过滤条件
type SystemLogFilter struct {
	LogType   string     `json:"log_type"`
	Operation string     `json:"operation"`
	Status    string     `json:"status"`
	UserID    *uint      `json:"user_id"`
	Action    string     `json:"action"`
	IPAddress string     `json:"ip_address"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Page      int        `json:"page"`
	Size      int        `json:"size"`
}

// SystemHealthLogFilter 系统健康日志过滤条件
type SystemHealthLogFilter struct {
	DBStatus  string     `json:"db_status"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Page      int        `json:"page"`
	Size      int        `json:"size"`
}

// BackupRecordFilter 备份记录过滤条件
type BackupRecordFilter struct {
	BackupType string     `json:"backup_type"`
	Status     string     `json:"status"`
	StartDate  *time.Time `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	Page       int        `json:"page"`
	Size       int        `json:"size"`
}

// SystemConfig 系统配置结构
type SystemConfig struct {
	General      GeneralConfig      `json:"general"`
	Upload       UploadConfig       `json:"upload"`
	Security     SecurityConfig     `json:"security"`
	Backup       BackupConfig       `json:"backup"`
	Email        EmailConfig        `json:"email"`
	Notification NotificationConfig `json:"notification"`
	System       SystemConfigInfo   `json:"system"`
	Monitoring   MonitoringConfig   `json:"monitoring"`
}

// GeneralConfig 通用配置
type GeneralConfig struct {
	SystemName    string `json:"system_name"`
	SystemVersion string `json:"system_version"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxFileSize      int64  `json:"max_file_size"`
	AllowedFileTypes string `json:"allowed_file_types"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	SessionTimeout   int `json:"session_timeout"`
	MaxLoginAttempts int `json:"max_login_attempts"`
}

// BackupConfig 备份配置
type BackupConfig struct {
	RetentionDays     int    `json:"backup_retention_days"`
	AutoBackupEnabled bool   `json:"auto_backup_enabled"`
	BackupSchedule    string `json:"backup_schedule"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
}

// NotificationConfig 通知配置
type NotificationConfig struct {
	EmailNotifications bool `json:"email_notifications"`
}

// SystemConfigInfo 系统配置信息
type SystemConfigInfo struct {
	MaintenanceMode    bool   `json:"maintenance_mode"`
	MaintenanceMessage string `json:"maintenance_message"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	HealthMonitorEnabled       bool   `json:"health_monitor_enabled"`
	HealthMonitorInterval      int    `json:"health_monitor_interval"`
	MaxHealthLogs              int    `json:"max_health_logs"`
	SystemAlertEmail           string `json:"system_alert_email"`
	PerformanceThresholdCPU    int    `json:"performance_threshold_cpu"`
	PerformanceThresholdMemory int    `json:"performance_threshold_memory"`
	PerformanceThresholdDisk   int    `json:"performance_threshold_disk"`
	LogRetentionDays           int    `json:"log_retention_days"`
}

// ==================== 新增系统管理模型 ====================

// SystemPerformanceLog 系统性能监控日志
type SystemPerformanceLog struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	CPUUsage          float64   `json:"cpu_usage" gorm:"type:decimal(5,2);not null;comment:CPU使用率(%)"`
	MemoryUsage       float64   `json:"memory_usage" gorm:"type:decimal(5,2);not null;comment:内存使用率(%)"`
	DiskUsage         float64   `json:"disk_usage" gorm:"type:decimal(5,2);not null;comment:磁盘使用率(%)"`
	NetworkIn         *float64  `json:"network_in" gorm:"type:decimal(10,2);comment:网络入流量(MB/s)"`
	NetworkOut        *float64  `json:"network_out" gorm:"type:decimal(10,2);comment:网络出流量(MB/s)"`
	ActiveConnections *int      `json:"active_connections" gorm:"comment:活跃连接数"`
	ResponseTimeAvg   *float64  `json:"response_time_avg" gorm:"type:decimal(10,3);comment:平均响应时间(ms)"`
	ErrorRate         *float64  `json:"error_rate" gorm:"type:decimal(5,2);comment:错误率(%)"`
	RecordTime        time.Time `json:"record_time" gorm:"default:CURRENT_TIMESTAMP;comment:记录时间"`
}

// SystemAlert 系统告警
type SystemAlert struct {
	ID             uint            `json:"id" gorm:"primaryKey"`
	AlertType      string          `json:"alert_type" gorm:"type:enum('cpu_high','memory_high','disk_full','error_rate_high','response_time_slow','backup_failed','security_breach','database_error','service_down');not null;comment:告警类型"`
	Severity       string          `json:"severity" gorm:"type:enum('low','medium','high','critical');not null;comment:告警级别"`
	Title          string          `json:"title" gorm:"type:varchar(200);not null;comment:告警标题"`
	Message        string          `json:"message" gorm:"type:text;not null;comment:告警消息"`
	Status         string          `json:"status" gorm:"type:enum('active','acknowledged','resolved');default:active;comment:告警状态"`
	TriggeredAt    time.Time       `json:"triggered_at" gorm:"default:CURRENT_TIMESTAMP;comment:触发时间"`
	AcknowledgedAt *time.Time      `json:"acknowledged_at" gorm:"comment:确认时间"`
	ResolvedAt     *time.Time      `json:"resolved_at" gorm:"comment:解决时间"`
	AcknowledgedBy *uint           `json:"acknowledged_by" gorm:"comment:确认人ID"`
	ResolvedBy     *uint           `json:"resolved_by" gorm:"comment:解决人ID"`
	Metadata       *datatypes.JSON `json:"metadata" gorm:"comment:告警元数据"`

	// 关联
	AcknowledgedByUser *User `json:"acknowledged_by_user" gorm:"foreignKey:AcknowledgedBy"`
	ResolvedByUser     *User `json:"resolved_by_user" gorm:"foreignKey:ResolvedBy"`
}

// SystemDiagnostic 系统诊断记录
type SystemDiagnostic struct {
	ID              uint            `json:"id" gorm:"primaryKey"`
	DiagnosticType  string          `json:"diagnostic_type" gorm:"type:enum('system_check','database_check','performance_check','security_check','backup_check','full_check');not null;comment:诊断类型"`
	Status          string          `json:"status" gorm:"type:enum('running','completed','failed');default:running;comment:诊断状态"`
	StartedAt       time.Time       `json:"started_at" gorm:"default:CURRENT_TIMESTAMP;comment:开始时间"`
	CompletedAt     *time.Time      `json:"completed_at" gorm:"comment:完成时间"`
	DurationSeconds *int            `json:"duration_seconds" gorm:"comment:执行时长(秒)"`
	ResultSummary   *string         `json:"result_summary" gorm:"type:text;comment:结果摘要"`
	DetailedReport  *datatypes.JSON `json:"detailed_report" gorm:"comment:详细报告"`
	IssuesFound     int             `json:"issues_found" gorm:"default:0;comment:发现的问题数量"`
	Recommendations *string         `json:"recommendations" gorm:"type:text;comment:建议"`
	ExecutedBy      *uint           `json:"executed_by" gorm:"comment:执行人ID"`

	// 关联
	ExecutedByUser *User `json:"executed_by_user" gorm:"foreignKey:ExecutedBy"`
}

// SystemConfigTemplate 系统配置模板
type SystemConfigTemplate struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	TemplateName string         `json:"template_name" gorm:"type:varchar(100);not null;unique;comment:模板名称"`
	TemplateType string         `json:"template_type" gorm:"type:enum('basic','performance','security','backup','monitoring','custom');not null;comment:模板类型"`
	Description  *string        `json:"description" gorm:"type:text;comment:模板描述"`
	ConfigData   datatypes.JSON `json:"config_data" gorm:"not null;comment:配置数据"`
	IsDefault    bool           `json:"is_default" gorm:"default:false;comment:是否默认模板"`
	CreatedAt    time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
	CreatedBy    *uint          `json:"created_by" gorm:"comment:创建人ID"`

	// 关联
	CreatedByUser *User `json:"created_by_user" gorm:"foreignKey:CreatedBy"`
}

// SystemScheduledTask 系统任务调度
type SystemScheduledTask struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	TaskName       string     `json:"task_name" gorm:"type:varchar(100);not null;unique;comment:任务名称"`
	TaskType       string     `json:"task_type" gorm:"type:enum('backup','cleanup','health_check','performance_monitor','alert_check','report_generation');not null;comment:任务类型"`
	CronExpression string     `json:"cron_expression" gorm:"type:varchar(100);not null;comment:Cron表达式"`
	IsActive       bool       `json:"is_active" gorm:"default:true;comment:是否激活"`
	LastRunAt      *time.Time `json:"last_run_at" gorm:"comment:上次运行时间"`
	NextRunAt      *time.Time `json:"next_run_at" gorm:"comment:下次运行时间"`
	RunCount       int        `json:"run_count" gorm:"default:0;comment:运行次数"`
	SuccessCount   int        `json:"success_count" gorm:"default:0;comment:成功次数"`
	FailureCount   int        `json:"failure_count" gorm:"default:0;comment:失败次数"`
	LastStatus     *string    `json:"last_status" gorm:"type:enum('success','failed','running');comment:上次运行状态"`
	LastError      *string    `json:"last_error" gorm:"type:text;comment:上次错误信息"`
	CreatedAt      time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
	CreatedBy      *uint      `json:"created_by" gorm:"comment:创建人ID"`

	// 关联
	CreatedByUser *User `json:"created_by_user" gorm:"foreignKey:CreatedBy"`
}

// ==================== 响应结构体 ====================

// SystemPerformanceResponse 系统性能响应
type SystemPerformanceResponse struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Data    SystemPerformanceData `json:"data"`
}

// SystemPerformanceData 系统性能数据
type SystemPerformanceData struct {
	Current SystemPerformanceLog   `json:"current"`
	History []SystemPerformanceLog `json:"history"`
	Summary PerformanceSummary     `json:"summary"`
	Trends  PerformanceTrends      `json:"trends"`
}

// PerformanceSummary 性能摘要
type PerformanceSummary struct {
	AvgCPUUsage     float64 `json:"avg_cpu_usage"`
	MaxCPUUsage     float64 `json:"max_cpu_usage"`
	AvgMemoryUsage  float64 `json:"avg_memory_usage"`
	MaxMemoryUsage  float64 `json:"max_memory_usage"`
	AvgDiskUsage    float64 `json:"avg_disk_usage"`
	MaxDiskUsage    float64 `json:"max_disk_usage"`
	AvgResponseTime float64 `json:"avg_response_time"`
	MaxResponseTime float64 `json:"max_response_time"`
	AvgErrorRate    float64 `json:"avg_error_rate"`
	MaxErrorRate    float64 `json:"max_error_rate"`
	TotalRecords    int     `json:"total_records"`
	TimeRange       string  `json:"time_range"`
}

// PerformanceTrends 性能趋势
type PerformanceTrends struct {
	CPUUsageTrend     []TrendPoint `json:"cpu_usage_trend"`
	MemoryUsageTrend  []TrendPoint `json:"memory_usage_trend"`
	DiskUsageTrend    []TrendPoint `json:"disk_usage_trend"`
	ResponseTimeTrend []TrendPoint `json:"response_time_trend"`
	ErrorRateTrend    []TrendPoint `json:"error_rate_trend"`
}

// TrendPoint 趋势点
type TrendPoint struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

// SystemAlertsResponse 系统告警响应
type SystemAlertsResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    SystemAlertsData `json:"data"`
}

// SystemAlertsData 系统告警数据
type SystemAlertsData struct {
	Alerts     []SystemAlert    `json:"alerts"`
	Summary    AlertsSummary    `json:"summary"`
	Statistics AlertsStatistics `json:"statistics"`
}

// AlertsSummary 告警摘要
type AlertsSummary struct {
	TotalAlerts    int `json:"total_alerts"`
	ActiveAlerts   int `json:"active_alerts"`
	CriticalAlerts int `json:"critical_alerts"`
	HighAlerts     int `json:"high_alerts"`
	MediumAlerts   int `json:"medium_alerts"`
	LowAlerts      int `json:"low_alerts"`
}

// AlertsStatistics 告警统计
type AlertsStatistics struct {
	ByType     map[string]int  `json:"by_type"`
	BySeverity map[string]int  `json:"by_severity"`
	ByStatus   map[string]int  `json:"by_status"`
	ByDay      []DayAlertCount `json:"by_day"`
}

// DayAlertCount 每日告警数量
type DayAlertCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// SystemDiagnosticsResponse 系统诊断响应
type SystemDiagnosticsResponse struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Data    SystemDiagnosticsData `json:"data"`
}

// SystemDiagnosticsData 系统诊断数据
type SystemDiagnosticsData struct {
	Diagnostics []SystemDiagnostic `json:"diagnostics"`
	Summary     DiagnosticSummary  `json:"summary"`
	Report      DiagnosticReport   `json:"report"`
}

// DiagnosticSummary 诊断摘要
type DiagnosticSummary struct {
	TotalDiagnostics int `json:"total_diagnostics"`
	Completed        int `json:"completed"`
	Running          int `json:"running"`
	Failed           int `json:"failed"`
	TotalIssues      int `json:"total_issues"`
}

// DiagnosticReport 诊断报告
type DiagnosticReport struct {
	OverallStatus    string                 `json:"overall_status"`
	IssuesFound      int                    `json:"issues_found"`
	Recommendations  []string               `json:"recommendations"`
	ComponentStatus  map[string]string      `json:"component_status"`
	PerformanceScore float64                `json:"performance_score"`
	SecurityScore    float64                `json:"security_score"`
	Details          map[string]interface{} `json:"details"`
}

// ==================== 请求结构体 ====================

// SystemPerformanceRequest 系统性能请求
type SystemPerformanceRequest struct {
	TimeRange string `json:"time_range" form:"time_range"` // 1h, 24h, 7d, 30d
	Interval  string `json:"interval" form:"interval"`     // 1m, 5m, 15m, 1h
	Limit     int    `json:"limit" form:"limit"`
}

// SystemAlertsRequest 系统告警请求
type SystemAlertsRequest struct {
	Status    string     `json:"status" form:"status"`
	Type      string     `json:"type" form:"type"`
	Severity  string     `json:"severity" form:"severity"`
	StartDate *time.Time `json:"start_date" form:"start_date"`
	EndDate   *time.Time `json:"end_date" form:"end_date"`
	Page      int        `json:"page" form:"page"`
	Size      int        `json:"size" form:"size"`
}

// SystemDiagnosticsRequest 系统诊断请求
type SystemDiagnosticsRequest struct {
	DiagnosticType string `json:"diagnostic_type" form:"diagnostic_type"`
	IncludeDetails bool   `json:"include_details" form:"include_details"`
}

// AlertActionRequest 告警操作请求
type AlertActionRequest struct {
	Action  string `json:"action" binding:"required"` // acknowledge, resolve
	Comment string `json:"comment"`
}

// ==================== 过滤器结构体 ====================

// SystemPerformanceFilter 系统性能过滤器
type SystemPerformanceFilter struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	MinCPU    *float64   `json:"min_cpu"`
	MaxCPU    *float64   `json:"max_cpu"`
	MinMemory *float64   `json:"min_memory"`
	MaxMemory *float64   `json:"max_memory"`
	MinDisk   *float64   `json:"min_disk"`
	MaxDisk   *float64   `json:"max_disk"`
	Page      int        `json:"page"`
	Size      int        `json:"size"`
}

// SystemAlertsFilter 系统告警过滤器
type SystemAlertsFilter struct {
	AlertType string     `json:"alert_type"`
	Severity  string     `json:"severity"`
	Status    string     `json:"status"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Page      int        `json:"page"`
	Size      int        `json:"size"`
}

// SystemDiagnosticsFilter 系统诊断过滤器
type SystemDiagnosticsFilter struct {
	DiagnosticType string     `json:"diagnostic_type"`
	Status         string     `json:"status"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	ExecutedBy     *uint      `json:"executed_by"`
	Page           int        `json:"page"`
	Size           int        `json:"size"`
}
