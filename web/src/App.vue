<template>
  <NConfigProvider
      :locale="zhCN"
      :date-locale="dateZhCN"
  >
    <AppProvider>
      <RouterView />
    </AppProvider>
  </NConfigProvider>
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted } from 'vue';
import { zhCN, dateZhCN, NConfigProvider } from 'naive-ui';
import { AppProvider } from '@/components/Application';

onMounted(() => {
  // null
})
onUnmounted(() => {
  // null
});
</script>
