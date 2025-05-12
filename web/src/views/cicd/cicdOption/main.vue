<template>
  <div class="odyssey-page">
    <div class="odyssey-inner-page">
     <el-container>
      <el-header>
        <el-tabs v-model="activeIndex" @tab-click="handleTabClick">
          <el-tab-pane label="集群" name="cluster" />
            <el-tab-pane label="项目" name="project"/>
            <el-tab-pane label="环境" name="env"/>
            <el-tab-pane label="凭证" name="credential"/>
            <el-tab-pane label="代码库" name="code_library"/>
            <el-tab-pane label="CICD工具" name="cicd_tool"/>
          </el-tabs>
        </el-header>
        <el-main>
          <router-view />
          <clusterTable v-if="activeIndex === 'cluster'" />
          <projectTable v-if="activeIndex === 'project'" />
          <envTable v-if="activeIndex === 'env'" />
          <credentialTable v-if="activeIndex === 'credential'" />
          <codeLibraryTable v-if="activeIndex === 'code_library'" />
          <cicdToolTable v-if="activeIndex === 'cicd_tool'" />
        </el-main>
      </el-container>
    </div> 
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import type { TabsPaneContext } from 'element-plus'
import { useRouter } from 'vue-router';
import { useRoute } from 'vue-router';
import clusterTable from '@/views/cicd/cicdOption/clusterTable/clusterTable.vue'
import projectTable from '@/views/cicd/cicdOption/projectTable/projectTable.vue'
import envTable from '@/views/cicd/cicdOption/envTable/envTable.vue'
import credentialTable from '@/views/cicd/cicdOption/credentialTable/credentialTable.vue'
import codeLibraryTable from '@/views/cicd/cicdOption/codeLibraryTable/codeLibraryTable.vue'
import cicdToolTable from '@/views/cicd/cicdOption/cicdToolTable/cicdToolTable.vue'
const router = useRouter();
const route = useRoute();
const activeIndex = ref(route.query.Index || "cluster");

const handleTabClick = (tab: TabsPaneContext) => {
  activeIndex.value = tab.props.name?.toString() || "cluster";
  // 合并现有的 query 参数，保留之前的参数
  router.push({
    query: { ...route.query, Index: activeIndex.value },  
  });
};

const navigateToDetails = async (project: string) => {
  // 将项目存储到 sessionStorage
  router.push({
    query: { project: project },  // 设置查询参数
  }).then(() => {
    // 刷新当前页面
    window.location.reload();  // 刷新页面并保留查询参数
  });
};


</script>