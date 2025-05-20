<template>
  <div>
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
        label="发布名称"
        :rules="[ { required: true, message: '发布名称不能为空', trigger: 'blur' } ]"
      >
        <el-input v-model="editForm.model.name" placeholder="输入发布名称" />
      </el-form-item>
      <div style="margin: 50px 0" />
      <el-form-item label="发布版本">
        <el-input v-model="editForm.model.build_record_name" placeholder="输入发布版本" disabled />
      </el-form-item>
      <el-form-item
        label="发布集群"
        :rules="[ { required: true, message: '发布集群不能为空', trigger: 'blur' } ]"
      >
        <el-select v-model="selectClusters" placeholder="选择集群" @focus="service.fetchClusters()" multiple>
          <ElOption  v-for="item in service.clusterData" :key="item.name" :label="item.name" :value="item.name" />
        </el-select>
      </el-form-item>
      <div style="margin: 50px 0" />
      <el-form-item label="发布需求">
        <el-input v-model="editForm.model.description" placeholder="输入版本描述" type="textarea" />
      </el-form-item>
    </el-form>
    <el-footer>
      <div class="button-container" style="position: fixed; bottom: 100px; right: 100px;" @click="service.createDeployRecord()">
        <el-button type="primary">提交</el-button>
      </div>
    </el-footer>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import useUserStore from '@/stores/useUserStore'
import http from '@/api'

const userStore = useUserStore()
const selectClusters = ref<string[]>([])

const props = defineProps<{
  buildRecordRow: any
}>()

const editForm = reactive({
  ref: null as any,
  show: false,
  detailShow: false,

  title: '',
  state: '',
  model: {
    name: props.buildRecordRow.name,
    env: props.buildRecordRow.env,
    project_name: props.buildRecordRow.project_name,
    build_record_name: props.buildRecordRow.name,
    tag: props.buildRecordRow.tag,
    status: "Pending",
    deploy_user: userStore.userInfo.name,
    cluster_names: "",
    description: '',
  },
  clear: () => {
    editForm.model.name = '';
    editForm.model.env = '';
    editForm.model.project_name = '';
    editForm.model.build_record_name = '';
    editForm.model.tag = '';
    editForm.model.cluster_names = "";
    editForm.model.description = '';
  }
});


const service = reactive({
  loading: false,
  data: [] as any[],
  clusterData: [] as any[],
  createDeployRecord: () => {
    service.loading = true
    editForm.model.cluster_names= selectClusters.value.join(",")
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/create/deploy_record`, editForm.model).then((res: any) => {
      service.loading = false
      service.data = res.data
      editForm.clear();
      //window.location.reload();
      ElMessage.success(res.message)
    })
  },
  fetchClusters: () => {
    service.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/cluster`).then((res: any) => {
      service.clusterData = res.data
      service.loading = false
    })
  }
})
</script>