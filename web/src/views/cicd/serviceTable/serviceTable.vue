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
        <ElSelect class="w-50" @change="table.selectEnv" :placeholder=table.activeEnv >
          <ElOption v-for="item in table.envData" :key="item.name" :label="item.name" :value="item.name" />
        </ElSelect>
        <ElButton icon="Refresh" round @click="table.request" class="ml-4"></ElButton>
        <ElButton icon="Search" round @click="search.show = !search.show"></ElButton>
      </div>
    </div>

    <!-- 列表 -->
    <el-table :data="table.filteredData" style="width: 100%" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="服务名称" sortable/>
      <el-table-column prop="project_name" label="所属项目" sortable/>
      <el-table-column prop="env_name" label="环境" sortable/>
      <el-table-column prop="clusters" label="发布集群" sortable/>
      <el-table-column label="操作" fixed="right" width="200">
        <template #default="{ row }">
          <el-button link type="primary" @click="editForm.toView(row)">查看</el-button>
          <el-divider direction="vertical" />
          <el-button link type="primary" @click="editForm.toEdit(row)">编辑</el-button>
          <el-divider direction="vertical" />
          <el-button link type="primary" @click="table.delete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <ElDialog v-model="editForm.show" :title="editForm.title" width="800px" top="10vh">
      <ElForm 
        :ref="(v: FormInstance | null) => (editForm.ref = v)" 
        :model="editForm.model" 
        label-width="80px"
        :rules="rules"
      >
        <ElFormItem label="服务名称" prop="name">
          <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.name" placeholder="请输入服务名称" />
        </ElFormItem>
        <ElFormItem label="所属项目" prop="project_name">
          <ElSelect disabled v-model="editForm.model.project_name" placeholder="请选择项目名称">
            <ElOption v-for="item in table.envData" :key="item.name" :label="item.name" :value="item.name" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="环境" prop="env_name">
          <ElSelect disabled v-model="editForm.model.env_name" placeholder="请选择环境名称">
            <ElOption v-for="item in table.envData" :key="item.name" :label="item.name" :value="item.name" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="发布集群" prop="clusters">
          <ElSelect multiple :disabled="editForm.state === 'view'" v-model="editForm.model.clusters" placeholder="请选择集群">
            <ElOption v-for="item in table.clusterData" :key="item.name" :label="item.name" :value="item.name" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="代码库" prop="code_library_name">
          <ElSelect :disabled="editForm.state === 'view'" v-model="editForm.model.code_library_name" placeholder="请选择代码库名称">
            <ElOption v-for="item in table.codeLibraryData" :key="item.name" :label="item.name" :value="item.name" />
          </ElSelect>
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton v-if="editForm.state === 'edit'" type="primary" @click="editForm.toBuildMap">构建信息</ElButton>
        <ElButton v-if="editForm.state === 'edit'" type="primary" @click="editForm.toDeployMap">发布信息</ElButton>
        <ElButton v-if="editForm.state === 'add'" type="primary" @click="editForm.submit">提交</ElButton>
        <ElButton v-else-if="editForm.state === 'edit'" type="primary" @click="editForm.editSubmit">提交</ElButton>
      </template>
    </ElDialog>
    <el-drawer v-model="editForm.buildMapShow" direction="btt" size="70%">
      <template #header>
        <h4>CI信息</h4>
      </template>
      <template #default>   
        <buildMap 
          :activeBuildMap = editForm.model.build_map
          :activeID = editForm.activeID
        />
      </template>
      <template #footer>
      </template>
    </el-drawer>
    <el-drawer v-model="editForm.deployMapShow" direction="btt" size="100%">
      <template #header>
        <h4>CD信息</h4>
      </template>
      <template #default>
        <el-tabs v-model="editForm.activeTab" class="demo-tabs" @tab-click="editForm.tabClick">
          <el-tab-pane v-for="item in editForm.model.clusters" :key="item" :label="item" :name="item">      
            <deployMap 
              :activeCluster = editForm.activeTab
              :activeDeployMap = editForm.model.deploy_map[item]
              :activeID = editForm.activeID
            />
          </el-tab-pane>
        </el-tabs>
      </template>
      <template #footer>
      </template>
    </el-drawer>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import http from '@/api'
