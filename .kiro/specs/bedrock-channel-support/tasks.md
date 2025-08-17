# Implementation Plan

- [x] 1. 添加 AWS SDK 依赖和基础结构

  - 添加 AWS SDK Go v2 依赖到 go.mod 文件
  - 创建 internal/channel/bedrock_channel.go 文件
  - 实现基础的 BedrockChannel 结构体和构造函数
  - 在 init() 函数中注册 "bedrock" 渠道类型到工厂
  - _Requirements: 1.1_

- [ ] 2. 实现 AWS 凭证解析和签名器
- [ ] 2.1 创建 AWS 凭证解析功能

  - 实现从 APIKey.KeyValue 解析 JSON 格式 AWS 凭证的功能
  - 添加凭证格式验证和错误处理
  - 支持 Access Key ID、Secret Access Key、Session Token 和 Region
  - _Requirements: 1.2, 1.5_

- [ ] 2.2 实现 AWS Signature Version 4 签名器

  - 创建 AWSSigner 结构体和签名生成方法
  - 实现 AWS Signature Version 4 签名算法
  - 处理请求头部的签名和时间戳设置
  - _Requirements: 2.1, 8.3_

- [ ] 3. 实现 BedrockChannel 核心接口方法
- [ ] 3.1 实现 ModifyRequest 方法

  - 解析 AWS 凭证并配置签名器
  - 生成 AWS Signature Version 4 认证头
  - 设置必需的 AWS 请求头部（Content-Type、Host 等）
  - _Requirements: 2.1, 2.4_

- [ ] 3.2 实现 BuildUpstreamURL 方法

  - 实现请求路径到 Bedrock API 端点的映射逻辑
  - 支持 Converse API 路径转换（/converse, /converse-stream）
  - 支持传统 Invoke API 路径转换（/model/{model-id}/invoke）
  - 处理不同 AWS 区域的端点构建
  - _Requirements: 2.3, 9.1, 9.2_

- [ ] 3.3 实现 IsStreamRequest 方法

  - 检测 Converse Stream API 请求（/converse-stream）
  - 检测传统流式 API 请求（/invoke-with-response-stream）
  - 支持通过 Accept 头和请求体参数检测流式请求
  - _Requirements: 5.2, 9.4_

- [ ] 4. 实现模型验证和缓存
- [ ] 4.1 创建模型缓存结构

  - 定义 ModelCache 结构体和 FoundationModel 数据结构
  - 实现模型缓存的基本操作（存储、获取、过期检查）
  - 添加线程安全的并发访问支持
  - _Requirements: 6.1, 6.2_

- [ ] 4.2 实现 ListFoundationModels API 调用

  - 集成 AWS Bedrock SDK 调用 ListFoundationModels API
  - 实现模型列表的获取和解析逻辑
  - 添加错误处理和重试机制
  - _Requirements: 3.2, 3.3, 6.2_

- [ ] 4.3 实现 ValidateKey 方法

  - 使用 ListFoundationModels API 验证 AWS 凭证
  - 更新模型缓存并返回验证结果
  - 实现详细的错误分类和消息处理
  - _Requirements: 3.1, 3.4, 3.5_

- [ ] 5. 实现错误处理和日志记录
- [ ] 5.1 创建 Bedrock 特定的错误处理

  - 定义 AWS 错误类型结构（AWSSignatureError、ModelAccessError）
  - 实现 AWS 错误代码到 HTTP 状态码的映射
  - 创建统一的错误消息格式转换器
  - _Requirements: 4.4, 4.5_

- [ ] 5.2 实现日志记录和监控

  - 添加 Bedrock 特定的请求日志字段
  - 实现 AWS 签名错误的详细日志记录
  - 确保敏感信息不被记录到日志中
  - _Requirements: 7.1, 7.2, 8.4_

- [ ] 6. 创建单元测试
- [ ] 6.1 创建 AWS 签名器单元测试

  - 测试签名生成的正确性和一致性
  - 测试不同凭证格式的处理
  - 测试签名错误场景的处理
  - _Requirements: 2.1, 8.3_

- [ ] 6.2 创建 BedrockChannel 接口方法测试

  - 测试 ModifyRequest 方法的签名生成
  - 测试 BuildUpstreamURL 方法的路径映射
  - 测试 IsStreamRequest 方法的流式检测
  - 测试 ValidateKey 方法的凭证验证
  - _Requirements: 2.1, 2.3, 5.2, 3.1_

