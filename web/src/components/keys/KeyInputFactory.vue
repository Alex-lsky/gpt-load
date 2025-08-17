<script setup lang="ts">
import { CHANNEL_CONFIGS } from "@/config/channels";
import { computed } from "vue";
import BedrockKeyInput from "./BedrockKeyInput.vue";
import StandardKeyInput from "./StandardKeyInput.vue";

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

const currentComponent = computed(() => {
  const config = CHANNEL_CONFIGS[props.channelType];
  return config?.keyInputComponent === "BedrockKeyInput" ? BedrockKeyInput : StandardKeyInput;
});

const handleUpdateValue = (value: string) => {
  emit("update:modelValue", value);
};

const handleValidate = (isValid: boolean) => {
  emit("validate", isValid);
};
</script>

<template>
  <component
    :is="currentComponent"
    :model-value="modelValue"
    :channel-type="channelType"
    :placeholder="placeholder"
    :disabled="disabled"
    :size="size"
    @update:model-value="handleUpdateValue"
    @validate="handleValidate"
  />
</template>
