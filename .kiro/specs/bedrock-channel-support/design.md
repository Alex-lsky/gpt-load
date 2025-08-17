# Design Document

## Overview

本设计文档描述了为 GPT-Load 代理服务添加 AWS Bedrock 渠道支持的技术实现方案。设计遵循现有的渠道架构模式，通过实现 ChannelProxy 接口来提供 AWS Bedrock 的透明代理功能。

AWS Bedrock 渠道将支持两种 API 模式：

1. **Converse API**（推荐）：统一的现代化接口，支持所有模型
2. **传统 Invoke API**：向后兼容的模型特定接口

## Architecture

### 系统架构图

```mermaid
graph TB
    Client[客户端应用] --> Proxy[GPT-Load代理服务]
    Proxy --> BedrockChannel[BedrockChannel]
    BedrockChannel --> AWSSigner[AWS签名器]
    BedrockChannel --> ModelCache[模型缓存]
    BedrockChannel --> BedrockAPI[AWS Bedrock API]

    subgraph "Bedrock Channel Components"
        BedrockChannel
        AWSSigner
        ModelCache
    end

    subgraph "AWS Bedrock APIs"
        BedrockAPI --> ConverseAPI[Converse API]
        BedrockAPI --> InvokeAPI[传统Invoke API]
        BedrockAPI --> ListModelsAPI[ListFoundationModels API]
    end
```

### 请求流程图

```mermaid
sequenceDiagram
    participant C as 客户端
    participant P as GPT-Load代理
    participant BC as BedrockChannel
    participant AS as AWS签名器
    participant BR as Bedrock API

    C->>P: HTTP请求 /proxy/bedrock/converse
    P->>BC: 处理请求
    BC->>BC: 检查请求类型（Converse/Invoke）
    BC->>AS: 生成AWS签名
    AS-->>BC: 返回签名头
    BC->>BR: 发送已签名请求
    BR-->>BC: 返回响应（流式/非流式）
    BC-->>P: 转发响应
    P-->>C: 返回最终响应
```

## Components and Interfaces

### 1. BedrockChannel 结构

```go
type BedrockChannel struct {
    *BaseChannel
    awsConfig    *AWSConfig
    signer       *AWSSigner
    modelCache   *ModelCache
    region       string
}

type AWSConfig struct {
    AccessKeyID     string
    SecretAccessKey string
    SessionToken    string // 可选
    Region          string
}
```

### 2. AWS 签名器组件

```go
type AWSSigner struct {
    accessKeyID     string
    secretAccessKey string
    sessionToken    string
    region          string
    service         string // "bedrock"
}

func (s *AWSSigner) SignRequest(req *http.Request) error
```

### 3. 模型缓存组件

```go
type ModelCache struct {
    models    map[string]*FoundationModel
    lastSync  time.Time
    ttl       time.Duration
    mutex     sync.RWMutex
}

type FoundationModel struct {
    ModelId           string `json:"modelId"`
    ModelName         string `json:"modelName"`
    ProviderName      string `json:"providerName"`
    InputModalities   []string `json:"inputModalities"`
    OutputModalities  []string `json:"outputModalities"`
    ResponseStreamingSupported bool `json:"responseStreamingSupported"`
}
```

### 4. 接口实现

#### ModifyRequest 方法

```go
func (ch *BedrockChannel) ModifyRequest(req *http.Request, apiKey *models.APIKey, group *models.Group) {
    // 1. 解析AWS凭证（从apiKey.KeyValue中提取JSON格式的凭证）
    // 2. 设置必需的AWS头部
    // 3. 调用AWS签名器生成签名
    // 4. 添加签名头到请求
}
```

#### IsStreamRequest 方法

```go
func (ch *BedrockChannel) IsStreamRequest(c *gin.Context, bodyBytes []byte) bool {
    // 1. 检查路径是否包含 "converse-stream" 或 "invoke-with-response-stream"
    // 2. 检查Accept头是否包含 "text/event-stream"
    // 3. 检查请求体中的stream参数
}
```

#### ValidateKey 方法

```go
func (ch *BedrockChannel) ValidateKey(ctx context.Context, key string) (bool, error) {
    // 1. 解析AWS凭证
    // 2. 调用ListFoundationModels API验证凭证
    // 3. 更新模型缓存
    // 4. 返回验证结果
}
```

