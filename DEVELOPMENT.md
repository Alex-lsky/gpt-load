# GPT-Load 开发指南

本文档介绍如何在开发环境中快速启动和调试 GPT-Load 项目。

## 🚀 快速开始

### 方式一：使用 VSCode 工作区（推荐）

1. 打开 `gpt-load.code-workspace` 文件
2. VSCode 会提示安装推荐的扩展，点击安装
3. 使用以下调试配置：
   - **启动后端服务器**: 构建前端并启动后端
   - **启动后端服务器 (开发模式)**: 启用竞态检测的开发模式
   - **全栈调试**: 同时启动前端开发服务器和后端

### 方式二：使用批处理脚本

#### Windows 用户

```bash
# 完整开发环境（构建前端 + 启动后端）
scripts\dev.bat

# 仅启动前端开发服务器
scripts\dev-frontend.bat

# 仅启动后端开发服务器
scripts\dev-backend.bat
```

### 方式三：手动启动

#### 前端开发

```bash
cd web
npm install
npm run dev  # 开发服务器: http://localhost:5173
```

#### 后端开发

```bash
# 开发模式（推荐）
go run -race ./main.go

# 或者使用 Makefile
make dev
```

## 🛠️ VSCode 配置说明

### 调试配置 (launch.json)

- **启动后端服务器**: 自动构建前端并启动后端，适合生产环境测试
- **启动后端服务器 (开发模式)**: 启用竞态检测，适合开发调试
- **调试密钥迁移**: 调试密钥迁移功能
- **附加到正在运行的进程**: 附加调试器到运行中的进程
- **全栈调试**: 组合配置，同时启动前后端

### 任务配置 (tasks.json)

- **安装前端依赖**: `npm install`
- **构建前端**: `npm run build`
- **启动前端开发服务器**: `npm run dev`（后台任务）
- **运行后端服务器**: 构建前端后启动后端
- **运行后端服务器 (开发模式)**: 开发模式启动后端
- **前端代码检查**: ESLint 检查
- **前端格式化**: Prettier 格式化
- **Go 模块整理**: `go mod tidy`
- **清理构建文件**: 清理前端构建文件

### 工作区设置

- **文件嵌套**: 自动组织相关文件
- **Go 工具**: 自动更新、格式化、代码检查
- **前端工具**: ESLint、Prettier、TypeScript 支持
- **推荐扩展**: 自动推荐必要的 VSCode 扩展

## 🔧 开发工作流

### 1. 前端开发

```bash
# 启动前端开发服务器
cd web
npm run dev
```

- 热重载: ✅
- TypeScript 检查: ✅
- ESLint 检查: ✅
- 访问地址: http://localhost:5173

### 2. 后端开发

```bash
# 开发模式启动
go run -race ./main.go
```

- 竞态检测: ✅
- 调试信息: ✅
- 访问地址: http://localhost:8080

### 3. 全栈开发

使用 VSCode 的"全栈调试"配置，或者分别在两个终端启动前后端。

## 📝 代码规范

### Go 代码

- 使用 `goimports` 自动格式化和导入整理
- 使用 `golangci-lint` 进行代码检查
- 测试时启用竞态检测 (`-race`)

### 前端代码

- 使用 Prettier 格式化
- 使用 ESLint 检查代码质量
- TypeScript 严格模式

## 🐛 调试技巧

### 后端调试

1. 在 VSCode 中设置断点
2. 使用 F5 启动调试配置
3. 或者使用 `dlv` 命令行调试器

### 前端调试

1. 使用浏览器开发者工具
2. VSCode 中安装 Vue DevTools 扩展
3. 使用 Vue DevTools 浏览器扩展

### 网络调试

- 后端 API: http://localhost:8080
- 前端开发服务器: http://localhost:5173
- 使用 Postman 或 curl 测试 API

## 📦 构建和部署

### 开发构建

```bash
# 构建前端
cd web && npm run build

# 运行后端（包含前端静态文件）
go run ./main.go
```

### 生产构建

```bash
# 使用 Makefile
make run

# 或者手动
cd web && npm install && npm run build
go build -o gpt-load ./main.go
./gpt-load
```

## 🔑 密钥迁移

```bash
# 启用加密
go run ./main.go migrate-keys --to new-key

# 禁用加密
go run ./main.go migrate-keys --from old-key

# 更换密钥
go run ./main.go migrate-keys --from old-key --to new-key
```

## 📋 环境要求

- **Go**: 1.20+
- **Node.js**: 18+
- **npm**: 8+

## 🆘 常见问题

### 1. 前端构建失败

```bash
cd web
rm -rf node_modules package-lock.json
npm install
```

### 2. Go 模块问题

```bash
go clean -modcache
go mod download
go mod tidy
```

### 3. 端口冲突

- 前端默认端口: 5173
- 后端默认端口: 8080
- 可以在配置文件中修改

### 4. VSCode Go 扩展问题

1. 打开命令面板 (Ctrl+Shift+P)
2. 运行 "Go: Install/Update Tools"
3. 安装所有推荐的工具
