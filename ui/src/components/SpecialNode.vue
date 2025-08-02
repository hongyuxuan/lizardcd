<!-- CustomTaskNode.vue -->
<template>
  <div 
    class="task-node"
    :style="{ 
      borderColor: statusColor,
      backgroundColor: statusBackground
    }">
    <Handle v-if="!data.end"
      id="source-right" 
      type="source" 
      class="custom-handle"
      :position="Position.Right" 
      :connectable="true" />
    <Handle v-if="!data.start"
      id="target-left" 
      type="target" 
      class="custom-handle"
      :position="Position.Left" />
    <div class="task-name">{{ data.label }} <el-icon v-if="['Running','Pending','pending'].includes(data.status)" :class="`is-loading text-${getColor(data.status)}`"><Refresh /></el-icon></div>
    <div class="task-status">{{ data.status }}</div>
  </div>
</template>

<script setup>
import { defineProps, computed } from 'vue';
import { Handle, Position } from '@vue-flow/core';
import { getColor, getPodClass } from '@/assets/util/common';
const props = defineProps(['data']);

const statusMap = {
  Pending: { color: 'var(--el-color-warning)', bg: 'var(--el-color-warning-light-9)' },
  pending: { color: 'var(--el-color-warning)', bg: 'var(--el-color-warning-light-9)' },
  Started: { color: 'var(--el-color-warning)', bg: 'var(--el-color-warning-light-9)' },
  Running: { color: '#2196F3', bg: '#e3f2fd' },
  Succeeded: { color: '#4CAF50', bg: '#e8f5e9' },
  approved: { color: '#4CAF50', bg: '#e8f5e9' },
  rejected: { color: '#F44336', bg: '#ffebee' },
  Failed: { color: '#F44336', bg: '#ffebee' },
  RunTimedOut: { color: '#F44336', bg: '#ffebee' },
  TaskRunTimeout: { color: '#F44336', bg: '#ffebee' },
  TaskRunCancelled: { color: '#F44336', bg: '#ffebee' },
};

const statusColor = computed(() => statusMap[props.data.status]?.color || '#ccc');
const statusBackground = computed(() => statusMap[props.data.status]?.bg || '#f5f5f5');

</script>

<style scoped>
.task-node {
  padding: 10px;
  border-radius: 5px;
  border: 2px solid;
  min-width: 150px;
  height: 80px;
  text-align: center;
}
.task-name {
  font-size: 15px;
  font-weight: bold;
}
.task-status {
  font-size: 13px;
  margin-top: 20px;
  color: var(--el-color-info)
}
.custom-handle {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}
</style>