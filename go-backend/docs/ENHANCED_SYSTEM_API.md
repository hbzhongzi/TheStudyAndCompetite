# 云梦高校科研竞赛管理系统 - 增强系统API文档

## 概述

本文档描述了系统增强后的API接口，包含根据用户建议新增的系统监控、日志管理、健康检查等功能。

## 数据库增强概览

### 1. ✅ 系统日志表 (system_logs) 增强
- **新增字段**: `log_type`、`operation`、`status`、`expire_time`
- **用途**: 日志类型分类、操作名称、执行状态、自动过期时间

### 2. ✅ 系统设置表 (system_settings) 增强
- **新增字段**: `description`、`update_time`、`updated_by`
- **用途**: 配置项说明、更新时间、修改人追踪

### 3. ✅ 备份记录表 (backup_records) 增强
- **完善字段**: `status`、`file_path`、`file_size`
- **用途**: 备份状态管理、文件路径、文件大小

### 4. ✅ 新增系统健康监控表 (system_health_logs)
- **字段**: `cpu_usage`、`memory_usage`、`disk_usage`、`db_status`、`record_time`等
- **用途**: 存储系统健康监控历史记录

---

## 新增API接口

### 1. 系统健康监控API

#### 1.1 获取系统健康日志
```http
GET /api/admin/logs/health
Authorization: Bearer <token>
```

**查询参数**:
- `page`: 页码 (默认: 1)
- `size`: 每页数量 (默认: 20)
- `db_status`: 数据库状态筛选 (可选: healthy, warning, error, offline)
- `start_date`: 开始日期 (可选: YYYY-MM-DD)
- `end_date`: 结束日期 (可选: YYYY-MM-DD)

**响应示例**:
```json
{
  "code": 200,
  "message": "获取系统健康日志成功",
  "data": {
    "list": [
      {
        "id": 1,
        "cpu_usage": 45.2,
        "memory_usage": 67.8,
        "disk_usage": 72.1,
        "db_status": "healthy",
        "db_connection_count": 25,
        "active_users": 45,
        "request_count": 1234,
        "error_count": 5,
        "response_time_avg": 125.5,
        "record_time": "2024-01-15T14:30:00Z",
        "created_at": "2024-01-15T14:30:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 20
  }
}
```

#### 1.2 获取系统健康统计
```http
GET /api/admin/logs/health/summary
Authorization: Bearer <token>
```

**查询参数**:
- `days`: 统计天数 (默认: 7)

**响应示例**:
```json
{
  "code": 200,
  "message": "获取系统健康统计成功",
  "data": [
    {
      "health_date": "2024-01-15",
      "avg_cpu_usage": 42.5,
      "avg_memory_usage": 65.3,
      "avg_disk_usage": 71.8,
      "avg_response_time": 128.7,
      "max_active_users": 67,
      "total_requests": 8567,
      "total_errors": 23,
      "db_error_count": 0,
      "db_warning_count": 2
    }
  ]
}
```

#### 1.3 记录系统健康状态
```http
POST /api/admin/logs/health
Authorization: Bearer <token>
Content-Type: application/json

{
  "cpu_usage": 45.2,
  "memory_usage": 67.8,
  "disk_usage": 72.1,
  "db_status": "healthy",
  "db_connection_count": 25,
  "active_users": 45,
  "request_count": 1234,
  "error_count": 5,
  "response_time_avg": 125.5
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "记录系统健康状态成功",
  "data": {
    "id": 1
  }
}
```

### 2. 增强的系统日志API

#### 2.1 获取系统日志（增强版）
```http
GET /api/admin/logs
Authorization: Bearer <token>
```

**新增查询参数**:
- `log_type`: 日志类型筛选 (可选: info, warning, error, debug, security)
- `operation`: 操作名称筛选 (可选)
- `status`: 执行状态筛选 (可选: success, failed, pending)

