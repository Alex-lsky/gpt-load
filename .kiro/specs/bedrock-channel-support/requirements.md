# Requirements Document

## Introduction

本功能旨在为 GPT-Load 代理服务添加 AWS Bedrock 渠道支持。AWS Bedrock 是亚马逊提供的托管服务，通过 Converse API 提供统一的多模型访问接口。Converse API 是 AWS Bedrock 的现代化接口，提供统一的请求/响应格式，支持所有 Bedrock 上的基础模型，包括 Anthropic Claude、Amazon Titan、Cohere Command、Meta Llama、Mistral 等。

AWS Bedrock Converse API 的主要特点：

- Runtime API 端点: `https://bedrock-runtime.{region}.amazonaws.com`
- 统一接口: `/converse` (非流式) 和 `/converse-stream` (流式)
- 认证方式: AWS Signature Version 4
- 统一的消息格式，类似 OpenAI 的 messages 结构
- 自动处理不同模型的格式转换
- 支持工具调用、系统提示等高级功能

## Requirements

### Requirement 1

**User Story:** 作为系统管理员，我希望能够配置 AWS Bedrock 渠道，以便用户可以通过代理服务访问 Bedrock 上的 AI 模型

#### Acceptance Criteria

1. WHEN 管理员在系统中创建新的渠道组时 THEN 系统 SHALL 支持选择"bedrock"作为渠道类型
2. WHEN 配置 Bedrock 渠道时 THEN 系统 SHALL 要求提供 AWS 访问凭证（Access Key ID、Secret Access Key、可选的 Session Token）
3. WHEN 配置 Bedrock 渠道时 THEN 系统 SHALL 要求指定 AWS 区域（如 us-east-1、us-west-2 等）
4. WHEN 配置 Bedrock 渠道时 THEN 系统 SHALL 自动构建正确的 Bedrock Runtime API 端点 URL 格式
5. IF 提供了无效的 AWS 凭证或区域 THEN 系统 SHALL 在验证时返回明确的错误信息

### Requirement 2

**User Story:** 作为开发者，我希望通过代理服务调用 Bedrock Converse API，使用统一的消息格式

#### Acceptance Criteria

1. WHEN 客户端向 Bedrock 渠道发送请求时 THEN 系统 SHALL 正确生成 AWS Signature Version 4 认证头
2. WHEN 请求路径为`/converse`时 THEN 系统 SHALL 处理为非流式调用
3. WHEN 请求路径为`/converse-stream`时 THEN 系统 SHALL 处理为流式调用
4. WHEN 请求包含统一的消息格式时 THEN 系统 SHALL 正确设置 Content-Type 和其他必需的 AWS 头部
5. WHEN Bedrock 返回错误响应时 THEN 系统 SHALL 解析并转换为统一的错误格式

### Requirement 3

**User Story:** 作为系统管理员，我希望能够验证 Bedrock API 密钥的有效性并获取可用模型列表

#### Acceptance Criteria

1. WHEN 管理员添加或更新 Bedrock API 密钥时 THEN 系统 SHALL 自动验证密钥的有效性
2. WHEN 验证 Bedrock 密钥时 THEN 系统 SHALL 使用 ListFoundationModels API 获取可用模型列表
3. WHEN 密钥验证成功时 THEN 系统 SHALL 返回成功状态和可用模型 ID 列表
4. WHEN 密钥验证失败时 THEN 系统 SHALL 返回具体的失败原因（如权限不足、区域不匹配、凭证无效等）
5. IF 验证过程中发生网络错误 THEN 系统 SHALL 区分网络问题和认证问题

### Requirement 4

**User Story:** 作为最终用户，我希望能够通过统一的接口访问 Bedrock 上的各种 AI 模型

#### Acceptance Criteria

1. WHEN 用户请求 Bedrock 渠道时 THEN 系统 SHALL 支持标准的 Bedrock 模型 ID 格式（如 anthropic.claude-3-5-sonnet-20241022-v2:0）
2. WHEN 用户发送 Converse API 请求时 THEN 系统 SHALL 支持统一的消息格式（messages 数组）
3. WHEN 请求包含系统提示、工具调用等参数时 THEN 系统 SHALL 透传所有兼容的参数到 Bedrock
4. WHEN 模型不可用或不存在时 THEN 系统 SHALL 返回 Bedrock 原始的错误信息
5. IF 用户权限不足访问特定模型 THEN 系统 SHALL 返回权限相关的错误信息

### Requirement 5

**User Story:** 作为开发者，我希望 Bedrock 渠道能够正确处理 Converse API 的流式和非流式响应

#### Acceptance Criteria

1. WHEN 请求使用`/converse`端点时 THEN 系统 SHALL 处理完整的 JSON 响应格式
2. WHEN 请求使用`/converse-stream`端点时 THEN 系统 SHALL 正确处理 Server-Sent Events 流
3. WHEN 处理流式响应时 THEN 系统 SHALL 保持连接直到流结束或客户端断开
4. WHEN 流式响应中包含不同事件类型时 THEN 系统 SHALL 正确传递所有事件（messageStart、contentBlockStart、contentBlockDelta、contentBlockStop、messageStop 等）
5. IF 流式连接意外中断 THEN 系统 SHALL 记录中断原因并清理资源

### Requirement 6

**User Story:** 作为系统管理员，我希望系统能够动态获取和缓存 Bedrock 可用模型列表