### 5. 路径映射逻辑

```go
func (ch *BedrockChannel) BuildUpstreamURL(originalURL *url.URL, group *models.Group) (string, error) {
    // 输入: /proxy/bedrock/converse
    // 输出: https://bedrock-runtime.{region}.amazonaws.com/converse

    // 输入: /proxy/bedrock/model/anthropic.claude-3-sonnet/invoke
    // 输出: https://bedrock-runtime.{region}.amazonaws.com/model/anthropic.claude-3-sonnet/invoke
}
```

## Data Models

### 1. AWS 凭证存储格式

API Key 的 KeyValue 字段将存储 JSON 格式的 AWS 凭证：

```json
{
  "access_key_id": "AKIAIOSFODNN7EXAMPLE",
  "secret_access_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
  "session_token": "optional-session-token",
  "region": "us-east-1"
}
```

### 2. 上游配置格式

Group 的 Upstreams 字段将存储 Bedrock 端点配置：

```json
[
  {
    "url": "https://bedrock-runtime.us-east-1.amazonaws.com",
    "weight": 1
  }
]
```

### 3. 模型缓存数据结构

```go
type CachedModels struct {
    Models    []FoundationModel `json:"models"`
    UpdatedAt time.Time         `json:"updated_at"`
    Region    string            `json:"region"`
}
```

## Error Handling

### 1. AWS 签名错误处理

```go
type AWSSignatureError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Region  string `json:"region"`
}
```

### 2. 模型访问错误处理

```go
type ModelAccessError struct {
    ModelId string `json:"model_id"`
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

### 3. 错误映射表

| AWS 错误代码              | HTTP 状态码 | 统一错误消息             |
| ------------------------- | ----------- | ------------------------ |
| UnauthorizedOperation     | 401         | AWS 凭证无效或权限不足   |
| ValidationException       | 400         | 请求参数验证失败         |
| ResourceNotFoundException | 404         | 指定的模型不存在         |
| ThrottlingException       | 429         | 请求频率超限，请稍后重试 |
| InternalServerException   | 500         | AWS 服务内部错误         |

## Testing Strategy

### 1. 单元测试

- **AWS 签名器测试**：验证签名生成的正确性
- **模型缓存测试**：验证缓存的更新和过期逻辑
- **路径映射测试**：验证 URL 构建的正确性
- **错误处理测试**：验证各种错误场景的处理

### 2. 集成测试

- **端到端 API 调用测试**：使用真实的 AWS 凭证测试完整流程
- **流式响应测试**：验证 Server-Sent Events 的正确处理
- **模型验证测试**：测试不同模型的调用和响应

### 3. 性能测试

- **并发请求测试**：验证高并发下的性能表现
- **签名生成性能测试**：确保签名生成不成为性能瓶颈
- **模型缓存性能测试**：验证缓存的命中率和更新效率

## Implementation Details

### 1. 依赖管理

需要添加的 Go 模块依赖：

```go
// go.mod 新增依赖
require (
    github.com/aws/aws-sdk-go-v2 v1.24.0
    github.com/aws/aws-sdk-go-v2/config v1.26.1
    github.com/aws/aws-sdk-go-v2/service/bedrock v1.7.2
    github.com/aws/aws-sdk-go-v2/service/bedrockruntime v1.7.1
    github.com/aws/aws-sdk-go-v2/credentials v1.16.12
    github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.15.7
    github.com/aws/smithy-go v1.19.0
)
```

### 2. 配置管理

```go
// 在Group配置中支持Bedrock特定配置
type BedrockConfig struct {
    DefaultRegion        string `json:"default_region"`
    ModelCacheTTL        int    `json:"model_cache_ttl_minutes"`
    SignatureVersion     string `json:"signature_version"` // 默认 "v4"
    EnableModelCaching   bool   `json:"enable_model_caching"`
    ValidationModel      string `json:"validation_model"` // 用于验证的默认模型
}
```

### 3. 监控和日志

```go
// 添加Bedrock特定的日志字段
type BedrockRequestLog struct {
    AWSRegion     string `json:"aws_region"`
    ModelId       string `json:"model_id"`
    APIType       string `json:"api_type"` // "converse" 或 "invoke"
    SignatureTime int64  `json:"signature_time_ms"`
    IsStreaming   bool   `json:"is_streaming"`
}
```

### 4. 安全考虑

- **凭证加密**：AWS 凭证在数据库中使用 AES 加密存储
- **签名安全**：确保时间戳在合理范围内，防止重放攻击
- **日志安全**：避免在日志中记录敏感的 AWS 凭证信息
- **网络安全**：强制使用 HTTPS 连接到 AWS Bedrock

### 5. 性能优化

- **连接池复用**：复用 HTTP 连接到 AWS Bedrock
- **签名缓存**：对相同请求的签名进行短期缓存
- **模型缓存**：缓存模型列表减少 API 调用
- **异步验证**：后台异步验证密钥状态

## Frontend Components Design

### 1. API Key Input Component Architecture

```typescript
// 基础密钥输入组件接口
interface KeyInputComponent {
  modelValue: string;
  channelType: string;
  placeholder?: string;
  disabled?: boolean;
  validate?: () => boolean;
  getFormattedValue?: () => string;
}

