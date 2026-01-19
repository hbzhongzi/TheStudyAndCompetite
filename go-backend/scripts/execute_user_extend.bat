@echo off
echo 正在执行用户表扩展迁移...
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
mysql -u root -p < migrations/extend_user_table.sql

if %errorlevel% equ 0 (
    echo.
    echo 用户表扩展迁移完成！
    echo 新添加的字段：
    echo - department: 部门/院系
    echo - title: 职称/职位  
    echo - grade: 年级
    echo - major: 专业
    echo - created_at: 创建时间
    echo - updated_at: 更新时间
) else (
    echo.
    echo 迁移失败，请检查错误信息
)

echo.
pause 