<template>
  <ElTabs v-model="activeTab">
    <!-- 构建 JSON 页 -->
    <ElTabPane label="构建配置" name="build">
      <div class="p-4 max-w-4xl mx-auto">
        <h3 class="text-lg font-semibold mb-2">构建任务变量</h3>
        <div v-for="(item, index) in buildForm.param" :key="index" class="flex gap-2 items-center mb-2">
          <ElInput v-model="item.key" placeholder="变量名" class="w-1/2" />
          <ElButton type="danger" size="small" @click="removeBuildVar(index)">删除</ElButton>
        </div>
        <ElButton type="primary" size="small" @click="addBuildVar">添加变量</ElButton>

        <div v-for="(stage, stageIndex) in buildForm.stages" :key="stageIndex" class="border rounded p-4 my-4">
          <div class="flex justify-between items-center mb-2">
            <h3 class="font-semibold text-base">阶段 {{ stageIndex + 1 }}</h3>
            <el-button type="danger" size="small" @click="removeBuildStage(stageIndex)">删除阶段</el-button>
          </div>
          <div v-for="(step, stepIndex) in stage.steps" :key="stepIndex" class="border rounded p-3 bg-gray-50 mb-2">
            <ElInput v-model="step.name" placeholder="Step 名称" class="mb-2" />
            <ElSelect v-model="step.type" placeholder="选择类型" class="mb-2 w-full">
              <ElOption v-for="item in buildStepTypes" :key="item.value" :label="item.label" :value="item.value" />
            </ElSelect>
            <ElSelect v-model="step.cicd_tool" placeholder="选择cicd工具" class="mb-2 w-full">
              <ElOption
                v-for="item in filteredCicdTools(step.type)"
                :key="item.name"
                :label="item.name"
                :value="item.type"
              />
            </ElSelect>
            <div v-if="step.type === 'jenkins'">
              <ElInput v-model="step.job_url" placeholder="Jenkins Job URL" class="mb-2" />
              <ElSelect v-model="step.job_param" multiple placeholder="参数" class="mb-2 w-full">
                <ElOption v-for="param in buildForm.param" :key="param.key" :label="param.key" :value="param.key" />
              </ElSelect>
            </div>
            <ElButton type="danger" size="small" @click="removeBuildStep(stageIndex, stepIndex)">删除步骤</ElButton>
          </div>
          <ElButton type="primary" @click="addBuildStep(stageIndex)">添加步骤</ElButton>
        </div>
        <ElButton type="success" @click="addBuildStage">添加阶段</ElButton>
      </div>
    </ElTabPane>

    <!-- 发布 JSON 页 -->
    <ElTabPane label="发布配置" name="release">
      <div class="p-4 max-w-4xl mx-auto">
        <h3 class="text-lg font-semibold mb-2">发布任务变量</h3>
        <!--<div v-for="(item, index) in releaseForm.param" :key="index" class="flex gap-2 items-center mb-2">
          <ElInput v-model="item.key" placeholder="变量名" class="w-1/2" />
          <ElButton type="danger" size="small" @click="removeReleaseVar(index)">删除</ElButton>
        </div>-->
        <!--<ElButton type="primary" size="small" @click="addReleaseVar">添加变量</ElButton>-->
        <ElSelect v-model="releaseForm.workload" placeholder="负载类型" class="mb-2 w-full">
          <ElOption v-for="item in releaseWorkloadTypes" :key="item.value" :label="item.label" :value="item.value" />
        </ElSelect>
        <ElSelect v-model="releaseForm.deployType" placeholder="发布类型" class="mb-2 w-full">
          <ElOption v-for="item in releaseDeployTypes" :key="item.value" :label="item.label" :value="item.value" />
        </ElSelect>
        <div v-for="(stage, stageIndex) in releaseForm.stages" :key="stageIndex" class="border rounded p-4 my-4">
          <div class="flex justify-between items-center mb-2">
            <h3 class="font-semibold text-base">阶段 {{ stageIndex + 1 }}</h3>
            <el-button type="danger" size="small" @click="removeReleaseStage(stageIndex)">删除阶段</el-button>
          </div>
          <div v-for="(step, stepIndex) in stage.steps" :key="stepIndex" class="border rounded p-3 bg-gray-50 mb-2">
            <ElInput v-model="step.name" placeholder="Step 名称" class="mb-2" />
            <ElSelect v-model="step.type" placeholder="选择类型" class="mb-2 w-full">
              <ElOption v-for="item in releaseStepTypes" :key="item.value" :label="item.label" :value="item.value" />
            </ElSelect>
            <ElSelect v-model="step.cicd_tool" placeholder="选择cicd工具" class="mb-2 w-full">
              <ElOption
                v-for="item in filteredCicdTools(step.type)"
                :key="item.name"
                :label="item.name"
                :value="item.type"
              />
            </ElSelect>
            <div v-if="step.type === 'argocd'">
              <ElInput v-model="step.argocd_application" placeholder="Argocd Application" class="mb-2" />
            </div>
            <ElButton type="danger" size="small" @click="removeReleaseStep(stageIndex, stepIndex)">删除步骤</ElButton>
          </div>
          <ElButton type="primary" @click="addReleaseStep(stageIndex)">添加步骤</ElButton>
        </div>
        <ElButton type="success" @click="addReleaseStage">添加阶段</ElButton>
      </div>
    </ElTabPane>

    <!-- JSON 预览 -->
    <ElTabPane label="JSON 预览" name="json">
      <div class="p-4">
        <ElButton type="primary" @click="exportJson">导出 JSON</ElButton>
        <pre class="bg-gray-100 mt-4 p-4 text-sm overflow-auto max-h-200">{{ JSON.stringify(exportData, null, 2) }}</pre>
      </div>
    </ElTabPane>
  </ElTabs>
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


