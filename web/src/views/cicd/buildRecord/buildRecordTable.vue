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
      <el-table-column prop="build_user" label="提交用户" sortable/>
      <el-table-column prop="created_at" label="创建时间" sortable/> 
      <el-table-column label="操作" fixed="right" width="200">
        <template #default="{ row }">
          <el-button link type="primary" @click="table.build(row)">构建</el-button>
          <el-divider direction="vertical" />
          <el-button link type="primary" @click="editForm.toDetail(row)">详情</el-button>
          <el-divider direction="vertical" />
          <el-button link type="primary" @click="editForm.toCreateDeploy(row)">发布</el-button>
          <el-divider direction="vertical" />
          <el-button link type="primary" @click="table.delete(row)">删除</el-button>
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
      size="40%"
    >
      <template #header>
        <h4>发布详情</h4>
      </template>
      <template #default>
        <el-table
          row-key="id"
          :data="table.buildServiceRecordData"
          table-layout="auto"
          max-height="900"
        >
          <el-table-column prop="service_name" label="服务名称" />
          <el-table-column prop="branch" label="构建分支" />
          <el-table-column prop="image" label="构建镜像" />
          <el-table-column prop="build_url" label="构建地址" />、
          <el-table-column prop="status" label="状态" />
        </el-table>
      </template>
      <template #footer>
        <div style="flex: auto">
          <el-button @click="editForm.show = false">关闭</el-button>
        </div>
      </template>
    </el-drawer>
    <el-drawer
      v-model="editForm.showCreateDeploy" 
      direction="rtl"
    >
      <template #header>
        <h4>创建发布</h4>
      </template>
      <template #default>
        <newDeployRecordIndex
          :buildRecordRow = table.activeRow
        />
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
import { ElMessage, type FormInstance } from 'element-plus'
import http from '@/api'
import newBuildRecordIndex from '@/views/cicd/buildRecord/newBuildRecord.vue'
import newDeployRecordIndex from '@/views/cicd/deployRecord/newDeployRecord.vue'

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
  showCreateDeploy: false,
  title: '',
  state: '',
  model: {
    name: '',
    project_name: '',
    env: '',
    tag: '',
    description: '',
    build_user: '',
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
      build_user: '',
    }
    table.fetchProject()
  },
  toDetail: (row: any) => {
    editForm.detailShow = true
    table.fetchBuildServiceRecord(row)
  },
  toCreateDeploy: (row: any) => {
    table.activeRow = row
    editForm.showCreateDeploy = true
  }
})
// 表格配置
const table = reactive({
  loading: false,
  border: true,
  activeProject: '请选择项目名称',
  activeRow: null as any,
  data: [] as any[],
  buildServiceRecordData: [] as any[],
  filteredData: [] as any[],
  projectData: [] as any[],
  request: () => {
    if (table.activeProject !== '请选择项目名称'){
      table.loading = true
        http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/build_record?project=${table.activeProject}`).then((res: any) => {
          table.data = res.data
          table.filteredData = res.data
          table.loading = false
      })
    }
  },
  fetchBuildServiceRecord: (form: any) => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/build_service_record?build_record=${form.name}`).then((res: any) => {
      table.buildServiceRecordData = res.data
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
    http.delete(import.meta.env.VITE_APP_BASE_URL + `/cicd/delete/build_record/${form.id}`).then((res: any) => {
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
