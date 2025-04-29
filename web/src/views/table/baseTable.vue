<template>
  <div class="flex flex-col overflow-auto">
    <!-- 查询 -->
    <ElForm :ref="(v: FormInstance | null) => (search.ref = v)" v-show="search.show" :model="search.model" label-width="80px">
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
      <template #operation="scope: { row: FormModel }">
        <ElButton link type="primary" @click="form.toView(scope.row)">查看</ElButton>
        <ElDivider direction="vertical" />
        <ElButton link type="primary" @click="form.toEdit(scope.row)">编辑</ElButton>
      </template>
    </AgelTable>
    <!-- 弹窗表单 -->
    <ElDialog v-model="form.show" :title="form.title" width="800px" top="10vh">
      <ElForm :ref="(v: FormInstance | null) => (form.ref = v)" :model="form.model" label-width="80px">
        <AgelFormDesc :items="form.items" :view-model="form.state == 'view'"></AgelFormDesc>
      </ElForm>
      <template #footer>
        <ElButton v-if="form.state !== 'view'" type="primary" @click="form.submit">提交</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script lang="ts" setup>
import { reactive } from 'vue'
import type { FormInstance } from 'element-plus'
import http from '@/api'

interface SearchModel {
  name: string
  age: string
  email: string
  date: string | null
}

interface FormModel {
  name: string
  age: string
  date: string | null
  email: string
  decs: string
}

interface TablePage {
  sortOrder: string | null
  sortProp: string
  currentPage: number
  pageSize: number
  total: number
}

// 搜索表单配置
const search = reactive({
  ref: null as FormInstance | null,
  show: false,
  model: { name: '', age: '', email: '', date: '' } as SearchModel,
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
  ref: null as FormInstance | null,
  show: false,
  title: '',
  state: 'edit' as 'edit' | 'add' | 'view',
  items: [
    { label: '姓名', prop: 'name' },
    { label: '年龄', prop: 'age' },
    { label: '出生日期', prop: 'date' },
    { label: '邮件', prop: 'email' },
    { label: '介绍', prop: 'decs' }
  ],
  model: { name: '', age: '', date: '', email: '', decs: '' } as FormModel,
  toAdd: () => {
    form.show = true
    form.title = '新增用户'
    form.state = 'add'
    form.ref?.resetFields()
  },
  toEdit: (row: FormModel) => {
    form.show = true
    form.title = '编辑用户资料'
    form.state = 'edit'
    form.model = { ...row }
  },
  toView: (row: FormModel) => {
    form.show = true
    form.title = '查看用户资料'
    form.state = 'view'
    form.model = { ...row }
  },
  submit: () => {
    form.ref?.validate().then(() => {
      form.show = false
      table.request()
    })
  }
})

// 表格配置
const table = reactive({
  loading: false,
  border: true,
  data: [] as FormModel[],
  page: {
    sortOrder: null,
    sortProp: '',
    currentPage: 1,
    pageSize: 10,
    total: 0
  } as TablePage,
  columns: [
    { label: '#', type: 'selection' },
    { label: '姓名', prop: 'name', sortable: 'custom', width: 100 },
    { label: '年龄', prop: 'age', sortable: 'custom', width: 100 },
    { label: '邮件', prop: 'email', width: 250 },
    { label: '出生日期', prop: 'date', width: 150 },
    { label: '介绍', prop: 'decs', showOverflowTooltip: true },
    { width: '120px', label: '操作', fixed: 'right', slot: 'operation' }
  ],
  request: () => {
    table.loading = true
    http.post('/mock/data', { size: table.page.pageSize }).then((res) => {
      table.data = res.data.list
      table.page.total = res.data.total
      table.loading = false
    })
  }
})

// 初始化加载数据
table.request()
</script>
