<template>
  <div class="flex flex-col overflow-auto">
    <!-- 按钮条 -->
    <div class="table-bar flex justify-between items-center mb-3">
      <div>
        <ElButton icon="Plus" @click="form.toAdd">新增</ElButton>
        <el-popconfirm title="Are you sure to delete this?" @confirm="deleteSelected">
          <template #reference>
            <ElButton icon="Delete">删除</ElButton>
          </template>
        </el-popconfirm>
      </div>
      <div>
        <ElButton icon="Refresh" round @click="table.request"></ElButton>
      </div>
    </div>
    <!-- 列表 -->
    <AgelTable class="flex-1" v-bind="table" @selection-change="handleSelectionChange"> </AgelTable>
    <!-- 弹窗表单 -->
    <ElDialog v-model="form.show" :title="form.title" width="800px" top="10vh">
      <ElForm :ref="(v) => (form.ref = v)" :model="form.model" label-width="80px">
        <ElFormItem label="集群名称" prop="name" required>
          <ElInput v-model="form.model.name" :disabled="form.state === 'view'" />
        </ElFormItem>
        <ElFormItem label="APIServer" prop="api_server" required>
          <ElInput v-model="form.model.api_server" :disabled="form.state === 'view'" />
        </ElFormItem>
        <ElFormItem label="Region" prop="region">
          <ElInput v-model="form.model.region" :disabled="form.state === 'view'" />
        </ElFormItem>
        <ElFormItem label="集群版本" prop="version">
          <ElInput v-model="form.model.version" :disabled="form.state === 'view'" />
        </ElFormItem>
        <ElFormItem label="介绍" prop="description">
          <ElInput 
            v-model="form.model.description" 
            type="textarea" 
            :rows="5"
            :disabled="form.state === 'view'" 
          />
        </ElFormItem>
        <ElFormItem 
          v-if="form.state !== 'edit'" 
          label="集群凭证" 
          prop="config" 
          required
        >
          <ElInput 
            v-model="form.model.config" 
            type="textarea" 
            :rows="20"
            :disabled="form.state === 'view'" 
          />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton v-if="form.state === 'add'" type="primary" @click="form.submit">提交</ElButton>
        <ElButton v-if="form.state === 'edit'" type="primary" @click="form.submitEdit">提交</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref, nextTick } from 'vue'
import http from '@/api'
import yaml from 'js-yaml'

type ClusterForm = {
  name: string
  api_server: string
  region: string
  version: string
  description: string
  config: string
}

type TableRow = ClusterForm & {
  id: string
}

const selectedRows = ref<TableRow[]>([])

const form = reactive({
  ref: null as FormInstance | null,
  show: false,
  title: '',
  state: '',
  model: {
    name: '',
    api_server: '',
    region: '',
    version: '',
    description: '',
    config: ''
  } as ClusterForm,
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

const table = reactive({
  loading: false,
  border: true,
  data: [] as TableRow[],
  columns: [
    { label: '#', type: 'selection' },
    { label: '集群名称', prop: 'name', width: 200 },
    { label: 'APIServer', prop: 'api_server', width: 200 },
    { label: 'Region', prop: 'region', width: 150 },
    { label: '集群版本', prop: 'version', width: 150 },
    { label: '介绍', prop: 'description', showOverflowTooltip: true },
    {
      width: '200px',
      label: '操作',
      fixed: 'right',
      slot: (scope: { row: TableRow }) => {
        return (
          <div>
            <el-button link type="primary" onClick={() => form.toEdit(scope.row)}>
              编辑
            </el-button>
          </div>
        )
      }
    }
  ],
  request: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/cluster`).then((res) => {
      table.data = res.result
      table.loading = false
    })
  },
  create: (form: ClusterForm) => {
    form.config = yaml.dump(form.config)
      .replace(/^ *\|-\n/, '')
      .replace(/\\n {2}/g, '\n')
      .toString()
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/create/cluster`, form).then((res) => {
      table.loading = false
      table.request()
    })
  },
  edit: (form: ClusterForm) => {
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/update/cluster`, form).then((res) => {
      table.loading = false
      table.request()
    })
  },
  delete: (form: TableRow) => {
    table.loading = true
    http.delete(import.meta.env.VITE_APP_BASE_URL + `/cicd/delete/cluster/${form.id}`).then((res) => {
      table.loading = false
      table.request()
    })
  }
})

const handleSelectionChange = (selected: TableRow[]) => {
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

table.request()
</script>
