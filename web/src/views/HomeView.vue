<template>
    <n-calendar
        v-model:value="value"
        #="{ year, month, date }"
        :is-date-disabled="isDateDisabled"
        @update:value="handleUpdateValue"
    >
      {{ year }}-{{ month }}-{{ date }}
    </n-calendar>
</template>

<script lang="ts">
import {defineComponent, ref} from 'vue'
import {useMessage, NCalendar} from 'naive-ui'
import {isYesterday, addDays} from 'date-fns/esm'

export default defineComponent({
  components: {
    NCalendar,
  },
  setup() {
    // const message = useMessage();
    // (window as any).$message = message;
    return {
      value: ref(addDays(Date.now(), 1).valueOf()),
      handleUpdateValue(
          _: number,
          {year, month, date}: { year: number; month: number; date: number }
      ) {
        // message.success(`${year}-${month}-${date}`)
      },
      isDateDisabled(timestamp: number) {
        return isYesterday(timestamp);
      }
    }
  }
})
</script>