// AWS 凭证输入组件数据结构
interface AWSCredentials {
  access_key_id: string;
  secret_access_key: string;
  session_token?: string;
  region: string;
}
```

### 2. 组件层次结构

```mermaid
graph TB
    GroupFormModal[GroupFormModal.vue] --> KeyInputFactory[KeyInputFactory.vue]
    KeyInputFactory --> StandardKeyInput[StandardKeyInput.vue]
    KeyInputFactory --> BedrockKeyInput[BedrockKeyInput.vue]
    KeyInputFactory --> AzureKeyInput[AzureKeyInput.vue]

    KeyCreateDialog[KeyCreateDialog.vue] --> KeyInputFactory

    BedrockKeyInput --> AWSCredentialsForm[AWS凭证表单]
    StandardKeyInput --> SimpleTextInput[简单文本输入]
```

### 3. Bedrock Key Input Component

```vue
<!-- BedrockKeyInput.vue -->
<template>
  <div class="bedrock-key-input">
    <div class="aws-credentials-form">
      <n-form-item label="Access Key ID" required>
        <n-input
          v-model:value="credentials.access_key_id"
          placeholder="AKIAIOSFODNN7EXAMPLE"
          @input="updateValue"
        />
      </n-form-item>

      <n-form-item label="Secret Access Key" required>
        <n-input
          v-model:value="credentials.secret_access_key"
          type="password"
          placeholder="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
          @input="updateValue"
        />
      </n-form-item>

      <n-form-item label="Session Token (可选)">
        <n-input
          v-model:value="credentials.session_token"
          type="password"
          placeholder="临时会话令牌（可选）"
          @input="updateValue"
        />
      </n-form-item>

      <n-form-item label="AWS Region" required>
        <n-select
          v-model:value="credentials.region"
          :options="regionOptions"
          placeholder="选择 AWS 区域"
          @update:value="updateValue"
        />
      </n-form-item>
    </div>
  </div>
