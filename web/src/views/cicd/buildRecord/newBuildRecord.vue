<template>
  <el-scrollbar>
    <div class="common-layout">
      <div class="odyssey-inner-page">
        <el-container>
          <el-header>
            <div style="margin: 80px 0" />
            <el-steps
              style="max-width: 900px; margin: 0 auto;"
              :space="500"
              :active="editForm.activeStep"
              align-center
              finish-status="success"
            >
              <el-step title="填写基本信息" />
              <el-step title="配置服务" />
              <el-step title="配置构建" />
            </el-steps>
          </el-header>
          <el-main>
            <div style="margin: 150px 0" />
            <el-main>
              <div v-if="editForm.activeStep == 0">
                <el-form
                  label-position="right"
                  label-width="auto"
                  :model="editForm.model"
                  style="max-width: 500px; margin: 0 auto;"
                >
                  <el-form-item label="项目名称">
                    <el-input v-model="editForm.model.project_name" placeholder="输入项目名称" disabled />
                  </el-form-item>
                  <div style="margin: 50px 0" />
                  <el-form-item
                    label="版本名称"
                    :rules="[ { required: true, message: '构建名称不能为空', trigger: 'blur' } ]"
                  >
                    <el-input v-model="editForm.model.name" placeholder="输入构建名称" />
                  </el-form-item>
                  <el-button @click="editForm.fastBuildRecordName">测试名称快速生成</el-button>
                  <div style="margin: 50px 0" />
                  <el-form-item
                    label="构建标签"
                    :rules="[ { required: true, message: '构建标签不能为空', trigger: 'blur' } ]"
                  >
                    <el-select v-model="editForm.model.tag" placeholder="选择标签">
                      <el-option label="测试发版" value="测试发版" />
                      <el-option label="上线" value="上线" />
                    </el-select>
                  </el-form-item>
                  <div style="margin: 50px 0" />
                  <el-form-item label="构建描述">
                    <el-input v-model="editForm.model.describe" placeholder="输入构建描述" type="textarea" />
                  </el-form-item>
                </el-form>
              </div>
              <div v-else-if="editForm.activeStep == 1">
                <el-form
                  label-position="right"
                  label-width="auto"
                  :model="editForm.model"
                  style="max-width: 900px; margin: 0 auto;"
                >
                  <el-form-item
                    label="构建环境"
                    :rules="[ { required: true, message: '构建环境不能为空', trigger: 'blur' } ]"
                  >
                    <el-select v-model="editForm.model.env" placeholder="选择环境" @focus="service.fetchEnv">
                      <el-option v-for="env in service.envData" :key="env.name" :label="env.label" :value="env.name" />
                    </el-select>
                  </el-form-item>
                  <el-form-item
                    label="构建服务"
                    :rules="[ { required: true, message: '构建服务不能为空', trigger: 'blur' } ]"
                  >
                    <el-select
                      v-model="editForm.selectedServices"
                      filterable
                      multiple
                      placeholder="选择服务"
                      @focus="service.fetchService"
                    >
                      <el-option
                        v-for="item in service.serviceData"
                        :key="item.name"
                        :label="item.name"
                        :value="item.name"
                      />
                    </el-select>
                  </el-form-item>
                </el-form>
              </div>
              <div v-else-if="editForm.activeStep == 2">
                <el-form
                  label-position="right"
                  label-width="auto"
                  :model="editForm.model"
                  style="max-width: 900px; margin: 0 auto;"
                >
                  <div class="mb-2 ml-4">
                    <el-radio-group v-model="radio1" @change="editForm.radioChange">
                      <el-radio value="branch" size="large">分支</el-radio>
                      <el-radio value="tag" size="large">标签</el-radio>
                    </el-radio-group>
                  </div>
                  <el-form-item
                    :rules="[ { required: true, message: '服务分支不能为空', trigger: 'blur' } ]"
                    v-for="(service, index) in editForm.model.services"
                    :key="index"
                    :label="service.service_name"
                  >
                    <el-select
                      v-model="editForm.model.services[index].branch"
                      filterable
                      placeholder="选择分支"
                    >
                      <el-option
                        v-for="tag in editForm.serviceTags[service.service_name] || []"
                        :key="tag"
                        :label="tag"
                        :value="tag"
                      />
                    </el-select>
                  </el-form-item>
                  <div style="margin-top: 20px;">
                    <el-button type="primary" @click="editForm.applyAllTags">应用到所有服务</el-button>
                  </div>
                </el-form>
              </div>
              <div
                v-else-if="editForm.activeStep == 3"
                style="display: flex; justify-content: center; align-items: center; text-align: center;"
              >
                <h2>提交成功</h2>
              </div>
            </el-main>
          </el-main>
          <el-footer>
            <div class="button-container" style="position: fixed; bottom: 100px; right: 100px;">
              <el-button v-if="editForm.activeStep < 3" :disabled="editForm.activeStep === 0" @click="editForm.prevStep">上一步</el-button>
              <el-button v-if="editForm.activeStep <= 1" type="primary" :disabled="editForm.activeStep === 3" @click="editForm.nextStep">下一步</el-button>
              <el-button v-else-if="editForm.activeStep === 2" type="primary" @click="editForm.nextStep">提交</el-button>
            </div>
          </el-footer>
        </el-container>
      </div>
    </div>
  </el-scrollbar>
