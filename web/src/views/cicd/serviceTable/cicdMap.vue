<template>
  <div class="p-6 max-w-7xl mx-auto space-y-8">
    <!-- YAML 配置 -->
    <ElCard shadow="never">
      <h5 class="text-lg font-semibold">配置集群: {{props.activeCluster}}</h5>
    </ElCard>
    <!-- YAML 配置 -->
    <ElCard shadow="never">
      <template #header>
        <h2 class="text-lg font-semibold">YAML 配置</h2>
      </template>
      <ElForm label-width="100px">
        <ElFormItem label="是否 GitOps">
          <ElSwitch v-model="editForm.model.yaml.isGitOps" />
        </ElFormItem>

        <template v-if="editForm.model.yaml.isGitOps">
          <ElFormItem label="GitOps 库">
            <ElSelect 
              v-model="editForm.model.yaml.gitopsrepo" 
              placeholder="请选择 GitOps 库" 
              class="w-full" 
              @focus="service.requestCodeSource()"
            >
              <ElOption  
                v-for="item in service.codeLibraryData" 
                :key="item.name" 
                :label="item.name" 
                :value="item.name" 
              />
            </ElSelect>
          </ElFormItem>

          <ElFormItem label="GitOps 类型">
            <ElSelect v-model="editForm.model.yaml.gitopsType" placeholder="选择 GitOps 类型" class="w-full">
              <ElOption label="Kustomize" value="kustomize" />
            </ElSelect>
          </ElFormItem>

          <ElFormItem label="GitOps 路径">
            <ElInput v-model="editForm.model.yaml.filePath" placeholder="请输入文件路径" class="w-full" />
          </ElFormItem>
        </template>

        <template v-else>
          <ElFormItem
            label="YAML 内容"
            :error="yamlError"
          >
            <ElInput
              v-model="editForm.model.yaml.content"
              type="textarea"
              :rows="6"
              placeholder="请输入合法的 YAML 内容"
              class="w-full"
              @blur="validateYaml"
            />
          </ElFormItem>
        </template>
      </ElForm>
    </ElCard>

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
            v-for="(item, index) in editForm.model.build.param"
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
            <ElSelect v-model="editForm.model.build.type" placeholder="选择类型" class="w-full">
              <ElOption label="Jenkins" value="jenkins" />
            </ElSelect>
          </ElFormItem>

          <ElFormItem label="CICD 工具" v-if="editForm.model.build.type === 'jenkins'" >
            <ElSelect placeholder="请选择工具" v-model="editForm.model.build.cicd_tool" @focus="service.requestCICDToolsData()">
              <ElOption  
                v-for="item in filteredCicdTools(editForm.model.build.type)" 
                :key="item.name" 
                :label="item.name" 
                :value="item.name" 
              />
            </ElSelect>
          </ElFormItem>

          <template v-if="editForm.model.build.type === 'jenkins'">
            <ElFormItem label="Job URL">
              <ElInput v-model="editForm.model.build.job_url" placeholder="Jenkins Job URL" class="w-full" />
            </ElFormItem>

            <ElFormItem label="参数选择">
              <ElSelect v-model="editForm.model.build.job_param" multiple placeholder="参数" class="w-full">
                <ElOption
                  v-for="param in editForm.model.build.param"
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
            <ElSelect v-model="editForm.model.release.workload" placeholder="负载类型" class="w-full">
              <ElOption label="Deployment" value="deployment" />
              <ElOption label="StatefulSet" value="statefulset" />
              <ElOption label="DaemonSet" value="daemonset" />
              <ElOption label="CloneSet" value="cloneset" />
            </ElSelect>
          </ElFormItem>

          <ElFormItem label="发布策略类型">
            <ElSelect v-model="editForm.model.release.deployType" placeholder="发布类型" class="w-full">
              <ElOption label="KruiseCanary" value="kruise_canary" />
              <ElOption label="KruiseBlueGreen" value="kruise_blue_green" />
            </ElSelect>
          </ElFormItem>

          <ElFormItem label="CD类型">
            <ElSelect v-model="editForm.model.release.type" placeholder="选择类型" class="w-full">
              <ElOption label="Argocd" value="argocd" />
            </ElSelect>
          </ElFormItem>

          <ElFormItem label="CICD 工具" v-if="editForm.model.release.type === 'argocd'" >
            <ElSelect placeholder="请选择工具" v-model="editForm.model.release.cicd_tool" @focus="service.requestCICDToolsData()">
              <ElOption  
                v-for="item in filteredCicdTools(editForm.model.release.type)" 
                :key="item.name" 
                :label="item.name" 
                :value="item.name" 
              />
            </ElSelect>
          </ElFormItem>

          <template v-if="editForm.model.release.type === 'argocd'">
            <ElFormItem label="Argocd 应用">
              <ElInput
                v-model="editForm.model.release.argocd_application"
                placeholder="Argocd Application"
                class="w-full"
              />
            </ElFormItem>
          </template>
        </ElForm>
      </ElCard>
    </div>
    <ElButton type="primary" @click="submitData">保存配置</ElButton>
    <!-- JSON 预览 -->
    <!--<ElCard shadow="never">
      <template #header>
        <h2 class="text-lg font-semibold">JSON 预览</h2>
      </template>
      <ElButton type="primary" @click="exportJson">导出 JSON</ElButton>
      <pre class="bg-gray-100 mt-4 p-4 text-sm overflow-auto max-h-96">
{{ JSON.stringify(props.activeCICDMap, null, 2) }}
      </pre>
    </ElCard>-->
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, PropType } from 'vue';
import yaml from 'yaml';
import http from '@/api'

const props = defineProps({
  activeID: {
    type: Number,
  },
  activeCluster: {
    type: String,
  },
  activeCICDMap: {
    type: Array as PropType<any[]>,
  }
})

const editForm = reactive({
  model: props.activeCICDMap?.[0] ||{
    yaml: {
      isGitOps: false,
      gitopsrepo: '',
      gitopsType: '',
      filePath: '',
      content: '',
    },
    build: {
      param: [{ key: '' }],
      type: '',
      cicd_tool: '',
      job_url: '',
      job_param: [],
    },
    release: {
      deployType: '',
      workload: '',
      type: '',
      cicd_tool: '',
      argocd_application: '',
    }
  }
});

const addBuildVar = () => {
  // Ensure 'param' is always an array
  if (Array.isArray(editForm.model.build.param)) {
    editForm.model.build.param.push({ key: '' });
  } else {
    editForm.model.build.param = [{ key: '' }]; // Initialize if not an array
  }
};
const removeBuildVar = (i: number) => editForm.model.build.param.splice(i, 1);

const submitData = () => {
  console.log(editForm.model)
  service.submitCICDMap()
}

const yamlError = ref('');

const validateYaml = () => {
  if (!editForm.model.yaml.isGitOps && editForm.model.yaml.content.trim()) {
    try {
      yaml.parse(editForm.model.yaml.content);
      yamlError.value = '';
    } catch (e) {
      yamlError.value = 'YAML 格式错误';
    }
  } else {
    yamlError.value = '';
  }
};

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
  submitCICDMap: () => {
    service.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/service/cicd_map/${props.activeID}?action=update&clusters=${props.activeCluster}`, editForm.model).then((res: any) => {
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