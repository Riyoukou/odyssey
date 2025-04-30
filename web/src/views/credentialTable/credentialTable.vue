<template>
  <div class="flex flex-col overflow-auto">
    <!-- 搜索表单 -->
    <ElForm v-if="search.show" :model="search.model" ref="search.ref" label-width="80px">
      <ElFormItem label="名称" prop="name">
        <ElInput v-model="search.model.name" placeholder="请输入名称" />
      </ElFormItem>
      <div class="flex justify-center mb-3">
        <ElButton icon="RefreshRight" @click="() => search.ref?.resetFields()">重置</ElButton>
        <ElButton type="primary" icon="Search" @click="credentialTable.request">查询</ElButton>
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
        <ElButton icon="Refresh" round @click="credentialTable.request"></ElButton>
        <ElButton icon="Search" round @click="search.show = !search.show"></ElButton>
      </div>
    </div>
    <!-- 列表 -->
    <el-table :data="credentialTable.data" style="width: 100%" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="凭证名称" sortable width="200" />
      <el-table-column prop="type" label="凭证类型" sortable width="200" />
      <el-table-column prop="description" label="凭证描述信息" />
      <el-table-column label="操作" fixed="right" width="120">
        <template #default="{ row }">
          <el-button link type="primary" @click="editForm.toView(row)">查看</el-button>
        </template>
      </el-table-column>
    </el-table>
    <ElDialog v-model="editForm.show" :title="editForm.title" width="800px" top="10vh">
      <ElForm :ref="(v: FormInstance | null) => (editForm.ref = v)" :model="editForm.model" label-width="80px" :rules="credentialRules">
        <ElFormItem label="凭证名称" prop="name">
          <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.name" placeholder="请输入凭证名称" />
        </ElFormItem>
        <ElFormItem label="凭证类型" prop="type">
          <ElSelect v-if="editForm.state !== 'view'" v-model="editForm.model.type" placeholder="请选择凭证类型">
            <ElOption label="无" value="none" />
            <ElOption label="Token" value="token" />
            <ElOption label="KubeConfig" value="kube_config" />
            <ElOption label="用户名密码" value="basic" />
          </ElSelect>
          <ElInput v-else v-model="editForm.model.type" placeholder="请选择凭证类型" disabled />
        </ElFormItem>
        <ElFormItem label="凭证内容" prop="data">
          <!-- 用户名密码形式 -->
          <template v-if="editForm.model.type === 'basic'">
            <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.username" placeholder="请输入用户名" style="margin-bottom: 8px" />
            <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.password" placeholder="请输入密码" show-password />
          </template>
          <template v-else-if="editForm.model.type === 'none'">
          </template>
          <!-- 其他类型 -->
          <template v-else-if="editForm.model.type === 'token' ">
            <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.data" placeholder="请输入Token内容" show-password />
          </template>
          <template v-else-if="editForm.model.type === 'kube_config' ">
            <ElInput v-if="editForm.state !== 'view'" v-model="editForm.model.data" placeholder="请输入凭证内容" type="textarea" :rows="5" />
            <ElInput v-else :disabled="editForm.state === 'view'" v-model="editForm.model.data" placeholder="请输入凭证内容" show-password disable/>
          </template>
        </ElFormItem>
        <ElFormItem label="描述" prop="description">
          <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.description" type="textarea" :rows="5" placeholder="请输入描述" />
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
const search = reactive({
  ref: null as FormInstance | null,
  show: false,
  model: {
    name: '',
    type: '',
    description: '',
    data: ''
  } ,
})

const credentialRules = computed(() => {
  if (editForm.model.type === 'basic') {
    return {
      name: [{ required: true, message: '请输入凭证名称', trigger: 'blur' }],
      type: [{ required: true, message: '请选择凭证类型', trigger: 'change' }],
      username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
      password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
    }
  } else if (editForm.model.type === 'none') {
    return {
      name: [{ required: true, message: '请输入凭证名称', trigger: 'blur' }],
      type: [{ required: true, message: '请选择凭证类型', trigger: 'blur' }],
    }
  } else {
    return {
      name: [{ required: true, message: '请输入凭证名称', trigger: 'blur' }],
      type: [{ required: true, message: '请选择凭证类型', trigger: 'blur' }],
      data: [{ required: true, message: '请输入凭证内容', trigger: 'blur' }],
    }
  }
})

// 表单配置
const editForm = reactive({
  ref: null as FormInstance | null,
  show: false,
  title: '',
  state: 'edit',
  model: {
    name: '',
    type: '',
    description: '',
    data: '',
    username: '',
    password: ''
  } ,
  toAdd: () => {
    editForm.ref?.resetFields()
    editForm.show = true
    editForm.title = '新增凭证'
    editForm.state = 'add'
    editForm.model = {
      name: '',
      type: '',
      description: '',
      data: '',
      username: '',
      password: ''
    }
  },
  toView: (row: any) => {
    editForm.show = true
    editForm.title = '查看凭证'
    editForm.state = 'view'
    editForm.model = { ...row }
    if (editForm.model.type === 'basic') {
      editForm.model.username =  JSON.parse(atob(row.data)).username
      editForm.model.password =  JSON.parse(atob(row.data)).password
    } else {
      editForm.model.data = atob(row.data)
    }
  },
  submit: () => {
    editForm.ref?.validate().then(() => {
      const { username, password, ...credential } = editForm.model
      // 如果是用户名密码类型，打包到 data 字段
      if (editForm.model.type === 'basic') {
        credential.data = JSON.stringify({ username, password })
      }  
      editForm.show = false
      credentialTable.create(credential)
    })
  }
})

// 表格配置
const credentialTable = reactive({
  loading: false,
  border: true,
  data: [],
  request: () => {
    credentialTable.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/credential`).then((res: any) => {
      credentialTable.data = res.data
      credentialTable.loading = false
    })
  },
  create: (form: any) => {
    credentialTable.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/create/credential`, form).then((res: any) => {
      credentialTable.loading = false
      credentialTable.request()
      ElMessage.success('创建成功')
    })
  },
  delete: (form: any) => {
    credentialTable.loading = true
    http.delete(import.meta.env.VITE_APP_BASE_URL + `/cicd/delete/credential/${form.id}`).then((res: any) => {
      credentialTable.loading = false
      credentialTable.request()
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
    credentialTable.delete(row)
  })
}

// 初始化加载数据
credentialTable.request()
</script>
