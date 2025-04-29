<template>
  <div class="flex flex-col overflow-auto">
    <!-- 搜索表单 -->
    <ElForm v-if="search.show" :model="search.model" ref="search.ref" label-width="80px">
      <ElFormItem label="姓名" prop="name">
        <ElInput v-model="search.model.name" placeholder="请输入姓名" />
      </ElFormItem>
      <ElFormItem label="年龄" prop="age">
        <ElSelect v-model="search.model.age" placeholder="请选择年龄" @change="(v: string) => console.log(v)">
          <ElOption label="12" value="12" />
          <ElOption label="13" value="13" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="邮件" prop="email">
        <ElInput v-model="search.model.email" placeholder="请输入邮件" />
      </ElFormItem>
      <ElFormItem label="出生日期" prop="date">
        <ElDatePicker v-model="search.model.date" type="datetime" placeholder="选择出生日期" />
      </ElFormItem>
      <div class="flex justify-center mb-3">
        <ElButton icon="RefreshRight" @click="() => search.ref?.resetFields()">重置</ElButton>
        <ElButton type="primary" icon="Search" @click="baseTable.request">查询</ElButton>
      </div>
    </ElForm>
    <!-- 按钮条 -->
    <div class="table-bar flex justify-between items-center mb-3">
      <div>
        <ElButton icon="Plus" @click="editForm.toAdd">新增</ElButton>
        <ElButton icon="Delete">删除</ElButton>
      </div>
      <div>
        <ElButton icon="Refresh" round @click="baseTable.request"></ElButton>
        <ElButton icon="Search" round @click="search.show = !search.show"></ElButton>
      </div>
    </div>
    <!-- 列表 -->
    <el-table :data="baseTable.data" style="width: 100%">
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
import { reactive } from 'vue'
import type { FormInstance } from 'element-plus'
import http from '@/api'

// 搜索表单配置
const search = reactive({
  ref: null as FormInstance | null,
  show: false,
  model: { name: '', age: '', email: '', date: '' },
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
      baseTable.request()
    })
  }
})

// 表格配置
const baseTable = reactive({
  loading: false,
  border: true,
  data: [],
  page: {
    sortOrder: null,
    sortProp: '',
    currentPage: 1,
    pageSize: 10,
    total: 0
  },
  request: () => {
    baseTable.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/cluster`).then((res) => {
      baseTable.data = res.data
      baseTable.loading = false
    })
  },
})

// 初始化加载数据
baseTable.request()
</script>
