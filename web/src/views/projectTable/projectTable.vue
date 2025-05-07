<template>
  <div class="flex flex-col overflow-auto">
    <!-- 搜索表单 -->
    <ElForm v-if="search.show" :model="search.model" ref="search.ref" label-width="80px">
      <ElFormItem label="名称" prop="name">
        <ElInput v-model="search.model.name" placeholder="请输入名称" />
      </ElFormItem>
      <div class="flex justify-center mb-3">
        <ElButton icon="RefreshRight" @click="() => search.ref?.resetFields()">重置</ElButton>
        <ElButton type="primary" icon="Search" @click="projectTable.request">查询</ElButton>
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
        <ElButton icon="Refresh" round @click="projectTable.request"></ElButton>
        <ElButton icon="Search" round @click="search.show = !search.show"></ElButton>
      </div>
    </div>

    <!-- 列表 -->
    <el-table :data="projectTable.data" style="width: 100%" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="项目名称" sortable/>
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
        :rules="projectRules"
      >
        <ElFormItem label="项目名称" prop="name">
          <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.name" placeholder="请输入项目名称" />
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
const search = reactive({
  ref: null as FormInstance | null,
  show: false,
  model: {
    name: '',
  } 
})

const projectRules = computed(() => {
  const baseRules = {
    name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
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
  } ,
  toAdd: () => {
    editForm.ref?.resetFields()
    editForm.show = true
    editForm.title = '新增项目'
    editForm.state = 'add'
    editForm.model = {
      name: '',
    }
  },
  toEdit: (row: any) => {
    editForm.show = true
    editForm.title = '编辑项目'
    editForm.state = 'edit'
    editForm.model = { ...row }
  },
  toView: (row: any) => {
    editForm.show = true
    editForm.title = '查看项目'
    editForm.state = 'view'
    editForm.model = { ...row }
  },
  submit: () => {
    editForm.ref?.validate().then(() => {
      editForm.show = false
      projectTable.create(editForm.model)
    })
  }
})
// 表格配置
const projectTable = reactive({
  loading: false,
  border: true,
  data: [],
  request: () => {
    projectTable.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/project`).then((res: any) => {
      projectTable.data = res.data
      projectTable.loading = false
    })
  },
  create: (form: any) => {
    projectTable.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/create/project`, form).then((res: any) => {
      projectTable.loading = false
      projectTable.request()
      ElMessage.success('新增成功')
    })
  },
  delete: (form: any) => {
    projectTable.loading = true
    http.delete(import.meta.env.VITE_APP_BASE_URL + `/cicd/delete/project/${form.id}`).then((res: any) => {
      projectTable.loading = false
      projectTable.request()
      ElMessage.success('删除成功')
    })
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
    projectTable.delete(row)
  })
}

// 初始化加载数据
projectTable.request()
</script>
