@echo off
REM 开发环境检查脚本
echo ========================================
echo   GPT-Load 开发环境检查
echo ========================================

set "error_count=0"

REM 检查 Go
echo 🔍 检查 Go...
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ Go 未安装或不在 PATH 中
    set /a error_count+=1
) else (
    for /f "tokens=3" %%i in ('go version') do set go_version=%%i
    echo ✅ Go 已安装: %go_version%
)

REM 检查 Node.js
echo.
echo 🔍 检查 Node.js...
where node >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ Node.js 未安装或不在 PATH 中
    set /a error_count+=1
) else (
    for /f "tokens=*" %%i in ('node --version') do set node_version=%%i
    echo ✅ Node.js 已安装: %node_version%
)

REM 检查 npm
echo.
echo 🔍 检查 npm...
where npm >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ npm 未安装或不在 PATH 中
    set /a error_count+=1
) else (
    for /f "tokens=*" %%i in ('npm --version') do set npm_version=%%i
    echo ✅ npm 已安装: %npm_version%
)

REM 检查 Git
echo.
echo 🔍 检查 Git...
where git >nul 2>nul
if %errorlevel% neq 0 (
    echo ⚠️  Git 未安装或不在 PATH 中（可选）
) else (
    for /f "tokens=3" %%i in ('git --version') do set git_version=%%i
    echo ✅ Git 已安装: %git_version%
)

REM 检查项目文件
echo.
echo 🔍 检查项目文件...
if not exist "go.mod" (
    echo ❌ go.mod 文件不存在
    set /a error_count+=1
) else (
    echo ✅ go.mod 文件存在
)

if not exist "web\package.json" (
    echo ❌ web/package.json 文件不存在
    set /a error_count+=1
) else (
    echo ✅ web/package.json 文件存在
)

if not exist "main.go" (
    echo ❌ main.go 文件不存在
    set /a error_count+=1
) else (
    echo ✅ main.go 文件存在
)

REM 检查 Go 模块
echo.
echo 🔍 检查 Go 模块...
go mod verify >nul 2>nul
if %errorlevel% neq 0 (
    echo ⚠️  Go 模块验证失败，建议运行 'go mod tidy'
) else (
    echo ✅ Go 模块验证通过
)

REM 总结
echo.
echo ========================================
if %error_count% equ 0 (
    echo ✅ 环境检查通过！可以开始开发了。
    echo.
    echo 💡 快速启动命令:
    echo    scripts\dev.bat          - 完整开发环境
    echo    scripts\dev-frontend.bat - 仅前端开发服务器
    echo    scripts\dev-backend.bat  - 仅后端开发服务器
) else (
    echo ❌ 发现 %error_count% 个问题，请先解决这些问题。
    echo.
    echo 📋 安装指南:
    echo    Go: https://golang.org/dl/
    echo    Node.js: https://nodejs.org/
    echo    Git: https://git-scm.com/
)
echo ========================================

pause
