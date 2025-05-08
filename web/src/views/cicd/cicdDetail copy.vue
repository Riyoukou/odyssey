<template>
  <ElTabs v-model="activeTab">
    <!-- Page 1: Service Variables -->
    <ElTabPane label="变量配置" name="service">
      <div class="p-4 max-w-5xl mx-auto">
        <div class="mb-4">
          <ElSelect v-model="selectedServices" multiple placeholder="请选择服务" style="width: 400px;">
            <ElOption v-for="svc in allServices" :key="svc" :label="svc" :value="svc" />
          </ElSelect>
        </div>

        <div v-for="(vars, svc) in serviceVars" :key="svc" class="mb-4 border p-4 rounded">
          <div class="mb-2 flex justify-between items-center">
            <span class="font-bold">{{ svc }}</span>
            <ElButton v-if="svc === selectedServices[0]" type="primary" size="small" @click="syncServiceVars">同步到其他服务</ElButton>
          </div>
          <div v-for="(item, index) in vars" :key="index" class="flex gap-2 items-center mb-2">
            <ElInput v-model="item.key" placeholder="变量名" class="w-1/2" />
            <ElInput v-model="item.value" placeholder="变量值" class="w-1/2" />
            <ElButton type="danger" size="small" @click="removeVar(svc, index)">删除</ElButton>
          </div>
          <ElButton type="primary" size="small" @click="addVar(svc)">添加变量</ElButton>
        </div>
      </div>
    </ElTabPane>

    <!-- Page 2: Stage & Step Editor -->
    <ElTabPane label="阶段配置" name="stages">
      <div class="p-4 max-w-6xl mx-auto">
        <div class="mb-4 flex gap-2">
          <ElInput v-model="templateName" placeholder="模板名称" class="w-64" />
          <ElButton type="primary" @click="saveTemplate">保存为模板</ElButton>
          <ElSelect v-model="selectedTemplate" placeholder="导入模板" class="w-64">
            <ElOption v-for="tpl in templateList" :key="tpl" :label="tpl" :value="tpl" />
          </ElSelect>
          <ElButton @click="loadTemplate">导入</ElButton>
        </div>

        <div class="flex gap-4 overflow-auto">
          <div
            v-for="(stage, stageIndex) in form.stages"
            :key="stageIndex"
            class="min-w-[300px] border rounded p-4"
          >
            <div class="flex justify-between items-center mb-2">
              <h3 class="font-semibold text-base">阶段 {{ stageIndex + 1 }}</h3>
              <ElButton type="danger" size="small" @click="removeStage(stageIndex)">删除阶段</ElButton>
            </div>

            <div class="flex flex-col gap-3">
              <div v-for="(step, stepIndex) in stage.steps" :key="stepIndex" class="border rounded p-3 bg-gray-50">
                <ElInput v-model="step.name" placeholder="Step 名称" class="mb-2" />
                <ElSelect v-model="step.type" placeholder="选择类型" class="mb-2 w-full">
                  <ElOption v-for="item in stepTypes" :key="item.value" :label="item.label" :value="item.value" />
                </ElSelect>

                <div v-if="step.type === 'jenkins'">
                  <ElInput v-model="step.job_url" placeholder="Jenkins Job URL" class="mb-2" />
                  <ElSelect v-model="step.job_param" multiple placeholder="参数" class="mb-2 w-full">
                    <ElOption
                      v-for="param in mergedVars"
                      :key="param.key"
                      :label="param.key"
                      :value="param.key"
                    />
                  </ElSelect>
                </div>

                <div v-else-if="step.type === 'shell'">
                  <ElInput type="textarea" v-model="step.script" placeholder="Shell 脚本内容" class="mb-2" />
                </div>

                <div v-else-if="step.type === 'http'">
                  <ElInput v-model="step.endpoint" placeholder="请求地址" class="mb-2" />
                  <ElSelect v-model="step.method" class="mb-2">
                    <ElOption label="GET" value="GET" />
                    <ElOption label="POST" value="POST" />
                  </ElSelect>
                </div>

                <ElButton type="danger" size="small" @click="removeStep(stageIndex, stepIndex)">删除步骤</ElButton>
              </div>
              <ElButton type="primary" size="small" @click="addStep(stageIndex)">添加并行步骤</ElButton>
            </div>
          </div>
          <ElButton type="success" @click="addStage">添加阶段</ElButton>
        </div>
      </div>
    </ElTabPane>

    <!-- Page 3: JSON -->
    <ElTabPane label="JSON 预览" name="json">
      <div class="p-4">
        <ElButton type="primary" @click="exportJson">导出 JSON</ElButton>
        <pre class="bg-gray-100 mt-4 p-4 text-sm overflow-auto max-h-96">{{ JSON.stringify(exportData, null, 2) }}</pre>
      </div>
    </ElTabPane>
  </ElTabs>