#### Acceptance Criteria

1. WHEN 系统启动或配置更新时 THEN 系统 SHALL 调用 ListFoundationModels API 获取最新的模型列表
2. WHEN 获取模型列表成功时 THEN 系统 SHALL 缓存模型信息包括模型 ID、名称、提供商等
3. WHEN 模型列表缓存过期时 THEN 系统 SHALL 自动刷新模型列表
4. WHEN 用户请求不存在的模型时 THEN 系统 SHALL 检查缓存的模型列表并返回适当的错误
5. IF 无法获取模型列表 THEN 系统 SHALL 记录错误但不阻止服务启动

### Requirement 7

**User Story:** 作为系统运维人员，我希望 Bedrock 渠道能够与现有的监控和日志系统集成

#### Acceptance Criteria

1. WHEN Bedrock 渠道处理请求时 THEN 系统 SHALL 记录请求和响应的关键信息到日志
2. WHEN 发生 AWS 签名错误时 THEN 系统 SHALL 记录详细的签名生成过程和错误信息
3. WHEN 渠道配置发生变化时 THEN 系统 SHALL 记录配置变更日志
4. WHEN 模型列表更新时 THEN 系统 SHALL 记录更新的模型数量和变化
5. IF 系统检测到配置过期 THEN 系统 SHALL 自动刷新渠道配置

### Requirement 8

**User Story:** 作为安全管理员，我希望 Bedrock 渠道遵循 AWS 安全最佳实践

#### Acceptance Criteria

1. WHEN 存储 AWS 凭证时 THEN 系统 SHALL 使用安全的加密存储方式
2. WHEN 传输请求到 Bedrock 时 THEN 系统 SHALL 强制使用 HTTPS 连接
3. WHEN 生成 AWS 签名时 THEN 系统 SHALL 使用当前时间戳和正确的签名算法防止重放攻击
4. WHEN 处理敏感信息时 THEN 系统 SHALL 避免在日志中记录 AWS 凭证信息
5. IF 检测到异常访问模式 THEN 系统 SHALL 记录安全事件日志

### Requirement 9

**User Story:** 作为开发者，我希望 Bedrock 渠道支持向后兼容的传统 invoke API

#### Acceptance Criteria

1. WHEN 用户请求传统的`/model/{model-id}/invoke`端点时 THEN 系统 SHALL 支持传统的模型特定格式
2. WHEN 用户请求传统的`/model/{model-id}/invoke-with-response-stream`端点时 THEN 系统 SHALL 支持传统的流式格式
3. WHEN 处理传统 API 请求时 THEN 系统 SHALL 保持原始的请求/响应格式不变
4. WHEN 同时支持 Converse API 和传统 API 时 THEN 系统 SHALL 根据请求路径自动选择处理方式
5. IF 传统 API 格式与 Converse API 冲突 THEN 系统 SHALL 优先使用请求路径确定的 API 类型

### Requirement 10

**User Story:** 作为系统管理员，我希望前端界面能够支持 Bedrock 渠道的特殊认证配置需求

#### Acceptance Criteria

1. WHEN 管理员选择 "bedrock" 渠道类型时 THEN 前端 SHALL 显示专门的 AWS 凭证输入组件
2. WHEN 输入 Bedrock API 密钥时 THEN 系统 SHALL 提供 Access Key ID、Secret Access Key、Session Token（可选）和 Region 的输入字段
3. WHEN 管理员填写 AWS 凭证时 THEN 系统 SHALL 实时验证输入格式的有效性
4. WHEN 保存 Bedrock 分组时 THEN 系统 SHALL 将 AWS 凭证转换为 JSON 格式存储
5. IF 用户输入的 AWS 凭证格式不正确 THEN 系统 SHALL 显示具体的错误提示信息

### Requirement 11

**User Story:** 作为前端开发者，我希望有一个可复用的 API 密钥输入组件来处理不同渠道的特殊需求

#### Acceptance Criteria

1. WHEN 创建新的渠道类型时 THEN 系统 SHALL 支持通过组件化方式扩展不同的密钥输入界面
2. WHEN 选择不同渠道类型时 THEN 系统 SHALL 自动切换到对应的密钥输入组件
3. WHEN 使用 OpenAI 渠道时 THEN 系统 SHALL 显示标准的单行密钥输入框
4. WHEN 使用 Bedrock 渠道时 THEN 系统 SHALL 显示多字段的 AWS 凭证输入表单
5. IF 将来需要支持其他特殊认证方式 THEN 系统 SHALL 能够轻松添加新的密钥输入组件

### Requirement 12

**User Story:** 作为系统管理员，我希望前端能够正确处理 Bedrock 渠道的默认配置和验证

#### Acceptance Criteria

1. WHEN 创建 Bedrock 分组时 THEN 系统 SHALL 提供合适的默认上游地址模板
2. WHEN 选择 Bedrock 渠道时 THEN 系统 SHALL 自动设置合适的测试模型默认值
3. WHEN 配置 Bedrock 分组时 THEN 系统 SHALL 隐藏不适用的配置选项（如验证端点）
4. WHEN 用户切换到 Bedrock 渠道时 THEN 系统 SHALL 更新相关字段的占位符文本和提示信息
5. IF Bedrock 渠道有特殊的配置要求 THEN 系统 SHALL 在界面上提供相应的说明和帮助文本
