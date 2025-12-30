@echo off
REM 开发环境启动脚本 - Windows批处理版本
echo ========================================
echo   GPT-Load 开发环境启动脚本
echo ========================================

REM 检查是否安装了必要的工具
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ 错误: 未找到 Go，请先安装 Go
    pause
    exit /b 1
)

where node >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ 错误: 未找到 Node.js，请先安装 Node.js
    pause
    exit /b 1
)

where npm >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ 错误: 未找到 npm，请先安装 npm
    pause
    exit /b 1
)

echo ✅ 环境检查通过

REM 安装前端依赖
echo.
echo 📦 安装前端依赖...
cd web
call npm install
if %errorlevel% neq 0 (
    echo ❌ 前端依赖安装失败
    pause
    exit /b 1
)

REM 构建前端
echo.
echo 🔨 构建前端...
call npm run build
if %errorlevel% neq 0 (
    echo ❌ 前端构建失败
    pause
    exit /b 1
)

cd ..

REM 整理Go模块
echo.
echo 📋 整理Go模块...
go mod tidy
if %errorlevel% neq 0 (
    echo ❌ Go模块整理失败
    pause
    exit /b 1
)

REM 启动后端服务器
echo.
echo 🚀 启动后端服务器...
echo 服务器将在 http://localhost:8080 启动
echo 按 Ctrl+C 停止服务器
echo.
go run -race ./main.go

pause
