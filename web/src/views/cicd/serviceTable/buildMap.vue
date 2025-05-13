<template>
  <div class="p-6 max-w-7xl mx-auto space-y-8">
    <div class="flex gap-6">
      <!-- 构建配置 -->
      <ElCard shadow="never" class="w-1/2">
        <template #header>
          <h2 class="text-lg font-semibold">构建配置</h2>
        </template>
        <ElForm label-width="100px">
          <h3 class="text-md font-semibold mb-2">构建任务变量</h3>

          <ElFormItem
            v-for="(item, index) in editForm.model.param"
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

          <ElFormItem label="CI类型">
            <ElSelect v-model="editForm.model.type" placeholder="选择类型" class="w-full">
              <ElOption label="Jenkins" value="jenkins" />
            </ElSelect>
          </ElFormItem>

          <ElFormItem label="CICD 工具" v-if="editForm.model.type === 'jenkins'" >
            <ElSelect placeholder="请选择工具" v-model="editForm.model.cicd_tool" @focus="service.requestCICDToolsData()">
              <ElOption  
                v-for="item in filteredCicdTools(editForm.model.type)" 
                :key="item.name" 
                :label="item.name" 
                :value="item.name" 
              />
            </ElSelect>
          </ElFormItem>

          <template v-if="editForm.model.type === 'jenkins'">
            <ElFormItem label="Job URL">
              <ElInput v-model="editForm.model.job_url" placeholder="Jenkins Job URL" class="w-full" />
            </ElFormItem>

            <ElFormItem label="参数选择">
              <ElSelect v-model="editForm.model.job_param" multiple placeholder="参数" class="w-full">
                <ElOption
                  v-for="param in editForm.model.param"
                  :key="param.key"
                  :label="param.key"
                  :value="param.key"
                />
              </ElSelect>
            </ElFormItem>
          </template>
        </ElForm>
      </ElCard>
    </div>
    <ElButton type="primary" @click="submitData">保存配置</ElButton>
  </div>
</template>

<script setup lang="ts">
import { reactive, PropType } from 'vue';
import http from '@/api'

const props = defineProps({
  activeID: {
    type: Number,
  },
  activeCluster: {
    type: String,
  },
  activeBuildMap: {
    type: Array as PropType<any>,
  }
})

const editForm = reactive({
  model: props.activeBuildMap || {
    param: [{ key: '' }],
    type: '',
    cicd_tool: '',
    job_url: '',
    job_param: [],
  }
});

const addBuildVar = () => {
  if (Array.isArray(editForm.model.param)) {
    editForm.model.param.push({ key: '' });
  } else {
    editForm.model.param = [{ key: '' }]; // Initialize if not an array
  }
};
const removeBuildVar = (i: number) => editForm.model.param.splice(i, 1);

const submitData = () => {
  console.log(editForm.model)
  service.submitBuildMap()
}

const service = reactive({
  loading: false,
  cicdToolsData: [] as any[],
  codeLibraryData: [] as any[],
  requestCICDToolsData: () => {
    service.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/cicd_tool`).then((res: any) => {
      service.cicdToolsData = res.data
      service.loading = false
    })
  },
  requestCodeSource: () => {
    service.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/code_library`).then((res: any) => {
      service.codeLibraryData = res.data
      service.codeLibraryData = service.codeLibraryData.filter( tool => tool.type === 'gitops')
      service.loading = false
    })
  },
  submitBuildMap: () => {
    service.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/update/service_build_map?id=${props.activeID}`, editForm.model).then((res: any) => {
      service.loading = false
    })
  }
});


const filteredCicdTools = (type: string) => {
  return service.cicdToolsData.filter(tool => tool.type === type);
};

</script>

<style scoped>
pre {
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>