</template>

<script lang="ts" setup>
import { ElMessage, type FormInstance } from 'element-plus'
import { ref, reactive } from 'vue'
import useUserStore from '@/stores/useUserStore'
import http from '@/api'
const userStore = useUserStore()
const radio1 = ref('branch');

const props = defineProps({
  activeProject: {
    type: String,
  }
})

// 表单配置
const editForm = reactive({
  ref: null as FormInstance | null,
  activeStep: 0,
  show: false,
  title: '',
  state: '',
  selectedServices: [] as string[],
  activeEnv: '',
  serviceTags: ref<{ [key: string]: string[] }>({}),
  model: {
    name: '',
    tag: '',
    env: '',
    project_name: props.activeProject,
    build_user: userStore.userInfo.name,
    describe: '',
    services: [] as { service_name: string; branch: string }[],
  },
  nextStep: () => {
    if (editForm.activeStep < 3) {
      switch (editForm.activeStep) {
        case 0:
          if (editForm.model.name) {
            editForm.activeStep++;
          } else {
            ElMessage.error("请正确填写表单");
          }
          break;
        case 1:
          if (editForm.model.env && editForm.selectedServices.length > 0) {
            editForm.activeStep++;
            editForm.goToSecondStep();
          } else {
            ElMessage.error("请正确填写表单");
          }
          break;
        case 2:
          if (editForm.model.services.length > 0 && editForm.model.services.every(service => service.branch !== '')) {
            editForm.activeStep++;
            service.createBuildRecord();
          } else {
            ElMessage.error("请正确填写表单");
          }
          break;
      }
    }    
  },
  prevStep: () => {
    if (editForm.activeStep > 0) {
      editForm.activeStep--;
    }
  },
  goToSecondStep: () => {
    if (editForm.selectedServices.length > 0) {
      editForm.selectedServices.forEach((selectedService) => {
        service.fetchGitLabBranch(radio1.value, selectedService);
      });
      editForm.model.services = editForm.selectedServices.map((name) => ({
        service_name: name,
        branch: "",
      }));
    }    
  },
  applyAllTags: () => {
    if (editForm.model.services.length > 0) {
      const firstTag = editForm.model.services[0].branch;
      editForm.model.services.forEach(service => {
        service.branch = firstTag;
      });
    }
  },
  radioChange: (value: string) => {
      radio1.value = value;
      editForm.selectedServices.forEach((selectedService) => {
      service.fetchGitLabBranch(radio1.value, selectedService);
    });    
  },
  fastBuildRecordName: () => {
    const now = new Date();
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, "0");
    const day = String(now.getDate()).padStart(2, "0");
    const hours = String(now.getHours()).padStart(2, "0");
    const minutes = String(now.getMinutes()).padStart(2, "0");
    editForm.model.name = `${props.activeProject}-dev-${year}${month}${day}${hours}${minutes}`;
  },
  clear: () => {
    editForm.model.name = '',
    editForm.model.tag ='',
    editForm.model.env = '',
    editForm.model.project_name = props.activeProject,
    editForm.model.build_user = userStore.userInfo.name,
    editForm.model.describe = '',
    editForm.model.services = [] as { service_name: string; branch: string }[],
    editForm.activeStep = 0; 
  }
})

const service = reactive({
  loading: false,
  data: [] as any[],
  serviceData: [] as any[],
  branchData: [] as any[],
  tagData: [] as any[],
  envData: [] as any[],
  fetchData: () => {
    service.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/build_record?project=${props.activeProject}`).then((res: any) => {
      service.data = res.data
      service.loading = false
    })
  },
  fetchEnv: () => {
    service.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/env?project=${props.activeProject}`).then((res: any) => {
      service.envData = res.data
      service.loading = false
    })
  },
  createBuildRecord: () => {
    service.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/create/build_record`, editForm.model).then((res: any) => {
      service.loading = false
      service.data = res.data
      editForm.clear();
      window.location.reload();
      ElMessage.success(res.message)
    })
  },
  fetchService: () => {
    service.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/service?project=${props.activeProject}&env=${editForm.model.env}`).then((res: any) => {
      service.serviceData = res.data
      service.loading = false
    })
  },
  fetchGitLabBranch: (type:string ,service_name: string) => {
    service.loading = true
    if (type === 'branch') {
      http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/gitlab_branch?project=${props.activeProject}&code_library=vehicle-platform`).then((res: any) => {
        service.branchData = res.data
        console.log(service.branchData)
        editForm.serviceTags[service_name] = service.branchData || [];
        service.loading = false
      })
    } else if (type === 'tag') {
      http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/gitlab_tag?project=${props.activeProject}&code_library=vehicle-platform`).then((res: any) => {
        service.tagData = res.data
        editForm.serviceTags[service_name] = service.tagData.sort().reverse() || [];
        if (editForm.serviceTags[service_name].length === 0) {
          editForm.serviceTags[service_name].push("master");
        }
        service.loading = false
      })
    }
  },
})
</script>

<style scoped>
.scrollbar-demo-item {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 50px;
  margin: 10px;
  text-align: center;
  border-radius: 4px;
  background: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
}
</style>
