<template>
  <div class="flex flex-col overflow-auto">
    <!-- 搜索表单 -->
    <ElForm v-if="search.show" :model="search.model" ref="searchRef" label-width="80px">
      <ElFormItem label="名称" prop="name">
        <ElInput v-model="search.model.name" placeholder="请输入名称" />
      </ElFormItem>
      <div class="flex justify-center mb-3">
        <ElButton icon="RefreshRight" @click="search.reset">重置</ElButton>
        <ElButton type="primary" icon="Search" @click="search.search">查询</ElButton>
      </div>
    </ElForm>

    <!-- 按钮条 -->
    <div class="table-bar flex justify-between items-center mb-3">
      <div>
        <ElButton icon="Plus" @click="editForm.toAdd">新增</ElButton>
        <el-popconfirm title="Are you sure to delete this?" @confirm="deleteSelected">
          <template #reference>
            <ElButton icon="Delete">删除</ElButton>
          </template>
        </el-popconfirm>
      </div>
      <div >
        <ElSelect class="w-50" @change="table.selectProject" :placeholder=table.activeProject >
          <ElOption v-for="item in table.projectData" :key="item.name" :label="item.name" :value="item.name" />
        </ElSelect>
        <ElButton icon="Refresh" round @click="table.request" class="ml-4"></ElButton>
        <ElButton icon="Search" round @click="search.show = !search.show"></ElButton>
      </div>
    </div>

    <!-- 列表 -->
    <el-table 
      :data="table.filteredData" 
      style="width: 100%" 
      @selection-change="handleSelectionChange"
      :default-sort="{ prop: 'created_at', order: 'descending' }"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="名称" sortable/>
      <el-table-column prop="status" label="状态" sortable/>
      <el-table-column prop="project_name" label="项目" sortable/>
      <el-table-column prop="env" label="环境" sortable/>
      <el-table-column prop="tag" label="标签" sortable/>
      <el-table-column prop="cluster_names" label="发布集群" sortable/>
      <el-table-column prop="deploy_user" label="提交用户" sortable/>
      <el-table-column prop="created_at" label="创建时间" sortable/> 
      <el-table-column label="操作" fixed="right" width="200">
        <template #default="{ row }">
          <el-button link type="primary" @click="editForm.toDetail(row)">详情</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-drawer
      v-model="editForm.show" 
      direction="rtl"
    >
      <template #header>
        <h4>创建构建</h4>
      </template>
      <template #default>
        <newBuildRecordIndex
          :activeProject = table.activeProject
        />
      </template>
      <template #footer>
        <div style="flex: auto">
          <el-button @click="editForm.show = false">关闭</el-button>
        </div>
      </template>
    </el-drawer>
    <el-drawer
      v-model="editForm.detailShow" 
      direction="rtl"
      size="60%"
    >
      <template #header>
        <h4>发布详情</h4>
      </template>
      <template #default>
        <el-button type="primary" round @click="table.fetchDeployServiceRecord(table.activeDeployRecord)" >刷新</el-button>
        <el-button type="primary" round @click="table.startDeploy(multipleDeployServiceRecordsSelection)">发布</el-button>
        <el-button type="primary" round @click="table.approveDeploy(multipleDeployServiceRecordsSelection)">批准</el-button>
        <el-table
          ref="multipleTableRef"
          row-key="id"
          :data="table.deployServiceRecordData"
          table-layout="auto"
          max-height="900"
          @selection-change="handleDeployServiceRecordsSelectionChange"
          :default-sort="{ prop: 'cluster_name', order: 'ascending' }"
        >
          <el-table-column type="selection" :selectable="selectable" width="55" />
          <el-table-column prop="service_name" label="服务名称" />
          <el-table-column prop="env" label="发布环境" />
          <el-table-column prop="image" label="发布镜像" />
          <el-table-column prop="status" label="状态" />
          <el-table-column prop="cluster_name" label="发布集群" />
        </el-table>
      </template>
      <template #footer>
        <div style="flex: auto">
          <el-button @click="editForm.show = false">关闭</el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage, type FormInstance, type TableInstance } from 'element-plus'
import http from '@/api'
import newBuildRecordIndex from '@/views/cicd/buildRecord/newBuildRecord.vue'

