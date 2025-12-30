@echo off
REM 后端开发服务器启动脚本
echo ========================================
echo   GPT-Load 后端开发服务器
echo ========================================

REM 检查Go
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ 错误: 未找到 Go，请先安装 Go
    pause
    exit /b 1
)

echo ✅ 环境检查通过

REM 整理Go模块
echo.
echo 📋 整理Go模块...
go mod tidy
if %errorlevel% neq 0 (
    echo ❌ Go模块整理失败
    pause
    exit /b 1
)

REM 设置开发环境变量
set GIN_MODE=debug
set GO_ENV=development

REM 启动后端服务器（开发模式）
echo.
echo 🚀 启动后端开发服务器...
echo 服务器将在 http://localhost:8080 启动
echo 开发模式已启用（包含竞态检测）
echo 按 Ctrl+C 停止服务器
echo.
go run -race ./main.go

pause
