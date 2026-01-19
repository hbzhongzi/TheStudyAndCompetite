@echo off
echo ========================================
echo Cloud Dream System - Quick Start All
echo ========================================
echo.

echo This script will:
echo 1. Initialize database (if needed)
echo 2. Create default users (if needed)
echo 3. Start backend service
echo 4. Start frontend service
echo.

set /p choice="Do you want to continue? (y/n): "
if /i "%choice%" neq "y" (
    echo Setup cancelled.
    pause
    exit /b 0
)

echo.
echo ========================================
echo Step 1: Checking Database
echo ========================================

REM Check if MySQL is available
mysql --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: MySQL not found. Please install MySQL first.
    pause
    exit /b 1
)

REM Check if database exists
echo Checking if database exists...
mysql -u root -p -e "USE cloud_dream_system;" >nul 2>&1
if %errorlevel% neq 0 (
    echo Database not found. Initializing database...
    call run_init_database.bat
    if %errorlevel% neq 0 (
        echo Database initialization failed!
        pause
        exit /b 1
    )
    
    echo Creating default users...
    call run_create_users.bat
    if %errorlevel% neq 0 (
        echo User creation failed!
        pause
        exit /b 1
    )
) else (
    echo Database exists. Checking users...
    mysql -u root -p -e "USE cloud_dream_system; SELECT COUNT(*) FROM users;" >nul 2>&1
    if %errorlevel% neq 0 (
        echo Users not found. Creating default users...
        call run_create_users.bat
    ) else (
        echo Users exist. Skipping user creation.
    )
)

echo.
echo ========================================
echo Step 2: Starting Backend Service
echo ========================================

echo Starting backend service...
echo Backend will be available at: http://localhost:8080
echo.

REM Check if Go is available
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Go not found. Please install Go first.
    pause
    exit /b 1
)

REM Start backend in background
start "Cloud Dream Backend" cmd /k "cd /d %CD% && go run main.go"

echo Backend service started in background.
echo.

echo ========================================
echo Step 3: Starting Frontend Service
echo ========================================

echo Starting frontend service...
echo Frontend will be available at: http://localhost:5173
echo.

REM Check if Node.js is available
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Node.js not found. Please install Node.js first.
    echo Backend is still running. You can start frontend manually.
    pause
    exit /b 1
)

REM Check if frontend directory exists
if not exist "..\yunmeng-frontend" (
    echo ERROR: Frontend directory not found.
    echo Backend is still running. You can start frontend manually.
    pause
    exit /b 1
)

REM Start frontend in background
start "Cloud Dream Frontend" cmd /k "cd /d %CD%\..\yunmeng-frontend && npm run dev"

echo Frontend service started in background.
echo.

echo ========================================
echo Setup Complete!
echo ========================================
echo.
echo Services started:
echo - Backend: http://localhost:8080
echo - Frontend: http://localhost:5173
echo.
echo Default login credentials:
echo - Admin: admin / admin123
echo - Teacher: teacher / teacher123
echo - Student: student / student123
echo.
echo API Documentation:
echo - Backend API: http://localhost:8080/api/docs (if available)
echo - System Health: http://localhost:8080/api/admin/logs/health
echo.
echo To stop services:
echo - Close the command windows that opened
echo - Or use Ctrl+C in each window
echo.
echo Press any key to open the frontend in your browser...
pause >nul

REM Open frontend in default browser
start http://localhost:5173

echo.
echo Enjoy using Cloud Dream System!
echo.
pause 