- [ ] 6.3 创建模型缓存单元测试

  - 测试缓存的增删改查操作
  - 测试缓存过期和自动刷新逻辑
  - 测试并发访问的线程安全性
  - _Requirements: 6.1, 6.2_

- [ ] 7. 创建集成测试
- [ ] 7.1 创建端到端 API 调用测试

  - 使用模拟的 AWS 凭证测试完整的请求流程
  - 测试 Converse API 的非流式和流式调用
  - 测试传统 Invoke API 的兼容性
  - _Requirements: 2.2, 2.3, 9.1, 9.2_

- [ ] 7.2 创建错误场景集成测试

  - 测试无效凭证的错误处理
  - 测试网络错误和超时的处理
  - 测试模型不存在等业务错误的处理
  - _Requirements: 3.4, 3.5, 4.4, 4.5_

- [ ] 8. 创建前端密钥输入组件架构
- [x] 8.1 创建可复用的密钥输入组件基础架构

  - 创建 KeyInputFactory.vue 组件工厂
  - 定义 KeyInputComponent 接口和基础类型
  - 创建渠道配置映射对象 CHANNEL_CONFIGS
  - 实现组件动态加载和切换逻辑
  - _Requirements: 11.1, 11.2_

- [x] 8.2 重构现有的标准密钥输入组件

  - 创建 StandardKeyInput.vue 组件用于 OpenAI/Anthropic/Gemini
  - 从 GroupFormModal 和 KeyCreateDialog 中提取密钥输入逻辑
  - 实现标准的单行文本输入界面
  - 添加输入验证和格式化功能
  - _Requirements: 11.3_

- [ ] 9. 实现 Bedrock 专用前端组件
- [x] 9.1 创建 BedrockKeyInput.vue 组件

  - 实现 AWS 凭证的多字段输入表单
  - 添加 Access Key ID、Secret Access Key、Session Token、Region 输入字段
  - 实现 AWS 区域选择下拉菜单
  - 添加实时输入验证和格式检查
  - _Requirements: 10.1, 10.2, 10.3_

- [x] 9.2 实现 AWS 凭证的 JSON 序列化

  - 实现 AWS 凭证对象到 JSON 字符串的转换
  - 实现 JSON 字符串到 AWS 凭证对象的解析
  - 添加凭证格式验证和错误处理
  - 确保与后端 API 密钥存储格式兼容
  - _Requirements: 10.4, 10.5_

- [ ] 10. 更新分组管理界面
- [x] 10.1 更新 GroupFormModal.vue 支持 Bedrock

  - 集成 KeyInputFactory 组件到分组表单
  - 更新渠道类型选项包含 "bedrock"
  - 实现 Bedrock 渠道的默认配置设置
  - 更新占位符文本和帮助提示信息
  - _Requirements: 12.1, 12.2, 12.4_

- [x] 10.2 更新 KeyCreateDialog.vue 支持 Bedrock

  - 集成新的密钥输入组件架构
  - 支持 Bedrock 渠道的多字段密钥输入
  - 实现密钥格式验证和错误提示
  - 确保与现有密钥创建流程兼容
  - _Requirements: 10.1, 10.2_

- [ ] 11. 更新类型定义和配置
- [x] 11.1 更新 TypeScript 类型定义

  - 在 types/models.ts 中添加 "bedrock" 渠道类型
  - 定义 AWSCredentials 接口类型
  - 更新 Group 接口支持 Bedrock 配置
  - 添加渠道配置映射的类型定义
  - _Requirements: 10.1, 11.1_

- [x] 11.2 实现渠道特定的配置逻辑

  - 实现 Bedrock 渠道的默认值设置
  - 隐藏 Bedrock 不适用的配置选项（验证端点）
  - 更新上游地址和测试模型的默认值
  - 添加 Bedrock 特定的帮助文本和说明
  - _Requirements: 12.3, 12.5_

- [ ] 12. 创建前端组件测试
- [ ] 12.1 创建密钥输入组件单元测试

  - 测试 KeyInputFactory 的组件切换逻辑
  - 测试 StandardKeyInput 的基础功能
  - 测试 BedrockKeyInput 的 AWS 凭证输入和验证
  - 测试组件间的数据传递和格式转换
  - _Requirements: 11.1, 11.2, 10.2_

- [ ] 12.2 创建分组管理界面集成测试

  - 测试 GroupFormModal 中 Bedrock 渠道的创建和编辑
  - 测试 KeyCreateDialog 中 Bedrock 密钥的添加
  - 测试渠道类型切换时的界面更新
  - 测试表单验证和错误处理
  - _Requirements: 12.1, 12.2, 10.1_