const selectable = (row: any) => !["Running","Failed","Deploying"].some(status => row.status.includes(status))

const multipleTableRef = ref<TableInstance>();
const multipleDeployServiceRecordsSelection = ref<any[]>([])

const handleDeployServiceRecordsSelectionChange = (val: any[]) => {
  multipleDeployServiceRecordsSelection.value = val
}

// 搜索表单配置
const searchRef = ref<FormInstance | null>(null);
const search = reactive({
  show: false,
  model: {
    name: '',
  },
  search: () => {
    table.filteredData = table.data.filter(
      (data: any) =>
        data.name?.toLowerCase().includes(search.model.name.toLowerCase())
    );
  },
  reset: () => {
    searchRef.value?.resetFields();
  },
})

const rules = computed(() => {
  const baseRules = {
    name: [{ required: true, message: '请输入环境名称', trigger: 'blur' }],
    type: [{ required: true, message: '请选择环境类型', trigger: 'change' }],
    project_name: [{ required: true, message: '请选择项目名称', trigger: 'change' }],
    namespace: [{ required: true, message: '请输入命名空间', trigger: 'blur' }],
  }
  return baseRules
})

// 表单配置
const editForm = reactive({
  ref: null as FormInstance | null,
  show: false,
  detailShow: false,
  title: '',
  state: '',
  model: {
    name: '',
    project_name: '',
    env: '',
    tag: '',
    description: '',
    deploy_user: '',
  },
  toAdd: () => {
    editForm.ref?.resetFields()
    editForm.show = true
    editForm.title = '新增环境'
    editForm.state = 'add'
    editForm.model = {
      name: '',
      project_name: '',
      env: '',
      tag: '',
      description: '',
      deploy_user: '',
    }
    table.fetchProject()
  },
  toDetail: (row: any) => {
    editForm.detailShow = true
    table.activeDeployRecord = row
    table.fetchDeployServiceRecord(table.activeDeployRecord)
  },
})
// 表格配置
const table = reactive({
  loading: false,
  border: true,
  activeProject: '请选择项目名称',
  activeDeployRecord: null as any,
  data: [] as any[],
  deployServiceRecordData: [] as any[],
  filteredData: [] as any[],
  projectData: [] as any[],
  request: () => {
    if (table.activeProject !== '请选择项目名称'){
      table.loading = true
        http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/deploy_record?project=${table.activeProject}`).then((res: any) => {
          table.data = res.data
          table.filteredData = res.data
          table.loading = false
      })
    }
  },
  fetchDeployServiceRecord: (form: any) => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/deploy_service_record?deploy_record=${form.name}`).then((res: any) => {
      table.deployServiceRecordData = res.data
      table.loading = false
    })
  },
  fetchProject: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/project`).then((res: any) => {
      table.projectData = res.data
      table.loading = false
    })
  },
  delete: (form: any) => {
    table.loading = true
    http.delete(import.meta.env.VITE_APP_BASE_URL + `/cicd/delete/deploy_record/${form.id}`).then((res: any) => {
      table.loading = false
      table.request()
      ElMessage.success(res.message)
    })
  },
  build: (form: any) => {
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/build/${form.id}`).then((res: any) => {
      table.loading = false
      table.request()
      ElMessage.success(res.message)
    })
  },
  selectProject: (project_name: string) => {
    table.activeProject = project_name
    table.request()
  },
  startDeploy: (deployService: any) => {
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/start_deploy`, deployService).then((res: any) => {
      table.loading = false
      table.fetchDeployServiceRecord(table.activeDeployRecord)
      ElMessage.success(res.message)
    })
  },
  approveDeploy: (deployService: any) => {
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/approve_deploy`, deployService).then((res: any) => {
      table.loading = false
      table.fetchDeployServiceRecord(table.activeDeployRecord)
      ElMessage.success(res.message)
    })
  }
})

const selectedRows = ref<any[]>([])

const handleSelectionChange = (selected: any) => {
  selectedRows.value = selected
}

const deleteSelected = () => {
  if (selectedRows.value.length === 0) {
    return alert('请选择至少一行进行删除！')
  }

  selectedRows.value.forEach(row => {
    table.delete(row)
  })
}

// 初始化加载数据
table.request()
table.fetchProject()
</script>
