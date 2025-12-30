# AWS Bedrock 渠道支持需求文档

## 简介

为 GPT-Load 项目添加 AWS Bedrock 渠道支持，使用 Bedrock 的 Converse API 形式接入，同时支持最新的 API Key 认证和传统的 IAM 用户认证两种方式，并在前端提供认证方式切换功能。

## 术语表

- **GPT-Load 系统**: 现有的 AI API 透明代理服务系统
- **Bedrock 渠道**: 新增的 AWS Bedrock 服务渠道实现
- **Converse API**: AWS Bedrock 提供的统一对话 API 接口
- **API Key 认证**: AWS Bedrock 最新支持的 API 密钥认证方式
- **IAM 认证**: 传统的 AWS 身份和访问管理认证方式
- **渠道工厂**: 系统中负责创建和管理不同渠道实例的工厂类
- **前端管理界面**: Vue 3 构建的 Web 管理控制台

## 需求

### 需求 1

**用户故事:** 作为系统管理员，我希望能够创建 AWS Bedrock 类型的渠道组，以便将 AI 请求代理到 AWS Bedrock 服务

#### 验收标准

1. 当管理员在前端创建新分组时，THE GPT-Load 系统 SHALL 在渠道类型选项中提供"bedrock"选项
2. 当管理员选择 bedrock 渠道类型时，THE GPT-Load 系统 SHALL 显示 Bedrock 特定的配置选项
3. 当管理员保存 Bedrock 分组配置时，THE GPT-Load 系统 SHALL 验证配置的有效性并存储到数据库
4. 当系统启动时，THE GPT-Load 系统 SHALL 自动注册 Bedrock 渠道到渠道工厂

### 需求 2

**用户故事:** 作为系统管理员，我希望能够选择 API Key 或 IAM 两种认证方式，以便根据不同的部署环境选择合适的认证方法

#### 验收标准

1. 当管理员配置 Bedrock 分组时，THE GPT-Load 系统 SHALL 提供认证方式选择器，包含"API Key"和"IAM"两个选项
2. 当管理员选择 API Key 认证时，THE GPT-Load 系统 SHALL 显示 API Key 输入字段和区域选择器
3. 当管理员选择 IAM 认证时，THE GPT-Load 系统 SHALL 显示 Access Key ID、Secret Access Key、Session Token（可选）和区域选择器
4. 当认证方式发生变化时，THE GPT-Load 系统 SHALL 动态切换显示对应的认证配置字段

### 需求 3

**用户故事:** 作为开发者，我希望系统能够使用 Converse API 与 Bedrock 服务通信，以便获得统一的 API 体验

#### 验收标准

1. 当接收到代理请求时，THE Bedrock 渠道 SHALL 将请求转换为 Converse API 格式
2. 当调用 Bedrock 服务时，THE Bedrock 渠道 SHALL 使用配置的认证方式进行身份验证
3. 当收到 Bedrock 响应时，THE Bedrock 渠道 SHALL 将响应转换为标准格式返回给客户端
4. 当请求包含流式标识时，THE Bedrock 渠道 SHALL 支持流式响应处理

### 需求 4

**用户故事:** 作为系统管理员，我希望能够验证 Bedrock API 密钥的有效性，以便确保配置的密钥可以正常工作

#### 验收标准

1. 当管理员添加 Bedrock API 密钥时，THE GPT-Load 系统 SHALL 提供密钥验证功能
2. 当执行密钥验证时，THE Bedrock 渠道 SHALL 使用配置的认证方式调用 Bedrock 测试接口
3. 当验证成功时，THE GPT-Load 系统 SHALL 将密钥状态标记为有效
4. 当验证失败时，THE GPT-Load 系统 SHALL 显示具体的错误信息并将密钥标记为无效

### 需求 5

**用户故事:** 作为系统用户，我希望能够通过代理接口调用 Bedrock 模型，以便在应用中使用 AWS Bedrock 的 AI 能力

#### 验收标准

1. 当客户端发送请求到 Bedrock 代理端点时，THE GPT-Load 系统 SHALL 正确路由请求到 Bedrock 渠道
2. 当处理 Bedrock 请求时，THE Bedrock 渠道 SHALL 从请求体中提取模型名称
3. 当转发请求到 Bedrock 时，THE Bedrock 渠道 SHALL 构建正确的上游 URL
4. 当请求完成时，THE GPT-Load 系统 SHALL 记录请求日志包含 Bedrock 特定信息

### 需求 6

**用户故事:** 作为系统管理员，我希望前端界面能够提供友好的 Bedrock 配置体验，以便轻松管理 Bedrock 渠道

#### 验收标准

1. 当管理员访问分组管理页面时，THE 前端管理界面 SHALL 在渠道类型下拉菜单中显示"AWS Bedrock"选项
2. 当选择 Bedrock 渠道类型时，THE 前端管理界面 SHALL 显示专门的 Bedrock 配置表单
3. 当配置认证信息时，THE 前端管理界面 SHALL 提供认证方式切换开关和对应的输入字段
4. 当保存配置时，THE 前端管理界面 SHALL 验证必填字段并提供清晰的错误提示