</template>
```

### 4. 渠道类型配置映射

```typescript
// 渠道配置映射
const CHANNEL_CONFIGS = {
  openai: {
    keyInputComponent: "StandardKeyInput",
    defaultUpstream: "https://api.openai.com",
    defaultTestModel: "gpt-4o-mini",
    validationEndpoint: "/v1/chat/completions",
    showValidationEndpoint: true,
  },
  anthropic: {
    keyInputComponent: "StandardKeyInput",
    defaultUpstream: "https://api.anthropic.com",
    defaultTestModel: "claude-3-haiku-20240307",
    validationEndpoint: "/v1/messages",
    showValidationEndpoint: true,
  },
  gemini: {
    keyInputComponent: "StandardKeyInput",
    defaultUpstream: "https://generativelanguage.googleapis.com",
    defaultTestModel: "gemini-1.5-flash",
    validationEndpoint: "",
    showValidationEndpoint: false,
  },
  bedrock: {
    keyInputComponent: "BedrockKeyInput",
    defaultUpstream: "https://bedrock-runtime.us-east-1.amazonaws.com",
    defaultTestModel: "anthropic.claude-3-haiku-20240307-v1:0",
    validationEndpoint: "",
    showValidationEndpoint: false,
  },
};
```

## Migration and Deployment

### 1. 数据库迁移

无需额外的数据库表，使用现有的 Group 和 APIKey 表结构。

### 2. 配置迁移

```sql
-- 添加bedrock渠道类型支持
INSERT INTO system_settings (setting_key, setting_value, description)
VALUES ('supported_channel_types', '["openai","gemini","anthropic","bedrock"]', '支持的渠道类型列表');
```

### 3. 前端组件迁移

1. 创建可复用的密钥输入组件架构
2. 重构现有的 GroupFormModal 和 KeyCreateDialog
3. 添加 Bedrock 特定的输入组件
4. 更新类型定义以支持新的渠道类型

### 4. 部署步骤

1. 更新 Go 依赖
2. 部署新的前端组件
3. 部署新的后端代码版本
4. 在管理界面中创建 Bedrock 渠道组
5. 配置 AWS 凭证和测试模型
6. 验证功能正常工作

### 5. 回滚计划

如果需要回滚：

1. 停止创建新的 Bedrock 渠道组
2. 将现有 Bedrock 组标记为不可用
3. 回滚到之前的代码版本
4. 清理相关配置数据

## Frontend Component Detailed Design

### 1. 渠道类型配置映射

```typescript
// 渠道配置映射
const CHANNEL_CONFIGS = {
  openai: {
    keyInputComponent: "StandardKeyInput",
    defaultUpstream: "https://api.openai.com",
    defaultTestModel: "gpt-4o-mini",
    validationEndpoint: "/v1/chat/completions",
    showValidationEndpoint: true,
    upstreamPlaceholder: "https://api.openai.com",
    testModelPlaceholder: "gpt-4o-mini",
  },
  anthropic: {
    keyInputComponent: "StandardKeyInput",
    defaultUpstream: "https://api.anthropic.com",
    defaultTestModel: "claude-3-haiku-20240307",
    validationEndpoint: "/v1/messages",
    showValidationEndpoint: true,
    upstreamPlaceholder: "https://api.anthropic.com",
    testModelPlaceholder: "claude-3-haiku-20240307",
  },
  gemini: {
    keyInputComponent: "StandardKeyInput",
    defaultUpstream: "https://generativelanguage.googleapis.com",
    defaultTestModel: "gemini-1.5-flash",
    validationEndpoint: "",
    showValidationEndpoint: false,
    upstreamPlaceholder: "https://generativelanguage.googleapis.com",
    testModelPlaceholder: "gemini-1.5-flash",
  },
  bedrock: {
    keyInputComponent: "BedrockKeyInput",
    defaultUpstream: "https://bedrock-runtime.us-east-1.amazonaws.com",
    defaultTestModel: "anthropic.claude-3-haiku-20240307-v1:0",
    validationEndpoint: "",
    showValidationEndpoint: false,
    upstreamPlaceholder: "https://bedrock-runtime.{region}.amazonaws.com",
    testModelPlaceholder: "anthropic.claude-3-haiku-20240307-v1:0",
    regions: [
      { label: "US East (N. Virginia)", value: "us-east-1" },
      { label: "US West (Oregon)", value: "us-west-2" },
      { label: "Europe (Ireland)", value: "eu-west-1" },
      { label: "Asia Pacific (Tokyo)", value: "ap-northeast-1" },
      { label: "Asia Pacific (Singapore)", value: "ap-southeast-1" },
    ],
  },
};
```

### 2. 组件接口定义

```typescript
// 基础密钥输入组件接口
interface KeyInputProps {
  modelValue: string;
  channelType: string;
  placeholder?: string;
  disabled?: boolean;
  size?: "small" | "medium" | "large";
}

interface KeyInputEmits {
  (e: "update:modelValue", value: string): void;
  (e: "validate", isValid: boolean): void;
}

// AWS 凭证数据结构
interface AWSCredentials {
  access_key_id: string;
  secret_access_key: string;
  session_token?: string;
  region: string;
}

// 渠道配置接口
interface ChannelConfig {
  keyInputComponent: string;
  defaultUpstream: string;
  defaultTestModel: string;
  validationEndpoint: string;
  showValidationEndpoint: boolean;
  upstreamPlaceholder: string;
  testModelPlaceholder: string;
  regions?: Array<{ label: string; value: string }>;
}
```

### 3. KeyInputFactory 组件

```vue
<!-- KeyInputFactory.vue -->
<template>
  <component
    :is="currentComponent"
    v-model="modelValue"
    :channel-type="channelType"
    :placeholder="placeholder"
    :disabled="disabled"
    :size="size"
    @validate="handleValidate"
  />