**响应示例**:
```json
{
  "code": 200,
  "message": "获取系统日志成功",
  "data": {
    "list": [
      {
        "id": 1,
        "log_type": "security",
        "operation": "user_login",
        "status": "success",
        "user_id": 3,
        "action": "user_login",
        "details": "用户登录成功，用户ID: 3",
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0...",
        "created_at": "2024-01-15T14:30:00Z",
        "expire_time": "2024-04-15T14:30:00Z",
        "user": {
          "id": 3,
          "username": "student001",
          "real_name": "张三",
          "email": "zhangsan@yunmeng.edu.cn"
        }
      }
    ],
    "total": 100,
    "page": 1,
    "size": 20
  }
}
```

#### 2.2 获取系统日志统计（增强版）
```http
GET /api/admin/logs/summary
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取系统日志统计成功",
  "data": [
    {
      "log_date": "2024-01-15",
      "log_type": "security",
      "operation": "user_login",
      "status": "success",
      "action_count": 45,
      "unique_users": 23,
      "unique_ips": 18,
      "failed_count": 2,
      "success_count": 43
    }
  ]
}
```

### 3. 增强的系统设置API

#### 3.1 获取系统设置（增强版）
```http
GET /api/admin/settings
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取系统设置成功",
  "data": [
    {
      "id": 1,
      "setting_key": "system_name",
      "setting_value": "云梦高校学生科研与竞赛项目管理系统",
      "description": "系统名称，用于显示在页面标题和登录页面",
      "category": "general",
      "is_public": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-15T10:30:00Z",
      "update_time": "2024-01-15T10:30:00Z",
      "updated_by": 1,
      "updated_by_user": {
        "id": 1,
        "username": "admin",
        "real_name": "系统管理员",
        "email": "admin@yunmeng.edu.cn"
      }
    }
  ]
}
```

#### 3.2 更新系统设置（增强版）
```http
PUT /api/admin/settings/:key
Authorization: Bearer <token>
Content-Type: application/json

{
  "setting_value": "新的系统名称",
  "description": "更新后的描述",
  "category": "general",
  "is_public": true
}
```

### 4. 增强的系统统计API

#### 4.1 获取系统统计（增强版）
```http
GET /api/admin/stats
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取系统统计成功",
  "data": {
    "total_users": 1234,
    "active_users": 1189,
    "total_projects": 567,
    "total_competitions": 89,
    "logs_24h": 2345,
    "backups_7d": 7,
    "health_logs_24h": 24,
    "error_logs_24h": 12
  }
}
```

### 5. 增强的系统健康检查API

#### 5.1 获取系统健康状态（增强版）
```http
GET /api/admin/health
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取系统健康状态成功",
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-15T14:30:00Z",
    "uptime": "7 days, 3 hours, 45 minutes",
    "database": {
      "status": "healthy",
      "version": "MySQL 8.0",
      "connections": 25,
      "size": "1.2GB"
    },
    "services": {
      "api": {
        "status": "healthy",
        "response_time_ms": 50,
        "last_check": "2024-01-15T14:30:00Z"
      },
      "database": {
        "status": "healthy",
        "response_time_ms": 10,
        "last_check": "2024-01-15T14:30:00Z"
      }
    },
    "metrics": {
      "cpu_usage": 45.2,
      "memory_usage": 67.8,
      "disk_usage": 72.1,
      "active_users": 45,
      "request_count": 1234,
      "error_count": 5,
      "response_time_avg": 125.5
    }
  }
}
```

### 6. 增强的日志清理API

#### 6.1 清理过期日志（增强版）
```http
POST /api/admin/logs/cleanup
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "清理过期日志成功",
  "data": {
    "deleted_system_logs": 1234,
    "deleted_health_logs": 567,
    "retention_days": 90,
    "health_retention_days": 30,
    "cleanup_time": "2024-01-15T14:30:00Z"
  }
}
```

---

## 新增系统配置项