const buildForm = reactive({
  param: [{ key: '' }],
  stages: [
    {
      stage: 1,
      steps: [
        {
          name: '',
          type: '',
          cicd_tool: '',
          job_url: '',
          job_param: [],
        }
      ]
    }
  ]
});

const releaseForm = reactive({
//  param: [{ key: '' }],
  workload: "",
  deployType : "",
  stages: [
    {
      stage: 1,
      steps: [
        {
          name: '',
          type: '',
          cicd_tool: '',
          argocd_application: '',
        }
      ]
    }
  ]
});

const addBuildVar = () => buildForm.param.push({ key: ''});
const removeBuildVar = (i: number) => buildForm.param.splice(i, 1);
const addBuildStage = () => buildForm.stages.push({ stage: buildForm.stages.length + 1, steps: [] });
const addBuildStep = (stageIdx: number) =>
  buildForm.stages[stageIdx].steps.push({ name: '', type: '', cicd_tool: '', job_url: '', job_param: []});
const removeBuildStep = (stageIdx: number, stepIdx: number) => buildForm.stages[stageIdx].steps.splice(stepIdx, 1);
const removeBuildStage = (stageIdx: number) => buildForm.stages.splice(stageIdx, 1);

//const addReleaseVar = () => releaseForm.param.push({ key: '' });
//const removeReleaseVar = (i: number) => releaseForm.param.splice(i, 1);
const addReleaseStage = () => releaseForm.stages.push({ stage: releaseForm.stages.length + 1, steps: [] });
const addReleaseStep = (stageIdx: number) =>
  releaseForm.stages[stageIdx].steps.push({ name: '', type: '', cicd_tool: '', argocd_application: ''});
const removeReleaseStep = (stageIdx: number, stepIdx: number) => releaseForm.stages[stageIdx].steps.splice(stepIdx, 1);
const removeReleaseStage = (stageIdx: number) => releaseForm.stages.splice(stageIdx, 1);
const exportData = ref({});
const exportJson = () => {
  exportData.value = {
    build: {
      stages: buildForm.stages
    },
    release: {
      deployType : releaseForm.deployType,
      workload: releaseForm.workload,
      stages: releaseForm.stages
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