</template>

<script setup lang="ts">
import { computed } from "vue";
import StandardKeyInput from "./StandardKeyInput.vue";
import BedrockKeyInput from "./BedrockKeyInput.vue";
import { CHANNEL_CONFIGS } from "@/config/channels";

const props = defineProps<KeyInputProps>();
const emit = defineEmits<KeyInputEmits>();

const currentComponent = computed(() => {
  const config = CHANNEL_CONFIGS[props.channelType];
  return config?.keyInputComponent === "BedrockKeyInput"
    ? BedrockKeyInput
    : StandardKeyInput;
});

const handleValidate = (isValid: boolean) => {
  emit("validate", isValid);
};
</script>
```

### 4. BedrockKeyInput 组件

```vue
<!-- BedrockKeyInput.vue -->
<template>
  <div class="bedrock-key-input">
    <div class="aws-credentials-grid">
      <n-form-item label="Access Key ID" required>
        <n-input
          v-model:value="credentials.access_key_id"
          placeholder="AKIAIOSFODNN7EXAMPLE"
          :disabled="disabled"
          @input="updateValue"
          @blur="validateField('access_key_id')"
        />
        <div v-if="errors.access_key_id" class="error-text">
          {{ errors.access_key_id }}
        </div>
      </n-form-item>

      <n-form-item label="Secret Access Key" required>
        <n-input
          v-model:value="credentials.secret_access_key"
          type="password"
          show-password-on="click"
          placeholder="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
          :disabled="disabled"
          @input="updateValue"
          @blur="validateField('secret_access_key')"
        />
        <div v-if="errors.secret_access_key" class="error-text">
          {{ errors.secret_access_key }}
        </div>
      </n-form-item>

      <n-form-item label="Session Token (可选)">
        <n-input
          v-model:value="credentials.session_token"
          type="password"
          show-password-on="click"
          placeholder="临时会话令牌（可选）"
          :disabled="disabled"
          @input="updateValue"
        />
      </n-form-item>

      <n-form-item label="AWS Region" required>
        <n-select
          v-model:value="credentials.region"
          :options="regionOptions"
          placeholder="选择 AWS 区域"
          :disabled="disabled"
          @update:value="updateValue"
          @blur="validateField('region')"
        />
        <div v-if="errors.region" class="error-text">
          {{ errors.region }}
        </div>
      </n-form-item>
    </div>

    <div class="help-text">
      <n-alert type="info" :show-icon="false">
        AWS Bedrock 需要有效的 AWS 凭证。请确保您的 Access Key 具有
        bedrock:InvokeModel 和 bedrock:ListFoundationModels 权限。
      </n-alert>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed, watch } from "vue";
import { CHANNEL_CONFIGS } from "@/config/channels";

const props = defineProps<KeyInputProps>();
const emit = defineEmits<KeyInputEmits>();

const credentials = reactive<AWSCredentials>({
  access_key_id: "",
  secret_access_key: "",
  session_token: "",
  region: "us-east-1",
});

const errors = reactive({
  access_key_id: "",
  secret_access_key: "",
  region: "",
});

const regionOptions = computed(() => {
  return CHANNEL_CONFIGS.bedrock.regions || [];
});

// 验证单个字段
const validateField = (field: keyof AWSCredentials) => {
  switch (field) {
    case "access_key_id":
      errors.access_key_id = !credentials.access_key_id
        ? "Access Key ID 不能为空"
        : !/^AKIA[0-9A-Z]{16}$/.test(credentials.access_key_id)
        ? "Access Key ID 格式不正确"
        : "";
      break;
    case "secret_access_key":
      errors.secret_access_key = !credentials.secret_access_key
        ? "Secret Access Key 不能为空"
        : credentials.secret_access_key.length !== 40
        ? "Secret Access Key 长度应为 40 个字符"
        : "";
      break;
    case "region":
      errors.region = !credentials.region ? "AWS Region 不能为空" : "";
      break;
  }
};

// 验证所有字段
const validateAll = () => {
  validateField("access_key_id");
  validateField("secret_access_key");
  validateField("region");

  const isValid =
    !errors.access_key_id && !errors.secret_access_key && !errors.region;
  emit("validate", isValid);
  return isValid;
};