import deployMap from '@/views/cicd/serviceTable/deployMap.vue'
import buildMap from '@/views/cicd/serviceTable/buildMap.vue'

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
    project_name: [{ required: true, message: '请选择项目名称', trigger: 'change' }],
    clusters: [{ required: true, message: '请选择集群', trigger: 'change' }],
    env_name: [{ required: true, message: '请选择环境名称', trigger: 'change' }],
    code_library_name: [{ required: true, message: '请选择代码库名称', trigger: 'change' }],
  }
  return baseRules
})

// 表单配置
const editForm = reactive({
  ref: null as FormInstance | null,
  show: false,
  buildMapShow: false,
  deployMapShow: false,
  title: '',
  state: '',
  activeTab: '',
  activeID: 0,
  model: {
    name: '',
    project_name: '',
    clusters: [] as string[],
    env_name: '',
    code_library_name: '',
    build_map: {} as Record<string, any>,
    deploy_map: {} as Record<string, any[]>,
  },
  toAdd: () => {
    editForm.ref?.resetFields()
    editForm.show = true
    editForm.title = '新增环境'
    editForm.state = 'add'
    editForm.model = {
      name: '',
      project_name: table.activeProject,
      code_library_name: '',
      clusters: [],
      env_name: table.activeEnv,
      build_map: {} as Record<string, any>,
      deploy_map: {} as Record<string, any[]>,
    }
    table.fetchCluster()
    table.fetchCodeLibrary()
  },
  toBuildMap: () => {
    editForm.buildMapShow = true
  },
  toDeployMap: () => {
    editForm.deployMapShow = true
    editForm.activeTab = editForm.model.clusters[0]
  },
  toView: (row: any) => {
    editForm.show = true
    editForm.title = '查看环境'
    editForm.state = 'view'
    editForm.activeID = row.id
    editForm.model = { ...row }
  },
  toEdit: (row: any) => {
    editForm.show = true
    editForm.title = '编辑环境'
    editForm.state = 'edit'
    editForm.model = { ...row }
    editForm.activeID = row.id
    table.fetchProject()
  },
  submit: () => {
    editForm.ref?.validate().then(() => {
      editForm.show = false
      table.create(editForm.model)
    })
  },
  editSubmit: () => {
    editForm.ref?.validate().then(() => {
      editForm.show = false
      table.edit(editForm.model)
    })
  },
  tabClick: (tab: any) => {
    editForm.activeTab = tab.props.name
  }
})
// 表格配置
const table = reactive({
  loading: false,
  border: true,
  activeProject: '请选择项目名称',
  activeEnv: '请选择环境名称',
  data: [] as any[],
  envData: [] as any[],
  filteredData: [] as any[],
  projectData: [] as any[],
  clusterData: [] as any[],
  codeLibraryData: [] as any[],
  request: () => {
    if (table.activeProject !== '请选择项目名称' && table.activeEnv !== '请选择环境名称') {
      table.loading = true
      http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/service?project=${table.activeProject}&env=${table.activeEnv}`).then((res: any) => {
        table.data = res.data
        table.filteredData = res.data
        table.loading = false
      })
    }
  },
  fetchProject: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/project`).then((res: any) => {
      table.projectData = res.data
      table.loading = false
    })
  },
  fetchEnv: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/env?project=${table.activeProject}`).then((res: any) => {
      table.envData = res.data
      table.loading = false
    })
  },
  fetchCluster: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/cluster`).then((res: any) => {
      table.clusterData = res.data
      table.loading = false
    })
  },
  fetchCodeLibrary: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/code_library`).then((res: any) => {
      table.codeLibraryData = res.data
      table.loading = false
    })
  },   
  create: (form: any) => {
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/create/service`, form).then((res: any) => {
      table.loading = false
      table.request()
      ElMessage.success(res.message)
    })
  },
  edit: (form: any) => {
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/update/service`, form).then((res: any) => {
      table.loading = false
      table.request()
      ElMessage.success(res.message)
    })
  },
  delete: (form: any) => {
    table.loading = true
    http.delete(import.meta.env.VITE_APP_BASE_URL + `/cicd/delete/service/${form.id}`).then((res: any) => {
      table.loading = false
      table.request()
      ElMessage.success(res.message)
    })
  },
  selectProject: (project_name: string) => {
    table.activeProject = project_name
    table.fetchEnv()
  },
  selectEnv: (env_name: string) => {
    table.activeEnv = env_name
    table.request()
  },
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
