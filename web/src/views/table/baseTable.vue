<template>
  <div class="flex flex-col overflow-auto">
    <!-- 搜索表单 -->
    <ElForm v-if="search.show" :model="search.model" ref="searchRef" label-width="80px">
      <ElFormItem label="姓名" prop="name">
        <ElInput v-model="search.model.name" placeholder="请输入姓名" />
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
        <ElButton icon="Delete">删除</ElButton>
      </div>
      <div>
        <ElButton icon="Refresh" round @click="table.request"></ElButton>
        <ElButton icon="Search" round @click="search.show = !search.show"></ElButton>
      </div>
    </div>

    <!-- 列表 -->
    <el-table :data="table.data" style="width: 100%">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="姓名" sortable />
      <el-table-column prop="age" label="年龄" sortable />
      <el-table-column prop="email" label="邮件" width="250" />
      <el-table-column prop="date" label="出生日期" width="150" />
      <el-table-column prop="desc" label="介绍" show-overflow-tooltip />
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
      >
        <ElFormItem label="姓名" prop="name" :rules="{ required: true, message: '请输入姓名', trigger: 'blur' }">
          <ElInput v-model="editForm.model.name" placeholder="请输入姓名" />
        </ElFormItem>
        <ElFormItem label="年龄" prop="age" :rules="{ required: true, message: '请输入年龄', trigger: 'blur' }">
          <ElInput v-model="editForm.model.age" placeholder="请输入年龄" />
        </ElFormItem>
        <ElFormItem label="出生日期" prop="date">
          <ElDatePicker v-model="editForm.model.date" type="date" placeholder="选择出生日期" />
        </ElFormItem>
        <ElFormItem label="邮件" prop="email">
          <ElInput v-model="editForm.model.email" placeholder="请输入邮件" />
        </ElFormItem>
        <ElFormItem label="介绍" prop="decs">
          <ElInput
            v-model="editForm.model.decs"
            type="textarea"
            :rows="5"
            placeholder="请输入介绍"
          />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton v-if="editForm.state !== 'view'" type="primary" @click="editForm.submit">提交</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script lang="ts" setup>
import { ref,reactive } from 'vue'
import type { FormInstance } from 'element-plus'
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

// 表单配置
const editForm = reactive({
  ref: null as FormInstance | null,
  show: false,
  title: '',
  state: 'edit',
  model: { name: '', age: '', date: '', email: '', decs: '' },
  toAdd: () => {
    editForm.show = true
    editForm.title = '新增用户'
    editForm.state = 'add'
    editForm.ref?.resetFields()
  },
  toEdit: (row: any) => {
    editForm.show = true
    editForm.title = '编辑用户资料'
    editForm.state = 'edit'
    editForm.model = { ...row }
  },
  toView: (row: any) => {
    editForm.show = true
    editForm.title = '查看用户资料'
    editForm.state = 'view'
    editForm.model = { ...row }
  },
  submit: () => {
    editForm.ref?.validate().then(() => {
      editForm.show = false
      table.request()
    })
  }
})

// 表格配置
const table = reactive({
  loading: false,
  border: true,
  data: [] as any[],
  filteredData: [] as any[],
  page: {
    sortOrder: null,
    sortProp: '',
    currentPage: 1,
    pageSize: 10,
    total: 0
  },
  request: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/cluster`).then((res) => {
      table.data = res.data
      table.loading = false
    })
  },
})

// 初始化加载数据
table.request()
</script>
