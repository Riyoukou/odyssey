<template>
  <div class="p-6 max-w-7xl mx-auto space-y-8">
    <!-- YAML 配置 -->
    <ElCard shadow="never">
      <h5 class="text-lg font-semibold">ID: {{props.activeID}}</h5>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  activeID: {
    type: Number,
  }
})
</script>