</template>

<script setup lang="ts">
import { reactive, ref, computed } from 'vue';

const activeTab = ref('service');

const allServices = ['app-server', 'feishu-server', 'data-service'];
const selectedServices = ref<string[]>([]);
const serviceVars = reactive<Record<string, { key: string; value: string }[]>>({});
const branchOptions = ['main', 'develop', 'release', 'feature/demo'];

const addVar = (svc: string) => {
  serviceVars[svc].push({ key: '', value: '' });
};

const removeVar = (svc: string, index: number) => {
  serviceVars[svc].splice(index, 1);
};

const syncServiceVars = () => {
  const base = serviceVars[selectedServices.value[0]];
  for (const svc of selectedServices.value.slice(1)) {
    serviceVars[svc] = JSON.parse(JSON.stringify(base));
  }
};

watch(selectedServices, (newVal) => {
  for (const svc of newVal) {
    if (!serviceVars[svc]) {
      serviceVars[svc] = [{ key: 'GIT_BRANCH', value: 'main' }];
    }
  }
});

const form = reactive({
  stages: [
    {
      stage: 1,
      steps: [
        {
          name: '',
          type: 'jenkins',
          job_url: '',
          job_param: [],
          script: '',
          endpoint: '',
          method: 'GET',
          headers: {}
        }
      ]
    }
  ]
});

const stepTypes = [
  { label: 'Jenkins', value: 'jenkins' },
  { label: 'Shell 脚本', value: 'shell' },
  { label: 'HTTP 请求', value: 'http' }
];

const mergedVars = computed(() => {
  const keys = new Set<string>();
  selectedServices.value.forEach(svc => {
    (serviceVars[svc] || []).forEach(item => keys.add(item.key));
  });
  return Array.from(keys).map(k => ({ key: k }));
});

const addStage = () => {
  form.stages.push({ stage: form.stages.length + 1, steps: [] });
};

const removeStage = (index: number) => {
  form.stages.splice(index, 1);
};

const addStep = (stageIndex: number) => {
  form.stages[stageIndex].steps.push({
    name: '',
    type: 'jenkins',
    job_url: '',
    job_param: [],
    script: '',
    endpoint: '',
    method: 'GET',
    headers: {}
  });
};

const removeStep = (stageIndex: number, stepIndex: number) => {
  form.stages[stageIndex].steps.splice(stepIndex, 1);
};

// 模板逻辑
const templateName = ref('');
const selectedTemplate = ref('');
const templateList = ref<string[]>([]);

const refreshTemplateList = () => {
  const keys = Object.keys(localStorage).filter(k => k.startsWith('stage-template:'));
  templateList.value = keys.map(k => k.replace('stage-template:', ''));
};
refreshTemplateList();

const saveTemplate = () => {
  if (!templateName.value) return alert('请输入模板名称');
  localStorage.setItem(`stage-template:${templateName.value}`, JSON.stringify(form.stages));
  refreshTemplateList();
  alert('保存成功');
};

const loadTemplate = () => {
  if (!selectedTemplate.value) return;
  const tpl = localStorage.getItem(`stage-template:${selectedTemplate.value}`);
  if (tpl) {
    form.stages = JSON.parse(tpl);
    alert('模板导入成功');
  }
};

const exportData = ref({});
const exportJson = () => {
  const finalService = {};
  for (const svc of selectedServices.value) {
    const vars = serviceVars[svc] || [];
    finalService[svc] = {};
    for (const { key, value } of vars) {
      if (key) finalService[svc][key] = value;
    }
  }

  exportData.value = {
    service: finalService,
    stages: form.stages
  };
};
</script>

<style scoped>
pre {
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>