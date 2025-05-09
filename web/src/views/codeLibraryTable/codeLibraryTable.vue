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
      <div>
        <ElButton icon="Refresh" round @click="table.request"></ElButton>
        <ElButton icon="Search" round @click="search.show = !search.show"></ElButton>
      </div>
    </div>
    <!-- 列表 -->
    <el-table :data="table.filteredData" style="width: 100%" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="代码库名称" sortable width="200" />
      <el-table-column prop="type" label="代码库类型" sortable width="200" />
      <el-table-column prop="url" label="代码库地址" />
      <el-table-column prop="code_source_name" label="代码源名称" />
      <el-table-column prop="project_name" label="所属项目" />
      <el-table-column label="操作" fixed="right" width="120">
        <template #default="{ row }">
          <el-button link type="primary" @click="editForm.toView(row)">查看</el-button>
        </template>
      </el-table-column>
    </el-table>
    <ElDialog v-model="editForm.show" :title="editForm.title" width="800px" top="10vh">
      <ElForm :ref="(v: FormInstance | null) => (editForm.ref = v)" :model="editForm.model" label-width="120px" :rules="rules">
        <ElFormItem label="项目名称" prop="project_name">
          <ElSelect 
            :disabled="editForm.state === 'view'"
            v-model="editForm.model.project_name" 
            placeholder="请选择项目" 
            @focus="table.requestProject"
            filterable
          >
            <ElOption  
                v-for="item in table.projectData" 
                :key="item.name" 
                :label="item.name" 
                :value="item.name" 
              />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="代码库类型" prop="type">
          <ElSelect :disabled="editForm.state === 'view'" v-model="editForm.model.type" placeholder="请选择类型">
            <el-option label="CodeLibrary" value="code_library"></el-option>
            <el-option label="GitOps" value="gitops"></el-option>
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="代码源名称" prop="code_source_name">
          <ElSelect 
            :disabled="editForm.state === 'view'"
            v-model="editForm.model.code_source_name" 
            placeholder="请选择代码源" 
            @focus="table.requestCodeSource"
            filterable
          >
            <ElOption  
                v-for="item in table.filtereCodeLibraryData" 
                :key="item.name" 
                :label="item.name" 
                :value="item.name" 
              />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="代码库" prop="name">
          <ElSelect
            :disabled="editForm.state === 'view'"
            v-model="editForm.model.name" 
            placeholder="请选择代码库" 
            @focus="table.requestGitProject"  
            @change="editForm.selectGitProject"
            value-key="name"
            filterable
          >
            <ElOption  
              v-for="item in table.gitProjectData" 
              :key="item.name" 
              :label="item.name" 
              :value="item" 
            />
          </ElSelect>
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton v-if="editForm.state !== 'view'" type="primary" @click="editForm.submit">提交</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import http from '@/api'

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
  return {
    project_name: [{ required: true, message: '请选择项目名称', trigger: 'blur' }],
    type: [{ required: true, message: '请选择代码库类型', trigger: 'blur' }],
    code_source_name: [{ required: true, message: '请选择代码源', trigger: 'blur' }],
    name: [{ required: true, message: '请选择代码库', trigger: 'blur' }],
  }
})

// 表单配置
const editForm = reactive({
  ref: null as FormInstance | null,
  show: false,
  title: '',
  state: '',
  model: {
    name: '',
    url: '',
    code_source_name: '',
    project_name: '',
    project_id: 0,
    type: '',
  } ,
  toAdd: () => {
    editForm.ref?.resetFields()
    editForm.show = true
    editForm.title = '新增代码库'
    editForm.state = 'add'
    editForm.model = {
      name: '',
      url: '',
      code_source_name: '',
      project_name: '',
      project_id: 0,
      type: '',
    }
  },
  toView: (row: any) => {
    editForm.show = true
    editForm.title = '查看代码库'
    editForm.state = 'view'
    editForm.model = { ...row }
  },
  selectGitProject: (item : any) => {
   editForm.model.name = item.name,
   editForm.model.url = item.web_url,
   editForm.model.project_id = item.id
   console.log(editForm.model)
  },
  submit: () => {
    editForm.ref?.validate().then(() => {
      editForm.show = false
      table.create(editForm.model)
    })
  }
})

// 表格配置
const table = reactive({
  loading: false,
  border: true,
  data: [] as any[],
  projectData: [] as any[],
  codeLibraryData: [] as any[],
  gitProjectData: [] as any[],
  filtereCodeLibraryData: [] as any[],
  filteredData: [] as any[],
  request: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/code_library`).then((res: any) => {
      table.data = res.data
      table.filteredData = res.data
      table.loading = false
    })
  },
  requestProject: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/project`).then((res: any) => {
      table.projectData = res.data
      table.loading = false
    })
  },
  requestCodeSource: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/cicd_tool`).then((res: any) => {
      table.codeLibraryData = res.data
      table.filtereCodeLibraryData = table.codeLibraryData.filter( tool => tool.type === 'git')
      table.loading = false
    })
  },
  requestGitProject: () => {
    table.loading = true
    if (editForm.model.code_source_name !== '') {
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/git_project?cicd_tool=${editForm.model.code_source_name}`).then((res: any) => {
      table.gitProjectData = res.data
      table.loading = false
    })
    } else {
      ElMessage.error('请选择代码源')
    }
  },
  create: (form: any) => {
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/create/code_library`, form).then((res: any) => {
      table.loading = false
      table.request()
      ElMessage.success('创建成功')
    })
  },
  delete: (form: any) => {
    table.loading = true
    http.delete(import.meta.env.VITE_APP_BASE_URL + `/cicd/delete/code_library/${form.id}`).then((res: any) => {
      table.loading = false
      table.request()
      ElMessage.success('删除成功')
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
</script>