// 更新值并发出事件
const updateValue = () => {
  const jsonValue = JSON.stringify(credentials);
  emit("update:modelValue", jsonValue);
  validateAll();
};

// 监听外部值变化
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue && newValue !== JSON.stringify(credentials)) {
      try {
        const parsed = JSON.parse(newValue);
        Object.assign(credentials, parsed);
      } catch (e) {
        // 忽略解析错误，保持当前状态
      }
    }
  },
  { immediate: true }
);
</script>

<style scoped>
.bedrock-key-input {
  width: 100%;
}

.aws-credentials-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 16px;
}

.aws-credentials-grid :deep(.n-form-item) {
  margin-bottom: 0;
}

.error-text {
  color: #d03050;
  font-size: 12px;
  margin-top: 4px;
}

.help-text {
  margin-top: 12px;
}

.help-text :deep(.n-alert) {
  --n-border-radius: 6px;
}
</style>
```

### 5. StandardKeyInput 组件

```vue
<!-- StandardKeyInput.vue -->
<template>
  <div class="standard-key-input">
    <n-input
      v-model:value="inputValue"
      :type="inputType"
      :placeholder="placeholder || '请输入 API 密钥'"
      :disabled="disabled"
      :size="size"
      :show-password-on="inputType === 'password' ? 'click' : undefined"
      @input="handleInput"
      @blur="handleBlur"
    />
    <div v-if="errorMessage" class="error-text">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";

const props = defineProps<KeyInputProps>();
const emit = defineEmits<KeyInputEmits>();

const inputValue = ref("");
const errorMessage = ref("");

const inputType = computed(() => {
  // 对于敏感的 API 密钥使用密码类型
  return props.channelType === "anthropic" ? "password" : "text";
});

const validateKey = (value: string) => {
  if (!value.trim()) {
    errorMessage.value = "API 密钥不能为空";
    return false;
  }

  // 基础格式验证
  switch (props.channelType) {
    case "openai":
      if (!value.startsWith("sk-")) {
        errorMessage.value = 'OpenAI API 密钥应以 "sk-" 开头';
        return false;
      }
      break;
    case "anthropic":
      if (!value.startsWith("sk-ant-")) {
        errorMessage.value = 'Anthropic API 密钥应以 "sk-ant-" 开头';
        return false;
      }
      break;
    case "gemini":
      if (value.length < 20) {
        errorMessage.value = "Gemini API 密钥长度不足";
        return false;
      }
      break;
  }

  errorMessage.value = "";
  return true;
};

const handleInput = (value: string) => {
  inputValue.value = value;
  emit("update:modelValue", value);

  // 清除之前的错误信息
  if (errorMessage.value) {
    errorMessage.value = "";
  }
};

const handleBlur = () => {
  const isValid = validateKey(inputValue.value);
  emit("validate", isValid);
};

// 监听外部值变化
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue !== inputValue.value) {
      inputValue.value = newValue || "";
    }
  },
  { immediate: true }
);
</script>

<style scoped>
.standard-key-input {
  width: 100%;
}

.error-text {
  color: #d03050;
  font-size: 12px;
  margin-top: 4px;
}
</style>
```

### 6. 类型定义更新

```typescript
// types/models.ts 更新
export interface Group {
  id?: number;
  name: string;
  display_name: string;
  description: string;
  sort: number;
  test_model: string;
  channel_type: "openai" | "gemini" | "anthropic" | "bedrock"; // 添加 bedrock
  upstreams: UpstreamInfo[];
  validation_endpoint: string;
  config: Record<string, unknown>;
  api_keys?: APIKey[];
  endpoint?: string;
  param_overrides: Record<string, unknown>;
  proxy_keys: string;
  created_at?: string;
  updated_at?: string;
}

// 新增 AWS 凭证类型
export interface AWSCredentials {
  access_key_id: string;
  secret_access_key: string;
  session_token?: string;
  region: string;
}

// 渠道配置类型
export interface ChannelConfig {
  keyInputComponent: string;
  defaultUpstream: string;
  defaultTestModel: string;
  validationEndpoint: string;
  showValidationEndpoint: boolean;
  upstreamPlaceholder: string;
  testModelPlaceholder: string;
  regions?: Array<{ label: string; value: string }>;
}
```
