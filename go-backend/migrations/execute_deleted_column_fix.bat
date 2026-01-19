@echo off
echo ========================================
echo 执行数据库迁移：添加deleted列
echo ========================================
echo.

REM 设置数据库连接参数
set DB_HOST=localhost
set DB_PORT=3306
set DB_USER=root
set DB_PASSWORD=123456
set DB_NAME=cloud_dream_system

echo 数据库连接信息：
echo 主机: %DB_HOST%
echo 端口: %DB_PORT%
echo 用户: %DB_USER%
echo 数据库: %DB_NAME%
echo.

echo 正在执行迁移脚本...
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASSWORD% %DB_NAME% < add_deleted_column.sql

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo 迁移执行成功！
    echo ========================================
    echo.
    echo 已添加以下列：
    echo - projects.deleted (软删除标记)
    echo - users.deleted (软删除标记)
    echo.
    echo 已添加以下索引：
    echo - idx_projects_deleted
    echo - idx_users_deleted
    echo.
) else (
    echo.
    echo ========================================
    echo 迁移执行失败！
    echo ========================================
    echo.
    echo 请检查：
    echo 1. 数据库连接参数是否正确
    echo 2. 数据库服务是否正在运行
    echo 3. 用户权限是否足够
    echo.
)

pause 