### 监控配置
- `health_monitor_enabled`: 是否启用系统健康监控
- `health_monitor_interval`: 健康监控间隔时间（秒）
- `max_health_logs`: 最大健康日志记录数
- `system_alert_email`: 系统告警邮件地址
- `performance_threshold_cpu`: CPU使用率告警阈值（%）
- `performance_threshold_memory`: 内存使用率告警阈值（%）
- `performance_threshold_disk`: 磁盘使用率告警阈值（%）
- `log_retention_days`: 日志保留天数

---

## 数据库视图

### 1. system_logs_summary（增强版）
```sql
SELECT 
    DATE(created_at) as log_date,
    log_type,
    operation,
    status,
    COUNT(*) as action_count,
    COUNT(DISTINCT user_id) as unique_users,
    COUNT(DISTINCT ip_address) as unique_ips,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_count,
    COUNT(CASE WHEN status = 'success' THEN 1 END) as success_count
FROM system_logs
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
GROUP BY DATE(created_at), log_type, operation, status
ORDER BY log_date DESC, action_count DESC;
```

### 2. system_health_summary（新增）
```sql
SELECT 
    DATE(record_time) as health_date,
    AVG(cpu_usage) as avg_cpu_usage,
    AVG(memory_usage) as avg_memory_usage,
    AVG(disk_usage) as avg_disk_usage,
    AVG(response_time_avg) as avg_response_time,
    MAX(active_users) as max_active_users,
    SUM(request_count) as total_requests,
    SUM(error_count) as total_errors,
    COUNT(CASE WHEN db_status = 'error' THEN 1 END) as db_error_count,
    COUNT(CASE WHEN db_status = 'warning' THEN 1 END) as db_warning_count
FROM system_health_logs
WHERE record_time >= DATE_SUB(NOW(), INTERVAL 7 DAY)
GROUP BY DATE(record_time)
ORDER BY health_date DESC;
```

---

## 自动清理机制

### 1. 日志自动过期
- 系统日志：默认保留90天，可通过配置调整
- 健康日志：默认保留30天
- 自动清理事件：每天凌晨3点执行

### 2. 备份记录清理
- 成功备份记录：保留90天
- 自动清理事件：每周日凌晨3点执行

---

## 使用示例

### 1. 监控系统健康状态
```bash
# 获取系统健康日志
curl -X GET "http://localhost:8080/api/admin/logs/health?page=1&size=20" \
  -H "Authorization: Bearer <token>"

# 记录系统健康状态
curl -X POST "http://localhost:8080/api/admin/logs/health" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "cpu_usage": 45.2,
    "memory_usage": 67.8,
    "disk_usage": 72.1,
    "db_status": "healthy",
    "db_connection_count": 25,
    "active_users": 45,
    "request_count": 1234,
    "error_count": 5,
    "response_time_avg": 125.5
  }'
```

### 2. 查看增强的系统日志
```bash
# 按日志类型筛选
curl -X GET "http://localhost:8080/api/admin/logs?log_type=error&status=failed" \
  -H "Authorization: Bearer <token>"

# 获取日志统计
curl -X GET "http://localhost:8080/api/admin/logs/summary?days=7" \
  -H "Authorization: Bearer <token>"
```

### 3. 管理系统设置
```bash
# 更新系统设置
curl -X PUT "http://localhost:8080/api/admin/settings/log_retention_days" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "setting_value": "60",
    "description": "日志保留60天"
  }'
```

---

## 注意事项

1. **权限要求**: 所有新增API都需要管理员权限
2. **数据量**: 健康日志会定期产生，建议定期清理
3. **性能影响**: 大量日志查询可能影响性能，建议使用分页
4. **配置管理**: 新增的监控配置项需要合理设置阈值
5. **自动清理**: 系统会自动清理过期数据，无需手动干预

---

## 迁移说明

1. 运行 `run_database_enhancement.bat` 执行数据库迁移
2. 重启后端服务以加载新的API功能
3. 检查数据库表结构是否正确更新
4. 验证新的API端点是否正常工作

这个增强版本提供了更完善的系统监控、日志管理和健康检查功能，有助于更好地维护和管理系统。 