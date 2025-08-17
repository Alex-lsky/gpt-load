<script setup lang="ts">
import { CHANNEL_CONFIGS } from "@/config/channels";
import { NAlert, NFormItem, NInput, NSelect } from "naive-ui";
import { computed, reactive, watch } from "vue";

interface Props {
  modelValue: string;
  channelType: string;
  placeholder?: string;
  disabled?: boolean;
  size?: "small" | "medium" | "large";
}

interface Emits {
  (e: "update:modelValue", value: string): void;
  (e: "validate", isValid: boolean): void;
}

interface AWSCredentials {
  access_key_id: string;
  secret_access_key: string;
  session_token?: string;
  region: string;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

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

  const isValid = !errors.access_key_id && !errors.secret_access_key && !errors.region;
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
  newValue => {
    if (newValue && newValue !== JSON.stringify(credentials)) {
      try {
        const parsed = JSON.parse(newValue);
        Object.assign(credentials, parsed);
      } catch (_e) {
        // 忽略解析错误，保持当前状态
      }
    }
  },
  { immediate: true }
);
</script>

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
        AWS Bedrock 需要有效的 AWS 凭证。请确保您的 Access Key 具有 bedrock:InvokeModel 和
        bedrock:ListFoundationModels 权限。
      </n-alert>
    </div>
  </div>
</template>

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
