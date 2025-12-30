# GPT-Load 开发环境启动脚本 - PowerShell 版本
param(
    [switch]$Frontend,
    [switch]$Backend,
    [switch]$Check,
    [switch]$Help
)

function Write-Header {
    param([string]$Title)
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "  $Title" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host "✅ $Message" -ForegroundColor Green
}

function Write-Error {
    param([string]$Message)
    Write-Host "❌ $Message" -ForegroundColor Red
}

function Write-Warning {
    param([string]$Message)
    Write-Host "⚠️  $Message" -ForegroundColor Yellow
}

function Write-Info {
    param([string]$Message)
    Write-Host "📋 $Message" -ForegroundColor Blue
}

function Test-Command {
    param([string]$Command)
    try {
        Get-Command $Command -ErrorAction Stop | Out-Null
        return $true
    }
    catch {
        return $false
    }
}

function Show-Help {
    Write-Header "GPT-Load 开发脚本帮助"
    Write-Host ""
    Write-Host "用法: .\scripts\dev.ps1 [选项]" -ForegroundColor White
    Write-Host ""
    Write-Host "选项:" -ForegroundColor White
    Write-Host "  -Frontend    仅启动前端开发服务器" -ForegroundColor Gray
    Write-Host "  -Backend     仅启动后端开发服务器" -ForegroundColor Gray
    Write-Host "  -Check       检查开发环境" -ForegroundColor Gray
    Write-Host "  -Help        显示此帮助信息" -ForegroundColor Gray
    Write-Host ""
    Write-Host "示例:" -ForegroundColor White
    Write-Host "  .\scripts\dev.ps1           # 完整开发环境" -ForegroundColor Gray
    Write-Host "  .\scripts\dev.ps1 -Frontend # 仅前端" -ForegroundColor Gray
    Write-Host "  .\scripts\dev.ps1 -Backend  # 仅后端" -ForegroundColor Gray
    Write-Host "  .\scripts\dev.ps1 -Check    # 环境检查" -ForegroundColor Gray
}

function Test-Environment {
    Write-Header "GPT-Load 开发环境检查"
    $errorCount = 0

    # 检查 Go
    Write-Host "🔍 检查 Go..." -ForegroundColor Blue
    if (Test-Command "go") {
        $goVersion = (go version).Split()[2]
        Write-Success "Go 已安装: $goVersion"
    } else {
        Write-Error "Go 未安装或不在 PATH 中"
        $errorCount++
    }

    # 检查 Node.js
    Write-Host ""
    Write-Host "🔍 检查 Node.js..." -ForegroundColor Blue
    if (Test-Command "node") {
        $nodeVersion = node --version
        Write-Success "Node.js 已安装: $nodeVersion"
    } else {
        Write-Error "Node.js 未安装或不在 PATH 中"
        $errorCount++
    }

    # 检查 npm
    Write-Host ""
    Write-Host "🔍 检查 npm..." -ForegroundColor Blue
    if (Test-Command "npm") {
        $npmVersion = npm --version
        Write-Success "npm 已安装: $npmVersion"
    } else {
        Write-Error "npm 未安装或不在 PATH 中"
        $errorCount++
    }

    # 检查 Git
    Write-Host ""
    Write-Host "🔍 检查 Git..." -ForegroundColor Blue
    if (Test-Command "git") {
        $gitVersion = (git --version).Split()[2]
        Write-Success "Git 已安装: $gitVersion"
    } else {
        Write-Warning "Git 未安装或不在 PATH 中（可选）"
    }

    # 检查项目文件
    Write-Host ""
    Write-Host "🔍 检查项目文件..." -ForegroundColor Blue

    if (Test-Path "go.mod") {
        Write-Success "go.mod 文件存在"
    } else {
        Write-Error "go.mod 文件不存在"
        $errorCount++
    }

    if (Test-Path "web\package.json") {
        Write-Success "web/package.json 文件存在"
    } else {
        Write-Error "web/package.json 文件不存在"
        $errorCount++
    }

    if (Test-Path "main.go") {
        Write-Success "main.go 文件存在"
    } else {
        Write-Error "main.go 文件不存在"
        $errorCount++
    }

    # 检查 Go 模块
    Write-Host ""
    Write-Host "🔍 检查 Go 模块..." -ForegroundColor Blue
    try {
        go mod verify | Out-Null
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Go 模块验证通过"
        } else {
            Write-Warning "Go 模块验证失败，建议运行 'go mod tidy'"
        }
    } catch {
        Write-Warning "无法验证 Go 模块"
    }

    # 总结
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    if ($errorCount -eq 0) {
        Write-Success "环境检查通过！可以开始开发了。"
        Write-Host ""
        Write-Info "快速启动命令:"
        Write-Host "   .\scripts\dev.ps1           - 完整开发环境" -ForegroundColor Gray
        Write-Host "   .\scripts\dev.ps1 -Frontend - 仅前端开发服务器" -ForegroundColor Gray
        Write-Host "   .\scripts\dev.ps1 -Backend  - 仅后端开发服务器" -ForegroundColor Gray
    } else {
        Write-Error "发现 $errorCount 个问题，请先解决这些问题。"
        Write-Host ""
        Write-Info "安装指南:"
        Write-Host "   Go: https://golang.org/dl/" -ForegroundColor Gray
        Write-Host "   Node.js: https://nodejs.org/" -ForegroundColor Gray
        Write-Host "   Git: https://git-scm.com/" -ForegroundColor Gray
    }
    Write-Host "========================================" -ForegroundColor Cyan

    return $errorCount -eq 0
}

