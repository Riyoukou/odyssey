<template>
  <div class="p-6 max-w-7xl mx-auto space-y-8">
    <!-- 构建 + 发布 并排 -->
    <div class="flex gap-6">
      <!-- 构建配置 -->
      <ElCard shadow="never" class="w-1/2">
        <template #header>
          <h2 class="text-lg font-semibold">构建配置</h2>
        </template>
        <ElForm label-width="100px">
          <h3 class="text-md font-semibold mb-2">构建任务变量</h3>

          <ElFormItem
            v-for="(item, index) in cicdForm.build.param"
            :key="index"
            :label="`变量 ${index + 1}`"
          >
            <div class="flex items-center gap-2 w-full">
              <ElInput v-model="item.key" placeholder="变量名" class="w-1/2" />
              <ElButton type="danger" size="small" @click="removeBuildVar(index)">删除</ElButton>
            </div>
          </ElFormItem>

          <ElFormItem>
            <ElButton type="primary" size="small" @click="addBuildVar">添加变量</ElButton>
          </ElFormItem>

          <ElFormItem label="步骤类型">
            <ElSelect v-model="cicdForm.build.type" placeholder="选择类型" class="w-full">
              <ElOption v-for="item in buildStepTypes" :key="item.value" :label="item.label" :value="item.value" />
            </ElSelect>
          </ElFormItem>

          <ElFormItem label="CICD 工具">
            <ElSelect v-model="cicdForm.build.cicd_tool" placeholder="选择 CICD 工具" class="w-full">
              <ElOption
                v-for="item in filteredCicdTools(cicdForm.build.type)"
                :key="item.name"
                :label="item.name"
                :value="item.type"
              />
            </ElSelect>
          </ElFormItem>

          <template v-if="cicdForm.build.type === 'jenkins'">
            <ElFormItem label="Job URL">
              <ElInput v-model="cicdForm.build.job_url" placeholder="Jenkins Job URL" class="w-full" />
            </ElFormItem>

            <ElFormItem label="参数选择">
              <ElSelect v-model="cicdForm.build.job_param" multiple placeholder="参数" class="w-full">
                <ElOption
                  v-for="param in cicdForm.build.param"
                  :key="param.key"
                  :label="param.key"
                  :value="param.key"
                />
              </ElSelect>
            </ElFormItem>
          </template>
        </ElForm>
      </ElCard>

      <!-- 发布配置 -->
      <ElCard shadow="never" class="w-1/2">
        <template #header>
          <h2 class="text-lg font-semibold">发布配置</h2>
        </template>
        <ElForm label-width="100px">
          <ElFormItem label="负载类型">
            <ElSelect v-model="cicdForm.release.workload" placeholder="负载类型" class="w-full">
              <ElOption
                v-for="item in releaseWorkloadTypes"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </ElSelect>
          </ElFormItem>

          <ElFormItem label="发布类型">
            <ElSelect v-model="cicdForm.release.deployType" placeholder="发布类型" class="w-full">
              <ElOption
                v-for="item in releaseDeployTypes"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </ElSelect>
          </ElFormItem>

          <ElFormItem label="步骤类型">
            <ElSelect v-model="cicdForm.release.type" placeholder="选择类型" class="w-full">
              <ElOption
                v-for="item in releaseStepTypes"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </ElSelect>
          </ElFormItem>

          <ElFormItem label="CICD 工具">
            <ElSelect v-model="cicdForm.release.cicd_tool" placeholder="选择 CICD 工具" class="w-full">
              <ElOption
                v-for="item in filteredCicdTools(cicdForm.release.type)"
                :key="item.name"
                :label="item.name"
                :value="item.type"
              />
            </ElSelect>
          </ElFormItem>

          <template v-if="cicdForm.release.type === 'argocd'">
            <ElFormItem label="Argocd 应用">
              <ElInput
                v-model="cicdForm.release.argocd_application"
                placeholder="Argocd Application"
                class="w-full"
              />
            </ElFormItem>
          </template>
        </ElForm>
      </ElCard>
    </div>

    <!-- JSON 预览 -->
    <ElCard shadow="never">
      <template #header>
        <h2 class="text-lg font-semibold">JSON 预览</h2>
      </template>
      <ElButton type="primary" @click="exportJson">导出 JSON</ElButton>
      <pre class="bg-gray-100 mt-4 p-4 text-sm overflow-auto max-h-96">
{{ JSON.stringify(exportData, null, 2) }}
      </pre>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const activeTab = ref('build');

const filteredCicdTools = (type: string) => {
  return cicdTools.filter(tool => tool.type === type);
};

const cicdTools = [
  { name: "ct-jenkins", type: "jenkins" },
  { name: "ct-argocd", type: "argocd" },
]

const buildStepTypes = [
  { label: 'Jenkins', value: 'jenkins' },
];

const releaseStepTypes = [
  { label: 'Argocd', value: 'argocd' },
];

const releaseWorkloadTypes = [
  { label: 'Deployment', value: 'deployment' },
  { label: 'StatefulSet', value:'statefulset' },
  { label: 'DaemonSet', value: 'daemonset' },
  { label: 'CloneSet', value: 'cloneset' }
];

const releaseDeployTypes = [
  { label: 'KruiseCanary', value: 'kruise_canary' },
  { label: 'KruiseBlueGreen', value: 'kruise_blue_green' },
];

const cicdForm = reactive({
  build: {
    param: [{ key: '' }],
    type: '',
    cicd_tool: '',
    job_url: '',
    job_param: [],
  },
  release: {
    deployType : "",
    workload: "",
    gitOPSURL: "",
    gitOPSPath: "",
    type: '',
    cicd_tool: '',
    argocd_application: '',
  }
});


const addBuildVar = () => cicdForm.build.param.push({ key: ''});
const removeBuildVar = (i: number) => cicdForm.build.param.splice(i, 1);

const exportData = ref({});
const exportJson = () => {
  exportData.value = {
    build: {
      param: cicdForm.build.param,
      type: cicdForm.build.type,
      cicd_tool: cicdForm.build.cicd_tool,
      job_url: cicdForm.build.job_url,
      job_param: cicdForm.build.job_param,
    },
    release: {
      deployType : cicdForm.release.deployType,
      workload: cicdForm.release.workload,
      type: cicdForm.release.type,
      cicd_tool: cicdForm.release.cicd_tool,
      argocd_application: cicdForm.release.argocd_application,
    }
  };
};
</script>

<style scoped>
pre {
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>