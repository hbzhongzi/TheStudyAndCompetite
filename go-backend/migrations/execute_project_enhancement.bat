@echo off
chcp 65001 >nul
echo =============================================
echo 项目管理模块功能增强 - 数据库更新脚本
echo =============================================
echo.

REM 检查MySQL是否安装
mysql --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误: 未检测到MySQL客户端，请确保MySQL已安装并添加到PATH环境变量
    pause
    exit /b 1
)

echo ✅ MySQL客户端检测成功
echo.

REM 设置数据库连接参数
set /p DB_HOST=请输入数据库主机地址 (默认: localhost): 
if "%DB_HOST%"=="" set DB_HOST=localhost

set /p DB_PORT=请输入数据库端口 (默认: 3306): 
if "%DB_PORT%"=="" set DB_PORT=3306

set /p DB_NAME=请输入数据库名称: 
if "%DB_NAME%"=="" (
    echo ❌ 错误: 数据库名称不能为空
    pause
    exit /b 1
)

set /p DB_USER=请输入数据库用户名: 
if "%DB_USER%"=="" (
    echo ❌ 错误: 数据库用户名不能为空
    pause
    exit /b 1
)

set /p DB_PASS=请输入数据库密码: 
if "%DB_PASS%"=="" (
    echo ❌ 错误: 数据库密码不能为空
    pause
    exit /b 1
)

echo.
echo =============================================
echo 数据库连接信息确认
echo =============================================
echo 主机: %DB_HOST%
echo 端口: %DB_PORT%
echo 数据库: %DB_NAME%
echo 用户名: %DB_USER%
echo =============================================
echo.

set /p CONFIRM=确认以上信息正确吗？(y/N): 
if /i not "%CONFIRM%"=="y" (
    echo 操作已取消
    pause
    exit /b 0
)

echo.
echo =============================================
echo 开始执行数据库更新...
echo =============================================

REM 创建备份
echo 📋 正在创建数据库备份...
set BACKUP_FILE=project_enhancement_backup_%date:~0,4%%date:~5,2%%date:~8,2%_%time:~0,2%%time:~3,2%%time:~6,2%.sql
set BACKUP_FILE=%BACKUP_FILE: =0%

mysqldump -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASS% --single-transaction --routines --triggers %DB_NAME% > "%BACKUP_FILE%"

if %errorlevel% neq 0 (
    echo ❌ 错误: 数据库备份失败，请检查连接信息
    pause
    exit /b 1
)

echo ✅ 数据库备份完成: %BACKUP_FILE%
echo.

REM 执行更新脚本
echo 🔄 正在执行数据库更新脚本...
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASS% %DB_NAME% < project_management_enhancement.sql

if %errorlevel% neq 0 (
    echo ❌ 错误: 数据库更新失败
    echo.
    echo 💡 建议操作:
    echo 1. 检查SQL脚本语法是否正确
    echo 2. 检查数据库用户权限是否足够
    echo 3. 检查是否有数据冲突
    echo.
    echo 📋 已创建备份文件: %BACKUP_FILE%
    echo 如需恢复，请执行: mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASS% %DB_NAME% < "%BACKUP_FILE%"
    pause
    exit /b 1
)

echo ✅ 数据库更新成功完成！
echo.

REM 验证更新结果
echo 🔍 正在验证更新结果...
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASS% %DB_NAME% -e "SELECT 'projects表新字段检查' as check_item, COUNT(*) as count FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='%DB_NAME%' AND TABLE_NAME='projects' AND COLUMN_NAME IN ('start_date', 'expected_end_date', 'actual_end_date', 'progress');"

mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASS% %DB_NAME% -e "SELECT '新表检查' as check_item, COUNT(*) as count FROM information_schema.TABLES WHERE TABLE_SCHEMA='%DB_NAME%' AND TABLE_NAME IN ('project_milestones', 'project_extensions', 'file_type_configs', 'project_type_paths', 'project_type_stats', 'project_review_flows', 'review_delegations', 'review_reminders', 'project_notifications', 'notification_templates');"

echo.
echo =============================================
echo 🎉 项目管理模块功能增强完成！
echo =============================================
echo.
echo 📋 更新内容:
echo ✅ 项目状态管理增强 (新增4个状态)
echo ✅ 项目生命周期管理 (进度跟踪、里程碑、延期申请)
echo ✅ 成果文件管理增强 (文件类型、版本管理、审核状态)
echo ✅ 项目分类管理增强 (层级管理、路径查询、统计)
echo ✅ 审核流程增强 (多级审核、委托、提醒)
echo ✅ 通知系统基础结构 (模板、优先级、已读状态)
echo.
echo 📁 备份文件: %BACKUP_FILE%
echo 📊 新增表数量: 10个
echo 🔧 修改表数量: 3个
echo.
echo 💡 后续操作建议:
echo 1. 更新Go模型文件以匹配新的数据库结构
echo 2. 更新服务层代码以支持新功能
echo 3. 更新控制器和路由以提供新的API接口
echo 4. 测试新功能是否正常工作
echo.
pause 