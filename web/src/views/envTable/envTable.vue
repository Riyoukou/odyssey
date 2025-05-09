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
    <el-table :data="table.filteredData" style="width: 100%" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="环境名称" sortable/>
      <el-table-column prop="type" label="环境类型" sortable/>
      <el-table-column prop="project_name" label="项目名称" sortable/>
      <el-table-column prop="namespace" label="命名空间" sortable/>
      <el-table-column label="操作" fixed="right" width="120">
        <template #default="{ row }">
          <el-button link type="primary" @click="editForm.toView(row)">查看</el-button>
          <el-divider direction="vertical" />
          <el-button link type="primary" @click="editForm.toEdit(row)">编辑</el-button>
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
        <ElFormItem label="环境名称" prop="name">
          <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.name" placeholder="请输入项目名称" />
        </ElFormItem>
        <ElFormItem label="环境类型" prop="type">
          <ElSelect :disabled="editForm.state === 'view'" v-model="editForm.model.type" placeholder="请选择环境类型">
            <ElOption label="测试环境" value="dev" />
            <ElOption label="预发布环境" value="pre" />
            <ElOption label="生产环境" value="prod" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="项目名称" prop="project_name">
          <ElSelect disabled v-model="editForm.model.project_name" placeholder="请选择项目名称" @change="table.selectProject">
            <ElOption v-for="item in table.projectData" :key="item.name" :label="item.name" :value="item.name" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="名称空间" prop="namespace">
          <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.namespace" placeholder="请输入项目名称" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton v-if="editForm.state === 'add'" type="primary" @click="editForm.submit">提交</ElButton>
        <ElButton v-else-if="editForm.state === 'edit'" type="primary" @click="editForm.editSubmit">提交</ElButton>
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
  title: '',
  state: '',
  model: {
    name: '',
    type: '',
    project_name: '',
    namespace: '',
  },
  toAdd: () => {
    editForm.ref?.resetFields()
    editForm.show = true
    editForm.title = '新增环境'
    editForm.state = 'add'
    editForm.model = {
      name: '',
      type: '',
      project_name: '',
      namespace: '',
    }
    table.fetchProject()
  },
  toView: (row: any) => {
    editForm.show = true
    editForm.title = '查看环境'
    editForm.state = 'view'
    editForm.model = { ...row }
  },
  toEdit: (row: any) => {
    editForm.show = true
    editForm.title = '编辑环境'
    editForm.state = 'edit'
    editForm.model = { ...row }
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
  }
})
// 表格配置
const table = reactive({
  loading: false,
  border: true,
  activeProject: '请选择项目名称',
  data: [] as any[],
  filteredData: [] as any[],
  projectData: [] as any[],
  request: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/env?project=${table.activeProject}`).then((res: any) => {
      table.data = res.data
      table.filteredData = res.data
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
  create: (form: any) => {
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/create/env`, form).then((res: any) => {
      table.loading = false
      table.request()
      ElMessage.success('新增成功')
    })
  },
  edit: (form: any) => {
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/update/env`, form).then((res: any) => {
      table.loading = false
      table.request()
      ElMessage.success('编辑成功')
    })
  },
  delete: (form: any) => {
    table.loading = true
    http.delete(import.meta.env.VITE_APP_BASE_URL + `/cicd/delete/env/${form.id}`).then((res: any) => {
      table.loading = false
      table.request()
      ElMessage.success('删除成功')
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
