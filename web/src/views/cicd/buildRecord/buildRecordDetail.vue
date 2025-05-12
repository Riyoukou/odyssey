<template>
  <div class="odyssey-inner-page">
    <el-table
      :data="buildServiceRecords"
      row-key="id"
      table-layout="auto"
      max-height="900"
    >
      <el-table-column prop="service_name" label="服务名称" />
      <el-table-column prop="branch" label="构建分支" />
      <el-table-column prop="image" label="构建镜像" />
      <el-table-column prop="build_url" label="构建地址" />、
      <el-table-column prop="status" label="状态" />
    </el-table>
  </div>
</template>
<script lang="ts" setup>
//import { ElMessage } from 'element-plus';
import { ref } from 'vue';
//import { onMounted, defineProps } from 'vue'

const props = defineProps<{
  buildRecordName: string
}>()

const buildServiceRecords = ref<BuildServiceRecord[]>([]);
interface BuildServiceRecord {
  [key: string]: any;
}

const isLoading = ref(false);


</script>