<template>
  <NConfigProvider>
    <AppProvider>
      <RouterView />
    </AppProvider>
  </NConfigProvider>
</template>
<script lang="ts">
import { defineComponent } from 'vue';
import { AppProvider } from '@/components/Application';
export default defineComponent({
  name: 'App',
  components: { AppProvider },
  setup() {
  }
})
</script>