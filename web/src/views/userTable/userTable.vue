<template>
  <div class="flex flex-col overflow-auto">
    <!-- 搜索表单 -->
    <ElForm v-if="search.show" :model="search.model" ref="search.ref" label-width="80px">
      <ElFormItem label="名称" prop="name">
        <ElInput v-model="search.model.name" placeholder="请输入名称" />
      </ElFormItem>
      <div class="flex justify-center mb-3">
        <ElButton icon="RefreshRight" @click="() => search.ref?.resetFields()">重置</ElButton>
        <ElButton type="primary" icon="Search" @click="userTable.request">查询</ElButton>
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
        <ElButton icon="Refresh" round @click="userTable.request"></ElButton>
        <ElButton icon="Search" round @click="search.show = !search.show"></ElButton>
      </div>
    </div>

    <!-- 列表 -->
    <el-table :data="userTable.data" style="width: 100%" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="用户名" sortable width="200" />
      <el-table-column prop="email" label="邮箱" sortable width="200" />
      <el-table-column prop="phone" label="手机号" width="150" />
      <el-table-column prop="role" label="角色" width="150" />
      <el-table-column prop="last_login" label="最后登陆时间" show-overflow-tooltip />
      <el-table-column label="操作" fixed="right" width="120">
        <template #default="{ row }">
          <el-button link type="primary" @click="editForm.toView(row)">查看</el-button>
          <el-divider direction="vertical" />
          <el-button link type="primary" @click="editForm.toEdit(row)">编辑</el-button>
          <el-divider direction="vertical" />
          <el-button link type="primary" @click="editForm.toEditPassword(row)">修改密码</el-button>
        </template>
      </el-table-column>
    </el-table>

    <ElDialog v-model="editForm.show" :title="editForm.title" width="800px" top="10vh">
      <ElForm 
        :ref="(v: FormInstance | null) => (editForm.ref = v)" 
        :model="editForm.model" 
        label-width="80px"
        :rules="clusterRules"
      >
        <template v-if="editForm.state === 'editPassword'">
          <ElFormItem label="当前密码" prop="old_password">
            <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.old_password" placeholder="请输入" />
          </ElFormItem>
          <ElFormItem label="新的密码" prop="new_password" >
            <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.new_password" placeholder="请输入" />
          </ElFormItem>
        </template>
        <template v-else>
          <ElFormItem label="用户名" prop="name">
            <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.name" placeholder="请输入用户名称" />
          </ElFormItem>
          <ElFormItem label="邮箱" prop="email" >
            <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.email" placeholder="请输入邮箱" />
          </ElFormItem>
          <ElFormItem label="手机号" prop="phone">
            <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.phone" placeholder="请输入手机号" />
          </ElFormItem>
          <ElFormItem label="角色" prop="role">
            <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.role" placeholder="请输入角色" />
          </ElFormItem>
        </template>
      </ElForm>
      <template #footer>
        <ElButton v-if="editForm.state === 'add'" type="primary" @click="editForm.submit">提交</ElButton>
        <ElButton v-else-if="editForm.state === 'edit'" type="primary" @click="editForm.editSubmit">提交</ElButton>
        <ElButton v-else-if="editForm.state === 'editPassword'" type="primary" @click="editForm.editPasswordSubmit">提交</ElButton>
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
    password: '',
    email: '',
    phone: '',
  } 
})

const clusterRules = computed(() => {
  if (editForm.state === 'editPassword') {
    return {
      old_password: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
      new_password: [{ required: true, message: '请输入新密码', trigger: 'blur' }],
    }
  } else {
    return {
      name: [{ required: true, message: '请输入集群名称', trigger: 'blur' }],
      password: [{ required: true, message: '请输入APIServer', trigger: 'blur' }],
      email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }] ,
      phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }],
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
    password: '',
    email: '',
    phone: '',
    old_password: '',
    new_password: '',
  },
  toAdd: () => {
    editForm.ref?.resetFields()
    editForm.show = true
    editForm.title = '新增集群'
    editForm.state = 'add'
    editForm.model = {
      name: '',
      password: '',
      email: '',
      phone: '',
      old_password: '',
      new_password: '',
    }
  },
  toEdit: (row: any) => {
    editForm.show = true
    editForm.title = '编辑用户'
    editForm.state = 'edit'
    editForm.model = { ...row }
  },
  toEditPassword: (row: any) => {
    editForm.show = true
    editForm.title = '修改密码'
    editForm.state = 'editPassword'
    editForm.model = { ...row }
  },
  toView: (row: any) => {
    editForm.show = true
    editForm.title = '查看用户'
    editForm.state = 'view'
    editForm.model = { ...row }
  },
  submit: () => {
    editForm.ref?.validate().then(() => {
      const { new_password, old_password, ...user } = editForm.model
      editForm.show = false
      userTable.create(user)
    })
  },
  editSubmit: () => {
    editForm.ref?.validate().then(() => {
      editForm.show = false
      userTable.edit(editForm.model)
    })
  },
  editPasswordSubmit: () => {
    editForm.ref?.validate().then(() => {
      editForm.show = false
      userTable.editPassword(editForm.model)
    })
  }
})
// 表格配置
const userTable = reactive({
  loading: false,
  border: true,
  data: [],
  credentialData:[],
  request: () => {
    userTable.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/user/fetch/user`).then((res: any) => {
      userTable.data = res.data
      userTable.loading = false
    })
  },
  create: (form: any) => {
    userTable.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/user/create/user`, form).then((res: any) => {
      userTable.loading = false
      userTable.request()
      ElMessage.success('新增成功')
    })
  },
  edit: (form: any) => {
    userTable.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/user/update/user`, form).then((res: any) => {
      userTable.loading = false
      userTable.request()
      ElMessage.success('编辑成功')
    })
  },
  editPassword: (form: any) => {
    userTable.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/user/update/user_password`, form).then((res: any) => {
      userTable.loading = false
      userTable.request()
      ElMessage.success('编辑成功')
    })
  },
  delete: (form: any) => {
    userTable.loading = true
    http.delete(import.meta.env.VITE_APP_BASE_URL + `/user/delete/user/${form.id}`).then((res: any) => {
      userTable.loading = false
      userTable.request()
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
    userTable.delete(row)
  })
}

// 初始化加载数据
userTable.request()
</script>
