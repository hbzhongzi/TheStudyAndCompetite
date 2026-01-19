@echo off
chcp 65001 >nul
echo ========================================
echo 竞赛数据库同步更新脚本
echo ========================================
echo.

:: 设置颜色
set "GREEN=[92m"
set "YELLOW=[93m"
set "RED=[91m"
set "BLUE=[94m"
set "CYAN=[96m"
set "RESET=[0m"

:: 检查MySQL客户端是否可用
echo %CYAN%检查MySQL客户端...%RESET%
mysql --version >nul 2>&1
if %errorlevel% neq 0 (
    echo %RED%错误：未找到MySQL客户端，请确保MySQL已安装并添加到PATH环境变量%RESET%
    pause
    exit /b 1
)
echo %GREEN%MySQL客户端检查通过%RESET%
echo.

:: 获取数据库连接信息
echo %YELLOW%请输入数据库连接信息：%RESET%
set /p DB_HOST=数据库主机地址 (默认: localhost): 
if "%DB_HOST%"=="" set DB_HOST=localhost

set /p DB_PORT=数据库端口 (默认: 3306): 
if "%DB_PORT%"=="" set DB_PORT=3306

set /p DB_NAME=数据库名称: 
if "%DB_NAME%"=="" (
    echo %RED%错误：数据库名称不能为空%RESET%
    pause
    exit /b 1
)

set /p DB_USER=数据库用户名: 
if "%DB_USER%"=="" (
    echo %RED%错误：数据库用户名不能为空%RESET%
    pause
    exit /b 1
)

set /p DB_PASS=数据库密码: 
if "%DB_PASS%"=="" (
    echo %RED%错误：数据库密码不能为空%RESET%
    pause
    exit /b 1
)

echo.
echo %CYAN%数据库连接信息：%RESET%
echo 主机: %DB_HOST%
echo 端口: %DB_PORT%
echo 数据库: %DB_NAME%
echo 用户: %DB_USER%
echo.

:: 测试数据库连接
echo %CYAN%测试数据库连接...%RESET%
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASS% -e "USE %DB_NAME%; SELECT '连接成功' as status;" 2>nul
if %errorlevel% neq 0 (
    echo %RED%错误：数据库连接失败，请检查连接信息%RESET%
    pause
    exit /b 1
)
echo %GREEN%数据库连接成功%RESET%
echo.

:: 备份数据库
echo %YELLOW%是否要备份数据库？(y/n，默认: y)%RESET%
set /p BACKUP_CHOICE=
if /i "%BACKUP_CHOICE%"=="" set BACKUP_CHOICE=y
if /i "%BACKUP_CHOICE%"=="y" (
    echo %CYAN%开始备份数据库...%RESET%
    set BACKUP_FILE=backup_%DB_NAME%_%date:~0,4%%date:~5,2%%date:~8,2%_%time:~0,2%%time:~3,2%%time:~6,2%.sql
    set BACKUP_FILE=%BACKUP_FILE: =0%
    
    mysqldump -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASS% --single-transaction --routines --triggers %DB_NAME% > "%BACKUP_FILE%"
    if %errorlevel% equ 0 (
        echo %GREEN%数据库备份成功：%BACKUP_FILE%%RESET%
    ) else (
        echo %RED%数据库备份失败%RESET%
        pause
        exit /b 1
    )
    echo.
)

:: 确认执行更新
echo %YELLOW%确认执行数据库更新？这将修改数据库结构 (y/n，默认: y)%RESET%
set /p CONFIRM_UPDATE=
if /i "%CONFIRM_UPDATE%"=="" set CONFIRM_UPDATE=y
if /i "%CONFIRM_UPDATE%"=="y" (
    echo %CYAN%开始执行数据库更新...%RESET%
    echo.
    
    :: 执行完整的数据库更新脚本
    echo %CYAN%执行竞赛表结构更新...%RESET%
    mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASS% %DB_NAME% < "update_competition_database_complete.sql"
    if %errorlevel% equ 0 (
        echo %GREEN%数据库更新执行成功%RESET%
    ) else (
        echo %RED%数据库更新执行失败，请检查错误信息%RESET%
        echo %YELLOW%如果更新失败，可以使用备份文件恢复数据库%RESET%
        pause
        exit /b 1
    )
    echo.
    
    :: 验证更新结果
    echo %CYAN%验证更新结果...%RESET%
    mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASS% %DB_NAME% -e "
    SELECT 'competitions' AS table_name, COUNT(*) AS total_records FROM competitions
    UNION ALL
    SELECT 'competition_registrations' AS table_name, COUNT(*) AS total_records FROM competition_registrations
    UNION ALL
    SELECT 'competition_submissions' AS table_name, COUNT(*) AS total_records FROM competition_submissions;
    "
    
    echo.
    echo %CYAN%检查新字段...%RESET%
    mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASS% %DB_NAME% -e "
    SELECT 
        'competitions' AS table_name,
        COUNT(registration_start) AS has_registration_start,
        COUNT(registration_end) AS has_registration_end,
        COUNT(location) AS has_location,
        COUNT(contact) AS has_contact
    FROM competitions;
    "
    
    echo.
    echo %GREEN%数据库更新完成！%RESET%
    echo.
    echo %CYAN%更新内容包括：%RESET%
    echo - 添加了报名时间相关字段
    echo - 添加了竞赛详细信息字段
    echo - 创建了竞赛关联表
    echo - 添加了性能优化索引
    echo - 设置了外键约束
    echo - 更新了现有数据
    echo.
    
) else (
    echo %YELLOW%取消执行数据库更新%RESET%
)

echo.
echo %CYAN%按任意键退出...%RESET%
pause >nul 