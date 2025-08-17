<script setup lang="ts">
import { keysApi } from "@/api/keys";
import type { Group } from "@/types/models";
import { appState } from "@/utils/app-state";
import { Close } from "@vicons/ionicons5";
import { NButton, NCard, NIcon, NInput, NModal } from "naive-ui";
import { computed, ref, watch } from "vue";
import KeyInputFactory from "./KeyInputFactory.vue";

interface Props {
  show: boolean;
  groupId: number;
  groupName?: string;
  group?: Group;
}

interface Emits {
  (e: "update:show", value: boolean): void;
  (e: "success"): void;
}

const props = defineProps<Props>();

const emit = defineEmits<Emits>();

const loading = ref(false);
const keysText = ref("");
const singleKeyValue = ref("");
const isValidKey = ref(false);

// 判断是否为 Bedrock 渠道（需要特殊的密钥输入）
const isBedrock = computed(() => props.group?.channel_type === "bedrock");

// 根据渠道类型决定使用哪种输入方式
const useKeyInputFactory = computed(() => isBedrock.value);

// 监听弹窗显示状态
watch(
  () => props.show,
  show => {
    if (show) {
      resetForm();
    }
  }
);

// 重置表单
function resetForm() {
  keysText.value = "";
  singleKeyValue.value = "";
  isValidKey.value = false;
}

// 关闭弹窗
function handleClose() {
  emit("update:show", false);
}

// 处理密钥验证
function handleKeyValidate(valid: boolean) {
  isValidKey.value = valid;
}

// 提交表单
async function handleSubmit() {
  if (loading.value) {
    return;
  }

  // 对于 Bedrock，使用单个密钥值
  if (useKeyInputFactory.value) {
    if (!singleKeyValue.value.trim() || !isValidKey.value) {
      return;
    }
  } else {
    // 对于其他渠道，使用多行文本
    if (!keysText.value.trim()) {
      return;
    }
  }

  try {
    loading.value = true;

    const keyData = useKeyInputFactory.value ? singleKeyValue.value : keysText.value;
    await keysApi.addKeysAsync(props.groupId, keyData);
    resetForm();
    handleClose();
    window.$message.success("密钥导入任务已开始，请稍后在下方查看进度。");
    appState.taskPollingTrigger++;
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <n-modal :show="show" @update:show="handleClose" class="form-modal">
    <n-card
      style="width: 800px"
      :title="`为 ${groupName || '当前分组'} 添加密钥`"
      :bordered="false"
      size="huge"
      role="dialog"
      aria-modal="true"
    >
      <template #header-extra>
        <n-button quaternary circle @click="handleClose">
          <template #icon>
            <n-icon :component="Close" />
          </template>
        </n-button>
      </template>

      <!-- Bedrock 渠道使用专用的密钥输入组件 -->
      <div v-if="useKeyInputFactory" style="margin-top: 20px">
        <key-input-factory
          v-model="singleKeyValue"
          :channel-type="group?.channel_type || 'openai'"
          @validate="handleKeyValidate"
        />
      </div>

      <!-- 其他渠道使用传统的多行文本输入 -->
      <n-input
        v-else
        v-model:value="keysText"
        type="textarea"
        placeholder="输入密钥，每行一个"
        :rows="8"
        style="margin-top: 20px"
      />

      <template #footer>
        <div style="display: flex; justify-content: flex-end; gap: 12px">
          <n-button @click="handleClose">取消</n-button>
          <n-button
            type="primary"
            @click="handleSubmit"
            :loading="loading"
            :disabled="useKeyInputFactory ? !singleKeyValue || !isValidKey : !keysText"
          >
            创建
          </n-button>
        </div>
      </template>
    </n-card>
  </n-modal>
</template>

<style scoped>
.form-modal {
  --n-color: rgba(255, 255, 255, 0.95);
}

:deep(.n-input) {
  --n-border-radius: 6px;
}

:deep(.n-card-header) {
  border-bottom: 1px solid rgba(239, 239, 245, 0.8);
  padding: 10px 20px;
}

:deep(.n-card__content) {
  max-height: calc(100vh - 68px - 61px - 50px);
  overflow-y: auto;
}

:deep(.n-card__footer) {
  border-top: 1px solid rgba(239, 239, 245, 0.8);
  padding: 10px 15px;
}
</style>
