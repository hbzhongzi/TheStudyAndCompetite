@echo off
echo 正在执行项目表字段修复迁移...
echo.

REM 检查MySQL是否运行
echo 检查MySQL服务状态...
net start | findstr "MySQL" >nul
if %errorlevel% neq 0 (
    echo MySQL服务未运行，请先启动MySQL服务
    pause
    exit /b 1
)

REM 执行迁移脚本
echo 执行数据库迁移脚本...
mysql -u root -p < migrations/add_project_missing_fields.sql

if %errorlevel% equ 0 (
    echo.
    echo 项目表字段修复迁移完成！
    echo 新添加的字段：
    echo - category_id: 项目分类ID
    echo - expected_end_date: 预计完成时间
    echo - actual_end_date: 实际完成时间
    echo - progress: 项目进度
    echo - is_extended: 是否延期
    echo - extension_count: 延期次数
    echo - force_status_reason: 强制状态变更原因
    echo - level: 项目级别
    echo - created_by: 创建者ID
    echo - created_at: 创建时间
    echo - updated_at: 更新时间
) else (
    echo.
    echo 迁移失败，请检查错误信息
)

echo.
pause 