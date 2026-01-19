@echo off
chcp 65001 >nul 2>&1
setlocal enabledelayedexpansion

echo ========================================
echo    云梦高校科研竞赛管理系统
echo        数据库初始化脚本 v2.0
echo ========================================
echo.

:: 设置颜色代码
set "RED=[91m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "BLUE=[94m"
set "CYAN=[96m"
set "RESET=[0m"

:: 设置项目路径
set PROJECT_ROOT=%~dp0..
cd /d "!PROJECT_ROOT!" 2>nul
if !ERRORLEVEL! NEQ 0 (
    echo %RED%[错误] 无法切换到项目目录%RESET%
    pause
    exit /b 1
)

:: 检查MySQL环境
echo %BLUE%检查MySQL环境...%RESET%
mysql --version >nul 2>&1
if errorlevel 1 (
    echo %RED%[错误] MySQL未安装或未添加到PATH环境变量%RESET%
    echo %YELLOW%请先安装MySQL: https://dev.mysql.com/downloads/mysql/%RESET%
    pause
    exit /b 1
)

echo %GREEN%[成功] MySQL环境正常%RESET%
echo.

:: 获取MySQL连接参数
echo %CYAN%请输入MySQL连接信息:%RESET%
set /p DB_USER=用户名 (默认: root): 
if "!DB_USER!"=="" set DB_USER=root

set /p DB_PASSWORD=密码: 
if "!DB_PASSWORD!"=="" (
    echo %RED%[错误] 密码不能为空%RESET%
    pause
    exit /b 1
)

set DB_HOST=localhost
set DB_PORT=3306
set DB_NAME=cloud_dream_system

echo.
echo %BLUE%数据库连接信息:%RESET%
echo   主机: !DB_HOST!
echo   端口: !DB_PORT!
echo   用户: !DB_USER!
echo   数据库: !DB_NAME!
echo.

:: 测试MySQL连接
echo %BLUE%测试MySQL连接...%RESET%
mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! -e "SELECT VERSION() as mysql_version;" 2>nul >nul
if !ERRORLEVEL! NEQ 0 (
    echo %RED%[错误] 无法连接到MySQL数据库%RESET%
    echo %YELLOW%可能的原因:%RESET%
    echo %YELLOW%1. MySQL服务未启动%RESET%
    echo %YELLOW%2. 连接参数错误（主机、端口、用户名、密码）%RESET%
    echo %YELLOW%3. 防火墙阻止连接%RESET%
    echo %YELLOW%4. MySQL配置不允许远程连接%RESET%
    pause
    exit /b 1
)

echo %GREEN%[成功] MySQL连接正常%RESET%
echo.

:: 检查数据库是否存在
echo %BLUE%检查数据库是否存在...%RESET%
mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! -e "USE !DB_NAME!;" 2>nul >nul
if !ERRORLEVEL! NEQ 0 (
    echo %YELLOW%[信息] 数据库不存在，将创建新数据库%RESET%
    mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! -e "CREATE DATABASE !DB_NAME! CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>nul >nul
    if !ERRORLEVEL! NEQ 0 (
        echo %RED%[错误] 创建数据库失败%RESET%
        echo %YELLOW%可能的原因:%RESET%
        echo %YELLOW%1. 用户权限不足%RESET%
        echo %YELLOW%2. 数据库名称已存在但无法访问%RESET%
        echo %YELLOW%3. MySQL配置问题%RESET%
        pause
        exit /b 1
    )
    echo %GREEN%[成功] 数据库创建完成%RESET%
) else (
    echo %GREEN%[信息] 数据库已存在%RESET%
    
    :: 检查是否要清空数据库
    echo %YELLOW%[警告] 数据库已存在，是否要清空并重新初始化? (y/n)%RESET%
    set /p CLEAR_DB_CHOICE=
    
    if /i "!CLEAR_DB_CHOICE!"=="y" (
        echo %BLUE%清空数据库...%RESET%
        mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! -e "DROP DATABASE !DB_NAME!;" 2>nul >nul
        mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! -e "CREATE DATABASE !DB_NAME! CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>nul >nul
        if !ERRORLEVEL! NEQ 0 (
            echo %RED%[错误] 清空数据库失败%RESET%
            pause
            exit /b 1
        )
        echo %GREEN%[成功] 数据库已清空%RESET%
    ) else (
        echo %YELLOW%[跳过] 数据库清空%RESET%
    )
)

:: 选择SQL文件
echo.
echo %BLUE%选择SQL初始化文件:%RESET%
echo 1) 使用兼容版SQL文件 (推荐MySQL 8.0以下版本)
echo 2) 使用标准版SQL文件 (推荐MySQL 8.0及以上版本)
echo 3) 自动选择 (根据MySQL版本)
echo.
set /p SQL_CHOICE=请选择 (1-3，默认3): 

if "!SQL_CHOICE!"=="" set SQL_CHOICE=3

