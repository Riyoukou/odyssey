<template>
  <div class="flex flex-col overflow-auto">
    <!-- 查询 -->
    <ElForm :ref="(v) => (search.ref = v)" v-show="search.show" :model="search.model" label-width="80px">
      <AgelFormGrid :items="search.items" responsive></AgelFormGrid>
      <div class="flex justify-center mb-3">
        <ElButton icon="RefreshRight" @click="() => search.ref?.resetFields()">重置</ElButton>
        <ElButton type="primary" icon="Search" @click="table.request">查询</ElButton>
      </div>
    </ElForm>
    <!-- 按钮条 -->
    <div class="table-bar flex justify-between items-center mb-3">
      <div>
        <ElButton icon="Plus" @click="form.toAdd">新增</ElButton>
        <ElButton icon="Delete">删除</ElButton>
      </div>
      <div>
        <ElButton icon="Refresh" round @click="table.request"></ElButton>
        <ElButton icon="Search" round @click="search.show = !search.show"></ElButton>
      </div>
    </div>
    <!-- 列表 -->
    <AgelTable class="flex-1" v-bind="table" v-model:page="table.page">
      <template #operation="scope">
        <ElButton link type="primary" @click="form.toView(scope.row)">查看</ElButton>
        <ElDivider direction="vertical" />
        <ElButton link type="primary" @click="form.toEdit(scope.row)">编辑</ElButton>
      </template>
    </AgelTable>
    <!-- 弹窗表单 -->
    <ElDialog v-model="form.show" :title="form.title" width="800px" top="10vh">
      <ElForm :ref="(v) => (form.ref = v)" :model="form.model" label-width="80px">
        <ElRow :gutter="20">
          <ElCol :span="8">
            <ElFormItem label="姓名" prop="name" required>
              <ElInput v-model="form.model.name" :disabled="form.state === 'view'" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="8">
            <ElFormItem label="年龄" prop="age" required>
              <ElInput v-model="form.model.age" :disabled="form.state === 'view'" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="8">
            <ElFormItem label="出生日期" prop="date">
              <ElDatePicker v-model="form.model.date" :disabled="form.state === 'view'" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="24">
            <ElFormItem label="邮件" prop="email">
              <ElInput v-model="form.model.email" :disabled="form.state === 'view'" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="24">
            <ElFormItem label="介绍" prop="decs">
              <ElInput type="textarea" v-model="form.model.decs" :rows="5" :disabled="form.state === 'view'" />
            </ElFormItem>
          </ElCol>
        </ElRow>
      </ElForm>
      <template #footer>
        <ElButton v-if="form.state !== 'view'" type="primary" @click="form.submit">提交</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script lang="ts" setup>
import { reactive } from 'vue'
import http from '@/api'

// 搜索表单配置
const search = reactive({
  ref: null,
  show: false,
  model: { name: '', age: '', email: '', date: '' },
  items: [
    { label: '姓名', prop: 'name' },
    {
      label: '年龄',
      prop: 'age',
      slot: 'agel-select',
      attrs: {
        options: ['12', '13'],
        onChange: (v: string) => console.log(v)
      }
    },
    { label: '邮件', prop: 'email' },
    {
      label: '出生日期',
      prop: 'date',
      slot: 'el-date-picker',
      attrs: { type: 'datetime' }
    }
  ]
})

// 表单配置
const form = reactive({
  ref: null,
  show: false,
  title: '',
  state: 'edit',
  model: {
    name: '',
    api_server: '',
    region: '',
    version: '',
    description: '',
    config: ''
  },
  toAdd: () => {
    form.show = true
    form.title = '新增集群'
    form.state = 'add'
    nextTick(() => {
      form.ref?.resetFields()
    })
  },
  toEdit: (row: TableRow) => {
    form.show = true
    form.title = '编辑集群'
    form.state = 'edit'
    nextTick(() => {
      form.model = { ...row }
    })
  },
  submit: () => {
    form.ref?.validate().then(() => {
      form.show = false
      table.create(form.model)
    })
  },
  submitEdit: () => {
    form.ref?.validate().then(() => {
      form.show = false
      table.edit(form.model)
    })
  }
})

// 表格配置
const table = reactive({
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
  columns: [
  { label: '#', type: 'selection' },
    { label: '集群名称', prop: 'name', width: 200 },
    { label: 'APIServer', prop: 'api_server', width: 200 },
    { label: 'Region', prop: 'region', width: 150 },
    { label: '集群版本', prop: 'version', width: 150 },
    { label: '介绍', prop: 'description', showOverflowTooltip: true },
    { width: '120px', label: '操作', fixed: 'right', slot: 'operation' }
  ],
  request: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/cluster`).then((res) => {
      table.data = res.data.result
      table.loading = false
    })
  },
})

// 初始化加载数据
table.request()
</script>
