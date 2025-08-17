<script setup lang="ts">
import { computed, ref, watch } from "vue";

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

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

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
  newValue => {
    if (newValue !== inputValue.value) {
      inputValue.value = newValue || "";
    }
  },
  { immediate: true }
);
</script>

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