if "!SQL_CHOICE!"=="1" (
    set SQL_FILE=sql\database_setup_compatible.sql
    echo %GREEN%[选择] 兼容版SQL文件%RESET%
) else if "!SQL_CHOICE!"=="2" (
    set SQL_FILE=sql\database_setup.sql
    echo %GREEN%[选择] 标准版SQL文件%RESET%
) else (
    :: 自动选择
    for /f "tokens=6" %%i in ('mysql --version 2^>nul') do set MYSQL_VERSION_NUM=%%i
    if "!MYSQL_VERSION_NUM!"=="" (
        echo %YELLOW%[警告] 无法检测MySQL版本，使用兼容版%RESET%
        set SQL_FILE=sql\database_setup_compatible.sql
    ) else (
        echo %BLUE%检测到MySQL版本: !MYSQL_VERSION_NUM!%RESET%
        if "!MYSQL_VERSION_NUM:~0,1!"=="8" (
            set SQL_FILE=sql\database_setup.sql
            echo %GREEN%[自动选择] 标准版SQL文件 (MySQL 8.x)%RESET%
        ) else (
            set SQL_FILE=sql\database_setup_compatible.sql
            echo %GREEN%[自动选择] 兼容版SQL文件 (MySQL 5.x/6.x/7.x)%RESET%
        )
    )
)

:: 检查SQL文件是否存在
if not exist "!SQL_FILE!" (
    echo %RED%[错误] SQL文件不存在: !SQL_FILE!%RESET%
    echo %YELLOW%请检查以下文件是否存在:%RESET%
    echo %YELLOW%1. sql\database_setup_compatible.sql%RESET%
    echo %YELLOW%2. sql\database_setup.sql%RESET%
    pause
    exit /b 1
)

echo %GREEN%[信息] 使用SQL文件: !SQL_FILE!%RESET%
echo.

:: 执行SQL脚本
echo %BLUE%执行数据库初始化...%RESET%
echo %YELLOW%这可能需要几分钟时间，请耐心等待...%RESET%
echo.

mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! !DB_NAME! < "!SQL_FILE!" 2>nul
if !ERRORLEVEL! EQU 0 (
    echo %GREEN%[成功] 数据库初始化完成%RESET%
) else (
    echo %RED%[错误] 数据库初始化失败%RESET%
    echo %YELLOW%可能的原因:%RESET%
    echo %YELLOW%1. SQL文件语法错误%RESET%
    echo %YELLOW%2. 用户权限不足%RESET%
    echo %YELLOW%3. 数据库连接问题%RESET%
    echo %YELLOW%4. MySQL版本兼容性问题%RESET%
    echo.
    echo %CYAN%请检查SQL文件并重试%RESET%
    pause
    exit /b 1
)

:: 验证初始化结果
echo.
echo %BLUE%验证初始化结果...%RESET%

:: 检查表数量
set TABLE_COUNT=0
for /f "tokens=1" %%a in ('mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! !DB_NAME! -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = '!DB_NAME!';" 2^>nul ^| findstr /r "^[0-9]"') do (
    set TABLE_COUNT=%%a
)

if "!TABLE_COUNT!"=="" (
    echo %YELLOW%[警告] 无法获取表数量%RESET%
) else (
    echo %GREEN%[信息] 成功创建 !TABLE_COUNT! 个表%RESET%
)

:: 检查关键表是否存在
echo %BLUE%检查关键表...%RESET%
mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! !DB_NAME! -e "SHOW TABLES LIKE 'users';" 2>nul | findstr "users" >nul && echo %GREEN%   ✓ users表%RESET% || echo %RED%   ✗ users表%RESET%
mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! !DB_NAME! -e "SHOW TABLES LIKE 'projects';" 2>nul | findstr "projects" >nul && echo %GREEN%   ✓ projects表%RESET% || echo %RED%   ✗ projects表%RESET%
mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! !DB_NAME! -e "SHOW TABLES LIKE 'competitions';" 2>nul | findstr "competitions" >nul && echo %GREEN%   ✓ competitions表%RESET% || echo %RED%   ✗ competitions表%RESET%

:: 检查默认数据
echo.
echo %BLUE%检查默认数据...%RESET%
set USER_COUNT=0
for /f "tokens=1" %%a in ('mysql -h!DB_HOST! -P!DB_PORT! -u!DB_USER! -p!DB_PASSWORD! !DB_NAME! -e "SELECT COUNT(*) FROM users;" 2^>nul ^| findstr /r "^[0-9]"') do (
    set USER_COUNT=%%a
)

if "!USER_COUNT!"=="" (
    echo %YELLOW%[警告] 无法获取用户数量%RESET%
) else (
    echo %GREEN%[信息] 默认用户数量: !USER_COUNT!%RESET%
)

echo.
echo %GREEN%========================================%RESET%
echo %GREEN%数据库初始化完成%RESET%
echo %GREEN%========================================%RESET%
echo.
echo %CYAN%数据库信息:%RESET%
echo   数据库名: !DB_NAME!
echo   字符集: utf8mb4
echo   排序规则: utf8mb4_unicode_ci
echo.
echo %CYAN%默认用户:%RESET%
echo   管理员: admin / 123456
echo   教师: teacher001 / 123456
echo   学生: student001 / 123456
echo.
echo %YELLOW%提示: 现在可以运行launch.bat启动后端服务%RESET%
echo.
pause 