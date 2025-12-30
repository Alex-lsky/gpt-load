@echo off
REM 前端开发服务器启动脚本
echo ========================================
echo   GPT-Load 前端开发服务器
echo ========================================

REM 检查Node.js和npm
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

REM 进入前端目录
cd web

REM 安装依赖
echo.
echo 📦 安装前端依赖...
call npm install
if %errorlevel% neq 0 (
    echo ❌ 前端依赖安装失败
    pause
    exit /b 1
)

REM 启动开发服务器
echo.
echo 🎨 启动前端开发服务器...
echo 开发服务器将在 http://localhost:5173 启动
echo 按 Ctrl+C 停止服务器
echo.
call npm run dev

pause