function Start-Frontend {
    Write-Header "GPT-Load 前端开发服务器"

    if (-not (Test-Command "node") -or -not (Test-Command "npm")) {
        Write-Error "Node.js 或 npm 未安装"
        return $false
    }

    Write-Success "环境检查通过"

    # 进入前端目录
    Push-Location "web"

    try {
        # 安装依赖
        Write-Host ""
        Write-Info "安装前端依赖..."
        npm install
        if ($LASTEXITCODE -ne 0) {
            Write-Error "前端依赖安装失败"
            return $false
        }

        # 启动开发服务器
        Write-Host ""
        Write-Info "启动前端开发服务器..."
        Write-Host "开发服务器将在 http://localhost:5173 启动" -ForegroundColor Green
        Write-Host "按 Ctrl+C 停止服务器" -ForegroundColor Yellow
        Write-Host ""
        npm run dev
    }
    finally {
        Pop-Location
    }

    return $true
}

function Start-Backend {
    Write-Header "GPT-Load 后端开发服务器"

    if (-not (Test-Command "go")) {
        Write-Error "Go 未安装"
        return $false
    }

    Write-Success "环境检查通过"

    # 整理 Go 模块
    Write-Host ""
    Write-Info "整理 Go 模块..."
    go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Go 模块整理失败"
        return $false
    }

    # 设置环境变量
    $env:GIN_MODE = "debug"
    $env:GO_ENV = "development"

    # 启动后端服务器
    Write-Host ""
    Write-Info "启动后端开发服务器..."
    Write-Host "服务器将在 http://localhost:8080 启动" -ForegroundColor Green
    Write-Host "开发模式已启用（包含竞态检测）" -ForegroundColor Yellow
    Write-Host "按 Ctrl+C 停止服务器" -ForegroundColor Yellow
    Write-Host ""
    go run -race ./main.go

    return $true
}

function Start-FullStack {
    Write-Header "GPT-Load 完整开发环境"

    # 环境检查
    if (-not (Test-Environment)) {
        return $false
    }

    # 安装前端依赖
    Write-Host ""
    Write-Info "安装前端依赖..."
    Push-Location "web"
    npm install
    if ($LASTEXITCODE -ne 0) {
        Write-Error "前端依赖安装失败"
        Pop-Location
        return $false
    }

    # 构建前端
    Write-Host ""
    Write-Info "构建前端..."
    npm run build
    if ($LASTEXITCODE -ne 0) {
        Write-Error "前端构建失败"
        Pop-Location
        return $false
    }
    Pop-Location

    # 整理 Go 模块
    Write-Host ""
    Write-Info "整理 Go 模块..."
    go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Go 模块整理失败"
        return $false
    }

    # 启动后端服务器
    Write-Host ""
    Write-Info "启动后端服务器..."
    Write-Host "服务器将在 http://localhost:8080 启动" -ForegroundColor Green
    Write-Host "按 Ctrl+C 停止服务器" -ForegroundColor Yellow
    Write-Host ""
    go run -race ./main.go

    return $true
}

# 主逻辑
if ($Help) {
    Show-Help
    exit 0
}

if ($Check) {
    Test-Environment
    exit 0
}

if ($Frontend) {
    Start-Frontend
    exit 0
}

if ($Backend) {
    Start-Backend
    exit 0
}

# 默认：完整开发环境
Start-FullStack
